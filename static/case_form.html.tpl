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
        {{range $key, $value := .Data.ConditionOptions}}
            <option value="{{ $key }}" {{if eq $key "message_contains"}}selected{{end}}>{{ $value }}</option>
        {{end}}
        </select>
        <br/>
        <button id="condition-select-add" class="btn btn-block" type="submit">Add</button>
        <hr/>
        <h5>Responses</h5>
        <select id="response-select-dropdown" name="type" class="form-control">
        {{range $key, $value := .Data.ResponseOptions}}
            <option value="{{ $key }}" {{if eq $key "message_channel"}}selected{{end}}>{{ $value }}</option>
        {{end}}
        </select>
        <br/>
        <button id="response-select-add" class="btn btn-block" type="submit">Add</button>
    </div>
    <div class="col-md-8 order-md-1">
        <h4 class="mb-3">General</h4>
        <form class="needs-validation" novalidate method="post">
            <div class="row">
                <div class="col-md-6 mb-3">
                    <label for="firstName">Case Name</label>
                    <input type="text" class="form-control" id="casename" name="casename" placeholder="" value="{{.Data.Case.Name}}"
                           required>
                    <div class="invalid-feedback">
                        Valid name is required.
                    </div>
                </div>
                <div class="form-group col-md-3">
                    <label for="form-type">Type</label>
                    <select id="form-type" name="type" class="form-control">
                    {{range $key, $value := .Data.CaseTypes}}
                        <option value="{{ $key }}" {{if eq $.Data.CaseType $key}}selected{{end}}>{{ $value }}</option>
                    {{end}}
                    </select>
                </div>
                <div class="form-group col-md-3">
                    <label for="form-type">Condition Matching</label>
                    <select id="form-type" name="type" class="form-control">
                        <option {{if eq .Data.Case.ConditionMatching "or"}}selected{{end}}>or</option>
                        <option {{if eq .Data.Case.ConditionMatching "and"}}selected{{end}}>and</option>
                    </select>
                    <small class="text-muted">If set to 'or', a single condition will trigger the responses. Otherwise
                        all conditions have to match.
                    </small>
                </div>
            </div>
            <hr class="mb-4">
            <h4 class="mb-3">Conditions</h4>
            <div id="conditions-container" class="dynamic-container">
                {{range .Data.Case.Conditions }}
                {{.CondType}}
                    {{if eq .CondType "channel_equals"}}
                    {{ template "condition-channel_equals" .}}
                    {{end}}
                    {{if eq .CondType "channel_is_type"}}
                    {{ template "condition-channel_is_type" .}}
                    {{end}}
                    {{if eq .CondType "message_contains"}}
                    {{ template "condition-message_contains" .}}
                    {{end}}
                    {{if eq .CondType "message_equals"}}
                    {{ template "condition-message_equals" .}}
                    {{end}}
                    {{if eq .CondType "message_matches"}}
                    {{ template "condition-message_matches" .}}
                    {{end}}
                    {{if eq .CondType "message_starts_with"}}
                    {{ template "condition-message_starts_with" .}}
                    {{end}}
                    {{if eq .CondType "random"}}
                    {{ template "condition-random" .}}
                    {{end}}
                    {{if eq .CondType "user_equals"}}
                    {{ template "condition-user_equals" .}}
                    {{end}}
                    {{if eq .CondType "user_is_role"}}
                    {{ template "condition-user_is_role" .}}
                    {{end}}
                {{end}}
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
                {{range .Data.Case.Responses }}
                    {{ if eq .Action "message_channel"}}
                    {{ template "response-message_channel" . }}
                    {{end}}
                    {{ if eq .Action "message_user"}}
                    {{ template "response-message_user" . }}
                    {{end}}
                    {{ if eq .Action "message_ephemeral"}}
                    {{ template "response-message_ephemeral" . }}
                    {{end}}
                    {{ if eq .Action "delete_message"}}
                    {{ template "response-delete_message" . }}
                    {{end}}
                    {{ if eq .Action "invite_user"}}
                    {{ template "response-invite_user" . }}
                    {{end}}
                    {{ if eq .Action "kick_user"}}
                    {{ template "response-kick_user" . }}
                    {{end}}
                    {{ if eq .Action "message_user"}}
                    {{ template "response-message_user" . }}
                    {{end}}
                {{end}}
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
    {{ template "condition-channel_equals" .Data.DummyStruct}}
    {{ template "condition-channel_is_type" .Data.DummyStruct}}
    {{ template "condition-message_contains" .Data.DummyStruct}}
    {{ template "condition-message_equals" .Data.DummyStruct}}
    {{ template "condition-message_matches" .Data.DummyStruct}}
    {{ template "condition-message_starts_with" .Data.DummyStruct}}
    {{ template "condition-random" .Data.DummyStruct}}
    {{ template "condition-user_equals" .Data.DummyStruct}}
    {{ template "condition-user_is_role" .Data.DummyStruct}}
</div>
<div id="responses-template-container">
    {{ template "response-message_channel" .Data.DummyStruct }}
    {{ template "response-message_user" .Data.DummyStruct }}
    {{ template "response-message_ephemeral" .Data.DummyStruct }}
    {{ template "response-delete_message" .Data.DummyStruct }}
    {{ template "response-invite_user" .Data.DummyStruct }}
    {{ template "response-kick_user" .Data.DummyStruct }}
    {{ template "response-intercept_event" .Data.DummyStruct }}
</div>
{{end}}
