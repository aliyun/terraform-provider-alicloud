package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/packet"
)

// TestFetchKeybasePubkeys is an integration test that calls the real Keybase API.
// It is skipped when the KEYBASE_INTEGRATION environment variable is unset to
// avoid flakiness in offline / CI environments.
func TestUnitCommonFetchKeybasePubkeys(t *testing.T) {
	if os.Getenv("KEYBASE_INTEGRATION") == "" {
		t.Skip("skipping Keybase integration test: set KEYBASE_INTEGRATION=1 to run")
	}
	testset := []string{"keybase:jefferai", "keybase:hashicorp"}
	ret, err := FetchKeybasePubkeys(testset)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	fingerprints := []string{}
	for _, user := range testset {
		data, err := base64.StdEncoding.DecodeString(ret[user])
		if err != nil {
			t.Fatalf("error decoding key for user %s: %v", user, err)
		}
		entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewBuffer(data)))
		if err != nil {
			t.Fatalf("error parsing key for user %s: %v", user, err)
		}
		fingerprints = append(fingerprints, hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
	}

	exp := []string{
		"0f801f518ec853daff611e836528efcac6caa3db",
		"c874011f0ab405110d02105534365d9472d7468f",
	}

	if !reflect.DeepEqual(fingerprints, exp) {
		t.Fatalf("fingerprints do not match; expected \n%#v\ngot\n%#v\n", exp, fingerprints)
	}
}

// ---------------------------------------------------------------------------
// FetchKeybasePubkeys edge-case tests (no network required)
// ---------------------------------------------------------------------------

// TestFetchKeybasePubkeys_NilInput verifies that a nil slice returns nil, nil.
func TestUnitCommonFetchKeybasePubkeys_NilInput(t *testing.T) {
	got, err := FetchKeybasePubkeys(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil map, got %v", got)
	}
}

// TestFetchKeybasePubkeys_EmptyInput verifies that an empty slice returns nil, nil.
func TestUnitCommonFetchKeybasePubkeys_EmptyInput(t *testing.T) {
	got, err := FetchKeybasePubkeys([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil map, got %v", got)
	}
}

// TestFetchKeybasePubkeys_NoKeybasePrefix verifies that input entries without
// the "keybase:" prefix are silently ignored and nil is returned.
func TestUnitCommonFetchKeybasePubkeys_NoKeybasePrefix(t *testing.T) {
	got, err := FetchKeybasePubkeys([]string{"justausername", "another"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil map for non-keybase inputs, got %v", got)
	}
}

// TestFetchKeybasePubkeys_MockServer tests the full HTTP + JSON parsing path
// against a local mock server, avoiding any real network calls.
func TestUnitCommonFetchKeybasePubkeys_MockServer(t *testing.T) {
	// Build a minimal Keybase-shaped JSON response with an empty bundle to
	// exercise the "missing key" error path without hitting the real API.
	mockResp := `{"status":{"name":"OK"},"them":[{"public_keys":{"primary":{"bundle":""}}}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, mockResp)
	}))
	defer srv.Close()

	// FetchKeybasePubkeys hardcodes the Keybase URL, so we cannot redirect it
	// without modifying the implementation.  Instead, exercise DecodeJSONFromReader
	// directly to confirm the JSON parsing logic used inside FetchKeybasePubkeys.
	type kbStatus struct {
		Name string
	}
	type kbPrimary struct {
		Bundle string
	}
	type kbPubKeys struct {
		Primary kbPrimary `json:"primary"`
	}
	type kbThem struct {
		kbPubKeys `json:"public_keys"`
	}
	type kbResp struct {
		Status kbStatus
		Them   []kbThem
	}

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatalf("mock GET failed: %v", err)
	}
	defer resp.Body.Close()

	out := &kbResp{}
	if err := DecodeJSONFromReader(resp.Body, out); err != nil {
		t.Fatalf("DecodeJSONFromReader: %v", err)
	}
	if out.Status.Name != "OK" {
		t.Errorf("status.name = %q; want OK", out.Status.Name)
	}
	if len(out.Them) != 1 {
		t.Fatalf("expected 1 'them' entry, got %d", len(out.Them))
	}
	if out.Them[0].Primary.Bundle != "" {
		t.Errorf("expected empty bundle, got %q", out.Them[0].Primary.Bundle)
	}
}

// ---------------------------------------------------------------------------
// DecodeJSONFromReader tests
// ---------------------------------------------------------------------------

func TestUnitCommonDecodeJSONFromReader_NilReader(t *testing.T) {
	var out interface{}
	err := DecodeJSONFromReader(nil, &out)
	if err == nil {
		t.Fatal("expected error for nil reader, got nil")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil, got: %v", err)
	}
}

func TestUnitCommonDecodeJSONFromReader_NilOutput(t *testing.T) {
	r := strings.NewReader(`{"key":"value"}`)
	err := DecodeJSONFromReader(r, nil)
	if err == nil {
		t.Fatal("expected error for nil output, got nil")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Errorf("error should mention nil, got: %v", err)
	}
}

func TestUnitCommonDecodeJSONFromReader_ValidJSON(t *testing.T) {
	r := strings.NewReader(`{"name":"alicloud","count":42}`)
	var out map[string]interface{}
	if err := DecodeJSONFromReader(r, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["name"] != "alicloud" {
		t.Errorf("name = %v; want alicloud", out["name"])
	}
	// UseNumber means integers come back as json.Number, not float64
	if fmt.Sprint(out["count"]) != "42" {
		t.Errorf("count = %v; want 42", out["count"])
	}
}

func TestUnitCommonDecodeJSONFromReader_InvalidJSON(t *testing.T) {
	r := strings.NewReader(`{not valid json}`)
	var out interface{}
	if err := DecodeJSONFromReader(r, &out); err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestUnitCommonDecodeJSONFromReader_EmptyJSON(t *testing.T) {
	r := strings.NewReader(`{}`)
	var out map[string]interface{}
	if err := DecodeJSONFromReader(r, &out); err != nil {
		t.Fatalf("unexpected error for empty JSON object: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestUnitCommonDecodeJSONFromReader_JSONArray(t *testing.T) {
	r := strings.NewReader(`[1,2,3]`)
	var out []interface{}
	if err := DecodeJSONFromReader(r, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 elements, got %d", len(out))
	}
}

func TestUnitCommonDecodeJSONFromReader_NestedStruct(t *testing.T) {
	type Inner struct {
		Value string
	}
	type Outer struct {
		Inner Inner `json:"inner"`
		Count int   `json:"count"`
	}
	r := strings.NewReader(`{"inner":{"Value":"test"},"count":7}`)
	var out Outer
	if err := DecodeJSONFromReader(r, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Inner.Value != "test" {
		t.Errorf("inner.value = %q; want test", out.Inner.Value)
	}
	if out.Count != 7 {
		t.Errorf("count = %d; want 7", out.Count)
	}
}
