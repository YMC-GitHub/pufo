# Pufo - Pure Font Utility & Operations Tool

**A command-line tool for font management and text measurement**, supporting font scanning, searching, information viewing, and text measurement. Designed for font developers and content creators.

## ✨ Core Features
1. **`info`** Display detailed font information
2. **`list`** List available fonts
3. **`measure`** Measure text dimensions
4. **`search`** Search for font files
5. **`version`** Show version information
6. **`help`** Show help information

## 📌 Complete Parameter Reference

### Global Parameters

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| --font-dir | Custom font directory (can be used multiple times) | - | --font-dir /fonts |
| --system | Include system fonts | true | --system |
| --no-system | Exclude system fonts | false | --no-system |

### info Command Parameters

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| --font | Font file path or system font name (required) | - | /fonts/arial.ttf, Arial |
| --size | Font size in pixels | 16 | 24 |
| --details | Show detailed font information | false | --details |

### list Command Parameters

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| --pattern | Filter fonts by pattern | - | "Arial", "YaHei" |
| --limit | Limit number of results | 20 | 10 |
| --font-dir | Custom font directory (overrides global) | - | /fonts |
| --system | Include system fonts | true | --system |
| --no-system | Exclude system fonts | false | --no-system |

### measure Command Parameters

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| --text | Text to measure (required) | - | "Hello World" |
| --font | Font file path or system font name (required) | - | /fonts/arial.ttf |
| --size | Font size in pixels | 16 | 24 |
| --spacing | Line spacing multiplier | 1.2 | 1.5 |

### search Command Parameters

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| --name | Font name to search for (required) | - | "Arial" |
| --font-dir | Custom font directory (overrides global) | - | /fonts |
| --system | Include system fonts | true | --system |
| --no-system | Exclude system fonts | false | --no-system |

## 🚀 Most Common Command Examples

### 1. Basic Usage
```bash
# Show version
pufo version

# Show help
pufo help

# List system fonts
pufo list
```

### 2. Font Information Viewing
```bash
# View system font information
pufo info --font Arial --size 24

# View detailed font file information
pufo info --font /fonts/arial.ttf --details

# View built-in font information
pufo info --font "" --size 16
```

### 3. Font Listing and Searching
```bash
# List system fonts (limit to 20)
pufo list --limit 20

# Search for specific font patterns
pufo list --pattern "YaHei"

# Search system fonts
pufo search --name "Arial"
```

### 4. Windows Font Directory Support
```bash
# List Windows font directory
pufo list --font-dir /fonts --no-system

# Search within Windows fonts
pufo search --name "Arial" --font-dir /fonts --no-system

# Multi-directory scanning
pufo list --font-dir /windows/fonts --font-dir /usr/share/fonts
```

### 5. Text Measurement
```bash
# Measure English text
pufo measure --text "Hello World" --font /fonts/arial.ttf --size 24

# Measure Chinese text
pufo measure --text "你好世界" --font /fonts/msyh.ttc --size 32

# Measure multi-line text
pufo measure --text "Line1\nLine2\nLine3" --font Arial --size 20 --spacing 1.5

# Measure mixed text
pufo measure --text "Hello 世界 123" --font /fonts/arial.ttf --size 24
```

### 6. Custom Font Directories
```bash
# Scan custom font directory
pufo list --font-dir /fonts --pattern "sim" --no-system

# Search for font files
pufo search --name "simhei" --font-dir /fonts --no-system

# Multiple font directories
pufo list --font-dir /fonts --font-dir /usr/local/share/fonts
```

### 7. Text Analysis
```bash
# Analyze text composition
pufo measure --text "Hello 世界 123" --font Arial --size 24

# View character type distribution
pufo measure --text "English 中文 123 !@#" --font Arial --size 16
```

## 📁 Supported Font Formats

The tool supports multiple font formats and loading methods:

| Type | Format | Description |
|------|--------|-------------|
| **Font Files** | .ttf, .ttc, .otf | Direct font file loading |
| **System Fonts** | Font name | System font lookup via fc-list (Linux) |
| **Built-in Font** | GoRegular | English only |

### System Font Lookup Example (Linux):
```bash
# Use system font name
pufo info --font "Arial" --size 16

# Search for Chinese fonts
pufo search --name "YaHei" --system
```

## 🎯 Text Measurement Details

### Measurement Algorithm

Text width calculation formula:
- **Latin characters**: `font size × 0.5`
- **CJK characters** (Chinese, Japanese, Korean): `font size × 0.75`
- **Digits**: `font size × 0.5`
- **Fullwidth characters**: `font size × 0.75`

## 🐳 Docker Container Deployment

```sh
# Build runtime stage
# docker build --progress=plain -f Dockerfile.pufo --target runtime -t pufo .
docker build --progress=plain -f Dockerfile.pufo --target runtime -t ymc/pufo .

# Run container with mounted directories
docker run -it --rm \
  -v "$(pwd):/app" \
  -v "/mnt/i/capture2:/data" \
  -v "/mnt/c/Windows/Fonts:/fonts:ro" \
  ymc/pufo bash

# Inside container
# ls /data
# ls /fonts
# ls script

# Test commands
pufo version
pufo help

# Install fonts (if needed)
# apt-get update && apt-get install -y fonts-wqy-microhei
# ls /usr/share/fonts/truetype/wqy/wqy-microhei.ttc

# Font paths
# --font /data/fonts/simhei.ttf
# --font /usr/share/fonts/truetype/wqy/wqy-microhei.ttc
# fonts-wqy-zenhei

# List Chinese fonts
# fc-list :lang=zh
# WenQuanYi Zen Hei

# Example commands
pufo list --font-dir /fonts --no-system
pufo list --pattern YaHei --font-dir /fonts
pufo search --name "YaHei" --font-dir /fonts
pufo search --name "Arial" --font-dir /fonts --no-system
pufo measure --text "The quick brown fox jumps over the lazy dog" --font /fonts/arial.ttf --size 24
pufo measure --text "你好世界，这是粗体测试" --font /fonts/arialbd.ttf --size 32
```

## ⚠️ Important Notes

1. **System Font Dependency**: System font lookup depends on the `fc-list` command, only supported on Linux/macOS
2. **Windows Fonts**: Windows users need to directly mount font directories or use font file paths
3. **Chinese Support**: Ensure font files contain Chinese character sets
4. **Performance Considerations**: Scanning large numbers of font files may take time; use `--limit` to restrict results
5. **File Permissions**: Ensure you have read permissions for font files
6. **Path Format**: Windows paths should use `/mnt/c/...` format (WSL) or forward slashes

## 🛡️ Feature Highlights

- ✅ **Cross-platform Support**: Linux, macOS, Windows (WSL)
- ✅ **Multi-language Text**: Accurate support for Chinese, Japanese, Korean, and other multi-language text measurement
- ✅ **Custom Directories**: Support for scanning custom font directories
- ✅ **Detailed Analysis**: Character type distribution statistics
- ✅ **Flexible Filtering**: Pattern matching and result limit support
- ✅ **Container Friendly**: Perfect Docker container deployment support
- ✅ **Zero Configuration**: Ready to use out of the box

## 🔗 Related Links

- Project Home: [https://github.com/ymc-github/pufo](https://github.com/ymc-github/pufo)
- Issue Tracker: [Issues](https://github.com/ymc-github/pufo/issues)
- Changelog: [CHANGELOG.md](CHANGELOG.md)
