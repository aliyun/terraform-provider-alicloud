// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Esa CustomResponseCodeRule. >>> Resource test cases, automatically generated.
// Case CustomResponseCodeRule_test 12096
func TestAccAliCloudEsaCustomResponseCodeRule_basic12096(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_custom_response_code_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudEsaCustomResponseCodeRuleMap12096)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomResponseCodeRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaCustomResponseCodeRuleBasicDependence12096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"site_id":     CHECKSET,
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
						"page_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"return_code": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"return_code": CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudEsaCustomResponseCodeRuleMap12096)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomResponseCodeRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaCustomResponseCodeRuleBasicDependence12096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"site_id":      CHECKSET,
						"return_code":  "400",
						"rule_enable":  "on",
						"rule":         CHECKSET,
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
	ra := resourceAttrInit(resourceId, AliCloudEsaCustomResponseCodeRuleMap12096)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaCustomResponseCodeRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccesa%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEsaCustomResponseCodeRuleBasicDependence12096)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"site_id":     CHECKSET,
						"return_code": "400",
						"rule_enable": "on",
						"rule":        CHECKSET,
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
						"page_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"return_code": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"return_code": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule": "(http.request.uri eq \\\"/content?page=1234\\\")",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule": CHECKSET,
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
	"config_id": CHECKSET,
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
