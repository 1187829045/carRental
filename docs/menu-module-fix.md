# 菜单管理模块修复说明

## 背景
菜单管理模块由 `menuManager.jsp` 通过 `frameset` 加载左侧树（`menuLeft.jsp`）与右侧表格（`menuRight.jsp`）组成。

## 问题 1：接口调用出现 “not found”

### 现象
- 菜单树/菜单表格请求接口时出现 “not found”，页面无法加载菜单数据。

### 根因
- Go 后端渲染的 JSP 中 `${yeqifu}` 默认为 `/carRental`，前端请求接口与静态资源路径为：`/carRental/...`。
- 但 Go 后端路由与静态资源仅挂载在根路径：`/menu/...` 与 `/static/...`，导致前端请求 `404`。

### 修复
- Go 后端按 `BASE_PATH`（默认 `/carRental`）将接口与页面路由统一挂载到 `/carRental` 下。
- 静态资源同时支持 `/static/**` 与 `/carRental/static/**`。
- `requireLogin` 中间件支持带 `basePath` 的登录白名单与重定向。

相关代码：
- `internal/handler/router.go`
- `internal/config/config.go`

## 问题 2：分页栏位置异常

### 现象
- 菜单管理右侧表格分页栏出现位置偏移/显示异常。

### 根因
- LayUI 表格在设置固定高度（`height: 'full-148'`）时，会启用内部滚动与绝对定位分页区域；在 `frameset/iframe` 环境下，窗口高度计算可能不稳定，导致分页区域位置异常。

### 修复
- 移除固定高度配置，让表格使用自然高度布局，分页栏跟随表格正常排版。

相关代码：
- `src/main/webapp/WEB-INF/view/system/menu/menuRight.jsp`

## 额外增强
- 在 `menuLeft.jsp` / `menuRight.jsp` 增加 `ajaxError` 提示，接口请求失败时给出可见提示。

## 本地验证步骤
1. 启动 MySQL 并保证 `carRental` 数据库可连接。
2. 启动 Go 服务：

```bash
cd /Users/bytedance/Desktop/carRental
ADDR=:8081 BASE_PATH=/carRental go run ./cmd/goserver
```

3. 浏览器访问：
- `http://localhost:8081/carRental/login/toLogin.action`
- 登录后进入：系统管理 → 菜单管理

验证点：
- 左侧菜单树正常加载
- 右侧表格正常分页、翻页控件位置稳定
- 浏览器 Network 中接口为 `200`，无 `404 not found`

## 测试环境验证建议
- 测试环境按同样方式设置 `BASE_PATH=/carRental`（或按部署实际 contextPath 配置）。
- 回归检查：客户管理、车辆管理、出租/入库/检查单、统计分析页面的静态资源与接口路径均正常。

