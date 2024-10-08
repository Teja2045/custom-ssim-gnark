package custom_ssim

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	"github.com/arka-labs/ssim/circuit"
	"golang.org/x/image/draw"
)

// Default SSIM constants
var (
	L  = 255.0
	K1 = 0.01
	K2 = 0.03
	C1 = math.Pow((K1 * L), 2.0)
	C2 = math.Pow((K2 * L), 2.0)
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Given a path to an image file, read and return as
// an image.Image
func readImage(fname string) image.Image {
	file, err := os.Open(fname)
	handleError(err)
	defer file.Close()

	img, _, err := image.Decode(file)
	handleError(err)
	return img
}

// Resize an image to the specified dimensions
func resizeImage(img image.Image, width, height int) image.Image {
	resizedImg := image.NewGray(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return resizedImg
}

// Convert an Image to grayscale which
// equalizes RGB values
func convertToGray(originalImg image.Image) image.Image {
	bounds := originalImg.Bounds()
	w, h := dim(originalImg)

	grayImg := image.NewGray(bounds)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			originalColor := originalImg.At(x, y)
			grayColor := color.GrayModel.Convert(originalColor)
			grayImg.Set(x, y, grayColor)
		}
	}

	return grayImg
}

// Write an image.Image to a jpg file of quality 100
func WriteImage(img image.Image, path string) {
	w, err := os.Create(path + ".jpg")
	handleError(err)
	defer w.Close()

	quality := jpeg.Options{Quality: 100}
	jpeg.Encode(w, img, &quality)
}

// Convert uint32 R value to a float. The returning
// float will have a range of 0-255
func getPixVal(c color.Color) float64 {
	r, _, _, _ := c.RGBA()
	return float64(r >> 8)
}

// Helper function that return the dimension of an image
func dim(img image.Image) (w, h int) {
	w, h = img.Bounds().Max.X, img.Bounds().Max.Y
	return
}

// Check if two images have the same dimension
func equalDim(img1, img2 image.Image) bool {
	w1, h1 := dim(img1)
	w2, h2 := dim(img2)
	return (w1 == w2) && (h1 == h2)
}

// Given an Image, calculate the mean of its
// pixel values
func mean(img image.Image) float64 {
	w, h := dim(img)

	fmt.Println("width, height: ", w, h)
	n := float64((w * h) - 1)
	sum := 0.0

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			sum += getPixVal(img.At(x, y))
		}
	}
	return sum / n
}

// Compute the standard deviation with pixel values of Image
func stdev(img image.Image) float64 {
	w, h := dim(img)

	n := float64((w * h) - 1)
	sum := 0.0
	avg := mean(img)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			pix := getPixVal(img.At(x, y))
			sum += math.Pow((pix - avg), 2.0)
		}
	}
	return math.Sqrt(sum / n)
}

// Calculate the covariance of 2 images
func covar(img1, img2 image.Image) (c float64, err error) {
	if !equalDim(img1, img2) {
		err = errors.New("images must have the same dimension")
		return
	}
	avg1 := mean(img1)
	avg2 := mean(img2)
	w, h := dim(img1)
	sum := 0.0
	n := float64((w * h) - 1)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			pix1 := getPixVal(img1.At(x, y))
			pix2 := getPixVal(img2.At(x, y))

			sum += (pix1 - avg1) * (pix2 - avg2)
		}
	}
	c = sum / n
	return
}

func ssim(x, y image.Image) float64 {
	avg_x := mean(x)
	avg_y := mean(y)

	stdev_x := stdev(x)
	stdev_y := stdev(y)

	cov, err := covar(x, y)
	handleError(err)

	fmt.Println("values", avg_x, avg_y, stdev_x, stdev_y, cov)
	numerator := ((2.0 * avg_x * avg_y) + C1) * ((2.0 * cov) + C2)
	denominator := (math.Pow(avg_x, 2.0) + math.Pow(avg_y, 2.0) + C1) *
		(math.Pow(stdev_x, 2.0) + math.Pow(stdev_y, 2.0) + C2)

	ssim := numerator / denominator
	check(C1, C2, ssim, avg_x, avg_y, stdev_x, stdev_y, cov)
	return ssim
}

// Compare two images by resizing them to 100x100, converting them to grayscale, and calculating SSIM
func Compare2(image1, image2 string) {
	fmt.Println("Loading and resizing images...")

	img1 := readImage(image1)
	img2 := readImage(image2)

	resizedImg1 := resizeImage(img1, 100, 100)
	resizedImg2 := resizeImage(img2, 100, 100)

	grayImg1 := convertToGray(resizedImg1)
	grayImg2 := convertToGray(resizedImg2)

	fmt.Println("Calculating SSIM...")
	index := ssim(grayImg1, grayImg2)

	fmt.Printf("SSIM = %f\n", index)
}

func check(C1, C2, ssim, avg_x, avg_y, stdev_x, stdev_y, cov float64) {
	C1_3 := int64(C1 * 1000)
	C2_3 := int64(C2 * 1000)
	C1_6 := int64(C1 * 1000000)
	ssim_3 := int64(ssim * 1000)
	cov_3 := int64(cov * 1000)

	avg_x_3 := int64(avg_x * 1000)
	avg_x_sq := avg_x * avg_x
	avg_x_sq_3 := int64(avg_x_sq * 1000)

	avg_y_3 := int64(avg_y * 1000)
	avg_y_sq := avg_y * avg_y
	avg_y_sq_3 := int64(avg_y_sq * 1000)

	stdev_x_sq := stdev_x * stdev_x
	stdev_x_sq_3 := int64(stdev_x_sq * 1000)

	stdev_y_sq := stdev_y * stdev_y
	stdev_y_sq_3 := int64(stdev_y_sq * 1000)

	circuitInputs := circuit.CircuitInputs{
		AvgX_3:      avg_x_3,
		AvgY_3:      avg_y_3,
		AvgX_Sq_3:   avg_x_sq_3,
		AvgY_Sq_3:   avg_y_sq_3,
		StdevX_Sq_3: stdev_x_sq_3,
		StdevY_Sq_3: stdev_y_sq_3,
		Cov_3:       cov_3,
		SSIM_3:      ssim_3,
		C1_3:        C1_3,
		C2_3:        C2_3,
		C1_6:        C1_6,
	}

	circuit.Verify(circuitInputs)

	lhs := (avg_x_3*avg_y_3*2 + C1_6) * (cov_3*2 + C2_3)
	rhs := ssim_3 * (avg_x_sq_3 + avg_y_sq_3 + C1_3) * (stdev_x_sq_3 + stdev_y_sq_3 + C2_3)

	fmt.Println("lhs & rhs", lhs, rhs)
}
