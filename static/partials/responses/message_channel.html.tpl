{{define "response-message_channel"}}
<div id="response-template-message_channel" class="card responses-response">
    <div class="card-body">
        <h5>Send Channel Message<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-6 mb-3">
                <label for="responses-{{.Id}}-message">Message</label>
                <input type="text" class="form-control" id="responses-{{.Id}}-message" name="responses[{{.Id}}][message]" placeholder="" required value="{{.Message}}">
                <small class="text-muted">Message sent towards the channel</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="responses[{{.Id}}][type]" value="message_channel">
            </div>
            <div class="col-md-6 mb-3">
                <label for="responses-{{.Id}}-channelid">ChannelID</label>
                <input type="text" class="form-control" id="responses-{{.Id}}-channelid" name="responses[{{.Id}}][channel]" value="{{.Channel}}">
                <small class="text-muted">Optional value to send the message always to a specific channel</small>
            </div>
        </div>
    </div>
</div>
{{end}}