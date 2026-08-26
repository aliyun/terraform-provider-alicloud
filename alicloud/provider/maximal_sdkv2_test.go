// Direct-call coverage: does the SDK v2 wrapper drop anything the fixture holds?
package provider_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/sdkv2"
)

const (
	maximalSDKv2Name           = "alicloud_unit_test_maximal"
	maximalSDKv2ContextName    = "alicloud_unit_test_maximal_context"
	maximalSDKv2NoTimeoutName  = "alicloud_unit_test_maximal_no_timeout"
	maximalSDKv2DataSourceName = "alicloud_unit_test_maximal_data_source"
)

var crudFieldNames = map[string]bool{
	"Create": true, "CreateContext": true, "CreateWithoutTimeout": true,
	"Read": true, "ReadContext": true, "ReadWithoutTimeout": true,
	"Update": true, "UpdateContext": true, "UpdateWithoutTimeout": true,
	"Delete": true, "DeleteContext": true, "DeleteWithoutTimeout": true,
}

func maximalSDKv2SchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":     {Type: schema.TypeString, Required: true},
			"image_id": {Type: schema.TypeString, Optional: true},
		},
	}
}

func maximalSDKv2SchemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":     {Type: schema.TypeString, Required: true},
			"image_id": {Type: schema.TypeString, Optional: true},
			"tags":     {Type: schema.TypeMap, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}

func maximalSDKv2Resource(log *hookLog) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name":     {Type: schema.TypeString, Required: true},
			"image_id": {Type: schema.TypeString, Optional: true},
			"tags":     {Type: schema.TypeMap, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"status":   {Type: schema.TypeString, Computed: true},
		},
		SchemaVersion: 2,
		Identity: &schema.ResourceIdentity{
			Version: 1,
			SchemaFunc: func() map[string]*schema.Schema {
				return map[string]*schema.Schema{
					"name": {Type: schema.TypeString, RequiredForImport: true},
				}
			},
			IdentityUpgraders: []schema.IdentityUpgrader{
				{
					Version: 0,
					Type:    tftypes.Object{AttributeTypes: map[string]tftypes.Type{"name": tftypes.String}},
					Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
						log.record("IdentityUpgrade0")
						return rawState, nil
					},
				},
			},
		},
		MigrateState: func(version int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
			log.record("MigrateState")
			return is, nil
		},
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    maximalSDKv2SchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					log.record("StateUpgrade0")
					rawState["tags"] = map[string]interface{}{}
					return rawState, nil
				},
			},
			{
				Version: 1,
				Type:    maximalSDKv2SchemaV1().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					log.record("StateUpgrade1")
					rawState["status"] = "Running"
					return rawState, nil
				},
			},
		},
		Create: func(d *schema.ResourceData, meta interface{}) error {
			log.record("Create")
			d.SetId("unit-test")
			return setMaximalSDKv2Identity(d)
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			log.record("Read")
			return setMaximalSDKv2Identity(d)
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			log.record("Update")
			return setMaximalSDKv2Identity(d)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			log.record("Delete")
			return nil
		},
		Exists: func(d *schema.ResourceData, meta interface{}) (bool, error) {
			log.record("Exists")
			return true, nil
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			log.record("CustomizeDiff")
			if diff.HasChange("image_id") {
				return diff.ForceNew("image_id")
			}
			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				log.record("ImportState")
				return []*schema.ResourceData{d}, nil
			},
		},
		DeprecationMessage: "This is a unit test fixture and is never registered in the real provider.",
		Timeouts: &schema.ResourceTimeout{
			Create:  durationPtr(11 * time.Minute),
			Read:    durationPtr(12 * time.Minute),
			Update:  durationPtr(13 * time.Minute),
			Delete:  durationPtr(14 * time.Minute),
			Default: durationPtr(15 * time.Minute),
		},
		Description:                       "Feature-maximal SDK v2 fixture for the interceptor layer.",
		UseJSONNumber:                     true,
		EnableLegacyTypeSystemApplyErrors: true,
		EnableLegacyTypeSystemPlanErrors:  true,
		ResourceBehavior:                  schema.ResourceBehavior{MutableIdentity: true},
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			func(ctx context.Context, req schema.ValidateResourceConfigFuncRequest, resp *schema.ValidateResourceConfigFuncResponse) {
				log.record("ValidateRawResourceConfig")
				resp.Diagnostics = append(resp.Diagnostics, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "the raw config validator ran",
				})
			},
		},
	}
}

func maximalSDKv2ContextResource(log *hookLog) *schema.Resource {
	return &schema.Resource{
		SchemaFunc: func() map[string]*schema.Schema {
			return map[string]*schema.Schema{
				"name": {Type: schema.TypeString, Required: true},
			}
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    maximalSDKv2SchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
					log.record("ContextStateUpgrade0")
					return rawState, nil
				},
			},
		},
		CreateContext: recordContext(log, "ContextCreate"),
		ReadContext:   recordContext(log, "ContextRead"),
		UpdateContext: recordContext(log, "ContextUpdate"),
		DeleteContext: recordContext(log, "ContextDelete"),
	}
}

func maximalSDKv2WithoutTimeoutResource(log *hookLog) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {Type: schema.TypeString, Required: true},
		},
		CreateWithoutTimeout: recordContext(log, "NoTimeoutCreate"),
		ReadWithoutTimeout:   recordContext(log, "NoTimeoutRead"),
		UpdateWithoutTimeout: recordContext(log, "NoTimeoutUpdate"),
		DeleteWithoutTimeout: recordContext(log, "NoTimeoutDelete"),
	}
}

func maximalSDKv2DataSource(log *hookLog) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) ([]string, []error) {
					log.record("ValidateName")
					return []string{"the data source config validator ran"}, nil
				},
			},
			"image_id": {Type: schema.TypeString, Optional: true},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":     {Type: schema.TypeString, Computed: true},
						"status": {Type: schema.TypeString, Computed: true},
					},
				},
			},
			"output_file": {Type: schema.TypeString, Optional: true},
		},
		DeprecationMessage: "This is a unit test fixture and is never registered in the real provider.",
		Description:        "Feature-maximal SDK v2 data source fixture for the interceptor layer.",
		Read: func(d *schema.ResourceData, meta interface{}) error {
			log.record("Read")
			d.SetId("unit-test")
			if err := d.Set("ids", []string{"unit-test"}); err != nil {
				return err
			}
			return d.Set("instances", []map[string]interface{}{
				{"id": "unit-test", "status": "Running"},
			})
		},
	}
}

func setMaximalSDKv2Identity(d *schema.ResourceData) error {
	identity, err := d.Identity()
	if err != nil {
		return err
	}
	return identity.Set("name", d.Get("name"))
}

func recordContext(log *hookLog, name string) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		log.record(name)
		d.SetId("unit-test")
		return nil
	}
}

func durationPtr(d time.Duration) *time.Duration { return &d }

func maximalSDKv2Provider(log *hookLog, chain []intercept.Interceptor) *schema.Provider {
	resources := map[string]*schema.Resource{
		maximalSDKv2Name:          maximalSDKv2Resource(log),
		maximalSDKv2ContextName:   maximalSDKv2ContextResource(log),
		maximalSDKv2NoTimeoutName: maximalSDKv2WithoutTimeoutResource(log),
	}
	dataSources := map[string]*schema.Resource{
		maximalSDKv2DataSourceName: maximalSDKv2DataSource(log),
	}
	if len(chain) > 0 {
		for name, r := range resources {
			resources[name] = sdkv2.WrapResource(name, r, chain)
		}
		for name, ds := range dataSources {
			dataSources[name] = sdkv2.WrapDataSource(name, ds, chain)
		}
	}
	return &schema.Provider{ResourcesMap: resources, DataSourcesMap: dataSources}
}

func TestUnitProviderMaximalSDKv2FixturesAreLegal(t *testing.T) {
	log := &hookLog{}

	if err := maximalSDKv2Provider(log, nil).InternalValidate(); err != nil {
		t.Fatalf("the unwrapped fixtures are not legal resources: %s", err)
	}
	if err := maximalSDKv2Provider(log, []intercept.Interceptor{&chainRecorder{}}).InternalValidate(); err != nil {
		t.Fatalf("wrapping produced an illegal resource: %s", err)
	}
}

func TestUnitProviderMaximalSDKv2FixturesCoverEveryField(t *testing.T) {
	log := &hookLog{}
	fixtures := []*schema.Resource{
		maximalSDKv2Resource(log),
		maximalSDKv2ContextResource(log),
		maximalSDKv2WithoutTimeoutResource(log),
	}

	rt := reflect.TypeOf(schema.Resource{})
	var uncovered []string
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}
		covered := false
		for _, fixture := range fixtures {
			if !reflect.ValueOf(*fixture).Field(i).IsZero() {
				covered = true
				break
			}
		}
		if !covered {
			uncovered = append(uncovered, f.Name)
		}
	}

	if len(uncovered) > 0 {
		t.Fatalf("no fixture populates %v: the wrapper sweep cannot observe a zero field, so add these to a fixture", uncovered)
	}
}

func TestUnitProviderMaximalSDKv2WrapCopiesEveryNonCRUDField(t *testing.T) {
	log := &hookLog{}
	chain := []intercept.Interceptor{&chainRecorder{}}

	for _, tc := range []struct {
		name  string
		build func(*hookLog) *schema.Resource
		wrap  func(string, *schema.Resource, []intercept.Interceptor) *schema.Resource
	}{
		{maximalSDKv2Name, maximalSDKv2Resource, sdkv2.WrapResource},
		{maximalSDKv2ContextName, maximalSDKv2ContextResource, sdkv2.WrapResource},
		{maximalSDKv2NoTimeoutName, maximalSDKv2WithoutTimeoutResource, sdkv2.WrapResource},
		{maximalSDKv2DataSourceName, maximalSDKv2DataSource, sdkv2.WrapDataSource},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := tc.build(log)
			wrapped := tc.wrap(tc.name, original, chain)

			if wrapped == original {
				t.Fatal("wrapping returned the input pointer: the shared map literal would be mutated")
			}

			ov, wv := reflect.ValueOf(*original), reflect.ValueOf(*wrapped)
			rt := ov.Type()
			for i := 0; i < rt.NumField(); i++ {
				f := rt.Field(i)
				if !f.IsExported() || crudFieldNames[f.Name] {
					continue
				}
				of, wf := ov.Field(i), wv.Field(i)

				if f.Type.Kind() == reflect.Func {
					if of.Pointer() != wf.Pointer() {
						t.Errorf("%s was not copied", f.Name)
					}
					continue
				}
				if !reflect.DeepEqual(of.Interface(), wf.Interface()) {
					t.Errorf("%s was not copied: original %v, wrapped %v", f.Name, of.Interface(), wf.Interface())
				}
			}
		})
	}
}

func TestUnitProviderMaximalSDKv2WrapPreservesEveryHook(t *testing.T) {
	ctx := context.Background()
	log := &hookLog{}
	wrapped := sdkv2.WrapResource(maximalSDKv2Name, maximalSDKv2Resource(log), []intercept.Interceptor{&chainRecorder{}})

	t.Run("CustomizeDiff", func(t *testing.T) {
		state := &terraform.InstanceState{
			ID: "unit-test",
			Attributes: map[string]string{
				"id":       "unit-test",
				"name":     "tf-testacc-maximal",
				"image_id": "image-one",
			},
		}
		config := terraform.NewResourceConfigRaw(map[string]interface{}{
			"name":     "tf-testacc-maximal",
			"image_id": "image-two",
		})

		diff, err := wrapped.Diff(ctx, state, config, nil)
		if err != nil {
			t.Fatalf("Diff: %s", err)
		}
		if !log.has("CustomizeDiff") {
			t.Fatal("CustomizeDiff did not run on the wrapped resource")
		}
		if diff == nil || diff.Attributes["image_id"] == nil || !diff.Attributes["image_id"].RequiresNew {
			t.Fatalf("CustomizeDiff ran but its ForceNew was lost: %#v", diff)
		}
	})

	t.Run("StateUpgraders", func(t *testing.T) {
		if len(wrapped.StateUpgraders) != 2 {
			t.Fatalf("got %d state upgraders, want 2", len(wrapped.StateUpgraders))
		}
		raw := map[string]interface{}{"name": "tf-testacc-maximal", "image_id": "image-one"}
		for i, u := range wrapped.StateUpgraders {
			out, err := u.Upgrade(ctx, raw, nil)
			if err != nil {
				t.Fatalf("upgrader %d: %s", i, err)
			}
			raw = out
		}
		if !log.has("StateUpgrade0") || !log.has("StateUpgrade1") {
			t.Fatalf("an upgrader did not run: %v", log.all())
		}
		if _, ok := raw["tags"]; !ok {
			t.Error("the v0 upgrader's output was discarded")
		}
		if raw["status"] != "Running" {
			t.Error("the v1 upgrader's output was discarded")
		}
	})

	t.Run("MigrateState", func(t *testing.T) {
		if wrapped.MigrateState == nil {
			t.Fatal("MigrateState was dropped")
		}
		in := &terraform.InstanceState{ID: "unit-test"}
		out, err := wrapped.MigrateState(0, in, nil)
		if err != nil {
			t.Fatalf("MigrateState: %s", err)
		}
		if out != in {
			t.Error("MigrateState's return value was not the fixture's")
		}
		if !log.has("MigrateState") {
			t.Fatal("MigrateState did not run")
		}
	})

	t.Run("Identity", func(t *testing.T) {
		if wrapped.Identity == nil {
			t.Fatal("Identity was dropped")
		}
		if got := wrapped.Identity.Version; got != 1 {
			t.Errorf("Identity.Version = %d, want 1", got)
		}
		if _, ok := wrapped.Identity.SchemaMap()["name"]; !ok {
			t.Error("Identity.SchemaFunc was dropped")
		}
		if len(wrapped.Identity.IdentityUpgraders) != 1 {
			t.Fatalf("got %d identity upgraders, want 1", len(wrapped.Identity.IdentityUpgraders))
		}
		if _, err := wrapped.Identity.IdentityUpgraders[0].Upgrade(ctx, map[string]interface{}{"name": "tf-testacc-maximal"}, nil); err != nil {
			t.Fatalf("identity upgrader: %s", err)
		}
		if !log.has("IdentityUpgrade0") {
			t.Fatal("the identity upgrader did not run")
		}
	})

	t.Run("Importer", func(t *testing.T) {
		if wrapped.Importer == nil || wrapped.Importer.StateContext == nil {
			t.Fatal("Importer was dropped")
		}
		d := wrapped.Data(&terraform.InstanceState{ID: "unit-test"})
		out, err := wrapped.Importer.StateContext(ctx, d, nil)
		if err != nil {
			t.Fatalf("ImportState: %s", err)
		}
		if len(out) != 1 || out[0] != d {
			t.Error("the importer's return value was not the fixture's")
		}
		if !log.has("ImportState") {
			t.Fatal("the importer did not run")
		}
	})

	t.Run("Exists", func(t *testing.T) {
		if wrapped.Exists == nil {
			t.Fatal("Exists was dropped")
		}
		ok, err := wrapped.Exists(wrapped.Data(&terraform.InstanceState{ID: "unit-test"}), nil)
		if err != nil || !ok {
			t.Fatalf("Exists returned (%t, %v), want (true, nil)", ok, err)
		}
		if !log.has("Exists") {
			t.Fatal("Exists did not run")
		}
	})

	t.Run("ValidateRawResourceConfigFuncs", func(t *testing.T) {
		if len(wrapped.ValidateRawResourceConfigFuncs) != 1 {
			t.Fatalf("got %d raw config validators, want 1", len(wrapped.ValidateRawResourceConfigFuncs))
		}
		var resp schema.ValidateResourceConfigFuncResponse
		wrapped.ValidateRawResourceConfigFuncs[0](ctx, schema.ValidateResourceConfigFuncRequest{}, &resp)
		if !log.has("ValidateRawResourceConfig") {
			t.Fatal("the raw config validator did not run")
		}
		if len(resp.Diagnostics) != 1 {
			t.Errorf("the validator's diagnostics were lost: %v", resp.Diagnostics)
		}
	})

	t.Run("scalar fields", func(t *testing.T) {
		if wrapped.SchemaVersion != 2 {
			t.Errorf("SchemaVersion = %d, want 2", wrapped.SchemaVersion)
		}
		if wrapped.DeprecationMessage == "" || wrapped.Description == "" {
			t.Error("DeprecationMessage or Description was dropped")
		}
		if !wrapped.UseJSONNumber || !wrapped.EnableLegacyTypeSystemApplyErrors || !wrapped.EnableLegacyTypeSystemPlanErrors {
			t.Error("a behaviour flag was dropped")
		}
		if !wrapped.ResourceBehavior.MutableIdentity {
			t.Error("ResourceBehavior was dropped")
		}
		if wrapped.Timeouts == nil || wrapped.Timeouts.Create == nil || *wrapped.Timeouts.Create != 11*time.Minute {
			t.Error("Timeouts was dropped")
		}
	})
}

func TestUnitProviderMaximalSDKv2WrapRunsChainOnlyOnCRUD(t *testing.T) {
	ctx := context.Background()
	log := &hookLog{}
	rec := &chainRecorder{}
	wrapped := sdkv2.WrapResource(maximalSDKv2Name, maximalSDKv2Resource(log), []intercept.Interceptor{rec})

	d := wrapped.Data(&terraform.InstanceState{ID: "unit-test"})
	if _, err := wrapped.Importer.StateContext(ctx, d, nil); err != nil {
		t.Fatalf("ImportState: %s", err)
	}
	if _, err := wrapped.Exists(d, nil); err != nil {
		t.Fatalf("Exists: %s", err)
	}
	if _, err := wrapped.MigrateState(0, &terraform.InstanceState{ID: "unit-test"}, nil); err != nil {
		t.Fatalf("MigrateState: %s", err)
	}
	if _, err := wrapped.StateUpgraders[0].Upgrade(ctx, map[string]interface{}{}, nil); err != nil {
		t.Fatalf("StateUpgrade0: %s", err)
	}
	var validateResp schema.ValidateResourceConfigFuncResponse
	wrapped.ValidateRawResourceConfigFuncs[0](ctx, schema.ValidateResourceConfigFuncRequest{}, &validateResp)

	if got := rec.snapshotBefore(); len(got) != 0 {
		t.Fatalf("the chain ran outside CRUD: %v", got)
	}

	client := &struct{ name string }{name: "stand-in for the client"}
	if err := wrapped.Create(d, client); err != nil {
		t.Fatalf("Create: %s", err)
	}
	if err := wrapped.Read(d, client); err != nil {
		t.Fatalf("Read: %s", err)
	}
	if err := wrapped.Update(d, client); err != nil {
		t.Fatalf("Update: %s", err)
	}
	if err := wrapped.Delete(d, client); err != nil {
		t.Fatalf("Delete: %s", err)
	}

	wantOps := []string{"Create", "Read", "Update", "Delete"}
	if got := rec.beforeOps(); !reflect.DeepEqual(got, wantOps) {
		t.Errorf("Before ops = %v, want %v", got, wantOps)
	}
	if got := rec.afterOps(); !reflect.DeepEqual(got, wantOps) {
		t.Errorf("After ops = %v, want %v", got, wantOps)
	}
	for _, name := range namesOf(rec.snapshotBefore()) {
		if name != maximalSDKv2Name {
			t.Errorf("the chain saw name %q, want %q", name, maximalSDKv2Name)
		}
	}
	for _, call := range rec.snapshotBefore() {
		if call.Meta != interface{}(client) {
			t.Errorf("Call.Meta = %v, want the meta the CRUD function was given", call.Meta)
		}
	}
	if got := log.all(); !reflect.DeepEqual(got[len(got)-4:], wantOps) {
		t.Errorf("the inner CRUD functions ran as %v, want %v", got[len(got)-4:], wantOps)
	}
}

func TestUnitProviderMaximalSDKv2WrapCoversEveryCRUDForm(t *testing.T) {
	ctx := context.Background()
	wantOps := []string{"Create", "Read", "Update", "Delete"}

	t.Run("context form", func(t *testing.T) {
		log := &hookLog{}
		rec := &chainRecorder{}
		wrapped := sdkv2.WrapResource(maximalSDKv2ContextName, maximalSDKv2ContextResource(log), []intercept.Interceptor{rec})

		if wrapped.Create != nil || wrapped.CreateWithoutTimeout != nil {
			t.Error("wrapping moved the create function to another declaration form")
		}
		d := wrapped.Data(&terraform.InstanceState{ID: "unit-test"})
		for _, call := range []func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics{
			wrapped.CreateContext, wrapped.ReadContext, wrapped.UpdateContext, wrapped.DeleteContext,
		} {
			if diags := call(ctx, d, nil); diags.HasError() {
				t.Fatalf("unexpected diags: %v", diags)
			}
		}

		if got := rec.beforeOps(); !reflect.DeepEqual(got, wantOps) {
			t.Errorf("Before ops = %v, want %v", got, wantOps)
		}
		want := []string{"ContextCreate", "ContextRead", "ContextUpdate", "ContextDelete"}
		if got := log.all(); !reflect.DeepEqual(got, want) {
			t.Errorf("inner functions ran as %v, want %v", got, want)
		}
	})

	t.Run("without-timeout form", func(t *testing.T) {
		log := &hookLog{}
		rec := &chainRecorder{}
		wrapped := sdkv2.WrapResource(maximalSDKv2NoTimeoutName, maximalSDKv2WithoutTimeoutResource(log), []intercept.Interceptor{rec})

		if wrapped.Create != nil || wrapped.CreateContext != nil {
			t.Error("wrapping moved the create function to another declaration form")
		}
		d := wrapped.Data(&terraform.InstanceState{ID: "unit-test"})
		for _, call := range []func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics{
			wrapped.CreateWithoutTimeout, wrapped.ReadWithoutTimeout, wrapped.UpdateWithoutTimeout, wrapped.DeleteWithoutTimeout,
		} {
			if diags := call(ctx, d, nil); diags.HasError() {
				t.Fatalf("unexpected diags: %v", diags)
			}
		}

		if got := rec.beforeOps(); !reflect.DeepEqual(got, wantOps) {
			t.Errorf("Before ops = %v, want %v", got, wantOps)
		}
		want := []string{"NoTimeoutCreate", "NoTimeoutRead", "NoTimeoutUpdate", "NoTimeoutDelete"}
		if got := log.all(); !reflect.DeepEqual(got, want) {
			t.Errorf("inner functions ran as %v, want %v", got, want)
		}
	})
}
