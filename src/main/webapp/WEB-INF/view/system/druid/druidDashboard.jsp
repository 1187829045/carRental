<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <title>数据源监控</title>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <link rel="stylesheet" href="${yeqifu}/static/layui/css/layui.css" media="all"/>
    <style>
        body{margin:0;background:#f5f7fb;color:#1f2937;font:14px/1.6 Arial,"PingFang SC","Microsoft YaHei",sans-serif;}
        .monitor-page{padding:20px;}
        .monitor-toolbar{display:flex;justify-content:space-between;align-items:center;gap:12px;flex-wrap:wrap;margin-bottom:16px;padding:14px 18px;background:#fff;border:1px solid #e5e7eb;border-radius:12px;box-shadow:0 8px 24px rgba(15,23,42,.04);}
        .monitor-title{font-size:22px;font-weight:700;}
        .monitor-sub{color:#6b7280;font-size:13px;}
        .monitor-actions{display:flex;align-items:center;gap:10px;flex-wrap:wrap;}
        .monitor-tabs{display:flex;gap:10px;flex-wrap:wrap;margin-bottom:16px;}
        .monitor-tab{padding:8px 16px;border-radius:999px;background:#e5e7eb;color:#374151;cursor:pointer;}
        .monitor-tab.active{background:#1677ff;color:#fff;}
        .panel{display:none;background:#fff;border:1px solid #e5e7eb;border-radius:12px;padding:18px;box-shadow:0 8px 24px rgba(15,23,42,.04);}
        .panel.active{display:block;}
        .summary-grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:12px;margin-bottom:16px;}
        .summary-card{padding:16px;border-radius:12px;background:linear-gradient(135deg,#eff6ff,#fff);border:1px solid #dbeafe;}
        .summary-card__label{color:#6b7280;font-size:12px;}
        .summary-card__value{margin-top:8px;font-size:26px;font-weight:700;color:#111827;}
        .monitor-table-wrap{overflow:auto;}
        table.monitor-table{width:100%;border-collapse:collapse;min-width:920px;}
        .monitor-table th,.monitor-table td{padding:10px 12px;border-bottom:1px solid #f1f5f9;text-align:left;vertical-align:top;}
        .monitor-table th{background:#f8fafc;color:#475569;font-weight:600;}
        .monitor-filter{display:flex;gap:10px;flex-wrap:wrap;margin-bottom:16px;}
        .monitor-filter input,.monitor-filter select{height:36px;padding:0 12px;border:1px solid #d1d5db;border-radius:8px;}
        .sql-text{max-width:520px;white-space:normal;word-break:break-all;}
        .hint{margin-top:10px;color:#6b7280;font-size:12px;}
        .banner{display:none;margin-bottom:16px;padding:12px 14px;border-radius:10px;background:#eff6ff;border:1px solid #bfdbfe;color:#1d4ed8;}
        .detail-box{margin-top:16px;border:1px solid #e5e7eb;border-radius:12px;padding:16px;background:#fafafa;}
        .detail-box pre{white-space:pre-wrap;word-break:break-all;background:#111827;color:#f9fafb;padding:14px;border-radius:10px;}
        .danger{color:#dc2626;}
        @media (max-width: 1024px){
            .monitor-page{padding:12px;}
        }
    </style>
</head>
<body>
<div class="monitor-page">
    <div id="banner" class="banner"></div>
    <div class="monitor-toolbar">
        <div>
            <div class="monitor-title">数据源监控</div>
            <div class="monitor-sub">对齐 Java Druid 的数据源、SQL、URI 与会话监控</div>
        </div>
        <div class="monitor-actions">
            <label><input type="checkbox" id="autoRefresh" checked> 自动刷新</label>
            <select id="refreshInterval">
                <option value="2">2秒</option>
                <option value="5" selected>5秒</option>
                <option value="10">10秒</option>
                <option value="30">30秒</option>
            </select>
            <span id="lastUpdated" class="monitor-sub">尚未刷新</span>
            <button class="layui-btn layui-btn-primary layui-btn-sm" onclick="resetData('all')">重置全部</button>
            <button class="layui-btn layui-btn-primary layui-btn-sm" onclick="resetData('sql')">重置SQL</button>
            <button class="layui-btn layui-btn-primary layui-btn-sm" onclick="resetData('web')">重置Web</button>
            <button class="layui-btn layui-btn-primary layui-btn-sm" onclick="resetData('datasource')">重置数据源</button>
            <form method="post" action="${yeqifu}/druid/logout.action" style="margin:0;">
                <button type="submit" class="layui-btn layui-btn-danger layui-btn-sm">退出</button>
            </form>
        </div>
    </div>

    <div class="monitor-tabs">
        <div class="monitor-tab active" data-panel="datasourcePanel">数据源</div>
        <div class="monitor-tab" data-panel="sqlPanel">SQL监控</div>
        <div class="monitor-tab" data-panel="webPanel">URI监控</div>
        <div class="monitor-tab" data-panel="sessionPanel">会话监控</div>
    </div>

    <div id="datasourcePanel" class="panel active">
        <div id="datasourceSummary" class="summary-grid"></div>
        <div class="monitor-table-wrap">
            <table class="monitor-table">
                <thead>
                <tr>
                    <th>名称</th>
                    <th>URL</th>
                    <th>用户</th>
                    <th>最大连接</th>
                    <th>当前连接</th>
                    <th>活跃</th>
                    <th>空闲</th>
                    <th>活跃峰值</th>
                    <th>等待次数</th>
                    <th>等待时长(ms)</th>
                    <th>建连次数</th>
                    <th>关闭次数</th>
                    <th>建连错误</th>
                </tr>
                </thead>
                <tbody id="datasourceBody"></tbody>
            </table>
        </div>
    </div>

    <div id="sqlPanel" class="panel">
        <div class="monitor-filter">
            <input id="sqlQuery" placeholder="搜索 SQL 片段"/>
            <select id="sqlSort">
                <option value="totalTimeMs">按总耗时</option>
                <option value="executeCount">按执行次数</option>
                <option value="maxTimeMs">按最大耗时</option>
                <option value="errorCount">按错误数</option>
            </select>
            <button class="layui-btn layui-btn-sm" onclick="loadSQL()">查询</button>
        </div>
        <div class="monitor-table-wrap">
            <table class="monitor-table">
                <thead>
                <tr>
                    <th>SQL ID</th>
                    <th>SQL</th>
                    <th>执行次数</th>
                    <th>错误数</th>
                    <th>总耗时(ms)</th>
                    <th>平均耗时(ms)</th>
                    <th>最大耗时(ms)</th>
                    <th>运行中</th>
                    <th>并发峰值</th>
                    <th>首次出现</th>
                    <th>最后出现</th>
                    <th>详情</th>
                </tr>
                </thead>
                <tbody id="sqlBody"></tbody>
            </table>
        </div>
        <div id="sqlDetail" class="detail-box" style="display:none;"></div>
    </div>

    <div id="webPanel" class="panel">
        <div class="monitor-filter">
            <input id="uriQuery" placeholder="搜索 URI"/>
            <button class="layui-btn layui-btn-sm" onclick="loadWebURI()">查询</button>
        </div>
        <div class="monitor-table-wrap">
            <table class="monitor-table">
                <thead>
                <tr>
                    <th>URI</th>
                    <th>请求数</th>
                    <th>错误数</th>
                    <th>总耗时(ms)</th>
                    <th>平均耗时(ms)</th>
                    <th>最大耗时(ms)</th>
                    <th>运行中</th>
                    <th>并发峰值</th>
                    <th>JDBC次数</th>
                    <th>JDBC错误</th>
                    <th>JDBC总耗时(ms)</th>
                </tr>
                </thead>
                <tbody id="webBody"></tbody>
            </table>
        </div>
    </div>

    <div id="sessionPanel" class="panel">
        <div class="monitor-table-wrap">
            <table class="monitor-table">
                <thead>
                <tr>
                    <th>Session Key</th>
                    <th>用户</th>
                    <th>创建时间</th>
                    <th>最后访问</th>
                    <th>请求数</th>
                    <th>客户端IP</th>
                </tr>
                </thead>
                <tbody id="sessionBody"></tbody>
            </table>
        </div>
    </div>
</div>

<script>
const base = '${yeqifu}/druid';
let refreshTimer = null;

document.querySelectorAll('.monitor-tab').forEach(function(tab){
    tab.addEventListener('click', function(){
        document.querySelectorAll('.monitor-tab').forEach(function(node){ node.classList.remove('active'); });
        document.querySelectorAll('.panel').forEach(function(node){ node.classList.remove('active'); });
        tab.classList.add('active');
        document.getElementById(tab.dataset.panel).classList.add('active');
    });
});

document.getElementById('autoRefresh').addEventListener('change', setupAutoRefresh);
document.getElementById('refreshInterval').addEventListener('change', setupAutoRefresh);

function showBanner(msg, isError) {
    var banner = document.getElementById('banner');
    banner.style.display = 'block';
    banner.textContent = msg;
    banner.style.background = isError ? '#fef2f2' : '#eff6ff';
    banner.style.borderColor = isError ? '#fecaca' : '#bfdbfe';
    banner.style.color = isError ? '#b91c1c' : '#1d4ed8';
}

function hideBanner() {
    document.getElementById('banner').style.display = 'none';
}

function safeFetch(url, options) {
    return fetch(url, options).then(function(res){
        if (!res.ok) {
            throw new Error('请求失败: ' + res.status);
        }
        return res.json();
    });
}

function fillEmpty(tbody, text, colspan) {
    tbody.innerHTML = '<tr><td colspan="' + colspan + '">' + text + '</td></tr>';
}

function loadDatasource() {
    return safeFetch(base + '/datasource.json').then(function(resp){
        if (resp.code !== 0) throw new Error(resp.msg || '加载失败');
        var data = resp.data || {};
        document.getElementById('datasourceSummary').innerHTML = [
            card('当前连接', data.openConnections || 0),
            card('活跃连接', data.inUse || 0),
            card('空闲连接', data.idle || 0),
            card('等待次数', data.waitCount || 0),
            card('等待时长(ms)', data.waitDurationMs || 0),
            card('建连错误', data.connectErrorCount || 0)
        ].join('');
        document.getElementById('datasourceBody').innerHTML =
            '<tr>' +
            td(data.name) + td(data.jdbcUrlMasked) + td(data.userMasked) + td(data.maxOpenConnections) +
            td(data.openConnections) + td(data.inUse) + td(data.idle) + td((data.activePeak || 0) + '<div class="hint">' + (data.activePeakTime || '') + '</div>') +
            td(data.waitCount) + td(data.waitDurationMs) + td(data.connectCount) + td(data.closeCount) + td(warn(data.connectErrorCount)) +
            '</tr>';
    });
}

function loadSQL() {
    var q = encodeURIComponent(document.getElementById('sqlQuery').value || '');
    var sortBy = encodeURIComponent(document.getElementById('sqlSort').value || 'totalTimeMs');
    return safeFetch(base + '/sql.json?page=1&pageSize=20&q=' + q + '&sortBy=' + sortBy).then(function(resp){
        if (resp.code !== 0) throw new Error(resp.msg || '加载失败');
        var rows = (resp.data && resp.data.items) || [];
        var body = document.getElementById('sqlBody');
        if (!rows.length) {
            fillEmpty(body, '暂无 SQL 监控数据', 12);
            return;
        }
        body.innerHTML = rows.map(function(item){
            return '<tr>' +
                td(item.id) +
                td('<div class="sql-text">' + escapeHtml(item.sql) + '</div>') +
                td(item.executeCount) +
                td(warn(item.errorCount)) +
                td(item.totalTimeMs) +
                td(item.avgTimeMs) +
                td(item.maxTimeMs) +
                td(item.runningCount) +
                td(item.concurrentMax) +
                td(item.firstSeen) +
                td(item.lastSeen) +
                td('<button class="layui-btn layui-btn-xs" onclick="loadSQLDetail(\'' + item.id + '\')">查看</button>') +
                '</tr>';
        }).join('');
    });
}

function loadSQLDetail(id) {
    safeFetch(base + '/sql-' + encodeURIComponent(id) + '.json').then(function(resp){
        if (resp.code !== 0) throw new Error(resp.msg || '加载详情失败');
        var detail = resp.data || {};
        var item = detail.item || {};
        var samples = detail.samples || [];
        var html = '<h3>SQL 详情</h3>' +
            '<pre>' + escapeHtml(item.sql || '') + '</pre>' +
            '<div class="hint">SQL ID: ' + escapeHtml(item.id || '') + '</div>' +
            '<div class="summary-grid">' +
            card('执行次数', item.executeCount || 0) +
            card('错误数', item.errorCount || 0) +
            card('总耗时(ms)', item.totalTimeMs || 0) +
            card('平均耗时(ms)', item.avgTimeMs || 0) +
            card('最大耗时(ms)', item.maxTimeMs || 0) +
            card('并发峰值', item.concurrentMax || 0) +
            '</div>' +
            '<div class="monitor-table-wrap"><table class="monitor-table"><thead><tr><th>时间</th><th>耗时(ms)</th><th>结果</th><th>影响行数</th><th>URI</th><th>错误信息</th></tr></thead><tbody>' +
            (samples.length ? samples.map(function(sample){
                return '<tr>' + td(sample.ts) + td(sample.durationMs) + td(sample.ok ? '成功' : '<span class="danger">失败</span>') + td(sample.rowsAffected || 0) + td(sample.uri || '') + td(escapeHtml(sample.errMsg || '')) + '</tr>';
            }).join('') : '<tr><td colspan="6">暂无样本</td></tr>') +
            '</tbody></table></div>';
        var box = document.getElementById('sqlDetail');
        box.style.display = 'block';
        box.innerHTML = html;
    }).catch(function(err){
        showBanner(err.message, true);
    });
}

function loadWebURI() {
    var q = encodeURIComponent(document.getElementById('uriQuery').value || '');
    return safeFetch(base + '/weburi.json?page=1&pageSize=20&q=' + q).then(function(resp){
        if (resp.code !== 0) throw new Error(resp.msg || '加载失败');
        var rows = (resp.data && resp.data.items) || [];
        var body = document.getElementById('webBody');
        if (!rows.length) {
            fillEmpty(body, '暂无 URI 监控数据', 11);
            return;
        }
        body.innerHTML = rows.map(function(item){
            return '<tr>' +
                td(item.uri) + td(item.requestCount) + td(warn(item.errorCount)) + td(item.totalTimeMs) + td(item.avgTimeMs) + td(item.maxTimeMs) +
                td(item.runningCount) + td(item.concurrentMax) + td(item.jdbcExecuteCount) + td(warn(item.jdbcErrorCount)) + td(item.jdbcTotalTimeMs) +
                '</tr>';
        }).join('');
    });
}

function loadSessions() {
    return safeFetch(base + '/websession.json?page=1&pageSize=20').then(function(resp){
        if (resp.code !== 0) throw new Error(resp.msg || '加载失败');
        var rows = (resp.data && resp.data.items) || [];
        var body = document.getElementById('sessionBody');
        if (!rows.length) {
            fillEmpty(body, '暂无会话监控数据', 6);
            return;
        }
        body.innerHTML = rows.map(function(item){
            return '<tr>' +
                td(item.sessionKey) + td(item.user || '-') + td(item.createTime || '') + td(item.lastAccessTime || '') + td(item.requestCount || 0) + td(item.remoteAddr || '') +
                '</tr>';
        }).join('');
    });
}

function resetData(type) {
    if (!confirm('确认重置 ' + type + ' 监控数据吗？')) return;
    safeFetch(base + '/reset-' + type + '.json', {method: 'POST'}).then(function(resp){
        if (resp.code !== 0) throw new Error(resp.msg || '重置失败');
        showBanner(resp.msg || '重置成功', false);
        loadAll();
    }).catch(function(err){
        showBanner(err.message, true);
    });
}

function loadAll() {
    Promise.all([loadDatasource(), loadSQL(), loadWebURI(), loadSessions()]).then(function(){
        hideBanner();
        document.getElementById('lastUpdated').textContent = '最近刷新：' + new Date().toLocaleString();
    }).catch(function(err){
        showBanner(err.message, true);
    });
}

function setupAutoRefresh() {
    if (refreshTimer) {
        clearInterval(refreshTimer);
        refreshTimer = null;
    }
    if (!document.getElementById('autoRefresh').checked) {
        return;
    }
    var seconds = parseInt(document.getElementById('refreshInterval').value || '5', 10);
    refreshTimer = setInterval(loadAll, seconds * 1000);
}

function card(label, value) {
    return '<div class="summary-card"><div class="summary-card__label">' + label + '</div><div class="summary-card__value">' + value + '</div></div>';
}

function td(value) {
    return '<td>' + (value == null ? '' : value) + '</td>';
}

function warn(value) {
    value = value || 0;
    return value > 0 ? '<span class="danger">' + value + '</span>' : value;
}

function escapeHtml(value) {
    return String(value || '')
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

loadAll();
setupAutoRefresh();
</script>
</body>
</html>
