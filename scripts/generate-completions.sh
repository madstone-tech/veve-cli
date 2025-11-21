#!/usr/bin/env bash
set -e

# Generate shell completions for veve
echo "Generating shell completions..."

# Create completions directory
mkdir -p completions

# Build veve if not already built
if [ ! -f ./veve ]; then
    echo "Building veve..."
    go build -o veve ./cmd/veve
fi

# Generate completions using veve's completion commands
echo "  • Generating bash completion..."
./veve completion bash > completions/veve.bash || true

echo "  • Generating zsh completion..."
./veve completion zsh > completions/_veve || true

echo "  • Generating fish completion..."
./veve completion fish > completions/veve.fish || true

echo "✅ Completions generated in completions/"
ls -lh completions/
