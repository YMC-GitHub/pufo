# pufo - Pure Font Utility & Operations Library

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ymc-github/pufo)](https://goreportcard.com/report/github.com/ymc-github/pufo)

pufo is a Go pure function library focused on font processing. It provides type-safe, side-effect-free functions for font loading, system font management, text measurement, custom font directory scanning, and other common font processing tasks.

## ✨ Features

- 🔤 **Font Loading** - Load from files, system, or built-in fonts
- 📁 **Custom Directories** - Scan custom font directories (e.g., Windows Fonts)
- 📏 **Text Measurement** - Accurate text width and height estimation with multi-language support (Chinese, Japanese, Korean)
- 🔍 **Font Search** - Search fonts in system or custom directories
- 📊 **Text Analysis** - Detailed text composition analysis (Latin, CJK, Digits, Fullwidth)
- 🎯 **Font Fallback** - Support font fallback chains for guaranteed availability
- 🖥️ **System Fonts** - Linux system font management (fc-list)
- ✅ **Type Safety** - Pure function design with no side effects
- 🧪 **Test Coverage** - Comprehensive unit tests and benchmarks

## 📦 Installation

```bash
go get github.com/ymc-github/pufo/ipkg/pufo
```

## 🚀 Quick Start

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    // Load font from file
    face, isBuiltin, err := pufo.LoadFontSmart("/path/to/font.ttf", 16)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    if isBuiltin {
        fmt.Println("Using built-in font")
    } else {
        fmt.Println("Font loaded successfully")
    }
    _ = face
    
    // Measure text width
    text := "Hello 世界"
    width := pufo.ApproximateTextWidth(text, 16)
    fmt.Printf("Text width: %.0fpx\n", width)
    
    // Analyze text composition
    analysis := pufo.AnalyzeText(text, 16)
    fmt.Printf("Characters: %d (Latin: %d, CJK: %d, Digits: %d)\n",
        analysis.TotalChars, analysis.LatinChars, analysis.CJKChars, analysis.DigitChars)
    
    // Scan custom font directory
    fonts, _ := pufo.ScanFontDirectories([]string{"/fonts"})
    fmt.Printf("Found %d font files\n", len(fonts))
}
```

## 📚 API Documentation

### Font Loading

#### LoadFontSmart
Smart font loading that automatically detects whether the input is a file path or system font name.

```go
func LoadFontSmart(fontSpec string, size float64) (font.Face, bool, error)
```

**Parameters:**
- `fontSpec`: Font file path or system font name
- `size`: Font size in pixels

**Returns:**
- `font.Face`: Font face object
- `bool`: Whether it's a built-in font
- `error`: Error information

**Examples:**
```go
// Load from file
face, isBuiltin, err := pufo.LoadFontSmart("/fonts/arial.ttf", 16)

// Load from system
face, isBuiltin, err := pufo.LoadFontSmart("Arial", 16)

// Use built-in font
face, isBuiltin, err := pufo.LoadFontSmart("", 16)
```

#### LoadFontFromFile
Load font from a file path.

```go
func LoadFontFromFile(fontPath string, size float64) (font.Face, error)
```

#### LoadSystemFontByName
Load font by system font name.

```go
func LoadSystemFontByName(fontName string, size float64) (font.Face, error)
```

#### LoadBuiltinFont
Load the built-in Go font (goregular).

```go
func LoadBuiltinFont(size float64) (font.Face, error)
```

### System Font Management

#### SearchSystemFont
Search for font file path in the system.

```go
func SearchSystemFont(fontName string) (string, error)
```

#### IsSystemFontAvailable
Check if a system font is available.

```go
func IsSystemFontAvailable(fontName string) bool
```

**Example:**
```go
if pufo.IsSystemFontAvailable("Microsoft YaHei") {
    fmt.Println("Microsoft YaHei is available")
}
```

### Custom Font Directories

#### ScanFontDirectories
Scan custom font directories and return all font file paths.

```go
func ScanFontDirectories(directories []string) ([]string, error)
```

**Example:**
```go
dirs := []string{"/fonts", "/usr/share/fonts"}
fonts, err := pufo.ScanFontDirectories(dirs)
for _, font := range fonts {
    fmt.Println(font)
}
```

#### SearchFontsInDirectories
Search for fonts in custom directories.

```go
func SearchFontsInDirectories(directories []string, pattern string) ([]FontFileInfo, error)
```

**Example:**
```go
results, err := pufo.SearchFontsInDirectories([]string{"/fonts"}, "Arial")
for _, font := range results {
    fmt.Printf("Found: %s (%s, %.2f KB)\n", font.Name, font.Path, float64(font.Size)/1024)
}
```

#### GetFontNameFromFile
Extract font name from font file path.

```go
func GetFontNameFromFile(fontPath string) string
```

### Text Measurement

#### ApproximateTextWidth
Estimate text width, distinguishing between Latin, CJK, digits, and fullwidth characters.

```go
func ApproximateTextWidth(text string, fontSize float64) float64
```

**Character width rules:**
- Latin characters: `fontSize × 0.5`
- CJK characters: `fontSize × 0.75`
- Digits: `fontSize × 0.5`
- Fullwidth characters: `fontSize × 0.75`

**Example:**
```go
width := pufo.ApproximateTextWidth("Hello 世界 123", 16)
// "Hello" (5×8=40) + " " (8) + "世界" (2×12=24) + " " (8) + "123" (3×8=24) = 104px
```

#### ApproximateTextWidthSimple
Simplified text width calculation where all characters have the same width.

```go
func ApproximateTextWidthSimple(text string, fontSize float64) float64
```

#### ApproximateTextHeight
Estimate text height.

```go
func ApproximateTextHeight(fontSize float64) float64
// Returns fontSize × 1.2
```

#### CalculateLineHeight
Calculate line height.

```go
func CalculateLineHeight(fontSize float64, lineSpacing float64) float64
```

### Text Analysis

#### AnalyzeText
Analyze text composition and return detailed statistics.

```go
func AnalyzeText(text string, fontSize float64) TextAnalysis
```

**TextAnalysis Structure:**
```go
type TextAnalysis struct {
    TotalChars      int     // Total character count
    CJKChars        int     // CJK character count
    LatinChars      int     // Latin character count
    DigitChars      int     // Digit character count
    FullwidthChars  int     // Fullwidth character count
    OtherChars      int     // Other character count
    EstimatedWidth  float64 // Estimated width
    EstimatedHeight float64 // Estimated height
}
```

**Example:**
```go
analysis := pufo.AnalyzeText("Hello 世界 123", 16)
fmt.Printf("Total: %d, Latin: %d, CJK: %d, Digits: %d\n",
    analysis.TotalChars, analysis.LatinChars, analysis.CJKChars, analysis.DigitChars)
```

#### GetCharacterType
Get the type of a character.

```go
func GetCharacterType(r rune) string
```

**Return types:**
- `"Latin"` - Latin letters
- `"CJK"` - Chinese, Japanese, Korean characters
- `"Digit"` - Digits
- `"Fullwidth"` - Fullwidth characters
- `"Other"` - Other characters

### Font Information

#### GetFontInfo
Get font file information.

```go
func GetFontInfo(fontPath string) (*FontInfo, error)
```

#### FontInfo Structure
```go
type FontInfo struct {
    Path     string     // Font path
    Format   FontFormat // Font format (ttf, ttc, otf)
    IsSystem bool       // Whether it's a system font
}
```

#### DetectFontFormat
Detect font file format.

```go
func DetectFontFormat(fontPath string) FontFormat
```

#### IsFontFile
Check if a file is a font file.

```go
func IsFontFile(fontPath string) bool
```

### Font Fallback

#### DefaultFallbackChain
Return the default font fallback chain.

```go
func DefaultFallbackChain() FontFallbackChain
```

**Default fallback chain:**
```go
[]string{
    "Microsoft YaHei",
    "SimHei",
    "Arial",
    "Helvetica",
    "sans-serif",
}
```

#### TryLoadFontChain
Attempt to load a font from a fallback chain.

```go
func TryLoadFontChain(chain FontFallbackChain, size float64) (font.Face, string, error)
```

**Example:**
```go
chain := pufo.DefaultFallbackChain()
face, loadedFont, err := pufo.TryLoadFontChain(chain, 16)
fmt.Printf("Loaded font: %s\n", loadedFont)
```

### Helper Functions

#### NormalizeFontName
Normalize font name for system lookup.

```go
func NormalizeFontName(fontName string) string
```

**Supported conversions:**
- `"yahei"` → `"Microsoft YaHei"`
- `"hei"` → `"SimHei"`
- `"song"` → `"SimSun"`
- `"kaiti"` → `"KaiTi"`
- `"fangsong"` → `"FangSong"`

#### ValidateFontPath
Validate if a font path is valid.

```go
func ValidateFontPath(fontPath string) error
```

## 💡 Complete Examples

### Scan Windows Font Directory

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    // Windows font directory
    windowsFonts := "/mnt/c/Windows/Fonts"
    
    // Scan directory
    fonts, err := pufo.ScanFontDirectories([]string{windowsFonts})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Found %d fonts\n", len(fonts))
    
    // Search for Arial fonts
    results, _ := pufo.SearchFontsInDirectories([]string{windowsFonts}, "Arial")
    for _, font := range results {
        fmt.Printf("  %s - %s (%.2f KB)\n", font.Name, font.Path, float64(font.Size)/1024)
    }
}
```

### Text Measurement Tool

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    text := "The quick brown fox jumps over the lazy dog"
    fontSize := 24.0
    
    // Accurate measurement
    width := pufo.ApproximateTextWidth(text, fontSize)
    height := pufo.ApproximateTextHeight(fontSize)
    
    fmt.Printf("Text: %s\n", text)
    fmt.Printf("Size: %.0fpx\n", fontSize)
    fmt.Printf("Width: %.0fpx\n", width)
    fmt.Printf("Height: %.0fpx\n", height)
    
    // Text analysis
    analysis := pufo.AnalyzeText(text, fontSize)
    fmt.Printf("\nText Analysis:\n")
    fmt.Printf("  Total characters: %d\n", analysis.TotalChars)
    fmt.Printf("  Latin: %d\n", analysis.LatinChars)
    fmt.Printf("  CJK: %d\n", analysis.CJKChars)
    fmt.Printf("  Digits: %d\n", analysis.DigitChars)
}
```

### Font Fallback Example

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    // Get default fallback chain
    chain := pufo.DefaultFallbackChain()
    
    // Attempt to load font
    face, loadedFont, err := pufo.TryLoadFontChain(chain, 16)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Successfully loaded font: %s\n", loadedFont)
    _ = face
    
    // Custom fallback chain
    customChain := pufo.FontFallbackChain{
        "MyCustomFont",
        "Arial",
        "sans-serif",
    }
    
    face, loadedFont, _ = pufo.TryLoadFontChain(customChain, 16)
    fmt.Printf("Loaded from custom chain: %s\n", loadedFont)
}
```

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test -v

# Run benchmarks
go test -bench=. -benchmem

# Generate coverage report
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 Benchmarks

```bash
go test -bench=. -benchmem
```

## 🙏 Acknowledgments

Thanks to all developers who have contributed to this project!

## 📧 Contact

- Project Home: [https://github.com/ymc-github/pufo](https://github.com/ymc-github/pufo)
- Issue Tracker: [Issues](https://github.com/ymc-github/pufo/issues)
