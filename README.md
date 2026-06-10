```markdown
# Pufo - 纯字体处理与操作工具

**字体管理和文本测量命令行工具**，支持字体扫描、搜索、信息查看和文本测量，专为字体开发者和内容创作者设计。

## ✨ 核心功能
1. **`info`** 显示字体详细信息
2. **`list`** 列出可用字体
3. **`measure`** 测量文本尺寸
4. **`search`** 搜索字体文件
5. **`version`** 显示版本信息
6. **`help`** 显示帮助信息

## 📌 完整参数说明

### 全局参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --font-dir | 自定义字体目录（可多次使用） | - | --font-dir /fonts |
| --system | 包含系统字体 | true | --system |
| --no-system | 排除系统字体 | false | --no-system |

### info 命令参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --font | 字体文件路径或系统字体名（必需） | - | /fonts/arial.ttf, Arial |
| --size | 字体大小（像素） | 16 | 24 |
| --details | 显示详细字体信息 | false | --details |

### list 命令参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --pattern | 按模式过滤字体 | - | "Arial", "YaHei" |
| --limit | 限制结果数量 | 20 | 10 |
| --font-dir | 自定义字体目录（覆盖全局） | - | /fonts |
| --system | 包含系统字体 | true | --system |
| --no-system | 排除系统字体 | false | --no-system |

### measure 命令参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --text | 要测量的文本（必需） | - | "Hello World" |
| --font | 字体文件路径或系统字体名（必需） | - | /fonts/arial.ttf |
| --size | 字体大小（像素） | 16 | 24 |
| --spacing | 行距倍数 | 1.2 | 1.5 |

### search 命令参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --name | 要搜索的字体名（必需） | - | "Arial" |
| --font-dir | 自定义字体目录（覆盖全局） | - | /fonts |
| --system | 包含系统字体 | true | --system |
| --no-system | 排除系统字体 | false | --no-system |

## 🚀 最常用命令示例

### 1. 基础用法
```bash
# 查看版本
pufo version

# 显示帮助
pufo help

# 列出系统字体
pufo list
```

### 2. 字体信息查看
```bash
# 查看系统字体信息
pufo info --font Arial --size 24

# 查看字体文件详细信息
pufo info --font /fonts/arial.ttf --details

# 查看内置字体信息
pufo info --font "" --size 16
```

### 3. 字体列表和搜索
```bash
# 列出系统字体（限制20个）
pufo list --limit 20

# 搜索特定字体模式
pufo list --pattern "YaHei"

# 搜索系统字体
pufo search --name "Arial"
```

### 4. Windows 字体目录支持
```bash
# 列出 Windows 字体目录
pufo list --font-dir /fonts --no-system

# 在 Windows 字体中搜索
pufo search --name "Arial" --font-dir /fonts --no-system

# 多目录扫描
pufo list --font-dir /windows/fonts --font-dir /usr/share/fonts
```

### 5. 文本测量
```bash
# 测量英文文本
pufo measure --text "Hello World" --font /fonts/arial.ttf --size 24

# 测量中文文本
pufo measure --text "你好世界" --font /fonts/msyh.ttc --size 32

# 测量多行文本
pufo measure --text "第一行\n第二行\n第三行" --font Arial --size 20 --spacing 1.5

# 测量混合文本
pufo measure --text "Hello 世界 123" --font /fonts/arial.ttf --size 24
```

### 6. 自定义字体目录
```bash
# 扫描自定义字体目录
pufo list --font-dir /fonts --pattern "sim" --no-system

# 搜索字体文件
pufo search --name "simhei" --font-dir /fonts --no-system

# 多个字体目录
pufo list --font-dir /fonts --font-dir /usr/local/share/fonts
```

### 7. 文本分析
```bash
# 分析文本组成
pufo measure --text "Hello 世界 123" --font Arial --size 24

# 查看字符类型分布
pufo measure --text "English 中文 123 !@#" --font Arial --size 16
```

## 📁 支持的字体格式

工具支持多种字体格式和加载方式：

| 类型 | 格式 | 说明 |
|------|------|------|
| **字体文件** | .ttf, .ttc, .otf | 直接加载字体文件 |
| **系统字体** | 字体名称 | 通过fc-list查找系统字体（Linux） |
| **内置字体** | GoRegular | 仅支持英文 |

### 系统字体查找示例（Linux）：
```bash
# 使用系统字体名称
pufo info --font "Arial" --size 16

# 查找中文字体
pufo search --name "YaHei" --system
```
## 🎯 文本测量详解

### 测量算法

文本宽度计算公式：
- **拉丁字符**：`字体大小 × 0.5`
- **CJK字符**（中日韩）：`字体大小 × 0.75`
- **数字**：`字体大小 × 0.5`
- **全角字符**：`字体大小 × 0.75`


## 🐳 容器化部署（Docker）

```sh
# 构建 运行阶段
# docker build --progress=plain -f Dockerfile.pufo --target runtime -t pufo .
docker build --progress=plain -f Dockerfile.pufo --target runtime -t ymc/pufo .

docker run -it --rm -v "$(pwd):/app" -v "/mnt/i/capture2:/data" -v "/mnt/c/Windows/Fonts:/fonts:ro" ymc/pufo bash
# ls /data
# ls /fonts
# ls script

pufo version
pufo help


# 安装字体
# apt-get update && apt-get install -y fonts-wqy-microhei
# ls /usr/share/fonts/truetype/wqy/wqy-microhei.ttc

# --font /data/fonts/simhei.ttf
# --font /usr/share/fonts/truetype/wqy/wqy-microhei.ttc
# fonts-wqy-zenhei

# fc-list :lang=zh
# WenQuanYi Zen Hei

pufo list --font-dir /fonts --no-system
pufo list --pattern YaHei --font-dir /fonts
pufo search --name "YaHei" --font-dir /fonts
pufo search --name "Arial" --font-dir /fonts --no-system
pufo measure --text "The quick brown fox jumps over the lazy dog" --font /fonts/arial.ttf --size 24
pufo measure --text "你好世界，这是粗体测试" --font /fonts/arialbd.ttf --size 32
```


## ⚠️ 注意事项

1. **系统字体依赖**：系统字体查找依赖 `fc-list` 命令，仅Linux/macOS支持
2. **Windows字体**：Windows用户需要直接挂载字体目录或使用字体文件路径
3. **中文支持**：需要确保字体文件包含中文字符集
4. **性能考虑**：扫描大量字体文件时可能需要较长时间，建议使用 `--limit` 限制结果
5. **文件权限**：确保有读取字体文件的权限
6. **路径格式**：Windows路径建议使用 `/mnt/c/...` 格式（WSL）或正斜杠

## 🛡️ 特性亮点

- ✅ **跨平台支持**：Linux、macOS、Windows（WSL）
- ✅ **多语言文本**：精确支持中、日、韩等多语言文本测量
- ✅ **自定义目录**：支持扫描自定义字体目录
- ✅ **详细分析**：提供字符类型分布统计
- ✅ **灵活过滤**：支持模式匹配和结果数量限制
- ✅ **容器友好**：完美支持Docker容器化部署
- ✅ **零配置**：开箱即用，无需复杂设置