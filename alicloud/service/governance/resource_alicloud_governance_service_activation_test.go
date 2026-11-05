package governance_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/governance"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	sdkterraform "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Cloud Governance Center is activated once per account in a single region, and its operations must
// target that region: GetGovernanceServiceStatus documents InvalidRegion.NotMatch for a request that
// names another one. The activation region comes from governanceServiceActivationTestRegionEnv and
// defaults to cn-hangzhou, the region used by the API documentation, the CLI examples and the shared
// test account.
//
// The provider block is pinned to that region because this resource requires region == provider
// region, while the resource attribute itself is resolved from alicloud_regions.current instead of
// repeating the literal. Set the environment variable to run the test against another activation
// region; the whole test then targets that region (the resource cannot move between regions).
const governanceServiceActivationTestRegionEnv = "ALICLOUD_GOVERNANCE_ACTIVATION_REGION"

func governanceServiceActivationTestRegion() string {
	if region := strings.TrimSpace(os.Getenv(governanceServiceActivationTestRegionEnv)); region != "" {
		return region
	}
	return "cn-hangzhou"
}

func TestAccAliCloudGovernanceServiceActivation_basic0(t *testing.T) {
	resourceId := "alicloud_governance_service_activation.default"
	region := governanceServiceActivationTestRegion()

	// CheckDestroy restores the baseline on the normal path; cleanup covers a run that fails
	// before the destroy step, which would otherwise leave the shared account deactivated.
	t.Cleanup(func() {
		client := governanceServiceActivationCapturedClient()
		if client == nil {
			return
		}
		if err := restoreGovernanceServiceActivationBaseline(client, region); err != nil {
			t.Errorf("restoring the pre-test governance activation baseline: %v", err)
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			captureGovernanceServiceActivationBaseline(t, region)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckGovernanceServiceActivationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudGovernanceServiceActivationBasic0(region),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("id"), knownvalue.StringExact(region)),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("region"), knownvalue.StringExact(region)),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("status"), knownvalue.StringExact("Enabled")),
				},
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// timeouts is the only mutable argument: changing it reaches Update, which must read the
			// live activation instead of persisting the Computed status as an unknown value.
			{
				Config: testAccAliCloudGovernanceServiceActivationBasic0WithCreateTimeout(region, "10m"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("id"), knownvalue.StringExact(region)),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("region"), knownvalue.StringExact(region)),
					statecheck.ExpectKnownValue(resourceId, tfjsonpath.New("status"), knownvalue.StringExact("Enabled")),
				},
			},
		},
	})
}

// testAccCheckGovernanceServiceActivationDestroy verifies the account-level service was
// deactivated: GetGovernanceServiceStatus keeps answering after decommission, so the
// resource being gone is expressed as a Disabled status rather than a missing resource.
func testAccCheckGovernanceServiceActivationDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_governance_service_activation" {
			continue
		}
		region := rs.Primary.Attributes["region"]
		if region == "" {
			region = rs.Primary.ID
		}
		client := acctest.Provider.Meta().(*connectivity.AliyunClient)
		// The service sends RegionId in the request body for this API family; mirror that call shape
		// with the region of the resource instead of the harness provider region.
		body := map[string]interface{}{"RegionId": region}
		response, err := client.RpcPost("governance", "2021-01-20", "GetGovernanceServiceStatus", nil, body, true)
		if err != nil {
			return fmt.Errorf("checking destroy of %s %s: %w", rs.Type, rs.Primary.ID, err)
		}
		if status, _ := response["ServiceStatus"].(string); status != "Disabled" {
			return fmt.Errorf("the resource %s %s was not destroyed: governance service status is %q", rs.Type, rs.Primary.ID, status)
		}
		if err := restoreGovernanceServiceActivationBaseline(client, region); err != nil {
			return err
		}
	}
	return nil
}

// The activation is account-wide, so the test records the pre-test status and puts an already
// activated service back after the destroy step: a green run must not change the shared account.
var (
	governanceServiceActivationBaselineMu sync.Mutex
	governanceServiceActivationBaseline   string
	governanceServiceActivationClient     *connectivity.AliyunClient
)

func governanceServiceActivationCapturedClient() *connectivity.AliyunClient {
	governanceServiceActivationBaselineMu.Lock()
	defer governanceServiceActivationBaselineMu.Unlock()
	return governanceServiceActivationClient
}

func governanceServiceActivationStatus(client *connectivity.AliyunClient, region string) (string, error) {
	body := map[string]interface{}{"RegionId": region}
	response, err := client.RpcPost("governance", "2021-01-20", "GetGovernanceServiceStatus", nil, body, true)
	if err != nil {
		return "", err
	}
	status, _ := response["ServiceStatus"].(string)
	return status, nil
}

func captureGovernanceServiceActivationBaseline(t *testing.T, region string) {
	t.Helper()
	// PreCheck runs before any Terraform command, so the shared provider has not been configured
	// yet and Meta() is nil; configure it from the environment once so the baseline can be read.
	if acctest.Provider.Meta() == nil {
		// Configure the shared provider for the region under test, not for whatever the environment
		// default is: the baseline read and the cleanup restore must address the same activation the
		// test manages.
		raw := sdkterraform.NewResourceConfigRaw(map[string]interface{}{"region": region})
		if diagnostics := acctest.Provider.Configure(context.Background(), raw); diagnostics.HasError() {
			t.Fatalf("configuring the provider for the pre-test baseline read: %v", diagnostics)
		}
	}
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	status, err := governanceServiceActivationStatus(client, region)
	if err != nil {
		t.Fatalf("reading the pre-test governance service status in %s: %v", region, err)
	}
	governanceServiceActivationBaselineMu.Lock()
	governanceServiceActivationBaseline = status
	governanceServiceActivationClient = client
	governanceServiceActivationBaselineMu.Unlock()
	t.Logf("pre-test governance service status in %s: %s", region, status)
}

// governanceServiceActivationNeedsRestore decides whether the pre-test status must be restored to
// Enabled. A Disabled or EnableFailed baseline already ends the test disabled; a Disabling baseline
// settles to Disabled or DisableFailed, and a DisableFailed service cannot be opened again (the
// published Open table rejects it with InvalidStatus.NotAllowed), so those baselines are left as-is
// and only the on / transitionally-on ones are restored.
func governanceServiceActivationNeedsRestore(baseline string) bool {
	switch baseline {
	case "Enabled", "Enabling":
		return true
	}
	return false
}

func restoreGovernanceServiceActivationBaseline(client *connectivity.AliyunClient, region string) error {
	governanceServiceActivationBaselineMu.Lock()
	baseline := governanceServiceActivationBaseline
	governanceServiceActivationBaselineMu.Unlock()
	if !governanceServiceActivationNeedsRestore(baseline) {
		return nil
	}
	// An in-flight activation change has to settle before the restore decision: opening while a
	// deactivation is still running would race the same singleton switch.
	settleDeadline := time.Now().Add(10 * time.Minute)
	status, err := governanceServiceActivationStatus(client, region)
	if err != nil {
		return fmt.Errorf("restoring the pre-test activation in %s: %w", region, err)
	}
	for (status == "Enabling" || status == "Disabling") && time.Now().Before(settleDeadline) {
		time.Sleep(10 * time.Second)
		status, err = governanceServiceActivationStatus(client, region)
		if err != nil {
			return fmt.Errorf("restoring the pre-test activation in %s: %w", region, err)
		}
	}
	if status == "Enabling" || status == "Disabling" {
		return fmt.Errorf("restoring the pre-test activation in %s: the activation change did not settle within 10m (last status %q)", region, status)
	}
	// Idempotent: the destroy check and the test cleanup both restore, and a failure that never
	// reached the destroy step leaves the service enabled.
	if status == "Enabled" {
		return nil
	}
	// The restore budget starts after the settle wait, so a slow activation change cannot starve
	// the enable poll.
	deadline := time.Now().Add(20 * time.Minute)
	request := map[string]interface{}{
		"ClientToken": helper.BuildClientToken("OpenGovernanceService"),
		"RegionId":    region,
	}
	if _, err := client.RpcPost("governance", "2021-01-20", "OpenGovernanceService", nil, request, true); err != nil {
		return fmt.Errorf("restoring the pre-test activation in %s: %w", region, err)
	}
	for time.Now().Before(deadline) {
		status, err := governanceServiceActivationStatus(client, region)
		if err != nil {
			return fmt.Errorf("restoring the pre-test activation in %s: %w", region, err)
		}
		if status == "Enabled" {
			return nil
		}
		time.Sleep(10 * time.Second)
	}
	return fmt.Errorf("restoring the pre-test activation in %s: status did not reach Enabled", region)
}

func testAccAliCloudGovernanceServiceActivationBasic0(region string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = %[1]q
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_governance_service_activation" "default" {
  region = data.alicloud_regions.default.regions.0.id
}
`, region)
}

func testAccAliCloudGovernanceServiceActivationBasic0WithCreateTimeout(region, createTimeout string) string {
	return fmt.Sprintf(`
provider "alicloud" {
  region = %[1]q
}

data "alicloud_regions" "default" {
  current = true
}

resource "alicloud_governance_service_activation" "default" {
  region = data.alicloud_regions.default.regions.0.id

  timeouts {
    create = %[2]q
  }
}
`, region, createTimeout)
}

func TestUnitAliCloudGovernanceServiceActivationSchema(t *testing.T) {
	r := governance.NewServiceActivationResource()

	var metadata fwresource.MetadataResponse
	r.Metadata(context.Background(), fwresource.MetadataRequest{ProviderTypeName: "alicloud"}, &metadata)
	if metadata.TypeName != "alicloud_governance_service_activation" {
		t.Fatalf("TypeName = %q, want alicloud_governance_service_activation", metadata.TypeName)
	}

	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}

	region, ok := schemaResp.Schema.Attributes["region"].(schema.StringAttribute)
	if !ok {
		t.Fatal("region attribute is missing or is not a string attribute")
	}
	// Activation is account-and-region scoped and the region is the resource ID, so a
	// region change must replace the resource rather than update it in place.
	if !region.Required || len(region.PlanModifiers) != 1 {
		t.Fatalf("region must stay Required with one plan modifier, got Required=%v modifiers=%d", region.Required, len(region.PlanModifiers))
	}
	// Drive the modifier the way Terraform does when state and plan disagree: it must request
	// replacement, otherwise a region change would be planned as an in-place update.
	objectType, ok := schemaResp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("resource schema does not describe an object")
	}
	object := func(region, status string) tftypes.Value {
		return tftypes.NewValue(objectType, map[string]tftypes.Value{
			"id":       tftypes.NewValue(tftypes.String, region),
			"region":   tftypes.NewValue(tftypes.String, region),
			"status":   tftypes.NewValue(tftypes.String, status),
			"timeouts": tftypes.NewValue(objectType.AttributeTypes["timeouts"], nil),
		})
	}
	var modifierResp planmodifier.StringResponse
	region.PlanModifiers[0].PlanModifyString(context.Background(), planmodifier.StringRequest{
		State:      tfsdk.State{Schema: schemaResp.Schema, Raw: object("cn-hangzhou", "Enabled")},
		Plan:       tfsdk.Plan{Schema: schemaResp.Schema, Raw: object("cn-beijing", "Enabled")},
		StateValue: types.StringValue("cn-hangzhou"),
		PlanValue:  types.StringValue("cn-beijing"),
	}, &modifierResp)
	if !modifierResp.RequiresReplace {
		t.Fatal("region plan modifier did not request replacement for a changed region")
	}

	status, ok := schemaResp.Schema.Attributes["status"].(schema.StringAttribute)
	if !ok {
		t.Fatal("status attribute is missing or is not a string attribute")
	}
	if !status.Computed || status.Required || status.Optional {
		t.Fatalf("status must stay Computed-only, got Computed=%v Required=%v Optional=%v", status.Computed, status.Required, status.Optional)
	}
}

func TestUnitAliCloudGovernanceServiceActivationModifyPlanRegion(t *testing.T) {
	r := governance.NewServiceActivationResource()

	configured, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatal("resource does not accept a provider client")
	}
	client := &connectivity.AliyunClient{RegionId: "cn-hangzhou"}
	var configureResp fwresource.ConfigureResponse
	configured.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}

	modifyPlan, ok := r.(fwresource.ResourceWithModifyPlan)
	if !ok {
		t.Fatal("resource does not implement ModifyPlan")
	}

	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objectType, ok := schemaResp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("resource schema does not describe an object")
	}

	plan := func(region string) tfsdk.Plan {
		return tfsdk.Plan{
			Schema: schemaResp.Schema,
			Raw: tftypes.NewValue(objectType, map[string]tftypes.Value{
				"id":       tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"region":   tftypes.NewValue(tftypes.String, region),
				"status":   tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"timeouts": tftypes.NewValue(objectType.AttributeTypes["timeouts"], nil),
			}),
		}
	}

	// All three operations address the activation through the provider region, so a plan that names
	// another region must be rejected instead of describing one region through another's client.
	var rejected fwresource.ModifyPlanResponse
	modifyPlan.ModifyPlan(context.Background(), fwresource.ModifyPlanRequest{Plan: plan("cn-beijing")}, &rejected)
	if !rejected.Diagnostics.HasError() {
		t.Fatal("ModifyPlan accepted a region that differs from the provider region")
	}
	if detail := rejected.Diagnostics[0].Detail(); !strings.Contains(detail, "region must be") {
		t.Fatalf("unexpected plan diagnostic: %s", detail)
	}

	var accepted fwresource.ModifyPlanResponse
	modifyPlan.ModifyPlan(context.Background(), fwresource.ModifyPlanRequest{Plan: plan("cn-hangzhou")}, &accepted)
	if accepted.Diagnostics.HasError() {
		t.Fatalf("ModifyPlan rejected the provider region: %v", accepted.Diagnostics)
	}
}

// The restore decision is the test's own cleanup policy: every status except an on/transitionally-on
// baseline ends the test in the same disabled state the test itself leaves behind.
func TestUnitGovernanceServiceActivationRestorePolicy(t *testing.T) {
	policy := map[string]bool{
		"Enabled":       true,
		"Enabling":      true,
		"DisableFailed": false,
		"Disabled":      false,
		"Disabling":     false,
		"EnableFailed":  false,
		"":              false,
	}
	for baseline, want := range policy {
		if got := governanceServiceActivationNeedsRestore(baseline); got != want {
			t.Errorf("needsRestore(%q) = %v, want %v", baseline, got, want)
		}
	}
}

// TestUnitGovernanceServiceActivationDeleteWaitsForInFlightChange drives Delete against a hermetic
// endpoint: the first deactivation is rejected with the published InvalidStatus.OperationInProgress
// while a change is still running, and delete must observe the in-flight status, wait for it to
// settle and reach Disabled instead of failing the first attempt.
func TestUnitGovernanceServiceActivationDeleteWaitsForInFlightChange(t *testing.T) {
	var mu sync.Mutex
	statusCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		switch req.Form.Get("Action") {
		case "DecommissionGovernance":
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"RequestId":"mock","Code":"InvalidStatus.OperationInProgress","Message":"mock operation in progress"}`))
		case "GetGovernanceServiceStatus":
			mu.Lock()
			statusCalls++
			current := statusCalls
			mu.Unlock()
			if current == 1 {
				_, _ = w.Write([]byte(`{"RequestId":"mock","ServiceStatus":"Disabling"}`))
				return
			}
			_, _ = w.Write([]byte(`{"RequestId":"mock","ServiceStatus":"Disabled"}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"RequestId":"mock","Code":"InvalidParameter","Message":"unexpected action"}`))
		}
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("governance", strings.TrimPrefix(server.URL, "http://"))

	r := governance.NewServiceActivationResource()
	configured, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatal("resource does not accept a provider client")
	}
	var configureResp fwresource.ConfigureResponse
	configured.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}

	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objectType, ok := schemaResp.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("resource schema does not describe an object")
	}
	raw := tftypes.NewValue(objectType, map[string]tftypes.Value{
		"id":       tftypes.NewValue(tftypes.String, "cn-hangzhou"),
		"region":   tftypes.NewValue(tftypes.String, "cn-hangzhou"),
		"status":   tftypes.NewValue(tftypes.String, "Enabled"),
		"timeouts": tftypes.NewValue(objectType.AttributeTypes["timeouts"], nil),
	})
	state := tfsdk.State{Schema: schemaResp.Schema, Raw: raw}
	resp := &fwresource.DeleteResponse{State: state}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: state}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("delete failed instead of waiting for the in-flight change: %v", resp.Diagnostics)
	}
	mu.Lock()
	calls := statusCalls
	mu.Unlock()
	if calls < 2 {
		t.Fatalf("expected the delete to poll the status at least twice, got %d", calls)
	}
}

// governanceUnitClient points the governance endpoint at a hermetic server so the resource
// lifecycle can be driven without cloud credentials or a real activation.
func governanceUnitClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	t.Setenv("NO_PROXY", "127.0.0.1,localhost")
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("governance", strings.TrimPrefix(server.URL, "http://"))
	return client
}

func governanceUnitResource(t *testing.T, client *connectivity.AliyunClient) (fwresource.Resource, schema.Schema) {
	t.Helper()
	r := governance.NewServiceActivationResource()
	configured, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatal("resource does not accept a provider client")
	}
	var configureResp fwresource.ConfigureResponse
	configured.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}
	var schemaResp fwresource.SchemaResponse
	r.Schema(context.Background(), fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	return r, schemaResp.Schema
}

func governanceUnitObjectType(t *testing.T, resourceSchema schema.Schema) tftypes.Object {
	t.Helper()
	objectType, ok := resourceSchema.Type().TerraformType(context.Background()).(tftypes.Object)
	if !ok {
		t.Fatal("resource schema does not describe an object")
	}
	return objectType
}

func governanceUnitTimeoutValue(objectType tftypes.Object, create, deleteTimeout string) tftypes.Value {
	optional := func(value string) tftypes.Value {
		if value == "" {
			return tftypes.NewValue(tftypes.String, nil)
		}
		return tftypes.NewValue(tftypes.String, value)
	}
	return tftypes.NewValue(objectType.AttributeTypes["timeouts"], map[string]tftypes.Value{
		"create": optional(create),
		"delete": optional(deleteTimeout),
	})
}

func governanceUnitTimeouts(t *testing.T, resourceSchema schema.Schema, create, deleteTimeout string) tftypes.Value {
	t.Helper()
	return governanceUnitTimeoutValue(governanceUnitObjectType(t, resourceSchema), create, deleteTimeout)
}

func governanceUnitState(t *testing.T, resourceSchema schema.Schema, status string, timeouts tftypes.Value) tfsdk.State {
	t.Helper()
	objectType := governanceUnitObjectType(t, resourceSchema)
	return tfsdk.State{
		Schema: resourceSchema,
		Raw: tftypes.NewValue(objectType, map[string]tftypes.Value{
			"id":       tftypes.NewValue(tftypes.String, "cn-hangzhou"),
			"region":   tftypes.NewValue(tftypes.String, "cn-hangzhou"),
			"status":   tftypes.NewValue(tftypes.String, status),
			"timeouts": timeouts,
		}),
	}
}

func governanceUnitPlan(t *testing.T, resourceSchema schema.Schema, timeouts tftypes.Value) tfsdk.Plan {
	t.Helper()
	objectType := governanceUnitObjectType(t, resourceSchema)
	return tfsdk.Plan{
		Schema: resourceSchema,
		Raw: tftypes.NewValue(objectType, map[string]tftypes.Value{
			"id":       tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"region":   tftypes.NewValue(tftypes.String, "cn-hangzhou"),
			"status":   tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			"timeouts": timeouts,
		}),
	}
}

func governanceUnitBadRequest(w http.ResponseWriter, code, message string) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"RequestId":"mock","Code":%q,"Message":%q}`, code, message)))
}

// A rolled-back activation (EnableFailed) has no activation left, so Read must drop it from state
// instead of persisting a status a repeated destroy cannot clear.
func TestUnitGovernanceServiceActivationReadEnableFailedRemovesState(t *testing.T) {
	client := governanceUnitClient(t, func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		if action := req.Form.Get("Action"); action != "GetGovernanceServiceStatus" {
			governanceUnitBadRequest(w, "InvalidParameter", "unexpected action "+action)
			return
		}
		_, _ = w.Write([]byte(`{"RequestId":"mock","ServiceStatus":"EnableFailed","FailureCode":"mock","FailureMessage":"rolled back"}`))
	})
	r, resourceSchema := governanceUnitResource(t, client)

	state := governanceUnitState(t, resourceSchema, "Enabling", governanceUnitTimeouts(t, resourceSchema, "", ""))
	readResp := &fwresource.ReadResponse{State: state}
	r.Read(context.Background(), fwresource.ReadRequest{State: state}, readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("reading an EnableFailed activation: %v", readResp.Diagnostics)
	}
	if !readResp.State.Raw.IsNull() {
		t.Fatal("Read kept a fully rolled-back (EnableFailed) activation in state")
	}
}

// Enabling -> EnableFailed -> destroy/reapply: DecommissionGovernance answers InvalidStatus.NotEnabled
// for a rolled-back activation, so Delete has to observe the in-flight Enabling, wait for the rollback
// to finish and treat EnableFailed as a completed cleanup. The reapply must then open the service
// again instead of failing on the previous EnableFailed status.
func TestUnitGovernanceServiceActivationDeleteEnableFailedThenReapply(t *testing.T) {
	var mu sync.Mutex
	status := "Enabling"
	decommissions := 0
	statusCalls := 0
	opens := 0
	client := governanceUnitClient(t, func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		defer mu.Unlock()
		switch action := req.Form.Get("Action"); action {
		case "DecommissionGovernance":
			decommissions++
			governanceUnitBadRequest(w, "InvalidStatus.NotEnabled", "mock not enabled")
		case "OpenGovernanceService":
			opens++
			status = "Enabled"
			_, _ = w.Write([]byte(`{"RequestId":"mock"}`))
		case "GetGovernanceServiceStatus":
			statusCalls++
			if statusCalls > 1 && decommissions > 0 && status != "Enabled" {
				status = "EnableFailed"
			}
			_, _ = w.Write([]byte(fmt.Sprintf(`{"RequestId":"mock","ServiceStatus":%q}`, status)))
		default:
			governanceUnitBadRequest(w, "InvalidParameter", "unexpected action "+action)
		}
	})
	r, resourceSchema := governanceUnitResource(t, client)

	// A short delete budget keeps the regression verdict fast: with the pre-fix behaviour the wait
	// only accepted Disabled, so this test would run into the deadline instead of passing.
	state := governanceUnitState(t, resourceSchema, "Enabling", governanceUnitTimeouts(t, resourceSchema, "", "1m"))
	deleteResp := &fwresource.DeleteResponse{State: state}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: state}, deleteResp)
	if deleteResp.Diagnostics.HasError() {
		t.Fatalf("Delete failed to clear a rolled-back activation: %v", deleteResp.Diagnostics)
	}
	mu.Lock()
	decommissionCalls, observed := decommissions, statusCalls
	mu.Unlock()
	if decommissionCalls != 1 {
		t.Fatalf("expected a single DecommissionGovernance call, got %d", decommissionCalls)
	}
	if observed < 2 {
		t.Fatalf("expected Delete to observe Enabling before EnableFailed, got %d status calls", observed)
	}

	plan := governanceUnitPlan(t, resourceSchema, governanceUnitTimeouts(t, resourceSchema, "", ""))
	createResp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema}}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: plan}, createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("reapplying after EnableFailed: %v", createResp.Diagnostics)
	}
	mu.Lock()
	openCalls := opens
	mu.Unlock()
	if openCalls != 1 {
		t.Fatalf("expected the reapply to reopen the service once, got %d OpenGovernanceService calls", openCalls)
	}
	var applied serviceActivationAppliedState
	createResp.Diagnostics.Append(createResp.State.Get(context.Background(), &applied)...)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("reading the reapplied state: %v", createResp.Diagnostics)
	}
	if applied.Status.ValueString() != "Enabled" {
		t.Fatalf("reapplied status = %q, want Enabled", applied.Status.ValueString())
	}
}

// serviceActivationAppliedState is the resource model as the tests read it back from state.
type serviceActivationAppliedState struct {
	Id         types.String `tfsdk:"id"`
	FwTimeouts types.Object `tfsdk:"timeouts"`
	Region     types.String `tfsdk:"region"`
	Status     types.String `tfsdk:"status"`
}

// One delete budget covers the write and every wait: a delayed first write must not leave a second
// full budget for a re-issued DecommissionGovernance.
func TestUnitGovernanceServiceActivationDeleteUsesOneBudget(t *testing.T) {
	var mu sync.Mutex
	decommissions := 0
	statusCalls := 0
	client := governanceUnitClient(t, func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		defer mu.Unlock()
		switch action := req.Form.Get("Action"); action {
		case "DecommissionGovernance":
			decommissions++
			first := decommissions == 1
			if first {
				mu.Unlock()
				time.Sleep(time.Second)
				mu.Lock()
			}
			governanceUnitBadRequest(w, "InvalidStatus.OperationInProgress", "mock in flight")
		case "GetGovernanceServiceStatus":
			statusCalls++
			current := "DisableFailed"
			if statusCalls == 1 {
				current = "Enabling"
			}
			_, _ = w.Write([]byte(fmt.Sprintf(`{"RequestId":"mock","ServiceStatus":%q}`, current)))
		default:
			governanceUnitBadRequest(w, "InvalidParameter", "unexpected action "+action)
		}
	})
	r, resourceSchema := governanceUnitResource(t, client)

	state := governanceUnitState(t, resourceSchema, "Enabled", governanceUnitTimeouts(t, resourceSchema, "", "2s"))
	start := time.Now()
	deleteResp := &fwresource.DeleteResponse{State: state}
	r.Delete(context.Background(), fwresource.DeleteRequest{State: state}, deleteResp)
	elapsed := time.Since(start)

	if !deleteResp.Diagnostics.HasError() {
		t.Fatal("Delete reported success while the activation never reached Disabled")
	}
	mu.Lock()
	decommissionCalls := decommissions
	mu.Unlock()
	if decommissionCalls != 1 {
		t.Fatalf("Delete started %d DecommissionGovernance writes for a 2s budget; a re-issued write must not get a fresh budget", decommissionCalls)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("Delete spent %s for a 2s budget, so the write and the wait did not share one budget", elapsed)
	}
}

// Create has the same single-budget rule: the open write must not be retried or waited on with a
// fresh budget once the operation deadline is gone.
func TestUnitGovernanceServiceActivationCreateUsesOneBudget(t *testing.T) {
	var mu sync.Mutex
	opens := 0
	statusCalls := 0
	client := governanceUnitClient(t, func(w http.ResponseWriter, req *http.Request) {
		_ = req.ParseForm()
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		switch action := req.Form.Get("Action"); action {
		case "OpenGovernanceService":
			opens++
			mu.Unlock()
			// The write outlasts the 2s budget, so the retry backoff can never legally start again.
			time.Sleep(2500 * time.Millisecond)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"RequestId":"mock","Code":"InternalError","Message":"mock internal error"}`))
			return
		case "GetGovernanceServiceStatus":
			statusCalls++
			_, _ = w.Write([]byte(`{"RequestId":"mock","ServiceStatus":"Disabling"}`))
		default:
			governanceUnitBadRequest(w, "InvalidParameter", "unexpected action "+action)
		}
		mu.Unlock()
	})
	r, resourceSchema := governanceUnitResource(t, client)

	plan := governanceUnitPlan(t, resourceSchema, governanceUnitTimeouts(t, resourceSchema, "2s", ""))
	createResp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema}}
	r.Create(context.Background(), fwresource.CreateRequest{Plan: plan}, createResp)

	if !createResp.Diagnostics.HasError() {
		t.Fatal("Create reported success while the service never reached Enabled")
	}
	mu.Lock()
	openCalls, statusQueries := opens, statusCalls
	mu.Unlock()
	if openCalls != 1 {
		t.Fatalf("Create issued %d OpenGovernanceService writes for a 2s budget; no write may start after the deadline", openCalls)
	}
	// The only status query is the pre-open read: once the write consumed the budget, the
	// follow-up query must not be issued with a fresh budget either.
	if statusQueries != 1 {
		t.Fatalf("Create issued %d GetGovernanceServiceStatus queries, want only the pre-open read", statusQueries)
	}
}
