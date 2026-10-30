package slb

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The diff-suppress helpers below were moved into the slb package together
// with the listener, rule and load-balancer resources that use them. The
// schema fields they guard (protocol, sticky_session, health_check, acl,
// ssl certificate, payment_type, ...) are SLB-specific, and the package is
// not allowed to import suppressors from the root package. The original root
// definitions in alicloud/diff_suppress_funcs.go remain for any unmigrated
// root callers; the bodies are 1:1 copies with only the shared enum types
// (Protocol/Http/Https/Tcp/Udp) qualified to the helper leaf package.

func httpHttpsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if listener_forward, ok := d.GetOk("listener_forward"); ok && listener_forward.(string) == string(OnFlag) {
		return true
	}
	if protocol, ok := d.GetOk("protocol"); ok && (helper.Protocol(protocol.(string)) == helper.Http || helper.Protocol(protocol.(string)) == helper.Https) {
		return false
	}
	return true
}

func httpDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && helper.Protocol(protocol.(string)) == helper.Http {
		return false
	}
	return true
}

func forwardPortDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	httpDiff := httpDiffSuppressFunc(k, old, new, d)
	if listenerForward, ok := d.GetOk("listener_forward"); !httpDiff && ok && listenerForward.(string) == string(OnFlag) {
		return false
	}
	return true
}

func httpsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && helper.Protocol(protocol.(string)) == helper.Https {
		return false
	}
	return true
}

func stickySessionTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	httpDiff := httpHttpsDiffSuppressFunc(k, old, new, d)
	if session, ok := d.GetOk("sticky_session"); !httpDiff && ok && session.(string) == string(OnFlag) {
		return false
	}
	return true
}

func cookieTimeoutDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	stickSessionTypeDiff := stickySessionTypeDiffSuppressFunc(k, old, new, d)
	if session_type, ok := d.GetOk("sticky_session_type"); !stickSessionTypeDiff && ok && session_type.(string) == string(InsertStickySessionType) {
		return false
	}
	return true
}

func cookieDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	stickSessionTypeDiff := stickySessionTypeDiffSuppressFunc(k, old, new, d)
	if session_type, ok := d.GetOk("sticky_session_type"); !stickSessionTypeDiff && ok && session_type.(string) == string(ServerStickySessionType) {
		return false
	}
	return true
}

func tcpUdpDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && (helper.Protocol(protocol.(string)) == helper.Tcp || helper.Protocol(protocol.(string)) == helper.Udp) {
		return false
	}
	return true
}

func healthCheckDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	httpDiff := httpHttpsDiffSuppressFunc(k, old, new, d)
	if health, ok := d.GetOk("health_check"); httpDiff || (ok && health.(string) == string(OnFlag)) {
		return false
	}
	return true
}

func healthCheckTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && helper.Protocol(protocol.(string)) == helper.Tcp {
		return false
	}
	return true
}

func establishedTimeoutDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && helper.Protocol(protocol.(string)) == helper.Tcp {
		return false
	}
	return true
}

func httpHttpsTcpDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	httpDiff := httpHttpsDiffSuppressFunc(k, old, new, d)
	health, okHc := d.GetOk("health_check")
	protocol, okPro := d.GetOk("protocol")
	checkType, okType := d.GetOk("health_check_type")
	if (!httpDiff && okHc && health.(string) == string(OnFlag)) ||
		(okPro && helper.Protocol(protocol.(string)) == helper.Tcp && okType && checkType.(string) == string(HTTPHealthCheckType)) {
		return false
	}
	return true
}

func sslCertificateIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && helper.Protocol(protocol.(string)) == helper.Https {
		return false
	}
	return true
}

func slbAclDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if status, ok := d.GetOk("acl_status"); ok && status.(string) == string(OnFlag) {
		return false
	}
	return true
}

func slbServerCertificateDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if alicloudCertificateId, ok := d.GetOk("alicloud_certificate_id"); !ok || alicloudCertificateId.(string) == "" {
		return false
	}
	return true
}

func slbRuleStickySessionTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	listenerSync := slbRuleListenerSyncDiffSuppressFunc(k, old, new, d)
	if session, ok := d.GetOk("sticky_session"); !listenerSync && ok && session.(string) == string(OnFlag) {
		return false
	}
	return true
}

func slbRuleCookieTimeoutDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	stickSessionTypeDiff := slbRuleStickySessionTypeDiffSuppressFunc(k, old, new, d)
	if session_type, ok := d.GetOk("sticky_session_type"); !stickSessionTypeDiff && ok && session_type.(string) == string(InsertStickySessionType) {
		return false
	}
	return true
}

func slbRuleCookieDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	stickSessionTypeDiff := slbRuleStickySessionTypeDiffSuppressFunc(k, old, new, d)
	if session_type, ok := d.GetOk("sticky_session_type"); !stickSessionTypeDiff && ok && session_type.(string) == string(ServerStickySessionType) {
		return false
	}
	return true
}

func slbRuleHealthCheckDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	listenerSync := slbRuleListenerSyncDiffSuppressFunc(k, old, new, d)
	if health, ok := d.GetOk("health_check"); !listenerSync && ok && health.(string) == string(OnFlag) {
		return false
	}
	return true
}

func slbRuleListenerSyncDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if listenerSync, ok := d.GetOk("listener_sync"); ok && listenerSync.(string) == string(OffFlag) {
		return false
	}
	return true
}

// PostPaidDiffSuppressFunc is a generic shared suppressor also used by other
// products; the root copy in alicloud/diff_suppress_funcs.go remains for
// unmigrated callers. This local copy preserves the same behavior so the slb
// load-balancer resource does not import the root package.
func PostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// payment_type is the instance_charge_type's replacement.
	// If both instance_charge_type and payment_type are "", it means hiding a default "PostPaid"
	if v, ok := d.GetOk("instance_charge_type"); ok && strings.ToLower(v.(string)) == "prepaid" {
		return false
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		return false
	}
	return true
}
