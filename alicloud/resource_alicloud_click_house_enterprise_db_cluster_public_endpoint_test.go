package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ClickHouse EnterpriseDbClusterPublicEndpoint. >>> Resource test cases, automatically generated.
// Case 企业版CK开公网1-线上 10561
func TestAccAliCloudClickHouseEnterpriseDbClusterPublicEndpoint_basic10561(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_click_house_enterprise_db_cluster_public_endpoint.default"
	ra := resourceAttrInit(resourceId, AlicloudClickHouseEnterpriseDbClusterPublicEndpointMap10561)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ClickHouseServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeClickHouseEnterpriseDbClusterPublicEndpoint")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccclickhouse%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudClickHouseEnterpriseDbClusterPublicEndpointBasicDependence10561)
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
					"db_instance_id":           "${alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id}",
					"net_type":                 "Public",
					"connection_string_prefix": "${alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_id":           CHECKSET,
						"net_type":                 "Public",
						"connection_string_prefix": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "${alicloud_click_house_enterprise_db_cluster.defaultaqnt22.id}8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": CHECKSET,
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

var AlicloudClickHouseEnterpriseDbClusterPublicEndpointMap10561 = map[string]string{}

func AlicloudClickHouseEnterpriseDbClusterPublicEndpointBasicDependence10561(name string) string {
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

resource "alicloud_click_house_enterprise_db_cluster" "defaultaqnt22" {
  zone_id    = var.zone_id_i
  vpc_id     = alicloud_vpc.defaultktKLuM.id
  scale_min  = "8"
  scale_max  = "16"
  vswitch_id = alicloud_vswitch.defaultTQWN3k.id
}


`, name)
}

// Test ClickHouse EnterpriseDbClusterPublicEndpoint. <<< Resource test cases, automatically generated.
