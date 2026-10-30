package dcdn_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The generic ResourceCheck/ResourceAttr/ResourceAttrCheck/config builder,
// sentinels and data-source test framework now reuse the exported
// alicloud/acctest package directly (see common.go). Only the dcdn-specific
// precheck, region and sweeper helpers remain here: they are not reachable
// from a service package test because the root suite lives in package
// alicloud, and package dcdn_test keeps its own defaultRegionToTest snapshot
// in sync with ALICLOUD_REGION.

// defaultRegionToTest mirrors the root suite's package-level var: it snapshots
// ALICLOUD_REGION at test-function-entry time, before testAccPreCheck runs.
// testAccPreCheck keeps it in sync with the final ALICLOUD_REGION.
var defaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}()

// testAccPreCheck mirrors the root suite's testAccPreCheck: it fails the test
// unless credentials are in the environment and defaults the region. It keeps
// this package's defaultRegionToTest in sync with the final ALICLOUD_REGION.
func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
		defaultRegionToTest = "cn-beijing"
	} else {
		defaultRegionToTest = v
	}
}

// testAccPreCheckEnterpriseAccountEnabled mirrors the root suite's precheck of
// the same name.
func testAccPreCheckEnterpriseAccountEnabled(t *testing.T) {
	enable := os.Getenv("ENTERPRISE_ACCOUNT_ENABLED")
	if enable != "true" {
		t.Skipf("Skipping the test case because the enterprise account is not enabled.")
		t.Skipped()
	}
}

// ---------------------------------------------------------------------------
// Sweeper / region helpers — the flat-suite sharedClientForRegion,
// testSweepPreCheckWithRegions, checkoutSupportedRegions and sweepAll are
// not reachable from this service test package, so package-local mirrors are
// provided here.
// ---------------------------------------------------------------------------

var dcdnTestEndpoints sync.Map
var dcdnTestSignVersion sync.Map
var dcdnTestRegionalEndpoints sync.Map

// dcdnTestSharedEndpointCacheForRegion mirrors the root suite's
// sharedEndpointCacheForRegion: SDK endpoint caches are keyed by product, so
// sharing one across regions can route a Hangzhou request to a cached Beijing
// endpoint.
func dcdnTestSharedEndpointCacheForRegion(region string) *sync.Map {
	cache, _ := dcdnTestRegionalEndpoints.LoadOrStore(region, new(sync.Map))
	return cache.(*sync.Map)
}

// sharedClientForRegion mirrors alicloud_sweeper_test.go's helper: it builds
// an AliyunClient for a sweeper callback from environment credentials.
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
		Endpoints:   dcdnTestSharedEndpointCacheForRegion(region),
		SignVersion: &dcdnTestSignVersion,
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

// sweepAll mirrors the root suite's sweepAll gate.
func sweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// testSweepPreCheckWithRegions mirrors the root suite's sweeper gate.
func testSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}

// checkoutSupportedRegions mirrors the root suite's region precheck. It skips
// unsupported regions, defaulting to a backup region or Hangzhou when the
// configured region is not supported.
func checkoutSupportedRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	region := os.Getenv("ALICLOUD_REGION")
	find := false
	backupRegion := string(connectivity.APSouthEast1)
	if region == string(connectivity.APSouthEast1) {
		backupRegion = string(connectivity.EUCentral1)
	}

	checkoutRegion := os.Getenv("CHECKOUT_REGION")
	if checkoutRegion == "true" {
		if region == string(connectivity.Hangzhou) {
			region = string(connectivity.EUCentral1)
			os.Setenv("ALICLOUD_REGION", region)
		}
	}
	backupRegionFind := false
	hangzhouRegionFind := false
	for _, r := range regions {
		if region == string(r) {
			find = true
		}
		if string(r) == backupRegion {
			backupRegionFind = true
		}
		if string(connectivity.Hangzhou) == string(r) {
			hangzhouRegionFind = true
		}
	}

	if (find && !supported) || (!find && supported) {
		if supported {
			// Try to use backup region first
			if backupRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			// Try to use Hangzhou region
			if hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			// If neither backup region nor Hangzhou region is available, use the first supported region
			if len(regions) > 0 {
				targetRegion := string(regions[0])
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, targetRegion)
				os.Setenv("ALICLOUD_REGION", targetRegion)
				defaultRegionToTest = targetRegion
				return
			}
			t.Skipf("Skipping unsupported region %s. Supported regions: %s.", region, regions)
		} else {
			if !backupRegionFind {
				t.Logf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				defaultRegionToTest = backupRegion
				return
			}
			if !hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				defaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			t.Skipf("Skipping unsupported region %s. Unsupported regions: %s.", region, regions)
		}
		t.Skipped()
	}
}

func TestMain(m *testing.M) { resource.TestMain(m) }
