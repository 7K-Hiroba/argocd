{{- define "argocd.fullname" -}}
{{ .Release.Name }}
{{- end }}
{{- define "argocd.labels" -}}
app.kubernetes.io/name: argocd
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
