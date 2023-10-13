package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Redis TairInstance. >>> Resource test cases, automatically generated.
// Case 3314
func TestAccAliCloudRedisTairInstance_basic3314(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap3314)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence3314)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":       "Subscription",
					"period":             "1",
					"instance_type":      "tair_rdb",
					"zone_id":            "${local.zone_id}",
					"instance_class":     "tair.rdb.2g",
					"shard_count":        "2",
					"vswitch_id":         "${local.vswitch_id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"tair_instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "Subscription",
						"period":             "1",
						"instance_type":      "tair_rdb",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"shard_count":        "2",
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"tair_instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "5.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tair_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tair_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "tair.rdb.4g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "tair.rdb.4g",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tair_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tair_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":       "Subscription",
					"instance_type":      "tair_rdb",
					"password":           "Pass!123456",
					"engine_version":     "5.0",
					"zone_id":            "${local.zone_id}",
					"instance_class":     "tair.rdb.2g",
					"tair_instance_name": name + "_update",
					"shard_count":        "2",
					"secondary_zone_id":  "${local.zone_id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"vswitch_id":         "${local.vswitch_id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"auto_renew_period":  "12",
					"period":             "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "Subscription",
						"instance_type":      "tair_rdb",
						"password":           "Pass!123456",
						"engine_version":     "5.0",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"tair_instance_name": name + "_update",
						"shard_count":        "2",
						"secondary_zone_id":  CHECKSET,
						"resource_group_id":  CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"auto_renew_period":  "12",
						"period":             "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "effective_time", "force_upgrade", "password", "period"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap3314 = map[string]string{
	"resource_group_id": CHECKSET,
	"port":              CHECKSET,
	"payment_type":      CHECKSET,
	"status":            CHECKSET,
	"engine_version":    CHECKSET,
	"create_time":       CHECKSET,
	"storage_size_gb":   CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence3314(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_kvstore_zones" "default" {
  product_type = "Tair_rdb"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
  count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vswitch_name = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
}


`, name)
}

// Case 3340
func TestAccAliCloudRedisTairInstance_basic3340(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap3340)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence3340)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":      "tair_rdb",
					"zone_id":            "${local.zone_id}",
					"instance_class":     "tair.rdb.2g",
					"shard_count":        "2",
					"vswitch_id":         "${local.vswitch_id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"tair_instance_name": name,
					"payment_type":       "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":      "tair_rdb",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"shard_count":        "2",
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"tair_instance_name": name,
						"payment_type":       "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"shard_count": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"shard_count": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "5.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "tair.rdb.4g",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "tair.rdb.4g",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tair_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tair_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Pass!123456Change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Pass!123456Change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":         "false",
					"port":               "6379",
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_rdb",
					"password":           "Pass!123456",
					"engine_version":     "5.0",
					"zone_id":            "${local.zone_id}",
					"instance_class":     "tair.rdb.1g",
					"tair_instance_name": name + "_update",
					"shard_count":        "2",
					"secondary_zone_id":  "${local.zone_id}",
					"vswitch_id":         "${local.vswitch_id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":         "false",
						"port":               "6379",
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_rdb",
						"password":           "Pass!123456",
						"engine_version":     "5.0",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.1g",
						"tair_instance_name": name + "_update",
						"shard_count":        "2",
						"secondary_zone_id":  CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "effective_time", "force_upgrade", "password", "period"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap3340 = map[string]string{
	"resource_group_id": CHECKSET,
	"port":              CHECKSET,
	"payment_type":      CHECKSET,
	"status":            CHECKSET,
	"engine_version":    CHECKSET,
	"create_time":       CHECKSET,
	"storage_size_gb":   CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence3340(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_kvstore_zones" "default" {
  product_type = "Tair_rdb"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
  count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vswitch_name = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
}


`, name)
}

// Case 3549
func TestAccAliCloudRedisTairInstance_basic3549(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap3549)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence3549)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_type":      "tair_rdb",
					"zone_id":            "${local.zone_id}",
					"instance_class":     "tair.rdb.2g",
					"vswitch_id":         "${local.vswitch_id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"tair_instance_name": name,
					"payment_type":       "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":      "tair_rdb",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"tair_instance_name": name,
						"payment_type":       "PayAsYouGo",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Pass!123456",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Pass!123456",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "5.0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tair_instance_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tair_instance_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"password": "Pass!123456!change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"password": "Pass!123456!change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.groups.1.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
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
					"auto_renew":         "false",
					"port":               "6379",
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_rdb",
					"password":           "Pass!123456",
					"engine_version":     "5.0",
					"zone_id":            "${local.zone_id}",
					"instance_class":     "tair.rdb.2g",
					"tair_instance_name": name + "_update",
					"secondary_zone_id":  "${local.zone_id}",
					"vswitch_id":         "${local.vswitch_id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":         "false",
						"port":               "6379",
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_rdb",
						"password":           "Pass!123456",
						"engine_version":     "5.0",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"tair_instance_name": name + "_update",
						"secondary_zone_id":  CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "effective_time", "force_upgrade", "password", "period"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap3549 = map[string]string{
	"resource_group_id": CHECKSET,
	"port":              CHECKSET,
	"payment_type":      CHECKSET,
	"status":            CHECKSET,
	"engine_version":    CHECKSET,
	"create_time":       CHECKSET,
	"storage_size_gb":   CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence3549(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_kvstore_zones" "default" {
  product_type         = "Tair_rdb"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_vswitch" "this" {
  count = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vswitch_name = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
  cidr_block = cidrsubnet(data.alicloud_vpcs.default.vpcs.0.cidr_block, 8, 4)
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : concat(alicloud_vswitch.this.*.id, [""])[0]
  zone_id = data.alicloud_kvstore_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

// Case 3314  twin
func TestAccAliCloudRedisTairInstance_basic3314_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap3314)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence3314)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":         "false",
					"port":               "6379",
					"payment_type":       "Subscription",
					"instance_type":      "tair_rdb",
					"password":           "Pass!123456",
					"engine_version":     "5.0",
					"zone_id":            "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_class":     "tair.rdb.2g",
					"tair_instance_name": name,
					"shard_count":        "2",
					"secondary_zone_id":  "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"vswitch_id":         "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":             "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"auto_renew_period":  "12",
					"period":             "1",
					"effective_time":     "Immediately",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":         "false",
						"port":               "6379",
						"payment_type":       "Subscription",
						"instance_type":      "tair_rdb",
						"password":           "Pass!123456",
						"engine_version":     "5.0",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"tair_instance_name": name,
						"shard_count":        "2",
						"secondary_zone_id":  CHECKSET,
						"resource_group_id":  CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"auto_renew_period":  "12",
						"period":             "1",
						"effective_time":     "Immediately",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "effective_time", "force_upgrade", "password", "period"},
			},
		},
	})
}

// Case 3340  twin
func TestAccAliCloudRedisTairInstance_basic3340_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap3340)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence3340)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.RedisTariInstanceSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":         "false",
					"port":               "6379",
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_rdb",
					"password":           "Pass!123456!change",
					"engine_version":     "5.0",
					"zone_id":            "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"instance_class":     "tair.rdb.1g",
					"tair_instance_name": name,
					"shard_count":        "2",
					"secondary_zone_id":  "${data.alicloud_vswitches.default.vswitches.0.zone_id}",
					"vswitch_id":         "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":             "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"effective_time":     "Immediately",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":         "false",
						"port":               "6379",
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_rdb",
						"password":           "Pass!123456!change",
						"engine_version":     "5.0",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.1g",
						"tair_instance_name": name,
						"shard_count":        "2",
						"secondary_zone_id":  CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
						"effective_time":     "Immediately",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "effective_time", "force_upgrade", "password", "period"},
			},
		},
	})
}

// Case 3549  twin
func TestAccAliCloudRedisTairInstance_basic3549_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap3549)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence3549)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"auto_renew":         "false",
					"port":               "6379",
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_rdb",
					"password":           "Pass!123456!change",
					"engine_version":     "5.0",
					"zone_id":            "${local.zone_id}",
					"instance_class":     "tair.rdb.2g",
					"tair_instance_name": name,
					"secondary_zone_id":  "${local.zone_id}",
					"vswitch_id":         "${data.alicloud_vswitches.default.vswitches.0.id}",
					"vpc_id":             "${data.alicloud_vswitches.default.vswitches.0.vpc_id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"effective_time":     "Immediately",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":         "false",
						"port":               "6379",
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_rdb",
						"password":           "Pass!123456!change",
						"engine_version":     "5.0",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"tair_instance_name": name,
						"secondary_zone_id":  CHECKSET,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
						"effective_time":     "Immediately",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "effective_time", "force_upgrade", "password", "period"},
			},
		},
	})
}

// Case 4491  twin
func TestAccAliCloudRedisTairInstance_basic4491_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap4491)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence4491)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DRDSPolarDbxSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"force_upgrade":             "false",
					"auto_renew":                "false",
					"port":                      "6379",
					"payment_type":              "PayAsYouGo",
					"instance_type":             "tair_essd",
					"zone_id":                   "${alicloud_vswitch.default.zone_id}",
					"period":                    "9",
					"instance_class":            "tair.essd.standard.xlarge",
					"tair_instance_name":        name,
					"auto_renew_period":         "1",
					"secondary_zone_id":         "${alicloud_vswitch.default.zone_id}",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"vpc_id":                    "${alicloud_vswitch.default.vpc_id}",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"storage_performance_level": "PL1",
					"storage_size_gb":           "20",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"auto_renew":                "false",
						"port":                      "6379",
						"payment_type":              "PayAsYouGo",
						"instance_type":             "tair_essd",
						"zone_id":                   CHECKSET,
						"instance_class":            "tair.essd.standard.xlarge",
						"tair_instance_name":        name,
						"auto_renew_period":         "1",
						"secondary_zone_id":         CHECKSET,
						"vswitch_id":                CHECKSET,
						"vpc_id":                    CHECKSET,
						"resource_group_id":         CHECKSET,
						"storage_performance_level": "PL1",
						"storage_size_gb":           "20",
						"tags.%":                    "2",
						"tags.Created":              "TF",
						"tags.For":                  "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "effective_time", "force_upgrade", "period", "force_upgrade"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap4491 = map[string]string{
	"port":           CHECKSET,
	"payment_type":   "PayAsYouGo",
	"status":         CHECKSET,
	"engine_version": CHECKSET,
	"create_time":    CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence4491(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_kvstore_zones" "default" {
  product_type = "Tair_essd"
  instance_charge_type = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/24"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "192.168.0.0/24"
  zone_id      = data.alicloud_kvstore_zones.default.zones.1.id
  vpc_id       = alicloud_vpc.default.id
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

// Test Redis TairInstance. <<< Resource test cases, automatically generated.
