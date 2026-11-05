package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for the tag helpers moved from the root package
// (alicloud/tags.go); test names are unchanged from
// alicloud/tags_test.go.

func TestTagsMapEqual(t *testing.T) {
	lmap := map[string]interface{}{
		"Created": "TF6666",
		"For":     "acceptance test",
	}
	rmap := map[string]string{
		"Created": "TF6666",
		"For":     "acceptance test",
	}
	if !TagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is not equal.")
	}

	rmap = map[string]string{
		"Created": "TF6666",
	}
	if TagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}

	rmap = map[string]string{
		"Created": "TF_Fake",
		"For":     "acceptance test",
	}
	if TagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}

	rmap = map[string]string{
		"Fake": "TF6666",
		"For":  "acceptance test",
	}
	if TagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}

	lmap = map[string]interface{}{
		"Created": 6666,
		"For":     "acceptance test",
	}
	if TagsMapEqual(lmap, rmap) {
		t.Fatal("Tag maps is equal.")
	}
}

// Tag converters moved from alicloud/common.go; test names are unchanged
// from alicloud/common_test.go.

func TestUnitMergeTags(t *testing.T) {
	// A resource tag overrides the matching provider default of the same key.
	defaults := map[string]interface{}{"env": "prod", "team": "a"}
	resource := map[string]interface{}{"env": "dev", "role": "web"}
	merged := MergeTags(defaults, resource)
	assert.Equal(t, "dev", merged["env"], "resource tag overrides default")
	assert.Equal(t, "a", merged["team"], "default-only key carried through")
	assert.Equal(t, "web", merged["role"], "resource-only key carried through")
	assert.Equal(t, 3, len(merged))

	// With no default_tags the merge is the resource tags themselves, so a
	// provider that does not configure default_tags behaves exactly as before.
	assert.Equal(t, resource, MergeTags(nil, resource))
	assert.Equal(t, 2, len(MergeTags(nil, resource)))
}

func TestUnitStripDefaultTags(t *testing.T) {
	defaults := map[string]interface{}{"env": "prod", "team": "a"}
	cloud := map[string]string{
		"env":  "prod", // matches the default exactly -> stripped
		"team": "b",    // overrides the default value -> kept
		"role": "web",  // resource-only -> kept
	}
	stripped := StripDefaultTags(cloud, defaults)
	_, hasEnv := stripped["env"]
	assert.False(t, hasEnv, "an entry matching a default is stripped from user tags")
	assert.Equal(t, "b", stripped["team"], "an overridden default value is kept")
	assert.Equal(t, "web", stripped["role"], "a resource-only tag is kept")
	assert.Equal(t, 2, len(stripped))

	// With no default_tags the cloud set is returned unchanged.
	assert.Equal(t, len(cloud), len(StripDefaultTags(cloud, nil)))
}

func TestUnitTagsMapsEqual(t *testing.T) {
	a := map[string]interface{}{"env": "prod", "team": "a"}
	b := map[string]interface{}{"team": "a", "env": "prod"}
	assert.True(t, tagsMapsEqual(a, b), "same pairs in a different order are equal")
	assert.False(t, tagsMapsEqual(a, map[string]interface{}{"env": "prod"}), "different length is not equal")
	assert.False(t, tagsMapsEqual(a, map[string]interface{}{"env": "dev", "team": "a"}), "a different value is not equal")
	assert.True(t, tagsMapsEqual(map[string]interface{}{}, map[string]interface{}{}))
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

	result := ExpandTagsToMap(originMap, tags)

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

	result := ExpandTagsToMapWithTags(originMap, tags)

	if result["Tags.1.Key"] != "env" {
		t.Error("Tags.1.Key not set correctly")
	}
	if result["Tags.1.Value"] != "prod" {
		t.Error("Tags.1.Value not set correctly")
	}
}

func TestUnitCommonConvertLowercaseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected []map[string]interface{}
	}{
		{
			"empty",
			map[string]interface{}{},
			[]map[string]interface{}{},
		},
		{
			"single",
			map[string]interface{}{
				"Created": "Terraform",
			},
			[]map[string]interface{}{
				{
					"key":   "Created",
					"value": "Terraform",
				},
			},
		},
		{
			"multiple",
			map[string]interface{}{
				"Created": "Terraform",
				"For":     "Provider",
				"Updated": "TF",
			},
			[]map[string]interface{}{
				{
					"key":   "Created",
					"value": "Terraform",
				},
				{
					"key":   "For",
					"value": "Provider",
				},
				{
					"key":   "Updated",
					"value": "TF",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertLowercaseTags(tt.input)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
