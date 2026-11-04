package scoped_test

import (
	"context"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest/scoped"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/fcv3"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

func TestUnitScopedSuiteIsolationAndAliases(t *testing.T) {
	a := scoped.New(fcv3.ServicePackage())
	b := scoped.New()
	if a.Provider == b.Provider || len(b.Provider.ResourcesMap) != 0 {
		t.Fatal("suite selections are not independent")
	}
	p, err := a.SDKProviderFactories["alicloud"]()
	if err != nil || p != a.Provider || a.SDKProviders["alicloud"] != p {
		t.Fatal("SDK factories must configure the client observed by test checks")
	}
	for name, factory := range a.ProtoV5ProviderFactories("alicloud", "alicloudalt") {
		server, err := factory()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		schema, err := server.GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})
		if err != nil {
			t.Fatal(err)
		}
		if len(schema.ResourceSchemas) != 10 || len(schema.DataSourceSchemas) != 2 {
			t.Fatalf("%s: wrong FCV3 registration counts: %d/%d", name, len(schema.ResourceSchemas), len(schema.DataSourceSchemas))
		}
	}
}

func TestUnitScopedPreChecksSynchronizeRegion(t *testing.T) {
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-secret")
	t.Setenv("ALICLOUD_REGION", "")
	t.Setenv("CHECKOUT_REGION", "")
	suite := scoped.New()
	suite.PreCheck(t)
	if suite.DefaultRegionToTest != "cn-beijing" {
		t.Fatalf("default region=%s", suite.DefaultRegionToTest)
	}
	suite.PreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
	if suite.DefaultRegionToTest != string(connectivity.Hangzhou) {
		t.Fatalf("fallback region=%s", suite.DefaultRegionToTest)
	}
}
