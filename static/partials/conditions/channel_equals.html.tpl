{{ define "condition-channel_equals"}}
<div id="condition-template-channel_equals" class="card conditions-condition">
    <div class="card-body conditions-condition">
        <h5>Channel Equals<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="conditions-{{.Id}}-condition">Condition</label>
                <input type="text" class="form-control" id="conditions-{{.Id}}-condition" name="conditions[{{ .Id }}][condition]" placeholder="" value="{{ .Condition }}" required>
                <small class="text-muted">Checks if Channel ID or Channel Name equals Condition</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="channel_equals">
            </div>
        </div>
    </div>
</div>
{{end}}