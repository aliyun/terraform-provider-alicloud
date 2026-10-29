package rds_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rds"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRdsDBProxy_Public_MySQL(t *testing.T) {
	var instance map[string]interface{}
	var proxy map[string]interface{}
	var proxypublic map[string]interface{}

	resourceId := "alicloud_rds_db_proxy_public.default"
	var DBProxyMap = map[string]string{
		"db_instance_id":                      acctest.CHECKSET,
		"db_proxy_endpoint_id":                acctest.CHECKSET,
		"connection_string_prefix":            acctest.CHECKSET,
		"db_proxy_connection_string_net_type": "Public",
		"db_proxy_new_connect_string_port":    acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, DBProxyMap)
	rc_proxypublic := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &proxypublic, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBProxy")
	rc_instance := acctest.ResourceCheckInitWithDescribeMethod("alicloud_db_instance.default", &instance, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBInstance")
	rc_proxy := acctest.ResourceCheckInitWithDescribeMethod("alicloud_rds_db_proxy.default", &proxy, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBProxy")

	rand := tfacctest.RandString(5)
	rac := acctest.ResourceAttrCheckInit(rc_proxypublic, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBProxyoublic%s", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceRdsDBProxyDependence_Public_MySQL)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":                      "${alicloud_db_instance.default.id}",
					"db_proxy_endpoint_id":                "${alicloud_rds_db_proxy.default.db_proxy_endpoint_id}",
					"connection_string_prefix":            "${alicloud_db_instance.default.id}abc",
					"db_proxy_connection_string_net_type": "Public",
					"db_proxy_new_connect_string_port":    "3307",
				}),
				Check: resource.ComposeTestCheckFunc(
					rc_instance.CheckResourceExists(),
					rc_proxy.CheckResourceExists(),
					testAccCheck(map[string]string{
						"db_instance_id":                      acctest.CHECKSET,
						"db_proxy_endpoint_id":                acctest.CHECKSET,
						"connection_string_prefix":            acctest.CHECKSET,
						"db_proxy_connection_string_net_type": acctest.CHECKSET,
						"db_proxy_new_connect_string_port":    acctest.CHECKSET,
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

func TestAccAliCloudRdsDBProxy_Public_MySQL_Update(t *testing.T) {
	var instance map[string]interface{}
	var proxy map[string]interface{}
	var proxypublic map[string]interface{}

	resourceId := "alicloud_rds_db_proxy_public.default"
	var DBProxyMap = map[string]string{
		"db_instance_id":                      acctest.CHECKSET,
		"db_proxy_endpoint_id":                acctest.CHECKSET,
		"connection_string_prefix":            acctest.CHECKSET,
		"db_proxy_connection_string_net_type": "Public",
		"db_proxy_new_connect_string_port":    acctest.CHECKSET,
	}
	ra := acctest.ResourceAttrInit(resourceId, DBProxyMap)
	rc_proxypublic := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &proxypublic, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBProxy")
	rc_instance := acctest.ResourceCheckInitWithDescribeMethod("alicloud_db_instance.default", &instance, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBInstance")
	rc_proxy := acctest.ResourceCheckInitWithDescribeMethod("alicloud_rds_db_proxy.default", &proxy, func() interface{} {
		return rds.NewRdsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDBProxy")

	rand := tfacctest.RandString(5)
	rac := acctest.ResourceAttrCheckInit(rc_proxypublic, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccDBProxyoublic%s", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceRdsDBProxyDependence_Public_MySQL)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,

		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"db_instance_id":                      "${alicloud_db_instance.default.id}",
					"db_proxy_endpoint_id":                "${alicloud_rds_db_proxy.default.db_proxy_endpoint_id}",
					"connection_string_prefix":            "${alicloud_db_instance.default.id}abc",
					"db_proxy_connection_string_net_type": "Public",
				}),
				Check: resource.ComposeTestCheckFunc(
					rc_instance.CheckResourceExists(),
					rc_proxy.CheckResourceExists(),
					testAccCheck(map[string]string{
						"db_instance_id":                      acctest.CHECKSET,
						"db_proxy_endpoint_id":                acctest.CHECKSET,
						"connection_string_prefix":            acctest.CHECKSET,
						"db_proxy_connection_string_net_type": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"connection_string_prefix": "${alicloud_db_instance.default.id}def",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"connection_string_prefix": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"db_proxy_new_connect_string_port": "3308",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_proxy_new_connect_string_port": "3308",
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

func resourceRdsDBProxyDependence_Public_MySQL(name string) string {
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

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_db_instance" "default" {
    engine = "MySQL"
	engine_version = "8.0"
 	db_instance_storage_type = "cloud_essd"
	instance_type = "mysql.x2.large.2c"
	instance_storage = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.0.min
	vswitch_id = data.alicloud_vswitches.default.ids.0
	instance_name = var.name
}

resource "alicloud_rds_db_proxy" "default" {
	instance_id = alicloud_db_instance.default.id
	instance_network_type = "VPC"
	db_proxy_instance_num = 2
 	vpc_id = alicloud_db_instance.default.vpc_id
	vswitch_id = alicloud_db_instance.default.vswitch_id
}

`, name)
}
