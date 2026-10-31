package gpdb_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/gpdb"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Gpdb StreamingDataService. >>> Resource test cases, automatically generated.
// Case 实时消费服务_资源依赖_case 6971
func TestAccAliCloudGpdbStreamingDataService_basic6971(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap6971)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeStreamingDataService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence6971)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":        "test",
					"db_instance_id":      "${alicloud_gpdb_instance.default1oSPzX.id}",
					"service_description": "test",
					"service_spec":        "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":        "test",
						"db_instance_id":      acctest.CHECKSET,
						"service_description": "test",
						"service_spec":        acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "test2",
					"service_spec":        "16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test2",
						"service_spec":        acctest.CHECKSET,
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

var AlicloudGpdbStreamingDataServiceMap6971 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"service_id":  acctest.CHECKSET,
}

func AlicloudGpdbStreamingDataServiceBasicDependence6971(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultTXZPBL" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultrJ5mmz" {
  vpc_id     = alicloud_vpc.defaultTXZPBL.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "default1oSPzX" {
  instance_spec              = "2C8G"
  description                = "创建流服务需要的实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  security_ip_list           = ["127.0.0.1"]
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-h"
  vswitch_id                 = alicloud_vswitch.defaultrJ5mmz.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultTXZPBL.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
  seg_disk_performance_level = "pl1"
}


`, name)
}

// Case 实时消费服务_资源依赖_case_002 7052
func TestAccAliCloudGpdbStreamingDataService_basic7052(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap7052)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return gpdb.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeStreamingDataService")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccgpdb%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence7052)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":        "test",
					"db_instance_id":      "${alicloud_gpdb_instance.default1oSPzX.id}",
					"service_description": "test",
					"service_spec":        "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":        "test",
						"db_instance_id":      acctest.CHECKSET,
						"service_description": "test",
						"service_spec":        acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "test2",
					"service_spec":        "16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test2",
						"service_spec":        acctest.CHECKSET,
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

var AlicloudGpdbStreamingDataServiceMap7052 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
	"service_id":  acctest.CHECKSET,
}

func AlicloudGpdbStreamingDataServiceBasicDependence7052(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultTXZPBL" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultrJ5mmz" {
  vpc_id     = alicloud_vpc.defaultTXZPBL.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "default1oSPzX" {
  instance_spec              = "2C8G"
  description                = "创建流服务需要的实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  security_ip_list           = ["127.0.0.1"]
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-h"
  vswitch_id                 = alicloud_vswitch.defaultrJ5mmz.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultTXZPBL.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
  seg_disk_performance_level = "pl1"
}


`, name)
}

// Test Gpdb StreamingDataService. <<< Resource test cases, automatically generated.
