package service

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"time"

	"carRental/internal/model"
)

type AuthService struct {
	DB *sql.DB
}

func (s *AuthService) Login(ctx context.Context, loginName, plainPwd string) (*model.User, error) {
	pwdMD5 := md5Hex(plainPwd)
	row := s.DB.QueryRowContext(ctx,
		`SELECT userid, loginname, realname, type
		 FROM sys_user
		 WHERE loginname=? AND pwd=? AND available=1
		 LIMIT 1`,
		loginName, pwdMD5,
	)

	var u model.User
	if err := row.Scan(&u.UserID, &u.LoginName, &u.RealName, &u.Type); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *AuthService) AddLoginLog(ctx context.Context, loginName, loginIP string, t time.Time) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO sys_log_login(loginname, loginip, logintime) VALUES (?,?,?)`,
		loginName, loginIP, t,
	)
	return err
}

func md5Hex(s string) string {
	sum := md5.Sum([]byte(s))
	return hex.EncodeToString(sum[:])
}

