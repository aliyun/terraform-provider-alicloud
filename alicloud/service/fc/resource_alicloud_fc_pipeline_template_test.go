package fc_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	fcsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/fc"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// PipelineTemplate names follow the FC k8s-style naming rule: lowercase
// letters, digits and hyphens, starting with a letter and not ending with a
// hyphen. tf-testacc-fcpt-<rand> satisfies it and stays well under the
// 40-character limit.
func pipelineTemplateName(rand int) string {
	return fmt.Sprintf("tf-testacc-fcpt-%d", rand)
}

func TestAccAlicloudFcPipelineTemplate_basic(t *testing.T) {
	resourceId := "alicloud_fc_pipeline_template.default"
	rand := tfacctest.RandIntRange(10000, 99999)
	name := pipelineTemplateName(rand)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             testAccCheckFcPipelineTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAlicloudFcPipelineTemplateBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "name", name),
					resource.TestCheckResourceAttr(resourceId, "description", "basic"),
					resource.TestCheckResourceAttr(resourceId, "labels.foo", "bar"),
					// spec is re-marshalled to canonical JSON on read, so only
					// assert it is populated rather than matching the literal
					// HCL jsonencode output.
					resource.TestCheckResourceAttrSet(resourceId, "spec"),
					resource.TestCheckResourceAttrSet(resourceId, "id"),
					resource.TestCheckResourceAttrSet(resourceId, "uid"),
					resource.TestCheckResourceAttr(resourceId, "kind", "PipelineTemplate"),
					resource.TestCheckResourceAttrSet(resourceId, "created_time"),
					resource.TestCheckResourceAttrSet(resourceId, "generation"),
					resource.TestCheckResourceAttrSet(resourceId, "resource_version"),
				),
			},
			{
				Config: testAccAlicloudFcPipelineTemplateBasicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceId, "description", "updated"),
					resource.TestCheckResourceAttr(resourceId, "labels.foo", "baz"),
					resource.TestCheckResourceAttr(resourceId, "labels.new", "value"),
					resource.TestCheckResourceAttrSet(resourceId, "spec"),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"generation", "resource_version", "deletion_time"},
			},
		},
	})
}

func testAccCheckFcPipelineTemplateDestroy(s *terraform.State) error {
	client := acctest.Provider.Meta().(*connectivity.AliyunClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_fc_pipeline_template" {
			continue
		}
		object, err := fcsvc.DescribePipelineTemplate(client, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("describing pipeline template %s failed: %s", rs.Primary.ID, err)
		}
		if object != nil {
			return fmt.Errorf("pipeline template %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccAlicloudFcPipelineTemplateBasic(name string) string {
	return fmt.Sprintf(`
resource "alicloud_fc_pipeline_template" "default" {
  name        = "%s"
  description = "basic"
  spec        = jsonencode({ context = { data = {} } })
  labels      = {
    foo = "bar"
  }
}
`, name)
}

func testAccAlicloudFcPipelineTemplateBasicUpdate(name string) string {
	return fmt.Sprintf(`
resource "alicloud_fc_pipeline_template" "default" {
  name        = "%s"
  description = "updated"
  spec        = jsonencode({ context = { data = { env = "prod" } } })
  labels      = {
    foo = "baz"
    new = "value"
  }
}
`, name)
}
