package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"carRental/internal/config"
	"carRental/internal/db"
	"carRental/internal/handler"
	"carRental/internal/service"
)

func main() {
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd: %v", err)
	}
	projectRoot, err = filepath.Abs(projectRoot)
	if err != nil {
		log.Fatalf("abs wd: %v", err)
	}

	cfg, err := config.Load(projectRoot)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	mysqlDB, err := db.OpenMySQL(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open mysql: %v", err)
	}
	defer mysqlDB.Close()

	if ensureErr := service.EnsureFranchiseeTable(context.Background(), mysqlDB); ensureErr != nil {
		log.Fatalf("ensure franchisee table: %v", ensureErr)
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
		Cfg: cfg,
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
