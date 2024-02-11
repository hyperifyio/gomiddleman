// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import "testing"

func TestGeneratePrivateKey(t *testing.T) {
	// Define test cases with different bit sizes
	tests := []struct {
		name string
		bits int
	}{
		{"2048 bits", 2048},
		{"4096 bits", 4096},
		// You can add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			privateKey, err := GeneratePrivateKey(tc.bits)
			if err != nil {
				t.Fatalf("Failed to generate private key: %v", err)
			}

			// Verify the private key is not nil
			if privateKey == nil {
				t.Errorf("Generated private key is nil")
			}

			// Verify the private key has the correct bit length
			if privateKey.N.BitLen() != tc.bits {
				t.Errorf("Expected bit length %d, got %d", tc.bits, privateKey.N.BitLen())
			}
		})
	}
}
