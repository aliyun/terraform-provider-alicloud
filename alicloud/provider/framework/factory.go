package framework

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ProtoV5ProviderServerFactory(ctx context.Context, primary *sdkschema.Provider) (func() tfprotov5.ProviderServer, error) {
	secondary := NewProvider(primary)

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
