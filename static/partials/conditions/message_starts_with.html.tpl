{{ define "condition-message_starts_with" }}
<div id="condition-template-message_starts_with" class="card conditions-condition">
    <div class="card-body">
        <h5>Message Starts With<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="conditions-{INDEX}-condition">Condition</label>
                <input type="text" class="form-control" id="conditions-{INDEX}-condition" name="conditions[{{ .Id }}][condition]" placeholder="" value="{{.Condition}}" required>
                <small class="text-muted">Value to check within start of the original message</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="message_starts_with">
            </div>
        </div>
    </div>
</div>
{{end}}