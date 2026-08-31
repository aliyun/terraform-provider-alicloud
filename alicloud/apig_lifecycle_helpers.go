package alicloud

import (
	"fmt"
	"sort"
	"strings"
)

func apigResponseData(response map[string]interface{}) (map[string]interface{}, error) {
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("APIG response does not contain an object in data")
	}
	return data, nil
}

func apigObjectSlice(value interface{}) []map[string]interface{} {
	values, ok := value.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(values))
	for _, value := range values {
		if object, ok := value.(map[string]interface{}); ok {
			result = append(result, object)
		}
	}
	return result
}

func apigStringSlice(value interface{}) []string {
	values, ok := value.([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if text, ok := value.(string); ok && text != "" {
			result = append(result, text)
		}
	}
	return result
}

func apigSortedStringSet(value interface{}) []string {
	values := convertToInterfaceArray(value)
	result := make([]string, 0, len(values))
	for _, value := range values {
		if text, ok := value.(string); ok && text != "" {
			result = append(result, text)
		}
	}
	sort.Strings(result)
	return result
}

func apigParseCompositeID(id string, parts int) ([]string, error) {
	values := strings.Split(id, ":")
	if len(values) != parts {
		return nil, fmt.Errorf("invalid import ID %q: expected %d colon-separated parts", id, parts)
	}
	for _, value := range values {
		if value == "" {
			return nil, fmt.Errorf("invalid import ID %q: parts must not be empty", id)
		}
	}
	return values, nil
}
