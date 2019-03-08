package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsDomain_importBasic(t *testing.T) {
	resourceName := "alicloud_dns.dns"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsConfig(acctest.RandInt()),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
