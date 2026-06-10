package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// "time"

	"github.com/ymc-github/apkm/ipkg/pufo"
)

const (
	Version = "v1.0.0"
	Usage   = `
pufo - Pure Font Utility & Operations

A command-line tool for font inspection, management, and text measurement.

Usage:
  pufo [command] [options]

Commands:
  info       Display font information
  list       List available system fonts
  measure    Measure text dimensions
  search     Search for fonts by name
  version    Show version information
  help       Show this help message

Global Options:
  --font-dir DIR  Custom font directory (can be used multiple times)
                  Examples: --font-dir /fonts
                           --font-dir /mnt/c/Windows/Fonts
  --system        Include system fonts (default: true)
  --no-system     Exclude system fonts

Options for 'info' command:
  --font PATH         Font file path or system font name (required)
  --size SIZE         Font size in pixels (default: 16)
  --details           Show detailed font information

Options for 'list' command:
  --pattern PATTERN   Filter fonts by pattern (e.g., "YaHei", "Arial")
  --limit N           Limit number of results (default: 20)
  --font-dir DIR      Custom font directory (overrides global)
  --system            Include system fonts (default: true)
  --no-system         Exclude system fonts

Options for 'measure' command:
  --text TEXT         Text to measure (required)
  --font PATH         Font file path or system font name (required)
  --size SIZE         Font size in pixels (default: 16)
  --spacing FLOAT     Line spacing multiplier (default: 1.2)

Options for 'search' command:
  --name NAME         Font name to search for (required)
  --font-dir DIR      Custom font directory (overrides global)
  --system            Include system fonts (default: true)
  --no-system         Exclude system fonts

Examples:
  # List system fonts
  pufo list
  
  # List fonts from Windows directory
  pufo list --font-dir /mnt/c/Windows/Fonts --no-system
  
  # Search in Windows fonts
  pufo search --name "YaHei" --font-dir /mnt/c/Windows/Fonts
  
  # Multiple font directories
  pufo list --font-dir /fonts --font-dir /usr/share/fonts
  
  # Display font info from file
  pufo info --font /mnt/c/Windows/Fonts/msyh.ttc
  
  # Measure text with custom font
  pufo measure --text "Hello 世界" --font /fonts/msyh.ttc --size 32
`
)

// Global font directories
var globalFontDirs []string
var includeSystem bool

func init() {
	// Parse global flags before command
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		if args[i] == "--font-dir" && i+1 < len(args) {
			globalFontDirs = append(globalFontDirs, args[i+1])
			i++
		} else if args[i] == "--system" {
			includeSystem = true
		} else if args[i] == "--no-system" {
			includeSystem = false
		} else if args[i] == "help" || args[i] == "version" {
			break
		} else if !strings.HasPrefix(args[i], "-") {
			// First non-flag is command, stop parsing global flags
			break
		}
	}
	
	// Default: include system fonts if not specified
	if !strings.Contains(strings.Join(args, " "), "--system") &&
	   !strings.Contains(strings.Join(args, " "), "--no-system") {
		includeSystem = true
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(Usage)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "info":
		infoCommand()
	case "list":
		listCommand()
	case "measure":
		measureCommand()
	case "search":
		searchCommand()
	case "version":
		fmt.Println(Version)
	case "help", "-h", "--help":
		fmt.Print(Usage)
	default:
		fmt.Printf("Error: Unknown command '%s'\n\n", command)
		fmt.Print(Usage)
		os.Exit(1)
	}
}

func infoCommand() {
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)
	fontSpec := infoCmd.String("font", "", "Font file path or system font name")
	fontSize := infoCmd.Float64("size", 16, "Font size in pixels")
	showDetails := infoCmd.Bool("details", false, "Show detailed font information")

	infoCmd.Parse(os.Args[2:])

	if *fontSpec == "" {
		fmt.Println("Error: --font is required")
		infoCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("🔍 Font Information\n")
	fmt.Printf("==================\n")
	fmt.Printf("Font: %s\n", *fontSpec)
	fmt.Printf("Size: %.0fpx\n", *fontSize)

	// Check if it's a file path
	if pufo.IsFontFile(*fontSpec) {
		fmt.Printf("Type: Font file\n")
		
		// Get file info
		if info, err := pufo.GetFontInfo(*fontSpec); err == nil {
			fmt.Printf("Format: %s\n", info.Format)
			if fileInfo, err := os.Stat(*fontSpec); err == nil {
				fmt.Printf("Size: %.2f KB\n", float64(fileInfo.Size())/1024)
			}
		}
	}

	// Try to load font
	face, isBuiltin, err := pufo.LoadFontSmart(*fontSpec, *fontSize)
	if err != nil {
		fmt.Printf("❌ Error loading font: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if face != nil {
			_ = face
		}
	}()

	if isBuiltin {
		fmt.Println("Source: Built-in font (goregular)")
	} else if pufo.IsFontFile(*fontSpec) {
		fmt.Printf("Source: Font file\n")
	} else {
		fmt.Printf("Source: System font\n")
	}

	// Get font metrics
	metrics := face.Metrics()
	fmt.Printf("\n📊 Font Metrics:\n")
	fmt.Printf("  Ascent:  %.0fpx\n", float64(metrics.Ascent)/64.0)
	fmt.Printf("  Descent: %.0fpx\n", float64(metrics.Descent)/64.0)
	fmt.Printf("  Height:  %.0fpx\n", float64(metrics.Height)/64.0)

	if *showDetails {
		fmt.Printf("\n📝 Additional Info:\n")
		fmt.Printf("  Can load: ✓ Yes\n")
	}
}

func listCommand() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	pattern := listCmd.String("pattern", "", "Filter fonts by pattern")
	limit := listCmd.Int("limit", 20, "Limit number of results")
	fontDirs := listCmd.String("font-dir", "", "Custom font directory (comma-separated)")
	useSystem := listCmd.Bool("system", includeSystem, "Include system fonts")
	noSystem := listCmd.Bool("no-system", false, "Exclude system fonts")

	listCmd.Parse(os.Args[2:])

	// Determine whether to include system fonts
	includeSys := *useSystem && !*noSystem
	if *noSystem {
		includeSys = false
	}

	// Collect font directories
	var dirs []string
	if *fontDirs != "" {
		dirs = append(dirs, strings.Split(*fontDirs, ",")...)
	}
	dirs = append(dirs, globalFontDirs...)

	fmt.Printf("📋 Font List\n")
	fmt.Printf("============\n")
	
	if len(dirs) > 0 {
		fmt.Printf("Custom directories:\n")
		for _, dir := range dirs {
			fmt.Printf("  - %s\n", dir)
		}
	}
	
	if includeSys {
		fmt.Printf("System fonts: ✓ Included\n")
	} else {
		fmt.Printf("System fonts: ✗ Excluded\n")
	}
	
	if *pattern != "" {
		fmt.Printf("Pattern: %s\n", *pattern)
	}
	fmt.Println()

	found := 0
	var results []string

	// Search custom directories
	if len(dirs) > 0 {
		fmt.Printf("🔍 Scanning custom directories...\n")
		fonts, err := pufo.SearchFontsInDirectories(dirs, *pattern)
		if err == nil {
			for _, font := range fonts {
				if found >= *limit {
					break
				}
				fmt.Printf("  📁 %s\n", font.Path)
				fmt.Printf("     Name: %s\n", font.Name)
				fmt.Printf("     Size: %.2f KB\n", float64(font.Size)/1024)
				fmt.Printf("     Format: %s\n", font.Format)
				found++
				results = append(results, font.Name)
			}
		}
	}

	// Search system fonts
	if includeSys && found < *limit {
		fmt.Printf("🔍 Scanning system fonts...\n")
		
		commonFonts := []string{
			"Arial", "Helvetica", "Times New Roman", "Courier New",
			"Microsoft YaHei", "SimHei", "SimSun", "KaiTi", "FangSong",
			"Verdana", "Tahoma", "Georgia", "Comic Sans MS",
		}

		for _, fontName := range commonFonts {
			if found >= *limit {
				break
			}
			
			if *pattern != "" && !strings.Contains(strings.ToLower(fontName), strings.ToLower(*pattern)) {
				continue
			}

			if pufo.IsSystemFontAvailable(fontName) {
				fmt.Printf("  ✓ %s (system)\n", fontName)
				found++
				results = append(results, fontName)
			}
		}
	}

	if found == 0 {
		fmt.Println("  No fonts found")
		fmt.Println("\n💡 Tips:")
		fmt.Println("  - Use --font-dir to specify custom font directories")
		fmt.Println("  - Mount Windows fonts: -v /mnt/c/Windows/Fonts:/fonts:ro")
		fmt.Println("  - Example: pufo list --font-dir /fonts")
	} else {
		fmt.Printf("\n✅ Found %d font(s)\n", found)
	}
}

func measureCommand() {
	measureCmd := flag.NewFlagSet("measure", flag.ExitOnError)
	text := measureCmd.String("text", "", "Text to measure")
	fontSpec := measureCmd.String("font", "", "Font file path or system font name")
	fontSize := measureCmd.Float64("size", 16, "Font size in pixels")
	lineSpacing := measureCmd.Float64("spacing", 1.2, "Line spacing multiplier")

	measureCmd.Parse(os.Args[2:])

	if *text == "" {
		fmt.Println("Error: --text is required")
		measureCmd.Usage()
		os.Exit(1)
	}

	if *fontSpec == "" {
		fmt.Println("Error: --font is required")
		measureCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("📏 Text Measurement\n")
	fmt.Printf("==================\n")
	fmt.Printf("Text: \"%s\"\n", *text)
	fmt.Printf("Font: %s\n", *fontSpec)
	fmt.Printf("Size: %.0fpx\n", *fontSize)

	// Check if font file exists
	if pufo.IsFontFile(*fontSpec) {
		if info, err := os.Stat(*fontSpec); err == nil {
			fmt.Printf("Font file: %s (%.2f KB)\n", filepath.Base(*fontSpec), float64(info.Size())/1024)
		}
	}

	// Load font
	face, _, err := pufo.LoadFontSmart(*fontSpec, *fontSize)
	if err != nil {
		fmt.Printf("❌ Error loading font: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if face != nil {
			_ = face
		}
	}()

	// Get approximate measurements
	approxWidth := pufo.ApproximateTextWidth(*text, *fontSize)
	approxHeight := pufo.ApproximateTextHeight(*fontSize)
	lineHeight := pufo.CalculateLineHeight(*fontSize, *lineSpacing)

	fmt.Printf("\n📊 Measurements:\n")
	fmt.Printf("  Approximate width:  %.0fpx\n", approxWidth)
	fmt.Printf("  Approximate height: %.0fpx\n", approxHeight)
	fmt.Printf("  Line height:        %.0fpx (spacing: %.1f)\n", lineHeight, *lineSpacing)

	// Character count analysis
	charCount := len([]rune(*text))
	wordCount := len(strings.Fields(*text))
	lineCount := strings.Count(*text, "\n") + 1

	fmt.Printf("\n📝 Text Analysis:\n")
	fmt.Printf("  Characters: %d\n", charCount)
	fmt.Printf("  Words:      %d\n", wordCount)
	fmt.Printf("  Lines:      %d\n", lineCount)

	// Estimated area
	if lineCount > 1 {
		totalHeight := lineHeight * float64(lineCount)
		fmt.Printf("\n📐 Estimated bounding box:\n")
		fmt.Printf("  Width:  %.0fpx\n", approxWidth)
		fmt.Printf("  Height: %.0fpx\n", totalHeight)
		fmt.Printf("  Area:   %.0f px²\n", approxWidth*totalHeight)
	}
}

func searchCommand() {
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	fontName := searchCmd.String("name", "", "Font name to search for")
	fontDirs := searchCmd.String("font-dir", "", "Custom font directory (comma-separated)")
	useSystem := searchCmd.Bool("system", includeSystem, "Include system fonts")
	noSystem := searchCmd.Bool("no-system", false, "Exclude system fonts")

	searchCmd.Parse(os.Args[2:])

	if *fontName == "" {
		fmt.Println("Error: --name is required")
		searchCmd.Usage()
		os.Exit(1)
	}

	fmt.Printf("🔎 Searching for font: %s\n", *fontName)
	fmt.Printf("========================\n")

	// Determine whether to include system fonts
	includeSys := *useSystem && !*noSystem
	if *noSystem {
		includeSys = false
	}

	// Collect font directories
	var dirs []string
	if *fontDirs != "" {
		dirs = append(dirs, strings.Split(*fontDirs, ",")...)
	}
	dirs = append(dirs, globalFontDirs...)

	found := false

	// Search in custom directories
	if len(dirs) > 0 {
		fmt.Printf("📁 Searching in custom directories:\n")
		for _, dir := range dirs {
			fmt.Printf("  - %s\n", dir)
		}
		
		fonts, err := pufo.SearchFontsInDirectories(dirs, *fontName)
		if err == nil && len(fonts) > 0 {
			fmt.Printf("\n✅ Found %d match(es) in custom directories:\n", len(fonts))
			for _, font := range fonts {
				fmt.Printf("\n  📄 %s\n", filepath.Base(font.Path))
				fmt.Printf("     Full path: %s\n", font.Path)
				fmt.Printf("     Name: %s\n", font.Name)
				fmt.Printf("     Size: %.2f KB\n", float64(font.Size)/1024)
				fmt.Printf("     Format: %s\n", font.Format)
				found = true
			}
		}
	}

	// Search in system fonts
	if includeSys && !found {
		fmt.Printf("\n🖥️  Searching in system fonts...\n")
		
		normalized := pufo.NormalizeFontName(*fontName)
		if normalized != *fontName {
			fmt.Printf("Normalized name: %s\n", normalized)
		}

		if pufo.IsSystemFontAvailable(normalized) {
			fmt.Printf("\n✅ Font '%s' is available on this system\n", normalized)
			
			if path, err := pufo.SearchSystemFont(normalized); err == nil {
				fmt.Printf("📁 Font path: %s\n", path)
			}
			
			// Try to load and show metrics
			face, _, err := pufo.LoadFontSmart(normalized, 16)
			if err == nil && face != nil {
				metrics := face.Metrics()
				fmt.Printf("\n📊 Metrics (at 16px):\n")
				fmt.Printf("  Ascent:  %.0fpx\n", float64(metrics.Ascent)/64.0)
				fmt.Printf("  Descent: %.0fpx\n", float64(metrics.Descent)/64.0)
				fmt.Printf("  Height:  %.0fpx\n", float64(metrics.Height)/64.0)
			}
			found = true
		}
	}

	if !found {
		fmt.Printf("\n❌ Font '%s' not found\n", *fontName)
		fmt.Printf("\n💡 Tips:\n")
		fmt.Printf("  - Use --font-dir to specify custom font directories\n")
		fmt.Printf("  - Mount Windows fonts: -v /mnt/c/Windows/Fonts:/fonts:ro\n")
		fmt.Printf("  - Example: pufo search --name YaHei --font-dir /fonts\n")
		
		if len(dirs) == 0 && includeSys {
			fmt.Printf("\n  Try searching in custom directories:\n")
			fmt.Printf("    docker run -v /mnt/c/Windows/Fonts:/fonts:ro pufo search --name %s --font-dir /fonts\n", *fontName)
		}
	}
}

// Helper function to format file size
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}