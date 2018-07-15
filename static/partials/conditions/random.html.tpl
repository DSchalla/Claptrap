{{ define "condition-random"}}
<div id="condition-template-random" class="card conditions-condition">
    <div class="card-body conditions-condition">
        <h5>Random<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="conditions-{{.Id}}-condition">Condition</label>
                <input type="number" class="form-control" id="conditions-{{.Id}}-likeness" name="conditions[{{ .Id }}][likeness]" placeholder="" value="{{ .Likeness }}" required>
                <small class="text-muted">Percentage of likeness that condition is true</small>
                <div class="invalid-feedback">
                    Likeness is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="random">
            </div>
        </div>
    </div>
</div>
{{end}}