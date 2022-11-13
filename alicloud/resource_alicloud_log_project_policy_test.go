package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogProjectPolicy_basic(t *testing.T) {
	var policy string
	resourceId := "alicloud_log_project_policy.default"
	ra := resourceAttrInit(resourceId, logProjectPolicyMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &policy, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogprojectpolicy-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogProjectPolicyDependence)

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
					"project": "${alicloud_log_project.default.name}",
					"policy":  `{\"Statement\":[{\"Action\":[\"log:Post*\"],\"Effect\":\"Deny\",\"Resource\":\"acs:log:*:*:project/test-project/*\"}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project": name,
						"policy":  "{\"Statement\":[{\"Action\":[\"log:Post*\"],\"Effect\":\"Deny\",\"Resource\":\"acs:log:*:*:project/test-project/*\"}],\"Version\":\"1\"}",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project": "${alicloud_log_project.default.name}",
					"policy":  `{\"Statement\":[{\"Action\":[\"log:Post*\"],\"Effect\":\"Allow\",\"Resource\":\"acs:log:*:*:*\"}],\"Version\":\"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project": name,
						"policy":  "{\"Statement\":[{\"Action\":[\"log:Post*\"],\"Effect\":\"Allow\",\"Resource\":\"acs:log:*:*:*\"}],\"Version\":\"1\"}",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"project": "${alicloud_log_project.default.name}",
					"policy":  ``,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project": name,
						"policy":  "",
					}),
				),
			},
		},
	})
}

var logProjectPolicyMap = map[string]string{}

func resourceLogProjectPolicyDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "default" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	`, name)
}
