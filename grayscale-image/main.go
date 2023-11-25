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

var numWorkers = 1

func main() {
	now := time.Now()
	file, err := os.Open("img-bigres.png")
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

	newImg := grayscale(img)
	png.Encode(outputFile, newImg)

	fmt.Printf("Elapsed: %v", time.Since(now))
}

type inputPixel struct {
	x int
	y int
	c color.Color
}

func grayscale(img image.Image) image.Image {

	newImg := image.NewRGBA(image.Rectangle{Min: img.Bounds().Min, Max: img.Bounds().Max})

	wg := sync.WaitGroup{}
	jobs := make(chan inputPixel, img.Bounds().Max.X*img.Bounds().Max.Y)
	results := make(chan inputPixel, img.Bounds().Max.X*img.Bounds().Max.Y)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	go func() {
		defer close(jobs)
		for x := 0; x < img.Bounds().Max.X; x++ {
			for y := 0; y < img.Bounds().Max.Y; y++ {
				r, g, b, a := img.At(x, y).RGBA()
				jobs <- inputPixel{
					x: x,
					y: y,
					c: color.RGBA{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: uint8(a),
					},
				}
			}
		}
	}()

	go func() {
		for pixel := range results {
			newImg.Set(pixel.x, pixel.y, pixel.c)
		}
	}()

	wg.Wait()

	return newImg
}

func processPixel(img image.Image) <-chan inputPixel {
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
					grey := rgbaToGrey(img.At(posX, posY))
					pixel := grey
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

func worker(i chan inputPixel, o chan inputPixel, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		pixel, more := <-i
		if !more {
			return
		}

		grey := rgbaToGrey(pixel.c)
		pixel.c = grey

		o <- pixel
	}
}

func rgbaToGrey(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	grey := luminosity(float64(r>>8), float64(g>>8), float64(b>>8))
	return color.RGBA{
		R: grey,
		G: grey,
		B: grey,
		A: uint8(a),
	}
}

func luminosity(r, g, b float64) uint8 {
	grey := uint8(r*0.21 + g*0.587 + b*0.114)
	return grey
}
