# 数据源监控/检查单/菜单管理 改进说明

## 1. 数据源监控：独立登录页 + 居中模态效果

### 目标
- 不再使用浏览器 BasicAuth 弹窗。
- 点击“数据源监控”后进入独立登录页。
- 登录框在视窗水平/垂直居中，背景半透明遮罩，符合系统蓝白风格。

### 实现
- 新增页面：`system/druid/druidLogin.jsp`
- 新增路由：
  - `GET  /carRental/druid/toLogin.action`：渲染登录页
  - `POST /carRental/druid/login.action`：校验账号密码，写入 session 标记
  - `POST /carRental/druid/logout.action`：清理 session 标记
  - `GET  /carRental/druid/`：未登录则跳转到登录页，登录后展示监控页
- 账号密码默认读取 `db.properties` 的 `jdbc.user/jdbc.password`（可通过环境变量 `MONITOR_USER/MONITOR_PASS` 覆盖）。

相关文件：
- `internal/handler/router.go`
- `internal/config/config.go`
- `src/main/webapp/WEB-INF/view/system/druid/druidLogin.jsp`

## 2. 检查单页面：修复初始加载自动滚动

### 现象
- 页面初次加载时内容自动向下滚动，导致用户进入页面后不在顶部。

### 修复
- 取消 LayUI 表格固定高度配置，避免在 iframe/嵌套布局下触发高度计算导致的滚动异常。
- 页面加载完成后显式回到顶部，并关闭浏览器自动滚动恢复。

相关文件：
- `src/main/webapp/WEB-INF/view/business/check/checkManager.jsp`

## 3. 菜单管理：左侧 error not found 修复与友好提示

### 现象
- 左侧菜单树出现 “error not found”，菜单数据无法加载。

### 根因
- 前端请求依赖 `${yeqifu}`（默认 `/carRental`），后端需同时兼容 `/carRental/*` 路径。

### 修复
- 后端按 `BASE_PATH` 挂载接口与静态资源。
- 前端在 Ajax 失败时显示友好提示，并替换树容器内容，避免展示技术错误信息。

相关文件：
- `internal/handler/router.go`
- `src/main/webapp/WEB-INF/view/system/menu/menuLeft.jsp`
- `src/main/webapp/WEB-INF/view/system/menu/menuRight.jsp`

## 验证

### 单测
```bash
cd /Users/bytedance/Desktop/carRental
go test ./...
```

### 本地运行
```bash
cd /Users/bytedance/Desktop/carRental
ADDR=:8081 BASE_PATH=/carRental go run ./cmd/goserver
```

访问：
- 菜单管理：系统管理 → 菜单管理
- 数据源监控：系统管理 → 数据源监控（会进入独立登录页）
- 检查单管理：业务管理 → 检查单管理（进入页面应停留顶部）

