package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMnsTopicSubscription_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mns_topic_subscription.default"
	ra := resourceAttrInit(resourceId, mnsTopicSubscriptionMap)
	serviceFunc := func() interface{} {
		return &MnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeMessageServiceSubscription")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMNSTopicSubscriptionConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMnsTopicSubscriptionConfigDependence)
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
					"name":                  name,
					"topic_name":            "${alicloud_mns_topic.default.name}",
					"endpoint":              "http://www.test.com/test",
					"notify_content_format": "SIMPLIFIED",
					"filter_tag":            "tf-test",
					"push_type":             "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       name,
						"topic_name": name,
						"filter_tag": "tf-test",
						"endpoint":   "http://www.test.com/test",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"push_type"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_strategy": "EXPONENTIAL_DECAY_RETRY",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_strategy": "EXPONENTIAL_DECAY_RETRY",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_strategy": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_strategy": "BACKOFF_RETRY",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudMnsTopicSubscription_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_mns_topic_subscription.default.4"
	ra := resourceAttrInit(resourceId, mnsTopicSubscriptionMap)
	serviceFunc := func() interface{} {
		return &MnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeMessageServiceSubscription")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccMNSTopicSubscriptionConfig-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMnsTopicSubscriptionConfigDependence)

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
					"name":                  name + "${count.index}",
					"topic_name":            "${alicloud_mns_topic.default.name}",
					"endpoint":              "http://www.test.com/test${count.index}",
					"filter_tag":            "tf-test",
					"notify_content_format": "SIMPLIFIED",
					"count":                 "5",
					"push_type":             "http",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func resourceMnsTopicSubscriptionConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	resource "alicloud_mns_topic" "default"{
		name="${var.name}"
		maximum_message_size=12357
		logging_enabled=true
	}
	`, name)
}

var mnsTopicSubscriptionMap = map[string]string{
	"name":                  CHECKSET,
	"topic_name":            CHECKSET,
	"endpoint":              CHECKSET,
	"filter_tag":            "tf-test",
	"notify_content_format": "SIMPLIFIED",
}
