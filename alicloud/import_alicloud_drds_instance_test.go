package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDRDSInstance_importBasic(t *testing.T) {
	resourceName := "alicloud_drds_instance.basic"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions)
			testAccPreCheckWithRegions(t, false, connectivity.DrdsClassicNoSupportedRegions)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDRDSInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDrdsInstance,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccAlicloudDRDSInstance_importVpc(t *testing.T) {
	resourceName := "alicloud_drds_instance.vpc"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckWithRegions(t, true, connectivity.DrdsSupportedRegions) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDRDSInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDrdsInstance_Vpc,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}
