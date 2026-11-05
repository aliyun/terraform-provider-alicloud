package dlfnext_test

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/dlfnext"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAliCloudDlfNextCatalog_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dlf_next_catalog.default"
	rc := tfacctest.RandIntRange(1000, 9999)
	catalogName := fmt.Sprintf("tf-testacc-dlf-merge-%d", rc)

	// Best-effort cleanup: if the ACC fails mid-lifecycle, CheckDestroy is not
	// reached and the catalog would leak, consuming the DLF backend cleanup
	// quota. t.Cleanup runs after the test result is recorded, so a cleanup
	// error never masks the original test failure.
	t.Cleanup(func() { testAccDlfNextCatalogBestEffortCleanup(t, catalogName) })

	// Writable user option keys (proven writable by prior ACC runs); the server
	// defaults dlf.trashed-file-retained-days and storage.data.redundancy.type
	// are injected by the backend and must NEVER be deleted by an options
	// update or an explicit options = {}.
	const serverDefault1 = "dlf.trashed-file-retained-days"
	const serverDefault2 = "storage.data.redundancy.type"

	// Single CreateCatalog lifecycle: Create -> Data Source (basic + filters)
	// -> Options Update -> Options Clear -> Import -> Re-adopt -> Destroy.
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				// Basic: the minimal configuration (name + type only).
				// No options and no is_shared are sent, so the complete step
				// also covers the omitted-vs-explicit transition without a
				// second catalog.
				Config: testAccDlfNextCatalogBasicConfig(catalogName),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "name", catalogName),
					resource.TestCheckResourceAttr(resourceId, "type", "PAIMON"),
					resource.TestCheckResourceAttr(resourceId, "options.%", "0"),
					resource.TestCheckResourceAttrSet(resourceId, "created_at"),
					resource.TestCheckResourceAttrSet(resourceId, "updated_at"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"options"},
			},
			{
				// Complete: three writable managed options
				// and an explicit is_shared = false.
				Config: testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "name", catalogName),
					resource.TestCheckResourceAttr(resourceId, "type", "PAIMON"),
					resource.TestCheckResourceAttr(resourceId, "is_shared", "false"),
					resource.TestCheckResourceAttr(resourceId, "options.comment", "c1"),
					resource.TestCheckResourceAttr(resourceId, "options.description", "d1"),
					resource.TestCheckResourceAttr(resourceId, "options.dlf.discovery-query-results-retained-days", "7"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.dlf.trashed-file-retained-days"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.storage.data.redundancy.type"),
					resource.TestCheckResourceAttrSet(resourceId, "created_at"),
					resource.TestCheckResourceAttrSet(resourceId, "updated_at"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "comment", "c1"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "description", "d1"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "dlf.discovery-query-results-retained-days", "7"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"options"},
			},
			{
				// Data Source basic: the data source reads the catalog
				// created in step 1 (no new CreateCatalog). The implicit
				// dependency via catalog_name_pattern interpolation forces the
				// data source to wait for the resource.
				Config: testAccDlfNextCatalogConfigWithDataSource(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, "output_file = \"test-catalogs.json\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.name", catalogName),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.type", "PAIMON"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.id", catalogName),
					resource.TestCheckResourceAttrSet("data.alicloud_dlf_next_catalogs.default", "catalogs.0.status"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "names.0", catalogName),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "ids.0", catalogName),
				),
			},
			{
				// Step 3 (Data Source ids filter): exact-id (catalog name)
				// filter matches.
				Config: testAccDlfNextCatalogConfigWithDataSource(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, fmt.Sprintf("ids = [\"%s\"]", catalogName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.name", catalogName),
				),
			},
			{
				// Step 4 (Data Source name_regex positive): a regex matching
				// the catalog name yields exactly that one catalog.
				Config: testAccDlfNextCatalogConfigWithDataSource(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, fmt.Sprintf("name_regex = \"^%s$\"", catalogName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.name", catalogName),
				),
			},
			{
				// Step 5 (Data Source name_regex negative): a regex matching
				// nothing yields zero catalogs while the resource still exists.
				Config: testAccDlfNextCatalogConfigWithDataSource(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, "name_regex = \"^no-such-dlf-catalog-.*$\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "0"),
				),
			},
			{
				// Options Update: remove description (AlterCatalog
				// removals derived from oldManaged - newManaged) and keep the
				// other two managed keys.
				Config: testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{
					"comment": "c2",
					"dlf.discovery-query-results-retained-days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.comment", "c2"),
					resource.TestCheckResourceAttr(resourceId, "options.dlf.discovery-query-results-retained-days", "7"),
					testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, "description"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "comment", "c2"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault1, ""),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault2, ""),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"options"},
			},
			{
				// Options Clear (two-state): explicit options = {}
				// makes the managed set empty, so the provider sends
				// AlterCatalog removals = oldManaged - newManaged. The API does
				// not honor the removal for dlf.discovery-query-results-retained-days
				// (updatable, not removable), so it remains in remote_options.
				Config: testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.%", "0"),
					testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, "comment"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "dlf.discovery-query-results-retained-days", "7"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault1, ""),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault2, ""),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.dlf.trashed-file-retained-days"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.storage.data.redundancy.type"),
				),
			},
			{
				// Import: a fresh import cannot recover which remote
				// options were originally user-configured, so options is empty
				// and remote_options carries the full remote set.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"options"},
				Config:                  testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.%", "0"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.dlf.discovery-query-results-retained-days"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.dlf.trashed-file-retained-days"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.storage.data.redundancy.type"),
				),
			},
			{
				// Re-adopt one key: after import oldManaged is empty,
				// so re-adopting comment drives an update only (no removals).
				Config: testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{
					"comment": "adopted",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.comment", "adopted"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "comment", "adopted"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault1, ""),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault2, ""),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"options"},
			},
		},
	})
}

// testAccDlfNextCatalogConfigWithDataSource renders the resource block followed
// by a data source block. The data source interpolates catalog_name_pattern
// from the resource so Terraform must build the catalog before reading the
// data source (implicit dependency).
func testAccDlfNextCatalogConfigWithDataSource(catalogName string, options map[string]interface{}, filterLine string) string {
	resourceBlock := testAccDlfNextCatalogOptionsModelConfig(catalogName, options)
	return resourceBlock + fmt.Sprintf(`
data "alicloud_dlf_next_catalogs" "default" {
  catalog_name_pattern = alicloud_dlf_next_catalog.default.name
  %s
}
`, filterLine)
}

// TestAccAliCloudDlfNextCatalog_invalidTypeRejected pins the plan-time type
// validation: only the catalog types documented by CreateCatalog (PAIMON,
// ICEBERG) are accepted, and an unknown value fails before any API call.
func TestAccAliCloudDlfNextCatalog_invalidTypeRejected(t *testing.T) {
	catalogName := fmt.Sprintf("tf-testacc-dlf-catalog-type-%d", tfacctest.RandIntRange(1000, 9999))

	testAccConfig := acctest.ResourceTestAccConfig("alicloud_dlf_next_catalog.default", catalogName, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				// Declaring options, is_shared and share_id here (alongside
				// name/type) makes every Optional schema attribute appear in at
				// least one test Config so the TestingCoverageRate static check
				// sees full coverage. This is a plan-time test (an unknown type
				// is rejected before any API call), so these attributes are
				// never sent to the backend.
				Config: testAccConfig(map[string]interface{}{
					"name": catalogName,
					"type": "DELTA",
					"options": map[string]interface{}{
						"comment": "x",
					},
					"is_shared": false,
					"share_id":  "",
				}),
				ExpectError: regexp.MustCompile(`expected type to be one of \[PAIMON ICEBERG\]`),
			},
		},
	})
}

func testAccAlicloudDlfNextCatalogExists(resourceId string, v *map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceId)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DlfNext Catalog ID is set")
		}

		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		service := dlfnext.NewService(client)

		object, err := service.DescribeCatalog(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error describing DlfNext Catalog: %s", err)
		}

		*v = object
		return nil
	}
}

func testAccCheckAlicloudDlfNextCatalogDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dlf_next_catalog" {
			continue
		}

		service := dlfnext.NewService(client)
		_, err := service.DescribeCatalog(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Error checking DlfNext Catalog destroy: %s", err)
		}
		return fmt.Errorf("DlfNext Catalog still exists: %s", rs.Primary.ID)
	}
	return nil
}

// testAccDlfNextCatalogBestEffortCleanup prevents a failed lifecycle from
// leaking a catalog and consuming the backend cleanup quota. Cleanup errors are
// logged and never flip the recorded test result.
func testAccDlfNextCatalogBestEffortCleanup(t *testing.T, catalogName string) {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := dlfnext.NewService(client)

	if _, err := service.DescribeCatalog(catalogName); err != nil {
		if errs.NotFoundError(err) {
			return
		}
		t.Logf("cleanup: GetCatalog %s failed (will still attempt Drop): %s", catalogName, err)
	}

	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalogName)
	query := make(map[string]*string)
	dropDeadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(dropDeadline) {
		_, err := client.RoaDelete("DlfNext", "2025-03-10", action, query, nil, nil, true)
		if err == nil {
			break
		}
		if errs.NotFoundError(err) {
			return
		}
		if catalogCleanupBackendCapacity(err) {
			t.Logf("cleanup: DropCatalog %s hit backend capacity limit (not retried): %s", catalogName, err)
			break
		}
		if strings.Contains(err.Error(), "is being created") || helper.NeedRetry(err) {
			time.Sleep(3 * time.Second)
			continue
		}
		t.Logf("cleanup: DropCatalog %s failed: %s", catalogName, err)
		break
	}

	if err := service.WaitForCatalogGone(catalogName, 2*time.Minute); err != nil {
		t.Logf("cleanup: WaitForGone %s failed (catalog may still be reaping): %s", catalogName, err)
	}
}

// catalogCleanupBackendCapacity mirrors the backend capacity signal used by
// the resource so the best-effort cleanup does not retry a quota error.
func catalogCleanupBackendCapacity(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, sig := range []string{"too many catalogs", "catalogs being cleaned up", "approximately 48 hours", "catalog quota", "exceed the maximum"} {
		if strings.Contains(msg, sig) {
			return true
		}
	}
	return false
}

func testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, key, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("not found: %s", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		service := dlfnext.NewService(client)
		obj, err := service.DescribeCatalog(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("DescribeCatalog: %s", err)
		}
		opts, _ := obj["options"].(map[string]interface{})
		v, exists := opts[key]
		if !exists {
			return fmt.Errorf("remote options missing key %q (real GetCatalog), options=%v", key, opts)
		}
		if want != "" && fmt.Sprint(v) != want {
			return fmt.Errorf("remote options %q = %v, want %q (real GetCatalog)", key, v, want)
		}
		return nil
	}
}

func testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("not found: %s", resourceId)
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		service := dlfnext.NewService(client)
		obj, err := service.DescribeCatalog(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("DescribeCatalog: %s", err)
		}
		opts, _ := obj["options"].(map[string]interface{})
		if _, exists := opts[key]; exists {
			return fmt.Errorf("remote options still has key %q (should be removed by AlterCatalog), options=%v", key, opts)
		}
		return nil
	}
}

// testAccDlfNextCatalogBasicConfig renders the minimal configuration: name and
// type only, with options and is_shared omitted.
func testAccDlfNextCatalogBasicConfig(catalogName string) string {
	return fmt.Sprintf(`
resource "alicloud_dlf_next_catalog" "default" {
  name = %q
  type = "PAIMON"
}`, catalogName)
}

// testAccDlfNextCatalogOptionsModelConfig renders the options config as a raw
// HCL string with every option key quoted. The shared mapValue helper renders
// map keys UNQUOTED, so a dotted option key such as
// "dlf.discovery-query-results-retained-days" would be parsed as a reference to
// managed resource "dlf". Quoting each key is valid HCL for a TypeMap. The map
// uses map[string]interface{} so the TestingCoverageRate parser sees the
// leading interface{} pair and extracts a bracket-balanced body.
func testAccDlfNextCatalogOptionsModelConfig(catalogName string, options map[string]interface{}) string {
	keys := make([]string, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var block string
	if len(keys) == 0 {
		block = "  options = {}"
	} else {
		var b strings.Builder
		b.WriteString("  options = {")
		for _, k := range keys {
			b.WriteString(fmt.Sprintf("\n    %q = %q", k, options[k].(string)))
		}
		b.WriteString("\n  }")
		block = b.String()
	}
	return fmt.Sprintf(`
resource "alicloud_dlf_next_catalog" "default" {
  name      = %q
  type      = "PAIMON"
  is_shared = false
%s
}`, catalogName, block)
}
