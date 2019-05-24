package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudRamRolePolicyAttachment_basic(t *testing.T) {
	var v *ram.Policy
	resourceId := "alicloud_ram_role_policy_attachment.default"
	ra := resourceAttrInit(resourceId, ramPolicyForRoleMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000000, 9999999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamRolePolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamRolePolicyAttachmentConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":   fmt.Sprintf("tf-testAcc%sRamRolePolicyAttachmentConfig-%d", defaultRegionToTest, rand),
						"policy_name": fmt.Sprintf("tf-testAcc%sRamRolePolicyAttachmentConfig-%d", defaultRegionToTest, rand),
					}),
				),
			},
		},
	})
}

var ramPolicyForRoleMap = map[string]string{
	"role_name":   CHECKSET,
	"policy_name": CHECKSET,
	"policy_type": "Custom",
}

func testAccRamRolePolicyAttachmentConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAcc%sRamRolePolicyAttachmentConfig-%d"
	}

	resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  statement = [
	    {
	      effect = "Deny"
	      action = [
		"oss:ListObjects",
		"oss:ListObjects"]
	      resource = [
		"acs:oss:*:*:mybucket",
		"acs:oss:*:*:mybucket/*"]
	    }]
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
	  description = "this is a test"
	  force = true
	}

	resource "alicloud_ram_role_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  role_name = "${alicloud_ram_role.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}`, defaultRegionToTest, rand)
}

func testAccCheckRamRolePolicyAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_role_policy_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListPoliciesForRoleRequest()
		request.RoleName = rs.Primary.Attributes["role_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForRole(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
		response, _ := raw.(*ram.ListPoliciesForRoleResponse)
		if len(response.Policies.Policy) > 0 {
			for _, v := range response.Policies.Policy {
				if v.PolicyName == rs.Primary.Attributes["policy_name"] && v.PolicyType == rs.Primary.Attributes["policy_type"] {
					return WrapError(Error("Error attachment still exist."))
				}
			}
		}
	}
	return nil
}
