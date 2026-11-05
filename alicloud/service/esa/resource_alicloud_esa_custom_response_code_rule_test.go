// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

// Test Esa CustomResponseCodeRule. >>> Resource test cases, automatically generated.
// Case CustomResponseCodeRule_test 12096
func TestAccAliCloudEsaCustomResponseCodeRule_basic12096(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_response_code_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaCustomResponseCodeRuleMap12096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCustomResponseCodeRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaCustomResponseCodeRuleBasicDependence12096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"page_id":     "0",
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"return_code": "400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"page_id":     "0",
						"site_id":     acctest.CHECKSET,
						"return_code": "400",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"page_id": "${alicloud_esa_page.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"page_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"return_code": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"return_code": acctest.CHECKSET,
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
					"rule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name,
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

func TestAccAliCloudEsaCustomResponseCodeRule_basic12096_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_response_code_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaCustomResponseCodeRuleMap12096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCustomResponseCodeRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaCustomResponseCodeRuleBasicDependence12096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"page_id":      "0",
					"site_id":      "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"return_code":  "400",
					"rule_enable":  "on",
					"rule":         "(http.host eq \\\"video.example.com\\\")",
					"rule_name":    name,
					"sequence":     "1",
					"site_version": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"page_id":      "0",
						"site_id":      acctest.CHECKSET,
						"return_code":  "400",
						"rule_enable":  "on",
						"rule":         acctest.CHECKSET,
						"rule_name":    name,
						"sequence":     "1",
						"site_version": "0",
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

func TestAccAliCloudEsaCustomResponseCodeRule_basic12098(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_response_code_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudEsaCustomResponseCodeRuleMap12096)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaCustomResponseCodeRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudEsaCustomResponseCodeRuleBasicDependence12096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"page_id":     "0",
					"site_id":     "${data.alicloud_esa_sites.default.sites.0.site_id}",
					"return_code": "400",
					"rule_enable": "on",
					"rule":        "(http.host eq \\\"video.example.com\\\")",
					"rule_name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"page_id":     "0",
						"site_id":     acctest.CHECKSET,
						"return_code": "400",
						"rule_enable": "on",
						"rule":        acctest.CHECKSET,
						"rule_name":   name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"page_id": "${alicloud_esa_page.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"page_id": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"return_code": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"return_code": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(http.request.uri eq \\\"/content?page=1234\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": acctest.CHECKSET,
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
					"rule_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"sequence": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"sequence": "2",
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

var AliCloudEsaCustomResponseCodeRuleMap12096 = map[string]string{
	"config_id": acctest.CHECKSET,
}

func AliCloudEsaCustomResponseCodeRuleBasicDependence12096(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_page" "default" {
  content_type = "text/html"
  content      = "PGh0bWw+aGVsbG8gcGFnZTwvaHRtbD4="
  page_name    = var.name
}
`, name)
}

// Test Esa CustomResponseCodeRule. <<< Resource test cases, automatically generated.
