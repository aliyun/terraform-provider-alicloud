package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Test Redis Backup. >>> Resource test cases, automatically generated.
// Case back测试3 11970
func TestAccAliCloudRedisBackup_basic11970(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_backup.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisBackupMap11970)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredis%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisBackupBasicDependence11970)
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
					"backup_retention_period": "7",
					"instance_id":             "${alicloud_kvstore_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "7",
						"instance_id":             CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"backup_retention_period", "cluster_backup_id"},
			},
		},
	})
}

var AlicloudRedisBackupMap11970 = map[string]string{
	"status":    CHECKSET,
	"backup_id": CHECKSET,
}

func AlicloudRedisBackupBasicDependence11970(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-h"
}

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

resource "alicloud_kvstore_instance" "default" {
  port           = "6379"
  payment_type   = "PrePaid"
  instance_type  = "Redis"
  password       = "123456_tf"
  engine_version = "5.0"
  zone_id        = var.zone_id
  vswitch_id     = local.vswitch_id
  period         = "1"
  instance_class = "redis.shard.small.2.ce"
}


`, name)
}

// Case cluster-delete: guards the fan-out Delete on cluster-architecture instances.
// Before the fix Delete only removed the state-recorded shard backup and left the other
// shards behind; the standard checkResourceDestroy would still pass because it only
// re-queries that one shard. This case captures instance_id + cluster_backup_id during
// Check, then in a custom CheckDestroy queries DescribeClusterBackupList for the captured
// cluster set and fails if any shard backup survives.
var AlicloudRedisBackupClusterDeleteMap = map[string]string{
	"status":    CHECKSET,
	"backup_id": CHECKSET,
}

func AlicloudRedisBackupClusterDeleteDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "zone_id" {
  default = "cn-hangzhou-h"
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

# Cloud-native cluster (shard_count >= 2) so the backup yields a cluster backup set (cb-*)
# with per-shard BackupIds.
resource "alicloud_redis_tair_instance" "cluster" {
  payment_type       = "PayAsYouGo"
  instance_type      = "tair_rdb"
  zone_id            = var.zone_id
  instance_class     = "tair.rdb.2g"
  shard_count        = 2
  tair_instance_name = var.name
  vswitch_id         = local.vswitch_id
  vpc_id             = data.alicloud_vpcs.default.ids.0
  resource_group_id  = data.alicloud_resource_manager_resource_groups.default.ids.0
  password           = "123456Tf"
  engine_version     = "5.0"
  port               = "6379"
}


`, name)
}

func TestAccAliCloudRedisBackup_clusterDelete(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_redis_backup.default"
	ra := resourceAttrInit(resourceId, AlicloudRedisBackupClusterDeleteMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RedisServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRedisBackup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccredisclusterdel%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRedisBackupClusterDeleteDependence)

	// Captured during Check so CheckDestroy can re-verify all shard backups are gone.
	// State-scoped access is not available at CheckDestroy time (the resource is already
	// removed from state), so we stash the two ids in the closure.
	var capturedInstanceID, capturedClusterBackupID string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy: func(s *terraform.State) error {
			if err := rac.checkResourceDestroy()(s); err != nil {
				return err
			}
			if capturedInstanceID == "" || capturedClusterBackupID == "" {
				// Cluster set never materialized; skip the extra assertion.
				return nil
			}
			client := testAccProvider.Meta().(*connectivity.AliyunClient)
			redisServiceV2 := &RedisServiceV2{client: client}
			shardIds, err := redisServiceV2.ListClusterBackupShardIds(capturedInstanceID, capturedClusterBackupID)
			if err != nil {
				if NotFoundError(err) {
					return nil
				}
				return fmt.Errorf("verify cluster backup fan-out delete for %s/%s: %v", capturedInstanceID, capturedClusterBackupID, err)
			}
			if len(shardIds) != 0 {
				return fmt.Errorf("cluster_backup_id %s still has %d shard backups after destroy: %v", capturedClusterBackupID, len(shardIds), shardIds)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"backup_retention_period": "7",
					"instance_id":             "${alicloud_redis_tair_instance.cluster.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backup_retention_period": "7",
						"instance_id":             CHECKSET,
					}),
					resource.TestCheckResourceAttrSet(resourceId, "cluster_backup_id"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceId]
						if !ok {
							return fmt.Errorf("resource %s not in state", resourceId)
						}
						capturedInstanceID = rs.Primary.Attributes["instance_id"]
						capturedClusterBackupID = rs.Primary.Attributes["cluster_backup_id"]
						return nil
					},
				),
			},
		},
	})
}

// Test Redis Backup. <<< Resource test cases, automatically generated.
