package rainfall

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func make2DMatrix(width, height int) [][]float64 {
	mat := make([][]float64, width) // width
	for x := range mat {
		mat[x] = make([]float64, height)
	}
	return mat

}

// mapRange map range to another range
func mapRange(v, v1, v2, min, max float64) float64 {
	return min + ((max-min)/(v2-v1))*(v-v1)
}

func imageToSlice(img image.Image) [][]float64 {

	// Convert image.Image to image.Gray
	grayImg := image.NewGray(img.Bounds())
	draw.Draw(grayImg, grayImg.Bounds(), img, image.Point{}, draw.Src)

	width := grayImg.Bounds().Size().X
	height := grayImg.Bounds().Size().Y

	matrix := make2DMatrix(width, height)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			grayValue := float64(grayImg.GrayAt(x, y).Y)
			matrix[x][y] = mapRange(grayValue, 0, 255, -1, 1)

		}
	}
	return matrix
}

func sliceToImage(m [][]float64) *image.Gray {
	width := len(m)
	height := len(m[0])
	grayImg := image.NewGray(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			v := m[x][y]
			grayValue := color.Gray{uint8(mapRange(v, -1, 1, 0, 255))}
			grayImg.SetGray(x, y, grayValue)
		}
	}
	return grayImg
}

func openImage(filename string) image.Image {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("os.Open failed: %v", err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("image.Decode failed: %v", err)
	}
	return img
}

func saveImage(filename string, img image.Image) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("os.Create failed: %v", err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Fatalf("png.Encode failed: %v", err)
	}
}
