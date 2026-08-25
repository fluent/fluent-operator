# Getting Started: Helm Charts Sync with GitHub App

**Updated to use GitHub App authentication** 🚀

This guide will get you from zero to a working automated Helm chart sync in ~15 minutes.

## 🎯 What You're Setting Up

An automated workflow that:
- Syncs Helm charts from `fluent/fluent-operator` → `fluent/helm-charts`
- Creates pull requests automatically
- Runs on releases, tags, or manually
- Uses secure GitHub App authentication

## ⚡ Quick Start (15 minutes)

### Step 1: Create GitHub App (5 min)

1. **Go to:** https://github.com/organizations/fluent/settings/apps/new

2. **Fill in:**
   - **Name:** `Fluent Helm Charts Sync`
   - **Homepage:** `https://github.com/fluent/fluent-operator`
   - **Webhook:** ❌ Uncheck "Active"

3. **Permissions:**
   - Contents: `Read and write`
   - Pull requests: `Read and write`

4. **Click:** "Create GitHub App"

5. **Generate private key:**
   - Scroll to "Private keys"
   - Click "Generate a private key"
   - Save the downloaded `.pem` file

6. **Note your App ID:**
   - At the top of the page: "App ID: 123456"
   - Write it down!

### Step 2: Install App (2 min)

1. **Click:** "Install App" (left sidebar)
2. **Select:** Only `fluent/helm-charts` repository
3. **Click:** "Install"

### Step 3: Add Secrets (3 min)

1. **Go to:** https://github.com/fluent/fluent-operator/settings/secrets/actions

2. **Add Secret #1:**
   - Click "New repository secret"
   - Name: `HELM_SYNC_APP_ID`
   - Value: Your App ID from Step 1.6
   - Click "Add secret"

3. **Add Secret #2:**
   - Click "New repository secret"
   - Name: `HELM_SYNC_APP_PRIVATE_KEY`
   - Value: Entire contents of `.pem` file (including BEGIN/END lines)
   - Click "Add secret"

### Step 4: Test It! (5 min)

```bash
# Test the workflow manually
gh workflow run sync-helm-charts.yaml \
  -f pr_title="[TEST] Initial sync test" \
  -f charts_to_sync="fluentd-crds"

# Watch it run
gh run watch

# Expected: PR created in fluent/helm-charts by "fluent-helm-charts-sync[bot]"
```

**Success indicators:**
- ✅ Workflow completes without errors
- ✅ PR created in fluent/helm-charts
- ✅ PR author shows as your app bot
- ✅ PR has properly formatted description

**If it works:** Close the test PR and celebrate! 🎉

**If it fails:** See [Troubleshooting](#-troubleshooting) below

## 📖 Detailed Documentation

### For Complete Setup Instructions
👉 **[GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md)** - Step-by-step with screenshots

### For Daily Usage
👉 **[HELM_SYNC_README.md](HELM_SYNC_README.md)** - How to use the workflow

### For Understanding the Flow
👉 **[WORKFLOW_FLOW.md](WORKFLOW_FLOW.md)** - Visual diagrams and explanations

### For Quick Reference
👉 **[SETUP_CHECKLIST.md](SETUP_CHECKLIST.md)** - Checklist format

### For Decision-Making
👉 **[GITHUB_APP_VS_PAT.md](GITHUB_APP_VS_PAT.md)** - Why GitHub App vs PAT

## 🚀 Daily Usage

Once set up, using the workflow is simple:

### Automatic Sync (Recommended)

```bash
# Just create a release - workflow triggers automatically!
gh release create v3.5.0 \
  --title "Release v3.5.0" \
  --notes-file RELEASE_NOTES.md

# Or push a version tag
git tag v3.5.0
git push origin v3.5.0
```

The workflow will:
1. Detect the release/tag
2. Sync all charts
3. Create a PR in fluent/helm-charts
4. Wait for your review

### Manual Sync

```bash
# Sync all charts
gh workflow run sync-helm-charts.yaml

# Sync specific charts only
gh workflow run sync-helm-charts.yaml \
  -f charts_to_sync="fluent-operator,fluentd-crds"

# Custom PR title
gh workflow run sync-helm-charts.yaml \
  -f pr_title="Update charts for v3.5.0"
```

### Monitoring

```bash
# List recent runs
gh run list --workflow=sync-helm-charts.yaml --limit 5

# Watch current run
gh run watch

# View specific run
gh run view <run-id> --log
```

## 🐛 Troubleshooting

### "Resource not accessible by integration"

**Fix:**
1. Check app is installed on helm-charts:
   ```
   https://github.com/organizations/fluent/settings/installations
   ```
2. Verify "fluent/helm-charts" is in the access list
3. Check permissions: Contents + Pull requests both set to "Write"

### "Bad credentials"

**Fix:**
1. Verify `HELM_SYNC_APP_ID` is just the number (e.g., `123456`)
2. Verify `HELM_SYNC_APP_PRIVATE_KEY` includes:
   ```
   -----BEGIN RSA PRIVATE KEY-----
   [key content]
   -----END RSA PRIVATE KEY-----
   ```
3. No extra spaces or formatting in secrets

### "Could not resolve to a Repository"

**Fix:**
1. App not installed on helm-charts repository
2. Go to installations → Configure → Add helm-charts

### Workflow runs but no PR created

**Possible causes:**
- No changes detected (charts already synced) ✅ This is OK!
- Missing "Pull requests: Write" permission
- Branch already exists in helm-charts

## 📊 Files Created

Here's everything that was created for you:

```
.github/workflows/
├── sync-helm-charts.yaml      (12K) ← The main workflow
├── test-sync.sh              (6.5K) ← Local testing script
│
├── GETTING_STARTED.md        (THIS FILE) ← Start here!
├── GITHUB_APP_SETUP.md       (9.9K) ← Detailed setup guide
├── SETUP_CHECKLIST.md        (5.5K) ← Quick checklist
│
├── HELM_SYNC_README.md       (9.7K) ← Complete usage docs
├── WORKFLOW_FLOW.md          (15K)  ← Visual flow diagrams
├── GITHUB_APP_VS_PAT.md      (7.8K) ← Why we chose GitHub App
└── README.md                 (5.6K) ← Workflows overview
```

**Total:** 72KB of documentation to make your life easier! 📚

## ✅ What Changed from PAT Approach

If you're familiar with Personal Access Tokens:

| Aspect | Old (PAT) | New (GitHub App) |
|--------|-----------|------------------|
| **Secrets** | 1 secret (token) | 2 secrets (app ID + key) |
| **Setup time** | 5 min | 15 min |
| **Expiration** | 1 year | Never |
| **User-tied** | Yes | No |
| **Security** | Good | Better |
| **Attribution** | User | Bot |
| **Maintenance** | Annual renewal | Optional key rotation |

**Bottom line:** 10 extra minutes now saves hours later + better security.

## 🎓 Learning Path

1. **Today:** Follow this quick start
2. **Tomorrow:** Read [HELM_SYNC_README.md](HELM_SYNC_README.md) for details
3. **Later:** Review [WORKFLOW_FLOW.md](WORKFLOW_FLOW.md) to understand internals
4. **Ongoing:** Reference [GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md) for management

## 🔒 Security Checklist

Before going to production:

- [ ] App has minimum permissions (only Contents + Pull requests)
- [ ] App installed only on helm-charts (not all repos)
- [ ] Private key stored as secret (never committed)
- [ ] Local `.pem` file deleted after adding to secrets
- [ ] Both secrets added correctly
- [ ] Test workflow completed successfully
- [ ] Team knows how to manage the app

## 🎉 Success Criteria

You know it's working when:

✅ Workflow runs without errors
✅ PR created automatically in helm-charts
✅ PR author shows as "fluent-helm-charts-sync[bot]"
✅ PR has formatted description with versions
✅ Charts are synced correctly
✅ Merging PR triggers helm-charts CI

## 🆘 Need Help?

**Read in order:**
1. This file's [Troubleshooting](#-troubleshooting) section
2. [GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md#troubleshooting)
3. [HELM_SYNC_README.md](HELM_SYNC_README.md#troubleshooting)
4. GitHub Actions logs in the failed workflow
5. Ask maintainers: @wenchajun @marcofranssen @joshuabaird

## 🚀 Ready?

Let's do this! Start with Step 1 above. ⬆️

---

**Created:** November 2025
**Status:** Production Ready
**Estimated Setup Time:** 15 minutes
**Difficulty:** Beginner-friendly 🟢

