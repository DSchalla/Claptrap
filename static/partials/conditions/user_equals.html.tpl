{{ define "condition-user_equals"}}
<div id="condition-template-user_equals" class="card conditions-condition">
    <div class="card-body conditions-condition">
        <h5>User Equals<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-8 mb-3">
                <label for="conditions-{{.Id}}-condition">Condition</label>
                <input type="text" class="form-control" id="conditions-{{.Id}}-condition" name="conditions[{{ .Id }}][condition]" value="{{ .Condition }}" required>
                <small class="text-muted">Checks if user id or username equals condition</small>
                <div class="invalid-feedback">
                    Condition value is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="user_equals">
            </div>
            <div class="col-md-4 mb-3">
                <label for="conditions-{{.Id}}-parameter">Parameter</label>
                <select id="conditions-{{.Id}}-parameter" name="conditions[{{ .Id }}][parameter]" class="form-control"
                        required>
                    <option value="user" {{ if eq .Parameter "user"}}selected{{end}}>User</option>
                    <option value="actor" {{ if eq .Parameter "actor"}}selected{{end}}>Actor</option>
                </select>
                <small class="text-muted">Check for user or actor role?</small>
                <div class="invalid-feedback">
                    Parameter is required
                </div>
            </div>
        </div>
    </div>
</div>
{{end}}