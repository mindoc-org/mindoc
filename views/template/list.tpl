<div class="table-responsive" id="templateListContainer">
    <table class="table table-hover">
        <thead>
        <tr>
            <td>#</td>
            <td class="col-sm-3">{{i18n $.Lang "doc.tpl_name"}}</td>
            <td class="col-sm-2">{{i18n $.Lang "doc.tpl_type"}}</td>
            <td class="col-sm-2">{{i18n $.Lang "doc.creator"}}</td>
            <td class="col-sm-3">{{i18n $.Lang "doc.create_time"}}</td>
            <td class="col-sm-2">{{i18n $.Lang "doc.operation"}}</td>
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
            <td>{{if $item.IsGlobal}}{{i18n .Lang "doc.global_tpl"}}{{else}}{{i18n .Lang "doc.project_tpl"}}{{end}}</td>
            <td>{{$item.CreateName}}</td>
            <td>{{date_format $item.CreateTime "2006-01-02 15:04:05"}}</td>
            <td>
                <button class="btn btn-primary btn-sm btn-insert" data-id="{{$item.TemplateId}}">
                    {{i18n .Lang "doc.insert"}}
                </button>
                <button class="btn btn-danger btn-sm btn-delete" data-id="{{$item.TemplateId}}" data-loading-text="删除中...">
                    {{i18n .Lang "doc.delete"}}
                </button>
            </td>
        </tr>
        {{else}}
        <tr>
            <td colspan="6" class="text-center">{{i18n .Lang "message.no_data"}}</td>
        </tr>
        {{end}}
        {{end}}
        </tbody>
    </table>
</div>
