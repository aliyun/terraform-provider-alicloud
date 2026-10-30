package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_redis_backup", &resource.Sweeper{
		Name: "alicloud_redis_backup",
		F:    testSweepRedisBackups,
	})
}

// testSweepRedisBackups sweeps KVStore (Redis) manual backups. The
// alicloud_redis_backup resource creates manual backups via CreateBackup, and
// the cloud only permits deleting manual backups via DeleteBackup (automatic
// backups are managed by the backup policy). Backups are a sub-resource of a
// KVStore instance, so we enumerate instances first (Redis + Memcache), then
// DescribeBackups per instance and delete each manual backup.
func testSweepRedisBackups(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
	}

	// 1. List all KVStore instances (Redis + Memcache) so we can enumerate backups per instance.
	listReq := map[string]interface{}{
		"RegionId":   client.RegionId,
		"PageSize":   PageSizeXLarge,
		"PageNumber": 1,
	}
	instanceIds := make([]string, 0)
	for _, instanceType := range []string{string(KVStoreRedis), string(KVStoreMemcache)} {
		listReq["InstanceType"] = instanceType
		for {
			response, err := client.RpcPost("R-kvstore", "2015-01-01", "DescribeInstances", nil, listReq, true)
			if err != nil {
				log.Printf("[ERROR] Failed to list KVStore instances for backup sweep: %s", err)
				return nil
			}
			resp, err := jsonpath.Get("$.Instances.KVStoreInstance", response)
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, "DescribeInstances", "$.Instances.KVStoreInstance", response)
			}
			result, _ := resp.([]interface{})
			for _, v := range result {
				item, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				instanceName := fmt.Sprint(item["InstanceName"])
				if protectedSweepVpcName(instanceName) {
					log.Printf("[INFO] Skipping KVStore Instance backups (NODELETING): %s (%s)", instanceName, item["InstanceId"])
					continue
				}
				if !sweepAll() {
					match := false
					for _, prefix := range prefixes {
						if strings.HasPrefix(strings.ToLower(instanceName), strings.ToLower(prefix)) {
							match = true
							break
						}
					}
					if !match {
						continue
					}
				}
				instanceIds = append(instanceIds, fmt.Sprint(item["InstanceId"]))
			}
			if len(result) < PageSizeXLarge {
				break
			}
			listReq["PageNumber"] = listReq["PageNumber"].(int) + 1
		}
	}

	// 2. For each qualifying instance, list backups and delete each manual backup.
	for _, instanceId := range instanceIds {
		backupReq := map[string]interface{}{
			"RegionId":   client.RegionId,
			"InstanceId":  instanceId,
			"StartTime":   "2010-01-01T10:00Z",
			"EndTime":     "2050-01-01T10:00Z",
			"PageSize":    PageSizeXLarge,
			"PageNumber":  1,
		}
		for {
			backupResp, err := client.RpcPost("R-kvstore", "2015-01-01", "DescribeBackups", nil, backupReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
					break
				}
				log.Printf("[ERROR] Failed to DescribeBackups for instance %s: %s", instanceId, err)
				break
			}
			backupObj, err := jsonpath.Get("$.Backups.Backup", backupResp)
			if err != nil {
				log.Printf("[ERROR] Failed to parse backups for instance %s: %s", instanceId, err)
				break
			}
			backupList, _ := backupObj.([]interface{})
			for _, b := range backupList {
				backup, ok := b.(map[string]interface{})
				if !ok {
					continue
				}
				// Only delete manual backups — DeleteBackup only supports
				// manually-backed-up backup sets; automatic backups are owned
				// by the backup policy and cannot be deleted here.
				if fmt.Sprint(backup["BackupMode"]) != "Manual" {
					continue
				}
				backupId := fmt.Sprint(backup["BackupId"])
				log.Printf("[INFO] Deleting KVStore Redis backup: %s:%s", instanceId, backupId)
				delReq := map[string]interface{}{
					"RegionId":   client.RegionId,
					"InstanceId": instanceId,
					"BackupId":   backupId,
				}
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(1*time.Minute, func() *resource.RetryError {
					_, err = client.RpcPost("R-kvstore", "2015-01-01", "DeleteBackup", nil, delReq, false)
					if err != nil {
						if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidBackupId.NotFound", "InvalidInstanceId.NotFound"}) {
							return nil
						}
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				if err != nil {
					log.Printf("[ERROR] Failed to delete KVStore Redis backup %s:%s: %s", instanceId, backupId, err)
				}
			}
			if len(backupList) < PageSizeXLarge {
				break
			}
			backupReq["PageNumber"] = backupReq["PageNumber"].(int) + 1
		}
	}
	return nil
}

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
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
				ImportStateVerifyIgnore: []string{"backup_retention_period"},
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

// Test Redis Backup. <<< Resource test cases, automatically generated.
