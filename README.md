# GameLife: QR Code and Conway's Game of Life Integration

- [中文版本](README_CN.md)
- [English version](README.md)

## Overview

GameLife is a project that combines QR code generation with Conway's Game of Life, a cellular automaton devised by mathematician John Conway. The program generates a QR code based on a user-provided signature, converts it into a grid representation, and applies the rules of Conway's Game of Life to evolve the grid until it reaches a stable state or completes a predefined number of iterations.

The project demonstrates how to integrate QR code generation, image processing, and cellular automata in Go (Golang). It also provides functionality to save intermediate states of the grid as images.

---

## Features

1. **QR Code Generation**:
   - Generates a QR code based on a user-provided signature.
   - Supports customization of QR code size and output format (e.g., PNG, JPEG).

2. **Image Processing**:
   - Converts the QR code image into a binary grid representation.
   - Detects the module size of the QR code for further processing.

3. **Conway's Game of Life**:
   - Applies the rules of Conway's Game of Life to the grid derived from the QR code.
   - Evolves the grid until it reaches a stable state or completes a maximum number of iterations.
   - Saves intermediate states of the grid as images for visualization.

4. **Concurrency and Stability Detection**:
   - Uses a thread-safe grid structure to handle concurrent operations.
   - Detects periodic oscillations or stable states using hashing techniques.

5. **Customizable Configuration**:
   - Allows users to configure parameters such as QR code size, output format, and maximum iterations.

---

## Installation

### Prerequisites

- Go 1.20 or higher
- Git

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/golrice/gamelife.git
   cd gamelife
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   go build -o gamelife cmd/main.go
   ```

4. Run the executable:
   ```bash
   ./gamelife -signature "YourSignature" -format "png" -size 255
   ```

---

## Usage

### Command-Line Arguments

| Flag         | Description                               | Default Value |
| ------------ | ----------------------------------------- | ------------- |
| `-signature` | User-defined signature for the QR code    | Required      |
| `-format`    | Output image format (e.g., `png`, `jpeg`) | `png`         |
| `-size`      | Size of the QR code                       | `255`         |

### Example

Generate a QR code with the signature "HelloWorld", save it as a PNG file, and apply Conway's Game of Life:

```bash
./gamelife -signature "HelloWorld" -format "png" -size 255
```

This will:
1. Generate a QR code for "HelloWorld".
2. Convert the QR code into a grid.
3. Apply Conway's Game of Life rules to the grid.
4. Save intermediate and final states as PNG files.

---

## Project Structure

The project is organized into the following directories:

- **`cmd/`**: Contains the main entry point (`main.go`).
- **`internal/config/`**: Handles configuration management.
- **`internal/engine/`**: Implements Conway's Game of Life logic and grid operations.
- **`internal/imageutil/`**: Provides utilities for image processing and saving.
- **`internal/qrcode/`**: Handles QR code generation.

---

## Key Components

### 1. QR Code Generation (`internal/qrcode/qrcode.go`)

Generates a QR code using the `skip2/go-qrcode` library. The QR code is represented as an `image.Image` object.

### 2. Image Processing (`internal/imageutil/util.go`)

Converts images to grids, detects module sizes, and saves images in various formats (PNG, JPEG).

### 3. Conway's Game of Life (`internal/engine/grid.go`)

Implements the rules of Conway's Game of Life:
- A live cell with 2 or 3 live neighbors survives.
- A dead cell with exactly 3 live neighbors becomes alive.
- All other cells die or remain dead.

The grid evolves iteratively, and stability is detected using hashing to identify periodic patterns.

### 4. Configuration Management (`internal/config/config.go`)

Provides a centralized configuration structure with default values for QR code size, output format, and iteration limits.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
