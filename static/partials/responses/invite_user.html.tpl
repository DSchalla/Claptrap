{{define "response-invite_user"}}
<div id="response-template-invite_user" class="card responses-response">
    <div class="card-body">
        <h5>Invite User<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-6 mb-3">
                <label for="responses-{{.Id}}-channelid">ChannelID</label>
                <input type="text" class="form-control" id="responses-{{.Id}}-channelid" name="responses[{{.Id}}][channel]" value="{{ .Channel }}" >
                <small class="text-muted">Channel the user should be invited to</small>
                <input type="hidden" name="responses[{{.Id}}][type]" value="invite_user">
            </div>
        </div>
    </div>
</div>
{{end}}