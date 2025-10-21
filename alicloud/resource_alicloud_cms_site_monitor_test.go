package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
		return fmt.Errorf("error getting AliCloud client: %s", err)
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

// Test CloudMonitorService SiteMonitor. >>> Resource test cases, automatically generated.
// Case pop3 5461
func TestAccAliCloudCloudMonitorServiceSiteMonitor_basic5461(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_site_monitor.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudMonitorServiceSiteMonitorMap5460)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceSiteMonitor")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudMonitorServiceSiteMonitorBasicDependence5460)
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
					"address":   "https://www.alibaba.com",
					"task_name": name,
					"task_type": "HTTP",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address":   "https://www.alibaba.com",
						"task_name": name,
						"task_type": "HTTP",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"address": "https://www.alibabacloud.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address": "https://www.alibabacloud.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_schedule": []map[string]interface{}{
						{
							"start_hour": "0",
							"days": []string{
								"2", "3"},
							"end_hour":  "2",
							"time_zone": "Local",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_schedule.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"interval": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"interval": "5",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"isp_cities": []map[string]interface{}{
						{
							"isp":  "232",
							"city": "641",
							"type": "IDC",
						},
						{
							"isp":  "5",
							"city": "738",
							"type": "LASTMILE",
						},
						{
							"isp":  "5",
							"city": "641",
							"type": "IDC",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"isp_cities.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"task_name": name + "_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"task_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"options_json": `{\"http_method\":\"get\",\"time_out\":5000}`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"options_json": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "1",
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

func TestAccAliCloudCloudMonitorServiceSiteMonitor_basic5461_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_site_monitor.default"
	ra := resourceAttrInit(resourceId, AliCloudCloudMonitorServiceSiteMonitorMap5460)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceSiteMonitor")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCloudMonitorServiceSiteMonitorBasicDependence5460)
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
					"address":     "https://www.alibaba.com",
					"agent_group": "PC",
					"task_name":   name,
					"task_type":   "HTTP",
					"custom_schedule": []map[string]interface{}{
						{
							"start_hour": "0",
							"days": []string{
								"2", "3"},
							"end_hour":  "2",
							"time_zone": "Local",
						},
					},
					"isp_cities": []map[string]interface{}{
						{
							"isp":  "232",
							"city": "641",
							"type": "IDC",
						},
						{
							"isp":  "5",
							"city": "738",
							"type": "LASTMILE",
						},
						{
							"isp":  "5",
							"city": "641",
							"type": "IDC",
						},
					},
					"interval":     "5",
					"options_json": `{\"http_method\":\"get\",\"time_out\":5000}`,
					"status":       "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"address":           "https://www.alibaba.com",
						"agent_group":       "PC",
						"task_name":         name,
						"task_type":         "HTTP",
						"isp_cities.#":      "3",
						"custom_schedule.#": "1",
						"interval":          "5",
						"options_json":      CHECKSET,
						"status":            "1",
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

var AliCloudCloudMonitorServiceSiteMonitorMap5460 = map[string]string{
	"agent_group":  CHECKSET,
	"interval":     CHECKSET,
	"options_json": CHECKSET,
	"status":       CHECKSET,
	"task_state":   CHECKSET,
}

func AliCloudCloudMonitorServiceSiteMonitorBasicDependence5460(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test CloudMonitorService SiteMonitor. <<< Resource test cases, automatically generated.
