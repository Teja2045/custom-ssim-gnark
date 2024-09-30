package custom_ssim

// import (
// 	"fmt"
// 	"image"

// 	"gocv.io/x/gocv"
// )

// // Function to calculate mean (μ) of a pixel window
// func mean1(pixels [][]float64) float64 {
// 	var sum float64
// 	count := len(pixels) * len(pixels[0])
// 	for i := 0; i < len(pixels); i++ {
// 		for j := 0; j < len(pixels[i]); j++ {
// 			sum += pixels[i][j]
// 		}
// 	}
// 	return sum / float64(count)
// }

// // Function to calculate variance (σ²) of a pixel window
// func variance1(pixels [][]float64, meanValue float64) float64 {
// 	var sum float64
// 	count := len(pixels) * len(pixels[0])
// 	for i := 0; i < len(pixels); i++ {
// 		for j := 0; j < len(pixels[i]); j++ {
// 			sum += (pixels[i][j] - meanValue) * (pixels[i][j] - meanValue)
// 		}
// 	}
// 	return sum / float64(count)
// }

// // Function to calculate covariance (σxy) between two windows of pixels
// func covariance1(pixels1, pixels2 [][]float64, mean1, mean2 float64) float64 {
// 	var sum float64
// 	count := len(pixels1) * len(pixels1[0])
// 	for i := 0; i < len(pixels1); i++ {
// 		for j := 0; j < len(pixels1[i]); j++ {
// 			sum += (pixels1[i][j] - mean1) * (pixels2[i][j] - mean2)
// 		}
// 	}
// 	return sum / float64(count)
// }

// // Extract a window of size WxW from an image, starting at (x, y)
// func extractWindow1(img gocv.Mat, x, y, windowSize int) [][]float64 {
// 	window := make([][]float64, windowSize)
// 	for i := 0; i < windowSize; i++ {
// 		window[i] = make([]float64, windowSize)
// 		for j := 0; j < windowSize; j++ {
// 			if x+i < img.Rows() && y+j < img.Cols() {
// 				window[i][j] = float64(img.GetUCharAt(x+i, y+j)) / 255.0
// 			}
// 		}
// 	}
// 	return window
// }

// // Calculate SSIM between two images
// func calculateSSIM1(image1, image2 gocv.Mat, windowSize int) float64 {
// 	C1 := 0.01 * 0.01
// 	C2 := 0.03 * 0.03

// 	var ssimSum float64
// 	var windowCount int

// 	for i := 0; i < image1.Rows(); i += windowSize {
// 		for j := 0; j < image1.Cols(); j += windowSize {
// 			// Extract windows from both images
// 			window1 := extractWindow1(image1, i, j, windowSize)
// 			window2 := extractWindow1(image2, i, j, windowSize)

// 			// Calculate mean (μ) of the windows
// 			mean_1 := mean1(window1)
// 			mean_2 := mean1(window2)

// 			// Calculate variance (σ²) of the windows
// 			var1 := variance1(window1, mean_1)
// 			var2 := variance1(window2, mean_2)

// 			// Calculate covariance (σxy) between the windows
// 			cov := covariance1(window1, window2, mean_1, mean_2)

// 			// Calculate SSIM for this window
// 			numerator := (2*mean_1*mean_2 + C1) * (2*cov + C2)
// 			denominator := (mean_1*mean_1 + mean_2*mean_2 + C1) * (var1 + var2 + C2)

// 			ssim := numerator / denominator
// 			ssimSum += ssim
// 			windowCount++
// 		}
// 	}

// 	return ssimSum / float64(windowCount)
// }

// // LoadImage loads an image from file and returns the gocv.Mat representation
// func LoadImage(filePath string) (gocv.Mat, error) {
// 	img := gocv.IMRead(filePath, gocv.IMReadColor)
// 	if img.Empty() {
// 		return img, fmt.Errorf("error loading image: %s", filePath)
// 	}
// 	return img, nil
// }

// func Compare(image1Path, image2Path string) {

// 	img1, err := LoadImage(image1Path)
// 	if err != nil {
// 		fmt.Println("Error loading image1:", err)
// 		return
// 	}
// 	img2, err := LoadImage(image2Path)
// 	if err != nil {
// 		fmt.Println("Error loading image2:", err)
// 		return
// 	}

// 	// Convert both images to grayscale
// 	gocv.CvtColor(img1, &img1, gocv.ColorBGRToGray)
// 	gocv.CvtColor(img2, &img2, gocv.ColorBGRToGray)

// 	// Resize images if necessary (should have same dimensions)
// 	targetSize := image.Point{100, 100}
// 	gocv.Resize(img1, &img1, targetSize, 0, 0, gocv.InterpolationDefault)
// 	gocv.Resize(img2, &img2, targetSize, 0, 0, gocv.InterpolationDefault)

// 	// Calculate SSIM between the two images
// 	ssim := calculateSSIM1(img1, img2, 11)
// 	fmt.Printf("SSIM: %.4f\n", ssim)
// }
