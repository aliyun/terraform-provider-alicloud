package sas_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccAliCloudSecurityCenterGroupedInstances_basic0 covers the query
// dimensions, the paging contract and the option pass-through of the data
// source. The account must have Security Center activated with ECS assets in
// more than one region; when it has fewer, the cross-page checks degrade to
// verifying the response shape instead of failing on account state.
//
// The paging assertions follow the behaviour the service really exhibits:
// no_page = true returns every group in one response and ignores page_size and
// current_page, while an unset or false no_page keeps the paged query. The
// public API description states the opposite; the cases below record the
// observed behaviour so that a future service or documentation change fails
// here instead of silently drifting.
//
// resource_directory_account_id is deliberately not exercised here: the
// authorized test account has no resource directory member account, which is
// the documented admission limitation of this change. The parameter keeps the
// contract shape (int64, request key ResourceDirectoryAccountId) and is passed
// through unchanged, but its member-account behaviour stays unverified until
// such an account is available.
func TestAccAliCloudSecurityCenterGroupedInstances_basic0(t *testing.T) {
	dataSourceId := "data.alicloud_security_center_grouped_instances.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	var firstPageAlias string
	var pagedTotalCount int
	var regionQueryId string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterGroupedInstancesConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "group_field", "regionId"),
					resource.TestCheckResourceAttr(dataSourceId, "machine_types", "ecs"),
					resource.TestCheckResourceAttr(dataSourceId, "lang", "zh"),
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
					resource.TestCheckResourceAttrSet(dataSourceId, "total_count"),
					resource.TestCheckResourceAttrSet(dataSourceId, "current_page"),
					testAccCheckSecurityCenterGroupedInstancesPageCount(dataSourceId),
					testAccCaptureSecurityCenterGroupedInstancesId(dataSourceId, &regionQueryId),
					testAccCheckSecurityCenterGroupedInstancesElements(dataSourceId, []string{
						"field_alias_name",
						"instance_count",
						"risk_instance_count",
						"unprotected_instance_count",
						"auth_version_check_count",
						"instance_core_count",
					}),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigPaging(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "page_size", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "current_page", "1"),
					testAccCheckSecurityCenterGroupedInstancesPageCount(dataSourceId),
					testAccCheckSecurityCenterGroupedInstancesPageSize(dataSourceId),
					testAccCaptureSecurityCenterGroupedInstancesAlias(dataSourceId, &firstPageAlias),
					testAccCaptureSecurityCenterGroupedInstancesTotal(dataSourceId, &pagedTotalCount),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigPaging(2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "current_page", "2"),
					testAccCheckSecurityCenterGroupedInstancesPageCount(dataSourceId),
					testAccCheckSecurityCenterGroupedInstancesNextPage(dataSourceId, &firstPageAlias),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigNoPage(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "no_page", "true"),
					resource.TestCheckResourceAttr(dataSourceId, "page_size", "1"),
					testAccCheckSecurityCenterGroupedInstancesAllOnOnePage(dataSourceId, &pagedTotalCount),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigNoPage(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "no_page", "false"),
					testAccCheckSecurityCenterGroupedInstancesPageSize(dataSourceId),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigFilters(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "vendors", "0,8"),
					resource.TestCheckResourceAttr(dataSourceId, "sale_version_check_code", "sas_gt_basic"),
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
					testAccCheckSecurityCenterGroupedInstancesPageCount(dataSourceId),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigVendor(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "vendor", "0"),
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
					testAccCheckSecurityCenterGroupedInstancesPageCount(dataSourceId),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigGroupId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "group_field", "groupId"),
					testAccCheckSecurityCenterGroupedInstancesElements(dataSourceId, []string{
						"field_alias_name",
						"instance_count",
						"group_flag",
					}),
				),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigVpcInstanceId(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "group_field", "vpcInstanceId"),
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
					testAccCheckSecurityCenterGroupedInstancesPageCount(dataSourceId),
					testAccCheckSecurityCenterGroupedInstancesIdDiffers(dataSourceId, &regionQueryId),
					// Assets outside a VPC or in the default group come back
					// without a group name, so only the mandatory counter is
					// required here.
					testAccCheckSecurityCenterGroupedInstancesElements(dataSourceId, []string{
						"instance_count",
					}),
				),
			},
			{
				Config:      testAccSecurityCenterGroupedInstancesConfigInvalidGroupField(),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			{
				Config: testAccSecurityCenterGroupedInstancesConfigEmptyResult(rand),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "instances.#", "0"),
					resource.TestCheckResourceAttr(dataSourceId, "page_item_count", "0"),
				),
			},
		},
	})
}

// testAccCheckSecurityCenterGroupedInstancesPageCount verifies that the
// reported page count matches the returned list, including the empty result.
func testAccCheckSecurityCenterGroupedInstancesPageCount(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}

		instanceCount, ok := attrs["instances.#"]
		if !ok {
			return fmt.Errorf("%s: instances is missing from state", name)
		}
		pageCount, ok := attrs["page_item_count"]
		if !ok {
			return fmt.Errorf("%s: count is missing from state", name)
		}
		if pageCount != instanceCount {
			return fmt.Errorf("%s: count %q does not match instances.# %q", name, pageCount, instanceCount)
		}
		return nil
	}
}

// testAccCheckSecurityCenterGroupedInstancesElements verifies that the
// attributes the API returns for every group are mapped into state. Attributes
// the API omits for a group (for example group_flag outside the group
// dimension, or os without operating system information) are not required.
func testAccCheckSecurityCenterGroupedInstancesElements(name string, requiredAttributes []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}

		instanceCount, err := securityCenterGroupedInstancesCount(attrs, "instances.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		for i := 0; i < instanceCount; i++ {
			prefix := fmt.Sprintf("instances.%d.", i)
			for _, attribute := range requiredAttributes {
				if _, ok := attrs[prefix+attribute]; !ok {
					return fmt.Errorf("%s: %s%s is missing from state", name, prefix, attribute)
				}
			}
		}
		return nil
	}
}

// testAccCheckSecurityCenterGroupedInstancesPageSize verifies that page_size
// is applied instead of the server default: a result set larger than one page
// must return exactly one entry.
func testAccCheckSecurityCenterGroupedInstancesPageSize(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}

		totalCount, err := securityCenterGroupedInstancesCount(attrs, "total_count")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		instanceCount, err := securityCenterGroupedInstancesCount(attrs, "instances.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if totalCount > 1 && instanceCount != 1 {
			return fmt.Errorf("%s: page_size is 1 but %d entries were returned for %d assets", name, instanceCount, totalCount)
		}
		return nil
	}
}

// testAccCaptureSecurityCenterGroupedInstancesAlias remembers the first group
// of page one so that the next page can be proven to be a different slice.
func testAccCaptureSecurityCenterGroupedInstancesAlias(name string, alias *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}
		*alias = attrs["instances.0.field_alias_name"]
		return nil
	}
}

// testAccCheckSecurityCenterGroupedInstancesNextPage verifies that
// current_page is applied by the API: with more than one group, page two must
// return the second group instead of repeating page one.
func testAccCheckSecurityCenterGroupedInstancesNextPage(name string, firstPageAlias *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}

		totalCount, err := securityCenterGroupedInstancesCount(attrs, "total_count")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if totalCount <= 1 || *firstPageAlias == "" {
			return nil
		}

		instanceCount, err := securityCenterGroupedInstancesCount(attrs, "instances.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if instanceCount != 1 {
			return fmt.Errorf("%s: page_size is 1 but page 2 returned %d entries", name, instanceCount)
		}
		alias := attrs["instances.0.field_alias_name"]
		if alias == *firstPageAlias {
			return fmt.Errorf("%s: page 2 repeated the group %q of page 1, so current_page was not applied", name, alias)
		}
		return nil
	}
}

// testAccCheckSecurityCenterGroupedInstancesAllOnOnePage verifies that
// no_page = true disables the paging: even with page_size = 1 the response
// carries every group the paged query counted before.
func testAccCheckSecurityCenterGroupedInstancesAllOnOnePage(name string, pagedTotal *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}

		instanceCount, err := securityCenterGroupedInstancesCount(attrs, "instances.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if *pagedTotal <= 1 {
			return nil
		}
		if instanceCount != *pagedTotal {
			return fmt.Errorf("%s: no_page is true but only %d of %d groups were returned", name, instanceCount, *pagedTotal)
		}
		return nil
	}
}

// testAccCaptureSecurityCenterGroupedInstancesId remembers the ID derived from
// one query condition so that another condition can be proven to differ.
func testAccCaptureSecurityCenterGroupedInstancesId(name string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}
		*id = attrs["id"]
		return nil
	}
}

// testAccCheckSecurityCenterGroupedInstancesIdDiffers verifies that the ID is
// derived from the query condition instead of being fixed.
func testAccCheckSecurityCenterGroupedInstancesIdDiffers(name string, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}
		if attrs["id"] == "" {
			return fmt.Errorf("%s: id is missing from state", name)
		}
		if *previous == "" {
			return fmt.Errorf("%s: the reference ID was not captured", name)
		}
		if attrs["id"] == *previous {
			return fmt.Errorf("%s: id %q is identical for different query conditions", name, attrs["id"])
		}
		return nil
	}
}

// testAccCaptureSecurityCenterGroupedInstancesTotal remembers the group count
// the paged query reported so the unpaged query can be compared with it.
func testAccCaptureSecurityCenterGroupedInstancesTotal(name string, total *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterGroupedInstancesAttributes(s, name)
		if err != nil {
			return err
		}
		*total, err = securityCenterGroupedInstancesCount(attrs, "total_count")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		return nil
	}
}

func securityCenterGroupedInstancesAttributes(s *terraform.State, name string) (map[string]string, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("data source %s is not found in state", name)
	}
	return rs.Primary.Attributes, nil
}

func securityCenterGroupedInstancesCount(attrs map[string]string, key string) (int, error) {
	value, ok := attrs[key]
	if !ok {
		return 0, fmt.Errorf("%s is missing from state", key)
	}
	count, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s is not a number: %v", key, err)
	}
	return count, nil
}

func testAccSecurityCenterGroupedInstancesConfigBasic() string {
	return `
data "alicloud_security_center_grouped_instances" "default" {
  group_field   = "regionId"
  machine_types = "ecs"
  lang          = "zh"
  page_size     = 10
  current_page  = 1
}
`
}

func testAccSecurityCenterGroupedInstancesConfigPaging(currentPage int) string {
	return fmt.Sprintf(`
data "alicloud_security_center_grouped_instances" "default" {
  group_field   = "regionId"
  machine_types = "ecs"
  page_size     = 1
  current_page  = %d
}
`, currentPage)
}

func testAccSecurityCenterGroupedInstancesConfigNoPage(noPage bool) string {
	return fmt.Sprintf(`
data "alicloud_security_center_grouped_instances" "default" {
  group_field   = "regionId"
  machine_types = "ecs"
  no_page       = %t
  page_size     = 1
}
`, noPage)
}

func testAccSecurityCenterGroupedInstancesConfigFilters() string {
	return `
data "alicloud_security_center_grouped_instances" "default" {
  group_field             = "regionId"
  machine_types           = "ecs"
  vendors                 = "0,8"
  sale_version_check_code = "sas_gt_basic"
}
`
}

func testAccSecurityCenterGroupedInstancesConfigVendor() string {
	return `
data "alicloud_security_center_grouped_instances" "default" {
  group_field   = "regionId"
  machine_types = "ecs"
  vendor        = 0
}
`
}

func testAccSecurityCenterGroupedInstancesConfigGroupId() string {
	return `
data "alicloud_security_center_grouped_instances" "default" {
  group_field   = "groupId"
  machine_types = "ecs"
}
`
}

func testAccSecurityCenterGroupedInstancesConfigVpcInstanceId() string {
	return `
data "alicloud_security_center_grouped_instances" "default" {
  group_field   = "vpcInstanceId"
  machine_types = "ecs"
}
`
}

func testAccSecurityCenterGroupedInstancesConfigEmptyResult(rand int) string {
	return fmt.Sprintf(`
data "alicloud_security_center_grouped_instances" "default" {
  group_field = "groupId"
  field_value = "tf-testacc-not-exist-%d"
}
`, rand)
}

func testAccSecurityCenterGroupedInstancesConfigInvalidGroupField() string {
	return `
data "alicloud_security_center_grouped_instances" "default" {
  group_field = "instanceId"
}
`
}
