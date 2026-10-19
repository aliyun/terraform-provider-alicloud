package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsAlertEventIntegrationPoliciesDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	resourceId := "data.alicloud_cms_alert_event_integration_policies.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlertEventIntegrationPoliciesBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{connectivity.Hangzhou})
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"workspace": "${alicloud_cms_workspace.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceId, "policies.#"),
				),
			},
		},
	})
}

func AliCloudCmsAlertEventIntegrationPoliciesBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_alert_event_integration_policy" "default" {
  alert_event_integration_policy_name = var.name
  workspace                           = alicloud_cms_workspace.default.id
  type                                = "SYS_EVENT"
  description                         = "tf-test-ds-description"
}
`, name)
}
