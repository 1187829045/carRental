package service

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"carRental/internal/model"

	"github.com/xuri/excelize/v2"
)

type StatService struct {
	DB         *sql.DB
	BusService *BusService
}

func (s *StatService) LoadDashboardMetrics(ctx context.Context, rangeType string) (*model.DashboardMetrics, error) {
	start, end, rangeKey, rangeLabel := resolveDashboardRange(time.Now(), rangeType)
	metrics := &model.DashboardMetrics{
		Range:             rangeKey,
		RangeLabel:        rangeLabel,
		RangeStart:        start,
		RangeEnd:          end,
		RefreshedAt:       time.Now(),
		VehicleScopeLabel: "实时库存快照",
	}

	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM bus_car WHERE COALESCE(isrenting, 0) <> 0`).Scan(&metrics.InRentCarCount); err != nil {
		return nil, err
	}
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM bus_car WHERE COALESCE(isrenting, 0) = 0`).Scan(&metrics.IdleCarCount); err != nil {
		return nil, err
	}
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM bus_rent WHERE COALESCE(rentflag, 0) <> 1`).Scan(&metrics.RentOrderCount); err != nil {
		return nil, err
	}
	if err := s.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM bus_customer`).Scan(&metrics.CustomerCount); err != nil {
		return nil, err
	}
	return metrics, nil
}

func resolveDashboardRange(now time.Time, raw string) (time.Time, time.Time, string, string) {
	raw = strings.ToLower(strings.TrimSpace(raw))
	switch raw {
	case "day":
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return start, start.AddDate(0, 0, 1), "day", "今日"
	case "week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -(weekday - 1))
		return start, start.AddDate(0, 0, 7), "week", "本周"
	default:
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		return start, start.AddDate(0, 1, 0), "month", "本月"
	}
}

func (s *StatService) LoadCustomerAreaStat(ctx context.Context) ([]model.BaseEntity, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT address AS name, COUNT(1) AS value FROM bus_customer GROUP BY address`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.BaseEntity
	for rows.Next() {
		var x model.BaseEntity
		if err := rows.Scan(&x.Name, &x.Value); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *StatService) LoadCustomerAreaSexStat(ctx context.Context, area string) ([]model.BaseEntity, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT sex AS name, COUNT(1) AS value FROM bus_customer WHERE address=? GROUP BY sex`, area)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.BaseEntity
	for rows.Next() {
		var sex int
		var value float64
		if err := rows.Scan(&sex, &value); err != nil {
			return nil, err
		}
		out = append(out, model.BaseEntity{Name: strconv.Itoa(sex), Value: value})
	}
	return out, rows.Err()
}

func (s *StatService) LoadOpernameYearGradeStat(ctx context.Context, year string) ([]model.BaseEntity, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT opername AS name, SUM(price) AS value FROM bus_rent WHERE DATE_FORMAT(createtime,"%Y")=? GROUP BY opername`, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.BaseEntity
	for rows.Next() {
		var x model.BaseEntity
		if err := rows.Scan(&x.Name, &x.Value); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (s *StatService) LoadCompanyYearGradeStat(ctx context.Context, year string) ([]float64, error) {
	out := make([]float64, 12)
	for i := 1; i <= 12; i++ {
		ym := fmt.Sprintf("%s%02d", year, i)
		var v sql.NullFloat64
		if err := s.DB.QueryRowContext(ctx, `SELECT SUM(price) FROM bus_rent WHERE DATE_FORMAT(createtime,"%Y%m")=?`, ym).Scan(&v); err != nil {
			return nil, err
		}
		if v.Valid {
			out[i-1] = v.Float64
		}
	}
	return out, nil
}

func ExportCustomersAsXLS(customers []model.Customer, sheetName string) ([]byte, error) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	defaultSheet := f.GetSheetName(0)
	if defaultSheet != sheetName {
		f.SetSheetName(defaultSheet, sheetName)
	}
	headers := []string{"身份证", "姓名", "性别", "地址", "电话", "职业", "创建时间"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheetName, cell, h)
	}
	for r, c := range customers {
		values := []any{
			c.Identity,
			c.CustName,
			c.Sex,
			c.Address,
			c.Phone,
			c.Career,
			c.CreateTime.Format("2006-01-02 15:04:05"),
		}
		for i, v := range values {
			cell, _ := excelize.CoordinatesToCellName(i+1, r+2)
			_ = f.SetCellValue(sheetName, cell, v)
		}
	}
	var b bytes.Buffer
	if err := f.Write(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func ExportRentAsXLS(rent *model.Rent, customer *model.Customer, sheetName string) ([]byte, error) {
	f := excelize.NewFile()
	defer func() { _ = f.Close() }()
	defaultSheet := f.GetSheetName(0)
	if defaultSheet != sheetName {
		f.SetSheetName(defaultSheet, sheetName)
	}
	row := 1
	writeKV := func(k string, v any) {
		_ = f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), k)
		_ = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), v)
		row++
	}
	if customer != nil {
		writeKV("客户姓名", customer.CustName)
		writeKV("身份证", customer.Identity)
		writeKV("电话", customer.Phone)
		writeKV("地址", customer.Address)
	}
	if rent != nil {
		writeKV("出租单号", rent.RentID)
		writeKV("车牌号", rent.CarNumber)
		writeKV("金额", rent.Price)
		writeKV("起租时间", rent.BeginDate.Format("2006-01-02 15:04:05"))
		if rent.ReturnDate != nil {
			writeKV("归还时间", rent.ReturnDate.Format("2006-01-02 15:04:05"))
		}
		writeKV("业务员", rent.OperName)
	}
	var b bytes.Buffer
	if err := f.Write(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
