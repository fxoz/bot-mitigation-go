package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

var folderPath = "captcha/dataset"
var ScaledWidth = 300
var ScaledHeight = 150

var fileTypes = []string{
	"*.jpg",
	"*.jpeg",
	"*.png",
}

type CaptchaImage struct {
	DataUri       string
	CorrectRegion image.Rectangle
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomInputFile() string {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	var validFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		for _, fileType := range fileTypes {
			if matched, _ := filepath.Match(fileType, file.Name()); matched {
				validFiles = append(validFiles, file.Name())
				break
			}
		}
	}

	if len(validFiles) == 0 {
		log.Fatal("No valid files found in directory.")
	}

	return filepath.Join(folderPath, validFiles[rand.Intn(len(validFiles))])
}

func GenerateImageCaptcha() CaptchaImage {
	inputFile := getRandomInputFile()

	src, err := imaging.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	bounds := src.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if imgWidth < ScaledWidth || imgHeight < ScaledHeight {
		log.Fatalf("Image is too small for a 240x480 crop.")
	}

	baseX := rand.Intn(imgWidth - ScaledWidth + 1)
	baseY := rand.Intn(imgHeight - ScaledHeight + 1)
	workingImg := imaging.Crop(src, image.Rect(baseX, baseY, baseX+ScaledWidth, baseY+ScaledHeight))
	wWidth, wHeight := ScaledWidth, ScaledHeight

	boxWidth := rand.Intn(21) + 20
	boxHeight := rand.Intn(21) + 20

	if boxWidth > wWidth {
		boxWidth = wWidth
	}
	if boxHeight > wHeight {
		boxHeight = wHeight
	}

	maxX := wWidth - boxWidth
	maxY := wHeight - boxHeight
	posX := rand.Intn(maxX + 1)
	posY := rand.Intn(maxY + 1)

	regionRect := image.Rect(posX, posY, posX+boxWidth, posY+boxHeight)
	region := imaging.Crop(workingImg, regionRect)

	inverted := imaging.Invert(region)
	adjusted := imaging.AdjustContrast(inverted, 45.0)
	adjusted = imaging.AdjustSaturation(adjusted, 50)

	result := imaging.Paste(workingImg, adjusted, image.Pt(posX, posY))

	var buf bytes.Buffer
	err = imaging.Encode(&buf, result, imaging.JPEG)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	dataURI := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	return CaptchaImage{
		DataUri:       dataURI,
		CorrectRegion: image.Rect(posX, posY, posX+boxWidth, posY+boxHeight),
	}
}
