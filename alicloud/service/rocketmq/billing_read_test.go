package rocketmq_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/rocketmq"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func billingQueryCount(cloud *rocketmqTraceOnFixture) int {
	count := 0
	for _, request := range cloud.captured(http.MethodPost, "/") {
		values, _ := url.ParseQuery(request.query + "&" + string(request.body))
		if values.Get("Action") == "QueryAvailableInstances" {
			count++
		}
	}
	return count
}

func assertNullRenewalState(t *testing.T, r *schema.Resource, d *schema.ResourceData, attributes ...string) {
	t.Helper()
	state, err := d.State().AttrsAsObjectValue(r.CoreConfigSchema().ImpliedType())
	if err != nil {
		t.Fatal(err)
	}
	for _, attribute := range attributes {
		if value := state.GetAttr(attribute); !value.IsNull() {
			t.Fatalf("%s must serialize as null for generated import configuration, got %#v", attribute, value)
		}
	}
}

func TestUnitRocketmqInstanceReadPaygDoesNotQueryRenewal(t *testing.T) {
	cloud, client := newRocketmqTraceOnFixture(t)
	cloud.paymentType = "PayAsYouGo"
	cloud.billingErrorCode = "InternalError"
	r := rocketmq.ResourceInstance()
	// Import has no configured renewal fields. The cloud's actual billing type
	// decides whether renewal information is applicable.
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"payment_type": "Subscription",
	})
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	assertNullRenewalState(t, r, d, "auto_renew_period", "auto_renew_period_unit")
	if err := r.Read(d, client); err != nil {
		t.Fatalf("PayAsYouGo refresh must not depend on the subscription billing API: %v", err)
	}
	if billingQueryCount(cloud) != 0 {
		t.Fatal("PayAsYouGo Read must not query BSS renewal information")
	}
	if d.Get("payment_type") != "PayAsYouGo" || d.Get("auto_renew") != false ||
		d.Get("auto_renew_period") != 0 || d.Get("auto_renew_period_unit") != "" {
		t.Fatalf("PayAsYouGo state contains inapplicable renewal attributes: %#v", d.State().Attributes)
	}
	if d.Id() != rocketmqTraceOnFixtureInstanceId {
		t.Fatal("an existing PayAsYouGo instance must retain its identity")
	}
	assertNullRenewalState(t, r, d, "auto_renew_period", "auto_renew_period_unit")
}

func TestUnitRocketmqInstanceReadSubscriptionRenewal(t *testing.T) {
	cloud, client := newRocketmqTraceOnFixture(t)
	cloud.paymentType = "Subscription"
	cloud.renewStatus = "AutoRenewal"
	cloud.renewalDuration = 6
	cloud.renewalDurationUnit = "Y"
	r := rocketmq.ResourceInstance()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	if err := r.Read(d, client); err != nil {
		t.Fatal(err)
	}
	if billingQueryCount(cloud) != 1 || d.Get("auto_renew") != true ||
		d.Get("auto_renew_period") != 6 || d.Get("auto_renew_period_unit") != "Year" {
		t.Fatalf("Subscription Read must preserve verified renewal settings: %#v", d.State().Attributes)
	}
}

func TestUnitRocketmqInstanceReadUnsetSubscriptionPeriodIsNull(t *testing.T) {
	cloud, client := newRocketmqTraceOnFixture(t)
	cloud.paymentType = "Subscription"
	r := rocketmq.ResourceInstance()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
	d.SetId(rocketmqTraceOnFixtureInstanceId)
	assertNullRenewalState(t, r, d, "auto_renew_period")
	if err := r.Read(d, client); err != nil {
		t.Fatal(err)
	}
	assertNullRenewalState(t, r, d, "auto_renew_period")
}

func TestUnitRocketmqInstanceReadSubscriptionBillingFailureRetainsState(t *testing.T) {
	for _, empty := range []bool{false, true} {
		t.Run(map[bool]string{false: "API error", true: "empty billing result"}[empty], func(t *testing.T) {
			cloud, client := newRocketmqTraceOnFixture(t)
			cloud.paymentType = "Subscription"
			cloud.billingEmpty = empty
			if !empty {
				cloud.billingErrorCode = "InternalError"
			}
			r := rocketmq.ResourceInstance()
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{})
			d.SetId(rocketmqTraceOnFixtureInstanceId)
			if err := r.Read(d, client); err == nil {
				t.Fatal("unconfirmed Subscription renewal information must remain an error")
			}
			if billingQueryCount(cloud) != 1 {
				t.Fatal("an empty result must not trigger a request to another account station")
			}
			if d.Id() != rocketmqTraceOnFixtureInstanceId {
				t.Fatal("billing lookup failure must not remove an instance confirmed by RocketMQ")
			}
		})
	}
}
