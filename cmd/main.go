package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golrice/gamelife/internal/config"
	"github.com/golrice/gamelife/internal/engine"
	"github.com/golrice/gamelife/internal/imageutil"
	"github.com/golrice/gamelife/internal/qrcode"
)

func main() {
	signatureFlag := flag.String("signature", "", "用户的签名")
	formatFlag := flag.String("format", "png", "保存格式")
	sizeFlag := flag.Int("size", 255, "大小")
	maxiterFlag := flag.Int("iter", 20, "大小")
	flag.Parse()

	fmt.Println("signature = ", *signatureFlag)

	config := config.NewConfig(*signatureFlag, *formatFlag, *sizeFlag, *maxiterFlag)

	bitmap, err := qrcode.GenerateQrcode(config)
	if err != nil {
		panic(err)
	}
	if bitmap == nil {
		fmt.Println("empty bitmap")
		os.Exit(0)
	}

	img := imageutil.NewImage(bitmap, len(bitmap), len(bitmap[0]))
	fmt.Println("img size: ", len(bitmap), len(bitmap[0]))

	err = imageutil.SaveImg(img, config)
	if err != nil {
		panic(err)
	}

	grid := engine.ImageToGrid(img)
	if ok := grid.ToStable(config); !ok {
		fmt.Println("can't iter to stable")
	}

	img = grid.ToImage()
	config.Signature = config.Signature + "_after"
	imageutil.SaveImg(img, config)
}
