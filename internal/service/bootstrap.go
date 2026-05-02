package service

import (
	"context"
	"database/sql"
	"time"
)

func EnsureFranchiseeTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS bus_franchisee (
  id INT NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) DEFAULT NULL,
  phone VARCHAR(255) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8`)
	if err != nil {
		return err
	}
	// 动态为 bus_car 增加 opername 列（如果不存在）
	db.ExecContext(ctx, "ALTER TABLE bus_car ADD COLUMN opername VARCHAR(255) DEFAULT 'admin'")
	return nil
}

func StartTempFileCleanup(fileService *FileService) func() {
	stop := make(chan struct{})
	ticker := time.NewTicker(6 * time.Hour)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_, _ = fileService.CleanupTempFiles(24 * time.Hour)
			case <-stop:
				return
			}
		}
	}()
	return func() { close(stop) }
}

