package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProviderSkipRegionValidation(t *testing.T) {
	testCases := []struct {
		name     string
		env      string
		raw      map[string]interface{}
		expected bool
	}{
		{"unset env defaults to false", "", map[string]interface{}{}, false},
		{"env true enables skip", "true", map[string]interface{}{}, true},
		{"env false disables skip", "false", map[string]interface{}{}, false},
		{"config overrides env", "true", map[string]interface{}{"skip_region_validation": false}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv("ALICLOUD_SKIP_REGION_VALIDATION", testCase.env)
			d := schema.TestResourceDataRaw(t, Provider().Schema, testCase.raw)
			if got := d.Get("skip_region_validation").(bool); got != testCase.expected {
				t.Fatalf("skip_region_validation: expected %v, got %v", testCase.expected, got)
			}
		})
	}
}

func TestProviderFeatures(t *testing.T) {
	testCases := []struct {
		name     string
		features []interface{}
		expected features.Features
	}{
		{
			name:     "no features block at all",
			features: []interface{}{},
			expected: features.Default(),
		},
		{
			name:     "a features block that is null",
			features: []interface{}{nil},
			expected: features.Default(),
		},
		{
			name:     "a features block without the ecs_instance block",
			features: []interface{}{map[string]interface{}{}},
			expected: features.Default(),
		},
		{
			name:     "an empty ecs_instance block",
			features: []interface{}{map[string]interface{}{"ecs_instance": []interface{}{}}},
			expected: features.Default(),
		},
		{
			name: "replace_on_image_update enabled",
			features: []interface{}{map[string]interface{}{
				"ecs_instance": []interface{}{map[string]interface{}{"replace_on_image_update": true}},
			}},
			expected: features.Features{
				EcsInstance: features.EcsInstance{ReplaceOnImageUpdate: true},
			},
		},
		{
			name: "replace_on_image_update disabled explicitly",
			features: []interface{}{map[string]interface{}{
				"ecs_instance": []interface{}{map[string]interface{}{"replace_on_image_update": false}},
			}},
			expected: features.Default(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := expandFeatures(testCase.features); got != testCase.expected {
				t.Fatalf("expandFeatures: expected %+v, got %+v", testCase.expected, got)
			}
		})
	}
}
