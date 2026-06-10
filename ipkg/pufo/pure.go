// Package pufo provides pure functions for font handling
// pufo = Pure Font Utility & Operations
package pufo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/gofont/goregular"
)

// ============================================================================
// Type Definitions
// ============================================================================

// FontFormat represents the font file format
type FontFormat string

const (
	FormatTTF FontFormat = "ttf"
	FormatTTC FontFormat = "ttc"
	FormatOTF FontFormat = "otf"
)

// FontLoadOptions contains options for loading fonts
type FontLoadOptions struct {
	Size    float64
	DPI     float64
	Hinting font.Hinting
}

// FontInfo contains information about a loaded font
type FontInfo struct {
	Path     string
	Format   FontFormat
	IsSystem bool
}

// FontFileInfo contains information about a font file
type FontFileInfo struct {
	Path    string
	Name    string
	Size    int64
	Format  FontFormat
	ModTime time.Time
}

// FontFallbackChain defines a chain of font fallbacks
type FontFallbackChain []string

// ============================================================================
// Unicode Ranges for CJK Characters
// ============================================================================

// isCJKCharacter checks if a rune is a CJK (Chinese, Japanese, Korean) character
func isCJKCharacter(r rune) bool {
	// CJK Unified Ideographs
	if r >= 0x4E00 && r <= 0x9FFF {
		return true
	}
	// CJK Unified Ideographs Extension A
	if r >= 0x3400 && r <= 0x4DBF {
		return true
	}
	// CJK Unified Ideographs Extension B
	if r >= 0x20000 && r <= 0x2A6DF {
		return true
	}
	// CJK Unified Ideographs Extension C
	if r >= 0x2A700 && r <= 0x2B73F {
		return true
	}
	// CJK Unified Ideographs Extension D
	if r >= 0x2B740 && r <= 0x2B81F {
		return true
	}
	// CJK Unified Ideographs Extension E
	if r >= 0x2B820 && r <= 0x2CEAF {
		return true
	}
	// CJK Unified Ideographs Extension F
	if r >= 0x2CEB0 && r <= 0x2EBEF {
		return true
	}
	// CJK Compatibility Ideographs
	if r >= 0xF900 && r <= 0xFAFF {
		return true
	}
	// CJK Compatibility Ideographs Supplement
	if r >= 0x2F800 && r <= 0x2FA1F {
		return true
	}
	// Hiragana (Japanese)
	if r >= 0x3040 && r <= 0x309F {
		return true
	}
	// Katakana (Japanese)
	if r >= 0x30A0 && r <= 0x30FF {
		return true
	}
	// Hangul (Korean)
	if r >= 0xAC00 && r <= 0xD7AF {
		return true
	}
	// Hangul Jamo
	if r >= 0x1100 && r <= 0x11FF {
		return true
	}
	return false
}

// isFullwidthCharacter checks if a rune is a fullwidth character
func isFullwidthCharacter(r rune) bool {
	// Fullwidth ASCII variants
	if r >= 0xFF01 && r <= 0xFF5E {
		return true
	}
	// Fullwidth brackets and punctuation
	if r >= 0x3000 && r <= 0x303F {
		return true
	}
	return false
}

// ============================================================================
// Font File Detection
// ============================================================================

// DetectFontFormat detects the font format from file extension
func DetectFontFormat(fontPath string) FontFormat {
	ext := strings.ToLower(filepath.Ext(fontPath))
	
	switch ext {
	case ".ttf":
		return FormatTTF
	case ".ttc":
		return FormatTTC
	case ".otf":
		return FormatOTF
	default:
		return ""
	}
}

// IsFontFile checks if the given path points to a font file
func IsFontFile(fontPath string) bool {
	ext := strings.ToLower(filepath.Ext(fontPath))
	return ext == ".ttf" || ext == ".ttc" || ext == ".otf"
}

// ValidateFontPath validates if a font file exists and is accessible
func ValidateFontPath(fontPath string) error {
	if fontPath == "" {
		return fmt.Errorf("font path is empty")
	}
	
	info, err := os.Stat(fontPath)
	if err != nil {
		return fmt.Errorf("font file not found: %w", err)
	}
	
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a font file")
	}
	
	return nil
}

// ============================================================================
// Font Parsing
// ============================================================================

// ParseFontData parses font data from bytes
func ParseFontData(fontData []byte, isTTC bool) (*opentype.Font, error) {
	if isTTC {
		collection, err := opentype.ParseCollection(fontData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse TTC collection: %w", err)
		}
		
		if collection.NumFonts() == 0 {
			return nil, fmt.Errorf("TTC collection contains no fonts")
		}
		
		return collection.Font(0)
	}
	
	return opentype.Parse(fontData)
}

// ParseFontFile parses a font file from the given path
func ParseFontFile(fontPath string) (*opentype.Font, FontFormat, error) {
	if err := ValidateFontPath(fontPath); err != nil {
		return nil, "", err
	}
	
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read font file: %w", err)
	}
	
	format := DetectFontFormat(fontPath)
	isTTC := format == FormatTTC
	
	parsedFont, err := ParseFontData(fontData, isTTC)
	if err != nil {
		return nil, "", err
	}
	
	return parsedFont, format, nil
}

// ============================================================================
// Font Face Creation
// ============================================================================

// DefaultFontLoadOptions returns default font loading options
func DefaultFontLoadOptions(size float64) FontLoadOptions {
	return FontLoadOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}
}

// CreateFontFace creates a font.Face from an opentype.Font
func CreateFontFace(parsedFont *opentype.Font, opts FontLoadOptions) (font.Face, error) {
	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    opts.Size,
		DPI:     opts.DPI,
		Hinting: opts.Hinting,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}
	
	return face, nil
}

// LoadBuiltinFont loads the built-in Go font (goregular)
func LoadBuiltinFont(size float64) (font.Face, error) {
	parsedFont, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse builtin font: %w", err)
	}
	
	opts := DefaultFontLoadOptions(size)
	return CreateFontFace(parsedFont, opts)
}

// LoadFontFromFile loads a font face from a file path
func LoadFontFromFile(fontPath string, size float64) (font.Face, error) {
	parsedFont, _, err := ParseFontFile(fontPath)
	if err != nil {
		return nil, err
	}
	
	opts := DefaultFontLoadOptions(size)
	return CreateFontFace(parsedFont, opts)
}

// ============================================================================
// System Font Management
// ============================================================================

// SearchSystemFont searches for a font in the system by name
// Returns the font file path if found
func SearchSystemFont(fontName string) (string, error) {
	if fontName == "" {
		return "", fmt.Errorf("font name is empty")
	}
	
	// Use fc-list on Linux systems
	cmd := exec.Command("fc-list", "-f", "%{file}", fontName)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to search font '%s': %w", fontName, err)
	}
	
	fontPath := strings.TrimSpace(string(output))
	if fontPath == "" {
		return "", fmt.Errorf("font '%s' not found in system", fontName)
	}
	
	return fontPath, nil
}

// LoadSystemFontByName loads a system font by its name
func LoadSystemFontByName(fontName string, size float64) (font.Face, error) {
	fontPath, err := SearchSystemFont(fontName)
	if err != nil {
		return nil, err
	}
	
	return LoadFontFromFile(fontPath, size)
}

// IsSystemFontAvailable checks if a system font is available
func IsSystemFontAvailable(fontName string) bool {
	_, err := SearchSystemFont(fontName)
	return err == nil
}

// ============================================================================
// Font Loading (Unified Interface)
// ============================================================================

// LoadFont loads a font from various sources
// Priority: file path > system font > builtin font
func LoadFont(fontPath string, size float64, preferSystem bool) (font.Face, bool, error) {
	if fontPath != "" {
		if IsFontFile(fontPath) {
			face, err := LoadFontFromFile(fontPath, size)
			if err == nil {
				return face, false, nil
			}
			if !preferSystem {
				face, err := LoadBuiltinFont(size)
				return face, true, err
			}
		}
		
		if preferSystem || !IsFontFile(fontPath) {
			face, err := LoadSystemFontByName(fontPath, size)
			if err == nil {
				return face, true, nil
			}
		}
	}
	
	face, err := LoadBuiltinFont(size)
	return face, true, err
}

// LoadFontSmart loads a font with smart detection
// Automatically detects if the path is a file or font name
func LoadFontSmart(fontSpec string, size float64) (font.Face, bool, error) {
	if fontSpec == "" {
		face, err := LoadBuiltinFont(size)
		return face, true, err
	}
	
	if IsFontFile(fontSpec) {
		if err := ValidateFontPath(fontSpec); err == nil {
			face, err := LoadFontFromFile(fontSpec, size)
			if err == nil {
				return face, false, nil
			}
		}
	}
	
	face, err := LoadSystemFontByName(fontSpec, size)
	if err == nil {
		return face, true, nil
	}
	
	face, err = LoadBuiltinFont(size)
	return face, true, err
}

// ============================================================================
// Font Information
// ============================================================================

// GetFontInfo returns information about a font file
func GetFontInfo(fontPath string) (*FontInfo, error) {
	if err := ValidateFontPath(fontPath); err != nil {
		return nil, err
	}
	
	info := &FontInfo{
		Path:   fontPath,
		Format: DetectFontFormat(fontPath),
	}
	
	return info, nil
}

// GetBuiltinFontInfo returns information about the builtin font
func GetBuiltinFontInfo() *FontInfo {
	return &FontInfo{
		Path:     "builtin",
		Format:   FormatTTF,
		IsSystem: false,
	}
}

// ============================================================================
// Font Metrics (Pure Functions)
// ============================================================================

// ApproximateTextWidth approximates text width based on font size and character types
// This function distinguishes between CJK characters, fullwidth characters, and Latin characters
func ApproximateTextWidth(text string, fontSize float64) float64 {
	if text == "" {
		return 0
	}
	
	// Character width factors
	latinWidth := fontSize * 0.5      // Latin characters (a, b, c)
	cjkWidth := fontSize * 0.75       // CJK characters (中, 文, 日, 韓)
	fullwidthWidth := fontSize * 0.75 // Fullwidth characters (！, ”, ＃)
	digitWidth := fontSize * 0.5      // Digits (0-9)
	
	var totalWidth float64
	
	for _, r := range text {
		switch {
		case isCJKCharacter(r):
			totalWidth += cjkWidth
		case isFullwidthCharacter(r):
			totalWidth += fullwidthWidth
		case r >= '0' && r <= '9':
			totalWidth += digitWidth
		default:
			totalWidth += latinWidth
		}
	}
	
	return totalWidth
}

// ApproximateTextWidthSimple is a simplified version that treats all characters equally
func ApproximateTextWidthSimple(text string, fontSize float64) float64 {
	avgCharWidth := fontSize * 0.5
	charCount := len([]rune(text))
	return float64(charCount) * avgCharWidth
}

// ApproximateTextHeight approximates text height based on font size
func ApproximateTextHeight(fontSize float64) float64 {
	// Line height is typically 1.2 to 1.5 times font size
	return fontSize * 1.2
}

// CalculateLineHeight calculates line height based on font size and line spacing
func CalculateLineHeight(fontSize float64, lineSpacing float64) float64 {
	return fontSize * lineSpacing
}

// GetCharacterType returns the type of a character for debugging purposes
func GetCharacterType(r rune) string {
	switch {
	case isCJKCharacter(r):
		return "CJK"
	case isFullwidthCharacter(r):
		return "Fullwidth"
	case r >= '0' && r <= '9':
		return "Digit"
	case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
		return "Latin"
	default:
		return "Other"
	}
}

// ============================================================================
// Font Name Normalization
// ============================================================================

// NormalizeFontName normalizes font name for system lookup
func NormalizeFontName(fontName string) string {
	// Remove extension if present
	fontName = strings.TrimSuffix(fontName, ".ttf")
	fontName = strings.TrimSuffix(fontName, ".ttc")
	fontName = strings.TrimSuffix(fontName, ".otf")
	
	// Convert common names
	switch strings.ToLower(fontName) {
	case "yahei", "microsoft yahei":
		return "Microsoft YaHei"
	case "simhei", "hei":
		return "SimHei"
	case "simsun", "song":
		return "SimSun"
	case "kaiti", "kai":
		return "KaiTi"
	case "fangsong":
		return "FangSong"
	default:
		return fontName
	}
}

// ============================================================================
// Font Fallback
// ============================================================================

// DefaultFallbackChain returns the default font fallback chain
func DefaultFallbackChain() FontFallbackChain {
	return FontFallbackChain{
		"Microsoft YaHei",
		"SimHei",
		"Arial",
		"Helvetica",
		"sans-serif",
	}
}

// TryLoadFontChain tries to load a font from a chain of fallbacks
func TryLoadFontChain(chain FontFallbackChain, size float64) (font.Face, string, error) {
	for _, fontName := range chain {
		if IsSystemFontAvailable(fontName) {
			face, err := LoadSystemFontByName(fontName, size)
			if err == nil {
				return face, fontName, nil
			}
		}
	}
	
	// Fallback to builtin
	face, err := LoadBuiltinFont(size)
	return face, "builtin", err
}

// ============================================================================
// Custom Font Directory Support
// ============================================================================

// ScanFontDirectories scans custom font directories for font files
func ScanFontDirectories(directories []string) ([]string, error) {
	var fonts []string
	
	for _, dir := range directories {
		if dir == "" {
			continue
		}
		
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // 忽略无法访问的文件
			}
			
			if !info.IsDir() && IsFontFile(path) {
				fonts = append(fonts, path)
			}
			return nil
		})
		
		if err != nil {
			return fonts, fmt.Errorf("error scanning directory %s: %w", dir, err)
		}
	}
	
	return fonts, nil
}

// GetFontNameFromFile extracts font name from font file (simplified version)
// For full implementation, you would need to parse the font file's name table
func GetFontNameFromFile(fontPath string) string {
	base := filepath.Base(fontPath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	
	// Common font file naming patterns
	commonNames := map[string]string{
		"msyh":    "Microsoft YaHei",
		"msyhbd":  "Microsoft YaHei Bold",
		"simhei":  "SimHei",
		"simsun":  "SimSun",
		"simkai":  "KaiTi",
		"arial":   "Arial",
		"times":   "Times New Roman",
		"cour":    "Courier New",
		"verdana": "Verdana",
		"tahoma":  "Tahoma",
		"georgia": "Georgia",
	}
	
	lower := strings.ToLower(name)
	if displayName, ok := commonNames[lower]; ok {
		return displayName
	}
	
	return name
}

// SearchFontsInDirectories searches for fonts in specified directories
func SearchFontsInDirectories(directories []string, pattern string) ([]FontFileInfo, error) {
	fonts, err := ScanFontDirectories(directories)
	if err != nil {
		return nil, err
	}
	
	var results []FontFileInfo
	patternLower := strings.ToLower(pattern)
	
	for _, fontPath := range fonts {
		info, err := os.Stat(fontPath)
		if err != nil {
			continue
		}
		
		fontName := GetFontNameFromFile(fontPath)
		
		// Filter by pattern
		if pattern != "" {
			if !strings.Contains(strings.ToLower(fontName), patternLower) &&
				!strings.Contains(strings.ToLower(fontPath), patternLower) {
				continue
			}
		}
		
		results = append(results, FontFileInfo{
			Path:    fontPath,
			Name:    fontName,
			Size:    info.Size(),
			Format:  DetectFontFormat(fontPath),
			ModTime: info.ModTime(),
		})
	}
	
	return results, nil
}

// ============================================================================
// Text Analysis
// ============================================================================

// TextAnalysis contains detailed information about text composition
type TextAnalysis struct {
	TotalChars      int
	CJKChars        int
	LatinChars      int
	DigitChars      int
	FullwidthChars  int
	OtherChars      int
	EstimatedWidth  float64
	EstimatedHeight float64
}

// AnalyzeText analyzes text composition and returns detailed metrics
func AnalyzeText(text string, fontSize float64) TextAnalysis {
	analysis := TextAnalysis{}
	
	for _, r := range text {
		analysis.TotalChars++
		
		switch {
		case isCJKCharacter(r):
			analysis.CJKChars++
		case isFullwidthCharacter(r):
			analysis.FullwidthChars++
		case r >= '0' && r <= '9':
			analysis.DigitChars++
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'):
			analysis.LatinChars++
		default:
			analysis.OtherChars++
		}
	}
	
	analysis.EstimatedWidth = ApproximateTextWidth(text, fontSize)
	analysis.EstimatedHeight = ApproximateTextHeight(fontSize)
	
	return analysis
}