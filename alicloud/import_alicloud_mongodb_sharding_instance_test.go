package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMongoDBShardingInstance_import(t *testing.T) {
	resourceName := "alicloud_mongodb_sharding_instance.default"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testMongoDBShardingInstance_vpc_base,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
