package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDatahubSubscription_basic(t *testing.T) {
	var v *datahub.GetSubscriptionResult

	resourceId := "alicloud_datahub_subscription.default"
	ra := resourceAttrInit(resourceId, datahubSubscriptionBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_project%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubSubscriptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name": "${alicloud_datahub_project.basic.name}",
					"topic_name":   "${alicloud_datahub_topic.basic.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": "subscription for basic.",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "subscription for basic.",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"comment": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"comment": "subscription added by terraform",
					}),
				),
			},
		},
	})
}
func TestAccAlicloudDatahubSubscription_multi(t *testing.T) {
	var v *datahub.GetSubscriptionResult

	resourceId := "alicloud_datahub_subscription.default.4"
	ra := resourceAttrInit(resourceId, datahubSubscriptionBasicMap)

	serviceFunc := func() interface{} {
		return &DatahubService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_project%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceDatahubSubscriptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"project_name": "${alicloud_datahub_project.basic.name}",
					"topic_name":   "${alicloud_datahub_topic.basic.name}",
					"count":        "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}
func resourceDatahubSubscriptionConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "project_name" {
	  default = "%s"
	}
	variable "topic_name" {
	  default = "tf_testacc_datahub_topic"
	}
	variable "record_type" {
	  default = "BLOB"
	}
	resource "alicloud_datahub_project" "basic" {
	  name = "${var.project_name}"
	  comment = "project for basic."
	}
	resource "alicloud_datahub_topic" "basic" {
	  project_name = "${alicloud_datahub_project.basic.name}"
	  name = "${var.topic_name}"
	  record_type = "${var.record_type}"
	  shard_count = 3
	  life_cycle = 7
	  comment = "topic for basic."
	}
	`, name)
}

var datahubSubscriptionBasicMap = map[string]string{
	"project_name": CHECKSET,
	"topic_name":   CHECKSET,
	"comment":      "subscription added by terraform",
}
