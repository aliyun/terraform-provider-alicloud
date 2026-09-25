package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEcsSavingPlans_basic(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set.")
	}

	resourceName := "data.alicloud_ecs_saving_plans.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudEcsSavingPlansConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID(resourceName),
				),
			},
		},
	})
}

func TestAccAliCloudEcsSavingPlans_unitTest(t *testing.T) {
	if v := os.Getenv("ALICLOUD_SAVING_PLAN_TEST"); v == "" {
		t.Skip("Skipping test because ALICLOUD_SAVING_PLAN_TEST is not set.")
	}

	resourceName := "data.alicloud_ecs_saving_plans.default"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudEcsSavingPlansConfigWithInstanceId(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID(resourceName),
				),
			},
		},
	})
}

func testAccAlicloudEcsSavingPlansConfigBasic() string {
	return fmt.Sprintf(`
data "alicloud_ecs_saving_plans" "default" {
  output_file = "saving_plans.txt"
}
`)
}

func testAccAlicloudEcsSavingPlansConfigWithInstanceId() string {
	return fmt.Sprintf(`
data "alicloud_ecs_saving_plans" "default" {
  instance_id = "spn-dummy"
}
`)
}
