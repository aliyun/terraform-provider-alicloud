package sls_test

import (
	"fmt"
	"testing"

	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudLogAlertResource_basic(t *testing.T) {
	var v *slsPop.AnalyzeProductLogResponse
	resourceId := "alicloud_log_alert_resource.default"
	ra := acctest.ResourceAttrInit(resourceId, logAlertResourceMap)
	serviceFunc := func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()

	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogalert-resource-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceAlertResourceConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      nil,
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
