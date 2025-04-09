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
var scaledWidth = 300
var scaledHeight = 150

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

	// Open the original image.
	src, err := imaging.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	bounds := src.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if imgWidth < scaledWidth || imgHeight < scaledHeight {
		log.Fatalf("Image is too small for a 240x480 crop.")
	}

	baseX := rand.Intn(imgWidth - scaledWidth + 1)
	baseY := rand.Intn(imgHeight - scaledHeight + 1)
	workingImg := imaging.Crop(src, image.Rect(baseX, baseY, baseX+scaledWidth, baseY+scaledHeight))
	wWidth, wHeight := scaledWidth, scaledHeight

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

	margin := 5
	blurX1 := posX - margin
	blurY1 := posY - margin
	blurX2 := posX + boxWidth + margin
	blurY2 := posY + boxHeight + margin

	if blurX1 < 0 {
		blurX1 = 0
	}
	if blurY1 < 0 {
		blurY1 = 0
	}
	if blurX2 > wWidth {
		blurX2 = wWidth
	}
	if blurY2 > wHeight {
		blurY2 = wHeight
	}
	blurRect := image.Rect(blurX1, blurY1, blurX2, blurY2)

	blurRegion := imaging.Crop(result, blurRect)
	blurred := imaging.Blur(blurRegion, 2.0)
	result = imaging.Paste(result, blurred, image.Pt(blurX1, blurY1))

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
