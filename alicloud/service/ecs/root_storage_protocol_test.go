// Package ecs contains the Framework-native resources of the ECS product.
package ecs

import (
	"context"
	"net/http"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// root_storage_protocol_test.go exercises the real resource through a public
// ProtoV6 protocol server (providerserver.NewProtocol6), not a constructed
// model or static self-assertion. It reproduces the root's
// ProtocolEmptyTagsNoop gate: a real DescribeStorageSets response with
// Tags.Tag=[] must Read into an empty KNOWN tags map, and a subsequent plan
// with that same config must be a no-op (prior == planned, no RequiresReplace).
// Before the fix, tagsFromObject returned MapNull for an empty array, so prior
// tags were null while planned tags were an empty known map — a perpetual
// non-no-op diff.

// storageSetProtocolTestProvider is a minimal framework provider that wires
// only the StorageSetResource and injects the mock client at Configure.
type storageSetProtocolTestProvider struct {
	client *connectivity.AliyunClient
}

var _ provider.Provider = (*storageSetProtocolTestProvider)(nil)

func (p *storageSetProtocolTestProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "alicloud"
}

func (p *storageSetProtocolTestProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerschema.Schema{}
}

func (p *storageSetProtocolTestProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.ResourceData = p.client
	resp.DataSourceData = p.client
}

func (p *storageSetProtocolTestProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return NewStorageSetResource() },
	}
}

func (p *storageSetProtocolTestProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

// storageSetTfType mirrors the resource schema attribute types for building
// raw protocol values.
var storageSetTfType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"id":                   tftypes.String,
	"description":          tftypes.String,
	"max_partition_number": tftypes.Number,
	"region_id":            tftypes.String,
	"resource_group_id":    tftypes.String,
	"storage_set_id":       tftypes.String,
	"storage_set_name":     tftypes.String,
	"tags":                 tftypes.Map{ElementType: tftypes.String},
	"zone_id":              tftypes.String,
}}

func protocolDynV6(t *testing.T, ty tftypes.Type, v tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	dv, err := tfprotov6.NewDynamicValue(ty, v)
	if err != nil {
		t.Fatalf("encoding a %s: %s", ty, err)
	}
	return &dv
}

func protocolObjectV6(t *testing.T, ty tftypes.Object, attrs map[string]tftypes.Value) tftypes.Value {
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

func protocolDiagHasError(diags []*tfprotov6.Diagnostic) bool {
	for _, d := range diags {
		if d == nil {
			continue
		}
		if d.Severity == tfprotov6.DiagnosticSeverityError {
			return true
		}
	}
	return false
}

func protocolAttrsV6(t *testing.T, ty tftypes.Type, dv *tfprotov6.DynamicValue) map[string]tftypes.Value {
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

// newStorageSetProtocolServer builds a configured ProtoV6 server over the real
// StorageSetResource backed by the httptest mock client.
func newStorageSetProtocolServer(t *testing.T, client *connectivity.AliyunClient) tfprotov6.ProviderServer {
	t.Helper()
	server := providerserver.NewProtocol6(&storageSetProtocolTestProvider{client: client})()
	emptyObj := tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}
	resp, err := server.ConfigureProvider(context.Background(), &tfprotov6.ConfigureProviderRequest{
		Config: protocolDynV6(t, emptyObj, tftypes.NewValue(emptyObj, map[string]tftypes.Value{})),
	})
	if err != nil {
		t.Fatalf("ConfigureProvider: %s", err)
	}
	if protocolDiagHasError(resp.Diagnostics) {
		t.Fatalf("ConfigureProvider diagnostics: %v", resp.Diagnostics)
	}
	return server
}

// TestUnitRootStorageProtocolEmptyTagsNoop drives the real resource over a public
// ProtoV6 server: ReadResource returns an empty KNOWN tags map from a real
// HTTP Tags.Tag=[] response, and a follow-up plan with that same config is a
// no-op (prior == planned, RequiresReplace empty). Before the fix the Read
// mapped the empty array to a null map, so prior(null) != planned(empty) and
// the next plan was not a no-op.
func TestUnitRootStorageProtocolEmptyTagsNoop(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DescribeStorageSets" {
			return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`)
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	server := newStorageSetProtocolServer(t, client)
	ctx := context.Background()

	// prior carries the id + zone_id so Read issues DescribeStorageSets.
	priorState := protocolDynV6(t, storageSetTfType, protocolObjectV6(t, storageSetTfType, map[string]tftypes.Value{
		"id":      tftypes.NewValue(tftypes.String, testStorageSetID),
		"zone_id": tftypes.NewValue(tftypes.String, "cn-hangzhou-i"),
	}))

	readResp, err := server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     storageSetTypeName,
		CurrentState: priorState,
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	if protocolDiagHasError(readResp.Diagnostics) {
		t.Fatalf("ReadResource diagnostics: %v", readResp.Diagnostics)
	}

	prior := protocolAttrsV6(t, storageSetTfType, readResp.NewState)
	priorTags := prior["tags"]
	if priorTags.IsNull() {
		t.Fatal("ReadResource mapped an empty Tag array to a null tags map; expected an empty KNOWN map")
	}
	if !priorTags.IsKnown() {
		t.Fatal("ReadResource mapped an empty Tag array to an unknown tags map; expected a KNOWN map")
	}

	// Plan with the refreshed state as both proposed and config: a no-op.
	planResp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         storageSetTypeName,
		PriorState:       readResp.NewState,
		ProposedNewState: readResp.NewState,
		Config:           readResp.NewState,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange: %s", err)
	}
	if protocolDiagHasError(planResp.Diagnostics) {
		t.Fatalf("PlanResourceChange diagnostics: %v", planResp.Diagnostics)
	}
	if len(planResp.RequiresReplace) != 0 {
		t.Fatalf("a no-op plan requires replacing %v", planResp.RequiresReplace)
	}

	planned := protocolAttrsV6(t, storageSetTfType, planResp.PlannedState)
	plannedTags := planned["tags"]
	if !priorTags.Equal(plannedTags) {
		t.Fatalf("no-op plan changed tags: prior=%s planned=%s", priorTags, plannedTags)
	}
}

// TestUnitRootStorageProtocolOmittedTagsPreservesState verifies that omitting tags
// in config (unknown via Framework plan) does not produce a destructive tag
// write and preserves the prior state: planned tags equal prior tags, no
// RequiresReplace. This is the whole-map unknown gate at the protocol layer.
func TestUnitRootStorageProtocolOmittedTagsPreservesState(t *testing.T) {
	srv, sent := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DescribeStorageSets" {
			return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[{"TagKey":"managed","TagValue":"keep"}]}`)
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	server := newStorageSetProtocolServer(t, client)
	ctx := context.Background()

	priorState := protocolDynV6(t, storageSetTfType, protocolObjectV6(t, storageSetTfType, map[string]tftypes.Value{
		"id":      tftypes.NewValue(tftypes.String, testStorageSetID),
		"zone_id": tftypes.NewValue(tftypes.String, "cn-hangzhou-i"),
	}))

	readResp, err := server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName:     storageSetTypeName,
		CurrentState: priorState,
	})
	if err != nil {
		t.Fatalf("ReadResource: %s", err)
	}
	if protocolDiagHasError(readResp.Diagnostics) {
		t.Fatalf("ReadResource diagnostics: %v", readResp.Diagnostics)
	}
	prior := protocolAttrsV6(t, storageSetTfType, readResp.NewState)
	priorTags := prior["tags"]
	if priorTags.IsNull() {
		t.Fatal("expected prior tags to carry the managed tag from the response")
	}

	// Plan with tags set to UNKNOWN (the Framework plan value for an omitted
	// Optional+Computed attribute). UseStateForUnknown must preserve prior.
	unknownTags := tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil)
	proposed := protocolDynV6(t, storageSetTfType, protocolObjectV6(t, storageSetTfType, map[string]tftypes.Value{
		"id":                   prior["id"],
		"zone_id":              prior["zone_id"],
		"max_partition_number": prior["max_partition_number"],
		"tags":                 unknownTags,
	}))
	planResp, err := server.PlanResourceChange(ctx, &tfprotov6.PlanResourceChangeRequest{
		TypeName:         storageSetTypeName,
		PriorState:       readResp.NewState,
		ProposedNewState: proposed,
		Config:           proposed,
	})
	if err != nil {
		t.Fatalf("PlanResourceChange: %s", err)
	}
	if protocolDiagHasError(planResp.Diagnostics) {
		t.Fatalf("PlanResourceChange diagnostics: %v", planResp.Diagnostics)
	}
	planned := protocolAttrsV6(t, storageSetTfType, planResp.PlannedState)
	plannedTags := planned["tags"]
	if !priorTags.Equal(plannedTags) {
		t.Fatalf("omitted (unknown) tags must preserve prior state: prior=%s planned=%s", priorTags, plannedTags)
	}
	// The plan itself does not issue RPC writes; only Apply does. Confirm no
	// Untag/Tag was recorded during planning (sanity).
	if _, ok := sent.Load("UntagResources"); ok {
		t.Fatal("planning an unknown tags map must not issue UntagResources")
	}
}
