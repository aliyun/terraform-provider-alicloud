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
			resource.TestStep{
				Config: testAccKVStoreInstance_vpc(DatabaseCommonTestCase, redisInstanceClassForTest, string(KVStoreRedis), string(KVStore4Dot0)),
			},

			resource.TestStep{
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
			resource.TestStep{
				Config: testAccKVStoreInstance_vpc(DatabaseCommonTestCase, memcacheInstanceClassForTest, string(KVStoreMemcache), string(KVStore2Dot8)),
			},

			resource.TestStep{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
