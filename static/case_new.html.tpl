{{define "content"}}
<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">New Case</h1>
</div>
<div class="row">
    <div class="col-md-4 order-md-2 mb-4">
        <h4 class="d-flex justify-content-between align-items-center mb-3">
            <span class="text-muted">Add Condition / Response</span>
        </h4>
        <h5>Conditions</h5>
        <select id="condition-select-dropdown" name="type" class="form-control">
            <option value="message_contains" selected>Message Contains</option>
            <option value="message_starts_with">Message Starts With</option>
            <option value="message_matches">Message Matches</option>
            <option value="message_equals">Message Equals</option>
        </select>
        <br/>
        <button id="condition-select-add" class="btn btn-block" type="submit">Add</button>
        <hr/>
        <h5>Responses</h5>
        <select id="response-select-dropdown" name="type" class="form-control">
            <option value="message_channel" selected>Send Channel Message</option>
            <option value="channel_join">Send Direct Message</option>
            <option value="channel_leave">Send Ephemeral Message</option>
            <option value="channel_leave">Message Equals</option>
        </select>
        <br/>
        <button id="response-select-add" class="btn btn-block" type="submit">Add</button>
    </div>
    <div class="col-md-8 order-md-1">
        <h4 class="mb-3">General</h4>
        <form class="needs-validation" novalidate method="post">
            <div class="row">
                <div class="col-md-12 mb-3">
                    <label for="firstName">Case Name</label>
                    <input type="text" class="form-control" id="casename" name="casename" placeholder="" value="" required>
                    <div class="invalid-feedback">
                        Valid name is required.
                    </div>
                </div>
            </div>
            <div class="row">
                <div class="form-group col-md-4">
                    <label for="form-type">Type</label>
                    <select id="form-type" name="type" class="form-control">
                        <option value="message" selected>Message</option>
                        <option value="channel_join">Channel Join</option>
                        <option value="channel_leave">Channel Leave</option>
                    </select>
                </div>
                <div class="form-group col-md-4">
                    <label for="form-type">Condition Matching</label>
                    <select id="form-type" name="type" class="form-control">
                        <option selected>or</option>
                        <option>and</option>
                    </select>
                    <small class="text-muted">If set to 'or', a single condition will trigger the responses. Otherwise all conditions have to match.</small>
                </div>
                <div class="form-group col-md-4">
                    <label for="form-intercept">Interception</label>
                    <select id="form-intercept" name="intercept" class="form-control">
                        <option selected>No</option>
                        <option>Yes</option>
                    </select>
                    <small class="text-muted">If turned on, original message is rejected if all conditions matched</small>
                </div>
            </div>
            <hr class="mb-4">

            <h4 class="mb-3">Conditions</h4>
            <div id="conditions-container" class="dynamic-container">
                <div class="card conditions-nocondition">
                    <div class="card-body">
                        <h5 class="card-title">No conditions</h5>
                        <p class="card-text">You didn't add any conditions yet - Add some using the right panel</p>
                    </div>
                </div>
            </div>
            <hr class="mb-4">

            <h4 class="mb-3">Responses</h4>
            <div id="responses-container" class="dynamic-container">
                <div class="card responses-noresponse">
                    <div class="card-body">
                        <h5 class="card-title">No responses</h5>
                        <p class="card-text">You didn't add any response yet - Add some using the right panel</p>
                    </div>
                </div>
            </div>
            <hr class="mb-4">
            <button class="btn btn-primary btn-lg btn-block" type="submit">Save</button>
        </form>
    </div>
</div>
<br/>
<br/>
<div id="conditions-template-container">
    <div id="condition-template-message_contains" class="card conditions-condition">
        <div class="card-body conditions-condition">
            <h5>Message Contains<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
            <div class="row">
                <div class="col-md-12 mb-3">
                    <label for="conditions-{INDEX}-condition">Condition</label>
                    <input type="text" class="form-control" id="conditions-{INDEX}-condition" name="conditions[{INDEX}][condition]" placeholder="" required>
                    <small class="text-muted">Value to check inside the original message</small>
                    <div class="invalid-feedback">
                        Condition is required
                    </div>
                    <input type="hidden" name="conditions[{INDEX}][type]" value="message_contains">
                </div>
            </div>
        </div>
    </div>
    <div id="condition-template-message_starts_with" class="card conditions-condition">
        <div class="card-body">
            <h5>Message Starts With<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
            <div class="row">
                <div class="col-md-12 mb-3">
                    <label for="conditions-{INDEX}-condition">Condition</label>
                    <input type="text" class="form-control" id="conditions-{INDEX}-condition" name="conditions[{INDEX}][condition]" placeholder="" required>
                    <small class="text-muted">Value to check within start of the original message</small>
                    <div class="invalid-feedback">
                        Condition is required
                    </div>
                    <input type="hidden" name="conditions[{INDEX}][type]" value="message_starts_with">
                </div>
            </div>
        </div>
    </div>
</div>
<div id="responses-template-container">
    <div id="response-template-message_channel" class="card responses-response">
        <div class="card-body">
            <h5>Send Channel Message<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
            <div class="row">
                <div class="col-md-6 mb-3">
                    <label for="responses-{INDEX}-message">Message</label>
                    <input type="text" class="form-control" id="responses-{INDEX}-message" placeholder="" required>
                    <small class="text-muted">Message sent towards the channel</small>
                    <div class="invalid-feedback">
                        Condition is required
                    </div>
                </div>
                <div class="col-md-6 mb-3">
                    <label for="responses-{INDEX}-channelid">ChannelID</label>
                    <input type="text" class="form-control" id="responses-{INDEX}-channelid" placeholder="">
                    <small class="text-muted">Optional value to send the message always to a specific channel</small>
                </div>
                <input type="hidden" name="responses[{INDEX}][type]" value="message_channel">
            </div>
        </div>
    </div>
</div>
{{end}}
