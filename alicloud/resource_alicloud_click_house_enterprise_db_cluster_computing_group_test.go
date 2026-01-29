// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ClickHouse EnterpriseDbClusterComputingGroup. >>> Resource test cases, automatically generated.
// Case 线上-企业版CK-多计算组 12417
func TestAccAliCloudClickHouseEnterpriseDbClusterComputingGroup_basic12417(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster_computing_group.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterComputingGroupMap12417)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbClusterComputingGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterComputingGroupBasicDependence12417)
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
					"node_scale_min":              "4",
					"computing_group_description": "test",
					"node_count":                  "2",
					"db_instance_id":              "${alicloud_click_house_enterprise_db_cluster.defaultQ5vukB.id}",
					"node_scale_max":              "4",
					"is_readonly":                 "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_scale_min":              "4",
						"computing_group_description": "test",
						"node_count":                  "2",
						"db_instance_id":              CHECKSET,
						"node_scale_max":              "4",
						"is_readonly":                 "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"computing_group_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"computing_group_description": "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"is_readonly": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"is_readonly": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"node_scale_min": "8",
					"node_count":     "3",
					"node_scale_max": "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"node_scale_min": "8",
						"node_count":     "3",
						"node_scale_max": "8",
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

var AlicloudClickHouseEnterpriseDbClusterComputingGroupMap12417 = map[string]string{
	"computing_group_endpoints.#":        CHECKSET,
	"computing_group_public_endpoints.#": CHECKSET,
	"computing_group_endpoint_names.#":   CHECKSET,
	"computing_group_id":                 CHECKSET,
	"computing_group_status":             CHECKSET,
}

func AlicloudClickHouseEnterpriseDbClusterComputingGroupBasicDependence12417(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vsw__ip_range_i" {
  default = "172.16.1.0/24"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "vpc__ip_range" {
  default = "172.16.0.0/12"
}

variable "vsw__ip_range_k" {
  default = "172.16.3.0/24"
}

variable "vsw__ip_range_l" {
  default = "172.16.2.0/24"
}

variable "zone_id_i" {
  default = "cn-beijing-i"
}

variable "zone_id_l" {
  default = "cn-beijing-l"
}

variable "zone_id_k" {
  default = "cn-beijing-k"
}

resource "alicloud_vpc" "defaultp2mwWM" {
  cidr_block = var.vpc__ip_range
}

resource "alicloud_vswitch" "defaultkCZhNu" {
  vpc_id     = alicloud_vpc.defaultp2mwWM.id
  zone_id    = var.zone_id_i
  cidr_block = var.vsw__ip_range_i
}

resource "alicloud_click_house_enterprise_db_cluster" "defaultQ5vukB" {
  zone_id    = alicloud_vswitch.defaultkCZhNu.zone_id
  vpc_id     = alicloud_vpc.defaultp2mwWM.id
  scale_min  = "8"
  scale_max  = "8"
  vswitch_id = alicloud_vswitch.defaultkCZhNu.id
}


`, name)
}

// Test ClickHouse EnterpriseDbClusterComputingGroup. <<< Resource test cases, automatically generated.
