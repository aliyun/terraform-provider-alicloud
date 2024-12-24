package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAliCloudDcdnKvNamespace_basic2278(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_kv_namespace.default"
	ra := resourceAttrInit(resourceId, AlicloudDcdnKvNamespaceMap2278)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DcdnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDcdnKvNamespace")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccDcdnKvNamespace%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDcdnKvNamespaceBasicDependence2278)
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
					"description": "wkmtest",
					"namespace":   "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "wkmtest",
						"namespace":   name,
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

var AlicloudDcdnKvNamespaceMap2278 = map[string]string{}

func AlicloudDcdnKvNamespaceBasicDependence2278(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_dcdn_kv_account" "current" {
  status = "online"
}

`, name)
}
