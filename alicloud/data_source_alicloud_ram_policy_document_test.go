package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudRamPolicyDocumentDataSource0(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPolicyDocumentDataSourceConfig0(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\",\"acs:oss:*:*:myphotos/*\"]}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamPolicyDocumentDataSource1(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPolicyDocumentDataSourceConfig1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"oss:ListBuckets\",\"oss:GetBucketStat\",\"oss:GetBucketInfo\",\"oss:GetBucketTagging\",\"oss:GetBucketAcl\"],\"Resource\":\"acs:oss:*:*:*\"},{\"Effect\":\"Allow\",\"Action\":[\"oss:GetObject\",\"oss:GetObjectAcl\"],\"Resource\":\"acs:oss:*:*:myphotos/hangzhou/2015/*\"},{\"Effect\":\"Allow\",\"Action\":\"oss:ListObjects\",\"Resource\":\"acs:oss:*:*:myphotos\",\"Condition\":{\"StringLike\":{\"oss:Delimiter\":\"/\",\"oss:Prefix\":[\"\",\"hangzhou/\",\"hangzhou/2015/*\"]}}}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamPolicyDocumentDataSource2(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPolicyDocumentDataSourceConfig2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": fmt.Sprintf("{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"sts:AssumeRole\",\"Principal\":{\"RAM\":[\"acs:ram::%s:root\"]}}],\"Version\":\"1\"}", os.Getenv("ALICLOUD_ACCOUNT_ID")),
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamPolicyDocumentDataSource3(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPolicyDocumentDataSourceConfig3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"sts:AssumeRole\",\"Principal\":{\"Service\":[\"ecs.aliyuncs.com\"]}}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudRamPolicyDocumentDataSource4(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPolicyDocumentDataSourceConfig4(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"sts:AssumeRole\",\"Principal\":{\"Federated\":[\"acs:ram::1511928242963727:saml-provider/testprovider\"]},\"Condition\":{\"StringEquals\":{\"saml:recipient\":\"https://signin.aliyun.com/saml-role/sso\"}}}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func testAccCheckAlicloudRamPolicyDocumentDataSourceConfig0() string {
	return fmt.Sprintf(`
data "alicloud_ram_policy_document" "default" {
  version = "1"
  statement {
    effect   = "Allow"
    action   = ["oss:*"]
    resource = ["acs:oss:*:*:myphotos", "acs:oss:*:*:myphotos/*"]
  }
}

resource "alicloud_ram_policy" "default" {
  policy_name     = "tf-test-no-condition"
  policy_document = data.alicloud_ram_policy_document.default.document
  force           = true
}
	`)
}

func testAccCheckAlicloudRamPolicyDocumentDataSourceConfig1() string {
	return fmt.Sprintf(`
data "alicloud_ram_policy_document" "default" {
  version = "1"
  statement {
    effect	= "Allow"
    action   = ["oss:ListBuckets","oss:GetBucketStat","oss:GetBucketInfo","oss:GetBucketTagging","oss:GetBucketAcl"]
    resource = ["acs:oss:*:*:*"]
  }
  statement {
    effect	= "Allow"
    action   = ["oss:GetObject","oss:GetObjectAcl"]
    resource = ["acs:oss:*:*:myphotos/hangzhou/2015/*"]
  }
  statement {
    effect	= "Allow"
    action   = ["oss:ListObjects"]
    resource = ["acs:oss:*:*:myphotos"]
    condition {
      operator     = "StringLike"
      variable = "oss:Delimiter"
      values   = ["/"]
    }
    condition {
      operator     = "StringLike"
      variable = "oss:Prefix"
      values   = ["","hangzhou/","hangzhou/2015/*"]
    }
  }
}

resource "alicloud_ram_policy" "policy" {
  policy_name     = "tf-test-condition"
  policy_document = data.alicloud_ram_policy_document.default.document
  force           = true
}
	`)
}

func testAccCheckAlicloudRamPolicyDocumentDataSourceConfig2() string {
	return fmt.Sprintf(`
data "alicloud_ram_policy_document" "default" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
    principal {
      entity      = "RAM"
      identifiers = ["acs:ram::%s:root"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-ram"
  document = data.alicloud_ram_policy_document.default.document
  force    = true
}
	`, os.Getenv("ALICLOUD_ACCOUNT_ID"))
}

func testAccCheckAlicloudRamPolicyDocumentDataSourceConfig3() string {
	return fmt.Sprintf(`
data "alicloud_ram_policy_document" "default" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
    principal {
      entity      = "Service"
      identifiers = ["ecs.aliyuncs.com"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-service"
  document = data.alicloud_ram_policy_document.default.document
  force    = true
}
	`)
}

func testAccCheckAlicloudRamPolicyDocumentDataSourceConfig4() string {
	return fmt.Sprintf(`
data "alicloud_ram_policy_document" "default" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
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
  document = data.alicloud_ram_policy_document.default.document
  force    = true
}
	`, os.Getenv("ALICLOUD_ACCOUNT_ID"))
}
