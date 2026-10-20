package alicloud

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAliCloudPolarDBGatewayBudgetPolicy_basic(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-budget-policy", resourcePolarDBGatewayAIConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBGatewayBudgetPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":          "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":         "GlobalTotal",
					"reset_day_of_month":  1,
					"budget_points":       "1000",
					"alert_threshold_pct": 80,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_type", "GlobalTotal"),
					resource.TestCheckResourceAttr(resourceID, "budget_points", "1000"),
					resource.TestCheckResourceAttr(resourceID, "reset_day_of_month", "1"),
					resource.TestCheckResourceAttr(resourceID, "alert_threshold_pct", "80"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":          "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":         "GlobalTotal",
					"reset_day_of_month":  2,
					"budget_points":       "2000",
					"alert_threshold_pct": 90,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_points", "2000"),
					resource.TestCheckResourceAttr(resourceID, "reset_day_of_month", "2"),
					resource.TestCheckResourceAttr(resourceID, "alert_threshold_pct", "90"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudPolarDBGatewayBudgetPolicy_dimension(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-budget-policy-dimension", resourcePolarDBGatewayAIConsumerGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBGatewayBudgetPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":              "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":             "ConsumerGroupTotal",
					"budget_dimension_ref_id": "${alicloud_polardb_gateway_consumer_group.dependency.consumer_group_id}",
					"reset_day_of_month":      1,
					"budget_points":           "1000",
					"alert_threshold_pct":     80,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_type", "ConsumerGroupTotal"),
					resource.TestCheckResourceAttrPair(resourceID, "budget_dimension_ref_id", "alicloud_polardb_gateway_consumer_group.dependency", "consumer_group_id"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudPolarDBGatewayBudgetPolicy_consumer(t *testing.T) {
	resourceID := "alicloud_polardb_gateway_budget_policy.default"
	testAccConfig := resourceTestAccConfigFunc(resourceID, "tf-acc-budget-policy-consumer", resourcePolarDBGatewayAIConsumerDependence)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckPolarDBGatewayAI(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPolarDBGatewayBudgetPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"gateway_id":              "${alicloud_polardb_gateway.dependency.id}",
					"budget_type":             "ConsumerTotal",
					"budget_dimension_ref_id": "${alicloud_polardb_gateway_consumer.dependency.consumer_id}",
					"reset_day_of_month":      1,
					"budget_points":           "1000",
					"alert_threshold_pct":     80,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceID, "budget_type", "ConsumerTotal"),
					resource.TestCheckResourceAttrPair(resourceID, "budget_dimension_ref_id", "alicloud_polardb_gateway_consumer.dependency", "consumer_id"),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUnitPolarDBGatewayBudgetPolicyDiff(t *testing.T) {
	cases := []struct {
		name, budgetType, dimensionRef, wantError string
	}{
		{"global without ref", "GlobalTotal", "", ""},
		{"global with ref", "GlobalTotal", "consumer-1", "must not be set"},
		{"consumer without ref", "ConsumerTotal", "", "must be set"},
		{"group without ref", "ConsumerGroupTotal", "", "must be set"},
		{"consumer with ref", "ConsumerTotal", "consumer-1", ""},
		{"group with ref", "ConsumerGroupTotal", "group-1", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := map[string]interface{}{
				"gateway_id":         "gateway-1",
				"budget_type":        tc.budgetType,
				"reset_day_of_month": 1,
				"budget_points":      "1000",
			}
			if tc.dimensionRef != "" {
				config["budget_dimension_ref_id"] = tc.dimensionRef
			}
			_, err := resourceAlicloudPolarDBGatewayBudgetPolicy().Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil)
			if tc.wantError == "" && err != nil {
				t.Fatalf("unexpected diff error: %v", err)
			}
			if tc.wantError != "" && (err == nil || !strings.Contains(err.Error(), tc.wantError)) {
				t.Fatalf("diff error = %v; want substring %q", err, tc.wantError)
			}
		})
	}
}

func TestUnitPolarDBGatewayBudgetPolicyPointsRoundTrip(t *testing.T) {
	resourceSchema := resourceAlicloudPolarDBGatewayBudgetPolicy().Schema
	data := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{"budget_points": "1000.50"})
	want := resourceSchema["budget_points"].StateFunc("1000.50")
	if want != "1000.5" {
		t.Fatalf("config normalization = %q; want 1000.5", want)
	}
	if err := setPolarDBGatewayBudgetPolicyPoints(data, "1000.50"); err != nil {
		t.Fatal(err)
	}
	if got := data.Get("budget_points").(string); got != want {
		t.Fatalf("read normalization = %q; config normalization = %q", got, want)
	}
}

func TestUnitPolarDBGatewayBudgetPolicyConsumerTypeMapping(t *testing.T) {
	resourceSchema := resourceAlicloudPolarDBGatewayBudgetPolicy().Schema
	for _, tc := range []struct {
		name       string
		budgetType string
		wireType   string
	}{
		{name: "consumer", budgetType: "ConsumerTotal", wireType: "Consumer"},
		{name: "consumer group", budgetType: "ConsumerGroupTotal", wireType: "ConsumerGroupTotal"},
		{name: "global", budgetType: "GlobalTotal", wireType: "GlobalTotal"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
				"gateway_id":              "pg-test",
				"budget_type":             tc.budgetType,
				"budget_dimension_ref_id": "test-ref",
				"reset_day_of_month":      1,
				"budget_points":           "1000",
			})
			request := polarDBGatewayBudgetPolicyRequest(data, "cn-beijing")
			if got := request["BudgetType"]; got != tc.wireType {
				t.Fatalf("BudgetType = %v; want %q", got, tc.wireType)
			}
			if got := flattenPolarDBGatewayBudgetType(tc.wireType); got != tc.budgetType {
				t.Fatalf("flattened BudgetType = %q; want %q", got, tc.budgetType)
			}
		})
	}
}
