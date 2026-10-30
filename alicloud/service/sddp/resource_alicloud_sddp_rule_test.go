package sddp_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sddp"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudSDDPRule_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.SddpSupportRegions)
	resourceId := "alicloud_sddp_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSDDPRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sddp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSddpRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddprule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSDDPRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name,
					"category":  "0",
					"content":   `[{\"rule\":[{\"operator\":\"contains\",\"target\":\"content\",\"value\":\"tf-testACCContent\"}],\"ruleRelation\":\"AND\"}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name,
						"category":  "0",
						"content":   acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"category": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"category": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content": `[{\"rule\":[{\"operator\":\"contains\",\"target\":\"content\",\"value\":\"tf-testACCContent-update\"}],\"ruleRelation\":\"AND\"}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"risk_level_id": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"risk_level_id": "5",
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
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rule_type", "lang"},
			},
		},
	})
}

func TestAccAliCloudSDDPRule_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.SddpSupportRegions)
	resourceId := "alicloud_sddp_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudSDDPRuleMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sddp.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSddpRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssddprule%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudSDDPRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name":        name,
					"category":         "0",
					"content":          `[{\"rule\":[{\"operator\":\"contains\",\"target\":\"content\",\"value\":\"tf-testACCContent\"}],\"ruleRelation\":\"AND\"}]`,
					"content_category": "106",
					"risk_level_id":    "5",
					"rule_type":        "1",
					"product_code":     "OSS",
					"product_id":       "2",
					"warn_level":       "1",
					"stat_express":     name,
					"target":           "tf-target",
					"status":           "1",
					"description":      name,
					"lang":             "en",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name":        name,
						"category":         "0",
						"content":          acctest.CHECKSET,
						"content_category": "106",
						"risk_level_id":    "5",
						"rule_type":        "1",
						"product_code":     "OSS",
						"product_id":       "2",
						"warn_level":       "1",
						"stat_express":     name,
						"target":           "tf-target",
						"status":           "1",
						"description":      name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rule_type", "lang"},
			},
		},
	})
}

var AliCloudSDDPRuleMap0 = map[string]string{
	"content_category": acctest.CHECKSET,
	"status":           acctest.CHECKSET,
	"custom_type":      acctest.CHECKSET,
}

func AliCloudSDDPRuleBasicDependence0(name string) string {
	return ""
}
