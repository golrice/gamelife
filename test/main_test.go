package test

import (
	"flag"
	"fmt"
	"testing"

	"github.com/golrice/gamelife/internal/config"
	"github.com/golrice/gamelife/internal/engine"
	"github.com/golrice/gamelife/internal/imageutil"
	"github.com/golrice/gamelife/internal/qrcode"
)

func TestMain(t *testing.T) {
	signatureFlag := "123"
	formatFlag := "png"
	sizeFlag := 255
	maxiterFlag := 100000
	flag.Parse()

	config := config.NewConfig(signatureFlag, formatFlag, sizeFlag, maxiterFlag)

	bitmap, err := qrcode.GenerateQrcode(config)
	if err != nil {
		panic(err)
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
