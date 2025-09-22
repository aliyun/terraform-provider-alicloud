package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudMaxComputeRolePolicyDocumentDataSource0(t *testing.T) {
	resourceId := "data.alicloud_max_compute_role_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudMaxComputeRolePolicyDocumentDataSourceConfig0(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"odps:*\"],\"Resource\":[\"acs:odps:*:projects/my_project/schemas/*\"]}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudMaxComputeRolePolicyDocumentDataSource1(t *testing.T) {
	resourceId := "data.alicloud_max_compute_role_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudMaxComputeRolePolicyDocumentDataSourceConfig1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"odps:CreateInstance\",\"odps:CreateSchema\",\"odps:List\"],\"Resource\":[\"acs:odps:*:projects/my_project\"]},{\"Effect\":\"Allow\",\"Action\":[\"odps:CreateTable\"],\"Resource\":[\"acs:odps:*:projects/my_project\",\"acs:odps:*:projects/my_project2\"]},{\"Effect\":\"Allow\",\"Action\":[\"odps:Read\",\"odps:Write\"],\"Resource\":[\"acs:odps:*:projects/my_project/instances/*\",\"acs:odps:*:projects/my_project2/instances/*\"]}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}


func testAccCheckAliCloudMaxComputeRolePolicyDocumentDataSourceConfig0() string {
	return fmt.Sprintf(`
data "alicloud_max_compute_role_policy_document" "default" {
  version = "1"
  statement {
    effect   = "Allow"
    action   = ["odps:*"]
    resource = ["acs:odps:*:projects/my_project/schemas/*"]
  }
}
	`)
}

func testAccCheckAliCloudMaxComputeRolePolicyDocumentDataSourceConfig1() string {
	return fmt.Sprintf(`
data "alicloud_max_compute_role_policy_document" "default" {
  version = "1"
  statement {
    effect	= "Allow"
    action   = ["odps:CreateInstance","odps:CreateSchema","odps:List"]
    resource = ["acs:odps:*:projects/my_project"]
  }
  statement {
    effect	= "Allow"
    action   = ["odps:CreateTable"]
    resource = ["acs:odps:*:projects/my_project", "acs:odps:*:projects/my_project2"]
  }
  statement {
    effect	= "Allow"
    action   = ["odps:Read","odps:Write"]
    resource = ["acs:odps:*:projects/my_project/instances/*", "acs:odps:*:projects/my_project2/instances/*"]
  }
}
	`)
}
