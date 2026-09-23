package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRamPolicyDocumentDataSource0(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudRamPolicyDocumentDataSourceConfig0(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"oss:*\",\"Resource\":[\"acs:oss:*:*:myphotos\",\"acs:oss:*:*:myphotos/*\"]}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRamPolicyDocumentDataSource1(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudRamPolicyDocumentDataSourceConfig1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"oss:ListBuckets\",\"oss:GetBucketStat\",\"oss:GetBucketInfo\",\"oss:GetBucketTagging\",\"oss:GetBucketAcl\"],\"Resource\":\"acs:oss:*:*:*\"},{\"Effect\":\"Allow\",\"Action\":[\"oss:GetObject\",\"oss:GetObjectAcl\"],\"Resource\":\"acs:oss:*:*:myphotos/hangzhou/2015/*\"},{\"Effect\":\"Allow\",\"Action\":\"oss:ListObjects\",\"Resource\":\"acs:oss:*:*:myphotos\",\"Condition\":{\"StringLike\":{\"oss:Delimiter\":\"/\",\"oss:Prefix\":[\"hangzhou/\",\"hangzhou/2015/*\"]}}}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRamPolicyDocumentDataSource2(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudRamPolicyDocumentDataSourceConfig2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRamPolicyDocumentDataSource3(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudRamPolicyDocumentDataSourceConfig3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"sts:AssumeRole\",\"Principal\":{\"Service\":[\"ecs.aliyuncs.com\"]}}],\"Version\":\"1\"}",
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRamPolicyDocumentDataSource4(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudRamPolicyDocumentDataSourceConfig4(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": CHECKSET,
					}),
				),
			},
		},
	})
}

func testAccCheckAliCloudRamPolicyDocumentDataSourceConfig0() string {
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

func testAccCheckAliCloudRamPolicyDocumentDataSourceConfig1() string {
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
      values   = ["hangzhou/","hangzhou/2015/*"]
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

func testAccCheckAliCloudRamPolicyDocumentDataSourceConfig2() string {
	return fmt.Sprintf(`
data "alicloud_account" "default" {}
	
data "alicloud_ram_policy_document" "default" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
    principal {
      entity      = "RAM"
      identifiers = ["acs:ram::${data.alicloud_account.default.id}:root"]
    }
  }
}

resource "alicloud_ram_role" "role" {
  name     = "test-role-ram"
  document = data.alicloud_ram_policy_document.default.document
  force    = true
}
	`)
}

func testAccCheckAliCloudRamPolicyDocumentDataSourceConfig3() string {
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

func testAccCheckAliCloudRamPolicyDocumentDataSourceConfig4() string {
	return fmt.Sprintf(`
data "alicloud_account" "default" {}

data "alicloud_ram_policy_document" "default" {
  statement {
    effect = "Allow"
    action = ["sts:AssumeRole"]
    principal {
      entity      = "Federated"
      identifiers = ["acs:ram::${data.alicloud_account.default.id}:saml-provider/testprovider"]
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
	`)
}
