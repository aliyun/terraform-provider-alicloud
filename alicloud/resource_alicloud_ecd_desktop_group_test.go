package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECDDesktopGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_desktop_group.default"
	checkoutSupportedRegions(t, true, connectivity.ECDSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECDDesktopGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdDesktopGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccecddesktopgroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDDesktopGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"allow_buffer_count": "1",
					"max_desktops_count": "1",
					"min_desktops_count": "1",
					"charge_type":        "PostPaid",
					"policy_group_id":    "pg-dmuionlvb0d3djeug",
					"office_site_id":     "cn-hangzhou+dir-7826873134",
					"bundle_id":          "bundle_eds_enterprise_office_2c4g_s8d5_win2019",
					"allow_auto_setup":   "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allow_buffer_count": "1",
						"max_desktops_count": "1",
						"min_desktops_count": "1",
						"charge_type":        "PostPaid",
						"policy_group_id":    CHECKSET,
						"office_site_id":     CHECKSET,
						"bundle_id":          CHECKSET,
						"allow_auto_setup":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"keep_duration": "6000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"keep_duration": "6000",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comments": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comments": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"desktop_group_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desktop_group_name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "default_init_desktop_count", "vpc_id", "period"},
			},
		},
	})
}

func TestAccAlicloudECDDesktopGroup_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecd_desktop_group.default"
	checkoutSupportedRegions(t, true, connectivity.ECDSupportRegions)
	ra := resourceAttrInit(resourceId, AlicloudECDDesktopGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdDesktopGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testaccecddesktopgroup%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDDesktopGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"desktop_group_name": "${var.name}",
					"allow_buffer_count": "1",
					"keep_duration":      "6000",
					"max_desktops_count": "1",
					"min_desktops_count": "1",
					"charge_type":        "PostPaid",
					"policy_group_id":    "pg-dmuionlvb0d3djeug",
					"office_site_id":     "cn-hangzhou+dir-7826873134",
					"comments":           "${var.name}",
					"own_type":           "0",
					"bundle_id":          "bundle_eds_enterprise_office_2c4g_s8d5_win2019",
					//"directory_id":      "tf-testAcc-FJgX3",
					//"scale_strategy_id": "tf-testAcc-ZAaKk",
					"allow_auto_setup": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"desktop_group_name": name,
						"allow_buffer_count": "1",
						"keep_duration":      "6000",
						"max_desktops_count": "1",
						"min_desktops_count": "1",
						"charge_type":        "PostPaid",
						"policy_group_id":    CHECKSET,
						"comments":           CHECKSET,
						"own_type":           "0",
						"office_site_id":     CHECKSET,
						"bundle_id":          CHECKSET,
						//"directory_id":      "tf-testAcc-FJgX3",
						//"scale_strategy_id": "tf-testAcc-ZAaKk",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_pay", "default_init_desktop_count", "vpc_id", "period"},
			},
		},
	})
}

var AlicloudECDDesktopGroupMap0 = map[string]string{
	"default_init_desktop_count": NOSET,
	"status":                     CHECKSET,
}

func AlicloudECDDesktopGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
