{{define "content"}}
<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Audit Log: {{.Data.Date}}</h1>
</div>
<div class="table-responsive">
    <table class="table table-striped table-sm">
        <tbody>
        {{range .Data.Events}}
        <tr>
            <td>{{.}}</td>
        </tr>
        {{end}}
        </tbody>
    </table>
</div>
{{end}}