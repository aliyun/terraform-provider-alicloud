package ecs_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// preCheckNeedRetry retries the precheck through the same classes as the read:
// the shared helper covers throttling, decodable 5xx and transport failures, and
// an undecodable 5xx is added here because the helper does not know that type.
func preCheckNeedRetry(err error) bool {
	var unreadable *connectivity.UnreadableResponseError
	if errors.As(err, &unreadable) && unreadable.StatusCode >= 500 {
		return true
	}
	return helper.NeedRetry(err)
}

// preCheckAutoProvisioningGroupsExist skips the test when the account holds no
// auto provisioning group. Every step is derived from `groups[0]` and this
// workflow must not create a billable group, so an empty account would otherwise
// fail at plan time with `Invalid index`, which names neither the cause nor the
// missing precondition.
func preCheckAutoProvisioningGroupsExist(t *testing.T) {
	region := acctest.DefaultRegionToTest
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		t.Fatalf("cannot initialize the auto provisioning group precheck: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	// The sweeper client carries no transport timeout of its own, so the
	// precheck bounds the call and retries transient failures — including a 5xx
	// whose body cannot be decoded, which the shared helper does not classify —
	// instead of failing the single acceptance evidence on one such response.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	response, err := helper.RetryRpcPost(ctx, "DescribeAutoProvisioningGroups", func() (map[string]interface{}, error) {
		return client.RpcPostWithContext(ctx, "Ecs", "2014-05-26", "DescribeAutoProvisioningGroups", nil, map[string]interface{}{
			"RegionId": region,
		})
	}, 2*time.Minute, preCheckNeedRetry)
	if err != nil {
		t.Fatalf("cannot query the auto provisioning groups this test reads: %s", err)
	}
	if count, ok := autoProvisioningGroupCount(response); ok {
		if count == 0 {
			t.Skip("this data source's acceptance test reads the auto provisioning groups that already exist in the test account; it needs at least one and must not create a billable group")
		}
		return
	}
	t.Fatal("cannot read the auto provisioning group count from the precheck response")
}

// autoProvisioningGroupCount reads the record count from the RPC envelope. An
// empty result omits the record list, so the envelope total is the fallback. It
// reports false when neither shape is readable, so a decode change fails the
// precheck instead of silently skipping the test.
func autoProvisioningGroupCount(response map[string]interface{}) (int, bool) {
	if container, ok := response["AutoProvisioningGroups"].(map[string]interface{}); ok {
		if list, ok := container["AutoProvisioningGroup"].([]interface{}); ok {
			return len(list), true
		}
		return 0, true
	}
	if raw, ok := response["TotalCount"]; ok {
		if count, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(raw))); err == nil {
			return count, true
		}
	}
	return 0, false
}

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
	resourceGroupMismatchId := "data.alicloud_ecs_auto_provisioning_groups.resource_group_mismatch"

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
			preCheckAutoProvisioningGroupsExist(t)
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
				// The positive query uses an existing resource group. The negative
				// query must return no rows while the unfiltered query remains
				// non-empty, so dropping or ignoring this filter fails the step.
				Config: autoProvisioningGroupsResourceGroupConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkAutoProvisioningGroupsResourceGroupFilter(t, resourceId, allResourceId),
					checkAutoProvisioningGroupsResourceGroupMismatch(allResourceId, resourceGroupMismatchId),
				),
			},
			{
				// The tag filter is derived from the tags of an existing group. When
				// that group carries no expressible tag, the filter falls back to a
				// key that cannot match, which still proves the filter is sent and
				// excludes the groups it must not return.
				Config: autoProvisioningGroupsTagConfig(),
				Check: resource.ComposeTestCheckFunc(
					checkAutoProvisioningGroupsTagFilter(t, resourceId, allResourceId),
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
				// The status filter is not part of this release: an argument the
				// schema does not declare must fail instead of being ignored.
				Config:      autoProvisioningGroupsConfig("  status = [\"active\"]\n"),
				ExpectError: regexp.MustCompile(`Unsupported argument|not expected here`),
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
// result without depending on the account holding more than one type: the
// unfiltered read must return at least one group in the same step, every row of
// each type-filtered read carries the requested type, no group appears under two
// types, the reads together cover the unfiltered result, and their row counts
// add up to it. The non-empty baseline is asserted here because four empty reads
// would otherwise satisfy every partition assertion. A server that ignores
// `types` returns the full result for each of the three reads, which fails the
// overlap and the sum assertions.
func checkAutoProvisioningGroupsTypesPartition(allResourceId string, typeResources map[string]string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		allCount, err := autoProvisioningGroupsStateCount(state, allResourceId, "groups.#")
		if err != nil {
			return err
		}
		if allCount < 1 {
			return fmt.Errorf("%s returned no auto provisioning group, so the type partition has no baseline to compare against", allResourceId)
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

const autoProvisioningGroupsNonMatchingResourceGroupId = "rg-apg-filter-no-match"

// autoProvisioningGroupsResourceGroupConfig pairs a filter derived from an
// existing group with an independent filter checked against every existing row.
func autoProvisioningGroupsResourceGroupConfig() string {
	return fmt.Sprintf(`
data "alicloud_ecs_auto_provisioning_groups" "all" {
}

data "alicloud_ecs_auto_provisioning_groups" "default" {
  resource_group_id = data.alicloud_ecs_auto_provisioning_groups.all.groups[0].resource_group_id
}

data "alicloud_ecs_auto_provisioning_groups" "resource_group_mismatch" {
  resource_group_id = %q
}
`, autoProvisioningGroupsNonMatchingResourceGroupId)
}

// checkAutoProvisioningGroupsResourceGroupFilter checks the positive half of
// the resource group query: every row carries the filtered resource group and
// the group the filter was derived from is returned. The separate mismatch
// check below proves that the filter also excludes the unfiltered rows.
func checkAutoProvisioningGroupsResourceGroupFilter(t *testing.T, resourceId, allResourceId string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		allCount, err := autoProvisioningGroupsStateCount(state, allResourceId, "groups.#")
		if err != nil {
			return err
		}
		target, err := autoProvisioningGroupsStateAttribute(state, allResourceId, "groups.0.resource_group_id")
		if err != nil {
			return err
		}
		if target == "" {
			t.Logf("gap: the first group of %s reports no resource group, so the resource_group_id filter was not exercised", allResourceId)
			return nil
		}
		count, err := autoProvisioningGroupsStateCount(state, resourceId, "groups.#")
		if err != nil {
			return err
		}
		if count < 1 {
			return fmt.Errorf("the resource group filter %q matched no group", target)
		}
		if count > allCount {
			return fmt.Errorf("the resource group filter returned %d groups while only %d exist", count, allCount)
		}
		targetId, err := autoProvisioningGroupsStateAttribute(state, allResourceId, "ids.0")
		if err != nil {
			return err
		}
		found := false
		for index := 0; index < count; index++ {
			id, err := autoProvisioningGroupsStateAttribute(state, resourceId, fmt.Sprintf("ids.%d", index))
			if err != nil {
				return err
			}
			if id == targetId {
				found = true
			}
			returned, err := autoProvisioningGroupsStateAttribute(state, resourceId, fmt.Sprintf("groups.%d.resource_group_id", index))
			if err != nil {
				return err
			}
			if returned != target {
				return fmt.Errorf("group %q is in resource group %q, want the filtered resource group %q", id, returned, target)
			}
		}
		if !found {
			return fmt.Errorf("the resource group filter %q did not return the group %q it was derived from", target, targetId)
		}
		return nil
	}
}

func checkAutoProvisioningGroupsResourceGroupMismatch(allResourceId, mismatchResourceId string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		allCount, err := autoProvisioningGroupsStateCount(state, allResourceId, "groups.#")
		if err != nil {
			return err
		}
		if allCount < 1 {
			return fmt.Errorf("the resource group mismatch check requires a non-empty unfiltered result")
		}
		for index := 0; index < allCount; index++ {
			resourceGroupId, err := autoProvisioningGroupsStateAttribute(state, allResourceId, fmt.Sprintf("groups.%d.resource_group_id", index))
			if err != nil {
				return err
			}
			if resourceGroupId == autoProvisioningGroupsNonMatchingResourceGroupId {
				return fmt.Errorf("the non-matching resource group probe %q unexpectedly belongs to unfiltered group index %d", resourceGroupId, index)
			}
		}
		for _, attribute := range []string{"groups.#", "ids.#", "names.#", "total_count"} {
			if err := resource.TestCheckResourceAttr(mismatchResourceId, attribute, "0")(state); err != nil {
				return fmt.Errorf("resource group mismatch filter was not applied: %w", err)
			}
		}
		return nil
	}
}

// autoProvisioningGroupsTagConfig filters by one tag of the first group of the
// unfiltered read. A tag key that the filter cannot express (a dot, which the
// state cannot address, or a reserved prefix, which the validator rejects) is
// skipped. When no usable tag remains, the filter is set to a key that cannot
// match, so the step still proves the tag filter is sent and excludes the
// groups it must not return.
func autoProvisioningGroupsTagConfig() string {
	return `
data "alicloud_ecs_auto_provisioning_groups" "all" {
}

locals {
  tag_keys = [
    for key in keys(data.alicloud_ecs_auto_provisioning_groups.all.groups[0].tags) : key
    if !can(regex("\\.", key)) && !can(regex("^(aliyun|acs:)", key))
  ]
  tag_key = length(local.tag_keys) > 0 ? local.tag_keys[0] : ""
}

data "alicloud_ecs_auto_provisioning_groups" "default" {
  tags = local.tag_key != "" ? {
    (local.tag_key) = data.alicloud_ecs_auto_provisioning_groups.all.groups[0].tags[local.tag_key]
    } : {
    "tf-acc-tag-probe" = "no-such-tag-value"
  }
}
`
}

// autoProvisioningGroupsStateTagKeys lists the tag keys of one row of the data
// source state. The map is stored flat, so the keys are the attribute names that
// follow `groups.<index>.tags.`.
func autoProvisioningGroupsStateTagKeys(state *terraform.State, resourceId string, index int) []string {
	instance, ok := state.RootModule().Resources[resourceId]
	if !ok || instance.Primary == nil {
		return nil
	}
	prefix := fmt.Sprintf("groups.%d.tags.", index)
	keys := make([]string, 0, 4)
	for name := range instance.Primary.Attributes {
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		key := strings.TrimPrefix(name, prefix)
		if key == "%" || key == "#" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// checkAutoProvisioningGroupsTagFilter proves the tag filter is sent and keeps
// only the groups that carry it: every row carries the filtered tag value and
// the group the filter was derived from is returned. Without a usable tag on the
// first group, the non-matching filter must return nothing; that case is
// reported as a gap instead of a pass, because the narrowing of a real tag was
// not observed.
func checkAutoProvisioningGroupsTagFilter(t *testing.T, resourceId, allResourceId string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		allCount, err := autoProvisioningGroupsStateCount(state, allResourceId, "groups.#")
		if err != nil {
			return err
		}
		if allCount < 1 {
			return fmt.Errorf("%s returned no auto provisioning group, so the tag filter has no baseline to compare against", allResourceId)
		}
		usable := ""
		for _, key := range autoProvisioningGroupsStateTagKeys(state, allResourceId, 0) {
			if strings.Contains(key, ".") || strings.HasPrefix(key, "aliyun") || strings.HasPrefix(key, "acs:") {
				continue
			}
			usable = key
			break
		}
		count, err := autoProvisioningGroupsStateCount(state, resourceId, "groups.#")
		if err != nil {
			return err
		}
		if usable == "" {
			t.Logf("gap: the first group of %s carries no tag the filter can express, so only the non-matching tag filter was exercised", allResourceId)
			if count != 0 {
				return fmt.Errorf("the tag filter that cannot match returned %d groups", count)
			}
			return nil
		}
		wanted, err := autoProvisioningGroupsStateAttribute(state, allResourceId, fmt.Sprintf("groups.0.tags.%s", usable))
		if err != nil {
			return err
		}
		if count < 1 {
			return fmt.Errorf("the tag filter %s=%q matched no group", usable, wanted)
		}
		if count > allCount {
			return fmt.Errorf("the tag filter returned %d groups while only %d exist", count, allCount)
		}
		targetId, err := autoProvisioningGroupsStateAttribute(state, allResourceId, "ids.0")
		if err != nil {
			return err
		}
		found := false
		for index := 0; index < count; index++ {
			id, err := autoProvisioningGroupsStateAttribute(state, resourceId, fmt.Sprintf("ids.%d", index))
			if err != nil {
				return err
			}
			if id == targetId {
				found = true
			}
			returned, err := autoProvisioningGroupsStateAttribute(state, resourceId, fmt.Sprintf("groups.%d.tags.%s", index, usable))
			if err != nil {
				return fmt.Errorf("group %q does not carry the filtered tag %q: %w", id, usable, err)
			}
			if returned != wanted {
				return fmt.Errorf("group %q has tag %s=%q, want the filtered value %q", id, usable, returned, wanted)
			}
		}
		if !found {
			return fmt.Errorf("the tag filter %s=%q did not return the group %q it was derived from", usable, wanted, targetId)
		}
		if count == allCount {
			t.Logf("gap: every group of the account carries %s=%q, so the filter's narrowing effect is not observable", usable, wanted)
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
