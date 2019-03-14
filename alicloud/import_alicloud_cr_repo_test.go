package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCRRepo_Import(t *testing.T) {
	resourceName := "alicloud_cr_repo.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCRRepoDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCRRepo_Basic(acctest.RandIntRange(1000, 9999)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
