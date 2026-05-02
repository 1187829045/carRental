package main

import (
	"context"
	"log"
	"path/filepath"
	"runtime"

	"carRental/internal/config"
	"carRental/internal/db"
	"carRental/internal/handler"
	"carRental/internal/monitor"
	"carRental/internal/service"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("failed to get caller info")
	}
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	projectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		log.Fatalf("abs project root: %v", err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	monitorCollector := monitor.NewCollector(monitor.MetaFromDSN(cfg.MySQLDSN))
	monitor.SetDefault(monitorCollector)

	mysqlDB, err := db.OpenMySQLAutoCreateDB(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	defer mysqlDB.Close()
	if err := (service.DBInit{ProjectRoot: projectRoot}).EnsureSchemaAndSeed(context.Background(), mysqlDB); err != nil {
		log.Fatalf("db init: %v", err)
	}

	// 执行动态表结构和安全字段迁移（如增加operid字段，分配默认属主）
	if err := service.EnsureFranchiseeTable(context.Background(), mysqlDB); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	fileService, err := service.NewFileService(projectRoot)
	if err != nil {
		log.Fatalf("init file service: %v", err)
	}
	stopCleanup := service.StartTempFileCleanup(fileService)
	defer stopCleanup()

	busService := &service.BusService{
		DB:          mysqlDB,
		FileService: fileService,
	}

	r := handler.NewRouter(handler.Deps{
		Cfg:     cfg,
		Monitor: monitorCollector,
		AuthService: &service.AuthService{
			DB: mysqlDB,
		},
		MenuService: &service.MenuService{
			DB: mysqlDB,
		},
		SystemService: &service.SystemService{
			DB: mysqlDB,
		},
		BusService:  busService,
		StatService: &service.StatService{DB: mysqlDB, BusService: busService},
		FileService: fileService,
	})

	log.Printf("Go server listening on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

// [GIN-debug] POST   /carRental/rent/updateRent.action --> carRental/internal/handler.registerBusRoutes.func16 (7 handlers)
// [GIN-debug] POST   /carRental/rent/checkRent.action --> carRental/internal/handler.registerBusRoutes.func17 (7 handlers)
// [GIN-debug] GET    /carRental/rent/loadAllRent.action --> carRental/internal/handler.registerBusRoutes.func18 (7 handlers)
// panic: handlers are already registered for path '/carRental/rent/loadAllRent.action'

// goroutine 1 [running]:
// github.com/gin-gonic/gin.(*node).addRoute(0x7ff6cc6ff7dd?, {0xc0003f8960, 0x22}, {0xc0003f1040, 0x7, 0x7})
//         C:/Users/11878/go/pkg/mod/github.com/gin-gonic/gin@v1.12.0/tree.go:243 +0x6e8
// github.com/gin-gonic/gin.(*Engine).addRoute(0xc0003a0380, {0x7ff6cc6ff7dd, 0x3}, {0xc0003f8960, 0x22}, {0xc0003f1040, 0x7, 0x7})
//         C:/Users/11878/go/pkg/mod/github.com/gin-gonic/gin@v1.12.0/gin.go:377 +0x25f
// github.com/gin-gonic/gin.(*RouterGroup).handle(0xc0003c0780, {0x7ff6cc6ff7dd, 0x3}, {0x7ff6cc729b4c?, 0xc0003982f0?}, {0xc000398558, 0x1, 0xc0003982f0?})
//         C:/Users/11878/go/pkg/mod/github.com/gin-gonic/gin@v1.12.0/routergroup.go:89 +0x13e
// github.com/gin-gonic/gin.(*RouterGroup).GET(...)
//         C:/Users/11878/go/pkg/mod/github.com/gin-gonic/gin@v1.12.0/routergroup.go:117
// carRental/internal/handler.registerBusRoutes(0xc0003c0780, {{{0x7ff6cc702622, 0x5}, {0x7ff6cc712862, 0xa}, {0xc0000328c0, 0x34}, {0xc000031140, 0x2e}, {0xc00003a480, ...}, ...}, ...})
//         E:/go_project/carRental/internal/handler/modules.go:694 +0xa98
// carRental/internal/handler.registerModules(0xc0003c0780, {{{0x7ff6cc702622, 0x5}, {0x7ff6cc712862, 0xa}, {0xc0000328c0, 0x34}, {0xc000031140, 0x2e}, {0xc00003a480, ...}, ...}, ...})
//         E:/go_project/carRental/internal/handler/modules.go:29 +0xe5
// carRental/internal/handler.NewRouter({{{0x7ff6cc702622, 0x5}, {0x7ff6cc712862, 0xa}, {0xc0000328c0, 0x34}, {0xc000031140, 0x2e}, {0xc00003a480, 0x54}, ...}, ...})
//         E:/go_project/carRental/internal/handler/router.go:360 +0x2ce9
// main.main()
//         E:/go_project/carRental/cmd/goserver/main.go:60 +0x758
// exit status 2
