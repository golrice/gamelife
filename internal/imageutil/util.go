package imageutil

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/golrice/gamelife/internal/config"
)

func SaveImg(img image.Image, config *config.Config) (err error) {
	file, err := os.Create(config.Signature + "." + config.Format)
	if err != nil {
		return err
	}

	switch config.Format {
	case "png":
		err = png.Encode(file, img)
		if err != nil {
			return fmt.Errorf("PNG 编码失败: %w", err)
		}
	case "jpeg":
		err = jpeg.Encode(file, img, nil) // 第三个参数是可选的编码选项
		if err != nil {
			return fmt.Errorf("JPEG 编码失败: %w", err)
		}
	default:
		return fmt.Errorf("不支持的格式: %s", config.Format)
	}

	return nil
}

func isBlack(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	gray := (r + g + b) / 3
	return gray < 32768 // 灰度值小于阈值视为黑色
}

func DetectModuleSize(img image.Image) int {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 检查左上角的定位图案
	x, y := 0, 0
	for ; x < width; x++ {
		find := false
		for tempy := 0; tempy < height-1; tempy++ {
			if isBlack(img.At(x, tempy)) {
				y = tempy
				find = true
				break
			}
		}
		if find {
			break
		}
	}

	tempx := x

	for x < width && isBlack(img.At(x, y)) {
		x++
	}

	tempx = (tempx + x) / 2
	moduleSize := y // 定位图案的宽度即为模块大小
	for tempy := y; y < height; tempy++ {
		if !isBlack(img.At(tempx, tempy)) {
			moduleSize = tempy - y
			break
		}
	}

	return moduleSize
}
