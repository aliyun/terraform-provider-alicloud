package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudEcsSavingPlan_basic(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set. Saving Plan creation involves financial commitment.")
	}

	resourceName := "alicloud_ecs_saving_plan.default"
	ra := fmt.Sprintf("tf-test-%d", acctest.RandIntRange(1000, 9999))
	testAccConfig := resourceTestAccConfigFunc(resourceName, ra, func(name string) string { return "" })

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudEcsSavingPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"committed_amount": "0.08",
					"plan_type":        "EcsCompute",
					"offering_type":    "AllUpfront",
					"purchase_method":  "ByInstanceFamily",
					"charge_type":      "PrePaid",
					"period":           1,
					"period_unit":      "Year",
					"instance_family":  "ecs.g5",
					"saving_plan_name": fmt.Sprintf("%s-name", ra),
					"description":      "test saving plan",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "committed_amount", "0.08"),
					resource.TestCheckResourceAttr(resourceName, "plan_type", "EcsCompute"),
					resource.TestCheckResourceAttr(resourceName, "offering_type", "AllUpfront"),
					resource.TestCheckResourceAttr(resourceName, "purchase_method", "ByInstanceFamily"),
					resource.TestCheckResourceAttr(resourceName, "charge_type", "PrePaid"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "Year"),
					resource.TestCheckResourceAttr(resourceName, "instance_family", "ecs.g5"),
					resource.TestCheckResourceAttr(resourceName, "saving_plan_name", fmt.Sprintf("%s-name", ra)),
					resource.TestCheckResourceAttr(resourceName, "description", "test saving plan"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"charge_type",
					"description",
					"instance_family_set",
					"period_unit",
					"purchase_method",
					"saving_plan_name",
				},
			},
		},
	})
}

func TestAccAliCloudEcsSavingPlan_unitTest(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set.")
	}

	resourceName := "alicloud_ecs_saving_plan.default"
	ra := fmt.Sprintf("tf-test-%d", acctest.RandIntRange(1000, 9999))
	testAccConfig := resourceTestAccConfigFunc(resourceName, ra, func(name string) string { return "" })

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAlicloudEcsSavingPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"committed_amount":    "0.08",
					"plan_type":           "General",
					"offering_type":       "HalfUpfront",
					"purchase_method":     "ByInstanceFamilySet",
					"instance_family_set": "ecs.g5",
					"saving_plan_name":    fmt.Sprintf("%s-general", ra),
					"description":         "general saving plan",
					"start_time":          "2026-01-01 00:00:00",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "plan_type", "General"),
					resource.TestCheckResourceAttr(resourceName, "offering_type", "HalfUpfront"),
				),
			},
		},
	})
}

func testAccCheckAlicloudEcsSavingPlanDestroy(s *terraform.State) error {
	// Saving Plans cannot be deleted via API; they are only removed from Terraform state.
	// No cloud resource destruction to verify here.
	return nil
}
