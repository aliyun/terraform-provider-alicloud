package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_cms_site_monitor", &resource.Sweeper{
		Name: "alicloud_cms_site_monitor",
		F:    testSweepCmsSiteMonitor,
	})
}

func testSweepCmsSiteMonitor(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}
	request := cms.CreateDescribeSiteMonitorListRequest()
	raw, err := client.WithCmsClient(func(CmsClient *cms.Client) (interface{}, error) {
		return CmsClient.DescribeSiteMonitorList(request)
	})
	if err != nil {
		log.Printf("[ERROR] Error retrieving Cms Site Monitor: %s", WrapError(err))
	}
	response, _ := raw.(*cms.DescribeSiteMonitorListResponse)
	sweeped := false
	for _, v := range response.SiteMonitors.SiteMonitor {
		id := v.TaskId
		name := v.TaskName
		skip := true
		if !sweepAll() {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Cms Site Monitors: %s (%s)", name, id)
				continue
			}
		}

		sweeped = true
		log.Printf("[INFO] Deleting Cms Site Monitors: %s (%s)", name, id)
		req := cms.CreateDeleteSiteMonitorsRequest()
		req.TaskIds = id
		_, err := client.WithCmsClient(func(CmsClient *cms.Client) (interface{}, error) {
			return CmsClient.DeleteSiteMonitors(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Cms Site Monitors (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 30 seconds to ensure these Cms Site Monitors have been deleted.
		time.Sleep(30 * time.Second)
	}
	return nil
}

func TestAccAliCloudCmsSiteMonitor_basic(t *testing.T) {
	resourceName := "alicloud_cms_site_monitor.basic"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCmsSiteMonitorDestroy,
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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCmsSiteMonitor_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_cms_site_monitor.update",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCmsSiteMonitorDestroy,
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

func TestAccAliCloudCmsSiteMonitor_basic1(t *testing.T) {
	resourceName := "alicloud_cms_site_monitor.basic"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceName,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCmsSiteMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCmsSiteMonitor_basic1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.basic", "task_name", "tf-testAccCmsSiteMonitor_basic"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.basic", "interval", "5"),
					resource.TestCheckResourceAttr("alicloud_cms_site_monitor.basic", "address", "http://www.alibabacloud.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
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
  		address   = "http://www.alibabacloud.com"
  		task_name = "tf-testAccCmsSiteMonitor_basic"
  		task_type = "HTTP"
  		interval  = 5
  		isp_cities {
    		city = "546"
    		isp  = "465"
  		}
	}
	`)
}

func testAccCmsSiteMonitor_update() string {
	return fmt.Sprintf(`
	resource "alicloud_cms_site_monitor" "update" {
  		address   = "http://www.alibabacloud.com"
  		task_name = "tf-testAccCmsSiteMonitor_update"
  		task_type = "HTTP"
  		interval  = 5
  		isp_cities {
    		city = "546"
    		isp  = "465"
  		}
		options_json = <<EOT
		{
			"http_method": "get",
			"waitTime_after_completion": null,
			"ipv6_task": false,
			"diagnosis_ping": false,
			"diagnosis_mtr": false,
			"assertions": [
				{
					"operator": "lessThan",
					"type": "response_time",
					"target": 1000
				}
			],
			"time_out": 30000
		}
EOT
	}
	`)
}

func testAccCmsSiteMonitor_updateAfter() string {
	return fmt.Sprintf(`
	resource "alicloud_cms_site_monitor" "update" {
  		address   = "http://www.alibaba.com"
  		task_name = "tf-testAccCmsSiteMonitor_updateafter"
  		task_type = "HTTP"
  		interval  = 1
  		isp_cities {
    		city = "546"
    		isp  = "465"
  		}
  		isp_cities {
    		city = "572"
    		isp  = "465"
  		}
		options_json = <<EOT
		{
			"http_method": "post",
			"waitTime_after_completion": null,
			"ipv6_task": true,
			"diagnosis_ping": false,
			"diagnosis_mtr": false,
			"assertions": [
				{
					"operator": "lessThan",
					"type": "response_time",
					"target": 1000
				}
			],
			"time_out": 30000
		}
		EOT
	}
	`)
}

func testAccCmsSiteMonitor_basic1() string {
	return fmt.Sprintf(`
	resource "alicloud_cms_site_monitor" "basic" {
  		address   = "http://www.alibabacloud.com"
  		task_name = "tf-testAccCmsSiteMonitor_basic"
  		task_type = "PING"
  		interval  = 5
  		isp_cities {
			city = "546"
    		isp  = "465"
  		}
	}
	`)
}
