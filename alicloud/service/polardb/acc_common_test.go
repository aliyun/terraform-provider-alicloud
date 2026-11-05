package polardb_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aws/smithy-go"
)

// Shared acceptance-test helpers re-homed into the split PolarDB external test
// package. The flat-layout originals live in the root alicloud package's
// internal _test.go files (service_alicloud_common_test.go / provider_test.go /
// alicloud_sweeper_test.go / common.go) and cannot be imported from a service
// package's external test package, so the subset the migrated PolarDB tests
// use is kept here with the same semantics. This mirrors the service/ecs
// acc_common_test.go pattern. Behavior is unchanged.

const (
	NOSET      = "#NOSET"       // be equivalent to method "TestCheckNoResourceAttrSet"
	CHECKSET   = "#CHECKSET"    // "TestCheckResourceAttrSet"
	REMOVEKEY  = "#REMOVEKEY"   // remove checkMap key
	CLEARLIST  = "#CLEARLIST"   // clear checkMap key and set empty list
	CLEARMAP   = "#CLEARMAP"    // clear checkMap key and set empty map
	REGEXMATCH = "#REGEXMATCH:" // "TestMatchResourceAttr" ,the map name/key like `"attribute" : REGEXMATCH + "attributeString"`
	ForceSleep = "force_sleep"
)

const PolarDBCommonTestCase = `
data "alicloud_polardb_zones" "default"{}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}
resource "alicloud_vpc" "default" {
    vpc_name = var.name
}
resource "alicloud_vswitch" "default" {
 vswitch_name = "tf_testAccPolarDB"
 vpc_id  = alicloud_vpc.default.id
 zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
 cidr_block = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
}
locals {
  vpc_id = alicloud_vpc.default.id
  vswitch_id = concat(alicloud_vswitch.default.*.id, [""])[0]
  zone_id = data.alicloud_polardb_zones.default.ids[length(data.alicloud_polardb_zones.default.ids)-1]
}
`

// mapMerge re-exports helper.MapMerge (root common.go alias) for the two
// mock-server unit tests that merge mock responses. Same target as the
// production alias in package polardb; the external test package holds its own
// reference rather than reverse-importing the production package's unexported
// symbol.
var mapMerge = helper.MapMerge

// defaultRegionToTest mirrors the root provider_test.go package-level variable
// (an IIFE returning the active region, falling back to cn-beijing). Read by
// PolarDB test fixtures that build region-prefixed names at test entry.
var defaultRegionToTest = func() string {
	if v := os.Getenv("ALICLOUD_REGION"); v != "" {
		return v
	}
	// Fallback to cn-beijing when ALICLOUD_REGION is unset so test fixtures
	// that read defaultRegionToTest at test-function-entry time (before
	// testAccPreCheck runs) don't produce malformed strings.
	return "cn-beijing"
}()

// sweepAll reports whether the sweeper is allowed to delete resources.
// Migrated verbatim from alicloud/alicloud_sweeper_test.go.
func sweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// checkoutSupportedRegions pins or swaps the active test region according to the
// supported/unsupported region list. Migrated verbatim from
// alicloud/provider_test.go; it assigns the local defaultRegionToTest so the
// region-prefixed fixtures stay consistent within this test package.
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

// isDynamoNotFoundError mirrors the production package's unexported helper
// (resource_dynamo_table.go) so the migrated Dynamo external tests can classify
// not-found errors without reverse-importing the production package's unexported
// symbol. Behavior is identical: smithy APIError code match plus a narrow
// string fallback for the PolarDB OpenAPI not-found wording.
func isDynamoNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "ResourceNotFoundException", "TableNotFound":
			return true
		}
		return false
	}
	// Errors from the PolarDB OpenAPI (not the DynamoDB-compatible endpoint) do
	// not implement smithy.APIError, so keep a narrow string match for them.
	return strings.Contains(err.Error(), "InvalidDBClusterId.NotFound")
}
