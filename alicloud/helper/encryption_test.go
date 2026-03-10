package helper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/packet"
)

// testNewPGPKey generates a throw-away PGP entity and returns:
//   - b64Key: the base64-encoded binary representation of the entity
//   - entity: the full entity (public + private) for fingerprint comparison
//
// 1024-bit RSA is used so key generation stays fast in unit tests.
// SerializePrivate is used because in keybase/go-crypto the self-signatures on
// a freshly created entity are lazy and Serialize() (public-only) fails until
// they have been explicitly signed; SerializePrivate performs the signing as a
// side-effect of writing the secret key material.
func testNewPGPKey(t *testing.T) (b64Key string, entity *openpgp.Entity) {
	t.Helper()
	cfg := &packet.Config{RSABits: 1024}
	var err error
	entity, err = openpgp.NewEntity("Test User", "unit-test", "unit@test.local", cfg)
	if err != nil {
		t.Fatalf("openpgp.NewEntity: %v", err)
	}
	buf := new(bytes.Buffer)
	if err = entity.SerializePrivate(buf, nil); err != nil {
		t.Fatalf("entity.SerializePrivate: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), entity
}

// ---------------------------------------------------------------------------
// GetEntities
// ---------------------------------------------------------------------------

// TestGetEntities_Valid verifies that a well-formed base64-encoded PGP public
// key is decoded into exactly one entity with the expected fingerprint.
func TestUnitGetEntities_Valid(t *testing.T) {
	b64Key, want := testNewPGPKey(t)

	got, err := GetEntities([]string{b64Key})
	if err != nil {
		t.Fatalf("GetEntities returned unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 entity, got %d", len(got))
	}
	wantFP := fmt.Sprintf("%x", want.PrimaryKey.Fingerprint)
	gotFP := fmt.Sprintf("%x", got[0].PrimaryKey.Fingerprint)
	if wantFP != gotFP {
		t.Errorf("fingerprint mismatch: want %s, got %s", wantFP, gotFP)
	}
}

// TestGetEntities_MultipleKeys verifies that multiple keys are decoded in order.
func TestUnitGetEntities_MultipleKeys(t *testing.T) {
	k1, e1 := testNewPGPKey(t)
	k2, e2 := testNewPGPKey(t)

	got, err := GetEntities([]string{k1, k2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(got))
	}
	fp1 := fmt.Sprintf("%x", e1.PrimaryKey.Fingerprint)
	fp2 := fmt.Sprintf("%x", e2.PrimaryKey.Fingerprint)
	if fmt.Sprintf("%x", got[0].PrimaryKey.Fingerprint) != fp1 {
		t.Errorf("entity[0] fingerprint mismatch")
	}
	if fmt.Sprintf("%x", got[1].PrimaryKey.Fingerprint) != fp2 {
		t.Errorf("entity[1] fingerprint mismatch")
	}
}

// TestGetEntities_Empty verifies that an empty input returns an empty slice.
func TestUnitGetEntities_Empty(t *testing.T) {
	got, err := GetEntities([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 entities, got %d", len(got))
	}
}

// TestGetEntities_InvalidBase64 verifies that an invalid base64 string yields
// an error that mentions the decoding failure.
func TestUnitGetEntities_InvalidBase64(t *testing.T) {
	_, err := GetEntities([]string{"not-valid-base64!!!"})
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
	if !strings.Contains(err.Error(), "decoding") {
		t.Errorf("error should mention decoding, got: %v", err)
	}
}

// TestGetEntities_InvalidPGP verifies that valid base64 whose payload is not
// a PGP packet returns an error.
func TestUnitGetEntities_InvalidPGP(t *testing.T) {
	notPGP := base64.StdEncoding.EncodeToString([]byte("this is definitely not a pgp packet"))
	_, err := GetEntities([]string{notPGP})
	if err == nil {
		t.Fatal("expected error for non-PGP data, got nil")
	}
}

// ---------------------------------------------------------------------------
// GetFingerprints
// ---------------------------------------------------------------------------

// TestGetFingerprints_WithEntities checks that GetFingerprints returns the
// correct hex fingerprint when entities are supplied and pgpKeys is nil.
func TestUnitGetFingerprints_WithEntities(t *testing.T) {
	_, entity := testNewPGPKey(t)

	fps, err := GetFingerprints(nil, []*openpgp.Entity{entity})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fps) != 1 {
		t.Fatalf("expected 1 fingerprint, got %d", len(fps))
	}
	want := fmt.Sprintf("%x", entity.PrimaryKey.Fingerprint)
	if fps[0] != want {
		t.Errorf("fingerprint mismatch: want %s, got %s", want, fps[0])
	}
}

// TestGetFingerprints_WithPGPKeys checks that GetFingerprints derives entities
// from pgpKeys when entities is nil.
func TestUnitGetFingerprints_WithPGPKeys(t *testing.T) {
	b64Key, entity := testNewPGPKey(t)

	fps, err := GetFingerprints([]string{b64Key}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fps) != 1 {
		t.Fatalf("expected 1 fingerprint, got %d", len(fps))
	}
	want := fmt.Sprintf("%x", entity.PrimaryKey.Fingerprint)
	if fps[0] != want {
		t.Errorf("fingerprint mismatch: want %s, got %s", want, fps[0])
	}
}

// TestGetFingerprints_MultipleEntities checks that multiple entities are handled
// in the correct order.
func TestUnitGetFingerprints_MultipleEntities(t *testing.T) {
	_, e1 := testNewPGPKey(t)
	_, e2 := testNewPGPKey(t)

	fps, err := GetFingerprints(nil, []*openpgp.Entity{e1, e2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fps) != 2 {
		t.Fatalf("expected 2 fingerprints, got %d", len(fps))
	}
	if fps[0] != fmt.Sprintf("%x", e1.PrimaryKey.Fingerprint) {
		t.Errorf("fingerprint[0] mismatch")
	}
	if fps[1] != fmt.Sprintf("%x", e2.PrimaryKey.Fingerprint) {
		t.Errorf("fingerprint[1] mismatch")
	}
}

// TestGetFingerprints_InvalidKey checks that an invalid key string causes an error.
func TestUnitGetFingerprints_InvalidKey(t *testing.T) {
	_, err := GetFingerprints([]string{"not-base64!!!"}, nil)
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

// TestGetFingerprints_EmptyEntities checks the empty case for both parameters.
func TestUnitGetFingerprints_EmptyEntities(t *testing.T) {
	fps, err := GetFingerprints(nil, []*openpgp.Entity{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fps) != 0 {
		t.Errorf("expected 0 fingerprints, got %d", len(fps))
	}
}

// ---------------------------------------------------------------------------
// EncryptShares
// ---------------------------------------------------------------------------

// TestEncryptShares_LengthMismatch checks that mismatched slice lengths return
// an error that mentions "mismatch".
func TestUnitEncryptShares_LengthMismatch(t *testing.T) {
	b64Key, _ := testNewPGPKey(t)

	_, _, err := EncryptShares(
		[][]byte{[]byte("a"), []byte("b")},
		[]string{b64Key},
	)
	if err == nil {
		t.Fatal("expected error for length mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "mismatch") {
		t.Errorf("error should mention mismatch, got: %v", err)
	}
}

// TestEncryptShares_Valid checks that valid inputs produce a non-empty
// ciphertext and the correct fingerprint.
func TestUnitEncryptShares_Valid(t *testing.T) {
	b64Key, entity := testNewPGPKey(t)
	plaintext := []byte("hello, PGP world!")

	fps, cts, err := EncryptShares([][]byte{plaintext}, []string{b64Key})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fps) != 1 || len(cts) != 1 {
		t.Fatalf("expected 1 fingerprint and 1 ciphertext, got fps=%d cts=%d", len(fps), len(cts))
	}
	want := fmt.Sprintf("%x", entity.PrimaryKey.Fingerprint)
	if fps[0] != want {
		t.Errorf("fingerprint mismatch: want %s, got %s", want, fps[0])
	}
	if len(cts[0]) == 0 {
		t.Error("ciphertext should not be empty")
	}
	// Ciphertext must differ from plaintext.
	if bytes.Equal(cts[0], plaintext) {
		t.Error("ciphertext equals plaintext – encryption did not run")
	}
}

// TestEncryptShares_EmptyInput checks that empty slices return empty results
// without error.
func TestUnitEncryptShares_EmptyInput(t *testing.T) {
	fps, cts, err := EncryptShares([][]byte{}, []string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fps) != 0 || len(cts) != 0 {
		t.Errorf("expected empty results, got fps=%v cts=%v", fps, cts)
	}
}

// TestEncryptShares_InvalidKey checks that an invalid PGP key causes an error.
func TestUnitEncryptShares_InvalidKey(t *testing.T) {
	_, _, err := EncryptShares([][]byte{[]byte("data")}, []string{"not-base64!!!"})
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}

// TestEncryptShares_MultipleShares verifies that multiple (plaintext, key) pairs
// are all encrypted independently.
func TestUnitEncryptShares_MultipleShares(t *testing.T) {
	k1, _ := testNewPGPKey(t)
	k2, _ := testNewPGPKey(t)
	p1 := []byte("secret one")
	p2 := []byte("secret two")

	fps, cts, err := EncryptShares([][]byte{p1, p2}, []string{k1, k2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fps) != 2 || len(cts) != 2 {
		t.Fatalf("expected 2 results, got fps=%d cts=%d", len(fps), len(cts))
	}
	// The two ciphertexts should differ (different keys, different plaintexts).
	if bytes.Equal(cts[0], cts[1]) {
		t.Error("ciphertexts from different keys should not be equal")
	}
}

// ---------------------------------------------------------------------------
// EncryptValue
// ---------------------------------------------------------------------------

// TestEncryptValue_Valid checks that EncryptValue returns a non-empty fingerprint
// and a valid base64-encoded ciphertext.
func TestUnitEncryptValue_Valid(t *testing.T) {
	b64Key, _ := testNewPGPKey(t)

	fp, encrypted, err := EncryptValue(b64Key, "super-secret-password", "test password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp == "" {
		t.Error("expected non-empty fingerprint")
	}
	if encrypted == "" {
		t.Error("expected non-empty encrypted value")
	}
	// The returned encrypted value must be valid base64.
	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Errorf("encrypted value is not valid base64: %v", err)
	}
	if len(decoded) == 0 {
		t.Error("decoded ciphertext should be non-empty")
	}
}

// TestEncryptValue_EmptyValue checks that encrypting an empty value still
// succeeds (PGP can encrypt zero-length messages).
func TestUnitEncryptValue_EmptyValue(t *testing.T) {
	b64Key, _ := testNewPGPKey(t)

	fp, encrypted, err := EncryptValue(b64Key, "", "empty value test")
	if err != nil {
		t.Fatalf("unexpected error encrypting empty value: %v", err)
	}
	if fp == "" {
		t.Error("expected non-empty fingerprint even for empty value")
	}
	if encrypted == "" {
		t.Error("expected non-empty ciphertext even for empty value")
	}
}

// TestEncryptValue_InvalidKey checks that an invalid key causes an error whose
// message includes the description parameter.
func TestUnitEncryptValue_InvalidKey(t *testing.T) {
	_, _, err := EncryptValue("not-base64!!!", "secret", "test description")
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
	if !strings.Contains(err.Error(), "test description") {
		t.Errorf("error should contain the description 'test description', got: %v", err)
	}
}

// TestEncryptValue_FingerprintConsistency verifies that the fingerprint returned
// by EncryptValue matches the one returned by GetFingerprints for the same key.
func TestUnitEncryptValue_FingerprintConsistency(t *testing.T) {
	b64Key, _ := testNewPGPKey(t)

	fp, _, err := EncryptValue(b64Key, "data", "fp consistency test")
	if err != nil {
		t.Fatalf("EncryptValue error: %v", err)
	}

	fps, err := GetFingerprints([]string{b64Key}, nil)
	if err != nil {
		t.Fatalf("GetFingerprints error: %v", err)
	}
	if len(fps) != 1 {
		t.Fatalf("expected 1 fingerprint, got %d", len(fps))
	}
	if fp != fps[0] {
		t.Errorf("EncryptValue fingerprint %q != GetFingerprints fingerprint %q", fp, fps[0])
	}
}

// ---------------------------------------------------------------------------
// RetrieveGPGKey
// ---------------------------------------------------------------------------

// TestRetrieveGPGKey_Direct checks that a key without the "keybase:" prefix is
// returned unchanged.
func TestUnitRetrieveGPGKey_Direct(t *testing.T) {
	b64Key, _ := testNewPGPKey(t)

	result, err := RetrieveGPGKey(b64Key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != b64Key {
		t.Error("expected key to be returned unchanged for non-keybase input")
	}
}

// TestRetrieveGPGKey_Empty checks that an empty string is returned as-is.
func TestUnitRetrieveGPGKey_Empty(t *testing.T) {
	result, err := RetrieveGPGKey("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

// TestRetrieveGPGKey_ArbitraryString checks that any non-keybase string is
// returned verbatim without triggering a network call.
func TestUnitRetrieveGPGKey_ArbitraryString(t *testing.T) {
	for _, input := range []string{
		"just-a-string",
		"BEGIN PGP PUBLIC KEY BLOCK",
		"base64data==",
		"no:prefix:here",
	} {
		result, err := RetrieveGPGKey(input)
		if err != nil {
			t.Errorf("RetrieveGPGKey(%q) error: %v", input, err)
			continue
		}
		if result != input {
			t.Errorf("RetrieveGPGKey(%q) = %q; want unchanged", input, result)
		}
	}
}
