package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/cmplx"
	"os"
	"path/filepath"
	"strings"
)

// TODO:
// [X] Implement FFT pseudocode
// [ ] Take pixel brightness for each row of an image (save to matrix)
// [ ] Perform FFT on every row (save to matrix)

func main() {

	grayscale()

	img := getPixels()
	tf := make([][]complex128, len(img))

	// create fourier transformed image
	for i := 0; i < 512; i++ {
		y := make([]complex128, len(img[i]))
		fft(img[i], y, len(img[i]), 1)
		tf[i] = y
	}

	fourierImage(tf)
}

func fourierImage(tf [][]complex128) {
	imgPath := "image.jpg"

	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()
	img, _, _ := image.Decode(f)

	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		// and now loop thorough all of this x's y
		for y := 0; y < size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			// Offset colors a little, adjust it to your taste
			r := float64(real(tf[x][y])) * 0.92126
			g := float64(real(tf[x][y])) * 0.97152
			b := float64(real(tf[x][y])) * 0.90722
			// average
			grey := uint8((r + g + b) / 3)
			c := color.RGBA{
				R: grey, G: grey, B: grey, A: originalColor.A,
			}
			wImg.Set(x, y, c)
		}
	}

	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_fourier%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	check(err)
	defer fg.Close()
	err = jpeg.Encode(fg, wImg, nil)
	check(err)
}

func getPixels() [][]float64 {

	file, err := os.Open("image_gray.jpg")
	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}
	defer file.Close()

	img, _, err := image.Decode(file)

	if err != nil {
		return nil
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]float64
	for y := 0; y < height; y++ {
		var row []float64
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels
}

func rgbaToPixel(r, g, b, a uint32) float64 {
	return float64(r / 257)
}

func grayscale() {
	imgPath := "image.jpg"

	f, err := os.Open(imgPath)
	check(err)
	defer f.Close()
	img, _, _ := image.Decode(f)

	size := img.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)
	wImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		// and now loop thorough all of this x's y
		for y := 0; y < size.Y; y++ {
			pixel := img.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)
			// Offset colors a little, adjust it to your taste
			r := float64(originalColor.R) * 0.92126
			g := float64(originalColor.G) * 0.97152
			b := float64(originalColor.B) * 0.90722
			// average
			grey := uint8((r + g + b) / 3)
			c := color.RGBA{
				R: grey, G: grey, B: grey, A: originalColor.A,
			}
			wImg.Set(x, y, c)
		}
	}

	ext := filepath.Ext(imgPath)
	name := strings.TrimSuffix(filepath.Base(imgPath), ext)
	newImagePath := fmt.Sprintf("%s/%s_gray%s", filepath.Dir(imgPath), name, ext)
	fg, err := os.Create(newImagePath)
	check(err)
	defer fg.Close()
	err = jpeg.Encode(fg, wImg, nil)
	check(err)
}

func fft(x []float64, y []complex128, n, s int) {
	if n == 1 {
		y[0] = complex(x[0], 0)
		return
	}

	fft(x, y, n/2, 2*s)
	fft(x[s:], y[n/2:], n/2, 2*s)

	for k := 0; k < n/2; k++ {
		tf := cmplx.Rect(1, -2*math.Pi*float64(k)/float64(n)) * y[k+n/2]
		y[k], y[k+n/2] = y[k]+tf, y[k]-tf
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
