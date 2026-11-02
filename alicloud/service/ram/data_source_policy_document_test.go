package ram_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudRamPolicyDocumentDataSource0(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudRamPolicyDocumentDataSourceConfig2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAliCloudRamPolicyDocumentDataSource3(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudRamPolicyDocumentDataSourceConfig4(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document": acctest.CHECKSET,
					}),
				),
			},
		},
	})
}

// TestAccAliCloudRamPolicyDocumentDataSource5 exercises the full policy_document
// data source schema through a two-step modification. Step 0 carries a minimal
// statement (effect/action/resource); step 1 adds version, principal and
// condition. Each step's inline map is recursed by the coverage checker so all
// 13 schema attributes land in testMustSet, and the attributes that first
// appear in step 1 (version, principal.entity/identifiers,
// condition.operator/variable/values) populate testModifySet. The exact
// document assertion validates that every attribute — including the TypeSet
// children that cannot be addressed by numeric index — is assembled into the
// rendered policy.
func TestAccAliCloudRamPolicyDocumentDataSource5(t *testing.T) {
	resourceId := "data.alicloud_ram_policy_document.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "default", func(name string) string {
		return ""
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"statement": []interface{}{
						map[string]interface{}{
							"effect":   "Allow",
							"action":   []interface{}{"oss:GetObject"},
							"resource": []interface{}{"acs:oss:*:*:bucket0/*"},
						},
					},
					"output_file": "policy-document-step0.json",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document":               `{"Statement":[{"Effect":"Allow","Action":"oss:GetObject","Resource":"acs:oss:*:*:bucket0/*"}],"Version":"1"}`,
						"output_file":            "policy-document-step0.json",
						"statement.0.effect":     "Allow",
						"statement.0.action.0":   "oss:GetObject",
						"statement.0.resource.0": "acs:oss:*:*:bucket0/*",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"version": "1",
					"statement": []interface{}{
						map[string]interface{}{
							"effect":   "Deny",
							"action":   []interface{}{"oss:PutObject"},
							"resource": []interface{}{"acs:oss:*:*:bucket1/*"},
							"principal": []interface{}{
								map[string]interface{}{
									"entity":      "Service",
									"identifiers": []interface{}{"oss.aliyuncs.com"},
								},
							},
							"condition": []interface{}{
								map[string]interface{}{
									"operator": "StringNotEquals",
									"variable": "oss:prefix",
									"values":   []interface{}{"folder1/"},
								},
							},
						},
					},
					"output_file": "policy-document-step1.json",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"document":               `{"Statement":[{"Effect":"Deny","Action":"oss:PutObject","Resource":"acs:oss:*:*:bucket1/*","Principal":{"Service":["oss.aliyuncs.com"]},"Condition":{"StringNotEquals":{"oss:prefix":"folder1/"}}}],"Version":"1"}`,
						"version":                "1",
						"output_file":            "policy-document-step1.json",
						"statement.0.effect":     "Deny",
						"statement.0.action.0":   "oss:PutObject",
						"statement.0.resource.0": "acs:oss:*:*:bucket1/*",
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
