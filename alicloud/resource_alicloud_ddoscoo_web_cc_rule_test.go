package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DdosCoo WebCcRule. >>> Resource test cases, automatically generated.
// Case 频率控制策略 12398
func TestAccAliCloudDdosCooWebCcRule_basic12398(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_web_cc_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosCooWebCcRuleMap12398)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosCooServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosCooWebCcRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosCooWebCcRuleBasicDependence12398)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_detail": []map[string]interface{}{
						{
							"action": "block",
							"rate_limit": []map[string]interface{}{
								{
									"interval":  "10",
									"threshold": "2",
									"ttl":       "860",
									"target":    "ip",
								},
							},
							"condition": []map[string]interface{}{
								{
									"match_method": "belong",
									"field":        "ip",
									"content":      "1.1.1.1",
								},
							},
						},
					},
					"name":   name,
					"domain": "${alicloud_ddoscoo_domain_resource.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"domain":        CHECKSET,
						"rule_detail.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_detail": []map[string]interface{}{
						{
							"action": "watch",
							"rate_limit": []map[string]interface{}{
								{
									"interval":  "600",
									"threshold": "70",
									"ttl":       "300",
									"target":    "header",
									"sub_key":   "30",
								},
							},
							"condition": []map[string]interface{}{
								{
									"match_method": "contain",
									"field":        "header",
									"header_name":  "222",
									"content":      "666",
								},
								{
									"match_method": "belong",
									"field":        "ip",
									"content":      "1.1.1.1",
								},
								{
									"match_method": "nbelong",
									"field":        "ip",
									"content":      "2.2.2.2",
								},
							},
							"statistics": []map[string]interface{}{
								{
									"mode":        "distinct",
									"field":       "header",
									"header_name": "222",
								},
							},
							"status_code": []map[string]interface{}{
								{
									"enabled":         "false",
									"code":            "100",
									"use_ratio":       "false",
									"count_threshold": "8",
									"ratio_threshold": "2",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_detail.#": "1",
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

func TestAccAliCloudDdosCooWebCcRule_basic12398_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_web_cc_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudDdosCooWebCcRuleMap12398)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdosCooServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdosCooWebCcRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudDdosCooWebCcRuleBasicDependence12398)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_detail": []map[string]interface{}{
						{
							"action": "watch",
							"rate_limit": []map[string]interface{}{
								{
									"interval":  "600",
									"threshold": "70",
									"ttl":       "300",
									"target":    "header",
									"sub_key":   "30",
								},
							},
							"condition": []map[string]interface{}{
								{
									"match_method": "contain",
									"field":        "header",
									"header_name":  "222",
									"content":      "666",
								},
								{
									"match_method": "belong",
									"field":        "ip",
									"content":      "1.1.1.1",
								},
								{
									"match_method": "nbelong",
									"field":        "ip",
									"content":      "2.2.2.2",
								},
							},
							"statistics": []map[string]interface{}{
								{
									"mode":        "distinct",
									"field":       "header",
									"header_name": "222",
								},
							},
							"status_code": []map[string]interface{}{
								{
									"enabled":         "false",
									"code":            "100",
									"use_ratio":       "false",
									"count_threshold": "8",
									"ratio_threshold": "2",
								},
							},
						},
					},
					"name":   name,
					"domain": "${alicloud_ddoscoo_domain_resource.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":          name,
						"domain":        CHECKSET,
						"rule_detail.#": "1",
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

var AliCloudDdosCooWebCcRuleMap12398 = map[string]string{}

func AliCloudDdosCooWebCcRuleBasicDependence12398(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_ddoscoo_instances" "default" {
}

resource "alicloud_ddoscoo_domain_resource" "default" {
  domain       = "tf-testacc%s.alibaba.com"
  instance_ids = [data.alicloud_ddoscoo_instances.default.ids.0]
  proxy_types {
    proxy_ports = [443]
    proxy_type  = "https"
  }
  real_servers = ["177.167.32.11"]
  rs_type      = 0
}

`, name, name)
}

// Test DdosCoo WebCcRule. <<< Resource test cases, automatically generated.
