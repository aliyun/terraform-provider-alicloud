package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudEcsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsDedicatedHostCreate,
		Read:   resourceAlicloudEcsDedicatedHostRead,
		Update: resourceAlicloudEcsDedicatedHostUpdate,
		Delete: resourceAlicloudEcsDedicatedHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_on_maintenance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Stop", "Migrate"}, false),
				Default:      "Stop",
			},
			"auto_placement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Default:      "on",
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_renew_period": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_renew_with_ecs": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewWithEcs", "StopRenewWithEcs", "NoOperation"}, false),
			},
			"cpu_over_commit_ratio": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"dedicated_host_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"detail_fee": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 12, 24, 36, 48, 60}),
			},
			"expired_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"min_quantity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// lintignore: S018
			"network_attributes": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"slb_udp_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Week", "Month", "Year"}, false),
			},
			"quantity": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ForceNew: true,
				// A Terraform resource maps to a single dedicated host.
				// AllocateDedicatedHosts with quantity > 1 creates multiple
				// hosts but only the first DedicatedHostId is tracked, leaving
				// the rest as unmanaged paid resources. Restrict to 1 until a
				// proper multi-host lifecycle design exists.
				ValidateFunc: validation.IntInSlice([]int{1}),
			},
			"renewal_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"AutoRenewal", "Normal", "NotRenewal"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sale_cycle": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "AllocateDedicatedHosts"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("action_on_maintenance"); ok {
		request["ActionOnMaintenance"] = v
	}

	if v, ok := d.GetOk("auto_placement"); ok {
		request["AutoPlacement"] = v
	}

	if v, ok := d.GetOk("auto_release_time"); ok {
		request["AutoReleaseTime"] = v
	}

	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}

	if v, ok := d.GetOk("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}

	if v, ok := d.GetOk("cpu_over_commit_ratio"); ok {
		request["CpuOverCommitRatio"] = v
	}

	if v, ok := d.GetOk("dedicated_host_cluster_id"); ok {
		request["DedicatedHostClusterId"] = v
	}

	if v, ok := d.GetOk("dedicated_host_name"); ok {
		request["DedicatedHostName"] = v
	}

	request["DedicatedHostType"] = d.Get("dedicated_host_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	} else if v, ok := d.GetOk("expired_time"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("min_quantity"); ok {
		request["MinQuantity"] = v
	}

	if v, ok := d.GetOk("network_attributes"); ok {
		networkAttributesMap := make(map[string]interface{})
		for _, networkAttributes := range v.(*schema.Set).List() {
			networkAttributesArg := networkAttributes.(map[string]interface{})
			networkAttributesMap["SlbUdpTimeout"] = requests.NewInteger(networkAttributesArg["slb_udp_timeout"].(int))
			networkAttributesMap["UdpTimeout"] = requests.NewInteger(networkAttributesArg["udp_timeout"].(int))
		}
		request["NetworkAttributes"] = networkAttributesMap

	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = v
	}

	// Quantity is restricted to 1 by schema validation: a Terraform resource
	// maps to a single dedicated host. AllocateDedicatedHosts with quantity > 1
	// creates multiple hosts but only the first DedicatedHostId is tracked.
	request["Quantity"] = d.Get("quantity")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	} else if v, ok := d.GetOk("sale_cycle"); ok {
		request["PeriodUnit"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_dedicated_host", action, AlibabaCloudSdkGoERROR)
	}
	responseDedicatedHostIdSets := response["DedicatedHostIdSets"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseDedicatedHostIdSets["DedicatedHostId"].([]interface{})[0]))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 15*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{"PermanentFailure"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	// Apply the auto-renewal attribute after the dedicated host becomes
	// Available. RenewalStatus/AutoRenewWithEcs/Duration are accepted by the
	// ModifyDedicatedHostAutoRenewAttribute API, which only applies to
	// subscription (PrePaid) dedicated hosts.
	autoRenewWithEcs := d.Get("auto_renew_with_ecs").(string)
	if d.Get("payment_type").(string) == "PrePaid" &&
		(d.Get("renewal_status").(string) != "" ||
			(autoRenewWithEcs != "" && autoRenewWithEcs != "NoOperation") ||
			d.Get("duration").(int) != 0) {
		if err := modifyDedicatedHostAutoRenewAttribute(d, meta); err != nil {
			return WrapError(err)
		}
	}

	return resourceAlicloudEcsDedicatedHostRead(d, meta)
}
func resourceAlicloudEcsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDedicatedHost(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_dedicated_host ecsService.DescribeEcsDedicatedHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("action_on_maintenance", object["ActionOnMaintenance"])
	d.Set("auto_placement", object["AutoPlacement"])
	d.Set("auto_release_time", object["AutoReleaseTime"])
	d.Set("cpu_over_commit_ratio", object["CpuOverCommitRatio"])
	d.Set("dedicated_host_name", object["DedicatedHostName"])
	d.Set("dedicated_host_type", object["DedicatedHostType"])
	d.Set("description", object["Description"])

	networkAttributesSli := make([]map[string]interface{}, 0)
	if len(object["NetworkAttributes"].(map[string]interface{})) > 0 {
		networkAttributes := object["NetworkAttributes"]
		networkAttributesMap := make(map[string]interface{})
		networkAttributesMap["slb_udp_timeout"] = networkAttributes.(map[string]interface{})["SlbUdpTimeout"]
		networkAttributesMap["udp_timeout"] = networkAttributes.(map[string]interface{})["UdpTimeout"]
		networkAttributesSli = append(networkAttributesSli, networkAttributesMap)
	}
	d.Set("network_attributes", networkAttributesSli)
	d.Set("payment_type", object["ChargeType"])
	// quantity is a create-only parameter not returned by DescribeDedicatedHosts;
	// unconditionally preserve the schema default (1) in state so that import
	// state comparison is stable and the attribute no longer needs to be in
	// ImportStateVerifyIgnore.
	d.Set("quantity", 1)
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("period_unit", object["SaleCycle"])
	d.Set("sale_cycle", object["SaleCycle"])
	d.Set("status", object["Status"])
	// The auto-renewal attributes (RenewalStatus, AutoRenewWithEcs, Duration)
	// are returned by the separate DescribeDedicatedHostAutoRenew API rather
	// than DescribeDedicatedHosts. They only apply to subscription (PrePaid)
	// dedicated hosts; for PostPaid hosts the values are cleared.
	if object["ChargeType"] == "PrePaid" {
		renewAttr, renewErr := ecsService.DescribeDedicatedHostAutoRenewAttribute(d.Id())
		if renewErr == nil && renewAttr != nil {
			d.Set("renewal_status", renewAttr["RenewalStatus"])
			d.Set("auto_renew_with_ecs", renewAttr["AutoRenewWithEcs"])
			d.Set("duration", renewAttr["Duration"])
		} else {
			log.Printf("[DEBUG] Resource alicloud_ecs_dedicated_host DescribeDedicatedHostAutoRenewAttribute failed: %s", renewErr)
		}
	} else {
		d.Set("renewal_status", nil)
		d.Set("auto_renew_with_ecs", nil)
		d.Set("duration", nil)
	}
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("zone_id", object["ZoneId"])
	return nil
}
func resourceAlicloudEcsDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var err error
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "ddh"); err != nil {
			return WrapError(err)
		}
	}
	if !d.IsNewResource() && d.HasChange("auto_release_time") {
		request := map[string]interface{}{
			"DedicatedHostId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		request["AutoReleaseTime"] = d.Get("auto_release_time")
		action := "ModifyDedicatedHostAutoReleaseTime"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := map[string]interface{}{
			"ResourceId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		request["ResourceGroupId"] = d.Get("resource_group_id")
		request["ResourceType"] = "ddh"
		action := "JoinResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update := false
	request := map[string]interface{}{
		"DedicatedHostIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}
	if !d.IsNewResource() && (d.HasChange("period") || d.HasChange("expired_time")) {
		update = true
		if d.HasChange("period") {
			request["Period"] = d.Get("period")
		} else {
			request["Period"] = d.Get("expired_time")
		}
	}
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && (d.HasChange("period_unit") || d.HasChange("sale_cycle")) {
		update = true
		if d.HasChange("period_unit") {
			request["PeriodUnit"] = d.Get("period_unit")
		} else {
			request["PeriodUnit"] = d.Get("sale_cycle")
		}
	}
	if update {
		action := "RenewDedicatedHosts"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	modifyDedicatedHostsChargeTypeReq := map[string]interface{}{
		"DedicatedHostIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}
	modifyDedicatedHostsChargeTypeReq["RegionId"] = client.RegionId
	// AutoPay is surfaced as the auto_pay argument and defaults to true:
	// when changing the charge type of a subscription dedicated host, the
	// order is paid automatically by default; set auto_pay to false to
	// generate an unpaid order that must be settled manually.
	if v, ok := d.GetOkExists("auto_pay"); ok {
		modifyDedicatedHostsChargeTypeReq["AutoPay"] = v
	} else {
		modifyDedicatedHostsChargeTypeReq["AutoPay"] = true
	}
	if d.HasChange("period") {
		modifyDedicatedHostsChargeTypeReq["Period"] = d.Get("period")
	} else if d.HasChange("expired_time") {
		modifyDedicatedHostsChargeTypeReq["Period"] = d.Get("expired_time")
	} else if v, ok := d.GetOk("period"); ok {
		modifyDedicatedHostsChargeTypeReq["Period"] = v
	} else {
		modifyDedicatedHostsChargeTypeReq["Period"] = d.Get("expired_time")
	}
	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
		modifyDedicatedHostsChargeTypeReq["DedicatedHostChargeType"] = d.Get("payment_type")
	}
	if d.HasChange("period_unit") {
		modifyDedicatedHostsChargeTypeReq["PeriodUnit"] = d.Get("period_unit")
	} else if d.HasChange("sale_cycle") {
		modifyDedicatedHostsChargeTypeReq["PeriodUnit"] = d.Get("sale_cycle")
	} else if v, ok := d.GetOk("period_unit"); ok {
		modifyDedicatedHostsChargeTypeReq["PeriodUnit"] = v
	} else {
		modifyDedicatedHostsChargeTypeReq["PeriodUnit"] = d.Get("sale_cycle")
	}
	if update {
		// ModifyDedicatedHostsChargeType only accepts PeriodUnit=Week or Month.
		// Year is valid for AllocateDedicatedHosts, RenewDedicatedHosts and
		// ModifyDedicatedHostAutoRenewAttribute but not for charge-type changes,
		// where the API rejects it with InvalidPeriod.UnitMismatch.
		if modifyDedicatedHostsChargeTypeReq["PeriodUnit"] == "Year" {
			return WrapError(fmt.Errorf("period_unit cannot be Year when changing payment_type: ModifyDedicatedHostsChargeType only accepts Week and Month"))
		}
		if _, ok := d.GetOkExists("detail_fee"); ok {
			modifyDedicatedHostsChargeTypeReq["DetailFee"] = d.Get("detail_fee")
		}
		if _, ok := d.GetOkExists("dry_run"); ok {
			modifyDedicatedHostsChargeTypeReq["DryRun"] = d.Get("dry_run")
		}
		action := "ModifyDedicatedHostsChargeType"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, modifyDedicatedHostsChargeTypeReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, modifyDedicatedHostsChargeTypeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	update = false
	modifyDedicatedHostAttributeReq := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}
	modifyDedicatedHostAttributeReq["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("action_on_maintenance") {
		update = true
		modifyDedicatedHostAttributeReq["ActionOnMaintenance"] = d.Get("action_on_maintenance")
	}
	if !d.IsNewResource() && d.HasChange("auto_placement") {
		update = true
		modifyDedicatedHostAttributeReq["AutoPlacement"] = d.Get("auto_placement")
	}
	if !d.IsNewResource() && d.HasChange("cpu_over_commit_ratio") {
		update = true
		modifyDedicatedHostAttributeReq["CpuOverCommitRatio"] = d.Get("cpu_over_commit_ratio")
	}
	if !d.IsNewResource() && d.HasChange("dedicated_host_name") {
		update = true
		modifyDedicatedHostAttributeReq["DedicatedHostName"] = d.Get("dedicated_host_name")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		modifyDedicatedHostAttributeReq["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("network_attributes") {
		update = true
		if d.Get("network_attributes") != nil {
			networkAttributesMap := make(map[string]interface{})
			for _, networkAttributes := range d.Get("network_attributes").(*schema.Set).List() {
				networkAttributesArg := networkAttributes.(map[string]interface{})
				networkAttributesMap["SlbUdpTimeout"] = requests.NewInteger(networkAttributesArg["slb_udp_timeout"].(int))
				networkAttributesMap["UdpTimeout"] = requests.NewInteger(networkAttributesArg["udp_timeout"].(int))
			}
			modifyDedicatedHostAttributeReq["NetworkAttributes"] = networkAttributesMap
		}
	}
	if update {
		if _, ok := d.GetOk("dedicated_host_cluster_id"); ok {
			modifyDedicatedHostAttributeReq["DedicatedHostClusterId"] = d.Get("dedicated_host_cluster_id")
		}
		action := "ModifyDedicatedHostAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, modifyDedicatedHostAttributeReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, modifyDedicatedHostAttributeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	// ModifyDedicatedHostAutoRenewAttribute: renewal_status, auto_renew_with_ecs
	// and duration are managed by the auto-renewal API and only apply to
	// subscription (PrePaid) dedicated hosts.
	if !d.IsNewResource() && d.Get("payment_type").(string) == "PrePaid" &&
		(d.HasChange("renewal_status") || d.HasChange("auto_renew_with_ecs") || d.HasChange("duration")) {
		if err := modifyDedicatedHostAutoRenewAttribute(d, meta); err != nil {
			return WrapError(err)
		}
	}
	d.Partial(false)
	return resourceAlicloudEcsDedicatedHostRead(d, meta)
}
func resourceAlicloudEcsDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "PrePaid" {
		log.Printf("[WARN] Cannot destroy Subscription resource: alicloud_ecs_dedicated_host. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	action := "ReleaseDedicatedHost"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectHostStatus.Initializing"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDedicatedHostId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}

// modifyDedicatedHostAutoRenewAttribute calls the
// ModifyDedicatedHostAutoRenewAttribute API to set the renewal status, the
// auto-renew-with-ecs behaviour and (when auto-renewal is enabled) the renewal
// duration of a subscription dedicated host. It is shared by the Create
// (post-create configuration) and Update paths so the wiring stays consistent.
func modifyDedicatedHostAutoRenewAttribute(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	var err error
	action := "ModifyDedicatedHostAutoRenewAttribute"
	request := map[string]interface{}{
		"DedicatedHostIds": convertListToJsonString(convertListStringToListInterface([]string{d.Id()})),
	}
	request["RegionId"] = client.RegionId
	renewalStatus := d.Get("renewal_status").(string)
	request["RenewalStatus"] = renewalStatus
	if v, ok := d.GetOkExists("auto_renew_with_ecs"); ok {
		request["AutoRenewWithEcs"] = v
	} else {
		request["AutoRenewWithEcs"] = "NoOperation"
	}
	// Duration and PeriodUnit are only honoured when auto-renewal is enabled.
	if renewalStatus == "AutoRenewal" {
		if v, ok := d.GetOk("duration"); ok {
			request["Duration"] = requests.NewInteger(v.(int))
		}
		periodUnit := d.Get("period_unit").(string)
		if d.HasChange("sale_cycle") && !d.HasChange("period_unit") {
			// When the deprecated sale_cycle alias is explicitly changed but
			// period_unit is not configured, honour the new sale_cycle value
			// instead of the stale period_unit populated by Read.
			periodUnit = d.Get("sale_cycle").(string)
		}
		if periodUnit == "" {
			// Fall back to the deprecated sale_cycle alias for backward
			// compatibility when period_unit is not configured.
			periodUnit = d.Get("sale_cycle").(string)
		}
		if periodUnit == "" {
			periodUnit = "Month"
		}
		request["PeriodUnit"] = periodUnit
		// AutoRenewWithEcs=AutoRenewWithEcs requires AutoRenew=true to take effect.
		request["AutoRenew"] = true
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
		if err != nil {
			// The charge type may not have fully settled on the billing subsystem
			// right after a charge-type change, so the auto-renew call can
			// transiently fail with ChargeTypeViolation.
			if NeedRetry(err) || IsExpectedErrors(err, []string{"ChargeTypeViolation", "IncorrectDedicatedHostStatus"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
