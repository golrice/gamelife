package engine

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os/exec"
	"sync"

	"github.com/disintegration/imaging"
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

	return NewGrid(width, height, BitmapCopy)
}

func (g *Grid) ToImage() *imageutil.Image {
	CellsCopy := make([][]bool, len(g.Cells))
	for i, row := range g.Cells {
		CellsCopy[i] = make([]bool, len(row))
		copy(CellsCopy[i], row)
	}

	return imageutil.NewImage(CellsCopy, g.Width, g.Height)
}

func (g *Grid) ToStandImage(config *config.Config) image.Image {
	originalContainer := image.NewRGBA(image.Rect(0, 0, g.Width, g.Height))
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Cells[y][x] == imageutil.WHITE {
				originalContainer.Set(x, y, color.White)
			}
		}
	}

	resizedImg := imaging.Resize(originalContainer, config.QRSize, config.QRSize, imaging.NearestNeighbor)
	return resizedImg
}

func (g *Grid) ToStable(config *config.Config) (bool, error) {
	historyHash := make(map[string]struct{}, 0)
	inputChannel := make(chan chanInArg, bufferSize)
	imageChannel := make(chan image.Image, bufferSize)
	var wg sync.WaitGroup
	var packet sync.WaitGroup

	wg.Add(consumerCount)
	for i := 0; i < consumerCount; i++ {
		go consumer(inputChannel, &wg, &packet)
	}

	if config.SaveVideo {
		wg.Add(1)
		go videoMaker(config.Signature, imageChannel, &wg)
	}

	for i := 0; i < config.MaxIter; i++ {
		if config.SaveVideo {
			standardImage := g.ToStandImage(config)
			imageChannel <- standardImage
		}

		g.NextGeneration(inputChannel, &packet)

		curHash := MatrixToHash(g.Cells)
		if _, ok := historyHash[curHash]; ok {
			close(inputChannel)
			close(imageChannel)
			wg.Wait()
			fmt.Println("after", i, "round, we reach stable")
			return true, nil
		}
		historyHash[curHash] = struct{}{}
	}

	close(inputChannel)
	close(imageChannel)
	wg.Wait()

	return false, nil
}

func videoMaker(videoName string, ch <-chan image.Image, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "image2pipe",
		"-vcodec", "png",
		"-r", "30",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-preset", "superfast",
		"-threads", "4",
		fmt.Sprintf(videoName+".mp4"),
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("无法获取 FFmpeg 标准输入管道:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("无法启动 FFmpeg:", err)
		return
	}

	writer := bufio.NewWriter(stdin)

	for img := range ch {
		if err := png.Encode(writer, img); err != nil {
			fmt.Println("写入图片失败")
			return
		}
	}

	if err := writer.Flush(); err != nil {
		fmt.Println("刷新缓冲区失败:", err)
		return
	}

	if err := stdin.Close(); err != nil {
		fmt.Println("关闭 FFmpeg 标准输入管道失败:", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("执行失败")
		return
	}
}

func consumer(input <-chan chanInArg, wg *sync.WaitGroup, packet *sync.WaitGroup) {
	defer wg.Done()

	for item := range input {
		item.OutCells[item.Y][item.X] = cellNextColor(item.PvCells, item.X, item.Y, item.Width, item.Height)
		packet.Done()
	}
}

func cellNextColor(cells [][]bool, x, y int, width, height int) bool {
	if !validPosition(x, y, cells) {
		return imageutil.WHITE
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

			if cells[y+j][x+i] == imageutil.BLACK {
				neighbors++
			}
		}
	}

	if cells[y][x] == imageutil.BLACK {
		neighbors--
	}

	if neighbors == 3 || (cells[y][x] && neighbors == 2) {
		return imageutil.BLACK
	}

	return imageutil.WHITE
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
