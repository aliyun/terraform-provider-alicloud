package threatdetection_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudThreatDetectionVulWhitelist_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_vul_whitelist.default"
	ra := acctest.ResourceAttrInit(resourceId, resourceAlicloudThreatDetectionVulWhitelistMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionVulWhitelist")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccThreatDetectionVulWhitelist-name%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlicloudThreatDetectionVulWhitelistBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
