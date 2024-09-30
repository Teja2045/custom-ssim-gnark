package main

import custom_ssim "github.com/arka-labs/ssim/custom"

func main() {
	image1 := "./images/dolphin.jpg"
	image2 := "./images/img1.jpg"

	custom_ssim.Compare2(image1, image2)
}
