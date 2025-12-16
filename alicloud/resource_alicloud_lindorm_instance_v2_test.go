// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Lindorm InstanceV2. >>> Resource test cases, automatically generated.
// Case InstanceV2修改删除白名单测试用例 11876
func TestAccAliCloudLindormInstanceV2_basic11876(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceV2Map11876)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstanceV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacclindorm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceV2BasicDependence11876)
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
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "yn-rg-ips",
								},
							},
						},
					},
					"cloud_storage_size": "320",
					"zone_id":            "cn-beijing-l",
					"cloud_storage_type": "PerformanceStorage",
					"arch_version":       "1.0",
					"vswitch_id":         "${alicloud_vswitch.default9rq7dN.id}",
					"vpc_id":             "${alicloud_vpc.defaultVLjfBs.id}",
					"instance_alias":     "preTest-modify-ips",
					"payment_type":       "POSTPAY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":      "1",
						"cloud_storage_size": "320",
						"zone_id":            "cn-beijing-l",
						"cloud_storage_type": "PerformanceStorage",
						"arch_version":       CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"instance_alias":     "preTest-modify-ips",
						"payment_type":       "POSTPAY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.c.2xlarge",
									"resource_group_name": "yn-rg-ips",
								},
							},
						},
					},
					"instance_alias": "preTest-modify-ips-x",
					"white_ip_list": []map[string]interface{}{
						{
							"group_name": "user001",
							"ip_list":    "127.0.0.1",
						},
						{
							"group_name": "user002",
							"ip_list":    "127.0.0.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":   "1",
						"instance_alias":  "preTest-modify-ips-x",
						"white_ip_list.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arch_version", "auto_renew_duration", "auto_renewal", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudLindormInstanceV2Map11876 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudLindormInstanceV2BasicDependence11876(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultVLjfBs" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default9rq7dN" {
  vpc_id     = alicloud_vpc.defaultVLjfBs.id
  cidr_block = "10.0.0.0/16"
  zone_id    = "cn-beijing-l"
}


`, name)
}

// Case InstanceV2单可用区预付费实例 11753
func TestAccAliCloudLindormInstanceV2_basic11753(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceV2Map11753)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstanceV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacclindorm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceV2BasicDependence11753)
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
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "chixiao-rg-test",
									"node_disk_type":      "cloud_essd",
									"node_disk_size":      "200",
								},
							},
						},
					},
					"arch_version":        "1.0",
					"vswitch_id":          "${alicloud_vswitch.defaultcEoant.id}",
					"vpc_id":              "${alicloud_vpc.defaultHgtTqh.id}",
					"instance_alias":      "preTest",
					"payment_type":        "PREPAY",
					"zone_id":             "cn-beijing-l",
					"pricing_cycle":       "Month",
					"duration":            "1",
					"auto_renew_duration": "1",
					"auto_renewal":        "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":       "1",
						"arch_version":        CHECKSET,
						"vswitch_id":          CHECKSET,
						"vpc_id":              CHECKSET,
						"instance_alias":      "preTest",
						"payment_type":        "PREPAY",
						"zone_id":             "cn-beijing-l",
						"pricing_cycle":       "Month",
						"duration":            "1",
						"auto_renew_duration": CHECKSET,
						"auto_renewal":        "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "chixiao-rg-test",
									"node_disk_size":      "400",
									"node_disk_type":      "cloud_essd",
								},
							},
						},
					},
					"instance_alias":      "preTest-cx",
					"payment_type":        "POSTPAY",
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":       "1",
						"instance_alias":      "preTest-cx",
						"payment_type":        "POSTPAY",
						"deletion_protection": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arch_version", "auto_renew_duration", "auto_renewal", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudLindormInstanceV2Map11753 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudLindormInstanceV2BasicDependence11753(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultHgtTqh" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultcEoant" {
  zone_id    = "cn-beijing-l"
  cidr_block = "10.0.0.0/16"
  vpc_id     = alicloud_vpc.defaultHgtTqh.id
}


`, name)
}

// Case InstanceV2多可用区测试用例 11176
func TestAccAliCloudLindormInstanceV2_basic11176(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceV2Map11176)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstanceV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacclindorm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceV2BasicDependence11176)
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
					"standby_zone_id": "cn-beijing-l",
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "4",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "cx-mz-rg",
								},
							},
						},
					},
					"cloud_storage_size": "400",
					"primary_zone_id":    "cn-beijing-h",
					"zone_id":            "cn-beijing-h",
					"cloud_storage_type": "PerformanceStorage",
					"arch_version":       "2.0",
					"vswitch_id":         "${alicloud_vswitch.default9umuzwH.id}",
					"standby_vswitch_id": "${alicloud_vswitch.defaultgOFAo3L.id}",
					"primary_vswitch_id": "${alicloud_vswitch.default9umuzwH.id}",
					"arbiter_vswitch_id": "${alicloud_vswitch.defaultTAbr2pJ.id}",
					"vpc_id":             "${alicloud_vpc.defaultR8vXlP.id}",
					"instance_alias":     "preTest-MZ",
					"payment_type":       "POSTPAY",
					"arbiter_zone_id":    "cn-beijing-j",
					"auto_renewal":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"standby_zone_id":    "cn-beijing-l",
						"engine_list.#":      "1",
						"cloud_storage_size": "400",
						"primary_zone_id":    "cn-beijing-h",
						"zone_id":            "cn-beijing-h",
						"cloud_storage_type": "PerformanceStorage",
						"arch_version":       CHECKSET,
						"vswitch_id":         CHECKSET,
						"standby_vswitch_id": CHECKSET,
						"primary_vswitch_id": CHECKSET,
						"arbiter_vswitch_id": CHECKSET,
						"vpc_id":             CHECKSET,
						"instance_alias":     "preTest-MZ",
						"payment_type":       "POSTPAY",
						"arbiter_zone_id":    "cn-beijing-j",
						"auto_renewal":       "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cloud_storage_size":  "800",
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cloud_storage_size":  "800",
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arch_version", "auto_renew_duration", "auto_renewal", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudLindormInstanceV2Map11176 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudLindormInstanceV2BasicDependence11176(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultR8vXlP" {
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default9umuzwH" {
  vpc_id     = alicloud_vpc.defaultR8vXlP.id
  zone_id    = "cn-beijing-h"
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_vswitch" "defaultgOFAo3L" {
  vpc_id     = alicloud_vpc.defaultR8vXlP.id
  zone_id    = "cn-beijing-l"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "defaultTAbr2pJ" {
  vpc_id     = alicloud_vpc.defaultR8vXlP.id
  zone_id    = "cn-beijing-j"
  cidr_block = "172.16.2.0/24"
}


`, name)
}

// Case InstanceV2单AZ变配测试用例 11191
func TestAccAliCloudLindormInstanceV2_basic11191(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceV2Map11191)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstanceV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacclindorm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceV2BasicDependence11191)
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
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "cx-rg-spec",
								},
							},
						},
					},
					"cloud_storage_size": "320",
					"zone_id":            "cn-beijing-l",
					"cloud_storage_type": "PerformanceStorage",
					"arch_version":       "1.0",
					"vswitch_id":         "${alicloud_vswitch.default9rq7dN.id}",
					"vpc_id":             "${alicloud_vpc.defaultVLjfBs.id}",
					"instance_alias":     "preTest-upgrade-spec",
					"payment_type":       "POSTPAY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":      "1",
						"cloud_storage_size": "320",
						"zone_id":            "cn-beijing-l",
						"cloud_storage_type": "PerformanceStorage",
						"arch_version":       CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"instance_alias":     "preTest-upgrade-spec",
						"payment_type":       "POSTPAY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.c.2xlarge",
									"resource_group_name": "cx-rg-spec",
								},
							},
						},
					},
					"instance_alias": "preTest-upgrade-spec-cx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":  "1",
						"instance_alias": "preTest-upgrade-spec-cx",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arch_version", "auto_renew_duration", "auto_renewal", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudLindormInstanceV2Map11191 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudLindormInstanceV2BasicDependence11191(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultVLjfBs" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default9rq7dN" {
  vpc_id     = alicloud_vpc.defaultVLjfBs.id
  cidr_block = "10.0.0.0/16"
  zone_id    = "cn-beijing-l"
}


`, name)
}

// Case InstanceV2单可用区测试用例 11034
func TestAccAliCloudLindormInstanceV2_basic11034(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_lindorm_instance_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudLindormInstanceV2Map11034)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &LindormServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeLindormInstanceV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacclindorm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudLindormInstanceV2BasicDependence11034)
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
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "chixiao-rg-test",
									"node_disk_type":      "cloud_essd",
									"node_disk_size":      "200",
								},
							},
						},
					},
					"arch_version":   "1.0",
					"vswitch_id":     "${alicloud_vswitch.defaultcEoant.id}",
					"vpc_id":         "${alicloud_vpc.defaultHgtTqh.id}",
					"instance_alias": "preTest",
					"payment_type":   "POSTPAY",
					"zone_id":        "cn-beijing-l",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":  "1",
						"arch_version":   CHECKSET,
						"vswitch_id":     CHECKSET,
						"vpc_id":         CHECKSET,
						"instance_alias": "preTest",
						"payment_type":   "POSTPAY",
						"zone_id":        "cn-beijing-l",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "2",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "chixiao-rg-test",
									"node_disk_size":      "400",
									"node_disk_type":      "cloud_essd",
								},
							},
						},
					},
					"instance_alias":      "preTest-cx",
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":       "1",
						"instance_alias":      "preTest-cx",
						"deletion_protection": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_list": []map[string]interface{}{
						{
							"engine_type": "TABLE",
							"node_group": []map[string]interface{}{
								{
									"node_count":          "4",
									"node_spec":           "lindorm.g.2xlarge",
									"resource_group_name": "chixiao-rg-test",
									"node_disk_size":      "400",
									"node_disk_type":      "cloud_essd",
								},
							},
						},
					},
					"instance_alias":      "preTest-cx",
					"deletion_protection": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_list.#":       "1",
						"instance_alias":      "preTest-cx",
						"deletion_protection": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"deletion_protection": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"deletion_protection": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"arch_version", "auto_renew_duration", "auto_renewal", "duration", "pricing_cycle"},
			},
		},
	})
}

var AlicloudLindormInstanceV2Map11034 = map[string]string{
	"region_id": CHECKSET,
}

func AlicloudLindormInstanceV2BasicDependence11034(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "defaultHgtTqh" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultcEoant" {
  zone_id    = "cn-beijing-l"
  cidr_block = "10.0.0.0/16"
  vpc_id     = alicloud_vpc.defaultHgtTqh.id
}


`, name)
}

// Test Lindorm InstanceV2. <<< Resource test cases, automatically generated.
