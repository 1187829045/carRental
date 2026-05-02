package service

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"carRental/internal/model"
)

const (
	defaultCarImg = "images/defaultcarimage.jpg"
	fileUploadTmp = "_temp"
)

type BusService struct {
	DB          *sql.DB
	FileService *FileService
}

func (s *BusService) QueryCars(ctx context.Context, q map[string]string) (int64, []model.Car, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := []any{}
	for _, item := range []struct{ key, col string }{
		{"carnumber", "carnumber"},
		{"cartype", "cartype"},
		{"color", "color"},
		{"description", "description"},
	} {
		if v := strings.TrimSpace(q[item.key]); v != "" {
			where = append(where, item.col+" LIKE ?")
			args = append(args, likeArg(v))
		}
	}
	if v := strings.TrimSpace(q["isrenting"]); v != "" {
		where = append(where, "isrenting=?")
		args = append(args, v)
	}
	if v := strings.TrimSpace(q["opername"]); v != "" {
		where = append(where, "su.realname=?")
		args = append(args, v)
	}
	if v := strings.TrimSpace(q["operid"]); v != "" {
		where = append(where, "bc.operid=?")
		args = append(args, v)
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM bus_car bc LEFT JOIN sys_user su ON bc.operid = su.userid WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT bc.carnumber,bc.cartype,bc.color,bc.price,bc.rentprice,bc.deposit,bc.isrenting,bc.description,bc.carimg,bc.createtime,bc.operid,su.realname FROM bus_car bc LEFT JOIN sys_user su ON bc.operid = su.userid WHERE `+wsql+` ORDER BY bc.createtime DESC LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.Car
	for rows.Next() {
		var x model.Car
		var opername sql.NullString
		if err := rows.Scan(&x.CarNumber, &x.CarType, &x.Color, &x.Price, &x.RentPrice, &x.Deposit, &x.IsRenting, &x.Description, &x.CarImg, &x.CreateTime, &x.OperId, &opername); err != nil {
			return 0, nil, err
		}
		if opername.Valid {
			x.OperName = opername.String
		}
		out = append(out, x)
	}
	return count, out, rows.Err()
}

func (s *BusService) GetCar(ctx context.Context, carnumber string) (*model.Car, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT bc.carnumber,bc.cartype,bc.color,bc.price,bc.rentprice,bc.deposit,bc.isrenting,bc.description,bc.carimg,bc.createtime,bc.operid,su.realname FROM bus_car bc LEFT JOIN sys_user su ON bc.operid = su.userid WHERE bc.carnumber=?`, carnumber)
	var x model.Car
	var opername sql.NullString
	if err := row.Scan(&x.CarNumber, &x.CarType, &x.Color, &x.Price, &x.RentPrice, &x.Deposit, &x.IsRenting, &x.Description, &x.CarImg, &x.CreateTime, &x.OperId, &opername); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if opername.Valid {
		x.OperName = opername.String
	}
	return &x, nil
}

func (s *BusService) normalizeCarImage(rel string) (string, error) {
	rel = strings.TrimSpace(rel)
	if rel == "" {
		return defaultCarImg, nil
	}
	if rel == defaultCarImg {
		return rel, nil
	}
	if !strings.HasSuffix(rel, fileUploadTmp) {
		return rel, nil
	}
	newRel := strings.TrimSuffix(rel, fileUploadTmp)
	oldPath, err := SafeJoin(s.FileService.UploadRoot, rel)
	if err != nil {
		return "", err
	}
	newPath, err := SafeJoin(s.FileService.UploadRoot, newRel)
	if err != nil {
		return "", err
	}
	if err := EnsureDir(filepath.Dir(newPath)); err != nil {
		return "", err
	}
	if _, statErr := os.Stat(oldPath); statErr == nil {
		if err := os.Rename(oldPath, newPath); err != nil {
			return "", err
		}
	}
	return newRel, nil
}

func (s *BusService) removeUpload(rel string) {
	if rel == "" || rel == defaultCarImg {
		return
	}
	if abs, err := SafeJoin(s.FileService.UploadRoot, rel); err == nil {
		_ = os.Remove(abs)
	}
}

func (s *BusService) AddCar(ctx context.Context, x model.Car) error {
	carimg, err := s.normalizeCarImage(x.CarImg)
	if err != nil {
		return err
	}
	if x.CreateTime.IsZero() {
		x.CreateTime = time.Now()
	}
	_, err = s.DB.ExecContext(ctx, `INSERT INTO bus_car(carnumber,cartype,color,price,rentprice,deposit,isrenting,description,carimg,createtime,operid) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
		x.CarNumber, x.CarType, x.Color, x.Price, x.RentPrice, x.Deposit, x.IsRenting, x.Description, carimg, x.CreateTime, x.OperId)
	return err
}

func (s *BusService) UpdateCar(ctx context.Context, x model.Car) error {
	if strings.HasSuffix(strings.TrimSpace(x.CarImg), fileUploadTmp) {
		carimg, err := s.normalizeCarImage(x.CarImg)
		if err != nil {
			return err
		}
		old, err := s.GetCar(ctx, x.CarNumber)
		if err != nil {
			return err
		}
		if old != nil {
			s.removeUpload(old.CarImg)
		}
		x.CarImg = carimg
	}
	_, err := s.DB.ExecContext(ctx, `UPDATE bus_car SET cartype=?,color=?,price=?,rentprice=?,deposit=?,isrenting=?,description=?,carimg=? WHERE carnumber=?`,
		x.CarType, x.Color, x.Price, x.RentPrice, x.Deposit, x.IsRenting, x.Description, x.CarImg, x.CarNumber)
	return err
}

func (s *BusService) UpdateCarRenting(ctx context.Context, carnumber string, isrenting int) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE bus_car SET isrenting=? WHERE carnumber=?`, isrenting, carnumber)
	return err
}

func (s *BusService) DeleteCar(ctx context.Context, carnumber string) error {
	old, err := s.GetCar(ctx, carnumber)
	if err != nil {
		return err
	}
	if old != nil {
		s.removeUpload(old.CarImg)
	}
	_, err = s.DB.ExecContext(ctx, `DELETE FROM bus_car WHERE carnumber=?`, carnumber)
	return err
}

func (s *BusService) QueryCustomers(ctx context.Context, q map[string]string) (int64, []model.Customer, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := []any{}
	for _, item := range []struct{ key, col string }{
		{"identity", "identity"},
		{"custname", "custname"},
		{"phone", "phone"},
		{"career", "career"},
		{"address", "address"},
	} {
		if v := strings.TrimSpace(q[item.key]); v != "" {
			where = append(where, item.col+" LIKE ?")
			args = append(args, likeArg(v))
		}
	}
	if v := strings.TrimSpace(q["sex"]); v != "" {
		where = append(where, "sex=?")
		args = append(args, v)
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM bus_customer WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT identity,custname,sex,address,phone,career,createtime FROM bus_customer WHERE `+wsql+` ORDER BY createtime DESC LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.Customer
	for rows.Next() {
		var x model.Customer
		if err := rows.Scan(&x.Identity, &x.CustName, &x.Sex, &x.Address, &x.Phone, &x.Career, &x.CreateTime); err != nil {
			return 0, nil, err
		}
		out = append(out, x)
	}
	return count, out, rows.Err()
}

func (s *BusService) ListCustomers(ctx context.Context, q map[string]string) ([]model.Customer, error) {
	_, data, err := s.QueryCustomers(ctx, map[string]string{
		"page":     "1",
		"limit":    "1000000",
		"identity": q["identity"],
		"custname": q["custname"],
		"phone":    q["phone"],
		"career":   q["career"],
		"address":  q["address"],
		"sex":      q["sex"],
	})
	return data, err
}

func (s *BusService) GetCustomer(ctx context.Context, identity string) (*model.Customer, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT identity,custname,sex,address,phone,career,createtime FROM bus_customer WHERE identity=?`, identity)
	var x model.Customer
	if err := row.Scan(&x.Identity, &x.CustName, &x.Sex, &x.Address, &x.Phone, &x.Career, &x.CreateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &x, nil
}

func (s *BusService) AddCustomer(ctx context.Context, x model.Customer) error {
	if x.CreateTime.IsZero() {
		x.CreateTime = time.Now()
	}
	_, err := s.DB.ExecContext(ctx, `INSERT INTO bus_customer(identity,custname,sex,address,phone,career,createtime) VALUES(?,?,?,?,?,?,?)`,
		x.Identity, x.CustName, x.Sex, x.Address, x.Phone, x.Career, x.CreateTime)
	return err
}
func (s *BusService) UpdateCustomer(ctx context.Context, x model.Customer) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE bus_customer SET custname=?,sex=?,address=?,phone=?,career=? WHERE identity=?`,
		x.CustName, x.Sex, x.Address, x.Phone, x.Career, x.Identity)
	return err
}
func (s *BusService) DeleteCustomer(ctx context.Context, identity string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM bus_customer WHERE identity=?`, identity)
	return err
}

func (s *BusService) QueryRents(ctx context.Context, q map[string]string) (int64, []model.Rent, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := []any{}
	for _, item := range []struct{ key, col string }{
		{"rentid", "rentid"},
		{"carnumber", "carnumber"},
		{"identity", "identity"},
	} {
		if v := strings.TrimSpace(q[item.key]); v != "" {
			where = append(where, item.col+" LIKE ?")
			args = append(args, likeArg(v))
		}
	}
	if t, err := parseDateTime(q["startTime"]); err == nil && t != nil {
		where = append(where, "createtime >= ?")
		args = append(args, *t)
	}
	if t, err := parseDateTime(q["endTime"]); err == nil && t != nil {
		where = append(where, "createtime <= ?")
		args = append(args, *t)
	}
	if v := strings.TrimSpace(q["rentflag"]); v != "" {
		where = append(where, "rentflag=?")
		args = append(args, v)
	}
	if v := strings.TrimSpace(q["opername"]); v != "" {
		where = append(where, "su.realname=?")
		args = append(args, v)
	}
	if v := strings.TrimSpace(q["operid"]); v != "" {
		where = append(where, "br.operid=?")
		args = append(args, v)
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM bus_rent br LEFT JOIN sys_user su ON br.operid = su.userid WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT br.rentid,br.price,br.begindate,br.returndate,br.rentflag,br.identity,br.carnumber,br.operid,su.realname,br.createtime FROM bus_rent br LEFT JOIN sys_user su ON br.operid = su.userid WHERE `+wsql+` ORDER BY br.createtime DESC LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.Rent
	for rows.Next() {
		var x model.Rent
		var opername sql.NullString
		if err := rows.Scan(&x.RentID, &x.Price, &x.BeginDate, &x.ReturnDate, &x.RentFlag, &x.Identity, &x.CarNumber, &x.OperId, &opername, &x.CreateTime); err != nil {
			return 0, nil, err
		}
		if opername.Valid {
			x.OperName = opername.String
		}
		out = append(out, x)
	}
	return count, out, rows.Err()
}

func (s *BusService) GetRent(ctx context.Context, rentid string) (*model.Rent, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT br.rentid,br.price,br.begindate,br.returndate,br.rentflag,br.identity,br.carnumber,br.operid,su.realname,br.createtime FROM bus_rent br LEFT JOIN sys_user su ON br.operid = su.userid WHERE br.rentid=?`, rentid)
	var x model.Rent
	var opername sql.NullString
	if err := row.Scan(&x.RentID, &x.Price, &x.BeginDate, &x.ReturnDate, &x.RentFlag, &x.Identity, &x.CarNumber, &x.OperId, &opername, &x.CreateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if opername.Valid {
		x.OperName = opername.String
	}
	return &x, nil
}

func (s *BusService) NewRentForm(ctx context.Context, identity string) (*model.Rent, error) {
	customer, err := s.GetCustomer(ctx, identity)
	if err != nil {
		return nil, err
	}
	if customer == nil {
		return nil, nil
	}
	now := time.Now()
	return &model.Rent{
		RentID:    randomOrderID("CZ", now),
		BeginDate: now,
		Identity:  identity,
		OperName:  customer.CustName,
	}, nil
}

func randomOrderID(prefix string, now time.Time) string {
	return fmt.Sprintf("%s_%s_%04d_%05d", prefix, now.Format("20060102_150405"), now.Nanosecond()%10000, now.UnixNano()%100000)
}

func (s *BusService) SaveRent(ctx context.Context, x model.Rent) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if x.CreateTime.IsZero() {
		x.CreateTime = time.Now()
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO bus_rent(rentid,price,begindate,returndate,rentflag,identity,carnumber,operid,createtime) VALUES(?,?,?,?,?,?,?,?,?)`,
		x.RentID, x.Price, x.BeginDate, x.ReturnDate, x.RentFlag, x.Identity, x.CarNumber, x.OperId, x.CreateTime); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE bus_car SET isrenting=? WHERE carnumber=?`, 2, x.CarNumber); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *BusService) UpdateRent(ctx context.Context, x model.Rent) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE bus_rent SET price=?,begindate=?,returndate=?,rentflag=?,identity=?,carnumber=?,operid=? WHERE rentid=?`,
		x.Price, x.BeginDate, x.ReturnDate, x.RentFlag, x.Identity, x.CarNumber, x.OperId, x.RentID)
	return err
}

func (s *BusService) DeleteRent(ctx context.Context, rentid string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	var carnumber string
	if err := tx.QueryRowContext(ctx, `SELECT carnumber FROM bus_rent WHERE rentid=?`, rentid).Scan(&carnumber); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE bus_car SET isrenting=? WHERE carnumber=?`, 0, carnumber); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM bus_rent WHERE rentid=?`, rentid); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *BusService) CheckRent(ctx context.Context, rentid, carnumber string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `UPDATE bus_rent SET rentflag=? WHERE rentid=?`, 0, rentid); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE bus_car SET isrenting=? WHERE carnumber=?`, 1, carnumber); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *BusService) QueryChecks(ctx context.Context, q map[string]string) (int64, []model.Check, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := []any{}
	for _, item := range []struct{ key, col string }{
		{"checkid", "checkid"},
		{"rentid", "rentid"},
		{"problem", "problem"},
		{"checkdesc", "checkdesc"},
	} {
		if v := strings.TrimSpace(q[item.key]); v != "" {
			where = append(where, item.col+" LIKE ?")
			args = append(args, likeArg(v))
		}
	}
	if t, err := parseDateTime(q["startTime"]); err == nil && t != nil {
		where = append(where, "checkdate >= ?")
		args = append(args, *t)
	}
	if t, err := parseDateTime(q["endTime"]); err == nil && t != nil {
		where = append(where, "checkdate <= ?")
		args = append(args, *t)
	}
	if v := strings.TrimSpace(q["opername"]); v != "" {
		where = append(where, "su.realname=?")
		args = append(args, v)
	}
	if v := strings.TrimSpace(q["operid"]); v != "" {
		where = append(where, "bc.operid=?")
		args = append(args, v)
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM bus_check bc LEFT JOIN sys_user su ON bc.operid = su.userid WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT bc.checkid,bc.checkdate,bc.checkdesc,bc.problem,bc.paymoney,bc.operid,su.realname,bc.rentid,bc.createtime FROM bus_check bc LEFT JOIN sys_user su ON bc.operid = su.userid WHERE `+wsql+` ORDER BY bc.createtime DESC LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.Check
	for rows.Next() {
		var x model.Check
		var opername sql.NullString
		if err := rows.Scan(&x.CheckID, &x.CheckDate, &x.CheckDesc, &x.Problem, &x.PayMoney, &x.OperId, &opername, &x.RentID, &x.CreateTime); err != nil {
			return 0, nil, err
		}
		if opername.Valid {
			x.OperName = opername.String
		}
		out = append(out, x)
	}
	return count, out, rows.Err()
}

func (s *BusService) InitCheckFormData(ctx context.Context, rentid string) (map[string]any, error) {
	rent, err := s.GetRent(ctx, rentid)
	if err != nil || rent == nil {
		return nil, err
	}
	customer, err := s.GetCustomer(ctx, rent.Identity)
	if err != nil {
		return nil, err
	}
	car, err := s.GetCar(ctx, rent.CarNumber)
	if err != nil {
		return nil, err
	}
	check := model.Check{
		CheckID:   randomOrderID("JC", time.Now()),
		RentID:    rentid,
		CheckDate: time.Now(),
		OperName:  rent.OperName,
	}
	return map[string]any{
		"rent":     rent,
		"customer": customer,
		"car":      car,
		"check":    check,
	}, nil
}

func (s *BusService) SaveCheck(ctx context.Context, x model.Check) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if x.CreateTime.IsZero() {
		x.CreateTime = time.Now()
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO bus_check(checkid,checkdate,checkdesc,problem,paymoney,operid,rentid,createtime) VALUES(?,?,?,?,?,?,?,?)`,
		x.CheckID, x.CheckDate, x.CheckDesc, x.Problem, x.PayMoney, x.OperId, x.RentID, x.CreateTime); err != nil {
		return err
	}
	var carnumber string
	if err = tx.QueryRowContext(ctx, `SELECT carnumber FROM bus_rent WHERE rentid=?`, x.RentID).Scan(&carnumber); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE bus_rent SET rentflag=? WHERE rentid=?`, 1, x.RentID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE bus_car SET isrenting=? WHERE carnumber=?`, 0, carnumber); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *BusService) UpdateCheck(ctx context.Context, x model.Check) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE bus_check SET checkdate=?,checkdesc=?,problem=?,paymoney=?,operid=?,rentid=? WHERE checkid=?`,
		x.CheckDate, x.CheckDesc, x.Problem, x.PayMoney, x.OperId, x.RentID, x.CheckID)
	return err
}
func (s *BusService) DeleteCheck(ctx context.Context, checkid string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM bus_check WHERE checkid=?`, checkid)
	return err
}

func (s *BusService) QueryFranchisees(ctx context.Context, q map[string]string) (int64, []model.Franchisee, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := []any{}
	if v := strings.TrimSpace(q["name"]); v != "" {
		where = append(where, "name LIKE ?")
		args = append(args, likeArg(v))
	}
	if v := strings.TrimSpace(q["phone"]); v != "" {
		where = append(where, "phone LIKE ?")
		args = append(args, likeArg(v))
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM bus_franchisee WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT id,name,phone FROM bus_franchisee WHERE `+wsql+` ORDER BY id DESC LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.Franchisee
	for rows.Next() {
		var x model.Franchisee
		if err := rows.Scan(&x.ID, &x.Name, &x.Phone); err != nil {
			return 0, nil, err
		}
		out = append(out, x)
	}
	return count, out, rows.Err()
}
func (s *BusService) AddFranchisee(ctx context.Context, x model.Franchisee) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO bus_franchisee(name,phone) VALUES(?,?)`, x.Name, x.Phone)
	return err
}
func (s *BusService) UpdateFranchisee(ctx context.Context, x model.Franchisee) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE bus_franchisee SET name=?,phone=? WHERE id=?`, x.Name, x.Phone, x.ID)
	return err
}
func (s *BusService) DeleteFranchisee(ctx context.Context, id int) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM bus_franchisee WHERE id=?`, id)
	return err
}
