package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
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
