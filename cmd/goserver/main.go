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
