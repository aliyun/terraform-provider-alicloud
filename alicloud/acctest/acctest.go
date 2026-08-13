package acctest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/framework"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

func ProtoV5ProviderFactories(names ...string) map[string]func() (tfprotov5.ProviderServer, error) {
	if len(names) == 0 {
		names = []string{"alicloud"}
	}

	factories := make(map[string]func() (tfprotov5.ProviderServer, error), len(names))
	for _, name := range names {
		factories[name] = func() (tfprotov5.ProviderServer, error) {
			serverFactory, err := framework.ProtoV5ProviderServerFactory(context.Background(), alicloud.Provider())
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
