package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb DBInstanceIPArray. >>> Resource test cases, automatically generated.
// Case IPArray测试 7637
func TestAccAliCloudGpdbDBInstanceIPArray_basic7637(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_db_instance_ip_array.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbDBInstanceIPArrayMap7637)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbDBInstanceIPArray")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbdbinstanceiparray%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbDBInstanceIPArrayBasicDependence7637)
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
					"db_instance_ip_array_attribute": "taffyFish",
					"security_ip_list": []string{
						"12.34.56.78", "11.45.14.0", "19.19.81.0"},
					"db_instance_ip_array_name": "taffy",
					"db_instance_id":            "${alicloud_gpdb_instance.defaultHKdDs3.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"db_instance_ip_array_attribute": "taffyFish",
						"security_ip_list.#":             "3",
						"db_instance_ip_array_name":      "taffy",
						"db_instance_id":                 CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_list": []string{
						"18.18.18.18", "19.19.19.19"},
					"modify_mode": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_list.#": "2",
						"modify_mode":        "0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"modify_mode"},
			},
		},
	})
}

var AlicloudGpdbDBInstanceIPArrayMap7637 = map[string]string{}

func AlicloudGpdbDBInstanceIPArrayBasicDependence7637(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultNpLRa1" {
	cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultwLA5v4" {
	vpc_id = "${alicloud_vpc.defaultNpLRa1.id}"
	zone_id = "cn-beijing-h"
	cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultHKdDs3" {
	instance_spec = "2C8G"
	seg_node_num = "2"
	seg_storage_type = "cloud_essd"
	instance_network_type = "VPC"
	db_instance_category = "Basic"
	payment_type = "PayAsYouGo"
	ssl_enabled = "0"
	engine_version = "6.0"
	zone_id = "cn-beijing-h"
	vswitch_id = "${alicloud_vswitch.defaultwLA5v4.id}"
	storage_size = "50"
	master_cu = "4"
	vpc_id = "${alicloud_vpc.defaultNpLRa1.id}"
	db_instance_mode = "StorageElastic"
  	engine           = "gpdb"
    description = var.name
}


`, name)
}

// Test Gpdb DBInstanceIPArray. <<< Resource test cases, automatically generated.
