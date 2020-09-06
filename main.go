package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"runtime"
	"sync"
	"time"
)

func init() {
	// damn important or else At(), Bounds() functions will
	// caused memory pointer error!!
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
}

var width int
var height int

var wg sync.WaitGroup

func main() {

	start := time.Now()

	imgfile, err := os.Open("src/imageBlur/gray-wolf_thumb.jpg")

	if err != nil {
		fmt.Println("img.jpg file not found!")
		os.Exit(1)
	}

	defer imgfile.Close()

	// get image height and width with image/jpeg
	// change accordinly if file is png or gif

	imgCfg, _, err := image.DecodeConfig(imgfile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	width = imgCfg.Width
	height = imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)

	imgfile.Seek(0, 0)

	// get the image
	img, _, err := image.Decode(imgfile)

	b := img.Bounds()

	imgRGBA := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	//imgRGBA1 := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))

	radius := 20

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			colors := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			imgRGBA.Set(x, y, color.RGBA{colors.R, colors.G, colors.B, 255})

		}
	}

	runtime.GOMAXPROCS(20)

	goNumber := 20

	j := -(height / goNumber)
	for i := 0; (height/goNumber)+i < height; i += (height / goNumber) {
		wg.Add(1)
		go setPixel(imgRGBA, 0, (height/goNumber)+j, width, (height/goNumber)+i, radius)
		fmt.Printf("from: %d, to: %d\n", (height/goNumber)+j, (height/goNumber)+i)
		j += (height / goNumber)

	}
	//go setPixel(imgRGBA, 0, 0, width, height/2, radius)
	//go setPixel(imgRGBA, 0, height/2, width, height, radius)
	wg.Wait()

	outputImage, errI := os.Create("src/imageBlur/output.png")

	defer outputImage.Close()

	if errI != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	err = png.Encode(outputImage, imgRGBA)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(time.Since(start))

}

func findAverage(img *image.RGBA, radius, x, y int) color.RGBA {
	var r int
	var g int
	var b int

	//var a uint32

	var sum int = 1

	var minX = x - radius
	var maxX = x + radius
	var minY = y - radius
	var maxY = y + radius

	if minX < 0 {
		minX = 0
	}
	if maxX > width {
		maxX = width
	}
	if minY < 0 {
		minY = 0
	}
	if maxY > height {
		maxY = height
	}

	for i := minX; i <= maxX; i++ {
		for j := minY; j <= maxY; j++ {
			colors := color.RGBAModel.Convert(img.At(i, j)).(color.RGBA)

			r += int(colors.R)
			g += int(colors.G)
			b += int(colors.B)

			sum++
		}
	}

	newV := color.RGBA{uint8(r / sum), uint8(g / sum), uint8(b / sum), 255}
	sum = 0

	return newV
}

func setPixel(img *image.RGBA, minX, minY, width, height, radius int) {
	for y := minY; y < height; y++ {
		for x := minX; x < width; x++ {
			img.Set(x, y, findAverage(img, radius, x, y))
		}
	}
	wg.Done()
}
