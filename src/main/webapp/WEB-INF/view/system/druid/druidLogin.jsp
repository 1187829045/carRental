<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <meta charset="utf-8">
    <title>数据源监控登录</title>
    <meta name="renderer" content="webkit">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <link rel="stylesheet" href="${yeqifu}/static/layui/css/layui.css" media="all"/>
    <link rel="stylesheet" href="${yeqifu}/static/css/public.css" media="all"/>
    <style>
        html, body { height: 100%; }
        body { background: #f5f8ff; margin: 0; }
        .monitor-mask {
            position: fixed;
            inset: 0;
            background: rgba(17, 47, 102, 0.45);
            backdrop-filter: blur(2px);
        }
        .monitor-modal {
            position: fixed;
            left: 50%;
            top: 50%;
            transform: translate(-50%, -50%);
            width: min(420px, calc(100vw - 40px));
            background: #fff;
            border: 1px solid #dbe7ff;
            border-radius: 14px;
            box-shadow: 0 18px 60px rgba(11, 92, 255, 0.18);
            overflow: hidden;
        }
        .monitor-modal__header {
            padding: 16px 18px;
            color: #fff;
            background: linear-gradient(90deg, #2f80ff 0%, #6bb6ff 100%);
        }
        .monitor-modal__title {
            font-size: 16px;
            font-weight: 700;
            letter-spacing: 0.5px;
        }
        .monitor-modal__sub {
            margin-top: 6px;
            font-size: 12px;
            opacity: 0.9;
        }
        .monitor-modal__body { padding: 18px; }
        .monitor-modal__footer {
            padding: 0 18px 18px;
        }
        .monitor-tip {
            margin-top: 10px;
            color: #e74c3c;
            font-size: 13px;
        }
    </style>
</head>
<body>
<div class="monitor-mask"></div>

<div class="monitor-modal">
    <div class="monitor-modal__header">
        <div class="monitor-modal__title">数据源监控登录</div>
        <div class="monitor-modal__sub">请输入账号密码后继续</div>
    </div>

    <div class="monitor-modal__body">
        <form class="layui-form" method="post" action="${yeqifu}/druid/login.action">
            <div class="layui-form-item">
                <label class="layui-form-label">账号</label>
                <div class="layui-input-block">
                    <input type="text" name="username" autocomplete="off" placeholder="请输入账号" class="layui-input" />
                </div>
            </div>
            <div class="layui-form-item">
                <label class="layui-form-label">密码</label>
                <div class="layui-input-block">
                    <input type="password" name="password" autocomplete="off" placeholder="请输入密码" class="layui-input" />
                </div>
            </div>
            <div class="monitor-tip">${error}</div>
            <div class="monitor-modal__footer">
                <button type="submit" class="layui-btn layui-btn-normal" style="width:100%">登录</button>
            </div>
        </form>
    </div>
</div>

</body>
</html>

