package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAmqpQueue_basic(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_queue.default"
	ra := resourceAttrInit(resourceId, AmqpQueueBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpQueuebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpQueueConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":             "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":       "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":       "true",
					"auto_expire_state":       "10000",
					"dead_letter_exchange":    "",
					"dead_letter_routing_key": "",
					"exclusive_state":         "false",
					"max_length":              "100",
					"maximum_priority":        "10",
					"message_ttl":             "100",
					"queue_name":              name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":             CHECKSET,
						"virtual_host_name":       name,
						"auto_delete_state":       "true",
						"auto_expire_state":       "10000",
						"dead_letter_exchange":    "",
						"dead_letter_routing_key": "",
						"exclusive_state":         "false",
						"max_length":              "100",
						"maximum_priority":        "10",
						"message_ttl":             "100",
						"queue_name":              name,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_expire_state", "dead_letter_exchange", "dead_letter_routing_key", "max_length", "maximum_priority", "message_ttl"},
			},
		},
	})

}

func resourceAmqpQueueConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
 			default = "%v"
		}
		data "alicloud_amqp_instances" "default" {
			status = "SERVING"
		}
		resource "alicloud_amqp_virtual_host" "default" {
		  instance_id       = data.alicloud_amqp_instances.default.ids.0
		  virtual_host_name = var.name
		}
		`, name)
}

var AmqpQueueBasicMap = map[string]string{}
