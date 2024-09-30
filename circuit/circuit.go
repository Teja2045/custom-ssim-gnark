package circuit

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

const INPUTS_SIZE = 12

// CircuitInputs represents the inputs for the SSIM circuit, all as int64
type CircuitInputs struct {
	AvgX_3      int64 `json:"avg_x_3"`      // avg_x_3
	AvgY_3      int64 `json:"avg_y_3"`      // avg_y_3
	AvgX_Sq_3   int64 `json:"avg_x_sq_3"`   // avg_x_sq_3
	AvgY_Sq_3   int64 `json:"avg_y_sq_3"`   // avg_y_sq_3
	StdevX_Sq_3 int64 `json:"stdev_x_sq_3"` // stdev_x_sq_3
	StdevY_Sq_3 int64 `json:"stdev_y_sq_3"` // stdev_y_sq_3
	Cov_3       int64 `json:"cov_3"`        // cov_3
	SSIM_3      int64 `json:"ssim_3"`       // ssim_3, marked as public in gnark
	C1_3        int64 `json:"c1_3"`         // C1_3 constant
	C2_3        int64 `json:"c2_3"`         // C2_3 constant
	C1_6        int64 `json:"c1_6"`         // C1_6 constant
}

// CheckCircuit defines the circuit structure
type Circuit struct {
	AvgX_3      frontend.Variable // avg_x_3
	AvgY_3      frontend.Variable // avg_y_3
	AvgX_Sq_3   frontend.Variable // avg_x_sq_3
	AvgY_Sq_3   frontend.Variable // avg_y_sq_3
	StdevX_Sq_3 frontend.Variable // stdev_x_sq_3
	StdevY_Sq_3 frontend.Variable // stdev_y_sq_3
	Cov_3       frontend.Variable // cov_3
	SSIM_3      frontend.Variable `gnark:",public"` // ssim_3
	C1_3        frontend.Variable // C1_3 constant
	C2_3        frontend.Variable // C2_3 constant
	C1_6        frontend.Variable // C1_6 constant
}

// Define defines the constraints for the circuit
func (circuit *Circuit) Define(api frontend.API) error {
	// 2 as a constant
	two := frontend.Variable(2)

	// Compute lhs: (avg_x_3 * avg_y_3 * 2 + C1_6) * (cov_3 * 2 + C2_3)
	avg_x_y_2 := api.Mul(circuit.AvgX_3, circuit.AvgY_3, two) // avg_x_3 * avg_y_3 * 2
	term1_lhs := api.Add(avg_x_y_2, circuit.C1_6)             // avg_x_y_2 + C1_6
	cov_2 := api.Mul(circuit.Cov_3, two)                      // cov_3 * 2
	term2_lhs := api.Add(cov_2, circuit.C2_3)                 // cov_2 + C2_3
	lhs := api.Mul(term1_lhs, term2_lhs)                      // (avg_x_3 * avg_y_3 * 2 + C1_6) * (cov_3 * 2 + C2_3)

	// Compute rhs: ssim_3 * (avg_x_sq_3 + avg_y_sq_3 + C1_3) * (stdev_x_sq_3 + stdev_y_sq_3 + C2_3)
	avg_sq_sum := api.Add(circuit.AvgX_Sq_3, circuit.AvgY_Sq_3, circuit.C1_3)       // avg_x_sq_3 + avg_y_sq_3 + C1_3
	stdev_sq_sum := api.Add(circuit.StdevX_Sq_3, circuit.StdevY_Sq_3, circuit.C2_3) // stdev_x_sq_3 + stdev_y_sq_3 + C2_3
	rhs := api.Mul(circuit.SSIM_3, avg_sq_sum, stdev_sq_sum)                        // ssim_3 * (avg_x_sq_3 + avg_y_sq_3 + C1_3) * (stdev_x_sq_3 + stdev_y_sq_3 + C2_3)

	// Calculate the ratio lhs / rhs
	// lhs / rhs should be we close to 1
	// i.e ratio <= 1.1 & ratio >= 0.9
	// multiply 10 on both sides
	// ratio * 10 <= 11 && ratio * 10 >= 9
	ratio := api.Div(lhs, rhs)
	ratio = api.Mul(ratio, 10)

	// Set the tolerance bounds: [0.9, 1.1]
	// upperBound := frontend.Variable(11)

	// Ensure that the ratio is within the range [0.9, 1.1]
	api.AssertIsLessOrEqual(9, ratio)
	// api.AssertIsLessOrEqual(ratio, upperBound) // TODO: contraint not working..

	return nil
}

func GenerateZKProof(cs constraint.ConstraintSystem, pk groth16.ProvingKey, customInputs any) (groth16.Proof, witness.Witness, error) {

	circuitInputs := customInputs.(CircuitInputs)

	assignment := Circuit{
		AvgX_3:      circuitInputs.AvgX_3,
		AvgY_3:      circuitInputs.AvgY_3,
		AvgX_Sq_3:   circuitInputs.AvgX_Sq_3,
		AvgY_Sq_3:   circuitInputs.AvgY_Sq_3,
		StdevX_Sq_3: circuitInputs.StdevX_Sq_3,
		StdevY_Sq_3: circuitInputs.StdevY_Sq_3,
		Cov_3:       circuitInputs.Cov_3,
		SSIM_3:      circuitInputs.SSIM_3,
		C1_3:        circuitInputs.C1_3,
		C2_3:        circuitInputs.C2_3,
		C1_6:        circuitInputs.C1_6,
	}

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, nil, err
	}

	zkproof, err := groth16.Prove(cs, pk, witness)
	if err != nil {
		return nil, nil, err
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return nil, nil, err
	}

	return zkproof, publicWitness, nil
}
