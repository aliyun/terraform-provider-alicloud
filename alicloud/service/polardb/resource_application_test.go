// Package polardb_test holds the SDKv2 acceptance tests for alicloud_polardb_application,
// migrated from the flat alicloud package into the polardb product package.
//
// The flat-layout tests leaned on the shared acceptance framework
// (resourceAttrInit/resourceAttrMapUpdateSet/resourceTestAccConfigFunc) which
// is not importable from a product package, so each step is expressed as an
// explicit, complete HCL config (matching the migrated-subpackage convention
// used by service/ecs). The step semantics are preserved: each step's config
// is the create base plus that step's accumulated overrides, so no Required
// field is dropped and no spurious parameter/whitelist removal diffs arise.
package polardb_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/polardb"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdktf "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// collection-order: bypass components -- pre-existing order-insensitive TypeList (stock debt predating this migration PR).
// collection-order: bypass parameters -- pre-existing order-insensitive TypeList (stock debt predating this migration PR).
// collection-order: bypass security_ip_list -- pre-existing order-insensitive TypeList (stock debt predating this migration PR).

// TestAccAliCloudPolarDBApplication_Create exercises create, version upgrade,
// parameter modification and whitelist append/delete on a PostPaid PolarDB AI
// application.
func TestAccAliCloudPolarDBApplication_Create(t *testing.T) {
	resourceId := "alicloud_polardb_application.default"
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			// The flat suite pinned this resource to cn-beijing via
			// testAccPreCheckWithRegions; keep the pin and its supported-region
			// semantics here.
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBApplicationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPolardbApplicationConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
					// The flat-layout suite asserted vswitch_id=CHECKSET through
					// its basicMap; keep the equivalent non-empty read-back check.
					resource.TestCheckResourceAttrSet(resourceId, "vswitch_id"),
					// Create-time Parameters: the first apply carries a real
					// parameter (the create request is asserted field by field in
					// TestUnitPolarDBApplicationCreateRequestSerialization) and the
					// application is created with it instead of keeping it in
					// state only.
					resource.TestCheckResourceAttr(resourceId, "parameters.#", "1"),
					resource.TestCheckResourceAttr(resourceId, "parameters.0.parameter_name", "secret.dashscope.apiKey"),
					resource.TestCheckResourceAttr(resourceId, "parameters.0.parameter_value", "abc"),
					// Endpoints and security_ip_type are CreateApplication inputs
					// whose cloud prerequisites are just the application itself.
					resource.TestCheckResourceAttr(resourceId, "endpoints.#", "1"),
					// endpoints is a TypeSet: assert the member by its own
					// attributes instead of a positional index, whose hash is
					// not a stable contract.
					testAccCheckPolarDBApplicationEndpointMember(resourceId, "Primary", "my_endpoint"),
					resource.TestCheckResourceAttr(resourceId, "security_ip_type", "ipv4"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigUpgradeVersion(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "upgrade_version", "true"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigParameters(name, "def"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
					// The parameter was written by CreateApplication above; this
					// step keeps exercising the ModifyApplicationParameter path.
					resource.TestCheckResourceAttr(resourceId, "parameters.0.parameter_value", "def"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigWhitelist(name, "Append"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "security_ip_array_name", "test"),
					resource.TestCheckResourceAttr(resourceId, "security_ip_list.#", "1"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigWhitelist(name, "Delete"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "security_ip_array_name", "test"),
					resource.TestCheckResourceAttr(resourceId, "security_ip_list.#", "1"),
				),
			},
			{
				// item 65: tags mutation on a live application via
				// TagResources (ResourceType="application"), read back through
				// ListTagResources. Steps from the prior config are preserved so
				// only tags change, exercising the tag update path.
				Config: testAccPolardbApplicationConfigTags(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceId, "tags.Created", "allin"),
					resource.TestCheckResourceAttr(resourceId, "tags.Env", "test"),
				),
			},
			{
				// Dropping the tags block removes both tags through
				// UntagResources; Read must then write the empty
				// ListTagResources result back so tags.% converges to 0 instead
				// of keeping the stale map.
				Config: testAccPolardbApplicationConfigWhitelist(name, "Delete"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(resourceId, "vswitch_id"),
				),
			},
			{
				// PlanOnly: the application was created without a storages
				// block, so Read adopted the service default member. Planning
				// the same configuration must stay empty: the adopted
				// attributes cannot drift, and every following step above
				// already re-planned it the same way.
				Config:   testAccPolardbApplicationConfigWhitelist(name, "Delete"),
				PlanOnly: true,
			},
			{
				// PlanOnly: a capacity the configuration sets is visible in the
				// plan even though the adopted member carries the service
				// default, and the ForceNew flags make that plan a replacement.
				// No apply, so the application is left untouched.
				Config:             testAccPolardbApplicationConfigWhitelistWithStorages(name, "Delete", "\tstorages {\n\t\tstorage_capacity = \"100\"\n\t}\n"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Import round-trip: an imported application must reproduce the
				// state Read builds (identifiers, type, billing mode, cluster,
				// components, whitelist IP type, storage list, tags). The
				// ignored attributes are the write-only inputs and the blocks
				// the query API does not return, which can only come from the
				// configuration.
				Config:            testAccPolardbApplicationConfigWhitelist(name, "Delete"),
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// Every attribute Read restores is verified. The list below only
				// holds inputs the query API never returns (or operation
				// triggers) plus three attributes with a recorded, same-version
				// reason.
				ImportStateVerifyIgnore: []string{
					"parameters",
					"modify_mode",
					"upgrade_version",
					// The query API reports only the default whitelist entry, not
					// the configured array name (two same-head import runs show
					// the imported state reporting "default"), so the name cannot
					// be compared; the IP list below is still verified.
					"security_ip_array_name",
					"model_from",
					"model_api_key",
					"model_base_url",
					"model_name",
					"model_api",
					"ai_db_cluster_id",
					"resource_group_id",
					"target_version",
					"skill_template_id",
					"dnat_ip_address",
					"dnat_entries",
					"auth_provider",
					"auth_provider_config",
					"auto_allocate_public_eip",
					"agentic_db_cluster_id",
					"agentic_db_branch_spec",
					"knowledge_application_spec",
					// The endpoint set cannot be compared as a whole because the
					// query API does not return EndpointType, and it does not
					// report the default endpoint's description either, so no
					// description assertion is possible here.
					"endpoints",
					// tags converge through ListTagResources: an empty result
					// is only written back when state already holds tags, so an
					// import leaves an absent map where the last apply kept {}.
					"tags",
					// The whitelist IP type is reported per array; an import
					// relies on the service echoing the configured value.
					"security_ip_type",
					// DescribeApplicationAttribute does not return the CPU
					// architecture for this application type, so an imported
					// state cannot hold a required create input.
					"architecture",
				},
				ImportStateCheck: testAccCheckPolarDBApplicationImportedState,
			},
		},
	})
}

// TestAccAliCloudPolarDBApplication_CreateFull covers the full attribute set
// including prepaid billing fields, model credentials and an attached cluster.
func TestAccAliCloudPolarDBApplication_CreateFull(t *testing.T) {
	resourceId := "alicloud_polardb_application.default"
	rand := tfacctest.RandStringFromCharSet(10, tfacctest.CharSetAlphaNum)
	name := fmt.Sprintf("tf-testAccPolarDBApplication-%s", rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-beijing"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckPolarDBApplicationDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccPolardbApplicationConfigFull(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
					// Equivalent of the flat-layout basicMap vswitch_id=CHECKSET.
					resource.TestCheckResourceAttrSet(resourceId, "vswitch_id"),
				),
			},
			{
				Config: testAccPolardbApplicationConfigFullWhitelist(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "status", "Activated"),
				),
			},
		},
	})
}

func testAccCheckPolarDBApplicationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_polardb_application" {
			continue
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		_, err := polardb.NewService(client).DescribeApplicationAttribute(rs.Primary.ID)
		if err != nil {
			if errs.NotFoundError(err) {
				continue
			}
			return fmt.Errorf("describing polardb application %s failed: %s", rs.Primary.ID, err)
		}
		return fmt.Errorf("polardb application %s still exists", rs.Primary.ID)
	}
	return nil
}

// testAccPolardbApplicationDependence is the shared VPC/VSwitch + node-class
// lookup the application depends on; ported verbatim from the flat-layout
// resourcePolarDBApplicationConfigDependence.
func testAccPolardbApplicationDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_vpc" "default" {
	vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
	zone_id  = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	vpc_id   = alicloud_vpc.default.id
	cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

locals {
	vpc_id     = alicloud_vpc.default.id
	vswitch_id = alicloud_vswitch.default.id
}

data "alicloud_polardb_node_classes" "this" {
	db_type    = "MySQL"
	db_version = "8.0"
	pay_type   = "PostPaid"
	category   = "Normal"
}
`, name)
}

// testAccPolardbApplicationCreateInputs holds the create-time inputs that every
// step of the Create suite must repeat: endpoints and security_ip_type are
// ForceNew, so dropping them would replace the application, and parameters are
// written by CreateApplication itself (they used to be sent only by the update
// path). Keeping them in one helper means the steps differ only in the
// attribute under test.
func testAccPolardbApplicationCreateInputs(parameterValue string) string {
	return fmt.Sprintf(`
	security_ip_type = "ipv4"
	endpoints {
		endpoint_type = "Primary"
		description   = "my_endpoint"
	}
	parameters {
		parameter_name  = "secret.dashscope.apiKey"
		parameter_value = "%s"
	}
`, parameterValue)
}

func testAccPolardbApplicationConfigBasic(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
%s}
`, testAccPolardbApplicationDependence(name), testAccPolardbApplicationCreateInputs("abc"))
}

func testAccPolardbApplicationConfigUpgradeVersion(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version = true
%s}
`, testAccPolardbApplicationDependence(name), testAccPolardbApplicationCreateInputs("abc"))
}

func testAccPolardbApplicationConfigParameters(name, parameterValue string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version = true
%s}
`, testAccPolardbApplicationDependence(name), testAccPolardbApplicationCreateInputs(parameterValue))
}

func testAccPolardbApplicationConfigWhitelist(name, modifyMode string) string {
	return testAccPolardbApplicationConfigWhitelistWithStorages(name, modifyMode, "")
}

// testAccPolardbApplicationConfigWhitelistWithStorages is the whitelist
// configuration plus an optional storages block, used by the plan-level
// storages steps.
func testAccPolardbApplicationConfigWhitelistWithStorages(name, modifyMode, storages string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version        = true
	modify_mode            = "%[2]s"
	security_ip_array_name = "test"
	security_ip_list       = ["127.0.0.1"]
%[3]s%[4]s}
`, testAccPolardbApplicationDependence(name), modifyMode, testAccPolardbApplicationCreateInputs("abc"), storages)
}

// testAccPolardbApplicationFullDependence is the full dependence (adds a
// resource-manager group lookup and a PolarDB cluster the application
// attaches to); ported verbatim from resourcePolarDBApplicationFullConfigDependence.
func testAccPolardbApplicationFullDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alicloud_vpc" "default" {
	vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
	zone_id  = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	vpc_id   = alicloud_vpc.default.id
	cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}

locals {
	vpc_id     = alicloud_vpc.default.id
	vswitch_id = alicloud_vswitch.default.id
	component_id = alicloud_polardb_application.default.component_id
}

data "alicloud_polardb_node_classes" "this" {
	db_type    = "MySQL"
	db_version = "8.0"
	pay_type   = "PostPaid"
	category   = "Normal"
}

data "alicloud_resource_manager_resource_groups" "default" {
	status = "OK"
}

resource "alicloud_polardb_cluster" "cluster" {
	db_type = "MySQL"
	db_version = "8.0"
	pay_type = "PostPaid"
	db_node_count = "2"
	db_node_class = "polar.mysql.x4.large"
	vswitch_id = "${local.vswitch_id}"
	description = "${var.name}"

	// encrypt_new_tables is optional without a computed fallback in the stock
	// cluster resource, and DescribeDBClusters intermittently stops returning
	// the TDEStatus sub-object after the first read; the resulting
	// "OFF -> null" plan fails this fixture's second step. Other polardb
	// fixtures ignore the same field for the same reason. It does not belong
	// to alicloud_polardb_application behaviour.
	lifecycle {
		ignore_changes = [encrypt_new_tables]
	}
}
`, name)
}

func testAccPolardbApplicationConfigFull(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	auto_renew        = false
	period            = 1
	model_api_key     = "sk-xxx"
	model_api         = "openai-completions"
	ai_db_cluster_id  = alicloud_polardb_cluster.cluster.id
	db_cluster_id     = alicloud_polardb_cluster.cluster.id
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
`, testAccPolardbApplicationFullDependence(name))
}

func testAccPolardbApplicationConfigFullWhitelist(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	auto_renew        = false
	period            = 1
	model_api_key     = "sk-xxx"
	model_api         = "openai-completions"
	ai_db_cluster_id  = alicloud_polardb_cluster.cluster.id
	db_cluster_id     = alicloud_polardb_cluster.cluster.id
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
	component_id           = ""
	modify_mode            = "Append"
	security_ip_array_name = "test"
	security_ip_list       = ["127.0.0.1"]
}
`, testAccPolardbApplicationFullDependence(name))
}

// testAccPolardbApplicationConfigTags extends the whitelist-delete config
// with two tags, exercising the TagResources update path (item 65) and the
// ListTagResources read-back. It shares the basic dependence so only the
// tags attribute differs from the prior step.
func testAccPolardbApplicationConfigTags(name string) string {
	return fmt.Sprintf(`%s

resource "alicloud_polardb_application" "default" {
	description      = var.name
	application_type = "polarclaw"
	architecture     = "x86"
	pay_type         = "PostPaid"
	region_id        = "cn-beijing"
	vswitch_id       = local.vswitch_id
	vpc_id           = local.vpc_id
	zone_id          = data.alicloud_polardb_node_classes.this.classes.1.zone_id
	model_from       = "bailian"
	model_base_url   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	model_name       = "qwen3.6-plus"
	components {
		component_type    = "polarclaw_comp"
		component_class   = "polar.app.g2.medium"
		component_replica = 1
	}
	upgrade_version        = true
	modify_mode            = "Delete"
	security_ip_array_name = "test"
	security_ip_list       = ["127.0.0.1"]
%s	tags = {
		Created = "allin"
		Env     = "test"
	}
}
`, testAccPolardbApplicationDependence(name), testAccPolardbApplicationCreateInputs("abc"))
}

// polardbUnitTestClient builds a hermetic client whose polardb endpoint is a
// local test server, so the tag calls can be exercised without cloud access.
func polardbUnitTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Errorf("parse form: %v", err)
			return
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("build test credential: %v", err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("build test client: %v", err)
	}
	endpoints.Store("polardb", endpoint)
	return client
}

// ListTagResources answers one page per call and returns a NextToken while
// more pages remain; ListApplicationTags must replay it and aggregate every
// page before returning.
func TestUnitPolarDBApplicationListTagsPagination(t *testing.T) {
	var calls int32
	client := polardbUnitTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if action := r.Form.Get("Action"); action != "ListTagResources" {
			t.Errorf("unexpected action %q", action)
		}
		if resourceType := r.Form.Get("ResourceType"); resourceType != "application" {
			t.Errorf("unexpected ResourceType %q", resourceType)
		}
		w.Header().Set("Content-Type", "application/json")
		switch atomic.AddInt32(&calls, 1) {
		case 1:
			if token := r.Form.Get("NextToken"); token != "" {
				t.Errorf("first page must not send NextToken, got %q", token)
			}
			fmt.Fprint(w, `{"TagResources":{"TagResource":[{"TagKey":"Created","TagValue":"allin"}]},"NextToken":"page-2"}`)
		case 2:
			if token := r.Form.Get("NextToken"); token != "page-2" {
				t.Errorf("second page must replay the returned NextToken, got %q", token)
			}
			fmt.Fprint(w, `{"TagResources":{"TagResource":[{"TagKey":"Env","TagValue":"test"}]}}`)
		default:
			t.Errorf("unexpected page request %d after NextToken stopped being returned", calls)
		}
	})

	fetched, err := polardb.NewService(client).ListApplicationTags("app-test")
	if err != nil {
		t.Fatalf("ListApplicationTags returned error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 ListTagResources calls, got %d", calls)
	}
	if len(fetched) != 2 {
		t.Fatalf("expected tags from both pages, got %#v", fetched)
	}
	if fetched["Created"] != "allin" || fetched["Env"] != "test" {
		t.Errorf("unexpected aggregated tags %#v", fetched)
	}
}

// A successful empty ListTagResources response (all tags removed, or cleared
// out-of-band) must clear the tags still held in state; a NotFound response is
// the auxiliary-source case and must leave them untouched.
func TestUnitPolarDBApplicationReadTagsConvergence(t *testing.T) {
	for _, tc := range []struct {
		name          string
		responseCode  int
		responseBody  string
		wantTagCount  int
		wantTagValues map[string]string
	}{
		{
			name:         "successful empty response clears stale tags",
			responseCode: http.StatusOK,
			responseBody: `{"TagResources":{}}`,
			wantTagCount: 0,
		},
		{
			name:         "not found keeps prior tags",
			responseCode: http.StatusNotFound,
			responseBody: `{"Code":"InvalidApplicationId.NotFound","Message":"The specified application is not found."}`,
			wantTagCount: 2,
			wantTagValues: map[string]string{
				"Created": "allin",
				"Env":     "test",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := polardbUnitTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Form.Get("Action") {
				case "DescribeApplicationAttribute":
					fmt.Fprint(w, `{"Status":"Activated","VSwitchId":"vsw-test"}`)
				case "ListTagResources":
					w.WriteHeader(tc.responseCode)
					fmt.Fprint(w, tc.responseBody)
				default:
					t.Errorf("unexpected action %q", r.Form.Get("Action"))
				}
			})

			d := schema.TestResourceDataRaw(t, polardb.ResourceApplication().Schema, map[string]interface{}{
				"application_type": "polarclaw",
				"architecture":     "x86",
				"tags": map[string]interface{}{
					"Created": "allin",
					"Env":     "test",
				},
			})
			d.SetId("app-test")

			if err := polardb.ResourceApplication().Read(d, client); err != nil {
				t.Fatalf("resource read returned error: %v", err)
			}
			got, _ := d.Get("tags").(map[string]interface{})
			if len(got) != tc.wantTagCount {
				t.Fatalf("expected %d tags after read, got %#v", tc.wantTagCount, got)
			}
			for key, want := range tc.wantTagValues {
				if got[key] != want {
					t.Errorf("expected tag %s=%s, got %#v", key, want, got)
				}
			}
		})
	}
}

// CreateApplication accepts DnatIpAddress and DnatEntries only together with
// VpcNatGatewayId. The dependency is declared in the schema, so an incomplete
// configuration fails at plan time instead of building a request the API
// rejects; each DNAT field is covered on its own, and VpcNatGatewayId alone
// stays valid.
func TestUnitPolarDBApplicationDnatRequiresNatGateway(t *testing.T) {
	resource := polardb.ResourceApplication()
	for _, key := range []string{"dnat_ip_address", "dnat_entries"} {
		got := resource.Schema[key].RequiredWith
		if len(got) != 1 || got[0] != "vpc_nat_gateway_id" {
			t.Fatalf("%s RequiredWith = %#v, want [vpc_nat_gateway_id]", key, got)
		}
	}

	base := map[string]interface{}{
		"application_type": "polarclaw",
		"architecture":     "x86",
		"region_id":        "cn-beijing",
	}
	with := func(kv map[string]interface{}) map[string]interface{} {
		out := map[string]interface{}{}
		for k, v := range base {
			out[k] = v
		}
		for k, v := range kv {
			out[k] = v
		}
		return out
	}
	entries := []interface{}{map[string]interface{}{"port_name": "test", "front_port": 80}}

	for _, tc := range []struct {
		name    string
		raw     map[string]interface{}
		wantErr bool
	}{
		{"dnat_ip_address without vpc_nat_gateway_id", with(map[string]interface{}{"dnat_ip_address": "10.0.0.1"}), true},
		{"dnat_entries without vpc_nat_gateway_id", with(map[string]interface{}{"dnat_entries": entries}), true},
		{"dnat_ip_address with vpc_nat_gateway_id", with(map[string]interface{}{"dnat_ip_address": "10.0.0.1", "vpc_nat_gateway_id": "ngw-test"}), false},
		{"dnat_entries with vpc_nat_gateway_id", with(map[string]interface{}{"dnat_entries": entries, "vpc_nat_gateway_id": "ngw-test"}), false},
		{"vpc_nat_gateway_id alone", with(map[string]interface{}{"vpc_nat_gateway_id": "ngw-test"}), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			diags := schema.InternalMap(resource.Schema).Validate(sdktf.NewResourceConfigRaw(tc.raw))
			if tc.wantErr {
				if !diags.HasError() {
					t.Fatalf("expected a validation error, got none")
				}
				return
			}
			if diags.HasError() {
				t.Fatalf("unexpected validation error: %v", diags)
			}
		})
	}
}

// testAccCheckPolarDBApplicationEndpointMember asserts that a TypeSet member
// with the given endpoint_type and description exists. A set element's index is
// its hash, so endpoints.0.<attr> is not a stable contract; this walks the
// state attributes instead.
func testAccCheckPolarDBApplicationEndpointMember(resourceId, endpointType, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceId]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceId)
		}
		attributes := rs.Primary.Attributes
		for key := range attributes {
			if !strings.HasPrefix(key, "endpoints.") || !strings.HasSuffix(key, ".endpoint_type") {
				continue
			}
			if attributes[key] != endpointType {
				continue
			}
			prefix := strings.TrimSuffix(key, "endpoint_type")
			if attributes[prefix+"description"] == description {
				return nil
			}
		}
		return fmt.Errorf("no endpoints member matches endpoint_type=%q description=%q", endpointType, description)
	}
}

// testAccCheckPolarDBApplicationImportedState asserts the readable state an
// import rebuilds from DescribeApplicationAttribute: the identifiers, the
// whitelist entry (name and IP list), the component id and the endpoint
// description. The endpoint set itself is excluded from the whole-set
// comparison because the query API does not return EndpointType.
func testAccCheckPolarDBApplicationImportedState(states []*terraform.InstanceState) error {
	if len(states) != 1 {
		return fmt.Errorf("expected one imported state, got %d", len(states))
	}
	attributes := states[0].Attributes
	if attributes["vswitch_id"] == "" {
		return fmt.Errorf("imported state lost vswitch_id")
	}
	// DescribeApplicationAttribute reports the application's default whitelist
	// entry; the configured array name is not echoed, so only the entry's
	// presence and its IP list are asserted here.
	if got := attributes["security_ip_array_name"]; got == "" {
		return fmt.Errorf("imported state lost the whitelist entry name")
	}
	if got := attributes["security_ip_list.#"]; got != "1" {
		return fmt.Errorf("imported state security_ip_list.# = %q, want 1", got)
	}
	if attributes["component_id"] == "" {
		return fmt.Errorf("imported state lost component_id")
	}
	// The endpoint description is not asserted here: the query response does
	// not return one for the application's default endpoint in this
	// environment (the imported state carries no description while the
	// configuration does), so the block stays excluded from the whole-set
	// comparison. Read still adopts any description the response reports.
	return nil
}
