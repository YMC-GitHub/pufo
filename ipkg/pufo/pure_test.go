package pufo

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDetectFontFormat(t *testing.T) {
	tests := []struct {
		name     string
		fontPath string
		want     FontFormat
	}{
		{"TTF file", "/path/to/font.ttf", FormatTTF},
		{"TTC file", "/path/to/font.ttc", FormatTTC},
		{"OTF file", "/path/to/font.otf", FormatOTF},
		{"Unknown extension", "/path/to/font.xyz", ""},
		{"No extension", "/path/to/font", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectFontFormat(tt.fontPath); got != tt.want {
				t.Errorf("DetectFontFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFontFile(t *testing.T) {
	tests := []struct {
		name     string
		fontPath string
		want     bool
	}{
		{"TTF file", "font.ttf", true},
		{"TTC file", "font.ttc", true},
		{"OTF file", "font.otf", true},
		{"Not a font", "file.txt", false},
		{"No extension", "font", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFontFile(tt.fontPath); got != tt.want {
				t.Errorf("IsFontFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateFontPath(t *testing.T) {
	tests := []struct {
		name     string
		fontPath string
		wantErr  bool
	}{
		{"Empty path", "", true},
		{"Non-existent file", "/nonexistent/font.ttf", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFontPath(tt.fontPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFontPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseFontData(t *testing.T) {
	// Test with invalid data
	invalidData := []byte("invalid font data")
	_, err := ParseFontData(invalidData, false)
	if err == nil {
		t.Error("ParseFontData() should fail with invalid data")
	}
}

func TestDefaultFontLoadOptions(t *testing.T) {
	size := 16.0
	opts := DefaultFontLoadOptions(size)

	if opts.Size != size {
		t.Errorf("DefaultFontLoadOptions().Size = %v, want %v", opts.Size, size)
	}

	if opts.DPI != 72 {
		t.Errorf("DefaultFontLoadOptions().DPI = %v, want 72", opts.DPI)
	}
}

func TestLoadBuiltinFont(t *testing.T) {
	face, err := LoadBuiltinFont(16)
	if err != nil {
		t.Errorf("LoadBuiltinFont() error = %v", err)
	}
	if face == nil {
		t.Error("LoadBuiltinFont() returned nil face")
	}
}

func TestNormalizeFontName(t *testing.T) {
	tests := []struct {
		name     string
		fontName string
		want     string
	}{
		{"YaHei", "yahei", "Microsoft YaHei"},
		{"Microsoft YaHei", "Microsoft YaHei", "Microsoft YaHei"},
		{"SimHei", "hei", "SimHei"},
		{"SimSun", "song", "SimSun"},
		{"KaiTi", "kaiti", "KaiTi"},
		{"FangSong", "fangsong", "FangSong"},
		{"Custom", "Custom Font", "Custom Font"},
		{"With extension", "arial.ttf", "arial"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeFontName(tt.fontName); got != tt.want {
				t.Errorf("NormalizeFontName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproximateTextWidth(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		fontSize float64
		want     float64
	}{
		{"Empty text", "", 16, 0},
		{"English letters", "Hello", 16, 40},
		{"English with spaces", "Hello World", 16, 88},
		{"Chinese text", "你好", 16, 24},
		{"Mixed Chinese and English", "Hello世界", 16, 64},
		{"Digits", "12345", 16, 40},
		{"Fullwidth characters", "１２３", 16, 36},
		{"Large font", "Hi", 32, 32},
		{"Japanese Hiragana", "こんにちは", 16, 60},
		{"Japanese Katakana", "コンニチハ", 16, 60},
		{"Korean Hangul", "안녕하세요", 16, 60},
		{"Mixed all", "Hello世界123", 16, 88},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproximateTextWidth(tt.text, tt.fontSize)
			if got-tt.want > 0.1 || tt.want-got > 0.1 {
				t.Errorf("ApproximateTextWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproximateTextWidthSimple(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		fontSize float64
		want     float64
	}{
		{"Empty text", "", 16, 0},
		{"English letters", "Hello", 16, 40},
		{"Chinese text", "你好", 16, 16},
		{"Mixed text", "Hello世界", 16, 56},
		{"Digits", "12345", 16, 40},
		{"Large font", "Hi", 32, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproximateTextWidthSimple(tt.text, tt.fontSize)
			if got != tt.want {
				t.Errorf("ApproximateTextWidthSimple() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApproximateTextHeight(t *testing.T) {
	tests := []struct {
		name     string
		fontSize float64
		want     float64
	}{
		{"Small font", 12, 14.4},
		{"Normal font", 16, 19.2},
		{"Large font", 32, 38.4},
		{"Zero font", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproximateTextHeight(tt.fontSize)
			if got-tt.want > 0.0001 || tt.want-got > 0.0001 {
				t.Errorf("ApproximateTextHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateLineHeight(t *testing.T) {
	tests := []struct {
		name        string
		fontSize    float64
		lineSpacing float64
		want        float64
	}{
		{"Single spacing", 16, 1.0, 16},
		{"1.5 spacing", 16, 1.5, 24},
		{"Double spacing", 16, 2.0, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateLineHeight(tt.fontSize, tt.lineSpacing)
			if got != tt.want {
				t.Errorf("CalculateLineHeight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultFallbackChain(t *testing.T) {
	chain := DefaultFallbackChain()

	expected := FontFallbackChain{
		"Microsoft YaHei",
		"SimHei",
		"Arial",
		"Helvetica",
		"sans-serif",
	}

	if len(chain) != len(expected) {
		t.Errorf("DefaultFallbackChain() length = %v, want %v", len(chain), len(expected))
	}

	for i := range chain {
		if chain[i] != expected[i] {
			t.Errorf("DefaultFallbackChain()[%d] = %v, want %v", i, chain[i], expected[i])
		}
	}
}

func TestGetFontInfo(t *testing.T) {
	_, err := GetFontInfo("/nonexistent/font.ttf")
	if err == nil {
		t.Error("GetFontInfo() should fail with non-existent file")
	}
}

func TestGetBuiltinFontInfo(t *testing.T) {
	info := GetBuiltinFontInfo()
	
	if info.Path != "builtin" {
		t.Errorf("GetBuiltinFontInfo().Path = %v, want 'builtin'", info.Path)
	}
	
	if info.Format != FormatTTF {
		t.Errorf("GetBuiltinFontInfo().Format = %v, want %v", info.Format, FormatTTF)
	}
}

func TestGetFontNameFromFile(t *testing.T) {
	tests := []struct {
		name     string
		fontPath string
		want     string
	}{
		{"MS YaHei", "msyh.ttf", "Microsoft YaHei"},
		{"SimHei", "simhei.ttf", "SimHei"},
		{"SimSun", "simsun.ttc", "SimSun"},
		{"Arial", "arial.ttf", "Arial"},
		{"Custom", "customfont.ttf", "customfont"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFontNameFromFile(tt.fontPath)
			if got != tt.want {
				t.Errorf("GetFontNameFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanFontDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fonttest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := []string{"test1.ttf", "test2.ttc", "test3.otf", "notfont.txt"}
	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	fonts, err := ScanFontDirectories([]string{tmpDir})
	if err != nil {
		t.Errorf("ScanFontDirectories() error = %v", err)
	}

	if len(fonts) != 3 {
		t.Errorf("ScanFontDirectories() found %d fonts, want 3", len(fonts))
	}
}

func TestSearchFontsInDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "fontsearch")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFonts := map[string]string{
		"msyh.ttf":   "Microsoft YaHei",
		"simhei.ttf": "SimHei",
		"arial.ttf":  "Arial",
	}
	
	for filename := range testFonts {
		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	results, err := SearchFontsInDirectories([]string{tmpDir}, "YaHei")
	if err != nil {
		t.Errorf("SearchFontsInDirectories() error = %v", err)
	}
	
	if len(results) != 1 {
		t.Errorf("SearchFontsInDirectories() found %d fonts, want 1", len(results))
	}
	
	if len(results) > 0 && results[0].Name != "Microsoft YaHei" {
		t.Errorf("SearchFontsInDirectories() name = %v, want 'Microsoft YaHei'", results[0].Name)
	}
}

func TestFontFileInfoStruct(t *testing.T) {
	now := time.Now()
	info := FontFileInfo{
		Path:    "/path/to/font.ttf",
		Name:    "Test Font",
		Size:    1024,
		Format:  FormatTTF,
		ModTime: now,
	}
	
	if info.Path != "/path/to/font.ttf" {
		t.Errorf("FontFileInfo.Path = %v, want '/path/to/font.ttf'", info.Path)
	}
	
	if info.Name != "Test Font" {
		t.Errorf("FontFileInfo.Name = %v, want 'Test Font'", info.Name)
	}
	
	if info.Size != 1024 {
		t.Errorf("FontFileInfo.Size = %v, want 1024", info.Size)
	}
	
	if info.Format != FormatTTF {
		t.Errorf("FontFileInfo.Format = %v, want %v", info.Format, FormatTTF)
	}
	
	if !info.ModTime.Equal(now) {
		t.Errorf("FontFileInfo.ModTime = %v, want %v", info.ModTime, now)
	}
}

func TestLoadFontSmart(t *testing.T) {
	face, isBuiltin, err := LoadFontSmart("", 16)
	if err != nil {
		t.Errorf("LoadFontSmart() with empty spec error = %v", err)
	}
	if !isBuiltin {
		t.Error("LoadFontSmart() with empty spec should return builtin font")
	}
	if face == nil {
		t.Error("LoadFontSmart() with empty spec returned nil face")
	}

	face, isBuiltin, err = LoadFontSmart("NonExistentFontNameXYZ123", 16)
	if err != nil {
		t.Errorf("LoadFontSmart() with invalid name error = %v", err)
	}
	if !isBuiltin {
		t.Error("LoadFontSmart() with invalid name should fallback to builtin")
	}
	if face == nil {
		t.Error("LoadFontSmart() with invalid name returned nil face")
	}
}

func TestTryLoadFontChain(t *testing.T) {
	chain := FontFallbackChain{"NonExistentFont1", "NonExistentFont2"}
	face, loadedFont, err := TryLoadFontChain(chain, 16)
	
	if err != nil {
		t.Errorf("TryLoadFontChain() error = %v", err)
	}
	if loadedFont != "builtin" {
		t.Errorf("TryLoadFontChain() loadedFont = %v, want 'builtin'", loadedFont)
	}
	if face == nil {
		t.Error("TryLoadFontChain() returned nil face")
	}
}

func TestIsSystemFontAvailable(t *testing.T) {
	available := IsSystemFontAvailable("")
	if available {
		t.Error("IsSystemFontAvailable() with empty string should return false")
	}
}

func TestGetCharacterType(t *testing.T) {
	tests := []struct {
		r    rune
		want string
	}{
		{'A', "Latin"},
		{'z', "Latin"},
		{'中', "CJK"},
		{'文', "CJK"},
		{'０', "Fullwidth"},
		{'１', "Fullwidth"},
		{'5', "Digit"},
		{'9', "Digit"},
		{'@', "Other"},
		{'こ', "CJK"},
		{'ン', "CJK"},
		{'안', "CJK"},
	}

	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			got := GetCharacterType(tt.r)
			if got != tt.want {
				t.Errorf("GetCharacterType(%c) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestAnalyzeText(t *testing.T) {
	tests := []struct {
		name           string
		text           string
		fontSize       float64
		wantTotal      int
		wantCJK        int
		wantLatin      int
		wantDigit      int
		wantFullwidth  int
		wantOther      int
	}{
		{
			name:      "English only",
			text:      "Hello",
			fontSize:  16,
			wantTotal: 5,
			wantLatin: 5,
		},
		{
			name:      "Chinese only",
			text:      "你好世界",
			fontSize:  16,
			wantTotal: 4,
			wantCJK:   4,
		},
		{
			name:      "Mixed",
			text:      "Hello世界123",
			fontSize:  16,
			wantTotal: 10,
			wantLatin: 5,
			wantCJK:   2,
			wantDigit: 3,
		},
		{
			name:          "Fullwidth",
			text:          "１２３",
			fontSize:      16,
			wantTotal:     3,
			wantFullwidth: 3,
		},
		{
			name:      "Japanese",
			text:      "こんにちは",
			fontSize:  16,
			wantTotal: 5,
			wantCJK:   5,
		},
		{
			name:      "Korean",
			text:      "안녕하세요",
			fontSize:  16,
			wantTotal: 5,
			wantCJK:   5,
		},
		{
			name:      "Empty text",
			text:      "",
			fontSize:  16,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeText(tt.text, tt.fontSize)
			
			if analysis.TotalChars != tt.wantTotal {
				t.Errorf("AnalyzeText().TotalChars = %v, want %v", analysis.TotalChars, tt.wantTotal)
			}
			if analysis.CJKChars != tt.wantCJK {
				t.Errorf("AnalyzeText().CJKChars = %v, want %v", analysis.CJKChars, tt.wantCJK)
			}
			if analysis.LatinChars != tt.wantLatin {
				t.Errorf("AnalyzeText().LatinChars = %v, want %v", analysis.LatinChars, tt.wantLatin)
			}
			if analysis.DigitChars != tt.wantDigit {
				t.Errorf("AnalyzeText().DigitChars = %v, want %v", analysis.DigitChars, tt.wantDigit)
			}
			if analysis.FullwidthChars != tt.wantFullwidth {
				t.Errorf("AnalyzeText().FullwidthChars = %v, want %v", analysis.FullwidthChars, tt.wantFullwidth)
			}
			if analysis.OtherChars != tt.wantOther {
				t.Errorf("AnalyzeText().OtherChars = %v, want %v", analysis.OtherChars, tt.wantOther)
			}
			
			if analysis.EstimatedWidth < 0 {
				t.Errorf("AnalyzeText().EstimatedWidth = %v, should be non-negative", analysis.EstimatedWidth)
			}
			if analysis.EstimatedHeight < 0 {
				t.Errorf("AnalyzeText().EstimatedHeight = %v, should be non-negative", analysis.EstimatedHeight)
			}
		})
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkDetectFontFormat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetectFontFormat("/path/to/font.ttf")
	}
}

func BenchmarkApproximateTextWidth(b *testing.B) {
	text := "Hello World"
	fontSize := 16.0

	for i := 0; i < b.N; i++ {
		ApproximateTextWidth(text, fontSize)
	}
}

func BenchmarkApproximateTextWidthChinese(b *testing.B) {
	text := "你好世界"
	fontSize := 16.0

	for i := 0; i < b.N; i++ {
		ApproximateTextWidth(text, fontSize)
	}
}

func BenchmarkApproximateTextWidthMixed(b *testing.B) {
	text := "Hello世界123"
	fontSize := 16.0

	for i := 0; i < b.N; i++ {
		ApproximateTextWidth(text, fontSize)
	}
}

func BenchmarkApproximateTextWidthSimple(b *testing.B) {
	text := "Hello世界123"
	fontSize := 16.0

	for i := 0; i < b.N; i++ {
		ApproximateTextWidthSimple(text, fontSize)
	}
}

func BenchmarkNormalizeFontName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NormalizeFontName("yahei")
	}
}

func BenchmarkIsFontFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsFontFile("font.ttf")
	}
}

func BenchmarkGetFontNameFromFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetFontNameFromFile("msyh.ttf")
	}
}

func BenchmarkScanFontDirectories(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchfont")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	
	path := filepath.Join(tmpDir, "test.ttf")
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ScanFontDirectories([]string{tmpDir})
	}
}

func BenchmarkAnalyzeText(b *testing.B) {
	text := "Hello世界123こんにちは안녕하세요"
	fontSize := 16.0

	for i := 0; i < b.N; i++ {
		AnalyzeText(text, fontSize)
	}
}

func BenchmarkGetCharacterType(b *testing.B) {
	chars := []rune{'A', '中', '０', '5', '@'}
	
	for i := 0; i < b.N; i++ {
		for _, r := range chars {
			GetCharacterType(r)
		}
	}
}

// ============================================================================
// Additional Edge Case Tests
// ============================================================================

func TestIsCJKCharacter(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'A', false},
		{'中', true},
		{'文', true},
		{'こ', true},
		{'ン', true},
		{'안', true},
		{'1', false},
		{'@', false},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			if got := isCJKCharacter(tt.r); got != tt.want {
				t.Errorf("isCJKCharacter(%c) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}

func TestIsFullwidthCharacter(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'A', false},
		{'１', true},
		{'２', true},
		{'３', true},
		{'！', true},
		{'＠', true},
		{'中', false},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.r), func(t *testing.T) {
			if got := isFullwidthCharacter(tt.r); got != tt.want {
				t.Errorf("isFullwidthCharacter(%c) = %v, want %v", tt.r, got, tt.want)
			}
		})
	}
}