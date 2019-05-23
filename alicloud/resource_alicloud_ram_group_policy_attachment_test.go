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

func TestAccAlicloudRamGroupPolicyAttachment_basic(t *testing.T) {
	var v *ram.Policy
	resourceId := "alicloud_ram_group_policy_attachment.default"
	ra := resourceAttrInit(resourceId, ramGroupMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckRamGroupPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamGroupPolicyAttachmentCreateConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var ramGroupMap = map[string]string{
	"group_name":  CHECKSET,
	"policy_name": CHECKSET,
	"policy_type": "Custom",
}

func testAccRamGroupPolicyAttachmentCreateConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamGroupPolicyAttachmentConfig-%d"
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

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}

	resource "alicloud_ram_group_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  group_name = "${alicloud_ram_group.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}`, defaultRegionToTest, rand)
}

func testAccCheckRamGroupPolicyAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_group_policy_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListPoliciesForGroupRequest()
		request.GroupName = rs.Primary.Attributes["group_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForGroup(request)
		})

		if err != nil && !RamEntityNotExist(err) {
			return WrapError(err)
		}
		response, _ := raw.(*ram.ListPoliciesForGroupResponse)
		if len(response.Policies.Policy) > 0 {
			for _, v := range response.Policies.Policy {
				if v.PolicyName == rs.Primary.Attributes["name"] && v.PolicyType == rs.Primary.Attributes["policy_type"] {
					return WrapError(Error("Error attachment still exist."))
				}
			}
		}
	}
	return nil
}
