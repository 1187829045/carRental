package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"carRental/internal/model"
)

type SystemService struct {
	DB *sql.DB
}

func (s *SystemService) DBStats(ctx context.Context) (map[string]string, error) {
	_ = ctx
	st := s.DB.Stats()
	return map[string]string{
		"open_connections":     fmt.Sprintf("%d", st.OpenConnections),
		"in_use":               fmt.Sprintf("%d", st.InUse),
		"idle":                 fmt.Sprintf("%d", st.Idle),
		"wait_count":           fmt.Sprintf("%d", st.WaitCount),
		"wait_duration":        st.WaitDuration.String(),
		"max_open_connections": fmt.Sprintf("%d", st.MaxOpenConnections),
	}, nil
}

func (s *SystemService) GetUserByID(ctx context.Context, userid int) (*model.SysUser, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT userid,loginname,identity,realname,sex,address,phone,pwd,position,type,available FROM sys_user WHERE userid=?`, userid)
	var u model.SysUser
	if err := row.Scan(&u.UserID, &u.LoginName, &u.Identity, &u.RealName, &u.Sex, &u.Address, &u.Phone, &u.Pwd, &u.Position, &u.Type, &u.Available); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (s *SystemService) QueryUsers(ctx context.Context, q map[string]string) (int64, []model.SysUser, error) {
	page, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	_ = page
	where := []string{"type<>1"}
	args := make([]any, 0, 16)
	addLike := func(k, col string) {
		if v := strings.TrimSpace(q[k]); v != "" {
			where = append(where, col+" LIKE ?")
			args = append(args, likeArg(v))
		}
	}
	addLike("realname", "realname")
	addLike("loginname", "loginname")
	addLike("identity", "identity")
	addLike("address", "address")
	addLike("phone", "phone")
	if v := strings.TrimSpace(q["sex"]); v != "" {
		where = append(where, "sex=?")
		args = append(args, v)
	}
	wsql := strings.Join(where, " AND ")

	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM sys_user WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx,
		`SELECT userid,loginname,identity,realname,sex,address,phone,pwd,position,type,available
		 FROM sys_user WHERE `+wsql+` ORDER BY userid DESC LIMIT ? OFFSET ?`,
		append(args, limit, offset)...,
	)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.SysUser
	for rows.Next() {
		var u model.SysUser
		if err := rows.Scan(&u.UserID, &u.LoginName, &u.Identity, &u.RealName, &u.Sex, &u.Address, &u.Phone, &u.Pwd, &u.Position, &u.Type, &u.Available); err != nil {
			return 0, nil, err
		}
		out = append(out, u)
	}
	return count, out, rows.Err()
}

func (s *SystemService) AddUser(ctx context.Context, u model.SysUser) error {
	defaultPwd := md5Hex("123456")
	userType := 2
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO sys_user(loginname,identity,realname,sex,address,phone,pwd,position,type,available)
		 VALUES(?,?,?,?,?,?,?,?,?,?)`,
		u.LoginName, u.Identity, u.RealName, u.Sex, u.Address, u.Phone, defaultPwd, u.Position, userType, u.Available)
	return err
}

func (s *SystemService) UpdateUser(ctx context.Context, u model.SysUser) error {
	_, err := s.DB.ExecContext(ctx,
		`UPDATE sys_user SET loginname=?,identity=?,realname=?,sex=?,address=?,phone=?,pwd=COALESCE(?,pwd),position=?,type=COALESCE(?,type),available=? WHERE userid=?`,
		u.LoginName, u.Identity, u.RealName, u.Sex, u.Address, u.Phone, u.Pwd, u.Position, u.Type, u.Available, u.UserID)
	return err
}

func (s *SystemService) DeleteUser(ctx context.Context, userid int) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `DELETE FROM sys_user WHERE userid=?`, userid); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM sys_role_user WHERE uid=?`, userid); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *SystemService) ResetUserPwd(ctx context.Context, userid int) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE sys_user SET pwd=? WHERE userid=?`, md5Hex("123456"), userid)
	return err
}

func (s *SystemService) QueryUserRoles(ctx context.Context, userid int) ([]map[string]any, error) {
	rows, err := s.DB.QueryContext(ctx, `SELECT roleid,rolename,roledesc,available FROM sys_role WHERE available=1 ORDER BY roleid ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	selected := map[int]bool{}
	rs, err := s.DB.QueryContext(ctx, `SELECT rid FROM sys_role_user WHERE uid=?`, userid)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		var rid int
		if err := rs.Scan(&rid); err != nil {
			_ = rs.Close()
			return nil, err
		}
		selected[rid] = true
	}
	_ = rs.Close()

	var data []map[string]any
	for rows.Next() {
		var r model.Role
		if err := rows.Scan(&r.RoleID, &r.RoleName, &r.RoleDesc, &r.Available); err != nil {
			return nil, err
		}
		data = append(data, map[string]any{
			"roleid":      r.RoleID,
			"rolename":    r.RoleName,
			"roledesc":    r.RoleDesc,
			"LAY_CHECKED": selected[r.RoleID],
		})
	}
	return data, rows.Err()
}

func (s *SystemService) SaveUserRole(ctx context.Context, userid int, roleIDs []int) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `DELETE FROM sys_role_user WHERE uid=?`, userid); err != nil {
		return err
	}
	for _, rid := range roleIDs {
		if _, err = tx.ExecContext(ctx, `INSERT INTO sys_role_user(uid,rid) VALUES(?,?)`, userid, rid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SystemService) ChangePassword(ctx context.Context, userid int, oldPassword, newPassword, confirmPassword string) ResultPassword {
	if newPassword != confirmPassword {
		return ResultPassword{Code: 500, Msg: "您输入的两次新密码不一致！"}
	}
	u, err := s.GetUserByID(ctx, userid)
	if err != nil || u == nil || u.Pwd == nil {
		return ResultPassword{Code: 500, Msg: "用户不存在！"}
	}
	if *u.Pwd != md5Hex(oldPassword) {
		return ResultPassword{Code: 500, Msg: "您输入的原密码错误！"}
	}
	if _, err := s.DB.ExecContext(ctx, `UPDATE sys_user SET pwd=? WHERE userid=?`, md5Hex(newPassword), userid); err != nil {
		return ResultPassword{Code: 500, Msg: "修改失败，请稍后重试！"}
	}
	return ResultPassword{Code: 200, Msg: "修改密码成功，请重新登陆！"}
}

type ResultPassword struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (s *SystemService) QueryRoles(ctx context.Context, q map[string]string) (int64, []model.Role, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := make([]any, 0, 8)
	if v := strings.TrimSpace(q["rolename"]); v != "" {
		where = append(where, "rolename LIKE ?")
		args = append(args, likeArg(v))
	}
	if v := strings.TrimSpace(q["roledesc"]); v != "" {
		where = append(where, "roledesc LIKE ?")
		args = append(args, likeArg(v))
	}
	if v := strings.TrimSpace(q["available"]); v != "" {
		where = append(where, "available=?")
		args = append(args, v)
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM sys_role WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT roleid,rolename,roledesc,available FROM sys_role WHERE `+wsql+` ORDER BY roleid ASC LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.Role
	for rows.Next() {
		var r model.Role
		if err := rows.Scan(&r.RoleID, &r.RoleName, &r.RoleDesc, &r.Available); err != nil {
			return 0, nil, err
		}
		out = append(out, r)
	}
	return count, out, rows.Err()
}

func (s *SystemService) AddRole(ctx context.Context, r model.Role) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO sys_role(rolename,roledesc,available) VALUES(?,?,?)`, r.RoleName, r.RoleDesc, r.Available)
	return err
}

func (s *SystemService) UpdateRole(ctx context.Context, r model.Role) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE sys_role SET rolename=?,roledesc=?,available=? WHERE roleid=?`, r.RoleName, r.RoleDesc, r.Available, r.RoleID)
	return err
}

func (s *SystemService) DeleteRole(ctx context.Context, roleid int) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	for _, q := range []string{
		`DELETE FROM sys_role WHERE roleid=?`,
		`DELETE FROM sys_role_menu WHERE rid=?`,
		`DELETE FROM sys_role_user WHERE rid=?`,
	} {
		if _, err = tx.ExecContext(ctx, q, roleid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SystemService) InitRoleMenuTree(ctx context.Context, roleid int) ([]model.MenuManagerNode, error) {
	allRows, err := s.DB.QueryContext(ctx, `SELECT id,pid,title,spread FROM sys_menu WHERE available=1 ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer allRows.Close()
	selected := map[int]bool{}
	rs, err := s.DB.QueryContext(ctx, `SELECT mid FROM sys_role_menu WHERE rid=?`, roleid)
	if err != nil {
		return nil, err
	}
	for rs.Next() {
		var mid int
		if err := rs.Scan(&mid); err != nil {
			_ = rs.Close()
			return nil, err
		}
		selected[mid] = true
	}
	_ = rs.Close()

	var out []model.MenuManagerNode
	for allRows.Next() {
		var id, pid, spread int
		var title string
		if err := allRows.Scan(&id, &pid, &title, &spread); err != nil {
			return nil, err
		}
		checkArr := "0"
		if selected[id] {
			checkArr = "1"
		}
		out = append(out, model.MenuManagerNode{
			ID:       id,
			ParentID: pid,
			Title:    title,
			Spread:   spread == 1,
			CheckArr: checkArr,
		})
	}
	return out, allRows.Err()
}

func (s *SystemService) SaveRoleMenu(ctx context.Context, roleid int, mids []int) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `DELETE FROM sys_role_menu WHERE rid=?`, roleid); err != nil {
		return err
	}
	for _, mid := range mids {
		if _, err = tx.ExecContext(ctx, `INSERT INTO sys_role_menu(rid,mid) VALUES(?,?)`, roleid, mid); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *SystemService) QueryLogInfos(ctx context.Context, q map[string]string) (int64, []model.LoginLog, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := []any{}
	if v := strings.TrimSpace(q["loginname"]); v != "" {
		where = append(where, "loginname LIKE ?")
		args = append(args, likeArg(v))
	}
	if v := strings.TrimSpace(q["loginip"]); v != "" {
		where = append(where, "loginip LIKE ?")
		args = append(args, likeArg(v))
	}
	if t, err := parseDateTime(q["startTime"]); err == nil && t != nil {
		where = append(where, "logintime >= ?")
		args = append(args, *t)
	}
	if t, err := parseDateTime(q["endTime"]); err == nil && t != nil {
		where = append(where, "logintime <= ?")
		args = append(args, *t)
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM sys_log_login WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT id,loginname,loginip,logintime FROM sys_log_login WHERE `+wsql+` ORDER BY logintime DESC LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.LoginLog
	for rows.Next() {
		var x model.LoginLog
		if err := rows.Scan(&x.ID, &x.LoginName, &x.LoginIP, &x.LoginTime); err != nil {
			return 0, nil, err
		}
		out = append(out, x)
	}
	return count, out, rows.Err()
}

func (s *SystemService) DeleteLogInfo(ctx context.Context, id int) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM sys_log_login WHERE id=?`, id)
	return err
}

func (s *SystemService) QueryNews(ctx context.Context, q map[string]string) (int64, []model.News, error) {
	return s.queryRichTextTable(ctx, "sys_news", q)
}

func (s *SystemService) QueryMessages(ctx context.Context, q map[string]string) (int64, []model.Message, error) {
	count, news, err := s.queryRichTextTable(ctx, "sys_message", q)
	if err != nil {
		return 0, nil, err
	}
	out := make([]model.Message, 0, len(news))
	for _, x := range news {
		out = append(out, model.Message(x))
	}
	return count, out, nil
}

func (s *SystemService) queryRichTextTable(ctx context.Context, table string, q map[string]string) (int64, []model.News, error) {
	_, limit, offset, err := parsePage(q["page"], q["limit"])
	if err != nil {
		return 0, nil, err
	}
	where := []string{"1=1"}
	args := []any{}
	if v := strings.TrimSpace(q["title"]); v != "" {
		where = append(where, "title LIKE ?")
		args = append(args, likeArg(v))
	}
	if v := strings.TrimSpace(q["content"]); v != "" {
		where = append(where, "content LIKE ?")
		args = append(args, likeArg(v))
	}
	if t, err := parseDateTime(q["startTime"]); err == nil && t != nil {
		where = append(where, "createtime >= ?")
		args = append(args, *t)
	}
	if t, err := parseDateTime(q["endTime"]); err == nil && t != nil {
		where = append(where, "createtime <= ?")
		args = append(args, *t)
	}
	wsql := strings.Join(where, " AND ")
	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM "+table+" WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	rows, err := s.DB.QueryContext(ctx, fmt.Sprintf(`SELECT id,title,content,createtime,opername FROM %s WHERE %s ORDER BY createtime DESC LIMIT ? OFFSET ?`, table, wsql), append(args, limit, offset)...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.News
	for rows.Next() {
		var x model.News
		if err := rows.Scan(&x.ID, &x.Title, &x.Content, &x.CreateTime, &x.OperName); err != nil {
			return 0, nil, err
		}
		out = append(out, x)
	}
	return count, out, rows.Err()
}

func (s *SystemService) AddNews(ctx context.Context, title, content, opername string) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO sys_news(title,content,createtime,opername) VALUES(?,?,?,?)`, title, content, time.Now(), opername)
	return err
}
func (s *SystemService) UpdateNews(ctx context.Context, id int, title, content string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE sys_news SET title=?,content=? WHERE id=?`, title, content, id)
	return err
}
func (s *SystemService) DeleteNews(ctx context.Context, id int) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM sys_news WHERE id=?`, id)
	return err
}
func (s *SystemService) GetNewsByID(ctx context.Context, id int) (*model.News, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT id,title,content,createtime,opername FROM sys_news WHERE id=?`, id)
	var x model.News
	if err := row.Scan(&x.ID, &x.Title, &x.Content, &x.CreateTime, &x.OperName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &x, nil
}

func (s *SystemService) AddMessage(ctx context.Context, title, content, opername string) error {
	_, err := s.DB.ExecContext(ctx, `INSERT INTO sys_message(title,content,createtime,opername) VALUES(?,?,?,?)`, title, content, time.Now(), opername)
	return err
}
func (s *SystemService) UpdateMessage(ctx context.Context, id int, title, content string) error {
	_, err := s.DB.ExecContext(ctx, `UPDATE sys_message SET title=?,content=? WHERE id=?`, title, content, id)
	return err
}
func (s *SystemService) DeleteMessage(ctx context.Context, id int) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM sys_message WHERE id=?`, id)
	return err
}
func (s *SystemService) GetMessageByID(ctx context.Context, id int) (*model.Message, error) {
	row := s.DB.QueryRowContext(ctx, `SELECT id,title,content,createtime,opername FROM sys_message WHERE id=?`, id)
	var x model.Message
	if err := row.Scan(&x.ID, &x.Title, &x.Content, &x.CreateTime, &x.OperName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &x, nil
}
