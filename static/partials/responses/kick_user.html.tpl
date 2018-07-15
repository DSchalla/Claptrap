{{define "response-kick_user"}}
<div id="response-template-kick_user" class="card responses-response">
    <div class="card-body">
        <h5>Kick User<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="responses-{{.Id}}-channelid">ChannelID</label>
                <input type="text" class="form-control" id="responses-{{.Id}}-channelid" name="responses[{{.Id}}][channel]" value="{{.Channel}}">
                <small class="text-muted">Optional channel the user should be kicked from (Default: channel event got triggered from)</small>
                <input type="hidden" name="responses[{{.Id}}][type]" value="invite_user">
            </div>
        </div>
    </div>
</div>
{{end}}