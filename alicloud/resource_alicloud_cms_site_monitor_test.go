package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCmsSiteMonitor_basic(t *testing.T) {
	resourceName := "alicloud_cms_site_monitor.basic"
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
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.basic", "task_name", "tf-testAccCmsSiteMonitor_basic"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.basic", "interval", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.basic", "address", "http://www.alibabacloud.com"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"interval"},
			},
		},
	})
}

func TestAccAlicloudCmsSiteMonitor_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: "alicloud_cms_site_monitor.update",

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCmsSiteMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSiteMonitor_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "task_name", "tf-testAccCmsSiteMonitor_update"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "interval", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "address", "http://www.alibabacloud.com"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "isp_cities.#", "1"),
				),
			},

			{
				Config: testAccCmsSiteMonitor_updateAfter(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "task_name", "tf-testAccCmsSiteMonitor_updateafter"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "interval", "1"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "address", "http://www.alibaba.com"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.update", "isp_cities.#", "2"),
				),
			},
		},
	})
}

func testAccCheckCmsSiteMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cms_site_monitor" {
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
	resource "alicloud_cms_site_monitor" "basic" {
	  address = "http://www.alibabacloud.com"
	  task_name = "tf-testAccCmsSiteMonitor_basic"
	  task_type = "HTTP"
	  interval = 5
	  isp_cities {
		  city = "546"
		  isp = "465"
	  }
	}
	`)
}

func testAccCmsSiteMonitor_update() string {
	return fmt.Sprintf(`
data "alicloud_account" "current"{
}
resource "alicloud_cms_site_monitor" "update" {
	address = "http://www.alibabacloud.com"
	task_name = "tf-testAccCmsSiteMonitor_update"
	task_type = "HTTP"
	interval = 5
	isp_cities {
		city = "546"
		isp = "465"
	}
}
`)
}

func testAccCmsSiteMonitor_updateAfter() string {
	return fmt.Sprintf(`
	data "alicloud_account" "current"{
	}
	
	resource "alicloud_cms_site_monitor" "update" {
		address = "http://www.alibaba.com"
		task_name = "tf-testAccCmsSiteMonitor_updateafter"
		task_type = "HTTP"
		interval = 1
		isp_cities {
			city = "546"
			isp = "465"
		}
		isp_cities {
			city = "572"
			isp = "465"
		}
	}
	`)
}
