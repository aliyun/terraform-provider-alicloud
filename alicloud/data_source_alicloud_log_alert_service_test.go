package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogAlertResourceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogAlertResourceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_log_alert_resource.current_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_alert_resource.current_user", "user"),
					resource.TestCheckResourceAttr("data.alicloud_log_alert_resource.current_user", "region", "cn-hangzhou"),
					resource.TestCheckResourceAttr("data.alicloud_log_alert_resource.current_user", "lang", "cn"),
					testAccCheckAlicloudDataSourceID("data.alicloud_log_alert_resource.current_project"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_alert_resource.current_project", "project"),
					resource.TestCheckResourceAttr("data.alicloud_log_alert_resource.current_project", "project", "test-alert-tf"),
					resource.TestCheckResourceAttr("data.alicloud_log_alert_resource.current_project", "region", "cn-heyuan"),
				),
			},
		},
	})
}

const testAccCheckAlicloudLogAlertResourceDataSource = `
resource "alicloud_log_alert_resource" "current_user" { 
  type          = "user" 
  region        = "cn-hangzhou" 
  lang          = "cn" 
} 
 
resource "alicloud_log_alert_resource" "current_project" {
  type          = "project"
  project       = "test-alert-tf"
  region        = "cn-heyuan"
}
`
