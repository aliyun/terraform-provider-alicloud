package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudROSChangeSet_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ros_change_set.default"
	ra := resourceAttrInit(resourceId, AlicloudRosChangeSetMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &RosService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRosChangeSet")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccRosChangeSet%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudRosChangeSetBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"change_set_name": name,
					"stack_name":      name + "stack",
					"change_set_type": "CREATE",
					"description":     "Test From Terraform",
					"template_body":   `{\"ROSTemplateFormatVersion\":\"2015-09-01\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"change_set_name": name,
						"stack_name":      name + "stack",
						"change_set_type": "CREATE",
						"description":     "Test From Terraform",
						"template_body":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"notification_urls", "ram_role_name", "replacement_option", "stack_policy_body", "stack_policy_during_update_body", "stack_policy_during_update_url", "stack_policy_url", "template_url", "use_previous_parameters", "template_body"},
			},
		},
	})
}

var AlicloudRosChangeSetMap = map[string]string{
	"status": CHECKSET,
}

func AlicloudRosChangeSetBasicDependence(name string) string {
	return ""
}
