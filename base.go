package base

import (
	"bytes"
	"image"
	"image/color"
	"image/png"

	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
)

type Encoder struct {
	AlphaTreshold int
	GreyTreshold  int
	QRLevel       qrcode.RecoveryLevel
}

var DefaultEncoder = Encoder{
	AlphaTreshold: 200,
	GreyTreshold:  30,
	QRLevel:       qrcode.Highest,
}

func Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	return DefaultEncoder.Encode(str, logo, size)
}

func (encoder Encoder) Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	logo = resizeImage(logo, 60, 60)

	code, err := qrcode.New(str, encoder.QRLevel)
	if err != nil {
		return nil, err
	}

	img := code.Image(size)
	encoder.overlayLogo(img, logo)

	err = png.Encode(&buffer, img)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}

func (encoder Encoder) overlayLogo(img, logo image.Image) {
	grey := uint32(^uint16(0)) * uint32(encoder.GreyTreshold) / 100
	alphaOffset := uint32(encoder.AlphaTreshold)
	offset := img.Bounds().Max.X/2 - logo.Bounds().Max.X/2

	for x := 0; x < logo.Bounds().Max.X; x++ {
		for y := 0; y < logo.Bounds().Max.Y; y++ {
			r, g, b, alpha := logo.At(x, y).RGBA()

			if alpha > alphaOffset {
				black := color.Black

				if r > grey && g > grey && b > grey {
					black = color.White
				}

				img.(*image.Paletted).Set(x+offset, y+offset, black)
			}

		}
	}
}

func resizeImage(src image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	return dst
}
