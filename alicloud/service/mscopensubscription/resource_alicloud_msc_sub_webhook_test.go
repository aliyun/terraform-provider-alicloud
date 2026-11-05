package mscopensubscription_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/mscopensubscription"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// mscSubTestClientForRegion delegates to the exported
// alicloud/sweep.SharedClientForRegion so the sweeper reuses the same
// per-region endpoint cache isolation (SharedEndpointCacheForRegion) and
// account metadata handling as the flat-suite sharedClientForRegion, instead
// of a lossy copy that shares one package-level endpoint cache across regions.
func mscSubTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	client, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return client.(*connectivity.AliyunClient), nil
}

func init() {
	resource.AddTestSweepers(
		"alicloud_msc_sub_webhook",
		&resource.Sweeper{
			Name: "alicloud_msc_sub_webhook",
			F:    testSweepMscSubWebhook,
		})
}

func testSweepMscSubWebhook(region string) error {
	client, err := mscSubTestClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	prefixes := []string{
		"tftest",
	}
	action := "ListWebhooks"
	request := map[string]interface{}{}

	request["MaxResults"] = helper.PageSizeLarge

	var response map[string]interface{}
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcGet("MscOpenSubscription", "2021-07-13", action, request, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
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
			_, err = client.RpcPost("MscOpenSubscription", "2021-07-13", action, nil, request, true)
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

func TestAccAliCloudMscSubWebhook_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_msc_sub_webhook.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudMscSubWebhookMap0)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return mscopensubscription.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeWebhook")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	name := "tftest"
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudMscSubWebhookBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithEnvVariable(t, "ALICLOUD_MSC_SUB_WEBHOOK_TOKEN")
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook_name": "${var.name}",
					"server_url":   "https://oapi.dingtalk.com/robot/send?access_token=${var.token}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook_name": name,
						"server_url":   "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("ALICLOUD_MSC_SUB_WEBHOOK_TOKEN"),
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
						"server_url": "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("ALICLOUD_MSC_SUB_WEBHOOK_TOKEN"),
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
						"server_url":   "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("ALICLOUD_MSC_SUB_WEBHOOK_TOKEN"),
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

var AlicloudMscSubWebhookMap0 = map[string]string{}

func AlicloudMscSubWebhookBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
variable "token" {
  default = "%s"
}
`, name, os.Getenv("ALICLOUD_MSC_SUB_WEBHOOK_TOKEN"))
}
