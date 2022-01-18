package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudNASLifecyclePolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_lifecycle_policy.default"
	checkoutSupportedRegions(t, true, connectivity.NASSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudNASLifecyclePolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasLifecyclePolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snaslifecyclepolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNASLifecyclePolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_policy_name": "${var.name}",
					"file_system_id":        "${alicloud_nas_file_system.default.id}",
					"lifecycle_rule_name":   "DEFAULT_ATIME_14",
					"storage_type":          "InfrequentAccess",
					"paths":                 []string{"/"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_policy_name": name,
						"file_system_id":        CHECKSET,
						"lifecycle_rule_name":   "DEFAULT_ATIME_14",
						"storage_type":          "InfrequentAccess",
						"paths.#":               "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lifecycle_rule_name": "DEFAULT_ATIME_30",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lifecycle_rule_name": "DEFAULT_ATIME_30",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudNASLifecyclePolicyMap0 = map[string]string{
	"file_system_id":        CHECKSET,
	"storage_type":          CHECKSET,
	"lifecycle_policy_name": CHECKSET,
	"paths.#":               CHECKSET,
}

func AlicloudNASLifecyclePolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

resource "alicloud_nas_file_system" "default" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}
`, name)
}
