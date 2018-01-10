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
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.0.policy_name", "ReadOnlyAccess"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.0.policy_type", "System"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.0.description", "只读访问所有阿里云资源的权限"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.0.default_version", "v2"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.0.policy_document", "{\n    \"Version\": \"1\", \n    \"Statement\": [\n        {\n            \"Action\": [\n                \"*:Describe*\", \n                \"*:List*\", \n                \"*:Get*\", \n                \"*:BatchGet*\", \n                \"*:Query*\", \n                \"*:BatchQuery*\", \n                \"actiontrail:LookupEvents\", \n                \"dm:Desc*\", \n                \"dm:SenderStatistics*\"\n            ], \n            \"Resource\": \"*\", \n            \"Effect\": \"Allow\"\n        }\n    ]\n}"),
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
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "2"),
				),
			},
		},
	})
}

func TestAccAlicloudRamPoliciesDataSource_for_all(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRamPoliciessDataSourceForAllConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ram_policies.policy"),
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "131"),
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
					resource.TestCheckResourceAttr("data.alicloud_ram_policies.policy", "policies.#", "74"),
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
data "alicloud_ram_policies" "policy" {
  group_name = "group2"
}`

const testAccCheckAlicloudRamPoliciessDataSourceForRoleConfig = `
data "alicloud_ram_policies" "policy" {
  role_name = "testrole"
  type = "Custom"
}`

const testAccCheckAlicloudRamPoliciessDataSourceForUserConfig = `
data "alicloud_ram_policies" "policy" {
  user_name = "user1"
}`

const testAccCheckAlicloudRamPoliciessDataSourceForAllConfig = `
data "alicloud_ram_policies" "policy" {
}`

const testAccCheckAlicloudRamPoliciessDataSourcePolicyNameRegexConfig = `
data "alicloud_ram_policies" "policy" {
  name_regex = ".*Full.*"
}`

const testAccCheckAlicloudRamPoliciessDataSourcePolicyTypeConfig = `
data "alicloud_ram_policies" "policy" {
  type = "Custom"
}`
