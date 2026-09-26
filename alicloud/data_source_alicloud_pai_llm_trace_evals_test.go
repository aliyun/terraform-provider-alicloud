package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudPaiLlmTraceEvals_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc-paillmtrace-evals-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceConfigForPaiLlmTraceEvals(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("evals_is_not_empty", "true"),
				),
			},
		},
	})
}

func testAccDataSourceConfigForPaiLlmTraceEvals(name string) string {
	return fmt.Sprintf(`
resource "alicloud_pai_llm_trace_eval" "default" {
  eval_name       = "%[1]s"
  description     = "tf-test-evals-description"
  app_name        = "tf-test-app-evals"
  data_source     = "tf-test-data-source-evals"
  evaluation_data = "{\"key\":\"value\"}"
  metadata        = "{\"env\":\"test\"}"
}

data "alicloud_pai_llm_trace_evals" "default" {
  ids = ["${alicloud_pai_llm_trace_eval.default.id}"]
  output_file = "evals.txt"
}

output "evals_is_not_empty" {
  value = length(data.alicloud_pai_llm_trace_evals.default.evals) > 0 ? "true" : "false"
}
`, name)
}
