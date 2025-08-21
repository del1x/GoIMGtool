package fileio

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"image/draw"

	"github.com/kolesa-team/go-webp/encoder"
)

type ImageEncoder interface {
	Encode(img image.Image, writer io.Writer, quality int) error
}

type JpegEncoder struct{}
type PngEncoder struct{}
type WebpEncoder struct{}

func (e *JpegEncoder) Encode(img image.Image, writer io.Writer, quality int) error {
	return jpeg.Encode(writer, img, &jpeg.Options{Quality: quality})
}

func (e *PngEncoder) Encode(img image.Image, writer io.Writer, quality int) error {
	imgNRGBA := image.NewNRGBA(img.Bounds())
	draw.Draw(imgNRGBA, imgNRGBA.Bounds(), img, image.Point{0, 0}, draw.Over)
	return png.Encode(writer, imgNRGBA)
}

func (e *WebpEncoder) Encode(img image.Image, writer io.Writer, quality int) error {
	qualityFloat := float32(quality)
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, qualityFloat)
	if err != nil {
		return fmt.Errorf("error creating encoder options: %v", err)
	}
	enc, err := encoder.NewEncoder(img, options)
	if err != nil {
		return fmt.Errorf("error creating encoder: %v", err)
	}
	return enc.Encode(writer)
}

func GetEncoder(format string) (ImageEncoder, error) {
	switch format {
	case "jpg", "jpeg":
		return &JpegEncoder{}, nil
	case "png":
		return &PngEncoder{}, nil
	case "webp":
		return &WebpEncoder{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
