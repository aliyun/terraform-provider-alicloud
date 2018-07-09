package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudLogStoreIndex_importFull(t *testing.T) {
	resourceName := "alicloud_log_store_index.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreIndexDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudLogStoreIndexFullText,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudLogStoreIndex_importField(t *testing.T) {
	resourceName := "alicloud_log_store_index.bar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudLogStoreIndexDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAlicloudLogStoreIndexField,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
