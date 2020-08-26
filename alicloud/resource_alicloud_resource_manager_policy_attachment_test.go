package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudResourceManagerPolicyAttachment_basic(t *testing.T) {
	var v resourcemanager.PolicyAttachment
	resourceId := "alicloud_resource_manager_policy_attachment.default"
	ra := resourceAttrInit(resourceId, ResourceManagerPolicyAttachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerPolicyAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerPolicyAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerPolicyAttachmentBasicdependence)
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
					"policy_name":       "${alicloud_resource_manager_policy.this.policy_name}",
					"policy_type":       "Custom",
					"principal_name":    "${local.principal_name}",
					"principal_type":    "IMSUser",
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.this.groups.0.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_name":       name,
						"policy_type":       "Custom",
						"principal_name":    CHECKSET,
						"principal_type":    "IMSUser",
						"resource_group_id": CHECKSET,
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

var ResourceManagerPolicyAttachmentMap = map[string]string{}

func ResourceManagerPolicyAttachmentBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_ram_user" "this" {
  name = "${var.name}"
}

resource "alicloud_resource_manager_policy" "this" {
  policy_name     = "${var.name}"
  description 	  = "policy_attachment"
  policy_document = <<EOF
        {
            "Statement": [{
                "Action": ["oss:*"],
                "Effect": "Allow",
                "Resource": ["acs:oss:*:*:*"]
            }],
            "Version": "1"
        }
    EOF
}

data "alicloud_account" "this" {}

data "alicloud_resource_manager_resource_groups" "this" {
  name_regex = "default"
}

locals{
	principal_name = format("%%s@%%s.onaliyun.com", alicloud_ram_user.this.name, data.alicloud_account.this.id)	
}
`, name)
}
