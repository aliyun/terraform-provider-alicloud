package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ClickHouse EnterpriseDbClusterBackupPolicy. >>> Resource test cases, automatically generated.
// Case CK企业版备份策略3-线上 10564
func TestAccAliCloudClickHouseEnterpriseDbClusterBackupPolicy_basic10564(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster_backup_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterBackupPolicyMap10564)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbClusterBackupPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterBackupPolicyBasicDependence10564)
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
					"preferred_backup_period": "Monday",
					"preferred_backup_time":   "04:00Z-05:00Z",
					"backup_retention_period": "7",
					"db_instance_id":          "${alicloud_click_house_enterprise_db_cluster.default1tTLwe.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period": "Monday",
						"preferred_backup_time":   "04:00Z-05:00Z",
						"backup_retention_period": "7",
						"db_instance_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"preferred_backup_period": "Tuesday",
					"preferred_backup_time":   "05:00Z-06:00Z",
					"backup_retention_period": "14",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"preferred_backup_period": "Tuesday",
						"preferred_backup_time":   "05:00Z-06:00Z",
						"backup_retention_period": "14",
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

var AlicloudClickHouseEnterpriseDbClusterBackupPolicyMap10564 = map[string]string{}

func AlicloudClickHouseEnterpriseDbClusterBackupPolicyBasicDependence10564(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vsw_ip_range_i" {
  default = "172.16.1.0/24"
}

variable "vpc_ip_range" {
  default = "172.16.0.0/12"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

resource "alicloud_vpc" "defaultktKLuM" {
  cidr_block = var.vpc_ip_range
}

resource "alicloud_vswitch" "defaultTQWN3k" {
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw_ip_range_i
}

resource "alicloud_click_house_enterprise_db_cluster" "default1tTLwe" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


`, name)
}

// Test ClickHouse EnterpriseDbClusterBackupPolicy. <<< Resource test cases, automatically generated.
