{{define "response-message_user"}}
<div id="response-template-message_user" class="card responses-response">
    <div class="card-body">
        <h5>Send User Message<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-6 mb-3">
                <label for="responses-{{.Id}}-message">Message</label>
                <input type="text" class="form-control" id="responses-{{.Id}}-message" name="responses[{{.Id}}][message]" required value="{{ .Message }}" >
                <small class="text-muted">Message sent towards the channel</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="responses[{{.Id}}][type]" value="message_user">
            </div>
            <div class="col-md-6 mb-3">
                <label for="responses-{{.Id}}-channelid">UserID</label>
                <input type="text" class="form-control" id="responses-{{.Id}}-channelid" name="responses[{{.Id}}][user]" value="{{ .User }}" >
                <small class="text-muted">Optional value to send the message always to a specific user</small>
            </div>
        </div>
    </div>
</div>
{{end}}