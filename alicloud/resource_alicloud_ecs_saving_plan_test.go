package alicloud

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccProtoV5ProviderFactories returns the muxed provider (SDK v2 + Framework)
// behind one protocol v5 server. It mirrors alicloud/acctest.ProtoV5ProviderFactories
// but is defined locally to avoid the import cycle (package alicloud → alicloud/acctest
// → alicloud) that would arise from importing alicloud/acctest from an internal test.
func testAccProtoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"alicloud": func() (tfprotov5.ProviderServer, error) {
			serverFactory, err := provider.ProtoV5ProviderServerFactory(context.Background(), Provider())
			if err != nil {
				return nil, err
			}
			return serverFactory(), nil
		},
	}
}

// testAccCheckEcsSavingPlanDestroy is a no-op. Saving Plans are financial
// commitments; the PurchaseSavingPlanOffering API cannot be reversed and no
// cancellation API exists. The Framework resource's Delete therefore removes
// only the Terraform state entry — it does not call any purchase or cancel
// API and does not clean up the cloud resource. The cloud Saving Plan remains
// and must be cancelled via the console if needed. This check does NOT verify
// cloud resource cleanup; it satisfies the TestAcc CheckDestroy contract
// required by tflint AT001 without claiming a deletion that cannot occur.
func testAccCheckEcsSavingPlanDestroy(s *terraform.State) error {
	return nil
}

// Saving Plan creation is a financial commitment. These tests are skipped
// unless ALICLOUD_SAVING_PLAN_TEST is set, matching the SDKv2 gate.
func TestAccEcsSavingPlan_basic(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set. Saving Plan creation involves financial commitment.")
	}

	resourceName := "alicloud_ecs_saving_plan.default"
	r := acctest.RandIntRange(1000, 9999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcsSavingPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsSavingPlanConfigBasic(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "committed_amount", "0.08"),
					resource.TestCheckResourceAttr(resourceName, "plan_type", "EcsCompute"),
					resource.TestCheckResourceAttr(resourceName, "offering_type", "AllUpfront"),
					resource.TestCheckResourceAttr(resourceName, "purchase_method", "ByInstanceFamily"),
					resource.TestCheckResourceAttr(resourceName, "charge_type", "PrePaid"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "Year"),
					resource.TestCheckResourceAttr(resourceName, "instance_family", "ecs.g5"),
					resource.TestCheckResourceAttr(resourceName, "saving_plan_name", fmt.Sprintf("tf-test-%d-name", r)),
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

func TestAccEcsSavingPlan_general(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set.")
	}

	resourceName := "alicloud_ecs_saving_plan.default"
	r := acctest.RandIntRange(1000, 9999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckEcsSavingPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsSavingPlanConfigGeneral(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "plan_type", "General"),
					resource.TestCheckResourceAttr(resourceName, "offering_type", "HalfUpfront"),
					resource.TestCheckResourceAttr(resourceName, "purchase_method", "ByInstanceFamilySet"),
				),
			},
		},
	})
}

func testAccEcsSavingPlanConfigBasic(r int) string {
	return fmt.Sprintf(`
resource "alicloud_ecs_saving_plan" "default" {
  committed_amount = "0.08"
  plan_type        = "EcsCompute"
  offering_type    = "AllUpfront"
  purchase_method  = "ByInstanceFamily"
  charge_type      = "PrePaid"
  period           = 1
  period_unit      = "Year"
  instance_family  = "ecs.g5"
  saving_plan_name = "tf-test-%[1]d-name"
  description      = "test saving plan"
}
`, r)
}

func testAccEcsSavingPlanConfigGeneral(r int) string {
	return fmt.Sprintf(`
resource "alicloud_ecs_saving_plan" "default" {
  committed_amount    = "0.08"
  plan_type           = "General"
  offering_type       = "HalfUpfront"
  purchase_method     = "ByInstanceFamilySet"
  instance_family_set = "ecs.g5"
  saving_plan_name    = "tf-test-%[1]d-general"
  description         = "general saving plan"
  start_time          = "2026-01-01 00:00:00"
}
`, r)
}
