package paifeaturestore_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// paiFeatureStoreInstanceEnv / paiFeatureStoreProjectEnv name the environment
// variables that supply a pre-existing FeatureStore instance and project for
// acceptance tests. FeatureEntity lives inside an instance+project that the
// test account provisions out of band; the tests do not create them.
const (
	paiFeatureStoreInstanceEnv = "PAI_FEATURE_STORE_INSTANCE_ID"
	paiFeatureStoreProjectEnv  = "PAI_FEATURE_STORE_PROJECT_ID"
)

func paiFeatureStoreAccPreCheck(t *testing.T) {
	acctest.PreCheck(t)
	instanceId := os.Getenv(paiFeatureStoreInstanceEnv)
	projectId := os.Getenv(paiFeatureStoreProjectEnv)
	if instanceId == "" || projectId == "" {
		t.Skipf("skipping acceptance test: %s and %s must reference a pre-existing FeatureStore instance and project",
			paiFeatureStoreInstanceEnv, paiFeatureStoreProjectEnv)
	}
}

func TestAccAliCloudPaiFeatureStoreFeatureEntity_basic(t *testing.T) {
	r := tfacctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			paiFeatureStoreAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiFeatureStoreFeatureEntityConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_feature_entity.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_feature_entity.default", "feature_entity_id"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_feature_entity.default", "feature_entity_name", fmt.Sprintf("tf-acc-fe-%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_feature_entity.default", "join_id", "user_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_feature_entity.default", "create_time"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_feature_entity.default", "region_id"),
				),
			},
			{
				ResourceName:      "alicloud_pai_feature_store_feature_entity.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudPaiFeatureStoreFeatureEntity_parent(t *testing.T) {
	r := tfacctest.RandIntRange(1000, 9999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			paiFeatureStoreAccPreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiFeatureStoreFeatureEntityParentConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_feature_entity.parent", "feature_entity_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_feature_entity.child", "feature_entity_id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_feature_entity.child", "parent_feature_entity_id"),
				),
			},
		},
	})
}

func testAccPaiFeatureStoreFeatureEntityConfig(r int) string {
	return fmt.Sprintf(`
variable "instance_id" {
  default = "%s"
}

variable "project_id" {
  default = "%s"
}

resource "alicloud_pai_feature_store_feature_entity" "default" {
  instance_id         = var.instance_id
  project_id          = var.project_id
  feature_entity_name = "tf-acc-fe-%[3]d"
  join_id             = "user_id"
}
`, os.Getenv(paiFeatureStoreInstanceEnv), os.Getenv(paiFeatureStoreProjectEnv), r)
}

func testAccPaiFeatureStoreFeatureEntityParentConfig(r int) string {
	return fmt.Sprintf(`
variable "instance_id" {
  default = "%s"
}

variable "project_id" {
  default = "%s"
}

resource "alicloud_pai_feature_store_feature_entity" "parent" {
  instance_id         = var.instance_id
  project_id          = var.project_id
  feature_entity_name = "tf-acc-fe-parent-%[3]d"
  join_id             = "user_id"
}

resource "alicloud_pai_feature_store_feature_entity" "child" {
  instance_id            = var.instance_id
  project_id             = var.project_id
  feature_entity_name    = "tf-acc-fe-child-%[3]d"
  join_id                = "user_id"
  parent_feature_entity_id = alicloud_pai_feature_store_feature_entity.parent.feature_entity_id
}
`, os.Getenv(paiFeatureStoreInstanceEnv), os.Getenv(paiFeatureStoreProjectEnv), r)
}
