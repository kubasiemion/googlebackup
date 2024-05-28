package main

import (
	"encoding/base64"
	"testing"
)

func TestEncoding(t *testing.T) {
	sample := `CAAF3jzUGsvFz_3pj6ME16sq1Hi28b35PUHx9PLpSC__HWt-9+A`
	plainhex, err := base64.URLEncoding.DecodeString(sample)
	if err != nil {
		t.Errorf("Error decoding base64: %v", err)
	}
	t.Logf("Decoded: %x", plainhex)
}
