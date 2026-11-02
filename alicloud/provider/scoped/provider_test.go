package scoped_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/scoped"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	fwschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type fixture struct{ fwadapt.ResourceBase }

func (*fixture) Metadata(_ context.Context, _ resource.MetadataRequest, r *resource.MetadataResponse) {
	r.TypeName = "alicloud_scoped_framework"
}
func (*fixture) Schema(_ context.Context, _ resource.SchemaRequest, r *resource.SchemaResponse) {
	r.Schema = fwschema.Schema{Attributes: map[string]fwschema.Attribute{"id": fwschema.StringAttribute{Computed: true}}}
}
func (*fixture) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {}
func (*fixture) Read(context.Context, resource.ReadRequest, *resource.ReadResponse)       {}
func (*fixture) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {}
func (*fixture) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {}
func declaration() conns.ServicePackage {
	return conns.ServicePackage{Name: "Fixture", SDKResources: []conns.SDKResource{{TypeName: "alicloud_scoped_sdk", Factory: func() *schema.Resource {
		return &schema.Resource{Schema: map[string]*schema.Schema{"value": {Type: schema.TypeString, Optional: true}}, CreateContext: schema.NoopContext, ReadContext: schema.NoopContext, DeleteContext: schema.NoopContext}
	}}}, Resources: []conns.Resource{{TypeName: "alicloud_scoped_framework", Factory: func() resource.Resource { return &fixture{} }}}}
}
func TestUnitScopedSDKAndFrameworkSelection(t *testing.T) {
	ctx := context.Background()
	sp := declaration()
	p := scoped.SDKProvider(sp)
	if len(p.ResourcesMap) != 1 || p.ResourcesMap["alicloud_scoped_sdk"] == nil {
		t.Fatalf("unexpected SDK resources: %v", p.ResourcesMap)
	}
	factory, err := scoped.ProtoV5ProviderServerFactory(ctx, p, sp)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := factory().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range resp.Diagnostics {
		if d.Severity == tfprotov5.DiagnosticSeverityError {
			t.Fatalf("schema error: %s: %s", d.Summary, d.Detail)
		}
	}
	if len(resp.ResourceSchemas) != 2 || resp.ResourceSchemas["alicloud_scoped_sdk"] == nil || resp.ResourceSchemas["alicloud_scoped_framework"] == nil {
		t.Fatalf("wrong muxed selection: %v", resp.ResourceSchemas)
	}
	other := scoped.SDKProvider()
	if len(other.ResourcesMap) != 0 {
		t.Fatal("independent providers leaked registrations")
	}
	other.Schema["region"].Description = "changed"
	if p.Schema["region"].Description == "changed" {
		t.Fatal("provider configurations share mutable schema")
	}
}
func TestUnitScopedRejectsDuplicateAndNilDeclarations(t *testing.T) {
	for _, test := range []struct {
		name     string
		packages []conns.ServicePackage
		want     string
	}{
		{"duplicate", []conns.ServicePackage{declaration(), declaration()}, "declared twice"},
		{"cross protocol duplicate", []conns.ServicePackage{{Name: "Fixture", SDKResources: declaration().SDKResources, Resources: []conns.Resource{{TypeName: "alicloud_scoped_sdk", Factory: func() resource.Resource { return &fixture{} }}}}}, "declared twice"},
		{"nil factory", []conns.ServicePackage{{Name: "Bad", SDKResources: []conns.SDKResource{{TypeName: "alicloud_bad"}}}}, "nil Factory"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := scoped.ProtoV5ProviderServerFactory(context.Background(), scoped.SDKProvider(), test.packages...); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("unexpected factory error: %v", err)
			}
			defer func() {
				r := recover()
				if r == nil || !strings.Contains(r.(string), test.want) {
					t.Fatalf("unexpected panic: %v", r)
				}
			}()
			scoped.SDKProvider(test.packages...)
		})
	}
}
