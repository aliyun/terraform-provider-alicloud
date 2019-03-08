package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCSSwarm_importBasic(t *testing.T) {
	resourceName := "alicloud_cs_swarm.cs_vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_basic,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "cidr_block", "image_id", "password", "release_eip", "need_slb"},
			},
		},
	})
}

func TestAccAlicloudCSSwarm_importBasicNoSlb(t *testing.T) {
	resourceName := "alicloud_cs_swarm.cs_vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.SwarmSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwarmClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCSSwarm_basic_noslb,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name_prefix", "cidr_block", "image_id", "password", "release_eip", "need_slb"},
			},
		},
	})
}
