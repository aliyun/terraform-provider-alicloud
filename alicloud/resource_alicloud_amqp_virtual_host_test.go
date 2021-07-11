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
	resource.AddTestSweepers("alicloud_amqp_virtual_host", &resource.Sweeper{
		Name: "alicloud_amqp_virtual_host",
		F:    testSweepAmqpVirtualHost,
	})
}

func testSweepAmqpVirtualHost(region string) error {
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
	action := "ListVirtualHosts"
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
			log.Println(WrapErrorf(err, DataDefaultErrorMsg, "alicloud_amqp_virtual_hosts", action, AlibabaCloudSdkGoERROR))
			return nil
		}
		resp, err := jsonpath.Get("$.Data.VirtualHosts", response)
		if err != nil {
			log.Println(WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.VirtualHosts", response))
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

			action := "DeleteVirtualHost"
			request := map[string]interface{}{
				"InstanceId":  instanceId,
				"VirtualHost": item["Name"],
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

func TestAccAlicloudAmqpVirtualHost_basic(t *testing.T) {

	var v map[string]interface{}
	resourceId := "alicloud_amqp_virtual_host.default"
	ra := resourceAttrInit(resourceId, AmqpVirtualHostBasicMap)
	serviceFunc := func() interface{} {
		return &AmqpOpenService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testacc-AmqpVirtualHostbasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceAmqpVirtualHostConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			//testAccPreCheckWithRegions(t, true, connectivity.AlikafkaSupportedRegions)
			testAccPreCheck(t)
			//testAccPreCheckWithNoDefaultVswitch(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_id":       os.Getenv("ALICLOUD_AMQP_INSTANCE_ID"),
					"virtual_host_name": "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_id":       os.Getenv("ALICLOUD_AMQP_INSTANCE_ID"),
						"virtual_host_name": name,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

}

func resourceAmqpVirtualHostConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
 			default = "%v"
		}
		`, name)
}

var AmqpVirtualHostBasicMap = map[string]string{}
