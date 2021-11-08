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
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers(
		"alicloud_msc_sub_webhook",
		&resource.Sweeper{
			Name: "alicloud_msc_sub_webhook",
			F:    testSweepMscSubWebhook,
		})
}

func testSweepMscSubWebhook(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	prefixes := []string{
		"tftest",
	}
	action := "ListWebhooks"
	request := map[string]interface{}{}

	request["MaxResults"] = PageSizeLarge

	var response map[string]interface{}
	conn, err := client.NewMscopensubscriptionClient()
	if err != nil {
		log.Printf("[ERROR] %s get an error: %#v", action, err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2021-07-13"), StringPointer("AK"), request, nil, &runtime)
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
			log.Printf("[ERROR] %s get an error: %#v", action, err)
			return nil
		}

		resp, err := jsonpath.Get("$.Webhooks", response)
		if err != nil {
			log.Printf("[ERROR] Getting resource %s attribute by path %s failed!!! Body: %v.", "$.Webhooks", action, err)
			return nil
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})

			skip := true
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(item["WebhookName"].(string)), strings.ToLower(prefix)) {
					skip = false
				}
			}
			if skip {
				log.Printf("[INFO] Skipping Msc Sub Webhook: %s", item["WebhookName"].(string))
				continue
			}
			action := "DeleteWebhook"
			request := map[string]interface{}{
				"WebhookId": item["WebhookId"],
			}
			_, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-07-13"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				log.Printf("[ERROR] Failed to delete Msc Sub Webhook (%s): %s", item["WebhookName"].(string), err)
			}
			log.Printf("[INFO] Delete Msc Sub Webhook success: %s ", item["WebhookName"].(string))
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	return nil
}

func TestAccAlicloudMscSubWebhook_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_msc_sub_webhook.default"
	ra := resourceAttrInit(resourceId, AlicloudMscSubWebhookMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &MscOpenSubscriptionService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeMscSubWebhook")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := "tftest"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudMscSubWebhookBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALICLOUD_WAF_MSC_SUB_TOKEN")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook_name": "${var.name}",
					"server_url":   "https://oapi.dingtalk.com/robot/send?access_token=${var.token}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook_name": name,
						"server_url":   "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("ALICLOUD_WAF_MSC_SUB_TOKEN"),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook_name": "${var.name}update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook_name": name + "update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"server_url": "https://oapi.dingtalk.com/robot/send?access_token=${var.token}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_url": "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("ALICLOUD_WAF_MSC_SUB_TOKEN"),
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook_name": "${var.name}",
					"server_url":   "https://oapi.dingtalk.com/robot/send?access_token=${var.token}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook_name": name,
						"server_url":   "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("ALICLOUD_WAF_MSC_SUB_TOKEN"),
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"locale"},
			},
		},
	})
}

var AlicloudMscSubWebhookMap0 = map[string]string{
	"locale": NOSET,
}

func AlicloudMscSubWebhookBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
variable "token" {
  default = "%s"
}
`, name, os.Getenv("ALICLOUD_WAF_MSC_SUB_TOKEN"))
}
