package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogDashboard_basic(t *testing.T) {
	var Dashboard *sls.Dashboard
	resourceId := "alicloud_log_dashboard.default"
	ra := resourceAttrInit(resourceId, logDashboardMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &Dashboard, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogdashboard-%d", rand)
	displayname := fmt.Sprintf("dashboard_displayname-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogDashboardDependence)

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
					"project_name":   name,
					"dashboard_name": "dashboard_name",
					"display_name":   displayname,
					"char_list":      `[{\"title\":\"new_title\",\"type\":\"map\",\"search\":{\"logstore\":\"new_logstore\",\"topic\":\"new_topic\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"end\":\"now\"},\"display\":{\"xAxis\":[\"province\"],\"yAxis\":[\"aini\"],\"xPos\":0,\"yPos\":0,\"width\":10,\"height\":12,\"displayName\":\"xixihaha911\"}}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"project_name":   name,
						"dashboard_name": "dashboard_name",
						"display_name":   displayname,
						"char_list":      "[{\"display\":{\"displayName\":\"xixihaha911\",\"height\":12,\"width\":10,\"xAxis\":[\"province\"],\"xPos\":0,\"yAxis\":[\"aini\"],\"yPos\":0},\"search\":{\"end\":\"now\",\"logstore\":\"new_logstore\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"topic\":\"new_topic\"},\"title\":\"new_title\",\"type\":\"map\"}]",
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
					"display_name": "new_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "new_name",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"char_list": `[{\"title\":\"upadte_title\",\"type\":\"map\",\"search\":{\"logstore\":\"new_logstore\",\"topic\":\"new_topic\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"end\":\"\"},\"display\":{\"xAxis\":[\"province\"],\"yAxis\":[\"aini\"],\"xPos\":0,\"yPos\":0,\"width\":10,\"height\":12,\"displayName\":\"xixihaha911\"}}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"char_list": "[{\"display\":{\"displayName\":\"xixihaha911\",\"height\":12,\"width\":10,\"xAxis\":[\"province\"],\"xPos\":0,\"yAxis\":[\"aini\"],\"yPos\":0},\"search\":{\"end\":\"\",\"logstore\":\"new_logstore\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"topic\":\"new_topic\"},\"title\":\"upadte_title\",\"type\":\"map\"}]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "update_name",
					"char_list":    `[{\"title\":\"upadte_title_2\",\"type\":\"map\",\"search\":{\"logstore\":\"new_logstore\",\"topic\":\"new_topic\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"end\":\"\"},\"display\":{\"xAxis\":[\"province\"],\"yAxis\":[\"aini\"],\"xPos\":0,\"yPos\":0,\"width\":10,\"height\":12,\"displayName\":\"xixihaha911\"}}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "update_name",
						"char_list":    "[{\"display\":{\"displayName\":\"xixihaha911\",\"height\":12,\"width\":10,\"xAxis\":[\"province\"],\"xPos\":0,\"yAxis\":[\"aini\"],\"yPos\":0},\"search\":{\"end\":\"\",\"logstore\":\"new_logstore\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"topic\":\"new_topic\"},\"title\":\"upadte_title_2\",\"type\":\"map\"}]",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "update_name",
					"char_list":    `[{\"action\":{},\"title\":\"upadte_title_3\",\"type\":\"map\",\"search\":{\"logstore\":\"new_logstore\",\"topic\":\"new_topic\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"end\":\"\"},\"display\":{\"xAxis\":[\"province\"],\"yAxis\":[\"aini\"],\"xPos\":0,\"yPos\":0,\"width\":10,\"height\":12,\"displayName\":\"xixihaha911\"}}]`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "update_name",
						"char_list":    "[{\"action\":{},\"display\":{\"displayName\":\"xixihaha911\",\"height\":12,\"width\":10,\"xAxis\":[\"province\"],\"xPos\":0,\"yAxis\":[\"aini\"],\"yPos\":0},\"search\":{\"end\":\"\",\"logstore\":\"new_logstore\",\"query\":\"method:  GET  | select  ip_to_province(remote_addr) as province , count(1) as pv group by province order by pv desc \",\"start\":\"-86400s\",\"topic\":\"new_topic\"},\"title\":\"upadte_title_3\",\"type\":\"map\"}]",
					}),
				),
			},
		},
	})
}

var logDashboardMap = map[string]string{}

func resourceLogDashboardDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_project" "default" {
	    name = "${var.name}"
	    description = "tf unit test"
	}
	resource "alicloud_log_store" "default" {
	    project = "${alicloud_log_project.default.name}"
	    name = "${var.name}"
	    retention_period = "3000"
	    shard_count = 1
	}
	`, name)
}
