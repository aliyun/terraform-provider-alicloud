package convert

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"testing"
)

// Tests for the format helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go. The helpers return (value, error) instead of
// panicking, so these cases pin the error semantics.

func TestUnitCommonFormatInt(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    int
		wantErr     bool
		expectedMsg string
	}{
		{"nil", nil, 0, false, ""},
		{"float64", 3.14, 3, false, ""},
		{"float32", float32(3.0), 3, false, ""},
		{"int64", int64(123), 123, false, ""},
		{"int32", int32(456), 456, false, ""},
		{"int", 789, 789, false, ""},
		{"string_number", "123", 123, false, ""},
		{"string_negative", "-456", -456, false, ""},
		{"string_empty", "", 0, false, ""},
		{"json_number", json.Number("123"), 123, false, ""},
		{"string_invalid", "abc", 0, true, `Converting string value "abc" to int got an error`},
		{"json_number_invalid", json.Number("12.5"), 0, true, "Converting json.Number value 12.5 to int got an error"},
		{"unsupported_type", true, 0, true, "Not support type bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatInt(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error expectation failed: wantErr %v, got %v", tt.wantErr, err)
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}

			if tt.wantErr && err != nil && tt.expectedMsg != "" {
				if !strings.Contains(err.Error(), tt.expectedMsg) {
					t.Errorf("Error message mismatch. Expected to contain: %s, got: %s",
						tt.expectedMsg, err.Error())
				}
			}
		})
	}
}

func FuzzCommonFormatInt(f *testing.F) {
	f.Add("123")
	f.Add("")
	f.Add("abc")
	f.Add("2147483648")
	f.Add("-1")
	f.Add(" 7")
	f.Add("0x10")
	f.Fuzz(func(t *testing.T, s string) {
		v, err := FormatInt(s)
		if s == "" {
			if err != nil || v != 0 {
				t.Fatalf("FormatInt(\"\") = %d, %v; want 0, nil", v, err)
			}
			return
		}
		parsed, atoiErr := strconv.Atoi(s)
		if atoiErr != nil {
			if err == nil {
				t.Fatalf("FormatInt(%q) = %d, want error like strconv.Atoi: %v", s, v, atoiErr)
			}
			if v != 0 {
				t.Fatalf("FormatInt(%q) returned error %v but non-zero value %d", s, err, v)
			}
			return
		}
		if err != nil || v != parsed {
			t.Fatalf("FormatInt(%q) = %d, %v; want %d, nil", s, v, err, parsed)
		}
	})
}

func TestUnitCommonFormatBool(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    bool
		wantErr     bool
		expectedMsg string
	}{
		{"nil", nil, false, false, ""},
		{"bool_true", true, true, false, ""},
		{"bool_false", false, false, false, ""},
		{"string_true", "true", true, false, ""},
		{"string_false", "false", false, false, ""},
		{"string_one", "1", true, false, ""},
		{"string_zero", "0", false, false, ""},
		{"string_empty", "", false, false, ""},
		{"string_invalid", "abc", false, true, `Converting string value "abc" to bool got an error`},
		{"unsupported_type", 123, false, true, "Not support type int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatBool(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error expectation failed: wantErr %v, got %v", tt.wantErr, err)
			}

			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}

			if tt.wantErr && err != nil && tt.expectedMsg != "" {
				if !strings.Contains(err.Error(), tt.expectedMsg) {
					t.Errorf("Error message mismatch. Expected to contain: %s, got: %s",
						tt.expectedMsg, err.Error())
				}
			}
		})
	}
}

func FuzzCommonFormatBool(f *testing.F) {
	f.Add("true")
	f.Add("false")
	f.Add("1")
	f.Add("0")
	f.Add("abc")
	f.Add("")
	f.Add("TRUE")
	f.Fuzz(func(t *testing.T, s string) {
		v, err := FormatBool(s)
		if s == "" {
			if err != nil || v {
				t.Fatalf("FormatBool(\"\") = %t, %v; want false, nil", v, err)
			}
			return
		}
		parsed, parseErr := strconv.ParseBool(s)
		if parseErr != nil {
			if err == nil {
				t.Fatalf("FormatBool(%q) = %t, want error like strconv.ParseBool: %v", s, v, parseErr)
			}
			return
		}
		if err != nil || v != parsed {
			t.Fatalf("FormatBool(%q) = %t, %v; want %t, nil", s, v, err, parsed)
		}
	})
}

func TestUnitCommonFormatFloat64(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expected    float64
		wantErr     bool
		expectedMsg string
	}{
		{"nil", nil, 0.0, false, ""},
		{"float64", 3.14, 3.14, false, ""},
		{"float32", float32(3.0), 3.0, false, ""},
		{"int64", int64(123), 123.0, false, ""},
		{"int32", int32(456), 456.0, false, ""},
		{"int", 789, 789.0, false, ""},
		{"string_number", "3.14", 3.14, false, ""},
		{"string_negative", "-2.5", -2.5, false, ""},
		{"string_empty", "", 0.0, false, ""},
		{"json_number", json.Number("123.45"), 123.45, false, ""},
		{"string_invalid", "abc", 0.0, true, `Converting string value "abc" to float64 got an error`},
		{"json_number_invalid", json.Number("abc"), 0.0, true, "Converting json.Number value abc to float64 got an error"},
		{"unsupported_type", true, 0.0, true, "Not support type bool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatFloat64(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error expectation failed: wantErr %v, got %v", tt.wantErr, err)
			}

			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}

			if tt.wantErr && err != nil && tt.expectedMsg != "" {
				if !strings.Contains(err.Error(), tt.expectedMsg) {
					t.Errorf("Error message mismatch. Expected to contain: %s, got: %s",
						tt.expectedMsg, err.Error())
				}
			}
		})
	}
}

func FuzzCommonFormatFloat64(f *testing.F) {
	f.Add("3.14")
	f.Add("")
	f.Add("abc")
	f.Add("1e309")
	f.Add("-0.0")
	f.Add("NaN")
	f.Add("Inf")
	f.Fuzz(func(t *testing.T, s string) {
		v, err := FormatFloat64(s)
		if s == "" {
			if err != nil || v != 0 {
				t.Fatalf("FormatFloat64(\"\") = %f, %v; want 0, nil", v, err)
			}
			return
		}
		parsed, parseErr := strconv.ParseFloat(s, 64)
		if parseErr != nil {
			if err == nil {
				t.Fatalf("FormatFloat64(%q) = %f, want error like strconv.ParseFloat: %v", s, v, parseErr)
			}
			return
		}
		if err != nil {
			t.Fatalf("FormatFloat64(%q) returned error %v but strconv.ParseFloat succeeded: %f", s, err, parsed)
		}
		sameNaN := math.IsNaN(v) && math.IsNaN(parsed)
		if !sameNaN && v != parsed {
			t.Fatalf("FormatFloat64(%q) = %f, want %f", s, v, parsed)
		}
	})
}
