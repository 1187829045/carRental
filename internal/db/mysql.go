package db

import (
	"database/sql"
	"strings"
	"time"

	"carRental/internal/monitor"

	mysql "github.com/go-sql-driver/mysql"
)

func OpenMySQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open(monitor.DriverName, dsn)
	if err != nil {
		return nil, err
	}
	// Conservative defaults.
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func OpenMySQLAutoCreateDB(dsn string) (*sql.DB, error) {
	db, err := OpenMySQL(dsn)
	if err == nil {
		return db, nil
	}
	if !strings.Contains(err.Error(), "Unknown database") {
		return nil, err
	}
	cfg, parseErr := mysql.ParseDSN(dsn)
	if parseErr != nil {
		return nil, err
	}
	dbName := cfg.DBName
	if dbName == "" {
		return nil, err
	}
	cfg.DBName = ""
	admin, adminErr := sql.Open(monitor.DriverName, cfg.FormatDSN())
	if adminErr != nil {
		return nil, err
	}
	defer admin.Close()
	if pingErr := admin.Ping(); pingErr != nil {
		return nil, err
	}
	if _, execErr := admin.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "` DEFAULT CHARACTER SET utf8"); execErr != nil {
		return nil, execErr
	}
	return OpenMySQL(dsn)
}
