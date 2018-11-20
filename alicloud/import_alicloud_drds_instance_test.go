package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudDRDSInstance_importBasic(t *testing.T) {
	resourceName := "alicloud_drds_instance.basic"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDRDSInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDrdsInstance,
			},
			resource.TestStep{
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDRDSInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDrdsInstance_Vpc,
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}