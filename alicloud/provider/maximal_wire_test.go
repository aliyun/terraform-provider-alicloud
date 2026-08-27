// The two fixtures driven through real protocol-v5 servers — the only place
// non-CRUD RPCs are testable. Not muxed: factory_test.go covers that.
package provider_test

import (
	"context"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
)

var (
	maximalSDKv2Type = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":       tftypes.String,
		"name":     tftypes.String,
		"image_id": tftypes.String,
		"status":   tftypes.String,
		"tags":     tftypes.Map{ElementType: tftypes.String},
		"timeouts": tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"create":  tftypes.String,
			"read":    tftypes.String,
			"update":  tftypes.String,
			"delete":  tftypes.String,
			"default": tftypes.String,
		}},
	}}
	maximalSDKv2IdentityType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"name": tftypes.String,
	}}
	maximalSDKv2DataSourceType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":       tftypes.String,
		"name":     tftypes.String,
		"image_id": tftypes.String,
		"ids":      tftypes.List{ElementType: tftypes.String},
		"instances": tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"id":     tftypes.String,
			"status": tftypes.String,
		}}},
		"output_file": tftypes.String,
	}}
	maximalFrameworkType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id":   tftypes.String,
		"name": tftypes.String,
	}}
	maximalFrameworkIdentityType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"id": tftypes.String,
	}}
	emptyObjectType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
)

var wantCRUDOps = []string{"Create", "Read", "Update", "Delete"}

type maximalFrameworkTestProvider struct {
	log    *hookLog
	chain  []intercept.Interceptor
	client *connectivity.AliyunClient
}

var _ fwprovider.Provider = (*maximalFrameworkTestProvider)(nil)

func (p *maximalFrameworkTestProvider) Metadata(ctx context.Context, req fwprovider.MetadataRequest, resp *fwprovider.MetadataResponse) {
	resp.TypeName = "alicloud"
}

func (p *maximalFrameworkTestProvider) Schema(ctx context.Context, req fwprovider.SchemaRequest, resp *fwprovider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *maximalFrameworkTestProvider) Configure(ctx context.Context, req fwprovider.ConfigureRequest, resp *fwprovider.ConfigureResponse) {
	resp.ResourceData = p.client
	resp.DataSourceData = p.client
}

func (p *maximalFrameworkTestProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return fwadapt.WrapResource(maximalFrameworkName, newMaximalFrameworkResource(p.log), p.chain)
		},
	}
}

func (p *maximalFrameworkTestProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return fwadapt.WrapDataSource(maximalFrameworkDataSourceName, newMaximalFrameworkDataSource(p.log), p.chain)
		},
	}
}

func TestUnitProviderMaximalWireSDKv2(t *testing.T) {
	ctx := context.Background()
	log, chain, client := &hookLog{}, &chainRecorder{}, &connectivity.AliyunClient{}
	server := newMaximalSDKv2Server(t, log, chain, client)

	t.Run("GetProviderSchema", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
		if err != nil {
			t.Fatalf("GetProviderSchema: %s", err)
		}
		requireNoDiagErrors(t, "GetProviderSchema", resp.Diagnostics)

		schema, ok := resp.ResourceSchemas[maximalSDKv2Name]
		if !ok {
			t.Fatalf("the wrapped fixture is not in the schema: %v", resp.ResourceSchemas)
		}
		if schema.Version != 2 {
			t.Errorf("the served schema version is %d, want 2", schema.Version)
		}
		if !schema.Block.Deprecated {
			t.Error("the served schema is not marked deprecated")
		}

		dataSchema, ok := resp.DataSourceSchemas[maximalSDKv2DataSourceName]
		if !ok {
			t.Fatalf("the wrapped data source fixture is not in the schema: %v", resp.DataSourceSchemas)
		}
		if got := dataSchema.ValueType(); !got.Equal(maximalSDKv2DataSourceType) {
			t.Errorf("the served data source type is %s, want %s", got, maximalSDKv2DataSourceType)
		}
	})

	t.Run("GetResourceIdentitySchemas", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.GetResourceIdentitySchemas(ctx, &tfprotov5.GetResourceIdentitySchemasRequest{})
		if err != nil {
			t.Fatalf("GetResourceIdentitySchemas: %s", err)
		}
		requireNoDiagErrors(t, "GetResourceIdentitySchemas", resp.Diagnostics)

		schema, ok := resp.IdentitySchemas[maximalSDKv2Name]
		if !ok {
			t.Fatalf("the wrapped fixture has no identity schema: %v", resp.IdentitySchemas)
		}
		if schema.Version != 1 {
			t.Errorf("the served identity version is %d, want 1", schema.Version)
		}
		if len(schema.IdentityAttributes) != 1 {
			t.Fatalf("the served identity has %d attributes, want 1", len(schema.IdentityAttributes))
		}
		if got := schema.IdentityAttributes[0]; got.Name != "name" || !got.RequiredForImport {
			t.Errorf("the served identity attribute is %+v, want name/RequiredForImport", got)
		}
	})

	t.Run("ValidateResourceTypeConfig", func(t *testing.T) {
		log.reset()
		chain.reset()

		config := dynamic(t, maximalSDKv2Type, object(t, maximalSDKv2Type, map[string]tftypes.Value{
			"name": tftypes.NewValue(tftypes.String, "fixture"),
		}))
		resp, err := server.ValidateResourceTypeConfig(ctx, &tfprotov5.ValidateResourceTypeConfigRequest{
			TypeName: maximalSDKv2Name,
			Config:   config,
		})
		if err != nil {
			t.Fatalf("ValidateResourceTypeConfig: %s", err)
		}
		requireNoDiagErrors(t, "ValidateResourceTypeConfig", resp.Diagnostics)

		if !log.has("ValidateRawResourceConfig") {
			t.Error("the raw config validator did not run")
		}
		if !hasWarning(resp.Diagnostics, "the raw config validator ran") {
			t.Errorf("the validator's warning did not survive: %v", resp.Diagnostics)
		}
		assertChainSilent(t, chain, "validating a config")
	})

	t.Run("UpgradeResourceState/json", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.UpgradeResourceState(ctx, &tfprotov5.UpgradeResourceStateRequest{
			TypeName: maximalSDKv2Name,
			Version:  0,
			RawState: &tfprotov5.RawState{
				JSON: []byte(`{"id":"unit-test","name":"fixture","image_id":"image-one"}`),
			},
		})
		if err != nil {
			t.Fatalf("UpgradeResourceState: %s", err)
		}
		requireNoDiagErrors(t, "UpgradeResourceState", resp.Diagnostics)

		if !log.has("StateUpgrade0") || !log.has("StateUpgrade1") {
			t.Errorf("the upgrader chain did not run to completion: %v", log.all())
		}
		attrs := attrsOf(t, maximalSDKv2Type, resp.UpgradedState)
		if got := stringOf(t, attrs["status"]); got != "Running" {
			t.Errorf("status after upgrading is %q, want %q — the second upgrader's write was lost", got, "Running")
		}
		if attrs["tags"].IsNull() {
			t.Error("tags after upgrading is null — the first upgrader's write was lost")
		}
		assertChainSilent(t, chain, "upgrading a JSON state")
	})

	t.Run("UpgradeResourceState/flatmap", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.UpgradeResourceState(ctx, &tfprotov5.UpgradeResourceStateRequest{
			TypeName: maximalSDKv2Name,
			Version:  0,
			RawState: &tfprotov5.RawState{
				Flatmap: map[string]string{
					"id":       "unit-test",
					"name":     "fixture",
					"image_id": "image-one",
				},
			},
		})
		if err != nil {
			t.Fatalf("UpgradeResourceState: %s", err)
		}
		requireNoDiagErrors(t, "UpgradeResourceState", resp.Diagnostics)

		if !log.has("StateUpgrade0") || !log.has("StateUpgrade1") {
			t.Errorf("the upgrader chain did not run to completion: %v", log.all())
		}
		if log.has("MigrateState") {
			t.Error("MigrateState ran, but no version is low enough to reach it")
		}
		if got := stringOf(t, attrsOf(t, maximalSDKv2Type, resp.UpgradedState)["status"]); got != "Running" {
			t.Errorf("status after upgrading a flatmap is %q, want %q", got, "Running")
		}
		assertChainSilent(t, chain, "upgrading a flatmap state")
	})

	t.Run("UpgradeResourceIdentity", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.UpgradeResourceIdentity(ctx, &tfprotov5.UpgradeResourceIdentityRequest{
			TypeName:    maximalSDKv2Name,
			Version:     0,
			RawIdentity: &tfprotov5.RawState{JSON: []byte(`{"name":"fixture"}`)},
		})
		if err != nil {
			t.Fatalf("UpgradeResourceIdentity: %s", err)
		}
		requireNoDiagErrors(t, "UpgradeResourceIdentity", resp.Diagnostics)

		if !log.has("IdentityUpgrade0") {
			t.Errorf("the identity upgrader did not run: %v", log.all())
		}
		if resp.UpgradedIdentity == nil {
			t.Fatal("the upgraded identity is missing")
		}
		if got := stringOf(t, attrsOf(t, maximalSDKv2IdentityType, resp.UpgradedIdentity.IdentityData)["name"]); got != "fixture" {
			t.Errorf("the upgraded identity name is %q, want %q", got, "fixture")
		}
		assertChainSilent(t, chain, "upgrading an identity")
	})

	t.Run("ImportResourceState", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.ImportResourceState(ctx, &tfprotov5.ImportResourceStateRequest{
			TypeName: maximalSDKv2Name,
			ID:       "unit-test",
		})
		if err != nil {
			t.Fatalf("ImportResourceState: %s", err)
		}
		requireNoDiagErrors(t, "ImportResourceState", resp.Diagnostics)

		if len(resp.ImportedResources) != 1 {
			t.Fatalf("importing produced %d resources, want 1", len(resp.ImportedResources))
		}
		if !log.has("ImportState") {
			t.Errorf("the importer did not run: %v", log.all())
		}
		if log.has("Read") {
			t.Errorf("importing also ran Read: %v", log.all())
		}
		assertChainSilent(t, chain, "importing")
	})

	t.Run("ValidateDataSourceConfig", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.ValidateDataSourceConfig(ctx, &tfprotov5.ValidateDataSourceConfigRequest{
			TypeName: maximalSDKv2DataSourceName,
			Config: dynamic(t, maximalSDKv2DataSourceType, object(t, maximalSDKv2DataSourceType, map[string]tftypes.Value{
				"name": tftypes.NewValue(tftypes.String, "fixture"),
			})),
		})
		if err != nil {
			t.Fatalf("ValidateDataSourceConfig: %s", err)
		}
		requireNoDiagErrors(t, "ValidateDataSourceConfig", resp.Diagnostics)

		if !log.has("ValidateName") {
			t.Errorf("the schema validator did not run: %v", log.all())
		}
		if !hasWarning(resp.Diagnostics, "the data source config validator ran") {
			t.Errorf("the validator's warning did not survive: %v", resp.Diagnostics)
		}
		if !hasWarning(resp.Diagnostics, "Deprecated Resource") {
			t.Errorf("the deprecation warning did not survive: %v", resp.Diagnostics)
		}
		if log.has("Read") {
			t.Errorf("validation also ran Read: %v", log.all())
		}
		assertChainSilent(t, chain, "validating a data source config")
	})

	t.Run("ReadDataSource", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.ReadDataSource(ctx, &tfprotov5.ReadDataSourceRequest{
			TypeName: maximalSDKv2DataSourceName,
			Config: dynamic(t, maximalSDKv2DataSourceType, object(t, maximalSDKv2DataSourceType, map[string]tftypes.Value{
				"name": tftypes.NewValue(tftypes.String, "fixture"),
			})),
		})
		if err != nil {
			t.Fatalf("ReadDataSource: %s", err)
		}
		requireNoDiagErrors(t, "ReadDataSource", resp.Diagnostics)

		attrs := attrsOf(t, maximalSDKv2DataSourceType, resp.State)
		if got := stringOf(t, attrs["id"]); got != "unit-test" {
			t.Errorf("the data source produced id %q, want %q", got, "unit-test")
		}
		var ids []tftypes.Value
		if err := attrs["ids"].As(&ids); err != nil {
			t.Fatalf("reading ids: %s", err)
		}
		if len(ids) != 1 || stringOf(t, ids[0]) != "unit-test" {
			t.Errorf("the data source produced ids %v, want one element %q", ids, "unit-test")
		}
		if attrs["instances"].IsNull() {
			t.Error("the nested list is null: the schema the read wrote to is not the one being served")
		}
		if got := chain.beforeOps(); !slices.Equal(got, []string{"Read"}) {
			t.Errorf("the chain saw %v for a data source read, want [Read]", got)
		}
		if got := chain.afterOps(); !slices.Equal(got, []string{"Read"}) {
			t.Errorf("the chain saw %v after a data source read, want [Read]", got)
		}
		for _, call := range chain.snapshotBefore() {
			if call.Name != maximalSDKv2DataSourceName {
				t.Errorf("the read was described as %q, want %q", call.Name, maximalSDKv2DataSourceName)
			}
			if call.Meta != client {
				t.Errorf("the read carried meta %v, want the configured client", call.Meta)
			}
		}
	})

	t.Run("lifecycle", func(t *testing.T) {
		log.reset()
		chain.reset()

		runMaximalSDKv2Lifecycle(t, server)

		assertChainSawLifecycle(t, chain, maximalSDKv2Name, client)
		if !log.has("CustomizeDiff") {
			t.Errorf("CustomizeDiff never ran during a full lifecycle: %v", log.all())
		}
		if !log.has("Exists") {
			t.Errorf("Exists never ran during a full lifecycle: %v", log.all())
		}
	})
}

func TestUnitProviderMaximalWireFramework(t *testing.T) {
	ctx := context.Background()
	log, chain, client := &hookLog{}, &chainRecorder{}, &connectivity.AliyunClient{}
	server := newMaximalFrameworkServer(t, log, chain, client)

	t.Run("GetProviderSchema", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
		if err != nil {
			t.Fatalf("GetProviderSchema: %s", err)
		}
		requireNoDiagErrors(t, "GetProviderSchema", resp.Diagnostics)

		schema, ok := resp.ResourceSchemas[maximalFrameworkName]
		if !ok {
			t.Fatalf("the wrapped fixture is not in the schema: %v", resp.ResourceSchemas)
		}
		if schema.Version != 2 {
			t.Errorf("the served schema version is %d, want 2", schema.Version)
		}
		if _, ok := resp.DataSourceSchemas[maximalFrameworkDataSourceName]; !ok {
			t.Errorf("the wrapped data source is not in the schema: %v", resp.DataSourceSchemas)
		}
	})

	t.Run("GetResourceIdentitySchemas", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.GetResourceIdentitySchemas(ctx, &tfprotov5.GetResourceIdentitySchemasRequest{})
		if err != nil {
			t.Fatalf("GetResourceIdentitySchemas: %s", err)
		}
		requireNoDiagErrors(t, "GetResourceIdentitySchemas", resp.Diagnostics)

		schema, ok := resp.IdentitySchemas[maximalFrameworkName]
		if !ok {
			t.Fatalf("the wrapped fixture has no identity schema: %v", resp.IdentitySchemas)
		}
		if schema.Version != 1 {
			t.Errorf("the served identity version is %d, want 1", schema.Version)
		}
		if len(schema.IdentityAttributes) != 1 {
			t.Fatalf("the served identity has %d attributes, want 1", len(schema.IdentityAttributes))
		}
		if got := schema.IdentityAttributes[0]; got.Name != "id" || !got.RequiredForImport {
			t.Errorf("the served identity attribute is %+v, want id/RequiredForImport", got)
		}
	})

	t.Run("ValidateResourceTypeConfig", func(t *testing.T) {
		log.reset()
		chain.reset()

		config := dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, map[string]tftypes.Value{
			"name": tftypes.NewValue(tftypes.String, "fixture"),
		}))
		resp, err := server.ValidateResourceTypeConfig(ctx, &tfprotov5.ValidateResourceTypeConfigRequest{
			TypeName: maximalFrameworkName,
			Config:   config,
		})
		if err != nil {
			t.Fatalf("ValidateResourceTypeConfig: %s", err)
		}
		requireNoDiagErrors(t, "ValidateResourceTypeConfig", resp.Diagnostics)

		for _, hook := range []string{"ConfigValidators", "ConfigValidator", "ValidateConfig"} {
			if !log.has(hook) {
				t.Errorf("%s did not run while validating: %v", hook, log.all())
			}
		}
		assertChainSilent(t, chain, "validating a config")
	})

	t.Run("UpgradeResourceState", func(t *testing.T) {
		log.reset()
		chain.reset()

		for version := int64(0); version <= 1; version++ {
			resp, err := server.UpgradeResourceState(ctx, &tfprotov5.UpgradeResourceStateRequest{
				TypeName: maximalFrameworkName,
				Version:  version,
				RawState: &tfprotov5.RawState{
					JSON: []byte(`{"id":"unit-test","name":"fixture"}`),
				},
			})
			if err != nil {
				t.Fatalf("UpgradeResourceState from version %d: %s", version, err)
			}
			requireNoDiagErrors(t, "UpgradeResourceState", resp.Diagnostics)

			if got := stringOf(t, attrsOf(t, maximalFrameworkType, resp.UpgradedState)["id"]); got != "unit-test" {
				t.Errorf("the state upgraded from version %d has id %q, want %q", version, got, "unit-test")
			}
		}
		if !log.has("UpgradeState0") || !log.has("UpgradeState1") {
			t.Errorf("not every upgrader was reachable: %v", log.all())
		}
		assertChainSilent(t, chain, "upgrading a state")
	})

	t.Run("UpgradeResourceIdentity", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.UpgradeResourceIdentity(ctx, &tfprotov5.UpgradeResourceIdentityRequest{
			TypeName:    maximalFrameworkName,
			Version:     0,
			RawIdentity: &tfprotov5.RawState{JSON: []byte(`{"id":"unit-test"}`)},
		})
		if err != nil {
			t.Fatalf("UpgradeResourceIdentity: %s", err)
		}
		requireNoDiagErrors(t, "UpgradeResourceIdentity", resp.Diagnostics)

		if !log.has("UpgradeIdentity0") {
			t.Errorf("the identity upgrader did not run: %v", log.all())
		}
		if resp.UpgradedIdentity == nil {
			t.Fatal("the upgraded identity is missing")
		}
		if got := stringOf(t, attrsOf(t, maximalFrameworkIdentityType, resp.UpgradedIdentity.IdentityData)["id"]); got != "unit-test" {
			t.Errorf("the upgraded identity id is %q, want %q", got, "unit-test")
		}
		assertChainSilent(t, chain, "upgrading an identity")
	})

	t.Run("MoveResourceState", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.MoveResourceState(ctx, &tfprotov5.MoveResourceStateRequest{
			SourceProviderAddress: "registry.terraform.io/hashicorp/null",
			SourceTypeName:        "null_resource",
			SourceSchemaVersion:   0,
			SourceState:           &tfprotov5.RawState{JSON: []byte(`{"id":"unit-test"}`)},
			TargetTypeName:        maximalFrameworkName,
		})
		if err != nil {
			t.Fatalf("MoveResourceState: %s", err)
		}
		requireNoDiagErrors(t, "MoveResourceState", resp.Diagnostics)

		if !log.has("MoveState0") {
			t.Errorf("the state mover did not run: %v", log.all())
		}
		if got := stringOf(t, attrsOf(t, maximalFrameworkType, resp.TargetState)["name"]); got != "moved" {
			t.Errorf("the moved state has name %q, want %q", got, "moved")
		}
		assertChainSilent(t, chain, "moving a state")
	})

	t.Run("ImportResourceState", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.ImportResourceState(ctx, &tfprotov5.ImportResourceStateRequest{
			TypeName: maximalFrameworkName,
			ID:       "unit-test",
		})
		if err != nil {
			t.Fatalf("ImportResourceState: %s", err)
		}
		requireNoDiagErrors(t, "ImportResourceState", resp.Diagnostics)

		if len(resp.ImportedResources) != 1 {
			t.Fatalf("importing produced %d resources, want 1", len(resp.ImportedResources))
		}
		if !log.has("ImportState") {
			t.Errorf("the importer did not run: %v", log.all())
		}
		if log.has("Read") {
			t.Errorf("importing also ran Read: %v", log.all())
		}
		assertChainSilent(t, chain, "importing")
	})

	t.Run("ValidateDataSourceConfig", func(t *testing.T) {
		log.reset()
		chain.reset()

		resp, err := server.ValidateDataSourceConfig(ctx, &tfprotov5.ValidateDataSourceConfigRequest{
			TypeName: maximalFrameworkDataSourceName,
			Config: dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, map[string]tftypes.Value{
				"name": tftypes.NewValue(tftypes.String, "fixture"),
			})),
		})
		if err != nil {
			t.Fatalf("ValidateDataSourceConfig: %s", err)
		}
		requireNoDiagErrors(t, "ValidateDataSourceConfig", resp.Diagnostics)

		for _, hook := range []string{"ConfigValidators", "ConfigValidator", "ValidateConfig"} {
			if !log.has(hook) {
				t.Errorf("%s did not run during validation: %v", hook, log.all())
			}
		}
		if log.has("Read") {
			t.Errorf("validation also ran Read: %v", log.all())
		}
		assertChainSilent(t, chain, "validating a data source config")
	})

	t.Run("ReadDataSource", func(t *testing.T) {
		log.reset()
		chain.reset()

		config := dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, map[string]tftypes.Value{
			"name": tftypes.NewValue(tftypes.String, "fixture"),
		}))
		resp, err := server.ReadDataSource(ctx, &tfprotov5.ReadDataSourceRequest{
			TypeName: maximalFrameworkDataSourceName,
			Config:   config,
		})
		if err != nil {
			t.Fatalf("ReadDataSource: %s", err)
		}
		requireNoDiagErrors(t, "ReadDataSource", resp.Diagnostics)

		if got := stringOf(t, attrsOf(t, maximalFrameworkType, resp.State)["id"]); got != "unit-test" {
			t.Errorf("the data source produced id %q, want %q", got, "unit-test")
		}
		if got := chain.beforeOps(); !slices.Equal(got, []string{"Read"}) {
			t.Errorf("the chain saw %v for a data source read, want [Read]", got)
		}
		if got := chain.afterOps(); !slices.Equal(got, []string{"Read"}) {
			t.Errorf("the chain saw %v after a data source read, want [Read]", got)
		}
		if !log.has("Read/client") {
			t.Errorf("the data source read without a client: %v", log.all())
		}
	})

	t.Run("lifecycle", func(t *testing.T) {
		log.reset()
		chain.reset()

		runMaximalFrameworkLifecycle(t, server)

		assertChainSawLifecycle(t, chain, maximalFrameworkName, client)
		for _, op := range wantCRUDOps {
			if !log.has(op + "/client") {
				t.Errorf("%s ran without a client: %v", op, log.all())
			}
		}
		if !log.has("ModifyPlan") {
			t.Errorf("ModifyPlan never ran during a full lifecycle: %v", log.all())
		}
	})
}

func TestUnitProviderMaximalWireOneInterceptorBothStacks(t *testing.T) {
	chain := &chainRecorder{}
	client := &connectivity.AliyunClient{}

	sdkServer := newMaximalSDKv2Server(t, &hookLog{}, chain, client)
	fwServer := newMaximalFrameworkServer(t, &hookLog{}, chain, client)

	runMaximalSDKv2Lifecycle(t, sdkServer)
	runMaximalFrameworkLifecycle(t, fwServer)

	want := slices.Concat(wantCRUDOps, wantCRUDOps)
	if got := chain.beforeOps(); !slices.Equal(got, want) {
		t.Errorf("Before saw %v across both stacks, want %v", got, want)
	}
	if got := chain.afterOps(); !slices.Equal(got, want) {
		t.Errorf("After saw %v across both stacks, want %v", got, want)
	}

	wantNames := []string{
		maximalSDKv2Name, maximalSDKv2Name, maximalSDKv2Name, maximalSDKv2Name,
		maximalFrameworkName, maximalFrameworkName, maximalFrameworkName, maximalFrameworkName,
	}
	if got := namesOf(chain.snapshotBefore()); !slices.Equal(got, wantNames) {
		t.Errorf("Before saw names %v, want %v", got, wantNames)
	}

	for _, call := range chain.snapshotBefore() {
		if call.Meta != client {
			t.Errorf("%s on %s carried meta %v, want the configured client", call.Op, call.Name, call.Meta)
		}
	}
}

func newMaximalSDKv2Server(t *testing.T, log *hookLog, chain *chainRecorder, client *connectivity.AliyunClient) tfprotov5.ProviderServer {
	t.Helper()

	p := maximalSDKv2Provider(log, []intercept.Interceptor{chain})
	p.SetMeta(client)

	return p.GRPCProvider()
}

func newMaximalFrameworkServer(t *testing.T, log *hookLog, chain *chainRecorder, client *connectivity.AliyunClient) tfprotov5.ProviderServer {
	t.Helper()

	server := providerserver.NewProtocol5(&maximalFrameworkTestProvider{
		log:    log,
		chain:  []intercept.Interceptor{chain},
		client: client,
	})()

	resp, err := server.ConfigureProvider(context.Background(), &tfprotov5.ConfigureProviderRequest{
		Config: dynamic(t, emptyObjectType, tftypes.NewValue(emptyObjectType, map[string]tftypes.Value{})),
	})
	if err != nil {
		t.Fatalf("ConfigureProvider: %s", err)
	}
	requireNoDiagErrors(t, "ConfigureProvider", resp.Diagnostics)

	return server
}

func runMaximalSDKv2Lifecycle(t *testing.T, server tfprotov5.ProviderServer) {
	t.Helper()
	ctx := context.Background()

	name := tftypes.NewValue(tftypes.String, "fixture")
	imageID := tftypes.NewValue(tftypes.String, "image-one")
	tagsBefore := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
		"Created": tftypes.NewValue(tftypes.String, "terraform"),
	})
	tagsAfter := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{
		"Created": tftypes.NewValue(tftypes.String, "terraform"),
		"Updated": tftypes.NewValue(tftypes.String, "yes"),
	})

	createConfig := dynamic(t, maximalSDKv2Type, object(t, maximalSDKv2Type, map[string]tftypes.Value{
		"name":     name,
		"image_id": imageID,
		"tags":     tagsBefore,
	}))
	createPlan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         maximalSDKv2Name,
		PriorState:       nullOf(t, maximalSDKv2Type),
		ProposedNewState: createConfig,
		Config:           createConfig,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (create): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (create)", createPlan.Diagnostics)

	created, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:        maximalSDKv2Name,
		PriorState:      nullOf(t, maximalSDKv2Type),
		PlannedState:    createPlan.PlannedState,
		Config:          createConfig,
		PlannedPrivate:  createPlan.PlannedPrivate,
		PlannedIdentity: createPlan.PlannedIdentity,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (create): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (create)", created.Diagnostics)

	if got := stringOf(t, attrsOf(t, maximalSDKv2Type, created.NewState)["id"]); got != "unit-test" {
		t.Fatalf("the created resource has id %q, want %q", got, "unit-test")
	}
	if created.NewIdentity == nil {
		t.Fatal("the created resource has no identity")
	}
	if got := stringOf(t, attrsOf(t, maximalSDKv2IdentityType, created.NewIdentity.IdentityData)["name"]); got != "fixture" {
		t.Errorf("the created identity has name %q, want %q", got, "fixture")
	}

	refreshed, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:        maximalSDKv2Name,
		CurrentState:    created.NewState,
		Private:         created.Private,
		CurrentIdentity: created.NewIdentity,
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	requireNoDiagErrors(t, "ReadResource", refreshed.Diagnostics)

	if got := stringOf(t, attrsOf(t, maximalSDKv2Type, refreshed.NewState)["id"]); got != "unit-test" {
		t.Fatalf("the refreshed resource has id %q, want %q", got, "unit-test")
	}

	prior := attrsOf(t, maximalSDKv2Type, refreshed.NewState)
	updateConfig := dynamic(t, maximalSDKv2Type, object(t, maximalSDKv2Type, map[string]tftypes.Value{
		"name":     name,
		"image_id": imageID,
		"tags":     tagsAfter,
	}))
	updateProposed := dynamic(t, maximalSDKv2Type, object(t, maximalSDKv2Type, map[string]tftypes.Value{
		"id":       prior["id"],
		"status":   prior["status"],
		"name":     name,
		"image_id": imageID,
		"tags":     tagsAfter,
	}))
	updatePlan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         maximalSDKv2Name,
		PriorState:       refreshed.NewState,
		ProposedNewState: updateProposed,
		Config:           updateConfig,
		PriorPrivate:     refreshed.Private,
		PriorIdentity:    refreshed.NewIdentity,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (update): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (update)", updatePlan.Diagnostics)

	if len(updatePlan.RequiresReplace) != 0 {
		t.Errorf("planning a tags-only change requires replacing %v", updatePlan.RequiresReplace)
	}

	updated, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:        maximalSDKv2Name,
		PriorState:      refreshed.NewState,
		PlannedState:    updatePlan.PlannedState,
		Config:          updateConfig,
		PlannedPrivate:  updatePlan.PlannedPrivate,
		PlannedIdentity: updatePlan.PlannedIdentity,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (update): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (update)", updated.Diagnostics)

	if got := stringOf(t, attrsOf(t, maximalSDKv2Type, updated.NewState)["id"]); got != "unit-test" {
		t.Fatalf("the updated resource has id %q, want %q", got, "unit-test")
	}

	destroyPlan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         maximalSDKv2Name,
		PriorState:       updated.NewState,
		ProposedNewState: nullOf(t, maximalSDKv2Type),
		Config:           nullOf(t, maximalSDKv2Type),
		PriorPrivate:     updated.Private,
		PriorIdentity:    updated.NewIdentity,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (destroy): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (destroy)", destroyPlan.Diagnostics)

	destroyed, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:        maximalSDKv2Name,
		PriorState:      updated.NewState,
		PlannedState:    destroyPlan.PlannedState,
		Config:          nullOf(t, maximalSDKv2Type),
		PlannedPrivate:  destroyPlan.PlannedPrivate,
		PlannedIdentity: destroyPlan.PlannedIdentity,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (destroy): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (destroy)", destroyed.Diagnostics)

	assertStateIsNull(t, maximalSDKv2Type, destroyed.NewState)
}

func runMaximalFrameworkLifecycle(t *testing.T, server tfprotov5.ProviderServer) {
	t.Helper()
	ctx := context.Background()

	name := tftypes.NewValue(tftypes.String, "fixture")
	renamed := tftypes.NewValue(tftypes.String, "fixture-updated")

	createConfig := dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, map[string]tftypes.Value{
		"name": name,
	}))
	createPlan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         maximalFrameworkName,
		PriorState:       nullOf(t, maximalFrameworkType),
		ProposedNewState: createConfig,
		Config:           createConfig,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (create): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (create)", createPlan.Diagnostics)

	created, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:        maximalFrameworkName,
		PriorState:      nullOf(t, maximalFrameworkType),
		PlannedState:    createPlan.PlannedState,
		Config:          createConfig,
		PlannedPrivate:  createPlan.PlannedPrivate,
		PlannedIdentity: createPlan.PlannedIdentity,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (create): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (create)", created.Diagnostics)

	if got := stringOf(t, attrsOf(t, maximalFrameworkType, created.NewState)["id"]); got != "unit-test" {
		t.Fatalf("the created resource has id %q, want %q", got, "unit-test")
	}
	if created.NewIdentity == nil {
		t.Fatal("the created resource has no identity")
	}
	if got := stringOf(t, attrsOf(t, maximalFrameworkIdentityType, created.NewIdentity.IdentityData)["id"]); got != "unit-test" {
		t.Errorf("the created identity has id %q, want %q", got, "unit-test")
	}

	refreshed, err := server.ReadResource(ctx, &tfprotov5.ReadResourceRequest{
		TypeName:        maximalFrameworkName,
		CurrentState:    created.NewState,
		Private:         created.Private,
		CurrentIdentity: created.NewIdentity,
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	requireNoDiagErrors(t, "ReadResource", refreshed.Diagnostics)

	if got := stringOf(t, attrsOf(t, maximalFrameworkType, refreshed.NewState)["id"]); got != "unit-test" {
		t.Fatalf("the refreshed resource has id %q, want %q", got, "unit-test")
	}

	prior := attrsOf(t, maximalFrameworkType, refreshed.NewState)
	updateConfig := dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, map[string]tftypes.Value{
		"name": renamed,
	}))
	updateProposed := dynamic(t, maximalFrameworkType, object(t, maximalFrameworkType, map[string]tftypes.Value{
		"id":   prior["id"],
		"name": renamed,
	}))
	updatePlan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         maximalFrameworkName,
		PriorState:       refreshed.NewState,
		ProposedNewState: updateProposed,
		Config:           updateConfig,
		PriorPrivate:     refreshed.Private,
		PriorIdentity:    refreshed.NewIdentity,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (update): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (update)", updatePlan.Diagnostics)

	if len(updatePlan.RequiresReplace) != 0 {
		t.Errorf("planning a name change requires replacing %v", updatePlan.RequiresReplace)
	}

	updated, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:        maximalFrameworkName,
		PriorState:      refreshed.NewState,
		PlannedState:    updatePlan.PlannedState,
		Config:          updateConfig,
		PlannedPrivate:  updatePlan.PlannedPrivate,
		PlannedIdentity: updatePlan.PlannedIdentity,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (update): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (update)", updated.Diagnostics)

	if got := stringOf(t, attrsOf(t, maximalFrameworkType, updated.NewState)["name"]); got != "fixture-updated" {
		t.Errorf("the updated resource has name %q, want %q", got, "fixture-updated")
	}

	destroyPlan, err := server.PlanResourceChange(ctx, &tfprotov5.PlanResourceChangeRequest{
		TypeName:         maximalFrameworkName,
		PriorState:       updated.NewState,
		ProposedNewState: nullOf(t, maximalFrameworkType),
		Config:           nullOf(t, maximalFrameworkType),
		PriorPrivate:     updated.Private,
		PriorIdentity:    updated.NewIdentity,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange (destroy): %s", err)
	}
	requireNoDiagErrors(t, "PlanResourceChange (destroy)", destroyPlan.Diagnostics)

	destroyed, err := server.ApplyResourceChange(ctx, &tfprotov5.ApplyResourceChangeRequest{
		TypeName:        maximalFrameworkName,
		PriorState:      updated.NewState,
		PlannedState:    destroyPlan.PlannedState,
		Config:          nullOf(t, maximalFrameworkType),
		PlannedPrivate:  destroyPlan.PlannedPrivate,
		PlannedIdentity: destroyPlan.PlannedIdentity,
	})
	if err != nil {
		t.Fatalf("ApplyResourceChange (destroy): %s", err)
	}
	requireNoDiagErrors(t, "ApplyResourceChange (destroy)", destroyed.Diagnostics)

	assertStateIsNull(t, maximalFrameworkType, destroyed.NewState)
}

func assertChainSawLifecycle(t *testing.T, chain *chainRecorder, name string, client *connectivity.AliyunClient) {
	t.Helper()

	if got := chain.beforeOps(); !slices.Equal(got, wantCRUDOps) {
		t.Errorf("Before saw %v, want %v", got, wantCRUDOps)
	}
	if got := chain.afterOps(); !slices.Equal(got, wantCRUDOps) {
		t.Errorf("After saw %v, want %v", got, wantCRUDOps)
	}
	for _, call := range chain.snapshotBefore() {
		if call.Name != name {
			t.Errorf("%s was described as %q, want %q", call.Op, call.Name, name)
		}
		if call.Meta != client {
			t.Errorf("%s carried meta %v, want the configured client", call.Op, call.Meta)
		}
	}
}

func assertChainSilent(t *testing.T, chain *chainRecorder, doing string) {
	t.Helper()

	if got := chain.beforeOps(); len(got) != 0 {
		t.Errorf("%s ran the chain: %v", doing, got)
	}
	if got := chain.afterOps(); len(got) != 0 {
		t.Errorf("%s ran the chain: %v", doing, got)
	}
}

func assertStateIsNull(t *testing.T, ty tftypes.Type, dv *tfprotov5.DynamicValue) {
	t.Helper()

	if dv == nil {
		t.Fatal("the destroy response carried no state at all")
	}
	val, err := dv.Unmarshal(ty)
	if err != nil {
		t.Fatalf("unmarshalling the destroyed state: %s", err)
	}
	if !val.IsNull() {
		t.Errorf("the state after destroy is %s, want null", val)
	}
}

func requireNoDiagErrors(t *testing.T, step string, diags []*tfprotov5.Diagnostic) {
	t.Helper()

	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("%s: %s: %s", step, d.Summary, d.Detail)
		}
	}
}

func hasWarning(diags []*tfprotov5.Diagnostic, summary string) bool {
	for _, d := range diags {
		if d.Severity == tfprotov5.DiagnosticSeverityWarning && d.Summary == summary {
			return true
		}
	}
	return false
}

func object(t *testing.T, ty tftypes.Object, attrs map[string]tftypes.Value) tftypes.Value {
	t.Helper()

	for name := range attrs {
		if _, ok := ty.AttributeTypes[name]; !ok {
			t.Fatalf("%q is not an attribute of %s", name, ty)
		}
	}
	full := make(map[string]tftypes.Value, len(ty.AttributeTypes))
	for name, attrType := range ty.AttributeTypes {
		if v, ok := attrs[name]; ok {
			full[name] = v
			continue
		}
		full[name] = tftypes.NewValue(attrType, nil)
	}
	return tftypes.NewValue(ty, full)
}

func dynamic(t *testing.T, ty tftypes.Type, v tftypes.Value) *tfprotov5.DynamicValue {
	t.Helper()

	dv, err := tfprotov5.NewDynamicValue(ty, v)
	if err != nil {
		t.Fatalf("encoding a %s: %s", ty, err)
	}
	return &dv
}

func nullOf(t *testing.T, ty tftypes.Type) *tfprotov5.DynamicValue {
	t.Helper()

	return dynamic(t, ty, tftypes.NewValue(ty, nil))
}

func attrsOf(t *testing.T, ty tftypes.Type, dv *tfprotov5.DynamicValue) map[string]tftypes.Value {
	t.Helper()

	if dv == nil {
		t.Fatal("the response carried no value")
	}
	val, err := dv.Unmarshal(ty)
	if err != nil {
		t.Fatalf("decoding a %s: %s", ty, err)
	}
	attrs := map[string]tftypes.Value{}
	if err := val.As(&attrs); err != nil {
		t.Fatalf("reading the attributes of a %s: %s", ty, err)
	}
	return attrs
}

func stringOf(t *testing.T, v tftypes.Value) string {
	t.Helper()

	if v.IsNull() {
		return ""
	}
	var s string
	if err := v.As(&s); err != nil {
		t.Fatalf("reading a string: %s", err)
	}
	return s
}
