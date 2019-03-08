package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSecurityGroup_importBasic(t *testing.T) {
	resourceName := "alicloud_security_group.foo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.EcsClassicSupportedRegions)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudSecurityGroup_importWithVpc(t *testing.T) {
	resourceName := "alicloud_security_group.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupConfig_withVpc,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
