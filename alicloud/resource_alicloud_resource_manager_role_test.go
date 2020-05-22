package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/resourcemanager"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudResourceManagerRole_basic(t *testing.T) {
	var v resourcemanager.Role
	resourceId := "alicloud_resource_manager_role.default"
	ra := resourceAttrInit(resourceId, ResourceManagerRoleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourcemanagerService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerRole")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccResourceManagerRole%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ResourceManagerRoleBasicdependence)
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
					"assume_role_policy_document": `{\n    \"Statement\": [{\n            \"Action\": \"sts:AssumeRole\",\n            \"Effect\": \"Allow\",\n            \"Principal\": {\"RAM\": [\"acs:ram::${data.alicloud_account.default.id}:root\"]}}],\n    \"Version\": \"1\"}`,
					"role_name":                   name,
					"description":                 "Test resourceManager role",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"assume_role_policy_document": CHECKSET,
						"role_name":                   name,
						"description":                 "Test resourceManager role",
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
					"assume_role_policy_document": `{\n    \"Statement\": [{\n            \"Action\": \"sts:AssumeRole\",\n            \"Effect\": \"Deny\",\n            \"Principal\": {\"RAM\": [\"acs:ram::${data.alicloud_account.default.id}:root\"]}}],\n    \"Version\": \"1\"}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"assume_role_policy_document": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"max_session_duration": "3600",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"max_session_duration": "3600",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"assume_role_policy_document": `{\n    \"Statement\": [{\n            \"Action\": \"sts:AssumeRole\",\n            \"Effect\": \"Allow\",\n            \"Principal\": {\"RAM\": [\"acs:ram::${data.alicloud_account.default.id}:root\"]}}],\n    \"Version\": \"1\"}`,
					"role_name":                   name,
					"description":                 "Test resourceManager role",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"assume_role_policy_document": CHECKSET,
						"role_name":                   name,
						"description":                 "Test resourceManager role",
					}),
				),
			},
		},
	})
}

var ResourceManagerRoleMap = map[string]string{}

func ResourceManagerRoleBasicdependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_account" "default" {}
`)
}
