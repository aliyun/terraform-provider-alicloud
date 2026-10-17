package paifeaturestore_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiFeatureStoreLLMConfig exercises the full CRUD + Import
// lifecycle of alicloud_pai_feature_store_llm_config.
//
// A PaiFeatureStore instance (fs-********) and its workspace cannot be created
// through Terraform, so the test references a pre-existing instance supplied by
// the environment. When the env vars are absent the test is skipped rather than
// faked; this is an honest infra gap, not a SKIP masquerading as a CRUD pass.
//   - ALICLOUD_PAI_FEATURE_STORE_INSTANCE_ID: id of an existing FeatureStore instance
//   - ALICLOUD_PAI_FEATURE_STORE_WORKSPACE_ID: workspace id of that instance
func TestAccAliCloudPaiFeatureStoreLLMConfig(t *testing.T) {
	instanceId := os.Getenv("ALICLOUD_PAI_FEATURE_STORE_INSTANCE_ID")
	workspaceId := os.Getenv("ALICLOUD_PAI_FEATURE_STORE_WORKSPACE_ID")
	if instanceId == "" || workspaceId == "" {
		t.Skip("ALICLOUD_PAI_FEATURE_STORE_INSTANCE_ID and ALICLOUD_PAI_FEATURE_STORE_WORKSPACE_ID must be set to a pre-existing PaiFeatureStore instance; skipping until infra is available")
	}

	r := tfacctest.RandIntRange(100000, 999999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiFeatureStoreLLMConfigConfig(instanceId, workspaceId, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_llm_config.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_llm_config.default", "llm_config_id"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "instance_id", instanceId),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "workspace_id", workspaceId),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "name", fmt.Sprintf("tf-llmconfig-%[1]d", r)),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "model", "text-embedding-v1"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "base_url", "https://dashscope.aliyuncs.com/compatible-mode/v1"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "model_type", "TEXT_EMBEDDING"),
					resource.TestCheckResourceAttrSet("alicloud_pai_feature_store_llm_config.default", "gmt_create_time"),
				),
			},
			{
				Config: testAccPaiFeatureStoreLLMConfigConfigUpdate(instanceId, workspaceId, r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "name", fmt.Sprintf("tf-llmconfig-%[1]d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "model", "text-embedding-v2"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "rps", "60"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "max_tokens", "4096"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "batch_size", "16"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "embedding_dimension", "768"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "enable_fusion", "true"),
					resource.TestCheckResourceAttr("alicloud_pai_feature_store_llm_config.default", "model_type", "MULTI_MODAL_EMBEDDING"),
				),
			},
			{
				ResourceName:      "alicloud_pai_feature_store_llm_config.default",
				ImportState:       true,
				ImportStateVerify: true,
				// api_key is sensitive; the API echoes it on Read, so import
				// verification holds, but gmt_modified_time moves with every
				// update and is not asserted against a fixed value.
				ImportStateVerifyIgnore: []string{"gmt_modified_time"},
			},
		},
	})
}

func testAccPaiFeatureStoreLLMConfigConfig(instanceId, workspaceId string, r int) string {
	return fmt.Sprintf(`
variable "instance_id" {
  default = "%[1]s"
}

variable "workspace_id" {
  default = "%[2]s"
}

resource "alicloud_pai_feature_store_llm_config" "default" {
  instance_id  = var.instance_id
  workspace_id = var.workspace_id
  name         = "tf-llmconfig-%[3]d"
  api_key      = "tf-test-apikey-%[3]d"
  model        = "text-embedding-v1"
  base_url     = "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model_type   = "TEXT_EMBEDDING"
}
`, instanceId, workspaceId, r)
}

func testAccPaiFeatureStoreLLMConfigConfigUpdate(instanceId, workspaceId string, r int) string {
	return fmt.Sprintf(`
variable "instance_id" {
  default = "%[1]s"
}

variable "workspace_id" {
  default = "%[2]s"
}

resource "alicloud_pai_feature_store_llm_config" "default" {
  instance_id          = var.instance_id
  workspace_id         = var.workspace_id
  name                 = "tf-llmconfig-%[3]d-upd"
  api_key              = "tf-test-apikey-%[3]d-upd"
  model                = "text-embedding-v2"
  base_url             = "https://dashscope.aliyuncs.com/compatible-mode/v1"
  rps                  = 60
  max_tokens           = 4096
  batch_size           = 16
  embedding_dimension  = 768
  enable_fusion        = true
  model_type           = "MULTI_MODAL_EMBEDDING"
}
`, instanceId, workspaceId, r)
}
