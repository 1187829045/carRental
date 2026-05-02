# Page Design — Go Druid Monitor Parity (Desktop-first)

## Global Styles (applies to all monitor pages)
- Layout system: **CSS Grid for main layout** + Flexbox for toolbars and table controls.
- Theme tokens (match existing admin look):
  - Background: `#F6F7FB` (page), `#FFFFFF` (cards)
  - Text: `#1F2937`, muted `#6B7280`
  - Accent: `#1677FF` (primary), danger `#DC2626`, warning `#D97706`, success `#16A34A`
  - Typography: base 14px; headings 16–20px; tables 12–13px
  - Buttons: primary filled; secondary outline; hover adds subtle shadow
- Responsive behavior:
  - Desktop-first at 1280px width.
  - At <1024px, tables become horizontally scrollable; filters collapse into a single row.

---

## 1) Monitor Login
### Layout
- Centered card (max-width 420px) with subtle shadow.

### Meta Information
- Title: "DB Monitor Login"
- Description: "Login to database monitoring dashboard"
- Open Graph: same as title/description; no image required.

### Page Structure
1. Header: product name "DB Monitor" + short subtitle.
2. Form card:
   - Username input
   - Password input
   - Submit button
   - Error message area (inline)
3. Footer: small security notice (optional).

### Interaction states
- Submit disabled while posting.
- On failure: keep username, clear password.

---

## 2) Monitor Dashboard
### Layout
- Top-level grid:
  - Row 1: sticky toolbar
  - Row 2: tab bar (DataSource / SQL / Web URI / Web Session)
  - Row 3: content area (table + filters)

### Meta Information
- Title: "DB Monitor"
- Description: "Druid-like datasource, SQL and web monitoring"

### Sections & Components

#### A. Sticky Toolbar
- Left: title + environment badge (optional).
- Right:
  - Auto-refresh toggle (on/off)
  - Interval select (2s/5s/10s/30s)
  - “Last updated” timestamp
  - Reset dropdown: Reset All / Reset SQL / Reset Web / Reset DataSource
  - Logout link/button

#### B. Tabs
- Tabs switch the visible panel without navigation (single-page feel).
- Each tab remembers its last filter/sort state.

#### C. DataSource Panel
- Summary cards (4–6 small cards): OpenConnections, InUse, Idle, WaitCount, WaitDuration, ConnectErrors.
- Detail table (single row for current datasource, but future-proof for multiple):
  - Columns: Name, URL (masked), User (masked), MaxActive, InUse, Idle, ActivePeak, WaitCount, WaitDuration, ConnectCount, CloseCount, ConnectErrorCount.
- Visual alerts:
  - InUse close to MaxActive → warning.
  - WaitDuration rising fast → warning.
  - ConnectErrorCount increases → danger.

#### D. SQL Panel
- Filter bar:
  - Search input (SQL substring)
  - Sort select (Total Time / Exec Count / Max Time / Error Count)
  - Page size select
- Table:
  - Columns: SQL ID, SQL (truncate + expand), Exec Count, Error Count, Total Time, Avg Time, Max Time, Running, Concurrent Max, First Seen, Last Seen.
  - Row click → opens SQL Detail (navigate or modal drawer).

#### E. Web URI Panel
- Filter bar:
  - Search input (URI)
  - Sort select (Total Time / Req Count / Error Count / Max Time)
- Table:
  - Columns: URI, Req Count, Error Count, Total Time, Avg Time, Max Time, Running, Concurrent Max, JDBC Exec Count, JDBC Error Count, JDBC Total Time.

#### F. Web Session Panel
- Table:
  - Columns: Session Key (hashed), User, Create Time, Last Access, Request Count, Remote Addr.
- Optional action column:
  - “Expire” (if implemented) with confirm.

### Realtime Refresh
- Default: polling every 5 seconds when auto-refresh is on.
- UX rules:
  - Do not reset scroll position on refresh.
  - Keep current filters/sort/page.
  - If refresh fails, show non-blocking banner and exponential backoff.

---

## 3) SQL Detail
### Layout
- Two-column grid:
  - Left (65%): SQL text + charts/tables
  - Right (35%): key metrics and last error

### Meta Information
- Title: "SQL Detail"
- Description: "Drilldown for a single SQL statement"

### Sections & Components
1. SQL Header
   - SQL ID (copy button)
   - Normalized SQL in code block with wrap toggle
2. KPI cards
   - Exec Count, Error Count, Total Time, Avg Time, Max Time, Running, Concurrent Max
3. Samples table (last N)
   - Columns: Time, Duration, OK/Err, Err Msg (truncate), Rows, URI
4. Actions
   - Reset this SQL (optional parity enhancement)
   - Back to Dashboard
