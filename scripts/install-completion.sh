#!/usr/bin/env bash
#
# veve shell completion installation script
#
# This script installs shell completions for veve in the appropriate system directories.
# Supports bash, zsh, and fish shells.
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Determine which shells to install completions for
SHELLS_TO_INSTALL=()

if [ $# -eq 0 ]; then
    # No arguments - auto-detect current shell
    CURRENT_SHELL=$(basename "$SHELL")
    case "$CURRENT_SHELL" in
        bash) SHELLS_TO_INSTALL=("bash") ;;
        zsh) SHELLS_TO_INSTALL=("zsh") ;;
        fish) SHELLS_TO_INSTALL=("fish") ;;
        *)
            echo -e "${YELLOW}Unknown shell: $CURRENT_SHELL${NC}"
            echo "Please specify shells: $0 bash zsh fish"
            exit 1
            ;;
    esac
else
    # Process arguments
    for shell in "$@"; do
        case "$shell" in
            all)
                SHELLS_TO_INSTALL=("bash" "zsh" "fish")
                break
                ;;
            bash|zsh|fish)
                SHELLS_TO_INSTALL+=("$shell")
                ;;
            *)
                echo -e "${RED}Unknown shell: $shell${NC}"
                echo "Supported shells: bash, zsh, fish, all"
                exit 1
                ;;
        esac
    done
fi

# Check if completions directory exists
if [ ! -d "completions" ]; then
    echo -e "${YELLOW}Completions directory not found. Generating completions...${NC}"
    ./scripts/generate-completions.sh
fi

# Function to install bash completion
install_bash() {
    if [ ! -f "completions/veve.bash" ]; then
        echo -e "${RED}Bash completion file not found${NC}"
        return 1
    fi

    # Try different installation paths
    if [ -d "$HOME/.local/share/bash-completion/completions" ]; then
        cp "completions/veve.bash" "$HOME/.local/share/bash-completion/completions/veve"
        echo -e "${GREEN}✓ Bash completion installed to ~/.local/share/bash-completion/completions/veve${NC}"
    elif [ -d "/usr/local/etc/bash_completion.d" ]; then
        sudo cp "completions/veve.bash" "/usr/local/etc/bash_completion.d/veve"
        echo -e "${GREEN}✓ Bash completion installed to /usr/local/etc/bash_completion.d/veve${NC}"
    else
        mkdir -p "$HOME/.bash_completion.d"
        cp "completions/veve.bash" "$HOME/.bash_completion.d/veve"
        echo -e "${GREEN}✓ Bash completion installed to ~/.bash_completion.d/veve${NC}"
        echo -e "${YELLOW}  Add to ~/.bashrc: source ~/.bash_completion.d/veve${NC}"
    fi
}

# Function to install zsh completion
install_zsh() {
    if [ ! -f "completions/_veve" ]; then
        echo -e "${RED}Zsh completion file not found${NC}"
        return 1
    fi

    # Check for Oh My Zsh
    if [ -d "$HOME/.oh-my-zsh" ]; then
        if [ -d "$HOME/.oh-my-zsh/completions" ]; then
            cp "completions/_veve" "$HOME/.oh-my-zsh/completions/_veve"
            echo -e "${GREEN}✓ Zsh completion installed to ~/.oh-my-zsh/completions/_veve${NC}"
        else
            mkdir -p "$HOME/.oh-my-zsh/custom/plugins/veve"
            cp "completions/_veve" "$HOME/.oh-my-zsh/custom/plugins/veve/_veve"
            echo -e "${GREEN}✓ Zsh completion installed to ~/.oh-my-zsh/custom/plugins/veve/_veve${NC}"
            echo -e "${YELLOW}  Add to ~/.zshrc: plugins=(... veve)${NC}"
        fi
    else
        # Standard zsh setup
        mkdir -p "$HOME/.zfunc"
        cp "completions/_veve" "$HOME/.zfunc/_veve"
        echo -e "${GREEN}✓ Zsh completion installed to ~/.zfunc/_veve${NC}"
        
        # Check if .zshrc has fpath
        if ! grep -q 'fpath.*\.zfunc' "$HOME/.zshrc" 2>/dev/null; then
            echo -e "${YELLOW}  Add to ~/.zshrc:${NC}"
            echo "    fpath=(~/.zfunc \$fpath)"
            echo "    autoload -U compinit && compinit"
        fi
    fi
}

# Function to install fish completion
install_fish() {
    if [ ! -f "completions/veve.fish" ]; then
        echo -e "${RED}Fish completion file not found${NC}"
        return 1
    fi

    mkdir -p "$HOME/.config/fish/completions"
    cp "completions/veve.fish" "$HOME/.config/fish/completions/veve.fish"
    echo -e "${GREEN}✓ Fish completion installed to ~/.config/fish/completions/veve.fish${NC}"
}

# Install completions for selected shells
echo "Installing veve shell completions..."
echo ""

for shell in "${SHELLS_TO_INSTALL[@]}"; do
    case "$shell" in
        bash) install_bash ;;
        zsh) install_zsh ;;
        fish) install_fish ;;
    esac
done

echo ""
echo -e "${GREEN}✓ Installation complete!${NC}"
echo ""
echo "Testing completions:"
echo "  bash: type 'veve <TAB>'"
echo "  zsh:  type 'veve <TAB>'"
echo "  fish: type 'veve <TAB>'"
