package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECDCommand_basic0(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.EcdUserSupportRegions)
	resourceId := "alicloud_ecd_command.default"
	ra := resourceAttrInit(resourceId, AlicloudECDCommandMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcdService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcdCommand")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secdcommand%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECDCommandBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"command_content":  "ipconfig",
					"command_type":     "RunBatScript",
					"desktop_id":       "${alicloud_ecd_desktop.default.id}",
					"timeout":          "3600",
					"content_encoding": "PlainText",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"command_content":  "ipconfig",
						"command_type":     "RunBatScript",
						"desktop_id":       CHECKSET,
						"timeout":          "3600",
						"content_encoding": "PlainText",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeout", "content_encoding"},
			},
		},
	})
}

var AlicloudECDCommandMap0 = map[string]string{
	"content_encoding": NOSET,
	"desktop_id":       CHECKSET,
}

func AlicloudECDCommandBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

data "alicloud_ecd_bundles" "default"{
	bundle_type = "SYSTEM"
	name_regex  = "windows"
}
resource "alicloud_ecd_simple_office_site" "default" {
  cidr_block = "172.16.0.0/12"
  desktop_access_type = "Internet"
  office_site_name    = var.name
}
resource "alicloud_ecd_policy_group" "default" {
  policy_group_name = var.name
  clipboard = "readwrite"
  local_drive = "read"
  authorize_access_policy_rules{
    description= var.name
    cidr_ip=     "1.2.3.4/24"
  }
  authorize_security_policy_rules  {
    type=        "inflow"
    policy=      "accept"
    description=  var.name
    port_range= "80/80"
    ip_protocol= "TCP"
    priority=    "1"
    cidr_ip=     "0.0.0.0/0"
  }
}
resource "alicloud_ecd_desktop" "default" {
	office_site_id  = alicloud_ecd_simple_office_site.default.id
	policy_group_id = alicloud_ecd_policy_group.default.id
	bundle_id 		= data.alicloud_ecd_bundles.default.bundles.0.id
	desktop_name 	= var.name
}

`, name)
}
