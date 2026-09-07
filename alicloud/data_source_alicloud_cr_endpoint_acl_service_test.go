package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCREndpointAclServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_endpoint_acl_service.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
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
data "alicloud_cr_endpoint_acl_service" "default" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = data.alicloud_cr_ee_instances.default.ids.0
  module_name   = "Registry"
}
`
