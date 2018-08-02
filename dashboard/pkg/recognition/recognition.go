package recognition

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/machinebox/sdk-go/facebox"
)

// Faces type
type Faces []facebox.Face

func (a Faces) Len() int           { return len(a) }
func (a Faces) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Faces) Less(i, j int) bool { return a[i].Confidence < a[j].Confidence }

// SaveFaceImage : save found face into a file
func SaveFaceImage(f facebox.Face, buf []byte, namePrefix string) error {
	srcImg, decodeErr := png.Decode(bytes.NewReader(buf))
	if decodeErr != nil {
		return decodeErr
	}

	imgName := fmt.Sprintf("./faces/%s-%d.png", namePrefix, time.Now().Unix())

	// Crop face only
	rect := image.Rect(f.Rect.Left, f.Rect.Top, f.Rect.Left+f.Rect.Width, f.Rect.Top+f.Rect.Height)
	cropped := imaging.Crop(srcImg, rect)

	os.Mkdir("faces", os.ModePerm)
	outputFile, createErr := os.Create(imgName)
	if createErr != nil {
		return createErr
	}

	encodeErr := png.Encode(outputFile, cropped)
	if encodeErr != nil {
		return encodeErr
	}

	outputFile.Close()

	return nil
}

// Photo type
type Photo struct {
	Filename string
	ID       string
}

// GetUserPhotos : returns user photos from ./faces-db
func GetUserPhotos(slack string) ([]Photo, error) {
	pattern := fmt.Sprintf("./faces-db/%s*.png", slack)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var photos []Photo
	for _, m := range matches {
		filename := strings.Replace(m, "faces-db/", "", 1)
		id := strings.Replace(filename, slack+"-", "", 1)
		id = strings.Replace(id, ".png", "", 1)

		photos = append(photos, Photo{
			Filename: filename,
			ID:       id,
		})
	}

	return photos, nil
}
