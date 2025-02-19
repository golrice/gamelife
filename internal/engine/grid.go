package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/golrice/gamelife/internal/config"
	"github.com/golrice/gamelife/internal/imageutil"
)

const (
	consumerCount = 5
	bufferSize    = 10
)

type Grid struct {
	Width   int
	Height  int
	Cells   [][]bool
	pvCells [][]bool
	mu      sync.RWMutex // 用于并发安全
}

type chanInArg struct {
	PvCells  [][]bool
	OutCells [][]bool
	X        int
	Y        int
	Width    int
	Height   int
}

func NewGrid(width int, height int, cells [][]bool) *Grid {
	cellsCopy := make([][]bool, len(cells))
	for i, row := range cells {
		cellsCopy[i] = make([]bool, len(row))
		copy(cellsCopy[i], row)
	}

	return &Grid{
		Width:   width,
		Height:  height,
		Cells:   cells,
		pvCells: cellsCopy,
		mu:      sync.RWMutex{},
	}
}

func ImageToGrid(img *imageutil.Image) *Grid {
	width, height := img.Width, img.Height

	BitmapCopy := make([][]bool, len(img.Bitmap))
	for i, row := range img.Bitmap {
		BitmapCopy[i] = make([]bool, len(row))
		copy(BitmapCopy[i], row)
	}

	for i := range BitmapCopy {
		for j := range BitmapCopy[i] {
			BitmapCopy[i][j] = !BitmapCopy[i][j]
		}
	}

	return NewGrid(width, height, BitmapCopy)
}

func (g *Grid) ToImage() *imageutil.Image {
	CellsCopy := make([][]bool, len(g.Cells))
	for i, row := range g.Cells {
		CellsCopy[i] = make([]bool, len(row))
		copy(CellsCopy[i], row)
	}

	for i := range CellsCopy {
		for j := range CellsCopy[i] {
			CellsCopy[i][j] = !CellsCopy[i][j]
		}
	}

	return imageutil.NewImage(CellsCopy, g.Width, g.Height)
}

func (g *Grid) ToStable(config *config.Config) bool {
	historyHash := make(map[string]struct{}, 0)
	inputChannel := make(chan chanInArg, bufferSize)
	var wg sync.WaitGroup
	var packet sync.WaitGroup
	defer close(inputChannel)

	wg.Add(consumerCount)
	for i := 0; i < consumerCount; i++ {
		go consumer(inputChannel, &wg, &packet)
	}

	for i := 0; i < config.MaxIter; i++ {
		g.NextGeneration(inputChannel, &packet)

		curHash := MatrixToHash(g.Cells)
		if _, ok := historyHash[curHash]; ok {
			fmt.Println("after", i, "round, we reach stable")
			return true
		}
		historyHash[curHash] = struct{}{}
	}

	return false
}

func consumer(input <-chan chanInArg, wg *sync.WaitGroup, packet *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		black := safeCell(item.PvCells, item.X, item.Y, item.Width, item.Height)
		item.OutCells[item.Y][item.X] = black
		packet.Done()
	}
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

	if cells[y][x] {
		neighbors--
	}

	return neighbors == 3 || (cells[y][x] && neighbors == 2)
}

func (g *Grid) NextGeneration(intput chan<- chanInArg, packet *sync.WaitGroup) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			packet.Add(1)
			intput <- chanInArg{
				PvCells:  g.pvCells,
				OutCells: g.Cells,
				X:        x,
				Y:        y,
				Width:    g.Width,
				Height:   g.Height,
			}
		}
	}

	packet.Wait()

	for i, row := range g.Cells {
		copy(g.pvCells[i], row)
	}

}

func MatrixToHash(matrix [][]bool) string {
	hasher := sha256.New()
	for _, row := range matrix {
		var packed uint64 = 0
		for i, value := range row {
			if value {
				packed |= 1 << (i % 64)
			}
			// 每 64 位写入一次
			if i%64 == 63 || i == len(row)-1 {
				hasher.Write([]byte{byte(packed), byte(packed >> 8), byte(packed >> 16), byte(packed >> 24),
					byte(packed >> 32), byte(packed >> 40), byte(packed >> 48), byte(packed >> 56)})
				packed = 0
			}
		}
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func validPosition(x, y int, cells [][]bool) bool {
	height := len(cells)
	if x < 0 || y < 0 || y >= height || (height > 0 && x >= len(cells[0])) {
		return false
	}

	return true
}
