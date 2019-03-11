<div class="table-responsive" id="templateListContainer">
    <table class="table table-hover">
        <thead>
        <tr>
            <td>#</td>
            <td class="col-sm-3">模板名称</td>
            <td class="col-sm-2">模板类型</td>
            <td class="col-sm-2">创建人</td>
            <td class="col-sm-3">创建时间</td>
            <td class="col-sm-2">操作</td>
        </tr>
        </thead>
        <tbody>
        {{if .ErrorMessage}}
        <tr>
            <td colspan="6" class="text-center">{{.ErrorMessage}}</td>
        </tr>
        {{else}}
        {{range $index,$item := .List}}
        <tr>
            <td>{{$item.TemplateId}}</td>
            <td>{{$item.TemplateName}}</td>
            <td>{{if $item.IsGlobal}}全局{{else}}项目{{end}}</td>
            <td>{{$item.CreateName}}</td>
            <td>{{date_format $item.CreateTime "2006-01-02 15:04:05"}}</td>
            <td>
                <button class="btn btn-primary btn-sm btn-insert" data-id="{{$item.TemplateId}}">
                    插入
                </button>
                <button class="btn btn-danger btn-sm btn-delete" data-id="{{$item.TemplateId}}" data-loading-text="删除中...">
                    删除
                </button>
            </td>
        </tr>
        {{else}}
        <tr>
            <td colspan="6" class="text-center">暂无数据</td>
        </tr>
        {{end}}
        {{end}}
        </tbody>
    </table>
</div>
