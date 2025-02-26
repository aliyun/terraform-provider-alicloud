package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ESA WaitingRoom. >>> Resource test cases, automatically generated.
// Case resource_WaitingRoom_event_test
func TestAccAliCloudESAWaitingRoomresource_WaitingRoom_event_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_waiting_room.default"
	ra := resourceAttrInit(resourceId, AliCloudESAWaitingRoomresource_WaitingRoom_event_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaWaitingRoom")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAWaitingRoom%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAWaitingRoomresource_WaitingRoom_event_testBasicDependence)
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
					"status": "off",
					"host_name_and_path": []map[string]interface{}{

						{
							"path":      "/test",
							"subdomain": "test_sub_domain.com.",
							"domain":    "sub_domain.com",
						},

						{
							"path":      "/test",
							"subdomain": "test_sub_domain1.com.",
							"domain":    "sub_domain.com",
						},

						{
							"path":      "/test",
							"subdomain": "test_sub_domain2.com.",
							"domain":    "sub_domain.com",
						},
					},
					"site_id":                        "${alicloud_esa_site.resource_Site_event_test.id}",
					"json_response_enable":           "off",
					"description":                    "测试1",
					"waiting_room_type":              "default",
					"disable_session_renewal_enable": "off",
					"cookie_name":                    "__aliwaitingroom_example",
					"waiting_room_name":              "waitingroom_example",
					"queue_all_enable":               "off",
					"queuing_status_code":            "200",
					"custom_page_html":               "",
					"new_users_per_minute":           "200",
					"session_duration":               "5",
					"language":                       "zhcn",
					"total_active_users":             "300",
					"queuing_method":                 "fifo",
				}),
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

var AliCloudESAWaitingRoomresource_WaitingRoom_event_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAWaitingRoomresource_WaitingRoom_event_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_event_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_WaitingRoom_test
func TestAccAliCloudESAWaitingRoomresource_WaitingRoom_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_waiting_room.default"
	ra := resourceAttrInit(resourceId, AliCloudESAWaitingRoomresource_WaitingRoom_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaWaitingRoom")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAWaitingRoom%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAWaitingRoomresource_WaitingRoom_testBasicDependence)
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
					"status": "off",
					"host_name_and_path": []map[string]interface{}{

						{
							"path":      "/test",
							"subdomain": "test_sub_domain.com.",
							"domain":    "sub_domain.com",
						},

						{
							"path":      "/test",
							"subdomain": "test_sub_domain1.com.",
							"domain":    "sub_domain.com",
						},

						{
							"path":      "/test",
							"subdomain": "test_sub_domain2.com.",
							"domain":    "sub_domain.com",
						},
					},
					"site_id":                        "${alicloud_esa_site.resource_Site_test.id}",
					"json_response_enable":           "off",
					"description":                    "测试1",
					"waiting_room_type":              "default",
					"disable_session_renewal_enable": "off",
					"cookie_name":                    "__aliwaitingroom_example",
					"waiting_room_name":              "waitingroom_example",
					"queue_all_enable":               "off",
					"queuing_status_code":            "200",
					"custom_page_html":               "",
					"new_users_per_minute":           "200",
					"session_duration":               "5",
					"language":                       "zhcn",
					"total_active_users":             "300",
					"queuing_method":                 "fifo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "off",
					"host_name_and_path": []map[string]interface{}{

						{
							"path":      "/test1",
							"subdomain": "test_sub_domain3.com.",
							"domain":    "sub_domain1.com",
						},

						{
							"path":      "/test1",
							"subdomain": "test_sub_domain4.com.",
							"domain":    "sub_domain1.com",
						},
					},
					"json_response_enable":           "off",
					"description":                    "测试1",
					"waiting_room_type":              "default",
					"disable_session_renewal_enable": "off",
					"cookie_name":                    "__aliwaitingroom_example1",
					"waiting_room_name":              "waitingroom_example1",
					"queue_all_enable":               "off",
					"queuing_status_code":            "202",
					"custom_page_html":               "",
					"new_users_per_minute":           "200",
					"session_duration":               "5",
					"language":                       "zhcn",
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

var AliCloudESAWaitingRoomresource_WaitingRoom_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAWaitingRoomresource_WaitingRoom_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Case resource_WaitingRoom_rule_test
func TestAccAliCloudESAWaitingRoomresource_WaitingRoom_rule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_waiting_room.default"
	ra := resourceAttrInit(resourceId, AliCloudESAWaitingRoomresource_WaitingRoom_rule_testMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EsaServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEsaWaitingRoom")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAWaitingRoom%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudESAWaitingRoomresource_WaitingRoom_rule_testBasicDependence)
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
					"status": "off",
					"host_name_and_path": []map[string]interface{}{

						{
							"path":      "/test",
							"subdomain": "test_sub_domain.com.",
							"domain":    "sub_domain.com",
						},

						{
							"path":      "/test",
							"subdomain": "test_sub_domain1.com.",
							"domain":    "sub_domain.com",
						},

						{
							"path":      "/test",
							"subdomain": "test_sub_domain2.com.",
							"domain":    "sub_domain.com",
						},
					},
					"site_id":                        "${alicloud_esa_site.resource_Site_rule_test.id}",
					"json_response_enable":           "off",
					"description":                    "测试1",
					"waiting_room_type":              "default",
					"disable_session_renewal_enable": "off",
					"cookie_name":                    "__aliwaitingroom_example",
					"waiting_room_name":              "waitingroom_example",
					"queue_all_enable":               "off",
					"queuing_status_code":            "200",
					"custom_page_html":               "",
					"new_users_per_minute":           "200",
					"session_duration":               "5",
					"language":                       "zhcn",
					"total_active_users":             "300",
					"queuing_method":                 "fifo",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "off",
					"host_name_and_path": []map[string]interface{}{

						{
							"path":      "/test1",
							"subdomain": "test_sub_domain3.com.",
							"domain":    "sub_domain1.com",
						},

						{
							"path":      "/test1",
							"subdomain": "test_sub_domain4.com.",
							"domain":    "sub_domain1.com",
						},
					},
					"json_response_enable":           "off",
					"description":                    "测试1",
					"waiting_room_type":              "default",
					"disable_session_renewal_enable": "off",
					"cookie_name":                    "__aliwaitingroom_example1",
					"waiting_room_name":              "waitingroom_example1",
					"queue_all_enable":               "off",
					"queuing_status_code":            "202",
					"custom_page_html":               "",
					"new_users_per_minute":           "200",
					"session_duration":               "5",
					"language":                       "zhcn",
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

var AliCloudESAWaitingRoomresource_WaitingRoom_rule_testMap = map[string]string{
	"id": CHECKSET,
}

func AliCloudESAWaitingRoomresource_WaitingRoom_rule_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_rule_test" {
  site_name   = "chenxin0116.site"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

`, name)
}

// Test ESA WaitingRoom. <<< Resource test cases, automatically generated.
