// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package ddoscoo_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ddoscoo"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test DdosCoo DomainPreciseAccessRule. >>> Resource test cases, automatically generated.
// Case 网站业务精确访问控制规则 12326
func TestAccAliCloudDdosCooDomainPreciseAccessRule_basic12326(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_domain_precise_access_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosCooDomainPreciseAccessRuleMap12326)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooDomainPreciseAccessRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosCooDomainPreciseAccessRuleBasicDependence12326)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"condition": []map[string]interface{}{
						{
							"match_method": "belong",
							"field":        "ip",
							"content":      "1.1.1.1",
						},
					},
					"action": "accept",
					"domain": "${alicloud_ddoscoo_domain_resource.default.id}",
					"name":   name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition.#": "1",
						"action":      "accept",
						"domain":      acctest.CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action": "block",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action": "block",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition": []map[string]interface{}{
						{
							"match_method": "contain",
							"field":        "header",
							"content":      "222",
							"header_name":  "122",
						},
						{
							"match_method": "contain",
							"field":        "referer",
							"content":      "22",
						},
						{
							"match_method": "belong",
							"field":        "ip",
							"content":      "1.1.1.1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"expires": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"expires": "600",
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

func TestAccAliCloudDdosCooDomainPreciseAccessRule_basic12326_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ddoscoo_domain_precise_access_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudDdosCooDomainPreciseAccessRuleMap12326)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ddoscoo.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeDdosCooDomainPreciseAccessRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1, 999)
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudDdosCooDomainPreciseAccessRuleBasicDependence12326)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheckWithAccountSiteType(t, acctest.DomesticSite)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"condition": []map[string]interface{}{
						{
							"match_method": "contain",
							"field":        "header",
							"content":      "222",
							"header_name":  "122",
						},
						{
							"match_method": "contain",
							"field":        "referer",
							"content":      "22",
						},
						{
							"match_method": "belong",
							"field":        "ip",
							"content":      "1.1.1.1",
						},
					},
					"action":  "block",
					"domain":  "${alicloud_ddoscoo_domain_resource.default.id}",
					"name":    name,
					"expires": "600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition.#": "3",
						"action":      "block",
						"domain":      acctest.CHECKSET,
						"name":        name,
						"expires":     "600",
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

var AliCloudDdosCooDomainPreciseAccessRuleMap12326 = map[string]string{}

func AliCloudDdosCooDomainPreciseAccessRuleBasicDependence12326(name string) string {
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

// Test DdosCoo DomainPreciseAccessRule. <<< Resource test cases, automatically generated.
