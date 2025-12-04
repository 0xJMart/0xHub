{{/*
Expand the name of the chart.
*/}}
{{- define "0xhub.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "0xhub.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "0xhub.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "0xhub.labels" -}}
helm.sh/chart: {{ include "0xhub.chart" . }}
{{ include "0xhub.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "0xhub.selectorLabels" -}}
app.kubernetes.io/name: {{ include "0xhub.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Backend fullname
*/}}
{{- define "0xhub.backend.fullname" -}}
{{- printf "%s-backend" (include "0xhub.fullname" .) }}
{{- end }}

{{/*
Frontend fullname
*/}}
{{- define "0xhub.frontend.fullname" -}}
{{- printf "%s-frontend" (include "0xhub.fullname" .) }}
{{- end }}

{{/*
Operator fullname
*/}}
{{- define "0xhub.operator.fullname" -}}
{{- printf "%s-operator" (include "0xhub.fullname" .) }}
{{- end }}

{{/*
Backend labels
*/}}
{{- define "0xhub.backend.labels" -}}
{{ include "0xhub.labels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "0xhub.backend.selectorLabels" -}}
{{ include "0xhub.selectorLabels" . }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Frontend labels
*/}}
{{- define "0xhub.frontend.labels" -}}
{{ include "0xhub.labels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "0xhub.frontend.selectorLabels" -}}
{{ include "0xhub.selectorLabels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Operator labels
*/}}
{{- define "0xhub.operator.labels" -}}
{{ include "0xhub.labels" . }}
app.kubernetes.io/component: operator
{{- end }}

{{/*
Operator selector labels
*/}}
{{- define "0xhub.operator.selectorLabels" -}}
{{ include "0xhub.selectorLabels" . }}
app.kubernetes.io/component: operator
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "0xhub.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "0xhub.operator.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Namespace
*/}}
{{- define "0xhub.namespace" -}}
{{- if .Values.namespace.create }}
{{- .Values.namespace.name }}
{{- else }}
{{- .Release.Namespace }}
{{- end }}
{{- end }}

