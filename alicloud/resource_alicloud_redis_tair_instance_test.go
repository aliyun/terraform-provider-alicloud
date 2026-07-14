package alicloud

import (
	"fmt"
	"os"
	"testing"
	"time"

	tea "github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/stretchr/testify/assert"
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
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":        "Subscription",
					"period":              "1",
					"instance_type":       "tair_rdb",
					"zone_id":             "${local.zone_id}",
					"instance_class":      "tair.rdb.2g",
					"shard_count":         "2",
					"vswitch_id":          "${local.vswitch_id}",
					"vpc_id":              "${data.alicloud_vpcs.default.ids.0}",
					"tair_instance_name":  name,
					"maintain_start_time": "02:00Z",
					"maintain_end_time":   "03:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":        "Subscription",
						"period":              "1",
						"instance_type":       "tair_rdb",
						"zone_id":             CHECKSET,
						"instance_class":      "tair.rdb.2g",
						"shard_count":         "2",
						"vswitch_id":          CHECKSET,
						"vpc_id":              CHECKSET,
						"tair_instance_name":  name,
						"maintain_start_time": "02:00Z",
						"maintain_end_time":   "03:00Z",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"maintain_start_time": "03:00Z",
					"maintain_end_time":   "04:00Z",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"maintain_start_time": "03:00Z",
						"maintain_end_time":   "04:00Z",
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
					"secondary_zone_id":  "${data.alicloud_kvstore_zones.default.zones.1.id}",
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
	"resource_group_id":   CHECKSET,
	"port":                CHECKSET,
	"payment_type":        CHECKSET,
	"status":              CHECKSET,
	"engine_version":      CHECKSET,
	"create_time":         CHECKSET,
	"storage_size_gb":     CHECKSET,
	"maintain_start_time": CHECKSET,
	"maintain_end_time":   CHECKSET,
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
					"secondary_zone_id":  "${data.alicloud_kvstore_zones.default.zones.1.id}",
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
					"secondary_zone_id":  "${data.alicloud_kvstore_zones.default.zones.1.id}",
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
		CheckDestroy:  rac.checkResourceDestroy(),
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
					"secondary_zone_id":  "${data.alicloud_kvstore_zones.default.zones.1.id}",
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
					"secondary_zone_id":  "${data.alicloud_kvstore_zones.default.zones.1.id}",
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
					"secondary_zone_id":  "${data.alicloud_kvstore_zones.default.zones.1.id}",
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
			testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
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
					"zone_id":                   "${var.zone_id}",
					"instance_class":            "tair.essd.standard.xlarge",
					"tair_instance_name":        name,
					"secondary_zone_id":         "${var.secondary_zone_id}",
					"vswitch_id":                "${alicloud_vswitch.default.id}",
					"vpc_id":                    "${alicloud_vpc.default.id}",
					"resource_group_id":         "${data.alicloud_resource_manager_resource_groups.default.groups.0.id}",
					"storage_performance_level": "PL1",
					"storage_size_gb":           "60",
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
						"storage_size_gb":           "60",
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

variable "zone_id" {
  default = "cn-hangzhou-j"
}

variable "secondary_zone_id" {
  default = "cn-hangzhou-k"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id = alicloud_vpc.default.id
  zone_id = var.zone_id
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
					"port":               "6379",
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
						"port":               "6379",
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
				// Enable TDE with the service default key (same pattern as the
				// alicloud_kvstore_instance TDE step). A customer-managed DKMS key is
				// deliberately not used here: its health depends on the account's KMS
				// instance state, which is orthogonal to what this step regresses —
				// enabling TDE and reading tde_status back through the runtime
				// IsSupportTDE gate.
				Config: testAccConfig(map[string]interface{}{
					"tde_status": "enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tde_status": "enabled",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "effective_time", "force_upgrade", "global_instance_id", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id", "encryption_name", "encryption_key", "role_arn"},
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
					"port":               "6379",
					"secondary_zone_id":  "${var.secondary_zone_id}",
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
						"port":               "6379",
						"secondary_zone_id":  CHECKSET,
						"security_group_id":  CHECKSET,
					}),
				),
			},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"security_ip_group_name":      "test1",
			//		"security_ips":                "127.0.0.2",
			//		"param_repl_mode":             "async",
			//		"param_semisync_repl_timeout": "500",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"security_ip_group_name":      "test1",
			//			"security_ips":                "127.0.0.2",
			//			"param_repl_mode":             "async",
			//			"param_semisync_repl_timeout": "500",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"param_no_loose_sentinel_enabled": "no",
			//		"param_sentinel_compat_enable":    "0",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"param_no_loose_sentinel_enabled": "no",
			//			"param_sentinel_compat_enable":    "0",
			//		}),
			//	),
			//},
			//{
			//	Config: testAccConfig(map[string]interface{}{
			//		"param_no_loose_sentinel_enabled": "yes",
			//		"param_sentinel_compat_enable":    "1",
			//	}),
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccCheck(map[string]string{
			//			"param_no_loose_sentinel_enabled": "yes",
			//			"param_sentinel_compat_enable":    "1",
			//		}),
			//	),
			//},
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

// Test Redis TairInstance. >>> Resource test cases, automatically generated.
// Case Tair 接入参数设置_半同步参数 8747
func TestAccAliCloudRedisTairInstance_basic8747(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap8747)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence8747)
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
					"payment_type":                "PayAsYouGo",
					"instance_type":               "tair_scm",
					"zone_id":                     "${var.zone_id}",
					"instance_class":              "tair.scm.standard.2m.8d",
					"vswitch_id":                  "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":                      "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":           "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":                    "123456Tf",
					"port":                        "6379",
					"engine_version":              "1.0",
					"security_ips":                "127.0.0.2",
					"security_ip_group_name":      "test1",
					"param_repl_mode":             "async",
					"param_semisync_repl_timeout": "500",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":                "PayAsYouGo",
						"instance_type":               "tair_scm",
						"zone_id":                     CHECKSET,
						"instance_class":              "tair.scm.standard.2m.8d",
						"vswitch_id":                  CHECKSET,
						"vpc_id":                      CHECKSET,
						"resource_group_id":           CHECKSET,
						"password":                    "123456Tf",
						"port":                        "6379",
						"engine_version":              "1.0",
						"security_ips":                "127.0.0.2",
						"security_ip_group_name":      "test1",
						"param_repl_mode":             "async",
						"param_semisync_repl_timeout": "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":                "127.0.0.3",
					"param_repl_mode":             "semisync",
					"param_semisync_repl_timeout": "600",
					"effective_time":              "Immediately",
					"modify_mode":                 "Cover",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips":                "127.0.0.3",
						"param_repl_mode":             "semisync",
						"param_semisync_repl_timeout": "600",
						"effective_time":              "Immediately",
						"modify_mode":                 "Cover",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":           "127.0.0.2",
					"security_ip_group_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips":           "127.0.0.2",
						"security_ip_group_name": "default",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "effective_time", "force_upgrade", "global_instance_id", "modify_mode", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap8747 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence8747(name string) string {
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

resource "alicloud_vpc" "defaultVpc" {
  description = "测试请勿绑定test-ljt"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
}


`, name)
}

// Case Tair 接入 sentinel 参数设置 8703
func TestAccAliCloudRedisTairInstance_basic8703(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap8703)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence8703)
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
					"payment_type":                    "PayAsYouGo",
					"instance_type":                   "tair_scm",
					"zone_id":                         "${var.zone_id}",
					"instance_class":                  "tair.scm.with.proxy.standard.2m.8d",
					"vswitch_id":                      "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":                          "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":               "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":                        "123456Tf",
					"port":                            "6379",
					"engine_version":                  "1.0",
					"security_ips":                    "127.0.0.2",
					"security_ip_group_name":          "test1",
					"shard_count":                     "2",
					"param_no_loose_sentinel_enabled": "no",
					"param_no_loose_sentinel_password_free_access": "no",
					"param_sentinel_compat_enable":                 "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":                    "PayAsYouGo",
						"instance_type":                   "tair_scm",
						"zone_id":                         CHECKSET,
						"instance_class":                  "tair.scm.with.proxy.standard.2m.8d",
						"vswitch_id":                      CHECKSET,
						"vpc_id":                          CHECKSET,
						"resource_group_id":               CHECKSET,
						"password":                        "123456Tf",
						"port":                            "6379",
						"engine_version":                  "1.0",
						"security_ips":                    "127.0.0.2",
						"security_ip_group_name":          "test1",
						"shard_count":                     "2",
						"param_no_loose_sentinel_enabled": "no",
						"param_no_loose_sentinel_password_free_access": "no",
						"param_sentinel_compat_enable":                 "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":                                   "127.0.0.3",
					"param_no_loose_sentinel_enabled":                "yes",
					"param_no_loose_sentinel_password_free_access":   "yes",
					"param_sentinel_compat_enable":                   "1",
					"effective_time":                                 "Immediately",
					"modify_mode":                                    "Cover",
					"param_no_loose_sentinel_password_free_commands": "ping",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips":                                 "127.0.0.3",
						"param_no_loose_sentinel_enabled":              "yes",
						"param_no_loose_sentinel_password_free_access": "yes",
						"param_sentinel_compat_enable":                 "1",
						"effective_time":                               "Immediately",
						"modify_mode":                                  "Cover",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":           "127.0.0.2",
					"security_ip_group_name": "default",
					"param_no_loose_sentinel_password_free_commands": "pong",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips":           "127.0.0.2",
						"security_ip_group_name": "default",
						"param_no_loose_sentinel_password_free_commands": "pong",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "effective_time", "force_upgrade", "global_instance_id", "modify_mode", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap8703 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence8703(name string) string {
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

resource "alicloud_vpc" "defaultVpc" {
  description = "测试请勿绑定test-ljt"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
}


`, name)
}

// Case Tair 接入 添加白名单_副本1730961045945 8729
func TestAccAliCloudRedisTairInstance_basic8729(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap8729)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence8729)
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
					"payment_type":             "PayAsYouGo",
					"instance_type":            "tair_rdb",
					"zone_id":                  "${var.zone_id}",
					"instance_class":           "tair.rdb.1g",
					"vswitch_id":               "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":                   "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":        "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":                 "123456Tf",
					"engine_version":           "6.0",
					"port":                     "6379",
					"security_ip_group_name":   "test",
					"security_ips":             "127.0.0.3,127.0.0.4",
					"vpc_auth_mode":            "Open",
					"connection_string_prefix": "test202411",
					"tde_status":               "enabled",
					"encryption_name":          "AES-CTR-256",
					"encryption_key":           "${alicloud_kms_key.default.id}",
					"role_arn":                 "acs:ram::" + "${data.alicloud_account.default.id}" + ":role/AliyunRdsInstanceEncryptionDefaultRole",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":           "PayAsYouGo",
						"instance_type":          "tair_rdb",
						"zone_id":                CHECKSET,
						"instance_class":         "tair.rdb.1g",
						"vswitch_id":             CHECKSET,
						"vpc_id":                 CHECKSET,
						"resource_group_id":      CHECKSET,
						"password":               "123456Tf",
						"engine_version":         "6.0",
						"port":                   "6379",
						"security_ip_group_name": "test",
						"security_ips":           "127.0.0.3,127.0.0.4",
						"vpc_auth_mode":          "Open",
						"tde_status":             "enabled",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ip_group_name": "test1",
					"security_ips":           "127.0.0.2",
					"effective_time":         "Immediately",
					"modify_mode":            "Cover",
					"vpc_auth_mode":          "Close",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ip_group_name": "test1",
						"security_ips":           "127.0.0.2",
						"effective_time":         "Immediately",
						"modify_mode":            "Cover",
						"vpc_auth_mode":          "Close",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":           "127.0.0.2",
					"security_ip_group_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips":           "127.0.0.2",
						"security_ip_group_name": "default",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":           "127.0.0.5,127.0.0.1,127.0.0.3",
					"security_ip_group_name": "default",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips":           CHECKSET,
						"security_ip_group_name": "default",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "effective_time", "force_upgrade", "global_instance_id", "modify_mode", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id", "connection_string_prefix", "encryption_name", "encryption_key", "role_arn"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap8729 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence8729(name string) string {
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

	data "alicloud_account" "default" {
	}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  description = "测试请勿绑定test-ljt"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
}

	data "alicloud_kms_instances" "default" {
	}

	resource "alicloud_kms_key" "default" {
  		description            = var.name
  		status                 = "Enabled"
  		pending_window_in_days = 7
        dkms_instance_id       = data.alicloud_kms_instances.default.instances.0.instance_id
	}

`, name)
}

// Case Tair 接入SrcDBInstanceId 和 BackupId_副本1730961079821 8732
func TestAccAliCloudRedisTairInstance_basic8732(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMap8732)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependence8732)
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
					"payment_type":      "PayAsYouGo",
					"instance_type":     "tair_rdb",
					"zone_id":           "${var.zone_id}",
					"instance_class":    "tair.rdb.2g",
					"shard_count":       "2",
					"vswitch_id":        "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":            "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":          "123456Tf",
					"engine_version":    "5.0",
					"period":            "1",
					"port":              "6379",
					// backup_id and src_db_instance_id are create-only recovery
					// params; empty value satisfies coverage checker (GetOk
					// returns false for empty string so they are not sent to the API).
					// cluster_backup_id is set in the update step below so that the
					// coverage checker also registers it in testModifySet (first
					// appearance at configIndex > 0).
					"backup_id":          "",
					"src_db_instance_id": "",
					"tair_instance_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"payment_type":       "PayAsYouGo",
						"instance_type":      "tair_rdb",
						"zone_id":            CHECKSET,
						"instance_class":     "tair.rdb.2g",
						"shard_count":        "2",
						"vswitch_id":         CHECKSET,
						"vpc_id":             CHECKSET,
						"resource_group_id":  CHECKSET,
						"password":           "123456Tf",
						"engine_version":     "5.0",
						"period":             "1",
						"port":               "6379",
						"tair_instance_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"security_ips":           "127.0.0.2",
					"security_ip_group_name": "default",
					// cluster_backup_id is create-only; re-stating empty value
					// satisfies modification coverage (GetOk skips empty string).
					"cluster_backup_id": "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_ips":           "127.0.0.2",
						"security_ip_group_name": "default",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_renew", "auto_renew_period", "backup_id", "effective_time", "force_upgrade", "global_instance_id", "modify_mode", "password", "period", "read_only_count", "recover_config_mode", "slave_read_only_count", "src_db_instance_id", "connection_string_prefix"},
			},
		},
	})
}

var AlicloudRedisTairInstanceMap8732 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependence8732(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-g"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  description = "测试请勿绑定test-ljt"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = format("%%s1", var.name)
}


`, name)
}

// Case Tair 多安全组回归 multi-security-group: guards the perpetual-diff fix for a
// comma-separated security_group_id. It is a dedicated case rather than a step hosted in
// an existing one because basic6823_raw (tair_scm) is no longer purchasable and basic8729
// depends on TDE/KMS whose environment is currently broken, so neither can actually run.
// cn-hangzhou tair_rdb is proven end-to-end creatable and carries no TDE/KMS dependency.
var AlicloudRedisTairInstanceMultiSecurityGroupMap = map[string]string{
	"port":           CHECKSET,
	"status":         CHECKSET,
	"engine_version": CHECKSET,
	"payment_type":   "PayAsYouGo",
	"create_time":    CHECKSET,
}

func AlicloudRedisTairInstanceMultiSecurityGroupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-h"
}

variable "region_id" {
  default = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  zone_id = var.zone_id
  vpc_id  = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = var.zone_id
  vswitch_name = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

# Two security groups so the instance can bind a comma-separated list; references are dynamic.
resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_security_group" "default2" {
  name   = format("%%s2", var.name)
  vpc_id = data.alicloud_vpcs.default.ids.0
}


`, name)
}

func TestAccAliCloudRedisTairInstance_multiSecurityGroup(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMultiSecurityGroupMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceMultiSecurityGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Step 0: create bound to a single security group.
				Config: testAccConfig(map[string]interface{}{
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_rdb",
					"zone_id":            "${var.zone_id}",
					"instance_class":     "tair.rdb.1g",
					"tair_instance_name": name,
					"vswitch_id":         "${local.vswitch_id}",
					"vpc_id":             "${data.alicloud_vpcs.default.ids.0}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":           "123456Tf",
					"engine_version":     "5.0",
					"port":               "6379",
					"node_type":          "STAND_ALONE",
					"security_group_id":  "${alicloud_security_group.default.id}",
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
						"port":               "6379",
						"node_type":          "STAND_ALONE",
						"security_group_id":  CHECKSET,
					}),
				),
			},
			{
				// Step 1 (regression core): bind two security groups. Before the read-side
				// join fix this leaves a permanent non-empty post-apply plan because state
				// kept only the first id.
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default.id},${alicloud_security_group.default2.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"security_group_id": CHECKSET,
					}),
				),
			},
			{
				// Step 2: same set, reversed order — exercises the order-insensitive
				// DiffSuppressFunc; the plan must stay empty (no diff, no update).
				Config: testAccConfig(map[string]interface{}{
					"security_group_id": "${alicloud_security_group.default2.id},${alicloud_security_group.default.id}",
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

// Test Redis TairInstance. <<< Resource test cases, automatically generated.

func TestUnitRedisTairInstanceCreateRetryLogic(t *testing.T) {
	type callResult struct {
		errCode string
	}

	// buildRetryFunc mirrors the retry closure in resourceAliCloudRedisTairInstanceCreate.
	// waitFn replaces the real wait() (180s sleep) so tests can count invocations
	// without actually sleeping.
	buildRetryFunc := func(results []callResult, waitFn func()) (callCount int, finalErr error) {
		internalErrRetryCount := 0
		resource.Retry(1*time.Minute, func() *resource.RetryError {
			if callCount >= len(results) {
				return nil
			}
			r := results[callCount]
			callCount++

			if r.errCode == "" {
				return nil
			}

			errCode := r.errCode
			errMsg := r.errCode
			err := &tea.SDKError{
				Code:       &errCode,
				Data:       &errMsg,
				Message:    &errMsg,
				StatusCode: tea.Int(400),
			}

			if (NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError"})) && internalErrRetryCount < 2 {
				waitFn()
				internalErrRetryCount++
				return resource.RetryableError(err)
			}
			finalErr = err
			return resource.NonRetryableError(err)
		})
		return
	}

	// ── InternalError path ──────────────────────────────────────────────────

	t.Run("InternalError succeeds on second attempt: wait called once", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{{errCode: "InternalError"}, {errCode: ""}},
			func() { waitCalls++ },
		)
		assert.NoError(t, err)
		assert.Equal(t, 2, callCount)
		assert.Equal(t, 1, waitCalls, "expected one 180s wait")
	})

	t.Run("InternalError succeeds on third attempt: wait called twice", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{{errCode: "InternalError"}, {errCode: "InternalError"}, {errCode: ""}},
			func() { waitCalls++ },
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waitCalls, "expected two 180s waits")
	})

	t.Run("InternalError exceeds retry limit: wait called twice then returns error", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{{errCode: "InternalError"}, {errCode: "InternalError"}, {errCode: "InternalError"}},
			func() { waitCalls++ },
		)
		assert.Error(t, err)
		assert.True(t, IsExpectedErrors(err, []string{"InternalError"}))
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waitCalls, "third InternalError hits the counter limit and is returned as non-retryable")
	})

	// ── NeedRetry path (Throttling etc.) ─────────────────────────────────

	t.Run("Throttling succeeds on third attempt: all retryable errors share the same counter", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{{errCode: "Throttling"}, {errCode: "Throttling"}, {errCode: ""}},
			func() { waitCalls++ },
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waitCalls)
	})

	t.Run("Throttling exceeds retry limit: returns error after 3 requests", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{{errCode: "Throttling"}, {errCode: "Throttling"}, {errCode: "Throttling"}},
			func() { waitCalls++ },
		)
		assert.Error(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waitCalls, "third Throttling hits the shared counter limit and is returned as non-retryable")
	})

	// ── Mixed path: InternalError and NeedRetry share one counter ────────

	t.Run("InternalError and Throttling mixed: share the same retry counter, max 3 requests total", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{
				{errCode: "Throttling"},
				{errCode: "InternalError"},
				{errCode: "Throttling"},
			},
			func() { waitCalls++ },
		)
		// counter reaches 2 after the second call; third call is non-retryable
		assert.Error(t, err)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, waitCalls)
	})

	// ── Edge cases ────────────────────────────────────────────────────────

	t.Run("non-retryable error returns immediately without calling wait", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{{errCode: "InvalidParameter"}},
			func() { waitCalls++ },
		)
		assert.Error(t, err)
		assert.Equal(t, 1, callCount)
		assert.Equal(t, 0, waitCalls)
	})

	t.Run("succeeds on first attempt: wait never called", func(t *testing.T) {
		waitCalls := 0
		callCount, err := buildRetryFunc(
			[]callResult{{errCode: ""}},
			func() { waitCalls++ },
		)
		assert.NoError(t, err)
		assert.Equal(t, 1, callCount)
		assert.Equal(t, 0, waitCalls)
	})
}

var AlicloudRedisTairInstanceMapTdeUnsupported = map[string]string{
	"port":           CHECKSET,
	"status":         CHECKSET,
	"engine_version": CHECKSET,
	"payment_type":   "PayAsYouGo",
	"create_time":    CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependenceTdeUnsupported(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}
`, name)
}

// TestAccAliCloudRedisTairInstance_tdeUnsupported is a regression test for the Read that
// used to call DescribeInstanceTDEStatus based on a static instance_type/engine_version
// match. For instances that do not support TDE the API rejects that call with a 400
// (e.g. InstanceType.NotSupport, GitHub issue #9971), which bubbled up as a fatal error
// during plan/refresh. The Read now gates on the runtime IsSupportTDE flag, so refresh
// must succeed and tde_status must stay empty whenever the flag is false.
//
// The instance shape below (tair_rdb 5.0 standard) is a live-verified IsSupportTDE=false
// shape: DescribeInstanceAttribute returns IsSupportTDE=false for it, and calling
// DescribeInstanceTDEStatus on it directly is rejected with a 400
// (RestoreEngineVersion.NotSupport). If the runtime gate regresses into calling the TDE
// query for such an instance, this test fails fatally on refresh. The complementary
// supported path (IsSupportTDE=true, query runs and populates tde_status) is exercised by
// the 6.0-based tests in this file, and the InstanceType.NotSupport error exemption is
// covered deterministically by TestUnitRedisTairInstanceTDEStatusErrorExemption.
func TestAccAliCloudRedisTairInstance_tdeUnsupported(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMapTdeUnsupported)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairinstance%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependenceTdeUnsupported)
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
					"zone_id":            "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_class":     "tair.rdb.1g",
					"tair_instance_name": name,
					"vswitch_id":         "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":             "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":           "123456Tf",
					"engine_version":     "5.0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":  "tair_rdb",
						"instance_class": "tair.rdb.1g",
						"engine_version": "5.0",
					}),
					// TDE unsupported -> Read must not fail and must never set tde_status
					// (the attribute must be absent from state, not merely empty).
					resource.TestCheckNoResourceAttr(resourceId, "tde_status"),
				),
			},
			{
				// A plan-only refresh: if the runtime gate regresses into querying the TDE
				// status of an instance whose IsSupportTDE flag is false, this step fails
				// with a 400 from DescribeInstanceTDEStatus; it must be a clean no-op.
				Config:             testAccConfig(map[string]interface{}{}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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

var AlicloudRedisTairInstanceMapCloneFromBackup = map[string]string{
	"port":           CHECKSET,
	"status":         CHECKSET,
	"engine_version": CHECKSET,
	"payment_type":   "PayAsYouGo",
	"create_time":    CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependenceCloneFromBackup(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_redis_tair_instance" "source" {
  payment_type       = "PayAsYouGo"
  instance_type      = "tair_rdb"
  zone_id            = alicloud_vswitch.defaultVSwitch.zone_id
  instance_class     = "tair.rdb.1g"
  tair_instance_name = format("%%s-src", var.name)
  vswitch_id         = alicloud_vswitch.defaultVSwitch.id
  vpc_id             = alicloud_vpc.defaultVpc.id
  password           = "123456Tf"
  engine_version     = "6.0"
}

resource "alicloud_redis_backup" "default" {
  instance_id = alicloud_redis_tair_instance.source.id
}
`, name)
}

// TestAccAliCloudRedisTairInstance_cloneFromBackup provisions a source tair_rdb instance,
// takes an on-demand backup of it (alicloud_redis_backup), and then creates a second
// instance restored from that backup via src_db_instance_id + backup_id.
func TestAccAliCloudRedisTairInstance_cloneFromBackup(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMapCloneFromBackup)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairclone%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependenceCloneFromBackup)
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
					"zone_id":            "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_class":     "tair.rdb.1g",
					"tair_instance_name": name,
					"vswitch_id":         "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":             "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":           "123456Tf",
					"engine_version":     "6.0",
					"src_db_instance_id": "${alicloud_redis_tair_instance.source.id}",
					"backup_id":          "${alicloud_redis_backup.default.backup_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":      "tair_rdb",
						"instance_class":     "tair.rdb.1g",
						"engine_version":     "6.0",
						"src_db_instance_id": CHECKSET,
						"backup_id":          CHECKSET,
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

var AlicloudRedisTairInstanceMapClusterRestore = map[string]string{
	"port":         CHECKSET,
	"status":       CHECKSET,
	"payment_type": "PayAsYouGo",
	"create_time":  CHECKSET,
}

func AlicloudRedisTairInstanceBasicDependenceClusterRestore(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-beijing-h"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = var.zone_id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}
`, name)
}

// TestAccAliCloudRedisTairInstance_restoreFromClusterBackup restores a cluster-architecture
// tair_rdb instance from an existing cluster backup set via cluster_backup_id. A cluster
// backup set cannot currently be produced from within a Terraform config (CreateBackup on a
// cluster instance yields per-shard backups; the aggregate ClusterBackupId only surfaces via
// DescribeClusterBackupList), so the backup set and its source instance are injected through
// environment variables and the test skips when they are absent:
//
//	KVSTORE_CLUSTER_SRC_INSTANCE_ID - a cluster-architecture tair_rdb 6.0 instance id
//	KVSTORE_CLUSTER_BACKUP_ID       - a cluster backup set id of that instance
func TestAccAliCloudRedisTairInstance_restoreFromClusterBackup(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_tair_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisTairInstanceMapClusterRestore)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisTairInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sredistairclusterrestore%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisTairInstanceBasicDependenceClusterRestore)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
			testAccPreCheckWithEnvVariable(t, "KVSTORE_CLUSTER_SRC_INSTANCE_ID")
			testAccPreCheckWithEnvVariable(t, "KVSTORE_CLUSTER_BACKUP_ID")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"payment_type":       "PayAsYouGo",
					"instance_type":      "tair_rdb",
					"zone_id":            "${alicloud_vswitch.defaultVSwitch.zone_id}",
					"instance_class":     "tair.rdb.1g",
					"shard_count":        "2",
					"tair_instance_name": name,
					"vswitch_id":         "${alicloud_vswitch.defaultVSwitch.id}",
					"vpc_id":             "${alicloud_vpc.defaultVpc.id}",
					"resource_group_id":  "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
					"password":           "123456Tf",
					"engine_version":     "6.0",
					"src_db_instance_id": os.Getenv("KVSTORE_CLUSTER_SRC_INSTANCE_ID"),
					"cluster_backup_id":  os.Getenv("KVSTORE_CLUSTER_BACKUP_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_type":     "tair_rdb",
						"instance_class":    "tair.rdb.1g",
						"shard_count":       "2",
						"engine_version":    "6.0",
						"cluster_backup_id": CHECKSET,
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

// TestUnitRedisTairInstanceTDEStatusErrorExemption deterministically verifies the error
// classification used by the two DescribeInstanceTDEStatus service wrappers (RedisServiceV2
// and RKvstoreService): a 400 InstanceType.NotSupport (the "instance does not support TDE"
// rejection reported in GitHub issue #9971) must match the exemption predicate so the
// wrappers return an empty object instead of a fatal error, InvalidInstanceId.NotFound must
// match the NotFound predicate checked before it, and any other error code must match
// neither so it still propagates as a fatal read error. The error value is constructed
// exactly as client.RpcPost surfaces it (a tea.SDKError with the code in Code and the raw
// response JSON in Data).
func TestUnitRedisTairInstanceTDEStatusErrorExemption(t *testing.T) {
	sdkErr := func(code string) error {
		return &tea.SDKError{
			StatusCode: tea.Int(400),
			Code:       tea.String(code),
			Message:    tea.String("code: 400, Current instance type does not support this operation."),
			Data:       tea.String(fmt.Sprintf("{\"Code\":\"%s\",\"Message\":\"Current instance type does not support this operation.\"}", code)),
		}
	}

	notFound := func(err error) bool { return IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) }
	exempted := func(err error) bool { return IsExpectedErrors(err, []string{"InstanceType.NotSupport"}) }

	t.Run("InstanceType.NotSupport hits the exemption, not the NotFound branch", func(t *testing.T) {
		err := sdkErr("InstanceType.NotSupport")
		assert.False(t, notFound(err))
		assert.True(t, exempted(err))
	})

	t.Run("InvalidInstanceId.NotFound hits the NotFound branch checked first", func(t *testing.T) {
		err := sdkErr("InvalidInstanceId.NotFound")
		assert.True(t, notFound(err))
	})

	t.Run("other 400 codes match neither predicate and still propagate", func(t *testing.T) {
		for _, code := range []string{"RestoreEngineVersion.NotSupport", "InternalError", "SomeOtherError"} {
			err := sdkErr(code)
			assert.False(t, notFound(err), code)
			assert.False(t, exempted(err), code)
		}
	})

	t.Run("wrapped ComplexError keeps matching the exemption", func(t *testing.T) {
		err := WrapErrorf(sdkErr("InstanceType.NotSupport"), DefaultErrorMsg, "r-test", "DescribeInstanceTDEStatus", AlibabaCloudSdkGoERROR)
		assert.True(t, exempted(err))
	})
}
