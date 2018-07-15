{{ define "condition-message_contains"}}
<div id="condition-template-message_contains" class="card conditions-condition">
    <div class="card-body conditions-condition">
        <h5>Message Contains<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="conditions-{INDEX}-condition">Condition</label>
                <input type="text" class="form-control" id="conditions-{INDEX}-condition" name="conditions[{{ .Id }}][condition]" placeholder="" value="{{ .Condition }}" required>
                <small class="text-muted">Value to check inside the original message</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="message_contains">
            </div>
        </div>
    </div>
</div>
{{end}}