{{ define "condition-user_is_role"}}
<div id="condition-template-user_is_role" class="card conditions-condition">
    <div class="card-body conditions-condition">
        <h5>User Is Role<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span>
        </h5>
        <div class="row">
            <div class="col-md-8 mb-3">
                <label for="conditions-{{.Id}}-condition">Condition</label>
                <select id="conditions-{{.Id}}-condition" name="conditions[{{ .Id }}][condition]" class="form-control"
                        required>
                    <option value="admin" {{ if eq .Condition "admin"}}selected{{end}}>Admin (Both System and Team)</option>
                    <option value="system_admin" {{ if eq .Condition "system_admin"}}selected{{end}}>System Admin</option>
                    <option value="team_admin" {{ if eq .Condition "team_admin"}}selected{{end}}>Team Admin</option>
                    <option value="user" {{ if eq .Condition "D"}}selected{{end}}>User</option>
                </select>
                <small class="text-muted">Only triggers if the user triggering the action has the correct type</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="user_is_role">
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