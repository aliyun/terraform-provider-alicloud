package provider

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProtoV5ProviderServerFactory muxes the SDK v2 provider and the framework provider behind
// one protocol v5 server. It is the only way either half should be served: main.go and the
// tests all go through it, so a resource behaves under test exactly as it does in a real
// terraform run, whichever half implements it. It lives here rather than in package
// framework because the muxed server is not the framework's — an SDK v2 resource is served
// by it too.
//
// primary comes first in the server list deliberately: tf5muxserver configures its servers
// in the order they are passed, and the framework half's Configure reads the parsed
// configuration back out of primary.Meta() instead of duplicating credential resolution.
func ProtoV5ProviderServerFactory(ctx context.Context, primary *sdkschema.Provider) (func() tfprotov5.ProviderServer, error) {
	secondary := framework.NewProvider(primary)

	servers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider,
		providerserver.NewProtocol5(secondary),
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, servers...)
	if err != nil {
		return nil, err
	}

	return muxServer.ProviderServer, nil
}
