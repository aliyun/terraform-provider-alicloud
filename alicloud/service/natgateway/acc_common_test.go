package natgateway_test

import (
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// The flat-layout generic acceptance-test engine (resourceCheck / resourceAttr /
// resourceAttrCheck / resourceConfig reflection machinery, its sentinels and the
// dataSourceAttr runner) is sunk once into the exported alicloud/acctest package,
// which split service packages import directly; no local copy is kept here. Only
// the NAT-Gateway-specific test fixtures that have no public equivalent remain.

// natGatewayDataSourcePreCheck pins the migrated NAT Gateway data-source
// acceptance tests to eu-central-1, the only region where the test account has
// NAT Gateway inventory, then runs the standard PreCheck. It preserves the
// original region guard and ordering of the local dataSourceTestCheck runner,
// now routed through acctest.DataSourceTestCheckWithPreCheck. This is a
// test-environment fixture only; it does not change production Schema/CRUD/defaults.
func natGatewayDataSourcePreCheck(t *testing.T) func() {
	return func() {
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
		acctest.PreCheck(t)
	}
}

// unitTestMockClient returns a self-contained AliyunClient that issues no real
// network calls: it uses fake access-key credentials, caches the "vpc" endpoint
// so loadApiEndpoint never triggers describeEndpointForService, and skips region
// validation. The TestUnit subtests patch client.Client.DoRequest to drive
// deterministic responses through the same rpcRequest path used by production CRUD.
func unitTestMockClient(t *testing.T) *connectivity.AliyunClient {
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoints.Store("vpc", "vpc.test-endpoint.example.com")
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	meta, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	return meta
}
