package alicloud

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_amqp_exchange", &resource.Sweeper{
		Name: "alicloud_amqp_exchange",
		F:    testSweepAmqpExchange,
	})
}

func testSweepAmqpExchange(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	instanceId := os.Getenv("ALICLOUD_AMQP_INSTANCE_ID")
	action := "ListExchanges"
	request := make(map[string]interface{})
	request["InstanceId"] = instanceId
	request["MaxResults"] = PageSizeLarge
	var response map[string]interface{}
	conn, err := client.NewOnsproxyClient()
	if err != nil {
		log.Println(WrapError(err))
		return nil
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-12-12"), StringPointer("AK"), request, nil, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_exchanges", action, AlibabaCloudSdkGoERROR))
			return nil
		}
		resp, err := jsonpath.Get("$.Data.Exchanges", response)
		if err != nil {
			log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Exchanges", response))
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			skip := true
			for _, prefixe := range prefixes {
				if strings.HasPrefix(fmt.Sprint(item["Name"]), prefixe) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[DEBUG] Skipping the resource %s", item["Name"])
			}

			action := "DeleteExchange"
			request := map[string]interface{}{
				"InstanceId": instanceId,
				"Exchange":   item["Name"],
			}

			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(3*time.Minute, func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-12"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			log.Println(WrapError(err))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	return nil
}

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
						"instance_id":        os.Getenv("ALICLOUD_AMQP_INSTANCE_ID"),
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
						"instance_id":        os.Getenv("ALICLOUD_AMQP_INSTANCE_ID"),
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
						"instance_id":        os.Getenv("ALICLOUD_AMQP_INSTANCE_ID"),
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
						"instance_id":        os.Getenv("ALICLOUD_AMQP_INSTANCE_ID"),
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
		resource "alicloud_amqp_virtual_host" "default" {
		  instance_id       = "%s"
		  virtual_host_name = var.name
		}
		`, name, os.Getenv("ALICLOUD_AMQP_INSTANCE_ID"))
}

var AmqpExchangeBasicMap = map[string]string{}
