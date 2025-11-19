package alicloud

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func TestUnitCommonGetOneStringOrAllStringSlice(t *testing.T) {
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
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: "{}\n",
			hasError: false,
		},
		{
			name:     "nil map",
			input:    nil,
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

func TestUnitCommonConvertToInterfaceArray(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: []interface{}{},
		},
		{
			name:     "already []interface{}",
			input:    []interface{}{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name:     "string slice",
			input:    []string{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name:     "int slice",
			input:    []int{1, 2, 3},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "mixed type slice",
			input:    []interface{}{1, "two", 3.0, true},
			expected: []interface{}{1, "two", 3.0, true},
		},
		{
			name:     "single value",
			input:    "single",
			expected: []interface{}{"single"},
		},
		{
			name:     "single int",
			input:    42,
			expected: []interface{}{42},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []interface{}{},
		},
		{
			name:     "bool slice",
			input:    []bool{true, false, true},
			expected: []interface{}{true, false, true},
		},
		{
			name:     "float slice",
			input:    []float64{1.1, 2.2, 3.3},
			expected: []interface{}{1.1, 2.2, 3.3},
		},
		{
			name:     "single value",
			input:    "single",
			expected: []interface{}{"single"},
		},
		{
			name:     "single int",
			input:    42,
			expected: []interface{}{42},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []interface{}{},
		},
		{
			name:     "bool slice",
			input:    []bool{true, false, true},
			expected: []interface{}{true, false, true},
		},
		{
			name:     "float slice",
			input:    []float64{1.1, 2.2, 3.3},
			expected: []interface{}{1.1, 2.2, 3.3},
		},
		{
			name:     "schema.Set with strings",
			input:    func() *schema.Set { s := schema.NewSet(schema.HashString, []interface{}{"a", "b", "c"}); return s }(),
			expected: []interface{}{"a", "b", "c"},
		},
		{
			name: "schema.Set with ints",
			input: func() *schema.Set {
				s := schema.NewSet(func(i interface{}) int { return schema.HashString(fmt.Sprintf("%v", i)) }, []interface{}{1, 2, 3})
				return s
			}(),
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "empty schema.Set",
			input:    func() *schema.Set { s := schema.NewSet(schema.HashString, []interface{}{}); return s }(),
			expected: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertToInterfaceArray(tt.input)
			if tt.name == "schema.Set with strings" || tt.name == "schema.Set with ints" {
				// For schema.Set, we need to compare elements regardless of order
				if len(result) != len(tt.expected) {
					t.Errorf("convertToInterfaceArray(%v) = %v (len %d), want %v (len %d)", tt.input, result, len(result), tt.expected, len(tt.expected))
				}
				// Create maps to count occurrences
				resultMap := make(map[interface{}]int)
				expectedMap := make(map[interface{}]int)
				for _, v := range result {
					resultMap[v]++
				}
				for _, v := range tt.expected {
					expectedMap[v]++
				}
				if !reflect.DeepEqual(resultMap, expectedMap) {
					t.Errorf("convertToInterfaceArray(%v) = %v, want %v (elements mismatch)", tt.input, result, tt.expected)
				}
			} else {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("convertToInterfaceArray(%v) = %v, want %v", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestUnitCommonConvertYamlToObject(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    map[string]interface{}
		expectError bool
	}{
		{
			"empty",
			"",
			map[string]interface{}{},
			false,
		},
		{
			"int",
			"min: 1",
			map[string]interface{}{
				"min": 1,
			},
			false,
		},
		{
			"float",
			"max: 10.01",
			map[string]interface{}{
				"max": 10.01,
			},
			false,
		},
		{
			"string",
			"rule: MayRunAs",
			map[string]interface{}{
				"rule": "MayRunAs",
			},
			false,
		},
		{
			"bool",
			"hostNetwork: true",
			map[string]interface{}{
				"hostNetwork": true,
			},
			false,
		},
		{
			"array",
			"repos:\n- registry-vpc.cn-hangzhou.aliyuncs.com/acs/\n- registry.cn-hangzhou.aliyuncs.com/acs/",
			map[string]interface{}{
				"repos": []interface{}{"registry-vpc.cn-hangzhou.aliyuncs.com/acs/", "registry.cn-hangzhou.aliyuncs.com/acs/"},
			},
			false,
		},
		{
			"object",
			"metadata:\n  name: bad\n  namespace: test-gatekeeper",
			map[string]interface{}{
				"metadata": map[string]interface{}{
					"name":      "bad",
					"namespace": "test-gatekeeper",
				},
			},
			false,
		},
		{
			"mix type",
			"object:\n    sub_object:\n      string: MustRunAs\n      second_object:\n        int: 100\n        float: 10.01\n        bool: true\n        array:\n          - \"1\"\n          - \"2\"\n          - \"3\"",
			map[string]interface{}{
				"object": map[string]interface{}{
					"sub_object": map[string]interface{}{
						"second_object": map[string]interface{}{
							"array": []interface{}{"1", "2", "3"},
							"bool":  true,
							"float": 10.01,
							"int":   100,
						},
						"string": "MustRunAs",
					},
				},
			},
			false,
		},
		{
			"error test",
			"error",
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertYamlToObject(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUnitCommonPointerFunctions(t *testing.T) {
	t.Run("StringPointer", func(t *testing.T) {
		str := "test"
		ptr := StringPointer(str)
		if ptr == nil || *ptr != str {
			t.Error("StringPointer failed")
		}
	})

	t.Run("BoolPointer", func(t *testing.T) {
		b := true
		ptr := BoolPointer(b)
		if ptr == nil || *ptr != b {
			t.Error("BoolPointer failed")
		}
	})

	t.Run("Int32Pointer", func(t *testing.T) {
		i := int32(42)
		ptr := Int32Pointer(i)
		if ptr == nil || *ptr != i {
			t.Error("Int32Pointer failed")
		}
	})

	t.Run("Int64Pointer", func(t *testing.T) {
		i := int64(100)
		ptr := Int64Pointer(i)
		if ptr == nil || *ptr != i {
			t.Error("Int64Pointer failed")
		}
	})
}

func TestUnitCommonIntMin(t *testing.T) {
	tests := []struct {
		name     string
		x        int
		y        int
		expected int
	}{
		{"x smaller", 1, 2, 1},
		{"y smaller", 5, 3, 3},
		{"equal", 4, 4, 4},
		{"negative", -5, -2, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntMin(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("IntMin(%d, %d) = %d, want %d", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonGetPagination(t *testing.T) {
	tests := []struct {
		name       string
		pageNumber int
		pageSize   int
	}{
		{"default", 1, 10},
		{"custom", 2, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPagination(tt.pageNumber, tt.pageSize)
			if result.PageNumber != tt.pageNumber || result.PageSize != tt.pageSize {
				t.Errorf("getPagination failed")
			}
		})
	}
}

func TestUnitCommonBuildClientToken(t *testing.T) {
	token1 := buildClientToken("CreateInstance")
	token2 := buildClientToken("CreateInstance")

	if token1 == "" {
		t.Error("buildClientToken returned empty string")
	}

	if len(token1) > 64 {
		t.Error("buildClientToken returned token longer than 64 characters")
	}

	if token1 == token2 {
		t.Error("buildClientToken should return different tokens")
	}

	if !strings.Contains(token1, "TF-CreateInstance") {
		t.Error("buildClientToken should contain action name")
	}
}

func TestUnitCommonGetNextpageNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    requests.Integer
		expected requests.Integer
		hasError bool
	}{
		{"page 1", requests.Integer("1"), requests.Integer("2"), false},
		{"page 10", requests.Integer("10"), requests.Integer("11"), false},
		{"invalid", requests.Integer("abc"), requests.Integer(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getNextpageNumber(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonTerraformToAPI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single word", "name", "Name"},
		{"two words", "user_name", "UserName"},
		{"multiple words", "vpc_id_test", "VpcIdTest"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := terraformToAPI(tt.input)
			if result != tt.expected {
				t.Errorf("terraformToAPI(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonCompareJsonTemplateAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		tem1     string
		tem2     string
		expected bool
		hasError bool
	}{
		{
			"equal templates",
			`{"key":"value","number":123}`,
			`{"number":123,"key":"value"}`,
			true,
			false,
		},
		{
			"different templates",
			`{"key":"value1"}`,
			`{"key":"value2"}`,
			false,
			false,
		},
		{
			"invalid json1",
			`{invalid}`,
			`{"key":"value"}`,
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compareJsonTemplateAreEquivalent(tt.tem1, tt.tem2)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonCompareArrayJsonTemplateAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		tem1     string
		tem2     string
		expected bool
		hasError bool
	}{
		{
			"equal arrays",
			`[{"key":"value"},{"number":123}]`,
			`[{"key":"value"},{"number":123}]`,
			true,
			false,
		},
		{
			"different arrays",
			`[{"key":"value1"}]`,
			`[{"key":"value2"}]`,
			false,
			false,
		},
		{
			"invalid json",
			`[{invalid}]`,
			`[{"key":"value"}]`,
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compareArrayJsonTemplateAreEquivalent(tt.tem1, tt.tem2)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonCompareYamlTemplateAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		tem1     string
		tem2     string
		expected bool
		hasError bool
	}{
		{
			"equal yaml",
			"key: value\nnumber: 123",
			"number: 123\nkey: value",
			true,
			false,
		},
		{
			"different yaml",
			"key: value1",
			"key: value2",
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compareYamlTemplateAreEquivalent(tt.tem1, tt.tem2)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonLoadFileContent(t *testing.T) {
	tmpFile := fmt.Sprintf("%s/test_load.txt", t.TempDir())
	content := "test content"
	os.WriteFile(tmpFile, []byte(content), 0644)

	result, err := loadFileContent(tmpFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(result) != content {
		t.Errorf("Expected %s, got %s", content, string(result))
	}

	// Test non-existent file
	_, err = loadFileContent("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestUnitCommonParseResourceIds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"single part", "id1", []string{"id1"}},
		{"two parts", "id1:id2", []string{"id1", "id2"}},
		{"three parts", "vpc:vsw:ecs", []string{"vpc", "vsw", "ecs"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := ParseResourceIds(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseResourceIds(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonParseResourceId(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected []string
		hasError bool
	}{
		{"valid two parts", "vpc:vsw", 2, []string{"vpc", "vsw"}, false},
		{"valid three parts", "a:b:c", 3, []string{"a", "b", "c"}, false},
		{"invalid length", "vpc:vsw", 3, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResourceId(tt.input, tt.length)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonParseResourceIdN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected []string
		hasError bool
	}{
		{"valid split", "a:b:c:d", 2, []string{"a", "b:c:d"}, false},
		{"exact match", "a:b", 2, []string{"a", "b"}, false},
		{"invalid length", "a:b", 3, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResourceIdN(tt.input, tt.length)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonParseResourceIdWithEscaped(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected []string
		hasError bool
	}{
		{"no escape", "a:b:c", 3, []string{"a", "b", "c"}, false},
		{"with escape", `a\:b:c`, 2, []string{"a:b", "c"}, false},
		{"multiple escape", `a\:b\:c:d`, 2, []string{"a:b:c", "d"}, false},
		{"invalid length", "a:b", 3, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResourceIdWithEscaped(tt.input, tt.length)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonEscapeColons(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no colon", "abc", "abc"},
		{"one colon", "a:b", `a\:b`},
		{"multiple colons", "a:b:c", `a\:b\:c`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeColons(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeColons(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonParseSlbListenerId(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		hasError bool
	}{
		{"valid 2 parts", "slb-123:http", 2, false},
		{"valid 3 parts", "slb-123:http:80", 3, false},
		{"invalid", "slb-123", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseSlbListenerId(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expected {
					t.Errorf("Expected %d parts, got %d", tt.expected, len(result))
				}
			}
		})
	}
}

func TestUnitCommonGetCenChildInstanceType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"vpc instance", "vpc-12345", "VPC", false},
		{"vbr instance", "vbr-12345", "VBR", false},
		{"ccn instance", "ccn-12345", "CCN", false},
		{"invalid", "ecs-12345", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetCenChildInstanceType(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonMapMerge(t *testing.T) {
	tests := []struct {
		name     string
		target   map[string]interface{}
		merged   map[string]interface{}
		expected map[string]interface{}
	}{
		{
			"simple merge",
			map[string]interface{}{"a": 1},
			map[string]interface{}{"b": 2},
			map[string]interface{}{"a": 1, "b": 2},
		},
		{
			"key exists not override",
			map[string]interface{}{"a": 1},
			map[string]interface{}{"a": 2},
			map[string]interface{}{"a": 2},
		},
		{
			"nested map",
			map[string]interface{}{"a": map[string]interface{}{"x": 1}},
			map[string]interface{}{"a": map[string]interface{}{"y": 2}},
			map[string]interface{}{"a": map[string]interface{}{"x": 1, "y": 2}},
		},
		{
			"merge arrays",
			map[string]interface{}{"arr": []interface{}{1, 2}},
			map[string]interface{}{"arr": []interface{}{3, 4}},
			map[string]interface{}{"arr": []interface{}{3, 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of target to avoid mutation affecting tests
			targetCopy := make(map[string]interface{})
			for k, v := range tt.target {
				targetCopy[k] = v
			}
			result := mapMerge(targetCopy, tt.merged)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonInterface2String(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 123, "123"},
		{"float", 3.14, "3.14"},
		{"bool", true, "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Interface2String(tt.input)
			if result != tt.expected {
				t.Errorf("Interface2String(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonInterface2StrSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []interface{}
		expected []string
	}{
		{"string slice", []interface{}{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"mixed slice", []interface{}{"a", 1, true}, []string{"a", "1", "true"}},
		{"empty", []interface{}{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Interface2StrSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Interface2StrSlice(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonStr2InterfaceSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []interface{}
	}{
		{"string slice", []string{"a", "b", "c"}, []interface{}{"a", "b", "c"}},
		{"empty", []string{}, []interface{}{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Str2InterfaceSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Str2InterfaceSlice(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonExpandSingletonToList(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []interface{}
	}{
		{"string", "test", []interface{}{"test"}},
		{"int", 123, []interface{}{123}},
		{"map", map[string]string{"key": "value"}, []interface{}{map[string]string{"key": "value"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandSingletonToList(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expandSingletonToList(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonMD5(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"test string", []byte("test"), "098f6bcd4621d373cade4e832627b4f6"},
		{"empty", []byte(""), "d41d8cd98f00b204e9800998ecf8427e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MD5(tt.input)
			if result != tt.expected {
				t.Errorf("MD5(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonConvertTags(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected []map[string]interface{}
	}{
		{
			"basic tags",
			map[string]interface{}{"env": "prod", "app": "test"},
			[]map[string]interface{}{
				{"Key": "env", "Value": "prod"},
				{"Key": "app", "Value": "test"},
			},
		},
		{
			"empty map",
			map[string]interface{}{},
			[]map[string]interface{}{},
		},
		{
			"nil value",
			map[string]interface{}{"key": nil},
			[]map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTags(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d tags, got %d", len(tt.expected), len(result))
			}
		})
	}
}

func TestUnitCommonConvertTagsForKms(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected int
	}{
		{"basic tags", map[string]interface{}{"env": "prod", "app": "test"}, 2},
		{"empty map", map[string]interface{}{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTagsForKms(tt.input)
			if len(result) != tt.expected {
				t.Errorf("Expected %d tags, got %d", tt.expected, len(result))
			}
			for _, tag := range result {
				if _, ok := tag["TagKey"]; !ok {
					t.Error("Missing TagKey field")
				}
				if _, ok := tag["TagValue"]; !ok {
					t.Error("Missing TagValue field")
				}
			}
		})
	}
}

func TestUnitCommonExpandTagsToMap(t *testing.T) {
	originMap := map[string]interface{}{"existing": "value"}
	tags := []map[string]interface{}{
		{"Key": "env", "Value": "prod"},
		{"Key": "app", "Value": "test"},
	}

	result := expandTagsToMap(originMap, tags)

	if result["Tag.1.Key"] != "env" {
		t.Error("Tag.1.Key not set correctly")
	}
	if result["Tag.1.Value"] != "prod" {
		t.Error("Tag.1.Value not set correctly")
	}
	if result["Tag.2.Key"] != "app" {
		t.Error("Tag.2.Key not set correctly")
	}
	if result["existing"] != "value" {
		t.Error("Original value was modified")
	}
}

func TestUnitCommonExpandTagsToMapWithTags(t *testing.T) {
	originMap := map[string]interface{}{"existing": "value"}
	tags := []map[string]interface{}{
		{"Key": "env", "Value": "prod"},
	}

	result := expandTagsToMapWithTags(originMap, tags)

	if result["Tags.1.Key"] != "env" {
		t.Error("Tags.1.Key not set correctly")
	}
	if result["Tags.1.Value"] != "prod" {
		t.Error("Tags.1.Value not set correctly")
	}
}

func TestUnitCommonConvertChargeTypeToPaymentType(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected interface{}
	}{
		{"PostPaid", "PayAsYouGo"},
		{"Postpaid", "PayAsYouGo"},
		{"PrePaid", "Subscription"},
		{"Prepaid", "Subscription"},
		{"Other", "Other"},
	}

	for _, tt := range tests {
		t.Run(tt.input.(string), func(t *testing.T) {
			result := convertChargeTypeToPaymentType(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonBytesToTB(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected float64
	}{
		{"1 TB", 1099511627776, 1.0},
		{"2 TB", 2199023255552, 2.0},
		{"0 bytes", 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bytesToTB(tt.input)
			if result != tt.expected {
				t.Errorf("bytesToTB(%d) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonIsPagingRequest(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		expected bool
	}{
		{
			"no page_number",
			map[string]interface{}{},
			false,
		},
		{
			"page_number is 0",
			map[string]interface{}{"page_number": 0},
			false,
		},
		{
			"page_number is positive",
			map[string]interface{}{"page_number": 1},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"page_number": {Type: schema.TypeInt, Optional: true},
			}, tt.data)

			result := isPagingRequest(d)
			if result != tt.expected {
				t.Errorf("isPagingRequest() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnitCommonSetPagingRequest(t *testing.T) {
	tests := []struct {
		name             string
		data             map[string]interface{}
		maxPageSize      int
		expectedPageNum  int
		expectedPageSize int
	}{
		{
			"default values",
			map[string]interface{}{},
			0,
			1,
			PageSizeLarge,
		},
		{
			"custom page_number",
			map[string]interface{}{"page_number": 2},
			0,
			2,
			PageSizeLarge,
		},
		{
			"custom page_size",
			map[string]interface{}{"page_size": 20},
			0,
			1,
			20,
		},
		{
			"custom max page size",
			map[string]interface{}{},
			30,
			1,
			30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"page_number": {Type: schema.TypeInt, Optional: true},
				"page_size":   {Type: schema.TypeInt, Optional: true},
			}, tt.data)

			request := make(map[string]interface{})
			setPagingRequest(d, request, tt.maxPageSize)

			if request["PageNumber"] != tt.expectedPageNum {
				t.Errorf("PageNumber = %v, want %v", request["PageNumber"], tt.expectedPageNum)
			}
			if request["PageSize"] != tt.expectedPageSize {
				t.Errorf("PageSize = %v, want %v", request["PageSize"], tt.expectedPageSize)
			}
		})
	}
}

func TestUnitCommonRpcParam(t *testing.T) {
	param := rpcParam("POST", "2020-01-01", "CreateInstance")

	if *param.Method != "POST" {
		t.Error("Method not set correctly")
	}
	if *param.Version != "2020-01-01" {
		t.Error("Version not set correctly")
	}
	if *param.Action != "CreateInstance" {
		t.Error("Action not set correctly")
	}
	if *param.Style != "RPC" {
		t.Error("Style should be RPC")
	}
}

func TestUnitCommonRoaParam(t *testing.T) {
	param := roaParam("GET", "2020-01-01", "ListInstances", "/instances")

	if *param.Method != "GET" {
		t.Error("Method not set correctly")
	}
	if *param.Version != "2020-01-01" {
		t.Error("Version not set correctly")
	}
	if *param.Action != "ListInstances" {
		t.Error("Action not set correctly")
	}
	if *param.Pathname != "/instances" {
		t.Error("Pathname not set correctly")
	}
	if *param.Style != "ROA" {
		t.Error("Style should be ROA")
	}
}

func TestUnitCommonXmlParam(t *testing.T) {
	param := xmlParam("PUT", "2020-01-01", "UpdateConfig", "/config")

	if *param.Method != "PUT" {
		t.Error("Method not set correctly")
	}
	if *param.ReqBodyType != "xml" {
		t.Error("ReqBodyType should be xml")
	}
	if *param.BodyType != "xml" {
		t.Error("BodyType should be xml")
	}
}

func TestUnitCommonJsonXmlParam(t *testing.T) {
	param := jsonXmlParam("POST", "2020-01-01", "CreateResource", "/resource")

	if *param.ReqBodyType != "json" {
		t.Error("ReqBodyType should be json")
	}
	if *param.BodyType != "xml" {
		t.Error("BodyType should be xml")
	}
}

func TestUnitCommonXmlJsonParam(t *testing.T) {
	param := xmlJsonParam("POST", "2020-01-01", "CreateResource", "/resource")

	if *param.ReqBodyType != "xml" {
		t.Error("ReqBodyType should be xml")
	}
	if *param.BodyType != "json" {
		t.Error("BodyType should be json")
	}
}

func TestUnitCommonMyMapMarshalXML(t *testing.T) {
	// Test empty map - should return nil error for empty map
	emptyMap := MyMap{}
	var emptyResult strings.Builder
	emptyEncoder := xml.NewEncoder(&emptyResult)
	err := emptyMap.MarshalXML(emptyEncoder, xml.StartElement{Name: xml.Name{Local: "root"}})
	if err != nil {
		t.Errorf("MarshalXML failed for empty map: %v", err)
	}

	// Test non-empty map
	m := MyMap{
		"key1": "value1",
		"key2": "value2",
	}
	var result strings.Builder
	encoder := xml.NewEncoder(&result)
	err = m.MarshalXML(encoder, xml.StartElement{Name: xml.Name{Local: "root"}})
	if err != nil {
		t.Errorf("MarshalXML failed: %v", err)
	}
	encoder.Flush()
	xmlStr := result.String()
	if !strings.Contains(xmlStr, "key1") || !strings.Contains(xmlStr, "value1") {
		t.Errorf("MarshalXML output doesn't contain expected keys: %s", xmlStr)
	}
}

func TestUnitCommonDebugOn(t *testing.T) {
	// Save original value
	originalDebug := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", originalDebug)

	// Test when DEBUG is not set
	os.Setenv("DEBUG", "")
	if debugOn() {
		t.Error("debugOn should return false when DEBUG is not set")
	}

	// Test when DEBUG contains terraform
	os.Setenv("DEBUG", "terraform")
	if !debugOn() {
		t.Error("debugOn should return true when DEBUG contains terraform")
	}

	// Test when DEBUG contains terraform with other values
	os.Setenv("DEBUG", "something,terraform,other")
	if !debugOn() {
		t.Error("debugOn should return true when DEBUG contains terraform")
	}

	// Test when DEBUG does not contain terraform
	os.Setenv("DEBUG", "other")
	if debugOn() {
		t.Error("debugOn should return false when DEBUG does not contain terraform")
	}
}

func TestUnitCommonTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with spaces", "  hello  ", "hello"},
		{"no spaces", "hello", "hello"},
		{"empty", "", ""},
		{"only spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Trim(tt.input)
			if result != tt.expected {
				t.Errorf("Trim(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"with duplicates", []string{"a", "b", "a", "c"}, []string{"a", "b", "c"}},
		{"no duplicates", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{"with empty strings", []string{"a", "", "b", ""}, []string{"a", "b"}},
		{"empty slice", []string{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Unique(%v) length = %d, want %d", tt.input, len(result), len(tt.expected))
			}
		})
	}
}

// Additional tests to improve coverage

func TestUnitCommonUserDataHashSumNonBase64(t *testing.T) {
	// Test non-base64 input
	input := "plain text"
	result := userDataHashSum(input)
	if result != input {
		t.Errorf("Expected %s, got %s", input, result)
	}
}

func TestUnitCommonConvertJsonStringToObjectNil(t *testing.T) {
	result := convertJsonStringToObject(`{"key":"value"}`)
	if result == nil {
		t.Error("Expected non-nil result")
	}
	if result["key"] != "value" {
		t.Errorf("Expected value, got %v", result["key"])
	}

	// Test invalid JSON returns nil
	result = convertJsonStringToObject(`{invalid}`)
	if result != nil {
		t.Error("Expected nil for invalid JSON")
	}
}

func TestUnitCommonConvertObjectToJsonStringError(t *testing.T) {
	// Test error case with channel (cannot be marshaled)
	input := make(chan int)
	result := convertObjectToJsonString(input)
	if result != "" {
		t.Errorf("Expected empty string for channel, got %s", result)
	}

	// Test valid case
	validInput := map[string]string{"key": "value"}
	result = convertObjectToJsonString(validInput)
	if !strings.Contains(result, "key") {
		t.Errorf("Expected JSON string, got %s", result)
	}
}

func TestUnitCommonConvertMaptoJsonStringError(t *testing.T) {
	// Test error case with channel
	input := map[string]interface{}{"ch": make(chan int)}
	_, err := convertMaptoJsonString(input)
	if err == nil {
		t.Error("Expected error for channel in map")
	}

	// Test valid case
	validInput := map[string]interface{}{"key": "value"}
	result, err := convertMaptoJsonString(validInput)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(result, "key") {
		t.Errorf("Expected JSON string, got %s", result)
	}
}

func TestUnitCommonConvertMapToJsonStringIgnoreErrorWithChannel(t *testing.T) {
	// Test with channel that causes error
	input := map[string]interface{}{"ch": make(chan int)}
	result := convertMapToJsonStringIgnoreError(input)
	if result != "" {
		t.Errorf("Expected empty string for error case, got %s", result)
	}

	// Test valid case
	validInput := map[string]interface{}{"key": "value"}
	result = convertMapToJsonStringIgnoreError(validInput)
	if !strings.Contains(result, "key") {
		t.Errorf("Expected JSON string, got %s", result)
	}
}

func TestUnitCommonConvertInterfaceToJsonStringError(t *testing.T) {
	// Test error case with channel
	input := make(chan int)
	_, err := convertInterfaceToJsonString(input)
	if err == nil {
		t.Error("Expected error for channel")
	}

	// Test valid case
	validInput := struct{ Key string }{"value"}
	result, err := convertInterfaceToJsonString(validInput)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(result, "Key") {
		t.Errorf("Expected JSON string, got %s", result)
	}
}

func TestUnitCommonWriteToFileWithHomeDir(t *testing.T) {
	// Test with ~ in path
	tmpDir := t.TempDir()
	// We can't directly test ~ expansion without mocking, but we can test the logic
	err := writeToFile(fmt.Sprintf("%s/test.txt", tmpDir), "content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test with nil data
	err = writeToFile(fmt.Sprintf("%s/test2.txt", tmpDir), nil)
	if err != nil {
		t.Errorf("Unexpected error for nil: %v", err)
	}

	// Test overwriting existing file
	err = writeToFile(fmt.Sprintf("%s/test.txt", tmpDir), "new content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestUnitCommonComputePeriodByUnitEdgeCases(t *testing.T) {
	now := time.Now()

	// Test with Year unit
	result, err := computePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(1, 0, 0).Format(time.RFC3339),
		0,
		"Year",
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1 year, got %d", result)
	}

	// Test with nil createTime
	_, err = computePeriodByUnit(nil, now.Format(time.RFC3339), 0, "Month")
	if err == nil {
		t.Error("Expected error for nil createTime")
	}

	// Test with nil endTime
	_, err = computePeriodByUnit(now.Format(time.RFC3339), nil, 0, "Month")
	if err == nil {
		t.Error("Expected error for nil endTime")
	}

	// Test with int64 time
	_, err = computePeriodByUnit(
		now.Unix(),
		now.AddDate(0, 1, 0).Unix(),
		0,
		"Month",
	)
	if err != nil {
		t.Errorf("Unexpected error with int64 time: %v", err)
	}

	// Test with invalid time format
	_, err = computePeriodByUnit("invalid", "invalid", 0, "Month")
	if err == nil {
		t.Error("Expected error for invalid time format")
	}

	// Test with unsupported period unit
	_, err = computePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(0, 1, 0).Format(time.RFC3339),
		0,
		"Invalid",
	)
	if err == nil {
		t.Error("Expected error for unsupported period unit")
	}

	// Test with currentPeriod set
	result, err = computePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(0, 2, 0).Format(time.RFC3339),
		5,
		"Month",
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5 (currentPeriod), got %d", result)
	}

	// Test UnStandardRFC3339 format
	UnStandardRFC3339 := "2006-01-02T15:04Z07:00"
	create := now.Format(UnStandardRFC3339)
	end := now.AddDate(0, 1, 0).Format(UnStandardRFC3339)
	_, err = computePeriodByUnit(create, end, 0, "Month")
	if err != nil {
		t.Errorf("Unexpected error with UnStandardRFC3339: %v", err)
	}

	// Test period > 12
	result, err = computePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(2, 0, 0).Format(time.RFC3339),
		0,
		"Month",
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 12 {
		t.Errorf("Expected 12 (max period), got %d", result)
	}
}

func TestUnitCommonCompressIPv6OrCIDREdgeCases(t *testing.T) {
	// Test empty string
	result, err := compressIPv6OrCIDR("")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}

	// Test invalid CIDR
	_, err = compressIPv6OrCIDR("invalid/cidr")
	if err == nil {
		t.Error("Expected error for invalid CIDR")
	}

	// Test IPv4
	result, err = compressIPv6OrCIDR("192.168.1.1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "192.168.1.1" {
		t.Errorf("Expected 192.168.1.1, got %s", result)
	}

	// Test invalid IP
	result, err = compressIPv6OrCIDR("not-an-ip")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "not-an-ip" {
		t.Errorf("Expected original string, got %s", result)
	}
}

func TestUnitCommonNormalizeValueEdgeCases(t *testing.T) {
	// Test string with JSON array that's invalid
	result, err := normalizeValue(`[invalid`)
	if err != nil {
		// Should return original value on error
		if result == nil {
			t.Error("Expected non-nil result on error")
		}
	}

	// Test nested map
	nestedMap := map[string]interface{}{
		"nested": map[string]interface{}{
			"key": "value",
		},
	}
	result, err = normalizeValue(nestedMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}

	// Test string true/false
	result, _ = normalizeValue("true")
	if result != true {
		t.Errorf("Expected true, got %v", result)
	}

	result, _ = normalizeValue("false")
	if result != false {
		t.Errorf("Expected false, got %v", result)
	}

	// Test numeric string
	result, _ = normalizeValue("123")
	if result != 123 {
		t.Errorf("Expected 123, got %v", result)
	}

	// Test array with nested values
	arr := []interface{}{"true", "false", "123", map[string]interface{}{"key": "value"}}
	result, err = normalizeValue(arr)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}

	// Test JSON array string
	result, err = normalizeValue(`["a", "b", "c"]`)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}
}

func TestUnitCommonConvertArrayObjectToJsonStringError(t *testing.T) {
	// Test with valid input
	input := []interface{}{1, 2, 3}
	result, err := convertArrayObjectToJsonString(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(result, "[1,2,3]") {
		t.Errorf("Expected [1,2,3], got %s", result)
	}

	// Test with error case (channels cannot be marshaled)
	input = []interface{}{make(chan int)}
	_, err = convertArrayObjectToJsonString(input)
	if err == nil {
		t.Error("Expected error for channel in array")
	}
}

func TestUnitCommonMapMergeNestedArrays(t *testing.T) {
	// Test nested array merging with maps
	target := map[string]interface{}{
		"arr": []interface{}{
			map[string]interface{}{"x": 1},
		},
	}
	merged := map[string]interface{}{
		"arr": []interface{}{
			map[string]interface{}{"y": 2},
		},
	}
	result := mapMerge(target, merged)
	arr := result["arr"].([]interface{})
	mergedMap := arr[0].(map[string]interface{})
	if mergedMap["x"] != 1 || mergedMap["y"] != 2 {
		t.Errorf("Expected merged map with x=1 and y=2, got %v", mergedMap)
	}

	// Test with non-map array element
	target2 := map[string]interface{}{
		"arr": []interface{}{1, 2},
	}
	merged2 := map[string]interface{}{
		"arr": []interface{}{3, 4},
	}
	result2 := mapMerge(target2, merged2)
	arr2 := result2["arr"].([]interface{})
	if arr2[0] != 3 || arr2[1] != 4 {
		t.Errorf("Expected [3, 4], got %v", arr2)
	}
}

func TestUnitCommonNewInstanceDiffComplex(t *testing.T) {
	// Test with simple map attribute changes
	state := &terraform.InstanceState{
		ID: "vpc-123",
		Attributes: map[string]string{
			"cidr_block": "192.168.0.0/16",
			"name":       "old-name",
		},
	}

	attributes := map[string]interface{}{
		"cidr_block": "192.168.0.0/16",
		"name":       "old-name",
	}

	attributesDiff := map[string]interface{}{
		"name": "new-name",
	}

	diff, err := newInstanceDiff("alicloud_vpc", attributes, attributesDiff, state)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if diff == nil {
		t.Fatal("Expected non-nil diff")
	}

	// Verify the diff contains the change
	if diff.Attributes["name"].Old != "old-name" || diff.Attributes["name"].New != "new-name" {
		t.Errorf("Expected name change from old-name to new-name")
	}
}

func TestUnitCommonConvertToJsonWithoutEscapeHTMLError(t *testing.T) {
	// Test with channel that causes error
	input := map[string]interface{}{"ch": make(chan int)}
	_, err := convertToJsonWithoutEscapeHTML(input)
	if err == nil {
		t.Error("Expected error for channel in map")
	}
}

func TestUnitCommonInterface2BoolNil(t *testing.T) {
	result := Interface2Bool(nil)
	if result != false {
		t.Errorf("Expected false for nil, got %v", result)
	}
}

func TestUnitCommonSplitMultiZoneIdNoSymbol(t *testing.T) {
	// Test without MAZ symbol and without parentheses
	result := splitMultiZoneId("simple-zone")
	if result != nil {
		t.Errorf("Expected nil for simple zone, got %v", result)
	}

	// Test with MAZ but without proper parentheses format
	// When MAZ is present but no "(", the function still tries to split
	// This tests the edge case where secondIndex is -1
	result = splitMultiZoneId("zone-MAZ-something")
	if result == nil {
		t.Error("Expected non-nil result when MAZ symbol is present")
	}
}

func TestUnitCommonWriteToFileEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with map data
	mapData := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	err := writeToFile(fmt.Sprintf("%s/map.json", tmpDir), mapData)
	if err != nil {
		t.Errorf("Unexpected error writing map: %v", err)
	}

	// Test with struct data
	structData := struct {
		Name string
		Age  int
	}{
		Name: "test",
		Age:  30,
	}
	err = writeToFile(fmt.Sprintf("%s/struct.json", tmpDir), structData)
	if err != nil {
		t.Errorf("Unexpected error writing struct: %v", err)
	}

	// Verify file was created and can be read
	content, err := os.ReadFile(fmt.Sprintf("%s/struct.json", tmpDir))
	if err != nil {
		t.Errorf("Failed to read written file: %v", err)
	}
	if len(content) == 0 {
		t.Error("Written file is empty")
	}
}

func TestUnitCommonNewInstanceDiffRemoveAttributes(t *testing.T) {
	// Test removing attributes (setting to empty)
	state := &terraform.InstanceState{
		ID: "vpc-123",
		Attributes: map[string]string{
			"description": "old description",
			"name":        "test-vpc",
		},
	}

	attributes := map[string]interface{}{
		"description": "old description",
		"name":        "test-vpc",
	}

	attributesDiff := map[string]interface{}{
		"description": "",
		"name":        "test-vpc",
	}

	diff, err := newInstanceDiff("alicloud_vpc", attributes, attributesDiff, state)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if diff == nil {
		t.Fatal("Expected non-nil diff")
	}

	// Verify description was removed
	if descDiff, ok := diff.Attributes["description"]; ok {
		if descDiff.New != "" {
			t.Errorf("Expected empty description, got %s", descDiff.New)
		}
	}
}
