package dbfs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dbfs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Case 1
func TestAccAliCloudDbfsAutoSnapShotPolicy_basic2601(t *testing.T) {
	var v map[string]interface{}
	t.Skip("This resource has been deprecated since v1.279.0 and will be removed in the future. See: https://help.aliyun.com/en/dbfs/")
	acctest.CheckoutSupportedRegions(t, true, connectivity.DBFSSystemSupportRegions)
	resourceId := "alicloud_dbfs_auto_snap_shot_policy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudDbfsAutoSnapShotPolicyMap2601)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return dbfs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDbfsAutoSnapShotPolicy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sDbfsAutoSnapShotPolicy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudDbfsAutoSnapShotPolicyBasicDependence2601)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
