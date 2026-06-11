package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test MaxCompute RoleUserAttachment. >>> Resource test cases, automatically generated.
// Case RoleUserAttachment_terraform测试(RamRole) 9962
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9962(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9962)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9962)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${alicloud_max_compute_role.default.role_name}",
					"user":         "${local.user}",
					"project_name": "${alicloud_maxcompute_project.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudMaxComputeRoleUserAttachmentMap9962 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9962(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "account_name" {
  # canonical MaxCompute account name of the ACube test account; RAM GetAccountAlias
  # returns a different alias, so this cannot be derived from a data source
  default = "openapiautomation_testcloud_com"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

resource "alicloud_max_compute_role" "default" {
  project_name = alicloud_maxcompute_project.default.id
  role_name    = "role_project_admin"
  type         = "admin"
  policy       = jsonencode({
    Statement = [{
      Action   = ["odps:*"]
      Effect   = "Allow"
      Resource = ["acs:odps:*:projects/project_name/authorization/roles", "acs:odps:*:projects/project_name/authorization/roles/*/*"]
    }]
    Version = "1"
  })
}

resource "alicloud_ram_role" "default" {
  name     = replace(var.name, "_", "-")
  document = <<EOF
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": ["ecs.aliyuncs.com"]
      }
    }
  ],
  "Version": "1"
}
EOF
  force    = true
}

locals {
  user = format("RAM$%%s:role/%%s", var.account_name, alicloud_ram_role.default.name)
}


`, name)
}

// Case RoleUserAttachment_terraform测试(RamUser) 9961
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9961(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9961)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9961)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${alicloud_max_compute_role.default.role_name}",
					"user":         "${local.user}",
					"project_name": "${alicloud_maxcompute_project.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudMaxComputeRoleUserAttachmentMap9961 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9961(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "account_name" {
  # canonical MaxCompute account name of the ACube test account; RAM GetAccountAlias
  # returns a different alias, so this cannot be derived from a data source
  default = "openapiautomation_testcloud_com"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

resource "alicloud_max_compute_role" "default" {
  project_name = alicloud_maxcompute_project.default.id
  role_name    = "role_project_admin"
  type         = "admin"
  policy       = jsonencode({
    Statement = [{
      Action   = ["odps:*"]
      Effect   = "Allow"
      Resource = ["acs:odps:*:projects/project_name/authorization/roles", "acs:odps:*:projects/project_name/authorization/roles/*/*"]
    }]
    Version = "1"
  })
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

locals {
  user = format("RAM$%%s:%%s", var.account_name, alicloud_ram_user.default.name)
}


`, name)
}

// Case RoleUserAttachment_terraform测试(adminRole) 9966
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9966(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9966)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9966)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "admin",
					"user":         "${local.user}",
					"project_name": "${alicloud_maxcompute_project.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    "admin",
						"user":         CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudMaxComputeRoleUserAttachmentMap9966 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9966(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "account_name" {
  # canonical MaxCompute account name of the ACube test account; RAM GetAccountAlias
  # returns a different alias, so this cannot be derived from a data source
  default = "openapiautomation_testcloud_com"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

locals {
  user = format("RAM$%%s:%%s", var.account_name, alicloud_ram_user.default.name)
}


`, name)
}

// Case RoleUserAttachment_terraform测试(super_administratorRole) 9967
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9967(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9967)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9967)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${alicloud_max_compute_role.default.role_name}",
					"user":         "${local.user}",
					"project_name": "${alicloud_maxcompute_project.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudMaxComputeRoleUserAttachmentMap9967 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9967(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "account_name" {
  # canonical MaxCompute account name of the ACube test account; RAM GetAccountAlias
  # returns a different alias, so this cannot be derived from a data source
  default = "openapiautomation_testcloud_com"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

resource "alicloud_max_compute_role" "default" {
  project_name = alicloud_maxcompute_project.default.id
  role_name    = "role_project_admin"
  type         = "admin"
  policy       = jsonencode({
    Statement = [{
      Action   = ["odps:*"]
      Effect   = "Allow"
      Resource = ["acs:odps:*:projects/project_name/authorization/roles", "acs:odps:*:projects/project_name/authorization/roles/*/*"]
    }]
    Version = "1"
  })
}

resource "alicloud_ram_user" "default" {
  name = var.name
}

locals {
  user = format("RAM$%%s:%%s", var.account_name, alicloud_ram_user.default.name)
}


`, name)
}

// Case RoleUserAttachment_terraform测试(AliyunUser) 9960
func TestAccAliCloudMaxComputeRoleUserAttachment_basic9960(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_max_compute_role_user_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudMaxComputeRoleUserAttachmentMap9960)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MaxComputeServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMaxComputeRoleUserAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMaxComputeRoleUserAttachmentBasicDependence9960)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"role_name":    "${alicloud_max_compute_role.default.role_name}",
					"user":         "${local.user}",
					"project_name": "${alicloud_maxcompute_project.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"role_name":    CHECKSET,
						"user":         CHECKSET,
						"project_name": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudMaxComputeRoleUserAttachmentMap9960 = map[string]string{}

func AlicloudMaxComputeRoleUserAttachmentBasicDependence9960(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "account_name" {
  # canonical MaxCompute account name of the ACube test account; RAM GetAccountAlias
  # returns a different alias, so this cannot be derived from a data source
  default = "openapiautomation_testcloud_com"
}

resource "alicloud_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

resource "alicloud_max_compute_role" "default" {
  project_name = alicloud_maxcompute_project.default.id
  role_name    = "role_project_admin"
  type         = "admin"
  policy       = jsonencode({
    Statement = [{
      Action   = ["odps:*"]
      Effect   = "Allow"
      Resource = ["acs:odps:*:projects/project_name/authorization/roles", "acs:odps:*:projects/project_name/authorization/roles/*/*"]
    }]
    Version = "1"
  })
}

locals {
  user = format("ALIYUN$%%s", var.account_name)
}


`, name)
}

// Test MaxCompute RoleUserAttachment. <<< Resource test cases, automatically generated.
