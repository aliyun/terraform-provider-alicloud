// Package waf_test holds the shared test helpers that the split-package waf
// tests cannot reach in the flat root layout: the package-level sweeper
// client and the clean-up-instance pre-check, both of which previously
// relied on the root sharedClientForRegion / testAccPreCheckForCleanUpInstances.
package waf_test

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

// wafTestClientForRegion delegates to the migrated shared sweeper client
// (alicloud/sweep.SharedClientForRegion) so the split-package waf tests reuse
// the same per-region endpoint cache (SharedEndpointCacheForRegion), account
// id/type metadata (ALICLOUD_ACCOUNT_ID/ALIBABA_CLOUD_ACCOUNT_ID and
// ALICLOUD_ACCOUNT_TYPE/ALIBABA_CLOUD_ACCOUNT_TYPE) and credential
// construction as the flat layout's sharedClientForRegion, instead of a lossy
// product-local copy that shared one endpoint map across regions and dropped
// the account metadata. The concrete type is returned as interface{} so callers
// retain the original type assertion.
func wafTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

// wafDingTalkMessage mirrors the root-package DingTalkMessage so the
// clean-up-instance alert keeps its original payload without referencing the
// flat layout.
type wafDingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// wafSendMessage mirrors the root-package sendMessage: it posts the alert to
// the dingtalk webhook configured via DINGTALK_WEBHOOK_ACCESS_TOKEN. The
// webhook token is optional; when absent the post is skipped silently, as in
// the flat layout.
func wafSendMessage(msg string) {
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("DINGTALK_WEBHOOK_ACCESS_TOKEN")
	if webhookURL == "https://oapi.dingtalk.com/robot/send?access_token=" {
		return
	}
	message := wafDingTalkMessage{
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

// wafPreCheckForCleanUpInstances mirrors the root-package
// testAccPreCheckForCleanUpInstances: it polls BssOpenApi for an available
// (Normal) waf subscription instance in instanceRegion and, while one is still
// present, alerts via dingtalk and waits so the acceptance tests do not race
// against a leftover paid instance. The BSS query is reproduced on the
// exported AliyunClient RPC entry because the root BssOpenApiService is not
// reachable from this package.
func wafPreCheckForCleanUpInstances(t *testing.T, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) {
	rawClient, err := wafTestClientForRegion(acctest.DefaultRegionToTest)
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
			instances, qerr := wafQueryAvailableInstanceList(client, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl)
			if qerr != nil {
				t.Errorf("error querying available instances: %s", qerr)
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
			wafSendMessage(fmt.Sprintf(`
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

// wafQueryAvailableInstanceList reproduces BssOpenApiService.QueryAvailableInstanceList
// on the exported AliyunClient RPC entry: it requests QueryAvailableInstances
// against BssOpenApi 2017-12-14, retrying throttling, and falls back to the
// international endpoint / product codes for a domestic account encountering
// NotApplicable, exactly as the root implementation does.
func wafQueryAvailableInstanceList(client *connectivity.AliyunClient, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) ([]interface{}, error) {
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
	var err error
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
		return nil, errs.WrapError(err)
	}
	v, err := jsonpath.Get("$.Data.InstanceList", response)
	if err != nil {
		return nil, errs.WrapError(err)
	}
	return v.([]interface{}), nil
}
