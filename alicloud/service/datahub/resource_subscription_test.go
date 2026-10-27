package datahub_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	datahubsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/datahub"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudDatahubSubscription_basic(t *testing.T) {
	var v *datahub.GetSubscriptionResult

	resourceId := "alicloud_datahub_subscription.default"
	ra := acctest.ResourceAttrInit(resourceId, datahubSubscriptionBasicMap)

	serviceFunc := func() interface{} {
		return datahubsvc.NewDatahubService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_project%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDatahubSubscriptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
					"comment": acctest.REMOVEKEY,
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
func TestAccAliCloudDatahubSubscription_multi(t *testing.T) {
	var v *datahub.GetSubscriptionResult

	resourceId := "alicloud_datahub_subscription.default.4"
	ra := acctest.ResourceAttrInit(resourceId, datahubSubscriptionBasicMap)

	serviceFunc := func() interface{} {
		return datahubsvc.NewDatahubService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)

	rac := acctest.ResourceAttrCheckInit(rc, ra)

	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100000, 999999)
	name := fmt.Sprintf("tf_testacc_datahub_project%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceDatahubSubscriptionConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, connectivity.DatahubSupportedRegions)
		},
		// module name
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
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
	"project_name": acctest.CHECKSET,
	"topic_name":   acctest.CHECKSET,
	"comment":      "subscription added by terraform",
}
