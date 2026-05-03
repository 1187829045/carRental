<%--
  Created by IntelliJ IDEA.
  Car: YQF
  Date: 2019/10/14
  Time: 18:50
  To change this template use File | Settings | File Templates.
--%>
<!DOCTYPE html>
<%@ page contentType="text/html;charset=UTF-8" language="java" %>
<html>
<head>
    <meta charset="utf-8">
    <title>出租汽车管理</title>
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
    <style>
        .rent-page-card{
            background: #fff;
            border: 1px solid rgba(47, 128, 255, 0.10);
            border-radius: 18px;
            box-shadow: 0 18px 38px rgba(15, 23, 42, 0.06);
            padding: 18px 20px;
            margin-bottom: 18px;
        }
        .rent-page-card .layui-field-title{
            margin-top: 0 !important;
            margin-bottom: 12px;
        }
        .rent-page-card .layui-field-title legend{
            font-size: 18px;
            font-weight: 700;
            color: #13223a;
        }
        .rent-search-row{
            display: flex;
            align-items: flex-end;
            gap: 14px;
            flex-wrap: wrap;
        }
        .rent-search-field{
            min-width: 280px;
            flex: 1 1 320px;
        }
        .rent-field-label{
            display: block;
            margin-bottom: 8px;
            font-size: 13px;
            font-weight: 600;
            color: #5f7392;
        }
        .rent-page-card .layui-input,
        .rent-page-card .layui-textarea{
            height: 42px;
            border-radius: 12px;
            border: 1px solid #d9e7ff;
            background: #fbfdff;
            padding-left: 14px;
        }
        .rent-page-card .layui-input:focus,
        .rent-page-card .layui-textarea:focus{
            border-color: rgba(47, 128, 255, 0.65) !important;
            box-shadow: 0 0 0 4px rgba(47, 128, 255, 0.12);
            background: #fff;
        }
        .rent-action-group{
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }
        .rent-table-card{
            background: #fff;
            border: 1px solid rgba(47, 128, 255, 0.10);
            border-radius: 18px;
            box-shadow: 0 18px 38px rgba(15, 23, 42, 0.06);
            padding: 16px;
        }
        .rent-table-card__head{
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 12px;
            flex-wrap: wrap;
            margin-bottom: 12px;
        }
        .rent-table-card__title{
            font-size: 18px;
            font-weight: 700;
            color: #13223a;
        }
        .rent-table-card__desc{
            font-size: 12px;
            color: #6b7b94;
        }
        .rent-form-grid{
            display: grid;
            grid-template-columns: repeat(2, minmax(0, 1fr));
            gap: 14px 16px;
        }
        .rent-form-grid .layui-form-item{
            margin-bottom: 0;
        }
        .rent-form-grid .layui-input-inline,
        .rent-form-grid .layui-input-block{
            width: 100%;
            margin-left: 0;
        }
        .rent-form-grid .layui-form-label{
            width: auto;
            padding: 0 0 8px;
            line-height: 1.2;
            float: none;
            font-weight: 600;
            color: #5f7392;
        }
        .rent-form-grid .layui-inline{
            width: 100%;
            margin-right: 0;
        }
        .rent-form-grid .layui-input-block{
            min-height: auto;
        }
        .rent-form-grid .layui-form-item.is-full{
            grid-column: 1 / -1;
        }
        .rent-submit-row{
            text-align: center;
            padding-top: 10px;
        }
        .rent-loading-tip{
            display: inline-flex;
            align-items: center;
            gap: 8px;
            color: #6b7b94;
            font-size: 12px;
        }
        .rent-loading-tip i{
            color: #2f80ff;
        }
        @media screen and (max-width: 768px){
            .rent-page-card,
            .rent-table-card{
                padding: 14px;
                border-radius: 16px;
            }
            .rent-search-field{
                min-width: 100%;
                flex-basis: 100%;
            }
            .rent-form-grid{
                grid-template-columns: 1fr;
            }
            #view_carimg{
                width: 100%;
                height: auto;
            }
        }
    </style>
</head>
<body class="childrenBody">

<!-- 搜索条件开始 -->
<div class="rent-page-card">
    <fieldset class="layui-elem-field layui-field-title">
        <legend>查询条件</legend>
    </fieldset>
    <form class="layui-form" method="post" id="searchFrm">
        <div class="rent-search-row">
            <div class="rent-search-field">
                <label class="rent-field-label">身份证号</label>
                <input type="text" name="identity" id="identity" autocomplete="off"
                       class="layui-input"
                       placeholder="请输入身份证号后查询可租车辆">
            </div>
            <div class="rent-action-group">
                <button type="button"
                        class="layui-btn layui-btn-normal layui-icon layui-icon-search layui-btn-radius"
                        id="doSearch">查询
                </button>
                <button type="reset"
                        class="layui-btn layui-btn-warm layui-icon layui-icon-refresh layui-btn-radius"
                        id="resetSearchBtn">重置
                </button>
            </div>
        </div>
    </form>
</div>

<!-- 数据表格开始 -->
<div id="content" class="rent-table-card" style="display: none;">
    <div class="rent-table-card__head">
        <div>
            <div class="rent-table-card__title">可租车辆列表</div>
            <div class="rent-table-card__desc">默认每页展示 20 条数据，支持切换到 25 / 30 条，按更新时间倒序展示</div>
        </div>
        <div class="rent-loading-tip" id="tableLoadingHint">
            <i class="layui-icon layui-icon-loading layui-anim layui-anim-rotate layui-anim-loop"></i>
            <span>等待查询</span>
        </div>
    </div>
    <table id="carTable" lay-filter="carTable"></table>
    <div id="carBar" style="display: none;">
        <a class="layui-btn layui-btn-warm layui-btn-xs layui-btn-radius" lay-event="rentCar">租赁汽车</a>
        <a class="layui-btn layui-btn-xs layui-btn-radius" lay-event="viewImage">查看车辆大图</a>
    </div>
</div>

<%--添加和修改的弹出层开始--%>
<div style="display: none;padding: 20px;" id="saveOrUpdateDiv">
    <form class="layui-form" lay-filter="dataFrm" id="dataFrm">
        <div class="rent-form-grid">
        <div class="layui-form-item is-full">
            <label class="layui-form-label">出租单号</label>
            <div class="layui-input-block">
                <input type="text" name="rentid" lay-verify="required" readonly="readonly" placeholder="请输入出租单号"
                       class="layui-input">
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">起租时间</label>
                <div class="layui-input-inline">
                    <input type="text" name="begindate" id="begindate" lay-verify="required" placeholder="请输入起租时间" class="layui-input">
                </div>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">预计还车时间</label>
                <div class="layui-input-inline">
                    <input type="text" name="returndate" id="returndate" lay-verify="required" placeholder="请输入预计还车时间" class="layui-input">
                </div>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">身份证号</label>
                <div class="layui-input-inline">
                    <input type="text" name="identity" id="rentIdentity" lay-verify="required" placeholder="请输入身份证号"
                           class="layui-input">
                </div>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">客户名称</label>
                <div class="layui-input-inline">
                    <input type="text" name="opername" id="opername" lay-verify="required" placeholder="请输入客户名称" class="layui-input">
                </div>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">车牌号</label>
                <div class="layui-input-inline">
                    <input type="text" name="carnumber" lay-verify="required" readonly="readonly"  placeholder="请输入车牌号" class="layui-input">
                </div>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">出租价格</label>
                <div class="layui-input-inline">
                    <input type="text" name="price" lay-verify="required" readonly="readonly" placeholder="请输入出租价格" class="layui-input">
                </div>
            </div>
        </div>
        <div class="layui-form-item is-full">
            <div class="layui-input-block rent-submit-row">
                <button type="button"
                        class="layui-btn layui-btn-normal layui-btn-md layui-icon layui-icon-release layui-btn-radius"
                        lay-filter="doSubmit" lay-submit="">提交
                </button>
                <button type="reset"
                        class="layui-btn layui-btn-warm layui-btn-md layui-icon layui-icon-refresh layui-btn-radius">重置
                </button>
            </div>
        </div>
        </div>
    </form>
</div>

<%--查看大图弹出的层开始--%>
<div id="viewCarImageDiv" style="display: none;text-align: center">
    <img alt="出租图片" width="700px" height="450px" id="view_carimg">
</div>

<script src="${yeqifu}/static/layui/layui.js"></script>
<script type="text/javascript">
    var tableIns;
    layui.use(['jquery', 'layer', 'form', 'table', 'laydate'], function () {
        var $ = layui.jquery;
        var layer = layui.layer;
        var form = layui.form;
        var table = layui.table;
        var dtree = layui.dtree;
        var laydate = layui.laydate;
        var currentSelectedCar = null;
        var carTableInitialized = false;
        var activeSearchIdentity = "";
        var tableLoadingIndex = null;

        laydate.render({
            elem:'#begindate',
            type:'datetime'
        });
        laydate.render({
            elem:'#returndate',
            type:'datetime'
        });

        var defaultCarImage = "https://images.pexels.com/photos/170811/pexels-photo-170811.jpeg?auto=compress&cs=tinysrgb&w=1200";

        function buildCarImageUrl(path) {
            var finalPath = path || defaultCarImage;
            return "${yeqifu}/file/downloadShowFile.action?path=" + encodeURIComponent(finalPath);
        }

        function setTableLoading(text) {
            $("#tableLoadingHint span").text(text || "加载中...");
        }

        function openTableLoading(text) {
            closeTableLoading();
            setTableLoading(text || "加载中...");
            tableLoadingIndex = layer.load(1, {shade: [0.12, '#fff']});
        }

        function closeTableLoading(text) {
            if (tableLoadingIndex !== null) {
                layer.close(tableLoadingIndex);
                tableLoadingIndex = null;
            }
            if (text) {
                setTableLoading(text);
            }
        }

        function setRentFormCarData(car) {
            if (!car) {
                return;
            }
            $("input[name='carnumber']").val(car.carnumber || "");
            $("input[name='price']").val(car.rentprice || "");
        }

        function applyRentFormData(data, selectedCar) {
            var merged = $.extend({}, data || {});
            if (selectedCar) {
                merged.carnumber = selectedCar.carnumber || "";
                merged.price = selectedCar.rentprice || "";
            }
            form.val("dataFrm", merged);
        }

        function fillRentCustomer(identity, selectedCar, silent) {
            var normalized = normalizeIdentity(identity);
            $("#rentIdentity").val(normalized);
            setRentFormCarData(selectedCar);
            if (!normalized) {
                $("#opername").val("");
                return;
            }
            if (!isValidIdentity(normalized)) {
                if (!silent) {
                    layer.msg("请输入正确的身份证号");
                }
                $("#opername").val("");
                return;
            }
            $.get("${yeqifu}/rent/initRentFrom.action", {
                identity: normalized,
                price: selectedCar ? selectedCar.rentprice : "",
                carnumber: selectedCar ? selectedCar.carnumber : ""
            }, function (obj) {
                if (!obj || !obj.rentid) {
                    $("#opername").val("");
                    if (!silent) {
                        layer.msg("未找到对应客户信息");
                    }
                    return;
                }
                applyRentFormData(obj, selectedCar);
            });
        }

        function normalizeIdentity(value) {
            return $.trim(value || "").replace(/[^0-9xX]/g, "").toUpperCase();
        }

        function isValidIdentity(value) {
            return /^(?:\d{15}|\d{17}[\dX])$/.test(value);
        }

        function initCarData(identity) {
            activeSearchIdentity = identity || activeSearchIdentity;
            if (carTableInitialized && tableIns) {
                openTableLoading("正在刷新车辆列表...");
                tableIns.reload({
                    url: '${yeqifu}/car/loadAllCar.action',
                    where: {
                        isrenting: 0
                    },
                    page: {
                        curr: 1,
                        limit: 20
                    }
                });
                return;
            }
            carTableInitialized = true;
            openTableLoading("正在加载车辆列表...");
            //渲染数据表格
            tableIns = table.render({
                elem: '#carTable'   //渲染的目标对象
                , url: '${yeqifu}/car/loadAllCar.action' //数据接口
                , where: {isrenting: 0}
                , title: '车辆列表'//数据导出来的标题
                , height: 'full-150'
                , page: true  //是否启用分页
                , limit: 20
                , limits: [20, 25, 30]
                , cols: [[   //列表数据
                    {field: 'carnumber', title: '车牌号', align: 'center', width: '104'}
                    , {field: 'cartype', title: '出租类型', align: 'center', width: '90'}
                    , {field: 'color', title: '出租颜色', align: 'center', width: '90'}
                    , {field: 'price', title: '汽车价格', align: 'center', width: '90'}
                    , {field: 'rentprice', title: '出租价格', align: 'center', width: '90'}
                    , {field: 'deposit', title: '出租押金', align: 'center', width: '90'}
                    , {
                        field: 'isrenting', title: '出租状态', align: 'center', width: '90', templet: function (d) {
                            return d.isrenting == '1' ? '<font color=blue>已出租</font>' : '<font color=red>未出租</font>';
                        }
                    }
                    , {field: 'description', title: '出租描述', align: 'center', width: '160'}
                    , {
                        field: 'carimg', title: '缩略图', align: 'center', width: '80', templet: function (d) {
                            return "<img width=40 height=40 style='object-fit:cover;border-radius:6px;' src='" + buildCarImageUrl(d.carimg) + "' onerror=\"this.onerror=null;this.src='" + buildCarImageUrl(defaultCarImage) + "';\"/>";
                        }
                    }
                    , {field: 'createtime', title: '录入时间', align: 'center', width: '170'}
                    , {fixed: 'right', title: '操作', toolbar: '#carBar', align: 'center', width: '220'}
                ]]
                , text: {none: '暂无可租赁车辆数据'}
                , done: function (res, curr, count) {
                    closeTableLoading('已加载 ' + count + ' 条可租赁车辆数据');
                }
            });

        }

        //模糊查询
        $("#doSearch").click(function () {
            var identity = normalizeIdentity($("#identity").val());
            $("#identity").val(identity);
            if (!identity) {
                layer.msg("请输入身份证号");
                $("#content").hide();
                setTableLoading("等待查询");
                return;
            }
            if (!isValidIdentity(identity)) {
                layer.msg("请输入正确的身份证号");
                $("#content").hide();
                setTableLoading("请输入有效身份证号");
                return;
            }
            var params = $("#searchFrm").serialize();
            openTableLoading("正在校验客户信息...");
            $.post("${yeqifu}/rent/checkCustomerExist.action", params, function (obj) {
                if (obj.code >= 0) { //此客户存在，code的返回值为0
                    $("#content").show();
                    initCarData(identity); //初始化未出租汽车的所有数据
                } else {
                    closeTableLoading("未找到可用客户信息");
                    layer.msg("客户身份证号不存在，请更正后在查询");
                    //隐藏数据表格
                    $("#content").hide();
                }
            }).fail(function () {
                closeTableLoading("客户校验失败");
                layer.msg("客户信息校验失败，请稍后重试");
                $("#content").hide();
            });
        });

        $("#resetSearchBtn").click(function () {
            $("#content").hide();
            setTableLoading("等待查询");
        });

        //监听行工具事件
        table.on('tool(carTable)', function (obj) {
            var data = obj.data; //获得当前行数据
            var layEvent = obj.event; //获得 lay-event 对应的值（也可以是表头的 event 参数对应的值）
            if (layEvent === 'rentCar') { //汽车出租
                //汽车出租，打开添加汽车出租页面
                openRentCar(data);
            } else if (layEvent === 'viewImage') { //查看大图
                showCarImage(data);
            }
        });

        var mainIndex;

        //打开添加页面
        function openRentCar(data) {
            currentSelectedCar = data;
            mainIndex = layer.open({
                type: 1,
                title: '添加汽车出租',
                content: $("#saveOrUpdateDiv"),
                area: ['690px', '380px'],
                success: function (index) {
                    //清空表单数据
                    $("#dataFrm")[0].reset();
                    setRentFormCarData(data);
                    fillRentCustomer($("#identity").val(), data, true);
                }
            });
        }

        $("#rentIdentity").on("input blur", function () {
            if (!currentSelectedCar) {
                return;
            }
            fillRentCustomer($(this).val(), currentSelectedCar, true);
        });

        //保存
        form.on("submit(doSubmit)", function (obj) {
            //序列化表单数据
            var params = $("#dataFrm").serialize();
            openTableLoading("正在提交出租单...");
            $.post("${yeqifu}/rent/saveRent.action", params, function (obj) {
                closeTableLoading("出租单提交完成");
                layer.msg(obj.msg);
                //关闭弹出层
                layer.close(mainIndex);
                $("#content").hide();
            }).fail(function () {
                closeTableLoading("出租单提交失败");
                layer.msg("出租单保存失败，请稍后重试");
            });
        });

        //查看大图
        function showCarImage(data) {
            mainIndex = layer.open({
                type: 1,
                title: "【" + data.carnumber + '】的出租图片',
                content: $("#viewCarImageDiv"),
                area: ['750px', '500px'],
                success: function (index) {
                    $("#view_carimg").attr("src", buildCarImageUrl(data.carimg));
                }
            });
        }

    });

</script>
</body>
</html>
