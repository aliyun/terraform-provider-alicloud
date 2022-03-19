package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAlicloudRAMUserPolicyAttachment_basic(t *testing.T) {
	var v *ram.PolicyInListPoliciesForUser
	var u *ram.UserInGetUser
	resourceId := "alicloud_ram_user_policy_attachment.default"
	ra := resourceAttrInit(resourceId, ramPolicyForUserMap)
	serviceFunc := func() interface{} {
		return &RamService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rcUser := resourceCheckInit("alicloud_ram_user.default", &u, serviceFunc)

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
		CheckDestroy:  testAccCheckRamUserPolicyAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRamUserPolicyAttachmentConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					rcUser.checkResourceExists(),
					testAccCheck(nil),
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

var ramPolicyForUserMap = map[string]string{
	"user_name":   CHECKSET,
	"policy_name": CHECKSET,
	"policy_type": "Custom",
}

func testAccRamUserPolicyAttachmentConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "tf-testAcc%sRamUserPolicyAttachmentConfig-%d"
	}
	resource "alicloud_ram_policy" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": [
				"oss:ListObjects",
				"oss:ListObjects"
			  ],
			  "Effect": "Deny",
			  "Resource": [
				"acs:oss:*:*:mybucket",
				"acs:oss:*:*:mybucket/*"
			  ]
			}
		  ],
			"Version": "1"
		}
	  EOF
	  description = "this is a policy test"
	  force = true
	}

	resource "alicloud_ram_user" "default" {
	  name = "${var.name}"
	  display_name = "displayname"
	  mobile = "86-18888888888"
	  email = "hello.uuu@aaa.com"
	  comments = "yoyoyo"
	}

	resource "alicloud_ram_user_policy_attachment" "default" {
	  policy_name = "${alicloud_ram_policy.default.name}"
	  user_name = "${alicloud_ram_user.default.name}"
	  policy_type = "${alicloud_ram_policy.default.type}"
	}`, defaultRegionToTest, rand)
}

func testAccCheckRamUserPolicyAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ram_user_policy_attachment" {
			continue
		}

		// Try to find the attachment
		client := testAccProvider.Meta().(*connectivity.AliyunClient)

		request := ram.CreateListPoliciesForUserRequest()
		request.UserName = rs.Primary.Attributes["user_name"]

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListPoliciesForUser(request)
		})

		if err != nil && !IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return WrapError(err)
		}
		response, _ := raw.(*ram.ListPoliciesForUserResponse)
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
