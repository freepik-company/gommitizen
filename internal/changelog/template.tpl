# {{ .Version }} ({{ .Date }})
{{- printf "\n" -}}

{{- range $commonChangeType, $commits := .GroupByChangeType }}
{{- printf "\n" -}}
## {{ $commonChangeType }}
{{- range $commits }}
{{- if eq $commonChangeType "Miscellaneous" }}
- {{ .ChangeType }}{{ if .Scope }}(**{{ .Scope }}**): {{ else }}: {{ end }}{{ .Subject }} (#{{ .ShortHash }})
{{- else }}
- {{ if .Scope }}**{{ .Scope }}**: {{ end }}{{ .Subject }} (#{{ .ShortHash }})
{{- end }}
{{- end }}
{{- printf "\n" -}}
{{- end }}
