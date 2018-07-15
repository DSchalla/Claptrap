{{ define "condition-message_matches"}}
<div id="condition-template-message_matches" class="card conditions-condition">
    <div class="card-body conditions-condition">
        <h5>Message Matches<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="conditions-{{.Id}}-condition">Condition</label>
                <input type="text" class="form-control" id="conditions-{{.Id}}-condition" name="conditions[{{ .Id }}][condition]" placeholder="" value="{{ .Condition }}" required>
                <small class="text-muted">Checks if message receives matches defined regex (RE2 syntax)</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="message_matches">
            </div>
        </div>
    </div>
</div>
{{end}}