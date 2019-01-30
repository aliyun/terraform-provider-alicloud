package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudOtsTableCapacity_import(t *testing.T) {
	resourceName := "alicloud_ots_table.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsCapacity), acctest.RandIntRange(10000, 999999)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudOtsTableHighPerformance_import(t *testing.T) {
	resourceName := "alicloud_ots_table.basic"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOtsTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOtsTableStore(string(OtsHighPerformance), acctest.RandIntRange(10000, 999999)),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
