package fwadapt

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

type stubResourceEmbeddingBase struct{ ResourceBase }

type stubDataSourceEmbeddingBase struct{ DataSourceBase }

func TestResourceBaseConfigureInjectsMeta(t *testing.T) {
	r := &stubResourceEmbeddingBase{}
	client := &connectivity.AliyunClient{}

	var resp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, &resp)

	if r.Meta() != client {
		t.Fatal("Configure did not inject the meta")
	}
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diags: %v", resp.Diagnostics)
	}
}

func TestDataSourceBaseConfigureInjectsMeta(t *testing.T) {
	ds := &stubDataSourceEmbeddingBase{}
	client := &connectivity.AliyunClient{}

	var resp datasource.ConfigureResponse
	ds.Configure(context.Background(), datasource.ConfigureRequest{ProviderData: client}, &resp)

	if ds.Meta() != client {
		t.Fatal("Configure did not inject the meta")
	}
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diags: %v", resp.Diagnostics)
	}
}

func TestResourceBaseConfigureNilProviderData(t *testing.T) {
	r := &stubResourceEmbeddingBase{}

	var resp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: nil}, &resp)

	if r.Meta() != nil {
		t.Fatal("nil ProviderData must not set the meta")
	}
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diags: %v", resp.Diagnostics)
	}
}

func TestResourceBaseConfigureWrongProviderDataType(t *testing.T) {
	r := &stubResourceEmbeddingBase{}

	var resp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: "not a client"}, &resp)

	if r.Meta() != nil {
		t.Fatal("wrong ProviderData type must not set the meta")
	}
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diag for wrong ProviderData type")
	}
}

func TestBasesPromoteConfigure(t *testing.T) {
	var _ interface {
		Configure(context.Context, resource.ConfigureRequest, *resource.ConfigureResponse)
	} = &stubResourceEmbeddingBase{}
	var _ interface {
		Configure(context.Context, datasource.ConfigureRequest, *datasource.ConfigureResponse)
	} = &stubDataSourceEmbeddingBase{}
}

func TestBasesImplementWithMeta(t *testing.T) {
	var _ WithMeta = &stubResourceEmbeddingBase{}
	var _ WithMeta = &stubDataSourceEmbeddingBase{}

	if _, ok := interface{}(stubResourceEmbeddingBase{}).(WithMeta); ok {
		t.Error("a value satisfies WithMeta: the registration guard would pass a value-returning factory whose Configure never runs")
	}
	if _, ok := interface{}(stubDataSourceEmbeddingBase{}).(WithMeta); ok {
		t.Error("a value satisfies WithMeta: the registration guard would pass a value-returning factory whose Configure never runs")
	}
}

type stubEphemeralEmbeddingBase struct{ EphemeralResourceBase }

type stubListEmbeddingBase struct{ ListResourceBase }

type stubActionEmbeddingBase struct{ ActionBase }

func TestRemainingBasesConfigure(t *testing.T) {
	cases := []struct {
		name      string
		configure func(providerData interface{}) (*connectivity.AliyunClient, diag.Diagnostics)
	}{
		{
			name: "ephemeral resource",
			configure: func(providerData interface{}) (*connectivity.AliyunClient, diag.Diagnostics) {
				e := &stubEphemeralEmbeddingBase{}
				var resp ephemeral.ConfigureResponse
				e.Configure(context.Background(), ephemeral.ConfigureRequest{ProviderData: providerData}, &resp)
				return e.Meta(), resp.Diagnostics
			},
		},
		{
			name: "list resource",
			configure: func(providerData interface{}) (*connectivity.AliyunClient, diag.Diagnostics) {
				l := &stubListEmbeddingBase{}
				var resp list.ConfigureResponse
				l.Configure(context.Background(), list.ConfigureRequest{ProviderData: providerData}, &resp)
				return l.Meta(), resp.Diagnostics
			},
		},
		{
			name: "action",
			configure: func(providerData interface{}) (*connectivity.AliyunClient, diag.Diagnostics) {
				a := &stubActionEmbeddingBase{}
				var resp action.ConfigureResponse
				a.Configure(context.Background(), action.ConfigureRequest{ProviderData: providerData}, &resp)
				return a.Meta(), resp.Diagnostics
			},
		},
	}

	client := &connectivity.AliyunClient{}
	for _, tc := range cases {
		t.Run(tc.name+" injects the meta", func(t *testing.T) {
			got, diags := tc.configure(client)
			if got != client {
				t.Fatal("Configure did not inject the meta")
			}
			if diags.HasError() {
				t.Fatalf("unexpected diags: %v", diags)
			}
		})

		t.Run(tc.name+" tolerates nil provider data", func(t *testing.T) {
			got, diags := tc.configure(nil)
			if got != nil {
				t.Fatal("nil ProviderData must not set the meta")
			}
			if diags.HasError() {
				t.Fatalf("unexpected diags: %v", diags)
			}
		})

		t.Run(tc.name+" reports a wrong provider data type", func(t *testing.T) {
			got, diags := tc.configure("not a client")
			if got != nil {
				t.Fatal("wrong ProviderData type must not set the meta")
			}
			if !diags.HasError() {
				t.Fatal("expected error diag for wrong ProviderData type")
			}
		})
	}
}

type notACarrier struct{}

func TestMetaNotInjected(t *testing.T) {
	client := &connectivity.AliyunClient{}
	injected := &stubResourceEmbeddingBase{}
	injected.meta = client

	cases := []struct {
		name         string
		inner        interface{}
		providerData interface{}
		existing     diag.Diagnostics
		wantError    bool
	}{
		{
			name:         "shadowed Configure left the meta nil",
			inner:        &stubResourceEmbeddingBase{},
			providerData: client,
			wantError:    true,
		},
		{
			name:         "meta arrived",
			inner:        injected,
			providerData: client,
		},
		{
			name:         "no provider data on offer",
			inner:        &stubResourceEmbeddingBase{},
			providerData: nil,
		},
		{
			name:         "existing error is not buried",
			inner:        &stubResourceEmbeddingBase{},
			providerData: client,
			existing:     diag.Diagnostics{diag.NewErrorDiagnostic("Unexpected Configure Type", "")},
			wantError:    false,
		},
		{
			name:         "not a carrier is the registration guard's case",
			inner:        &notACarrier{},
			providerData: client,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := metaNotInjected("alicloud_unit_test", tc.inner, tc.providerData, tc.existing)
			if got.HasError() != tc.wantError {
				t.Fatalf("HasError() = %t, want %t (diags: %v)", got.HasError(), tc.wantError, got)
			}
		})
	}
}
