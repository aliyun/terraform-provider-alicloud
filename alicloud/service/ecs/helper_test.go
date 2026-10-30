package ecs_test

// Local behavior-preserving compat copies of root-package test helpers that the
// migrated ECS tests reference bare. The ECS sub-package external tests
// (package ecs_test) may import the root alicloud package directly: they already
// depend on alicloud transitively via alicloud/acctest (acctest.Provider =
// alicloud.Provider()), and importing root here forms no production import
// cycle (only the internal package ecs cannot import root alicloud, because
// alicloud -> service/ecs). These helpers mirror root's alicloud_sweeper_test.go
// and provider_test.go verbatim, delegating to the shared alicloud/acctest
// package for precheck and to alicloud.Provider() for the provider factory.

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns the fully-wired alicloud provider. ECS resources are
// registered through service/ecs.ServicePackage delegation, so ResourcesMap
// lookups (p := Provider().ResourcesMap) resolve the migrated ECS types.
//
// This is the thinnest possible delegation to the real root alicloud.Provider()
// factory, mirroring root's provider_test.go Provider(): each call returns a
// freshly-built *schema.Provider so per-call factories (e.g. the image_copy
// Basic factory's p := Provider() that captures the providers list, and the
// multi-region ExistsWithProviders/DestroyWithProviders sweepers that need a
// distinct Meta/client per instance) cannot be folded into one shared pointer.
// Returning acctest.Provider instead would turn the fresh factory into a
// singleton (acctest.Provider is initialized once at acctest package init).
// The intentionally-shared testAccProvider is still constructed once in init()
// below, matching root provider_test.go's init().
func Provider() *schema.Provider {
	return alicloud.Provider()
}

var endpoints sync.Map
var signVersion sync.Map
var regionalEndpoints sync.Map

// SDK endpoint caches are keyed by product. Sharing one across regions can route
// a Hangzhou request to a previously cached Beijing endpoint.
func sharedEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := regionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// sharedClientForRegion returns a common AliyunClient setup needed for the
// sweeper functions for a given region. Verbatim from root alicloud_sweeper_test.go.
func sharedClientForRegion(region string) (interface{}, error) {
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
		Endpoints:   sharedEndpointCacheForRegion(region),
		SignVersion: &signVersion,
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

// sweepAll returns true when the sweeper should run across all resources.
// Verbatim from root alicloud_sweeper_test.go.
func sweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// testAccProvider mirrors root provider_test.go's provider instance; the
// dataSourceAttr test engine (copied into acc_common_test.go) builds
// testAccProviderFactory from it exactly as the flat layout's init() did.
var testAccProvider *schema.Provider

// testAccProviderFactory mirrors root provider_test.go's SDKv2 factory map;
// referenced by dataSourceAttr.dataSourceTestCheck. Initialized in init()
// below, matching the root init() that wires "alicloud" -> testAccProvider.
var testAccProviderFactory map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviderFactory = map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

// testAccPreCheck mirrors root provider_test.go's precheck; delegates to the
// shared acctest precheck (same ALICLOUD_ACCESS_KEY/SECRET_KEY/REGION checks
// and cn-beijing default). Referenced by dataSourceAttr.dataSourceTestCheck.
func testAccPreCheck(t *testing.T) {
	acctest.PreCheck(t)
}
