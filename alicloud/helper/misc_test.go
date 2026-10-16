package helper

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// Tests for the misc helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

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
			result, err := CompressIPv6OrCIDR(tt.input)
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
			result := IsProtocolValid(tc.input)
			if result != tc.expected {
				t.Errorf("input: %s, except: %v, result: %v", tc.input, tc.expected, result)
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
			result := SplitMultiZoneId(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitMultiZoneId() = %v, want %v", result, tt.expected)
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
			result, err := GetOneStringOrAllStringSlice(tt.input, tt.name)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
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
			result := GetPagination(tt.pageNumber, tt.pageSize)
			if result.PageNumber != tt.pageNumber || result.PageSize != tt.pageSize {
				t.Errorf("GetPagination failed")
			}
		})
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
			result, err := GetNextpageNumber(tt.input)
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
			result := TerraformToAPI(tt.input)
			if result != tt.expected {
				t.Errorf("TerraformToAPI(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonLoadFileContent(t *testing.T) {
	tmpFile := fmt.Sprintf("%s/test_load.txt", t.TempDir())
	content := "test content"
	os.WriteFile(tmpFile, []byte(content), 0644)

	result, err := LoadFileContent(tmpFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(result) != content {
		t.Errorf("Expected %s, got %s", content, string(result))
	}

	// Test non-existent file
	_, err = LoadFileContent("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
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
			result := MapMerge(targetCopy, tt.merged)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
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
			result := BytesToTB(tt.input)
			if result != tt.expected {
				t.Errorf("BytesToTB(%d) = %f, want %f", tt.input, result, tt.expected)
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

			result := IsPagingRequest(d)
			if result != tt.expected {
				t.Errorf("IsPagingRequest() = %v, want %v", result, tt.expected)
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
			SetPagingRequest(d, request, tt.maxPageSize)

			if request["PageNumber"] != tt.expectedPageNum {
				t.Errorf("PageNumber = %v, want %v", request["PageNumber"], tt.expectedPageNum)
			}
			if request["PageSize"] != tt.expectedPageSize {
				t.Errorf("PageSize = %v, want %v", request["PageSize"], tt.expectedPageSize)
			}
		})
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

func TestUnitCommonCompressIPv6OrCIDREdgeCases(t *testing.T) {
	// Test empty string
	result, err := CompressIPv6OrCIDR("")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string, got %s", result)
	}

	// Test invalid CIDR
	_, err = CompressIPv6OrCIDR("invalid/cidr")
	if err == nil {
		t.Error("Expected error for invalid CIDR")
	}

	// Test IPv4
	result, err = CompressIPv6OrCIDR("192.168.1.1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "192.168.1.1" {
		t.Errorf("Expected 192.168.1.1, got %s", result)
	}

	// Test invalid IP
	result, err = CompressIPv6OrCIDR("not-an-ip")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != "not-an-ip" {
		t.Errorf("Expected original string, got %s", result)
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
	result := MapMerge(target, merged)
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
	result2 := MapMerge(target2, merged2)
	arr2 := result2["arr"].([]interface{})
	if arr2[0] != 3 || arr2[1] != 4 {
		t.Errorf("Expected [3, 4], got %v", arr2)
	}
}

func TestUnitCommonSplitMultiZoneIdNoSymbol(t *testing.T) {
	// Test without MAZ symbol and without parentheses
	result := SplitMultiZoneId("simple-zone")
	if result != nil {
		t.Errorf("Expected nil for simple zone, got %v", result)
	}

	// Test with MAZ but without proper parentheses format
	// When MAZ is present but no "(", the function still tries to split
	// This tests the edge case where secondIndex is -1
	result = SplitMultiZoneId("zone-MAZ-something")
	if result == nil {
		t.Error("Expected non-nil result when MAZ symbol is present")
	}
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
			result := UserDataHashSum(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnitCommonUserDataHashSumNonBase64(t *testing.T) {
	// Test non-base64 input
	input := "plain text"
	result := UserDataHashSum(input)
	if result != input {
		t.Errorf("Expected %s, got %s", input, result)
	}
}
