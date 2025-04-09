package captcha

import (
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

var (
	folderPath   = "captcha/dataset"
	ScaledWidth  = 300
	ScaledHeight = 150

	fileTypes = []string{
		"*.jpg",
		"*.jpeg",
		"*.png",
	}

	// Cache file list and ensure it is only loaded once.
	fileList     []string
	fileListOnce sync.Once
)

type CaptchaImage struct {
	DataUri       string
	CorrectRegion image.Rectangle
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getFileList() {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		for _, fileType := range fileTypes {
			if matched, _ := filepath.Match(fileType, file.Name()); matched {
				fileList = append(fileList, filepath.Join(folderPath, file.Name()))
				break
			}
		}
	}

	if len(fileList) == 0 {
		log.Fatal("No valid files found in directory.")
	}
}

func getRandomInputFile() string {
	fileListOnce.Do(getFileList)
	return fileList[rand.Intn(len(fileList))]
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
		log.Fatalf("Image is too small for a %dx%d crop", ScaledWidth, ScaledHeight)
	}

	baseX := rand.Intn(imgWidth - ScaledWidth + 1)
	baseY := rand.Intn(imgHeight - ScaledHeight + 1)
	workingImg := imaging.Crop(src, image.Rect(baseX, baseY, baseX+ScaledWidth, baseY+ScaledHeight))

	boxWidth := rand.Intn(21) + 20
	boxHeight := rand.Intn(21) + 20

	if boxWidth > ScaledWidth {
		boxWidth = ScaledWidth
	}
	if boxHeight > ScaledHeight {
		boxHeight = ScaledHeight
	}

	maxX := ScaledWidth - boxWidth
	maxY := ScaledHeight - boxHeight
	posX := rand.Intn(maxX + 1)
	posY := rand.Intn(maxY + 1)

	regionRect := image.Rect(posX, posY, posX+boxWidth, posY+boxHeight)
	region := imaging.Crop(workingImg, regionRect)

	inverted := imaging.Invert(region)
	adjusted := imaging.AdjustContrast(inverted, 45.0)
	adjusted = imaging.AdjustSaturation(adjusted, 50)

	result := imaging.Paste(workingImg, adjusted, image.Pt(posX, posY))

	var buf bytes.Buffer
	if err = imaging.Encode(&buf, result, imaging.JPEG); err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	dataURI := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	return CaptchaImage{
		DataUri:       dataURI,
		CorrectRegion: image.Rect(posX, posY, posX+boxWidth, posY+boxHeight),
	}
}
