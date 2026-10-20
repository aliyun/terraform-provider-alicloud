package acctest

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DefaultRegionToTest snapshots ALICLOUD_REGION at package init, falling back
// to cn-beijing, mirroring the root suite's defaultRegionToTest: test
// fixtures read it at test-function-entry time, before PreCheck runs. PreCheck,
// PreCheckWithRegions and CheckoutSupportedRegions write it back when the
// region changes, but code that captured the old value before a region jump
// keeps seeing it — the root suite has the same entry-time gap.
var DefaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	return "cn-beijing"
}()

var SDKProviders = map[string]*schema.Provider{
	"alicloud": Provider,
}

var SDKProviderFactories = map[string]func() (*schema.Provider, error){
	"alicloud": func() (*schema.Provider, error) {
		return Provider, nil
	},
}

// Skip automatically the testcases which does not support some known regions.
// If supported is true, the regions should a list of supporting the service regions.
// If supported is false, the regions should a list of unsupporting the service regions.
// If the region is unsupported and has backend region, the backend region will instead
func PreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		t.Logf("[WARNING] The region is not set and using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	}
	CheckoutSupportedRegions(t, supported, regions)
	DefaultRegionToTest = os.Getenv("ALICLOUD_REGION")
}

func CheckoutSupportedRegions(t *testing.T, supported bool, regions []connectivity.Region) {
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
				DefaultRegionToTest = backupRegion
				return
			}
			// Try to use Hangzhou region
			if hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				DefaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			// If neither backup region nor Hangzhou region is available, use the first supported region
			if len(regions) > 0 {
				targetRegion := string(regions[0])
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, targetRegion)
				os.Setenv("ALICLOUD_REGION", targetRegion)
				DefaultRegionToTest = targetRegion
				return
			}
			t.Skipf("Skipping unsupported region %s. Supported regions: %s.", region, regions)
		} else {
			if !backupRegionFind {
				t.Logf("Skipping unsupported region %s. Unsupported regions: %s. Using %s as this test region", region, regions, backupRegion)
				os.Setenv("ALICLOUD_REGION", backupRegion)
				DefaultRegionToTest = backupRegion
				return
			}
			if !hangzhouRegionFind {
				t.Logf("Skipping unsupported region %s. Supported regions: %s. Using %s as this test region", region, regions, connectivity.Hangzhou)
				os.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))
				DefaultRegionToTest = string(connectivity.Hangzhou)
				return
			}
			t.Skipf("Skipping unsupported region %s. Unsupported regions: %s.", region, regions)
		}
		t.Skipped()
	}
}
