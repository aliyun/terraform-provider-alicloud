package acctest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ProtoV5ProviderFactories serves the real muxed provider — the SDK v2 half and the
// framework half behind one server — under the name given, defaulting to "alicloud".
//
// It works from a test on either harness, terraform-plugin-sdk/v2/helper/resource or
// terraform-plugin-testing/helper/resource, but only because no non-test file in package
// alicloud imports the former any more: both register a -sweep flag in their init, and a
// test binary linking both dies at startup with "flag redefined: sweep" before any test
// runs. This package imports package alicloud for Provider(), so a single non-test file
// in package alicloud going back to helper/resource resurrects that panic for every
// terraform-plugin-testing test that reaches this factory. Non-test code retries with
// helper/retry instead.
func ProtoV5ProviderFactories(names ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	if len(names) == 0 {
		names = []string{"alicloud"}
	}

	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(names))
	for _, name := range names {
		factories[name] = func() (tfprotov5.ProviderServer, error) {
			serverFactory, err := provider.ProtoV5ProviderServerFactory(context.Background(), alicloud.Provider())
			if err != nil {
				return nil, err
			}

			return serverFactory(), nil
		}
	}

	return factories
}

// PreCheck is the exported form of the suite's testAccPreCheck: it fails the test unless
// credentials are in the environment, and defaults the region the way the suite does.
//
// It does not set package alicloud's defaultRegionToTest, which is unexported and only
// read by tests in that package. Read the region back from ALICLOUD_REGION if a test in a
// service package needs it.
func PreCheck(t *testing.T) {
	if v := os.Getenv("ALICLOUD_ACCESS_KEY"); v == "" {
		t.Fatal("ALICLOUD_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_SECRET_KEY"); v == "" {
		t.Fatal("ALICLOUD_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALICLOUD_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
		os.Setenv("ALICLOUD_REGION", "cn-beijing")
	}
}
