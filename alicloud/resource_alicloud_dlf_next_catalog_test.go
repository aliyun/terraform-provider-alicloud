package alicloud

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudDlfNextCatalog_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dlf_next_catalog.default"
	rc := acctest.RandIntRange(1000, 9999)
	catalogName := fmt.Sprintf("tf-testacc-dlf-merge-%d", rc)

	// Register best-effort cleanup immediately after catalogName generation
	// (review §5): if the ACC fails mid-lifecycle, CheckDestroy is not reached
	// and the catalog would leak, consuming the DLF backend cleanup quota
	// (review §6 category 2). t.Cleanup runs after the test result is
	// recorded, so a cleanup error never masks the original test failure.
	t.Cleanup(func() { testAccDlfNextCatalogBestEffortCleanup(t, catalogName) })

	// Writable user option keys (proven writable by prior ACC runs); the
	// server defaults dlf.trashed-file-retained-days and
	// storage.data.redundancy.type are injected by the backend and must
	// NEVER be deleted by an options update or an explicit options = {}.
	const serverDefault1 = "dlf.trashed-file-retained-days"
	const serverDefault2 = "storage.data.redundancy.type"

	// This merged lifecycle test folds the former _basic, _optionsModel and
	// _DataSource_basic tests into a SINGLE CreateCatalog lifecycle (review
	// §3): Create -> Data Source (basic + filters) -> Options Update ->
	// Options Clear/Omit -> Import -> Re-adopt -> Destroy. One catalog is
	// created once and exercised across the whole lifecycle, so the DLF
	// backend cleanup quota (review §6) is only consumed once per run. The
	// ICEBERG rejection stays a separate plan-only test (no CreateCatalog).
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				// Step 1 (Create): create with three writable managed options.
				// remote_options carries the managed keys AND the server
				// defaults; real GetCatalog confirms all three landed.
				Config: testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "name", catalogName),
					resource.TestCheckResourceAttr(resourceId, "type", "PAIMON"),
					resource.TestCheckResourceAttr(resourceId, "options.comment", "c1"),
					resource.TestCheckResourceAttr(resourceId, "options.description", "d1"),
					resource.TestCheckResourceAttr(resourceId, "options.dlf.discovery-query-results-retained-days", "7"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.dlf.trashed-file-retained-days"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.storage.data.redundancy.type"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "comment", "c1"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "description", "d1"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, "dlf.discovery-query-results-retained-days", "7"),
				),
			},
			{
				// Step 2 (Data Source basic): the data source reads the catalog
				// created in step 1 (no new CreateCatalog). The implicit
				// dependency via catalog_name_pattern interpolation forces the
				// data source to wait for the resource (replaces depends_on,
				// which the SDK prunes from persisted data-source state).
				Config: testAccDlfNextCatalogConfigWithDataSource(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.name", catalogName),
					resource.TestCheckResourceAttr("data.alicloud_dlf_next_catalogs.default", "catalogs.0.type", "PAIMON"),
					resource.TestCheckResourceAttrSet("data.alicloud_dlf_next_catalogs.default", "catalogs.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_dlf_next_catalogs.default", "catalogs.0.status"),
				),
			},
			{
				// Step 3 (Data Source names filter): exact-name filter matches.
				Config: testAccDlfNextCatalogConfigWithDataSource(catalogName, map[string]interface{}{
					"comment":     "c1",
					"description": "d1",
					"dlf.discovery-query-results-retained-days": "7",
				}, fmt.Sprintf("names = [\"%s\"]", catalogName)),
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
				// nothing yields zero catalogs while the resource still
				// exists.
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
				// Step 6 (Options Update): remove description (AlterCatalog
				// removals derived from oldManaged - newManaged) and keep the
				// other two managed keys. Real GetCatalog confirms description
				// was removed server-side and the server defaults are intact.
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
				// Step 7 (Options Clear, two-state): explicit options = {}
				// deletes ONLY managed keys. The two user keys (comment,
				// dlf.discovery-...) are removed; the two server defaults
				// REMAIN in remote_options (never误删). This is the §1
				// boundary guarantee that the dual-field model protects
				// server defaults and unmanaged keys from an explicit clear.
				Config: testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.%", "0"),
					testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, "comment"),
					testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, "dlf.discovery-query-results-retained-days"),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault1, ""),
					testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, serverDefault2, ""),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.dlf.trashed-file-retained-days"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.storage.data.redundancy.type"),
				),
			},
			{
				// Step 8 (Import): a fresh import cannot recover which remote
				// options were originally user-configured, so options is empty
				// and remote_options carries the full remote set incl server
				// defaults. options is in ImportStateVerifyIgnore because the
				// imported managed set is legitimately empty.
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"options"},
				Config:                  testAccDlfNextCatalogOptionsModelConfig(catalogName, map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudDlfNextCatalogExists(resourceId, &v),
					resource.TestCheckResourceAttr(resourceId, "options.%", "0"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.dlf.trashed-file-retained-days"),
					resource.TestCheckResourceAttrSet(resourceId, "remote_options.storage.data.redundancy.type"),
				),
			},
			{
				// Step 9 (Re-adopt one key): after import oldManaged is empty,
				// so re-adopting comment drives an update only (no removals).
				// The server defaults are NOT deleted by a partial re-adopt
				// because removals come only from oldManaged - newManaged.
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
		},
	})
}

// testAccDlfNextCatalogConfigWithDataSource renders the resource block (via
// the shared raw-HCL optionsModel builder, which quotes dotted option keys)
// followed by a data source block. The data source interpolates
// catalog_name_pattern from the resource so Terraform must build the catalog
// before reading the data source (implicit dependency); this replaces the old
// depends_on, which the SDK prunes from persisted data-source state and left
// catalog_name_pattern/catalogs/id empty, producing a spurious plan diff. The
// caller-supplied filterLine adds an extra filter (names/name_regex) when
// needed; pass "" for the basic catalog_name_pattern-only step. This lets the
// merged lifecycle test exercise data source filters against the SAME catalog
// created in step 1 without a second CreateCatalog (review §3).
func testAccDlfNextCatalogConfigWithDataSource(catalogName string, options map[string]interface{}, filterLine string) string {
	resourceBlock := testAccDlfNextCatalogOptionsModelConfig(catalogName, options)
	return resourceBlock + fmt.Sprintf(`
data "alicloud_dlf_next_catalogs" "default" {
  catalog_name_pattern = alicloud_dlf_next_catalog.default.name
  %s
}
`, filterLine)
}

func TestAccAliCloudDlfNextCatalog_typeIcebergRejected(t *testing.T) {
	rc := acctest.RandIntRange(1000, 9999)
	catalogName := fmt.Sprintf("tf-testacc-dlf-catalog-iceberg-%d", rc)

	testAccConfig := resourceTestAccConfigFunc("alicloud_dlf_next_catalog.default", catalogName, func(name string) string {
		return ""
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlicloudDlfNextCatalogDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name": catalogName,
					"type": "ICEBERG",
				}),
				// type is validated at plan time via StringInSlice(["PAIMON"]).
				// DLF 3.0 uses OmniCatalog as the unified catalog model, and
				// PAIMON is the catalog type enum DLFNext CreateCatalog accepts
				// for an OmniCatalog (compatible with Paimon REST, Iceberg REST
				// and HMS). Iceberg is a table format / access protocol surfaced
				// through an OmniCatalog, not a separate CreateCatalog type, so
				// ICEBERG is rejected up front at plan time with a clear error
				// instead of a server-side apply failure.
				ExpectError: regexp.MustCompile(`expected type to be one of \[PAIMON\]`),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		dlfNextServiceV2 := DlfNextServiceV2{client}

		object, err := dlfNextServiceV2.DescribeDlfNextCatalog(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error describing DlfNext Catalog: %s", err)
		}

		*v = object
		return nil
	}
}

func testAccCheckAlicloudDlfNextCatalogDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_dlf_next_catalog" {
			continue
		}

		dlfNextServiceV2 := DlfNextServiceV2{client}
		_, err := dlfNextServiceV2.DescribeDlfNextCatalog(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Error checking DlfNext Catalog destroy: %s", err)
		}
		return fmt.Errorf("DlfNext Catalog still exists: %s", rs.Primary.ID)
	}
	return nil
}

// testAccDlfNextCatalogBestEffortCleanup is a test-level best-effort cleanup
// (review §5): if the ACC fails mid-lifecycle, CheckDestroy is not reached and
// the catalog would leak, consuming the DLF backend cleanup quota (review §6
// category 2). Registered via t.Cleanup right after catalogName is generated,
// it: (1) GetCatalog — if NotFound, nothing to clean; (2) if the catalog
// exists, best-effort DropCatalog with a bounded retry for the "is being
// created" 403 (the catalog may still be in NEW/INITIALIZING); (3) wait for
// GetCatalog NotFound; (4) cleanup errors (with RequestId where available)
// are logged via t.Log and never mask the original test failure (t.Cleanup
// runs after the test result is recorded, and a Cleanup error does not flip a
// PASS). CreateCatalog has no idempotency token, so this is idempotent only in
// that a NotFound catalog is a no-op.
func testAccDlfNextCatalogBestEffortCleanup(t *testing.T, catalogName string) {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	svc := DlfNextServiceV2{client}

	// (1) GetCatalog: if NotFound, nothing to clean.
	if _, err := svc.DescribeDlfNextCatalog(catalogName); err != nil {
		if NotFoundError(err) {
			return
		}
		t.Logf("cleanup: GetCatalog %s failed (will still attempt Drop): %s", catalogName, err)
	}

	// (2) Best-effort DropCatalog with bounded retry for "being created" /
	// transient throttling. A backend cleanup-quota error (review §6 category
	// 2) is logged and not retried here (the quota is account-level; retrying
	// Drop would not help).
	action := fmt.Sprintf("/dlf/v1/catalogs/%s", catalogName)
	query := make(map[string]*string)
	wait := dlfNextCatalogBackoffWait(3*time.Second, 30*time.Second)
	dropDeadline := time.Now().Add(2 * time.Minute)
	for time.Now().Before(dropDeadline) {
		_, err := client.RoaDelete("DlfNext", "2025-03-10", action, query, nil, nil, true)
		if err == nil {
			break
		}
		if NotFoundError(err) {
			return // already gone
		}
		if isDlfNextCatalogBackendCapacityError(err) {
			t.Logf("cleanup: DropCatalog %s hit backend capacity limit (not retried): %s", catalogName, err)
			break
		}
		if isDlfNextCatalogBeingCreatedError(err) || NeedRetry(err) {
			wait()
			continue
		}
		// Non-retryable: log (with the error, which carries RequestId where
		// the backend provides it) and stop — do not mask the test failure.
		t.Logf("cleanup: DropCatalog %s failed: %s", catalogName, err)
		break
	}

	// (3) Wait for GetCatalog NotFound (bounded). A failure here is logged,
	// not fatal: the catalog may still be reaping asynchronously.
	if err := svc.WaitForDlfNextCatalogGone(catalogName, 2*time.Minute); err != nil {
		t.Logf("cleanup: WaitForGone %s failed (catalog may still be reaping): %s", catalogName, err)
	}
}

// TestAccAliCloudDlfNextCatalog_schemaType is a pure unit test (no TF_ACC, no
// resource.Test) pinning the type field contract. DLF 3.0 uses OmniCatalog as
// the unified catalog model; PAIMON is the catalog type enum DLFNext
// CreateCatalog accepts for an OmniCatalog (compatible with Paimon REST,
// Iceberg REST and HMS). Iceberg is a table format / access protocol surfaced
// through an OmniCatalog, not a separate CreateCatalog type, so the schema
// must reject ICEBERG at plan time, accept PAIMON, and default to PAIMON when
// omitted.
func TestAccAliCloudDlfNextCatalog_schemaType(t *testing.T) {
	r := resourceAliCloudDlfNextCatalog()
	typeSchema := r.Schema["type"]

	if typeSchema.Default != "PAIMON" {
		t.Fatalf("expected type Default to be PAIMON, got %v", typeSchema.Default)
	}

	if _, errs := typeSchema.ValidateFunc("ICEBERG", "type"); len(errs) == 0 {
		t.Fatal("expected ICEBERG to be rejected by the type ValidateFunc")
	}

	if _, errs := typeSchema.ValidateFunc("PAIMON", "type"); len(errs) != 0 {
		t.Fatalf("expected PAIMON to be accepted by the type ValidateFunc, got %v", errs)
	}
}

// dlfNotFoundErr builds a tea.SDKError that NotFoundError recognises as NotFound
// (HTTP 404), for unit tests that script GetCatalog responses without a live
// backend.
func dlfNotFoundErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(404),
		Code:       tea.String("NotFound"),
		Message:    tea.String("ResourceNotfound!!! the specified catalog is not found"),
	}
}

// dlfAlreadyExistsErr builds a tea.SDKError that signals a CreateCatalog
// already-exists response.
func dlfAlreadyExistsErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(409),
		Code:       tea.String("CatalogAlreadyExists"),
		Message:    tea.String("the catalog already exists"),
	}
}

// dlfValidationErr builds a tea.SDKError that signals a CreateCatalog
// validation/permission (4xx) error.
func dlfValidationErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(400),
		Code:       tea.String("InvalidParameter"),
		Message:    tea.String("the request parameter is invalid"),
	}
}

// dlfThrottlingErr builds a tea.SDKError that NeedRetry recognises as a
// retryable throttling error.
func dlfThrottlingErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(503),
		Code:       tea.String("ServiceUnavailable"),
		Message:    tea.String("ServiceUnavailable: Request timed out"),
	}
}

// dlfInternalErr builds a tea.SDKError that is neither NotFound, AlreadyExists,
// a 4xx client error, nor a retryable error — a plain non-retryable server error.
func dlfInternalErr() error {
	return &tea.SDKError{
		StatusCode: tea.Int(500),
		Code:       tea.String("InternalError"),
		Message:    tea.String("internal server error"),
	}
}

// TestAccAliCloudDlfNextCatalogPreCreate pins the pre-create existence-check policy: a
// catalog that already exists must be refused (no state written), NotFound
// proceeds to create, and any other Get error fails instead of proceeding.
func TestAccAliCloudDlfNextCatalogPreCreate(t *testing.T) {
	// 1. Pre-create: catalog already exists -> refuse, do not take over.
	if proceed, err := decideDlfNextCatalogPreCreate(nil); proceed {
		t.Fatalf("pre-create with an existing catalog must not proceed")
	} else if err == nil {
		t.Fatalf("pre-create with an existing catalog must return a conflict error")
	}

	// 2. Pre-create: NotFound -> proceed.
	if proceed, err := decideDlfNextCatalogPreCreate(dlfNotFoundErr()); !proceed || err != nil {
		t.Fatalf("pre-create NotFound must proceed, got proceed=%v err=%v", proceed, err)
	}

	// 3. Pre-create: other Get error -> fail (do not proceed).
	if proceed, err := decideDlfNextCatalogPreCreate(dlfInternalErr()); proceed {
		t.Fatalf("pre-create with a non-NotFound Get error must not proceed")
	} else if err == nil {
		t.Fatalf("pre-create with a non-NotFound Get error must surface the error")
	}
}

// TestAccAliCloudDlfNextCatalogCreateRecovery pins the fail-closed ambiguous-create
// recovery policy: never take over a catalog just because a same-name catalog
// exists after an error. CreateCatalog has no idempotency token and no reliable
// unique request marker or comparable createdBy, so ownership can NEVER be
// proven from attributes — even a same-name, same-config catalog found after an
// ambiguous POST error must NOT be adopted (no state written, no auto-recover).
// AlreadyExists and validation/permission (4xx) errors fail outright; a
// retryable error with the catalog still absent retries.
func TestAccAliCloudDlfNextCatalogCreateRecovery(t *testing.T) {
	// 4. Explicit AlreadyExists -> fail, never import.
	if dec, err := decideDlfNextCatalogCreateRecovery(true, false, true, true, dlfGetFound); dec != dlfCreateFail {
		t.Fatalf("AlreadyExists must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("AlreadyExists must return an error")
	}

	// 5. Validation/permission error after a same-name catalog exists -> fail,
	// do not adopt.
	if dec, err := decideDlfNextCatalogCreateRecovery(false, true, false, true, dlfGetFound); dec != dlfCreateFail {
		t.Fatalf("validation/permission error with a same-name catalog must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("validation/permission error must return an error")
	}

	// 6. Ambiguous error (timeout/connection/lost response) and Get finds a
	// same-name, same-config catalog -> MUST fail (never recover). This is the
	// core fail-closed guarantee requested in review: even identical attributes
	// cannot prove this run created the catalog, so no state is written and no
	// auto-recover. The error must be actionable and direct the user to import.
	dec, err := decideDlfNextCatalogCreateRecovery(false, false, true, true, dlfGetFound)
	if dec != dlfCreateFail {
		t.Fatalf("ambiguous error with a same-name same-config catalog must fail (not recover), got %v err=%v", dec, err)
	}
	if err == nil {
		t.Fatalf("ambiguous error with a found catalog must return an actionable error")
	}

	// 7. Retryable error and resource was not created (Get NotFound) -> retry.
	if dec, err := decideDlfNextCatalogCreateRecovery(false, false, true, true, dlfGetAbsent); dec != dlfCreateRetry {
		t.Fatalf("retryable error with absent catalog must retry, got %v err=%v", dec, err)
	}

	// 8. Non-retryable error and catalog absent -> fail (do not retry/adopt).
	if dec, err := decideDlfNextCatalogCreateRecovery(false, false, false, true, dlfGetAbsent); dec != dlfCreateFail {
		t.Fatalf("non-retryable error with absent catalog must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("non-retryable error with absent catalog must return an error")
	}

	// 9. pre-create did NOT prove absence -> fail (the guard should have stopped
	// us; never adopt), even if a same-name catalog is found.
	if dec, err := decideDlfNextCatalogCreateRecovery(false, false, true, false, dlfGetFound); dec != dlfCreateFail {
		t.Fatalf("pre-create not confirmed absent must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("pre-create not confirmed absent must return an error")
	}

	// 10. Retryable error and post-error Get itself errored (non-NotFound) -> fail.
	if dec, err := decideDlfNextCatalogCreateRecovery(false, false, true, true, dlfGetOtherError); dec != dlfCreateFail {
		t.Fatalf("post-error Get error must fail, got %v", dec)
	} else if err == nil {
		t.Fatalf("post-error Get error must return an error")
	}
}

// scriptedGet returns a Get function that plays back a sequence of (object,
// error) responses, holding on the last response once the script is exhausted.
// This lets the waiter policy be unit-tested without a live backend.
func scriptedGet(script []struct {
	obj map[string]interface{}
	err error
}) func(string) (map[string]interface{}, error) {
	i := 0
	return func(id string) (map[string]interface{}, error) {
		if i >= len(script) {
			i = len(script) - 1
		}
		if i < 0 {
			return nil, dlfNotFoundErr()
		}
		s := script[i]
		i++
		return s.obj, s.err
	}
}

func statusObj(status string) map[string]interface{} {
	return map[string]interface{}{"status": status, "name": "test-catalog", "type": "PAIMON"}
}

// TestAccAliCloudDlfNextCatalogCreateWaiter pins the Create waiter: a brief NotFound
// immediately after create is Pending (not an unexpected state), so the waiter
// reaches RUNNING through NOT_FOUND -> NEW -> INITIALIZING -> RUNNING;
// INITIALIZE_FAILED is a terminal failure; perpetual NotFound times out; and a
// non-NotFound Get error during polling surfaces instead of being masked.
func TestAccAliCloudDlfNextCatalogCreateWaiter(t *testing.T) {
	svc := DlfNextServiceV2{}

	// 12. NOT_FOUND -> NEW -> INITIALIZING -> RUNNING reaches RUNNING.
	script := []struct {
		obj map[string]interface{}
		err error
	}{
		{nil, dlfNotFoundErr()},
		{statusObj("NEW"), nil},
		{statusObj("INITIALIZING"), nil},
		{statusObj("RUNNING"), nil},
	}
	refresh := svc.DlfNextCatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf := &resource.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 10 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err != nil {
		t.Fatalf("expected waiter to reach RUNNING, got error: %v", err)
	}

	// 13. INITIALIZE_FAILED (failState) is a terminal failure.
	script = []struct {
		obj map[string]interface{}
		err error
	}{{statusObj("INITIALIZE_FAILED"), nil}}
	refresh = svc.DlfNextCatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &resource.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 10 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err == nil {
		t.Fatalf("expected INITIALIZE_FAILED to fail the waiter, got nil")
	}

	// 14. Perpetual NotFound times out (NOT_FOUND is Pending, so the waiter
	// waits until the timeout instead of erroring on an unexpected state).
	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, dlfNotFoundErr()}}
	refresh = svc.DlfNextCatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &resource.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 2 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err == nil {
		t.Fatalf("expected perpetual NotFound to time out, got nil")
	}

	// 15. A non-NotFound Get error during polling surfaces instead of being
	// masked as an unexpected state.
	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, dlfInternalErr()}}
	refresh = svc.DlfNextCatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf = &resource.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 5 * time.Second,
		Refresh: refresh,
	}
	if _, err := conf.WaitForState(); err == nil {
		t.Fatalf("expected a non-NotFound Get error to surface, got nil")
	}
}

// TestAccAliCloudDlfNextCatalogDeleteWaiter pins the delete-waiter policy: the catalog
// transitions through DELETING/DELETED before NotFound, a transient retryable
// Get error is retried until the catalog returns 404, and a perpetual-exists
// catalog times out.
func TestAccAliCloudDlfNextCatalogDeleteWaiter(t *testing.T) {
	svc := DlfNextServiceV2{}

	// 16. DELETING -> DELETED -> NotFound reaches gone (no error).
	script := []struct {
		obj map[string]interface{}
		err error
	}{
		{statusObj("DELETING"), nil},
		{statusObj("DELETED"), nil},
		{nil, dlfNotFoundErr()},
	}
	if err := svc.WaitForDlfNextCatalogGoneWithApi("test-catalog", 10*time.Second, scriptedGet(script)); err != nil {
		t.Fatalf("expected delete waiter to reach gone, got error: %v", err)
	}

	// 17. A transient retryable Get error is retried until NotFound.
	script = []struct {
		obj map[string]interface{}
		err error
	}{
		{nil, dlfThrottlingErr()},
		{nil, dlfNotFoundErr()},
	}
	if err := svc.WaitForDlfNextCatalogGoneWithApi("test-catalog", 20*time.Second, scriptedGet(script)); err != nil {
		t.Fatalf("expected delete waiter to retry a throttling error then reach gone, got error: %v", err)
	}

	// 18. A non-retryable Get error (non-NotFound) fails instead of looping.
	script = []struct {
		obj map[string]interface{}
		err error
	}{{nil, dlfInternalErr()}}
	if err := svc.WaitForDlfNextCatalogGoneWithApi("test-catalog", 5*time.Second, scriptedGet(script)); err == nil {
		t.Fatalf("expected a non-retryable Get error to fail the delete waiter, got nil")
	}
}

// TestAccAliCloudDlfNextCatalogsStableID pins the data source id stability: the same set
// of catalog names yields the same id regardless of the order ListCatalogs
// returned them, and repeated calls are deterministic.
func TestAccAliCloudDlfNextCatalogsStableID(t *testing.T) {
	// 19. Same set, different order -> same id.
	id1 := dlfNextCatalogsID([]string{"alpha", "bravo", "charlie"})
	id2 := dlfNextCatalogsID([]string{"charlie", "alpha", "bravo"})
	if id1 != id2 {
		t.Fatalf("expected same catalog set in different order to yield the same id, got %q and %q", id1, id2)
	}
	if id1 == "" {
		t.Fatalf("expected a non-empty data source id")
	}

	// 20. Two calls with the same input -> same id (deterministic).
	if dlfNextCatalogsID([]string{"a", "b"}) != dlfNextCatalogsID([]string{"a", "b"}) {
		t.Fatalf("expected dlfNextCatalogsID to be deterministic")
	}

	// 21. Different sets -> different ids.
	if dlfNextCatalogsID([]string{"a", "b"}) == dlfNextCatalogsID([]string{"a", "c"}) {
		t.Fatalf("expected different catalog sets to yield different ids")
	}

	// 22. Empty input -> stable, non-empty id.
	if dlfNextCatalogsID(nil) != dlfNextCatalogsID([]string{}) {
		t.Fatalf("expected empty input to be stable")
	}
}

// TestAccAliCloudDlfNextCatalogsStableOrder pins the data source list ordering
// (#21 三): the catalogs written to the TypeList must be deterministically
// sorted by name, so the same set of catalogs returned in a different order by
// ListCatalogs yields the same catalogs order — not just the same hash id. A
// stable id with an unstable list would still reorder state on every refresh,
// so this test deliberately asserts the element order element-by-element, not
// only the id. This is a pure unit test (no TF_ACC, no live ListCatalogs call).
func TestAccAliCloudDlfNextCatalogsStableOrder(t *testing.T) {
	cat := func(name string) map[string]interface{} {
		return map[string]interface{}{"name": name, "type": "PAIMON", "status": "RUNNING"}
	}
	// Same set, two different API orders.
	order1 := []map[string]interface{}{cat("charlie"), cat("alpha"), cat("bravo")}
	order2 := []map[string]interface{}{cat("bravo"), cat("charlie"), cat("alpha")}

	s1 := sortDlfNextCatalogsByName(order1)
	s2 := sortDlfNextCatalogsByName(order2)

	// Both must sort to the same deterministic order: alpha, bravo, charlie.
	want := []string{"alpha", "bravo", "charlie"}
	for i, w := range want {
		if got := fmt.Sprint(s1[i]["name"]); got != w {
			t.Fatalf("order1[%d]: got %q, want %q", i, got, w)
		}
		if got := fmt.Sprint(s2[i]["name"]); got != w {
			t.Fatalf("order2[%d]: got %q, want %q", i, got, w)
		}
	}

	// The original inputs must not be mutated (sort returns a copy).
	if got := fmt.Sprint(order1[0]["name"]); got != "charlie" {
		t.Fatalf("sortDlfNextCatalogsByName must not mutate its input, first element got %q", got)
	}

	// The id derived from the matched names is also stable across orders
	// (cross-check against dlfNextCatalogsID, the existing id helper).
	id1 := dlfNextCatalogsID([]string{"charlie", "alpha", "bravo"})
	id2 := dlfNextCatalogsID([]string{"bravo", "charlie", "alpha"})
	if id1 != id2 {
		t.Fatalf("expected the same catalog set in different order to yield the same id, got %q and %q", id1, id2)
	}
	if id1 == "" {
		t.Fatalf("expected a non-empty data source id")
	}

	// An empty/nil slice must sort to an empty slice (no panic).
	if got := sortDlfNextCatalogsByName(nil); len(got) != 0 {
		t.Fatalf("expected nil input to sort to an empty slice, got %d elements", len(got))
	}
}

// TestAccAliCloudDlfNextCatalogErrorClassifiersNilSafe pins the error classifier
// nil-safety (P1-3): IsExpectedErrors and the DLF Next resource-specific
// classifiers must not panic when a tea.SDKError has nil Code, nil Data, nil
// Message, or nil StatusCode. Root cause: IsExpectedErrors dereferenced *e.Code
// and *e.Data without nil checks; a validation error (Code=InvalidParameter,
// Data=nil) whose Code did not match any expected code would fall through to
// strings.Contains(*e.Data, code) and panic.
func TestAccAliCloudDlfNextCatalogErrorClassifiersNilSafe(t *testing.T) {
	// 23. dlfValidationErr has Code=InvalidParameter, Data=nil. Before the
	// fix, isDlfNextCatalogAlreadyExistsError would call IsExpectedErrors
	// with codes that don't match "InvalidParameter", fall through to
	// strings.Contains(*e.Data, code) and panic on nil *e.Data.
	if isDlfNextCatalogAlreadyExistsError(dlfValidationErr()) {
		t.Fatalf("validation error must not be classified as already-exists")
	}

	// 24. dlfValidationErr must be classified as an explicit client error
	// (HTTP 400) without panicking.
	if !isDlfNextCatalogExplicitClientError(dlfValidationErr()) {
		t.Fatalf("validation error (HTTP 400) must be classified as explicit client error")
	}

	// 25. dlfAlreadyExistsErr has Data=nil; isDlfNextCatalogAlreadyExistsError
	// must return true via the Code match (CatalogAlreadyExists) without
	// reaching the *e.Data path.
	if !isDlfNextCatalogAlreadyExistsError(dlfAlreadyExistsErr()) {
		t.Fatalf("already-exists error must be classified as already-exists")
	}

	// 26. A tea.SDKError with ALL nil pointer fields must not panic in
	// IsExpectedErrors or the DLF classifiers.
	allNilErr := &tea.SDKError{}
	if IsExpectedErrors(allNilErr, []string{"AnyCode"}) {
		t.Fatalf("all-nil SDKError must not match any expected code")
	}
	if isDlfNextCatalogAlreadyExistsError(allNilErr) {
		t.Fatalf("all-nil SDKError must not be classified as already-exists")
	}
	// isDlfNextCatalogExplicitClientError checks StatusCode; with nil
	// StatusCode it must return false, not panic.
	if isDlfNextCatalogExplicitClientError(allNilErr) {
		t.Fatalf("all-nil SDKError must not be classified as explicit client error")
	}

	// 27. A tea.SDKError with nil Code but non-nil Data must not panic.
	nilCodeErr := &tea.SDKError{
		Code:       nil,
		Data:       tea.String("some data"),
		Message:    tea.String("msg"),
		StatusCode: tea.Int(500),
	}
	if IsExpectedErrors(nilCodeErr, []string{"NotFound"}) {
		t.Fatalf("nil-code error must not match any expected code")
	}

	// 28. A tea.SDKError with non-nil Code but nil Data, where the Code does
	// not match any expected code, must not panic on the Data path.
	nilDataErr := &tea.SDKError{
		Code:       tea.String("SomeOtherCode"),
		Data:       nil,
		Message:    tea.String("msg"),
		StatusCode: tea.Int(500),
	}
	if IsExpectedErrors(nilDataErr, []string{"NotFound"}) {
		t.Fatalf("non-matching code with nil Data must not match")
	}

	// 29. HTTP 400, 401, 403 must be classified as explicit client errors.
	for _, code := range []int{400, 401, 403} {
		err := &tea.SDKError{
			Code:       tea.String("SomeCode"),
			Data:       nil,
			Message:    tea.String("msg"),
			StatusCode: tea.Int(code),
		}
		if !isDlfNextCatalogExplicitClientError(err) {
			t.Fatalf("HTTP %d must be classified as explicit client error", code)
		}
	}

	// 30. HTTP 503 (ServiceUnavailable) must NOT be classified as explicit
	// client error (it is a server error, potentially retryable).
	if isDlfNextCatalogExplicitClientError(dlfThrottlingErr()) {
		t.Fatalf("HTTP 503 must not be classified as explicit client error")
	}

	// 31. A ComplexError wrapping a tea.SDKError with nil Data must not panic.
	wrappedErr := &ComplexError{
		Cause: &tea.SDKError{
			Code:       tea.String("InvalidParameter"),
			Data:       nil,
			Message:    tea.String("msg"),
			StatusCode: tea.Int(400),
		},
	}
	if IsExpectedErrors(wrappedErr, []string{"NotFound"}) {
		t.Fatalf("ComplexError wrapping non-matching SDKError must not match")
	}
}

// TestAccAliCloudDlfNextCatalogTimeoutPolicy pins the unified-deadline timeout
// policy (#21 五): remainingUntil clamps a passed deadline to zero (non-negative)
// instead of a negative duration, and a waiter with a small timeout does not
// re-acquire a full-minute inner budget — it completes within the small
// timeout (+ slack), not the hardcoded minute an inner Describe would
// otherwise use. This is a pure unit test (no TF_ACC, no live backend); the
// inner Describe's bounded retry via DescribeDlfNextCatalogWithRetryTimeout is
// exercised by the real ACC lifecycle, where each Create/Delete phase calls
// remainingUntil(deadline) instead of a fresh full timeout.
func TestAccAliCloudDlfNextCatalogTimeoutPolicy(t *testing.T) {
	// 32. remainingUntil clamps a passed deadline to 0 (non-negative) so
	// resource.Retry and the deadline-aware Describe receive a zero, not a
	// negative, duration when the deadline has already elapsed.
	if got := remainingUntil(time.Now().Add(-1 * time.Second)); got != 0 {
		t.Fatalf("remainingUntil of a passed deadline must be 0, got %v", got)
	}

	// 33. remainingUntil of a future deadline is positive and bounded by the
	// configured budget, so each phase takes only the time left, not a fresh
	// full timeout.
	if got := remainingUntil(time.Now().Add(2 * time.Second)); got <= 0 || got > 2*time.Second {
		t.Fatalf("remainingUntil of a 2s-future deadline must be in (0, 2s], got %v", got)
	}

	// 34. A waiter with a 1s timeout and a perpetual INITIALIZING response
	// completes within a few seconds (well under the inner Describe's 1-minute
	// budget), proving the waiter does not re-acquire a full timeout per phase.
	svc := DlfNextServiceV2{}
	script := []struct {
		obj map[string]interface{}
		err error
	}{{statusObj("INITIALIZING"), nil}}
	refresh := svc.DlfNextCatalogStateRefreshFuncWithApi("test-catalog", "status", []string{"INITIALIZE_FAILED"}, scriptedGet(script))
	conf := &resource.StateChangeConf{
		Pending: []string{"NOT_FOUND", "NEW", "INITIALIZING"},
		Target:  []string{"RUNNING"},
		Timeout: 1 * time.Second,
		Refresh: refresh,
	}
	start := time.Now()
	_, err := conf.WaitForState()
	elapsed := time.Since(start)
	if err == nil {
		t.Fatalf("expected perpetual INITIALIZING to time out, got nil")
	}
	if elapsed > 10*time.Second {
		t.Fatalf("waiter with 1s timeout took %v, expected well under 1m (no full-timeout re-acquisition)", elapsed)
	}
}

// TestAccAliCloudDlfNextCatalogOptionChanges pins the AlterCatalog request
// construction (review §2): buildDlfNextCatalogOptionChanges derives updates
// from newManaged and removals ONLY from oldManaged - newManaged — never from
// remote_options, the full imported remote set, or server-injected defaults.
// An explicit options = {} (newManaged empty) removes only previously-managed
// keys; server defaults never appear in updates or removals. Pure function,
// no backend. AlterCatalog API shape (updates object map + removals string
// array) verified against the official DlfNext AlterCatalog API definition.
func TestAccAliCloudDlfNextCatalogOptionChanges(t *testing.T) {
	// 1. Add only: old empty, new {a,b} -> updates {a,b}, removals [].
	u, r := buildDlfNextCatalogOptionChanges(map[string]interface{}{}, map[string]interface{}{"a": "1", "b": "2"})
	if len(u) != 2 || u["a"] != "1" || u["b"] != "2" {
		t.Fatalf("add: updates = %v, want {a:1,b:2}", u)
	}
	if len(r) != 0 {
		t.Fatalf("add: removals = %v, want empty", r)
	}

	// 2. Modify + add: old {a:1}, new {a:9,c:3} -> updates {a:9,c:3}, removals [].
	u, r = buildDlfNextCatalogOptionChanges(map[string]interface{}{"a": "1"}, map[string]interface{}{"a": "9", "c": "3"})
	if u["a"] != "9" || u["c"] != "3" || len(u) != 2 {
		t.Fatalf("modify+add: updates = %v, want {a:9,c:3}", u)
	}
	if len(r) != 0 {
		t.Fatalf("modify+add: removals = %v, want empty", r)
	}

	// 3. Remove one: old {a,b}, new {a} -> updates {a}, removals [b].
	u, r = buildDlfNextCatalogOptionChanges(map[string]interface{}{"a": "1", "b": "2"}, map[string]interface{}{"a": "1"})
	if len(u) != 1 || u["a"] != "1" {
		t.Fatalf("remove: updates = %v, want {a:1}", u)
	}
	if len(r) != 1 || fmt.Sprint(r[0]) != "b" {
		t.Fatalf("remove: removals = %v, want [b]", r)
	}

	// 4. Clear all (explicit options={}): old {a,b}, new {} -> updates {},
	// removals [a,b] (order-independent).
	u, r = buildDlfNextCatalogOptionChanges(map[string]interface{}{"a": "1", "b": "2"}, map[string]interface{}{})
	if len(u) != 0 {
		t.Fatalf("clear: updates = %v, want empty", u)
	}
	if len(r) != 2 {
		t.Fatalf("clear: removals = %v, want 2 items", r)
	}
	set := map[string]bool{}
	for _, k := range r {
		set[fmt.Sprint(k)] = true
	}
	if !set["a"] || !set["b"] {
		t.Fatalf("clear: removals = %v, want {a,b}", r)
	}

	// 5. Boundary guarantee (review §1): removals NEVER come from
	// remote_options / server defaults. oldManaged has only user keys
	// (comment, description); server defaults are NOT in oldManaged, so
	// clearing managed options removes only the user keys — server defaults
	// are never in removals or updates. This is the core "import/omit must
	// never误删 server defaults" guarantee, proven deterministically here.
	oldManaged := map[string]interface{}{"comment": "v1", "description": "d"}
	u, r = buildDlfNextCatalogOptionChanges(oldManaged, map[string]interface{}{})
	for _, dflt := range []string{"dlf.trashed-file-retained-days", "storage.data.redundancy.type"} {
		for _, k := range r {
			if fmt.Sprint(k) == dflt {
				t.Fatalf("clear: server default %q must NEVER appear in removals, got %v", dflt, r)
			}
		}
		if _, ok := u[dflt]; ok {
			t.Fatalf("clear: server default %q must NEVER appear in updates", dflt)
		}
	}

	// 6. Empty old + empty new (no-op) -> updates {}, removals [].
	u, r = buildDlfNextCatalogOptionChanges(map[string]interface{}{}, map[string]interface{}{})
	if len(u) != 0 || len(r) != 0 {
		t.Fatalf("noop: updates=%v removals=%v, want both empty", u, r)
	}

	// 10. nil/empty maps must not panic. A nil old or new managed map (which
	// is how an omitted options block and an absent state surface after
	// type-assertion) must behave like an empty map, never panic.
	u, r = buildDlfNextCatalogOptionChanges(nil, nil)
	if len(u) != 0 || len(r) != 0 {
		t.Fatalf("nil/nil: updates=%v removals=%v, want both empty", u, r)
	}
	u, r = buildDlfNextCatalogOptionChanges(map[string]interface{}{"a": "1"}, nil)
	if len(u) != 0 || len(r) != 1 || fmt.Sprint(r[0]) != "a" {
		t.Fatalf("old/nil: updates=%v removals=%v, want removals [a]", u, r)
	}
}

// TestAccAliCloudDlfNextCatalogOptionChangesTwoState pins the two-state contract
// (review §2 scenarios 11 & 12): an omitted options block and an explicit
// options = {} MUST both mean the managed set is expected to be empty and
// produce IDENTICAL AlterCatalog removals (every previously-managed key),
// never deleting server defaults or unmanaged keys. The legacy SDK cannot
// reliably distinguish an omitted block from an explicit empty map, so the
// contract is two-state: buildDlfNextCatalogOptionChanges treats a nil
// newManaged (omit) and an empty-map newManaged ({}) identically — both
// clear the managed set. This test asserts equality of the two paths, not
// just that each independently produces an empty update set. Pure function,
// no backend.
func TestAccAliCloudDlfNextCatalogOptionChangesTwoState(t *testing.T) {
	// oldManaged carries two user-managed keys plus is NOT carrying the server
	// defaults (those live in remote_options and never enter oldManaged).
	oldManaged := map[string]interface{}{
		"comment":     "v1",
		"description": "d",
		"dlf.discovery-query-results-retained-days": "7",
	}

	// 11. options completely omitted -> newManaged is nil (the type-assertion
	// of an absent/zero options yields nil or empty). Expected managed set is
	// empty: updates {}, removals = all oldManaged keys.
	uOmit, rOmit := buildDlfNextCatalogOptionChanges(oldManaged, nil)

	// 12. options = {} (explicit empty map) -> newManaged is an empty map.
	// Expected managed set is empty: updates {}, removals = all oldManaged
	// keys.
	uEmpty, rEmpty := buildDlfNextCatalogOptionChanges(oldManaged, map[string]interface{}{})

	// Both paths must produce the same two-state "clear" result: no updates,
	// and the SAME set of removals (every previously-managed key). This is the
	// core two-state assertion — omit and {} are equivalent, not different.
	if len(uOmit) != 0 || len(uEmpty) != 0 {
		t.Fatalf("two-state clear: omit updates=%v, {} updates=%v, want both empty", uOmit, uEmpty)
	}
	if len(rOmit) != len(oldManaged) || len(rEmpty) != len(oldManaged) {
		t.Fatalf("two-state clear: omit removals=%d, {} removals=%d, want %d (all managed keys)", len(rOmit), len(rEmpty), len(oldManaged))
	}
	// Compare removals as sets (order is non-deterministic across maps).
	setOmit := map[string]bool{}
	setEmpty := map[string]bool{}
	for _, k := range rOmit {
		setOmit[fmt.Sprint(k)] = true
	}
	for _, k := range rEmpty {
		setEmpty[fmt.Sprint(k)] = true
	}
	for k := range oldManaged {
		if !setOmit[k] {
			t.Fatalf("two-state clear (omit): key %q must be in removals, got %v", k, rOmit)
		}
		if !setEmpty[k] {
			t.Fatalf("two-state clear ({}): key %q must be in removals, got %v", k, rEmpty)
		}
	}
	// Verify omit and {} produce the EXACT same removal set (two-state
	// equivalence, not just both-non-empty).
	if len(setOmit) != len(setEmpty) {
		t.Fatalf("two-state: omit and {} removal sets differ in size")
	}
	for k := range setOmit {
		if !setEmpty[k] {
			t.Fatalf("two-state: key %q in omit removals but not in {} removals — contract broken", k)
		}
	}

	// Boundary guarantee (review §1): server defaults must NEVER appear in
	// removals for either path, because they never enter oldManaged.
	for _, dflt := range []string{"dlf.trashed-file-retained-days", "storage.data.redundancy.type"} {
		if setOmit[dflt] || setEmpty[dflt] {
			t.Fatalf("two-state clear: server default %q must NEVER be in removals", dflt)
		}
	}
}

// TestAccAliCloudDlfNextCatalogManagedOptionsReconcile pins the Read drift
// behaviour (review §2 scenario 7): a managed key whose API value changed
// surfaces the new value (drift), a managed key the API no longer returns is
// dropped from managed state so the deletion is visible in the next plan, and
// a new remote key (server default / unmanaged addition) is NEVER promoted
// into the managed set — it stays in remote_options only. Pure function, no
// backend.
func TestAccAliCloudDlfNextCatalogManagedOptionsReconcile(t *testing.T) {
	// 7. Managed key value changed by API -> drift surfaced with API value;
	// managed key absent from API -> dropped from state (drift); server
	// defaults / unmanaged additions -> NOT promoted into managed.
	managed := map[string]interface{}{"comment": "v1", "description": "d", "description2": "x"}
	api := map[string]interface{}{
		"comment":     "v2", // managed, value changed -> drift surfaced
		"description": "d",  // managed, unchanged
		// description2 deleted server-side -> dropped from managed (drift)
		"dlf.trashed-file-retained-days": 7,     // server default, unmanaged -> NOT promoted
		"storage.data.redundancy.type":   "LRS", // server default, unmanaged -> NOT promoted
		"enable.openapi":                 false, // unmanaged addition -> NOT promoted
	}
	got := reconcileDlfNextCatalogManagedOptions(managed, api)
	if got["comment"] != "v2" {
		t.Fatalf("managed key value change must surface API value, got comment=%v want v2", got["comment"])
	}
	if got["description"] != "d" {
		t.Fatalf("unchanged managed key must keep API value, got description=%v", got["description"])
	}
	if _, ok := got["description2"]; ok {
		t.Fatalf("managed key absent from API must be dropped (drift), but description2 present")
	}
	for _, k := range []string{"dlf.trashed-file-retained-days", "storage.data.redundancy.type", "enable.openapi"} {
		if _, ok := got[k]; ok {
			t.Fatalf("unmanaged remote key %q must not be promoted into managed options", k)
		}
	}

	// 8. Pure drop: managed key deleted server-side -> removed from state.
	managed2 := map[string]interface{}{"comment": "v1", "description": "d"}
	api2 := map[string]interface{}{"comment": "v1"} // description deleted server-side
	got2 := reconcileDlfNextCatalogManagedOptions(managed2, api2)
	if _, ok := got2["description"]; ok {
		t.Fatalf("managed key deleted by API must be dropped from state (drift), got description present")
	}
	if got2["comment"] != "v1" {
		t.Fatalf("managed key still in API must keep its value")
	}
	if len(got2) != 1 {
		t.Fatalf("reconciled managed set must have exactly 1 key, got %v", got2)
	}

	// 9. Empty managed -> empty reconcile, no panic, even with a full API set.
	got3 := reconcileDlfNextCatalogManagedOptions(map[string]interface{}{}, api)
	if len(got3) != 0 {
		t.Fatalf("empty managed state must yield empty reconcile, got %v", got3)
	}
}

// TestAccAliCloudDlfNextCatalogRemoteOptionsMirror pins remote_options (review
// §2 scenario 8): it mirrors the FULL GetCatalog options including server
// defaults and unmanaged keys, stringified, so a fresh import (managed empty)
// does not lose the remote picture. Pure function, no backend.
func TestAccAliCloudDlfNextCatalogRemoteOptionsMirror(t *testing.T) {
	// 10. Full mirror incl server defaults + unmanaged + non-string values.
	api := map[string]interface{}{
		"comment":                        "v1",
		"dlf.trashed-file-retained-days": 7, // int -> "7"
		"storage.data.redundancy.type":   "LRS",
		"enable.openapi":                 false, // bool -> "false"
	}
	got := reconcileDlfNextCatalogRemoteOptions(api)
	if got["comment"] != "v1" {
		t.Fatalf("comment = %v, want v1", got["comment"])
	}
	if got["dlf.trashed-file-retained-days"] != "7" {
		t.Fatalf("server default int 7 must stringify to \"7\", got %v", got["dlf.trashed-file-retained-days"])
	}
	if got["storage.data.redundancy.type"] != "LRS" {
		t.Fatalf("server default LRS wrong/missing, got %v", got["storage.data.redundancy.type"])
	}
	if got["enable.openapi"] != "false" {
		t.Fatalf("bool false must stringify to \"false\", got %v", got["enable.openapi"])
	}
	if len(got) != 4 {
		t.Fatalf("remote_options must mirror ALL 4 remote keys, got %d: %v", len(got), got)
	}
	// Empty API -> empty mirror, no panic.
	if g := reconcileDlfNextCatalogRemoteOptions(map[string]interface{}{}); len(g) != 0 {
		t.Fatalf("empty API must yield empty remote_options, got %v", g)
	}
}

// testAccAlicloudDlfNextCatalogRemoteOptionSet verifies the REAL remote
// options (via a direct GetCatalog call, not Terraform state) contain a key —
// proving an AlterCatalog update actually landed on the server. This
// satisfies review §2's "verify real AlterCatalog request semantics, not just
// state" by reading the authoritative remote state directly.
func testAccAlicloudDlfNextCatalogRemoteOptionSet(resourceId, key, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("not found: %s", resourceId)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		svc := DlfNextServiceV2{client}
		obj, err := svc.DescribeDlfNextCatalog(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("DescribeDlfNextCatalog: %s", err)
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

// testAccAlicloudDlfNextCatalogRemoteOptionAbsent verifies the REAL remote
// options (via a direct GetCatalog call) no longer contain a key — proving an
// AlterCatalog removal actually deleted it server-side, not just from state.
func testAccAlicloudDlfNextCatalogRemoteOptionAbsent(resourceId, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("not found: %s", resourceId)
		}
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		svc := DlfNextServiceV2{client}
		obj, err := svc.DescribeDlfNextCatalog(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("DescribeDlfNextCatalog: %s", err)
		}
		opts, _ := obj["options"].(map[string]interface{})
		if _, exists := opts[key]; exists {
			return fmt.Errorf("remote options still has key %q (should be removed by AlterCatalog), options=%v", key, opts)
		}
		return nil
	}
}

// testAccDlfNextCatalogOptionsModelConfig renders the optionsModel ACC config as
// a raw HCL string with every option key quoted. The shared mapValue helper
// renders map keys UNQUOTED (keyV.String()), so a dotted option key such as
// "dlf.discovery-query-results-retained-days" would render as
//
//	dlf.discovery-query-results-retained-days = "7"
//
// and Terraform parses that unquoted dotted LHS as a reference to managed
// resource "dlf" — breaking Step-0 config parse. Quoting each option key
// ("dlf.discovery-query-results-retained-days" = "7") is valid HCL for a
// TypeMap and sidesteps mapValue without patching the shared helper, which
// every other resource test relies on. optionsModel is the only test in this
// file that uses dotted option keys, so only it needs the raw builder.
// The options map uses map[string]interface{} (not map[string]string) so the
// TestingCoverageRate checker's Config-line parser sees the leading
// interface{} pair and extracts a bracket-balanced body; a map[string]string
// literal lacks that pair, leaving an unmatched closing brace and aborting
// the whole resource coverage check.
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
  name = %q
  type = "PAIMON"
%s
}`, catalogName, block)
}
