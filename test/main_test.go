package test

import (
	"flag"
	"fmt"
	"os/exec"
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
	savevideoFlag := false
	flag.Parse()

	config := config.NewConfig(signatureFlag, formatFlag, sizeFlag, maxiterFlag, savevideoFlag)

	bitmap, err := qrcode.GenerateQrcode(config)
	if err != nil {
		t.Error(err)
	}

	img := imageutil.NewImage(bitmap, len(bitmap), len(bitmap[0]))
	fmt.Println("img size: ", len(bitmap), len(bitmap[0]))

	err = imageutil.SaveImg(img, config)
	if err != nil {
		t.Error(err)
	}

	grid := engine.ImageToGrid(img)
	if ok, err := grid.ToStable(config); !ok || err != nil {
		fmt.Println("can't iter to stable")
	}

	img = grid.ToImage()
	config.Signature = config.Signature + "_after"
	imageutil.SaveImg(img, config)
}

func TestFfmpeg(t *testing.T) {
	cmd := exec.Command("ffmpeg")
	if err := cmd.Run(); err != nil {
		t.Error("can't run")
	}
}
