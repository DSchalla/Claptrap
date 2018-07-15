{{define "response-intercept_event"}}
<div id="response-template-delete_message" class="card responses-response">
    <div class="card-body">
        <h5>Intercept Event<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <small class="text-muted">Intercept the original event and prohibit the message to be send, channel to be joined etc.
                    This response has no options.</small>
                <input type="hidden" name="responses[{{.Id}}][type]" value="intercept_event">
            </div>
        </div>
    </div>
</div>
{{end}}