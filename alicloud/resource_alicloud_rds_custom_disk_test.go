package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Rds CustomDisk. >>> Resource test cases, automatically generated.
// Case CustomDisk cloud_auto 10802
func TestAccAliCloudRdsCustomDisk_basic10802(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsCustomDiskMap10802)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsCustomDiskBasicDependence10802)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":       "${data.alicloud_vswitches.default.zone_id}",
					"size":          "500",
					"disk_category": "cloud_auto",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":       CHECKSET,
						"size":          "500",
						"disk_category": "cloud_auto",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "510",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "510",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "type", "dry_run"},
			},
		},
	})
}

func TestAccAliCloudRdsCustomDisk_basic10802_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsCustomDiskMap10802)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsCustomDiskBasicDependence10802)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":              "${data.alicloud_vswitches.default.zone_id}",
					"size":                 "500",
					"disk_category":        "cloud_auto",
					"auto_pay":             "true",
					"bursting_enabled":     "false",
					"delete_with_instance": "true",
					"description":          name,
					"disk_name":            name,
					"instance_charge_type": "Postpaid",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":              CHECKSET,
						"size":                 "500",
						"disk_category":        "cloud_auto",
						"bursting_enabled":     "false",
						"delete_with_instance": "true",
						"description":          name,
						"disk_name":            name,
						"instance_charge_type": "Postpaid",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "type", "dry_run"},
			},
		},
	})
}

var AliCloudRdsCustomDiskMap10802 = map[string]string{
	"bursting_enabled":     CHECKSET,
	"disk_name":            CHECKSET,
	"instance_charge_type": CHECKSET,
	"region_id":            CHECKSET,
	"resource_group_id":    CHECKSET,
	"create_time":          CHECKSET,
	"status":               CHECKSET,
}

func AliCloudRdsCustomDiskBasicDependence10802(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-i"
}
`, name)
}

// Case CustomDisk cloud_essd 10806
func TestAccAliCloudRdsCustomDisk_basic10806(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsCustomDiskMap10806)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsCustomDiskBasicDependence10806)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":              "${data.alicloud_vswitches.default.zone_id}",
					"size":                 "500",
					"disk_category":        "cloud_essd",
					"instance_charge_type": "Prepaid",
					"instance_id":          "${alicloud_rds_custom.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":              CHECKSET,
						"size":                 "500",
						"disk_category":        "cloud_essd",
						"instance_charge_type": "Prepaid",
						"instance_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delete_with_instance": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delete_with_instance": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"performance_level": "PL2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"performance_level": "PL2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disk_category":     "cloud_auto",
					"performance_level": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disk_category":     "cloud_auto",
						"performance_level": REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bursting_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bursting_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"size": "510",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"size": "510",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "type", "dry_run"},
			},
		},
	})
}

func TestAccAliCloudRdsCustomDisk_basic10806_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_rds_custom_disk.default"
	ra := resourceAttrInit(resourceId, AliCloudRdsCustomDiskMap10806)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RdsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRdsCustomDisk")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccrds%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRdsCustomDiskBasicDependence10806)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"zone_id":              "${data.alicloud_vswitches.default.zone_id}",
					"size":                 "500",
					"disk_category":        "cloud_essd",
					"instance_charge_type": "Prepaid",
					"instance_id":          "${alicloud_rds_custom.default.id}",
					"auto_pay":             "true",
					"auto_renew":           "true",
					"delete_with_instance": "false",
					"description":          name,
					"disk_name":            name,
					"performance_level":    "PL2",
					"resource_group_id":    "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_id":              CHECKSET,
						"size":                 "500",
						"disk_category":        "cloud_essd",
						"instance_charge_type": "Prepaid",
						"instance_id":          CHECKSET,
						"delete_with_instance": "false",
						"description":          name,
						"disk_name":            name,
						"performance_level":    "PL2",
						"resource_group_id":    CHECKSET,
						"tags.%":               "2",
						"tags.Created":         "TF",
						"tags.For":             "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "auto_renew", "type", "dry_run"},
			},
		},
	})
}

var AliCloudRdsCustomDiskMap10806 = map[string]string{
	"disk_name":            CHECKSET,
	"instance_charge_type": CHECKSET,
	"region_id":            CHECKSET,
	"resource_group_id":    CHECKSET,
	"create_time":          CHECKSET,
	"status":               CHECKSET,
}

func AliCloudRdsCustomDiskBasicDependence10806(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-hangzhou-i"
}

data "alicloud_security_groups" "default" {
  vpc_id     = data.alicloud_vpcs.default.ids.0
  name_regex = "default-NODELETING"
}

resource "alicloud_rds_custom" "default" {
  zone_id              = data.alicloud_vswitches.default.zone_id
  instance_charge_type = "Prepaid"
  vswitch_id           = data.alicloud_vswitches.default.ids.0
  amount               = "1"
  security_group_ids   = [data.alicloud_security_groups.default.ids.0]
  system_disk {
    size = "40"
  }
  force         = true
  instance_type = "mysql.x4.xlarge.6cm"
  spot_strategy = "NoSpot"
}
`, name)
}

// Test Rds CustomDisk. <<< Resource test cases, automatically generated.
