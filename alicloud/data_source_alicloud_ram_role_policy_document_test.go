package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"os"
	"testing"
)

func TestAccAlicloudRamRolePolicyDocumentDataSource00(t *testing.T) {
	resourceId := "data.alicloud_ram_role_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolePolicyDocumentDataSourceConfig00(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": fmt.Sprintf("{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"sts:AssumeRole\",\"Principal\":{\"RAM\":[\"acs:ram::%s:root\"]}}],\"Version\":\"1\"}", os.Getenv("ALICLOUD_ACCOUNT_ID")),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamRolePolicyDocumentDataSource01(t *testing.T) {
	resourceId := "data.alicloud_ram_role_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolePolicyDocumentDataSourceConfig01(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"sts:AssumeRole\",\"Principal\":{\"Service\":[\"ecs.aliyuncs.com\"]}}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamRolePolicyDocumentDataSource02(t *testing.T) {
	resourceId := "data.alicloud_ram_role_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamRolePolicyDocumentDataSourceConfig02(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"sts:AssumeRole\",\"Principal\":{\"Federated\":[\"acs:ram::1511928242963727:saml-provider/testprovider\"]},\"Condition\":{\"StringEquals\":{\"saml:recipient\":\"https://signin.aliyun.com/saml-role/sso\"}}}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func testAccCheckAlicloudRamRolePolicyDocumentDataSourceConfig00() string {
	return fmt.Sprintf(`
data "alicloud_ram_role_policy_document" "default" {
  statement {
    effect = "Allow"
    action = "sts:AssumeRole"
    principal {
      entity      = "RAM"
      identifiers = ["acs:ram::%s:root"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-ram"
  document = data.alicloud_ram_role_policy_document.default.document
  force    = true
}
	`, os.Getenv("ALICLOUD_ACCOUNT_ID"))
}

func testAccCheckAlicloudRamRolePolicyDocumentDataSourceConfig01() string {
	return fmt.Sprintf(`
data "alicloud_ram_role_policy_document" "default" {
  statement {
    effect = "Allow"
    action = "sts:AssumeRole"
    principal {
      entity      = "Service"
      identifiers = ["ecs.aliyuncs.com"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-service"
  document = data.alicloud_ram_role_policy_document.default.document
  force    = true
}
	`)
}

func testAccCheckAlicloudRamRolePolicyDocumentDataSourceConfig02() string {
	return fmt.Sprintf(`
data "alicloud_ram_role_policy_document" "default" {
  statement {
    effect = "Allow"
    action = "sts:AssumeRole"
    principal {
      entity      = "Federated"
      identifiers = ["acs:ram::%s:saml-provider/testprovider"]
    }
    condition {
      operator = "StringEquals"
      variable = "saml:recipient"
      values   = ["https://signin.aliyun.com/saml-role/sso"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-federated"
  document = data.alicloud_ram_role_policy_document.default.document
  force    = true
}
	`, os.Getenv("ALICLOUD_ACCOUNT_ID"))
}
