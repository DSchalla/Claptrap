{{define "response-delete_message"}}
<div id="response-template-delete_message" class="card responses-response">
    <div class="card-body">
        <h5>Delete Message<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <small class="text-muted">This response has no options.</small>
                <input type="hidden" name="responses[{{.Id}}][type]" value="delete_message">
            </div>
        </div>
    </div>
</div>
{{end}}