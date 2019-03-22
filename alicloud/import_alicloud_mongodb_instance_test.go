package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudMongoDBInstance_import(t *testing.T) {
	resourceName := "alicloud_mongodb_instance.foo"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudMongoDBInstance_import_config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
