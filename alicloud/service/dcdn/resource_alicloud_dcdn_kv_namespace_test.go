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
func TestAccAliCloudDcdnKvNamespace_basic2278(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dcdn_kv_namespace.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDcdnKvNamespaceMap2278)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dcdn.NewDcdnService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDcdnKvNamespace")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccDcdnKvNamespace%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDcdnKvNamespaceBasicDependence2278)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
