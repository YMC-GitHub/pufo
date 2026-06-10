# pufo - 纯字体处理与操作库

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ymc-github/pufo)](https://goreportcard.com/report/github.com/ymc-github/pufo)

pufo 是一个专注于字体处理的 Go 语言纯函数库。它提供了类型安全、无副作用的函数，用于字体加载、系统字体管理、文本测量、自定义字体目录扫描等常见字体处理任务。

## ✨ 特性

- 🔤 **字体加载** - 支持从文件、系统、内置字体加载
- 📁 **自定义目录** - 支持扫描自定义字体目录（如 Windows 字体目录）
- 📏 **文本测量** - 精确的文本宽度和高度估算，支持中、日、韩等多语言
- 🔍 **字体搜索** - 在系统或自定义目录中搜索字体
- 📊 **文本分析** - 详细的文本组成分析（拉丁、CJK、数字、全角等）
- 🎯 **字体回退** - 支持字体回退链，确保字体可用
- 🖥️ **系统字体** - Linux 系统字体管理（fc-list）
- ✅ **类型安全** - 纯函数设计，无副作用
- 🧪 **测试覆盖** - 完整的单元测试和基准测试

## 📦 安装

```bash
go get github.com/ymc-github/pufo/ipkg/pufo
```

## 🚀 快速开始

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    // 从文件加载字体
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
    
    // 测量文本宽度
    text := "Hello 世界"
    width := pufo.ApproximateTextWidth(text, 16)
    fmt.Printf("Text width: %.0fpx\n", width)
    
    // 分析文本组成
    analysis := pufo.AnalyzeText(text, 16)
    fmt.Printf("Characters: %d (Latin: %d, CJK: %d, Digits: %d)\n",
        analysis.TotalChars, analysis.LatinChars, analysis.CJKChars, analysis.DigitChars)
    
    // 扫描自定义字体目录
    fonts, _ := pufo.ScanFontDirectories([]string{"/fonts"})
    fmt.Printf("Found %d font files\n", len(fonts))
}
```

## 📚 API 文档

### 字体加载

#### LoadFontSmart
智能加载字体，自动检测是文件路径还是系统字体名。

```go
func LoadFontSmart(fontSpec string, size float64) (font.Face, bool, error)
```

**参数：**
- `fontSpec`: 字体文件路径或系统字体名称
- `size`: 字体大小（像素）

**返回：**
- `font.Face`: 字体对象
- `bool`: 是否为内置字体
- `error`: 错误信息

**示例：**
```go
// 从文件加载
face, isBuiltin, err := pufo.LoadFontSmart("/fonts/arial.ttf", 16)

// 从系统加载
face, isBuiltin, err := pufo.LoadFontSmart("Arial", 16)

// 使用内置字体
face, isBuiltin, err := pufo.LoadFontSmart("", 16)
```

#### LoadFontFromFile
从文件路径加载字体。

```go
func LoadFontFromFile(fontPath string, size float64) (font.Face, error)
```

#### LoadSystemFontByName
通过系统字体名称加载字体。

```go
func LoadSystemFontByName(fontName string, size float64) (font.Face, error)
```

#### LoadBuiltinFont
加载内置的 Go 字体（goregular）。

```go
func LoadBuiltinFont(size float64) (font.Face, error)
```

### 系统字体管理

#### SearchSystemFont
在系统中搜索字体文件路径。

```go
func SearchSystemFont(fontName string) (string, error)
```

#### IsSystemFontAvailable
检查系统字体是否可用。

```go
func IsSystemFontAvailable(fontName string) bool
```

**示例：**
```go
if pufo.IsSystemFontAvailable("Microsoft YaHei") {
    fmt.Println("Microsoft YaHei is available")
}
```

### 自定义字体目录

#### ScanFontDirectories
扫描自定义字体目录，返回所有字体文件路径。

```go
func ScanFontDirectories(directories []string) ([]string, error)
```

**示例：**
```go
dirs := []string{"/fonts", "/usr/share/fonts"}
fonts, err := pufo.ScanFontDirectories(dirs)
for _, font := range fonts {
    fmt.Println(font)
}
```

#### SearchFontsInDirectories
在自定义目录中搜索字体。

```go
func SearchFontsInDirectories(directories []string, pattern string) ([]FontFileInfo, error)
```

**示例：**
```go
results, err := pufo.SearchFontsInDirectories([]string{"/fonts"}, "Arial")
for _, font := range results {
    fmt.Printf("Found: %s (%s, %.2f KB)\n", font.Name, font.Path, float64(font.Size)/1024)
}
```

#### GetFontNameFromFile
从字体文件名提取字体名称。

```go
func GetFontNameFromFile(fontPath string) string
```

### 文本测量

#### ApproximateTextWidth
估算文本宽度，区分拉丁字符、CJK字符、数字和全角字符。

```go
func ApproximateTextWidth(text string, fontSize float64) float64
```

**字符宽度规则：**
- 拉丁字符：`fontSize × 0.5`
- CJK 字符：`fontSize × 0.75`
- 数字：`fontSize × 0.5`
- 全角字符：`fontSize × 0.75`

**示例：**
```go
width := pufo.ApproximateTextWidth("Hello 世界 123", 16)
// "Hello" (5×8=40) + " " (8) + "世界" (2×12=24) + " " (8) + "123" (3×8=24) = 104px
```

#### ApproximateTextWidthSimple
简化版文本宽度计算，所有字符按相同宽度处理。

```go
func ApproximateTextWidthSimple(text string, fontSize float64) float64
```

#### ApproximateTextHeight
估算文本高度。

```go
func ApproximateTextHeight(fontSize float64) float64
// 返回 fontSize × 1.2
```

#### CalculateLineHeight
计算行高。

```go
func CalculateLineHeight(fontSize float64, lineSpacing float64) float64
```

### 文本分析

#### AnalyzeText
分析文本组成，返回详细的统计信息。

```go
func AnalyzeText(text string, fontSize float64) TextAnalysis
```

**TextAnalysis 结构：**
```go
type TextAnalysis struct {
    TotalChars      int     // 总字符数
    CJKChars        int     // CJK 字符数
    LatinChars      int     // 拉丁字符数
    DigitChars      int     // 数字字符数
    FullwidthChars  int     // 全角字符数
    OtherChars      int     // 其他字符数
    EstimatedWidth  float64 // 估算宽度
    EstimatedHeight float64 // 估算高度
}
```

**示例：**
```go
analysis := pufo.AnalyzeText("Hello 世界 123", 16)
fmt.Printf("Total: %d, Latin: %d, CJK: %d, Digits: %d\n",
    analysis.TotalChars, analysis.LatinChars, analysis.CJKChars, analysis.DigitChars)
```

#### GetCharacterType
获取字符类型。

```go
func GetCharacterType(r rune) string
```

**返回类型：**
- `"Latin"` - 拉丁字母
- `"CJK"` - 中日韩文字
- `"Digit"` - 数字
- `"Fullwidth"` - 全角字符
- `"Other"` - 其他字符

### 字体信息

#### GetFontInfo
获取字体文件信息。

```go
func GetFontInfo(fontPath string) (*FontInfo, error)
```

#### FontInfo 结构
```go
type FontInfo struct {
    Path     string     // 字体路径
    Format   FontFormat // 字体格式 (ttf, ttc, otf)
    IsSystem bool       // 是否为系统字体
}
```

#### DetectFontFormat
检测字体文件格式。

```go
func DetectFontFormat(fontPath string) FontFormat
```

#### IsFontFile
检查文件是否为字体文件。

```go
func IsFontFile(fontPath string) bool
```

### 字体回退

#### DefaultFallbackChain
返回默认字体回退链。

```go
func DefaultFallbackChain() FontFallbackChain
```

**默认回退链：**
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
尝试从回退链中加载字体。

```go
func TryLoadFontChain(chain FontFallbackChain, size float64) (font.Face, string, error)
```

**示例：**
```go
chain := pufo.DefaultFallbackChain()
face, loadedFont, err := pufo.TryLoadFontChain(chain, 16)
fmt.Printf("Loaded font: %s\n", loadedFont)
```

### 辅助函数

#### NormalizeFontName
规范化字体名称。

```go
func NormalizeFontName(fontName string) string
```

**支持的转换：**
- `"yahei"` → `"Microsoft YaHei"`
- `"hei"` → `"SimHei"`
- `"song"` → `"SimSun"`
- `"kaiti"` → `"KaiTi"`
- `"fangsong"` → `"FangSong"`

#### ValidateFontPath
验证字体路径是否有效。

```go
func ValidateFontPath(fontPath string) error
```

## 💡 完整示例

### 扫描 Windows 字体目录

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    // Windows 字体目录
    windowsFonts := "/mnt/c/Windows/Fonts"
    
    // 扫描目录
    fonts, err := pufo.ScanFontDirectories([]string{windowsFonts})
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Found %d fonts\n", len(fonts))
    
    // 搜索 Arial 字体
    results, _ := pufo.SearchFontsInDirectories([]string{windowsFonts}, "Arial")
    for _, font := range results {
        fmt.Printf("  %s - %s (%.2f KB)\n", font.Name, font.Path, float64(font.Size)/1024)
    }
}
```

### 文本测量工具

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    text := "The quick brown fox 跳过了 lazy dog 123"
    fontSize := 24.0
    
    // 精确测量
    width := pufo.ApproximateTextWidth(text, fontSize)
    height := pufo.ApproximateTextHeight(fontSize)
    
    fmt.Printf("Text: %s\n", text)
    fmt.Printf("Size: %.0fpx\n", fontSize)
    fmt.Printf("Width: %.0fpx\n", width)
    fmt.Printf("Height: %.0fpx\n", height)
    
    // 文本分析
    analysis := pufo.AnalyzeText(text, fontSize)
    fmt.Printf("\nText Analysis:\n")
    fmt.Printf("  Total characters: %d\n", analysis.TotalChars)
    fmt.Printf("  Latin: %d\n", analysis.LatinChars)
    fmt.Printf("  CJK: %d\n", analysis.CJKChars)
    fmt.Printf("  Digits: %d\n", analysis.DigitChars)
}
```

### 字体回退示例

```go
package main

import (
    "fmt"
    "github.com/ymc-github/pufo/ipkg/pufo"
)

func main() {
    // 获取默认回退链
    chain := pufo.DefaultFallbackChain()
    
    // 尝试加载字体
    face, loadedFont, err := pufo.TryLoadFontChain(chain, 16)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Successfully loaded font: %s\n", loadedFont)
    _ = face
    
    // 自定义回退链
    customChain := pufo.FontFallbackChain{
        "MyCustomFont",
        "Arial",
        "sans-serif",
    }
    
    face, loadedFont, _ = pufo.TryLoadFontChain(customChain, 16)
    fmt.Printf("Loaded from custom chain: %s\n", loadedFont)
}
```

## 🧪 测试

运行测试套件：

```bash
# 运行所有测试
go test -v

# 运行基准测试
go test -bench=. -benchmem

# 生成覆盖率报告
go test -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 📊 性能基准

```bash
go test -bench=. -benchmem
```


## 🙏 致谢

感谢所有为这个项目做出贡献的开发者！

## 📧 联系方式

- 项目主页: [https://github.com/ymc-github/pufo](https://github.com/ymc-github/pufo)
- 问题反馈: [Issues](https://github.com/ymc-github/pufo/issues)
