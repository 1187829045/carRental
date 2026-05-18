<%--
  Created by IntelliJ IDEA.
  User: YQF
  Date: 2019/9/30
  Time: 22:57
  To change this template use File | Settings | File Templates.
--%>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>角色管理</title>
    <meta name="renderer" content="webkit">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta http-equiv="Access-Control-Allow-Origin" content="*">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <meta name="apple-mobile-web-app-status-bar-style" content="black">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="format-detection" content="telephone=no">
    <%--<link rel="icon" href="favicon.ico">--%>
    <link rel="stylesheet" href="${yeqifu}/static/layui/css/layui.css" media="all"/>
    <link rel="stylesheet" href="${yeqifu}/static/css/public.css" media="all"/>
    <link rel="stylesheet" href="${yeqifu}/static/layui_ext/dtree/dtree.css">
    <link rel="stylesheet" href="${yeqifu}/static/layui_ext/dtree/font/dtreefont.css">
    <style>
        html, body {
            height: 100%;
            box-sizing: border-box;
        }
        body.childrenBody {
            overflow: hidden;
        }
        .layui-table-view {
            margin-top: 0;
        }
    </style>
</head>
<body class="childrenBody">

<!-- 搜索条件开始 -->
<fieldset class="layui-elem-field layui-field-title" style="margin-top: 20px;">
    <legend>查询条件</legend>
</fieldset>
<form class="layui-form" method="post" id="searchFrm">
    <div class="layui-form-item">
        <div class="layui-inline">
            <label class="layui-form-label">角色名称:</label>
            <div class="layui-input-inline" style="padding: 5px">
                <input type="text" name="rolename" autocomplete="off" class="layui-input layui-input-inline"
                       placeholder="请输入角色名称" style="height: 30px;border-radius: 10px">
            </div>
        </div>
        <div class="layui-inline">
            <label class="layui-form-label">角色备注:</label>
            <div class="layui-input-inline" style="padding: 5px">
                <input type="text" name="roledesc" autocomplete="off" class="layui-input layui-input-inline"
                       placeholder="请输入角色备注" style="height: 30px;border-radius: 10px">
            </div>
        </div>
        <div class="layui-inline">
            <label class="layui-form-label">是否可用:</label>
            <div class="layui-input-inline">
                <input type="radio" name="available" value="1" title="可用">
                <input type="radio" name="available" value="0" title="不可用">
            </div>
        </div>
        <div class="layui-inline">
            <button type="button" class="layui-btn layui-btn-normal layui-icon layui-icon-search layui-btn-radius layui-btn-sm" id="doSearch">查询</button>
            <button type="reset" class="layui-btn layui-btn-warm layui-icon layui-icon-refresh layui-btn-radius layui-btn-sm">重置</button>
        </div>
    </div>
</form>

<!-- 数据表格开始 -->
<table class="layui-hide" id="roleTable" lay-filter="roleTable"></table>
<div style="display: none;" id="roleToolBar">
    <button type="button" class="layui-btn layui-btn-sm layui-btn-radius" lay-event="add">增加</button>
    <button type="button" class="layui-btn layui-btn-danger layui-btn-sm layui-btn-radius" lay-event="deleteBatch">批量删除</button>
</div>
<div id="roleBar" style="display: none;">
    <a class="layui-btn layui-btn-xs layui-btn-radius" lay-event="edit">编辑</a>
    <a class="layui-btn layui-btn-warm layui-btn-xs layui-btn-radius" lay-event="selectRoleUser">分配用户</a>
    <a class="layui-btn layui-btn-danger layui-btn-xs layui-btn-radius" lay-event="del">删除</a>
</div>

<!-- 添加和修改的弹出层-->
<div style="display: none;padding: 20px" id="saveOrUpdateDiv">
    <form class="layui-form" lay-filter="dataFrm" id="dataFrm">

        <div class="layui-form-item">
            <label class="layui-form-label">角色名称:</label>
            <div class="layui-input-block">
                <input type="hidden" name="roleid">
                <input type="text" name="rolename" placeholder="请输入角色名称" autocomplete="off" class="layui-input">
            </div>
        </div>
        <div class="layui-form-item">
            <label class="layui-form-label">角色备注:</label>
            <div class="layui-input-block">
                <input type="text" name="roledesc" placeholder="请输入角色备注" autocomplete="off" class="layui-input">
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">是否可用:</label>
                <div class="layui-input-inline">
                    <input type="radio" name="available" value="1" checked="checked" title="可用">
                    <input type="radio" name="available" value="0" title="不可用">
                </div>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-input-block" style="text-align: center;padding-right: 120px">
                <button type="button"
                        class="layui-btn layui-btn-normal layui-btn-md layui-icon layui-icon-release layui-btn-radius"
                        lay-filter="doSubmit" lay-submit="">提交
                </button>
                <button type="reset"
                        class="layui-btn layui-btn-warm layui-btn-md layui-icon layui-icon-refresh layui-btn-radius">重置
                </button>
            </div>
        </div>
    </form>
</div>

<%--角色分配用户的弹出层开始--%>
<div style="display: none;padding: 15px" id="selectRoleUser">
    <form class="layui-form" id="roleUserSearchFrm">
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">用户姓名:</label>
                <div class="layui-input-inline">
                    <input type="text" name="realname" autocomplete="off" class="layui-input" placeholder="请输入用户姓名">
                </div>
            </div>
            <div class="layui-inline">
                <label class="layui-form-label">登录名:</label>
                <div class="layui-input-inline">
                    <input type="text" name="loginname" autocomplete="off" class="layui-input" placeholder="请输入登录名">
                </div>
            </div>
            <div class="layui-inline">
                <button type="button" class="layui-btn layui-btn-sm layui-btn-normal" id="doUserSearch">查询</button>
                <button type="reset" class="layui-btn layui-btn-sm layui-btn-warm" id="resetUserSearch">重置</button>
            </div>
        </div>
    </form>
    <table class="layui-hide" id="roleUserTable" lay-filter="roleUserTable"></table>
</div>


<script src="${yeqifu}/static/layui/layui.js"></script>
<script type="text/javascript">
    var tableIns;
    layui.use(['jquery', 'layer', 'form', 'table'], function () {
        var $ = layui.jquery;
        var layer = layui.layer;
        var form = layui.form;
        var table = layui.table;
        var roleUserTableIns;
        var selectedRoleUserIds = {};
        var currentRoleId;
        //渲染数据表格
        tableIns = table.render({
            elem: '#roleTable'   //渲染的目标对象
            , url: '${yeqifu}/role/loadAllRole.action' //数据接口
            , title: '用户数据表'//数据导出来的标题
            , toolbar: "#roleToolBar"   //表格的工具条
            , height: 'full-150'
            , cellMinWidth: 100 //设置列的最小默认宽度
            , page: true  //是否启用分页
            , cols: [[   //列表数据
                {type: 'checkbox', fixed: 'left'}
                , {field: 'roleid', title: 'ID', align: 'center'}
                , {field: 'rolename', title: '角色名称', align: 'center'}
                , {field: 'roledesc', title: '角色备注', align: 'center'}
                , {field: 'assignedusers', title: '已分配用户', align: 'center', minWidth: 180, templet: function (d) {
                    return d.assignedusers ? d.assignedusers : '<span class="layui-badge layui-bg-gray">暂无用户</span>';
                }}
                , {
                    field: 'available', title: '是否可用', align: 'center', templet: function (d) {
                        return d.available == '1' ? '<font color=blue>可用</font>' : '<font color=red>不可用</font>';
                    }
                }
                , {fixed: 'right', title: '操作', toolbar: '#roleBar', align: 'center'}
            ]],
            done:function (data, curr, count) {
                //不是第一页时，如果当前返回的数据为0那么就返回上一页
                if(data.data.length==0&&curr!=1){
                    tableIns.reload({
                        page:{
                            curr:curr-1
                        }
                    })
                }
            }
        })

        //模糊查询
        $("#doSearch").click(function () {
            var params = $("#searchFrm").serialize();
            //alert(params);
            tableIns.reload({
                url: "${yeqifu}/role/loadAllRole.action?" + params,
                page:{curr:1}
            })
        });

        //监听头部工具栏事件
        table.on("toolbar(roleTable)", function (obj) {
            switch (obj.event) {
                case 'add':
                    openAddRole();
                    break;
                case 'deleteBatch':
                    deleteBatch();
                    break;
            }
        });

        //监听行工具事件
        table.on('tool(roleTable)', function (obj) {
            var data = obj.data; //获得当前行数据
            var layEvent = obj.event; //获得 lay-event 对应的值（也可以是表头的 event 参数对应的值）
            if (layEvent === 'del') { //删除
                layer.confirm('真的删除【' + data.rolename + '】这个角色么？', function (index) {
                    //向服务端发送删除指令
                    $.post("${yeqifu}/role/deleteRole.action", {roleid: data.roleid}, function (res) {
                        layer.msg(res.msg);
                        //刷新数据表格
                        tableIns.reload();
                    })
                });
            } else if (layEvent === 'edit') { //编辑
                //编辑，打开修改界面
                openUpdateRole(data);
            }else if(layEvent === 'selectRoleUser'){//分配用户
                openselectRoleUser(data);
            }
        });

        var url;
        var mainIndex;

        //打开添加页面
        function openAddRole() {
            mainIndex = layer.open({
                type: 1,
                title: '添加角色',
                content: $("#saveOrUpdateDiv"),
                area: ['600px', '300px'],
                success: function (index) {
                    //清空表单数据
                    $("#dataFrm")[0].reset();
                    url = "${yeqifu}/role/addRole.action";
                }
            });
        }

        //打开修改页面
        function openUpdateRole(data) {
            mainIndex = layer.open({
                type: 1,
                title: '修改角色',
                content: $("#saveOrUpdateDiv"),
                area: ['600px', '300px'],
                success: function (index) {
                    form.val("dataFrm", data);
                    url = "${yeqifu}/role/updateRole.action";
                }
            });
        }

        //保存
        form.on("submit(doSubmit)", function (obj) {
            //序列化表单数据
            var params = $("#dataFrm").serialize();
            $.post(url, params, function (obj) {
                layer.msg(obj.msg);
                //关闭弹出层
                layer.close(mainIndex)
                //刷新数据 表格
                tableIns.reload();
            })
        });

        //批量删除
        function deleteBatch() {
            //得到选中的数据行
            var checkStatus = table.checkStatus('roleTable');
            var data = checkStatus.data;
            layer.alert(data.length);
            var params="";
            $.each(data,function(i,item){
                if (i==0){
                    params+="ids="+item.roleid;
                }else{
                    params+="&ids="+item.roleid;
                }
            });
            layer.confirm('真的要删除这些角色么？', function (index) {
                //向服务端发送删除指令
                $.post("${yeqifu}/role/deleteBatchRole.action",params, function (res) {
                    layer.msg(res.msg);
                    //刷新数据表格
                    tableIns.reload();
                })
            });
        }

        //打开分配用户的弹出层
        function openselectRoleUser(data) {
            selectedRoleUserIds = {};
            currentRoleId = data.roleid;
            mainIndex=layer.open({
                type:1,
                title:'给【'+data.rolename+'】分配用户',
                content:$("#selectRoleUser"),
                area:['860px','560px'],
                btnAlign:'c',
                btn:['<div class="layui-icon layui-icon-release">保存分配</div>','<div class="layui-icon layui-icon-close">取消</div>'],
                yes:function (index, layero) {
                    var params="roleid="+data.roleid;
                    $.each(selectedRoleUserIds,function (userid, checked) {
                        if (checked) {
                            params+="&ids="+userid;
                        }
                    });
                    //保存角色和用户的关系，空选择时表示清空该角色下的用户
                    $.post("${yeqifu}/role/saveRoleUser.action",params,function (obj) {
                        layer.msg(obj.msg);
                        if(obj.code === 0){
                            layer.close(mainIndex);
                            tableIns.reload();
                        }
                    })
                },
                success:function (index) {
                    $("#roleUserSearchFrm")[0].reset();
                    roleUserTableIns = table.render({
                        elem: '#roleUserTable',
                        url: '${yeqifu}/role/initRoleUser.action?roleid='+data.roleid+'&page=1&limit=1000',
                        title: '用户列表',
                        height: 360,
                        page: false,
                        cols: [[
                            {type: 'checkbox', fixed: 'left'},
                            {field: 'userid', title: 'ID', align: 'center', width: 70},
                            {field: 'realname', title: '用户姓名', align: 'center', width: 120},
                            {field: 'loginname', title: '登录名', align: 'center', width: 120},
                            {field: 'phone', title: '手机号', align: 'center', width: 140},
                            {field: 'identity', title: '身份证号', align: 'center', minWidth: 180},
                            {field: 'available', title: '状态', align: 'center', width: 90, templet: function (d) {
                                return d.available == '1' ? '<font color=blue>可用</font>' : '<font color=red>不可用</font>';
                            }}
                        ]],
                        done:function(res){
                            $.each(res.data || [], function(i, item){
                                if(selectedRoleUserIds[item.userid] === undefined && item.LAY_CHECKED){
                                    selectedRoleUserIds[item.userid] = true;
                                }
                            });
                        }
                    });
                }
            })
        }

        table.on('checkbox(roleUserTable)', function(obj){
            if(obj.type === 'all'){
                $.each(table.cache.roleUserTable || [], function(i, item){
                    if(item && item.userid){
                        selectedRoleUserIds[item.userid] = obj.checked;
                    }
                });
                if(obj.checked){
                    var checkStatus = table.checkStatus('roleUserTable');
                    $.each(checkStatus.data || [], function(i, item){
                        selectedRoleUserIds[item.userid] = true;
                    });
                }
            }else if(obj.data && obj.data.userid){
                selectedRoleUserIds[obj.data.userid] = obj.checked;
            }
        });

        $("#doUserSearch").click(function () {
            if(!roleUserTableIns){
                return;
            }
            roleUserTableIns.reload({
                url:"${yeqifu}/role/initRoleUser.action?roleid="+currentRoleId+"&page=1&limit=1000",
                where: {
                    realname: $("#roleUserSearchFrm input[name='realname']").val(),
                    loginname: $("#roleUserSearchFrm input[name='loginname']").val()
                }
            });
        });

        $("#resetUserSearch").click(function () {
            if(roleUserTableIns){
                setTimeout(function(){
                    roleUserTableIns.reload({
                        url:"${yeqifu}/role/initRoleUser.action?roleid="+currentRoleId+"&page=1&limit=1000",
                        where: {realname: '', loginname: ''}
                    });
                }, 0);
            }
        });
    });

</script>
</body>
</html>
