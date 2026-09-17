package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudPolarDBGatewayModelService_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_model_service.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-model-service", func(string) string { return "" })
	resource.Test(t, resource.TestCase{Providers: testAccProviders, CheckDestroy: func(*terraform.State) error { return nil }, Steps: []resource.TestStep{
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "pg-test",
				"name":                           "service-one",
				"model_category":                 "text",
				"protocol":                       "openai",
				"base_url":                       "https://example.com/openai/v1",
				"api_key":                        "placeholder-one",
				"vendor":                         "custom",
				"input_cost_points_per_million":  "10",
				"output_cost_points_per_million": "20",
				"request_cost_points":            "1",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
		{
			Config: testAccConfig(map[string]interface{}{
				"gateway_id":                     "pg-test",
				"name":                           "service-two",
				"model_category":                 "embedding",
				"protocol":                       "anthropic",
				"base_url":                       "https://example.com/anthropic/v1",
				"api_key":                        "placeholder-two",
				"vendor":                         "custom",
				"input_cost_points_per_million":  "11",
				"output_cost_points_per_million": "21",
				"request_cost_points":            "2",
			}),
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		},
	}})
}
