package ram_test

import (
	"fmt"
	"testing"

	ramsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ram"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRAMGroupPolicyAttachment_basic(t *testing.T) {
	var v *ramsdk.Policy
	resourceId := "alicloud_ram_group_policy_attachment.default"
	ra := acctest.ResourceAttrInit(resourceId, ramGroupMap)
	serviceFunc := func() interface{} {
		return ram.NewRamService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)

	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAcc%sRamGroupPolicyAttachmentConfig-%d", acctest.DefaultRegionToTest, rand)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceRamGroupPolicyAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},

		// module name
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_name": "${alicloud_ram_policy.default.name}",
					"group_name":  "${alicloud_ram_group.default.name}",
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

var ramGroupMap = map[string]string{
	"group_name":  acctest.CHECKSET,
	"policy_name": acctest.CHECKSET,
	"policy_type": "Custom",
}

func resourceRamGroupPolicyAttachmentConfigDependence(name string) string {
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

	resource "alicloud_ram_group" "default" {
	  name = "${var.name}"
	  comments = "group comments"
	  force=true
	}
`, name)
}
