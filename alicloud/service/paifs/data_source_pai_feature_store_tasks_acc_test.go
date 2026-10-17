package paifs_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiFeatureStoreTasksDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_pai_feature_store_tasks.default", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_feature_store_tasks.default", "tasks.#"),
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
