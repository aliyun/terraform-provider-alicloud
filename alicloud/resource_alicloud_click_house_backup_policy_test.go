package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudClickHouseBackupPolicy_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_backup_policy.default"
	checkoutSupportedRegions(t, true, connectivity.ClickHouseSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudClickHouseBackupPolicyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickhouseService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sclickhousebackuppolicy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseBackupPolicyBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Monday", "Tuesday"},
					"preferred_backup_time":   "00:00Z-01:00Z",
					"db_cluster_id":           "${alicloud_click_house_db_cluster.default.id}",
					"backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "00:00Z-01:00Z",
						"db_cluster_id":             CHECKSET,
						"backup_retention_period":   "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Wednesday", "Monday", "Saturday"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_time": "01:00Z-02:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_time": "01:00Z-02:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "14",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "14",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": []string{"Monday", "Tuesday"},
					"preferred_backup_time":   "00:00Z-01:00Z",
					"backup_retention_period": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period.#": "2",
						"preferred_backup_time":     "00:00Z-01:00Z",
						"backup_retention_period":   "7",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudClickHouseBackupPolicyMap0 = map[string]string{
	"db_cluster_id": CHECKSET,
	"status":        CHECKSET,
}

func AlicloudClickHouseBackupPolicyBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_click_house_regions" "default" {	
  current = true
}

data "alicloud_vpcs" "default"	{
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id = "${data.alicloud_vpcs.default.ids.0}"
  zone_id = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}

resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  status                  = "Running"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = "1"
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}
`, name)
}
