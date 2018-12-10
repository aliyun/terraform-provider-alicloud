package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudKVStoreParameter_import(t *testing.T) {
	resourceName := "alicloud_kvstore_parameter.compat"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKVStoreParameterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKVStoreParameter_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
