package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
		{"MultipleContainsNilStrings", []interface{}{"a", nil, "c"}, "a,c"},
		{"MultipleContainsBooleanIntegerStrings", []interface{}{"a", false, 1, "c"}, "a,false,1,c"},
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

func TestUnitCommonConvertJsonStringToStringList(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []interface{}
		panics   bool
	}{
		{
			name:     "ValidIntegerSlice",
			input:    []interface{}{1, 2, 3},
			expected: []interface{}{"1", "2", "3"},
			panics:   false,
		},
		{
			name:     "ValidMixedSlice",
			input:    []interface{}{"10", 5.5, 1},
			expected: []interface{}{"10", "5", "1"}, // formatInt(true)=1
			panics:   false,
		},
		{
			name:     "EmptySlice",
			input:    []interface{}{},
			expected: []interface{}{},
			panics:   false,
		},
		{
			name:     "InvalidType",
			input:    "not a slice",
			expected: nil,
			panics:   true,
		},
		{
			name:     "JsonNumber",
			input:    []interface{}{json.Number("10")},
			expected: []interface{}{"10"},
			panics:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.panics {
					t.Errorf("Unexpected panic: %v", r)
				} else if r == nil && tt.panics {
					t.Error("Expected panic did not occur")
				}
			}()

			result := convertJsonStringToStringList(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("Length mismatch: got %d, want %d", len(result), len(tt.expected))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Index %d: got %v, want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestUnitCommonConvertListMapToJsonString(t *testing.T) {
	tests := []struct {
		name     string
		input    []map[string]interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "EmptySlice",
			input:    []map[string]interface{}{},
			expected: "[]",
			wantErr:  false,
		},
		{
			name: "SingleMap",
			input: []map[string]interface{}{
				{"key1": "value1", "key2": 42},
			},
			expected: `[{"key1":"value1","key2":42}]`,
			wantErr:  false,
		},
		{
			name: "MultipleMaps",
			input: []map[string]interface{}{
				{"a": "b"},
				{"c": 123, "d": true},
			},
			expected: `[{"a":"b"},{"c":123,"d":true}]`,
			wantErr:  false,
		},
		{
			name: "WithNilElement",
			input: []map[string]interface{}{
				{"valid": "data"},
				nil,
				{"another": "map"},
			},
			expected: `[{"valid":"data"},{"another":"map"}]`,
			wantErr:  false,
		},
		{
			name: "ComplexDataTypes",
			input: []map[string]interface{}{
				{
					"string": "text",
					"int":    100,
					"float":  3.14,
					"bool":   true,
					"slice":  []interface{}{1, 2, 3},
					"map":    map[string]interface{}{"nested": "value"},
				},
			},
			expected: `[{"bool":true,"float":3.14,"int":100,"map":{"nested":"value"},"slice":[1,2,3],"string":"text"}]`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertListMapToJsonString(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error expectation failed: wantErr %v, got %v", tt.wantErr, err)
			}

			var expectedObj, resultObj interface{}
			if err := json.Unmarshal([]byte(tt.expected), &expectedObj); err != nil {
				t.Fatal("Invalid expected JSON:", err)
			}
			if err := json.Unmarshal([]byte(result), &resultObj); err != nil {
				t.Fatal("Invalid result JSON:", err)
			}

			if !reflect.DeepEqual(expectedObj, resultObj) {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonConvertIntegerToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.Integer
		expected int
		wantErr  bool
	}{
		{
			name:     "ValidNumber",
			input:    requests.Integer("123"),
			expected: 123,
			wantErr:  false,
		},
		{
			name:     "EmptyString",
			input:    requests.Integer(""),
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "WhitespaceString",
			input:    requests.Integer("   "),
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "NegativeNumber",
			input:    requests.Integer("-456"),
			expected: -456,
			wantErr:  false,
		},
		{
			name:     "MaxInt",
			input:    requests.Integer("2147483647"),
			expected: 2147483647,
			wantErr:  false,
		},
		{
			name:     "MinInt",
			input:    requests.Integer("-2147483648"),
			expected: -2147483648,
			wantErr:  false,
		},
		{
			name:     "NonNumericString",
			input:    requests.Integer("abc"),
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "PartialNumeric",
			input:    requests.Integer("123abc"),
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertIntegerToInt(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error expectation failed: wantErr %v, got %v", tt.wantErr, err)
			}

			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}

			if tt.wantErr && err != nil {
				expectedErrMsg := fmt.Sprintf("Converting integer %s to int got an error", tt.input)
				if !strings.Contains(err.Error(), expectedErrMsg) {
					t.Errorf("Error message mismatch. Expected to contain: %s, got: %s",
						expectedErrMsg, err.Error())
				}
			}
		})
	}
}

func TestUnitCommonFormatInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
		panic    bool
	}{
		{"nil", nil, 0, false},
		{"float64", 3.14, 3, false},
		{"float32", float32(3.0), 3, false},
		{"int64", int64(123), 123, false},
		{"int32", int32(456), 456, false},
		{"int", 789, 789, false},
		{"string_number", "123", 123, false},
		{"string_empty", "", 0, false},
		{"json_number", json.Number("123"), 123, false},
		{"string_invalid", "abc", 0, true},
		{"unsupported_type", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.panic {
					t.Error("Expected panic but no panic occurred")
				}
			}()

			result := formatInt(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonFormatBool(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
		panic    bool
	}{
		{"nil", nil, false, false},
		{"bool_true", true, true, false},
		{"bool_false", false, false, false},
		{"string_true", "true", true, false},
		{"string_false", "false", false, false},
		{"string_invalid", "abc", false, true},
		{"unsupported_type", 123, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.panic {
					t.Error("Expected panic but no panic occurred")
				}
			}()

			result := formatBool(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonFormatFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		panic    bool
	}{
		{"nil", nil, 0.0, false},
		{"float64", 3.14, 3.14, false},
		{"float32", float32(3.0), 3.0, false},
		{"int64", int64(123), 123.0, false},
		{"int32", int32(456), 456.0, false},
		{"int", 789, 789.0, false},
		{"string_number", "3.14", 3.14, false},
		{"string_empty", "", 0.0, false},
		{"json_number", json.Number("123.45"), 123.45, false},
		{"string_invalid", "abc", 0.0, true},
		{"unsupported_type", true, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.panic {
					t.Error("Expected panic but no panic occurred")
				}
			}()

			result := formatFloat64(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonConvertArrayObjectToJsonString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "integer array",
			input:    []int{1, 2, 3},
			expected: "[1,2,3]",
			wantErr:  false,
		},
		{
			name:     "string array",
			input:    []string{"a", "b", "c"},
			expected: `["a","b","c"]`,
			wantErr:  false,
		},
		{
			name:     "mixed array",
			input:    []interface{}{1, "a", true},
			expected: `[1,"a",true]`,
			wantErr:  false,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "null",
			wantErr:  false,
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			expected: "[]",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertArrayObjectToJsonString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertArrayObjectToJsonString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("convertArrayObjectToJsonString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonConvertArrayToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		sep      string
		expected string
	}{
		{
			name:     "string array with comma",
			input:    []interface{}{"a", "b", "c"},
			sep:      ",",
			expected: "a,b,c",
		},
		{
			name:     "integer array with dash",
			input:    []interface{}{1, 2, 3},
			sep:      "-",
			expected: "1-2-3",
		},
		{
			name:     "mixed array with space",
			input:    []interface{}{1, "a", true},
			sep:      " ",
			expected: "1 a true",
		},
		{
			name:     "nil input",
			input:    nil,
			sep:      ",",
			expected: "",
		},
		{
			name:     "empty array",
			input:    []interface{}{},
			sep:      ",",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertArrayToString(tt.input, tt.sep)
			if result != tt.expected {
				t.Errorf("convertArrayToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonSplitMultiZoneId(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "valid multi-zone",
			input:    "prefixMAZ(a,b,c)",
			expected: []string{"prefixa", "prefixb", "prefixc"},
		},
		{
			name:     "single zone",
			input:    "prefixMAZ(zone)",
			expected: []string{"prefixzone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitMultiZoneId(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitMultiZoneId() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonCase2Camel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single word", "hello", "Hello"},
		{"two words", "hello_world", "HelloWorld"},
		{"multiple words", "this_is_a_test", "ThisIsATest"},
		{"empty string", "", ""},
		{"no underscores", "alreadyCamel", "AlreadyCamel"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Case2Camel(tt.input)
			if result != tt.expected {
				t.Errorf("Case2Camel() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonFirstLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"capitalized", "Hello", "hello"},
		{"already lowercase", "hello", "hello"},
		{"single character", "A", "a"},
		{"empty string", "", ""},
		{"mixed case", "GoLang", "goLang"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FirstLower(tt.input)
			if result != tt.expected {
				t.Errorf("FirstLower() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonSplitSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		chunk    int
		expected [][]interface{}
	}{
		{
			name:     "exact division",
			input:    []interface{}{1, 2, 3, 4},
			chunk:    2,
			expected: [][]interface{}{{1, 2}, {3, 4}},
		},
		{
			name:     "remainder",
			input:    []interface{}{1, 2, 3, 4, 5},
			chunk:    2,
			expected: [][]interface{}{{1, 2}, {3, 4}, {5}},
		},
		{
			name:     "chunk larger than slice",
			input:    []interface{}{1, 2, 3},
			chunk:    5,
			expected: [][]interface{}{{1, 2, 3}},
		},
		{
			name:     "empty slice",
			input:    []interface{}{},
			chunk:    3,
			expected: nil,
		},
		{
			name:     "single chunk",
			input:    []interface{}{1, 2, 3},
			chunk:    3,
			expected: [][]interface{}{{1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitSlice(tt.input, tt.chunk)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonNewInstanceDiff(t *testing.T) {

	tests := []struct {
		name           string
		state          *terraform.InstanceState
		attributes     map[string]interface{}
		attributesDiff map[string]interface{}
		expectChanges  map[string]terraform.ResourceAttrDiff
	}{
		{
			name: "simple attribute change",
			state: &terraform.InstanceState{
				ID: "vpc-123456",
				Attributes: map[string]string{
					"name":       "old_vpc_name",
					"cidr_block": "192.168.0.0/16",
				},
			},
			attributes: map[string]interface{}{
				"name":       "old_vpc_name",
				"cidr_block": "192.168.0.0/16",
			},
			attributesDiff: map[string]interface{}{
				"name": "new_vpc_name",
			},
			expectChanges: map[string]terraform.ResourceAttrDiff{
				"name": {Old: "old_vpc_name", New: "new_vpc_name"},
			},
		},
		{
			name: "add new attribute",
			state: &terraform.InstanceState{
				ID: "vpc-234567",
				Attributes: map[string]string{
					"cidr_block": "10.0.0.0/16",
				},
			},
			attributes: map[string]interface{}{
				"cidr_block": "10.0.0.0/16",
			},
			attributesDiff: map[string]interface{}{
				"description": "new description",
			},
			expectChanges: map[string]terraform.ResourceAttrDiff{
				"description": {Old: "", New: "new description"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := newInstanceDiff("alicloud_vpc", tt.attributes, tt.attributesDiff, tt.state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for attr, expectedDiff := range tt.expectChanges {
				actualDiff, ok := diff.Attributes[attr]
				if !ok {
					t.Errorf("missing expected diff for attribute: %s", attr)
					continue
				}

				if actualDiff.Old != expectedDiff.Old || actualDiff.New != expectedDiff.New {
					t.Errorf("attribute %s mismatch:\n  expected: Old=%q, New=%q\n  actual:   Old=%q, New=%q",
						attr, expectedDiff.Old, expectedDiff.New, actualDiff.Old, actualDiff.New)
				}
			}

			for attr := range diff.Attributes {
				found := false
				for expectedAttr := range tt.expectChanges {
					if attr == expectedAttr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("unexpected diff attribute: %s (Old=%q, New=%q)",
						attr, diff.Attributes[attr].Old, diff.Attributes[attr].New)
				}
			}

		})
	}
}

func TestUnitCommonCompareMapWithIgnoreEquivalent(t *testing.T) {
	tests := []struct {
		name   string
		map1   map[string]interface{}
		map2   map[string]interface{}
		ignore []string
		expect bool
	}{
		{
			name:   "EqualMapsNoIgnore",
			map1:   map[string]interface{}{"a": 1, "b": "test"},
			map2:   map[string]interface{}{"a": 1, "b": "test"},
			ignore: []string{},
			expect: true,
		},
		{
			name:   "EqualMapsWithIgnore",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 1, "b": "test", "c": false},
			ignore: []string{"c"},
			expect: true,
		},
		{
			name:   "DifferentMapsIgnoreNotApplied",
			map1:   map[string]interface{}{"a": 1, "b": "test"},
			map2:   map[string]interface{}{"a": 1, "b": "different"},
			ignore: []string{},
			expect: false,
		},
		{
			name:   "DifferentLengthMaps",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 1, "b": "test"},
			ignore: []string{"c"},
			expect: false,
		},
		{
			name:   "IgnoreKeyMissingInOneMap",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 1, "b": "test"},
			ignore: []string{"c"},
			expect: false, // 长度不同
		},
		{
			name:   "IgnoreMultipleKeys",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true, "d": 3.14},
			map2:   map[string]interface{}{"a": 1, "b": "test", "c": false, "d": 6.28},
			ignore: []string{"c", "d"},
			expect: true,
		},
		{
			name:   "DifferentValueNonIgnoredKey",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 2, "b": "test", "c": true},
			ignore: []string{"c"},
			expect: false,
		},
		{
			name:   "NestedMapsIgnored",
			map1:   map[string]interface{}{"a": 1, "config": map[string]interface{}{"key": "value"}},
			map2:   map[string]interface{}{"a": 1, "config": map[string]interface{}{"key": "different"}},
			ignore: []string{"config"},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareMapWithIgnoreEquivalent(tt.map1, tt.map2, tt.ignore)
			if result != tt.expect {
				t.Errorf(
					"compareMapWithIgnoreEquivalent() = %v, want %v\nMap1: %#v\nMap2: %#v\nIgnore: %v",
					result, tt.expect, tt.map1, tt.map2, tt.ignore,
				)
			}
		})
	}
}

func TestUnitCommonIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"Nil", nil, true},
		{"EmptyString", "", true},
		{"NonEmptyString", "test", false},
		{"ZeroInt", 0, true},
		{"PositiveInt", 42, false},
		{"NegativeInt", -5, true},
		{"ZeroInt8", int8(0), true},
		{"PositiveInt8", int8(1), false},
		{"ZeroInt16", int16(0), true},
		{"PositiveInt16", int16(10), false},
		{"ZeroInt32", int32(0), true},
		{"PositiveInt32", int32(20), false},
		{"ZeroInt64", int64(0), true},
		{"PositiveInt64", int64(30), false},
		{"ZeroFloat32", float32(0), true},
		{"PositiveFloat32", float32(3.14), false},
		{"ZeroFloat64", 0.0, true},
		{"PositiveFloat64", 6.28, false},
		{"EmptyMap", map[string]interface{}{}, true},
		{"NonEmptyMap", map[string]interface{}{"key": "value"}, false},
		{"NilPointer", (*int)(nil), true},
		{"ValidPointer", new(int), false},
		{"BoolFalse", false, false}, // 注意：布尔值false不被视为empty
		{"BoolTrue", true, false},
		{"Slice", []int{}, false}, // 注意：切片类型不在IsEmpty处理范围内
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmpty(%#v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonIsNil(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"Nil", nil, true},
		{"EmptyString", "", true},
		{"NonEmptyString", "test", false},
		{"EmptySlice", []interface{}{}, true},
		{"NonEmptySlice", []interface{}{1}, false},
		{"EmptyMap", map[string]interface{}{}, true},
		{"NonEmptyMap", map[string]interface{}{"key": "value"}, false},
		{"NilPointer", (*int)(nil), true},
		{"ValidPointer", new(int), false},
		{"Integer", 42, false}, // 整数类型不为nil
		{"Float", 3.14, false},
		{"Bool", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNil(tt.input)
			if result != tt.expected {
				t.Errorf("IsNil(%#v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonGetDaysBetween2Date(t *testing.T) {
	format := "2006-01-02"
	tests := []struct {
		name     string
		date1    string
		date2    string
		expected int
		wantErr  bool
	}{
		{"SameDate", "2023-01-01", "2023-01-01", 0, false},
		{"OneDayDifference", "2023-01-01", "2023-01-02", 1, false},
		{"ReverseOrder", "2023-01-02", "2023-01-01", -1, false},
		{"OneMonthDifference", "2023-01-01", "2023-02-01", 31, false},
		{"LeapYear", "2020-02-28", "2020-03-01", 2, false},
		{"NonLeapYear", "2021-02-28", "2021-03-01", 1, false},
		{"CrossYear", "2022-12-31", "2023-01-01", 1, false},
		{"InvalidDate1", "invalid", "2023-01-01", 0, true},
		{"InvalidDate2", "2023-01-01", "invalid", 0, true},
		{"DifferentFormat", "01 Jan 23", "02 Jan 23", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetDaysBetween2Date(format, tt.date1, tt.date2)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetDaysBetween2Date() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("GetDaysBetween2Date() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonConvertMapFloat64ToJsonString(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "EmptyMap",
			input:    map[string]interface{}{},
			expected: "{}",
			wantErr:  false,
		},
		{
			name: "SingleValue",
			input: map[string]interface{}{
				"key": json.Number("123.45"),
			},
			expected: `{"key":123.45}`,
			wantErr:  false,
		},
		{
			name: "MultipleValues",
			input: map[string]interface{}{
				"intValue":    json.Number("42"),
				"floatValue":  json.Number("3.14159"),
				"stringValue": json.Number("100.0"),
			},
			expected: `{"floatValue":3.14159,"intValue":42,"stringValue":100.0}`,
			wantErr:  false,
		},
		{
			name: "NegativeAndZero",
			input: map[string]interface{}{
				"negative": json.Number("-15.75"),
				"zero":     json.Number("0"),
			},
			expected: `{"negative":-15.75,"zero":0}`,
			wantErr:  false,
		},
		{
			name: "LargeNumbers",
			input: map[string]interface{}{
				"largeInt":   json.Number("1234567890"),
				"largeFloat": json.Number("123456.789012"),
			},
			expected: `{"largeFloat":123456.789012,"largeInt":1234567890}`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertMapFloat64ToJsonString(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("convertMapFloat64ToJsonString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("convertMapFloat64ToJsonString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonIsSubCollection(t *testing.T) {
	tests := []struct {
		name     string
		sub      []string
		full     []string
		expected bool
	}{
		{
			name:     "EmptySubset",
			sub:      []string{},
			full:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "ExactMatch",
			sub:      []string{"a", "b"},
			full:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "PartialMatch",
			sub:      []string{"a", "d"},
			full:     []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "CaseSensitive",
			sub:      []string{"A", "b"},
			full:     []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "DuplicateInSub",
			sub:      []string{"a", "a"},
			full:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "EmptyFull",
			sub:      []string{"a"},
			full:     []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSubCollection(tt.sub, tt.full)
			if result != tt.expected {
				t.Errorf("IsSubCollection(%v, %v) = %v, want %v", tt.sub, tt.full, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonMergeMaps(t *testing.T) {
	tests := []struct {
		name     string
		maps     []map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "SimpleMerge",
			maps: []map[string]interface{}{
				{"a": 1},
				{"b": 2},
			},
			expected: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		},
		{
			name: "NestedMapMerge",
			maps: []map[string]interface{}{
				{"a": []map[string]interface{}{{"x": 1}}},
				{"a": []map[string]interface{}{{"y": 2}}},
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"x": 1,
					"y": 2,
				},
			},
		},
		{
			name: "ComplexNestedMerge",
			maps: []map[string]interface{}{
				{
					"config": []map[string]interface{}{
						{"port": 80, "enabled": true},
					},
				},
				{
					"config": []map[string]interface{}{
						{"protocol": "http"},
					},
				},
			},
			expected: map[string]interface{}{
				"config": map[string]interface{}{
					"port":     80,
					"enabled":  true,
					"protocol": "http",
				},
			},
		},
		{
			name:     "EmptyInput",
			maps:     []map[string]interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeMaps(tt.maps...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MergeMaps() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonInArray(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		strArray []string
		expected bool
	}{
		{
			name:     "ElementExists",
			target:   "b",
			strArray: []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "ElementNotExists",
			target:   "d",
			strArray: []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "CaseSensitive",
			target:   "A",
			strArray: []string{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "EmptyArray",
			target:   "a",
			strArray: []string{},
			expected: false,
		},
		{
			name:     "DuplicateElements",
			target:   "a",
			strArray: []string{"a", "a", "b"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InArray(tt.target, tt.strArray)
			if result != tt.expected {
				t.Errorf("InArray(%s, %v) = %v, want %v", tt.target, tt.strArray, result, tt.expected)
			}
		})
	}
}

func TestUnitGetOneStringOrAllStringSlice(t *testing.T) {
	tests := []struct {
		name        string
		input       []interface{}
		expected    interface{}
		expectError bool
	}{
		{
			"one string",
			[]interface{}{"a"},
			"a",
			false,
		},
		{
			"string slice",
			[]interface{}{"a", "b", "c"},
			[]string{"a", "b", "c"},
			false,
		},
		{
			"contain empty strings",
			[]interface{}{""},
			nil,
			true,
		},
		{
			"contain empty strings",
			[]interface{}{"   "},
			nil,
			true,
		},
		{
			"contain empty strings",
			[]interface{}{"a", "", "c"},
			nil,
			true,
		},
		{
			"contain empty strings",
			[]interface{}{"a", "    ", "c"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getOneStringOrAllStringSlice(tt.input, tt.name)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUnitCommonNormalizeAndMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "basic conversion",
			input: map[string]interface{}{
				"name":   "john",
				"age":    "25",
				"active": "true",
				"admin":  "false",
			},
			expected: map[string]interface{}{
				"name":   "john",
				"age":    25,
				"active": true,
				"admin":  false,
			},
		},
		{
			name: "with nested map",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name":   "alice",
					"age":    "30",
					"active": "true",
				},
				"id": "123",
			},
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name":   "alice",
					"age":    30,
					"active": true,
				},
				"id": 123,
			},
		},
		{
			name: "with array",
			input: map[string]interface{}{
				"tags":  []interface{}{"true", "false", "42", "hello"},
				"count": "10",
			},
			expected: map[string]interface{}{
				"tags":  []interface{}{true, false, 42, "hello"},
				"count": 10,
			},
		},
		{
			name: "mixed types",
			input: map[string]interface{}{
				"name":     "bob",
				"age":      25,   // already int
				"active":   true, // already bool
				"score":    "95",
				"premium":  "true",
				"balance":  "0",
				"nickname": "bobby",
			},
			expected: map[string]interface{}{
				"name":     "bob",
				"age":      25,
				"active":   true,
				"score":    95,
				"premium":  true,
				"balance":  0,
				"nickname": "bobby",
			},
		},
		{
			name: "non-convertible strings",
			input: map[string]interface{}{
				"description": "this is a test",
				"version":     "1.2.3",
				"code":        "ABC123",
			},
			expected: map[string]interface{}{
				"description": "this is a test",
				"version":     "1.2.3",
				"code":        "ABC123",
			},
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "nil map",
			input:    nil,
			expected: nil,
		},
		{
			name: "zero values",
			input: map[string]interface{}{
				"zero_str":   "0",
				"zero_int":   0,
				"false_str":  "false",
				"false_bool": false,
			},
			expected: map[string]interface{}{
				"zero_str":   0,
				"zero_int":   0,
				"false_str":  false,
				"false_bool": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeMap(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("NormalizeMap() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUnitCommonConvertToJsonWithoutEscapeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
		hasError bool
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"key": "value",
			},
			expected: "{\"key\":\"value\"}\n",
			hasError: false,
		},
		{
			name: "map with html characters",
			input: map[string]interface{}{
				"html": "<div>Hello & World</div>",
			},
			expected: "{\"html\":\"<div>Hello & World</div>\"}\n",
			hasError: false,
		},
		{
			name: "map with special characters",
			input: map[string]interface{}{
				"special": "<>&",
			},
			expected: "{\"special\":\"<>&\"}\n",
			hasError: false,
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "john",
					"profile": map[string]interface{}{
						"description": "<p>This is a <b>bold</b> statement</p>",
					},
				},
			},
			expected: "{\"user\":{\"name\":\"john\",\"profile\":{\"description\":\"<p>This is a <b>bold</b> statement</p>\"}}}\n",
			hasError: false,
		},
		{
			name: "map with array",
			input: map[string]interface{}{
				"items": []interface{}{"<item1>", "<item2>", "& special chars"},
			},
			expected: "{\"items\":[\"<item1>\",\"<item2>\",\"& special chars\"]}\n",
			hasError: false,
		},
		{
			name: "empty map",
			input: map[string]interface{}{},
			expected: "{}\n",
			hasError: false,
		},
		{
			name: "nil map",
			input: nil,
			expected: "null\n",
			hasError: false,
		},
		{
			name: "map with various data types",
			input: map[string]interface{}{
				"string": "text",
				"int":    42,
				"float":  3.14,
				"bool":   true,
				"nil":    nil,
			},
			expected: "{\"bool\":true,\"float\":3.14,\"int\":42,\"nil\":null,\"string\":\"text\"}\n",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertToJsonWithoutEscapeHTML(tt.input)
			
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			// Compare the results
			if result != tt.expected {
				t.Errorf("Expected: %s, Got: %s", tt.expected, result)
			}
		})
	}
}
