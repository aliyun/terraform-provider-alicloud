package sas_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// The following input is not exercised with a valid value by this test:
// resource_directory_account_id needs a member account of the resource directory
// that the acceptance account does not provide, so the query is only verified to
// transmit it and to surface the permission error of the service.
func TestAccAliCloudSecurityCenterScaDetails_basic0(t *testing.T) {
	dataSourceId := "data.alicloud_security_center_sca_details.default"
	seedId := "data.alicloud_security_center_sca_details.seed"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterScaDetailsConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceId, "id"),
					testAccCheckSecurityCenterScaDetailsPage(dataSourceId, 10),
					testAccCheckSecurityCenterScaDetailsFields(dataSourceId),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigPaging(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_security_center_sca_details.first", "propertys.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_security_center_sca_details.first", "page_item_count", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "propertys.#", "1"),
					resource.TestCheckResourceAttr(dataSourceId, "page_item_count", "1"),
					testAccCheckSecurityCenterScaDetailsPagesDiffer("data.alicloud_security_center_sca_details.first", dataSourceId),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigSearchMiss(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "propertys.#", "0"),
					resource.TestCheckResourceAttr(dataSourceId, "page_item_count", "0"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigSearchMatch(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsName(dataSourceId, "data.alicloud_security_center_sca_details.seed"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigCriteriaMatch(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					resource.TestCheckResourceAttr(dataSourceId, "search_criteria_list.#", "1"),
					testAccCheckSecurityCenterScaDetailsName(dataSourceId, "data.alicloud_security_center_sca_details.seed"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigScaNameMatch(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsName(dataSourceId, "data.alicloud_security_center_sca_details.seed"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigSubConditionMiss(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "propertys.#", "0"),
					resource.TestCheckResourceAttr(dataSourceId, "page_item_count", "0"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigUuidFilter(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsUuid(dataSourceId, "data.alicloud_security_center_sca_details.seed"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigUseNextToken(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					resource.TestCheckResourceAttrSet(dataSourceId, "next_token_result"),
					testAccCheckSecurityCenterScaDetailsPageSize(dataSourceId),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigNextTokenPage(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsPagesDiffer("data.alicloud_security_center_sca_details.first", dataSourceId),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigLang(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsDiffers(dataSourceId, "biz_type_dispaly", seedId, "biz_type_dispaly"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigBiz(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsAttrMatches(dataSourceId, "biz_type", "database"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigBizType(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsEqualsSeed(dataSourceId, "biz_type", seedId, "biz_type"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigRemark(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsEqualsSeed(dataSourceId, "ip", seedId, "ip"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigScaVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsEqualsSeed(dataSourceId, "version", seedId, "version"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigPort(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "propertys.#", "0"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigUser(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsEqualsSeed(dataSourceId, "process_user", seedId, "process_user"),
					testAccCheckSecurityCenterScaDetailsFewerThanSeed(dataSourceId, seedId),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigNamePattern(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsEqualsSeed(dataSourceId, "name", seedId, "name"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigPid(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsEqualsSeed(dataSourceId, "pid", "data.alicloud_security_center_sca_details.seed_db", "pid"),
				),
			},
			{
				Config: testAccSecurityCenterScaDetailsConfigProcessStarted(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityCenterScaDetailsNonEmpty(dataSourceId),
					testAccCheckSecurityCenterScaDetailsProcessStarted(dataSourceId, seedId, 60),
				),
			},
			{
				Config:      testAccSecurityCenterScaDetailsConfigSearchItemOnly(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      testAccSecurityCenterScaDetailsConfigSearchInfoOnly(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      testAccSecurityCenterScaDetailsConfigSearchItemSubOnly(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      testAccSecurityCenterScaDetailsConfigSearchInfoSubOnly(),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				Config:      testAccSecurityCenterScaDetailsConfigInvalidBiz(),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			{
				Config:      testAccSecurityCenterScaDetailsConfigResourceDirectory(),
				ExpectError: regexp.MustCompile("RdCheckNoPermission"),
			},
		},
	})
}

// testAccCheckSecurityCenterScaDetailsPage verifies that the page holds at most
// maxPageSize asset fingerprints, that the reported page count matches the returned
// list, and that the total count covers the current page.
func testAccCheckSecurityCenterScaDetailsPage(name string, maxPageSize int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		count, err := securityCenterScaDetailsCount(attrs, "propertys.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if count < 1 || count > maxPageSize {
			return fmt.Errorf("%s: propertys.# is %d, expected 1 to %d", name, count, maxPageSize)
		}
		pageCount, err := securityCenterScaDetailsCount(attrs, "page_item_count")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if count != pageCount {
			return fmt.Errorf("%s: page_item_count %d does not match propertys.# %d", name, pageCount, count)
		}
		total, err := securityCenterScaDetailsCount(attrs, "total_count")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if total < pageCount {
			return fmt.Errorf("%s: total_count %d is smaller than page_item_count %d", name, total, pageCount)
		}
		return nil
	}
}

// testAccCheckSecurityCenterScaDetailsPageSize verifies that the token query
// reports the entries of the page it returned.
func testAccCheckSecurityCenterScaDetailsPageSize(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		count, err := securityCenterScaDetailsCount(attrs, "propertys.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		pageCount, err := securityCenterScaDetailsCount(attrs, "page_item_count")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if count != pageCount {
			return fmt.Errorf("%s: page_item_count %d does not match propertys.# %d", name, pageCount, count)
		}
		return nil
	}
}

// testAccCheckSecurityCenterScaDetailsNonEmpty requires the query to return at
// least one asset fingerprint.
func testAccCheckSecurityCenterScaDetailsNonEmpty(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		count, err := securityCenterScaDetailsCount(attrs, "propertys.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if count < 1 {
			return fmt.Errorf("%s: propertys is empty", name)
		}
		return nil
	}
}

// securityCenterScaDetailsPropertyFields lists every field the API contract
// documents for one asset fingerprint.
var securityCenterScaDetailsPropertyFields = []string{
	"type", "biz_type_dispaly", "process_started", "internet_ip", "config_path",
	"pid", "port", "cmdline", "biz_type", "listen_ip", "version", "instance_name",
	"listen_status", "name", "ip", "process_user", "instance_id", "web_path",
	"ppid", "intranet_ip", "uuid", "listen_protocol", "image_name", "path",
	"container_name", "proof", "create_timestamp", "runtime_env_version", "pod_name",
}

// testAccCheckSecurityCenterScaDetailsFields verifies that the documented fields
// of the returned asset fingerprints are mapped into state, including the fields
// the API returns empty and the misspelled biz_type_dispaly key.
func testAccCheckSecurityCenterScaDetailsFields(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		count, err := securityCenterScaDetailsCount(attrs, "propertys.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if count < 1 {
			return fmt.Errorf("%s: propertys is empty; the acceptance account must have collected asset fingerprints", name)
		}
		for _, field := range securityCenterScaDetailsPropertyFields {
			if _, ok := attrs["propertys.0."+field]; !ok {
				return fmt.Errorf("%s: propertys.0.%s is missing from state", name, field)
			}
		}
		for i := 1; i < count; i++ {
			prefix := fmt.Sprintf("propertys.%d.", i)
			for _, field := range securityCenterScaDetailsCoreFields {
				if _, ok := attrs[prefix+field]; !ok {
					return fmt.Errorf("%s: %s%s is missing from state", name, prefix, field)
				}
			}
		}
		if attrs["propertys.0.biz_type_dispaly"] == "" {
			return fmt.Errorf("%s: propertys.0.biz_type_dispaly is empty", name)
		}
		return nil
	}
}

// securityCenterScaDetailsCoreFields lists the fields every asset fingerprint is
// expected to carry.
var securityCenterScaDetailsCoreFields = []string{"type", "biz_type", "name", "uuid", "instance_id", "pid"}

// testAccCheckSecurityCenterScaDetailsPagesDiffer verifies that two queries that
// differ only in the page read different asset fingerprints, so the pagination is
// applied by the service instead of returning the first page twice.
func testAccCheckSecurityCenterScaDetailsPagesDiffer(firstName, secondName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		first, err := securityCenterScaDetailsIdentity(s, firstName)
		if err != nil {
			return err
		}
		second, err := securityCenterScaDetailsIdentity(s, secondName)
		if err != nil {
			return err
		}
		if first == second {
			return fmt.Errorf("%s and %s returned the same asset fingerprint %q; the page did not change", firstName, secondName, first)
		}
		return nil
	}
}

// testAccCheckSecurityCenterScaDetailsUuid verifies that every returned asset
// fingerprint belongs to the server the uuid filter names.
func testAccCheckSecurityCenterScaDetailsUuid(name, seedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		seedAttrs, err := securityCenterScaDetailsAttributes(s, seedName)
		if err != nil {
			return err
		}
		uuid := seedAttrs["propertys.0.uuid"]
		if uuid == "" {
			return fmt.Errorf("%s: propertys.0.uuid is empty", seedName)
		}
		return testAccCheckSecurityCenterScaDetailsAttrMatches(name, "uuid", uuid)(s)
	}
}

// testAccCheckSecurityCenterScaDetailsName verifies that the search condition on
// the fingerprint name only returns fingerprints of that name.
func testAccCheckSecurityCenterScaDetailsName(name, seedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		seedAttrs, err := securityCenterScaDetailsAttributes(s, seedName)
		if err != nil {
			return err
		}
		expected := seedAttrs["propertys.0.name"]
		if expected == "" {
			return fmt.Errorf("%s: propertys.0.name is empty", seedName)
		}
		return testAccCheckSecurityCenterScaDetailsAttrMatches(name, "name", expected)(s)
	}
}

func testAccCheckSecurityCenterScaDetailsAttrMatches(name, attribute, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		count, err := securityCenterScaDetailsCount(attrs, "propertys.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		for i := 0; i < count; i++ {
			if got := attrs[fmt.Sprintf("propertys.%d.%s", i, attribute)]; got != expected {
				return fmt.Errorf("%s: propertys.%d.%s is %q, expected the filtered value %q", name, i, attribute, got, expected)
			}
		}
		return nil
	}
}

// testAccCheckSecurityCenterScaDetailsEqualsSeed verifies that every returned asset
// fingerprint carries the value the unfiltered seed query reported, so the filter
// narrowed the result set to that value instead of being ignored by the service.
func testAccCheckSecurityCenterScaDetailsEqualsSeed(name, attribute, seedName, seedAttribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		seedAttrs, err := securityCenterScaDetailsAttributes(s, seedName)
		if err != nil {
			return err
		}
		expected := seedAttrs["propertys.0."+seedAttribute]
		if expected == "" {
			return fmt.Errorf("%s: propertys.0.%s is empty", seedName, seedAttribute)
		}
		return testAccCheckSecurityCenterScaDetailsAttrMatches(name, attribute, expected)(s)
	}
}

// testAccCheckSecurityCenterScaDetailsDiffers verifies that the returned value
// changed with the request, which a different language of the response does.
func testAccCheckSecurityCenterScaDetailsDiffers(name, attribute, seedName, seedAttribute string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		seedAttrs, err := securityCenterScaDetailsAttributes(s, seedName)
		if err != nil {
			return err
		}
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		got := attrs["propertys.0."+attribute]
		if got == "" {
			return fmt.Errorf("%s: propertys.0.%s is empty", name, attribute)
		}
		if got == seedAttrs["propertys.0."+seedAttribute] {
			return fmt.Errorf("%s: propertys.0.%s is still %q; the request did not change the response", name, attribute, got)
		}
		return nil
	}
}

// testAccCheckSecurityCenterScaDetailsFewerThanSeed verifies that the filter reduced
// the reported total, so it selected a subset instead of being ignored.
func testAccCheckSecurityCenterScaDetailsFewerThanSeed(name, seedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		total, err := securityCenterScaDetailsCount(attrs, "total_count")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		seedAttrs, err := securityCenterScaDetailsAttributes(s, seedName)
		if err != nil {
			return err
		}
		seedTotal, err := securityCenterScaDetailsCount(seedAttrs, "total_count")
		if err != nil {
			return fmt.Errorf("%s: %v", seedName, err)
		}
		if total >= seedTotal {
			return fmt.Errorf("%s: total_count %d is not smaller than the unfiltered total %d", name, total, seedTotal)
		}
		return nil
	}
}

// testAccCheckSecurityCenterScaDetailsProcessStarted verifies that every returned
// asset fingerprint started within window seconds of the seed fingerprint. The
// service compares these bounds as milliseconds while the published reference
// states seconds, so the configuration converts the anchor to milliseconds.
func testAccCheckSecurityCenterScaDetailsProcessStarted(name, seedName string, window int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		seedAttrs, err := securityCenterScaDetailsAttributes(s, seedName)
		if err != nil {
			return err
		}
		anchor, err := strconv.ParseInt(seedAttrs["propertys.0.process_started"], 10, 64)
		if err != nil {
			return fmt.Errorf("%s: process_started is not a number: %v", seedName, err)
		}
		lower, upper := anchor-window, anchor+window
		attrs, err := securityCenterScaDetailsAttributes(s, name)
		if err != nil {
			return err
		}
		count, err := securityCenterScaDetailsCount(attrs, "propertys.#")
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		for i := 0; i < count; i++ {
			value := attrs[fmt.Sprintf("propertys.%d.process_started", i)]
			started, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("%s: propertys.%d.process_started is not a number: %v", name, i, err)
			}
			if started < lower || started > upper {
				return fmt.Errorf("%s: propertys.%d.process_started %d is outside %d..%d", name, i, started, lower, upper)
			}
		}
		return nil
	}
}

func securityCenterScaDetailsIdentity(s *terraform.State, name string) (string, error) {
	attrs, err := securityCenterScaDetailsAttributes(s, name)
	if err != nil {
		return "", err
	}
	count, err := securityCenterScaDetailsCount(attrs, "propertys.#")
	if err != nil {
		return "", fmt.Errorf("%s: %v", name, err)
	}
	if count < 1 {
		return "", fmt.Errorf("%s: propertys is empty", name)
	}
	identity := attrs["propertys.0.name"] + "|" + attrs["propertys.0.pid"] + "|" + attrs["propertys.0.process_user"]
	if identity == "||" {
		return "", fmt.Errorf("%s: propertys.0 has no name, pid or process_user", name)
	}
	return identity, nil
}

func securityCenterScaDetailsAttributes(s *terraform.State, name string) (map[string]string, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("data source %s is not found in state", name)
	}
	return rs.Primary.Attributes, nil
}

func securityCenterScaDetailsCount(attrs map[string]string, key string) (int, error) {
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

func testAccSecurityCenterScaDetailsConfigBasic() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigPaging() string {
	return `
data "alicloud_security_center_sca_details" "first" {
  page_size    = 1
  current_page = 1
}

data "alicloud_security_center_sca_details" "default" {
  page_size    = 1
  current_page = 2
}
`
}

func testAccSecurityCenterScaDetailsConfigSearchMiss() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  search_item = "name"
  search_info = "tf-testacc-not-exist"
  page_size   = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigSearchMatch() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  search_item = "name"
  search_info = data.alicloud_security_center_sca_details.seed.propertys[0].name
  page_size   = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigCriteriaMatch() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  search_criteria_list = [
    {
      name  = "name"
      value = data.alicloud_security_center_sca_details.seed.propertys[0].name
    }
  ]
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigScaNameMatch() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  sca_name  = data.alicloud_security_center_sca_details.seed.propertys[0].name
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigSubConditionMiss() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  search_item_sub = "pid"
  search_info_sub = "999999"
  page_size       = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigUuidFilter() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  uuid      = data.alicloud_security_center_sca_details.seed.propertys[0].uuid
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigUseNextToken() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  use_next_token = true
}
`
}

func testAccSecurityCenterScaDetailsConfigNextTokenPage() string {
	return `
data "alicloud_security_center_sca_details" "first" {
  use_next_token = true
}

data "alicloud_security_center_sca_details" "default" {
  use_next_token = true
  next_token     = data.alicloud_security_center_sca_details.first.next_token_result
}
`
}

func testAccSecurityCenterScaDetailsConfigLang() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  lang      = "en"
  page_size = 1
}
`
}

func testAccSecurityCenterScaDetailsConfigBiz() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  biz       = "sca_database"
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigBizType() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  biz_type  = data.alicloud_security_center_sca_details.seed.propertys[0].biz_type
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigRemark() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  remark    = data.alicloud_security_center_sca_details.seed.propertys[0].ip
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigScaVersion() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  sca_version = data.alicloud_security_center_sca_details.seed.propertys[0].version
  page_size   = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigPort() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  biz       = "sca_database"
  port      = "65535"
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigUser() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 10
}

data "alicloud_security_center_sca_details" "default" {
  user      = data.alicloud_security_center_sca_details.seed.propertys[0].process_user
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigNamePattern() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  sca_name_pattern = data.alicloud_security_center_sca_details.seed.propertys[0].name
  page_size        = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigPid() string {
	return `
data "alicloud_security_center_sca_details" "seed_db" {
  biz       = "sca_database"
  page_size = 1
}

data "alicloud_security_center_sca_details" "default" {
  pid       = data.alicloud_security_center_sca_details.seed_db.propertys[0].pid
  page_size = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigProcessStarted() string {
	return `
data "alicloud_security_center_sca_details" "seed" {
  page_size = 10
}

data "alicloud_security_center_sca_details" "default" {
  process_started_start = (data.alicloud_security_center_sca_details.seed.propertys[0].process_started - 60) * 1000
  process_started_end   = (data.alicloud_security_center_sca_details.seed.propertys[0].process_started + 60) * 1000
  page_size             = 10
}
`
}

func testAccSecurityCenterScaDetailsConfigSearchItemOnly() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  search_item = "name"
}
`
}

func testAccSecurityCenterScaDetailsConfigSearchInfoOnly() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  search_info = "liblzma"
}
`
}

func testAccSecurityCenterScaDetailsConfigSearchItemSubOnly() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  search_item_sub = "pid"
}
`
}

func testAccSecurityCenterScaDetailsConfigSearchInfoSubOnly() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  search_info_sub = "1"
}
`
}

func testAccSecurityCenterScaDetailsConfigInvalidBiz() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  biz = "sca_invalid"
}
`
}

func testAccSecurityCenterScaDetailsConfigResourceDirectory() string {
	return `
data "alicloud_security_center_sca_details" "default" {
  resource_directory_account_id = 1234567890123456
  page_size                     = 1
}
`
}
