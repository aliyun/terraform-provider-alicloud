package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVODDomain_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vod_domain.default"
	ra := resourceAttrInit(resourceId, AlicloudVODDomainMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VodService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVodDomain")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svoddomain%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVODDomainBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.VodSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"domain_name": "136.chat",
					"sources": []map[string]interface{}{
						{
							"source_type":    "oss",
							"source_content": "outin-c7405446108111ec9a7100163e0eb78b.oss-cn-beijing.aliyuncs.com",
							"source_port":    "80",
						},
					},
					"scope": "domestic",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": "136.chat",
						"scope":       "domestic",
						"sources.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"tftestacc":    "TFTEST",
						"Tftestacc123": "Tftest123",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":            "2",
						"tags.tftestacc":    "TFTEST",
						"tags.Tftestacc123": "Tftest123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sources": []map[string]interface{}{
						{
							"source_type":    "domain",
							"source_content": "outin-c7405446108111ec9a7100163e0eb78b.oss-cn-beijing.aliyuncs.com",
							"source_port":    "443",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain_name": "136.chat",
						"scope":       "domestic",
						"sources.#":   "1",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"top_level_domain", "tags", "check_url"},
			},
		},
	})
}

var AlicloudVODDomainMap0 = map[string]string{}

func AlicloudVODDomainBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
