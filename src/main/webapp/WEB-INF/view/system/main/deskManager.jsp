<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <title>工作台</title>
    <meta name="renderer" content="webkit">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <meta name="apple-mobile-web-app-status-bar-style" content="black">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="format-detection" content="telephone=no">
    <link rel="stylesheet" href="${yeqifu}/static/layui/css/layui.css" media="all" />
    <link rel="stylesheet" href="${yeqifu}/static/css/public.css" media="all" />
    <style>
        body.childrenBody.dashboard-body{
            padding: 18px;
            background: linear-gradient(180deg, #f4f8ff 0%, #eef4ff 100%);
            color: #1f2d3d;
        }
        .dashboard-shell{
            max-width: 1440px;
            margin: 0 auto;
        }
        .dashboard-grid{
            display: grid;
            gap: 16px;
        }
        .dashboard-hero{
            display: grid;
            grid-template-columns: minmax(0, 1.65fr) minmax(320px, 0.95fr);
            gap: 16px;
            align-items: stretch;
        }
        .dashboard-card{
            position: relative;
            overflow: hidden;
            background: rgba(255,255,255,0.96);
            border: 1px solid rgba(47,128,255,0.12);
            border-radius: 20px;
            box-shadow: 0 12px 35px rgba(47, 128, 255, 0.08);
            animation: fadeUp .45s ease both;
        }
        .hero-panel{
            padding: 28px;
            min-height: 250px;
            background: linear-gradient(135deg, #1f6bff 0%, #55a4ff 58%, #8fc4ff 100%);
            color: #fff;
        }
        .hero-panel:before,
        .hero-panel:after{
            content: "";
            position: absolute;
            border-radius: 50%;
            background: rgba(255,255,255,0.12);
        }
        .hero-panel:before{
            width: 220px;
            height: 220px;
            right: -60px;
            top: -70px;
        }
        .hero-panel:after{
            width: 160px;
            height: 160px;
            right: 90px;
            bottom: -90px;
        }
        .hero-panel__top{
            position: relative;
            z-index: 1;
            display: flex;
            justify-content: space-between;
            gap: 16px;
            align-items: flex-start;
            flex-wrap: wrap;
        }
        .hero-badge{
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 7px 14px;
            background: rgba(255,255,255,0.16);
            border: 1px solid rgba(255,255,255,0.2);
            border-radius: 999px;
            font-size: 13px;
        }
        .hero-title{
            margin: 18px 0 10px;
            font-size: 30px;
            font-weight: 700;
            line-height: 1.3;
        }
        .hero-desc{
            max-width: 620px;
            font-size: 14px;
            line-height: 1.8;
            color: rgba(255,255,255,0.9);
        }
        .hero-actions{
            position: relative;
            z-index: 1;
            margin-top: 26px;
            display: flex;
            gap: 12px;
            flex-wrap: wrap;
        }
        .hero-action{
            display: inline-flex;
            align-items: center;
            gap: 8px;
            padding: 11px 18px;
            border-radius: 12px;
            background: #fff;
            color: #1f6bff;
            font-weight: 600;
            transition: transform .25s ease, box-shadow .25s ease;
        }
        .hero-action.secondary{
            background: rgba(255,255,255,0.16);
            color: #fff;
            border: 1px solid rgba(255,255,255,0.26);
        }
        .hero-action:hover{
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(15, 23, 42, 0.15);
        }
        .profile-panel{
            padding: 24px;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            min-height: 250px;
        }
        .profile-top{
            display: flex;
            align-items: center;
            gap: 16px;
        }
        .profile-avatar{
            width: 72px;
            height: 72px;
            border-radius: 18px;
            background: linear-gradient(135deg, #2f80ff 0%, #8fc4ff 100%);
            color: #fff;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 28px;
            font-weight: 700;
            box-shadow: 0 12px 24px rgba(47, 128, 255, 0.18);
        }
        .profile-name{
            font-size: 24px;
            font-weight: 700;
            color: #13223a;
        }
        .profile-role{
            margin-top: 6px;
            color: #5b6b85;
        }
        .profile-metrics{
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 12px;
            margin-top: 20px;
        }
        .metric-chip{
            padding: 14px 16px;
            border-radius: 16px;
            background: #f6f9ff;
            border: 1px solid #e4eeff;
        }
        .metric-chip__label{
            color: #6b7b94;
            font-size: 12px;
        }
        .metric-chip__value{
            margin-top: 6px;
            font-size: 20px;
            font-weight: 700;
            color: #13223a;
        }
        .section-title{
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
            margin-bottom: 14px;
        }
        .section-title h3{
            margin: 0;
            font-size: 18px;
            font-weight: 700;
            color: #13223a;
        }
        .section-title span{
            color: #6b7b94;
            font-size: 12px;
        }
        .stats-grid{
            display: grid;
            grid-template-columns: repeat(4, minmax(0, 1fr));
            gap: 16px;
        }
        .stats-card{
            padding: 22px 20px;
        }
        .stats-card__icon{
            width: 46px;
            height: 46px;
            border-radius: 14px;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            font-size: 22px;
            color: #fff;
            margin-bottom: 18px;
        }
        .stats-card__icon.blue{background: linear-gradient(135deg, #2f80ff 0%, #6db5ff 100%);}
        .stats-card__icon.green{background: linear-gradient(135deg, #34c38f 0%, #76dfba 100%);}
        .stats-card__icon.orange{background: linear-gradient(135deg, #ff9f43 0%, #ffc371 100%);}
        .stats-card__icon.purple{background: linear-gradient(135deg, #7a5cff 0%, #b07cff 100%);}
        .stats-card__label{
            color: #6b7b94;
            font-size: 13px;
        }
        .stats-card__value{
            margin-top: 10px;
            font-size: 30px;
            font-weight: 700;
            color: #13223a;
        }
        .stats-card__desc{
            margin-top: 6px;
            color: #8a97ab;
            font-size: 12px;
        }
        .content-grid{
            display: grid;
            grid-template-columns: 1.15fr 0.85fr;
            gap: 16px;
        }
        .module-card{
            padding: 20px;
        }
        .shortcut-grid,
        .recommend-grid{
            display: grid;
            gap: 12px;
        }
        .shortcut-grid{
            grid-template-columns: repeat(4, minmax(0, 1fr));
        }
        .recommend-grid{
            grid-template-columns: repeat(2, minmax(0, 1fr));
        }
        .shortcut-item,
        .recommend-item{
            display: block;
            padding: 16px;
            border-radius: 16px;
            background: #f7faff;
            border: 1px solid #e7efff;
            transition: transform .25s ease, box-shadow .25s ease, border-color .25s ease;
        }
        .shortcut-item:hover,
        .recommend-item:hover{
            transform: translateY(-3px);
            box-shadow: 0 12px 24px rgba(47, 128, 255, 0.08);
            border-color: #cdddff;
        }
        .shortcut-item__icon,
        .recommend-item__icon{
            width: 42px;
            height: 42px;
            border-radius: 14px;
            display: inline-flex;
            align-items: center;
            justify-content: center;
            background: #eaf2ff;
            color: #2f80ff;
            font-size: 20px;
        }
        .shortcut-item__title,
        .recommend-item__title{
            margin-top: 14px;
            font-size: 16px;
            font-weight: 700;
            color: #13223a;
        }
        .shortcut-item__desc,
        .recommend-item__desc{
            margin-top: 6px;
            color: #6b7b94;
            line-height: 1.7;
            min-height: 42px;
        }
        .list-card{
            padding: 20px;
        }
        .list-table{
            width: 100%;
            border-collapse: collapse;
        }
        .list-table tr{
            transition: background-color .2s ease;
        }
        .list-table tr:hover{
            background: #f7faff;
        }
        .list-table td{
            padding: 14px 8px;
            border-bottom: 1px solid #eef3ff;
            color: #1f2d3d;
        }
        .list-table td:first-child{
            padding-left: 0;
        }
        .list-table td:last-child{
            padding-right: 0;
        }
        .list-title-link{
            display: inline-block;
            max-width: 100%;
            color: #13223a;
            font-weight: 500;
        }
        .list-date{
            color: #8a97ab;
            font-size: 12px;
            white-space: nowrap;
            text-align: right;
        }
        .empty-block{
            padding: 26px 10px;
            text-align: center;
            color: #8a97ab;
        }
        .tips-list{
            margin: 0;
            padding-left: 18px;
            color: #62738d;
        }
        .tips-list li{
            margin-bottom: 10px;
            line-height: 1.8;
        }
        @keyframes fadeUp {
            from {
                opacity: 0;
                transform: translateY(14px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
        @media screen and (max-width: 1200px){
            .dashboard-hero,
            .content-grid{
                grid-template-columns: 1fr;
            }
            .shortcut-grid{
                grid-template-columns: repeat(2, minmax(0, 1fr));
            }
        }
        @media screen and (max-width: 768px){
            body.childrenBody.dashboard-body{
                padding: 12px;
            }
            .hero-panel,
            .profile-panel,
            .module-card,
            .list-card,
            .stats-card{
                padding: 18px;
            }
            .hero-title{
                font-size: 24px;
            }
            .stats-grid,
            .recommend-grid,
            .profile-metrics{
                grid-template-columns: 1fr 1fr;
            }
        }
        @media screen and (max-width: 520px){
            .stats-grid,
            .shortcut-grid,
            .recommend-grid,
            .profile-metrics{
                grid-template-columns: 1fr;
            }
            .hero-actions{
                flex-direction: column;
            }
            .hero-action{
                justify-content: center;
            }
        }
    </style>
</head>
<body class="childrenBody dashboard-body">
<div class="dashboard-shell dashboard-grid">
    <div class="dashboard-hero">
        <div class="dashboard-card hero-panel">
            <div class="hero-panel__top">
                <div class="hero-badge">
                    <i class="layui-icon">&#xe68e;</i>
                    <span>智能工作台</span>
                </div>
                <div class="hero-badge">
                    <i class="layui-icon">&#xe637;</i>
                    <span id="nowTimeText">正在同步时间...</span>
                </div>
            </div>
            <div class="hero-title">欢迎回来，${user.realname}</div>
            <div class="hero-desc" id="heroDesc">
                工作台将优先展示业务概览、常用快捷入口和与你角色相关的重点功能，帮助你在更少操作下快速进入工作状态。
            </div>
            <div class="hero-actions">
                <a href="javascript:;" class="hero-action js-open-tab" data-url="${yeqifu}/bus/toRentCarManager.action" data-title="汽车出租" data-icon="&#xe609;">
                    <i class="layui-icon">&#xe609;</i><span>发起出租</span>
                </a>
                <a href="javascript:;" class="hero-action secondary js-open-tab" data-url="${yeqifu}/stat/toCompanyYearGradeStat.action" data-title="公司年度销售额" data-icon="&#xe62c;">
                    <i class="layui-icon">&#xe62c;</i><span>查看经营趋势</span>
                </a>
            </div>
        </div>
        <div class="dashboard-card profile-panel">
            <div>
                <div class="profile-top">
                    <div class="profile-avatar" id="profileAvatar">租</div>
                    <div>
                        <div class="profile-name">${user.realname}</div>
                        <div class="profile-role" id="profileRole">正在识别当前角色...</div>
                    </div>
                </div>
                <div class="profile-metrics">
                    <div class="metric-chip">
                        <div class="metric-chip__label">当前身份</div>
                        <div class="metric-chip__value" id="metricRole">管理员</div>
                    </div>
                    <div class="metric-chip">
                        <div class="metric-chip__label">推荐模块</div>
                        <div class="metric-chip__value" id="metricRecommend">4项</div>
                    </div>
                    <div class="metric-chip">
                        <div class="metric-chip__label">快捷入口</div>
                        <div class="metric-chip__value" id="metricShortcut">6项</div>
                    </div>
                    <div class="metric-chip">
                        <div class="metric-chip__label">界面状态</div>
                        <div class="metric-chip__value">流畅</div>
                    </div>
                </div>
            </div>
            <div class="tips-list">
                <li>顶部菜单与左侧导航保持现有系统结构，避免影响你的业务习惯。</li>
                <li>工作台中的快捷操作支持直接打开新标签页，减少多级菜单查找。</li>
            </div>
        </div>
    </div>

    <div class="section-title">
        <h3>业务概览</h3>
        <span>关键数据一屏可见，帮助快速掌握系统运行情况</span>
    </div>
    <div class="stats-grid">
        <div class="dashboard-card stats-card">
            <div class="stats-card__icon blue"><i class="layui-icon">&#xe770;</i></div>
            <div class="stats-card__label">系统用户数</div>
            <div class="stats-card__value" id="userCount">--</div>
            <div class="stats-card__desc">包含系统管理与业务操作用户</div>
        </div>
        <div class="dashboard-card stats-card">
            <div class="stats-card__icon green"><i class="layui-icon">&#xe657;</i></div>
            <div class="stats-card__label">车辆档案数</div>
            <div class="stats-card__value" id="carCount">--</div>
            <div class="stats-card__desc">用于反映当前车辆资源规模</div>
        </div>
        <div class="dashboard-card stats-card">
            <div class="stats-card__icon orange"><i class="layui-icon">&#xe770;</i></div>
            <div class="stats-card__label">客户数量</div>
            <div class="stats-card__value" id="customerCount">--</div>
            <div class="stats-card__desc">便于快速了解客户沉淀情况</div>
        </div>
        <div class="dashboard-card stats-card">
            <div class="stats-card__icon purple"><i class="layui-icon">&#xe63c;</i></div>
            <div class="stats-card__label">出租订单数</div>
            <div class="stats-card__value" id="rentCount">--</div>
            <div class="stats-card__desc">当前累计出租订单统计</div>
        </div>
    </div>

    <div class="content-grid">
        <div class="dashboard-grid">
            <div class="dashboard-card module-card">
                <div class="section-title">
                    <h3>快捷操作</h3>
                    <span>高频入口直达，减少跳转层级</span>
                </div>
                <div class="shortcut-grid" id="shortcutGrid"></div>
            </div>
            <div class="dashboard-card list-card">
                <div class="section-title">
                    <h3>系统公告</h3>
                    <span>双击或点击标题查看详情</span>
                </div>
                <table class="list-table">
                    <tbody class="hot_news"></tbody>
                </table>
            </div>
            <div class="dashboard-card list-card">
                <div class="section-title">
                    <h3>最新留言</h3>
                    <span>跟进近期反馈与沟通内容</span>
                </div>
                <table class="list-table">
                    <tbody class="hot_message"></tbody>
                </table>
            </div>
        </div>

        <div class="dashboard-grid">
            <div class="dashboard-card module-card">
                <div class="section-title">
                    <h3>个性化推荐</h3>
                    <span>根据当前角色推荐更适合的工作模块</span>
                </div>
                <div class="recommend-grid" id="recommendGrid"></div>
            </div>
            <div class="dashboard-card module-card">
                <div class="section-title">
                    <h3>使用建议</h3>
                    <span>帮助更快完成日常管理操作</span>
                </div>
                <ul class="tips-list" id="tipsList">
                    <li>建议优先通过快捷操作进入高频页面，减少重复展开菜单。</li>
                    <li>如需查看经营表现，可直接进入“统计分析”查看年度趋势。</li>
                    <li>重要通知与业务反馈会优先显示在右侧动态模块中。</li>
                </ul>
            </div>
        </div>
    </div>
</div>

<div id="desk_viewNewsDiv" style="padding: 10px;display: none;">
    <h2 id="view_title" align="center"></h2>
    <hr>
    <div style="text-align: right;">
        发布人:<span id="view_opername"></span>
        <span style="display: inline-block;width: 20px" ></span>
        发布时间:<span id="view_createtime"></span>
    </div>
    <hr>
    <div id="view_content"></div>
</div>

<div id="desk_viewMessageDiv" style="padding: 10px;display: none;">
    <h2 id="view_title_message" align="center"></h2>
    <hr>
    <div style="text-align: right;">
        发布人:<span id="view_opername_message"></span>
        <span style="display: inline-block;width: 20px" ></span>
        发布时间:<span id="view_createtime_message"></span>
    </div>
    <hr>
    <div id="view_content_message"></div>
</div>

<script type="text/javascript" src="${yeqifu}/static/layui/layui.js"></script>
<script>
    var $;
    var layer;
    var userName = '${user.realname}';
    var userType = Number('${user.type}');
    var roleLabel = userType === 1 ? '超级管理员' : '业务操作员';
    var shortcutConfig = [
        {title: '车辆管理', desc: '快速维护车辆档案与状态', icon: '&#xe657;', url: '${yeqifu}/bus/toCarManager.action'},
        {title: '客户管理', desc: '查看和维护客户基础信息', icon: '&#xe770;', url: '${yeqifu}/bus/toCustomerManager.action'},
        {title: '汽车出租', desc: '进入出租流程处理新订单', icon: '&#xe609;', url: '${yeqifu}/bus/toRentCarManager.action'},
        {title: '出租管理', desc: '跟踪历史出租记录与状态', icon: '&#xe63c;', url: '${yeqifu}/bus/toRentManager.action'},
        {title: '检查单管理', desc: '处理还车检查与赔付记录', icon: '&#xe60a;', url: '${yeqifu}/bus/toCheckManager.action'},
        {title: '统计分析', desc: '查看经营趋势与业务分析', icon: '&#xe62c;', url: '${yeqifu}/stat/toCompanyYearGradeStat.action'}
    ];
    var adminRecommend = [
        {title: '用户管理', desc: '优先处理账号、角色与权限分配。', icon: '&#xe770;', url: '${yeqifu}/sys/toUserManager.action'},
        {title: '角色管理', desc: '检查角色授权，保证菜单与权限一致。', icon: '&#xe66f;', url: '${yeqifu}/sys/toRoleManager.action'},
        {title: '系统公告', desc: '更新公告内容，提升通知触达效率。', icon: '&#xe645;', url: '${yeqifu}/sys/toNewsManager.action'},
        {title: '日志管理', desc: '关注登录日志与审计记录变化。', icon: '&#xe60e;', url: '${yeqifu}/sys/toLogInfoManager.action'}
    ];
    var operatorRecommend = [
        {title: '汽车出租', desc: '优先处理当天出租业务，缩短业务链路。', icon: '&#xe609;', url: '${yeqifu}/bus/toRentCarManager.action'},
        {title: '出租管理', desc: '及时查看待归还订单与历史记录。', icon: '&#xe63c;', url: '${yeqifu}/bus/toRentManager.action'},
        {title: '检查单录入', desc: '处理还车检查与损耗信息。', icon: '&#xe60a;', url: '${yeqifu}/bus/toCheckCarManager.action'},
        {title: '客户信息', desc: '跟进高频客户资料，提升接待效率。', icon: '&#xe770;', url: '${yeqifu}/bus/toCustomerManager.action'}
    ];

    function dateFilter(date){
        return date < 10 ? '0' + date : date;
    }

    function updateClock(){
        var dateObj = new Date();
        var year = dateObj.getFullYear();
        var month = dateObj.getMonth() + 1;
        var date = dateObj.getDate();
        var day = dateObj.getDay();
        var weeks = ['星期日','星期一','星期二','星期三','星期四','星期五','星期六'];
        var week = weeks[day];
        var hour = dateObj.getHours();
        var minute = dateObj.getMinutes();
        var second = dateObj.getSeconds();
        var period = hour >= 18 ? '晚上' : hour >= 12 ? '下午' : '上午';
        var nowText = dateFilter(year) + '年' + dateFilter(month) + '月' + dateFilter(date) + '日 ' + dateFilter(hour) + ':' + dateFilter(minute) + ':' + dateFilter(second);
        $('#nowTimeText').text(nowText + ' ' + week);
        $('#heroDesc').text('现在是' + period + '，' + userName + '，欢迎回到汽车租赁系统。工作台已为你整合常用入口、关键数据和动态信息，帮助你更快进入工作状态。');
    }

    function openTab(title, url, icon) {
        var link = $('<a data-url="' + url + '"><i class="layui-icon" data-icon="' + icon + '">' + icon + '</i><cite>' + title + '</cite></a>');
        parent.addTab(link);
    }

    function escapeHtml(value) {
        return String(value || '')
            .replace(/&/g, '&amp;')
            .replace(/</g, '&lt;')
            .replace(/>/g, '&gt;')
            .replace(/"/g, '&quot;')
            .replace(/'/g, '&#39;');
    }

    function formatCount(value) {
        if (value === undefined || value === null || value === '') {
            return '--';
        }
        return value;
    }

    function renderShortcutCards() {
        $('#metricShortcut').text(shortcutConfig.length + '项');
        $('#shortcutGrid').html($.map(shortcutConfig, function(item) {
            return '<a href="javascript:;" class="shortcut-item js-open-tab" data-url="' + item.url + '" data-title="' + item.title + '" data-icon="' + item.icon + '">' +
                '<span class="shortcut-item__icon"><i class="layui-icon">' + item.icon + '</i></span>' +
                '<div class="shortcut-item__title">' + item.title + '</div>' +
                '<div class="shortcut-item__desc">' + item.desc + '</div>' +
            '</a>';
        }).join(''));
    }

    function renderRecommendCards() {
        var list = userType === 1 ? adminRecommend : operatorRecommend;
        $('#metricRecommend').text(list.length + '项');
        $('#recommendGrid').html($.map(list, function(item) {
            return '<a href="javascript:;" class="recommend-item js-open-tab" data-url="' + item.url + '" data-title="' + item.title + '" data-icon="' + item.icon + '">' +
                '<span class="recommend-item__icon"><i class="layui-icon">' + item.icon + '</i></span>' +
                '<div class="recommend-item__title">' + item.title + '</div>' +
                '<div class="recommend-item__desc">' + item.desc + '</div>' +
            '</a>';
        }).join(''));
    }

    function renderNews(list) {
        if (!list.length) {
            $('.hot_news').html('<tr><td class="empty-block" colspan="2">暂无公告内容</td></tr>');
            return;
        }
        $('.hot_news').html($.map(list, function(item) {
            var title = escapeHtml(item.title || '');
            var date = item.createtime ? item.createtime.substring(0, 10) : '';
            return '<tr ondblclick="viewNews(' + item.id + ')">' +
                '<td><a href="javascript:;" class="list-title-link" onclick="viewNews(' + item.id + ')">' + title + '</a></td>' +
                '<td class="list-date">' + date + '</td>' +
            '</tr>';
        }).join(''));
    }

    function renderMessages(list) {
        if (!list.length) {
            $('.hot_message').html('<tr><td class="empty-block" colspan="2">暂无留言信息</td></tr>');
            return;
        }
        $('.hot_message').html($.map(list, function(item) {
            var title = escapeHtml(item.title || '');
            var date = item.createtime ? item.createtime.substring(0, 10) : '';
            return '<tr ondblclick="viewMessage(' + item.id + ')">' +
                '<td><a href="javascript:;" class="list-title-link" onclick="viewMessage(' + item.id + ')">' + title + '</a></td>' +
                '<td class="list-date">' + date + '</td>' +
            '</tr>';
        }).join(''));
    }

    function setCount(id, response) {
        var count = response && response.count !== undefined ? response.count : '--';
        $(id).text(formatCount(count));
    }

    function loadOverview() {
        $.when(
            $.get('${yeqifu}/user/loadAllUser.action?page=1&limit=1'),
            $.get('${yeqifu}/car/loadAllCar.action?page=1&limit=1'),
            $.get('${yeqifu}/customer/loadAllCustomer.action?page=1&limit=1'),
            $.get('${yeqifu}/rent/loadAllRent.action?page=1&limit=1'),
            $.get('${yeqifu}/news/loadAllNews.action?page=1&limit=6'),
            $.get('${yeqifu}/message/loadAllMessage.action?page=1&limit=6')
        ).done(function(userResp, carResp, customerResp, rentResp, newsResp, messageResp) {
            setCount('#userCount', userResp[0]);
            setCount('#carCount', carResp[0]);
            setCount('#customerCount', customerResp[0]);
            setCount('#rentCount', rentResp[0]);
            renderNews((newsResp[0].data || []).slice(0, 5));
            renderMessages((messageResp[0].data || []).slice(0, 5));
        }).fail(function() {
            $('#userCount,#carCount,#customerCount,#rentCount').text('--');
            $('.hot_news').html('<tr><td class="empty-block" colspan="2">公告加载失败，请稍后刷新重试</td></tr>');
            $('.hot_message').html('<tr><td class="empty-block" colspan="2">留言加载失败，请稍后刷新重试</td></tr>');
        });
    }

    layui.use(['layer', 'jquery'], function(){
        layer = parent.layer === undefined ? layui.layer : top.layer;
        $ = layui.jquery;

        $('#profileAvatar').text(userName ? userName.charAt(0) : '租');
        $('#profileRole').text(roleLabel + ' · 已登录工作台');
        $('#metricRole').text(roleLabel);

        updateClock();
        setInterval(updateClock, 1000);
        renderShortcutCards();
        renderRecommendCards();
        loadOverview();

        $('body').on('click', '.js-open-tab', function() {
            openTab($(this).data('title'), $(this).data('url'), $(this).data('icon'));
        });
    });

    function viewNews(id){
        $.get('${yeqifu}/news/loadNewsById.action', {id:id}, function(news){
            layer.open({
                type:1,
                title:'查看公告',
                content:$('#desk_viewNewsDiv'),
                area:['800px','550px'],
                success:function(){
                    $('#view_title').html(news.title || '');
                    $('#view_opername').html(news.opername || '');
                    $('#view_createtime').html(news.createtime || '');
                    $('#view_content').html(news.content || '');
                }
            });
        });
    }

    function viewMessage(id){
        $.get('${yeqifu}/message/loadMessageById.action', {id:id}, function(message){
            layer.open({
                type:1,
                title:'查看留言',
                content:$('#desk_viewMessageDiv'),
                area:['800px','550px'],
                success:function(){
                    $('#view_title_message').html(message.title || '');
                    $('#view_opername_message').html(message.opername || '');
                    $('#view_createtime_message').html(message.createtime || '');
                    $('#view_content_message').html(message.content || '');
                }
            });
        });
    }
</script>
</body>
</html>
