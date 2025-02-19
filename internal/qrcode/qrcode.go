package qrcode

import (
	"github.com/golrice/gamelife/internal/config"
	qrcode "github.com/skip2/go-qrcode"
)

func GenerateQrcode(config *config.Config) ([][]bool, error) {
	qrc, err := qrcode.NewWithForcedVersion(config.Signature, 10, qrcode.High)

	return qrc.Bitmap(), err
}
