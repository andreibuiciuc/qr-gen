package img

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"qr/qr-gen/util"

	"golang.org/x/exp/constraints"
)

type Image[T constraints.Integer] interface {
	CreateImage(filename string, encoded [][]T) image.Image
}

type QrImage struct{}

func New() Image[util.Module] {
	return &QrImage{}
}

func (qi *QrImage) CreateImage(filename string, encoded [][]util.Module) image.Image {
	topLeftPoint := image.Point{0, 0}
	bottomRightPoint := image.Point{len(encoded), len(encoded[0])}

	img := image.NewNRGBA(image.Rectangle{topLeftPoint, bottomRightPoint})

	for i := 0; i < len(encoded); i++ {
		for j := 0; j < len(encoded[i]); j++ {
			if util.IsModuleLighten(encoded[i][j]) {
				img.Set(j, i, color.White)
			} else {
				img.Set(j, i, color.Black)
			}
		}
	}

	f, _ := os.Create(filename)
	png.Encode(f, img)

	return img
}
