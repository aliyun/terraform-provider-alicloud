package paidsw_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/paidsw"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// PAI DSW snapshots are created against an existing DSW instance and reference an
// image URL; the provider has no alicloud_pai_dsw_instance resource to build one
// inline (and creating one is a paid compute resource outside this task's
// authorization). The two values are therefore borrowed prerequisites supplied
// via environment, the same pattern other heavy-dependency ACC tests in this
// repo use. When they are absent the lifecycle test skips — this is a genuine
// account-access dependency reported as a verification gap, not a fake pass.
const (
	envPaiDswInstanceId = "PAI_DSW_INSTANCE_ID"
	envPaiDswImageUrl   = "PAI_DSW_IMAGE_URL"
)

func testAccPaiDswApiInstanceSnapshotPreCheck(t *testing.T) {
	acctest.PreCheck(t)
	if v := os.Getenv(envPaiDswInstanceId); v == "" {
		t.Skipf("%s must reference an existing PAI DSW instance for this acceptance test", envPaiDswInstanceId)
	}
	if v := os.Getenv(envPaiDswImageUrl); v == "" {
		t.Skipf("%s must reference a valid image URL for this acceptance test", envPaiDswImageUrl)
	}
}

// TestAccAliCloudPaiDswApiInstanceSnapshot_basic covers the full create/read/
// no-op/update-via-replace/import lifecycle, including the optional labels and
// exclude_paths collections and the create-only snapshot_description/overwrite
// fields. The snapshot API has no Update, so every input change forces
// replacement.
func TestAccAliCloudPaiDswApiInstanceSnapshot_basic(t *testing.T) {
	r := tfacctest.RandInt()
	testAccPaiDswApiInstanceSnapshotPreCheck(t)

	instanceId := os.Getenv(envPaiDswInstanceId)
	imageUrl := os.Getenv(envPaiDswImageUrl)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiDswApiInstanceSnapshotConfigBasic(instanceId, imageUrl, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_api_instance_snapshot.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_dsw_api_instance_snapshot.default", "status"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "instance_id", instanceId),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "image_url", imageUrl),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_name", fmt.Sprintf("tf-testacc-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_description", fmt.Sprintf("tf-testacc description %d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "overwrite", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "labels.#", "2"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "exclude_paths.#", "2"),
				),
			},
			{
				// No-op: re-applying the same config must produce an empty plan.
				Config: testAccPaiDswApiInstanceSnapshotConfigBasic(instanceId, imageUrl, r),
			},
			{
				// Update via replacement: snapshot_name, description and the
				// labels/exclude_paths collections all change; overwrite is
				// omitted to cover the unset branch.
				Config: testAccPaiDswApiInstanceSnapshotConfigUpdate(instanceId, imageUrl, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_name", fmt.Sprintf("tf-testacc-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "snapshot_description", fmt.Sprintf("tf-testacc description %d upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "labels.#", "1"),
					resource.TestCheckResourceAttr("alicloud_pai_dsw_api_instance_snapshot.default", "exclude_paths.#", "1"),
				),
			},
			{
				ResourceName:            "alicloud_pai_dsw_api_instance_snapshot.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"snapshot_description", "overwrite"},
			},
		},
	})
}

// TestAccAliCloudPaiDswApiInstanceSnapshot_invalidName covers the failure path
// where the snapshot_name violates the documented charset (lowercase letters,
// digits, hyphens). It is a schema-level validation that fires during plan
// with no API call, so it only needs credentials, not a borrowed instance.
func TestAccAliCloudPaiDswApiInstanceSnapshot_invalidName(t *testing.T) {
	acctest.PreCheck(t)
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiDswApiInstanceSnapshotConfigInvalidName(r),
				// The regex uses \s+ between tokens because the Terraform CLI
				// line-wraps the diagnostic detail (~80 cols), inserting a
				// newline + indentation mid-message. A literal-spaces pattern
				// would fail to match the wrapped text.
				ExpectError: regexp.MustCompile(`must\s+only\s+contain\s+lowercase\s+letters,\s+digits\s+and\s+hyphens`),
			},
		},
	})
}

// TestAccAliCloudPaiDswApiInstanceSnapshot_invalidInstance covers the cloud
// failure path: Create against a non-existent instance must surface an error
// rather than succeed silently. It needs credentials only — the bogus instance
// and dummy image make the Create call fail before any real resource is built.
func TestAccAliCloudPaiDswApiInstanceSnapshot_invalidInstance(t *testing.T) {
	acctest.PreCheck(t)
	r := tfacctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccPaiDswApiInstanceSnapshotConfigInvalidInstance(r),
				ExpectError: regexp.MustCompile(`(?i)(not\s*found|not\s*exist|invalid|forbidden|endpoint|404|400|instance)`),
			},
		},
	})
}

func testAccPaiDswApiInstanceSnapshotConfigBasic(instanceId, imageUrl string, r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id          = "%[1]s"
  image_url            = "%[2]s"
  snapshot_name        = "tf-testacc-%[3]d"
  snapshot_description = "tf-testacc description %[3]d"
  overwrite            = true

  labels = [
    { key = "env",    value = "tf-testacc" },
    { key = "purpose", value = "snapshot" },
  ]

  exclude_paths = ["/home/admin/data", "/tmp/cache"]
}
`, instanceId, imageUrl, r)
}

func testAccPaiDswApiInstanceSnapshotConfigUpdate(instanceId, imageUrl string, r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id          = "%[1]s"
  image_url            = "%[2]s"
  snapshot_name        = "tf-testacc-%[3]d-upd"
  snapshot_description = "tf-testacc description %[3]d upd"

  labels = [
    { key = "env", value = "tf-testacc-upd" },
  ]

  exclude_paths = ["/home/admin/workspace"]
}
`, instanceId, imageUrl, r)
}

func testAccPaiDswApiInstanceSnapshotConfigInvalidName(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id   = "i-tf-testacc-%[1]d"
  image_url     = "registry.cn-hangzhou.aliyuncs.com/pai-dsw/nonexistent:latest"
  snapshot_name = "TF-Testacc-UPPER-%[1]d"
}
`, r)
}

func testAccPaiDswApiInstanceSnapshotConfigInvalidInstance(r int) string {
	return fmt.Sprintf(`
resource "alicloud_pai_dsw_api_instance_snapshot" "default" {
  instance_id   = "i-tf-testacc-nonexistent-%[1]d"
  image_url     = "registry.cn-hangzhou.aliyuncs.com/pai-dsw/nonexistent:latest"
  snapshot_name = "tf-testacc-invalidinst-%[1]d"
}
`, r)
}

// The unit tests below exercise the Framework CRUD methods against an in-process
// HTTP server. Connectivity resolves the pai-dsw endpoint to that server, so no
// request leaves the test process and no cloud credentials or paid resources are
// involved. They are the cloudless fallback for the acceptance lifecycle that
// needs an existing DSW instance and image URL.
const (
	snapshotUnitInstanceId = "i-unit123"
	snapshotUnitSnapshotId = "s-unit456"
	snapshotUnitName       = "tf-testacc-unit"
	snapshotUnitImageUrl   = "registry.example.com/pai-dsw/image:latest"
)

// snapshotUnitResourceAPI is the exported surface the resource publishes through
// the framework: CRUD plus the Configure and ImportState extensions. The tests
// deliberately go through this interface (and the factory) instead of the
// unexported implementation so they live in the external test package, where
// they can share a file with the acceptance tests.
type snapshotUnitResourceAPI interface {
	fwresource.Resource
	fwresource.ResourceWithConfigure
	fwresource.ResourceWithImportState
}

// snapshotUnitClientCarrier exposes the injected connectivity client, the same
// exported accessor the provider asserts on at registration.
type snapshotUnitClientCarrier interface {
	Client() *connectivity.AliyunClient
}

// snapshotUnitLabel mirrors one element of the labels list attribute.
type snapshotUnitLabel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

// snapshotUnitLabelType must match the labels ElementType declared by the
// resource schema; it is used to build typed lists without importing the
// unexported type.
var snapshotUnitLabelType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	},
}

// snapshotUnitStateModel mirrors the resource schema attribute names so plans
// and states can be built and decoded in this external package.
type snapshotUnitStateModel struct {
	Id                  types.String `tfsdk:"id"`
	InstanceId          types.String `tfsdk:"instance_id"`
	SnapshotName        types.String `tfsdk:"snapshot_name"`
	SnapshotDescription types.String `tfsdk:"snapshot_description"`
	ImageUrl            types.String `tfsdk:"image_url"`
	Overwrite           types.Bool   `tfsdk:"overwrite"`
	Labels              types.List   `tfsdk:"labels"`
	ExcludePaths        types.List   `tfsdk:"exclude_paths"`
	SnapshotId          types.String `tfsdk:"snapshot_id"`
	Status              types.String `tfsdk:"status"`
	ImageId             types.String `tfsdk:"image_id"`
	ReasonCode          types.String `tfsdk:"reason_code"`
	ReasonMessage       types.String `tfsdk:"reason_message"`
	GmtCreateTime       types.String `tfsdk:"gmt_create_time"`
	GmtModifiedTime     types.String `tfsdk:"gmt_modified_time"`
}

// snapshotUnitRecorder collects the requests the resource sends. The HTTP
// handler runs on the server's goroutine, so every access is mutex-guarded.
type snapshotUnitRecorder struct {
	mu       sync.Mutex
	requests []string
	body     map[string]interface{}
}

func (rec *snapshotUnitRecorder) record(method, path string) {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.requests = append(rec.requests, method+" "+path)
}

func (rec *snapshotUnitRecorder) decodeBody(r *http.Request) error {
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return err
	}
	rec.mu.Lock()
	defer rec.mu.Unlock()
	rec.body = body
	return nil
}

func (rec *snapshotUnitRecorder) lastBody() map[string]interface{} {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	return rec.body
}

func (rec *snapshotUnitRecorder) recorded() []string {
	rec.mu.Lock()
	defer rec.mu.Unlock()
	return append([]string(nil), rec.requests...)
}

func snapshotUnitClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)

	cred, err := credential.NewCredential(new(credential.Config).
		SetType("access_key").
		SetAccessKeyId("unit-test-access-key").
		SetAccessKeySecret("unit-test-secret-key"))
	if err != nil {
		t.Fatalf("building the static test credential: %s", err)
	}

	endpoints := new(sync.Map)
	// roaRequest snake-cases the product code before resolving the endpoint.
	endpoints.Store("pai_dsw", endpoint)
	config := &connectivity.Config{
		AccessKey: "unit-test-access-key", SecretKey: "unit-test-secret-key", Credential: cred,
		RegionId: "cn-hangzhou", AccountType: "ResourceOwner", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
		ClientReadTimeout: 5000, ClientConnectTimeout: 5000,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("building the connectivity client: %s", err)
	}
	return client
}

func snapshotUnitResource(t *testing.T, client *connectivity.AliyunClient) (snapshotUnitResourceAPI, fwschema.Schema) {
	t.Helper()
	ctx := context.Background()
	factory := paidsw.NewInstanceSnapshotResource()
	r, ok := factory.(snapshotUnitResourceAPI)
	if !ok {
		t.Fatalf("the resource factory returned %T, want schema, configure and import support", factory)
	}

	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building the resource schema: %v", schemaResp.Diagnostics)
	}

	var configureResp fwresource.ConfigureResponse
	r.Configure(ctx, fwresource.ConfigureRequest{ProviderData: client}, &configureResp)
	if configureResp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", configureResp.Diagnostics)
	}
	carrier, ok := r.(snapshotUnitClientCarrier)
	if !ok || carrier.Client() == nil {
		t.Fatal("the resource did not receive the test client")
	}
	return r, schemaResp.Schema
}

func snapshotUnitModel(t *testing.T, labels []snapshotUnitLabel, excludePaths []string) snapshotUnitStateModel {
	t.Helper()
	ctx := context.Background()

	labelsList := types.ListNull(snapshotUnitLabelType)
	if labels != nil {
		v, diags := types.ListValueFrom(ctx, snapshotUnitLabelType, labels)
		if diags.HasError() {
			t.Fatalf("building the labels list: %v", diags)
		}
		labelsList = v
	}

	excludeList := types.ListNull(types.StringType)
	if excludePaths != nil {
		v, diags := types.ListValueFrom(ctx, types.StringType, excludePaths)
		if diags.HasError() {
			t.Fatalf("building the exclude_paths list: %v", diags)
		}
		excludeList = v
	}

	return snapshotUnitStateModel{
		Id:                  types.StringNull(),
		InstanceId:          types.StringValue(snapshotUnitInstanceId),
		SnapshotName:        types.StringValue(snapshotUnitName),
		SnapshotDescription: types.StringValue("unit test description"),
		ImageUrl:            types.StringValue(snapshotUnitImageUrl),
		Overwrite:           types.BoolValue(true),
		Labels:              labelsList,
		ExcludePaths:        excludeList,
		SnapshotId:          types.StringNull(),
		Status:              types.StringNull(),
		ImageId:             types.StringNull(),
		ReasonCode:          types.StringNull(),
		ReasonMessage:       types.StringNull(),
		GmtCreateTime:       types.StringNull(),
		GmtModifiedTime:     types.StringNull(),
	}
}

func snapshotUnitPlan(t *testing.T, s fwschema.Schema, model snapshotUnitStateModel) tfsdk.Plan {
	t.Helper()
	plan := tfsdk.Plan{Schema: s}
	if diags := plan.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("building the plan: %v", diags)
	}
	return plan
}

func snapshotUnitState(t *testing.T, s fwschema.Schema, model snapshotUnitStateModel) tfsdk.State {
	t.Helper()
	state := tfsdk.State{Schema: s}
	if diags := state.Set(context.Background(), model); diags.HasError() {
		t.Fatalf("building the state: %v", diags)
	}
	return state
}

func snapshotUnitReadState(t *testing.T, state tfsdk.State) snapshotUnitStateModel {
	t.Helper()
	var model snapshotUnitStateModel
	if diags := state.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("reading the state: %v", diags)
	}
	return model
}

func snapshotUnitLabelModels(t *testing.T, list types.List) []snapshotUnitLabel {
	t.Helper()
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var out []snapshotUnitLabel
	if diags := list.ElementsAs(context.Background(), &out, false); diags.HasError() {
		t.Fatalf("decoding labels from state: %v", diags)
	}
	return out
}

func snapshotUnitStrings(t *testing.T, list types.List) []string {
	t.Helper()
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	var out []string
	if diags := list.ElementsAs(context.Background(), &out, false); diags.HasError() {
		t.Fatalf("decoding strings from state: %v", diags)
	}
	return out
}

// snapshotUnitAssertInputsMatch compares every configured input between the plan
// and the refreshed state. All inputs carry RequiresReplace, so equal values mean
// Terraform plans no change (and therefore no replacement) after a refresh;
// computed attributes cannot force a diff because they use UseStateForUnknown.
func snapshotUnitAssertInputsMatch(t *testing.T, plan, state snapshotUnitStateModel) {
	t.Helper()
	inputs := []struct {
		name         string
		plan, actual attr.Value
	}{
		{"instance_id", plan.InstanceId, state.InstanceId},
		{"snapshot_name", plan.SnapshotName, state.SnapshotName},
		{"snapshot_description", plan.SnapshotDescription, state.SnapshotDescription},
		{"image_url", plan.ImageUrl, state.ImageUrl},
		{"overwrite", plan.Overwrite, state.Overwrite},
		{"labels", plan.Labels, state.Labels},
		{"exclude_paths", plan.ExcludePaths, state.ExcludePaths},
	}
	for _, input := range inputs {
		if !input.plan.Equal(input.actual) {
			t.Errorf("%s drifted across refresh: plan %s, state %s", input.name, input.plan, input.actual)
		}
	}
}

// snapshotUnitAssertModelEqual compares every attribute of two models. Refresh
// must be idempotent: a second Read of the same API response may not move any
// value, otherwise the next plan would show a diff Terraform never configured.
func snapshotUnitAssertModelEqual(t *testing.T, want, got snapshotUnitStateModel) {
	t.Helper()
	values := []struct {
		name      string
		want, got attr.Value
	}{
		{"id", want.Id, got.Id},
		{"instance_id", want.InstanceId, got.InstanceId},
		{"snapshot_name", want.SnapshotName, got.SnapshotName},
		{"snapshot_description", want.SnapshotDescription, got.SnapshotDescription},
		{"image_url", want.ImageUrl, got.ImageUrl},
		{"overwrite", want.Overwrite, got.Overwrite},
		{"labels", want.Labels, got.Labels},
		{"exclude_paths", want.ExcludePaths, got.ExcludePaths},
		{"snapshot_id", want.SnapshotId, got.SnapshotId},
		{"status", want.Status, got.Status},
		{"image_id", want.ImageId, got.ImageId},
		{"reason_code", want.ReasonCode, got.ReasonCode},
		{"reason_message", want.ReasonMessage, got.ReasonMessage},
		{"gmt_create_time", want.GmtCreateTime, got.GmtCreateTime},
		{"gmt_modified_time", want.GmtModifiedTime, got.GmtModifiedTime},
	}
	for _, value := range values {
		if !value.want.Equal(value.got) {
			t.Errorf("%s changed across refresh: %s -> %s", value.name, value.want, value.got)
		}
	}
}

// snapshotUnitImportState builds the state the framework hands to ImportState: a
// null object of the resource schema, not a zero tfsdk value, so attribute
// writes resolve their path against the schema type.
func snapshotUnitImportState(ctx context.Context, t *testing.T, s fwschema.Schema) tfsdk.State {
	t.Helper()
	return tfsdk.State{
		Schema: s,
		Raw:    tftypes.NewValue(s.Type().TerraformType(ctx), nil),
	}
}

func snapshotUnitGetBody(labels, excludePaths string) string {
	return fmt.Sprintf(`{
  "InstanceId": %[1]q,
  "SnapshotId": %[2]q,
  "SnapshotName": %[3]q,
  "ImageUrl": %[4]q,
  "Status": "Creating",
  "ImageId": "img-unit",
  "Labels": %[5]s,
  "ExcludePaths": %[6]s,
  "GmtCreateTime": "2026-10-21T00:00:00Z",
  "GmtModifiedTime": "2026-10-21T00:00:00Z"
}`, snapshotUnitInstanceId, snapshotUnitSnapshotId, snapshotUnitName, snapshotUnitImageUrl, labels, excludePaths)
}

// TestUnitPaiDswInstanceSnapshotCreateRead covers the successful Create -> Read
// path: the create request is serialized to the ROA body the API documents, the
// composite ID is rebuilt from the Get response, and the create-only fields
// survive a refresh because Get does not return them.
func TestUnitPaiDswInstanceSnapshotCreateRead(t *testing.T) {
	rec := &snapshotUnitRecorder{}

	client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		rec.record(r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/instances/"+snapshotUnitInstanceId+"/snapshots":
			if err := rec.decodeBody(r); err != nil {
				t.Errorf("decoding the create request body: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, `{"SnapshotId":%q}`, snapshotUnitSnapshotId)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/instances/"+snapshotUnitInstanceId+"/snapshots/"+snapshotUnitSnapshotId:
			fmt.Fprint(w, snapshotUnitGetBody(
				`[{"Key":"env","Value":"tf-testacc"},{"Key":"purpose","Value":"snapshot"}]`,
				`["/home/admin/data","/tmp/cache"]`,
			))
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	r, s := snapshotUnitResource(t, client)
	ctx := context.Background()

	plan := snapshotUnitModel(t, []snapshotUnitLabel{
		{Key: types.StringValue("env"), Value: types.StringValue("tf-testacc")},
		{Key: types.StringValue("purpose"), Value: types.StringValue("snapshot")},
	}, []string{"/home/admin/data", "/tmp/cache"})

	createResp := fwresource.CreateResponse{State: tfsdk.State{Schema: s}}
	r.Create(ctx, fwresource.CreateRequest{Plan: snapshotUnitPlan(t, s, plan)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", createResp.Diagnostics)
	}
	created := snapshotUnitReadState(t, createResp.State)

	if got, want := created.Id.ValueString(), snapshotUnitInstanceId+":"+snapshotUnitSnapshotId; got != want {
		t.Errorf("composite id = %q, want %q", got, want)
	}
	if got := created.SnapshotId.ValueString(); got != snapshotUnitSnapshotId {
		t.Errorf("snapshot_id = %q, want %q", got, snapshotUnitSnapshotId)
	}
	if got := created.InstanceId.ValueString(); got != snapshotUnitInstanceId {
		t.Errorf("instance_id = %q, want %q", got, snapshotUnitInstanceId)
	}
	if got := created.Status.ValueString(); got != "Creating" {
		t.Errorf("status = %q, want Creating", got)
	}
	if got := created.SnapshotDescription.ValueString(); got != "unit test description" {
		t.Errorf("create-only snapshot_description = %q, want it preserved from the plan", got)
	}
	if !created.Overwrite.ValueBool() {
		t.Error("create-only overwrite = false after Create, want the planned true")
	}

	createBody := rec.lastBody()
	if createBody == nil {
		t.Fatal("Create did not send a request body")
	}
	if got := createBody["SnapshotName"]; got != snapshotUnitName {
		t.Errorf("create body SnapshotName = %v, want %q", got, snapshotUnitName)
	}
	if got := createBody["ImageUrl"]; got != snapshotUnitImageUrl {
		t.Errorf("create body ImageUrl = %v, want %q", got, snapshotUnitImageUrl)
	}
	if got := createBody["SnapshotDescription"]; got != "unit test description" {
		t.Errorf("create body SnapshotDescription = %v, want it serialized", got)
	}
	if got, ok := createBody["Overwrite"].(bool); !ok || !got {
		t.Errorf("create body Overwrite = %v, want a JSON true", createBody["Overwrite"])
	}
	rawLabels, ok := createBody["Labels"].([]interface{})
	if !ok || len(rawLabels) != 2 {
		t.Fatalf("create body Labels = %v, want 2 entries", createBody["Labels"])
	}
	firstLabel, ok := rawLabels[0].(map[string]interface{})
	if !ok || firstLabel["Key"] != "env" || firstLabel["Value"] != "tf-testacc" {
		t.Errorf("create body Labels[0] = %v, want {Key: env, Value: tf-testacc}", rawLabels[0])
	}
	rawPaths, ok := createBody["ExcludePaths"].([]interface{})
	if !ok || len(rawPaths) != 2 || rawPaths[0] != "/home/admin/data" || rawPaths[1] != "/tmp/cache" {
		t.Errorf("create body ExcludePaths = %v, want the configured paths in order", createBody["ExcludePaths"])
	}

	stateLabels := snapshotUnitLabelModels(t, created.Labels)
	if len(stateLabels) != 2 {
		t.Fatalf("state labels = %v, want 2 entries", stateLabels)
	}
	if stateLabels[0].Key.ValueString() != "env" || stateLabels[0].Value.ValueString() != "tf-testacc" {
		t.Errorf("state labels[0] = %v, want {env tf-testacc}", stateLabels[0])
	}
	if got := snapshotUnitStrings(t, created.ExcludePaths); len(got) != 2 || got[0] != "/home/admin/data" || got[1] != "/tmp/cache" {
		t.Errorf("state exclude_paths = %v, want the values returned by the API in order", got)
	}
	snapshotUnitAssertInputsMatch(t, plan, created)

	// Refresh twice: Get never returns snapshot_description/overwrite, so both
	// reads must keep the planned values instead of nulling them out, and the
	// Get-derived attributes must not move on a repeated refresh either.
	state := createResp.State
	for i := 0; i < 2; i++ {
		readResp := fwresource.ReadResponse{State: state}
		r.Read(ctx, fwresource.ReadRequest{State: state}, &readResp)
		if readResp.Diagnostics.HasError() {
			t.Fatalf("Read #%d returned diagnostics: %v", i+1, readResp.Diagnostics)
		}
		refreshed := snapshotUnitReadState(t, readResp.State)
		snapshotUnitAssertModelEqual(t, created, refreshed)
		snapshotUnitAssertInputsMatch(t, plan, refreshed)
		state = readResp.State
	}

	wantRequests := []string{
		"POST /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots",
		"GET /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots/" + snapshotUnitSnapshotId,
		"GET /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots/" + snapshotUnitSnapshotId,
		"GET /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots/" + snapshotUnitSnapshotId,
	}
	requests := rec.recorded()
	if len(requests) != len(wantRequests) {
		t.Fatalf("requests = %v, want %v", requests, wantRequests)
	}
	for i := range wantRequests {
		if requests[i] != wantRequests[i] {
			t.Errorf("request #%d = %q, want %q", i+1, requests[i], wantRequests[i])
		}
	}
}

// TestUnitPaiDswInstanceSnapshotCollectionsRoundTrip checks the labels and
// exclude_paths contract in both directions: request serialization (capital
// Key/Value, a value-less label omits Value) and response-to-state mapping,
// including a label the API returns without a value.
func TestUnitPaiDswInstanceSnapshotCollectionsRoundTrip(t *testing.T) {
	rec := &snapshotUnitRecorder{}
	client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			if err := rec.decodeBody(r); err != nil {
				t.Errorf("decoding the create request body: %s", err)
			}
			fmt.Fprintf(w, `{"SnapshotId":%q}`, snapshotUnitSnapshotId)
		case http.MethodGet:
			fmt.Fprint(w, snapshotUnitGetBody(
				`[{"Key":"env","Value":"tf-testacc"},{"Key":"purpose"}]`,
				`["/home/admin/data"]`,
			))
		default:
			t.Errorf("unexpected method %s", r.Method)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	r, s := snapshotUnitResource(t, client)
	ctx := context.Background()

	plan := snapshotUnitModel(t, []snapshotUnitLabel{
		{Key: types.StringValue("env"), Value: types.StringValue("tf-testacc")},
		{Key: types.StringValue("purpose"), Value: types.StringNull()},
	}, []string{"/home/admin/data"})

	createResp := fwresource.CreateResponse{State: tfsdk.State{Schema: s}}
	r.Create(ctx, fwresource.CreateRequest{Plan: snapshotUnitPlan(t, s, plan)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", createResp.Diagnostics)
	}

	createBody := rec.lastBody()
	rawLabels, ok := createBody["Labels"].([]interface{})
	if !ok || len(rawLabels) != 2 {
		t.Fatalf("create body Labels = %v, want 2 entries", createBody["Labels"])
	}
	secondLabel, ok := rawLabels[1].(map[string]interface{})
	if !ok {
		t.Fatalf("create body Labels[1] = %v, want an object", rawLabels[1])
	}
	if _, present := secondLabel["Value"]; present {
		t.Errorf("a label without a value must omit Value in the request, got %v", secondLabel)
	}
	if secondLabel["Key"] != "purpose" {
		t.Errorf("create body Labels[1].Key = %v, want purpose", secondLabel["Key"])
	}

	created := snapshotUnitReadState(t, createResp.State)
	stateLabels := snapshotUnitLabelModels(t, created.Labels)
	if len(stateLabels) != 2 {
		t.Fatalf("state labels = %v, want 2 entries", stateLabels)
	}
	if stateLabels[1].Key.ValueString() != "purpose" || !stateLabels[1].Value.IsNull() {
		t.Errorf("state labels[1] = %v, want {purpose null}", stateLabels[1])
	}
	if got := snapshotUnitStrings(t, created.ExcludePaths); len(got) != 1 || got[0] != "/home/admin/data" {
		t.Errorf("state exclude_paths = %v, want [/home/admin/data]", got)
	}
}

// TestUnitPaiDswInstanceSnapshotEmptyCollectionsOmitted checks that unset
// optional inputs are left out of the create body (rather than serialized as
// empty values) and come back as null collections when the API omits them.
func TestUnitPaiDswInstanceSnapshotEmptyCollectionsOmitted(t *testing.T) {
	rec := &snapshotUnitRecorder{}
	client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodPost:
			if err := rec.decodeBody(r); err != nil {
				t.Errorf("decoding the create request body: %s", err)
			}
			fmt.Fprintf(w, `{"SnapshotId":%q}`, snapshotUnitSnapshotId)
		case http.MethodGet:
			fmt.Fprint(w, snapshotUnitGetBody(`null`, `null`))
		default:
			t.Errorf("unexpected method %s", r.Method)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	r, s := snapshotUnitResource(t, client)
	ctx := context.Background()

	plan := snapshotUnitModel(t, nil, nil)
	plan.SnapshotDescription = types.StringNull()
	plan.Overwrite = types.BoolNull()

	createResp := fwresource.CreateResponse{State: tfsdk.State{Schema: s}}
	r.Create(ctx, fwresource.CreateRequest{Plan: snapshotUnitPlan(t, s, plan)}, &createResp)
	if createResp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", createResp.Diagnostics)
	}

	createBody := rec.lastBody()
	for _, key := range []string{"Labels", "ExcludePaths", "SnapshotDescription", "Overwrite"} {
		if _, present := createBody[key]; present {
			t.Errorf("unset %s must be omitted from the create body, got %v", key, createBody[key])
		}
	}

	created := snapshotUnitReadState(t, createResp.State)
	if !created.Labels.IsNull() {
		t.Errorf("labels = %v, want null when the API returns none", created.Labels)
	}
	if !created.ExcludePaths.IsNull() {
		t.Errorf("exclude_paths = %v, want null when the API returns none", created.ExcludePaths)
	}
	snapshotUnitAssertInputsMatch(t, plan, created)
}

// TestUnitPaiDswInstanceSnapshotReadNotFoundRemovesState checks that a 404 from
// GetInstanceSnapshot removes the resource from state instead of returning an
// error, which is what lets Terraform converge after an out-of-band deletion.
func TestUnitPaiDswInstanceSnapshotReadNotFoundRemovesState(t *testing.T) {
	client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"Code":"NotFound","Message":"The snapshot %s not found"}`, snapshotUnitSnapshotId)
	})
	r, s := snapshotUnitResource(t, client)
	ctx := context.Background()

	model := snapshotUnitModel(t, nil, nil)
	model.Id = types.StringValue(snapshotUnitInstanceId + ":" + snapshotUnitSnapshotId)
	model.SnapshotId = types.StringValue(snapshotUnitSnapshotId)
	state := snapshotUnitState(t, s, model)

	readResp := fwresource.ReadResponse{State: state}
	r.Read(ctx, fwresource.ReadRequest{State: state}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read on a missing snapshot returned diagnostics: %v", readResp.Diagnostics)
	}
	if !readResp.State.Raw.IsNull() {
		t.Errorf("Read on a missing snapshot left state populated: %s", readResp.State.Raw)
	}
}

// TestUnitPaiDswInstanceSnapshotCreateIdentitySurvivesReadFailure locks the
// create failure mode where CreateInstanceSnapshot has already created the
// snapshot but the read-back cannot see it. The composite identity must be in
// state before the read-back so Terraform keeps tracking (and can destroy) the
// created snapshot instead of orphaning it and creating a duplicate next run.
func TestUnitPaiDswInstanceSnapshotCreateIdentitySurvivesReadFailure(t *testing.T) {
	for _, tc := range []struct {
		name      string
		getStatus int
		getBody   string
	}{
		{
			name:      "get_not_found",
			getStatus: http.StatusNotFound,
			getBody:   `{"Code":"NotFound","Message":"The snapshot not found"}`,
		},
		{
			name:      "get_forbidden",
			getStatus: http.StatusForbidden,
			getBody:   `{"Code":"Forbidden","Message":"denied by policy"}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var getFails atomic.Bool
			getFails.Store(true)
			rec := &snapshotUnitRecorder{}

			client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				rec.record(r.Method, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodPost:
					fmt.Fprintf(w, `{"SnapshotId":%q}`, snapshotUnitSnapshotId)
				case http.MethodGet:
					if getFails.Load() {
						w.WriteHeader(tc.getStatus)
						fmt.Fprint(w, tc.getBody)
						return
					}
					fmt.Fprint(w, snapshotUnitGetBody(
						`[{"Key":"env","Value":"tf-testacc"}]`,
						`["/home/admin/data"]`,
					))
				case http.MethodDelete:
					fmt.Fprint(w, `{}`)
				default:
					t.Errorf("unexpected method %s", r.Method)
					w.WriteHeader(http.StatusNotFound)
				}
			})
			r, s := snapshotUnitResource(t, client)
			ctx := context.Background()

			plan := snapshotUnitModel(t, []snapshotUnitLabel{
				{Key: types.StringValue("env"), Value: types.StringValue("tf-testacc")},
			}, []string{"/home/admin/data"})

			createResp := fwresource.CreateResponse{State: tfsdk.State{Schema: s}}
			r.Create(ctx, fwresource.CreateRequest{Plan: snapshotUnitPlan(t, s, plan)}, &createResp)
			if !createResp.Diagnostics.HasError() {
				t.Fatal("Create reported success even though the read-back failed")
			}
			if createResp.State.Raw.IsNull() {
				t.Fatal("Create failed without persisting the snapshot identity; the created snapshot would be orphaned")
			}
			created := snapshotUnitReadState(t, createResp.State)
			if got, want := created.Id.ValueString(), snapshotUnitInstanceId+":"+snapshotUnitSnapshotId; got != want {
				t.Errorf("id after failed read-back = %q, want the composite identity %q", got, want)
			}
			if got := created.SnapshotId.ValueString(); got != snapshotUnitSnapshotId {
				t.Errorf("snapshot_id after failed read-back = %q, want %q", got, snapshotUnitSnapshotId)
			}
			if got := created.InstanceId.ValueString(); got != snapshotUnitInstanceId {
				t.Errorf("instance_id after failed read-back = %q, want %q", got, snapshotUnitInstanceId)
			}
			if got := created.SnapshotName.ValueString(); got != snapshotUnitName {
				t.Errorf("snapshot_name after failed read-back = %q, want the planned value", got)
			}

			// Recovery: once Get works again, Read must rebuild the resource from
			// the saved identity without a second Create.
			getFails.Store(false)
			readResp := fwresource.ReadResponse{State: createResp.State}
			r.Read(ctx, fwresource.ReadRequest{State: createResp.State}, &readResp)
			if readResp.Diagnostics.HasError() {
				t.Fatalf("Read after the failed read-back returned diagnostics: %v", readResp.Diagnostics)
			}
			recovered := snapshotUnitReadState(t, readResp.State)
			if got := recovered.Status.ValueString(); got != "Creating" {
				t.Errorf("status after recovery = %q, want Creating", got)
			}
			if got := len(snapshotUnitLabelModels(t, recovered.Labels)); got != 1 {
				t.Errorf("labels after recovery = %d entries, want 1", got)
			}

			// The saved identity is also everything Delete needs.
			deleteResp := fwresource.DeleteResponse{State: readResp.State}
			r.Delete(ctx, fwresource.DeleteRequest{State: readResp.State}, &deleteResp)
			if deleteResp.Diagnostics.HasError() {
				t.Fatalf("Delete using the saved identity returned diagnostics: %v", deleteResp.Diagnostics)
			}

			wantRequests := []string{
				"POST /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots",
				"GET /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots/" + snapshotUnitSnapshotId,
				"GET /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots/" + snapshotUnitSnapshotId,
				"DELETE /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots/" + snapshotUnitSnapshotId,
			}
			requests := rec.recorded()
			if len(requests) != len(wantRequests) {
				t.Fatalf("requests = %v, want exactly one create attempt followed by read-back, recovery and delete: %v", requests, wantRequests)
			}
			for i := range wantRequests {
				if requests[i] != wantRequests[i] {
					t.Errorf("request #%d = %q, want %q", i+1, requests[i], wantRequests[i])
				}
			}
		})
	}
}

// TestUnitPaiDswInstanceSnapshotEmptyCollectionsPreserved locks the Optional
// (non-Computed) empty-collection contract: an explicitly configured [] must
// stay an empty list in state when the API returns no elements. Flipping it to
// null would make the next plan differ, and RequiresReplace would turn that
// diff into a replacement.
func TestUnitPaiDswInstanceSnapshotEmptyCollectionsPreserved(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		labelsResponse       string
		excludePathsResponse string
	}{
		{name: "null_response", labelsResponse: `null`, excludePathsResponse: `null`},
		{name: "empty_array_response", labelsResponse: `[]`, excludePathsResponse: `[]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := &snapshotUnitRecorder{}
			client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch r.Method {
				case http.MethodPost:
					if err := rec.decodeBody(r); err != nil {
						t.Errorf("decoding the create request body: %s", err)
					}
					fmt.Fprintf(w, `{"SnapshotId":%q}`, snapshotUnitSnapshotId)
				case http.MethodGet:
					fmt.Fprint(w, snapshotUnitGetBody(tc.labelsResponse, tc.excludePathsResponse))
				default:
					t.Errorf("unexpected method %s", r.Method)
					w.WriteHeader(http.StatusNotFound)
				}
			})
			r, s := snapshotUnitResource(t, client)
			ctx := context.Background()

			plan := snapshotUnitModel(t, []snapshotUnitLabel{}, []string{})

			createResp := fwresource.CreateResponse{State: tfsdk.State{Schema: s}}
			r.Create(ctx, fwresource.CreateRequest{Plan: snapshotUnitPlan(t, s, plan)}, &createResp)
			if createResp.Diagnostics.HasError() {
				t.Fatalf("Create returned diagnostics: %v", createResp.Diagnostics)
			}

			createBody := rec.lastBody()
			for _, key := range []string{"Labels", "ExcludePaths"} {
				if _, present := createBody[key]; present {
					t.Errorf("an explicitly empty %s must be omitted from the create body, got %v", key, createBody[key])
				}
			}

			created := snapshotUnitReadState(t, createResp.State)
			if created.Labels.IsNull() {
				t.Error("labels = null, want the configured empty list preserved")
			} else if got := len(created.Labels.Elements()); got != 0 {
				t.Errorf("labels = %d entries, want 0", got)
			}
			if created.ExcludePaths.IsNull() {
				t.Error("exclude_paths = null, want the configured empty list preserved")
			} else if got := len(created.ExcludePaths.Elements()); got != 0 {
				t.Errorf("exclude_paths = %d entries, want 0", got)
			}
			snapshotUnitAssertInputsMatch(t, plan, created)

			// Repeated refresh must keep the empty collections and every other
			// planned input stable, producing no diff (and no replacement).
			state := createResp.State
			for i := 0; i < 2; i++ {
				readResp := fwresource.ReadResponse{State: state}
				r.Read(ctx, fwresource.ReadRequest{State: state}, &readResp)
				if readResp.Diagnostics.HasError() {
					t.Fatalf("Read #%d returned diagnostics: %v", i+1, readResp.Diagnostics)
				}
				refreshed := snapshotUnitReadState(t, readResp.State)
				snapshotUnitAssertModelEqual(t, created, refreshed)
				snapshotUnitAssertInputsMatch(t, plan, refreshed)
				state = readResp.State
			}
		})
	}
}

// TestUnitPaiDswInstanceSnapshotDelete covers Delete success and its not-found
// idempotency: a re-run after an out-of-band deletion converges silently, while
// other API errors still surface.
func TestUnitPaiDswInstanceSnapshotDelete(t *testing.T) {
	for _, tc := range []struct {
		name       string
		statusCode int
		body       string
		wantError  bool
	}{
		{
			name:       "success",
			statusCode: http.StatusOK,
			body:       `{}`,
		},
		{
			name:       "not_found_is_idempotent",
			statusCode: http.StatusNotFound,
			body:       `{"Code":"NotFound","Message":"The snapshot not found"}`,
		},
		{
			name:       "server_error_surfaces",
			statusCode: http.StatusInternalServerError,
			body:       `{"Code":"InternalError","Message":"internal error"}`,
			wantError:  true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := &snapshotUnitRecorder{}
			client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
				rec.record(r.Method, r.URL.Path)
				if r.Method != http.MethodDelete {
					t.Errorf("unexpected method %s", r.Method)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				fmt.Fprint(w, tc.body)
			})
			r, s := snapshotUnitResource(t, client)
			ctx := context.Background()

			model := snapshotUnitModel(t, nil, nil)
			model.Id = types.StringValue(snapshotUnitInstanceId + ":" + snapshotUnitSnapshotId)
			model.SnapshotId = types.StringValue(snapshotUnitSnapshotId)
			state := snapshotUnitState(t, s, model)

			deleteResp := fwresource.DeleteResponse{State: state}
			r.Delete(ctx, fwresource.DeleteRequest{State: state}, &deleteResp)

			// The ROA client retries 5xx responses, so only the target of every
			// attempt is asserted, not the attempt count.
			requests := rec.recorded()
			if len(requests) == 0 {
				t.Fatal("Delete sent no request")
			}
			want := "DELETE /api/v2/instances/" + snapshotUnitInstanceId + "/snapshots/" + snapshotUnitSnapshotId
			for i, request := range requests {
				if request != want {
					t.Errorf("delete request #%d = %q, want %q", i+1, request, want)
				}
			}
			if tc.wantError && !deleteResp.Diagnostics.HasError() {
				t.Error("Delete swallowed an unexpected API error")
			}
			if !tc.wantError && deleteResp.Diagnostics.HasError() {
				t.Errorf("Delete returned diagnostics: %v", deleteResp.Diagnostics)
			}
		})
	}
}

// TestUnitPaiDswInstanceSnapshotImport parses the composite ID, seeds the
// identity attributes and reads the rest back from GetInstanceSnapshot.
func TestUnitPaiDswInstanceSnapshotImport(t *testing.T) {
	client := snapshotUnitClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, snapshotUnitGetBody(
			`[{"Key":"env","Value":"tf-testacc"}]`,
			`["/home/admin/data"]`,
		))
	})
	r, s := snapshotUnitResource(t, client)
	ctx := context.Background()

	importResp := fwresource.ImportStateResponse{State: snapshotUnitImportState(ctx, t, s)}
	r.ImportState(ctx, fwresource.ImportStateRequest{ID: snapshotUnitInstanceId + ":" + snapshotUnitSnapshotId}, &importResp)
	if importResp.Diagnostics.HasError() {
		t.Fatalf("ImportState returned diagnostics: %v", importResp.Diagnostics)
	}
	imported := snapshotUnitReadState(t, importResp.State)
	if got, want := imported.Id.ValueString(), snapshotUnitInstanceId+":"+snapshotUnitSnapshotId; got != want {
		t.Errorf("imported id = %q, want %q", got, want)
	}
	if got := imported.InstanceId.ValueString(); got != snapshotUnitInstanceId {
		t.Errorf("imported instance_id = %q, want %q", got, snapshotUnitInstanceId)
	}
	if got := imported.SnapshotId.ValueString(); got != snapshotUnitSnapshotId {
		t.Errorf("imported snapshot_id = %q, want %q", got, snapshotUnitSnapshotId)
	}

	readResp := fwresource.ReadResponse{State: importResp.State}
	r.Read(ctx, fwresource.ReadRequest{State: importResp.State}, &readResp)
	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read after import returned diagnostics: %v", readResp.Diagnostics)
	}
	readBack := snapshotUnitReadState(t, readResp.State)
	if got := readBack.SnapshotName.ValueString(); got != snapshotUnitName {
		t.Errorf("snapshot_name after import = %q, want %q", got, snapshotUnitName)
	}
	if got := readBack.ImageUrl.ValueString(); got != snapshotUnitImageUrl {
		t.Errorf("image_url after import = %q, want %q", got, snapshotUnitImageUrl)
	}
	if got := len(snapshotUnitLabelModels(t, readBack.Labels)); got != 1 {
		t.Errorf("labels after import = %d entries, want 1", got)
	}
	if got := snapshotUnitStrings(t, readBack.ExcludePaths); len(got) != 1 || got[0] != "/home/admin/data" {
		t.Errorf("exclude_paths after import = %v, want [/home/admin/data]", got)
	}
	// GetInstanceSnapshot does not return the create-only fields, so an import
	// must leave them null rather than inventing values.
	if !readBack.SnapshotDescription.IsNull() {
		t.Errorf("snapshot_description after import = %q, want null", readBack.SnapshotDescription)
	}
	if !readBack.Overwrite.IsNull() {
		t.Errorf("overwrite after import = %v, want null", readBack.Overwrite)
	}
}

// TestUnitPaiDswInstanceSnapshotImportRejectsMalformedId checks the composite ID
// contract: exactly one non-empty "<instance_id>:<snapshot_id>" pair is accepted.
func TestUnitPaiDswInstanceSnapshotImportRejectsMalformedId(t *testing.T) {
	r, s := snapshotUnitResource(t, snapshotUnitClient(t, func(w http.ResponseWriter, req *http.Request) {
		t.Errorf("a malformed import ID must not call the API, got %s %s", req.Method, req.URL.Path)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	ctx := context.Background()

	for _, id := range []string{"", "i-only", "i:a:b", ":snapshot", "instance:"} {
		importResp := fwresource.ImportStateResponse{State: tfsdk.State{Schema: s}}
		r.ImportState(ctx, fwresource.ImportStateRequest{ID: id}, &importResp)
		if !importResp.Diagnostics.HasError() {
			t.Errorf("ImportState(%q) accepted a malformed composite ID", id)
		}
	}
}
