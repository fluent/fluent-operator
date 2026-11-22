#!/usr/bin/env bash
# Syncs the charts from a fluent-operator release tag into a fluent/helm-charts checkout.
# Bumps chart versions, adds CHANGELOG entries and regenerates READMEs. Git and PR steps are left to the caller.
#
# Usage: hack/sync-helm-charts.sh <tag> <helm-charts-dir>

set -o errexit
set -o nounset
set -o pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <tag> <helm-charts-dir>" >&2
  exit 1
fi

TAG=$1
HELM_CHARTS_DIR=$(cd "$2" && pwd -P)
FLUENT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
REPO_URL=https://github.com/fluent/fluent-operator
CHARTS=(fluent-operator fluent-operator-fluent-bit-crds fluent-operator-fluentd-crds)
DATE=$(date -u +%Y-%m-%d)

semver_re='^v?([0-9]+)\.([0-9]+)\.([0-9]+)$'

chart_field() {
  awk -v key="$2:" '$1 == key { print $2; exit }' "$1"
}

set_chart_field() {
  local file=$1 key=$2 value=$3
  awk -v key="$key:" -v value="$value" '$1 == key && !done { print key " " value; done = 1; next } { print }' \
    "$file" > "$file.tmp"
  mv "$file.tmp" "$file"
}

# Prints the bump level (major, minor or patch) between two versions.
bump_level() {
  [[ $1 =~ $semver_re ]] && local a=("${BASH_REMATCH[@]:1}")
  [[ $2 =~ $semver_re ]] && local b=("${BASH_REMATCH[@]:1}")
  if [[ ${b[0]} -gt ${a[0]} ]]; then
    echo major
  elif [[ ${b[0]} -eq ${a[0]} && ${b[1]} -gt ${a[1]} ]]; then
    echo minor
  elif [[ ${b[0]} -eq ${a[0]} && ${b[1]} -eq ${a[1]} && ${b[2]} -gt ${a[2]} ]]; then
    echo patch
  fi
}

bump_version() {
  [[ $1 =~ $semver_re ]]
  local major=${BASH_REMATCH[1]} minor=${BASH_REMATCH[2]} patch=${BASH_REMATCH[3]}
  case $2 in
    major) echo "$((major + 1)).0.0" ;;
    minor) echo "${major}.$((minor + 1)).0" ;;
    patch) echo "${major}.${minor}.$((patch + 1))" ;;
  esac
}

# Prints a changelog section for conventional commits between two tags that touched a path.
# Section is one of: Added, Fixed, Changed.
changelog_items() {
  local range=$1 path=$2 section=$3
  git -C "$FLUENT_ROOT" log --no-merges --format='%h%x09%s' "$range" -- "$path" | while IFS=$'\t' read -r sha subject; do
    [[ $subject =~ ^([a-z]+)(\([^\)]*\))?(!)?:\ (.*)$ ]] || continue
    local type=${BASH_REMATCH[1]} breaking=${BASH_REMATCH[3]} desc=${BASH_REMATCH[4]} want
    case $type in
      feat | feature) want=Added ;;
      fix) want=Fixed ;;
      build | deps | perf | revert) want=Changed ;;
      *) continue ;;
    esac
    [[ $want == "$section" ]] || continue

    local link="[${sha}](${REPO_URL}/commit/${sha})"
    if [[ $desc =~ ^(.*)\ \(#([0-9]+)\)$ ]]; then
      desc=${BASH_REMATCH[1]}
      link="[#${BASH_REMATCH[2]}](${REPO_URL}/pull/${BASH_REMATCH[2]})"
    fi
    desc="$(echo "${desc:0:1}" | tr '[:lower:]' '[:upper:]')${desc:1}"
    [[ -n $breaking ]] && desc="**BREAKING:** ${desc}"
    echo "- ${desc} (${link})"
  done
}

# Prints the full CHANGELOG entry for a chart.
changelog_entry() {
  local chart=$1 version=$2 range=$3 app_version=$4 section items printed=false
  echo "## [v${version}] - ${DATE}"
  echo
  for section in Added Fixed Changed; do
    items=$(changelog_items "$range" "charts/${chart}" "$section")
    if [[ $chart == fluent-operator && $section == Changed ]]; then
      items=$(printf '%s\n- Bumped fluent-operator to v%s' "$items" "$app_version" | sed '/^$/d')
    fi
    [[ -n $items ]] || continue
    printf '### %s\n\n%s\n\n' "$section" "$items"
    printed=true
  done
  if [[ $printed == false ]]; then
    printf '### Changed\n\n- No chart changes; version aligned with fluent-operator v%s\n\n' "$app_version"
  fi
}

update_changelog() {
  local file=$1 entry_file=$2 link=$3
  # Fail rather than guess where manually added unreleased entries belong.
  if awk '/^## \[UNRELEASED\]/ { f = 1; next } f && /^## / { exit } f && NF { found = 1 } END { exit !found }' "$file"; then
    echo "error: ${file} has entries under [UNRELEASED]; release them manually" >&2
    exit 1
  fi
  awk -v entry_file="$entry_file" -v link="$link" '
    /^## \[UNRELEASED\]/ {
      print; print ""
      while ((getline line < entry_file) > 0) print line
      skip = 1; next
    }
    skip && !NF { next }
    { skip = 0 }
    /^\[v[0-9]/ && !linked { print link; linked = 1 }
    { print }
  ' "$file" > "$file.tmp"
  mv "$file.tmp" "$file"
}

if [[ ! $TAG =~ $semver_re ]]; then
  echo "error: ${TAG} is not a release tag (expected vX.Y.Z)" >&2
  exit 1
fi
APP_VERSION=${TAG#v}

PREV_APP_VERSION=$(chart_field "${HELM_CHARTS_DIR}/charts/fluent-operator/Chart.yaml" appVersion)
PREV_TAG=v${PREV_APP_VERSION}
LEVEL=$(bump_level "$PREV_APP_VERSION" "$APP_VERSION")
if [[ -z $LEVEL ]]; then
  echo "error: ${TAG} is not newer than the published appVersion ${PREV_APP_VERSION}" >&2
  exit 1
fi
git -C "$FLUENT_ROOT" rev-parse --verify --quiet "refs/tags/${PREV_TAG}" >/dev/null || {
  echo "error: previous tag ${PREV_TAG} not found; fetch tags first" >&2
  exit 1
}

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT
git -C "$FLUENT_ROOT" archive "$TAG" charts | tar -x -C "$tmpdir"

for chart in "${CHARTS[@]}"; do
  dest=${HELM_CHARTS_DIR}/charts/${chart}
  chart_yaml=${dest}/Chart.yaml
  version=$(bump_version "$(chart_field "$chart_yaml" version)" "$LEVEL")

  # Files maintained only in helm-charts are excluded, which also protects them from --delete.
  rsync -a --delete --exclude Chart.yaml --exclude CHANGELOG.md --exclude README.md --exclude _description.md \
    "${tmpdir}/charts/${chart}/" "${dest}/"

  set_chart_field "$chart_yaml" version "$version"
  set_chart_field "$chart_yaml" appVersion "$APP_VERSION"

  changelog_entry "$chart" "$version" "${PREV_TAG}..${TAG}" "$APP_VERSION" > "${tmpdir}/entry.md"
  update_changelog "${dest}/CHANGELOG.md" "${tmpdir}/entry.md" \
    "[v${version}]: https://github.com/fluent/helm-charts/releases/tag/${chart}-${version}"

  echo "${chart}: ${version} (appVersion ${APP_VERSION})"
done

(
  cd "$HELM_CHARTS_DIR"
  helm-docs --template-files=./_templates.gotmpl --template-files=./_chart-readme.md.gotmpl \
    --chart-to-generate "$(printf 'charts/%s,' "${CHARTS[@]}" | sed 's/,$//')"
)
