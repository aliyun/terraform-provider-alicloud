package wafopenapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// dingTalkMessage mirrors the root suite's DingTalkMessage so the cleanup
// pre-check can notify humans exactly as the original does.
type dingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func sendMessage(msg string) {
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("DINGTALK_WEBHOOK_ACCESS_TOKEN")
	message := dingTalkMessage{
		MsgType: "text",
	}
	message.Text.Content = msg
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("[ERROR] send dingTalk message failed. Error:", err)
		return
	}
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("[ERROR] send dingTalk message failed. Error:", err)
		return
	}
	defer resp.Body.Close()
}

// queryAvailableInstanceList mirrors BssOpenApiService.QueryAvailableInstanceList,
// which is root-private and not importable from a service test package. It
// reproduces the BssOpenApi QueryAvailableInstances call with the same
// international fallback semantics.
func queryAvailableInstanceList(client *connectivity.AliyunClient, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) (object []interface{}, err error) {
	var response map[string]interface{}
	var endpoint string
	action := "QueryAvailableInstances"
	request := map[string]interface{}{
		"ProductCode":      productCode,
		"ProductType":      productType,
		"SubscriptionType": "Subscription",
	}
	if client.IsInternationalAccount() {
		request = map[string]interface{}{
			"ProductCode": productCodeIntl,
			"ProductType": productTypeIntl,
		}
	}
	if instanceRegion != "" {
		request["Region"] = instanceRegion
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			if !client.IsInternationalAccount() && errs.IsExpectedErrors(err, []string{"NotApplicable"}) {
				request["ProductType"] = productTypeIntl
				request["ProductCode"] = productCodeIntl
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return object, errs.WrapError(err)
	}
	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return object, errs.WrapError(err)
	}
	return v.([]interface{}), nil
}

// testAccPreCheckForCleanUpInstances mirrors the root suite's cleanup pre-check,
// which is root-private and not importable from a service test package. It
// waits for pre-existing paid instances of the product to be released before
// running the acceptance test, notifying via DingTalk exactly as the original.
func testAccPreCheckForCleanUpInstances(t *testing.T, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) {
	rawClient, err := sweep.SharedClientForRegion(acctest.DefaultRegionToTest)
	if err != nil {
		t.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	accountId, err := client.AccountId()
	if err != nil {
		t.Errorf("error getting AliCloud client: %s", err)
	}
	deadline := time.Now().Add(30 * time.Minute)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if time.Now().After(deadline) {
				log.Println("Deadline reached, stopping waiting.")
				return
			}
			instances, err := queryAvailableInstanceList(client, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl)
			if err != nil {
				t.Errorf("error querying available instances: %s", err)
				return
			}
			instanceId := ""
			for _, instance := range instances {
				v := instance.(map[string]interface{})
				if v["Status"].(string) != "Normal" {
					continue
				}
				instanceId = v["InstanceID"].(string)
			}
			if instanceId == "" {
				return
			}
			t.Logf("waiting for cleaning up existed instances %s in region %s...", instanceId, instanceRegion)
			sendMessage(fmt.Sprintf(`
[Critical] Please cleaning up instance before running integration test.

AccountId: %s
ProductCode: %s
InstanceId: %s
Time: %s
`, accountId, productCode, instanceId, time.Now().String()))
			time.Sleep(3 * time.Minute)
		}
	}
}
