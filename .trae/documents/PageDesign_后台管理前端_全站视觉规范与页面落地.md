# 全站视觉规范与页面落地（Desktop-first）

## Global Styles（全站设计令牌）
- 色彩（建议与现有 `--ui-*` 对齐并补齐语义色）
  - 主色 Primary：`--ui-primary`（按钮/链接/选中态）
  - 主色渐变：`--ui-primary-2`（顶部栏渐变）
  - 背景：`--ui-bg`（页面底色） / `--ui-card`（卡片/容器底色）
  - 边框：`--ui-border`
  - 文本：`--ui-text`（主文本） / `--ui-muted`（辅助文本）
  - 语义色：Success/Warning/Danger/Info（用于提示、状态标签、危险按钮）
- 字体与排版
  - 基准字号：14px；标题 16/18px；数字/关键指标可 20–28px。
  - 行高：1.4–1.6；表格行高与表单控件高度统一（建议 32/36px 两档）。
- 圆角与阴影
  - 圆角：容器 12px、控件 10px、标签 8px。
  - 阴影：卡片轻阴影（与现有 `--ui-shadow` 风格一致）。
- 交互状态
  - Hover：浅底高亮（同主色 6%–12% 透明度）。
  - Focus：主色描边 + 轻投影，确保可见。
  - Disabled：降低不透明度 + 禁用光标，不改变布局。

## Layout（布局原则）
- 框架页采用“顶部栏 + 左侧导航 + 右侧内容区（tab + iframe）”的既有布局，不改变结构。
- 内容区统一“卡片化容器”：筛选区、表格区、详情区用卡片分组，卡片间距 12–16px。
- 响应式：桌面优先；移动端沿用现有媒体查询，确保不新增破坏性断点。

---

## Page 1：登录页（`system/main/login.jsp`）
### Meta Information
- Title：登录--汽车出租系统（保持不变）
- Description：后台登录页（可补充但不强制）

### Page Structure
- 全屏背景（可继续使用现有背景图），中心登录卡片。

### Sections & Components
1. 登录卡片
   - 宽度：300–360px；内边距：20–24px；圆角：12px；阴影：中等。
2. 头像区
   - 保持现有圆形头像；边框与阴影统一为主色系轻阴影。
3. 表单项（用户名/密码/验证码）
   - 输入框：统一高度、边框色、focus 反馈。
   - Label 浮动效果保持现有 JS 逻辑，仅调整颜色与字号。
4. 主按钮
   - 主按钮使用主色；hover 更深；disabled 明确。
5. 错误提示区
   - `${error}` 区域使用 Danger 语义色与更清晰的间距。

---

## Page 2：后台框架页（`system/main/index.jsp`）
### Meta Information
- Title：首页-汽车租赁系统（保持不变）

### Page Structure
- 顶部栏：左 logo + 折叠菜单 + 右侧用户菜单。
- 侧边栏：菜单树（Layui nav-tree）。
- 内容区：tab 标题栏 + iframe 内容区 + 右侧“页面操作”下拉。

### Sections & Components
1. 顶部栏（Header）
   - 背景：主色渐变；高度维持 50px（现有）。
   - 用户菜单：下拉面板圆角与阴影统一；hover 背景与文本对比清晰。
2. 侧边栏（Side Nav）
   - 背景卡片化：浅底 + 右侧分割线；选中态使用主色引导条。
   - 菜单项：统一 padding、圆角；hover/active/展开态一致。
3. Tab 栏与“页面操作”
   - Tab 选中态：浅主色底 + 主色文字；关闭按钮 hover 显示危险态。
   - 操作下拉：列表项 hover 与分割线统一。
4. 内容背景
   - 页面底色使用 `--ui-bg`；iframe 内页面建议统一 `.childrenBody` 容器卡片化。

---

## Page 3：工作台/首页（`desk/toDeskManager.action` 对应页面）
### Page Structure
- 欢迎区 + 概览卡片 + 快捷入口/公告/留言等模块（遵循 `docs/home-dashboard-redesign.md` 的卡片化方向）。

### Sections & Components
1. 欢迎横幅
   - 左信息、右用户卡；主按钮与次按钮区分明确。
2. 概览卡片
   - 四卡一致高度；数字字体更大；图标统一风格与配色。
3. 列表模块（公告/留言）
   - 统一列表行高、hover、空态；标题与“更多”入口样式统一。

---

## Page 4：通用列表/表单页（覆盖主要业务/系统页面）
> 目标：让所有“筛选 + 表格 + 弹窗表单”页面呈现一致，不要求逐页重构。

### Layout
- 顶部筛选区：一行/两行自适应；按钮组靠右。
- 表格区：卡片容器内放 Layui table；分页与工具条对齐。

### Components（需覆盖的公共组件）
1. 按钮：`.layui-btn` 主次/危险/禁用/小号尺寸
2. 表单：`.layui-input`、`.layui-select`、`.layui-textarea`、校验态
3. 表格：表头、行 hover、斑马纹、空数据提示、操作列按钮密度
4. 分页：当前页、hover、禁用态
5. 弹窗（Layer）：标题区、关闭按钮、底部按钮区、危险操作确认
6. 提示：成功/警告/错误 toast 的颜色与 icon

---

## Page 5：数据源监控（Druid 登录/看板）
### Page Structure
- 登录页：居中模态 + 遮罩（保持 `docs/ux-module-improvements.md` 描述的体验）。
- 看板页：仅做外层背景/字体/容器统一，避免覆盖 Druid 内部控件导致功能风险。

### Interaction
- 未登录跳转、登录/退出流程保持不变；仅视觉层（CSS）介入。
