package oss

import "testing"

// TestUnitNormalizeContentMD5 locks the normalization and validation of the
// OSS PutObject Content-MD5 header. It is a pure-function unit test (no ACC
// env required).
func TestUnitNormalizeContentMD5(t *testing.T) {
	cases := []struct {
		input   string
		want    string
		wantErr bool
	}{
		// A 32-character hex MD5 digest (as produced by Terraform's md5() and
		// filemd5() functions) is converted to base64.
		{"45d8d11bed20b619add195b98d4b1775", "RdjRG+0gthmt0ZW5jUsXdQ==", false},
		// An already base64-encoded MD5 digest is passed through unchanged.
		{"RdjRG+0gthmt0ZW5jUsXdQ==", "RdjRG+0gthmt0ZW5jUsXdQ==", false},
		{"ewBv9NcPaMxlBhrPL4Aubw==", "ewBv9NcPaMxlBhrPL4Aubw==", false},
		// A short hex string is valid base64 but decodes to fewer than 16 bytes,
		// so it is not a valid MD5 digest; the provider rejects it instead of
		// forwarding it to OSS, which would otherwise surface an opaque error.
		{"deadbeef", "", true},
		// Input that is neither a valid 32-character hex MD5 digest nor valid
		// base64 is rejected by the provider instead of being forwarded.
		{"not-hex-at-all==", "", true},
		{"plain garbage value", "", true},
	}
	for _, c := range cases {
		got, err := normalizeContentMD5(c.input)
		if c.wantErr {
			if err == nil {
				t.Errorf("normalizeContentMD5(%q) expected an error, got %q with nil err", c.input, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("normalizeContentMD5(%q) unexpected error: %v", c.input, err)
			continue
		}
		if got != c.want {
			t.Errorf("normalizeContentMD5(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
