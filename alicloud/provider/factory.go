package provider

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/scoped"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProtoV5ProviderServerFactory retains the complete SDK/Framework selection.
// Explicit selections use the same mux and configuration order through scoped.
func ProtoV5ProviderServerFactory(ctx context.Context, primary *sdkschema.Provider) (func() tfprotov5.ProviderServer, error) {
	return scoped.ProtoV5ProviderServerFactory(ctx, primary, registry.ServicePackages()...)
}
