package ga_test

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain hands control to the terraform-plugin-testing test runner, which
// drives acceptance tests and the per-region sweeper framework.
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// gaTestSignVersion is shared across regions; endpoint discovery caches are
// kept per region so that a request for one region is not routed to an
// endpoint cached for another (the SDK endpoint cache is keyed by product).
var gaTestSignVersion sync.Map
var gaTestRegionalEndpoints sync.Map

// gaTestEndpointCacheForRegion returns the per-region endpoint cache, creating
// one on first use; mirrors the flat-suite sharedEndpointCacheForRegion helper
// that is not reachable from this external test package.
func gaTestEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := gaTestRegionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// gaTestClientForRegion rebuilds the sweeper client on the exported
// connectivity API. It returns the concrete *connectivity.AliyunClient boxed as
// interface{} so the call sites that assert rawClient.(*connectivity.AliyunClient)
// keep working unchanged.
func gaTestClientForRegion(region string) (interface{}, error) {
	var accessKey, secretKey, securityToken string
	if accessKey = os.Getenv("ALICLOUD_ACCESS_KEY"); accessKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_ACCESS_KEY")
	}

	if secretKey = os.Getenv("ALICLOUD_SECRET_KEY"); secretKey == "" {
		return nil, fmt.Errorf("empty ALICLOUD_SECRET_KEY")
	}

	securityToken = os.Getenv("ALICLOUD_SECURITY_TOKEN")

	conf := connectivity.Config{
		Region:      connectivity.Region(region),
		RegionId:    region,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   gaTestEndpointCacheForRegion(region),
		SignVersion: &gaTestSignVersion,
	}
	if securityToken != "" {
		conf.SecurityToken = securityToken
	}
	accountId := os.Getenv("ALICLOUD_ACCOUNT_ID")
	if accountId == "" {
		accountId = os.Getenv("ALIBABA_CLOUD_ACCOUNT_ID")
	}
	if accountId != "" {
		conf.AccountId = accountId
	}
	accountType := os.Getenv("ALICLOUD_ACCOUNT_TYPE")
	if accountType == "" {
		accountType = os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")
	}
	if accountType != "" {
		conf.AccountType = accountType
	}
	credentialConfig := new(credentials.Config).SetType("access_key").SetAccessKeyId(accessKey).SetAccessKeySecret(secretKey)
	if v := strings.TrimSpace(securityToken); v != "" {
		credentialConfig.SetType("sts").SetSecurityToken(v)
	}
	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		return nil, err
	}
	conf.Credential = credential

	// configures a default client for the region, using the above env vars
	client, err := conf.Client()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// gaTestRegion selects the region to scope a test against, honoring the
// ALICLOUD_REGION env var and falling back to cn-hangzhou.
func gaTestRegion() string {
	if region := os.Getenv("ALICLOUD_REGION"); region != "" {
		return region
	}
	return "cn-hangzhou"
}
