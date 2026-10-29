package ga_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudGaDomain_basic2188(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ga_domain.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGaDomainMap2188)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaDomain")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaDomain%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGaDomainBasicDependence2188)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain":         "changes.com.cn",
					"accelerator_id": "${local.accelerator_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":         "changes.com.cn",
						"accelerator_id": acctest.CHECKSET,
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
