package paillmtrace_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudPaiLLMTraceEvaluationTemplatesDataSource exercises the data source
// against the live GetEvaluationTemplates API. The API exposes public information
// and takes no input parameters, so the count of returned templates depends on the
// account; only the identifier and the list shape are asserted, not the count.
func TestAccAliCloudPaiLLMTraceEvaluationTemplatesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPaiLLMTraceEvaluationTemplatesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_pai_llm_trace_evaluation_templates.default", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_pai_llm_trace_evaluation_templates.default", "templates.#"),
				),
			},
		},
	})
}

func testAccPaiLLMTraceEvaluationTemplatesDataSourceConfig() string {
	return `
data "alicloud_pai_llm_trace_evaluation_templates" "default" {
}
`
}
