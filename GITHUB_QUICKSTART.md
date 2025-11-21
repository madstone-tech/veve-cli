# GitHub Repository Setup - Quick Start Guide

Complete checklist for setting up your veve-cli repository on GitHub with best practices for branch protection, SEO, and automated management.

## Prerequisites

- GitHub account
- GitHub Personal Access Token (PAT) with `repo`, `admin:repo_hook`, `workflow` scopes
- `gh` CLI installed: `brew install gh`
- `curl` installed (usually pre-installed)

## Quick Setup (5 minutes)

### 1. Create GitHub Repository

Visit [github.com/new](https://github.com/new) and create:
- **Repository name**: `veve-cli`
- **Description**: Leave blank (we'll update via API)
- **Public**: Yes
- **Initialize**: No (we have existing code)

### 2. Push Existing Repository

```bash
cd /Users/andhi/code/mdstn/veve-cli

# Add remote
git remote add origin https://github.com/YOUR_USERNAME/veve-cli.git

# Rename branch to main (if needed)
git branch -M main

# Push code and tags
git push -u origin main
git push origin v0.2.0
```

### 3. Setup GitHub Token

```bash
# Create token: https://github.com/settings/tokens/new
# Scopes: repo, admin:repo_hook, workflow

# Export to environment
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_OWNER="YOUR_USERNAME"
export GITHUB_REPO="veve-cli"
```

### 4. Run Automated Setup

```bash
# Make script executable
chmod +x .github/scripts/setup-repo.sh

# Run setup
./.github/scripts/setup-repo.sh

# Expected output:
# ✓ Repository metadata updated
# ✓ Branch protection configured for main
# ✓ Merge strategies configured
# ✓ Created 14 issue labels
# ✓ GitHub Discussions enabled
```

### 5. Verify Setup

Visit your repository on GitHub and verify:

- [ ] Description updated
- [ ] 10 topics displayed
- [ ] Branch protection enabled on main
- [ ] Issue labels created
- [ ] GitHub Discussions enabled

## What Gets Configured

### Repository Metadata
```
Description: Fast, themeable markdown-to-PDF converter built with Go
Homepage: https://github.com/YOUR_USERNAME/veve-cli
Topics: markdown, pdf, converter, cli, golang, pandoc, themes, pdf-generation, 
        command-line-tool, document-conversion
```

### Branch Protection (main branch)
- ✓ Require status checks to pass
- ✓ Require 1 pull request review
- ✓ Dismiss stale reviews
- ✓ Enforce for administrators
- ✓ Require linear history
- ✓ Require conversation resolution

### Merge Strategies
- ✓ Allow squash merge
- ✓ Allow rebase merge
- ✓ Allow merge commit
- ✓ Auto-delete branches on merge
- ✓ Auto-merge enabled

### Issue Labels (14 total)
- Priority: `high`, `medium`, `low`
- Type: `feature`, `refactor`, `test`, `enhancement`
- Status: `bug`, `documentation`, `question`
- Community: `good first issue`, `help wanted`, `wontfix`

### GitHub Discussions
- ✓ Enabled for community engagement
- ✓ Better for long-form discussions vs issues

## Documentation Included

### For Repository Administrators
- **docs/GITHUB_SETUP.md** - Complete setup and branch protection guide
- **docs/GITHUB_API.md** - API reference with curl examples
- **docs/SEO.md** - Search engine optimization strategy

### For Contributors
- **.github/ISSUE_TEMPLATE/bug_report.md** - Bug report template
- **.github/ISSUE_TEMPLATE/feature_request.md** - Feature request template
- **.github/pull_request_template.md** - PR submission guide

### For Maintenance
- **.github/workflows/repo-management.yml** - Daily health checks

## Verify Everything Works

### Test 1: Check Token Access
```bash
curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/user
# Should return your user info
```

### Test 2: Verify Repository Settings
```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO | jq '{
    description,
    topics,
    has_discussions,
    allow_squash_merge
  }'
```

### Test 3: Check Branch Protection
```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/branches/main/protection \
  | jq '.required_pull_request_reviews'
```

## SEO Setup for Discoverability

The automated setup configures your repository for:

1. **GitHub Search**: Optimized metadata and topics
2. **Google Search**: Keywords in description and documentation
3. **Developer Communities**: Professional templates and guidelines
4. **Package Managers**: Clear installation documentation

Expected results:
- Higher ranking for "markdown to PDF converter"
- More visible in GitHub trending for Go projects
- Better search engine discoverability
- Increased community contributions

## Common Next Steps

### After Initial Setup

1. **Create Issues** (optional)
   ```bash
   gh issue create --title "First Issue" --body "Initial setup complete"
   ```

2. **Test Workflows**
   - Verify CI/CD runs on commits
   - Monitor release workflow
   - Check repo-management workflow runs daily

3. **Announce Release**
   - Create release notes
   - Share on social media
   - Post in developer communities

### Long-term Maintenance

- **Monthly**: Review stale issues, update dependencies
- **Quarterly**: Refresh documentation, audit security
- **Annually**: Major documentation update, strategy review

## Troubleshooting

### Setup Script Fails with 401
```bash
# Token is invalid or expired
# Create new token at: https://github.com/settings/tokens/new
export GITHUB_TOKEN="new_token_here"
./.github/scripts/setup-repo.sh
```

### Setup Script Fails with 403
```bash
# Token doesn't have required scopes
# Regenerate with: repo, admin:repo_hook, workflow
```

### Branch Protection Not Appearing
```bash
# Wait 30 seconds and refresh GitHub page
# Check repository settings > Branches
```

### Topics Not Showing
```bash
# Manually verify:
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO | jq '.topics'
```

## Environment Variable Persistence

To make GitHub token persistent across terminal sessions:

**For Bash:**
```bash
# Add to ~/.bashrc
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_OWNER="andhi"
export GITHUB_REPO="veve-cli"

# Reload
source ~/.bashrc
```

**For Zsh:**
```bash
# Add to ~/.zshrc
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_OWNER="andhi"
export GITHUB_REPO="veve-cli"

# Reload
source ~/.zshrc
```

**For Fish:**
```bash
# Run once
set -Ux GITHUB_TOKEN "ghp_your_token_here"
set -Ux GITHUB_OWNER "andhi"
set -Ux GITHUB_REPO "veve-cli"
```

## Manual Setup (Alternative)

If you prefer manual setup instead of running the script:

### 1. Update Repository Metadata
Visit: GitHub → Settings → General
- Description: "Fast, themeable markdown-to-PDF converter built with Go"
- Website: https://github.com/andhi/veve-cli
- Topics: markdown, pdf, converter, cli, golang

### 2. Configure Branch Protection
Visit: GitHub → Settings → Branches
- Create rule for `main`
- Require PR reviews: 1
- Dismiss stale reviews: Yes
- Require status checks: Yes
- Include administrators: Yes

### 3. Enable Discussions
Visit: GitHub → Settings → Features
- Check "Discussions"

### 4. Create Labels
Visit: GitHub → Issues → Labels
- Create: bug, enhancement, documentation, good first issue, help wanted, etc.

## Quick Command Reference

```bash
# View repository settings
gh repo view --json description,topics,hasDiscussions

# List open issues
gh issue list --state=open

# Create label
gh label create bug -c fc2929 -d "Something is not working"

# View branch protection
gh api repos/{owner}/{repo}/branches/main/protection

# List pull requests
gh pr list --state=open
```

## Next: Monitor and Maintain

Once setup is complete, you'll get:

1. **Automated Daily Checks** - repo-management.yml workflow
2. **Build Status Reporting** - CI/CD workflows
3. **GitHub Insights** - Traffic, activity, contributions
4. **Community Feedback** - Issues, discussions, PRs

Track success metrics:
- Stars and forks
- Clone statistics
- Issue quality and response time
- PR review feedback
- Documentation views

## Support Resources

- **GitHub Setup Guide**: docs/GITHUB_SETUP.md
- **API Reference**: docs/GITHUB_API.md
- **SEO Strategy**: docs/SEO.md
- **Contributing**: CONTRIBUTING.md
- **Release Process**: docs/RELEASE.md

## ✅ Checklist

- [ ] GitHub repository created
- [ ] Code pushed to main
- [ ] v0.2.0 tag pushed
- [ ] GITHUB_TOKEN environment variable set
- [ ] setup-repo.sh script executed
- [ ] Repository verified on GitHub
- [ ] Branch protection confirmed
- [ ] Topics visible on repository
- [ ] Issue templates displayed
- [ ] GitHub Discussions enabled

---

**You're all set!** Your repository is now professionally configured with:
- ✅ Branch protection and code quality rules
- ✅ SEO optimization for discoverability  
- ✅ Automated repository management
- ✅ Community contribution guidelines
- ✅ Security best practices
- ✅ Release management infrastructure

Start getting contributions and growing your open-source project! 🚀
