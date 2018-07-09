package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRamPoliciesDataSource_for_group(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPoliciessDataSourceForGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_policies.policy"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.0.name", "testAccCheckAlicloudRamPoliciessDataSourceForGroupConfig"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.0.type", "Custom"),
				),
			},
		},
	})
}

func TestAccAlicloudRamPoliciesDataSource_for_role(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPoliciessDataSourceForRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_policies.policy"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudRamPoliciesDataSource_for_user(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPoliciessDataSourceForUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_policies.policy"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudRamPoliciesDataSource_policy_name_regex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPoliciessDataSourcePolicyNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_policies.policy"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "1"),
				),
			},
		},
	})
}

func TestAccAlicloudRamPoliciesDataSource_policy_type(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPoliciessDataSourcePolicyTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_policies.policy"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "1"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRamPoliciessDataSourceForGroupConfig = `
variable "name" {
  default = "testAccCheckAlicloudRamPoliciessDataSourceForGroupConfig"
}
resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects",
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"]
    }]
  description = "this is a policy test"
  force = true
}

resource "alicloud_ram_group" "group" {
  name = "${var.name}"
  comments = "group comments"
  force=true
}

resource "alicloud_ram_group_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  group_name = "${alicloud_ram_group.group.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}
data "alicloud_ram_policies" "policy" {
  group_name = "${alicloud_ram_group_policy_attachment.attach.group_name}"
}`

const testAccCheckAlicloudRamPoliciessDataSourceForRoleConfig = `
variable "name" {
  default = "testAccCheckAlicloudRamPoliciessDataSourceForRoleConfig"
}
resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects",
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"]
    }]
  description = "this is a policy test"
  force = true
}

resource "alicloud_ram_role" "role" {
  name = "${var.name}"
  services = ["apigateway.aliyuncs.com", "ecs.aliyuncs.com"]
  description = "this is a test"
  force = true
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  role_name = "${alicloud_ram_role.role.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}
data "alicloud_ram_policies" "policy" {
  role_name = "${alicloud_ram_role_policy_attachment.attach.role_name}"
  type = "Custom"
}`

const testAccCheckAlicloudRamPoliciessDataSourceForUserConfig = `
variable "name" {
  default = "testAccCheckAlicloudRamPoliciessDataSourceForUserConfig"
}
resource "alicloud_ram_policy" "policy" {
  name = "${var.name}"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects",
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket",
        "acs:oss:*:*:mybucket/*"]
    }]
  description = "this is a policy test"
  force = true
}

resource "alicloud_ram_user" "user" {
  name = "${var.name}"
  display_name = "displayname"
  mobile = "86-18888888888"
  email = "hello.uuu@aaa.com"
  comments = "yoyoyo"
}

resource "alicloud_ram_user_policy_attachment" "attach" {
  policy_name = "${alicloud_ram_policy.policy.name}"
  user_name = "${alicloud_ram_user.user.name}"
  policy_type = "${alicloud_ram_policy.policy.type}"
}

data "alicloud_ram_policies" "policy" {
  user_name = "${alicloud_ram_user_policy_attachment.attach.user_name}"
}`

const testAccCheckAlicloudRamPoliciessDataSourcePolicyNameRegexConfig = `
resource "alicloud_ram_policy" "policy" {
  name = "testAccCheckAlicloudRamPoliciessDataSourcePolicyNameRegexConfig"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket"]
    }]
  description = "this is a policy test"
  force = true
}
data "alicloud_ram_policies" "policy" {
  name_regex = "${alicloud_ram_policy.policy.name}"
}`

const testAccCheckAlicloudRamPoliciessDataSourcePolicyTypeConfig = `
resource "alicloud_ram_policy" "policy" {
  name = "testAccCheckAlicloudRamPoliciessDataSourcePolicyNameRegexConfig"
  statement = [
    {
      effect = "Deny"
      action = [
        "oss:ListObjects"]
      resource = [
        "acs:oss:*:*:mybucket"]
    }]
  description = "this is a policy test"
  force = true
}
data "alicloud_ram_policies" "policy" {
  name_regex = "${alicloud_ram_policy.policy.name}"
  type = "Custom"
}`
