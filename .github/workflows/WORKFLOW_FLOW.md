# Helm Charts Sync - Workflow Flow

Visual guide to understand how the sync workflow operates.

## 📊 High-Level Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    TRIGGER (One of three ways)                  │
├─────────────────────────────────────────────────────────────────┤
│  1. Release Published     2. Version Tag Pushed    3. Manual    │
│     gh release create        git push v3.5.0         GitHub UI  │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                   GITHUB ACTIONS WORKFLOW                       │
│              (sync-helm-charts.yaml executes)                   │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
        ┌────────────────┐
        │ Checkout Repos │
        │  • source      │
        │  • helm-charts │
        └────────┬───────┘
                 │
                 ▼
     ┌───────────────────────┐
     │ Extract Chart Info    │
     │  • Versions           │
     │  • Metadata           │
     │  • Create branch name │
     └──────────┬────────────┘
                │
                ▼
     ┌───────────────────────┐
     │ Sync Charts           │
     │  ✓ fluent-operator    │
     │  ✓ fluent-bit-crds    │
     │  ✓ fluentd-crds       │
     └──────────┬────────────┘
                │
                ▼
     ┌───────────────────────┐
     │ Update References     │
     │  • Fix dependencies   │
     │  • Update repo URLs   │
     └──────────┬────────────┘
                │
                ▼
     ┌───────────────────────┐
     │ Validate Charts       │
     │  • Helm lint          │
     │  • File checks        │
     └──────────┬────────────┘
                │
                ▼
     ┌───────────────────────┐
     │ Commit & Push         │
     │  • Create commit      │
     │  • Push to branch     │
     └──────────┬────────────┘
                │
                ▼
     ┌───────────────────────┐
     │ Create Pull Request   │
     │  • Generate PR body   │
     │  • Add labels         │
     │  • Target main branch │
     └──────────┬────────────┘
                │
                ▼
┌─────────────────────────────────────────────────────────────────┐
│                   PR CREATED IN HELM-CHARTS                     │
│              https://github.com/fluent/helm-charts              │
└────────────────┬────────────────────────────────────────────────┘
                 │
                 ▼
        ┌────────────────┐
        │ Manual Review  │
        │  by maintainer │
        └────────┬───────┘
                 │
                 ▼
        ┌────────────────┐
        │   Merge PR     │
        └────────┬───────┘
                 │
                 ▼
┌─────────────────────────────────────────────────────────────────┐
│              HELM-CHARTS CI PUBLISHES CHARTS                    │
│          Charts available at fluent.github.io/helm-charts       │
└─────────────────────────────────────────────────────────────────┘
```

## 🔄 Detailed Step-by-Step

### Phase 1: Trigger

**Automatic Triggers:**
```bash
# Release trigger
gh release create v3.5.0 --title "v3.5.0" --notes "..."
# → Workflow starts automatically

# Tag trigger
git tag v3.5.0 && git push origin v3.5.0
# → Workflow starts automatically
```

**Manual Trigger:**
```bash
gh workflow run sync-helm-charts.yaml
# → Workflow starts on demand
```

### Phase 2: Repository Checkout

```
fluent-operator/            helm-charts/
    ├── charts/                ├── charts/
    │   ├── fluent-operator/   │   ├── fluent-operator/
    │   └── fluentd-crds/      │   ├── fluent-bit-crds/
    └── .git/                  │   └── fluentd-crds/
                               └── .git/

         Both cloned into workflow runner
```

### Phase 3: Chart Analysis

```yaml
# Workflow extracts:
fluent-operator: v4.0.0-rc1   # from charts/fluent-operator/Chart.yaml
fluent-bit-crds: v3.5.0       # from charts/fluent-operator/charts/...
fluentd-crds: v0.1.0          # from charts/fluentd-crds/Chart.yaml

# Creates branch name:
sync/fluent-operator-4.0.0-rc1-20251120-143022
```

### Phase 4: Sync Operations

```
SOURCE                                    DESTINATION
charts/fluent-operator/                   helm-charts/charts/fluent-operator/
    ├── Chart.yaml            ─────►          ├── Chart.yaml (modified)
    ├── values.yaml           ─────►          ├── values.yaml
    ├── README.md             ─────►          ├── README.md
    ├── templates/            ─────►          ├── templates/
    │   ├── deployment.yaml                   │   ├── deployment.yaml
    │   └── ...                               │   └── ...
    └── charts/                               └── (dependencies remote)
        └── fluent-bit-crds/  ─────►
                                          helm-charts/charts/fluent-bit-crds/
charts/fluentd-crds/          ─────►      helm-charts/charts/fluentd-crds/
```

### Phase 5: Transformations

**Before (in fluent-operator):**
```yaml
# charts/fluent-operator/Chart.yaml
dependencies:
  - name: fluent-bit-crds
    repository: "file://charts/fluent-bit-crds"  # Local reference
    version: 3.5.0
```

**After (in helm-charts):**
```yaml
# charts/fluent-operator/Chart.yaml
dependencies:
  - name: fluent-bit-crds
    repository: "https://fluent.github.io/helm-charts"  # Remote reference
    version: 3.5.0
```

### Phase 6: Validation

```bash
helm lint charts/fluent-operator/
# ==> Linting charts/fluent-operator
# [INFO] Chart.yaml: icon is recommended
# 1 chart(s) linted, 0 chart(s) failed

helm lint charts/fluent-bit-crds/
# ==> Linting charts/fluent-bit-crds
# 1 chart(s) linted, 0 chart(s) failed

helm lint charts/fluentd-crds/
# ==> Linting charts/fluentd-crds
# 1 chart(s) linted, 0 chart(s) failed
```

### Phase 7: Git Operations

```bash
cd helm-charts
git checkout main
git checkout -b sync/fluent-operator-4.0.0-rc1-20251120-143022

git add charts/
git commit -m "Sync charts from fluent-operator

- fluent-operator: v4.0.0-rc1
- fluent-bit-crds: v3.5.0
- fluentd-crds: v0.1.0

Source: fluent/fluent-operator@abc1234"

git push origin sync/fluent-operator-4.0.0-rc1-20251120-143022
```

### Phase 8: PR Creation

**Generated PR:**

```markdown
Title: Sync Helm charts (fluent-operator v4.0.0-rc1)
Labels: automated, helm-sync
Base: main
Head: sync/fluent-operator-4.0.0-rc1-20251120-143022

## 🔄 Helm Charts Sync

This PR syncs Helm charts from the development repository.

### 📦 Charts Updated
- **fluent-operator**: `v4.0.0-rc1`
- **fluent-bit-crds**: `v3.5.0`
- **fluentd-crds**: `v0.1.0`

### 📝 Details
- **Source Repository**: `fluent/fluent-operator`
- **Source Commit**: `abc1234`
- **Triggered By**: `release`

### ✅ Checklist
- [ ] Chart versions bumped
- [ ] Release notes reviewed
- [ ] Breaking changes documented
- [ ] CI tests pass
```

## 🔀 Decision Flow

```
                    ┌─────────────┐
                    │  Workflow   │
                    │   Starts    │
                    └──────┬──────┘
                           │
                           ▼
                    ┌──────────────┐
                    │ Extract      │
                    │ Chart Info   │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │ Charts to    │
                    │ sync = "all" │
                    │     OR       │
                    │  specific?   │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │                         │
         "all" or empty            specific list
              │                         │
              ▼                         ▼
    ┌─────────────────┐       ┌────────────────┐
    │  Sync all 3     │       │  Sync only     │
    │  charts         │       │  listed charts │
    └────────┬────────┘       └────────┬───────┘
             │                         │
             └──────────┬──────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │  Any changes?   │
              └────────┬────────┘
                       │
           ┌───────────┼───────────┐
           │                       │
          Yes                      No
           │                       │
           ▼                       ▼
    ┌──────────────┐      ┌──────────────┐
    │ Commit &     │      │ Exit with    │
    │ Create PR    │      │ "No changes" │
    └──────────────┘      └──────────────┘
```

## 📊 Timeline Example

Real-world example of workflow execution:

```
Time    Action                               Duration
─────────────────────────────────────────────────────
00:00   Workflow triggered (release)         -
00:01   Checkout source repo                 ~10s
00:11   Checkout helm-charts repo            ~8s
00:19   Setup Git config                     ~1s
00:20   Install Helm                         ~5s
00:25   Extract chart versions               ~2s
00:27   Create sync branch                   ~1s
00:28   Sync fluent-operator                 ~3s
00:31   Sync fluent-bit-crds                 ~2s
00:33   Sync fluentd-crds                    ~2s
00:35   Update Chart.yaml refs               ~1s
00:36   Validate charts (helm lint)          ~8s
00:44   Generate PR body                     ~1s
00:45   Commit changes                       ~2s
00:47   Push to branch                       ~5s
00:52   Create Pull Request                  ~3s
00:55   Post summary                         ~1s
─────────────────────────────────────────────────────
Total: ~55-60 seconds
```

## 🔄 Idempotency

The workflow is safe to run multiple times:

```
Run 1: Changes detected → PR created ✓
Run 2: No changes       → Exit gracefully ✓
Run 3: No changes       → Exit gracefully ✓

(Make chart changes)

Run 4: Changes detected → PR created ✓
```

## 🎯 Success Criteria

Workflow succeeds when:

✅ All specified charts are synced
✅ Chart dependencies are updated correctly
✅ Helm lint passes (or expected warnings only)
✅ Commit is created with proper message
✅ Branch is pushed successfully
✅ PR is created with complete information
✅ Labels are applied correctly
✅ Summary is generated

## 🔗 Integration Points

```
fluent-operator repo
    │
    ├─► GitHub Actions (this workflow)
    │       │
    │       ├─► Checks out both repos
    │       ├─► Syncs files
    │       └─► Creates PR
    │
    └─► fluent/helm-charts repo
            │
            ├─► PR awaits review
            │
            ├─► Maintainer approves & merges
            │
            └─► helm-charts CI runs
                    │
                    └─► Charts published to
                        https://fluent.github.io/helm-charts
                            │
                            └─► Users can helm install
```

## 📚 Related Documentation

- [SETUP_CHECKLIST.md](SETUP_CHECKLIST.md) - Initial setup
- [HELM_SYNC_README.md](HELM_SYNC_README.md) - Complete guide
- [test-sync.sh](test-sync.sh) - Local testing
- [README.md](README.md) - Workflows overview

---

**Last Updated**: November 2025
**Questions?** Check [HELM_SYNC_README.md](HELM_SYNC_README.md) or ask maintainers

