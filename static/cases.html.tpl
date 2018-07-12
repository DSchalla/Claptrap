{{define "content"}}
<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">Cases: {{.Data.Type}}</h1>
</div>
<table class="table table-striped table-sm">
    <thead class="thead-dark">
    <tr>
        <th scope="col">Case ID</th>
        <th scope="col">#Conditions</th>
        <th scope="col">#Responses</th>
        <th scope="col">Actions</th>
    </tr>
    </thead>
    <tbody>
    {{range .Data.Cases}}
    <tr>
        <td>{{.Name}}</td>
        <td>{{.NumConditions}}</td>
        <td>{{.NumResponses}}</td>
        <td>
            <form action="/plugins/com.dschalla.claptrap/cases/{{.Type}}/{{.Name}}/delete" method="POST">
                <button type="submit">Delete</button>
            </form>
        </td>
    </tr>
    {{end}}
    </tbody>
</table>
{{end}}
