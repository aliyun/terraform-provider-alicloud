package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCREndpointAclPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cr_endpoint_acl_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCREndpointAclPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CrService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCrEndpointAclPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scrEndpointaclpolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCREndpointAclPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithTime(t, []int{1})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   name,
					"entry":         "192.168.1.0/24",
					"instance_id":   "${local.instance_id}",
					"module_name":   "Registry",
					"endpoint_type": "internet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   name,
						"entry":         "192.168.1.0/24",
						"instance_id":   CHECKSET,
						"module_name":   "Registry",
						"endpoint_type": "internet",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"module_name"},
			},
		},
	})
}

var AlicloudCREndpointAclPolicyMap0 = map[string]string{}

func AlicloudCREndpointAclPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
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
`, name)
}
