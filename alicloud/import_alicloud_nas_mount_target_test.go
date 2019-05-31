package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudNasMountTarget_importBasic(t *testing.T) {
	resourceName := "alicloud_nas_mount_target.default"
	rand1 := acctest.RandIntRange(10000, 499999)
	rand2 := acctest.RandIntRange(500000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMountTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNasMountTargetVpcConfig(rand1, rand2),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
