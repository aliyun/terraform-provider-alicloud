package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSDDPRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sddp_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudSDDPRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SddpService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSddpRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddprule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSDDPRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.SddpSupportRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"category":  "2",
					"content":   "tf-testACC",
					"rule_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":  "2",
						"content":   "tf-testACC",
						"rule_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "${var.name}-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": "tf-testACCContent",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content": "tf-testACCContent",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"lang": "zh",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"lang": "zh",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_type": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_type": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "DescriptionAlone",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "DescriptionAlone",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level_id": "4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level_id": "4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"stat_express": "StatExpress",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"stat_express": "StatExpress",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"warn_level": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"warn_level": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code": "OSS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_code": "OSS",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"product_id": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"product_id": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target": "Content",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target": "Content",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content_category": "106",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_category": "106",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category":         "0",
					"content":          "tf-testACC-all",
					"rule_name":        "${var.name}-all",
					"description":      "DescriptionAlone",
					"risk_level_id":    "4",
					"stat_express":     "StatExpress",
					"warn_level":       "1",
					"product_code":     "OSS",
					"product_id":       "2",
					"lang":             "en",
					"status":           "1",
					"rule_type":        "1",
					"target":           "Content",
					"content_category": "104",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category":         "0",
						"content":          "tf-testACC-all",
						"rule_name":        name + "-all",
						"description":      "DescriptionAlone",
						"risk_level_id":    "4",
						"stat_express":     "StatExpress",
						"warn_level":       "1",
						"product_code":     "OSS",
						"product_id":       "2",
						"status":           "1",
						"lang":             "en",
						"rule_type":        "1",
						"target":           "Content",
						"content_category": "104",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true, ImportStateVerifyIgnore: []string{"ids", "lang", "rule_type"},
			},
		},
	})
}

var AlicloudSDDPRuleMap0 = map[string]string{}

func AlicloudSDDPRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
