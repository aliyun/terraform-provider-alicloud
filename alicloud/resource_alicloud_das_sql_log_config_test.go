package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAliCloudDasSqlLogConfig_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	resourceId := "alicloud_das_sql_log_config.default"
	ra := resourceAttrInit(resourceId, AliCloudDasSqlLogConfigMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DasServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDasSqlLogConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccDasSqlLogConfig-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDasSqlLogConfigBasicDependence0)
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
					"instance_id":    "${alicloud_polardb_cluster.default.id}",
					"enable":         "true",
					"request_enable": "true",
					"retention":      "30",
					"hot_retention":  "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":    CHECKSET,
						"enable":         "true",
						"request_enable": "true",
						"retention":      "30",
						"hot_retention":  "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention":     "180",
					"hot_retention": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention":     "180",
						"hot_retention": "7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"retention":     "180",
					"hot_retention": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"retention":     "180",
						"hot_retention": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable", "request_enable"},
			},
		},
	})
}

var AliCloudDasSqlLogConfigMap0 = map[string]string{
	"instance_id": CHECKSET,
}

func AliCloudDasSqlLogConfigBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_polardb_zones" "default" {
	}

	data "alicloud_polardb_node_classes" "default" {
  		db_type    = "MySQL"
  		db_version = "8.0"
  		pay_type   = "PostPaid"
  		zone_id    = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids) - 1]
  		category   = "Normal"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		zone_id      = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids) - 1]
  		vpc_id       = alicloud_vpc.default.id
  		vswitch_name = var.name
  		cidr_block   = "172.16.0.0/24"
	}

	resource "alicloud_polardb_cluster" "default" {
  		db_type       = "MySQL"
  		db_version    = "8.0"
  		pay_type      = "PostPaid"
  		db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  		vswitch_id    = alicloud_vswitch.default.id
  		description   = var.name
	}
`, name)
}
