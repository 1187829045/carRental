package monitor

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestNormalizeSQLAndIDStable(t *testing.T) {
	sql1 := "SELECT *  FROM sys_user WHERE userid = 123 AND loginname = 'admin'"
	sql2 := "SELECT * FROM sys_user WHERE userid=999 AND loginname='guest'"
	n1 := NormalizeSQL(sql1)
	n2 := NormalizeSQL(sql2)
	if n1 == "" || n2 == "" {
		t.Fatalf("normalized sql should not be empty")
	}
	if SQLID(n1) != SQLID(n2) {
		t.Fatalf("expected stable sql id for normalized statements: %q vs %q", n1, n2)
	}
}

func TestCollectorAggregatesSQLWebAndDatasource(t *testing.T) {
	c := NewCollector(DataSourceMeta{Name: "test", DBType: "mysql"})
	info := &RequestInfo{URI: "/bus/toCarManager.action"}
	ctx := WithRequestInfo(context.Background(), info)

	done := c.BeginSQL("SELECT * FROM sys_user WHERE userid = 1")
	done(ctx, nil, 1)
	doneErr := c.BeginSQL("SELECT * FROM sys_user WHERE userid = 2")
	doneErr(ctx, context.DeadlineExceeded, 0)

	sqlList := c.SQLList(1, 10, "", "totalTimeMs")
	if sqlList.Total != 1 {
		t.Fatalf("expected merged sql count 1, got %d", sqlList.Total)
	}
	if sqlList.Items[0].ExecuteCount != 2 || sqlList.Items[0].ErrorCount != 1 {
		t.Fatalf("unexpected sql aggregate: %+v", sqlList.Items[0])
	}

	c.RecordRequestStart("/bus/toCarManager.action")
	c.RecordRequestFinish("/bus/toCarManager.action", 200, 120*time.Millisecond, info)
	uriList := c.WebURIList(1, 10, "")
	if uriList.Total != 1 || uriList.Items[0].JDBCExecute != 2 {
		t.Fatalf("unexpected web uri aggregate: %+v", uriList.Items)
	}

	c.RecordSession(HashSessionKey("abc"), "admin", "127.0.0.1")
	sessionList := c.WebSessionList(1, 10)
	if sessionList.Total != 1 || sessionList.Items[0].User != "admin" {
		t.Fatalf("unexpected session aggregate: %+v", sessionList.Items)
	}

	db := &sql.DB{}
	_ = c.DataSourceSnapshot(db)
}
