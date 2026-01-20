# Helm Charts Sync - Setup Checklist

Follow these steps to set up the automated Helm chart sync workflow.

## ✅ Setup Steps

### Step 1: Create GitHub App

⭐ **See [GITHUB_APP_SETUP.md](GITHUB_APP_SETUP.md) for detailed instructions**

- [ ] Create GitHub App in fluent organization
- [ ] Name: `Fluent Helm Charts Sync`
- [ ] Set permissions:
  - [ ] Contents: Read and write
  - [ ] Pull requests: Read and write
- [ ] Disable webhook (not needed)
- [ ] Generate private key (downloads `.pem` file)
- [ ] **Save App ID** (shown at top of page)
- [ ] **Save private key file** securely

### Step 2: Install GitHub App
- [ ] Go to "Install App" tab in app settings
- [ ] Install on fluent organization
- [ ] Grant access to: `fluent/helm-charts` (only)
- [ ] Complete installation

### Step 3: Add Secrets to Repository
- [ ] Go to https://github.com/fluent/fluent-operator/settings/secrets/actions
- [ ] Add first secret:
  - [ ] Name: `HELM_SYNC_APP_ID`
  - [ ] Value: The App ID from Step 1
- [ ] Add second secret:
  - [ ] Name: `HELM_SYNC_APP_PRIVATE_KEY`
  - [ ] Value: Complete contents of `.pem` file (including BEGIN/END lines)

### Step 3: Test Locally (Optional but Recommended)
```bash
# Run the test script
cd /Users/jbaird/code/fluent-operator
./.github/workflows/test-sync.sh
```

- [ ] Test script runs successfully
- [ ] No missing files reported
- [ ] Helm lint passes (or expected warnings only)

### Step 4: Test the Workflow
```bash
# Test with manual trigger
gh workflow run sync-helm-charts.yaml \
  -f pr_title="[TEST] Chart sync test" \
  -f charts_to_sync="fluentd-crds"

# Watch the run
gh run watch
```

- [ ] Workflow starts successfully
- [ ] GitHub App authentication succeeds (check logs)
- [ ] Check workflow run: https://github.com/fluent/fluent-operator/actions
- [ ] PR created in fluent/helm-charts
- [ ] PR shows app as author (e.g., "fluent-helm-charts-sync[bot]")
- [ ] PR content looks correct
- [ ] Close the test PR after verification

### Step 5: Set Up Notifications (Optional)
- [ ] Watch the fluent/helm-charts repository
- [ ] Configure notification preferences
- [ ] Consider setting up Slack/Discord webhook for PR notifications

### Step 6: Document for Team
- [ ] Share setup instructions with maintainers
- [ ] Add workflow to release checklist
- [ ] Update CONTRIBUTING.md if needed

## 🔐 Security Verification

- [ ] App credentials stored as repository secrets (not committed to code)
- [ ] App has minimum required permissions (Contents + Pull requests only)
- [ ] App installed only on `fluent/helm-charts` (not all repos)
- [ ] Private key file deleted from local machine after adding to secrets
- [ ] Only trusted maintainers have access to secrets
- [ ] App ownership documented for team

## 📋 Testing Checklist

After setup, verify:
- [ ] Can sync all charts manually
- [ ] Can sync individual charts
- [ ] PR is created with correct format
- [ ] PR labels are applied
- [ ] Workflow summary is generated
- [ ] No sensitive data exposed in logs

## 🚀 First Real Sync

Ready to do your first production sync:

```bash
# 1. Bump chart versions if needed
vim charts/fluent-operator/Chart.yaml
# Update version: 3.5.1

# 2. Commit changes
git add charts/
git commit -m "chore: bump fluent-operator chart to 3.5.1"
git push origin master

# 3. Create a release (triggers workflow automatically)
gh release create v3.5.1 \
  --title "Release v3.5.1" \
  --notes "See CHANGELOG.md for details"

# 4. Monitor workflow
gh run watch

# 5. Review PR in helm-charts
# Check: https://github.com/fluent/helm-charts/pulls

# 6. Merge PR once CI passes
# Charts will be automatically released!
```

## 📚 Documentation

After setup is complete, reference these documents:

- [HELM_SYNC_README.md](./HELM_SYNC_README.md) - Complete documentation
- [sync-helm-charts.yaml](./sync-helm-charts.yaml) - The workflow file
- [test-sync.sh](./test-sync.sh) - Local testing script

## ⚠️ Troubleshooting

If something doesn't work:

1. **Check GitHub App installation**
   - Go to: https://github.com/organizations/fluent/settings/installations
   - Verify "Fluent Helm Charts Sync" is installed
   - Check it has access to `fluent/helm-charts`

2. **Verify app permissions**
   - Go to app settings
   - Ensure "Contents: Read and write" is set
   - Ensure "Pull requests: Read and write" is set

3. **Check workflow logs**
   - Go to Actions tab
   - Click on failed workflow
   - Look for "Generate token" step errors
   - Review step logs for authentication errors

4. **Verify secrets are set**
   - Settings → Secrets → Actions
   - Ensure `HELM_SYNC_APP_ID` exists (should be a number like 123456)
   - Ensure `HELM_SYNC_APP_PRIVATE_KEY` exists (should be large text)

5. **Common errors**
   - "Resource not accessible" → App not installed on helm-charts
   - "Bad credentials" → Check App ID or private key is correct
   - "Could not resolve to a Repository" → App not granted access to helm-charts

## 🎉 Success Indicators

You'll know everything is working when:
- ✅ Workflow runs without errors
- ✅ PR is created in fluent/helm-charts automatically
- ✅ PR has correct content and labels
- ✅ Charts pass validation
- ✅ Merging PR triggers helm-charts CI
- ✅ Charts are published to https://fluent.github.io/helm-charts

## 📞 Support

Need help?
- Check workflow logs first
- Review HELM_SYNC_README.md
- Ask in #fluent-operator channel
- Tag maintainers: @wenchajun @marcofranssen @joshuabaird

---

**Setup Date**: _________________

**Setup By**: _________________

**GitHub App**: _________________

**Next Key Rotation**: _________________ (recommended annually)

