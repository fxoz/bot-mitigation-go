package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

var fileTypes = []string{
	"*jpg",
	"*jpeg",
	"*png",
}

type CaptchaImage struct {
	dataUri       string
	correctRegion image.Rectangle
}

func getRandomInputFile() string {
	files, err := os.ReadDir("dataset")
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
				break // avoid duplicate entries for multiple matches
			}
		}
	}

	if len(validFiles) == 0 {
		log.Fatal("No valid files found in directory.")
	}

	return filepath.Join("dataset", validFiles[rand.Intn(len(validFiles))])
}

func GenerateImageCaptcha() CaptchaImage {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())

	inputFile := getRandomInputFile()

	src, err := imaging.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open image: %v\n", err)
		os.Exit(1)
	}

	bounds := src.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	fmt.Printf("Image dimensions: %d x %d\n", imgWidth, imgHeight)

	// Ensure image is large enough to pick a region.
	if imgWidth < 100 || imgHeight < 100 {
		fmt.Fprintln(os.Stderr, "Image is too small for a minimum random box of 50x50 pixels.")
		os.Exit(1)
	}

	// Choose a random region size between 50 and 100 pixels.
	boxWidth := rand.Intn(51) + 50  // 50 to 100 pixels width
	boxHeight := rand.Intn(51) + 50 // 50 to 100 pixels height

	// Calculate a random top-left position for the region.
	maxX := imgWidth - boxWidth
	maxY := imgHeight - boxHeight
	posX := rand.Intn(maxX + 1)
	posY := rand.Intn(maxY + 1)
	fmt.Printf("Random region: top-left (%d, %d) with size %dx%d\n", posX, posY, boxWidth, boxHeight)

	regionRect := image.Rect(posX, posY, posX+boxWidth, posY+boxHeight)
	region := imaging.Crop(src, regionRect)

	inverted := imaging.Invert(region)

	brightnessChange := rand.Intn(20) + 30
	if rand.Intn(2) == 0 {
		brightnessChange = -brightnessChange
	}
	adjusted := imaging.AdjustBrightness(inverted, float64(brightnessChange))

	result := imaging.Paste(src, adjusted, image.Pt(posX, posY))

	margin := 10
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
	if blurX2 > imgWidth {
		blurX2 = imgWidth
	}
	if blurY2 > imgHeight {
		blurY2 = imgHeight
	}
	blurRect := image.Rect(blurX1, blurY1, blurX2, blurY2)

	blurRegion := imaging.Crop(result, blurRect)
	blurred := imaging.Blur(blurRegion, 5.0)
	result = imaging.Paste(result, blurred, image.Pt(blurX1, blurY1))

	// err = imaging.Save(result, outputFile)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Failed to save image: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println("Processed image saved as:", outputFile)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, result, imaging.JPEG)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to encode image: %v\n", err)
		os.Exit(1)
	}

	dataURI := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())

	fmt.Printf("\nTime taken: %v\n", time.Since(start))

	return CaptchaImage{
		dataUri:       dataURI,
		correctRegion: image.Rect(blurX1, blurY1, blurX2, blurY2),
	}

}
