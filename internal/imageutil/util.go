package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
	"github.com/golrice/gamelife/internal/config"
)

const (
	BLACK = true
	WHITE = false
)

type Image struct {
	Bitmap [][]bool
	Width  int
	Height int
}

func NewImage(bitmap [][]bool, width, height int) *Image {
	return &Image{
		Bitmap: bitmap,
		Width:  width,
		Height: height,
	}
}

func ToStandardImg(img *Image, config *config.Config) (image.Image, error) {
	if img == nil {
		return nil, fmt.Errorf("img is nil")
	}
	originContainer := image.NewRGBA(image.Rect(0, 0, img.Width, img.Height))
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			if img.Bitmap[y][x] == WHITE {
				originContainer.Set(x, y, color.White)
			}
		}
	}

	resizedImg := imaging.Resize(originContainer, config.QRSize, config.QRSize, imaging.NearestNeighbor)

	return resizedImg, nil
}

func SaveImg(img *Image, config *config.Config) (err error) {
	standardImg, err := ToStandardImg(img, config)
	if err != nil {
		return err
	}

	file, err := os.Create(config.Signature + "." + config.Format)
	if err != nil {
		return err
	}
	defer file.Close()

	switch config.Format {
	case "png":
		err = png.Encode(file, standardImg)
		if err != nil {
			return fmt.Errorf("PNG 编码失败: %w", err)
		}
	case "jpeg":
		err = jpeg.Encode(file, standardImg, nil) // 第三个参数是可选的编码选项
		if err != nil {
			return fmt.Errorf("JPEG 编码失败: %w", err)
		}
	default:
		return fmt.Errorf("不支持的格式: %s", config.Format)
	}

	return nil
}
