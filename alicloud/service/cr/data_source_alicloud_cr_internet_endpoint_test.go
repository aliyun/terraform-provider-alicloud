package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCrInternetEndpointDataSource_basic0(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccr%d", rand)
	resourceId := "data.alicloud_cr_internet_endpoint.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Providers: acctest.SDKProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudCrInternetEndpointDataSourceConfigBasic0(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":          acctest.CHECKSET,
						"instance_id": acctest.CHECKSET,
						"status":      acctest.CHECKSET,
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
