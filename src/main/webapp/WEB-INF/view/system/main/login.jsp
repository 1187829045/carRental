<%--
  Created by IntelliJ IDEA.
  User: YQF
  Date: 2019/9/28
  Time: 16:35
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<!DOCTYPE html>
<html class="loginHtml" lang="zh-CN">
<head>
    <meta charset="utf-8">
    <title>登录--汽车出租系统</title>
    <meta name="renderer" content="webkit">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <meta name="apple-mobile-web-app-status-bar-style" content="black">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="format-detection" content="telephone=no">
    <link rel="icon" href="${yeqifu}/static/favicon.ico">
    <link rel="stylesheet" href="${yeqifu}/static/layui/css/layui.css" media="all" />
    <link rel="stylesheet" href="${yeqifu}/static/css/public.css" media="all" />
</head>
<body class="loginBody login-modern">
<div class="login-shell">
    <div class="login-brand">
        <div class="login-brand__name">汽车租赁系统</div>
        <div class="login-brand__desc">安全登录 · 统一管理 · 高效运营</div>
    </div>
    <div class="login-card">
        <div class="login-card__head">
            <div class="login-avatar"><img src="${yeqifu}/static/images/face.jpg" alt="avatar"></div>
            <div class="login-title">欢迎回来</div>
            <div class="login-sub">请使用账号密码完成登录</div>
        </div>

        <form class="layui-form" id="loginFrm" method="post" action="${yeqifu}/login/login.action" autocomplete="on">
            <div class="login-error" id="loginError">${error}</div>

            <div class="layui-form-item login-field">
                <div class="login-field__label">用户名</div>
                <div class="login-field__control">
                    <i class="layui-icon login-field__icon">&#xe66f;</i>
                    <input type="text" name="loginname" id="loginname" placeholder="请输入用户名" class="layui-input" lay-verify="required" autocomplete="username" />
                </div>
            </div>

            <div class="layui-form-item login-field">
                <div class="login-field__label">密码</div>
                <div class="login-field__control">
                    <i class="layui-icon login-field__icon">&#xe673;</i>
                    <input type="password" name="pwd" id="pwd" placeholder="请输入密码" class="layui-input" lay-verify="required" autocomplete="current-password" />
                    <button type="button" class="login-field__toggle" id="togglePwd" aria-label="toggle password">显示</button>
                </div>
            </div>

            <div class="layui-form-item login-field" id="imgCode">
                <div class="login-field__label">验证码</div>
                <div class="login-field__control">
                    <i class="layui-icon login-field__icon">&#xe64c;</i>
                    <input type="text" name="code" id="code" placeholder="请输入验证码" class="layui-input" lay-verify="required" autocomplete="off" inputmode="numeric" />
                    <img id="captchaImg" src="${yeqifu}/login/getCode.action" alt="captcha" />
                </div>
            </div>

            <div class="layui-form-item login-actions">
                <label class="login-remember"><input type="checkbox" id="rememberMe" lay-skin="primary" title="" />记住我</label>
                <a class="login-forgot" href="javascript:;" id="forgotPwd">忘记密码？</a>
            </div>

            <div class="layui-form-item">
                <button class="layui-btn layui-block login-btn" id="loginBtn" lay-filter="login" lay-submit>登录</button>
            </div>

            <div class="login-third">
                <div class="login-third__line"><span>其他方式</span></div>
                <div class="login-third__icons">
                    <a href="javascript:;" class="login-third__icon" data-provider="wechat">微信</a>
                    <a href="javascript:;" class="login-third__icon" data-provider="qq">QQ</a>
                </div>
            </div>
        </form>
    </div>
</div>
<script type="text/javascript" src="${yeqifu}/static/layui/layui.js"></script>
<script type="text/javascript" src="${yeqifu}/static/js/cache.js"></script>
<script type="text/javascript">
    layui.use(['form','layer','jquery'],function(){
        var form = layui.form,
            layer = parent.layer === undefined ? layui.layer : top.layer;
        $ = layui.jquery;

        /*$(".loginBody .seraph").click(function(){
            layer.msg("这只是做个样式，至于功能，你见过哪个后台能这样登录的？还是老老实实的找管理员去注册吧",{
                time:5000
            });
        })*/

        function refreshCaptcha(){
            $('#captchaImg').attr('src','${yeqifu}/login/getCode.action?ts=' + Date.now());
        }

        $('#captchaImg').on('click', refreshCaptcha);

        var errText = $.trim($('#loginError').text() || '');
        if (errText) {
            $('#loginError').show();
        } else {
            $('#loginError').hide();
        }

        $('#togglePwd').on('click', function(){
            var input = document.getElementById('pwd');
            if (!input) return;
            if (input.type === 'password') {
                input.type = 'text';
                this.textContent = '隐藏';
            } else {
                input.type = 'password';
                this.textContent = '显示';
            }
        });

        $('#forgotPwd').on('click', function(){
            layer.msg('请联系管理员重置密码');
        });

        $('.login-third__icon').on('click', function(){
            layer.msg('第三方登录暂未开放');
        });

        var saved = window.localStorage.getItem('remember_loginname');
        if (saved) {
            $('#loginname').val(saved);
            $('#rememberMe').prop('checked', true);
            form.render('checkbox');
        }

        form.on('submit(login)', function(){
            var btn = $('#loginBtn');
            btn.text('登录中...').attr('disabled','disabled').addClass('layui-disabled');
            layer.load(2, {shade: [0.18,'#000']});
            if ($('#rememberMe').is(':checked')) {
                window.localStorage.setItem('remember_loginname', $('#loginname').val() || '');
            } else {
                window.localStorage.removeItem('remember_loginname');
            }
            setTimeout(function(){
                document.getElementById('loginFrm').submit();
            }, 200);
            return false;
        });
    })

</script>
</body>
</html>
