package paifs_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The Datasource is a child of a FeatureStoreInstance, which is not yet
// implemented as a Terraform resource. The acceptance test therefore relies on
// a pre-provisioned instance, workspace and a type-specific uri/config supplied
// through the environment; it skips when those fixtures are absent rather than
// fabricating values the API would reject.
const (
	envPaifsInstanceId  = "TF_ACC_PAI_FEATURE_STORE_INSTANCE_ID"
	envPaifsWorkspaceId = "TF_ACC_PAI_FEATURE_STORE_WORKSPACE_ID"
	envPaifsType        = "TF_ACC_PAI_FEATURE_STORE_DATASOURCE_TYPE"
	envPaifsUri         = "TF_ACC_PAI_FEATURE_STORE_DATASOURCE_URI"
	envPaifsConfig      = "TF_ACC_PAI_FEATURE_STORE_DATASOURCE_CONFIG"
)

func paifsDatasourceFixtures(t *testing.T) (instanceId, workspaceId, dsType, uri, config string) {
	instanceId = os.Getenv(envPaifsInstanceId)
	workspaceId = os.Getenv(envPaifsWorkspaceId)
	dsType = os.Getenv(envPaifsType)
	uri = os.Getenv(envPaifsUri)
	config = os.Getenv(envPaifsConfig)
	if dsType == "" {
		dsType = "FeatureDB"
	}
	if instanceId == "" || workspaceId == "" || uri == "" {
		t.Skipf("skipping acceptance test: set %s, %s and %s to a pre-provisioned FeatureStore instance, workspace and a valid datasource uri",
			envPaifsInstanceId, envPaifsWorkspaceId, envPaifsUri)
	}
	return
}

// paifsConfigBlock renders the config attribute only when a fixture is supplied.
// Omitting it keeps the attribute null in both steps, so the API never receives
// an empty Config and Update never has to clear a value the server persists.
func paifsConfigBlock(config string) string {
	if config == "" {
		return ""
	}
	return fmt.Sprintf(`  config = <<~EOT
%s
EOT
`, config)
}

func TestAccAliCloudPaiFeatureStoreDatasource_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)
	instanceId, workspaceId, dsType, uri, config := paifsDatasourceFixtures(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiFeatureStoreDatasourceConfig(r, instanceId, workspaceId, dsType, uri, config, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_datasource.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_datasource.default", "datasource_id"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_datasource.default", "instance_id", instanceId),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_datasource.default", "workspace_id", workspaceId),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_datasource.default", "type", dsType),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_datasource.default", "datasource_name", fmt.Sprintf("tf-testacc%d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_datasource.default", "uri", uri),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_datasource.default", "create_time"),
				),
			},
			{
				// Update DatasourceName and Uri; config, when supplied, is carried
				// unchanged so Update exercises body.Config without relying on a
				// second valid config value or on server-side clearing.
				Config: testAccAliCloudPaiFeatureStoreDatasourceConfig(r, instanceId, workspaceId, dsType, uri, config, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_datasource.default", "datasource_name", fmt.Sprintf("tf-testacc%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_datasource.default", "uri", fmt.Sprintf("%s-upd", uri)),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_datasource.default", "create_time"),
				),
			},
			{
				ResourceName:      "alicloud_pai_feature_store_datasource.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudPaiFeatureStoreDatasourceConfig(r int, instanceId, workspaceId, dsType, uri, config string, updated bool) string {
	name := fmt.Sprintf("tf-testacc%d", r)
	dsUri := uri
	if updated {
		name += "-upd"
		dsUri += "-upd"
	}
	return fmt.Sprintf(`
variable "instance_id" {
  type    = string
  default = "%[2]s"
}
variable "workspace_id" {
  type    = string
  default = "%[3]s"
}

resource "alicloud_pai_feature_store_datasource" "default" {
  instance_id     = var.instance_id
  workspace_id    = var.workspace_id
  type            = "%[4]s"
  datasource_name = "%[6]s"
  uri             = "%[7]s"
%[8]s}
`, r, instanceId, workspaceId, dsType, config, name, dsUri, paifsConfigBlock(config))
}
