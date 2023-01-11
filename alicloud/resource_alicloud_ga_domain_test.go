package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudGaDomain_basic2188(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudGaDomainMap2188)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GaService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGaDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaDomain%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGaDomainBasicDependence2188)
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
					"domain":         "changes.com.cn",
					"accelerator_id": "${local.accelerator_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":         "changes.com.cn",
						"accelerator_id": CHECKSET,
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudGaDomainMap2188 = map[string]string{}

func AlicloudGaDomainBasicDependence2188(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_ga_accelerators" "default" {
  status = "active"
}

resource "alicloud_ga_accelerator" "default" {
  count           = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? 0 : 1
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

locals {
  accelerator_id = length(data.alicloud_ga_accelerators.default.accelerators) > 0 ? data.alicloud_ga_accelerators.default.accelerators.0.id : alicloud_ga_accelerator.default.0.id
}

`, name)
}
