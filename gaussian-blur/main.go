package main

import (
	"cmp"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

type matrix [][]float64

var gaussianKernel3x3 = matrix{
	[]float64{1.0 / 16.0, 1.0 / 8.0, 1.0 / 16.0},
	[]float64{1.0 / 8.0, 1.0 / 4.0, 1.0 / 8.0},
	[]float64{1.0 / 16.0, 1.0 / 8.0, 1.0 / 16.0},
}

var gaussianKernel8x8 = matrix{
	[]float64{1.0 / 256.0, 1.0 / 128.0, 1.0 / 256.0, 1.0 / 128.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 256.0, 1.0 / 128.0},
	[]float64{1.0 / 128.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0, 1.0 / 32.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0},
	[]float64{1.0 / 256.0, 1.0 / 128.0, 1.0 / 256.0, 1.0 / 128.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 256.0, 1.0 / 128.0},
	[]float64{1.0 / 128.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0, 1.0 / 32.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0},
	[]float64{1.0 / 64.0, 1.0 / 32.0, 1.0 / 64.0, 1.0 / 32.0, 1.0 / 16.0, 1.0 / 32.0, 1.0 / 64.0, 1.0 / 32.0},
	[]float64{1.0 / 128.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0, 1.0 / 32.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0},
	[]float64{1.0 / 256.0, 1.0 / 128.0, 1.0 / 256.0, 1.0 / 128.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 256.0, 1.0 / 128.0},
	[]float64{1.0 / 128.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0, 1.0 / 32.0, 1.0 / 64.0, 1.0 / 128.0, 1.0 / 64.0},
}

// TODO: Bruh, copy pasted from internet, find a way to generate this kernel
var gaussianKernel16x16 = matrix{
	{2.0e-06, 8.0e-06, 2.0e-05, 4.0e-05, 5.0e-05, 4.0e-05, 2.0e-05, 8.0e-06, 2.0e-06, 8.0e-06, 2.0e-05, 3.0e-05, 2.0e-05, 8.0e-06, 2.0e-06, 1.0e-07},
	{8.0e-06, 3.2e-05, 8.0e-05, 0.00016, 0.0002, 0.00016, 8.0e-05, 3.2e-05, 8.0e-06, 3.2e-05, 8.0e-05, 0.00012, 8.0e-05, 3.2e-05, 8.0e-06, 4.0e-07},
	{2.0e-05, 8.0e-05, 0.0002, 0.0004, 0.0005, 0.0004, 0.0002, 8.0e-05, 2.0e-05, 8.0e-05, 0.0002, 0.0003, 0.0002, 8.0e-05, 2.0e-05, 1.0e-06},
	{4.0e-05, 0.00016, 0.0004, 0.0008, 0.001, 0.0008, 0.0004, 0.00016, 4.0e-05, 0.00016, 0.0004, 0.0006, 0.0004, 0.00016, 4.0e-05, 2.0e-06},
	{5.0e-05, 0.0002, 0.0005, 0.001, 0.00125, 0.001, 0.0005, 0.0002, 5.0e-05, 0.0002, 0.0005, 0.00075, 0.0005, 0.0002, 5.0e-05, 2.5e-06},
	{4.0e-05, 0.00016, 0.0004, 0.0008, 0.001, 0.0008, 0.0004, 0.00016, 4.0e-05, 0.00016, 0.0004, 0.0006, 0.0004, 0.00016, 4.0e-05, 2.0e-06},
	{2.0e-05, 8.0e-05, 0.0002, 0.0004, 0.0005, 0.0004, 0.0002, 8.0e-05, 2.0e-05, 8.0e-05, 0.0002, 0.0003, 0.0002, 8.0e-05, 2.0e-05, 1.0e-06},
	{8.0e-06, 3.2e-05, 8.0e-05, 0.00016, 0.0002, 0.00016, 8.0e-05, 3.2e-05, 8.0e-06, 3.2e-05, 8.0e-05, 0.00012, 8.0e-05, 3.2e-05, 8.0e-06, 4.0e-07},
	{2.0e-06, 8.0e-06, 2.0e-05, 4.0e-05, 5.0e-05, 4.0e-05, 2.0e-05, 8.0e-06, 2.0e-06, 8.0e-06, 2.0e-05, 3.0e-05, 2.0e-05, 8.0e-06, 2.0e-06, 1.0e-07},
	{8.0e-06, 3.2e-05, 8.0e-05, 0.00016, 0.0002, 0.00016, 8.0e-05, 3.2e-05, 8.0e-06, 3.2e-05, 8.0e-05, 0.00012, 8.0e-05, 3.2e-05, 8.0e-06, 4.0e-07},
	{2.0e-05, 8.0e-05, 0.0002, 0.0004, 0.0005, 0.0004, 0.0002, 8.0e-05, 2.0e-05, 8.0e-05, 0.0002, 0.0003, 0.0002, 8.0e-05, 2.0e-05, 1.0e-06},
	{4.0e-05, 0.00016, 0.0004, 0.0008, 0.001, 0.0008, 0.0004, 0.00016, 4.0e-05, 0.00016, 0.0004, 0.0006, 0.0004, 0.00016, 4.0e-05, 2.0e-06},
	{5.0e-05, 0.0002, 0.0005, 0.001, 0.00125, 0.001, 0.0005, 0.0002, 5.0e-05, 0.0002, 0.0005, 0.00075, 0.0005, 0.0002, 5.0e-05, 2.5e-06},
	{4.0e-05, 0.00016, 0.0004, 0.0008, 0.001, 0.0008, 0.0004, 0.00016, 4.0e-05, 0.00016, 0.0004, 0.0006, 0.0004, 0.00016, 4.0e-05, 2.0e-06},
	{2.0e-05, 8.0e-05, 0.0002, 0.0004, 0.0005, 0.0004, 0.0002, 8.0e-05, 2.0e-05, 8.0e-05, 0.0002, 0.0003, 0.0002, 8.0e-05, 2.0e-05, 1.0e-06},
	{8.0e-06, 3.2e-05, 8.0e-05, 0.00016, 0.0002, 0.00016, 8.0e-05, 3.2e-05, 8.0e-06, 3.2e-05, 8.0e-05, 0.00012, 8.0e-05, 3.2e-05, 8.0e-06, 4.0e-07},
	{2.0e-06, 8.0e-06, 2.0e-05, 4.0e-05, 5.0e-05, 4.0e-05, 2.0e-05, 8.0e-06, 2.0e-06, 8.0e-06, 2.0e-05, 3.0e-05, 2.0e-05, 8.0e-06, 2.0e-06, 1.0e-07},
}

func logError(err error) {
	fmt.Printf("%v\n", err)
}

func testConvolution() {
	m := matrix{
		[]float64{0, 0, 0, 0, 0},
		[]float64{0, 0, 0, 0, 0},
		[]float64{0, 0, 0, 0, 0},
		[]float64{0, 0, 255, 0, 0},
		[]float64{0, 0, 0, 0, 0},
		[]float64{0, 0, 0, 0, 0},
		[]float64{0, 0, 0, 0, 0},
	}
	v := convolution(m, gaussianKernel3x3)

	fmt.Printf("%v\n", v)
}

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("Elapsed time: %v\n", time.Since(now))
	}()

	file, err := os.Open("in.png")
	if err != nil {
		logError(err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		logError(err)
		return
	}

	fmt.Printf("Dimension => x: %v y: %v\n", img.Bounds().Dx(), img.Bounds().Dy())

	newImg := gaussianBlur(img, gaussianKernel8x8)
	outputFile, err := os.Create("out.png")
	if err != nil {
		logError(err)
		return
	}
	png.Encode(outputFile, newImg)
}

func gaussianBlur(img image.Image, kernel matrix) image.Image {
	newImg := image.NewRGBA(img.Bounds())
	in := make(matrix, img.Bounds().Dx())
	for i := range in {
		in[i] = make([]float64, img.Bounds().Dy())
	}

	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grey := luminosity(float64(r>>8), float64(g>>8), float64(b>>8))
			in[x][y] = float64(grey)
		}
	}

	mat := convolution(in, kernel)
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			// _, _, _, a := img.At(x, y).RGBA()
			gr := uint8(mat[x][y])
			newImg.Set(x, y, color.Gray{Y: gr})
		}
	}

	return newImg
}

func convolution(in, k matrix) matrix {
	res := make(matrix, len(in))
	offset := len(k) / 2

	for x := 0; x < len(in); x++ {
		res[x] = make([]float64, len(in[x]))
		for y := 0; y < len(in[x]); y++ {
			var value float64 = 0
			for i := 0; i < len(k); i++ {
				for j := 0; j < len(k[i]); j++ {
					kx := clamp(x+i-offset, 0, len(in)-1)
					ky := clamp(y+j-offset, 0, len(in[x])-1)
					value += in[kx][ky] * k[i][j]
				}
			}
			res[x][y] = value
		}
	}

	fmt.Println(res)

	return res
}

func luminosity(r, g, b float64) uint8 {
	grey := uint8(r*0.21 + g*0.587 + b*0.114)
	return grey
}

func clamp[T cmp.Ordered](v T, min T, max T) T {
	switch {
	case v < min:
		return min
	case v > max:
		return max
	default:
		return v
	}
}
