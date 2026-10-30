package dcdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dcdn"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudDcdnKv_basic2277(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_kv.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDcdnKvMap2277)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnKv")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDcdnKv%d", defaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDcdnKvBasicDependence2277)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
