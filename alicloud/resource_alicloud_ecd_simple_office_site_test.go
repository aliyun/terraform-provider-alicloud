package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECDSimpleOfficeSite_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_simple_office_site.default"
	ra := resourceAttrInit(resourceId, AlicloudECDSimpleOfficeSiteMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdSimpleOfficeSite")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdsimpleofficesite%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDSimpleOfficeSiteBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.EcdSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cidr_block": "172.16.0.0/12",
					// todo: need to check the `bandwidth` and `enable_internet_access` after fixing the issue occurred in ap-southeast-1
					//"enable_internet_access": "true",
					"enable_admin_access": "true",
					"office_site_name":    name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cidr_block":          "172.16.0.0/12",
						"enable_admin_access": "true",
						//"enable_internet_access": "true",
						"office_site_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sso_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sso_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mfa_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mfa_enabled": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_cross_desktop_access": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_cross_desktop_access": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"office_site_name":    name + "_update",
					"desktop_access_type": "Any",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"office_site_name":    name + "_update",
						"desktop_access_type": "Any",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"office_site_name":    name + "_update1",
					"desktop_access_type": "VPC",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"office_site_name":    name + "_update1",
						"desktop_access_type": "VPC",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sso_enabled":                 "false",
					"mfa_enabled":                 "false",
					"enable_cross_desktop_access": "true",
					"desktop_access_type":         "Internet",
					"office_site_name":            name + "_update2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sso_enabled":                 "false",
						"mfa_enabled":                 "false",
						"enable_cross_desktop_access": "true",
						"desktop_access_type":         "Internet",
						"office_site_name":            name + "_update2",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudECDSimpleOfficeSiteMap0 = map[string]string{
	"bandwidth": NOSET,
}

func AlicloudECDSimpleOfficeSiteBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
