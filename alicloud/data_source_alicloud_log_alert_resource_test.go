package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogAlertResourceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogalert-resource-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlertResourceConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_log_alert_resource.current_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_alert_resource.current_user", "id"),
					testAccCheckAlicloudDataSourceID("data.alicloud_log_alert_resource.current_project"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_alert_resource.current_project", "id"),
				),
			},
		},
	})
}

func dataSourceAlertResourceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "alert_resoure_project" {
  default = "%s"
}

resource "alicloud_log_project" "alert_resource_default"{
	name = "${var.alert_resoure_project}"
	description = "create by terraform"
}

data "alicloud_log_alert_resource" "current_user" { 
  type          = "user" 
  lang          = "cn" 
} 
 
data "alicloud_log_alert_resource" "current_project" {
  type          = "project"
  project       = "${alicloud_log_project.alert_resource_default.name}"
}

`, name)
}
