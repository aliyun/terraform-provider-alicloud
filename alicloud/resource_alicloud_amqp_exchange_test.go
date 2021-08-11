package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAmqpExchange_DIRECT(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

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
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "DIRECT",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "DIRECT",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func TestAccAlicloudAmqpExchange_FANOUT(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

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
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "FANOUT",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "FANOUT",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func TestAccAlicloudAmqpExchange_HEADERS(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

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
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "HEADERS",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "HEADERS",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func TestAccAlicloudAmqpExchange_TOPIC(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_exchange.default"
	ra := resourceAttrInit(resourceId, AmqpExchangeBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpExchangebasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpExchangeConfigDependence)

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
					"alternate_exchange": name + "alternate_exchange",
					"instance_id":        "${alicloud_amqp_virtual_host.default.instance_id}",
					"virtual_host_name":  "${alicloud_amqp_virtual_host.default.virtual_host_name}",
					"auto_delete_state":  "true",
					"exchange_name":      name,
					"exchange_type":      "TOPIC",
					"internal":           "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":        CHECKSET,
						"virtual_host_name":  name,
						"alternate_exchange": name + "alternate_exchange",
						"auto_delete_state":  "true",
						"exchange_name":      name,
						"exchange_type":      "TOPIC",
						"internal":           "false",
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"alternate_exchange", "internal"},
			},
		},
	})
}

func resourceAmqpExchangeConfigDependence(name string) string {
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

var AmqpExchangeBasicMap = map[string]string{}
