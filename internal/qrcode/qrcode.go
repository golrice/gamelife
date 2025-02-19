package qrcode

import (
	"image"

	"github.com/golrice/gamelife/internal/config"
	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQrcode(config *config.Config) (image.Image, error) {
	qrc, err := qrcode.New(config.Signature, qrcode.Low)

	return qrc.Image(config.QRSize), err
}
