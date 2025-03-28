package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rds CustomDisk. >>> Resource test cases, automatically generated.
// Case CustomDisk_test1 10599
func TestAccAliCloudRdsCustomDisk_basic10599(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_disk.default"
	ra := resourceAttrInit(resourceId, AlicloudRdsCustomDiskMap10599)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRdsCustomDiskBasicDependence10599)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":          "zcc测试用例",
					"zone_id":              "cn-beijing-i",
					"size":                 "40",
					"performance_level":    "PL1",
					"instance_charge_type": "Postpaid",
					"disk_category":        "cloud_essd",
					"disk_name":            "custom_disk_001",
					"auto_renew":           "false",
					"period":               "1",
					"auto_pay":             "true",
					"period_unit":          "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "zcc测试用例",
						"zone_id":              "cn-beijing-i",
						"size":                 "40",
						"performance_level":    "PL1",
						"instance_charge_type": "Postpaid",
						"disk_category":        "cloud_essd",
						"disk_name":            "custom_disk_001",
						"auto_renew":           "false",
						"period":               "1",
						"auto_pay":             "true",
						"period_unit":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size":    "50",
					"type":    "offline",
					"dry_run": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size":    "50",
						"type":    "offline",
						"dry_run": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "disk_category", "dry_run", "instance_charge_type", "period", "period_unit", "snapshot_id", "type"},
			},
		},
	})
}

var AlicloudRdsCustomDiskMap10599 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"region_id":   CHECKSET,
}

func AlicloudRdsCustomDiskBasicDependence10599(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-beijing"
}


`, name)
}

// Test Rds CustomDisk. <<< Resource test cases, automatically generated.
