package dlfnext_test

import (
	"context"
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
	catalogName := fmt.Sprintf("tf-testAccDlfCatalog%d", tfacctest.RandIntRange(1000, 9999))

	t.Cleanup(func() { testAccDlfNextCatalogBestEffortCleanup(t, catalogName) })

	const serverDefault1 = "dlf.trashed-file-retained-days"
	const serverDefault2 = "storage.data.redundancy.type"

	// Single complete CreateCatalog lifecycle (the backend cleanup quota makes a
	// second CreateCatalog expensive): complete configuration -> data source
	// (basic, ids, regex) -> options update -> clear -> import -> re-adopt.
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				// Complete Create: options and an explicit is_shared = false.
				Config: testAccDlfNextCatalogCompleteConfig(catalogName, map[string]string{
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
				// Data source basic: pattern bound to the resource plus output_file.
				Config: testAccDlfNextCatalogWithDataSource(catalogName, map[string]string{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, `output_file = "catalogs.json"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.name", catalogName),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.type", "PAIMON"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.id", catalogName),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "names.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "names.0", catalogName),
				),
			},
			{
				// ids filter WITHOUT the server-side pattern, so the filter is
				// proven independently of catalogNamePattern.
				Config: testAccDlfNextCatalogWithDataSource(catalogName, map[string]string{
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
				Config: testAccDlfNextCatalogWithDataSource(catalogName, map[string]string{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, fmt.Sprintf("name_regex = \"^%s$\"", catalogName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "1"),
				),
			},
			{
				Config: testAccDlfNextCatalogWithDataSource(catalogName, map[string]string{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, "name_regex = \"^no-such-dlf-catalog-.*$\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "0"),
				),
			},
			{
				// Options update: remove description via AlterCatalog removals.
				Config: testAccDlfNextCatalogCompleteConfig(catalogName, map[string]string{
					"comment": "c2",
					"dlf.discovery-query-results-retained-days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.comment", "c2"),
					resource.TestCheckResourceAttr(resourceId, "options.dlf.discovery-query-results-retained-days", "7"),
					testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, "description"),
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
				// Options clear (two-state).
				Config: testAccDlfNextCatalogCompleteConfig(catalogName, map[string]string{}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.%", "0"),
					testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, "comment"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "dlf.discovery-query-results-retained-days", "7"),
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
				Config: testAccDlfNextCatalogCompleteConfig(catalogName, map[string]string{"comment": "adopted"}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.comment", "adopted"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "comment", "adopted"),
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

// TestAccAliCloudDlfNextCatalog_invalidTypeRejected is a plan-only check that an
// undocumented type is rejected by the schema validator before any API call.
func TestAccAliCloudDlfNextCatalog_invalidTypeRejected(t *testing.T) {
	catalogName := fmt.Sprintf("tf-testAccDlfCatalogType%d", tfacctest.RandIntRange(1000, 9999))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "alicloud_dlf_next_catalog" "default" {
  name = %q
  type = "DELTA"
}`, catalogName),
				ExpectError: regexp.MustCompile(`must be one of`),
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
		object, err := dlfnext.NewService(client).DescribeCatalog(context.Background(), rs.Primary.ID)
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
		_, err := dlfnext.NewService(client).DescribeCatalog(context.Background(), rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("Error checking DlfNext Catalog destroy: %s", err)
		}
		return fmt.Errorf("DlfNext Catalog still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccDlfNextCatalogBestEffortCleanup(t *testing.T, catalogName string) {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	service := dlfnext.NewService(client)
	ctx := t.Context()

	if _, err := service.DescribeCatalog(ctx, catalogName); err != nil {
		if errs.NotFoundError(err) {
			return
		}
		t.Logf("cleanup: GetCatalog %s failed (will still attempt Drop): %s", catalogName, err)
	}

	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalogName)
	query := make(map[string]*string)
	deadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(deadline) {
		_, err := client.RoaDelete("DlfNext", "2025-03-10", action, query, nil, nil, false)
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

	if err := service.WaitForCatalogGone(ctx, catalogName, 2*time.Minute); err != nil {
		t.Logf("cleanup: WaitForGone %s failed (catalog may still be reaping): %s", catalogName, err)
	}
}

func catalogCleanupBackendCapacity(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, sig := range []string{"too many catalogs", "catalogs being cleaned up", "approximately 48 hours", "catalog quota"} {
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
		obj, err := dlfnext.NewService(client).DescribeCatalog(context.Background(), rs.Primary.ID)
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
		obj, err := dlfnext.NewService(client).DescribeCatalog(context.Background(), rs.Primary.ID)
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

// testAccDlfNextCatalogCompleteConfig renders the complete resource config with
// quoted dotted option keys and an explicit is_shared = false.
func testAccDlfNextCatalogCompleteConfig(catalogName string, options map[string]string) string {
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
			b.WriteString(fmt.Sprintf("\n    %q = %q", k, options[k]))
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

// testAccDlfNextCatalogWithDataSource renders the resource and a data source
// with the extra filter line. The resource block is complete so every step
// keeps the same managed configuration.
func testAccDlfNextCatalogWithDataSource(catalogName string, options map[string]string, filterLine string) string {
	return testAccDlfNextCatalogCompleteConfig(catalogName, options) + fmt.Sprintf(`
data "alicloud_dlf_next_catalogs" "default" {
  %s
}
`, filterLine)
}
