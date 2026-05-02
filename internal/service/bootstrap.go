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
	
	// Migrate OperName to OperId for semantic accuracy and robust permission control
	migrateSQLs := []string{
		"ALTER TABLE bus_car ADD COLUMN operid INT DEFAULT 1",
		"ALTER TABLE bus_rent ADD COLUMN operid INT DEFAULT 1",
		"ALTER TABLE bus_check ADD COLUMN operid INT DEFAULT 1",
		"UPDATE bus_car bc JOIN sys_user su ON bc.opername = su.realname SET bc.operid = su.userid WHERE bc.opername IS NOT NULL",
		"UPDATE bus_rent br JOIN sys_user su ON br.opername = su.realname SET br.operid = su.userid WHERE br.opername IS NOT NULL",
		"UPDATE bus_check bc JOIN sys_user su ON bc.opername = su.realname SET bc.operid = su.userid WHERE bc.opername IS NOT NULL",
	}
	for _, query := range migrateSQLs {
		db.ExecContext(ctx, query) // Ignore errors as columns might already exist
	}
	
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

