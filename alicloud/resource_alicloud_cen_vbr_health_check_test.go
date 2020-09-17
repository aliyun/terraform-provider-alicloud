package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCenVbrHealthCheck_basic(t *testing.T) {
	var v cbn.VbrHealthCheck
	resourceId := "alicloud_cen_vbr_health_check.default"
	ra := resourceAttrInit(resourceId, CenVbrHealthCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CbnService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCenVbrHealthCheck")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCenVbrHealthCheck%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, CenVbrHealthCheckBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithCenVbrHealthCheckSetting(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cen_id":                 "${alicloud_cen_instance.default.id}",
					"health_check_source_ip": "192.168.1.2",
					"health_check_target_ip": "10.0.0.2",
					"vbr_instance_id":        os.Getenv("VBR_INSTANCE_ID"),
					"vbr_instance_region_id": os.Getenv("VBR_INSTANCE_REGION_ID"),
					"health_check_interval":  "2",
					"healthy_threshold":      "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":                 CHECKSET,
						"health_check_source_ip": "192.168.1.2",
						"health_check_target_ip": "10.0.0.2",
						"vbr_instance_id":        os.Getenv("VBR_INSTANCE_ID"),
						"vbr_instance_region_id": os.Getenv("VBR_INSTANCE_REGION_ID"),
						"health_check_interval":  "2",
						"healthy_threshold":      "8",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_source_ip": "192.168.1.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_source_ip": "192.168.1.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_target_ip": "10.0.0.3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_target_ip": "10.0.0.3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_interval": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_interval": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"healthy_threshold": "3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"healthy_threshold": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"health_check_source_ip": "192.168.1.2",
					"health_check_target_ip": "10.0.0.2",
					"health_check_interval":  "2",
					"healthy_threshold":      "8",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"health_check_source_ip": "192.168.1.2",
						"health_check_target_ip": "10.0.0.2",
						"health_check_interval":  "2",
						"healthy_threshold":      "8",
					}),
				),
			},
		},
	})
}

var CenVbrHealthCheckMap = map[string]string{
	"cen_id": CHECKSET,
}

func CenVbrHealthCheckBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alicloud_cen_instance" "default" {
  name = "${var.name}"
}
resource "alicloud_cen_instance_attachment" "default" {
  instance_id = "${alicloud_cen_instance.default.id}"
  child_instance_id = "%s"
  child_instance_type = "VBR"
  child_instance_region_id = "%s"
}
`, name, os.Getenv("VBR_INSTANCE_ID"), os.Getenv("VBR_INSTANCE_REGION_ID"))
}
