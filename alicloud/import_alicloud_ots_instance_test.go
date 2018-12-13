package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudOtsInstanceCapacity_import(t *testing.T) {
	resourceName := "alicloud_ots_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsInstance(string(OtsCapacity), acctest.RandIntRange(10000, 999999)),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlicloudOtsInstanceHighPerformance_import(t *testing.T) {
	resourceName := "alicloud_ots_instance.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOtsInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOtsInstance(string(OtsHighPerformance), acctest.RandIntRange(10000, 999999)),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
