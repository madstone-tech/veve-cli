# GitHub API Management for veve-cli

Quick reference guide for using GitHub API with `$GITHUB_TOKEN` to manage the veve-cli repository.

## Setup

### 1. Create GitHub Personal Access Token

```bash
# Visit: https://github.com/settings/tokens/new

# Required scopes:
# - repo (full control of private/public repositories)
# - admin:repo_hook (write access to hooks)
# - admin:org_hook (organization hooks)
# - workflow (GitHub Actions workflows)
```

### 2. Set Environment Variables

```bash
# Add to ~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish
export GITHUB_TOKEN="ghp_your_token_here"
export GITHUB_OWNER="andhi"
export GITHUB_REPO="veve-cli"

# Reload shell
source ~/.bashrc  # or ~/.zshrc
```

### 3. Verify Token

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/user

# Should show your GitHub user info
```

---

## Common Tasks with GITHUB_TOKEN

### Update Repository Metadata

```bash
curl -X PATCH \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO \
  -d '{
    "description": "New description",
    "homepage": "https://github.com/madstone-tech/veve-cli",
    "topics": ["markdown", "pdf", "cli"]
  }'
```

### Add Repository Topics

```bash
curl -X PUT \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/topics \
  -d '{
    "names": [
      "markdown",
      "pdf",
      "converter",
      "cli",
      "golang",
      "pandoc",
      "themes"
    ]
  }'
```

### Create a Release

```bash
curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/releases \
  -d '{
    "tag_name": "v0.2.0",
    "target_commitish": "main",
    "name": "Release v0.2.0",
    "body": "See CHANGELOG.md for details",
    "draft": false,
    "prerelease": false
  }'
```

### Get Release Information

```bash
# List all releases
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/releases

# Get specific release
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/releases/tags/v0.2.0
```

### Create Issue Label

```bash
curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/labels \
  -d '{
    "name": "bug",
    "color": "fc2929",
    "description": "Something is not working"
  }'
```

### List Issue Labels

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/labels
```

### Enable GitHub Discussions

```bash
curl -X PATCH \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO \
  -d '{
    "has_discussions": true
  }'
```

### Configure Branch Protection

```bash
curl -X PUT \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/branches/main/protection \
  -d '{
    "required_status_checks": {
      "strict": true,
      "contexts": []
    },
    "required_pull_request_reviews": {
      "dismiss_stale_reviews": true,
      "required_approving_review_count": 1
    },
    "enforce_admins": true,
    "require_linear_history": true,
    "required_conversation_resolution": true
  }'
```

### Get Branch Protection Status

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/branches/main/protection
```

### Create a Deployment

```bash
curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/deployments \
  -d '{
    "ref": "main",
    "environment": "production",
    "description": "Deploying v0.2.0",
    "auto_merge": false
  }'
```

### List Repository Issues

```bash
# Open issues
curl -H "Authorization: token $GITHUB_TOKEN" \
  "https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/issues?state=open&per_page=100"

# With label filter
curl -H "Authorization: token $GITHUB_TOKEN" \
  "https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/issues?state=open&labels=bug"
```

### Create an Issue

```bash
curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/issues \
  -d '{
    "title": "Issue Title",
    "body": "Issue description",
    "labels": ["bug", "high priority"]
  }'
```

### List Pull Requests

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  "https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/pulls?state=open"
```

### Get Repository Statistics

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO | jq '{
    stargazers_count,
    forks_count,
    watchers_count,
    open_issues_count,
    language,
    created_at,
    updated_at
  }'
```

---

## Automated Setup Script

Run the automated setup script after pushing to GitHub:

```bash
export GITHUB_TOKEN="your_token"
export GITHUB_OWNER="andhi"
export GITHUB_REPO="veve-cli"

./.github/scripts/setup-repo.sh
```

This script automatically:
- Updates repository description
- Adds 10 topics for SEO
- Configures branch protection
- Sets merge strategies
- Creates issue labels
- Enables discussions

---

## GitHub CLI Alternative

Instead of raw API calls, use the `gh` CLI:

```bash
# Install
brew install gh

# Authenticate
gh auth login

# Update repository
gh repo edit --description "New description"

# Add topics
gh repo edit --add-topic markdown --add-topic pdf

# Create release
gh release create v0.2.0 -F CHANGELOG.md

# List issues
gh issue list

# Create issue
gh issue create --title "Bug" --body "Description"

# Create label
gh label create bug -c fc2929 -d "Something is not working"
```

---

## Security Best Practices

### Protect Your Token

```bash
# ❌ Never commit token to git
# ❌ Never hardcode in scripts
# ✅ Use environment variables
# ✅ Use GitHub Secrets in Actions
# ✅ Rotate token regularly
# ✅ Use minimal required scopes
```

### Token Scopes Explained

| Scope | Use Case |
|-------|----------|
| `repo` | Full access to repositories |
| `admin:repo_hook` | Manage webhooks |
| `admin:org_hook` | Organization webhooks |
| `workflow` | GitHub Actions workflows |
| `gist` | Manage gists |
| `read:packages` | Read packages |
| `write:packages` | Publish packages |

### Revoke Token If Compromised

```bash
# Visit: https://github.com/settings/tokens
# Click "Delete" on compromised token
# Create a new one with same scopes
```

---

## Useful curl Flags

```bash
# Pretty print JSON response
curl ... | jq '.'

# Save response to file
curl ... -o response.json

# Show response headers
curl -i ...

# Verbose output (shows request/response)
curl -v ...

# Set timeout (30 seconds)
curl --max-time 30 ...

# Follow redirects
curl -L ...
```

---

## Error Handling

### Common Errors

**401 Unauthorized:**
- Token is invalid or expired
- Check `echo $GITHUB_TOKEN`
- Regenerate token if needed

**403 Forbidden:**
- Token doesn't have required scopes
- Need to add scopes to token
- Check token permissions: https://github.com/settings/tokens

**404 Not Found:**
- Repository doesn't exist
- Check `$GITHUB_OWNER` and `$GITHUB_REPO`
- Verify you have access

**422 Unprocessable Entity:**
- Invalid request body
- Check JSON syntax
- Verify field names and types

### Check Token Permissions

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/user/repos \
  | jq '.[0] | {name, owner, permissions}'
```

---

## Rate Limiting

GitHub API has rate limits:

```bash
# Check rate limit
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/rate_limit

# Shows:
# - limit: 5000 requests/hour
# - remaining: Current remaining
# - reset: When limit resets
```

### Avoid Rate Limits

- Batch operations when possible
- Use conditional requests (ETag)
- Cache responses locally
- Space out requests

---

## Testing API Calls

### Test Token

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/user
```

### Test Repository Access

```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO
```

### Dry Run (without -X POST/PATCH)

```bash
# Test without actually creating
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO/labels | jq 'length'
```

---

## Resources

- [GitHub REST API Documentation](https://docs.github.com/en/rest)
- [GitHub CLI Documentation](https://cli.github.com/manual)
- [GitHub API Rate Limiting](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api)
- [GitHub GraphQL API](https://docs.github.com/en/graphql)
- [curl Documentation](https://curl.se/)
- [jq Documentation](https://stedolan.github.io/jq/)

---

## Quick Setup Command

```bash
# All-in-one setup (after setting env vars)
echo "Setting up repository..."
./.github/scripts/setup-repo.sh

# Verify setup
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/$GITHUB_OWNER/$GITHUB_REPO | jq '{description, topics, has_discussions}'
```

---

By using the GitHub API with `$GITHUB_TOKEN`, you can fully automate 
repository management and keep your veve-cli project professionally configured!
