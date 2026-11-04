package cloudfirewall_test

// This file ports two root-package test helpers that the migrated Cloud
// Firewall tests depend on but that were not migrated to alicloud/acctest or
// alicloud/sweep:
//   - sendMessage / DingTalkMessage: posts a cleanup reminder to a DingTalk
//     robot when a leftover instance blocks the integration test.
//   - testAccPreCheckForCleanUpInstances: polls BSS QueryAvailableInstances
//     until no Normal-status Cloud Firewall instance remains, alerting via
//     sendMessage otherwise. It mirrors the root helper but reaches the BSS
//     query through the local cloudfirewall.Service so it does not import the
//     root package.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	cloudfirewall "github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudfirewall"
	sweep "github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// DingTalkMessage is the text payload accepted by a DingTalk robot webhook.
type DingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func sendMessage(msg string) {
	// 钉钉机器人的 Webhook 地址
	webhookURL := "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("DINGTALK_WEBHOOK_ACCESS_TOKEN")

	// 构建消息内容
	message := DingTalkMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: msg,
		},
	}

	// 将消息内容转换为 JSON 格式
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("[ERROR] send dingTalk message failed. Error:", err)
		return
	}

	// 发送 POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("[ERROR] send dingTalk message failed. Error:", err)
		return
	}
	defer resp.Body.Close()
}

func testAccPreCheckForCleanUpInstances(t *testing.T, instanceRegion, productCode, productType, productCodeIntl, productTypeIntl string) {
	rawClient, err := sweep.SharedClientForRegion(acctest.DefaultRegionToTest)
	if err != nil {
		t.Errorf("error getting AliCloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)
	cloudfirewallService := cloudfirewall.NewService(client)
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
			instances, err := cloudfirewallService.QueryAvailableInstanceList(instanceRegion, productCode, productType, productCodeIntl, productTypeIntl)
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
