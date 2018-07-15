{{ define "condition-channel_is_type"}}
<div id="condition-template-channel_is_type" class="card conditions-condition">
    <div class="card-body conditions-condition">
        <h5>Channel Is Type<span class="align-items-center badge badge-pill badge-danger dynamic-container-remove">x</span></h5>
        <div class="row">
            <div class="col-md-12 mb-3">
                <label for="conditions-{{.Id}}-condition">Condition</label>
                <select id="conditions-{{.Id}}-condition" name="conditions[{{ .Id }}][condition]" class="form-control" required>
                    <option value="O" {{ if eq .Condition "O"}}selected{{end}}>Public/Open</option>
                    <option value="G" {{ if eq .Condition "G"}}selected{{end}}>Group</option>
                    <option value="P" {{ if eq .Condition "P"}}selected{{end}}>Private</option>
                    <option value="D" {{ if eq .Condition "D"}}selected{{end}}>Direct Message</option>
                </select>
                <small class="text-muted">Only triggers where event took place has the correct type</small>
                <div class="invalid-feedback">
                    Condition is required
                </div>
                <input type="hidden" name="conditions[{{ .Id }}][type]" value="channel_is_type">
            </div>
        </div>
    </div>
</div>
{{end}}