package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudDcdnKv_basic2277(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_kv.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnKvMap2277)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnKv")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDcdnKv%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnKvBasicDependence2277)
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
					"value":     "${var.name}",
					"key":       "${var.name}",
					"namespace": "${alicloud_dcdn_kv_namespace.default.namespace}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value":     name,
						"key":       name,
						"namespace": name,
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"value": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"value": name + "_update",
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

var AlicloudDcdnKvMap2277 = map[string]string{}

func AlicloudDcdnKvBasicDependence2277(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_dcdn_kv_account" "default" {
  status = "online"
}

resource "alicloud_dcdn_kv_namespace" "default" {
  description = "wkmtest"
  namespace   = var.name
}
`, name)
}
