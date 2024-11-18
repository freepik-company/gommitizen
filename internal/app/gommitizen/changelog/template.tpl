# {{ .Version }} ({{ .Date }})
{{- printf "\n" -}}

{{- if .BreakingChanges }}
## Breaking changes
{{- range .BreakingChanges }}
- {{ if .Scope }}**{{ .Scope }}**: {{ end }}{{ .Subject }} (#{{ .ShortHash }})
{{- end }}
{{- printf "\n" -}}
{{- end }}

{{- if .Features }}
## Features
{{- range .Features }}
- {{ if .Scope }}**{{ .Scope }}**: {{ end }}{{ .Subject }} (#{{ .ShortHash }})
{{- end }}
{{- printf "\n" -}}
{{- end }}

{{- if .BugFixes }}
## Fixes
{{- range .BugFixes }}
- {{ if .Scope }}**{{ .Scope }}**: {{ end }}{{ .Subject }} (#{{ .ShortHash }})
{{- end }}
{{- printf "\n" -}}
{{- end }}

{{- if .Refactors }}
## Refactors
{{- range .Refactors }}
- {{ if .Scope }}**{{ .Scope }}**: {{ end }}{{ .Subject }} (#{{ .ShortHash }})
{{- end }}
{{- printf "\n" -}}
{{- end }}

{{- if .Miscellaneous }}
## Miscellaneous
{{- range .Miscellaneous }}
- {{ .ChangeType }}{{ if .Scope }}(**{{ .Scope }}**): {{ else }}: {{ end }}{{ .Subject }} (#{{ .ShortHash }})
{{- end }}
{{- printf "\n" -}}
{{- end }}
