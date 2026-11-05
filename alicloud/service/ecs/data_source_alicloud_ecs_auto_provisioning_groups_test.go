package ecs_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The query is exercised against the auto provisioning groups that already
// exist in the test account: the assertions prove a non-empty result and check
// that the returned IDs and names belong to the same group, and the follow-up
// steps prove that the ID and name filters select exactly the target group.
// Creating a fixture through alicloud_auto_provisioning_group is not used here:
// that resource drops a just-created group from state while the service is
// still converging, which fails the apply for reasons outside this data source.
func TestAccAliCloudECSAutoProvisioningGroupsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_auto_provisioning_groups.default"
	allResourceId := "data.alicloud_ecs_auto_provisioning_groups.all"

	checkNonEmpty := func(name, key string) resource.TestCheckFunc {
		return resource.TestCheckResourceAttrWith(resourceId, key, func(value string) error {
			count, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("%s is not a count: %w", key, err)
			}
			if count < 1 {
				return fmt.Errorf("%s expected at least one auto provisioning group, got %d", name, count)
			}
			return nil
		})
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: autoProvisioningGroupsConfig(""),
				Check: resource.ComposeTestCheckFunc(
					checkNonEmpty("all groups", "groups.#"),
					checkNonEmpty("all ids", "ids.#"),
					checkNonEmpty("all names", "names.#"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", resourceId, "groups.0.auto_provisioning_group_id"),
					resource.TestCheckResourceAttrPair(resourceId, "names.0", resourceId, "groups.0.auto_provisioning_group_name"),
					resource.TestCheckResourceAttrSet(resourceId, "groups.0.status"),
				),
			},
			{
				Config: autoProvisioningGroupsConfig("  ids = [data.alicloud_ecs_auto_provisioning_groups.all.groups.0.auto_provisioning_group_id]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", allResourceId, "groups.0.auto_provisioning_group_id"),
				),
			},
			{
				Config: autoProvisioningGroupsConfig("  auto_provisioning_group_name = data.alicloud_ecs_auto_provisioning_groups.all.groups.0.auto_provisioning_group_name\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttrPair(resourceId, "names.0", allResourceId, "groups.0.auto_provisioning_group_name"),
				),
			},
			{
				Config: autoProvisioningGroupsConfig("  ids = [\"apg-does-not-exist\"]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      autoProvisioningGroupsConfig("  types = [\"not-a-type\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      autoProvisioningGroupsConfig("  status = [\"not-a-status\"]\n"),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match|expected value to be one of`),
			},
			{
				Config:      autoProvisioningGroupsConfig("  tags = { \"env\" = \"https://example.com\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must not contain a URL`),
			},
			{
				Config:      autoProvisioningGroupsConfig("  tags = { \"aliyun-env\" = \"test\" }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Key|must not start with`),
			},
			{
				// The destroy phase re-reads the last configuration, so a valid
				// step must follow the expected validation failures.
				Config: autoProvisioningGroupsConfig(""),
				Check: resource.ComposeTestCheckFunc(
					checkNonEmpty("all groups", "groups.#"),
				),
			},
		},
	})
}

func autoProvisioningGroupsConfig(dataArguments string) string {
	return fmt.Sprintf(`
data "alicloud_ecs_auto_provisioning_groups" "all" {
}

data "alicloud_ecs_auto_provisioning_groups" "default" {
%s}
`, dataArguments)
}
