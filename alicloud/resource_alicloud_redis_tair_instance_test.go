package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

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
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
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
					"instance_type":       "tair_rdb",
					"zone_id":             "${local.zone_id}",
					"instance_class":      "tair.rdb.2g",
					"shard_count":         "2",
					"vswitch_id":          "${local.vswitch_id}",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"tair_instance_name":  name,
					"payment_type":        "PayAsYouGo",
					"global_instance_id":  "true",
					"recover_config_mode": "whitelist,config",
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
					"intranet_bandwidth": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"intranet_bandwidth": "200",
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
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
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
						"tags.%":             "2",
						"tags.Created":       "TF",
						"tags.For":           "Test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
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
					"intranet_bandwidth": "200",
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
						"intranet_bandwidth": "200",
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
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
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
				Config: testAccConfig(map[string]interface{}{
					"payment_type": "PayAsYouGo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type": "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
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
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
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
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
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
					"port":                      "6379",
					"payment_type":              "PayAsYouGo",
					"instance_type":             "tair_essd",
					"zone_id":                   "${data.alicloud_kvstore_zones.default.zones.1.id}",
					"instance_class":            "tair.essd.standard.xlarge",
					"tair_instance_name":        name,
					"secondary_zone_id":         "${data.alicloud_kvstore_zones.default.zones.1.id}",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"vpc_id":                    "${alicloud_vpc.default.id}",
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
						"port":                      "6379",
						"payment_type":              "PayAsYouGo",
						"instance_type":             "tair_essd",
						"zone_id":                   CHECKSET,
						"instance_class":            "tair.essd.standard.xlarge",
						"tair_instance_name":        name,
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
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
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
  product_type         = "Tair_essd"
  instance_charge_type = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id = alicloud_vpc.default.id
  zone_id = data.alicloud_kvstore_zones.default.zones.1.id
  cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

data "alicloud_resource_manager_resource_groups" "default" {
}

`, name)
}

var AlicloudRedisTairInstanceMap6500 = map[string]string{
	"port":           CHECKSET,
	"status":         CHECKSET,
	"engine_version": CHECKSET,
	"payment_type":   "PayAsYouGo",
	"create_time":    CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence6500(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}

variable "secondary_zone_id" {
  default = "cn-beijing-i"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  description       = "测试请勿绑定test-ljt"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  cidr_block        = "172.16.0.0/12"
  vpc_name          = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_vswitch" "secondaryvsw" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.secondary_zone_id
  cidr_block   = "172.16.1.0/24"
  vswitch_name = var.name
}


`, name)
}

// Case Tair_rdb_单副本变高可用 6639
var AlicloudRedisTairInstanceMap6639 = map[string]string{
	"port":           CHECKSET,
	"status":         CHECKSET,
	"engine_version": CHECKSET,
	"payment_type":   "PayAsYouGo",
	"create_time":    CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence6639(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  description       = "测试请勿绑定test-ljt"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  cidr_block        = "172.16.0.0/12"
  vpc_name          = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}


`, name)
}

var AlicloudRedisTairInstanceMap6473 = map[string]string{
	"port":           CHECKSET,
	"status":         CHECKSET,
	"engine_version": CHECKSET,
	"payment_type":   "PayAsYouGo",
	"create_time":    CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence6473(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "secondary_zone_id" {
  default = "cn-beijing-i"
}

variable "region_id" {
  default = "cn-beijing"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultVpc" {
  description       = "测试请勿绑定test-ljt"
  cidr_block        = "172.16.0.0/12"
  vpc_name          = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name = var.name
  vpc_id = alicloud_vpc.defaultVpc.id
}

resource "alicloud_security_group" "change" {
  name = var.name
  vpc_id = alicloud_vpc.defaultVpc.id
}
`, name)
}

// Case Tair_rdb_rw_双可用区_修改slave只读节点 6500  raw
func TestAccAliCloudRedisTairInstance_basic6500_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap6500)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence6500)
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
					"payment_type":          "PayAsYouGo",
					"instance_type":         "tair_rdb",
					"zone_id":               "${var.zone_id}",
					"instance_class":        "tair.rdb.with.proxy.1g",
					"tair_instance_name":    name,
					"vswitch_id":            "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":                "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":              "123456Tf",
					"engine_version":        "5.0",
					"port":                  "6379",
					"slave_read_only_count": "2",
					"secondary_zone_id":     "${var.secondary_zone_id}",
					"shard_count":           "1",
					"read_only_count":       "2",
					"node_type":             "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":          "PayAsYouGo",
						"instance_type":         "tair_rdb",
						"zone_id":               CHECKSET,
						"instance_class":        "tair.rdb.with.proxy.1g",
						"tair_instance_name":    name,
						"vswitch_id":            CHECKSET,
						"vpc_id":                CHECKSET,
						"resource_group_id":     CHECKSET,
						"password":              "123456Tf",
						"engine_version":        "5.0",
						"port":                  "6379",
						"slave_read_only_count": "2",
						"secondary_zone_id":     CHECKSET,
						"shard_count":           "1",
						"read_only_count":       "2",
						"node_type":             "MASTER_SLAVE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"slave_read_only_count": "3",
					"read_only_count":       "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"slave_read_only_count": "3",
						"read_only_count":       "3",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
			},
		},
	})
}

// Case Tair_rdb_单副本变高可用 6639  raw
func TestAccAliCloudRedisTairInstance_basic6639_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap6639)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence6639)
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
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_rdb",
					"zone_id":            "${var.zone_id}",
					"instance_class":     "tair.rdb.1g",
					"tair_instance_name": name,
					"vswitch_id":         "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":             "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":           "123456Tf",
					"engine_version":     "5.0",
					"period":             "12",
					"auto_renew_period":  "12",
					"port":               "6379",
					"cluster_backup_id":  "cb-hyxdof5x9kqb333",
					"secondary_zone_id":  "${var.zone_id}",
					"node_type":          "STAND_ALONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_rdb",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.1g",
						"tair_instance_name": name,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
						"password":           "123456Tf",
						"engine_version":     "5.0",
						"period":             "12",
						"auto_renew_period":  "12",
						"port":               "6379",
						"cluster_backup_id":  "cb-hyxdof5x9kqb333",
						"secondary_zone_id":  CHECKSET,
						"node_type":          "STAND_ALONE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_class": "tair.rdb.2g",
					"node_type":      "MASTER_SLAVE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_class": "tair.rdb.2g",
						"node_type":      "MASTER_SLAVE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
			},
		},
	})
}

// Case Tair_rdb_升级大版本 6473  raw
func TestAccAliCloudRedisTairInstance_basic6473_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap6473)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence6473)
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
					"payment_type":          "PayAsYouGo",
					"instance_type":         "tair_rdb",
					"zone_id":               "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_class":        "tair.rdb.with.proxy.1g",
					"tair_instance_name":    name,
					"vswitch_id":            "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":                "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":     "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":              "123456Tf",
					"engine_version":        "5.0",
					"period":                "12",
					"port":                  "6379",
					"auto_renew":            "false",
					"slave_read_only_count": "2",
					"secondary_zone_id":     "${var.secondary_zone_id}",
					"shard_count":           "1",
					"read_only_count":       "2",
					"node_type":             "MASTER_SLAVE",
					"ssl_enabled":           "Disable",
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_rdb",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.with.proxy.1g",
						"tair_instance_name": name,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
						"password":           "123456Tf",
						"engine_version":     "5.0",
						"period":             "12",
						"port":               "6379",
						"auto_renew":         "false",
						"node_type":          "MASTER_SLAVE",
						"ssl_enabled":        "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ssl_enabled": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ssl_enabled": "Enable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"engine_version": "6.0",
					"force_upgrade":  "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"engine_version": "6.0",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap6823 = map[string]string{
	"port":            CHECKSET,
	"storage_size_gb": CHECKSET,
	"status":          CHECKSET,
	"engine_version":  CHECKSET,
	"node_type":       CHECKSET,
	"payment_type":    "PayAsYouGo",
	"create_time":     CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence6823(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

variable "region_id" {
  default = "cn-beijing"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_resource_manager_resource_group" "defaultRg" {
  display_name = var.name
  resource_group_name = var.name
}

resource "alicloud_vpc" "defaultVpc" {
  description       = "测试请勿绑定test-ljt"
  resource_group_id = alicloud_resource_manager_resource_group.defaultRg.id
  cidr_block        = "172.16.0.0/12"
  vpc_name          = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_security_group" "defaultEcsSg" {
  name = var.name
  vpc_id              = alicloud_vpc.defaultVpc.id
  resource_group_id   = alicloud_resource_manager_resource_group.defaultRg.id
  security_group_type = "normal"
}

resource "alicloud_security_group" "defaultsg2" {
  name = var.name
  vpc_id              = alicloud_vpc.defaultVpc.id
  resource_group_id   = alicloud_resource_manager_resource_group.defaultRg.id
  security_group_type = "normal"
}


`, name)
}

// Case Tair_scm_白名单安全组 6823  raw
func TestAccAliCloudRedisTairInstance_basic6823_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap6823)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence6823)
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
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_scm",
					"zone_id":            "${var.zone_id}",
					"instance_class":     "tair.scm.standard.1m.4d",
					"tair_instance_name": name,
					"vswitch_id":         "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":             "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":  "${alicloud_resource_manager_resource_group.defaultRg.id}",
					"password":           "123456Tf",
					"engine_version":     "1.0",
					"period":             "12",
					"port":               "6379",
					"secondary_zone_id":  "${var.zone_id}",
					"auto_renew":         "false",
					"security_group_id":  "${alicloud_security_group.defaultEcsSg.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_scm",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.scm.standard.1m.4d",
						"tair_instance_name": name,
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
						"password":           "123456Tf",
						"engine_version":     "1.0",
						"period":             "12",
						"port":               "6379",
						"secondary_zone_id":  CHECKSET,
						"auto_renew":         "false",
						"security_group_id":  CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.defaultsg2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "cluster_backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
			},
		},
	})
}
