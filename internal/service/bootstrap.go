package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"carRental/internal/model"
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
		"UPDATE bus_car SET carimg = 'https://images.pexels.com/photos/170811/pexels-photo-170811.jpeg?auto=compress&cs=tinysrgb&w=1200' WHERE carimg IS NULL OR TRIM(carimg) = '' OR carimg = 'images/defaultcarimage.jpg' OR carimg = 'static/images/cars/placeholder-1.svg'",
		// 1.1 修复"章三"笔误，将其统一修正为"张三"以确保匹配一致性
		"UPDATE bus_car SET opername = '张三' WHERE opername = '章三'",
		"UPDATE bus_rent SET opername = '张三' WHERE opername = '章三'",
		"UPDATE bus_check SET opername = '张三' WHERE opername = '章三'",
		// 1.2 根据现有的 opername (中文名) 映射到正确的 userid
		"UPDATE bus_car bc JOIN sys_user su ON bc.opername = su.realname SET bc.operid = su.userid WHERE bc.opername IS NOT NULL",
		"UPDATE bus_rent br JOIN sys_user su ON br.opername = su.realname SET br.operid = su.userid WHERE br.opername IS NOT NULL",
		"UPDATE bus_check bc JOIN sys_user su ON bc.opername = su.realname SET bc.operid = su.userid WHERE bc.opername IS NOT NULL",
		// 2. 针对匹配失败、历史数据遗留的空值或 0 值，统一 fallback 给超级管理员 (userid=1)
		"UPDATE bus_car SET operid = 1 WHERE operid IS NULL OR operid = 0",
		"UPDATE bus_rent SET operid = 1 WHERE operid IS NULL OR operid = 0",
		"UPDATE bus_check SET operid = 1 WHERE operid IS NULL OR operid = 0",
	}
	for _, query := range migrateSQLs {
		db.ExecContext(ctx, query) // Ignore errors as columns might already exist
	}

	ensureOperationalIndexes(ctx, db)
	if err := ensureMockOperationalData(ctx, db); err != nil {
		return err
	}

	return nil
}

const (
	mockCheckSeedCount      = 50
	mockZhangsanSeedCount   = 20
	mockCheckIDPrefix       = "MOCKCHK_"
	mockRentIDPrefix        = "MOCKRENT_"
	mockCarNumberPrefix     = "MOCKCAR-"
	mockCustomerNamePrefix  = "MockCustomer"
	defaultMockCarImagePath = "https://images.pexels.com/photos/170811/pexels-photo-170811.jpeg?auto=compress&cs=tinysrgb&w=1200"
)

type mockOperationalSeed struct {
	Customer       model.Customer
	Car            model.Car
	Rent           model.Rent
	Check          model.Check
	LegacyOperName string
}

func ensureOperationalIndexes(ctx context.Context, db *sql.DB) {
	indexSQLs := []string{
		"ALTER TABLE bus_car ADD INDEX idx_bus_car_operid (operid)",
		"ALTER TABLE bus_rent ADD INDEX idx_bus_rent_operid (operid)",
		"ALTER TABLE bus_rent ADD INDEX idx_bus_rent_identity (identity)",
		"ALTER TABLE bus_check ADD INDEX idx_bus_check_operid (operid)",
		"ALTER TABLE bus_check ADD INDEX idx_bus_check_rentid (rentid)",
	}
	for _, query := range indexSQLs {
		db.ExecContext(ctx, query) // Ignore duplicate-index errors
	}
}

func ensureMockOperationalData(ctx context.Context, db *sql.DB) error {
	zhangsanID, err := lookupUserID(ctx, db, "zhangsan")
	if err != nil {
		return err
	}
	adminID, err := lookupUserID(ctx, db, "admin")
	if err != nil {
		return err
	}

	var mockCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM bus_check WHERE checkid LIKE ?`, mockCheckIDPrefix+"%").Scan(&mockCount); err != nil {
		return err
	}
	var zhangsanCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM bus_check WHERE checkid LIKE ? AND operid = ?`, mockCheckIDPrefix+"%", zhangsanID).Scan(&zhangsanCount); err != nil {
		return err
	}
	if mockCount == mockCheckSeedCount && zhangsanCount == mockZhangsanSeedCount {
		return nil
	}

	seeds := buildMockOperationalSeeds(time.Now(), zhangsanID, adminID)
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM bus_check WHERE checkid LIKE ?`, mockCheckIDPrefix+"%"); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM bus_rent WHERE rentid LIKE ?`, mockRentIDPrefix+"%"); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM bus_car WHERE carnumber LIKE ?`, mockCarNumberPrefix+"%"); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM bus_customer WHERE custname LIKE ?`, mockCustomerNamePrefix+"%"); err != nil {
		return err
	}

	for _, seed := range seeds {
		if _, err := tx.ExecContext(ctx, `INSERT INTO bus_customer(identity,custname,sex,address,phone,career,createtime) VALUES(?,?,?,?,?,?,?)`,
			seed.Customer.Identity, seed.Customer.CustName, seed.Customer.Sex, seed.Customer.Address, seed.Customer.Phone, seed.Customer.Career, seed.Customer.CreateTime); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO bus_car(carnumber,cartype,color,price,rentprice,deposit,isrenting,description,carimg,createtime,operid) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
			seed.Car.CarNumber, seed.Car.CarType, seed.Car.Color, seed.Car.Price, seed.Car.RentPrice, seed.Car.Deposit, seed.Car.IsRenting, seed.Car.Description, seed.Car.CarImg, seed.Car.CreateTime, seed.Car.OperId); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO bus_rent(rentid,price,begindate,returndate,rentflag,identity,carnumber,opername,operid,createtime) VALUES(?,?,?,?,?,?,?,?,?,?)`,
			seed.Rent.RentID, seed.Rent.Price, seed.Rent.BeginDate, seed.Rent.ReturnDate, seed.Rent.RentFlag, seed.Rent.Identity, seed.Rent.CarNumber, seed.LegacyOperName, seed.Rent.OperId, seed.Rent.CreateTime); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO bus_check(checkid,checkdate,checkdesc,problem,paymoney,opername,operid,rentid,createtime) VALUES(?,?,?,?,?,?,?,?,?)`,
			seed.Check.CheckID, seed.Check.CheckDate, seed.Check.CheckDesc, seed.Check.Problem, seed.Check.PayMoney, seed.LegacyOperName, seed.Check.OperId, seed.Check.RentID, seed.Check.CreateTime); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func lookupUserID(ctx context.Context, db *sql.DB, loginName string) (int, error) {
	var userID int
	if err := db.QueryRowContext(ctx, `SELECT userid FROM sys_user WHERE loginname = ? LIMIT 1`, loginName).Scan(&userID); err != nil {
		return 0, fmt.Errorf("lookup user %s: %w", loginName, err)
	}
	return userID, nil
}

func buildMockOperationalSeeds(now time.Time, zhangsanID, adminID int) []mockOperationalSeed {
	problems := []string{"违停", "外观划痕", "轮胎磨损", "超速", "内饰污损", "逾期还车", "补漆", "玻璃破损"}
	carTypes := []string{"SUV", "轿车", "商务车", "新能源"}
	colors := []string{"白色", "黑色", "银色", "灰色", "蓝色"}
	seeds := make([]mockOperationalSeed, 0, mockCheckSeedCount)
	for i := 0; i < mockCheckSeedCount; i++ {
		operatorID := adminID
		operatorName := "超级管理员"
		if i < mockZhangsanSeedCount {
			operatorID = zhangsanID
			operatorName = "张三"
		}

		dayOffset := (i * 89) / (mockCheckSeedCount - 1)
		checkDate := now.AddDate(0, 0, -dayOffset).Add(time.Duration(i%5) * time.Hour)
		beginDate := checkDate.AddDate(0, 0, -(i%5 + 1))
		returnDate := checkDate
		createTime := checkDate.Add(30 * time.Minute)
		identity := NormalizeIdentity(fmt.Sprintf("5201011990%02d%02d%04d", i%12+1, i%28+1, i+1))
		carNumber := fmt.Sprintf("%s%03d", mockCarNumberPrefix, i+1)
		rentID := fmt.Sprintf("%s%03d", mockRentIDPrefix, i+1)
		checkID := fmt.Sprintf("%s%03d", mockCheckIDPrefix, i+1)
		customerName := fmt.Sprintf("%s%03d", mockCustomerNamePrefix, i+1)
		problem := problems[i%len(problems)]

		seeds = append(seeds, mockOperationalSeed{
			Customer: model.Customer{
				Identity:   identity,
				CustName:   customerName,
				Sex:        i % 2,
				Address:    fmt.Sprintf("测试地址%02d号", i+1),
				Phone:      fmt.Sprintf("138%08d", i+1),
				Career:     "测试客户",
				CreateTime: beginDate,
			},
			Car: model.Car{
				CarNumber:   carNumber,
				CarType:     carTypes[i%len(carTypes)],
				Color:       colors[i%len(colors)],
				Price:       100000 + float64(i)*3200,
				RentPrice:   300 + float64(i%10)*40,
				Deposit:     2000 + float64(i%8)*500,
				IsRenting:   0,
				Description: fmt.Sprintf("检查单联动测试车辆%03d", i+1),
				CarImg:      defaultMockCarImagePath,
				CreateTime:  beginDate,
				OperId:      operatorID,
				OperName:    operatorName,
			},
			Rent: model.Rent{
				RentID:     rentID,
				Price:      500 + float64(i%12)*80,
				BeginDate:  beginDate,
				ReturnDate: &returnDate,
				RentFlag:   1,
				Identity:   identity,
				CarNumber:  carNumber,
				OperId:     operatorID,
				OperName:   operatorName,
				CreateTime: beginDate,
			},
			Check: model.Check{
				CheckID:    checkID,
				CheckDate:  checkDate,
				CheckDesc:  fmt.Sprintf("%s处理记录-%03d", problem, i+1),
				Problem:    problem,
				PayMoney:   float64((i%6)+1) * 100,
				OperId:     operatorID,
				OperName:   operatorName,
				RentID:     rentID,
				CreateTime: createTime,
			},
			LegacyOperName: operatorName,
		})
	}
	return seeds
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
