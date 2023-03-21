package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudDbfsAutoSnapShotPolicy_basic2601(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.DBFSSystemSupportRegions)
	resourceId := "alicloud_dbfs_auto_snap_shot_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudDbfsAutoSnapShotPolicyMap2601)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DbfsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDbfsAutoSnapShotPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDbfsAutoSnapShotPolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDbfsAutoSnapShotPolicyBasicDependence2601)
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
					"time_points":     []string{"01"},
					"policy_name":     "${var.name}",
					"retention_days":  "1",
					"repeat_weekdays": []string{"2"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"time_points.#":     "1",
						"policy_name":       name,
						"retention_days":    "1",
						"repeat_weekdays.#": "1",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"policy_name":     "${var.name}_update",
					"time_points":     []string{"05"},
					"retention_days":  "2",
					"repeat_weekdays": []string{"1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name":       name + "_update",
						"time_points.#":     "1",
						"retention_days":    "2",
						"repeat_weekdays.#": "1",
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

var AlicloudDbfsAutoSnapShotPolicyMap2601 = map[string]string{}

func AlicloudDbfsAutoSnapShotPolicyBasicDependence2601(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
