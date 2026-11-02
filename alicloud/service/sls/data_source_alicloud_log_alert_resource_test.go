package sls_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudLogAlertResourceDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogalert-resource-%d", rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: dataSourceAlertResourceConfigDependence(name),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_log_alert_resource.current_user"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_alert_resource.current_user", "id"),
					acctest.CheckAlicloudDataSourceID("data.alicloud_log_alert_resource.current_project"),
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
