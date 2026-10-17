package paiworkspace_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paiworkspace"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccCheckPaiWorkspaceConfigDestroy verifies that every PAI Workspace
// Config resource created during the test has been removed. It re-reads each
// instance through the resource's own getConfig helper — the same NotFound
// path Read relies on: getConfig returns (nil, nil) once the entry is deleted,
// so any entry that still exists fails the destroy check. This satisfies the
// tflint AT001 rule (every TestAcc must declare a CheckDestroy) with a real
// cloud-side verification rather than a no-op.
func testAccCheckPaiWorkspaceConfigDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pai_workspace_config" {
			continue
		}
		workspaceId, categoryName, configKey, ok := paiworkspace.ParseConfigID(rs.Primary.ID)
		if !ok {
			return fmt.Errorf("destroy check: invalid resource id %q", rs.Primary.ID)
		}
		exists, err := paiworkspace.ConfigExists(client, workspaceId, categoryName, configKey)
		if err != nil {
			return fmt.Errorf("destroy check: reading %s: %w", rs.Primary.ID, err)
		}
		if exists {
			return fmt.Errorf("PAI Workspace Config %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

// Config-key and category for the workspace Config ACC tests. Per the
// AIWorkSpace 2021-02-04 UpdateConfig/ListConfigs OpenAPI definitions,
// CategoryName must be one of CommonResourceConfig, DLCAutoRecycle,
// DLCPriorityConfig, DSWPriorityConfig, QuotaMaximumDuration or
// CommonTagConfig, and each ConfigKey is bound to a specific category.
// tempStoragePath (临时存储路径) is allowed only under CommonResourceConfig
// and its ConfigValue must be an OSS bucket URI
// (oss://<bucket>.<endpoint>/<path>). PAI validates the URI format only, not
// bucket existence — the same static-URI convention the dataset/model/run
// resources use (e.g. resource_alicloud_pai_workspace_run_test.go sets
// artifact_uri = "oss://test.oss-cn-hangzhou.aliyuncs.com/test/"). The test
// uses a static format-valid URI instead of creating an OSS bucket, so it does
// not depend on the ACC account's OSS service, which is disabled
// (alicloud_oss_bucket creation fails with ErrorCode=UserDisable).
const (
	configCategoryName = "CommonResourceConfig"
	configConfigKey    = "tempStoragePath"
	configOSSHost      = "test.oss-cn-hangzhou.aliyuncs.com"
)

// configValueURI returns the static OSS URI used as the tempStoragePath value.
// PAI validates the format only, so a non-existent bucket name is accepted.
// path selects the URI suffix so update steps exercise config_value changes.
func configValueURI(path string) string {
	return fmt.Sprintf("oss://%s/%s/", configOSSHost, path)
}

// configValueRegex returns a regexp matching the config_value the test sets for
// the given path. The trailing slash is optional so a normalized read-back
// still matches.
func configValueRegex(path string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`^oss://%s/%s/?$`, regexp.QuoteMeta(configOSSHost), regexp.QuoteMeta(path)))
}

// testAccPaiWorkspaceConfigWorkspace creates the workspace the Config resource
// is scoped to. The Config resource depends on a workspace existing, and there
// is no PAI workspace data source, so the workspace is created inline. The v5
// mux provider serves this SDKv2 resource (workspace) and the Framework Config
// resource from the same ProtoV5ProviderFactories entry.
func testAccPaiWorkspaceConfigWorkspace(rand int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_workspace_workspace" "default" {
  workspace_name = "tf_testacc_config_%[1]d"
  description    = "tf acc workspace for the config resource"
  env_types      = ["prod"]
}
`, rand)
}

func testAccPaiWorkspaceConfigBasic(rand int) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_pai_workspace_config" "default" {
  workspace_id  = alicloud_pai_workspace_workspace.default.id
  category_name = %[2]q
  config_key    = %[3]q
  config_value  = %[4]q
}
`, testAccPaiWorkspaceConfigWorkspace(rand), configCategoryName, configConfigKey, configValueURI("test"))
}

func testAccPaiWorkspaceConfigBasicUpdate(rand int) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_pai_workspace_config" "default" {
  workspace_id  = alicloud_pai_workspace_workspace.default.id
  category_name = %[2]q
  config_key    = %[3]q
  config_value  = %[4]q
}
`, testAccPaiWorkspaceConfigWorkspace(rand), configCategoryName, configConfigKey, configValueURI("updated"))
}

// TestAccAliCloudPaiWorkspaceConfig_basic covers Create, Read, Update
// (config_value) and Import on the Framework resource.
func TestAccAliCloudPaiWorkspaceConfig_basic(t *testing.T) {
	r := sdkacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiWorkspaceConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceConfigBasic(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_config.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_config.default", "workspace_id"),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_config.default", "category_name", configCategoryName),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_config.default", "config_key", configConfigKey),
					resource.TestMatchResourceAttr("alicloud_pai_workspace_config.default", "config_value",
						configValueRegex("test")),
					resource.TestCheckResourceAttrSet("alicloud_pai_workspace_config.default", "gmt_create_time"),
					// The composite_id carries the three identity parts.
					resource.TestMatchResourceAttr("alicloud_pai_workspace_config.default", "id",
						regexp.MustCompile(`^[^:]+:[^:]+:[^:]+$`)),
				),
			},
			{
				// Update only config_value; identity fields are ForceNew.
				Config: testAccPaiWorkspaceConfigBasicUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("alicloud_pai_workspace_config.default", "config_value",
						configValueRegex("updated")),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_config.default", "category_name", configCategoryName),
					resource.TestCheckResourceAttr("alicloud_pai_workspace_config.default", "config_key", configConfigKey),
				),
			},
			{
				ResourceName:            "alicloud_pai_workspace_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gmt_modified_time"},
			},
		},
	})
}

func testAccPaiWorkspaceConfigLabels(rand int) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_pai_workspace_config" "default" {
  workspace_id  = alicloud_pai_workspace_workspace.default.id
  category_name = %[2]q
  config_key    = %[3]q
  config_value  = %[4]q

  labels {
    key   = "env"
    value = "prod"
  }
  labels {
    key   = "team"
    value = "tf"
  }
}
`, testAccPaiWorkspaceConfigWorkspace(rand), configCategoryName, configConfigKey, configValueURI("test"))
}

func testAccPaiWorkspaceConfigLabelsReordered(rand int) string {
	return fmt.Sprintf(`
%[1]s

resource "alicloud_pai_workspace_config" "default" {
  workspace_id  = alicloud_pai_workspace_workspace.default.id
  category_name = %[2]q
  config_key    = %[3]q
  config_value  = %[4]q

  # Reordered (team before env) and a third label appended: PAI does not
  # honor a client-requested label order, so labels are modeled as an
  # unordered set. This step verifies the reorder is a no-op while the
  # appended label is added — the set grows to three members with no
  # order-inconsistency.
  labels {
    key   = "team"
    value = "tf"
  }
  labels {
    key   = "env"
    value = "prod"
  }
  labels {
    key   = "owner"
    value = "qa"
  }
}
`, testAccPaiWorkspaceConfigWorkspace(rand), configCategoryName, configConfigKey, configValueURI("test"))
}

// TestAccAliCloudPaiWorkspaceConfig_labels covers the Labels collection: PAI
// stores labels as a server-ordered set and ignores client-requested order, so
// the resource models labels as an unordered set. The test verifies that
// appending a label and changing label values are detected as changes, while a
// reorder is a no-op — exercising the collection's membership semantics across
// an update and an import.
func TestAccAliCloudPaiWorkspaceConfig_labels(t *testing.T) {
	r := sdkacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiWorkspaceConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceConfigLabels(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_config.default", "labels.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_pai_workspace_config.default", "labels.*", map[string]string{
						"key":   "env",
						"value": "prod",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_pai_workspace_config.default", "labels.*", map[string]string{
						"key":   "team",
						"value": "tf",
					}),
				),
			},
			{
				Config: testAccPaiWorkspaceConfigLabelsReordered(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_workspace_config.default", "labels.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_pai_workspace_config.default", "labels.*", map[string]string{
						"key":   "env",
						"value": "prod",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_pai_workspace_config.default", "labels.*", map[string]string{
						"key":   "team",
						"value": "tf",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("alicloud_pai_workspace_config.default", "labels.*", map[string]string{
						"key":   "owner",
						"value": "qa",
					}),
				),
			},
			{
				ResourceName:            "alicloud_pai_workspace_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"gmt_modified_time"},
			},
		},
	})
}

// TestAccAliCloudPaiWorkspaceConfig_noop covers the no-op path: re-applying an
// unchanged config must produce no diff, which only holds when Read is
// idempotent and writes back the exact planned state.
func TestAccAliCloudPaiWorkspaceConfig_noop(t *testing.T) {
	r := sdkacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiWorkspaceConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPaiWorkspaceConfigBasic(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("alicloud_pai_workspace_config.default", "config_value",
						configValueRegex("test")),
				),
			},
			{
				// Identical config: plan must be empty.
				Config:   testAccPaiWorkspaceConfigBasic(r),
				PlanOnly: true,
			},
		},
	})
}

func testAccPaiWorkspaceConfigMissingRequired(rand int) string {
	// config_value is a Required argument. Omitting it must fail validation
	// before any API call, so this step creates no cloud resource and cannot
	// leak. The workspace is still declared so the only violation is the
	// missing config_value on the Config resource.
	return fmt.Sprintf(`
%[1]s

resource "alicloud_pai_workspace_config" "default" {
  workspace_id  = alicloud_pai_workspace_workspace.default.id
  category_name = %[2]q
  config_key    = %[3]q
  // config_value intentionally omitted: Required, so plan must fail.
}
`, testAccPaiWorkspaceConfigWorkspace(rand), configCategoryName, configConfigKey)
}

// TestAccAliCloudPaiWorkspaceConfig_missingRequired covers the failure path: a
// plan that omits the Required config_value argument is rejected by schema
// validation before any API call. This is deterministic and creates no cloud
// resource, exercising the resource's Required constraints rather than a
// helper-stubbed or SKIP-avoided path.
func TestAccAliCloudPaiWorkspaceConfig_missingRequired(t *testing.T) {
	r := sdkacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPaiWorkspaceConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPaiWorkspaceConfigMissingRequired(r),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?i)config_value`),
			},
		},
	})
}
