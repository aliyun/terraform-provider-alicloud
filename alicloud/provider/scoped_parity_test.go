package provider_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/registry"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/scoped"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

func TestUnitScopedFullRegistryCompatibility(t *testing.T) {
	ctx := context.Background()
	full := alicloud.Provider()
	selected := scoped.SDKProvider(registry.ServicePackages()...)
	legacy, err := provider.ProtoV5ProviderServerFactory(ctx, full)
	if err != nil {
		t.Fatal(err)
	}
	explicit, err := scoped.ProtoV5ProviderServerFactory(ctx, selected, registry.ServicePackages()...)
	if err != nil {
		t.Fatal(err)
	}
	left, err := legacy().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	right, err := explicit().GetProviderSchema(ctx, &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(left, right) {
		t.Fatal("explicit full selection differs from the legacy SDK/Framework protocol schema")
	}
}
