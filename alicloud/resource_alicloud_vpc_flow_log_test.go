package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVpcFlowLog_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpc_flow_log.default"
	ra := resourceAttrInit(resourceId, AlicloudVpcFlowLogMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpcFlowLog")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcflowlog%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpcFlowLogBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"flow_log_name":  "${var.name}",
					"log_store_name": "${alicloud_log_store.default.name}",
					"project_name":   "${alicloud_log_project.default.name}",
					"resource_id":    "${alicloud_vpc.default.id}",
					"resource_type":  "VPC",
					"traffic_type":   "All",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flow_log_name":  name,
						"log_store_name": CHECKSET,
						"project_name":   CHECKSET,
						"resource_id":    CHECKSET,
						"resource_type":  "VPC",
						"traffic_type":   "All",
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
					"description": "tf-testaccflowlogchange",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-testaccflowlogchange",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"flow_log_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"flow_log_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Active",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":   "tf-testaccflowlog",
					"flow_log_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "tf-testaccflowlog",
						"flow_log_name": name,
					}),
				),
			},
		},
	})
}

var AlicloudVpcFlowLogMap0 = map[string]string{
	"status":         "Active",
	"log_store_name": CHECKSET,
	"project_name":   CHECKSET,
	"resource_id":    CHECKSET,
	"resource_type":  "VPC",
	"traffic_type":   "All",
	"flow_log_name":  "",
	"description":    "",
}

func AlicloudVpcFlowLogBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%v"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  vpc_name       = var.name
}

resource "alicloud_log_project" "default" {
  name       = var.name
}

resource "alicloud_log_store" "default" {
  project  = alicloud_log_project.default.name
  name     = var.name
}

`, name)
}
