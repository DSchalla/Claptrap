{{define "response-message_ephemeral"}}
<div id="response-template-message_ephemeral" class="card responses-response">
    <div class="card-body">
        <h5>Send Ephemeral Message<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="responses-{{.Id}}-message">Message</label>
                <input type="text" class="form-control" id="responses-{{.Id}}-message" name="responses[{{.Id}}][message]" value="{{.Message}}" required>
                <small class="text-muted">Message sent towards the channel</small>
                <div class="invalid-feedback">
                    Message is required
                </div>
                <input type="hidden" name="responses[{{.Id}}][type]" value="message_ephemeral">
            </div>
        </div>
    </div>
</div>
{{end}}