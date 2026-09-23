package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Gpdb ExternalDataService. >>> Resource test cases, automatically generated.
// Case 外部数据服务_资源依赖_case 6969
func TestAccAliCloudGpdbExternalDataService_basic6969(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_external_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbExternalDataServiceMap6969)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbExternalDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbexternaldataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbExternalDataServiceBasicDependence6969)
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
					"service_name":   "test6",
					"db_instance_id": "${alicloud_gpdb_instance.defaultZ7DPgB.id}",
					"service_spec":   "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":   "test6",
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
					"service_description": "test6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test6",
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
					"service_name":        "test6",
					"db_instance_id":      "${alicloud_gpdb_instance.defaultZ7DPgB.id}",
					"service_description": "test",
					"service_spec":        "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":        "test6",
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

var AlicloudGpdbExternalDataServiceMap6969 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudGpdbExternalDataServiceBasicDependence6969(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultrple4a" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultnYWSkl" {
  vpc_id     = alicloud_vpc.defaultrple4a.id
  zone_id    = "cn-beijing-i"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaultZ7DPgB" {
  instance_spec              = "2C8G"
  description                = "创建数据服务依赖的实例"
  seg_node_num               = "2"
  seg_storage_type           = "cloud_essd"
  instance_network_type      = "VPC"
  db_instance_category       = "Basic"
  payment_type               = "PayAsYouGo"
  ssl_enabled                = "0"
  engine_version             = "6.0"
  zone_id                    = "cn-beijing-i"
  vswitch_id                 = alicloud_vswitch.defaultnYWSkl.id
  storage_size               = "50"
  master_cu                  = "4"
  vpc_id                     = alicloud_vpc.defaultrple4a.id
  db_instance_mode           = "StorageElastic"
  engine                     = "gpdb"
}


`, name)
}

// Case 外部数据服务_资源依赖_case 6969  twin
func TestAccAliCloudGpdbExternalDataService_basic6969_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_external_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbExternalDataServiceMap6969)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbExternalDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbexternaldataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbExternalDataServiceBasicDependence6969)
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
					"service_name":        "test6",
					"db_instance_id":      "${alicloud_gpdb_instance.defaultZ7DPgB.id}",
					"service_description": "test",
					"service_spec":        "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":        "test6",
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

// Case 外部数据服务_资源依赖_case 6969  raw
func TestAccAliCloudGpdbExternalDataService_basic6969_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_gpdb_external_data_service.default"
	ra := resourceAttrInit(resourceId, AlicloudGpdbExternalDataServiceMap6969)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &GpdbServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeGpdbExternalDataService")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sgpdbexternaldataservice%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudGpdbExternalDataServiceBasicDependence6969)
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
					"service_name":        "test6",
					"db_instance_id":      "${alicloud_gpdb_instance.defaultZ7DPgB.id}",
					"service_description": "test",
					"service_spec":        "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_name":        "test6",
						"db_instance_id":      CHECKSET,
						"service_description": "test",
						"service_spec":        "8",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"service_description": "test6",
					"service_spec":        "16",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"service_description": "test6",
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

// Test Gpdb ExternalDataService. <<< Resource test cases, automatically generated.
