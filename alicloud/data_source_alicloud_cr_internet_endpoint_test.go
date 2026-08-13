package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCrInternetEndpointDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	resourceId := "data.alicloud_cr_internet_endpoint.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCrInternetEndpointDataSourceConfigBasic0(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":          CHECKSET,
						"instance_id": CHECKSET,
						"status":      CHECKSET,
					}),
					resource.TestCheckResourceAttr(resourceId, "entries.#", "1"),
				),
			},
		},
	})
}

func testAccAlicloudCrInternetEndpointDataSourceConfigBasic0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Economy"
  instance_name  = var.name
  image_scanner  = "DISABLE"
}

locals {
  cr_endpoint_instance_id = alicloud_cr_ee_instance.default.id
}

resource "alicloud_cr_internet_endpoint" "default" {
  instance_id = local.cr_endpoint_instance_id

  entries {
    entry   = "192.168.1.0/24"
    comment = "entry-1"
  }
}

data "alicloud_cr_internet_endpoint" "default" {
  instance_id = alicloud_cr_internet_endpoint.default.id
}
`, name)
}
