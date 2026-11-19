package rocketmq_test

import (
	"net/http"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rocketmq"
	"github.com/hashicorp/go-cty/cty"
)

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
