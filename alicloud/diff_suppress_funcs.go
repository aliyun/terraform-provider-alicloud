package alicloud

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func httpHttpsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if listener_forward, ok := d.GetOk("listener_forward"); ok && listener_forward.(string) == string(OnFlag) {
		return true
	}
	if protocol, ok := d.GetOk("protocol"); ok && (Protocol(protocol.(string)) == Http || Protocol(protocol.(string)) == Https) {
		return false
	}
	return true
}

func httpDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if protocol, ok := d.GetOk("protocol"); ok && Protocol(protocol.(string)) == Http {
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
	if protocol, ok := d.GetOk("protocol"); ok && Protocol(protocol.(string)) == Https {
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
	if protocol, ok := d.GetOk("protocol"); ok && (Protocol(protocol.(string)) == Tcp || Protocol(protocol.(string)) == Udp) {
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
	if protocol, ok := d.GetOk("protocol"); ok && Protocol(protocol.(string)) == Tcp {
		return false
	}
	return true
}

func establishedTimeoutDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
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
	if (!httpDiff && okHc && health.(string) == string(OnFlag)) ||
		(okPro && Protocol(protocol.(string)) == Tcp && okType && checkType.(string) == string(HTTPHealthCheckType)) {
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

func dnsValueDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	switch d.Get("type") {
	case "NS", "MX", "CNAME", "SRV":
		new = strings.TrimSuffix(strings.TrimSpace(new), ".")
	}
	return old == new
}

func dnsPriorityDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("type").(string) != "MX"
}

func slbInternetDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if internet, ok := d.GetOkExists("internet"); ok && internet.(bool) {
		return true
	}
	if internet, ok := d.GetOkExists("address_type"); ok && internet.(string) == "internet" {
		return true
	}
	return false
}

func slbInternetChargeTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// Uniform all internet chare type value and be compatible with previous lower value.
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return !slbInternetDiffSuppressFunc(k, old, new, d)
}

func slbBandwidthDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if slbInternetDiffSuppressFunc(k, old, new, d) && strings.ToLower(d.Get("internet_charge_type").(string)) == strings.ToLower(string(PayByBandwidth)) {
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

func ecsInternetDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if max, ok := d.GetOk("internet_max_bandwidth_out"); ok && max.(int) > 0 {
		return false
	}
	return true
}

func csKubernetesMasterPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("master_instance_charge_type").(string) == "PostPaid" || !(d.Id() == "") && !d.Get("force_update").(bool)
}

func csKubernetesWorkerPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("worker_instance_charge_type").(string) == "PostPaid" || !(d.Id() == "") && !d.Get("force_update").(bool)
}

func csForceUpdate(k, old, new string, d *schema.ResourceData) bool {
	if d.Id() == "" {
		return false
	}
	return !d.Get("force_update").(bool)
}

func csForceUpdateSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// many cs args are not returning from the server
	// if this is a new resource, allow the diff
	// args with this suppress func will always suppress the diff, unless user specified force_update
	log.Printf("key %s, old %s, new %s, isnew %v, id %s", k, old, new, d.IsNewResource(), d.Id())
	return !(d.Id() == "") && !d.Get("force_update").(bool)
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
	if d.Get("enable_backup_log").(bool) {
		return false
	}
	if d.Get("log_backup").(bool) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v > d.Get("backup_retention_period").(int) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v > d.Get("retention_period").(int) {
		return false
	}
	return true
}

func enableBackupLogDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("enable_backup_log").(bool) {
		return false
	}
	if d.Get("log_backup").(bool) {
		return false
	}

	return true
}

func archiveBackupPeriodDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("enable_backup_log").(bool) {
		return false
	}
	if d.Get("log_backup").(bool) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v+730 >= d.Get("backup_retention_period").(int) {
		return false
	}
	if v, err := strconv.Atoi(new); err != nil && v+730 >= d.Get("retention_period").(int) {
		return false
	}

	return true
}

func PostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return strings.ToLower(d.Get("instance_charge_type").(string)) == "postpaid"
}
func PostPaidAndRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(d.Get("instance_charge_type").(string)) == "prepaid" && d.Get("auto_renew").(bool) {
		return false
	}
	return true
}
func ecsNotAutoRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("instance_charge_type").(string) == "PostPaid" {
		return true
	}
	if RenewalStatus(d.Get("renewal_status").(string)) == RenewAutoRenewal {
		return false
	}
	return true
}

func polardbPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("pay_type").(string) == "PrePaid" {
		return false
	}
	return true
}

func polardbPostPaidAndRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("pay_type").(string) == "PrePaid" && d.Get("renewal_status").(string) != string(RenewNotRenewal) {
		return false
	}
	return true
}

func adbPostPaidAndRenewDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("pay_type").(string) == "PrePaid" && d.Get("renewal_status").(string) != string(RenewNotRenewal) {
		return false
	}
	return true
}
func adbPostPaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("pay_type").(string) == "PrePaid" {
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
	if d.Get("instance_charge_type").(string) == "PostPaid" && d.Get("spot_strategy").(string) == "SpotWithPriceLimit" {
		return false
	}
	return true
}

func ecsSecurityGroupRulePortRangeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	protocol := d.Get("ip_protocol").(string)
	if protocol == "tcp" || protocol == "udp" {
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

func routerInterfaceAcceptsideDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("role").(string) == string(AcceptingSide)
}

func routerInterfaceVBRTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("role").(string) == string(AcceptingSide) {
		return true
	}
	if d.Get("router_type").(string) == string(VRouter) {
		return true
	}
	return false
}

func workerDataDiskSizeSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	_, ok := d.GetOk("worker_data_disk_category")
	return !ok || !(d.Id() == "") && !d.Get("force_update").(bool)
}

func imageIdSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// setting image_id is not recommended, but is needed by some users.
	// when image_id is left blank, server will set a random default to it, we only know the default value after creation.
	// we suppress diff here to prevent unintentional force new action.

	// if we want to change cluster's image_id to default, we have to find out what the default image_id is,
	// then fill that image_id in this field.
	return new == "" || !(d.Id() == "") && !d.Get("force_update").(bool)
}

func esVersionDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {

	oldVersion := strings.Split(strings.Split(old, "_")[0], ".")
	newVersion := strings.Split(strings.Split(new, "_")[0], ".")

	if len(oldVersion) >= 2 && len(newVersion) >= 2 {
		if oldVersion[0] == newVersion[0] && oldVersion[1] == newVersion[1] {
			return true
		}
	}

	return false
}

func vpnSslConnectionsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if enable_ssl, ok := d.GetOk("enable_ssl"); !ok || !enable_ssl.(bool) {
		return true
	}
	return false
}

func actiontrailRoleNmaeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if !d.IsNewResource() && strings.ToLower(old) != strings.ToLower(new) {
		return false
	}
	return true
}

func slbDeleteProtectionSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if PayType(d.Get("instance_charge_type").(string)) == PrePaid {
		return true
	}
	return false
}

func slbAddressIpVersionSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("internet"); ok && v.(bool) {
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

func slbAddressDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("internet"); ok && v.(bool) {
		return true
	}
	return false
}

func kmsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		return true
	}
	if v, ok := d.GetOk("account_password"); ok && v.(string) != "" {
		return true
	}
	return false
}

func sagDnatEntryTypeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("type").(string) != "Intranet" {
		return true
	}
	return false
}

func sagClientUserPasswordSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("user_name").(string) == "" {
		return true
	}
	return false
}
