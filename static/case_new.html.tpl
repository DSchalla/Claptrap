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
        <select id="form-type" name="type" class="form-control">
            <option value="message" selected>Message Contains</option>
            <option value="channel_join">Message Starts With</option>
            <option value="channel_leave">Message Matches</option>
            <option value="channel_leave">Message Equals</option>
        </select>
        <br/>
        <button class="btn btn-block" type="submit">Add</button>
        <hr/>
        <h5>Responses</h5>
        <select id="form-type" name="type" class="form-control">
            <option value="message" selected>Send Channel Message</option>
            <option value="channel_join">Send Direct Message</option>
            <option value="channel_leave">Send Ephemeral Message</option>
            <option value="channel_leave">Message Equals</option>
        </select>
        <br/>
        <button class="btn btn-block" type="submit">Add</button>
    </div>
    <div class="col-md-8 order-md-1">
        <h4 class="mb-3">General</h4>
        <form class="needs-validation" novalidate>
            <div class="row">
                <div class="col-md-12 mb-3">
                    <label for="firstName">Case Name</label>
                    <input type="text" class="form-control" id="casename" placeholder="" value="" required>
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
                        <option>And</option>
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
            <div class="conditions-container">

                <div class="card dynamic-container">
                    <div class="card-body">
                        <h5>Message Contains<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
                        <div class="row">
                            <div class="col-md-12 mb-3">
                                <label for="cc-name">Condition</label>
                                <input type="text" class="form-control" id="cc-name" placeholder="" required>
                                <small class="text-muted">Value to check inside the original message</small>
                                <div class="invalid-feedback">
                                    Condition is required
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card dynamic-container">
                    <div class="card-body">
                        <h5>User Equals<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
                        <div class="row">
                            <div class="col-md-12 mb-3">
                                <label for="cc-name">Condition</label>
                                <input type="text" class="form-control" id="cc-name" placeholder="" required>
                                <small class="text-muted">Condition can be either the UserID or Username</small>
                                <div class="invalid-feedback">
                                    Condition is required
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="card conditions-nocondition">
                    <div class="card-body">
                        <h5 class="card-title">No conditions</h5>
                        <p class="card-text">You didn't add any conditions yet - Add some using the right panel</p>
                    </div>
                </div>
            </div>
            <hr class="mb-4">

            <h4 class="mb-3">Responses</h4>

            <h5>Message Channel</h5>
            <div class="row">
                <div class="col-md-6 mb-3">
                    <label for="cc-name">Message</label>
                    <input type="text" class="form-control" id="cc-name" placeholder="" required>
                    <small class="text-muted">Message sent towards the channel</small>
                    <div class="invalid-feedback">
                        Condition is required
                    </div>
                </div>
                <div class="col-md-6 mb-3">
                    <label for="cc-name">ChannelID</label>
                    <input type="text" class="form-control" id="cc-name" placeholder="">
                    <small class="text-muted">Optional value to send the message always to a specific channel</small>
                </div>
            </div>
            <hr class="mb-4">
            <button class="btn btn-primary btn-lg btn-block" type="submit">Save</button>
        </form>
    </div>
</div>
<br/>
<br/>
{{end}}
