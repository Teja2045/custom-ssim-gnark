package circuit

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func Verify(circuitInputs CircuitInputs) error {
	var circuit Circuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return fmt.Errorf("error while compiling circuit: %v", err)
	}

	// Generate prover key and verifier key using groth16
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return fmt.Errorf("error during circuit setup: %v", err)
	}

	zkProof, publicWitness, err := GenerateZKProof(ccs, pk, circuitInputs)
	if err != nil {
		return fmt.Errorf("error while generating zk proof: %v", err)
	}
	return groth16.Verify(zkProof, vk, publicWitness)
}
