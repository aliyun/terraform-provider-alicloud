package alicloud

import (
	"github.com/denverdino/aliyungo/dns"
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/schema"
)

func httpHttpsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && (Protocol(protocol.(string)) == Http || Protocol(protocol.(string)) == Https) {
		return false
	}
	return true
}

func stickySessionTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	httpDiff := httpHttpsDiffSuppressFunc(k, old, new, d)
	if session, ok := d.GetOk("sticky_session"); !httpDiff && ok && slb.FlagType(session.(string)) == slb.OnFlag {
		return false
	}
	return true
}

func cookieTimeoutDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	stickSessionTypeDiff := stickySessionTypeDiffSuppressFunc(k, old, new, d)
	if session_type, ok := d.GetOk("sticky_session_type"); !stickSessionTypeDiff && ok && slb.StickySessionType(session_type.(string)) == slb.InsertStickySessionType {
		return false
	}
	return true
}

func cookieDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	stickSessionTypeDiff := stickySessionTypeDiffSuppressFunc(k, old, new, d)
	if session_type, ok := d.GetOk("sticky_session_type"); !stickSessionTypeDiff && ok && slb.StickySessionType(session_type.(string)) == slb.ServerStickySessionType {
		return false
	}
	return true
}

func tcpUdpDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && (Protocol(protocol.(string)) == Tcp || Protocol(protocol.(string)) == Udp) {
		return false
	}
	return true
}

func healthCheckDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	httpDiff := httpHttpsDiffSuppressFunc(k, old, new, d)
	if health, ok := d.GetOk("health_check"); httpDiff || (ok && slb.FlagType(health.(string)) == slb.OnFlag) {
		return false
	}
	return true
}

func healthCheckTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && Protocol(protocol.(string)) == Tcp {
		return false
	}
	return true
}
func httpHttpsTcpDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	httpDiff := httpHttpsDiffSuppressFunc(k, old, new, d)
	health, okHc := d.GetOk("health_check")
	protocol, okPro := d.GetOk("protocol")
	checkType, okType := d.GetOk("health_check_type")
	if (!httpDiff && okHc && slb.FlagType(health.(string)) == slb.OnFlag) ||
		(okPro && Protocol(protocol.(string)) == Tcp && okType && slb.HealthCheckType(checkType.(string)) == slb.HTTPHealthCheckType) {
		return false
	}
	return true
}
func sslCertificateIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && Protocol(protocol.(string)) == Https {
		return false
	}
	return true
}

func dnsPriorityDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if recordType, ok := d.GetOk("type"); ok && recordType.(string) == dns.MXRecord {
		return false
	}
	return true
}
