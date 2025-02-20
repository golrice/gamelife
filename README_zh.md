# 生命游戏

[![Go Version](https://img.shields.io/badge/go-1.20+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

GameLife 是一个基于 Conway's Game of Life 的项目，结合了二维码生成和图像处理功能。用户可以通过输入签名生成二维码，并将其转换为 Game of Life 的初始状态，观察其演化过程。

## 功能特性

- **二维码生成**：根据用户提供的签名生成二维码。
- **Game of Life 演化**：将二维码作为初始状态，运行 Conway's Game of Life 模拟。
- **图像保存**：支持将结果保存为 PNG 或 JPEG 格式的图像文件。
- **视频录制**（可选）：记录 Game of Life 的演化过程并生成视频文件。

---

## 目录结构

```
gamelife/
├── cmd/                # 主程序入口
│   └── main.go         # 主函数
├── internal/           # 内部模块
│   ├── config/         # 配置管理
│   ├── engine/         # Game of Life 引擎
│   ├── imageutil/      # 图像处理工具
│   └── qrcode/         # 二维码生成模块
└── Makefile            # 构建和管理任务
```

---

## 安装与运行

### 前置条件

- Go 1.20 或更高版本
- FFmpeg（如果需要生成视频）

### 执行指令

1. **基本命令**：
   ```bash
   make build    # 编译项目
   make run      # 运行程序（需要参数时使用 ARGS="..."）
   make clean    # 清理构建文件和生成的媒体文件
   make test     # 运行测试
   ```

2. **带参数运行**：
   ```bash
   make run ARGS="-signature=myseed -size=300 -video"
   ```

3. **交叉编译**：
   ```bash
   make build-linux    # 编译 Linux 版本
   make build-windows  # 编译 Windows 版本
   ```

4. **依赖检查**：
   ```bash
   make check-ffmpeg   # 检查视频生成依赖
   ```

5. **安装到系统路径**：
   ```bash
   make install        # 安装到 $GOPATH/bin
   ```

#### 参数说明

| 参数         | 默认值  | 描述                                  |
| ------------ | ------- | ------------------------------------- |
| `-signature` | 必填    | 用户签名，用于生成二维码              |
| `-format`    | `"png"` | 输出图像格式（支持 `png` 和 `jpeg`）  |
| `-size`      | `255`   | 图像大小（二维码尺寸）                |
| `-iter`      | `20`    | 最大迭代次数                          |
| `-video`     | `false` | 是否保存演化过程的视频（需要 FFmpeg） |

---

## 示例

生成一个带有签名的二维码，并运行 Game of Life 演化：

```bash
./bin/gamelife -signature "HelloGameLife" -format "png" -size 512 -iter 30 -video
```

输出文件：
- 初始二维码图像：`HelloGameLife.png`
- 演化后的图像：`HelloGameLife_after.png`
- 演化过程视频（如果启用）：`HelloGameLife.mp4`

---

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

希望你喜欢这个项目！🎮