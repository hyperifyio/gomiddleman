// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package httputils

import (
	"io"
	"net/http"
	"testing"
)

// MakeRequestAndVerifyResponse makes an HTTP request to the specified target address using the provided client
// and verifies if the response body matches the expected content.
func MakeRequestAndVerifyResponse(t *testing.T, client *http.Client, targetAddr, expected string) {
	resp, err := client.Get(targetAddr)
	if err != nil {
		t.Fatalf("[MakeRequestAndVerifyResponse]: Failed to make request through proxyutils: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("[MakeRequestAndVerifyResponse]: Failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("[MakeRequestAndVerifyResponse]: Failed to read response body: %v", err)
	}

	if string(body) != expected {
		t.Errorf("[MakeRequestAndVerifyResponse]: Expected response body %q, got %q", expected, body)
	}
}
