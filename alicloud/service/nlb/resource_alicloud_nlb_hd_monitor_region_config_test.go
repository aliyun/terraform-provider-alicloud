package nlb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/nlb"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Nlb HdMonitorRegionConfig. >>> Resource test cases, automatically generated.
// Case hdMonitor 9477
func TestAccAliCloudNlbHdMonitorRegionConfig_basic9477(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nlb_hd_monitor_region_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudNlbHdMonitorRegionConfigMap9477)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return nlb.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeNlbHdMonitorRegionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccnlb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudNlbHdMonitorRegionConfigBasicDependence9477)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_store": "test",
					"log_project":  "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_store": "test",
						"log_project":  "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_store": "test2",
					"log_project":  "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_store": "test2",
						"log_project":  "test2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudNlbHdMonitorRegionConfigMap9477 = map[string]string{
	"region_id": acctest.CHECKSET,
}

func AlicloudNlbHdMonitorRegionConfigBasicDependence9477(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Nlb HdMonitorRegionConfig. <<< Resource test cases, automatically generated.
