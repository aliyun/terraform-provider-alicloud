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

// Test ESA HttpResponseHeaderModificationRule. >>> Resource test cases, automatically generated.
// Case httpResponseHeaderModificationRule_test
func TestAccAliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_http_response_header_modification_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaHttpResponseHeaderModificationRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAHttpResponseHeaderModificationRule%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testBasicDependence)
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
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.id}",
					"rule_enable": "on",
					"response_header_modification": []map[string]interface{}{

						{
							"value":     "add",
							"operation": "add",
							"name":      "testadd",
						},

						{
							"operation": "del",
							"name":      "testdel",
						},

						{
							"value":     "modify",
							"operation": "modify",
							"name":      "testmodify",
						},
					},
					"rule":         "(http.host eq \\\"video.example.com\\\")",
					"site_version": "0",
					"rule_name":    "testResponseHeader",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "testResponseHeader_modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(http.request.uri eq \\\"/content?page=1234\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"response_header_modification": []map[string]interface{}{

						{
							"value":     "add1",
							"operation": "add",
							"name":      "testadd1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
					"rule":        "(http.host eq \\\"api.example.com\\\")",
					"rule_name":   "test_httpResponseHeader_last",
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
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudESAHttpResponseHeaderModificationRulehttpResponseHeaderModificationRule_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

`, name)
}

// Test ESA HttpResponseHeaderModificationRule. <<< Resource test cases, automatically generated.
