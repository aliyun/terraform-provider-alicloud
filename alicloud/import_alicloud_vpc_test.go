package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudVpc_importBasic(t *testing.T) {
	resourceName := "alicloud_vpc.default"
	rand := acctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVpcConfigBasic(rand),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
