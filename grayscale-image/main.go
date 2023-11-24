package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"time"
)

func logError(err error) {
	fmt.Printf("%v\n", err)
}

func main() {
	now := time.Now()
	file, err := os.Open("img.png")
	if err != nil {
		logError(err)
		return
	}
	defer file.Close()

	outputFile, err := os.Create("gray.png")
	if err != nil {
		logError(err)
		return
	}
	defer outputFile.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		logError(err)
		return
	}

	newImg := image.NewRGBA(image.Rectangle{Min: img.Bounds().Min, Max: img.Bounds().Max})
	for i := range grayscale(img) {
		newImg.Set(i.x, i.y, i.c)
	}
	png.Encode(outputFile, newImg)

	fmt.Printf("Elapsed: %v", time.Since(now))
}

type inputPixel struct {
	x int
	y int
	c color.Color
}

func grayscale(img image.Image) <-chan inputPixel {
	out := make(chan inputPixel)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for x := 0; x < img.Bounds().Max.X; x++ {
			for y := 0; y < img.Bounds().Max.Y; y++ {
				wg.Add(1)
				posX, posY := x, y
				go func() {
					defer wg.Done()
					r, g, b, a := img.At(posX, posY).RGBA()
					grey := luminosity(float64(r>>8), float64(g>>8), float64(b>>8))
					pixel := color.RGBA{
						R: grey,
						G: grey,
						B: grey,
						A: uint8(a),
					}

					out <- inputPixel{
						x: posX,
						y: posY,
						c: pixel,
					}
				}()
			}
		}
	}()

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func luminosity(r, g, b float64) uint8 {
	grey := uint8(r*0.21 + g*0.587 + b*0.114)
	return grey
}
