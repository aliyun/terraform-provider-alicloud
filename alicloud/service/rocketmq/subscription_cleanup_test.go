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

// ACUBE_ROCKETMQ_REFUND_V1: the production RocketMQ instance Delete now
// unsubscribes prepaid (Subscription) instances through BssOpenApi
// RefundInstance (ProductCode=ons, ImmediatelyRelease=1, ProductType read
// from the DescribeInstance commodityCode) and confirms release before
// clearing state, instead of rejecting Subscription destroy. PayAsYouGo
// destroy keeps using the RocketMQ RoaDelete. Acceptance tests for real
// Subscription instances still install this provider so the durable cleanup
// coordinator can reconcile a refund when the Provider delete cannot confirm
// release: the production delete runs first, and the coordinator refund is
// only a fallback.
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
	originalDelete := r.Delete
	r.Delete = func(d *schema.ResourceData, meta interface{}) error {
		// The production delete now refunds Subscription instances via
		// BssOpenApi RefundInstance and releases PayAsYouGo instances via
		// RoaDelete; only fall back to the durable coordinator when the
		// Provider delete cannot confirm release.
		if err := originalDelete(d, meta); err == nil {
			return nil
		}
		return rocketmqCleanupCall("refund", d.Id())
	}
	p.ResourcesMap["alicloud_rocketmq_instance"] = r
	return p
}

func TestUnitRocketmqSubscriptionCoordinator(t *testing.T) {
	// The coordinator still guards account/region admission on Create and
	// records the begin lifecycle so acceptance cleanup can reconcile a refund
	// when the Provider delete cannot confirm release.
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
		})
	}
}

func TestUnitRocketmqSubscriptionCoordinatorFallback(t *testing.T) {
	// When the production refund cannot confirm release, the coordinator
	// refund is the fallback so acceptance cleanup still reconciles the
	// instance instead of leaking a prepaid resource.
	cloud, client := newRocketmqTraceOnFixture(t)
	cloud.paymentType = "Subscription"
	cloud.commodityCode = "ons_rmqpost_public_intl"
	cloud.refundErrorCode = "CommodityNotSupported"
	calls := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		calls = append(calls, body["action"])
		json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	}))
	defer server.Close()
	t.Setenv("ALICLOUD_ROCKETMQ_CLEANUP_URL", server.URL+"/")
	t.Setenv("ALICLOUD_ROCKETMQ_CLEANUP_TOKEN", strings.Repeat("a", 64))
	t.Setenv("ALICLOUD_ROCKETMQ_TEST_NAME", "rmq-fixture-fallback")
	t.Setenv("ALICLOUD_ROCKETMQ_TEST_ACCOUNT", "1234567890123456")
	t.Setenv("ALICLOUD_ROCKETMQ_TEST_REGION", "cn-hangzhou")
	t.Setenv("ALICLOUD_ROCKETMQ_TEST_CASE", "TestAccAliCloudRocketmqInstance_basic4128")
	r := rocketmqSubscriptionProvider().ResourcesMap["alicloud_rocketmq_instance"]
	raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture-fallback", nil)
	raw["payment_type"] = "Subscription"
	d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	if err := r.Delete(d, client); err != nil {
		t.Fatalf("coordinator fallback should reconcile the refund, got %v", err)
	}
	if len(cloud.capturedRefund()) != 1 {
		t.Fatal("production delete should attempt RefundInstance before the coordinator fallback")
	}
	if len(calls) == 0 || calls[len(calls)-1] != "refund" {
		t.Fatalf("expected coordinator refund fallback, got %v", calls)
	}
	if cloud.refunded {
		t.Fatal("production refund failure must not release the instance before the coordinator reconciles")
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

func TestUnitRocketmqInstanceSubscriptionDeleteRefund(t *testing.T) {
	// A Subscription destroy refunds through BssOpenApi RefundInstance with
	// ProductCode=ons, ImmediatelyRelease=1 and ProductType read from the
	// DescribeInstance commodityCode, then confirms release.
	cloud, client := newRocketmqTraceOnFixture(t)
	cloud.paymentType = "Subscription"
	cloud.commodityCode = "ons_rmqpost_public_intl"
	r := rocketmq.ResourceInstance()
	raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture-subscription", nil)
	raw["payment_type"] = "Subscription"
	d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	if err := r.Delete(d, client); err != nil {
		t.Fatalf("subscription delete should refund via RefundInstance, got %v", err)
	}
	refunds := cloud.capturedRefund()
	if len(refunds) != 1 {
		t.Fatalf("expected one RefundInstance call, got %d", len(refunds))
	}
	q := cloud.refundQuery()
	if got := q.Get("ProductCode"); got != "ons" {
		t.Errorf("ProductCode=%q, want ons", got)
	}
	if got := q.Get("ImmediatelyRelease"); got != "1" {
		t.Errorf("ImmediatelyRelease=%q, want 1", got)
	}
	if got := q.Get("ProductType"); got != "ons_rmqpost_public_intl" {
		t.Errorf("ProductType=%q, want commodityCode", got)
	}
	if got := q.Get("InstanceId"); got != rocketmqTraceOnFixtureInstanceId {
		t.Errorf("InstanceId=%q, want %q", got, rocketmqTraceOnFixtureInstanceId)
	}
	if q.Get("ClientToken") == "" {
		t.Error("ClientToken should be set for idempotent refund")
	}
	if !cloud.refunded {
		t.Error("instance should be released after a successful refund")
	}
	if len(cloud.captured(http.MethodDelete, "/instances/"+rocketmqTraceOnFixtureInstanceId)) != 0 {
		t.Error("subscription delete must not call the PayAsYouGo RoaDelete")
	}
}

func TestUnitRocketmqInstanceSubscriptionRefundIdempotent(t *testing.T) {
	for _, code := range []string{"ResourceNotExists", "ExistRefundingOrderError"} {
		t.Run(code, func(t *testing.T) {
			// An already-absent or already-refunding order is treated as a
			// successful delete so a repeated terraform destroy is idempotent.
			cloud, client := newRocketmqTraceOnFixture(t)
			cloud.paymentType = "Subscription"
			cloud.commodityCode = "ons_rmqpost_public_intl"
			cloud.refundErrorCode = code
			r := rocketmq.ResourceInstance()
			raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture-subscription", nil)
			raw["payment_type"] = "Subscription"
			d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
			d.SetId(rocketmqTraceOnFixtureInstanceId)
			if err := r.Delete(d, client); err != nil {
				t.Fatalf("idempotent refund code %s should succeed, got %v", code, err)
			}
			if len(cloud.capturedRefund()) != 1 {
				t.Fatalf("expected one RefundInstance call, got %d", len(cloud.capturedRefund()))
			}
			if !cloud.refunded {
				t.Errorf("instance should be released after idempotent refund %s", code)
			}
		})
	}
}

func TestUnitRocketmqInstanceSubscriptionRefundInternalErrorRetries(t *testing.T) {
	// A transient BssOpenApi InternalError (HTTP 500, 系统异常请重试) is
	// retried within the delete budget instead of failing the refund; the
	// second attempt succeeds and releases the instance.
	cloud, client := newRocketmqTraceOnFixture(t)
	cloud.paymentType = "Subscription"
	cloud.commodityCode = "ons_rmqpost_public_intl"
	cloud.refundErrorCode = "InternalError"
	r := rocketmq.ResourceInstance()
	raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture-internal", nil)
	raw["payment_type"] = "Subscription"
	d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	if err := r.Delete(d, client); err != nil {
		t.Fatalf("transient InternalError should retry to a successful refund, got %v", err)
	}
	refunds := cloud.capturedRefund()
	if len(refunds) != 2 {
		t.Fatalf("expected the transient InternalError to trigger a retry, got %d refund calls", len(refunds))
	}
	if !cloud.refunded {
		t.Error("instance should be released after the retried refund succeeds")
	}
}

func TestUnitRocketmqInstanceSubscriptionRefundFailureRetainsState(t *testing.T) {
	// A refund failure retains the Terraform state so a subsequent destroy can
	// retry the refund; the instance is not released and no RoaDelete is issued.
	cloud, client := newRocketmqTraceOnFixture(t)
	cloud.paymentType = "Subscription"
	cloud.commodityCode = "ons_rmqpost_public_intl"
	cloud.refundErrorCode = "CommodityNotSupported"
	r := rocketmq.ResourceInstance()
	raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture-subscription", nil)
	raw["payment_type"] = "Subscription"
	d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	if err := r.Delete(d, client); err == nil {
		t.Fatal("refund failure should retain the Terraform state")
	}
	if cloud.refunded {
		t.Fatal("a failed refund must not release the instance")
	}
	if d.Id() != rocketmqTraceOnFixtureInstanceId {
		t.Fatalf("state identity should be retained, got %q", d.Id())
	}
	if len(cloud.capturedRefund()) != 1 {
		t.Fatalf("expected one RefundInstance attempt, got %d", len(cloud.capturedRefund()))
	}
	if len(cloud.captured(http.MethodDelete, "/instances/"+rocketmqTraceOnFixtureInstanceId)) != 0 {
		t.Fatal("a failed refund must not fall through to the PayAsYouGo RoaDelete")
	}
}

func TestUnitRocketmqInstanceDeletePayg(t *testing.T) {
	// A PayAsYouGo destroy keeps using the RocketMQ RoaDelete; the RefundInstance
	// path is not taken.
	cloud, client := newRocketmqTraceOnFixture(t)
	r := rocketmq.ResourceInstance()
	raw := rocketmqTraceOnFixtureRawConfig("rmq-fixture-payg", nil)
	d := rocketmqTraceOnFixtureNewResourceData(t, r, raw, cty.NullVal(cty.Bool))
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	if err := r.Delete(d, client); err != nil {
		t.Fatalf("PayAsYouGo delete should release via RoaDelete, got %v", err)
	}
	if len(cloud.captured(http.MethodDelete, "/instances/"+rocketmqTraceOnFixtureInstanceId)) != 1 {
		t.Fatal("PayAsYouGo delete must call the RocketMQ RoaDelete")
	}
	if len(cloud.capturedRefund()) != 0 {
		t.Fatal("PayAsYouGo delete must not call RefundInstance")
	}
	if !cloud.refunded {
		t.Error("instance should be released after RoaDelete")
	}
}
