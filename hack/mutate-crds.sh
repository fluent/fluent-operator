#!/usr/bin/env bash

# Support additional annotations for CRDs

CRDS=(
  charts/fluentd-crds/templates/*.yaml
  charts/fluent-bit-crds/templates/*.yaml
)
for CRD in "${CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  awk '{print} /  annotations:/ && !n {print "    {{- with .Values.additionalAnnotations }}\n      {{- toYaml . | nindent 4 }}\n    {{- end }}"; n++}' "$CRD" > "$CRD.new" && mv "$CRD.new" "$CRD"
done
