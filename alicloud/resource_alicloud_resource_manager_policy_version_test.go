package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudResourceManagerPolicyVersion_basic(t *testing.T) {
	var v resourcemanager.PolicyVersion
	resourceId := "alicloud_resource_manager_policy_version.example"
	ra := resourceAttrInit(resourceId, ResourceManagerPolicyVersionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerPolicyVersion")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerPolicyVersion-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerPolicyVersionBasicdependence)
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
					"policy_document": `{\n\"Statement\": [{\n\"Action\": [\"oss:*\"],\n\"Effect\": \"Allow\",\n\"Resource\": [\"acs:oss:*:*:*\"]\n}],\n\"Version\": \"1\"\n}`,
					"policy_name":     "${alicloud_resource_manager_policy.example.policy_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_document": CHECKSET,
						"policy_name":     CHECKSET,
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

var ResourceManagerPolicyVersionMap = map[string]string{
	"is_default_version": "false",
}

func ResourceManagerPolicyVersionBasicdependence(name string) string {
	return fmt.Sprintf(`
	
	resource "alicloud_resource_manager_policy" "example" {
	  policy_name     = "%s"
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
	
	`, name)
}
