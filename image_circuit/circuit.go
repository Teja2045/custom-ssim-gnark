package custom_ssim

import (
	"math"

	"github.com/consensys/gnark/frontend"
)

// SSIMCircuit defines the SSIM circuit
type SSIMCircuit struct {
	Image1 []frontend.Variable // Pixel values for image 1
	Image2 []frontend.Variable // Pixel values for image 2

	ExpectedSSIM frontend.Variable // Expected SSIM value for asserting
}

// Constants for SSIM formula
var (
	L  = 255.0
	K1 = 0.01
	K2 = 0.03
	C1 = math.Pow((K1 * L), 2.0)
	C2 = math.Pow((K2 * L), 2.0)
)

// Define defines the circuit's constraints
func (c *SSIMCircuit) Define(api frontend.API) error {
	n := len(c.Image1) // Number of pixels

	// Step 1: Calculate Mean for both images
	sumImage1 := frontend.Variable(0)
	sumImage2 := frontend.Variable(0)
	for i := 0; i < n; i++ {
		sumImage1 = api.Add(sumImage1, c.Image1[i])
		sumImage2 = api.Add(sumImage2, c.Image2[i])
	}

	// Calculate the means
	Mean1 := api.Div(sumImage1, frontend.Variable(n))
	Mean2 := api.Div(sumImage2, frontend.Variable(n))

	// Step 2: Calculate Variance (Stdev^2) for both images
	variance1 := frontend.Variable(0)
	variance2 := frontend.Variable(0)
	for i := 0; i < n; i++ {
		diff1 := api.Sub(c.Image1[i], Mean1)
		diff2 := api.Sub(c.Image2[i], Mean2)
		variance1 = api.Add(variance1, api.Mul(diff1, diff1))
		variance2 = api.Add(variance2, api.Mul(diff2, diff2))
	}

	// Calculate the standard deviations
	// c.Stdev1 = api.Sqrt(api.Div(variance1, frontend.Variable(n)))
	// c.Stdev2 = api.Sqrt(api.Div(variance2, frontend.Variable(n)))

	// we will be using standanrd deviation square, so we won't need above sqaure root
	Stdev1 := api.Div(variance1, frontend.Variable(n))
	Stdev2 := api.Div(variance2, frontend.Variable(n))

	// Step 3: Calculate Covariance between the two images
	covariance := frontend.Variable(0)
	for i := 0; i < n; i++ {
		diff1 := api.Sub(c.Image1[i], Mean1)
		diff2 := api.Sub(c.Image2[i], Mean2)
		covariance = api.Add(covariance, api.Mul(diff1, diff2))
	}
	Covariance := api.Div(covariance, frontend.Variable(n))

	// Step 4: Calculate SSIM using the formula
	numerator := api.Mul(
		api.Add(api.Mul(Mean1, Mean2), frontend.Variable(C1)),
		api.Add(api.Mul(frontend.Variable(2), Covariance), frontend.Variable(C2)),
	)

	denominator := api.Mul(
		api.Add(api.Mul(Mean1, Mean1), api.Mul(Mean2, Mean2), frontend.Variable(C1)),
		api.Add(Stdev1, Stdev2, frontend.Variable(C2)),
	)

	SSIMResult := api.Div(numerator, denominator)

	// Assert that the computed SSIMResult is equal to the expected value
	api.AssertIsEqual(SSIMResult, c.ExpectedSSIM)

	return nil
}
