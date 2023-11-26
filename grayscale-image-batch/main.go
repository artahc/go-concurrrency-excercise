package main

import (
	"fmt"
	"image"
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
	defer func() {
		fmt.Printf("Elapsed: %v\n", time.Since(now))
	}()

	file, err := os.Open("in.png")
	if err != nil {
		logError(err)
		return
	}
	defer file.Close()

	outputFile, err := os.Create("out.png")
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
}

func grayscale(img image.Image) image.Image {
	newImg := image.NewNRGBA(img.Bounds())

	rgbaImg := img.(*image.NRGBA)
	fmt.Printf("Pixel count: %v\n", len(rgbaImg.Pix))
	fmt.Printf("Dimension: %v\n", img.Bounds().Max.X*img.Bounds().Dy()*4)

	numWorkers := 8

	in := make(chan []uint8, numWorkers)
	out := make(chan []uint8, numWorkers)
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	// m := sync.Mutex{}

	for i := 0; i < numWorkers; i++ {
		go worker(in, out, &wg)
	}

	pix := rgbaImg.Pix
	resultPix := make([]uint8, 0)

	batchSize := len(pix) / numWorkers
	fmt.Printf("Batch Size: %v\n", batchSize)
	fmt.Printf("-----------------------------------\n")

	go func() {
		defer close(in)
		for i := 0; i < numWorkers; i++ {
			fmt.Printf("Batch[%v] = [%v ... %v]\n", i, i*batchSize, (i*batchSize)+batchSize)
			pixels := pix[i*batchSize : (i*batchSize)+batchSize]

			in <- pixels
		}
	}()

	go func() {
		for pixels := range out {
			resultPix = append(resultPix, pixels...)
			wg.Done()
		}

	}()

	wg.Wait()
	newImg.Pix = resultPix
	fmt.Printf("Result length: %v\n", len(resultPix))
	return newImg
}

func worker(in <-chan []uint8, out chan []uint8, wg *sync.WaitGroup) {
	for pixels := range in {
		res := make([]uint8, 0)

		for i := 0; i < len(pixels)/4; i++ {
			r := pixels[i]
			g := pixels[i*4+1]
			b := pixels[i*4+2]
			a := pixels[i*4+3]
			p := []uint8{r, g, b, a}
			res = append(res, p...)
		}
		out <- res
	}
}
