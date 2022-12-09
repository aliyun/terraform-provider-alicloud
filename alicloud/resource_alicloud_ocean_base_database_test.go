package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAlicloudOceanBase_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ocean_base_database.default"
	checkoutSupportedRegions(t, true, connectivity.OceanBaseSupportRegions)
	ra := resourceAttrInit(resourceId, OceanBaseBasicMap)
	serviceFunc := func() interface{} {
		return &OceanBaseProService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tfacc%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOceanBaseDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"collation":     "utf8mb4_general_ci",
					"database_name": name,
					"description":   name,
					"encoding":      "utf8mb4",
					"instance_id":   "ob4a01r9i3er4w",
					"tenant_id":     "t4a03bchxk69s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"collation":     "utf8mb4_general_ci",
						"database_name": name,
						"description":   name,
						"encoding":      "utf8mb4",
						"tenant_id":     CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"users": `[{\"UserName\":\"omstest\",\"Role\":\"readwrite\"}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "users"},
			},
		},
	})

}

func resourceOceanBaseDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}

var OceanBaseBasicMap = map[string]string{}
