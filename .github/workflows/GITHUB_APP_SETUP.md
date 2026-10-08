# GitHub App Setup for Helm Charts Sync

This guide walks you through creating and configuring a GitHub App for automated chart synchronization.

## 🎯 Why GitHub App?

✅ **More secure** - Fine-grained permissions
✅ **Not tied to user** - Survives team changes
✅ **Better audit trail** - Clear attribution
✅ **Longer lived** - No expiration concerns
✅ **Organization-level** - Centralized management

## 📋 Prerequisites

- Admin access to the `fluent` GitHub organization (or both repos)
- Ability to create GitHub Apps in the organization
- Access to repository secrets in `fluent/fluent-operator`

## 🚀 Step-by-Step Setup

### Step 1: Create the GitHub App

1. **Navigate to GitHub App creation page:**
   - **For Organization**: https://github.com/organizations/fluent/settings/apps/new
   - **For Personal**: https://github.com/settings/apps/new

2. **Fill in basic information:**

   | Field | Value |
   |-------|-------|
   | **GitHub App name** | `Fluent Helm Charts Sync` |
   | **Homepage URL** | `https://github.com/fluent/fluent-operator` |
   | **Webhook** | ❌ Uncheck "Active" (not needed) |

3. **Set Repository permissions:**

   Navigate to "Repository permissions" section and set:

   | Permission | Access Level | Why Needed |
   |------------|--------------|------------|
   | **Contents** | `Read and write` | Clone repos, push commits |
   | **Pull requests** | `Read and write` | Create and manage PRs |
   | **Metadata** | `Read-only` | Required (auto-selected) |

   > **Note:** We only need these 3 permissions - keep everything else at "No access"

4. **Set "Where can this GitHub App be installed?"**
   - ✅ Select: **"Only on this account"**
   - This restricts the app to the `fluent` organization

5. **Click "Create GitHub App"**

### Step 2: Generate and Save Private Key

After creating the app:

1. **Scroll to "Private keys" section** (bottom of app settings page)

2. **Click "Generate a private key"**
   - A `.pem` file will download automatically
   - **⚠️ IMPORTANT**: This is the only time you'll see this key!

3. **Save the file securely:**
   ```bash
   # The file will be named something like:
   # fluent-helm-charts-sync.2024-11-20.private-key.pem

   # Store it securely - you'll need it in the next step
   ```

4. **Note the App ID:**
   - At the top of the page, you'll see "App ID: 123456"
   - Copy this number - you'll need it

### Step 3: Install the App on Repositories

1. **Go to "Install App" tab** (left sidebar of app settings)

2. **Click "Install" next to your organization/account**

3. **Select repository access:**
   - ✅ Select: **"Only select repositories"**
   - Choose: **`fluent/helm-charts`** (the target repository)

   > **Important**: The app only needs access to `helm-charts`, not `fluent-operator`

4. **Click "Install"**

5. **Note the Installation ID (Optional but helpful):**
   - Look at the URL after installation: `https://github.com/organizations/fluent/settings/installations/XXXXX`
   - The number is your installation ID (useful for debugging)

### Step 4: Add Secrets to fluent-operator Repository

Now add the App credentials as secrets in the `fluent-operator` repository:

1. **Navigate to:**
   ```
   https://github.com/fluent/fluent-operator/settings/secrets/actions
   ```

2. **Add Secret #1: App ID**
   - Click "New repository secret"
   - **Name**: `HELM_SYNC_APP_ID`
   - **Value**: The App ID from Step 2.4 (e.g., `123456`)
   - Click "Add secret"

3. **Add Secret #2: Private Key**
   - Click "New repository secret"
   - **Name**: `HELM_SYNC_APP_PRIVATE_KEY`
   - **Value**: The complete contents of the `.pem` file from Step 2.3

   ```bash
   # To get the contents on macOS/Linux:
   cat fluent-helm-charts-sync.2024-11-20.private-key.pem | pbcopy
   # (copies to clipboard)

   # Or simply open the file and copy all contents including:
   # -----BEGIN RSA PRIVATE KEY-----
   # [all the key content]
   # -----END RSA PRIVATE KEY-----
   ```

   - Click "Add secret"

### Step 5: Verify Secrets

Confirm both secrets are added:

```
Repository secrets for fluent/fluent-operator:
✓ HELM_SYNC_APP_ID
✓ HELM_SYNC_APP_PRIVATE_KEY
```

## ✅ Verification Checklist

Before proceeding, verify:

- [ ] GitHub App created successfully
- [ ] App has correct permissions (Contents: Write, Pull requests: Write)
- [ ] Private key generated and downloaded
- [ ] App installed on `fluent/helm-charts` repository
- [ ] `HELM_SYNC_APP_ID` secret added to fluent-operator
- [ ] `HELM_SYNC_APP_PRIVATE_KEY` secret added to fluent-operator
- [ ] Both secrets are visible in repository settings

## 🧪 Testing the Setup

Test that the app authentication works:

```bash
# Run the workflow manually
gh workflow run sync-helm-charts.yaml

# Watch the run
gh run watch

# Check for any authentication errors
```

**Expected behavior:**
- ✅ Workflow authenticates successfully
- ✅ Can clone helm-charts repository
- ✅ Can create branch in helm-charts
- ✅ Can create PR in helm-charts

**If you see authentication errors:**
- Verify app is installed on helm-charts repository
- Check that secrets are correctly named
- Ensure private key includes BEGIN/END lines
- Verify App ID is correct

## 🔄 App Management

### View App Activity

Monitor app usage:
```
https://github.com/organizations/fluent/settings/installations
```

Click on "Fluent Helm Charts Sync" to see:
- Recent deliveries (API calls)
- Events
- Repository access

### Modify Permissions

If you need to change permissions:

1. Go to: `https://github.com/organizations/fluent/settings/apps`
2. Click on "Fluent Helm Charts Sync"
3. Click "Edit" (top right)
4. Modify "Repository permissions"
5. Click "Save changes"

**⚠️ Note:** Permission changes require reinstallation on repositories

### Rotate Private Key

To rotate the key for security:

1. Go to app settings
2. Scroll to "Private keys"
3. Click "Generate a private key" (new key)
4. Update `HELM_SYNC_APP_PRIVATE_KEY` secret in fluent-operator
5. Test workflow
6. Revoke old key (click "Revoke" next to old key)

### Add to More Repositories

To grant access to additional repositories:

1. Go to: `https://github.com/organizations/fluent/settings/installations`
2. Click "Configure" next to "Fluent Helm Charts Sync"
3. Update "Repository access"
4. Save

## 🔒 Security Best Practices

### ✅ Do's

- ✅ Keep private key secure (never commit to git)
- ✅ Use organization-level apps when possible
- ✅ Grant minimum required permissions
- ✅ Regular security audits of app activity
- ✅ Document app purpose and owners
- ✅ Rotate keys annually
- ✅ Monitor app events for unusual activity

### ❌ Don'ts

- ❌ Never commit `.pem` file to repository
- ❌ Don't share private key via unsecured channels
- ❌ Don't grant more permissions than needed
- ❌ Don't install on repositories that don't need it
- ❌ Don't forget to revoke old keys after rotation

## 🆚 GitHub App vs PAT

| Feature | GitHub App | Personal Access Token |
|---------|------------|----------------------|
| **User-independent** | ✅ Yes | ❌ No |
| **Fine-grained permissions** | ✅ Yes | ⚠️ Limited |
| **Audit trail** | ✅ Excellent | ⚠️ Good |
| **Rate limits** | ✅ Higher | ⚠️ Standard |
| **Expiration** | ✅ Never | ❌ Max 1 year |
| **Setup complexity** | ⚠️ Medium | ✅ Easy |
| **Revocation** | ✅ Instant | ✅ Instant |
| **Organization control** | ✅ Centralized | ❌ Per-user |

## 📚 Additional Resources

- [GitHub Apps Documentation](https://docs.github.com/en/apps)
- [About GitHub App permissions](https://docs.github.com/en/apps/creating-github-apps/setting-up-a-github-app/choosing-permissions-for-a-github-app)
- [GitHub Actions: create-github-app-token](https://github.com/actions/create-github-app-token)
- [Authenticating with GitHub Apps](https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/about-authentication-with-a-github-app)

## 🐛 Troubleshooting

### Error: "Resource not accessible by integration"

**Cause:** App doesn't have required permissions or isn't installed on target repo

**Solution:**
1. Verify app is installed on `fluent/helm-charts`
2. Check app has "Contents: Write" and "Pull requests: Write" permissions
3. If permissions changed, reinstall app on repository

### Error: "Bad credentials"

**Cause:** App ID or private key is incorrect

**Solution:**
1. Verify `HELM_SYNC_APP_ID` matches App ID in GitHub
2. Verify `HELM_SYNC_APP_PRIVATE_KEY` is complete (including BEGIN/END lines)
3. Ensure no extra whitespace in secrets
4. Try regenerating private key

### Error: "Could not resolve to a Repository"

**Cause:** App not installed on target repository

**Solution:**
1. Go to app installations
2. Verify `fluent/helm-charts` is in the access list
3. If not, add it via "Configure" → "Repository access"

### Workflow runs but can't create PR

**Cause:** Missing "Pull requests: Write" permission

**Solution:**
1. Edit app permissions
2. Add "Pull requests: Read and write"
3. Reinstall app on helm-charts repository

## 👥 Team Handoff

When sharing this setup with team:

1. **Document the app:**
   - App name: `Fluent Helm Charts Sync`
   - Purpose: Automated Helm chart synchronization
   - Repositories: fluent-operator → helm-charts
   - Owner: [Your team/person]

2. **Share access:**
   - Organization admins can manage the app
   - App settings: `https://github.com/organizations/fluent/settings/apps`

3. **Backup information:**
   - App ID (can regenerate if lost)
   - Private key (generate new if lost)
   - Installation URLs
   - This documentation!

## 📞 Support

Questions? Check:
1. This guide's troubleshooting section
2. [HELM_SYNC_README.md](HELM_SYNC_README.md) for workflow details
3. GitHub's [App troubleshooting guide](https://docs.github.com/en/apps/creating-github-apps/troubleshooting/troubleshooting-github-apps)
4. Contact maintainers: @wenchajun @marcofranssen @joshuabaird

---

**Created**: November 2025
**Last Updated**: November 2025
**Next Review**: November 2026 (or when rotating keys)

