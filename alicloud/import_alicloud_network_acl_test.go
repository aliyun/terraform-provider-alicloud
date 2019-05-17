package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudNetworkAcl_importBasic(t *testing.T) {
	resourceName := "alicloud_network_acl.default"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.NetworkAclSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAcl_create(acctest.RandIntRange(1000, 9999)),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
