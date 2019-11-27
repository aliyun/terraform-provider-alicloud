package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/hashicorp/terraform/helper/acctest"
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
	name := fmt.Sprintf("tf-testAcc%sRamRolePolicyAttachmentConfig-%d", defaultRegionToTest, rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceRamRolePolicyAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.default.name}",
					"role_name":   "${alicloud_ram_role.default.name}",
					"policy_type": "${alicloud_ram_policy.default.type}",
				}),
				Check: resource.ComposeTestCheckFunc(
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

var ramPolicyForRoleMap = map[string]string{
	"role_name":   CHECKSET,
	"policy_name": CHECKSET,
	"policy_type": "Custom",
}

func resourceRamRolePolicyAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
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

	resource "alicloud_ram_role" "default" {
	  name = "${var.name}"
	  document = <<EOF
		{
		  "Statement": [
			{
			  "Action": "sts:AssumeRole",
			  "Effect": "Allow",
			  "Principal": {
				"Service": [
				  "apigateway.aliyuncs.com", 
				  "ecs.aliyuncs.com"
				]
			  }
			}
		  ],
		  "Version": "1"
		}
	  EOF
	  description = "this is a test"
	  force = true
	}
`, name)
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
