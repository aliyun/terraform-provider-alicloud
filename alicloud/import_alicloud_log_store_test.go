package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudLogStore_import(t *testing.T) {
	resourceName := "alicloud_log_store.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudLogStoreBasic(acctest.RandIntRange(10000, 999999)),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
