package worker

import (
	"bytes"
	"image"
	"github.com/disintegration/imaging"
)

func ResizeImage(img image.Image, width, height int) ([]byte, error) {
	
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	var buf bytes.Buffer
	err := imaging.Encode(&buf, resizedImg, imaging.JPEG)
	if err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}