package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/color"
	"sync"

	"github.com/golrice/gamelife/internal/config"
)

type Grid struct {
	Width  int
	Height int
	Cells  [][]bool
	mu     sync.RWMutex // 用于并发安全
}

func NewGrid(width int, height int, cells [][]bool) *Grid {
	return &Grid{
		Width:  width,
		Height: height,
		Cells:  cells,
		mu:     sync.RWMutex{},
	}
}

func ImageToGrid(img image.Image) *Grid {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	cells := make([][]bool, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]bool, width)
		for x := 0; x < width; x++ {
			// 获取像素颜色
			r, g, b, _ := img.At(x, y).RGBA()
			gray := (r + g + b) / 3 // 计算灰度值
			// 黑色像素映射为 true，白色像素映射为 false
			cells[y][x] = gray < 32768
		}
	}

	return NewGrid(width, height, cells)
}

func (g *Grid) ToImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.Width, g.Height))

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Cells[y][x] {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	return img
}

func (g *Grid) ToStable(config *config.Config) bool {
	history := make([][][]bool, 0)

	for i := 0; i < config.MaxIter; i++ {
		if CheckStable(history) {
			return true
		}

		g.NextGeneration()
		history = append(history, g.Cells)
	}

	return false
}

func (g *Grid) GetCell(x, y int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if validPosition(x, y, g.Cells) {
		return g.Cells[y][x]
	}
	return false
}

func (g *Grid) SetCell(x, y int, value bool) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if y >= 0 && y < g.Height && x >= 0 && x < g.Width {
		g.Cells[y][x] = value
	}
}

func safeCell(cells [][]bool, x, y int, width, height int) bool {
	if !validPosition(x, y, cells) {
		return false
	}

	neighbors := 0

	for i := -1; i <= 1; i += 1 {
		if x+i < 0 || x+i >= width {
			continue
		}

		for j := -1; j <= 1; j += 1 {
			if y+j < 0 || y+j >= height {
				continue
			}

			if cells[y+j][x+i] {
				neighbors++
			}
		}
	}

	return 2 <= neighbors && neighbors <= 3
}

func (g *Grid) NextGeneration() {
	g.mu.Lock()
	defer g.mu.Unlock()

	cells := make([][]bool, g.Height)
	for y := range g.Cells {
		cells[y] = make([]bool, g.Width)
		copy(cells[y], g.Cells[y])
	}

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			g.Cells[y][x] = safeCell(cells, x, y, g.Width, g.Height)
		}
	}
}

func CheckStable(history [][][]bool) bool {
	// 检查周期性震荡
	hm := make(map[string]struct{})
	for _, g := range history {
		hs := MatrixToHash(g)
		if _, ok := hm[hs]; ok {
			return true
		}

		hm[hs] = struct{}{}
	}

	return false
}

func MatrixToHash(matrix [][]bool) string {
	var data []byte

	for _, row := range matrix {
		for _, value := range row {
			v := 0
			if value {
				v = 1
			}
			data = append(data, byte(v))
		}
	}

	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:])
}

func validPosition(x, y int, cells [][]bool) bool {
	height := len(cells)
	if x < 0 || y < 0 || y >= height || (height > 0 && x >= len(cells[0])) {
		return false
	}

	return true
}
