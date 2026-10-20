package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Tests for the conversion helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

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
			result := ExpandStringList(tt.input)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestUnitCommonConvertListStringToListInterface(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := []interface{}{"a", "b", "c"}
	result := ConvertListStringToListInterface(input)
	assert.Equal(t, expected, result)
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
			result := ConvertPaymentTypeToChargeType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
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
			result := ExpandIntList(tc.input)
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
			result := ConvertListToJsonString(tc.input)
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
			encoded := EncodeToBase64String(tc.input)
			decoded, err := DecodeFromBase64String(encoded)
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
			result, err := ConvertJsonStringToMap(tc.input)

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
			result := ConvertListToCommaSeparate(tc.input)
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
			result := FilterEmptyStrings(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("input: %v, except: %v, result: %v", tc.input, tc.expected, result)
			}
		})
	}
}

func TestUnitCommonTypeConversionFunctions(t *testing.T) {
	t.Run("BoolToString", func(t *testing.T) {
		if ConvertBoolToString(true) != "true" || ConvertBoolToString(false) != "false" {
			t.Error("BoolToString failed")
		}
	})

	t.Run("StringToBool", func(t *testing.T) {
		if !ConvertStringToBool("true") || ConvertStringToBool("false") {
			t.Error("StringToBool failed")
		}
	})

	t.Run("IntToString", func(t *testing.T) {
		if ConvertIntergerToString(42) != "42" {
			t.Error("IntToString failed")
		}
	})

	t.Run("FloatToString", func(t *testing.T) {
		result := ConvertFloat64ToString(3.14)
		if !strings.Contains(result, "3.14") {
			t.Errorf("ConvertFloat64ToString failed: %s", result)
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
			result, err := ConvertJsonStringToList(tc.input)

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

	result := ExpandArrayToMap(inputMap, array, key)

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
		result := ConvertJsonStringToObject(input)
		expected := map[string]interface{}{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("input: %q, except: %v, result: %v", input, expected, result)
		}
	})

	t.Run("ObjectToJSONString", func(t *testing.T) {
		input := map[string]interface{}{"key": "value"}
		result := ConvertObjectToJsonString(input)
		if result != `{"key":"value"}` {
			t.Errorf("ObjectToJSONString error: %s", result)
		}
	})

	t.Run("MapToJSONString", func(t *testing.T) {
		input := map[string]interface{}{"key": "value"}
		result, err := ConvertMaptoJsonString(input)
		if err != nil || result != `{"key":"value"}` {
			t.Errorf("err: %s, 错误: %v", result, err)
		}
	})

	t.Run("MapToJSONIgnoreError", func(t *testing.T) {
		input := map[string]interface{}{"key": make(chan int)}
		result := ConvertMapToJsonStringIgnoreError(input)
		if result != "" {
			t.Error("err:", result)
		}
	})

	t.Run("InterfaceToJSON", func(t *testing.T) {
		input := struct{ Key string }{"value"}
		result, err := ConvertInterfaceToJsonString(input)
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
		wantErr  bool
	}{
		{
			name:     "ValidIntegerSlice",
			input:    []interface{}{1, 2, 3},
			expected: []interface{}{"1", "2", "3"},
			wantErr:  false,
		},
		{
			name:     "ValidMixedSlice",
			input:    []interface{}{"10", 5.5, 1},
			expected: []interface{}{"10", "5", "1"},
			wantErr:  false,
		},
		{
			name:     "EmptySlice",
			input:    []interface{}{},
			expected: []interface{}{},
			wantErr:  false,
		},
		{
			name:     "InvalidType",
			input:    "not a slice",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "NilInput",
			input:    nil,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "JsonNumber",
			input:    []interface{}{json.Number("10")},
			expected: []interface{}{"10"},
			wantErr:  false,
		},
		{
			name:     "IllegalElement",
			input:    []interface{}{1, "abc"},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertJsonStringToStringList(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Error expectation failed: wantErr %v, got %v", tt.wantErr, err)
			}

			if tt.wantErr {
				if result != nil {
					t.Errorf("Expected nil result on error, got %v", result)
				}
				return
			}

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
			result, err := ConvertListMapToJsonString(tt.input)

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
			result, err := ConvertArrayObjectToJsonString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertArrayObjectToJsonString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ConvertArrayObjectToJsonString() = %v, want %v", result, tt.expected)
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
			result := ConvertArrayToString(tt.input, tt.sep)
			if result != tt.expected {
				t.Errorf("ConvertArrayToString() = %v, want %v", result, tt.expected)
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
			result, err := ConvertMapFloat64ToJsonString(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertMapFloat64ToJsonString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result != tt.expected {
				t.Errorf("ConvertMapFloat64ToJsonString() = %v, want %v", result, tt.expected)
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
			result, err := ConvertToJsonWithoutEscapeHTML(tt.input)

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
			result := ConvertToInterfaceArray(tt.input)
			if tt.name == "schema.Set with strings" || tt.name == "schema.Set with ints" {
				// For schema.Set, we need to compare elements regardless of order
				if len(result) != len(tt.expected) {
					t.Errorf("ConvertToInterfaceArray(%v) = %v (len %d), want %v (len %d)", tt.input, result, len(result), tt.expected, len(tt.expected))
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
					t.Errorf("ConvertToInterfaceArray(%v) = %v, want %v (elements mismatch)", tt.input, result, tt.expected)
				}
			} else {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("ConvertToInterfaceArray(%v) = %v, want %v", tt.input, result, tt.expected)
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
			result, err := ConvertYamlToObject(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, result)
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
			result := ExpandSingletonToList(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ExpandSingletonToList(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
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
			result := ConvertChargeTypeToPaymentType(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonConvertJsonStringToObjectNil(t *testing.T) {
	result := ConvertJsonStringToObject(`{"key":"value"}`)
	if result == nil {
		t.Error("Expected non-nil result")
	}
	if result["key"] != "value" {
		t.Errorf("Expected value, got %v", result["key"])
	}

	// Test invalid JSON returns nil
	result = ConvertJsonStringToObject(`{invalid}`)
	if result != nil {
		t.Error("Expected nil for invalid JSON")
	}
}

func TestUnitCommonConvertObjectToJsonStringError(t *testing.T) {
	// Test error case with channel (cannot be marshaled)
	input := make(chan int)
	result := ConvertObjectToJsonString(input)
	if result != "" {
		t.Errorf("Expected empty string for channel, got %s", result)
	}

	// Test valid case
	validInput := map[string]string{"key": "value"}
	result = ConvertObjectToJsonString(validInput)
	if !strings.Contains(result, "key") {
		t.Errorf("Expected JSON string, got %s", result)
	}
}

func TestUnitCommonConvertMaptoJsonStringError(t *testing.T) {
	// Test error case with channel
	input := map[string]interface{}{"ch": make(chan int)}
	_, err := ConvertMaptoJsonString(input)
	if err == nil {
		t.Error("Expected error for channel in map")
	}

	// Test valid case
	validInput := map[string]interface{}{"key": "value"}
	result, err := ConvertMaptoJsonString(validInput)
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
	result := ConvertMapToJsonStringIgnoreError(input)
	if result != "" {
		t.Errorf("Expected empty string for error case, got %s", result)
	}

	// Test valid case
	validInput := map[string]interface{}{"key": "value"}
	result = ConvertMapToJsonStringIgnoreError(validInput)
	if !strings.Contains(result, "key") {
		t.Errorf("Expected JSON string, got %s", result)
	}
}

func TestUnitCommonConvertInterfaceToJsonStringError(t *testing.T) {
	// Test error case with channel
	input := make(chan int)
	_, err := ConvertInterfaceToJsonString(input)
	if err == nil {
		t.Error("Expected error for channel")
	}

	// Test valid case
	validInput := struct{ Key string }{"value"}
	result, err := ConvertInterfaceToJsonString(validInput)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(result, "Key") {
		t.Errorf("Expected JSON string, got %s", result)
	}
}

func TestUnitCommonConvertArrayObjectToJsonStringError(t *testing.T) {
	// Test with valid input
	input := []interface{}{1, 2, 3}
	result, err := ConvertArrayObjectToJsonString(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.Contains(result, "[1,2,3]") {
		t.Errorf("Expected [1,2,3], got %s", result)
	}

	// Test with error case (channels cannot be marshaled)
	input = []interface{}{make(chan int)}
	_, err = ConvertArrayObjectToJsonString(input)
	if err == nil {
		t.Error("Expected error for channel in array")
	}
}

func TestUnitCommonConvertToJsonWithoutEscapeHTMLError(t *testing.T) {
	// Test with channel that causes error
	input := map[string]interface{}{"ch": make(chan int)}
	_, err := ConvertToJsonWithoutEscapeHTML(input)
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

func TestUnitCommonNormalizeValueEdgeCases(t *testing.T) {
	// Test string with JSON array that's invalid
	result, err := NormalizeValue(`[invalid`)
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
	result, err = NormalizeValue(nestedMap)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}

	// Test string true/false
	result, _ = NormalizeValue("true")
	if result != true {
		t.Errorf("Expected true, got %v", result)
	}

	result, _ = NormalizeValue("false")
	if result != false {
		t.Errorf("Expected false, got %v", result)
	}

	// Test numeric string
	result, _ = NormalizeValue("123")
	if result != 123 {
		t.Errorf("Expected 123, got %v", result)
	}

	// Test array with nested values
	arr := []interface{}{"true", "false", "123", map[string]interface{}{"key": "value"}}
	result, err = NormalizeValue(arr)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}

	// Test JSON array string
	result, err = NormalizeValue(`["a", "b", "c"]`)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == nil {
		t.Error("Expected non-nil result")
	}
}
