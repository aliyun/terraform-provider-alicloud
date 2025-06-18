package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitCommonExpandStringList(t *testing.T) {
	tests := []struct {
		name   string
		input  []interface{}
		expect []string
	}{
		{"empty list", []interface{}{}, []string{}},
		{"list contains items", []interface{}{"a", "b"}, []string{"a", "b"}},
		{"list contains nil", []interface{}{nil, "c"}, []string{"c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandStringList(tt.input)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestUnitCommonConvertListStringToListInterface(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := []interface{}{"a", "b", "c"}
	result := convertListStringToListInterface(input)
	assert.Equal(t, expected, result)
}

func TestUnitCommonUserDataHashSum(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Base64 input", "dGVzdA==", "test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := userDataHashSum(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitCommonWriteToFile(t *testing.T) {
	tmpFile := fmt.Sprintf("%s/testfile.txt", t.TempDir())

	t.Run("write content", func(t *testing.T) {
		err := writeToFile(tmpFile, "test content")
		assert.NoError(t, err)
		data, _ := os.ReadFile(tmpFile)
		assert.Equal(t, "test content", string(data))
	})

	t.Run("write json", func(t *testing.T) {
		content := map[string]interface{}{"key": "value"}
		err := writeToFile(tmpFile, content)
		assert.NoError(t, err)
		data, _ := os.ReadFile(tmpFile)
		var result map[string]interface{}
		json.Unmarshal(data, &result)
		assert.Equal(t, content["key"], result["key"])
	})
}

func TestUnitCommonComputePeriodByUnit(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		createTime  string
		endTime     string
		unit        string
		expect      int
		expectError bool
	}{
		{
			"cal month",
			now.Format(time.RFC3339),
			now.AddDate(0, 2, 0).Format(time.RFC3339),
			"Month",
			2,
			false,
		},
		{
			"cal week",
			now.Format(time.RFC3339),
			now.AddDate(0, 0, 14).Format(time.RFC3339),
			"Week",
			2,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			period, err := computePeriodByUnit(
				tt.createTime,
				tt.endTime,
				0,
				tt.unit,
			)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.InDelta(t, tt.expect, period, 1)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnitCommonPaymentTypeConversion(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"PayAsYouGo", "PostPaid"},
		{"Subscription", "PrePaid"},
		{"Other", "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.input.(string), func(t *testing.T) {
			result := convertPaymentTypeToChargeType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitCommonIPv6Compression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"2001:0db8:85a3:0000:0000:8a2e:0370:7334/128",
			"2001:db8:85a3::8a2e:370:7334/128",
		},
		{
			"2001:0db8:85a3::8a2e:0370:7334",
			"2001:db8:85a3::8a2e:370:7334",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := compressIPv6OrCIDR(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitCommonUniqueStrings(t *testing.T) {
	input := []string{"a", "b", "a", "c"}
	expected := []string{"a", "b", "c"}
	result := Unique(input)
	assert.ElementsMatch(t, expected, result)
}

func TestUnitCommonInterfaceToBool(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"true", true},
		{true, true},
		{"false", false},
		{1, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.input), func(t *testing.T) {
			result := Interface2Bool(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitCommonIsProtocolValid(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"ValidHTTP", "http", true},
		{"ValidHTTPS", "https", true},
		{"ValidTCP", "tcp", true},
		{"ValidUDP", "udp", true},
		{"InvalidProtocol", "icmp", false},
		{"UpperCaseHTTP", "HTTP", false},
		{"MixedCaseHttP", "HttP", false},
		{"EmptyString", "", false},
		{"PartialMatch", "htt", false},
		{"ExtraCharacters", "http ", false},
		{"NumberSuffix", "http1", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isProtocolValid(tc.input)
			if result != tc.expected {
				t.Errorf("input: %s, except: %v, result: %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestUnitCommonExpandIntList(t *testing.T) {
	testCases := []struct {
		name     string
		input    []interface{}
		expected []int
	}{
		{"EmptySlice", []interface{}{}, []int{}},
		{"SingleInt", []interface{}{42}, []int{42}},
		{"MixedTypes", []interface{}{1, 2, 3}, []int{1, 2, 3}},
		{"NegativeValues", []interface{}{-5, 0, 10}, []int{-5, 0, 10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := expandIntList(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("input: %v, except: %v, result: %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestUnitCommonConvertListToJsonString(t *testing.T) {
	testCases := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{"EmptySlice", []interface{}{}, ""},
		{"StringValues", []interface{}{"a", "b"}, `["a","b"]`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convertListToJsonString(tc.input)
			if result != tc.expected {
				t.Errorf("input: %v, except: %q, result: %q", tc.input, tc.expected, result)
			}
		})
	}
}

func TestUnitCommonBase64Encoding(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"SingleValue", []string{"hello"}, []string{"hello"}},
		{"MultipleValues", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded := encodeToBase64String(tc.input)
			decoded, err := decodeFromBase64String(encoded)
			if err != nil {
				t.Fatalf("decode err: %v", err)
			}

			if !reflect.DeepEqual(decoded, tc.expected) {
				t.Errorf("input: %v, encode: %s, decode: %v", tc.input, encoded, decoded)
			}
		})
	}
}

func TestUnitCommonConvertJsonStringToMap(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]interface{}
		hasError bool
	}{
		{"ValidJSON", `{"key":"value"}`, map[string]interface{}{"key": "value"}, false},
		{"EmptyJSON", `{}`, map[string]interface{}{}, false},
		{"InvalidJSON", `{invalid}`, nil, true},
		{"NestedJSON", `{"a":1,"b":{"c":2}}`, map[string]interface{}{"a": 1.0, "b": map[string]interface{}{"c": 2.0}}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertJsonStringToMap(tc.input)

			if tc.hasError {
				if err == nil {
					t.Error("except error but return nil")
				}
			} else {
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("input: %q, except: %v, result: %v", tc.input, tc.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonConvertListToCommaSeparate(t *testing.T) {
	testCases := []struct {
		name     string
		input    []interface{}
		expected string
	}{
		{"EmptySlice", []interface{}{}, ""},
		{"SingleString", []interface{}{"a"}, "a"},
		{"MultipleStrings", []interface{}{"a", "b", "c"}, "a,b,c"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := convertListToCommaSeparate(tc.input)
			if result != tc.expected {
				t.Errorf("input: %v, except: %q, result: %q", tc.input, tc.expected, result)
			}
		})
	}
}

func TestUnitCommonFilterEmptyStrings(t *testing.T) {
	testCases := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{"AllEmpty", []interface{}{"", nil, 0}, []interface{}{nil, 0}},
		{"MixedValues", []interface{}{"a", "", "c", nil}, []interface{}{"a", "c", nil}},
		{"NoEmpty", []interface{}{1, "b", true}, []interface{}{1, "b", true}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filterEmptyStrings(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("input: %v, except: %v, result: %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestUnitCommonTypeConversionFunctions(t *testing.T) {
	t.Run("BoolToString", func(t *testing.T) {
		if convertBoolToString(true) != "true" || convertBoolToString(false) != "false" {
			t.Error("BoolToString failed")
		}
	})

	t.Run("StringToBool", func(t *testing.T) {
		if !convertStringToBool("true") || convertStringToBool("false") {
			t.Error("StringToBool failed")
		}
	})

	t.Run("IntToString", func(t *testing.T) {
		if convertIntergerToString(42) != "42" {
			t.Error("IntToString failed")
		}
	})

	t.Run("FloatToString", func(t *testing.T) {
		result := convertFloat64ToString(3.14)
		if !strings.Contains(result, "3.14") {
			t.Errorf("convertFloat64ToString failed: %s", result)
		}
	})
}

func TestUnitCommonConvertJsonStringToList(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []interface{}
		hasError bool
	}{
		{"ValidArray", `[1, "a", true]`, []interface{}{1.0, "a", true}, false},
		{"EmptyArray", `[]`, []interface{}{}, false},
		{"InvalidJSON", `[1, "a"`, nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertJsonStringToList(tc.input)

			if tc.hasError {
				if err == nil {
					t.Error("except error but return nil")
				}
			} else {
				if err != nil {
					t.Fatalf("err: %v", err)
				}
				if !reflect.DeepEqual(result, tc.expected) {
					t.Errorf("input: %q, except: %v, result: %v", tc.input, tc.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonExpandArrayToMap(t *testing.T) {
	inputMap := map[string]interface{}{"existing": "value"}
	array := []interface{}{"a", "b", "c"}
	key := "testKey"

	result := expandArrayToMap(inputMap, array, key)

	expected := map[string]interface{}{
		"existing":  "value",
		"testKey.1": "a",
		"testKey.2": "b",
		"testKey.3": "c",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("except %v, result %v", expected, result)
	}
}

func TestUnitCommonJsonConversionFunctions(t *testing.T) {
	t.Run("JSONStringToObject", func(t *testing.T) {
		input := `{"key":"value"}`
		result := convertJsonStringToObject(input)
		expected := map[string]interface{}{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("input: %q, except: %v, result: %v", input, expected, result)
		}
	})

	t.Run("ObjectToJSONString", func(t *testing.T) {
		input := map[string]interface{}{"key": "value"}
		result := convertObjectToJsonString(input)
		if result != `{"key":"value"}` {
			t.Errorf("ObjectToJSONString error: %s", result)
		}
	})

	t.Run("MapToJSONString", func(t *testing.T) {
		input := map[string]interface{}{"key": "value"}
		result, err := convertMaptoJsonString(input)
		if err != nil || result != `{"key":"value"}` {
			t.Errorf("err: %s, 错误: %v", result, err)
		}
	})

	t.Run("MapToJSONIgnoreError", func(t *testing.T) {
		input := map[string]interface{}{"key": make(chan int)}
		result := convertMapToJsonStringIgnoreError(input)
		if result != "" {
			t.Error("err:", result)
		}
	})

	t.Run("InterfaceToJSON", func(t *testing.T) {
		input := struct{ Key string }{"value"}
		result, err := convertInterfaceToJsonString(input)
		if err != nil || result != `{"Key":"value"}` {
			t.Errorf("result: %s, err: %v", result, err)
		}
	})
}
