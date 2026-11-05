package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA WafRuleset. >>> Resource test cases, automatically generated.
// Case resource_WafRuleset_test
func TestAccAliCloudESAWafRulesetresource_WafRuleset_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_waf_ruleset.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAWafRulesetresource_WafRuleset_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaWafRuleset")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAWafRuleset%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAWafRulesetresource_WafRuleset_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"phase":        "http_custom",
					"site_version": "0",
					"name":         "test_waf_ruleSet",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":       "Off",
					"site_version": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":       "On",
					"site_version": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site_version"},
			},
		},
	})
}

var AliCloudESAWafRulesetresource_WafRuleset_testMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudESAWafRulesetresource_WafRuleset_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}
`, name)
}

// Test ESA WafRuleset. <<< Resource test cases, automatically generated.
