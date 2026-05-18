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
    <style>
        .login-error-banner{
            display: none;
            margin: 0 0 14px;
            padding: 12px 14px;
            border: 1px solid #ffb3b3;
            border-radius: 10px;
            background: #fff1f0;
            color: #cf1322;
            font-size: 14px;
            line-height: 1.5;
            box-shadow: 0 6px 16px rgba(207,19,34,0.12);
        }
        .login-error-banner.is-visible{
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .login-error-banner .layui-icon{
            font-size: 18px;
        }
        .login-error-state .layui-input{
            border-color: #ff4d4f !important;
            box-shadow: 0 0 0 3px rgba(255,77,79,0.12);
        }
        .captcha-flash{
            animation: loginFlash 0.7s ease-in-out 2;
        }
        .field-shake{
            animation: loginShake 0.42s ease-in-out 2;
        }
        @keyframes loginShake{
            0%,100%{transform:translateX(0);}
            20%{transform:translateX(-6px);}
            40%{transform:translateX(6px);}
            60%{transform:translateX(-4px);}
            80%{transform:translateX(4px);}
        }
        @keyframes loginFlash{
            0%,100%{opacity:1; transform:scale(1);}
            50%{opacity:.35; transform:scale(1.06);}
        }
        .sr-live-region{
            position: absolute;
            width: 1px;
            height: 1px;
            padding: 0;
            margin: -1px;
            overflow: hidden;
            clip: rect(0, 0, 0, 0);
            white-space: nowrap;
            border: 0;
        }
        @media (max-width: 768px){
            .login-error-banner{
                margin-bottom: 12px;
                padding: 10px 12px;
                font-size: 13px;
            }
        }
    </style>
</head>
<body class="loginBody">
<form style="width: 400px; height: 600px;"  class="layui-form" id="loginFrm" method="post" action="${yeqifu}/login/login.action" data-error-type="${errorType}" data-error-message="${error}">
    <div class="login-brand-title">汽车租赁系统</div>
    <div class="login_face"><img src="${yeqifu}/static/images/face.jpg" class="userAvatar"></div>
    <div id="loginAriaLive" class="sr-live-region" aria-live="assertive" aria-atomic="true"></div>
    <div id="loginErrorBanner" class="login-error-banner" role="alert" aria-live="assertive" aria-atomic="true">
        <i class="layui-icon layui-icon-close-fill" aria-hidden="true"></i>
        <span id="loginErrorText"></span>
    </div>
    <div class="layui-form-item input-item">
        <label for="loginname">用户名</label>
        <input type="text" placeholder="请输入用户名" autocomplete="off" name="loginname" id="loginname" value="${loginname}" class="layui-input" lay-verify="required">
    </div>
    <div class="layui-form-item input-item" id="pwdField">
        <label for="pwd">密码</label>
        <input type="password" placeholder="请输入密码" autocomplete="off" name="pwd" id="pwd" class="layui-input" lay-verify="required">
    </div>
    <div class="layui-form-item input-item" id="imgCode">
        <label for="code">验证码1</label>
        <input type="text" placeholder="请输入验证码" autocomplete="off" name="code" id="code" class="layui-input">
        <img style="width: 200px; height: 30px;" id="captchaImg" src="${yeqifu}/login/getCode.action" onclick="this.src='${yeqifu}/login/getCode.action?ts='+Date.now()" alt="captcha">
    </div>
    <div class="layui-form-item">
        <button class="layui-btn layui-block" id="loginBtn" lay-filter="login" lay-submit>登录</button>
    </div>
</form>
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

        function showBanner(message){
            $('#loginErrorText').text(message || '');
            $('#loginErrorBanner').addClass('is-visible');
            $('#loginAriaLive').text(message || '');
        }

        function clearErrorState(){
            $('#loginErrorBanner').removeClass('is-visible');
            $('#loginAriaLive').text('');
            $('#pwdField, #imgCode').removeClass('login-error-state field-shake');
            $('#captchaImg').removeClass('captcha-flash');
            $('#pwd').removeAttr('aria-invalid');
            $('#code').removeAttr('aria-invalid');
        }

        function triggerCaptchaError(message){
            showBanner(message);
            $('#imgCode').addClass('login-error-state field-shake');
            $('#captchaImg').addClass('captcha-flash');
            $('#captchaImg').attr('src', '${yeqifu}/login/getCode.action?ts=' + Date.now());
            $('#code').val('').attr('aria-invalid', 'true').trigger('focus');
            setTimeout(function(){
                $('#imgCode').removeClass('field-shake');
                $('#captchaImg').removeClass('captcha-flash');
            }, 900);
        }

        function triggerPasswordError(message){
            $('#loginAriaLive').text(message || '');
            $('#pwdField').addClass('login-error-state');
            $('#pwd').attr('aria-invalid', 'true');
            layer.alert(message, {
                title: '登录失败',
                icon: 2,
                shadeClose: false,
                closeBtn: 0,
                anim: 6
            }, function(index){
                layer.close(index);
                $('#pwd').trigger('focus');
            });
        }

        form.on('submit(login)', function(){
            clearErrorState();
            var btn = $('#loginBtn');
            btn.text('登录中...').attr('disabled','disabled').addClass('layui-disabled');
            setTimeout(function(){
                document.getElementById('loginFrm').submit();
            }, 500);
            return false;
        });

        $(".loginBody .input-item").click(function(e){
            e.stopPropagation();
            $(this).addClass("layui-input-focus").find(".layui-input").focus();
        })
        $(".loginBody .layui-form-item .layui-input").focus(function(){
            $(this).parent().addClass("layui-input-focus");
        })
        $(".loginBody .layui-form-item .layui-input").blur(function(){
            $(this).parent().removeClass("layui-input-focus");
            if($(this).val() != ''){
                $(this).parent().addClass("layui-input-active");
            }else{
                $(this).parent().removeClass("layui-input-active");
            }
        })
        $("#pwd").on("input", function(){
            $('#pwdField').removeClass('login-error-state');
            $(this).removeAttr('aria-invalid');
        });
        $("#code").on("input", function(){
            $('#imgCode').removeClass('login-error-state');
            $(this).removeAttr('aria-invalid');
        });

        var errorType = $('#loginFrm').data('error-type') || '';
        var errorMessage = $('#loginFrm').data('error-message') || '';
        if(errorType === 'captcha' && errorMessage){
            triggerCaptchaError(errorMessage);
        }else if(errorType === 'password' && errorMessage){
            triggerPasswordError(errorMessage);
        }
    })

</script>
</body>
</html>
