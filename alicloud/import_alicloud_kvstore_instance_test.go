package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKVStoreRedisInstance_import(t *testing.T) {
	resourceName := "alicloud_kvstore_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccAlicloudKVStoreMemcacheInstance_import(t *testing.T) {
	resourceName := "alicloud_kvstore_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKVStoreInstance_vpc(KVStoreCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
