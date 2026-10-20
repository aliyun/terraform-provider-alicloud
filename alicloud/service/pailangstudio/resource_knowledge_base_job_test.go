package pailangstudio_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// KnowledgeBaseJob sits under a pre-existing workspace and knowledge base
// (KnowledgeBaseId is a Create path parameter; no knowledge-base resource exists
// in this provider to provision one in-test). The two parent IDs are therefore
// supplied via environment variables; when either is absent the test is skipped,
// which keeps the real ACC unverified rather than fabricated.
func testAccPreCheckPaiLangStudioKnowledgeBaseJob(t *testing.T) {
	acctest.PreCheck(t)
	if v := os.Getenv("ALICLOUD_PAILANGSTUDIO_WORKSPACE_ID"); v == "" {
		t.Skip("environment variable ALICLOUD_PAILANGSTUDIO_WORKSPACE_ID is not set; skipping KnowledgeBaseJob ACC")
	}
	if v := os.Getenv("ALICLOUD_PAILANGSTUDIO_KNOWLEDGE_BASE_ID"); v == "" {
		t.Skip("environment variable ALICLOUD_PAILANGSTUDIO_KNOWLEDGE_BASE_ID is not set; skipping KnowledgeBaseJob ACC")
	}
}

func TestAccAliCloudPaiLangStudioKnowledgeBaseJob(t *testing.T) {
	testAccPreCheckPaiLangStudioKnowledgeBaseJob(t)
	workspaceId := os.Getenv("ALICLOUD_PAILANGSTUDIO_WORKSPACE_ID")
	knowledgeBaseId := os.Getenv("ALICLOUD_PAILANGSTUDIO_KNOWLEDGE_BASE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliCloudPaiLangStudioKnowledgeBaseJobConfig(workspaceId, knowledgeBaseId, "tf-acc-kbjob"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_pailangstudio_knowledge_base_job.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_pailangstudio_knowledge_base_job.default", "knowledge_base_job_id"),
					resource.TestCheckResourceAttr("alicloud_pailangstudio_knowledge_base_job.default", "workspace_id", workspaceId),
					resource.TestCheckResourceAttr("alicloud_pailangstudio_knowledge_base_job.default", "knowledge_base_id", knowledgeBaseId),
					resource.TestCheckResourceAttr("alicloud_pailangstudio_knowledge_base_job.default", "description", "tf-acc-kbjob"),
					resource.TestCheckResourceAttr("alicloud_pailangstudio_knowledge_base_job.default", "accessibility", "PUBLIC"),
				),
			},
			{
				Config: testAccAliCloudPaiLangStudioKnowledgeBaseJobConfig(workspaceId, knowledgeBaseId, "tf-acc-kbjob-upd"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_pailangstudio_knowledge_base_job.default", "description", "tf-acc-kbjob-upd"),
				),
			},
			{
				// ImportState uses the composite id (`workspace:knowledgebase:job`)
				// from the prior step's state; no explicit ImportStateId is needed.
				ResourceName:      "alicloud_pailangstudio_knowledge_base_job.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliCloudPaiLangStudioKnowledgeBaseJobConfig(workspaceId, knowledgeBaseId, description string) string {
	return fmt.Sprintf(`
resource "alicloud_pailangstudio_knowledge_base_job" "default" {
  workspace_id     = "%[1]s"
  knowledge_base_id = "%[2]s"
  description      = "%[3]s"
}
`, workspaceId, knowledgeBaseId, description)
}
