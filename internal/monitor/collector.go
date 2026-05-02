package monitor

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

type ctxKey string

const requestInfoKey ctxKey = "monitor_request_info"

type RequestInfo struct {
	mu             sync.Mutex
	URI            string
	JDBCExecute    int64
	JDBCError      int64
	JDBCTotalTime  time.Duration
}

func (r *RequestInfo) addSQL(duration time.Duration, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.JDBCExecute++
	r.JDBCTotalTime += duration
	if err != nil {
		r.JDBCError++
	}
}

func (r *RequestInfo) snapshot() (int64, int64, time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.JDBCExecute, r.JDBCError, r.JDBCTotalTime
}

func WithRequestInfo(ctx context.Context, info *RequestInfo) context.Context {
	return context.WithValue(ctx, requestInfoKey, info)
}

func RequestInfoFromContext(ctx context.Context) (*RequestInfo, bool) {
	v := ctx.Value(requestInfoKey)
	info, ok := v.(*RequestInfo)
	return info, ok
}

type DataSourceMeta struct {
	Name          string
	DBType        string
	JDBCURLMasked string
	UserMasked    string
}

type DataSourceJSON struct {
	Name                string `json:"name"`
	DBType              string `json:"dbType"`
	JDBCURLMasked       string `json:"jdbcUrlMasked"`
	UserMasked          string `json:"userMasked"`
	MaxOpenConnections  int    `json:"maxOpenConnections"`
	OpenConnections     int    `json:"openConnections"`
	InUse               int    `json:"inUse"`
	Idle                int    `json:"idle"`
	WaitCount           int64  `json:"waitCount"`
	WaitDurationMs      int64  `json:"waitDurationMs"`
	ConnectCount        int64  `json:"connectCount"`
	CloseCount          int64  `json:"closeCount"`
	ConnectErrorCount   int64  `json:"connectErrorCount"`
	ActivePeak          int    `json:"activePeak"`
	ActivePeakTime      string `json:"activePeakTime,omitempty"`
	LastUpdated         string `json:"lastUpdated"`
}

type SQLItem struct {
	ID            string  `json:"id"`
	SQL           string  `json:"sql"`
	ExecuteCount  int64   `json:"executeCount"`
	ErrorCount    int64   `json:"errorCount"`
	RunningCount  int64   `json:"runningCount"`
	ConcurrentMax int64   `json:"concurrentMax"`
	TotalTimeMs   float64 `json:"totalTimeMs"`
	AvgTimeMs     float64 `json:"avgTimeMs"`
	MaxTimeMs     float64 `json:"maxTimeMs"`
	LastErrorMsg  string  `json:"lastErrorMsg,omitempty"`
	LastErrorTime string  `json:"lastErrorTime,omitempty"`
	FirstSeen     string  `json:"firstSeen"`
	LastSeen      string  `json:"lastSeen"`
}

type SQLSample struct {
	TS           string  `json:"ts"`
	DurationMs   float64 `json:"durationMs"`
	OK           bool    `json:"ok"`
	ErrMsg       string  `json:"errMsg,omitempty"`
	RowsAffected int64   `json:"rowsAffected,omitempty"`
	URI          string  `json:"uri,omitempty"`
}

type SQLListJSON struct {
	Items    []SQLItem `json:"items"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
	SortBy   string    `json:"sortBy"`
	Query    string    `json:"query,omitempty"`
}

type SQLDetailJSON struct {
	Item    SQLItem     `json:"item"`
	Samples []SQLSample `json:"samples"`
}

type WebURIItem struct {
	URI            string  `json:"uri"`
	RequestCount   int64   `json:"requestCount"`
	ErrorCount     int64   `json:"errorCount"`
	RunningCount   int64   `json:"runningCount"`
	ConcurrentMax  int64   `json:"concurrentMax"`
	TotalTimeMs    float64 `json:"totalTimeMs"`
	AvgTimeMs      float64 `json:"avgTimeMs"`
	MaxTimeMs      float64 `json:"maxTimeMs"`
	JDBCExecute    int64   `json:"jdbcExecuteCount"`
	JDBCError      int64   `json:"jdbcErrorCount"`
	JDBCTotalTimeMs float64 `json:"jdbcTotalTimeMs"`
}

type WebURIJSON struct {
	Items    []WebURIItem `json:"items"`
	Total    int          `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}

type WebSessionItem struct {
	SessionKey     string `json:"sessionKey"`
	User           string `json:"user,omitempty"`
	CreateTime     string `json:"createTime"`
	LastAccessTime string `json:"lastAccessTime"`
	RequestCount   int64  `json:"requestCount"`
	RemoteAddr     string `json:"remoteAddr,omitempty"`
}

type WebSessionJSON struct {
	Items    []WebSessionItem `json:"items"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
}

type sqlSampleState struct {
	ts           time.Time
	duration     time.Duration
	ok           bool
	errMsg       string
	rowsAffected int64
	uri          string
}

type sqlState struct {
	id            string
	sql           string
	executeCount  int64
	errorCount    int64
	runningCount  int64
	concurrentMax int64
	totalTime     time.Duration
	maxTime       time.Duration
	lastErrorMsg  string
	lastErrorTime time.Time
	firstSeen     time.Time
	lastSeen      time.Time
	samples       []sqlSampleState
}

type uriState struct {
	uri           string
	requestCount  int64
	errorCount    int64
	runningCount  int64
	concurrentMax int64
	totalTime     time.Duration
	maxTime       time.Duration
	jdbcExecute   int64
	jdbcError     int64
	jdbcTotalTime time.Duration
}

type sessionState struct {
	sessionKey string
	user       string
	createTime time.Time
	lastAccess time.Time
	requestCount int64
	remoteAddr string
}

type Collector struct {
	meta       DataSourceMeta
	maxSamples int

	mu                sync.RWMutex
	connectCount      int64
	connectErrorCount int64
	closeCount        int64
	activePeak        int
	activePeakTime    time.Time
	sqlItems          map[string]*sqlState
	webURIs           map[string]*uriState
	webSessions       map[string]*sessionState
}

var (
	defaultMu        sync.RWMutex
	defaultCollector *Collector
	singleQuoteRe    = regexp.MustCompile(`'[^']*'`)
	doubleQuoteRe    = regexp.MustCompile(`"[^"]*"`)
	numberRe         = regexp.MustCompile(`\b\d+(\.\d+)?\b`)
	spaceRe          = regexp.MustCompile(`\s+`)
	equalsRe         = regexp.MustCompile(`\s*=\s*`)
)

func SetDefault(c *Collector) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultCollector = c
}

func Default() *Collector {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultCollector
}

func NewCollector(meta DataSourceMeta) *Collector {
	return &Collector{
		meta:        meta,
		maxSamples:  50,
		sqlItems:    make(map[string]*sqlState),
		webURIs:     make(map[string]*uriState),
		webSessions: make(map[string]*sessionState),
	}
}

func MetaFromDSN(dsn string) DataSourceMeta {
	meta := DataSourceMeta{
		Name:   "carRental",
		DBType: "mysql",
	}
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return meta
	}
	meta.UserMasked = maskUser(cfg.User)
	host := cfg.Addr
	if host == "" {
		host = "127.0.0.1:3306"
	}
	u := url.URL{
		Scheme: "jdbc:mysql",
		Host:   host,
		Path:   "/" + cfg.DBName,
	}
	meta.JDBCURLMasked = u.String()
	return meta
}

func (c *Collector) RecordConnect(err error) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		c.connectErrorCount++
		return
	}
	c.connectCount++
}

func (c *Collector) RecordClose() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closeCount++
}

func (c *Collector) RecordSQL(ctx context.Context, query string, duration time.Duration, err error, rowsAffected int64) {
	if c == nil {
		return
	}
	norm := NormalizeSQL(query)
	if norm == "" {
		norm = "<empty>"
	}
	id := SQLID(norm)
	now := time.Now()
	uri := ""
	if info, ok := RequestInfoFromContext(ctx); ok {
		info.addSQL(duration, err)
		uri = info.URI
	}

	c.mu.Lock()
	state, ok := c.sqlItems[id]
	if !ok {
		state = &sqlState{
			id:        id,
			sql:       norm,
			firstSeen: now,
			samples:   make([]sqlSampleState, 0, c.maxSamples),
		}
		c.sqlItems[id] = state
	}
	state.executeCount++
	state.totalTime += duration
	state.lastSeen = now
	if duration > state.maxTime {
		state.maxTime = duration
	}
	if err != nil {
		state.errorCount++
		state.lastErrorMsg = err.Error()
		state.lastErrorTime = now
	}
	state.samples = append(state.samples, sqlSampleState{
		ts:           now,
		duration:     duration,
		ok:           err == nil,
		errMsg:       errString(err),
		rowsAffected: rowsAffected,
		uri:          uri,
	})
	if len(state.samples) > c.maxSamples {
		state.samples = state.samples[len(state.samples)-c.maxSamples:]
	}
	c.mu.Unlock()
}

func (c *Collector) BeginSQL(query string) func(context.Context, error, int64) {
	if c == nil {
		return func(context.Context, error, int64) {}
	}
	norm := NormalizeSQL(query)
	if norm == "" {
		norm = "<empty>"
	}
	id := SQLID(norm)
	start := time.Now()

	c.mu.Lock()
	state, ok := c.sqlItems[id]
	if !ok {
		state = &sqlState{
			id:        id,
			sql:       norm,
			firstSeen: start,
			samples:   make([]sqlSampleState, 0, c.maxSamples),
		}
		c.sqlItems[id] = state
	}
	state.runningCount++
	if state.runningCount > state.concurrentMax {
		state.concurrentMax = state.runningCount
	}
	c.mu.Unlock()

	return func(ctx context.Context, err error, rowsAffected int64) {
		c.mu.Lock()
		if state := c.sqlItems[id]; state != nil && state.runningCount > 0 {
			state.runningCount--
		}
		c.mu.Unlock()
		c.RecordSQL(ctx, query, time.Since(start), err, rowsAffected)
	}
}

func (c *Collector) DataSourceSnapshot(db *sql.DB) DataSourceJSON {
	st := db.Stats()
	c.mu.Lock()
	if st.InUse > c.activePeak {
		c.activePeak = st.InUse
		c.activePeakTime = time.Now()
	}
	activePeak := c.activePeak
	activePeakTime := c.activePeakTime
	connectCount := c.connectCount
	connectErrorCount := c.connectErrorCount
	closeCount := c.closeCount
	c.mu.Unlock()

	return DataSourceJSON{
		Name:               c.meta.Name,
		DBType:             c.meta.DBType,
		JDBCURLMasked:      c.meta.JDBCURLMasked,
		UserMasked:         c.meta.UserMasked,
		MaxOpenConnections: st.MaxOpenConnections,
		OpenConnections:    st.OpenConnections,
		InUse:              st.InUse,
		Idle:               st.Idle,
		WaitCount:          st.WaitCount,
		WaitDurationMs:     st.WaitDuration.Milliseconds(),
		ConnectCount:       connectCount,
		CloseCount:         closeCount,
		ConnectErrorCount:  connectErrorCount,
		ActivePeak:         activePeak,
		ActivePeakTime:     formatTime(activePeakTime),
		LastUpdated:        formatTime(time.Now()),
	}
}

func (c *Collector) RecordRequestStart(uri string) {
	if c == nil || uri == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	state, ok := c.webURIs[uri]
	if !ok {
		state = &uriState{uri: uri}
		c.webURIs[uri] = state
	}
	state.runningCount++
	if state.runningCount > state.concurrentMax {
		state.concurrentMax = state.runningCount
	}
}

func (c *Collector) RecordRequestFinish(uri string, status int, duration time.Duration, info *RequestInfo) {
	if c == nil || uri == "" {
		return
	}
	jdbcExec, jdbcErr, jdbcTotal := int64(0), int64(0), time.Duration(0)
	if info != nil {
		jdbcExec, jdbcErr, jdbcTotal = info.snapshot()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	state, ok := c.webURIs[uri]
	if !ok {
		state = &uriState{uri: uri}
		c.webURIs[uri] = state
	}
	state.requestCount++
	if status >= 400 {
		state.errorCount++
	}
	state.totalTime += duration
	if duration > state.maxTime {
		state.maxTime = duration
	}
	if state.runningCount > 0 {
		state.runningCount--
	}
	state.jdbcExecute += jdbcExec
	state.jdbcError += jdbcErr
	state.jdbcTotalTime += jdbcTotal
}

func (c *Collector) RecordSession(sessionKey, user, remoteAddr string) {
	if c == nil || sessionKey == "" {
		return
	}
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	state, ok := c.webSessions[sessionKey]
	if !ok {
		state = &sessionState{
			sessionKey: sessionKey,
			createTime: now,
		}
		c.webSessions[sessionKey] = state
	}
	state.lastAccess = now
	state.requestCount++
	if user != "" {
		state.user = user
	}
	if remoteAddr != "" {
		state.remoteAddr = remoteAddr
	}
}

func (c *Collector) SQLList(page, pageSize int, query, sortBy string) SQLListJSON {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	c.mu.RLock()
	items := make([]SQLItem, 0, len(c.sqlItems))
	for _, v := range c.sqlItems {
		item := toSQLItem(v)
		if query == "" || strings.Contains(strings.ToLower(item.SQL), strings.ToLower(query)) {
			items = append(items, item)
		}
	}
	c.mu.RUnlock()
	sortSQLItems(items, sortBy)
	total := len(items)
	items = paginateSQLItems(items, page, pageSize)
	return SQLListJSON{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		Query:    query,
	}
}

func (c *Collector) SQLDetail(id string) (SQLDetailJSON, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	state, ok := c.sqlItems[id]
	if !ok {
		return SQLDetailJSON{}, false
	}
	samples := make([]SQLSample, 0, len(state.samples))
	for i := len(state.samples) - 1; i >= 0; i-- {
		s := state.samples[i]
		samples = append(samples, SQLSample{
			TS:           formatTime(s.ts),
			DurationMs:   durationMs(s.duration),
			OK:           s.ok,
			ErrMsg:       s.errMsg,
			RowsAffected: s.rowsAffected,
			URI:          s.uri,
		})
	}
	return SQLDetailJSON{
		Item:    toSQLItem(state),
		Samples: samples,
	}, true
}

func (c *Collector) WebURIList(page, pageSize int, query string) WebURIJSON {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	c.mu.RLock()
	items := make([]WebURIItem, 0, len(c.webURIs))
	for _, v := range c.webURIs {
		if query == "" || strings.Contains(strings.ToLower(v.uri), strings.ToLower(query)) {
			items = append(items, WebURIItem{
				URI:             v.uri,
				RequestCount:    v.requestCount,
				ErrorCount:      v.errorCount,
				RunningCount:    v.runningCount,
				ConcurrentMax:   v.concurrentMax,
				TotalTimeMs:     durationMs(v.totalTime),
				AvgTimeMs:       avgDurationMs(v.totalTime, v.requestCount),
				MaxTimeMs:       durationMs(v.maxTime),
				JDBCExecute:     v.jdbcExecute,
				JDBCError:       v.jdbcError,
				JDBCTotalTimeMs: durationMs(v.jdbcTotalTime),
			})
		}
	}
	c.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].TotalTimeMs > items[j].TotalTimeMs
	})
	total := len(items)
	start, end := paginate(total, page, pageSize)
	return WebURIJSON{
		Items:    items[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

func (c *Collector) WebSessionList(page, pageSize int) WebSessionJSON {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	c.mu.RLock()
	items := make([]WebSessionItem, 0, len(c.webSessions))
	for _, v := range c.webSessions {
		items = append(items, WebSessionItem{
			SessionKey:     v.sessionKey,
			User:           v.user,
			CreateTime:     formatTime(v.createTime),
			LastAccessTime: formatTime(v.lastAccess),
			RequestCount:   v.requestCount,
			RemoteAddr:     v.remoteAddr,
		})
	}
	c.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].LastAccessTime > items[j].LastAccessTime
	})
	total := len(items)
	start, end := paginate(total, page, pageSize)
	return WebSessionJSON{
		Items:    items[start:end],
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

func (c *Collector) ResetSQL() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sqlItems = make(map[string]*sqlState)
}

func (c *Collector) ResetWeb() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.webURIs = make(map[string]*uriState)
	c.webSessions = make(map[string]*sessionState)
}

func (c *Collector) ResetDataSource() {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connectCount = 0
	c.connectErrorCount = 0
	c.closeCount = 0
	c.activePeak = 0
	c.activePeakTime = time.Time{}
}

func (c *Collector) ResetAll() {
	c.ResetSQL()
	c.ResetWeb()
	c.ResetDataSource()
}

func HashSessionKey(raw string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(raw))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}

func SQLID(s string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum64())
}

func NormalizeSQL(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	s = singleQuoteRe.ReplaceAllString(s, "?")
	s = doubleQuoteRe.ReplaceAllString(s, "?")
	s = numberRe.ReplaceAllString(s, "?")
	s = spaceRe.ReplaceAllString(s, " ")
	s = equalsRe.ReplaceAllString(s, "=")
	return strings.TrimSpace(s)
}

func toSQLItem(v *sqlState) SQLItem {
	return SQLItem{
		ID:            v.id,
		SQL:           v.sql,
		ExecuteCount:  v.executeCount,
		ErrorCount:    v.errorCount,
		RunningCount:  v.runningCount,
		ConcurrentMax: v.concurrentMax,
		TotalTimeMs:   durationMs(v.totalTime),
		AvgTimeMs:     avgDurationMs(v.totalTime, v.executeCount),
		MaxTimeMs:     durationMs(v.maxTime),
		LastErrorMsg:  v.lastErrorMsg,
		LastErrorTime: formatTime(v.lastErrorTime),
		FirstSeen:     formatTime(v.firstSeen),
		LastSeen:      formatTime(v.lastSeen),
	}
}

func sortSQLItems(items []SQLItem, sortBy string) {
	switch sortBy {
	case "executeCount":
		sort.Slice(items, func(i, j int) bool { return items[i].ExecuteCount > items[j].ExecuteCount })
	case "maxTimeMs":
		sort.Slice(items, func(i, j int) bool { return items[i].MaxTimeMs > items[j].MaxTimeMs })
	case "errorCount":
		sort.Slice(items, func(i, j int) bool { return items[i].ErrorCount > items[j].ErrorCount })
	default:
		sort.Slice(items, func(i, j int) bool { return items[i].TotalTimeMs > items[j].TotalTimeMs })
	}
}

func paginateSQLItems(items []SQLItem, page, pageSize int) []SQLItem {
	start, end := paginate(len(items), page, pageSize)
	return items[start:end]
}

func paginate(total, page, pageSize int) (int, int) {
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return start, end
}

func durationMs(d time.Duration) float64 {
	return float64(d.Microseconds()) / 1000
}

func avgDurationMs(total time.Duration, count int64) float64 {
	if count <= 0 {
		return 0
	}
	return durationMs(total) / float64(count)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func maskUser(user string) string {
	user = strings.TrimSpace(user)
	if user == "" {
		return ""
	}
	if len(user) <= 2 {
		return user[:1] + "*"
	}
	return user[:1] + strings.Repeat("*", len(user)-2) + user[len(user)-1:]
}
