package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKVStoreRedisBackupPolicy_import(t *testing.T) {
	resourceName := "alicloud_kvstore_backup_policy.policy"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_classic(string(KVStoreRedis), redisInstanceClassForTest, string(KVStore4Dot0)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudKVStoreMemcacheBackupPolicy_import(t *testing.T) {
	resourceName := "alicloud_kvstore_backup_policy.policy"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreBackupPolicy_vpc(DatabaseCommonTestCase, string(KVStoreMemcache), memcacheInstanceClassForTest, string(KVStore2Dot8)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
