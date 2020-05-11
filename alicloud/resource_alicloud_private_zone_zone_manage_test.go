package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPrivateZoneZoneManage_basic(t *testing.T) {
	var v pvtz.DescribeZoneInfoResponse
	resourceId := "alicloud_private_zone_zone_manage.default"
	ra := resourceAttrInit(resourceId, PrivateZoneZoneManageMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &PvtzService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribePrivateZoneZoneManage")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccPrivateZoneZoneManage%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, PrivateZoneZoneManageBasicdependence)
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
					"zone_name":     "demo.com",
					"proxy_pattern": "ZONE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"zone_name":     "demo.com",
						"proxy_pattern": "ZONE",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ResourceGroupId"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "en",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_pattern": "RECORD",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_pattern": "RECORD",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remark": "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"remark": "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_client_ip": "1.1.1.1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_client_ip": "1.1.1.1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"user_client_ip": "1.1.1.2",
					"proxy_pattern":  "ZONE",
					"lang":           "zh",
					"remark":         "lastupdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"user_client_ip": "1.1.1.2",
						"proxy_pattern":  "ZONE",
						"lang":           "zh",
						"remark":         "lastupdate",
					}),
				),
			},
		},
	})
}

var PrivateZoneZoneManageMap = map[string]string{}

func PrivateZoneZoneManageBasicdependence(name string) string {
	return ""
}
