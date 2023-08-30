package filters

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

func ApplySobelFilter(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	grayImg := image.NewGray(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			grayImg.SetGray(x, y, gray)
		}
	}

	resultImg := image.NewGray(bounds)

	var wg sync.WaitGroup
	pixelCh := make(chan struct{ x, y int }, runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for pixel := range pixelCh {
				x, y := pixel.x, pixel.y
				if x > 0 && x < width-1 && y > 0 && y < height-1 {
					gx := float64(grayImg.GrayAt(x-1, y-1).Y) - float64(grayImg.GrayAt(x+1, y-1).Y) +
						2*float64(grayImg.GrayAt(x-1, y).Y) - 2*float64(grayImg.GrayAt(x+1, y).Y) +
						float64(grayImg.GrayAt(x-1, y+1).Y) - float64(grayImg.GrayAt(x+1, y+1).Y)

					gy := float64(grayImg.GrayAt(x-1, y-1).Y) + 2*float64(grayImg.GrayAt(x, y-1).Y) +
						float64(grayImg.GrayAt(x+1, y-1).Y) - float64(grayImg.GrayAt(x-1, y+1).Y) -
						2*float64(grayImg.GrayAt(x, y+1).Y) - float64(grayImg.GrayAt(x+1, y+1).Y)

					magnitude := math.Sqrt(gx*gx + gy*gy)
					c := uint8(magnitude)
					resultImg.SetGray(x, y, color.Gray{Y: c})
				}
			}
		}()
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			pixelCh <- struct{ x, y int }{x, y}
		}
	}

	close(pixelCh)
	wg.Wait()

	return resultImg
}
