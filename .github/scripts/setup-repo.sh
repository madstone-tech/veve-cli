#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

# Configuration
OWNER="${GITHUB_OWNER:-andhi}"
REPO="${GITHUB_REPO:-veve-cli}"
TOKEN="$GITHUB_TOKEN"

# Validate token
if [ -z "$TOKEN" ]; then
  echo -e "${RED}Error: GITHUB_TOKEN environment variable not set${NC}"
  echo "Set it with: export GITHUB_TOKEN='your_token_here'"
  exit 1
fi

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  veve-cli GitHub Repository Setup      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo "Repository: $OWNER/$REPO"
echo ""

# 1. Update repository metadata
echo -e "${BLUE}→ Step 1: Updating repository metadata...${NC}"
curl -s -X PATCH \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO \
  -d '{
    "description": "Fast, themeable markdown-to-PDF converter built with Go. Convert your markdown files to beautiful PDFs with built-in themes and custom CSS styling.",
    "homepage": "https://github.com/'"$OWNER"'/'"$REPO"'#readme",
    "topics": ["markdown", "pdf", "converter", "cli", "golang", "pandoc", "themes", "pdf-generation", "command-line-tool", "document-conversion"],
    "has_wiki": false,
    "has_projects": true,
    "has_downloads": true,
    "is_template": false
  }' > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ Repository metadata updated${NC}"
else
  echo -e "${RED}✗ Failed to update metadata${NC}"
fi

echo ""

# 2. Setup branch protection for main
echo -e "${BLUE}→ Step 2: Setting up branch protection for main...${NC}"
curl -s -X PUT \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO/branches/main/protection \
  -d '{
    "required_status_checks": {
      "strict": true,
      "contexts": []
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
  }' > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ Branch protection configured for main${NC}"
else
  echo -e "${RED}✗ Failed to configure branch protection${NC}"
fi

echo ""

# 3. Configure merge strategies
echo -e "${BLUE}→ Step 3: Configuring merge strategies...${NC}"
curl -s -X PATCH \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO \
  -d '{
    "allow_squash_merge": true,
    "allow_merge_commit": true,
    "allow_rebase_merge": true,
    "allow_auto_merge": true,
    "delete_branch_on_merge": true
  }' > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ Merge strategies configured${NC}"
else
  echo -e "${RED}✗ Failed to configure merge strategies${NC}"
fi

echo ""

# 4. Create issue labels
echo -e "${BLUE}→ Step 4: Creating issue labels...${NC}"

declare -A labels=(
  ["bug"]="fc2929"
  ["enhancement"]="84b6eb"
  ["documentation"]="0075ca"
  ["good first issue"]="7057ff"
  ["help wanted"]="008672"
  ["question"]="d876e3"
  ["wontfix"]="ffffff"
  ["type: feature"]="a2eeef"
  ["type: refactor"]="fbca04"
  ["type: test"]="cccccc"
  ["priority: high"]="ff0000"
  ["priority: medium"]="ffff00"
  ["priority: low"]="0366d6"
)

labels_created=0
for label in "${!labels[@]}"; do
  curl -s -X POST \
    -H "Accept: application/vnd.github.v3+json" \
    -H "Authorization: token $TOKEN" \
    https://api.github.com/repos/$OWNER/$REPO/labels \
    -d '{
      "name": "'"$label"'",
      "color": "'"${labels[$label]}"'",
      "description": "'"$label"'"
    }' 2>/dev/null

  ((labels_created++))
done

echo -e "${GREEN}✓ Created $labels_created issue labels${NC}"

echo ""

# 5. Enable discussions
echo -e "${BLUE}→ Step 5: Enabling GitHub Discussions...${NC}"
curl -s -X PATCH \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: token $TOKEN" \
  https://api.github.com/repos/$OWNER/$REPO \
  -d '{
    "has_discussions": true
  }' > /dev/null 2>&1

if [ $? -eq 0 ]; then
  echo -e "${GREEN}✓ GitHub Discussions enabled${NC}"
else
  echo -e "${RED}✗ Failed to enable discussions${NC}"
fi

echo ""

# Summary
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}✓ Repository setup complete!${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo "Configuration applied:"
echo "  • Repository description updated"
echo "  • 10 topics added for SEO"
echo "  • Branch protection enabled on main"
echo "  • Merge strategies configured"
echo "  • $labels_created issue labels created"
echo "  • GitHub Discussions enabled"
echo ""
echo "Next steps:"
echo "  1. Push your code: git push origin main"
echo "  2. Create tags for releases: git push origin v0.2.0"
echo "  3. Monitor at: https://github.com/$OWNER/$REPO"
echo ""
