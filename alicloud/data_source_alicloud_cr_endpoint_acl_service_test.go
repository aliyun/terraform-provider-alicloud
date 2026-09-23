package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCREndpointAclServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_endpoint_acl_service.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	// CR EE instance names are validated server-side and reject uppercase
	// letters; keep the name lowercase and hyphen-separated.
	name := fmt.Sprintf("tf-testacc-cr-aclsvc-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCrEndpointAclServiceDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":            CHECKSET,
						"enable":        "true",
						"status":        "RUNNING",
						"instance_id":   CHECKSET,
						"module_name":   CHECKSET,
						"endpoint_type": "internet",
					}),
				),
			},
		},
	})
}

// The acl service data source must target a dedicated instance created by the
// test: picking an arbitrary account instance (data.alicloud_cr_ee_instances
// ids.0) can select one whose status does not support endpoint operations and
// UpdateInstanceEndpointStatus then fails with INSTANCE_STATUS_NOT_SUPPORT.
func testAccCheckAlicloudCrEndpointAclServiceDataSource(name string) string {
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

data "alicloud_cr_endpoint_acl_service" "default" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = alicloud_cr_ee_instance.default.id
  module_name   = "Registry"
}
`, name)
}
