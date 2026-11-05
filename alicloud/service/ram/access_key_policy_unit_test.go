package ram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// These TestUnit cases exercise the unexported access-key-policy helpers that
// live in package ram. They are kept in an internal test file (package ram)
// rather than the external ram_test package so the lower-case helpers remain
// reachable; the TestAcc cases that drive the managed resource stay in
// resource_access_key_policy_test.go.

func TestUnitRamAccessKeyPolicyStripVersion(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"removes version", `{"Version":1,"Status":"Active","Statements":[]}`, `{"Statements":[],"Status":"Active"}`},
		{"no version keeps document", `{"Status":"Active","Statements":[]}`, `{"Statements":[],"Status":"Active"}`},
		{"empty object", `{}`, `{}`},
		{"invalid json returned as-is", `not-json`, `not-json`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, stripAccessKeyPolicyVersion(c.in))
		})
	}
}

func TestUnitRamAccessKeyPolicyEquivalent(t *testing.T) {
	// Same policy, one with the server-managed Version field and one without.
	a := `{"Version":1,"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`
	b := `{"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`
	assert.True(t, accessKeyPolicyEquivalent(a, b), "documents differing only by Version should be equivalent")

	// Different IP list -> not equivalent.
	c := `{"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.2/32"]}]}`
	assert.False(t, accessKeyPolicyEquivalent(a, c), "documents with different statements should not be equivalent")
}

func TestUnitRamAccessKeyPolicyIsEmpty(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"blank string", "   ", true},
		{"empty object", `{}`, true},
		{"disabled without statements", `{"Version":1,"Status":"Inactive","Statements":[]}`, true},
		{"disabled without statements field", `{"Status":"Inactive"}`, true},
		{"disabled with statements is not empty", `{"Status":"Inactive","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`, false},
		{"active policy is not empty", `{"Status":"Active","Statements":[{"Type":"ClassicWhiteList","IPList":["10.0.0.1/32"]}]}`, false},
		{"invalid json is not empty", `not-json`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, isEmptyAccessKeyPolicy(c.in))
		})
	}
}

func TestUnitRamAccessKeyPolicyParseId(t *testing.T) {
	cases := []struct {
		name    string
		id      string
		wantUpn string
		wantAk  string
	}{
		{"composite id", "alice@1234567890.onaliyun.com:LTAI0123456789", "alice@1234567890.onaliyun.com", "LTAI0123456789"},
		{"access key only", "LTAI0123456789", "", "LTAI0123456789"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			upn, ak := parseRamAccessKeyPolicyId(c.id)
			assert.Equal(t, c.wantUpn, upn)
			assert.Equal(t, c.wantAk, ak)
		})
	}
}
