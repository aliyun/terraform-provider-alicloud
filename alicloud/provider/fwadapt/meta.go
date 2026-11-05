package fwadapt

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// WithMeta makes "this resource can receive the provider meta" assertable at
// registration; a field would not be.
type WithMeta interface {
	Meta() *connectivity.AliyunClient
}

// withMeta carries the provider meta: the configured provider state, of which
// the API client is what the resource sees today.
type withMeta struct {
	meta *connectivity.AliyunClient
}

// A resource declaring its own Configure must call this one from it, or the
// promoted method is shadowed and the meta never arrives.
type ResourceBase struct {
	withMeta
}

type DataSourceBase struct {
	withMeta
}

// These three have no wrapper yet, so metaNotInjected does not cover them: a
// shadowed Configure stays invisible until Meta() returns nil.
type EphemeralResourceBase struct {
	withMeta
}

type ListResourceBase struct {
	withMeta
}

type ActionBase struct {
	withMeta
}

var (
	_ WithMeta = (*ResourceBase)(nil)
	_ WithMeta = (*DataSourceBase)(nil)
	_ WithMeta = (*EphemeralResourceBase)(nil)
	_ WithMeta = (*ListResourceBase)(nil)
	_ WithMeta = (*ActionBase)(nil)
	_ interface {
		Configure(context.Context, resource.ConfigureRequest, *resource.ConfigureResponse)
	} = (*ResourceBase)(nil)
	_ interface {
		Configure(context.Context, datasource.ConfigureRequest, *datasource.ConfigureResponse)
	} = (*DataSourceBase)(nil)
	_ interface {
		Configure(context.Context, ephemeral.ConfigureRequest, *ephemeral.ConfigureResponse)
	} = (*EphemeralResourceBase)(nil)
	_ interface {
		Configure(context.Context, list.ConfigureRequest, *list.ConfigureResponse)
	} = (*ListResourceBase)(nil)
	_ interface {
		Configure(context.Context, action.ConfigureRequest, *action.ConfigureResponse)
	} = (*ActionBase)(nil)
)

// The pointer receiver matters: a factory returning a value then fails the
// registration guard instead of keeping a nil meta forever.
func (w *withMeta) Meta() *connectivity.AliyunClient { return w.meta }

func (b *ResourceBase) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	meta, diags := metaFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if meta != nil {
		b.meta = meta
	}
}

func (b *DataSourceBase) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	meta, diags := metaFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if meta != nil {
		b.meta = meta
	}
}

func (b *EphemeralResourceBase) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	meta, diags := metaFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if meta != nil {
		b.meta = meta
	}
}

func (b *ListResourceBase) Configure(ctx context.Context, req list.ConfigureRequest, resp *list.ConfigureResponse) {
	meta, diags := metaFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if meta != nil {
		b.meta = meta
	}
}

func (b *ActionBase) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	meta, diags := metaFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if meta != nil {
		b.meta = meta
	}
}

// Reports a resource that came out of Configure with no provider meta although
// one was on offer. A nil providerData is normal; an existing error is more
// specific.
func metaNotInjected(name string, inner interface{}, providerData interface{}, existing diag.Diagnostics) diag.Diagnostics {
	if providerData == nil || existing.HasError() {
		return nil
	}
	carrier, ok := inner.(WithMeta)
	if !ok || carrier.Meta() != nil {
		return nil
	}
	return diag.Diagnostics{
		diag.NewErrorDiagnostic(
			"Provider Meta Not Injected",
			fmt.Sprintf("%s (%T) ended Configure without the provider meta. This happens when a resource declares its own "+
				"Configure and does not call the embedded fwadapt base's. Please report this issue to the provider developers.",
				name, inner),
		),
	}
}

// The provider hands every framework kind the same meta today, the configured
// *connectivity.AliyunClient.
func metaFromProviderData(providerData interface{}) (*connectivity.AliyunClient, diag.Diagnostics) {
	if providerData == nil {
		return nil, nil
	}
	meta, ok := providerData.(*connectivity.AliyunClient)
	if !ok {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Unexpected Configure Type",
				fmt.Sprintf("Expected *connectivity.AliyunClient, got %T. Please report this issue to the provider developers.", providerData),
			),
		}
	}
	return meta, nil
}
