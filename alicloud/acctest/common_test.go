package acctest

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func Test_ResourceTestAccConfigBuildsHcl(t *testing.T) {
	configFunc := ResourceTestAccConfig("alicloud_instance.default", "default", func(name string) string {
		return fmt.Sprintf("variable \"name\" {\n  default = %q\n}", name)
	})
	config := configFunc(map[string]interface{}{
		"instance_name":   "tf-test",
		"count":           2,
		"enabled":         true,
		"tags":            map[string]interface{}{"k": "v"},
		"security_groups": []interface{}{"sg-1", "sg-2"},
	})

	for _, want := range []string{
		"variable \"name\"",
		"resource \"alicloud_instance\" \"default\" {",
		"instance_name = \"tf-test\"",
		"count = 2",
		"enabled = true",
		"k = \"v\"",
		"\"sg-1\"",
		"\"sg-2\"",
	} {
		if !strings.Contains(config, want) {
			t.Fatalf("built config missing %q:\n%s", want, config)
		}
	}
	if !strings.HasSuffix(config, "\n}") {
		t.Fatalf("built config must close the block, got tail %q", config[len(config)-5:])
	}
}

func Test_DataSourceTestAccConfigFuncBuildsHcl(t *testing.T) {
	configFunc := DataSourceTestAccConfigFunc("data.alicloud_zones.zones", "default", func(name string) string {
		return ""
	})
	config := configFunc(map[string]interface{}{"enable_file_md5": true})
	if !strings.Contains(config, "data \"alicloud_zones\" \"zones\" {") {
		t.Fatalf("built config missing data block:\n%s", config)
	}
}

func Test_PreCheckSyncsDefaultRegionToTest(t *testing.T) {
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-access-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-secret-key")

	original := DefaultRegionToTest
	t.Cleanup(func() { DefaultRegionToTest = original })

	t.Setenv("ALICLOUD_REGION", "")
	DefaultRegionToTest = ""
	PreCheck(t)
	if DefaultRegionToTest != "cn-beijing" {
		t.Fatalf("DefaultRegionToTest = %q, want cn-beijing after region defaulting", DefaultRegionToTest)
	}

	t.Setenv("ALICLOUD_REGION", "cn-shanghai")
	PreCheck(t)
	if DefaultRegionToTest != "cn-shanghai" {
		t.Fatalf("DefaultRegionToTest = %q, want cn-shanghai after explicit region", DefaultRegionToTest)
	}
}

func Test_PreCheckWithRegionsCheckoutOverride(t *testing.T) {
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-access-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-secret-key")
	t.Setenv("CHECKOUT_REGION", "true")
	t.Setenv("ALICLOUD_REGION", string(connectivity.Hangzhou))

	original := DefaultRegionToTest
	t.Cleanup(func() { DefaultRegionToTest = original })

	PreCheckWithRegions(t, true, []connectivity.Region{connectivity.EUCentral1})

	if got := os.Getenv("ALICLOUD_REGION"); got != string(connectivity.EUCentral1) {
		t.Fatalf("ALICLOUD_REGION = %q, want %q", got, connectivity.EUCentral1)
	}
	if DefaultRegionToTest != string(connectivity.EUCentral1) {
		t.Fatalf("DefaultRegionToTest = %q, want %q", DefaultRegionToTest, connectivity.EUCentral1)
	}
}

func Test_CheckoutSupportedRegionsFallsBackToHangzhou(t *testing.T) {
	t.Setenv("CHECKOUT_REGION", "")
	t.Setenv("ALICLOUD_REGION", "cn-zhangjiakou")

	original := DefaultRegionToTest
	t.Cleanup(func() { DefaultRegionToTest = original })

	CheckoutSupportedRegions(t, true, []connectivity.Region{connectivity.Hangzhou})

	if got := os.Getenv("ALICLOUD_REGION"); got != string(connectivity.Hangzhou) {
		t.Fatalf("ALICLOUD_REGION = %q, want %q", got, connectivity.Hangzhou)
	}
	if DefaultRegionToTest != string(connectivity.Hangzhou) {
		t.Fatalf("DefaultRegionToTest = %q, want %q", DefaultRegionToTest, connectivity.Hangzhou)
	}
}

func Test_PreCheckPrePaidResourcesGate(t *testing.T) {
	t.Run("skips when unset", func(t *testing.T) {
		t.Setenv("ENABLE_CHECKING_PRE_PAID", "")
		PreCheckPrePaidResources(t)
		t.Fatal("expected skip")
	})
	t.Run("passes when enabled", func(t *testing.T) {
		t.Setenv("ENABLE_CHECKING_PRE_PAID", "true")
		PreCheckPrePaidResources(t)
	})
	t.Run("trims surrounding whitespace", func(t *testing.T) {
		t.Setenv("ENABLE_CHECKING_PRE_PAID", " true ")
		PreCheckPrePaidResources(t)
	})
}

func Test_ClassicNetworkResourcesGate(t *testing.T) {
	t.Run("skips when unset", func(t *testing.T) {
		t.Setenv("ENABLE_CLASSIC_NETWORK", "")
		ClassicNetworkResources(t)
		t.Fatal("expected skip")
	})
	t.Run("passes when enabled", func(t *testing.T) {
		t.Setenv("ENABLE_CLASSIC_NETWORK", "true")
		ClassicNetworkResources(t)
	})
}

func Test_PreCheckWithEnvVariableGate(t *testing.T) {
	t.Run("skips when unset", func(t *testing.T) {
		t.Setenv("TF_ACCTEST_HELPER_VAR", "")
		PreCheckWithEnvVariable(t, "TF_ACCTEST_HELPER_VAR")
		t.Fatal("expected skip")
	})
	t.Run("passes when set", func(t *testing.T) {
		t.Setenv("TF_ACCTEST_HELPER_VAR", "value")
		PreCheckWithEnvVariable(t, "TF_ACCTEST_HELPER_VAR")
	})
	t.Run("skips on whitespace-only value", func(t *testing.T) {
		t.Setenv("TF_ACCTEST_HELPER_VAR", "  ")
		PreCheckWithEnvVariable(t, "TF_ACCTEST_HELPER_VAR")
		t.Fatal("expected skip")
	})
}

func Test_PreCheckWithTimeGate(t *testing.T) {
	t.Run("passes on today", func(t *testing.T) {
		PreCheckWithTime(t, []int{time.Now().Day()})
	})
	t.Run("skips on no matching day", func(t *testing.T) {
		PreCheckWithTime(t, nil)
		t.Fatal("expected skip")
	})
}

func Test_PreCheckEnterpriseAccountEnabledGate(t *testing.T) {
	t.Run("skips when unset", func(t *testing.T) {
		t.Setenv("ENTERPRISE_ACCOUNT_ENABLED", "")
		PreCheckEnterpriseAccountEnabled(t)
		t.Fatal("expected skip")
	})
	t.Run("passes when enabled", func(t *testing.T) {
		t.Setenv("ENTERPRISE_ACCOUNT_ENABLED", "true")
		PreCheckEnterpriseAccountEnabled(t)
	})
	t.Run("does not trim surrounding whitespace", func(t *testing.T) {
		t.Setenv("ENTERPRISE_ACCOUNT_ENABLED", " true")
		PreCheckEnterpriseAccountEnabled(t)
		t.Fatal("expected skip")
	})
}

func Test_PreCheckWithMultipleAccountGate(t *testing.T) {
	t.Run("skips when unset", func(t *testing.T) {
		t.Setenv("ALICLOUD_ACCESS_KEY_2", "")
		PreCheckWithMultipleAccount(t)
		t.Fatal("expected skip")
	})
	t.Run("passes when set", func(t *testing.T) {
		t.Setenv("ALICLOUD_ACCESS_KEY_2", "second-account-key")
		PreCheckWithMultipleAccount(t)
	})
}

func TestUnitPreCheckWithAccountSiteTypeGate(t *testing.T) {
	// Exercise fallback selection without inheriting a developer or runner account.
	t.Setenv("ALICLOUD_ACCESS_KEY", "test-default-key")
	t.Setenv("ALICLOUD_SECRET_KEY", "test-default-secret")
	t.Setenv("ALICLOUD_ACCTEST_CREDENTIAL_MODE", "")
	for _, key := range []string{"ALICLOUD_INTL_ACCESS_KEY", "ALICLOUD_INTL_SECRET_KEY", "ALICLOUD_INTL_SECURITY_TOKEN", "ALICLOUD_INTL_ACCOUNT_ID", "ALICLOUD_ACCESS_KEY_1", "ALICLOUD_SECRET_KEY_1", "ALICLOUD_SECURITY_TOKEN_1", "ALICLOUD_ACCOUNT_ID_1"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("ALICLOUD_ACCOUNT_SITE", "")
	t.Setenv("ALIBABA_CLOUD_ACCOUNT_TYPE", "")
	t.Setenv("ALICLOUD_ACCOUNT_TYPE", "")

	t.Run("defaults to Domestic", func(t *testing.T) {
		PreCheckWithAccountSiteType(t, DomesticSite)
	})
	t.Run("skips mismatched account site", func(t *testing.T) {
		PreCheckWithAccountSiteType(t, IntlSite)
		t.Fatal("expected skip")
	})
	t.Run("ALICLOUD_ACCOUNT_SITE wins over ALIBABA_CLOUD_ACCOUNT_TYPE", func(t *testing.T) {
		t.Setenv("ALICLOUD_ACCOUNT_SITE", "Domestic")
		t.Setenv("ALIBABA_CLOUD_ACCOUNT_TYPE", "International")
		PreCheckWithAccountSiteType(t, DomesticSite)
	})
	t.Run("ALIBABA_CLOUD_ACCOUNT_TYPE fallback", func(t *testing.T) {
		t.Setenv("ALIBABA_CLOUD_ACCOUNT_TYPE", "International")
		PreCheckWithAccountSiteType(t, IntlSite)
	})
	t.Run("ALICLOUD_ACCOUNT_TYPE fallback", func(t *testing.T) {
		t.Setenv("ALICLOUD_ACCOUNT_TYPE", "International")
		PreCheckWithAccountSiteType(t, IntlSite)
	})
}

func Test_CheckAlicloudDataSourceID(t *testing.T) {
	s := terraform.NewState()
	s.RootModule().Resources["data.alicloud_foo.bar"] = &terraform.ResourceState{
		Primary: &terraform.InstanceState{ID: "d-123"},
	}

	if err := CheckAlicloudDataSourceID("data.alicloud_foo.bar")(s); err != nil {
		t.Fatalf("CheckAlicloudDataSourceID with set ID: unexpected error %s", err)
	}
	if err := CheckAlicloudDataSourceID("data.alicloud_foo.absent")(s); err == nil {
		t.Fatal("CheckAlicloudDataSourceID with absent resource must error")
	}

	s.RootModule().Resources["data.alicloud_foo.empty"] = &terraform.ResourceState{
		Primary: &terraform.InstanceState{ID: ""},
	}
	if err := CheckAlicloudDataSourceID("data.alicloud_foo.empty")(s); err == nil {
		t.Fatal("CheckAlicloudDataSourceID with empty ID must error")
	}
}

func Test_SDKProviderFactoriesAlternate(t *testing.T) {
	factories := SDKProviderFactoriesAlternate()
	if len(factories) != 2 {
		t.Fatalf("SDKProviderFactoriesAlternate returned %d factories, want 2", len(factories))
	}
	for _, name := range []string{"alicloud", "alicloudalt"} {
		f, ok := factories[name]
		if !ok {
			t.Fatalf("missing factory %q", name)
		}
		if p, err := f(); err != nil || p == nil {
			t.Fatalf("factory %q returned provider %v, err %v", name, p, err)
		}
	}
}

func Test_SDKProviderFactoriesMultiRegion(t *testing.T) {
	if factories := SDKProviderFactoriesMultiRegion(1); len(factories) != 1 {
		t.Fatalf("n=1 returned %d factories, want 1", len(factories))
	}
	factories := SDKProviderFactoriesMultiRegion(5)
	if len(factories) != 4 {
		t.Fatalf("n=5 returned %d factories, want 4 (alt names cap the count)", len(factories))
	}
	for _, name := range []string{"alicloud", "alicloudalt", "alicloudalt2", "alicloudalt3"} {
		if _, ok := factories[name]; !ok {
			t.Fatalf("missing factory %q", name)
		}
	}
}

func Test_ConfigAlternateRegionProvider(t *testing.T) {
	expected := "\nprovider \"alicloudalt\" {\n  region = \"cn-hangzhou\"\n}\n"
	if got := ConfigAlternateRegionProvider("cn-hangzhou"); got != expected {
		t.Fatalf("ConfigAlternateRegionProvider = %q, want %q", got, expected)
	}
}

func Test_ConfigNamedRegionalProvider(t *testing.T) {
	expected := "\nprovider \"alicloudalt2\" {\n  region = \"eu-central-1\"\n}\n"
	if got := ConfigNamedRegionalProvider("alicloudalt2", "eu-central-1"); got != expected {
		t.Fatalf("ConfigNamedRegionalProvider = %q, want %q", got, expected)
	}
}

func Test_CheckoutAccount(t *testing.T) {
	t.Setenv("ALICLOUD_ACCESS_KEY", "key-main")
	t.Setenv("ALICLOUD_SECRET_KEY", "secret-main")
	t.Setenv("ALICLOUD_ACCESS_KEY_MASTER", "")
	t.Setenv("ALICLOUD_SECRET_KEY_MASTER", "")
	t.Setenv("ALICLOUD_ACCESS_KEY_SLAVE", "")
	t.Setenv("ALICLOUD_SECRET_KEY_SLAVE", "")

	CheckoutAccount(t, true)
	if got := os.Getenv("ALICLOUD_ACCESS_KEY"); got != "key-main" {
		t.Fatalf("slave path without slave credentials must not swap, got ALICLOUD_ACCESS_KEY %q", got)
	}

	t.Setenv("ALICLOUD_ACCESS_KEY_SLAVE", "key-slave")
	t.Setenv("ALICLOUD_SECRET_KEY_SLAVE", "secret-slave")
	CheckoutAccount(t, true)
	if got := os.Getenv("ALICLOUD_ACCESS_KEY"); got != "key-slave" {
		t.Fatalf("slave path must swap in the slave key, got %q", got)
	}
	if got := os.Getenv("ALICLOUD_ACCESS_KEY_MASTER"); got != "key-main" {
		t.Fatalf("slave path must back up the main key, got %q", got)
	}

	t.Setenv("ALICLOUD_ACCESS_KEY_MASTER", "key-master")
	t.Setenv("ALICLOUD_SECRET_KEY_MASTER", "secret-master")
	CheckoutAccount(t, false)
	if got := os.Getenv("ALICLOUD_ACCESS_KEY"); got != "key-master" {
		t.Fatalf("master path must restore the master key, got %q", got)
	}

	t.Setenv("ALICLOUD_ACCESS_KEY_MASTER", "")
	CheckoutAccount(t, false)
	if got := os.Getenv("ALICLOUD_ACCESS_KEY"); got != "key-master" {
		t.Fatalf("master path without master credentials must not swap, got %q", got)
	}
}
