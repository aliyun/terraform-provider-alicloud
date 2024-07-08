package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb StreamingDataService. >>> Resource test cases, automatically generated.
// Case 实时消费服务_资源依赖_case_002 7052
func TestAccAliCloudGpdbStreamingDataService_basic7052(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap7052)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence7052)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":   "test",
					"db_instance_id": "${alicloud_gpdb_instance.default1oSPzX.id}",
					"service_spec":   "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":   "test",
						"db_instance_id": CHECKSET,
						"service_spec":   "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_spec": "16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_spec": "16",
					}),
				),
			},
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
						"db_instance_id":      CHECKSET,
						"service_description": "test",
						"service_spec":        "8",
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"service_id":  CHECKSET,
}

func AlicloudGpdbStreamingDataServiceBasicDependence7052(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
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
}


`, name)
}

// Case 实时消费服务_资源依赖_case 6971
func TestAccAliCloudGpdbStreamingDataService_basic6971(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap6971)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence6971)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":   "test",
					"db_instance_id": "${alicloud_gpdb_instance.default1oSPzX.id}",
					"service_spec":   "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":   "test",
						"db_instance_id": CHECKSET,
						"service_spec":   "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "test2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_spec": "16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_spec": "16",
					}),
				),
			},
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
						"db_instance_id":      CHECKSET,
						"service_description": "test",
						"service_spec":        "8",
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
	"status":      CHECKSET,
	"create_time": CHECKSET,
	"service_id":  CHECKSET,
}

func AlicloudGpdbStreamingDataServiceBasicDependence6971(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
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
}


`, name)
}

// Case 实时消费服务_资源依赖_case_002 7052  twin
func TestAccAliCloudGpdbStreamingDataService_basic7052_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap7052)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence7052)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"db_instance_id":      CHECKSET,
						"service_description": "test",
						"service_spec":        "8",
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

// Case 实时消费服务_资源依赖_case 6971  twin
func TestAccAliCloudGpdbStreamingDataService_basic6971_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap6971)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence6971)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"db_instance_id":      CHECKSET,
						"service_description": "test",
						"service_spec":        "8",
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

// Case 实时消费服务_资源依赖_case_002 7052  raw
func TestAccAliCloudGpdbStreamingDataService_basic7052_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap7052)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence7052)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"db_instance_id":      CHECKSET,
						"service_description": "test",
						"service_spec":        "8",
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
						"service_spec":        "16",
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

// Case 实时消费服务_资源依赖_case 6971  raw
func TestAccAliCloudGpdbStreamingDataService_basic6971_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_streaming_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbStreamingDataServiceMap6971)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbStreamingDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbstreamingdataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbStreamingDataServiceBasicDependence6971)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"db_instance_id":      CHECKSET,
						"service_description": "test",
						"service_spec":        "8",
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
						"service_spec":        "16",
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

// Test Gpdb StreamingDataService. <<< Resource test cases, automatically generated.
