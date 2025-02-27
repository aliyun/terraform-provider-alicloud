package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Tag AssociatedRule. >>> Resource test cases, automatically generated.
// Case 关联资源测试用例_副本1737426969003 10094
func TestAccAliCloudTagAssociatedRule_basic10094(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_associated_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudTagAssociatedRuleMap10094)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &TagServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTagAssociatedRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacctag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudTagAssociatedRuleBasicDependence10094)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                  "Enable",
					"associated_setting_name": "rule:AttachEni-DetachEni-TagInstance:Ecs-Instance:Ecs-Eni",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                  "Enable",
						"associated_setting_name": "rule:AttachEni-DetachEni-TagInstance:Ecs-Instance:Ecs-Eni",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Disable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Disable",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_keys": []string{"user1", "user2", "user3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_keys.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_keys": []string{"user1"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_keys.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_keys": []string{},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_keys.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tag_keys": []string{"user1", "user2", "user3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tag_keys.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Enable",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Enable",
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

func TestAccAliCloudTagAssociatedRule_basic10094_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_associated_rule.default"
	ra := resourceAttrInit(resourceId, AliCloudTagAssociatedRuleMap10094)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &TagServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeTagAssociatedRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacctag%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudTagAssociatedRuleBasicDependence10094)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                  "Enable",
					"associated_setting_name": "rule:AttachEni-DetachEni-TagInstance:Ecs-Instance:Ecs-Eni",
					"tag_keys":                []string{"user1", "user2", "user3"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                  "Enable",
						"associated_setting_name": "rule:AttachEni-DetachEni-TagInstance:Ecs-Instance:Ecs-Eni",
						"tag_keys.#":              "3",
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

var AliCloudTagAssociatedRuleMap10094 = map[string]string{}

func AliCloudTagAssociatedRuleBasicDependence10094(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

// Test Tag AssociatedRule. <<< Resource test cases, automatically generated.
