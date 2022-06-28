package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVideoSurveillanceSystemGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_video_surveillance_system_group.default"
	ra := resourceAttrInit(resourceId, AlicloudVideoSurveillanceSystemGroupMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVideoSurveillanceSystemGroup")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testacc%sVsGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVideoSurveillanceSystemGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SurveillanceSystemSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":   name,
					"in_protocol":  "rtmp",
					"out_protocol": "flv",
					"play_domain":  "ultron.pub",
					"push_domain":  "lyzhuoan.cn",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":   name,
						"in_protocol":  "rtmp",
						"out_protocol": "flv",
						"play_domain":  "ultron.pub",
						"push_domain":  "lyzhuoan.cn",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"callback": "http://play.aliyunlive.com/notify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"callback": "http://play.aliyunlive.com/notify",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_desc",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_desc",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"in_protocol": "gb28181",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"in_protocol": "gb28181",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"out_protocol": "hls",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"out_protocol": "hls",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"group_name":   name + "updatall",
					"callback":     "http://play.aliyunlive.com/notify" + "update",
					"description":  name + "descAll",
					"enabled":      "false",
					"in_protocol":  "rtmp",
					"out_protocol": "flv",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_name":   name + "updatall",
						"callback":     "http://play.aliyunlive.com/notify" + "update",
						"description":  name + "descAll",
						"enabled":      "false",
						"in_protocol":  "rtmp",
						"out_protocol": "flv",
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

var AlicloudVideoSurveillanceSystemGroupMap0 = map[string]string{}

func AlicloudVideoSurveillanceSystemGroupBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
