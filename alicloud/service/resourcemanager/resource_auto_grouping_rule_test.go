package resourcemanager_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/resourcemanager"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order stock-debt declarations (pure rename migration; see scripts/collection-order/README.md)
// collection-order: bypass rule_contents -- pre-existing stock TypeList; this is a pure package-rename migration (alicloud/ -> alicloud/service/resourcemanager/), schema structurally unchanged vs merge-base 2f8f887c52e916d4e87f2e7ed34bdf2aba355259

// Test ResourceManager AutoGroupingRule. >>> Resource test cases, automatically generated.
// Case 自动分组规则用例_自定义条件_2 9946
func TestAccAliCloudResourceManagerAutoGroupingRule_basic9952(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_auto_grouping_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerAutoGroupingRuleMap9952)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerAutoGroupingRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerAutoGroupingRuleBasicDependence9952)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r1.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
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
					"exclude_resource_group_ids_scope": "${alicloud_resource_manager_resource_group.r3.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_group_ids_scope": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_group_ids_scope": acctest.REMOVEKEY,
					"resource_group_ids_scope":         "${alicloud_resource_manager_resource_group.r3.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_group_ids_scope": acctest.REMOVEKEY,
						"resource_group_ids_scope":         acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_group_ids_scope": "${alicloud_resource_manager_resource_group.r3.id}",
					"resource_group_ids_scope":         acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_group_ids_scope": acctest.CHECKSET,
						"resource_group_ids_scope":         acctest.REMOVEKEY,
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
					"exclude_region_ids_scope": acctest.REMOVEKEY,
					"region_ids_scope":         "cn-beijing",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_region_ids_scope": acctest.REMOVEKEY,
						"region_ids_scope":         "cn-beijing",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_region_ids_scope": "cn-beijing",
					"region_ids_scope":         acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_region_ids_scope": "cn-beijing",
						"region_ids_scope":         acctest.REMOVEKEY,
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
					"exclude_resource_ids_scope": acctest.REMOVEKEY,
					"resource_ids_scope":         "dmock-xxxx",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_ids_scope": acctest.REMOVEKEY,
						"resource_ids_scope":         "dmock-xxxx",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_ids_scope": "dmock-xxxx",
					"resource_ids_scope":         acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_ids_scope": "dmock-xxxx",
						"resource_ids_scope":         acctest.REMOVEKEY,
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
					"exclude_resource_types_scope": acctest.REMOVEKEY,
					"resource_types_scope":         "ecs.instance",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_types_scope": acctest.REMOVEKEY,
						"resource_types_scope":         "ecs.instance",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"exclude_resource_types_scope": "ecs.instance",
					"resource_types_scope":         acctest.REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"exclude_resource_types_scope": "ecs.instance",
						"resource_types_scope":         acctest.REMOVEKEY,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r1.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_b\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_c\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r2.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_d\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r3.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerAutoGroupingRuleMap9952)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerAutoGroupingRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerAutoGroupingRuleBasicDependence9952)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r1.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_b\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_c\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r2.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_d\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r3.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
					},
					"rule_desc":                        "资源用例测试规则",
					"rule_type":                        "custom_condition",
					"exclude_region_ids_scope":         "cn-beijing",
					"exclude_resource_group_ids_scope": "${alicloud_resource_manager_resource_group.r3.id}",
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
						"exclude_resource_group_ids_scope": acctest.CHECKSET,
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudResourceManagerAutoGroupingRuleMap9952)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return resourcemanager.NewResourceManagerServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeResourceManagerAutoGroupingRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudResourceManagerAutoGroupingRuleBasicDependence9952)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_contents": []map[string]interface{}{
						{
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r1.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_b\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_c\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r2.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
						{
							"auto_grouping_scope_condition":   "{\\\"children\\\":[{\\\"desired\\\":\\\"name_d\\\",\\\"featurePath\\\":\\\"$.resourceName\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEqualsAny\\\"}],\\\"operator\\\":\\\"and\\\"}",
							"target_resource_group_condition": "{\\\"children\\\":[{\\\"desired\\\":\\\"${alicloud_resource_manager_resource_group.r3.id}\\\",\\\"featurePath\\\":\\\"$.resourceGroupId\\\",\\\"featureSource\\\":\\\"RESOURCE\\\",\\\"operator\\\":\\\"StringEquals\\\"}],\\\"operator\\\":\\\"and\\\"}",
						},
					},
					"rule_desc":                "资源用例测试规则",
					"rule_type":                "custom_condition",
					"region_ids_scope":         "cn-beijing",
					"resource_group_ids_scope": "${alicloud_resource_manager_resource_group.r3.id}",
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
						"resource_group_ids_scope": acctest.CHECKSET,
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

resource "alicloud_resource_manager_resource_group" "r1" {
  resource_group_name = "${var.name}-1"
  display_name        = "${var.name}-1"
}
resource "alicloud_resource_manager_resource_group" "r2" {
  resource_group_name = "${var.name}-2"
  display_name        = "${var.name}-2"
}
resource "alicloud_resource_manager_resource_group" "r3" {
  resource_group_name = "${var.name}-3"
  display_name        = "${var.name}-3"
}
`, name)
}

// Test ResourceManager AutoGroupingRule. <<< Resource test cases, automatically generated.
