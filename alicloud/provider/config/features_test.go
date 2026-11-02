package config

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/features"
)

func TestUnitProviderFeatures(t *testing.T) {
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
