package alicloud

import (
	"fmt"
	"log"
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
	resource.AddTestSweepers("alicloud_arms_alert_contact", &resource.Sweeper{
		Name: "alicloud_arms_alert_contact",
		F:    testSweepArmsAlertContact,
	})
}

func testSweepArmsAlertContact(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return WrapErrorf(err, "error getting Alicloud client.")
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testacc",
	}

	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		action := "SearchAlertContact"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(2*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
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
			log.Printf("[ERROR] %s got an error: %s", action, err)
			return nil
		}
		resp, err := jsonpath.Get("$.PageBean.Contacts", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PageBean.Contacts", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			name := fmt.Sprint(item["ContactName"])
			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
					skip = false
					break
				}
			}
			if skip {
				log.Printf("[INFO] Skipping arms alert contact: %s ", name)
				continue
			}
			log.Printf("[INFO] delete arms alert contact: %s ", name)
			action = "DeleteAlertContact"
			request := map[string]interface{}{
				"ContactId": fmt.Sprint(item["ContactId"]),
				"RegionId":  client.RegionId,
			}

			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete alarm contact (%s): %s", name, err)
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	return nil
}

func TestAccAlicloudArmsAlertContact_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_arms_alert_contact.default"
	ra := resourceAttrInit(resourceId, ArmsAlertContactMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ArmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeArmsAlertContact")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccArmsAlertContact%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ArmsAlertContactBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name": "${var.name}",
					"email":              "hello.uuuu@aaa.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name": name,
						"email":              "hello.uuuu@aaa.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name": "${var.name}_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name": name + "_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email": "helloupdate.uuuu@aaa.com",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email": "helloupdate.uuuu@aaa.com",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"phone_num": "12345678900",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"phone_num": "12345678900",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"system_noc": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"system_noc": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_contact_name":     "${var.name}",
					"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f7",
					"email":                  "hello.uuuu@aaa.com",
					"phone_num":              "12345678901",
					"system_noc":             "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_contact_name":     name,
						"ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f7",
						"email":                  "hello.uuuu@aaa.com",
						"phone_num":              "12345678901",
						"system_noc":             "false",
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

var ArmsAlertContactMap = map[string]string{
	"ding_robot_webhook_url": "",
	"email":                  "",
	"phone_num":              "",
	"system_noc":             "false",
}

func ArmsAlertContactBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
