#!/usr/bin/env bash

# install-completion.sh - Install shell completion for veve CLI
#
# This script installs shell completion for veve in the user's shell configuration.
# It automatically detects the user's primary shell and installs to the appropriate location.
#
# Usage:
#   ./scripts/install-completion.sh [bash|zsh|fish]
#
# If no shell is specified, auto-detects the current shell.

set -e

SHELL_TYPE=${1:-}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the veve binary path
VEVE_BIN="${VEVE_BIN:-$(command -v veve || echo './veve')}"

if [[ ! -f "$VEVE_BIN" && ! -x "$VEVE_BIN" ]]; then
    echo -e "${RED}Error: veve binary not found at $VEVE_BIN${NC}"
    echo "Please ensure veve is installed and in your PATH, or set VEVE_BIN environment variable."
    exit 1
fi

# Detect shell if not specified
if [[ -z "$SHELL_TYPE" ]]; then
    # Get the shell from SHELL environment variable
    CURRENT_SHELL=$(basename "$SHELL")
    case "$CURRENT_SHELL" in
        bash) SHELL_TYPE="bash" ;;
        zsh) SHELL_TYPE="zsh" ;;
        fish) SHELL_TYPE="fish" ;;
        *)
            echo -e "${YELLOW}Unable to auto-detect shell. Please specify: bash, zsh, or fish${NC}"
            echo "Usage: $0 [bash|zsh|fish]"
            exit 1
            ;;
    esac
    echo -e "${GREEN}Detected shell: $SHELL_TYPE${NC}"
fi

# Function to install bash completion
install_bash_completion() {
    echo "Installing bash completion..."
    
    # Determine bash completion directory
    BASH_COMPLETION_DIR=""
    
    # Check common locations
    if [[ -d "/usr/local/etc/bash_completion.d" ]]; then
        BASH_COMPLETION_DIR="/usr/local/etc/bash_completion.d"
    elif [[ -d "/etc/bash_completion.d" ]]; then
        BASH_COMPLETION_DIR="/etc/bash_completion.d"
    elif [[ -d "$HOME/.bash_completion.d" ]]; then
        BASH_COMPLETION_DIR="$HOME/.bash_completion.d"
    else
        # Use bashrc for inline completion
        BASH_COMPLETION_DIR=""
    fi
    
    if [[ -n "$BASH_COMPLETION_DIR" && -w "$BASH_COMPLETION_DIR" ]]; then
        # Install to system/user completion directory
        echo "Installing to $BASH_COMPLETION_DIR..."
        "$VEVE_BIN" completion bash > "$BASH_COMPLETION_DIR/veve"
        chmod 644 "$BASH_COMPLETION_DIR/veve"
        echo -e "${GREEN}  ✓ Installed to: $BASH_COMPLETION_DIR/veve${NC}"
    else
        # Install to bashrc
        BASHRC_FILE="$HOME/.bashrc"
        if [[ ! -f "$BASHRC_FILE" ]]; then
            BASHRC_FILE="$HOME/.bash_profile"
        fi
        
        if [[ ! -f "$BASHRC_FILE" ]]; then
            echo -e "${YELLOW}  Warning: Neither ~/.bashrc nor ~/.bash_profile found${NC}"
            echo "  Creating ~/.bashrc..."
            touch "$BASHRC_FILE"
        fi
        
        # Check if already installed
        if grep -q "veve completion bash" "$BASHRC_FILE" 2>/dev/null; then
            echo -e "${YELLOW}  Completion already installed in $BASHRC_FILE${NC}"
        else
            echo "" >> "$BASHRC_FILE"
            echo "# veve shell completion" >> "$BASHRC_FILE"
            echo "eval \"\$(veve completion bash)\"" >> "$BASHRC_FILE"
            echo -e "${GREEN}  ✓ Added to: $BASHRC_FILE${NC}"
        fi
    fi
}

# Function to install zsh completion
install_zsh_completion() {
    echo "Installing zsh completion..."
    
    # Determine zsh completion directory
    ZSH_COMPLETION_DIR=""
    
    # Check for oh-my-zsh custom completions
    if [[ -d "$HOME/.oh-my-zsh/custom/completions" ]]; then
        ZSH_COMPLETION_DIR="$HOME/.oh-my-zsh/custom/completions"
    elif [[ -d "$HOME/.zsh/completions" ]]; then
        ZSH_COMPLETION_DIR="$HOME/.zsh/completions"
    else
        # Fall back to fpath[1] or create ~/.zsh/completions
        ZSH_COMPLETION_DIR="$HOME/.zsh/completions"
        mkdir -p "$ZSH_COMPLETION_DIR"
    fi
    
    if [[ -w "$ZSH_COMPLETION_DIR" ]]; then
        echo "Installing to $ZSH_COMPLETION_DIR..."
        "$VEVE_BIN" completion zsh > "$ZSH_COMPLETION_DIR/_veve"
        chmod 644 "$ZSH_COMPLETION_DIR/_veve"
        echo -e "${GREEN}  ✓ Installed to: $ZSH_COMPLETION_DIR/_veve${NC}"
        
        # Add to fpath if not already there
        ZSHRC_FILE="$HOME/.zshrc"
        if [[ -f "$ZSHRC_FILE" ]]; then
            if ! grep -q "fpath.*$ZSH_COMPLETION_DIR" "$ZSHRC_FILE" 2>/dev/null; then
                echo "" >> "$ZSHRC_FILE"
                echo "# veve completions" >> "$ZSHRC_FILE"
                echo "fpath=($ZSH_COMPLETION_DIR \$fpath)" >> "$ZSHRC_FILE"
                echo -e "${GREEN}  ✓ Added fpath to: $ZSHRC_FILE${NC}"
            fi
        fi
    else
        echo -e "${RED}  Error: Cannot write to $ZSH_COMPLETION_DIR${NC}"
        exit 1
    fi
}

# Function to install fish completion
install_fish_completion() {
    echo "Installing fish completion..."
    
    FISH_COMPLETION_DIR="$HOME/.config/fish/completions"
    mkdir -p "$FISH_COMPLETION_DIR"
    
    if [[ -w "$FISH_COMPLETION_DIR" ]]; then
        echo "Installing to $FISH_COMPLETION_DIR..."
        "$VEVE_BIN" completion fish > "$FISH_COMPLETION_DIR/veve.fish"
        chmod 644 "$FISH_COMPLETION_DIR/veve.fish"
        echo -e "${GREEN}  ✓ Installed to: $FISH_COMPLETION_DIR/veve.fish${NC}"
    else
        echo -e "${RED}  Error: Cannot write to $FISH_COMPLETION_DIR${NC}"
        exit 1
    fi
}

# Install completion based on shell type
case "$SHELL_TYPE" in
    bash)
        install_bash_completion
        ;;
    zsh)
        install_zsh_completion
        ;;
    fish)
        install_fish_completion
        ;;
    *)
        echo -e "${RED}Error: Unknown shell '$SHELL_TYPE'${NC}"
        echo "Supported shells: bash, zsh, fish"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo "To use completions immediately, reload your shell configuration:"
case "$SHELL_TYPE" in
    bash) echo "  source ~/.bashrc" ;;
    zsh) echo "  source ~/.zshrc" ;;
    fish) echo "  source ~/.config/fish/config.fish" ;;
esac
