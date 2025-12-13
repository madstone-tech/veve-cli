// Package engines provides PDF engine detection, validation, and selection logic.
package engines

import (
	"bytes"
	"unicode"
)

// DetectUnicodeInContent checks if markdown content contains non-ASCII unicode characters
// Returns true if content contains any characters outside ASCII range (0-127)
// This helps determine if we need a unicode-capable engine
func DetectUnicodeInContent(content string) bool {
	for _, r := range content {
		// Check if rune is outside ASCII range (0-127)
		if r > 127 {
			return true
		}
	}
	return false
}

// DetectUnicodeInBytes checks if byte content contains non-ASCII unicode
func DetectUnicodeInBytes(content []byte) bool {
	for _, b := range content {
		if b > 127 {
			return true
		}
	}
	return false
}

// ContainsEmoji checks if content contains emoji characters
// Emoji are in unicode ranges:
// - U+1F300-U+1F9FF (Emoticons, Symbols, Pictographs)
// - U+2600-U+27BF (Miscellaneous Symbols)
// - U+1F680-U+1F6FF (Transport and Map Symbols)
func ContainsEmoji(content string) bool {
	for _, r := range content {
		// Emoji ranges
		if (r >= 0x1F300 && r <= 0x1F9FF) || // Emoticons, Symbols, Pictographs
			(r >= 0x2600 && r <= 0x27BF) || // Miscellaneous Symbols
			(r >= 0x1F680 && r <= 0x1F6FF) { // Transport Symbols
			return true
		}
	}
	return false
}

// ContainsCJK checks if content contains Chinese, Japanese, or Korean characters
// CJK Unified Ideographs: U+4E00-U+9FFF
// CJK Compatibility Ideographs: U+F900-U+FAFF
// Hiragana: U+3040-U+309F
// Katakana: U+30A0-U+30FF
// Hangul Syllables: U+AC00-U+D7AF
func ContainsCJK(content string) bool {
	for _, r := range content {
		if (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
			(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
			(r >= 0x3040 && r <= 0x309F) || // Hiragana
			(r >= 0x30A0 && r <= 0x30FF) || // Katakana
			(r >= 0xAC00 && r <= 0xD7AF) { // Hangul Syllables
			return true
		}
	}
	return false
}

// ContainsMathSymbols checks if content contains mathematical symbols
// Mathematical Alphanumeric Symbols: U+1D400-U+1D7FF
// Mathematical Operators: U+2200-U+22FF
// Arrows: U+2190-U+21FF
// Miscellaneous Technical: U+2300-U+243F
func ContainsMathSymbols(content string) bool {
	for _, r := range content {
		if (r >= 0x2200 && r <= 0x22FF) || // Mathematical Operators
			(r >= 0x2190 && r <= 0x21FF) || // Arrows
			(r >= 0x2300 && r <= 0x243F) || // Miscellaneous Technical
			(r >= 0x1D400 && r <= 0x1D7FF) { // Mathematical Alphanumeric
			return true
		}
	}
	return false
}

// ContainsDiacritics checks if content contains combining diacritical marks
// Combining Diacritical Marks: U+0300-U+036F
// Latin Extended-A: U+0100-U+017F (includes accented characters)
// Latin Extended-B: U+0180-U+024F
func ContainsDiacritics(content string) bool {
	for _, r := range content {
		if (r >= 0x0300 && r <= 0x036F) || // Combining Diacritical Marks
			(r >= 0x0100 && r <= 0x017F) || // Latin Extended-A
			(r >= 0x0180 && r <= 0x024F) { // Latin Extended-B
			return true
		}
	}
	return false
}

// UnicodeCharacterProfile analyzes unicode content in markdown
// Returns detailed information about what types of unicode are present
type UnicodeCharacterProfile struct {
	HasUnicode       bool
	HasEmoji         bool
	HasCJK           bool
	HasMathSymbols   bool
	HasDiacritics    bool
	NeedsUnicodeFont bool // True if any unicode detected
}

// AnalyzeContent examines markdown content and returns unicode profile
func AnalyzeContent(content string) *UnicodeCharacterProfile {
	profile := &UnicodeCharacterProfile{
		HasUnicode:     DetectUnicodeInContent(content),
		HasEmoji:       ContainsEmoji(content),
		HasCJK:         ContainsCJK(content),
		HasMathSymbols: ContainsMathSymbols(content),
		HasDiacritics:  ContainsDiacritics(content),
	}

	// If any unicode type is present, we need unicode support
	profile.NeedsUnicodeFont = profile.HasUnicode || profile.HasEmoji ||
		profile.HasCJK || profile.HasMathSymbols || profile.HasDiacritics

	return profile
}

// ShouldUseUnicodeEngine determines if a unicode-capable engine is required
// Returns true if content requires unicode support
func ShouldUseUnicodeEngine(content string) bool {
	profile := AnalyzeContent(content)
	return profile.NeedsUnicodeFont
}

// IsHighComplexityUnicode checks if content has complex unicode requirements
// (e.g., CJK or emoji which typically need special font support)
func IsHighComplexityUnicode(content string) bool {
	profile := AnalyzeContent(content)
	return profile.HasCJK || profile.HasEmoji
}

// NormalizeLineEndings normalizes different line endings to \n
// Handles \r\n (Windows), \r (old Mac), and \n (Unix)
func NormalizeLineEndings(content string) string {
	// Replace \r\n with \n
	normalized := bytes.ReplaceAll([]byte(content), []byte("\r\n"), []byte("\n"))
	// Replace remaining \r with \n
	normalized = bytes.ReplaceAll(normalized, []byte("\r"), []byte("\n"))
	return string(normalized)
}

// StripNonPrintable removes control characters while preserving unicode
func StripNonPrintable(content string) string {
	var result []rune
	for _, r := range content {
		// Keep if printable or unicode (>127)
		if unicode.IsPrint(r) || r > 127 {
			result = append(result, r)
		}
	}
	return string(result)
}
