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
	saveVideoFlag := flag.Bool("video", false, "保存过程视频")
	flag.Parse()

	if signatureFlag == nil || *signatureFlag == "" {
		fmt.Println("no signature!")
		os.Exit(0)
	}

	fmt.Println("signature = ", *signatureFlag)

	config := config.NewConfig(*signatureFlag, *formatFlag, *sizeFlag, *maxiterFlag, *saveVideoFlag)

	bitmap, err := qrcode.GenerateQrcode(config)
	if err != nil {
		fmt.Println(err.Error())
	}
	if bitmap == nil {
		fmt.Println("empty bitmap")
		os.Exit(0)
	}

	img := imageutil.NewImage(bitmap, len(bitmap), len(bitmap[0]))
	fmt.Println("img size: ", len(bitmap), len(bitmap[0]))

	err = imageutil.SaveImg(img, config)
	if err != nil {
		fmt.Println(err.Error())
	}

	grid := engine.ImageToGrid(img)
	ok, err := grid.ToStable(config)
	if !ok {
		fmt.Println("can't iter to stable")
	}
	if err != nil {
		fmt.Println("some error occur when evolving")
		os.Exit(0)
	}

	img = grid.ToImage()
	config.Signature = config.Signature + "_after"
	imageutil.SaveImg(img, config)
}
