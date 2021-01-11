package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudBrainIndustrialPidOrganization_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_brain_industrial_pid_organization.default"
	ra := resourceAttrInit(resourceId, AlicloudBrainIndustrialPidOrganizationMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &Brain_industrialService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeBrainIndustrialPidOrganization")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudBrainIndustrialPidOrganizationBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.BrainIndustrialRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organization_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent_pid_organization_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"pid_organization_name": "tf-testAccUp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"pid_organization_name": "tf-testAccUp",
					}),
				),
			},
		},
	})
}

var AlicloudBrainIndustrialPidOrganizationMap = map[string]string{}

func AlicloudBrainIndustrialPidOrganizationBasicDependence(name string) string {
	return ""
}
