package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMaxcompute_basic(t *testing.T) {
	resourceId := "alicloud_maxcompute_project.default"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := "tf_testAccMCProject"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
			testAccPreCheckWithRegions(t, true, connectivity.MaxComputeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMaxcomputeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":       name,
						"specification_type": "OdpsStandard",
						"order_type":         "PayAsYouGo",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"specification_type"},
			},
		},
	})
}

const testAccMaxcomputeConfigBasic = `
resource "alicloud_maxcompute_project" "default"{
  project_name      = "tf_testAccWWQProject"
  specification_type = "OdpsStandard"
  order_type = "PayAsYouGo"
}
`

func TestAccAlicloudMaxcompute_multi(t *testing.T) {
	resourceId := "alicloud_maxcompute_project.default.4"
	ra := resourceAttrInit(resourceId, nil)
	testAccCheck := ra.resourceAttrMapUpdateSet()
	name := "tf_testAccMCProject"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			// Currently does not support creating projects with sub-accounts
			testAccPreCheckWithRegions(t, true, connectivity.MaxComputeSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMaxcomputeConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":       name + "4",
						"specification_type": "OdpsStandard",
						"order_type":         "PayAsYouGo",
					}),
				),
			},
		},
	})
}

const testAccMaxcomputeConfigMulti = `
resource "alicloud_maxcompute_project" "default"{
  count = "5"
  project_name      = "tf_testAccMCProject${count.index}"
  specification_type = "OdpsStandard"
  order_type = "PayAsYouGo"
}
`
