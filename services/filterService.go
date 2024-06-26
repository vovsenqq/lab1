package services

import (
	"app/utils"
	"fmt"
	"image"
	"image/color"
	"math"
)

// ApplyFilter применяет фильтр к изображению
func ApplyFilter(img image.Image, filterType string, maskSize int) image.Image {
	// Получение пикселей изображения
	pixels, err := utils.GetPixels(img)
	if err != nil {
		fmt.Println("Error: Image could not be converted to pixels")
	}

	// Применение свертки к пикселям изображения
	convolvedPixels, err := convolve(pixels, filterType, maskSize)
	if err != nil {
		fmt.Println("Error: Image could not be convolved")
	}

	// Создание нового изображения
	newImg := image.NewRGBA(img.Bounds())

	// Заполнение нового изображения свернутыми пикселями
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			p := convolvedPixels[y][x]
			newImg.Set(x, y, color.RGBA{uint8(p.R), uint8(p.G), uint8(p.B), uint8(p.A)})
		}
	}

	return newImg
}

// convolve применяет свертку к пикселям изображения
func convolve(pixels [][]utils.Pixel, filterType string, maskSize int) ([][]utils.Pixel, error) {
	height := len(pixels)
	width := len(pixels[0])

	// Создание массива для свернутых пикселей
	convolvedPixels := make([][]utils.Pixel, height)
	for i := range convolvedPixels {
		convolvedPixels[i] = make([]utils.Pixel, width)
	}

	// Применение свертки к каждому пикселю изображения
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sumR1, sumG1, sumB1 := 0.0, 0.0, 0.0
			sumR2, sumG2, sumB2 := 0.0, 0.0, 0.0
			count := 0
			for ky := -maskSize / 2; ky <= maskSize/2; ky++ {
				for kx := -maskSize / 2; kx <= maskSize/2; kx++ {
					px, py := x+kx, y+ky
					if px >= 0 && px < width && py >= 0 && py < height {
						pixel := pixels[py][px]
						// Применение различных типов фильтров
						if filterType == "arithmetic" {
							sumR1 += float64(pixel.R)
							sumG1 += float64(pixel.G)
							sumB1 += float64(pixel.B)
						} else if filterType == "contraharmonic_black" {
							sumR1 += math.Pow(float64(pixel.R), 2)
							sumG1 += math.Pow(float64(pixel.G), 2)
							sumB1 += math.Pow(float64(pixel.B), 2)
							sumR2 += float64(pixel.R)
							sumG2 += float64(pixel.G)
							sumB2 += float64(pixel.B)
						} else if filterType == "contraharmonic_white" {
							if pixel.R == 0 {
								sumR2 += math.Pow(float64(3), -1)
							} else {
								sumR2 += math.Pow(float64(pixel.R), -1)
							}
							if pixel.G == 0 {
								sumG2 += math.Pow(float64(3), -1)
							} else {
								sumG2 += math.Pow(float64(pixel.G), -1)
							}
							if pixel.B == 0 {
								sumB2 += math.Pow(float64(3), -1)
							} else {
								sumB2 += math.Pow(float64(pixel.B), -1)
							}
							sumR1 += 1
							sumG1 += 1
							sumB1 += 1
						}
						count++
					}
				}
			}
			// Вычисление свернутого пикселя
			if filterType == "arithmetic" {
				convolvedPixels[y][x] = utils.Pixel{
					R: int(sumR1 / float64(count)),
					G: int(sumG1 / float64(count)),
					B: int(sumB1 / float64(count)),
					A: 255,
				}
			} else if filterType == "contraharmonic_black" || filterType == "contraharmonic_white" {
				convolvedPixels[y][x] = utils.Pixel{
					R: int(math.Min(math.Max(sumR1/sumR2, 0), 255)),
					G: int(math.Min(math.Max(sumG1/sumG2, 0), 255)),
					B: int(math.Min(math.Max(sumB1/sumB2, 0), 255)),
					A: 255,
				}
			}
		}
	}

	return convolvedPixels, nil
}
