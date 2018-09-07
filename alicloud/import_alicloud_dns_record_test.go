package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDnsDomainRecord_importBasic(t *testing.T) {
	resourceName := "alicloud_dns_record.record"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsRecordConfig(acctest.RandInt()),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
