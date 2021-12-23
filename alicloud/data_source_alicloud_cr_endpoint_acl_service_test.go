package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCrEndpointAclServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_endpoint_acl_service.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCrEndpointAclServiceDataSource,
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

const testAccCheckAlicloudCrEndpointAclServiceDataSource = `
variable "name" {
  default = "tf-testacc-CrEndpointAclService"
}
data "alicloud_cr_ee_instances" "default" {}
resource "alicloud_cr_ee_instance" "default" {
  count          = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? 0 : 1
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}
locals {
  instance_id = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? data.alicloud_cr_ee_instances.default.ids[0] : concat(alicloud_cr_ee_instance.default.*.id, [""])[0]
}
data "alicloud_cr_endpoint_acl_service" "default" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = local.instance_id
  module_name   = "Registry"
}
`
