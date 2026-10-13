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

// WithClient makes "this resource can receive a client" assertable at
// registration; a field would not be.
type WithClient interface {
	Client() *connectivity.AliyunClient
}

type injectedClient struct {
	client *connectivity.AliyunClient
}

// A resource declaring its own Configure must call this one from it, or the
// promoted method is shadowed and the client never arrives.
type ResourceBase struct {
	injectedClient
}

type DataSourceBase struct {
	injectedClient
}

// These three have no wrapper yet, so clientNotInjected does not cover them: a
// shadowed Configure stays invisible until Client() returns nil.
type EphemeralResourceBase struct {
	injectedClient
}

type ListResourceBase struct {
	injectedClient
}

type ActionBase struct {
	injectedClient
}

var (
	_ WithClient = (*ResourceBase)(nil)
	_ WithClient = (*DataSourceBase)(nil)
	_ WithClient = (*EphemeralResourceBase)(nil)
	_ WithClient = (*ListResourceBase)(nil)
	_ WithClient = (*ActionBase)(nil)
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
// registration guard instead of keeping a nil client forever.
func (w *injectedClient) Client() *connectivity.AliyunClient { return w.client }

func (b *ResourceBase) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if client != nil {
		b.client = client
	}
}

func (b *DataSourceBase) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if client != nil {
		b.client = client
	}
}

func (b *EphemeralResourceBase) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	client, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if client != nil {
		b.client = client
	}
}

func (b *ListResourceBase) Configure(ctx context.Context, req list.ConfigureRequest, resp *list.ConfigureResponse) {
	client, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if client != nil {
		b.client = client
	}
}

func (b *ActionBase) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	client, diags := clientFromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if client != nil {
		b.client = client
	}
}

// Reports a resource that came out of Configure with no client although one was on
// offer. A nil providerData is normal; an existing error is more specific.
func clientNotInjected(name string, inner interface{}, providerData interface{}, existing diag.Diagnostics) diag.Diagnostics {
	if providerData == nil || existing.HasError() {
		return nil
	}
	carrier, ok := inner.(WithClient)
	if !ok || carrier.Client() != nil {
		return nil
	}
	return diag.Diagnostics{
		diag.NewErrorDiagnostic(
			"Provider Client Not Injected",
			fmt.Sprintf("%s (%T) ended Configure without a provider client. This happens when a resource declares its own "+
				"Configure and does not call the embedded fwadapt base's. Please report this issue to the provider developers.",
				name, inner),
		),
	}
}

func clientFromProviderData(providerData interface{}) (*connectivity.AliyunClient, diag.Diagnostics) {
	if providerData == nil {
		return nil, nil
	}
	client, ok := providerData.(*connectivity.AliyunClient)
	if !ok {
		return nil, diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Unexpected Configure Type",
				fmt.Sprintf("Expected *connectivity.AliyunClient, got %T. Please report this issue to the provider developers.", providerData),
			),
		}
	}
	return client, nil
}
