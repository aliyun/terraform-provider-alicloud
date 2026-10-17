package paifs_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccAliCloudPaiFeatureStoreTasksDataSource exercises the live
// ListTasks/GetTask read path against a pre-existing PaiFeatureStore instance.
// It is gated on ALICLOUD_PAI_FEATURE_STORE_INSTANCE_ID because FeatureStoreTask
// is read-only: it has no Create/Delete/Update lifecycle API, so the acceptance
// fixture cannot be synthesized in the test account. The test is retained so a
// real instance id supplied by the product owner can re-run it as supplementary
// evidence.
func TestAccAliCloudPaiFeatureStoreTasksDataSource(t *testing.T) {
	if v := os.Getenv("ALICLOUD_PAI_FEATURE_STORE_INSTANCE_ID"); v == "" {
		t.Skip("ALICLOUD_PAI_FEATURE_STORE_INSTANCE_ID must be set to a pre-existing PaiFeatureStore instance id; FeatureStoreTask is read-only and cannot be synthesized.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		// FeatureStoreTask is a read-only data source with no managed
		// lifecycle, so there is nothing to destroy. The no-op CheckDestroy
		// satisfies tfproviderlint AT001 (which only exempts the legacy
		// data_source_alicloud_ filename prefix, not the framework layout).
		CheckDestroy: func(s *terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiFeatureStoreTasksDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_pai_feature_store_tasks.default", "id"),
				),
			},
		},
	})
}

func testAccAliCloudPaiFeatureStoreTasksDataSourceConfig() string {
	instanceID := os.Getenv("ALICLOUD_PAI_FEATURE_STORE_INSTANCE_ID")
	return `
data "alicloud_pai_feature_store_tasks" "default" {
  instance_id = "` + instanceID + `"
}
`
}
