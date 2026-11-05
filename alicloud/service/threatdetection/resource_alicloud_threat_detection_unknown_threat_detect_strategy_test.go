package threatdetection_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-framework/path"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const unknownThreatDetectStrategyResourceId = "alicloud_threat_detection_unknown_threat_detect_strategy.default"

// ---- acceptance test ----

// TestAccAliCloudThreatDetectionUnknownThreatDetectStrategy_basic covers the
// create/read/update/delete lifecycle of the strategy. The strategy requires
// Agentic-EDR to be purchased and activated on the test account; without it
// CreateUnknownThreatDetectStrategy returns EdrSaleAdmissionDenied.
//
// lintignore: AT001
func TestAccAliCloudThreatDetectionUnknownThreatDetectStrategy_basic(t *testing.T) {
	rand := tfacctest.RandIntRange(10000, 99999)
	region := os.Getenv("ALICLOUD_REGION")
	if region == "" {
		region = "cn-hangzhou"
	}
	name := fmt.Sprintf("tf-testacc%sutds%d", region, rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		IDRefreshName:            unknownThreatDetectStrategyResourceId,
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckUnknownThreatDetectStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUnknownThreatDetectStrategyBasicConfig(name, "hash", 1, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "name", name),
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "study_mode", "hash"),
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "duration_days_after_init", "1"),
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "duration_days_after_stop", "1"),
					resource.TestCheckResourceAttrSet(unknownThreatDetectStrategyResourceId, "id"),
					resource.TestCheckResourceAttrSet(unknownThreatDetectStrategyResourceId, "machine_count"),
				),
			},
			{
				Config: testAccUnknownThreatDetectStrategyBasicConfig(name, "path", 3, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "name", name),
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "study_mode", "path"),
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "duration_days_after_init", "3"),
					resource.TestCheckResourceAttr(unknownThreatDetectStrategyResourceId, "duration_days_after_stop", "5"),
				),
			},
			{
				ResourceName:      unknownThreatDetectStrategyResourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUnknownThreatDetectStrategyBasicConfig(name, studyMode string, initDays, stopDays int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_threat_detection_unknown_threat_detect_strategy" "default" {
  name                      = var.name
  asset_selection_type      = "UNKNOWN_THREAT_DETECT_CONFIG_tf_example"
  study_mode                = "%s"
  duration_days_after_init  = %d
  duration_days_after_stop  = %d
}
`, name, studyMode, initDays, stopDays)
}

func testAccCheckUnknownThreatDetectStrategyDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	return checkUnknownThreatDetectStrategyDestroyed(s, func(id string) error {
		_, err := threatdetection.NewService(client).DescribeUnknownThreatDetectStrategy(id)
		return err
	})
}

// checkUnknownThreatDetectStrategyDestroyed accepts only the NotFound that
// DescribeUnknownThreatDetectStrategy constructs itself. A permission or other
// non-404 failure whose message contains "NotFound" must surface as a cleanup
// failure instead of pretending the strategy was deleted.
func checkUnknownThreatDetectStrategyDestroyed(s *terraform.State, describe func(string) error) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_threat_detection_unknown_threat_detect_strategy" {
			continue
		}
		err := describe(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("strategy %s still exists", rs.Primary.ID)
		}
		if errs.IsProviderNotFound(err) {
			continue
		}
		return fmt.Errorf("describing strategy %s failed: %s", rs.Primary.ID, err)
	}
	return nil
}

// ---- sweeper ----

func init() {
	resource.AddTestSweepers("alicloud_threat_detection_unknown_threat_detect_strategy", &resource.Sweeper{
		Name: "alicloud_threat_detection_unknown_threat_detect_strategy",
		F:    testSweepUnknownThreatDetectStrategy,
	})
}

func testSweepUnknownThreatDetectStrategy(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	return sweepUnknownThreatDetectStrategy(client)
}

// sweepUnknownThreatDetectStrategy lists every page before deleting anything,
// treats a confirmed typed NotFound as an already-finished cleanup, and reports
// list/deletion failures instead of logging them as a successful run.
func sweepUnknownThreatDetectStrategy(client *connectivity.AliyunClient) error {
	prefixes := []string{
		"tf-testacc",
		"tf_testacc",
		"terraform-example",
	}

	type strategyRef struct {
		id   string
		name string
	}
	var targets []strategyRef

	const listAction = "ListUnknownThreatDetectStrategy"
	query := map[string]interface{}{
		"CurrentPage": 1,
		"PageSize":    helper.PageSizeLarge,
	}
	for {
		var response map[string]interface{}
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err := retry.RetryContext(context.Background(), 2*time.Minute, func() *retry.RetryError {
			var callErr error
			response, callErr = client.RpcPost("Sas", "2018-12-03", listAction, query, map[string]interface{}{}, true)
			if callErr != nil {
				if helper.NeedRetry(callErr) {
					wait()
					return retry.RetryableError(callErr)
				}
				return retry.NonRetryableError(callErr)
			}
			return nil
		})
		logs.AddDebug(listAction, response, query)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_threat_detection_unknown_threat_detect_strategy", listAction, errs.AlibabaCloudSdkGoERROR)
		}
		raw, ok := response["Data"].([]interface{})
		if !ok {
			return errs.WrapErrorf(fmt.Errorf("unexpected %s response type %T for $.Data", listAction, response["Data"]), errs.FailedGetAttributeMsg, listAction, "$.Data", response)
		}
		for _, v := range raw {
			item, ok := v.(map[string]interface{})
			if !ok {
				return errs.WrapErrorf(fmt.Errorf("unexpected %s element type %T in $.Data", listAction, v), errs.FailedGetAttributeMsg, listAction, "$.Data", response)
			}
			name := fmt.Sprint(item["Name"])
			skip := true
			if !sweep.SweepAll() {
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						skip = false
						break
					}
				}
			} else {
				skip = false
			}
			if skip {
				continue
			}
			id, idErr := threatdetection.UnknownThreatDetectStrategyId(item["Id"])
			if idErr != nil {
				return errs.WrapErrorf(idErr, errs.FailedGetAttributeMsg, listAction, "$.Data[].Id", response)
			}
			targets = append(targets, strategyRef{id: id, name: name})
		}
		if len(raw) < helper.PageSizeLarge {
			break
		}
		query["CurrentPage"] = query["CurrentPage"].(int) + 1
	}

	service := threatdetection.NewService(client)
	var cleanupErrs []error
	for _, target := range targets {
		if err := service.DeleteUnknownThreatDetectStrategy(target.id); err != nil {
			cleanupErrs = append(cleanupErrs, err)
		}
	}
	if len(cleanupErrs) > 0 {
		return errors.Join(cleanupErrs...)
	}
	return nil
}

// ---- unit test scaffolding ----

type unknownThreatDetectStrategyRecorder struct {
	mu      sync.Mutex
	actions []string
	forms   []url.Values
}

func (r *unknownThreatDetectStrategyRecorder) record(action string, form url.Values) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.actions = append(r.actions, action)
	r.forms = append(r.forms, form)
}

func (r *unknownThreatDetectStrategyRecorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.actions)
}

func (r *unknownThreatDetectStrategyRecorder) actionList() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.actions...)
}

func (r *unknownThreatDetectStrategyRecorder) formsByAction(action string) []url.Values {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []url.Values
	for i, recorded := range r.actions {
		if recorded == action {
			result = append(result, r.forms[i])
		}
	}
	return result
}

func unknownThreatDetectStrategyTestClient(t *testing.T, handler func(action string, form url.Values) (int, interface{})) (*connectivity.AliyunClient, *unknownThreatDetectStrategyRecorder) {
	t.Helper()
	recorder := &unknownThreatDetectStrategyRecorder{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		action := r.Form.Get("Action")
		recorder.record(action, r.Form)
		status, body := handler(action, r.Form)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if body != nil {
			if err := json.NewEncoder(w).Encode(body); err != nil {
				t.Error(err)
			}
		}
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	endpoints.Store("sas", endpoint)
	return client, recorder
}

func newUnknownThreatDetectStrategyUnitResource(t *testing.T, apiClient *connectivity.AliyunClient) fwresource.Resource {
	t.Helper()
	r := threatdetection.NewUnknownThreatDetectStrategyResource()
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatalf("%T does not implement resource.ResourceWithConfigure", r)
	}
	var resp fwresource.ConfigureResponse
	configurable.Configure(context.Background(), fwresource.ConfigureRequest{ProviderData: apiClient}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configuring the resource: %v", resp.Diagnostics)
	}
	return r
}

func unknownThreatDetectStrategySchemaAndType(t *testing.T, r fwresource.Resource) (fwschema.Schema, tftypes.Object) {
	t.Helper()
	ctx := context.Background()
	var schemaResp fwresource.SchemaResponse
	r.Schema(ctx, fwresource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("building schema: %v", schemaResp.Diagnostics)
	}
	objType, ok := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", schemaResp.Schema.Type().TerraformType(ctx))
	}
	return schemaResp.Schema, objType
}

func unknownThreatDetectStrategyValue(objType tftypes.Object, overrides map[string]tftypes.Value) tftypes.Value {
	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, attrType := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(attrType, nil)
	}
	for name, value := range overrides {
		values[name] = value
	}
	return tftypes.NewValue(objType, values)
}

func unknownThreatDetectStrategyBaseOverrides() map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"name":                 tftypes.NewValue(tftypes.String, "test-strategy"),
		"asset_selection_type": tftypes.NewValue(tftypes.String, "UNKNOWN_THREAT_DETECT_CONFIG_test"),
	}
}

func unknownThreatDetectStrategyCreatePlanValue(t *testing.T, r fwresource.Resource, extra map[string]tftypes.Value) (fwschema.Schema, tftypes.Value) {
	t.Helper()
	resourceSchema, objType := unknownThreatDetectStrategySchemaAndType(t, r)
	overrides := unknownThreatDetectStrategyBaseOverrides()
	overrides["id"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	for name, value := range extra {
		overrides[name] = value
	}
	return resourceSchema, unknownThreatDetectStrategyValue(objType, overrides)
}

func unknownThreatDetectStrategyStateString(t *testing.T, state tfsdk.State, name string) string {
	t.Helper()
	var value types.String
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return value.ValueString()
}

func unknownThreatDetectStrategyStateInt64(t *testing.T, state tfsdk.State, name string) string {
	t.Helper()
	var value types.Int64
	if diags := state.GetAttribute(context.Background(), path.Root(name), &value); diags.HasError() {
		t.Fatalf("reading state attribute %q: %v", name, diags)
	}
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return fmt.Sprint(value.ValueInt64())
}

func unknownThreatDetectStrategyStrategyResponse(id interface{}, name, assetSelectionType string, extras map[string]interface{}) map[string]interface{} {
	item := map[string]interface{}{
		"Id":                 id,
		"Name":               name,
		"AssetSelectionType": assetSelectionType,
	}
	for k, v := range extras {
		item[k] = v
	}
	return map[string]interface{}{
		"Data":     []interface{}{item},
		"PageInfo": map[string]interface{}{"TotalCount": 1},
	}
}

// ---- Q-1: canonical Id identity ----

func TestUnitUnknownThreatDetectStrategyIdConversion(t *testing.T) {
	cases := map[string]struct {
		raw     interface{}
		want    string
		wantErr bool
	}{
		"string":             {raw: "123", want: "123"},
		"emptyString":        {raw: "", wantErr: true},
		"float":              {raw: float64(123), want: "123"},
		"floatExactLimit":    {raw: float64(1 << 53), want: "9007199254740992"},
		"floatBeyondExact":   {raw: float64(1<<53 + 2), wantErr: true},
		"floatFraction":      {raw: 1.5, wantErr: true},
		"jsonNumber":         {raw: json.Number("1183029123456789012"), want: "1183029123456789012"},
		"jsonNumberFraction": {raw: json.Number("1.5"), wantErr: true},
		"int":                {raw: 123, want: "123"},
		"intNegative":        {raw: -1, wantErr: true},
		"int64":              {raw: int64(123), want: "123"},
		"map":                {raw: map[string]interface{}{"id": "123"}, wantErr: true},
		"nil":                {raw: nil, wantErr: true},
		"bool":               {raw: true, wantErr: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := threatdetection.UnknownThreatDetectStrategyId(c.raw)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("UnknownThreatDetectStrategyId(%v) = %q, want %q", c.raw, got, c.want)
			}
		})
	}
}

// ---- Q-2: schema enum ----

func TestUnitUnknownThreatDetectStrategySchemaEnum(t *testing.T) {
	r := threatdetection.NewUnknownThreatDetectStrategyResource()
	resourceSchema, _ := unknownThreatDetectStrategySchemaAndType(t, r)
	studyMode, ok := resourceSchema.Attributes["study_mode"].(fwschema.StringAttribute)
	if !ok {
		t.Fatalf("study_mode attribute is %T, want schema.StringAttribute", resourceSchema.Attributes["study_mode"])
	}
	rejected := func(value string) bool {
		var resp validator.StringResponse
		for _, v := range studyMode.Validators {
			v.ValidateString(context.Background(), validator.StringRequest{
				Path:        path.Root("study_mode"),
				ConfigValue: types.StringValue(value),
			}, &resp)
		}
		return resp.Diagnostics.HasError()
	}
	for _, value := range []string{"hash", "path"} {
		if rejected(value) {
			t.Fatalf("%s must be accepted", value)
		}
	}
	for _, value := range []string{"Hash", "other", ""} {
		if !rejected(value) {
			t.Fatalf("%q must be rejected", value)
		}
	}
}

// ---- Q-3: create request shape ----

func TestUnitUnknownThreatDetectStrategyCreateRequestShape(t *testing.T) {
	cases := map[string]struct {
		studyMode     tftypes.Value
		initDays      tftypes.Value
		stopDays      tftypes.Value
		wantStudy     string
		wantInit      string
		wantStop      string
		expectNoStudy bool
	}{
		"allSet": {
			studyMode: tftypes.NewValue(tftypes.String, "hash"),
			initDays:  tftypes.NewValue(tftypes.Number, int64(1)),
			stopDays:  tftypes.NewValue(tftypes.Number, int64(2)),
			wantStudy: "hash", wantInit: "1", wantStop: "2",
		},
		"onlyRequired": {
			studyMode:     tftypes.NewValue(tftypes.String, nil),
			initDays:      tftypes.NewValue(tftypes.Number, nil),
			stopDays:      tftypes.NewValue(tftypes.Number, nil),
			expectNoStudy: true,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, recorder := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
				switch action {
				case "CreateUnknownThreatDetectStrategy":
					return http.StatusOK, map[string]interface{}{"Id": float64(123456)}
				case "ListUnknownThreatDetectStrategy":
					return http.StatusOK, unknownThreatDetectStrategyStrategyResponse(float64(123456), "test-strategy", "UNKNOWN_THREAT_DETECT_CONFIG_test", map[string]interface{}{
						"StudyMode":             c.wantStudy,
						"DurationDaysAfterInit": float64(1),
						"DurationDaysAfterStop": float64(2),
						"MachineCount":          float64(0),
					})
				}
				return http.StatusOK, map[string]interface{}{}
			})
			r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
			resourceSchema, planValue := unknownThreatDetectStrategyCreatePlanValue(t, r, map[string]tftypes.Value{
				"study_mode":               c.studyMode,
				"duration_days_after_init": c.initDays,
				"duration_days_after_stop": c.stopDays,
			})
			req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
			resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

			r.Create(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
			}
			createForms := recorder.formsByAction("CreateUnknownThreatDetectStrategy")
			if len(createForms) != 1 {
				t.Fatalf("expected one Create call, got %v", recorder.actionList())
			}
			form := createForms[0]
			if form.Get("Name") != "test-strategy" {
				t.Fatalf("Name = %q", form.Get("Name"))
			}
			if form.Get("AssetSelectionType") != "UNKNOWN_THREAT_DETECT_CONFIG_test" {
				t.Fatalf("AssetSelectionType = %q", form.Get("AssetSelectionType"))
			}
			if c.expectNoStudy && form.Get("StudyMode") != "" {
				t.Fatalf("StudyMode must not be sent when unset, got %q", form.Get("StudyMode"))
			}
			if c.wantStudy != "" && form.Get("StudyMode") != c.wantStudy {
				t.Fatalf("StudyMode = %q, want %q", form.Get("StudyMode"), c.wantStudy)
			}
			if got := unknownThreatDetectStrategyStateString(t, resp.State, "id"); got != "123456" {
				t.Fatalf("id = %q, want 123456", got)
			}
		})
	}
}

func TestUnitUnknownThreatDetectStrategyCreateNumericId(t *testing.T) {
	apiClient, recorder := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
		switch action {
		case "CreateUnknownThreatDetectStrategy":
			return http.StatusOK, map[string]interface{}{"Id": float64(1189784)}
		case "ListUnknownThreatDetectStrategy":
			return http.StatusOK, unknownThreatDetectStrategyStrategyResponse(float64(1189784), "test-strategy", "UNKNOWN_THREAT_DETECT_CONFIG_test", nil)
		}
		return http.StatusOK, map[string]interface{}{}
	})
	r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
	resourceSchema, planValue := unknownThreatDetectStrategyCreatePlanValue(t, r, nil)
	req := fwresource.CreateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.CreateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Create(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Create returned diagnostics: %v", resp.Diagnostics)
	}
	if got := unknownThreatDetectStrategyStateString(t, resp.State, "id"); got != "1189784" {
		t.Fatalf("id = %q, want canonical 1189784", got)
	}
	listForms := recorder.formsByAction("ListUnknownThreatDetectStrategy")
	if len(listForms) != 1 || listForms[0].Get("Id") != "1189784" {
		t.Fatalf("unexpected Id filter: %v", listForms)
	}
}

// ---- Q-4: read valid response mapping ----

func TestUnitUnknownThreatDetectStrategyReadValidResponse(t *testing.T) {
	apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusOK, unknownThreatDetectStrategyStrategyResponse(float64(42), "test-strategy", "UNKNOWN_THREAT_DETECT_CONFIG_test", map[string]interface{}{
			"StudyMode":             "path",
			"DurationDaysAfterInit": float64(3),
			"DurationDaysAfterStop": float64(5),
			"MachineCount":          float64(10),
		})
	})
	r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
	resourceSchema, objType := unknownThreatDetectStrategySchemaAndType(t, r)
	stateValue := unknownThreatDetectStrategyValue(objType, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "42"),
	})
	readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("Read returned diagnostics: %v", readResp.Diagnostics)
	}
	if got := unknownThreatDetectStrategyStateString(t, readResp.State, "name"); got != "test-strategy" {
		t.Fatalf("name = %q", got)
	}
	if got := unknownThreatDetectStrategyStateString(t, readResp.State, "asset_selection_type"); got != "UNKNOWN_THREAT_DETECT_CONFIG_test" {
		t.Fatalf("asset_selection_type = %q", got)
	}
	if got := unknownThreatDetectStrategyStateString(t, readResp.State, "study_mode"); got != "path" {
		t.Fatalf("study_mode = %q", got)
	}
	if got := unknownThreatDetectStrategyStateInt64(t, readResp.State, "duration_days_after_init"); got != "3" {
		t.Fatalf("duration_days_after_init = %q", got)
	}
	if got := unknownThreatDetectStrategyStateInt64(t, readResp.State, "machine_count"); got != "10" {
		t.Fatalf("machine_count = %q", got)
	}
}

// ---- Q-5: read rejects missing required fields ----

func TestUnitUnknownThreatDetectStrategyReadRejectsMissingRequired(t *testing.T) {
	for name, item := range map[string]map[string]interface{}{
		"missingName":               {"Id": float64(42), "AssetSelectionType": "cfg"},
		"missingAssetSelectionType": {"Id": float64(42), "Name": "n"},
	} {
		t.Run(name, func(t *testing.T) {
			apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return http.StatusOK, map[string]interface{}{"Data": []interface{}{item}}
			})
			r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
			resourceSchema, objType := unknownThreatDetectStrategySchemaAndType(t, r)
			stateValue := unknownThreatDetectStrategyValue(objType, map[string]tftypes.Value{
				"id": tftypes.NewValue(tftypes.String, "42"),
			})
			readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
			readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			r.Read(context.Background(), readReq, readResp)

			if !readResp.Diagnostics.HasError() {
				t.Fatal("a response missing a required field must fail the read")
			}
		})
	}
}

// ---- Q-6: collection read resolves the target by exact Id ----

func TestUnitUnknownThreatDetectStrategyDescribeIdentity(t *testing.T) {
	target := map[string]interface{}{
		"Id":                 float64(210),
		"Name":               "target",
		"AssetSelectionType": "cfg-target",
	}
	decoy := map[string]interface{}{
		"Id":                 float64(999),
		"Name":               "decoy",
		"AssetSelectionType": "cfg-decoy",
	}
	cases := map[string]struct {
		items        []interface{}
		wantName     string
		wantNotFound bool
	}{
		"decoyFirstTargetSecond": {items: []interface{}{decoy, target}, wantName: "target"},
		"targetFirstDecoySecond": {items: []interface{}{target, decoy}, wantName: "target"},
		"decoyOnly":              {items: []interface{}{decoy}, wantNotFound: true},
		"emptyList":              {items: []interface{}{}, wantNotFound: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, recorder := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return http.StatusOK, map[string]interface{}{"Data": c.items}
			})
			object, err := threatdetection.NewService(apiClient).DescribeUnknownThreatDetectStrategy("210")

			if c.wantNotFound {
				if err == nil {
					t.Fatalf("expected NotFound, got object %v", object)
				}
				if !errs.NotFoundError(err) {
					t.Fatalf("expected a NotFound error, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("Describe returned an error: %v", err)
				}
				if got := fmt.Sprint(object["Name"]); got != c.wantName {
					t.Fatalf("resolved Name = %q, want %q", got, c.wantName)
				}
			}
			forms := recorder.formsByAction("ListUnknownThreatDetectStrategy")
			if len(forms) != 1 || forms[0].Get("Id") != "210" {
				t.Fatalf("unexpected Id filter: %v", forms)
			}
		})
	}
}

// ---- Q-7: not-found versus malformed responses ----

func TestUnitUnknownThreatDetectStrategyNotFoundClassification(t *testing.T) {
	t.Run("typed404IsNotFound", func(t *testing.T) {
		apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusNotFound, map[string]interface{}{"Code": "404", "Message": "not found"}
		})
		_, err := threatdetection.NewService(apiClient).DescribeUnknownThreatDetectStrategy("gone")
		if err == nil || !errs.NotFoundError(err) {
			t.Fatalf("a typed 404 must be NotFound, got %v", err)
		}
	})

	t.Run("responseCode404WithEmptyArrayIsNotFound", func(t *testing.T) {
		apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusOK, map[string]interface{}{"Code": "404", "Data": []interface{}{}}
		})
		_, err := threatdetection.NewService(apiClient).DescribeUnknownThreatDetectStrategy("gone")
		if err == nil || !errs.NotFoundError(err) {
			t.Fatalf("a confirmed empty array with Code 404 must be NotFound, got %v", err)
		}
	})

	t.Run("forbiddenIsNotNotFound", func(t *testing.T) {
		apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusForbidden, map[string]interface{}{"Code": "APINotAuthorize", "Message": "denied"}
		})
		_, err := threatdetection.NewService(apiClient).DescribeUnknownThreatDetectStrategy("live")
		if err == nil {
			t.Fatal("403 must be returned as an error")
		}
		if errs.NotFoundError(err) {
			t.Fatalf("403 must not be NotFound: %v", err)
		}
	})

	t.Run("malformedContainersAreErrors", func(t *testing.T) {
		containers := map[string]interface{}{
			"null":   nil,
			"scalar": "not-an-array",
			"object": map[string]interface{}{},
		}
		for name, container := range containers {
			t.Run(name, func(t *testing.T) {
				apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
					return http.StatusOK, map[string]interface{}{"Data": container}
				})
				_, err := threatdetection.NewService(apiClient).DescribeUnknownThreatDetectStrategy("live")
				if err == nil {
					t.Fatalf("malformed container %v must be an error", container)
				}
				if errs.NotFoundError(err) {
					t.Fatalf("malformed container %v must not be NotFound: %v", container, err)
				}
			})
		}
	})
}

// ---- Q-8: read removes missing resource, keeps state on non-404 ----

func TestUnitUnknownThreatDetectStrategyReadRemovesMissing(t *testing.T) {
	apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusNotFound, map[string]interface{}{"Code": "404", "Message": "not found"}
	})
	r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
	resourceSchema, objType := unknownThreatDetectStrategySchemaAndType(t, r)
	stateValue := unknownThreatDetectStrategyValue(objType, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "gone"),
	})
	readReq := fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

	r.Read(context.Background(), readReq, readResp)

	if readResp.Diagnostics.HasError() {
		t.Fatalf("a typed 404 must be treated as NotFound: %v", readResp.Diagnostics)
	}
	if !readResp.State.Raw.IsNull() {
		t.Fatal("a missing strategy must be removed from state")
	}
}

func TestUnitUnknownThreatDetectStrategyReadKeepsStateOn403(t *testing.T) {
	apiClient, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
		return http.StatusForbidden, map[string]interface{}{"Code": "APINotAuthorize", "Message": "denied"}
	})
	r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
	resourceSchema, objType := unknownThreatDetectStrategySchemaAndType(t, r)
	stateValue := unknownThreatDetectStrategyValue(objType, map[string]tftypes.Value{
		"id": tftypes.NewValue(tftypes.String, "live"),
	})
	readResp := &fwresource.ReadResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
	r.Read(context.Background(), fwresource.ReadRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}, readResp)

	if !readResp.Diagnostics.HasError() {
		t.Fatal("a 403 must be returned as an error")
	}
	if readResp.State.Raw.IsNull() {
		t.Fatal("a failed describe must keep the strategy in state")
	}
}

// ---- Q-9: delete sends IdList and treats typed 404 as idempotent ----

func TestUnitUnknownThreatDetectStrategyDelete(t *testing.T) {
	cases := map[string]struct {
		status    int
		body      map[string]interface{}
		wantError bool
	}{
		"confirmed": {status: http.StatusOK, body: map[string]interface{}{"RequestId": "ok"}},
		"typed404":  {status: http.StatusNotFound, body: map[string]interface{}{"Code": "404", "Message": "not found"}},
		"forbidden": {status: http.StatusForbidden, body: map[string]interface{}{"Code": "denied"}, wantError: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			apiClient, recorder := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
				return c.status, c.body
			})
			r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
			resourceSchema, objType := unknownThreatDetectStrategySchemaAndType(t, r)
			stateValue := unknownThreatDetectStrategyValue(objType, map[string]tftypes.Value{
				"id": tftypes.NewValue(tftypes.String, "42"),
			})
			deleteReq := fwresource.DeleteRequest{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}
			deleteResp := &fwresource.DeleteResponse{State: tfsdk.State{Schema: resourceSchema, Raw: stateValue}}

			r.Delete(context.Background(), deleteReq, deleteResp)

			if c.wantError && !deleteResp.Diagnostics.HasError() {
				t.Fatal("expected a delete error")
			}
			if !c.wantError && deleteResp.Diagnostics.HasError() {
				t.Fatalf("Delete returned diagnostics: %v", deleteResp.Diagnostics)
			}
			forms := recorder.formsByAction("DeleteUnknownThreatDetectStrategy")
			if len(forms) != 1 || forms[0].Get("IdList.1") != "42" {
				t.Fatalf("unexpected delete request: %v", recorder.actionList())
			}
		})
	}
}

// ---- Q-10: update request shape ----

func TestUnitUnknownThreatDetectStrategyUpdateRequestShape(t *testing.T) {
	apiClient, recorder := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
		switch action {
		case "UpdateUnknownThreatDetectStrategy":
			return http.StatusOK, map[string]interface{}{"RequestId": "ok"}
		case "ListUnknownThreatDetectStrategy":
			return http.StatusOK, unknownThreatDetectStrategyStrategyResponse(float64(42), "updated-name", "UNKNOWN_THREAT_DETECT_CONFIG_test", map[string]interface{}{
				"StudyMode":             "path",
				"DurationDaysAfterInit": float64(7),
				"DurationDaysAfterStop": float64(9),
				"MachineCount":          float64(2),
			})
		}
		return http.StatusOK, map[string]interface{}{}
	})
	r := newUnknownThreatDetectStrategyUnitResource(t, apiClient)
	resourceSchema, objType := unknownThreatDetectStrategySchemaAndType(t, r)
	planValue := unknownThreatDetectStrategyValue(objType, map[string]tftypes.Value{
		"id":                       tftypes.NewValue(tftypes.String, "42"),
		"name":                     tftypes.NewValue(tftypes.String, "updated-name"),
		"asset_selection_type":     tftypes.NewValue(tftypes.String, "UNKNOWN_THREAT_DETECT_CONFIG_test"),
		"study_mode":               tftypes.NewValue(tftypes.String, "path"),
		"duration_days_after_init": tftypes.NewValue(tftypes.Number, int64(7)),
		"duration_days_after_stop": tftypes.NewValue(tftypes.Number, int64(9)),
	})
	req := fwresource.UpdateRequest{Plan: tfsdk.Plan{Schema: resourceSchema, Raw: planValue}}
	resp := &fwresource.UpdateResponse{State: tfsdk.State{Schema: resourceSchema, Raw: planValue}}

	r.Update(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Update returned diagnostics: %v", resp.Diagnostics)
	}
	forms := recorder.formsByAction("UpdateUnknownThreatDetectStrategy")
	if len(forms) != 1 {
		t.Fatalf("expected one Update call, got %v", recorder.actionList())
	}
	form := forms[0]
	if form.Get("Id") != "42" {
		t.Fatalf("Id = %q, want 42", form.Get("Id"))
	}
	if form.Get("Name") != "updated-name" {
		t.Fatalf("Name = %q", form.Get("Name"))
	}
	if form.Get("StudyMode") != "path" {
		t.Fatalf("StudyMode = %q", form.Get("StudyMode"))
	}
	if form.Get("DurationDaysAfterInit") != "7" {
		t.Fatalf("DurationDaysAfterInit = %q, want 7", form.Get("DurationDaysAfterInit"))
	}
	if got := unknownThreatDetectStrategyStateString(t, resp.State, "name"); got != "updated-name" {
		t.Fatalf("name state = %q", got)
	}
}

// ---- Q-11: check destroy ----

func TestUnitUnknownThreatDetectStrategyCheckDestroy(t *testing.T) {
	state := unknownThreatDetectStrategyDestroyState("live")

	t.Run("non404_fails", func(t *testing.T) {
		client, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusForbidden, map[string]interface{}{"Code": "denied"}
		})
		err := checkUnknownThreatDetectStrategyDestroyed(state, func(id string) error {
			_, describeErr := threatdetection.NewService(client).DescribeUnknownThreatDetectStrategy(id)
			return describeErr
		})
		if err == nil {
			t.Fatal("a non-404 describe failure must fail the cleanup check")
		}
	})

	t.Run("typed404_succeeds", func(t *testing.T) {
		client, _ := unknownThreatDetectStrategyTestClient(t, func(action string, form url.Values) (int, interface{}) {
			return http.StatusNotFound, map[string]interface{}{"Code": "404", "Message": "not found"}
		})
		err := checkUnknownThreatDetectStrategyDestroyed(state, func(id string) error {
			_, describeErr := threatdetection.NewService(client).DescribeUnknownThreatDetectStrategy(id)
			return describeErr
		})
		if err != nil {
			t.Fatalf("a typed 404 must count as destroyed: %s", err)
		}
	})
}

func unknownThreatDetectStrategyDestroyState(id string) *terraform.State {
	return &terraform.State{Modules: []*terraform.ModuleState{{
		Path: []string{"root"},
		Resources: map[string]*terraform.ResourceState{
			"alicloud_threat_detection_unknown_threat_detect_strategy.default": {
				Type:    "alicloud_threat_detection_unknown_threat_detect_strategy",
				Primary: &terraform.InstanceState{ID: id},
			},
		},
	}}}
}
