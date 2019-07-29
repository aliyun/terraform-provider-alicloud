package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCmsSiteMonitor_basic(t *testing.T) {
	resourceName := "alicloud_cms_sitemonitor.basic"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceName,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsSiteMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSiteMonitor_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.basic", "task_name", "tf-testAccCmsSiteMonitor_basic"),
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.basic", "interval", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.basic", "address", "http://www.alibabacloud.com"),
				),
			},
		},
	})
}

func TestAccAlicloudCmsSiteMonitor_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cms_sitemonitor.update",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsSiteMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSiteMonitor_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.update", "task_name", "tf-testAccCmsSiteMonitor_update"),
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.update", "interval", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.update", "address", "http://www.alibabacloud.com"),
				),
			},

			{
				Config: testAccCmsSiteMonitor_updateAfter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.update", "task_name", "tf-testAccCmsSiteMonitor_updateafter"),
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.update", "interval", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_sitemonitor.update", "address", "http://www.alibaba.com"),
				),
			},
		},
	})
}

func testAccCheckCmsSiteMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_sitemonitor" {
			continue
		}

		request := cms.CreateDescribeSiteMonitorListRequest()
		request.TaskId = rs.Primary.ID

		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DescribeSiteMonitorList(request)
		})
		list := raw.(*cms.DescribeSiteMonitorListResponse)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		if list.TotalCount > 0 {
			return fmt.Errorf("Site Monitor %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCmsSiteMonitor_basic() string {
	return fmt.Sprintf(`
	resource "alicloud_cms_sitemonitor" "basic" {
	  address = "http://www.alibabacloud.com"
	  task_name = "tf-testAccCmsSiteMonitor_basic"
	  task_type = "HTTP"
	  interval = 5
	}
	`)
}

func testAccCmsSiteMonitor_update() string {
	return fmt.Sprintf(`
data "alicloud_account" "current"{
}

resource "alicloud_cms_sitemonitor" "update" {
	address = "http://www.alibabacloud.com"
	task_name = "tf-testAccCmsSiteMonitor_update"
	task_type = "HTTP"
	interval = 5
}
`)
}

func testAccCmsSiteMonitor_updateAfter() string {
	return fmt.Sprintf(`
	data "alicloud_account" "current"{
	}
	
	resource "alicloud_cms_sitemonitor" "update" {
		address = "http://www.alibaba.com"
		task_name = "tf-testAccCmsSiteMonitor_updateafter"
		task_type = "HTTP"
		interval = 1
	}
	`)
}
