// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package tlsutils

import "testing"

func TestCreateCertificateRequest(t *testing.T) {
	// First, generate a private key to use for the CSR
	privateKey, err := GeneratePrivateKey(2048)
	if err != nil {
		t.Fatalf("Failed to generate private key for testing: %v", err)
	}

	// Now, create the CSR
	csr, err := CreateCertificateRequest(privateKey)
	if err != nil {
		t.Fatalf("CreateCertificateRequest failed: %v", err)
	}

	// Check that the CSR is not nil
	if csr == nil {
		t.Fatal("Expected non-nil CSR, got nil")
	}

	// Verify the subject details in the CSR
	expectedCN := "Client CN"
	if csr.Subject.CommonName != expectedCN {
		t.Errorf("Expected CommonName %s, got %s", expectedCN, csr.Subject.CommonName)
	}

	expectedOrg := "Client Org"
	if len(csr.Subject.Organization) == 0 || csr.Subject.Organization[0] != expectedOrg {
		t.Errorf("Expected Organization %s, got %v", expectedOrg, csr.Subject.Organization)
	}

	// Optionally, check the CSR's signature to ensure it's valid
	if err := csr.CheckSignature(); err != nil {
		t.Errorf("CSR signature check failed: %v", err)
	}
}
