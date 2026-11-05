package esa_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	esa "github.com/aliyun/terraform-provider-alicloud/alicloud/service/esa"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ESA WaitingRoomRule. >>> Resource test cases, automatically generated.
// Case resource_WaitingRoomRule_test
func TestAccAliCloudESAWaitingRoomRuleresource_WaitingRoomRule_test(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_esa_waiting_room_rule.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudESAWaitingRoomRuleresource_WaitingRoomRule_testMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return esa.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEsaWaitingRoomRule")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sESAWaitingRoomRule%d", acctest.DefaultRegionToTest, rand)

	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudESAWaitingRoomRuleresource_WaitingRoomRule_testBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"status":          "off",
					"site_id":         "${alicloud_esa_site.resource_Site_rule_test.id}",
					"rule":            "(http.host eq \\\"video.example.com\\\")",
					"waiting_room_id": "${alicloud_esa_waiting_room.resource_WaitingRoom_rule_test.waiting_room_id}",
					"rule_name":       "WaitingRoomRule_example1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":    "off",
					"rule":      "(http.host eq \\\"video.example1.com\\\")",
					"rule_name": "WaitingRoomRule_example1",
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

var AliCloudESAWaitingRoomRuleresource_WaitingRoomRule_testMap = map[string]string{
	"id": acctest.CHECKSET,
}

func AliCloudESAWaitingRoomRuleresource_WaitingRoomRule_testBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "resource_Site_rule_test" {
  site_name   = "${lower(var.name)}.com"
  instance_id = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage    = "overseas"
  access_type = "NS"
}

resource "alicloud_esa_waiting_room" "resource_WaitingRoom_rule_test" {
  status                         = "off"
  site_id                        = alicloud_esa_site.resource_Site_rule_test.id
  json_response_enable           = "off"
  description                    = "测试1"
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

// Test ESA WaitingRoomRule. <<< Resource test cases, automatically generated.
