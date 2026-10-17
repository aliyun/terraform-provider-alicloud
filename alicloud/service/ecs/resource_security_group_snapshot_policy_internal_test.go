package ecs

import (
	"encoding/json"
	"testing"
	"time"
)

// F1: extractSecurityGroupSnapshotPolicyId must not coerce a missing or nil id
// into the pseudo-string "<nil>" the way fmt.Sprint(nil) does. A create response
// that carries no usable id fails the create instead of chasing a phantom id.
func TestUnitExtractSecurityGroupSnapshotPolicyId_RejectsNilAndMissing(t *testing.T) {
	cases := []struct {
		name     string
		response map[string]interface{}
		wantId   string
		wantOk   bool
	}{
		{
			name:     "string id present",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-abc123"},
			wantId:   "sp-abc123",
			wantOk:   true,
		},
		{
			name:     "numeric id encoded as float64",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": float64(42)},
			wantId:   "42",
			wantOk:   true,
		},
		// D3: RpcPost delivers JSON numbers as json.Number (tea-utils ReadAsJSON
		// uses json.Decoder.UseNumber()), not float64. extract must normalize a
		// json.Number via String() the same way it normalizes a float64, so a
		// numeric id returned by CreateSecurityGroupSnapshotPolicy is not
		// rejected as missing.
		{
			name:     "numeric id as json.Number (RpcPost path)",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": json.Number("42")},
			wantId:   "42",
			wantOk:   true,
		},
		{
			name:     "missing key is absent, not <nil>",
			response: map[string]interface{}{"Other": 1},
			wantId:   "",
			wantOk:   false,
		},
		{
			name:     "nil value is absent, not <nil>",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": nil},
			wantId:   "",
			wantOk:   false,
		},
		{
			name:     "empty string is absent",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": ""},
			wantId:   "",
			wantOk:   false,
		},
		{
			name:     "bare response is absent",
			response: map[string]interface{}{},
			wantId:   "",
			wantOk:   false,
		},
		// D1: a wrong-type (map/slice/bool) id must be rejected (ok=false) by the
		// shared securityGroupSnapshotPolicyIdOf default branch, not coerced into
		// a fmt.Sprintf pseudo-id such as "map[x:sp-1]" / "[sp-1]" / "true". This
		// is the unit-level guard for the Create asymmetry fixed in D1.
		{
			name:     "map id rejected, not coerced",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": map[string]interface{}{"x": "sp-1"}},
			wantId:   "",
			wantOk:   false,
		},
		{
			name:     "slice id rejected, not coerced",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": []interface{}{"sp-1"}},
			wantId:   "",
			wantOk:   false,
		},
		{
			name:     "bool id rejected, not coerced",
			response: map[string]interface{}{"SecurityGroupSnapshotPolicyId": true},
			wantId:   "",
			wantOk:   false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotId, gotOk := extractSecurityGroupSnapshotPolicyId(c.response)
			if gotId != c.wantId || gotOk != c.wantOk {
				t.Fatalf("extractSecurityGroupSnapshotPolicyId(%v) = (%q, %v); want (%q, %v)", c.response, gotId, gotOk, c.wantId, c.wantOk)
			}
			if gotOk && gotId == "<nil>" {
				t.Fatalf("regression: id coerced to <nil> for %v", c.response)
			}
		})
	}
}

// F2: policyIdOfOk must not coerce a missing, nil or wrong-type id into a
// pseudo-string (the fmt.Sprintf(map) path that produced "map[x:sp-1]" and let a
// malformed element slip past as a non-match, eventually dropping the resource).
// Such an element returns ("", false) so the caller treats it as malformed
// rather than silently skipping it.
func TestUnitPolicyIdOfOk_RejectsNilMissingAndWrongType(t *testing.T) {
	cases := []struct {
		name   string
		item   map[string]interface{}
		wantId string
		wantOk bool
	}{
		{"string id", map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"}, "sp-1", true},
		{"numeric id", map[string]interface{}{"SecurityGroupSnapshotPolicyId": float64(7)}, "7", true},
		{"numeric id as json.Number", map[string]interface{}{"SecurityGroupSnapshotPolicyId": json.Number("7")}, "7", true},
		{"missing key", map[string]interface{}{"Status": "enable"}, "", false},
		{"nil value", map[string]interface{}{"SecurityGroupSnapshotPolicyId": nil}, "", false},
		{"empty string", map[string]interface{}{"SecurityGroupSnapshotPolicyId": ""}, "", false},
		{"wrong-type object id", map[string]interface{}{"SecurityGroupSnapshotPolicyId": map[string]interface{}{"x": "sp-1"}}, "", false},
		{"wrong-type slice id", map[string]interface{}{"SecurityGroupSnapshotPolicyId": []interface{}{"sp-1"}}, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotId, gotOk := policyIdOfOk(c.item)
			if gotId != c.wantId || gotOk != c.wantOk {
				t.Fatalf("policyIdOfOk(%v) = (%q, %v); want (%q, %v)", c.item, gotId, gotOk, c.wantId, c.wantOk)
			}
			if gotOk && (gotId == "" || gotId == "<nil>") {
				t.Fatalf("regression: id coerced to %q for %v", gotId, c.item)
			}
		})
	}
}

// F2: securityGroupSnapshotPolicyItems distinguishes a well-formed list (which
// may be empty — a confirmed NotFound) from a malformed response (a missing
// root or an unrecognized type). Malformed returns ok=false so describe reports
// an error and Read/Update keep the state instead of dropping the resource.
func TestUnitSecurityGroupSnapshotPolicyItems_MalformedVsEmpty(t *testing.T) {
	cases := []struct {
		name      string
		response  map[string]interface{}
		wantItems int
		wantOk    bool
	}{
		{
			name:      "top-level array well-formed",
			response:  map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"}}},
			wantItems: 1,
			wantOk:    true,
		},
		{
			name:      "top-level array empty is well-formed (NotFound)",
			response:  map[string]interface{}{"SecurityGroupSnapshotPolicySets": []interface{}{}},
			wantItems: 0,
			wantOk:    true,
		},
		{
			name:      "inner object array well-formed",
			response:  map[string]interface{}{"SecurityGroupSnapshotPolicySets": map[string]interface{}{"SecurityGroupSnapshotPolicy": []interface{}{map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-2"}}}},
			wantItems: 1,
			wantOk:    true,
		},
		{
			name:      "inner object without array is malformed",
			response:  map[string]interface{}{"SecurityGroupSnapshotPolicySets": map[string]interface{}{"Other": 1}},
			wantItems: 0,
			wantOk:    false,
		},
		{
			name:      "root missing entirely is malformed",
			response:  map[string]interface{}{"Other": 1},
			wantItems: 0,
			wantOk:    false,
		},
		{
			name:      "root wrong scalar type is malformed",
			response:  map[string]interface{}{"SecurityGroupSnapshotPolicySets": "unexpected"},
			wantItems: 0,
			wantOk:    false,
		},
		{
			name:      "bare response is malformed",
			response:  map[string]interface{}{},
			wantItems: 0,
			wantOk:    false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotItems, gotOk := securityGroupSnapshotPolicyItems(c.response)
			if gotOk != c.wantOk {
				t.Fatalf("securityGroupSnapshotPolicyItems(%v) ok = %v; want %v", c.response, gotOk, c.wantOk)
			}
			if gotOk && len(gotItems) != c.wantItems {
				t.Fatalf("securityGroupSnapshotPolicyItems(%v) items = %d; want %d", c.response, len(gotItems), c.wantItems)
			}
		})
	}
}

// F1: describeSecurityGroupSnapshotPolicyVisible retries across the
// eventual-consistency window only for an empty-but-well-formed list. A real
// error or a visible object returns immediately; an exhausted window returns
// (nil, nil) so Create reports not-yet-visible instead of an indefinite hang.
func TestUnitDescribeSecurityGroupSnapshotPolicyVisible_RetrySemantics(t *testing.T) {
	object := map[string]interface{}{"SecurityGroupSnapshotPolicyId": "sp-1"}
	wantErr := errBoom{}

	t.Run("returns object once visible after empty attempts", func(t *testing.T) {
		calls := 0
		describe := func() (map[string]interface{}, error) {
			calls++
			if calls < 2 {
				return nil, nil // not yet visible
			}
			return object, nil
		}
		got, err := describeSecurityGroupSnapshotPolicyVisible(describe, 5, 0)
		if err != nil || got == nil {
			t.Fatalf("expected visible object after retry, got (%v, %v) after %d calls", got, err, calls)
		}
		if calls != 2 {
			t.Fatalf("expected 2 calls, got %d", calls)
		}
	})

	t.Run("returns nil after exhausting attempts on a still-empty list", func(t *testing.T) {
		calls := 0
		describe := func() (map[string]interface{}, error) {
			calls++
			return nil, nil
		}
		got, err := describeSecurityGroupSnapshotPolicyVisible(describe, 3, 0)
		if err != nil || got != nil {
			t.Fatalf("expected (nil, nil) after exhausting, got (%v, %v)", got, err)
		}
		if calls != 3 {
			t.Fatalf("expected 3 calls, got %d", calls)
		}
	})

	t.Run("real error returns immediately without retry", func(t *testing.T) {
		calls := 0
		describe := func() (map[string]interface{}, error) {
			calls++
			return nil, wantErr
		}
		got, err := describeSecurityGroupSnapshotPolicyVisible(describe, 5, 0)
		if err != wantErr || got != nil {
			t.Fatalf("expected immediate error, got (%v, %v)", got, err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d (real errors must not be retried)", calls)
		}
	})

	t.Run("visible on first attempt makes no retries", func(t *testing.T) {
		calls := 0
		describe := func() (map[string]interface{}, error) {
			calls++
			return object, nil
		}
		got, err := describeSecurityGroupSnapshotPolicyVisible(describe, 5, 0)
		if err != nil || got == nil {
			t.Fatalf("expected object, got (%v, %v)", got, err)
		}
		if calls != 1 {
			t.Fatalf("expected 1 call, got %d", calls)
		}
	})

	// Guards the sleep budget: a non-zero interval on the retry path must not
	// make the exhausted case longer than attempts*interval by sleeping past the
	// last attempt (the last attempt sleeps nothing).
	t.Run("does not sleep after the final attempt", func(t *testing.T) {
		start := time.Now()
		calls := 0
		describe := func() (map[string]interface{}, error) {
			calls++
			return nil, nil
		}
		_, _ = describeSecurityGroupSnapshotPolicyVisible(describe, 3, 5*time.Millisecond)
		elapsed := time.Since(start)
		if elapsed > 15*time.Millisecond {
			t.Fatalf("slept too long (%v); the final attempt must not sleep", elapsed)
		}
	})
}

type errBoom struct{}

func (errBoom) Error() string { return "boom" }
