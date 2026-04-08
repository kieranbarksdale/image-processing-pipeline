package worker

import (
	"bytes"
	"github.com/disintegration/imaging"
)

func ResizeImage(img imaging.Image, width, height int) ([]byte, error) {
	// img, err := imaging.Decode(bytes.NewReader(data))
	// if err != nil {
	// 	return nil, err
	// }
	
	// Resize the photo to 800x600px using the Lanczos filter
	// for the best quality.
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)
	
	// The resulting image will have the same dimensions as the original photo.
	// Save it to a file.
	err := imaging.Save(resizedImg, "resized-photo.jpg")
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}