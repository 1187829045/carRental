package service

import (
	"context"
	"database/sql"
	"strings"

	"carRental/internal/model"
)

type MenuService struct {
	DB *sql.DB
}

func (s *MenuService) QueryMenuList(ctx context.Context, title string, available *int, id *int, page, limit int) (int64, []model.Menu, error) {
	where := []string{"1=1"}
	args := make([]any, 0, 8)
	if title != "" {
		where = append(where, "title LIKE ?")
		args = append(args, likeArg(title))
	}
	if available != nil {
		where = append(where, "available=?")
		args = append(args, *available)
	}
	if id != nil {
		where = append(where, "(id=? OR pid=?)")
		args = append(args, *id, *id)
	}
	wsql := strings.Join(where, " AND ")

	var count int64
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(1) FROM sys_menu WHERE "+wsql, args...).Scan(&count); err != nil {
		return 0, nil, err
	}
	offset := (page - 1) * limit
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id,pid,title,href,spread,target,icon,available
		 FROM sys_menu WHERE `+wsql+` ORDER BY id ASC LIMIT ? OFFSET ?`,
		append(args, limit, offset)...,
	)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()
	var out []model.Menu
	for rows.Next() {
		var m model.Menu
		if err := rows.Scan(&m.ID, &m.Pid, &m.Title, &m.Href, &m.Spread, &m.Target, &m.Icon, &m.Available); err != nil {
			return 0, nil, err
		}
		out = append(out, m)
	}
	return count, out, rows.Err()
}

func (s *MenuService) AddMenu(ctx context.Context, m model.Menu) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO sys_menu(pid,title,href,spread,target,icon,available) VALUES(?,?,?,?,?,?,?)`,
		m.Pid, m.Title, m.Href, m.Spread, m.Target, m.Icon, m.Available)
	return err
}

func (s *MenuService) UpdateMenu(ctx context.Context, m model.Menu) error {
	_, err := s.DB.ExecContext(ctx,
		`UPDATE sys_menu SET pid=?,title=?,href=?,spread=?,target=?,icon=?,available=? WHERE id=?`,
		m.Pid, m.Title, m.Href, m.Spread, m.Target, m.Icon, m.Available, m.ID)
	return err
}

func (s *MenuService) CountByPid(ctx context.Context, pid int) (int, error) {
	var n int
	err := s.DB.QueryRowContext(ctx, `SELECT COUNT(1) FROM sys_menu WHERE pid=?`, pid).Scan(&n)
	return n, err
}

func (s *MenuService) DeleteMenu(ctx context.Context, id int) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.ExecContext(ctx, `DELETE FROM sys_menu WHERE id=?`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM sys_role_menu WHERE mid=?`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *MenuService) QueryAllMenus(ctx context.Context, available int) ([]model.Menu, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, pid, title, href, spread, target, icon, available
		 FROM sys_menu
		 WHERE (? = -1 OR available = ?)
		 ORDER BY id ASC`,
		available, available,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Menu
	for rows.Next() {
		var m model.Menu
		if err := rows.Scan(&m.ID, &m.Pid, &m.Title, &m.Href, &m.Spread, &m.Target, &m.Icon, &m.Available); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *MenuService) QueryMenusByUID(ctx context.Context, available int, uid int) ([]model.Menu, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT DISTINCT t1.id, t1.pid, t1.title, t1.href, t1.spread, t1.target, t1.icon, t1.available
		 FROM sys_menu t1
		 INNER JOIN sys_role_menu t2 ON (t1.id = t2.mid)
		 INNER JOIN sys_role_user t3 ON (t2.rid = t3.rid)
		 WHERE t3.uid = ? AND t1.available = ?
		 ORDER BY t1.id ASC`,
		uid, available,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Menu
	for rows.Next() {
		var m model.Menu
		if err := rows.Scan(&m.ID, &m.Pid, &m.Title, &m.Href, &m.Spread, &m.Target, &m.Icon, &m.Available); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func MenusToTree(menus []model.Menu, rootID int) []model.TreeNode {
	nodes := make([]model.TreeNode, 0, len(menus))
	for _, m := range menus {
		nodes = append(nodes, model.TreeNode{
			ID:       m.ID,
			ParentID: m.Pid,
			Title:    m.Title,
			Icon:     m.Icon,
			Href:     m.Href,
			Spread:   m.Spread == 1,
			Target:   m.Target,
			CheckArr: "0",
		})
	}

	byPID := make(map[int][]model.TreeNode, len(nodes))
	for _, n := range nodes {
		byPID[n.ParentID] = append(byPID[n.ParentID], n)
	}

	var build func(pid int) []model.TreeNode
	build = func(pid int) []model.TreeNode {
		children := byPID[pid]
		for i := range children {
			children[i].Child = build(children[i].ID)
		}
		return children
	}
	return build(rootID)
}
