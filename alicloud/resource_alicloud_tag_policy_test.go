package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccAlicloudTagPolicy_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_policy.default"
	checkoutSupportedRegions(t, true, connectivity.TagPolicySupportRegions)
	ra := resourceAttrInit(resourceId, TagPolicyMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &TagService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTagPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("testAccTagPolicy%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, TagPolicyBasicdependence)
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
					"policy_content": `{\"tags\":{\"CostCenter\":{\"tag_value\":{\"@@assign\":[\"Beijing\",\"Shanghai\"]},\"tag_key\":{\"@@assign\":\"CostCenter\"}}}}`,
					"policy_name":    name,
					"policy_desc":    "test_desc",
					"user_type":      "USER",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_content": CHECKSET,
						"policy_name":    name,
						"policy_desc":    "test_desc",
						"user_type":      "USER",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_desc": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_desc": "test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_content": `{\"tags\":{\"CostCenter\":{\"tag_value\":{\"@@assign\":[\"Shanghai\"]},\"tag_key\":{\"@@assign\":\"CostCenter\"}}}}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_content": `{\"tags\":{\"CostCenter\":{\"tag_value\":{\"@@assign\":[\"Shanghai\"]},\"tag_key\":{\"@@assign\":\"CostCenter\"}}}}`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": "test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name": "test",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var TagPolicyMap = map[string]string{}

func TagPolicyBasicdependence(name string) string {
	return ""
}
