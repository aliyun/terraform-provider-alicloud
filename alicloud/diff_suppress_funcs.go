package alicloud

import (
	"strconv"

	"strings"

	"reflect"
	"sort"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/rds"
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

func slbInternetDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if internet, ok := d.GetOk("internet"); ok && internet.(bool) {
		return true
	}
	return false
}

func slbInternetChargeTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return !slbInternetDiffSuppressFunc(k, old, new, d)
}

func slbInstanceSpecDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if old != "" && new == "" {
		return true
	}
	return false
}

func slbBandwidthDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if slbInternetDiffSuppressFunc(k, old, new, d) && slb.InternetChargeType(d.Get("internet_charge_type").(string)) == slb.PayByBandwidth {
		return false
	}
	return true
}

func ecsPrivateIpDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	vswitch := ""
	if vsw, ok := d.GetOk("vswitch_id"); ok && vsw.(string) != "" {
		vswitch = vsw.(string)
	} else if subnet, ok := d.GetOk("subnet_id"); ok && subnet.(string) != "" {
		vswitch = subnet.(string)
	}

	if vswitch != "" {
		return false
	}
	return true
}
func ecsInternetDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if max, ok := d.GetOk("internet_max_bandwidth_out"); ok && max.(int) > 0 {
		return false
	}
	return true
}

func ecsPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if common.InstanceChargeType(d.Get("instance_charge_type").(string)) == common.PrePaid {
		return false
	}
	return true
}

func ecsNotAutoRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if common.InstanceChargeType(d.Get("instance_charge_type").(string)) == common.PostPaid {
		return true
	}
	if RenewalStatus(d.Get("renewal_status").(string)) == RenewAutoRenewal {
		return false
	}
	return true
}

func zoneIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if vsw, ok := d.GetOk("vswitch_id"); ok && vsw.(string) != "" {
		return true
	} else if multi, ok := d.GetOk("multi_az"); ok && multi.(bool) {
		return true
	}
	return false
}

func logRetentionPeriodDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("log_backup").(bool) {
		return false
	}

	if v, err := strconv.Atoi(new); err != nil && v > d.Get("retention_period").(int) {
		return false
	}

	return true
}

func rdsPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if rds.DBPayType(d.Get("instance_charge_type").(string)) == rds.Prepaid {
		return false
	}
	return true
}

func ecsSpotStrategyDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("instance_charge_type").(string) == string(PostPaid) {
		return false
	}
	return true
}

func ecsSpotPriceLimitDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if common.InstanceChargeType(d.Get("instance_charge_type").(string)) == common.PostPaid &&
		ecs.SpotStrategyType(d.Get("spot_strategy").(string)) == ecs.SpotWithPriceLimit {
		return false
	}
	return true
}

func ecsSecurityGroupRulePortRangeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	protocol := ecs.IpProtocol(d.Get("ip_protocol").(string))
	if protocol == ecs.IpProtocolTCP || protocol == ecs.IpProtocolUDP {
		if new == AllPortRange {
			return true
		}
		return false
	}
	if new == AllPortRange {
		return false
	}
	return true
}

func vpcTypeResourceDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if len(Trim(d.Get("vswitch_id").(string))) > 0 {
		return false
	}
	return true
}

func cmsDimensionsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.IsNewResource() {
		return false
	}
	olds := strings.Split(old, COMMA_SEPARATED)
	sort.Strings(olds)
	news := strings.Split(new, COMMA_SEPARATED)
	sort.Strings(news)
	return reflect.DeepEqual(olds, news)
}
