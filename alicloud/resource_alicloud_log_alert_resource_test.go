package alicloud

import (
	"fmt"
	"testing"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudLogAlertResource_basic(t *testing.T) {
	var v *slsPop.AnalyzeProductLogResponse
	resourceId := "alicloud_log_alert_resource.default"
	ra := resourceAttrInit(resourceId, logAlertResourceMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogalert-resource-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAlertResourceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"type":    "project",
					"project": "${alicloud_log_project.alert_resource_default.name}",
					"lang":    "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type": "project",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var logAlertResourceMap = map[string]string{}

func resourceAlertResourceConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "alert_resoure_project" {
  default = "%s"
}

resource "alicloud_log_project" "alert_resource_default"{
	name = "${var.alert_resoure_project}"
	description = "create by terraform"
}
`, name)
}
