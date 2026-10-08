# GitHub Actions Workflows

This directory contains automated workflows for the Fluent Operator project.

## 📋 Available Workflows

### Active Workflows

| Workflow | Description | Trigger |
|----------|-------------|---------|
| [sync-helm-charts.yaml](sync-helm-charts.yaml) | Syncs Helm charts to fluent/helm-charts for release | Release, Tags, Manual |
| [helm-ci.yaml](helm-ci.yaml) | Helm chart testing and validation | PR to master/release-* |
| [build-op-image.yaml](build-op-image.yaml) | Builds operator Docker images | Push to master |
| [build-fluentbit-image.yaml](build-fluentbit-image.yaml) | Builds Fluent Bit Docker images | Push to master |
| [build-fluentd-image.yaml](build-fluentd-image.yaml) | Builds Fluentd Docker images | Push to master |
| [test-e2e.yml](test-e2e.yml) | End-to-end integration tests | PR |
| [lint.yml](lint.yml) | Code linting and formatting | PR |

## 🔄 Helm Charts Sync Workflow

The **sync-helm-charts** workflow automates the process of syncing development Helm charts to the release repository at https://github.com/fluent/helm-charts.

### Quick Start

1. **Setup** (one-time):
   ```bash
   # Follow the setup checklist
   cat SETUP_CHECKLIST.md
   ```

2. **Test locally**:
   ```bash
   # Run the test script to validate sync process
   ./test-sync.sh
   ```

3. **Trigger sync**:
   ```bash
   # Automatic on release
   gh release create v3.5.0

   # Or manual trigger
   gh workflow run sync-helm-charts.yaml
   ```

### Documentation

- **[SETUP_CHECKLIST.md](SETUP_CHECKLIST.md)** - Step-by-step setup guide ⭐ Start here!
- **[HELM_SYNC_README.md](HELM_SYNC_README.md)** - Complete documentation with troubleshooting
- **[test-sync.sh](test-sync.sh)** - Local testing script

### Key Features

✅ **Automatic Sync**: Triggered on releases and version tags  
✅ **Manual Control**: Run manually with custom options  
✅ **Validation**: Helm lint checks on all charts  
✅ **PR Automation**: Creates formatted PRs with changelogs  
✅ **Safe**: No changes without review - always creates PR first  
✅ **Idempotent**: Safe to run multiple times  

### What Gets Synced

| Chart | Source | Destination |
|-------|--------|-------------|
| fluent-operator | `charts/fluent-operator/` | `fluent/helm-charts` |
| fluent-bit-crds | `charts/fluent-operator/charts/fluent-bit-crds/` | `fluent/helm-charts` |
| fluentd-crds | `charts/fluentd-crds/` | `fluent/helm-charts` |

## 🔐 Secrets Required

### GitHub App Credentials (Recommended)

Required for: `sync-helm-charts.yaml`

**HELM_SYNC_APP_ID**: GitHub App ID number  
**HELM_SYNC_APP_PRIVATE_KEY**: GitHub App private key (`.pem` file contents)

**What it is**: GitHub App authentication for cross-repository operations

**Permissions needed**:
- Contents: Read and write
- Pull requests: Read and write

**Setup**: See [GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md) for detailed instructions

**Why GitHub App?**
- ✅ More secure than Personal Access Tokens
- ✅ Not tied to individual user accounts
- ✅ Better audit trail and attribution
- ✅ No expiration concerns
- ✅ Organization-level management

## 📚 Additional Resources

### Testing

```bash
# Test helm chart sync locally
./test-sync.sh

# Run helm lint
helm lint charts/fluent-operator/
helm lint charts/fluentd-crds/

# Test workflow syntax (requires actionlint)
actionlint sync-helm-charts.yaml
```

### Monitoring

```bash
# List workflow runs
gh run list --workflow=sync-helm-charts.yaml

# Watch a running workflow
gh run watch

# View workflow logs
gh run view <run-id> --log
```

### Common Commands

```bash
# Manually trigger helm sync
gh workflow run sync-helm-charts.yaml

# Sync specific charts only
gh workflow run sync-helm-charts.yaml \
  -f charts_to_sync="fluent-operator,fluentd-crds"

# Use custom PR title
gh workflow run sync-helm-charts.yaml \
  -f pr_title="Release v3.5.0 charts"

# Target different branch
gh workflow run sync-helm-charts.yaml \
  -f target_branch="main"
```

## 🐛 Troubleshooting

### Workflow Fails

1. Check workflow logs in Actions tab
2. Verify secrets exist:
   - `HELM_SYNC_APP_ID` (should be a number)
   - `HELM_SYNC_APP_PRIVATE_KEY` (should be large text with BEGIN/END lines)
3. Verify GitHub App is installed on `fluent/helm-charts`
4. Check app has required permissions (Contents + Pull requests)
5. See [GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md#troubleshooting) for detailed troubleshooting

### PR Not Created

1. Verify branch doesn't already exist in helm-charts
2. Check that there are actual changes to commit
3. Ensure token has PR creation permissions
4. Review workflow logs for error messages

### Helm Lint Fails

```bash
# Test locally first
helm lint charts/fluent-operator/
helm lint charts/fluentd-crds/

# CRD-only charts may show warnings (expected)
```

## 🤝 Contributing

When adding or modifying workflows:

1. Test locally when possible
2. Use pinned action versions (with SHA)
3. Add documentation to this README
4. Update CONTRIBUTING.md if process changes
5. Test on a fork first for major changes

## 📞 Support

- **Documentation**: Start with [HELM_SYNC_README.md](HELM_SYNC_README.md)
- **Issues**: Check existing issues or create new one
- **Maintainers**: @wenchajun @marcofranssen @joshuabaird
- **Community**: Fluent Slack / GitHub Discussions

## 🔗 Related Links

- [Fluent Operator Docs](https://github.com/fluent/fluent-operator/blob/master/README.md)
- [Helm Charts Repository](https://github.com/fluent/helm-charts)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Helm Documentation](https://helm.sh/docs/)

---

**Last Updated**: November 2025  
**Maintained By**: Fluent Operator Team

