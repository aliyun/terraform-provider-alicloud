package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ResourceManager AutoGroupingRule. >>> Resource test cases, automatically generated.
// Case 自动分组规则用例_自定义条件_2 9946
func TestAccAliCloudResourceManagerAutoGroupingRule_basic9952(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_auto_grouping_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerAutoGroupingRuleMap9952)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerAutoGroupingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerAutoGroupingRuleBasicDependence9952)
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
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******zj2ob\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
					},
					"rule_type": "custom_condition",
					"rule_name": "资源用例测试规则",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_contents.#": "1",
						"rule_type":       "custom_condition",
						"rule_name":       "资源用例测试规则",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_group_ids_scope": "rg-aekz******4b5ea",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_group_ids_scope": "rg-aekz******4b5ea",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_group_ids_scope": REMOVEKEY,
					"resource_group_ids_scope":         "rg-aekz******4b5ea",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_group_ids_scope": REMOVEKEY,
						"resource_group_ids_scope":         "rg-aekz******4b5ea",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_group_ids_scope": "rg-aekz******4b5ea",
					"resource_group_ids_scope":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_group_ids_scope": "rg-aekz******4b5ea",
						"resource_group_ids_scope":         REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_region_ids_scope": "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_region_ids_scope": "cn-beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_region_ids_scope": REMOVEKEY,
					"region_ids_scope":         "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_region_ids_scope": REMOVEKEY,
						"region_ids_scope":         "cn-beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_region_ids_scope": "cn-beijing",
					"region_ids_scope":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_region_ids_scope": "cn-beijing",
						"region_ids_scope":         REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": "dmock-xxxx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_ids_scope": "dmock-xxxx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": REMOVEKEY,
					"resource_ids_scope":         "dmock-xxxx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_ids_scope": REMOVEKEY,
						"resource_ids_scope":         "dmock-xxxx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": "dmock-xxxx",
					"resource_ids_scope":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_ids_scope": "dmock-xxxx",
						"resource_ids_scope":         REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_types_scope": "ecs.instance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_types_scope": "ecs.instance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_types_scope": REMOVEKEY,
					"resource_types_scope":         "ecs.instance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_types_scope": REMOVEKEY,
						"resource_types_scope":         "ecs.instance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_types_scope": "ecs.instance",
					"resource_types_scope":         REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_types_scope": "ecs.instance",
						"resource_types_scope":         REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******zj2ob\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_b\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_c\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******r62ua\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_d\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******4b5ea\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_contents.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_name": "资源用例测试规则_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": "资源用例测试规则_2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_desc": "资源用例测试规则_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_desc": "资源用例测试规则_2",
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

func TestAccAliCloudResourceManagerAutoGroupingRule_basic9952_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_auto_grouping_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerAutoGroupingRuleMap9952)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerAutoGroupingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerAutoGroupingRuleBasicDependence9952)
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
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******zj2ob\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_b\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_c\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******r62ua\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_d\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******4b5ea\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
					},
					"rule_desc":                        "资源用例测试规则",
					"rule_type":                        "custom_condition",
					"exclude_region_ids_scope":         "cn-beijing",
					"exclude_resource_group_ids_scope": "rg-aekz******4b5ea",
					"exclude_resource_ids_scope":       "dmock-xxxx",
					"exclude_resource_types_scope":     "ecs.instance",
					"rule_name":                        "资源用例测试规则",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_contents.#":                  "3",
						"rule_desc":                        "资源用例测试规则",
						"rule_type":                        "custom_condition",
						"exclude_region_ids_scope":         "cn-beijing",
						"exclude_resource_group_ids_scope": "rg-aekz******4b5ea",
						"exclude_resource_ids_scope":       "dmock-xxxx",
						"exclude_resource_types_scope":     "ecs.instance",
						"rule_name":                        "资源用例测试规则",
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

func TestAccAliCloudResourceManagerAutoGroupingRule_basic9958_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_auto_grouping_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudResourceManagerAutoGroupingRuleMap9952)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerAutoGroupingRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudResourceManagerAutoGroupingRuleBasicDependence9952)
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
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******zj2ob\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_b\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_c\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******r62ua\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_d\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"rg-aekz******4b5ea\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
					},
					"rule_desc":                "资源用例测试规则",
					"rule_type":                "custom_condition",
					"region_ids_scope":         "cn-beijing",
					"resource_group_ids_scope": "rg-aekz******4b5ea",
					"resource_ids_scope":       "dmock-xxxx",
					"resource_types_scope":     "ecs.instance",
					"rule_name":                "资源用例测试规则",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_contents.#":          "3",
						"rule_desc":                "资源用例测试规则",
						"rule_type":                "custom_condition",
						"region_ids_scope":         "cn-beijing",
						"resource_group_ids_scope": "rg-aekz******4b5ea",
						"resource_ids_scope":       "dmock-xxxx",
						"resource_types_scope":     "ecs.instance",
						"rule_name":                "资源用例测试规则",
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

var AliCloudResourceManagerAutoGroupingRuleMap9952 = map[string]string{}

func AliCloudResourceManagerAutoGroupingRuleBasicDependence9952(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ResourceManager AutoGroupingRule. <<< Resource test cases, automatically generated.
