package main

import (
	"flag"
	"fmt"

	"github.com/golrice/gamelife/internal/config"
	"github.com/golrice/gamelife/internal/engine"
	"github.com/golrice/gamelife/internal/imageutil"
	"github.com/golrice/gamelife/internal/qrcode"
)

func main() {
	signatureFlag := flag.String("signature", "", "用户的签名")
	formatFlag := flag.String("format", "png", "保存格式")
	sizeFlag := flag.Int("size", 255, "大小")
	flag.Parse()

	config := config.NewConfig(*signatureFlag, *formatFlag, *sizeFlag)

	img, err := qrcode.GenerateQrcode(config)
	if err != nil {
		panic(err)
	}

	err = imageutil.SaveImg(img, config)
	if err != nil {
		panic(err)
	}

	grid := engine.ImageToGrid(img)
	if ok := grid.ToStable(config); !ok {
		if err = fmt.Errorf("can't iter to stable"); err != nil {
			panic(err)
		}
	}

	img = grid.ToImage()
	config.Signature = config.Signature + "_after"
	imageutil.SaveImg(img, config)
}
