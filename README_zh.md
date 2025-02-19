# GameLife: QR码与康威生命游戏的结合

- [中文版本](README_zh.md)
- [English version](README.md)

## 概述

GameLife 是一个将 QR 码生成与康威生命游戏（Conway's Game of Life）相结合的项目。程序根据用户提供的签名生成 QR 码，将其转换为网格表示，并应用康威生命游戏的规则演化网格，直到达到稳定状态或完成预定义的迭代次数。

该项目展示了如何在 Go（Golang）中集成 QR 码生成、图像处理和细胞自动机。它还提供了保存网格中间状态为图像的功能。

---

## 功能

1. **QR 码生成**：
   - 根据用户提供的签名生成 QR 码。
   - 支持自定义 QR 码尺寸和输出格式（如 PNG、JPEG）。

2. **图像处理**：
   - 将 QR 码图像转换为二进制网格表示。
   - 检测 QR 码的模块大小以便进一步处理。

3. **康威生命游戏**：
   - 对从 QR 码生成的网格应用康威生命游戏的规则。
   - 演化网格，直到达到稳定状态或完成最大迭代次数。
   - 保存网格的中间状态为图像以供可视化。

4. **并发与稳定性检测**：
   - 使用线程安全的网格结构处理并发操作。
   - 使用哈希技术检测周期性振荡或稳定状态。

5. **可配置参数**：
   - 允许用户配置 QR 码尺寸、输出格式和最大迭代次数等参数。

---

## 安装

### 前置条件

- Go 1.20 或更高版本
- Git

### 步骤

1. 克隆仓库：
   ```bash
   git clone https://github.com/golrice/gamelife.git
   cd gamelife
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```

3. 构建项目：
   ```bash
   go build -o gamelife cmd/main.go
   ```

4. 运行可执行文件：
   ```bash
   ./gamelife -signature "您的签名" -format "png" -size 255
   ```

---

## 使用方法

### 命令行参数

| 参数         | 描述                             | 默认值 |
| ------------ | -------------------------------- | ------ |
| `-signature` | 用于生成 QR 码的用户签名         | 必填   |
| `-format`    | 输出图像格式（如 `png`、`jpeg`） | `png`  |
| `-size`      | QR 码的尺寸                      | `255`  |

### 示例

生成一个带有签名 "HelloWorld" 的 QR 码，保存为 PNG 文件，并应用康威生命游戏：

```bash
./gamelife -signature "HelloWorld" -format "png" -size 255
```

该命令将：
1. 为 "HelloWorld" 生成 QR 码。
2. 将 QR 码转换为网格。
3. 对网格应用康威生命游戏规则。
4. 保存中间状态和最终状态为 PNG 文件。

---

## 项目结构

项目按以下目录组织：

- **`cmd/`**: 包含主入口点 (`main.go`)。
- **`internal/config/`**: 负责配置管理。
- **`internal/engine/`**: 实现康威生命游戏逻辑和网格操作。
- **`internal/imageutil/`**: 提供图像处理和保存的工具。
- **`internal/qrcode/`**: 负责 QR 码生成。

---

## 关键组件

### 1. QR 码生成 (`internal/qrcode/qrcode.go`)

使用 [`skip2/go-qrcode`](https://github.com/skip2/go-qrcode) 库生成 QR 码。QR 码被表示为 `image.Image` 对象。

### 2. 图像处理 (`internal/imageutil/util.go`)

将图像转换为网格，检测模块大小，并以各种格式（如 PNG、JPEG）保存图像。

### 3. 康威生命游戏 (`internal/engine/grid.go`)

实现康威生命游戏的规则：
- 活细胞有 2 或 3 个活邻居时存活。
- 死细胞有恰好 3 个活邻居时复活。
- 所有其他细胞死亡或保持死亡。

网格会迭代演化，使用哈希技术检测周期性模式以判断是否达到稳定状态。

### 4. 配置管理 (`internal/config/config.go`)

提供集中式配置结构，默认值包括 QR 码尺寸、输出格式和迭代限制。

---

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。
