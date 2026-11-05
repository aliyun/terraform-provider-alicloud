package rds_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRdsDBProxy_MySQL(t *testing.T) {
	var connection map[string]interface{}
	var primary map[string]interface{}
	var readonly map[string]interface{}

	resourceId := "alicloud_rds_db_proxy.default"
	var DBProxyMap = map[string]string{
		"instance_id":           acctest.CHECKSET,
		"instance_network_type": "VPC",
		"db_proxy_instance_num": "2",
		"vswitch_id":            acctest.CHECKSET,
		"vpc_id":                acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, DBProxyMap)

	rc_connection := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &connection, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeRdsProxyEndpoint")
	rc_primary := acctest.ResourceCheckInitWithDescribeMethod("alicloud_db_instance.default", &primary, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBInstance")
	rc_readonly := acctest.ResourceCheckInitWithDescribeMethod("alicloud_db_readonly_instance.default", &readonly, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBReadonlyInstance")
	rand := tfacctest.RandString(5)
	effectiveSpecificTime := time.Now().UTC().Add(45 * time.Minute).Format("2006-01-02T15:04:05Z")
	switchTime := time.Now().UTC().Add(45 * time.Minute).Format("2006-01-02T15:04:05Z")
	rac := acctest.ResourceAttrCheckInit(rc_connection, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBProxy%s", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceRdsDBProxyDependence_MySQL)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":            "${alicloud_db_instance.default.id}",
					"instance_network_type":  "VPC",
					"db_proxy_instance_num":  "2",
					"db_proxy_instance_type": "common",
					"vswitch_id":             "${alicloud_db_instance.default.vswitch_id}",
					"vpc_id":                 "${alicloud_db_instance.default.vpc_id}",
					"depends_on":             []string{"alicloud_db_readonly_instance.default"},
					"resource_group_id":      "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
				}),
				Check: resource.ComposeTestCheckFunc(
					rc_primary.CheckResourceExists(),
					rc_readonly.CheckResourceExists(),
					testAccCheck(map[string]string{
						"instance_id":            acctest.CHECKSET,
						"db_proxy_instance_type": acctest.CHECKSET,
						"instance_network_type":  "VPC",
						"db_proxy_instance_num":  "2",
						"vswitch_id":             acctest.CHECKSET,
						"vpc_id":                 acctest.CHECKSET,
						"resource_group_id":      acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"upgrade_time": "Immediate",
					"switch_time":  switchTime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"upgrade_time": "Immediate",
						"switch_time":  acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_num":   "2",
					"db_proxy_instance_type":  "exclusive",
					"effective_time":          "SpecificTime",
					"effective_specific_time": effectiveSpecificTime,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":   "2",
						"db_proxy_instance_type":  "exclusive",
						"effective_time":          "SpecificTime",
						"effective_specific_time": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_endpoint_read_write_mode": "ReadWrite",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_endpoint_read_write_mode": "ReadWrite",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_instance_distribution_type": "Standard",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_instance_distribution_type": "Standard",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_instance_max_delay_time": "90",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_instance_max_delay_time": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_features": "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_features": "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"read_only_instance_distribution_type": "Custom",
					"db_proxy_features":                    "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:1",
					"read_only_instance_weight": []map[string]interface{}{
						{
							"instance_id": "${alicloud_db_instance.default.id}",
							"weight":      "100",
						},
						{
							"instance_id": "${alicloud_db_readonly_instance.default.id}",
							"weight":      "500",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"read_only_instance_distribution_type": "Custom",
						"db_proxy_features":                    "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:1",
						"read_only_instance_weight.#":          "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_instance_num":  "3",
					"effective_time":         "Immediate",
					"db_proxy_instance_type": "exclusive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_instance_num":  "3",
						"effective_time":         "Immediate",
						"db_proxy_instance_type": "exclusive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_connect_string_port": "3306",
					"db_proxy_connection_prefix":   fmt.Sprintf("tf-testacc%s", rand),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_connect_string_port": "3306",
						"db_proxy_connection_prefix":   acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_ssl_enabled": "Open",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_ssl_enabled": "Open",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"switch_time", "effective_specific_time"},
			},
		},
	})
}

func resourceRdsDBProxyDependence_MySQL(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_db_zones" "default"{
	engine = "MySQL"
	engine_version = "8.0"
	instance_charge_type = "PostPaid"
	category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
    zone_id = data.alicloud_db_zones.default.zones.0.id
	engine = "MySQL"
	engine_version = "8.0"
    category = "HighAvailability"
 	db_instance_storage_type = "cloud_essd"
	instance_charge_type = "PostPaid"
}

data "alicloud_vpcs" "default" {
    name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
 count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
 vswitch_name = var.name
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_db_zones.default.ids.0
 cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_db_zones.default.ids[length(data.alicloud_db_zones.default.ids)-1]
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = "mysql.x2.large.2c"
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = local.vswitch_id
	instance_name = var.name
}

resource "alicloud_db_readonly_instance" "default" {
	master_db_instance_id = alicloud_db_instance.default.id
	zone_id = alicloud_db_instance.default.zone_id
	engine_version = alicloud_db_instance.default.engine_version
	instance_type = "mysqlro.n4.medium.1c"
	instance_storage = alicloud_db_instance.default.instance_storage
	instance_name = "${var.name}_ro"
	vswitch_id = alicloud_db_instance.default.vswitch_id
}

`, name)
}
