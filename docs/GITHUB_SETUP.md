# GitHub Repository Setup Guide

This guide provides best practices for setting up your veve-cli GitHub repository with branch protections, SEO optimization, and automated management using the GitHub API.

## Table of Contents

1. [Initial Repository Setup](#initial-repository-setup)
2. [Branch Protection Rules](#branch-protection-rules)
3. [Repository Metadata & SEO](#repository-metadata--seo)
4. [GitHub API Setup](#github-api-setup)
5. [Automation Scripts](#automation-scripts)

---

## Initial Repository Setup

### Prerequisites

- GitHub account with repository created
- `GITHUB_TOKEN` with `admin:repo_hook`, `repo` scopes
- `gh` CLI installed: `brew install gh`

### 1. Authenticate with GitHub

```bash
gh auth login
# Choose: GitHub.com
# Choose: HTTPS
# Authenticate with your browser
```

### 2. Clone and Configure

```bash
git clone https://github.com/USERNAME/veve-cli.git
cd veve-cli
git remote add origin https://github.com/USERNAME/veve-cli.git
git branch -M main  # Rename default branch to main
git push -u origin main
```

---

## Branch Protection Rules

Branch protection rules ensure code quality and prevent accidental commits to important branches.

### Main Branch Protection

Protect the `main` branch with these rules:

```bash
# Using GitHub CLI
gh repo edit --enable-auto-merge --enable-squash-merge

# Or via API (see scripts below)
```

### Best Practice Rules for Main Branch

1. **Require Status Checks to Pass**
   - CI tests must pass
   - Code coverage must be maintained
   - Linting must pass

2. **Require Code Reviews**
   - Minimum 1 approval required
   - Dismiss stale pull request approvals
   - Require review from code owners

3. **Require Branches to be Up to Date**
   - Branch must be up to date before merging
   - Prevents conflicts and ensures all checks pass

4. **Require Signed Commits**
   - All commits must be signed
   - Increases security and trust

5. **Include Administrators**
   - Rules apply to admins too
   - Prevents accidental bypasses

### Apply Rules with Script

Create `.github/scripts/setup-branch-protection.sh`:

```bash
#!/bin/bash

OWNER="USERNAME"
REPO="veve-cli"
BRANCH="main"
TOKEN="$GITHUB_TOKEN"

# Enable branch protection
curl -X PUT \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO/branches/$BRANCH/protection \
  -d '{
    "required_status_checks": {
      "strict": true,
      "contexts": ["ci/build", "ci/test", "ci/lint"]
    },
    "required_pull_request_reviews": {
      "dismiss_stale_reviews": true,
      "require_code_owner_reviews": true,
      "required_approving_review_count": 1
    },
    "enforce_admins": true,
    "require_linear_history": true,
    "required_conversation_resolution": true,
    "restrictions": null
  }'

echo "Branch protection rules applied to $BRANCH"
```

### Develop Branch Protection

For feature development, protect `develop` with slightly relaxed rules:

```bash
curl -X PUT \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO/branches/develop/protection \
  -d '{
    "required_status_checks": {
      "strict": true,
      "contexts": ["ci/build", "ci/test"]
    },
    "required_pull_request_reviews": {
      "required_approving_review_count": 1
    },
    "enforce_admins": false
  }'
```

---

## Repository Metadata & SEO

### 1. Repository Description

**About Section** (GitHub UI):

```
📝 Fast, themeable markdown-to-PDF converter built with Go

Convert your markdown files to beautiful PDFs with built-in themes
(default, dark, academic) and custom CSS styling. Perfect for
documentation, reports, technical writing, and more.

Features:
• Cross-platform (macOS, Linux, Windows)
• 3+ built-in professional themes
• Custom theme support with CSS
• stdin/stdout piping for Unix integration
• Shell completions (bash, zsh, fish)
• Production-ready with 100+ tests

Get started: https://github.com/andhi/veve-cli#installation
```

### 2. Topics for Discoverability

Add these topics to improve SEO:

```bash
# Using GitHub CLI
gh repo edit --add-topic markdown
gh repo edit --add-topic pdf
gh repo edit --add-topic converter
gh repo edit --add-topic cli
gh repo edit --add-topic golang
gh repo edit --add-topic pandoc
gh repo edit --add-topic themes
gh repo edit --add-topic pdf-generation
gh repo edit --add-topic command-line-tool
gh repo edit --add-topic document-conversion
```

Topics list:
- `markdown`
- `pdf`
- `converter`
- `cli`
- `golang`
- `pandoc`
- `themes`
- `pdf-generation`
- `command-line-tool`
- `document-conversion`
- `theming`
- `css`

### 3. Homepage and Documentation Links

Set in repository settings:

**Website**: `https://github.com/andhi/veve-cli` (or custom domain)

**Documentation**: Links in README.md and sidebar

### 4. Social Media Preview

Create `.github/GITHUB_PROFILE.md`:

```markdown
# veve-cli

**Fast, themeable markdown-to-PDF converter**

Convert markdown to beautiful PDFs with professional themes and Unix integration.

- 🚀 Cross-platform (macOS, Linux, Windows)
- 🎨 Built-in + custom themes
- ⚡ Fast conversion powered by Pandoc
- 🔧 Unix composable (stdin/stdout)
- 📚 100% documented
- ✅ Production-ready

[📖 Documentation](https://github.com/andhi/veve-cli#readme)
[🚀 Quick Start](https://github.com/andhi/veve-cli#installation)
[📝 Contributing](CONTRIBUTING.md)
```

---

## GitHub API Setup

### Set Environment Variables

```bash
# Save to ~/.bashrc, ~/.zshrc, or similar
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_OWNER="andhi"
export GITHUB_REPO="veve-cli"
```

### Verify Token Permissions

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/user

# Should return your user info
```

### API Endpoints Reference

**Repository Info:**
```bash
GET /repos/$OWNER/$REPO
```

**Update Repository:**
```bash
PATCH /repos/$OWNER/$REPO
```

**Branch Protection:**
```bash
PUT /repos/$OWNER/$REPO/branches/$BRANCH/protection
```

**Topics:**
```bash
PUT /repos/$OWNER/$REPO/topics
```

---

## Automation Scripts

### 1. Complete Setup Script

Create `.github/scripts/setup-repo.sh`:

```bash
#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

OWNER="${GITHUB_OWNER:-andhi}"
REPO="${GITHUB_REPO:-veve-cli}"
TOKEN="$GITHUB_TOKEN"

if [ -z "$TOKEN" ]; then
  echo "Error: GITHUB_TOKEN not set"
  exit 1
fi

echo -e "${BLUE}Setting up GitHub repository...${NC}"

# 1. Update repository metadata
echo -e "${BLUE}1. Updating repository metadata...${NC}"
curl -s -X PATCH \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO \
  -d '{
    "description": "Fast, themeable markdown-to-PDF converter built with Go",
    "homepage": "https://github.com/'"$OWNER"'/'"$REPO"'",
    "topics": ["markdown", "pdf", "converter", "cli", "golang", "pandoc", "themes"],
    "has_wiki": false,
    "has_projects": true,
    "has_downloads": true
  }' > /dev/null

echo -e "${GREEN}✓ Repository metadata updated${NC}"

# 2. Setup branch protection for main
echo -e "${BLUE}2. Setting up branch protection for main...${NC}"
curl -s -X PUT \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO/branches/main/protection \
  -d '{
    "required_status_checks": {
      "strict": true,
      "contexts": ["ci/build", "ci/test"]
    },
    "required_pull_request_reviews": {
      "dismiss_stale_reviews": true,
      "require_code_owner_reviews": false,
      "required_approving_review_count": 1
    },
    "enforce_admins": true,
    "require_linear_history": true,
    "required_conversation_resolution": true,
    "restrictions": null
  }' > /dev/null

echo -e "${GREEN}✓ Branch protection configured${NC}"

# 3. Enable auto-merge
echo -e "${BLUE}3. Enabling auto-merge...${NC}"
curl -s -X PATCH \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO \
  -d '{
    "allow_squash_merge": true,
    "allow_merge_commit": false,
    "allow_rebase_merge": true,
    "delete_branch_on_merge": true
  }' > /dev/null

echo -e "${GREEN}✓ Auto-merge configured${NC}"

# 4. Create default labels
echo -e "${BLUE}4. Creating issue labels...${NC}"
declare -A labels=(
  ["bug"]="fc2929"
  ["enhancement"]="84b6eb"
  ["documentation"]="0075ca"
  ["good first issue"]="7057ff"
  ["help wanted"]="008672"
  ["question"]="d876e3"
  ["wontfix"]="ffffff"
)

for label in "${!labels[@]}"; do
  curl -s -X POST \
    -H "Accept: application/vnd.github.v3+json" \
    -H "Authorization: token $TOKEN" \
    https://api.github.com/repos/$OWNER/$REPO/labels \
    -d '{
      "name": "'"$label"'",
      "color": "'"${labels[$label]}"'",
      "description": "'"$label"'"
    }' 2>/dev/null || true
done

echo -e "${GREEN}✓ Labels created${NC}"

echo -e "${GREEN}✓ Repository setup complete!${NC}"
```

### 2. Run Setup Script

```bash
chmod +x .github/scripts/setup-repo.sh
export GITHUB_TOKEN="your_token"
export GITHUB_OWNER="andhi"
export GITHUB_REPO="veve-cli"

./.github/scripts/setup-repo.sh
```

### 3. GitHub Actions for Maintenance

Create `.github/workflows/repo-management.yml`:

```yaml
name: Repository Management

on:
  schedule:
    # Daily at 2 AM UTC
    - cron: '0 2 * * *'
  workflow_dispatch:

jobs:
  repo-health:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Check repository health
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # Check for stale issues
          curl -s -H "Authorization: token $GITHUB_TOKEN" \
            "https://api.github.com/repos/${{ github.repository }}/issues?state=open&sort=updated&direction=asc" \
            | jq '.[] | select(.updated_at < now - 30*24*3600 | @base64d)'
          
          # Check PR status
          curl -s -H "Authorization: token $GITHUB_TOKEN" \
            "https://api.github.com/repos/${{ github.repository }}/pulls?state=open" \
            | jq 'length'
```

---

## SEO Best Practices

### 1. README Optimization

Structure your README for search engines:

```markdown
# veve-cli

**Fast, themeable markdown-to-PDF converter**

Convert your markdown files to beautiful PDFs with professional themes.
Built with Go, powered by Pandoc.

[Installation](#installation) | 
[Documentation](#documentation) | 
[Contributing](#contributing)

## Key Features

- Cross-platform (macOS, Linux, Windows)
- Built-in professional themes
- Custom CSS theme support
- stdin/stdout piping
- Shell completions

## Installation

## Usage

## Documentation

## Contributing
```

### 2. GitHub Discussions

Enable GitHub Discussions for SEO:

```bash
# Enable in repository settings
# This increases engagement and helps with SEO
```

### 3. Releases and Tags

Always create detailed release notes:

```bash
# When creating a release
gh release create v0.2.0 \
  --title "Release v0.2.0: Theme Management" \
  --notes "See CHANGELOG.md for details"
```

### 4. Semantic Keywords

Use in:
- Repository description
- README headings
- Commit messages
- Issue templates
- Release notes

Key phrases:
- "markdown to PDF converter"
- "cross-platform CLI tool"
- "themeable document conversion"
- "Pandoc wrapper"
- "Go command-line tool"

---

## Issue Templates

Create `.github/ISSUE_TEMPLATE/bug_report.md`:

```markdown
---
name: Bug Report
about: Report a bug
labels: bug
---

## Description
Clear description of the bug.

## Steps to Reproduce
1. 
2. 
3. 

## Expected Behavior
What should happen.

## Actual Behavior
What actually happens.

## System Information
- OS: [macOS/Linux/Windows]
- Version: [version]
- Pandoc Version: [version]
```

Create `.github/ISSUE_TEMPLATE/feature_request.md`:

```markdown
---
name: Feature Request
about: Suggest a feature
labels: enhancement
---

## Description
Clear description of the feature.

## Use Case
Why this feature is needed.

## Proposed Solution
How it might be implemented.
```

---

## Pull Request Template

Create `.github/pull_request_template.md`:

```markdown
## Description
Brief description of changes.

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation
- [ ] Code refactoring

## Related Issues
Closes #123

## Testing
How to verify the changes work.

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes
```

---

## Workflow for New Contributors

1. **Fork repository**
2. **Create feature branch**: `git checkout -b feature/name`
3. **Make changes** with tests
4. **Create pull request** using template
5. **Pass CI checks**
6. **Get 1 approval**
7. **Merge with squash** (configured)
8. **Branch auto-deletes** (configured)

---

## Monitoring and Maintenance

### Weekly Checks

```bash
# List open issues
gh issue list --state=open

# List open PRs
gh pr list --state=open

# Check recent activity
gh repo view --web
```

### Monthly Tasks

1. Review and close stale issues
2. Update dependencies
3. Check test coverage
4. Review security alerts
5. Update documentation

---

## Quick Reference

| Task | Command |
|------|---------|
| Add topics | `gh repo edit --add-topic markdown` |
| View settings | `gh repo view --json description` |
| Enable auto-merge | `gh repo edit --enable-auto-merge` |
| Create label | `gh label create bug -c fc2929` |
| List open issues | `gh issue list --state=open` |

---

## Resources

- [GitHub API Documentation](https://docs.github.com/en/rest)
- [GitHub CLI Reference](https://cli.github.com/manual)
- [Branch Protection Rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
- [GitHub SEO Best Practices](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features)

---

Once you push to GitHub, run the setup script to automatically configure all these settings!

