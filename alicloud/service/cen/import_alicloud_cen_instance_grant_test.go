package cen_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// This testcase can not work with multi account.
func SkipTestAccAlicloudCenInstanceGrant_importBasic(t *testing.T) {
	resourceName := "alicloud_cen_instance_grant.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	// multi provideris
	var providers []*schema.Provider
	providerFactories := map[string]func() (*schema.Provider, error){
		"alicloud": func() (*schema.Provider, error) {
			p := acctest.Provider
			providers = append(providers, p)
			return p, nil
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithMultipleAccount(t)
		},

		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckCenInstanceGrantDestroyWithProviders(&providers),
		Steps: []resource.TestStep{
			{
				Config: testAccCenInstanceGrantBasic(os.Getenv("ALICLOUD_ACCESS_KEY_2"), os.Getenv("ALICLOUD_SECRET_KEY_2"), os.Getenv("ALICLOUD_ACCOUNT_ID_1"), os.Getenv("ALICLOUD_ACCOUNT_ID_2"), rand),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
