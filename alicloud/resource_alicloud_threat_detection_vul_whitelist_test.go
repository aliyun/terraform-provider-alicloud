package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudThreatDetectionVulWhitelist_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_vul_whitelist.default"
	ra := resourceAttrInit(resourceId, resourceAlicloudThreatDetectionVulWhitelistMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionVulWhitelist")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccThreatDetectionVulWhitelist-name%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlicloudThreatDetectionVulWhitelistBasicDependence)
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
					"whitelist":   `[{\"aliasName\":\"RHSA-2021:2260: libwebp 安全更新\",\"name\":\"RHSA-2021:2260: libwebp 安全更新\",\"type\":\"cve\"}]`,
					"target_info": `{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678]}`,
					"reason":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"whitelist":   "[{\"aliasName\":\"RHSA-2021:2260: libwebp 安全更新\",\"name\":\"RHSA-2021:2260: libwebp 安全更新\",\"type\":\"cve\"}]",
						"target_info": "{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678]}",
						"reason":      name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_info": `{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678,10782677]}`,
					"reason":      name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_info": "{\"type\":\"GroupId\",\"uuids\":[],\"groupIds\":[10782678,10782677]}",
						"reason":      name + "-update",
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

var resourceAlicloudThreatDetectionVulWhitelistMap = map[string]string{}

func resourceAlicloudThreatDetectionVulWhitelistBasicDependence(name string) string {
	return ""
}
