package ecs_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
				// A repeated input ID is sent once, so the result and the count
				// do not depend on the server's duplicate handling.
				Config: autoProvisioningGroupsConfig("  ids = [data.alicloud_ecs_auto_provisioning_groups.all.groups.0.auto_provisioning_group_id, data.alicloud_ecs_auto_provisioning_groups.all.groups.0.auto_provisioning_group_id]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "ids.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "1"),
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
				// DescribeAutoProvisioningGroups documents a maximum of 20 IDs
				// even though the published schema metadata still declares 100.
				Config: autoProvisioningGroupsConfig(fmt.Sprintf("  ids = [%s]\n", autoProvisioningGroupsIdList(20))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				Config:      autoProvisioningGroupsConfig(fmt.Sprintf("  ids = [%s]\n", autoProvisioningGroupsIdList(21))),
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Size|at most 20`),
			},
			{
				// The name filter has to narrow the result to the known target:
				// every returned group must carry that name, and the count must
				// stay below the unfiltered count when the account holds more
				// than one group. A filter the server ignored, or one the
				// provider dropped, fails these assertions.
				Config: autoProvisioningGroupsConfig("  auto_provisioning_group_name = data.alicloud_ecs_auto_provisioning_groups.all.groups.0.auto_provisioning_group_name\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", allResourceId, "groups.0.auto_provisioning_group_id"),
					checkAutoProvisioningGroupsNameFilter(resourceId, allResourceId),
				),
			},
			{
				// The type filter must narrow the result as the request asks.
				// One read per documented type is compared with the unfiltered
				// read: every row carries the requested type, no group is
				// returned by two different type filters, and the three reads
				// together cover the unfiltered result exactly. A server that
				// ignores the filter returns every group for each type, so the
				// sum and the overlap assertions fail even when the account
				// holds groups of a single type.
				Config: autoProvisioningGroupsTypesConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkAutoProvisioningGroupsTypesPartition(allResourceId, map[string]string{
						"request":   "data.alicloud_ecs_auto_provisioning_groups.type_request",
						"maintain":  "data.alicloud_ecs_auto_provisioning_groups.type_maintain",
						"candidate": "data.alicloud_ecs_auto_provisioning_groups.type_candidate",
					}),
				),
			},
			{
				// A name that cannot match must not return the groups that do.
				Config: autoProvisioningGroupsConfig("  auto_provisioning_group_name = \"apg-does-not-exist\"\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "groups.#", "0"),
					resource.TestCheckResourceAttr(resourceId, "total_count", "0"),
				),
			},
			{
				// Reading from the second page must return the rows that follow
				// the first page, measured against the readable suffix of the
				// total: the two pages together cover the complete result once.
				Config: autoProvisioningGroupsPagedConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkAutoProvisioningGroupsLaterPage(allResourceId, "data.alicloud_ecs_auto_provisioning_groups.first", "data.alicloud_ecs_auto_provisioning_groups.second"),
				),
			},
			{
				// A null value is not the same request as an empty value.
				Config:      autoProvisioningGroupsConfig("  tags = { env = null }\n"),
				ExpectError: regexp.MustCompile(`Invalid Tag Value|must be set`),
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

func autoProvisioningGroupsIdList(count int) string {
	ids := make([]string, 0, count)
	for i := 0; i < count; i++ {
		ids = append(ids, fmt.Sprintf("%q", fmt.Sprintf("apg-does-not-exist-%d", i)))
	}
	return strings.Join(ids, ", ")
}

func autoProvisioningGroupsConfig(dataArguments string) string {
	return fmt.Sprintf(`
data "alicloud_ecs_auto_provisioning_groups" "all" {
}

data "alicloud_ecs_auto_provisioning_groups" "default" {
%s}
`, dataArguments)
}

// autoProvisioningGroupsPagedConfig declares the unfiltered read together with
// two reads that start at page 1 and page 2 with a page size of one, so the
// later-page result can be compared against the complete result.
func autoProvisioningGroupsPagedConfig() string {
	return `
data "alicloud_ecs_auto_provisioning_groups" "all" {
}

data "alicloud_ecs_auto_provisioning_groups" "first" {
  page_number = 1
  page_size   = 1
}

data "alicloud_ecs_auto_provisioning_groups" "second" {
  page_number = 2
  page_size   = 1
}
`
}

func autoProvisioningGroupsStateAttribute(state *terraform.State, resourceId, key string) (string, error) {
	instance, ok := state.RootModule().Resources[resourceId]
	if !ok || instance.Primary == nil {
		return "", fmt.Errorf("%s is missing from the state", resourceId)
	}
	value, ok := instance.Primary.Attributes[key]
	if !ok {
		return "", fmt.Errorf("%s.%s is missing from the state", resourceId, key)
	}
	return value, nil
}

func autoProvisioningGroupsStateCount(state *terraform.State, resourceId, key string) (int, error) {
	value, err := autoProvisioningGroupsStateAttribute(state, resourceId, key)
	if err != nil {
		return 0, err
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s.%s is not a count: %w", resourceId, key, err)
	}
	return count, nil
}

// autoProvisioningGroupsTypesConfig declares the unfiltered read together with
// one read per documented auto provisioning group type. The type filter cannot
// be derived from the account contents, so the partition of the unfiltered
// result is what proves the filter narrows the query.
func autoProvisioningGroupsTypesConfig() string {
	return `
data "alicloud_ecs_auto_provisioning_groups" "all" {
}

data "alicloud_ecs_auto_provisioning_groups" "type_request" {
  types = ["request"]
}

data "alicloud_ecs_auto_provisioning_groups" "type_maintain" {
  types = ["maintain"]
}

data "alicloud_ecs_auto_provisioning_groups" "type_candidate" {
  types = ["candidate"]
}
`
}

// checkAutoProvisioningGroupsTypesPartition proves the type filter narrows the
// result without depending on the account holding more than one type: every row
// of each type-filtered read carries the requested type, no group appears under
// two types, the reads together cover the unfiltered result, and their row
// counts add up to it. A server that ignores `types` returns the full result for
// each of the three reads, which fails the overlap and the sum assertions.
func checkAutoProvisioningGroupsTypesPartition(allResourceId string, typeResources map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		allCount, err := autoProvisioningGroupsStateCount(state, allResourceId, "groups.#")
		if err != nil {
			return err
		}
		allIds := make(map[string]struct{}, allCount)
		for index := 0; index < allCount; index++ {
			id, err := autoProvisioningGroupsStateAttribute(state, allResourceId, fmt.Sprintf("ids.%d", index))
			if err != nil {
				return err
			}
			allIds[id] = struct{}{}
		}

		seenTypes := make(map[string]string, allCount)
		filteredRows := 0
		for groupType, resourceId := range typeResources {
			count, err := autoProvisioningGroupsStateCount(state, resourceId, "groups.#")
			if err != nil {
				return err
			}
			filteredRows += count
			for index := 0; index < count; index++ {
				returnedType, err := autoProvisioningGroupsStateAttribute(state, resourceId, fmt.Sprintf("groups.%d.auto_provisioning_group_type", index))
				if err != nil {
					return err
				}
				if returnedType != groupType {
					return fmt.Errorf("%s returned group %d of type %q for the type filter %q", resourceId, index, returnedType, groupType)
				}
				id, err := autoProvisioningGroupsStateAttribute(state, resourceId, fmt.Sprintf("ids.%d", index))
				if err != nil {
					return err
				}
				if previousType, duplicate := seenTypes[id]; duplicate {
					return fmt.Errorf("group %q was returned by both the %q and the %q type filters", id, previousType, groupType)
				}
				seenTypes[id] = groupType
				if _, ok := allIds[id]; !ok {
					return fmt.Errorf("%s returned %q, which is not in the unfiltered result", resourceId, id)
				}
			}
		}
		if filteredRows != allCount {
			return fmt.Errorf("the three type filters returned %d rows in total while the unfiltered read returned %d; a filter that the server ignored duplicates rows", filteredRows, allCount)
		}
		if len(seenTypes) != allCount {
			return fmt.Errorf("the three type filters covered %d groups while the unfiltered read returned %d", len(seenTypes), allCount)
		}
		return nil
	}
}

// checkAutoProvisioningGroupsNameFilter proves the name filter actually narrows
// the result: every row carries the target name, and when the account holds more
// than one group the filtered count stays below the unfiltered one.
func checkAutoProvisioningGroupsNameFilter(resourceId, allResourceId string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		allCount, err := autoProvisioningGroupsStateCount(state, allResourceId, "groups.#")
		if err != nil {
			return err
		}
		target, err := autoProvisioningGroupsStateAttribute(state, allResourceId, "groups.0.auto_provisioning_group_name")
		if err != nil {
			return err
		}
		if target == "" {
			return fmt.Errorf("%s did not return a group to filter by", allResourceId)
		}
		count, err := autoProvisioningGroupsStateCount(state, resourceId, "groups.#")
		if err != nil {
			return err
		}
		if count < 1 {
			return fmt.Errorf("the name filter %q matched no group", target)
		}
		if count > allCount {
			return fmt.Errorf("the name filter returned %d groups while only %d exist", count, allCount)
		}
		if allCount > 1 && count >= allCount {
			return fmt.Errorf("the name filter did not narrow %d groups down", allCount)
		}
		for index := 0; index < count; index++ {
			name, err := autoProvisioningGroupsStateAttribute(state, resourceId, fmt.Sprintf("groups.%d.auto_provisioning_group_name", index))
			if err != nil {
				return err
			}
			if name != target {
				return fmt.Errorf("group %d is %q, want the filtered name %q", index, name, target)
			}
		}
		return nil
	}
}

// checkAutoProvisioningGroupsLaterPage proves the read starts at the configured
// page and measures completeness against the readable suffix of the total. A
// read with page_size = 1 still follows every page, so page 1 returns the whole
// result while page 2 returns exactly the rows the server did not return first.
// It does not assume the server returns the groups in ID order.
func checkAutoProvisioningGroupsLaterPage(allResourceId, pageOneId, pageTwoId string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		total, err := autoProvisioningGroupsStateCount(state, allResourceId, "groups.#")
		if err != nil {
			return err
		}
		for _, pageId := range []string{pageOneId, pageTwoId} {
			reported, err := autoProvisioningGroupsStateCount(state, pageId, "total_count")
			if err != nil {
				return err
			}
			if reported != total {
				return fmt.Errorf("%s reports total_count %d, want %d", pageId, reported, total)
			}
		}
		firstRows, err := autoProvisioningGroupsStateCount(state, pageOneId, "groups.#")
		if err != nil {
			return err
		}
		if firstRows != total {
			return fmt.Errorf("the read from page 1 returned %d rows with page_size = 1, want the complete result of %d rows", firstRows, total)
		}
		rows, err := autoProvisioningGroupsStateCount(state, pageTwoId, "groups.#")
		if err != nil {
			return err
		}
		if total == 0 {
			if rows != 0 {
				return fmt.Errorf("page 2 of an empty account must be empty, got %d rows", rows)
			}
			return nil
		}
		if rows != total-1 {
			return fmt.Errorf("the read from page 2 returned %d rows, want the %d rows that follow the first page", rows, total-1)
		}

		allIds := make(map[string]struct{}, total)
		for index := 0; index < total; index++ {
			id, err := autoProvisioningGroupsStateAttribute(state, allResourceId, fmt.Sprintf("ids.%d", index))
			if err != nil {
				return err
			}
			allIds[id] = struct{}{}
		}
		secondIds := make(map[string]struct{}, rows)
		for index := 0; index < rows; index++ {
			id, err := autoProvisioningGroupsStateAttribute(state, pageTwoId, fmt.Sprintf("ids.%d", index))
			if err != nil {
				return err
			}
			if _, ok := allIds[id]; !ok {
				return fmt.Errorf("page 2 returned %q, which is not in the complete result", id)
			}
			secondIds[id] = struct{}{}
		}
		missing := make([]string, 0, 1)
		for id := range allIds {
			if _, ok := secondIds[id]; !ok {
				missing = append(missing, id)
			}
		}
		if len(missing) != 1 {
			return fmt.Errorf("page 2 differs from the complete result by %d groups, want exactly the one row read on page 1: %v", len(missing), missing)
		}
		return nil
	}
}
