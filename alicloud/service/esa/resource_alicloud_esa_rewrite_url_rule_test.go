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

// Test ESA RewriteUrlRule. >>> Resource test cases, automatically generated.
// Case resource_RewriteUrlRule_test

func TestAccAliCloudEsaRewriteUrlRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rewrite_url_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRewriteUrlRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRewriteUrlRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARewriteUrlRule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRewriteUrlRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":     "${alicloud_esa_rewrite_url_rule.pre.site_id}",
					"rule_enable": "on",
					"rule":        "http.host eq \\\"video.example.com\\\"",
					"rule_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":     acctest.CHECKSET,
						"rule_enable": "on",
						"rule":        "http.host eq \"video.example.com\"",
						"rule_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_string": "example=123",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_string": "example=123",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rewrite_query_string_type": "static",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rewrite_query_string_type": "static",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rewrite_uri_type": "static",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rewrite_uri_type": "static",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "http.host eq \\\"example.com\\\"",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": "http.host eq \"example.com\"",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "off",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "off",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_enable": "on",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_enable": "on",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "重写URL规则名称示例",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "重写URL规则名称示例",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sequence": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sequence": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"uri": "/image/example.jpg",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"uri": "/image/example.jpg",
					}),
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

func TestAccAliCloudEsaRewriteUrlRule_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_rewrite_url_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaRewriteUrlRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaRewriteUrlRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESARewriteUrlRule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaRewriteUrlRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"site_id":                   "${alicloud_esa_rewrite_url_rule.pre.site_id}",
					"rule_enable":               "on",
					"rule":                      "http.host eq \\\"video.example.com\\\"",
					"rule_name":                 name,
					"query_string":              "example=123",
					"rewrite_query_string_type": "static",
					"rewrite_uri_type":          "static",
					"sequence":                  "1",
					"site_version":              "0",
					"uri":                       "/image/example.jpg",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"site_id":                   acctest.CHECKSET,
						"rule_enable":               "on",
						"rule":                      "http.host eq \"video.example.com\"",
						"rule_name":                 name,
						"query_string":              "example=123",
						"rewrite_query_string_type": "static",
						"rewrite_uri_type":          "static",
						"sequence":                  "1",
						"site_version":              "0",
						"uri":                       "/image/example.jpg",
					}),
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

var AliCloudEsaRewriteUrlRuleMap0 = map[string]string{
	"config_id": acctest.CHECKSET,
}

func AliCloudEsaRewriteUrlRuleBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name = "gositecdn.cn"
}

resource "alicloud_esa_rewrite_url_rule" "pre" {
  rewrite_uri_type          = "static"
  rewrite_query_string_type = "static"
  site_id                   = data.alicloud_esa_sites.default.sites.0.id
  rule_name                 = "{var.name}pre"
  rule_enable               = "on"
  query_string              = "example=123"
  site_version              = "0"
  rule                      = "http.host eq \"video.example.com\""
  uri                       = "/image/example.jpg"
}
`, name)
}

// Test ESA RewriteUrlRule. <<< Resource test cases, automatically generated.
