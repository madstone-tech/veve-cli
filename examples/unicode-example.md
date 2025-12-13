# Unicode & Emoji Support Examples

This document demonstrates veve's comprehensive unicode and emoji support. It includes various character sets and special symbols that require unicode-capable PDF rendering engines.

## Emoji Examples

### Standard Emoji
🎉 🚀 📄 ✅ ❤️ 🌟 💡 🔧 🎨 📊

### Emoji with Skin Tone Modifiers
👋 👋🏻 👋🏼 👋🏽 👋🏾 👋🏿

### Family Emoji (Zero-Width Joiner Sequences)
👨‍👩‍👧‍👦 👨‍💻 👩‍🔬 👨‍🎓 👩‍🎨

### Nature & Weather Emoji
🌈 🌸 🌺 🌻 🌷 🌹 ⛅ 🌤️ ⛈️ 🌊

## CJK (Chinese, Japanese, Korean) Characters

### Chinese (Simplified)
世界 中国 北京 上海 学习 理解 美好 希望

### Chinese (Traditional)
世界 中國 北京 上海 學習 理解 美好 希望

### Japanese (Hiragana)
あいうえお かきくけこ さしすせそ たちつてと

### Japanese (Katakana)
アイウエオ カキクケコ サシスセソ タチツテト

### Japanese (Kanji)
日本 東京 京都 学校 会社 友達 家族 先生

### Korean (Hangul)
한국 서울 안녕하세요 감사합니다 좋은 아침입니다

## Mathematical Symbols

### Basic Operators
+ − × ÷ = ≠ < > ≤ ≥

### Mathematical Functions
∑ ∏ ∫ √ ∛ ∜ ± ∓ ∞

### Set Theory
∈ ∉ ⊂ ⊃ ⊆ ⊇ ∩ ∪ ∅

### Logic
∧ ∨ ¬ ⟹ ⟺ ∀ ∃ ∄

### Greek Letters
α β γ δ ε ζ η θ ι κ λ μ ν ξ ο π ρ σ τ υ φ χ ψ ω
Α Β Γ Δ Ε Ζ Η Θ Ι Κ Λ Μ Ν Ξ Ο Π Ρ Σ Τ Υ Φ Χ Ψ Ω

## Latin Extended Characters & Diacritics

### Accented Vowels
à á â ã ä å  è é ê ë  ì í î ï  ò ó ô õ ö  ù ú û ü

### Accented Consonants
ç ñ š ž đ ĝ ĥ ĵ ķ ł ń ř

### Ligatures & Special Letters
æ œ ß ð þ ħ

## Examples in Context

### A Mathematical Expression
The quadratic formula is: x = (−b ± √(b² − 4ac)) / 2a

### A Multilingual Sentence
Hello (English) •你好 (Chinese) • こんにちは (Japanese) • 안녕하세요 (Korean) • مرحبا (Arabic)

### A Scientific Notation
The Planck constant: h ≈ 6.626 × 10⁻³⁴ J·s
Avogadro's number: Nₐ ≈ 6.022 × 10²³ mol⁻¹

### Currency Symbols
$ £ € ¥ ₹ ₽ ₩ ₪ ₦ ₨ ₱ ₡ ₲ ₴ ₵

## Special Unicode Blocks

### Arrows
← → ↑ ↓ ↔ ↕ ⇐ ⇒ ⇑ ⇓ ⇔ ⇕

### Box Drawing
┌─┬─┐  ├─┼─┤  └─┴─┘

### Geometric Shapes
◉ ◎ ○ ◔ ◕ ◖ ◗ ◘ ◙ ◚ ◛
□ ◻ ▯ ◼ ◽ ▪ ▫
△ ▲ ▴ ▵ ▶ ▷ ▸ ▹ ► ▻
▼ ▽ ▾ ◀ ◁ ◂ ◃

### Stars & Symbols
★ ☆ ✦ ✧ ✥ ✤ ✢ ✣ ✡ ✢ ⭐ 🌟

## Usage Notes

To convert this markdown file to PDF with automatic unicode support:

```bash
# Using default engine (automatically selected for unicode content)
veve convert unicode-example.md -o unicode-example.pdf

# Using a specific unicode-capable engine
veve convert --engine xelatex unicode-example.md -o unicode-example.pdf

# With verbose output to see which engine was selected
veve convert --verbose unicode-example.md -o unicode-example.pdf
```

## Supported Engines

veve automatically selects from these unicode-capable engines:

1. **XeLaTeX** (Recommended) - Excellent unicode support, native unicode fonts
2. **LuaLaTeX** - Full unicode support with Lua scripting
3. **WeasyPrint** - HTML/CSS to PDF, good unicode support
4. **Prince** - Commercial engine with excellent unicode handling

If none of these engines are installed, veve will provide installation instructions for your platform.

## Platform-Specific Notes

### macOS
Install via Homebrew:
```bash
brew install mactex  # For XeLaTeX and LuaLaTeX
brew install weasyprint  # Alternative
```

### Linux (Ubuntu/Debian)
```bash
sudo apt-get install texlive-xetex  # For XeLaTeX
sudo apt-get install weasyprint  # Alternative
```

### Windows
- Download and install MiKTeX from https://miktex.org/
- Or use WSL2 with Linux instructions
