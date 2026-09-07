package alicloud

import (
	"fmt"
	"testing"

	ali_mns "github.com/aliyun/aliyun-mns-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudMnsTopicSubscription_basic(t *testing.T) {
	var v *ali_mns.SubscriptionAttribute
	resourceId := "alicloud_mns_topic_subscription.default"
	ra := resourceAttrInit(resourceId, mnsTopicSubscriptionMap)
	serviceFunc := func() interface{} {
		return &MnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
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
					"filter_tag":            "tf-test",
					"notify_content_format": "SIMPLIFIED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       name,
						"topic_name": name,

						"endpoint": "http://www.test.com/test",
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

func TestAccAliCloudMnsTopicSubscription_queue(t *testing.T) {
	var v *ali_mns.SubscriptionAttribute
	resourceId := "alicloud_mns_topic_subscription.default"
	ra := resourceAttrInit(resourceId, mnsTopicSubscriptionMap)
	serviceFunc := func() interface{} {
		return &MnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
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
					"endpoint":              "acs:mns:" + defaultRegionToTest + ":1511928242963727:queues/${alicloud_mns_queue.default.name}",
					"filter_tag":            "tf-test",
					"notify_content_format": "SIMPLIFIED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":       name,
						"topic_name": name,
						"endpoint":   CHECKSET,
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

func TestAccAliCloudMnsTopicSubscription_multi(t *testing.T) {
	var v *ali_mns.SubscriptionAttribute
	resourceId := "alicloud_mns_topic_subscription.default.4"
	ra := resourceAttrInit(resourceId, mnsTopicSubscriptionMap)
	serviceFunc := func() interface{} {
		return &MnsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
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

	resource "alicloud_mns_queue" "default" {
	  name                     = "${var.name}"
	  delay_seconds            = 0
	  maximum_message_size     = 65536
	  message_retention_period = 345600
	  visibility_timeout       = 30
	  polling_wait_seconds     = 0
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
