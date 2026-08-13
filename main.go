package main

import (
	"context"
	"flag"
	"log"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
)

func main() {

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	serverFactory, err := framework.ProtoV5ProviderServerFactory(ctx, alicloud.Provider())
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	if err := tf5server.Serve("registry.terraform.io/aliyun/alicloud", serverFactory, serveOpts...); err != nil {
		log.Fatal(err)
	}
}
