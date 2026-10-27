// Package tag_test.
package tag_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/tag"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// collection-order: bypass tag_keys -- pre-existing uncovered TypeList attribute outside this pure directory migration (file moved from alicloud/ root; the attribute and its test coverage are unchanged stock)
// Test Tag AssociatedRule. >>> Resource test cases, automatically generated.
// Case 关联资源测试用例_副本1737426969003 10094
func TestAccAliCloudTagAssociatedRule_basic10094(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_tag_associated_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudTagAssociatedRuleMap10094)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return tag.NewTagServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTagAssociatedRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacctag%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudTagAssociatedRuleBasicDependence10094)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudTagAssociatedRuleMap10094)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return tag.NewTagServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeTagAssociatedRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacctag%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudTagAssociatedRuleBasicDependence10094)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
