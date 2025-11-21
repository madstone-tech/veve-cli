# veve-cli Release Guide

This guide explains how to create releases for veve-cli using GoReleaser and GitHub Actions.

## Prerequisites

- Go 1.20 or later
- Git
- GitHub repository with write access
- (Optional) Homebrew tap repository for distribution

## Release Process

### 1. Prepare the Release

Before releasing, ensure:

1. All changes are committed and pushed to main
2. Tests pass: `go test ./...`
3. Code is formatted: `gofmt -w ./...`
4. Linting passes: `golangci-lint run ./...`
5. CHANGELOG is updated with release notes (optional but recommended)

### 2. Create a Git Tag

Create a semantic version tag following [Semantic Versioning](https://semver.org/):

```bash
# Example: Creating version 0.2.0
git tag -a v0.2.0 -m "Release v0.2.0: Add new features"

# Push the tag to GitHub
git push origin v0.2.0
```

Tags should follow the pattern `v*.*.* ` (e.g., `v0.1.0`, `v1.2.3`).

### 3. GitHub Actions Automatically Creates Release

When you push a version tag:

1. GitHub Actions automatically detects the tag
2. Triggers the `.github/workflows/release.yml` workflow
3. GoReleaser builds binaries for all platforms:
   - macOS (amd64, arm64)
   - Linux (amd64, arm64)
   - Windows (amd64)
4. Creates a GitHub Release with:
   - Pre-built binaries
   - Checksums (SHA256)
   - Changelog based on commit history
   - Homebrew formula (if configured)

### 4. Monitor the Release

You can monitor the release process:

1. Go to your GitHub repository
2. Click "Actions" tab
3. Look for the "Release" workflow run
4. Once complete, the release appears in "Releases" section

## Release Configuration Files

### `.goreleaser.yaml`

Main configuration file that defines:

- **Builds**: Platforms and architectures to build
- **Archives**: How to package binaries
- **Checksums**: SHA256 verification
- **Changelog**: Automatic changelog generation from commits
- **Release Notes**: GitHub Release template
- **Homebrew**: macOS installation via Homebrew

Key features:

```yaml
builds:
  - goos: [linux, darwin, windows]  # Operating systems
    goarch: [amd64, arm64]          # CPU architectures
    ignore:                          # Skip unsupported combinations
      - goos: windows
        goarch: arm64

universal_binaries:                 # Universal macOS binaries
  - replace: true

before:                             # Pre-build hooks
  hooks:
    - go mod tidy
    - ./scripts/generate-completions.sh

archives:                           # Packaging format
  - format: tar.gz                  # tar.gz for Unix
    format_overrides:               # ZIP for Windows
      - goos: windows
        format: zip
```

### `.github/workflows/release.yml`

Triggers on version tags and runs GoReleaser:

```yaml
on:
  push:
    tags:
      - 'v*.*.*'  # Only trigger on version tags
```

## Local Testing

### Test Release Build Locally

```bash
# Build snapshot (test build without publishing)
goreleaser build --snapshot

# This creates:
# - dist/veve_Linux_x86_64/veve
# - dist/veve_Darwin_x86_64/veve
# - dist/veve_Darwin_arm64/veve
# - dist/veve_Windows_x86_64/veve.exe
# - etc.
```

### Dry Run Release

```bash
# Test release process without publishing
goreleaser release --snapshot --clean

# Creates distribution packages in dist/ without uploading
```

## Installation Methods After Release

### From Pre-built Binaries

Users can download pre-built binaries from the GitHub Releases page:

```bash
# macOS (Intel)
curl -L https://github.com/madstone-tech/veve-cli/releases/download/v0.2.0/veve_Darwin_x86_64.tar.gz | tar xz
sudo mv veve /usr/local/bin/

# macOS (Apple Silicon)
curl -L https://github.com/madstone-tech/veve-cli/releases/download/v0.2.0/veve_Darwin_arm64.tar.gz | tar xz
sudo mv veve /usr/local/bin/

# Linux (AMD64)
curl -L https://github.com/madstone-tech/veve-cli/releases/download/v0.2.0/veve_Linux_x86_64.tar.gz | tar xz
sudo mv veve /usr/local/bin/

# Linux (ARM64)
curl -L https://github.com/madstone-tech/veve-cli/releases/download/v0.2.0/veve_Linux_arm64.tar.gz | tar xz
sudo mv veve /usr/local/bin/

# Windows (Download .zip from releases page)
```

### From Go

```bash
# Install specific version
go install github.com/madstone-tech/veve-cli/cmd/veve@v0.2.0

# Install latest
go install github.com/madstone-tech/veve-cli/cmd/veve@latest
```

### From Homebrew (Optional)

Once a Homebrew tap is configured:

```bash
brew tap andhi/tap
brew install veve
```

## Shell Completions in Releases

Release artifacts include generated shell completions:

```bash
# Extract release
tar -xzf veve_Linux_x86_64.tar.gz

# Install completions
./install-completion.sh  # Auto-detects current shell
./install-completion.sh bash zsh fish  # Install for multiple shells
```

## Verifying Downloads

Each release includes SHA256 checksums:

```bash
# Download binary
curl -O https://github.com/madstone-tech/veve-cli/releases/download/v0.2.0/veve_Linux_x86_64.tar.gz

# Download checksums
curl -O https://github.com/madstone-tech/veve-cli/releases/download/v0.2.0/checksums.txt

# Verify
sha256sum -c checksums.txt
```

## Troubleshooting

### Release workflow fails

1. Check GitHub Actions logs: Repository → Actions → Release workflow
2. Common issues:
   - Missing `GITHUB_TOKEN` (should be automatic)
   - Tag format incorrect (must be `v*.*.* `)
   - Go version mismatch (must be 1.20+)

### Binary doesn't work on target system

1. Verify architecture: `uname -m` (amd64/arm64)
2. Verify OS: `uname -s` (Darwin/Linux/Windows)
3. Check dependencies: `veve --version` should work

### Homebrew installation issues

1. Ensure `HOMEBREW_TAP_TOKEN` is configured in GitHub secrets
2. Verify Homebrew tap repository exists
3. Check formula syntax in `.goreleaser.yaml`

## Best Practices

1. **Use Semantic Versioning**: Follow SemVer for version numbers
2. **Write Good Commit Messages**: Changelog is auto-generated from commits
3. **Test Before Release**: Run full test suite locally
4. **Update Documentation**: Keep README and guides current
5. **Create Release Notes**: Add release highlights in GitHub UI
6. **Announce Release**: Share on social media, forums, etc.

## Common Release Tasks

### Patch Release (bug fixes)

```bash
git tag -a v0.2.1 -m "Release v0.2.1: Bug fixes"
git push origin v0.2.1
```

### Minor Release (new features)

```bash
git tag -a v0.3.0 -m "Release v0.3.0: New theme features"
git push origin v0.3.0
```

### Major Release (breaking changes)

```bash
git tag -a v1.0.0 -m "Release v1.0.0: Production ready"
git push origin v1.0.0
```

## Automatic Version Bumping

To automate version bumping, you can use tools like:

- `bump2version` - Bumps version in files and creates tags
- `standard-version` - Follows Conventional Commits
- Manual tagging (recommended for now)

## Release Channels

Consider supporting multiple release channels:

1. **Stable** (`v*.*.* `) - Production releases
2. **Beta** (`v*.*.*-beta`) - Pre-release testing
3. **Alpha** (`v*.*.*-alpha`) - Experimental features

GoReleaser automatically handles prerelease status based on version tag.

## Next Steps

After releasing:

1. Test installation on different platforms
2. Create GitHub Release notes with highlights
3. Update official website with new version
4. Announce in project channels (Discord, mailing list, etc.)
5. Update package managers (Homebrew, apt, etc.)
6. Monitor for bug reports and issues

---

For more information, see:

- [GoReleaser Documentation](https://goreleaser.com)
- [GitHub Actions Workflows](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)
- [Conventional Commits](https://www.conventionalcommits.org/)
