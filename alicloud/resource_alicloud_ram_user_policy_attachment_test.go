package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ram UserPolicyAttachment. >>> Resource test cases, automatically generated.
// Case UserPolicyAttachment_UserId 9642
func TestAccAliCloudRamUserPolicyAttachment_basic9642(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_user_policy_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudRamUserPolicyAttachmentMap9642)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamUserPolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamUserPolicyAttachmentBasicDependence9642)
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
					"policy_type": "Custom",
					"user_name":   "${alicloud_ram_user.default4nWJxj.name}",
					"policy_name": "${alicloud_ram_policy.defaultVBNCbB.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type": "Custom",
						"user_name":   CHECKSET,
						"policy_name": CHECKSET,
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

func TestAccAliCloudRamUserPolicyAttachment_basic9643(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ram_user_policy_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudRamUserPolicyAttachmentMap9642)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RamServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRamUserPolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccram%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRamUserPolicyAttachmentBasicDependence9642)
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
					"policy_type": "System",
					"user_name":   "${alicloud_ram_user.default4nWJxj.name}",
					"policy_name": "AliyunECSFullAccess",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_type": "System",
						"user_name":   CHECKSET,
						"policy_name": "AliyunECSFullAccess",
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

var AliCloudRamUserPolicyAttachmentMap9642 = map[string]string{}

func AliCloudRamUserPolicyAttachmentBasicDependence9642(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_ram_user" "default4nWJxj" {
  name         = var.name
  display_name = "fortRamUPTest"
}

resource "alicloud_ram_policy" "defaultVBNCbB" {
  policy_name     = var.name
  policy_document = "{\"Statement\": [{\"Effect\": \"Allow\",\"Action\": \"ecs:Describe*\",\"Resource\": \"acs:ecs:cn-qingdao:*:instance/*\"}],\"Version\": \"1\"}"
}

`, name)
}

// Test Ram UserPolicyAttachment. <<< Resource test cases, automatically generated.
