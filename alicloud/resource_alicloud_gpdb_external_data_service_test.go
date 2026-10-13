package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
			testAccPreCheckWithRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":   "test6",
					"db_instance_id": "${alicloud_gpdb_instance.default.id}",
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

data "alicloud_gpdb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids.0
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
			testAccPreCheckWithRegions(t, true, connectivity.GPDBDBInstancePlanSupportRegions)
		},
		IDRefreshName:     resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"service_name":        "test6",
					"db_instance_id":      "${alicloud_gpdb_instance.default.id}",
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

// Test Gpdb ExternalDataService. <<< Resource test cases, automatically generated.
