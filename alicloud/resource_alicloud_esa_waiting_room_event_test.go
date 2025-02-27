package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA WaitingRoomEvent. >>> Resource test cases, automatically generated.
// Case resource_WaitingRoomEvent_test
func TestAccAliCloudESAWaitingRoomEventresource_WaitingRoomEvent_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_waiting_room_event.default"
	ra := resourceAttrInit(resourceId, AliCloudESAWaitingRoomEventresource_WaitingRoomEvent_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaWaitingRoomEvent")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAWaitingRoomEvent%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAWaitingRoomEventresource_WaitingRoomEvent_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                         "off",
					"site_id":                        "${alicloud_esa_site.default.id}",
					"json_response_enable":           "off",
					"description":                    "测试1",
					"waiting_room_type":              "default",
					"end_time":                       "1719863200",
					"disable_session_renewal_enable": "off",
					"pre_queue_start_time":           "",
					"start_time":                     "1719763200",
					"random_pre_queue_enable":        "off",
					"waiting_room_event_name":        "WaitingRoomEvent_example",
					"queuing_status_code":            "200",
					"custom_page_html":               "",
					"new_users_per_minute":           "200",
					"session_duration":               "5",
					"language":                       "zhcn",
					"pre_queue_enable":               "off",
					"total_active_users":             "300",
					"waiting_room_id":                "${alicloud_esa_waiting_room.default.waiting_room_id}",
					"queuing_method":                 "fifo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                         "off",
					"json_response_enable":           "off",
					"description":                    "测试1",
					"waiting_room_type":              "default",
					"end_time":                       "1719963200",
					"disable_session_renewal_enable": "off",
					"pre_queue_start_time":           "",
					"start_time":                     "1719763200",
					"random_pre_queue_enable":        "off",
					"waiting_room_event_name":        "WaitingRoomEvent_example1",
					"queuing_status_code":            "202",
					"custom_page_html":               "",
					"new_users_per_minute":           "200",
					"session_duration":               "5",
					"language":                       "zhcn",
					"pre_queue_enable":               "on",
					"total_active_users":             "300",
					"queuing_method":                 "fifo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AliCloudESAWaitingRoomEventresource_WaitingRoomEvent_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAWaitingRoomEventresource_WaitingRoomEvent_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_waiting_room" "default" {
  status                         = "off"
  site_id                        = alicloud_esa_site.default.id
  json_response_enable           = "off"
  description                    = "example"
  waiting_room_type              = "default"
  disable_session_renewal_enable = "off"
  cookie_name                    = "__aliwaitingroom_example"
  waiting_room_name              = "waitingroom_example"
  queue_all_enable               = "off"
  queuing_status_code            = "200"
  custom_page_html               = ""
  new_users_per_minute           = "200"
  session_duration               = "5"
  language                       = "zhcn"
  total_active_users             = "300"
  queuing_method                 = "fifo"
  host_name_and_path {
    domain    = "sub_domain.com"
    path      = "/example"
    subdomain = "example_sub_domain.com."
  }
}

`, name)
}

// Test ESA WaitingRoomEvent. <<< Resource test cases, automatically generated.
