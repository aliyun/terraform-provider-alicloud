package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDdosCooWebRule_basic(t *testing.T) {
	var v ddoscoo.WebRule
	resourceId := "alicloud_ddos_coo_web_rule.default"
	ra := resourceAttrInit(resourceId, DdosCooWebRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosCooWebRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccDdosCooWebRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DdosCooWebRuleBasicdependence)
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
					"domain":       "sojson.com",
					"rs_type":      "0",
					"rules":        `[{\"ProxyRules\":[{\"ProxyPort\":80,\"RealServers\":[\"1.1.1.1\"]}],\"ProxyType\":\"http\"},{\"ProxyRules\":[{\"ProxyPort\":443,\"RealServers\":[\"2.2.2.2\"]}],\"ProxyType\":\"https\"}]`,
					"real_servers": `[\"1.1.1.1\",\"2.2.2.2\"]`,
					"proxy_types":  `[{\"ProxyType\":\"http\",\"ProxyPorts\":[80]},{\"ProxyType\":\"https\",\"ProxyPorts\":[443]}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"domain":       "sojson.com",
						"rs_type":      "0",
						"rules":        CHECKSET,
						"real_servers": CHECKSET,
						"proxy_types":  CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rs_type", "real_servers", "proxy_types", "rules"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": os.Getenv("RESOURCE_GROUP_ID"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": os.Getenv("RESOURCE_GROUP_ID"),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"real_servers": `[\"1.1.1.1\",\"2.2.2.2\",\"3.3.3.3\"]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"real_servers": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_types": `[{\"ProxyType\":\"https\",\"ProxyPorts\":[443]}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_types": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"proxy_types":  `[{\"ProxyType\":\"http\",\"ProxyPorts\":[80]},{\"ProxyType\":\"https\",\"ProxyPorts\":[443]}]`,
					"real_servers": `[\"1.1.1.1\",\"2.2.2.2\"]`,

				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"proxy_types": CHECKSET,
						"real_servers": CHECKSET,
					}),
				),
			},
		},
	})
}

var DdosCooWebRuleMap = map[string]string{}

func DdosCooWebRuleBasicdependence(name string) string {
	return ""
}
