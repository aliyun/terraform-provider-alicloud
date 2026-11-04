package acctest

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest/checks"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
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

func PreCheckWithRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	checks.PreCheckWithRegions(t, supported, regions)
	DefaultRegionToTest = os.Getenv("ALICLOUD_REGION")
}
func CheckoutSupportedRegions(t *testing.T, supported bool, regions []connectivity.Region) {
	checks.CheckoutSupportedRegions(t, supported, regions)
	DefaultRegionToTest = os.Getenv("ALICLOUD_REGION")
}

type AccountSite = checks.AccountSite

const (
	DomesticSite = checks.DomesticSite
	IntlSite     = checks.IntlSite
)

var PreCheckWithAccountSiteType = checks.PreCheckWithAccountSiteType

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

const (
	ProviderNameAlicloudAlt  = "alicloudalt"
	ProviderNameAlicloudAlt2 = "alicloudalt2"
	ProviderNameAlicloudAlt3 = "alicloudalt3"
)

func SDKProviderFactoriesAlternate() map[string]func() (*schema.Provider, error) {
	altProvider := alicloud.Provider()
	return map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			return Provider, nil
		},
		ProviderNameAlicloudAlt: func() (*schema.Provider, error) {
			return altProvider, nil
		},
	}
}

func SDKProviderFactoriesMultiRegion(n int) map[string]func() (*schema.Provider, error) {
	factories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			return Provider, nil
		},
	}
	altNames := []string{ProviderNameAlicloudAlt, ProviderNameAlicloudAlt2, ProviderNameAlicloudAlt3}
	for i := 0; i < n-1 && i < len(altNames); i++ {
		p := alicloud.Provider()
		factories[altNames[i]] = func() (*schema.Provider, error) {
			return p, nil
		}
	}
	return factories
}

func ConfigAlternateRegionProvider(region string) string {
	return fmt.Sprintf(`
provider "alicloudalt" {
  region = "%s"
}
`, region)
}

func ConfigNamedRegionalProvider(providerName string, region string) string {
	return fmt.Sprintf(`
provider "%s" {
  region = "%s"
}
`, providerName, region)
}

func CheckoutAccount(t *testing.T, SLAVE bool) {
	if SLAVE {
		if os.Getenv("ALICLOUD_ACCESS_KEY_SLAVE") == "" || os.Getenv("ALICLOUD_SECRET_KEY_SLAVE") == "" {
			t.Logf("\nALICLOUD_ACCESS_KEY_SLAVE or ALICLOUD_SECRET_KEY_SLAVE is empty and please add them.")
		} else {
			os.Setenv("ALICLOUD_ACCESS_KEY_MASTER", os.Getenv("ALICLOUD_ACCESS_KEY"))
			os.Setenv("ALICLOUD_SECRET_KEY_MASTER", os.Getenv("ALICLOUD_SECRET_KEY"))
			os.Setenv("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY_SLAVE"))
			os.Setenv("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY_SLAVE"))
			t.Logf("%s is using the slave account", t.Name())
		}
	} else {
		if os.Getenv("ALICLOUD_ACCESS_KEY_MASTER") == "" || os.Getenv("ALICLOUD_SECRET_KEY_MASTER") == "" {
			t.Logf("\nALICLOUD_ACCESS_KEY_MASTER or ALICLOUD_SECRET_KEY_MASTER is empty and please add them.")
		} else {
			os.Setenv("ALICLOUD_ACCESS_KEY", os.Getenv("ALICLOUD_ACCESS_KEY_MASTER"))
			os.Setenv("ALICLOUD_SECRET_KEY", os.Getenv("ALICLOUD_SECRET_KEY_MASTER"))
		}
	}
}
