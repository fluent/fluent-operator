# Helm Charts Sync Workflow

This document explains how to set up and use the automated Helm charts sync workflow that synchronizes charts from this development repository to the [fluent/helm-charts](https://github.com/fluent/helm-charts) release repository.

## Overview

The workflow automatically:
1. ✅ Syncs Helm charts from `charts/` to fluent/helm-charts
2. ✅ Creates a new branch in the helm-charts repository
3. ✅ Opens a pull request with detailed changelog
4. ✅ Validates charts with `helm lint`
5. ✅ Updates chart dependencies to use proper repository URLs

## Setup Instructions

### 1. Create and Configure GitHub App

The workflow uses a GitHub App for authentication, which is more secure than Personal Access Tokens.

**Quick Steps:**
1. Create GitHub App in the `fluent` organization
2. Grant permissions: **Contents** (Write) and **Pull requests** (Write)
3. Generate and download private key
4. Install app on `fluent/helm-charts` repository
5. Add app credentials as secrets

**📖 Detailed Guide:** See [GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md) for complete step-by-step instructions.

### 2. Add Secrets to This Repository

1. Go to this repository's Settings → Secrets and variables → Actions
2. Add two secrets:

   **Secret 1:**
   - Name: `HELM_SYNC_APP_ID`
   - Value: Your GitHub App ID (e.g., `123456`)

   **Secret 2:**
   - Name: `HELM_SYNC_APP_PRIVATE_KEY`
   - Value: Complete contents of the `.pem` file (including BEGIN/END lines)

### 3. Verify Setup

Ensure the app is configured correctly:
- ✅ App created in `fluent` organization
- ✅ App has **Contents: Write** and **Pull requests: Write** permissions
- ✅ App installed on `fluent/helm-charts` repository
- ✅ Both secrets added to `fluent/fluent-operator`
- ✅ Private key file includes `-----BEGIN RSA PRIVATE KEY-----` and `-----END RSA PRIVATE KEY-----`

## Usage

### Automatic Triggers

The workflow runs automatically on:

1. **New Release**: When you publish a release via GitHub Releases
   ```bash
   gh release create v3.5.0 --title "Release v3.5.0" --notes "Release notes..."
   ```

2. **Version Tags**: When you push tags matching version patterns
   ```bash
   git tag v3.5.0
   git push origin v3.5.0
   ```
   
   Or chart-specific tags:
   ```bash
   git tag chart-fluent-operator-3.5.0
   git push origin chart-fluent-operator-3.5.0
   ```

### Manual Trigger

You can also run the workflow manually with custom options:

#### Via GitHub UI

1. Go to Actions → "Sync Helm Charts to Release Repository"
2. Click "Run workflow"
3. Configure options:
   - **Target branch**: Branch in helm-charts repo (default: `main`)
   - **PR title**: Custom title for the PR (optional)
   - **Charts to sync**: `all`, or comma-separated list (e.g., `fluent-operator,fluentd-crds`)
4. Click "Run workflow"

#### Via GitHub CLI

```bash
# Sync all charts (default)
gh workflow run sync-helm-charts.yaml

# Sync specific charts
gh workflow run sync-helm-charts.yaml \
  -f charts_to_sync="fluent-operator,fluentd-crds"

# Use custom PR title and target branch
gh workflow run sync-helm-charts.yaml \
  -f target_branch="main" \
  -f pr_title="Update charts for v3.5.0 release"
```

## Workflow Details

### Charts Synced

The workflow syncs these charts:

| Chart | Source Path | Destination |
|-------|-------------|-------------|
| fluent-operator | `charts/fluent-operator/` | `fluent/helm-charts/charts/fluent-operator/` |
| fluent-bit-crds | `charts/fluent-operator/charts/fluent-bit-crds/` | `fluent/helm-charts/charts/fluent-bit-crds/` |
| fluentd-crds | `charts/fluentd-crds/` | `fluent/helm-charts/charts/fluentd-crds/` |

### What Gets Synced

✅ **Included:**
- Chart.yaml
- values.yaml
- templates/
- crds/
- README.md
- All chart files and directories

❌ **Excluded:**
- .git directories
- Local chart dependencies (replaced with proper repo URLs)

### Chart Modifications

The workflow automatically:
1. Updates `fluent-operator/Chart.yaml` dependency references
   - Changes `file://charts/fluent-bit-crds` → `https://fluent.github.io/helm-charts`
2. Validates all charts with `helm lint`
3. Extracts version numbers for PR metadata

### Pull Request Content

The created PR includes:
- 📦 Chart versions being synced
- 🔗 Link to source commit
- 📝 Workflow run details
- ✅ Review checklist
- 🏷️ Labels: `automated`, `helm-sync`

## Troubleshooting

### Issue: "Resource not accessible by integration"

**Cause**: The GitHub App doesn't have required permissions or isn't installed

**Solution**:
1. Verify app is installed on `fluent/helm-charts`:
   - Go to: https://github.com/organizations/fluent/settings/installations
   - Check "Fluent Helm Charts Sync" shows `helm-charts` in access list
2. Verify app has required permissions:
   - Go to app settings
   - Ensure "Contents: Read and write" is enabled
   - Ensure "Pull requests: Read and write" is enabled
3. If permissions were changed, reinstall the app on helm-charts repository

### Issue: "No changes detected"

**Cause**: Charts in helm-charts repo are already up to date

**Solution**:
- This is normal! The workflow is idempotent
- Bump chart versions in Chart.yaml if you need to force an update
- Check that you've made actual changes to chart files

### Issue: Helm lint failures

**Cause**: Chart validation errors

**Solution**:
- CRD-only charts may show lint warnings (expected)
- Fix actual errors before syncing:
  ```bash
  helm lint charts/fluent-operator/
  helm lint charts/fluentd-crds/
  ```

### Issue: PR creation fails

**Cause**: Branch or PR might already exist

**Solution**:
1. Check if a similar PR is already open in helm-charts
2. Close/merge existing PR before re-running
3. Delete the branch in helm-charts if it exists:
   ```bash
   cd /tmp && git clone https://github.com/fluent/helm-charts.git
   cd helm-charts
   git push origin --delete sync/fluent-operator-X.X.X-YYYYMMDD
   ```

## Monitoring

### Check Workflow Status

```bash
# List recent workflow runs
gh run list --workflow=sync-helm-charts.yaml

# View specific run
gh run view <run-id>

# Watch live run
gh run watch
```

### Notifications

Configure GitHub notifications:
1. Watch the helm-charts repository
2. Settings → Notifications → Actions
3. Enable "Pull requests" notifications

## Best Practices

### 1. Version Bumping

Always bump chart versions before syncing:
```yaml
# charts/fluent-operator/Chart.yaml
version: 3.5.1  # Increment according to semver
```

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: Backwards-compatible functionality
- **PATCH**: Backwards-compatible bug fixes

### 2. Testing Before Release

Test charts locally before syncing:
```bash
# Install chart locally
helm install test-release ./charts/fluent-operator \
  --dry-run --debug

# Run helm tests
make helm-e2e
```

### 3. Review PR in helm-charts

Even though automated:
- ✅ Review the PR created in helm-charts
- ✅ Ensure CI passes
- ✅ Check for unintended changes
- ✅ Get approval from maintainers

### 4. Sync Timing

Best times to sync:
- ✅ After successful CI runs in this repo
- ✅ After release notes are prepared
- ✅ During normal working hours (easier to respond to issues)
- ❌ Not right before weekends/holidays

## Release Process

Recommended end-to-end process:

```bash
# 1. Update versions in Chart.yaml files
vim charts/fluent-operator/Chart.yaml
vim charts/fluentd-crds/Chart.yaml

# 2. Test locally
make helm-e2e

# 3. Commit and push
git add charts/
git commit -m "chore: bump chart versions to 3.5.1"
git push origin master

# 4. Create release (triggers workflow automatically)
gh release create v3.5.1 \
  --title "Release v3.5.1" \
  --notes-file RELEASE_NOTES.md

# 5. Monitor workflow
gh run watch

# 6. Review and merge PR in helm-charts
# Once merged, CI there will publish to Helm repo
```

## Manual Sync (Fallback)

If the workflow fails, you can sync manually:

```bash
# Clone both repos
git clone https://github.com/fluent/fluent-operator.git
git clone https://github.com/fluent/helm-charts.git

# Copy charts
cd fluent-operator
rsync -av --exclude='charts/' charts/fluent-operator/ ../helm-charts/charts/fluent-operator/
rsync -av charts/fluent-operator/charts/fluent-bit-crds/ ../helm-charts/charts/fluent-bit-crds/
rsync -av charts/fluentd-crds/ ../helm-charts/charts/fluentd-crds/

# Update Chart.yaml
cd ../helm-charts
sed -i 's|file://charts/fluent-bit-crds|https://fluent.github.io/helm-charts|g' \
  charts/fluent-operator/Chart.yaml

# Create PR
git checkout -b sync/manual-$(date +%Y%m%d)
git add .
git commit -m "Sync charts from fluent-operator"
git push origin sync/manual-$(date +%Y%m%d)
gh pr create --title "Manual chart sync" --body "Manual sync of charts"
```

## Security Considerations

### GitHub App Security
- 🔒 Never commit private key to git
- 🔒 Store credentials as repository secrets only
- 🔒 Rotate private key annually (see [GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md#rotate-private-key))
- 🔒 Use minimum required permissions (Contents + Pull requests only)
- 🔒 Install app only on required repositories
- 🔒 Monitor app activity via GitHub's installation logs

### Audit Trail
- All syncs are tracked in workflow runs
- PR descriptions include source commit SHA
- Git history preserved in helm-charts
- GitHub App provides attribution: commits/PRs show as "app-name[bot]"
- App activity visible in organization installations page

## Support

If you encounter issues:

1. Check this documentation
2. Review workflow logs in Actions tab
3. Check existing Issues in this repository
4. Contact maintainers: @wenchajun, @marcofranssen, @joshuabaird

## References

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Helm Documentation](https://helm.sh/docs/)
- [fluent/helm-charts Repository](https://github.com/fluent/helm-charts)
- [Semantic Versioning](https://semver.org/)

