package alicloud

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/denverdino/aliyungo/cdn"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/dns"
	"github.com/denverdino/aliyungo/ram"
	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/yaml.v2"
)

// common
func validateInstancePort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 65535 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid port between 1 and 65535",
			k))
		return
	}
	return
}

func validateInstanceProtocol(v interface{}, k string) (ws []string, errors []error) {
	protocol := v.(string)
	if !isProtocolValid(protocol) {
		errors = append(errors, fmt.Errorf(
			"%q is an invalid value. Valid values are either http, https, tcp or udp",
			k))
		return
	}
	return
}

// ecs
func validateDiskCategory(v interface{}, k string) (ws []string, errors []error) {
	category := DiskCategory(v.(string))
	if _, ok := SupportedDiskCategory[category]; !ok {
		var valid []string
		for key := range SupportedDiskCategory {
			valid = append(valid, string(key))
		}
		errors = append(errors, fmt.Errorf("%s must be one of %s", k, strings.Join(valid, ", ")))
	}

	return
}

func validateInstanceName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateInstanceDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

func validateDiskName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value == "" {
		return
	}

	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateDiskDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

//security group
func validateSecurityGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateSecurityGroupDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 256 characters", k))

	}
	return
}

func validateSecurityRuleType(v interface{}, k string) (ws []string, errors []error) {
	rt := Direction(v.(string))
	if rt != DirectionIngress && rt != DirectionEgress {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, DirectionIngress, DirectionEgress))
	}

	return
}

func validateSecurityRuleIpProtocol(v interface{}, k string) (ws []string, errors []error) {
	pt := Protocol(v.(string))
	if pt != Tcp && pt != Udp && pt != Icmp && pt != Gre && pt != All {
		errors = append(errors, fmt.Errorf("%s must be one of %s, %s, %s, %s and %s", k,
			Tcp, Udp, Icmp, Gre, All))
	}

	return
}

func validateSecurityRuleNicType(v interface{}, k string) (ws []string, errors []error) {
	pt := GroupRuleNicType(v.(string))
	if pt != GroupRuleInternet && pt != GroupRuleIntranet {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, GroupRuleInternet, GroupRuleIntranet))
	}

	return
}

func validateSecurityRulePolicy(v interface{}, k string) (ws []string, errors []error) {
	pt := GroupRulePolicy(v.(string))
	if pt != GroupRulePolicyAccept && pt != GroupRulePolicyDrop {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s", k, GroupRulePolicyAccept, GroupRulePolicyDrop))
	}

	return
}

func validateSecurityPriority(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 100 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid authorization policy priority between 1 and 100",
			k))
		return
	}
	return
}

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
	}

	return
}

func validateRouteEntryNextHopType(v interface{}, k string) (ws []string, errors []error) {
	nht := NextHopType(v.(string))
	if nht != NextHopIntance && nht != NextHopRouterInterface && nht != NextHopHaVip &&
		nht != NextHopTunnel && nht != NextHopVpnGateway {
		errors = append(errors, fmt.Errorf("%s must be one of %s %s %s %s %s", k,
			NextHopIntance, NextHopRouterInterface, NextHopTunnel, NextHopHaVip, NextHopVpnGateway))
	}

	return
}

func validateSwitchCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, expected %q, got %q",
			k, ipnet, value))
		return
	}

	mark, _ := strconv.Atoi(strings.Split(ipnet.String(), "/")[1])
	if mark < 16 || mark > 29 {
		errors = append(errors, fmt.Errorf(
			"%q must contain a network CIDR which mark between 16 and 29",
			k))
	}

	return
}

// validateIoOptimized ensures that the string value is a valid IoOptimized that
// represents a IoOptimized - it adds an error otherwise
func validateIoOptimized(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		ioOptimized := OptimizedType(value)
		if ioOptimized != NoneOptimized &&
			ioOptimized != IOOptimized {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid IoOptimized, expected %s or %s, got %q",
				k, IOOptimized, NoneOptimized, ioOptimized))
		}
	}

	return
}

// validateInstanceNetworkType ensures that the string value is a classic or vpc
func validateInstanceNetworkType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		network := InstanceNetWork(value)
		if network != ClassicNet &&
			network != VpcNet {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceNetworkType, expected %s or %s, go %q",
				k, ClassicNet, VpcNet, network))
		}
	}
	return
}

func validateInstanceChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InstanceChargeType(value)
		if chargeType != common.PrePaid &&
			chargeType != common.PostPaid {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, common.PrePaid, common.PostPaid, chargeType))
		}
	}

	return
}

func validateInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InternetChargeType(value)
		if chargeType != common.PayByBandwidth &&
			chargeType != common.PayByTraffic {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, common.PayByBandwidth, common.PayByTraffic, chargeType))
		}
	}

	return
}

func validateInternetMaxBandWidthOut(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 100 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid internet bandwidth out between 0 and 100",
			k))
		return
	}
	return
}

func validateInstanceChargeTypePeriod(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if (value > 0 && value < 10) || (value > 11 && value < 61 && value%12 == 0) {
		return
	}
	errors = append(errors, fmt.Errorf(
		"%q must be a valid period, expected [1-9], 12, 24, 36, 48 or 60, got %d.", k, value))
	return
}
func validateInstanceChargeTypePeriodUnit(v interface{}, k string) (ws []string, errors []error) {
	unit := common.TimeType(v.(string))
	if unit != common.Week && unit != common.Month {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid PeriodUnit, expected %s or %s, got %s.",
			k, common.Week, common.Month, unit))
	}
	return
}

func validateInstanceStatus(v interface{}, k string) (ws []string, errors []error) {
	status := Status(v.(string))
	if status != Running && status != Stopped && status != Creating &&
		status != Starting && status != Stopping {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid status, expected %s or %s or %s or %s or %s, got %s.",
			k, Creating, Starting, Running, Stopping, Stopped, status))
	}
	return
}

// SLB
func validateSlbName(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if len(value) < 1 || len(value) > 80 {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid load balancer name characters between 1 and 80",
				k))
			return
		}
	}

	return
}

func validateSlbInternetChargeType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		chargeType := common.InternetChargeType(value)

		if chargeType != "paybybandwidth" &&
			chargeType != "paybytraffic" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid InstanceChargeType, expected %s or %s, got %q",
				k, "paybybandwidth", "paybytraffic", value))
		}
	}

	return
}

func validateSlbInstanceSpecType(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		specType := slb.LoadBalancerSpecType(value)
		validLoadBalancerSpec := []slb.LoadBalancerSpecType{slb.S1Small, slb.S2Small,
			slb.S2Medium, slb.S3Small, slb.S3Medium, slb.S3Large}

		for _, s := range validLoadBalancerSpec {
			if s == specType {
				return
			}
		}
		errors = append(errors, fmt.Errorf("%q must contain a valid LoadBalancerSpecType,"+
			" expected %#v, got %q", k, validLoadBalancerSpec, value))
	}
	return
}

func validateSlbListenerBandwidth(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if (value < 1 || value > 1000) && value != -1 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer bandwidth between 1 and 1000 or -1",
			k))
		return
	}
	return
}

func validateSlbListenerScheduler(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		scheduler := slb.SchedulerType(value)

		if scheduler != "wrr" && scheduler != "wlc" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid SchedulerType, expected %s or %s, got %q",
				k, "wrr", "wlc", value))
		}
	}

	return
}

func validateSlbListenerCookie(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if len(value) < 1 || len(value) > 200 {
			errors = append(errors, fmt.Errorf("%q cannot be longer than 200 characters", k))
		}
	}
	return
}

func validateSlbListenerCookieTimeout(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 86400 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer cookie timeout between 1 and 86400",
			k))
		return
	}
	return
}

func validateSlbListenerPersistenceTimeout(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 3600 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid load balancer persistence timeout between 0 and 86400",
			k))
		return
	}
	return
}

func validateSlbListenerHealthCheckDomain(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if value == "$_ip" {
			errors = append(errors, fmt.Errorf("%q value '$_ip' has been deprecated, and empty string will replace it.", k))
		}
		if reg := regexp.MustCompile(`^[\w\-.]{1,80}$`); !reg.MatchString(value) {
			errors = append(errors, fmt.Errorf("%q length is limited to 1-80 and only characters such as letters, digits, '-' and '.' are allowed", k))
		}
	}
	return
}

func validateSlbListenerHealthCheckUri(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if len(value) < 1 || len(value) > 80 {
			errors = append(errors, fmt.Errorf("%q cannot be longer than 80 characters", k))
		}
	}
	return
}

func validateSlbListenerHealthCheckConnectPort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < 1 || value > 65535 {
		if value != -520 {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid load balancer health check connect port between 1 and 65535 or -520",
				k))
			return
		}

	}
	return
}

func validateDBBackupPeriod(v interface{}, k string) (ws []string, errors []error) {
	days := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	value := v.(string)
	exist := false
	for _, d := range days {
		if value == d {
			exist = true
			break
		}
	}
	if !exist {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid backup period value should in array %#v, got %q",
			k, days, value))
	}

	return
}

func validateAllowedStringValue(ss []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := Trim(v.(string))
		existed := false
		for _, s := range ss {
			if s == value {
				existed = true
				break
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid string value should be in array %#v, got %q",
				k, ss, value))
		}
		return

	}
}

func validateAllowedSplitStringValue(ss []string, splitStr string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		existed := false
		tsList := strings.Split(value, splitStr)

		for _, ts := range tsList {
			existed = false
			for _, s := range ss {
				if ts == s {
					existed = true
					break
				}
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid string value should in %#v, got %q",
				k, ss, value))
		}
		return

	}
}

func validateAllowedIntValue(is []int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		existed := false
		for _, i := range is {
			if i == value {
				existed = true
				break
			}
		}
		if !existed {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid int value should be in array %#v, got %q",
				k, is, value))
		}
		return

	}
}

func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min || value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d and larger than %d. Current value is %d.", k, min, max, value))
		}
		return
	}
}

func validateStringLengthInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if len(value) < min || len(value) > max {
			errors = append(errors, fmt.Errorf(
				"%q length cannot be lower than %d and larger than %d. Current length is %d.", k, min, max, len(value)))
		}
		return
	}
}

//data source validate func
//data_source_alicloud_image
func validateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}
	return
}

func validateImageOwners(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		owners := ImageOwnerAlias(value)
		if owners != ImageOwnerSystem &&
			owners != ImageOwnerSelf &&
			owners != ImageOwnerOthers &&
			owners != ImageOwnerMarketplace &&
			owners != ImageOwnerDefault {
			errors = append(errors, fmt.Errorf(
				"%q must contain a valid Image owner , expected %s, %s, %s, %s or %s, got %q",
				k, ImageOwnerSystem, ImageOwnerSelf, ImageOwnerOthers, ImageOwnerMarketplace, ImageOwnerDefault, owners))
		}
	}
	return
}

func validateRegion(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		region := common.Region(value)
		var valid string
		for _, re := range common.ValidRegions {
			if region == re {
				return
			}
			valid = valid + ", " + string(re)
		}
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid Region ID , expected %#v, got %q",
			k, valid, value))

	}
	return
}

func validateForwardPort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "any" {
		valueConv, err := strconv.Atoi(value)
		if err != nil || valueConv < 1 || valueConv > 65535 {
			errors = append(errors, fmt.Errorf("%q must be a valid port between 1 and 65535 or any ", k))
		}
	}
	return
}

func validateOssBucketName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 3 || len(value) > 63 {
		errors = append(errors, fmt.Errorf("%q cannot be less than 3 and longer than 63 characters", k))
	}
	return
}

func validateOssBucketAcl(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		acls := oss.ACLType(value)
		if acls != oss.ACLPrivate && acls != oss.ACLPublicRead && acls != oss.ACLPublicReadWrite {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid ACL value , expected %s, %s or %s, got %q",
				k, oss.ACLPrivate, oss.ACLPublicRead, oss.ACLPublicReadWrite, acls))
		}
	}
	return
}

func validateOssBucketLifecycleRuleId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 255 characters", k))
	}
	return
}

func validateOssBucketDateTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}
	return
}

func validateOssBucketObjectServerSideEncryption(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if ServerSideEncryptionAes256 != value {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid value, expected %s", k, ServerSideEncryptionAes256))
	}
	return
}

func validateDomainName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if vp := strings.Split(value, "."); len(vp) > 1 {
		mainDomain := strings.Join(vp[:len(vp)-1], ".")
		if len(mainDomain) > 63 || len(mainDomain) < 1 {
			errors = append(errors, fmt.Errorf("Main domain cannot be longer than 63 characters or less than 1 character"))
		}
	}

	if strings.HasSuffix(value, ".sh") || strings.HasSuffix(value, ".tel") {
		errors = append(errors, fmt.Errorf("Domain ends with .sh or .tel is not supported."))
	}

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("Domain name is invalid, it can not starts or ends with '-'"))
	}
	return
}

func validateDomainRecordType(v interface{}, k string) (ws []string, errors []error) {
	// Valid Record types
	// A, NS, MX, TXT, CNAME, SRV, AAAA, REDIRECT_URL, FORWORD_URL
	validTypes := map[string]string{
		dns.ARecord:           "",
		dns.NSRecord:          "",
		dns.MXRecord:          "",
		dns.TXTRecord:         "",
		dns.CNAMERecord:       "",
		dns.SRVRecord:         "",
		dns.AAAARecord:        "",
		dns.RedirectURLRecord: "",
		dns.ForwordURLRecord:  "",
	}

	value := v.(string)
	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf("%q must be one of [A, NS, MX, TXT, CNAME, SRV, AAAA, REDIRECT_URL, FORWORD_URL]", k))
	}
	return
}

func validateDomainRecordPriority(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value > 10 || value < 1 {
		errors = append(errors, fmt.Errorf("%q value is 1-10.", k))
	}
	return
}

func validateRR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
		errors = append(errors, fmt.Errorf("RR is invalid, it can not starts or ends with '-'"))
	}

	if len(value) > 253 {
		errors = append(errors, fmt.Errorf("RR can not longer than 253 characters."))
	}

	for _, part := range strings.Split(value, ".") {
		if len(part) > 63 {
			errors = append(errors, fmt.Errorf("Each part of RR split with . can not longer than 63 characters."))
			return
		}
	}
	return
}

func validateDomainRecordLine(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "default" && value != "telecom" && value != "unicom" && value != "mobile" && value != "oversea" && value != "edu" {
		errors = append(errors, fmt.Errorf("Record parsing line must be one of [default, telecom, unicom, mobile, oversea, edu]."))
	}
	return
}

func validateKeyPairName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters and less than 2", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}

	return
}

func validateKeyPairPrefix(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 100 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 100 characters, name is limited to 128", k))
	}

	return
}

func validateRamName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q can not be longer than 64 characters.", k))
	}

	pattern := `^[a-zA-Z0-9\.@\-_]+$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%q is invalid.", k))
	}
	return
}

func validateRamDisplayName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	pattern := `^[a-zA-Z0-9\.@\-\p{Han}]{1,12}$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%q is invalid.", k))
	}
	return
}

func validateComment(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q can not be longer than 128 characters.", k))
	}
	return
}

func validateRamDesc(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 1024 {
		errors = append(errors, fmt.Errorf("%q can not be longer than 1024 characters.", k))
	}
	return
}

func validateRamPolicyName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q can not be longer than 128 characters.", k))
	}

	pattern := `^[a-zA-Z0-9\-]+$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%q is invalid.", k))
	}
	return
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeYamlString(yamlString interface{}) (string, error) {
	var j interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)

	err := yaml.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := yaml.Marshal(j)

	return string(bytes[:]), nil
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := json.Marshal(j)

	return string(bytes[:]), nil
}

func validateYamlString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeYamlString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
	}

	return
}

func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func validatePolicyType(v interface{}, k string) (ws []string, errors []error) {
	value := ram.Type(v.(string))

	if value != ram.System && value != ram.Custom {
		errors = append(errors, fmt.Errorf("%q must be '%s' or '%s'.", k, ram.System, ram.Custom))
	}
	return
}

func validateRamGroupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q can not be longer than 64 characters.", k))
	}

	pattern := `^[a-zA-Z0-9\-]+$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%q is invalid.", k))
	}
	return
}

func validateRamAlias(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 32 || len(value) < 2 {
		errors = append(errors, fmt.Errorf("%q can not be longer than 32 or less than 2 characters.", k))
	}
	return
}

func validateRamAKStatus(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value != string(Active) && value != string(Inactive) {
		errors = append(errors, fmt.Errorf("%q must be 'Active' or 'Inactive'.", k))
	}
	return
}

func validateContainerName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 1 || len(value) > 63 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 63 characters and less than 1", k))
	}
	reg := regexp.MustCompile("^[a-zA-Z0-9\u4E00-\u9FA5]{1}[a-zA-Z0-9\u4E00-\u9FA5-]{0,62}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s should be 1-63 characters long, and can contain numbers, Chinese characters, English letters and hyphens, but cannot start with hyphens.", k))
	}

	return
}

func validateContainerNamePrefix(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 37 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 37 characters, name is limited to 63", k))
	}

	reg := regexp.MustCompile("^[a-zA-Z0-9\u4E00-\u9FA5]?[a-zA-Z0-9\u4E00-\u9FA5-]{0,36}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s should be 0-37 characters long, and can contain numbers, Chinese characters, English letters and hyphens, but cannot start with hyphens.", k))
	}

	return
}

func validateContainerAppName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	reg := regexp.MustCompile("^[a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,63}$")
	if !reg.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s should be 1-64 characters long, and can contain numbers, English letters and hyphens, but cannot start with hyphens.", k))
	}

	return
}

func validateCdnChargeType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if value != "PayByTraffic" && value != "PayByBandwidth" {
		errors = append(errors, fmt.Errorf("%q must be 'PayByTraffic' or 'PayByBandwidth'.", k))
	}
	return
}

func validateCdnType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, val := range cdn.CdnTypes {
		if val == value {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v.", k, cdn.CdnTypes))
	return
}

func validateCdnSourceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, val := range cdn.SourceTypes {
		if val == value {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v.", k, cdn.SourceTypes))
	return
}

func validateCdnScope(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, val := range cdn.Scopes {
		if val == value {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v.", k, cdn.Scopes))
	return
}

func validateCdnSourcePort(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value != 80 && value != 443 {
		errors = append(errors, fmt.Errorf("%q must be one 80 or 443.", k))
	}
	return
}

func validateCdnHttpHeader(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, val := range cdn.HeaderKeys {
		if val == value {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v.", k, cdn.HeaderKeys))
	return
}

func validateCacheType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "suffix" && value != "path" {
		errors = append(errors, fmt.Errorf("%q must be 'suffix' or 'path'.", k))
	}
	return
}

func validateCdnEnable(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "on" && value != "off" {
		errors = append(errors, fmt.Errorf("%q must be 'on' or 'off'.", k))
	}
	return
}

func validateCdnHashKeyArg(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.Contains(value, ",") {
		errors = append(errors, fmt.Errorf("%q can not contains any ','.", k))
	}
	return
}

func validateCdnPage404Type(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "default" && value != "charity" && value != "other" {
		errors = append(errors, fmt.Errorf("%q must be one of ['default', 'charity', 'other'].", k))
	}
	return
}

func validateCdnRedirectType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "Off" && value != "Http" && value != "Https" {
		errors = append(errors, fmt.Errorf("%q must be one of ['Off', 'Http', 'Https'].", k))
	}
	return
}

func validateCdnReferType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "block" && value != "allow" {
		errors = append(errors, fmt.Errorf("%q must be 'block' or 'allow'.", k))
	}
	return
}

func validateCdnAuthType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "no_auth" && value != "type_a" && value != "type_b" && value != "type_c" {
		errors = append(errors, fmt.Errorf("%q must be one of ['no_auth', 'type_a', 'type_b', 'type_c']", k))
	}
	return
}

func validateCdnAuthKey(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	pattern := `^[a-zA-Z0-9]{6,32}$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%q can only consists of alphanumeric characters and can not be longer than 32 or less than 6 characters.", k))
	}
	return
}

func validatePolicyDocVersion(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "1" {
		errors = append(errors, fmt.Errorf("%q can only be '1' so far.", k))
	}
	return
}

func validateRouterInterfaceDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 2 || len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q cannot be less than 2 characters or longer than 256 characters", k))
	}

	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		errors = append(errors, fmt.Errorf("%s cannot starts with http:// or https://", k))
	}
	return
}

func validateInstanceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !strings.HasPrefix(value, "ecs.") {
		errors = append(errors, fmt.Errorf("Invalid %q: %s. It must be 'ecs.' as prefix.", k, value))
	}
	return
}

func validateDBConnectionPort(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil {
			errors = append(errors, err)
		}
		if port < 3001 || len(value) > 3999 {
			errors = append(errors, fmt.Errorf("%q cannot be less than 3001 and larger than 3999.", k))
		}
	}
	return
}

func validateInstanceSpotStrategy(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		spot := SpotStrategyType(value)
		if spot != NoSpot && spot != SpotAsPriceGo && spot != SpotWithPriceLimit {
			errors = append(errors, fmt.Errorf(
				"%q must be a valid Spot Strategy value , expected %s, %s or %s, got %q",
				k, NoSpot, SpotAsPriceGo, SpotWithPriceLimit, spot))
		}
	}
	return
}

func validateDBConnectionPrefix(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if len(value) < 1 || len(value) > 31 {
			errors = append(errors, fmt.Errorf("%q cannot be less than 1 and larger than 30.", k))
		}
	}
	return
}

func validateDBInstanceName(v interface{}, k string) (ws []string, errors []error) {
	if value := v.(string); value != "" {
		if len(value) < 2 || len(value) > 256 {
			errors = append(errors, fmt.Errorf("%q cannot be less than 1 and larger than 30.", k))
		}
	}
	return
}

func validateKmsKeyStatus(v interface{}, k string) (ws []string, errors []error) {
	status := KeyState(v.(string))
	if status != Enabled && status != Disabled && status != PendingDeletion {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid status, expected %s or %s or %s, got %s.",
			k, Enabled, Disabled, PendingDeletion, status))
	}
	return
}

func validateNatGatewaySpec(v interface{}, k string) (ws []string, errors []error) {
	spec := NatGatewaySpec(v.(string))
	if spec != NatGatewaySmallSpec && spec != NatGatewayMiddleSpec && spec != NatGatewayLargeSpec {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid specification, expected %s or %s or %s, got %s.",
			k, NatGatewaySmallSpec, NatGatewayMiddleSpec, NatGatewayLargeSpec, spec))
	}
	return
}
