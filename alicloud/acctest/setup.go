package acctest

import (
	"os"
	"strings"
	"testing"
	"time"

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

type AccountSite string

const (
	DomesticSite = AccountSite("Domestic")
	IntlSite     = AccountSite("International")
)

// currently not all account site type support create PostPaid resources, PayByBandwidth and other limits.
// The setting of account site type can skip some unsupported cases automatically.

func PreCheckWithAccountSiteType(t *testing.T, account AccountSite) {
	defaultAccount := string(DomesticSite)
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_ACCOUNT_SITE")); v != "" {
		defaultAccount = v
	} else if v = strings.TrimSpace(os.Getenv("ALIBABA_CLOUD_ACCOUNT_TYPE")); v != "" {
		defaultAccount = v
	} else if v = strings.TrimSpace(os.Getenv("ALICLOUD_ACCOUNT_TYPE")); v != "" {
		defaultAccount = v
	}
	if defaultAccount != string(account) {
		t.Skipf("Skipping unsupported account type %s-Site. It only supports %s-Site.", defaultAccount, account)
		t.Skipped()
	}
}

func PreCheckPrePaidResources(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ENABLE_CHECKING_PRE_PAID")); v != "true" {
		t.Skip("Skipping testing PrePaid resources, otherwise setting environment parameter 'ENABLE_CHECKING_PRE_PAID'.")
		t.Skipped()
	}
}

func ClassicNetworkResources(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ENABLE_CLASSIC_NETWORK")); v != "true" {
		t.Skip("Skipping testing classic resources, otherwise setting environment parameter 'ENABLE_CLASSIC_NETWORK'.")
		t.Skipped()
	}
}

func PreCheckWithEnvVariable(t *testing.T, envVariableName string) {
	if v := strings.TrimSpace(os.Getenv(envVariableName)); v == "" {
		t.Skipf("Skipping the test case with no env variable %s", envVariableName)
		t.Skipped()
	}
}

func PreCheckWithTime(t *testing.T, days []int) {
	skipped := true
	for _, d := range days {
		if time.Now().Day() == d {
			skipped = false
			break
		}
	}
	if skipped {
		t.Skipf("Skipping the test case when not in specified days %#v of every month", days)
		t.Skipped()
	}
}

func PreCheckEnterpriseAccountEnabled(t *testing.T) {
	enable := os.Getenv("ENTERPRISE_ACCOUNT_ENABLED")
	if enable != "true" {
		t.Skipf("Skipping the test case because the enterprise account is not enabled.")
		t.Skipped()
	}
}

func PreCheckWithMultipleAccount(t *testing.T) {
	if v := strings.TrimSpace(os.Getenv("ALICLOUD_ACCESS_KEY_2")); v == "" {
		t.Skipf("Skipping unsupported test with multiple account")
		t.Skipped()
	}
}
