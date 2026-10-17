package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEcsSavingPlans_basic(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set.")
	}

	resourceName := "data.alicloud_ecs_saving_plans.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEcsSavingPlansConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccEcsSavingPlans_withInstanceId(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set.")
	}

	resourceName := "data.alicloud_ecs_saving_plans.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEcsSavingPlansConfigWithInstanceId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccEcsSavingPlansConfigBasic() string {
	return `
data "alicloud_ecs_saving_plans" "default" {
  output_file = "saving_plans.txt"
}
`
}

func testAccEcsSavingPlansConfigWithInstanceId() string {
	return fmt.Sprintf(`
data "alicloud_ecs_saving_plans" "default" {
  instance_id = "spn-dummy"
}
`)
}
