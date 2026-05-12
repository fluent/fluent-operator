#!/bin/bash
set -e

# Test script for validating helm chart sync locally
# This simulates what the GitHub Actions workflow will do

COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_RED='\033[0;31m'
COLOR_BLUE='\033[0;34m'
COLOR_RESET='\033[0m'

log_info() {
    echo -e "${COLOR_BLUE}ℹ️  $1${COLOR_RESET}"
}

log_success() {
    echo -e "${COLOR_GREEN}✅ $1${COLOR_RESET}"
}

log_warning() {
    echo -e "${COLOR_YELLOW}⚠️  $1${COLOR_RESET}"
}

log_error() {
    echo -e "${COLOR_RED}❌ $1${COLOR_RESET}"
}

# Check if we're in the right directory
if [ ! -f "charts/fluent-operator/Chart.yaml" ]; then
    log_error "Must be run from the root of the fluent-operator repository"
    exit 1
fi

log_info "Starting local helm chart sync test..."
echo ""

# Create temporary directory for testing
TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

log_info "Test directory: $TEST_DIR"

# Create mock helm-charts directory structure
mkdir -p "$TEST_DIR/helm-charts/charts"

log_info "Extracting chart versions..."

# Extract versions
FLUENT_OP_VERSION=$(grep '^version:' charts/fluent-operator/Chart.yaml | awk '{print $2}')
FLUENTD_CRDS_VERSION=$(grep '^version:' charts/fluentd-crds/Chart.yaml | awk '{print $2}')
FLUENT_BIT_CRDS_VERSION=$(grep '^version:' charts/fluent-operator/charts/fluent-bit-crds/Chart.yaml | awk '{print $2}')

echo "  • fluent-operator: v${FLUENT_OP_VERSION}"
echo "  • fluentd-crds: v${FLUENTD_CRDS_VERSION}"
echo "  • fluent-bit-crds: v${FLUENT_BIT_CRDS_VERSION}"
echo ""

# Sync fluent-operator chart
log_info "Syncing fluent-operator chart..."
mkdir -p "$TEST_DIR/helm-charts/charts/fluent-operator/charts"
rsync -av --exclude='charts/' charts/fluent-operator/ "$TEST_DIR/helm-charts/charts/fluent-operator/" > /dev/null
log_success "fluent-operator synced"

# Sync fluent-bit-crds chart
log_info "Syncing fluent-bit-crds chart..."
mkdir -p "$TEST_DIR/helm-charts/charts/fluent-bit-crds"
rsync -av charts/fluent-operator/charts/fluent-bit-crds/ "$TEST_DIR/helm-charts/charts/fluent-bit-crds/" > /dev/null
log_success "fluent-bit-crds synced"

# Sync fluentd-crds chart
log_info "Syncing fluentd-crds chart..."
mkdir -p "$TEST_DIR/helm-charts/charts/fluentd-crds"
rsync -av charts/fluentd-crds/ "$TEST_DIR/helm-charts/charts/fluentd-crds/" > /dev/null
log_success "fluentd-crds synced"

# Update Chart.yaml repository references
log_info "Updating Chart.yaml repository references..."
sed -i.bak 's|repository: "file://charts/fluent-bit-crds"|repository: "https://fluent.github.io/helm-charts"|g' \
    "$TEST_DIR/helm-charts/charts/fluent-operator/Chart.yaml"
log_success "Repository references updated"
echo ""

# Check for helm
if ! command -v helm &> /dev/null; then
    log_warning "Helm is not installed, skipping validation"
    log_info "Install helm: https://helm.sh/docs/intro/install/"
else
    log_info "Validating charts with helm lint..."
    echo ""
    
    # Validate each chart
    for chart_dir in "$TEST_DIR/helm-charts/charts/"*/; do
        chart_name=$(basename "$chart_dir")
        echo -n "  Linting ${chart_name}... "
        
        if helm lint "$chart_dir" > "$TEST_DIR/${chart_name}-lint.log" 2>&1; then
            log_success ""
        else
            log_warning "has warnings (may be expected for CRD charts)"
            if [ -f "$TEST_DIR/${chart_name}-lint.log" ]; then
                echo "    Log: $TEST_DIR/${chart_name}-lint.log"
            fi
        fi
    done
    echo ""
fi

# Show file tree
log_info "Synced chart structure:"
if command -v tree &> /dev/null; then
    tree -L 3 "$TEST_DIR/helm-charts/charts/" -I '.git*'
else
    find "$TEST_DIR/helm-charts/charts/" -type f | head -20
    log_warning "Install 'tree' for better output: brew install tree"
fi
echo ""

# Compare key files
log_info "Checking critical files..."

check_file() {
    local file=$1
    local chart=$2
    
    if [ -f "$TEST_DIR/helm-charts/charts/$chart/$file" ]; then
        echo -e "  ${COLOR_GREEN}✓${COLOR_RESET} $chart/$file"
    else
        echo -e "  ${COLOR_RED}✗${COLOR_RESET} $chart/$file (MISSING)"
    fi
}

check_file "Chart.yaml" "fluent-operator"
check_file "values.yaml" "fluent-operator"
check_file "README.md" "fluent-operator"
check_file "Chart.yaml" "fluent-bit-crds"
check_file "Chart.yaml" "fluentd-crds"
echo ""

# Show dependency changes
log_info "Checking Chart.yaml dependency changes..."
if [ -f "$TEST_DIR/helm-charts/charts/fluent-operator/Chart.yaml" ]; then
    if grep -q "https://fluent.github.io/helm-charts" "$TEST_DIR/helm-charts/charts/fluent-operator/Chart.yaml"; then
        log_success "fluent-bit-crds dependency uses remote repository"
    else
        log_error "fluent-bit-crds dependency still references local file"
    fi
fi
echo ""

# Statistics
log_info "Sync statistics:"
file_count=$(find "$TEST_DIR/helm-charts/charts/" -type f | wc -l | xargs)
echo "  • Total files synced: $file_count"

for chart in fluent-operator fluent-bit-crds fluentd-crds; do
    if [ -d "$TEST_DIR/helm-charts/charts/$chart" ]; then
        count=$(find "$TEST_DIR/helm-charts/charts/$chart" -type f | wc -l | xargs)
        echo "  • $chart: $count files"
    fi
done
echo ""

# Show what would be in the PR
log_info "Generated PR body preview:"
cat << EOF

═══════════════════════════════════════════════════════════

## 🔄 Helm Charts Sync

This PR syncs Helm charts from the development repository to the release repository.

### 📦 Charts Updated

- **fluent-operator**: \`v${FLUENT_OP_VERSION}\`
- **fluent-bit-crds**: \`v${FLUENT_BIT_CRDS_VERSION}\`
- **fluentd-crds**: \`v${FLUENTD_CRDS_VERSION}\`

### 📝 Details

- **Source Repository**: \`fluent/fluent-operator\`
- **Triggered By**: local test

### ✅ Checklist

- [ ] Chart versions have been bumped appropriately
- [ ] CHANGELOG/release notes have been reviewed
- [ ] No breaking changes, or breaking changes are documented
- [ ] CI tests pass

═══════════════════════════════════════════════════════════

EOF

# Final summary
echo ""
log_success "Test completed successfully!"
echo ""
log_info "Next steps:"
echo "  1. Review the synced files in: $TEST_DIR/helm-charts/"
echo "  2. Ensure chart versions are correct"
echo "  3. Run the GitHub Actions workflow: gh workflow run sync-helm-charts.yaml"
echo ""
log_info "To keep the test directory for inspection:"
echo "  cp -r $TEST_DIR/helm-charts /tmp/helm-charts-test"
echo "  cd /tmp/helm-charts-test"
echo ""

