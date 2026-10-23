package tag_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The CreatedBy switch is an account-scoped singleton with no configurable arguments: the
// resource exists only while the switch is enabled. The case therefore needs no random
// cloud-side name. Create adopts an already-enabled switch (OpenCreatedBy reports
// CreatedByAlreadyOpen), and destroy verifies the switch is off both through the
// resource's own post-CloseCreatedBy readback and through an independent
// CheckCreatedByEnabled call in CheckDestroy, so a passing destroy is evidence that the
// switch was actually turned off rather than just dropped from state. The case leaves the
// switch disabled, which is the only cleanup possible for an account-wide singleton.
//
// Tag serves the three CreatedBy operations only from cn-shanghai, so the resource ID is
// that control region regardless of ALICLOUD_REGION; the checks below call the same
// control endpoint directly.
const testAccAliCloudTagCreatedByID = "cn-shanghai"

func TestAccAliCloudTagCreatedBy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy: func(s *terraform.State) error {
			client := acctest.Provider.Meta().(*connectivity.AliyunClient)
			endpoint := connectivity.LoadRegionalEndpoint(testAccAliCloudTagCreatedByID, "tag")
			for _, rs := range s.RootModule().Resources {
				if rs.Type != "alicloud_tag_created_by" {
					continue
				}
				response, err := client.RpcPostWithEndpoint("Tag", "2018-08-28", "CheckCreatedByEnabled", nil, map[string]interface{}{
					"RegionId": testAccAliCloudTagCreatedByID,
				}, true, endpoint)
				if err != nil {
					return fmt.Errorf("checking alicloud_tag_created_by %s after destroy: %w", rs.Primary.ID, err)
				}
				if enabled, ok := response["OpenStatus"].(bool); ok && enabled {
					return fmt.Errorf("alicloud_tag_created_by %s is still enabled after destroy", rs.Primary.ID)
				}
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudTagCreatedByConfig,
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["alicloud_tag_created_by.default"]
						if !ok {
							return fmt.Errorf("alicloud_tag_created_by.default is absent from state")
						}
						if rs.Primary.ID != testAccAliCloudTagCreatedByID {
							return fmt.Errorf("id = %q, want the CreatedBy control region %q", rs.Primary.ID, testAccAliCloudTagCreatedByID)
						}
						return nil
					},
					resource.TestCheckResourceAttr("alicloud_tag_created_by.default", "open_status", "true"),
				),
			},
			{
				Config: testAccAliCloudTagCreatedByConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_tag_created_by.default", "open_status", "true"),
				),
			},
			{
				ResourceName:      "alicloud_tag_created_by.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Changing only timeouts is the one update this resource accepts. It must keep the
			// switch enabled and leave no diff afterwards: the Update path has to write the
			// live open_status back, not the unknown value carried by the plan.
			{
				Config: testAccAliCloudTagCreatedByConfigWithUpdateTimeout,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_tag_created_by.default", "open_status", "true"),
					resource.TestCheckResourceAttr("alicloud_tag_created_by.default", "timeouts.update", "10m"),
				),
			},
			{
				Config: testAccAliCloudTagCreatedByConfigWithUpdateTimeout,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_tag_created_by.default", "open_status", "true"),
					resource.TestCheckResourceAttr("alicloud_tag_created_by.default", "timeouts.update", "10m"),
				),
			},
		},
	})
}

const testAccAliCloudTagCreatedByConfig = `
resource "alicloud_tag_created_by" "default" {}
`

const testAccAliCloudTagCreatedByConfigWithUpdateTimeout = `
resource "alicloud_tag_created_by" "default" {
  timeouts {
    update = "10m"
  }
}
`
