package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ClickHouse EnterpriseDbClusterSecurityIP. >>> Resource test cases, automatically generated.
// Case CK企业版白名单-线上 10562
func TestAccAliCloudClickHouseEnterpriseDbClusterSecurityIP_basic10562(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster_security_ip.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterSecurityIPMap10562)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbClusterSecurityIP")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterSecurityIPBasicDependence10562)
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
					"group_name":       "test_group",
					"security_ip_list": "127.0.0.2",
					"db_instance_id":   "${alicloud_click_house_enterprise_db_cluster.defaultn0nVrN.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":       "test_group",
						"security_ip_list": "127.0.0.2",
						"db_instance_id":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": "127.0.0.3,127.0.0.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list": "127.0.0.3,127.0.0.2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": "127.0.0.2,127.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list": CHECKSET,
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

var AlicloudClickHouseEnterpriseDbClusterSecurityIPMap10562 = map[string]string{}

func AlicloudClickHouseEnterpriseDbClusterSecurityIPBasicDependence10562(name string) string {
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

resource "alicloud_click_house_enterprise_db_cluster" "defaultn0nVrN" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


`, name)
}

// Test ClickHouse EnterpriseDbClusterSecurityIP. <<< Resource test cases, automatically generated.
