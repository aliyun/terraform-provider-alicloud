package alicloud

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenInstance_importBasic(t *testing.T) {
	resourceName := "alicloud_cen_instance.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceConfig(rand),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
