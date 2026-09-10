package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 标签策略测试用例_1_副本1737105427623 10066
func TestAccAliCloudTagPolicy_basic10066(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TagPolicySupportRegions)
	resourceId := "alicloud_tag_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudTagPolicyMap10066)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &TagServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTagPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacctag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudTagPolicyBasicDependence10066)
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
					"policy_name":    name,
					"policy_content": `{\"tags\":{\"CostCenter\":{\"tag_value\":{\"@@assign\":[\"Beijing\",\"Shanghai\"]},\"tag_key\":{\"@@assign\":\"CostCenter\"}}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name":    name,
						"policy_content": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_content": `{\"tags\":{\"CostCenter\":{\"tag_value\":{\"@@assign\":[\"Shanghai\"]},\"tag_key\":{\"@@assign\":\"CostCenter\"}}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_content": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_desc": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_desc": name,
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

var AliCloudTagPolicyMap10066 = map[string]string{
	"user_type": CHECKSET,
}

func AliCloudTagPolicyBasicDependence10066(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
`, name)
}

// Case Policy 10066  twin
func TestAccAliCloudTagPolicy_basic10066_twin(t *testing.T) {
	var v map[string]interface{}
	checkoutSupportedRegions(t, true, connectivity.TagPolicySupportRegions)
	resourceId := "alicloud_tag_policy.default"
	ra := resourceAttrInit(resourceId, AliCloudTagPolicyMap10066)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &TagServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTagPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacctag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudTagPolicyBasicDependence10066)
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
					"policy_desc":    name,
					"policy_name":    name,
					"policy_content": `{\"tags\":{\"CostCenter\":{\"tag_value\":{\"@@assign\":[\"Beijing\",\"Shanghai\"]},\"tag_key\":{\"@@assign\":\"CostCenter\"}}}}`,
					"user_type":      "USER",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_desc":    name,
						"policy_name":    name,
						"policy_content": CHECKSET,
						"user_type":      "USER",
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
