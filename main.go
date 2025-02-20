package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Pix struct {
	x  int
	y  int
	cr uint8
	cg uint8
	cb uint8
}

type WorkItem struct {
	initialX int
	finalX   int
	initialY int
	finalY   int
}

const (
	posX   = -2
	posY   = -1.2
	height = 2.5

	imgWidth   = 1023
	imgHeight  = 1023
	pixelTotal = imgWidth * imgHeight

	maxIter = 500
	samples = 100 // 200

	numBlocks  = 128
	numThreads = 16

	ratio = float64(imgWidth) / float64(imgHeight)

	showProgress = true
	closeOnEnd   = false
)

var (
	img        *image.RGBA
	pixelCount int
)

func main() {
	pixelgl.Run(run)
}

func run() {
	log.Println("Initialise processing...")
	cfg := pixelgl.WindowConfig{
		Title:  "Parallel Mandelbrot in Go",
		Bounds: pixel.R(0, 0, imgWidth, imgHeight),
		VSync:  true,
		// Invisible: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pixelCount = 0
	img = image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	setWhiteBG(img)

	log.Println("Rendering...")
	start := time.Now()
	workBuffer := make(chan WorkItem, numBlocks)
	threadBuffer := make(chan bool, numThreads)
	drawBuffer := make(chan Pix, pixelTotal)

	go workBufferInit(workBuffer)
	go workersInit(drawBuffer, workBuffer, threadBuffer)
	go drawThread(drawBuffer)

	for !win.Closed() {
		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		if showProgress {
			fmt.Printf("\r%d/%d (%d%%)", pixelCount, pixelTotal, int(100*(float64(pixelCount)/float64(pixelTotal))))
		}

		if pixelCount == pixelTotal {
			end := time.Now()
			fmt.Println("\nFinished with time = ", end.Sub(start))
			pixelCount++

			if closeOnEnd {
				break
			}
		}
	}
}

func setWhiteBG(img *image.RGBA) {
	for x := 0; x < imgWidth-1; x++ {
		for y := 0; y < imgHeight-1; y++ {
			img.Set(x, y, color.White)
		}
	}
}

func workBufferInit(workBuffer chan WorkItem) {
	var sqrt = int(math.Sqrt(numBlocks))

	for i := 0; i < sqrt; i++ {
		for j := 0; j < sqrt; j++ {
			workItem := WorkItem{
				initialX: i * (imgWidth / sqrt),
				finalX:   (i + 1) * (imgWidth / sqrt),
				initialY: j * (imgHeight / sqrt),
				finalY:   (j + 1) * (imgHeight / sqrt),
			}
			workBuffer <- workItem
		}
	}
}

func workersInit(drawBuffer chan Pix, workBuffer chan WorkItem, threadBuffer chan bool) {
	// initialise threads
	for i := 1; i <= numThreads; i++ {
		threadBuffer <- true
	}

	// start for thread
	// and wait for another thread to finish
	for range threadBuffer {
		workItem := <-workBuffer
		go workerThread(workItem, drawBuffer, threadBuffer)
	}
}

func workerThread(workItem WorkItem, drawBuffer chan Pix, threadBuffer chan bool) {
	for x := workItem.initialX; x < workItem.finalX; x++ {
		for y := workItem.initialY; y < workItem.finalY; y++ {
			var colorR, colorG, colorB int
			for k := 0; k < samples; k++ {
				a := height*ratio*((float64(x)+RandFloat64())/float64(imgWidth)) + posX
				b := height*((float64(y)+RandFloat64())/float64(imgHeight)) + posY
				c := pixelColor(mandelbrotIteration(a, b, maxIter))
				colorR += int(c.R)
				colorG += int(c.G)
				colorB += int(c.B)
			}
			var cr, cg, cb uint8
			cr = uint8(float64(colorR) / float64(samples))
			cg = uint8(float64(colorG) / float64(samples))
			cb = uint8(float64(colorB) / float64(samples))

			drawBuffer <- Pix{
				x, y, cr, cg, cb,
			}

		}
	}
	threadBuffer <- true
}

func drawThread(drawBuffer chan Pix) {
	for i := range drawBuffer {
		img.SetRGBA(i.x, i.y, color.RGBA{R: i.cr, G: i.cg, B: i.cb, A: 255})
		pixelCount++
	}
}

func mandelbrotIteration(a, b float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + a
		y = 2*xy + b
	}

	return xx + yy, maxIter
}

func pixelColor(r float64, iter int) color.RGBA {
	insideSet := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// validar se está dentro do conjunto
	// https://pt.wikipedia.org/wiki/Conjunto_de_Mandelbrot
	if r > 4 {
		return hslToRGB(float64(iter)/100*r, 1, 0.5)
	}

	return insideSet
}
