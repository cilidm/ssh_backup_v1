package template

const (
	Top = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Layui</title>
    <meta name="renderer" content="webkit">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <link rel="stylesheet" type="text/css" href="https://www.layuicdn.com/layui/css/layui.css" />
	<style>.layui-table-cell{font-size:14px;padding:0 5px;height:auto;overflow:visible;text-overflow:inherit;white-space:normal;word-break:break-all}</style>
</head>
<body>
<table class="layui-hide" id="demo"></table>
<script src="https://www.layuicdn.com/layui/layui.js"></script>
`
	Script = `<script>
    layui.use('table', function(){
        var table = layui.table;
        table.render({
            elem: '#demo'
            ,toolbar: true
            ,cols: [[
                {field: 'file_name', title: '文件名', width: 150}
                ,{field: 'file_source', title: '源路径'}
                ,{field: 'file_target', title: '目标路径'}
                ,{field: 'file_size', title: '文件大小', width: 150,sort: true}
                ,{field: 'status', title: '状态', width: 80,sort: true, templet: function(res){
					if(res.status == 1)return '等待中';
					if(res.status == 2)return '成功';
					if(res.status == 3)return '失败';
				  }}
                ,{field: 'created_at', title: '创建时间', width:200, sort: true}
            ]]
            ,skin: 'line' //表格风格
            ,even: true
            ,page: true //是否显示分页
            //,limits: [5, 7, 10]
            //,limit: 5 //每页默认显示的数量
            ,data:`
	Bottom = `,
        });
    });
</script>`
)
