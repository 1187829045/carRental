package service

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DBInit struct {
	ProjectRoot string
}

func (d DBInit) EnsureSchemaAndSeed(ctx context.Context, db *sql.DB) error {
	extraSQL := filepath.Join(d.ProjectRoot, "scripts", "extra.sql")
	if err := execSQLFileIfExists(ctx, db, extraSQL); err != nil && !os.IsNotExist(err) {
		return err
	}

	has, err := hasTable(ctx, db, "sys_menu")
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	baseSQL := filepath.Join(d.ProjectRoot, "sql")
	if err := execSQLFileIfExists(ctx, db, baseSQL); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("db init: init sql file not found: %s", baseSQL)
		}
		return err
	}
	if err := execSQLFileIfExists(ctx, db, extraSQL); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func hasTable(ctx context.Context, db *sql.DB, name string) (bool, error) {
	var n int
	err := db.QueryRowContext(ctx, `SELECT COUNT(1) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?`, name).Scan(&n)
	return n > 0, err
}

func execSQLFile(ctx context.Context, db *sql.DB, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	statements, err := parseSQLStatements(f)
	if err != nil {
		return err
	}
	for _, s := range statements {
		if _, err := db.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("exec %s: %w", shortSQL(s), err)
		}
	}
	return nil
}

func execSQLFileIfExists(ctx context.Context, db *sql.DB, filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return err
	}
	return execSQLFile(ctx, db, filePath)
}

func parseSQLStatements(r *os.File) ([]string, error) {
	var b strings.Builder
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if strings.HasPrefix(trim, "--") {
			continue
		}
		b.WriteString(line)
		b.WriteString("\n")
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	clean := stripBlockComments(b.String())
	return splitSQL(clean), nil
}

func stripBlockComments(s string) string {
	for {
		i := strings.Index(s, "/*")
		if i < 0 {
			return s
		}
		j := strings.Index(s[i+2:], "*/")
		if j < 0 {
			return s[:i]
		}
		s = s[:i] + s[i+2+j+2:]
	}
}

func splitSQL(s string) []string {
	var out []string
	var cur strings.Builder
	inSingle := false
	esc := false
	for _, ch := range s {
		if esc {
			cur.WriteRune(ch)
			esc = false
			continue
		}
		if ch == '\\' {
			cur.WriteRune(ch)
			esc = true
			continue
		}
		if ch == '\'' {
			inSingle = !inSingle
			cur.WriteRune(ch)
			continue
		}
		if ch == ';' && !inSingle {
			stmt := strings.TrimSpace(cur.String())
			if stmt != "" {
				out = append(out, stmt)
			}
			cur.Reset()
			continue
		}
		cur.WriteRune(ch)
	}
	stmt := strings.TrimSpace(cur.String())
	if stmt != "" {
		out = append(out, stmt)
	}
	return out
}

func shortSQL(s string) string {
	s = strings.TrimSpace(s)
	if len(s) <= 80 {
		return s
	}
	return s[:80]
}
