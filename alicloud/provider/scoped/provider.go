// Package scoped assembles explicitly selected service packages without importing
// the complete provider registry.
package scoped

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/conns"
	fw "github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework/scoped"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SDKProvider uses the production configuration and interception paths for the
// supplied products. Invalid declarations fail before any resource is configured.
func SDKProvider(packages ...conns.ServicePackage) *schema.Provider {
	if err := conns.Validate(packages); err != nil {
		panic("provider registration error: " + err.Error())
	}
	p := config.Provider()
	if diags := Assemble(p, packages); diags.HasError() {
		panic("provider registration error: " + FormatRegistrationDiags(diags))
	}
	return p
}

// ProtoV5ProviderServerFactory serves SDK and Framework declarations from the
// same selection. The SDK provider is configured first, just as in production.
func ProtoV5ProviderServerFactory(ctx context.Context, primary *schema.Provider, packages ...conns.ServicePackage) (func() tfprotov5.ProviderServer, error) {
	if err := conns.Validate(packages); err != nil {
		return nil, err
	}
	secondary := fw.NewProvider(primary, packages)
	mux, err := tf5muxserver.NewMuxServer(ctx, primary.GRPCProvider, providerserver.NewProtocol5(secondary))
	if err != nil {
		return nil, err
	}
	return mux.ProviderServer, nil
}
