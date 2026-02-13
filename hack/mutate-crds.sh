#!/usr/bin/env bash

# Support additional annotations for CRDs
# This script is idempotent - it can be run multiple times safely

# Function to add annotations templating to a CRD if not already present
add_annotations() {
  local CRD="$1"
  # Check if additionalAnnotations templating already exists
  if ! grep -q "{{- with .Values.additionalAnnotations }}" "$CRD"; then
    awk '{print} /  annotations:/ && !n {print "    {{- with .Values.additionalAnnotations }}\n      {{- toYaml . | nindent 4 }}\n    {{- end }}"; n++}' "$CRD" > "$CRD.new" && mv "$CRD.new" "$CRD"
  fi
}

# Function to wrap CRD with conditional if not already present
wrap_conditional() {
  local CRD="$1"
  local CONDITION="$2"

  # Check if conditional already exists
  if ! head -1 "$CRD" | grep -q "{{- if"; then
    {
      echo "$CONDITION"
      cat "$CRD"
      echo "{{- end }}"
    } > "$CRD.new" && mv "$CRD.new" "$CRD"
  fi
}

# Handle fluent-operator-crds chart - Fluent Bit CRDs
FLUENT_BIT_CRDS=(charts/fluent-operator-crds/templates/fluent-bit/*.yaml)
for CRD in "${FLUENT_BIT_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  [[ "$(basename "$CRD")" == ".gitkeep" ]] && continue

  # Add annotations first (before conditional wrapper)
  add_annotations "$CRD"

  # Wrap with conditional
  wrap_conditional "$CRD" "{{- if .Values.fluent-bit.enabled }}"
done

# Handle fluent-operator-crds chart - Fluentd CRDs
FLUENTD_CRDS=(charts/fluent-operator-crds/templates/fluentd/*.yaml)
for CRD in "${FLUENTD_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  [[ "$(basename "$CRD")" == ".gitkeep" ]] && continue

  # Add annotations first (before conditional wrapper)
  add_annotations "$CRD"

  # Wrap with conditional
  wrap_conditional "$CRD" "{{- if .Values.fluentd.enabled }}"
done
