package ecs

// Local behavior-preserving compat copies of root-package shared declarations
// that ECS references bare. Root alicloud cannot be imported (circular); these
// stay package-private in ecs per stock-method-disposition.

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	connectivity "github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type NetworkType string

const (
	Classic         = NetworkType("Classic")
	Vpc             = NetworkType("Vpc")
	ClassicInternet = NetworkType("classic_internet")
	ClassicIntranet = NetworkType("classic_intranet")
	PUBLIC          = NetworkType("PUBLIC")
	PRIVATE         = NetworkType("PRIVATE")
)

type TimeType string

const (
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type PayType string

const (
	PrePaid    = PayType("PrePaid")
	PostPaid   = PayType("PostPaid")
	Prepaid    = PayType("Prepaid")
	Postpaid   = PayType("Postpaid")
	Serverless = PayType("Serverless")
)

type OptimizedType string

const (
	IOOptimized   = OptimizedType("optimized")
	NoneOptimized = OptimizedType("none")
)

type KubernetesNodeType string

const (
	KubernetesNodeMaster = ResourceType("Master")
	KubernetesNodeWorker = ResourceType("Worker")
)

const (
	DefaultTimeout       = 120
	DefaultTimeoutMedium = 500
	COLON_SEPARATED      = ":"
)

type FlagType string

const (
	OnFlag  = FlagType("on")
	OffFlag = FlagType("off")
)

const EcsInstance = "EcsInstance"

// descriptions mirrors root's provider descriptions lookup map. In root it is
// populated in init(); the only ECS reference (resource_alicloud_instance.go)
// looks up a key that is absent from root's map, so a nil map here yields the
// same empty-string result, preserving the original Description field value.
var descriptions map[string]string

func outputInstancesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"image_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"region_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"availability_zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vswitch_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"public_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"private_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func ecsTagsToMap(tags []ecs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !ecsTagIgnored(t) {
			result[t.TagKey] = t.TagValue
		}
	}

	return result
}

func ecsTagIgnored(t ecs.Tag) bool {
	filter := []string{"^aliyun", "^acs:", "^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.TagKey)
		ok, _ := regexp.MatchString(v, t.TagKey)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific tag %s (val: %s), ignoring.\n", t.TagKey, t.TagValue)
			return true
		}
	}
	return false
}

func setVolumeTags(client *connectivity.AliyunClient, resourceType TagResourceType, d *schema.ResourceData) error {
	ecsService := Service{client}
	if d.HasChange("volume_tags") {
		request := ecs.CreateDescribeDisksRequest()
		request.InstanceId = d.Id()
		var response *ecs.DescribeDisksResponse
		wait := incrementalWait(1*time.Second, 1*time.Second)
		err := retry.Retry(10*time.Minute, func() *retry.RetryError {
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.DescribeDisks(request)
			})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)

				}
				return retry.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
			response, _ = raw.(*ecs.DescribeDisksResponse)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		if len(response.Disks.Disk) == 0 {
			return WrapError(fmt.Errorf("the system disk cannot be queried in this instance %s. Please check whether you have permission to access the API DescribeDisks. Last response is: %v", d.Id(), response))
		}

		var ids []string
		systemDiskTag := make(map[string]interface{})
		for _, disk := range response.Disks.Disk {
			ids = append(ids, disk.DiskId)
			if disk.Type == "system" {
				for _, t := range disk.Tags.Tag {
					if !ecsService.ecsTagIgnored(t) {
						systemDiskTag[t.TagKey] = t.TagValue
					}
				}
			}
		}

		oraw, nraw := d.GetChange("volume_tags")
		if d.IsNewResource() {
			oraw = systemDiskTag
		}
		return updateTags(client, ids, resourceType, oraw, nraw)
	}

	return nil
}

func convertChargeType(originType string) string {
	if string("PostpayByFlow") == originType {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

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

func PrePaidDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	// payment_type is the instance_charge_type's replacement.
	// If both instance_charge_type and payment_type are "", it means hiding a default "PrePaid"
	if v, ok := d.GetOk("instance_charge_type"); ok && strings.ToLower(v.(string)) == "prepaid" {
		return true
	}
	if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
		return true
	}
	return false
}

func ecsInternetDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if max, ok := d.GetOk("internet_max_bandwidth_out"); ok && max.(int) > 0 {
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

func ecsSecurityGroupRulePreFixListIdDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	_, cidrIpExist := d.GetOk("cidr_ip")
	_, ipv6CidrIpExist := d.GetOk("ipv6_cidr_ip")
	_, SourceSecurityGroupIdExist := d.GetOk("source_security_group_id")
	if cidrIpExist || SourceSecurityGroupIdExist || ipv6CidrIpExist {
		return true
	}
	return false
}

func ecsSpotPriceLimitDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	value, ok := d.GetOk("instance_charge_type")
	if (!ok || value == "PostPaid") && d.Get("spot_strategy").(string) == "SpotWithPriceLimit" {
		return false
	}
	return true
}

func ecsSpotStrategyDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	value, ok := d.GetOk("instance_charge_type")
	if !ok || value == string(PostPaid) {
		return false
	}
	return true
}

func ecsSystemDiskPerformanceLevelSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	if d.Get("system_disk_category").(string) == string(DiskCloudESSD) {
		return false
	}
	return true
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
