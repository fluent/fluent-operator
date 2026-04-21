#!/usr/bin/env bash

# Support additional annotations for CRDs
# This script is idempotent - it can be run multiple times safely

# Cross-platform sed in-place edit (BSD sed requires '' after -i; GNU sed does not)
sedi() {
  if sed --version 2>/dev/null | grep -q GNU; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

# Function to strip the leading YAML document separator added by controller-gen
strip_doc_separator() {
  local CRD="$1" 
  sedi '/^---$/d' "$CRD"
}

# Function to strip the controller-gen version annotation added by controller-gen
strip_controller_gen_annotation() {
  local CRD="$1"
  sedi '/controller-gen\.kubebuilder\.io\/version:/d' "$CRD"
}

# Function to add annotations templating to a CRD if not already present.
# Wraps the entire annotations field in {{- with }} so that an empty
# additionalAnnotations value produces no annotations key at all, avoiding
# a null annotations field that fails kubeconform / Kubernetes OpenAPI validation.
add_annotations() {
  local CRD="$1"
  # Check if additionalAnnotations templating already exists (idempotent)
  if ! grep -q "{{- with .Values.additionalAnnotations }}" "$CRD"; then
    awk '/  annotations:/ && !n {
      print "  {{- with .Values.additionalAnnotations }}"
      print "  annotations:"
      print "    {{- toYaml . | nindent 4 }}"
      print "  {{- end }}"
      n++; next
    } {print}' "$CRD" > "$CRD.new" && mv "$CRD.new" "$CRD"
  fi
}

# Handle fluent-operator chart - bundled CRDs (crds/ directory)
OPERATOR_CRDS=(charts/fluent-operator/crds/*.yaml)
for CRD in "${OPERATOR_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  strip_doc_separator "$CRD"
done

# Handle fluent-operator-crds-fluent-bit chart - Fluent Bit CRDs
FLUENT_BIT_CRDS=(charts/fluent-operator-crds-fluent-bit/templates/*.yaml)
for CRD in "${FLUENT_BIT_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  [[ "$(basename "$CRD")" == ".gitkeep" ]] && continue

  strip_doc_separator "$CRD"
  strip_controller_gen_annotation "$CRD"
  add_annotations "$CRD"
done

# Handle fluent-operator-crds-fluentd chart - Fluentd CRDs
FLUENTD_CRDS=(charts/fluent-operator-crds-fluentd/templates/*.yaml)
for CRD in "${FLUENTD_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  [[ "$(basename "$CRD")" == ".gitkeep" ]] && continue

  strip_doc_separator "$CRD"
  strip_controller_gen_annotation "$CRD"
  add_annotations "$CRD"
done
