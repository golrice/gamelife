# Game of Life

[![Go Version](https://img.shields.io/badge/go-1.20+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

GameLife is a project based on Conway's Game of Life, incorporating QR code generation and image processing functionalities. Users can generate a QR code by inputting a signature, convert it into the initial state of the Game of Life, and observe its evolution process.

## Features

- **QR Code Generation**: Generate QR codes based on user-provided signatures.
- **Game of Life Evolution**: Use the QR code as the initial state to run Conway's Game of Life simulation.
- **Image Saving**: Supports saving results as PNG or JPEG image files.
- **Video Recording** (Optional): Record the evolution process of the Game of Life and generate a video file.

---

## Directory Structure

```
gamelife/
├── cmd/                # Main program entry
│   └── main.go         # Main function
├── internal/           # Internal modules
│   ├── config/         # Configuration management
│   ├── engine/         # Game of Life engine
│   ├── imageutil/      # Image processing utilities
│   └── qrcode/         # QR code generation module
└── Makefile            # Build and management tasks
```

---

## Installation and Execution

### Prerequisites

- Go 1.20 or higher
- FFmpeg (if video generation is required)

### Execution Instructions

1. **Basic Commands**:
   ```bash
   make build    # Compile the project
   make run      # Run the program (use ARGS="..." when parameters are needed)
   make clean    # Clean build files and generated media files
   make test     # Run tests
   ```

2. **Run with Parameters**:
   ```bash
   make run ARGS="-signature=myseed -size=300 -video"
   ```

3. **Cross-Compilation**:
   ```bash
   make build-linux    # Compile for Linux
   make build-windows  # Compile for Windows
   ```

4. **Dependency Check**:
   ```bash
   make check-ffmpeg   # Check video generation dependencies
   ```

5. **Install to System Path**:
   ```bash
   make install        # Install to $GOPATH/bin
   ```

#### Parameter Description

| Parameter    | Default  | Description                                                        |
| ------------ | -------- | ------------------------------------------------------------------ |
| `-signature` | Required | User signature, used to generate QR code                           |
| `-format`    | `"png"`  | Output image format (supports `png` and `jpeg`)                    |
| `-size`      | `255`    | Image size (QR code dimensions)                                    |
| `-iter`      | `20`     | Maximum number of iterations                                       |
| `-video`     | `false`  | Whether to save the evolution process as a video (requires FFmpeg) |

---

## Example

Generate a QR code with a signature and run the Game of Life evolution:

```bash
./bin/gamelife -signature "HelloGameLife" -format "png" -size 512 -iter 30 -video
```

Output files:
- Initial QR code image: `HelloGameLife.png`
- Evolved image: `HelloGameLife_after.png`
- Evolution process video (if enabled): `HelloGameLife.mp4`

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

Hope you enjoy this project! 🎮
