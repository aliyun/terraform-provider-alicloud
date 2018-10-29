package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKVStoreBackupPolicy_import(t *testing.T) {
	resourceName := "alicloud_kvstore_backup_policy.policy"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreBackupPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreBackupPolicy_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
