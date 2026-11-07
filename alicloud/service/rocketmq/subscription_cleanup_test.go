package rocketmq_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"

	"github.com/aliyun/terraform-provider-alicloud/alicloud"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rocketmq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ACUBE_ROCKETMQ_REFUND_V1: only these tests use the runner's durable cleanup
// coordinator. The production Provider continues to reject subscription delete.
func rocketmqSubscriptionName(fallback string) string {
	if name := os.Getenv("ALICLOUD_ROCKETMQ_TEST_NAME"); name != "" {
		return name
	}
	return fallback
}

func preCheckRocketmqSubscription(t *testing.T) {
	t.Helper()
	if os.Getenv("ALICLOUD_ROCKETMQ_CLEANUP_URL") == "" {
		acctest.PreCheckPrePaidResources(t)
		t.Fatal("subscription acceptance requires the ACube refund coordinator; do not purchase without verified cleanup")
	}
	if err := rocketmqCleanupCall("ready", ""); err != nil {
		t.Fatal(err)
	}
}

func rocketmqCleanupCall(action, id string) error {
	endpoint, err := url.Parse(os.Getenv("ALICLOUD_ROCKETMQ_CLEANUP_URL"))
	if err != nil || endpoint.Scheme != "http" || endpoint.Hostname() != "127.0.0.1" || endpoint.User != nil {
		return fmt.Errorf("invalid subscription cleanup coordinator")
	}
	token := os.Getenv("ALICLOUD_ROCKETMQ_CLEANUP_TOKEN")
	if len(token) < 32 {
		return fmt.Errorf("missing subscription cleanup authorization")
	}
	body, _ := json.Marshal(map[string]string{"action": action, "case": os.Getenv("ALICLOUD_ROCKETMQ_TEST_CASE"), "instanceId": id})
	req, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot build subscription cleanup request")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 20 * time.Minute, CheckRedirect: func(*http.Request, []*http.Request) error { return fmt.Errorf("redirect forbidden") }}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("subscription cleanup coordinator unavailable; retain original receipt")
	}
	defer resp.Body.Close()
	var result struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4096)).Decode(&result); err != nil || resp.StatusCode != http.StatusOK || !result.OK {
		return fmt.Errorf("subscription cleanup unconfirmed: %s", result.Error)
	}
	return nil
}

func rocketmqSubscriptionProvider() *schema.Provider {
	p := alicloud.Provider()
	r := rocketmq.ResourceInstance()
	originalCreate := r.Create
	r.Create = func(d *schema.ResourceData, meta interface{}) error {
		client := meta.(*connectivity.AliyunClient)
		account, err := client.AccountId()
		if err != nil || account != os.Getenv("ALICLOUD_ROCKETMQ_TEST_ACCOUNT") || client.RegionId != os.Getenv("ALICLOUD_ROCKETMQ_TEST_REGION") {
			return fmt.Errorf("subscription cleanup account or region does not match the Provider")
		}
		if d.Get("payment_type") != "Subscription" || d.Get("instance_name") != os.Getenv("ALICLOUD_ROCKETMQ_TEST_NAME") {
			return fmt.Errorf("subscription cleanup intent does not match the test resource")
		}
		if err := rocketmqCleanupCall("begin", ""); err != nil {
			return err
		}
		err = originalCreate(d, meta)
		if d.Id() != "" {
			if receiptErr := rocketmqCleanupCall("created", d.Id()); receiptErr != nil && err == nil {
				return receiptErr
			}
		}
		return err
	}
	r.Delete = func(d *schema.ResourceData, _ interface{}) error {
		return rocketmqCleanupCall("refund", d.Id())
	}
	p.ResourcesMap["alicloud_rocketmq_instance"] = r
	return p
}

func TestUnitRocketmqSubscriptionCoordinator(t *testing.T) {
	for _, admitted := range []bool{false, true} {
		t.Run(fmt.Sprint("admitted=", admitted), func(t *testing.T) {
			cloud, client := newRocketmqTraceOnFixture(t)
			calls := []string{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer "+strings.Repeat("a", 64) {
					t.Error("missing coordinator authentication")
				}
				var body map[string]string
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				calls = append(calls, body["action"])
				if !admitted {
					w.WriteHeader(http.StatusConflict)
				}
				json.NewEncoder(w).Encode(map[string]interface{}{"ok": admitted, "error": "fixture rejection"})
			}))
			defer server.Close()
			t.Setenv("ALICLOUD_ROCKETMQ_CLEANUP_URL", server.URL+"/")
			t.Setenv("ALICLOUD_ROCKETMQ_CLEANUP_TOKEN", strings.Repeat("a", 64))
			t.Setenv("ALICLOUD_ROCKETMQ_TEST_NAME", "rmq-fixture-subscription")
			t.Setenv("ALICLOUD_ROCKETMQ_TEST_ACCOUNT", "1234567890123456")
			t.Setenv("ALICLOUD_ROCKETMQ_TEST_REGION", "cn-hangzhou")
			t.Setenv("ALICLOUD_ROCKETMQ_TEST_CASE", "TestAccAliCloudRocketmqInstance_basic4128")
			r := rocketmqSubscriptionProvider().ResourcesMap["alicloud_rocketmq_instance"]
			raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture-subscription", nil)
			raw["payment_type"] = "Subscription"
			d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
			d.SetId("")
			if err := r.Create(d, client); err == nil {
				t.Fatal("fixture create should fail")
			}
			count := len(cloud.captured(http.MethodPost, "/instances"))
			if admitted && count != 1 || !admitted && count != 0 {
				t.Fatalf("create requests=%d, admitted=%t", count, admitted)
			}
			if len(calls) != 1 || calls[0] != "begin" {
				t.Fatalf("unexpected coordinator calls %v", calls)
			}
			d.SetId(rocketmqTraceOnFixtureInstanceId)
			err := r.Delete(d, client)
			if (err == nil) != admitted {
				t.Fatalf("delete success=%t, admitted=%t", err == nil, admitted)
			}
			if calls[len(calls)-1] != "refund" {
				t.Fatal("delete did not verify refund")
			}
			if d.Id() != rocketmqTraceOnFixtureInstanceId {
				t.Fatal("unconfirmed delete must retain state identity")
			}
			ordinary := rocketmq.ResourceInstance()
			if err := ordinary.Delete(d, client); err == nil {
				t.Fatal("production subscription delete protection changed")
			}
		})
	}
}

func TestUnitRocketmqSubscriptionRejectsForeignAccount(t *testing.T) {
	fixture, client := newRocketmqTraceOnFixture(t)
	t.Setenv("ALICLOUD_ROCKETMQ_TEST_ACCOUNT", "9999999999999999")
	t.Setenv("ALICLOUD_ROCKETMQ_TEST_REGION", "cn-hangzhou")
	r := rocketmqSubscriptionProvider().ResourcesMap["alicloud_rocketmq_instance"]
	raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture", nil)
	raw["payment_type"] = "Subscription"
	d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
	d.SetId("")
	if err := r.Create(d, client); err == nil || !strings.Contains(err.Error(), "account or region") {
		t.Fatalf("expected account guard, got %v", err)
	}
	if len(fixture.captured(http.MethodPost, "/instances")) != 0 {
		t.Fatal("mismatched account must not buy a subscription")
	}
}
