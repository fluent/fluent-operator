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

# Function to add additionalAnnotations Helm templating to a CRD.
# Handles three cases robustly in a single awk pass:
#
#   Case 1 — annotations: key is present but empty (current controller-gen output
#             after strip_controller_gen_annotation):
#             Wrap the entire field in {{- with }} so an empty additionalAnnotations
#             value produces no annotations key, avoiding a null field that fails
#             kubeconform / Kubernetes OpenAPI validation.
#
#   Case 2 — annotations: key is present with existing entries:
#             Keep the field unconditional (it already has real content) and merge
#             additionalAnnotations inside the block, preserving existing entries.
#
#   Case 3 — annotations: key is absent entirely (future-proofing for controller-gen
#             versions that may omit the field):
#             Insert the conditional block under metadata: before the name: field.
#
# The function is idempotent — it is safe to run multiple times.
add_annotations() {
  local CRD="$1"
  if grep -q "{{- with .Values.additionalAnnotations }}" "$CRD"; then
    return
  fi

  awk '
    /^metadata:/ { in_meta = 1 }
    /^spec:/     { in_meta = 0 }

    in_meta && /^  annotations:/ && !done {
      done = 1
      # Read ahead to collect any existing annotation entries (4-space-indented lines).
      existing = ""
      while ((getline nxt) > 0) {
        if (nxt ~ /^    [^ ]/) { existing = existing nxt "\n"; continue }
        break
      }
      if (existing == "") {
        # Case 1: empty annotations block — make the entire field conditional.
        print "  {{- with .Values.additionalAnnotations }}"
        print "  annotations:"
        print "    {{- toYaml . | nindent 4 }}"
        print "  {{- end }}"
      } else {
        # Case 2: existing entries — keep field unconditional, merge inside.
        print "  annotations:"
        printf "%s", existing
        print "    {{- with .Values.additionalAnnotations }}"
        print "    {{- toYaml . | nindent 4 }}"
        print "    {{- end }}"
      }
      print nxt  # first line after the annotations block (already read by getline)
      next
    }

    # Case 3: no annotations: key seen under metadata — insert before name:.
    in_meta && !done && /^  name:/ {
      print "  {{- with .Values.additionalAnnotations }}"
      print "  annotations:"
      print "    {{- toYaml . | nindent 4 }}"
      print "  {{- end }}"
      done = 1
    }

    { print }
  ' "$CRD" > "$CRD.new" && mv "$CRD.new" "$CRD"
}

# Handle fluent-operator chart - bundled CRDs (crds/ directory)
OPERATOR_CRDS=(charts/fluent-operator/crds/*.yaml)
for CRD in "${OPERATOR_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  strip_doc_separator "$CRD"
done

# Handle fluent-operator-fluent-bit-crds chart - Fluent Bit CRDs
FLUENT_BIT_CRDS=(charts/fluent-operator-fluent-bit-crds/templates/*.yaml)
for CRD in "${FLUENT_BIT_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  [[ "$(basename "$CRD")" == ".gitkeep" ]] && continue

  strip_doc_separator "$CRD"
  strip_controller_gen_annotation "$CRD"
  add_annotations "$CRD"
done

# Handle fluent-operator-fluentd-crds chart - Fluentd CRDs
FLUENTD_CRDS=(charts/fluent-operator-fluentd-crds/templates/*.yaml)
for CRD in "${FLUENTD_CRDS[@]}"
do
  [[ -f "$CRD" ]] || continue
  [[ "$(basename "$CRD")" == ".gitkeep" ]] && continue

  strip_doc_separator "$CRD"
  strip_controller_gen_annotation "$CRD"
  add_annotations "$CRD"
done
