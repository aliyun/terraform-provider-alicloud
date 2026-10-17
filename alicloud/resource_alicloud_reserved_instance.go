package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// autoRenewFromRawConfig checks whether the user explicitly set `auto_renew` in
// their HCL configuration, as opposed to a value populated by Read (computed).
//
// GetOkExists("auto_renew") cannot make this distinction: because auto_renew is
// Optional+Computed and Read always calls d.Set("auto_renew", …), the merged
// writer carries the Read-derived value, causing GetOkExists to return
// exists=true even when the user never wrote auto_renew in their config.  This
// leads to a control-flow inversion in Update: the computed auto_renew value
// overrides the user's explicit renewal_status, making it impossible to switch
// auto-renewal on or off via the legacy field alone.
//
// GetRawConfig() returns only what the user wrote; attributes not present in the
// HCL are null.  This lets us correctly distinguish "user set auto_renew=false"
// from "user did not set auto_renew at all".
//
// During plan, computed-but-unknown values are represented as
// cty.UnknownVal(bool).  Such a value is NOT an explicit user choice and must
// be treated the same as absent: IsKnown() returns false, so we return
// (false, false).  Without this check, UnknownVal(bool) would fall through to
// RawEquals(cty.True) (which is false) and incorrectly return (false, true),
// making the provider believe the user explicitly set auto_renew=false.
//
// Returns (value, true) when auto_renew is explicit in config (true or false),
// (false, false) when it is absent or unknown (computed).
func autoRenewFromRawConfig(rawConfig cty.Value) (bool, bool) {
	if rawConfig.IsNull() {
		return false, false
	}
	attr := rawConfig.GetAttr("auto_renew")
	if !attr.IsKnown() || attr.IsNull() {
		return false, false
	}
	return attr.RawEquals(cty.True), true
}

// autoRenewFromConfig is the production wrapper around autoRenewFromRawConfig.
func autoRenewFromConfig(d *schema.ResourceData) (bool, bool) {
	return autoRenewFromRawConfig(d.GetRawConfig())
}

// renewalStatusFromConfig returns the explicitly-configured renewal_status and
// whether the user set it in HCL, mirroring autoRenewFromConfig's distinction
// between an explicit value and one populated by Read (computed).
//
// In Update, d.GetOk("renewal_status") returns exists=true for the computed
// state value (Read always calls d.Set("renewal_status", …)), so a period-only
// update — where the user set neither auto_renew nor renewal_status — would
// otherwise re-send the stale state RenewalStatus and overwrite a server-side
// renewal-status change. ModifyReservedInstanceAutoRenewAttribute treats
// RenewalStatus as optional (omitting it preserves the cloud's current renewal
// status), so a period-only update must omit it rather than re-assert state.
// https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-modifyreservedinstanceautorenewattribute
func renewalStatusFromConfig(d *schema.ResourceData) (string, bool) {
	rawConfig := d.GetRawConfig()
	if rawConfig.IsNull() {
		return "", false
	}
	attr := rawConfig.GetAttr("renewal_status")
	if !attr.IsKnown() || attr.IsNull() {
		return "", false
	}
	return attr.AsString(), true
}

// renewalStatusDiffSuppress is the pure logic of the renewal_status
// DiffSuppressFunc, extracted so it can be unit-tested with cty values
// without relying on the SDK's GetRawConfig (which is null in test
// ResourceData built via Data(state, diff)).
//
// When auto_renew is explicitly set in config it takes precedence over
// renewal_status.  Suppress the renewal_status diff so that the computed
// value (set by Read from the auto_renew-derived RenewalStatus) does not
// produce a perpetual diff against the user's renewal_status value.
func renewalStatusDiffSuppress(rawConfig cty.Value) bool {
	if _, explicit := autoRenewFromRawConfig(rawConfig); explicit {
		return true
	}
	return false
}

// autoRenewPeriodDiffSuppress is the pure logic of the auto_renew_period
// DiffSuppressFunc, extracted for the same testability reason as above.
//
// The period is meaningful only when auto-renewal is active.  When auto_renew
// is explicit, its value determines whether the period diff is suppressed.
// When auto_renew is not in config, the effective renewal_status (from config
// or state) determines suppression.
func autoRenewPeriodDiffSuppress(rawConfig cty.Value, renewalStatus string) bool {
	if v, explicit := autoRenewFromRawConfig(rawConfig); explicit {
		return !v
	}
	if renewalStatus == string(RenewAutoRenewal) {
		return false
	}
	return true
}

func resourceAliCloudReservedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudReservedInstanceCreate,
		Read:   resourceAliCloudReservedInstanceRead,
		Update: resourceAliCloudReservedInstanceUpdate,
		Delete: resourceAliCloudReservedInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allocation_status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"renewal_status": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{string(RenewAutoRenewal), string(RenewNormal)}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return renewalStatusDiffSuppress(d.GetRawConfig())
				},
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"auto_renew_period": {
				Computed:     true,
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 12, 36, 60}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return autoRenewPeriodDiffSuppress(d.GetRawConfig(), d.Get("renewal_status").(string))
				},
			},
			"create_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"expired_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"operation_locks": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lock_reason": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"reserved_instance_name": {
				Optional:      true,
				Computed:      true,
				Type:          schema.TypeString,
				ConflictsWith: []string{"name"},
			},
			"start_time": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"status": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"tags": tagsSchema(),
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Region", "Zone"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux"}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Month", "Year"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 5}),
			},
			"offering_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"No Upfront", "Partial Upfront", "All Upfront"}, false),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from provider version 1.194.0. New field 'reserved_instance_name' instead.",
				ConflictsWith: []string{"reserved_instance_name"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}
func resourceAliCloudReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}
	var err error

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("instance_amount"); ok {
		request["InstanceAmount"] = v
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request["InstanceType"] = v
	}
	if v, ok := d.GetOk("offering_type"); ok {
		request["OfferingType"] = v
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	if v, ok := d.GetOk("platform"); ok {
		request["Platform"] = v
	}
	if v, ok := d.GetOk("reserved_instance_name"); ok {
		request["ReservedInstanceName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["ReservedInstanceName"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if scope, ok := d.GetOk("scope"); ok {
		request["Scope"] = scope
		if v, ok := d.GetOk("zone_id"); ok {
			if scope == "Zone" && v == "" {
				return WrapError(Error("Required when Scope is Zone."))
			}
			request["ZoneId"] = v
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}

	if v, explicit := autoRenewFromConfig(d); explicit {
		request["AutoRenew"] = v
	} else if v, ok := d.GetOk("renewal_status"); ok {
		request["AutoRenew"] = v.(string) == string(RenewAutoRenewal)
	}
	if v, ok := d.GetOkExists("auto_renew_period"); ok {
		request["AutoRenewPeriod"] = v
	}

	request["ClientToken"] = buildClientToken("PurchaseReservedInstancesOffering")
	var response map[string]interface{}
	action := "PurchaseReservedInstancesOffering"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response = resp

		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_reserved_instance", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.ReservedInstanceIdSets.ReservedInstanceId[0]", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_ecs_reserved_instance")
	} else {
		d.SetId(fmt.Sprint(v))
	}
	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecsService.EcsReservedInstanceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudReservedInstanceRead(d, meta)
}
func resourceAliCloudReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	ecsService := EcsService{client}
	d.Partial(true)
	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "reservedinstance"); err != nil {
			return WrapError(err)
		}
	}

	update := false
	request := ecs.CreateModifyReservedInstanceAttributeRequest()
	request.ReservedInstanceId = d.Id()
	request.RegionId = client.RegionId
	if d.HasChange("name") {
		update = true
		if v, ok := d.GetOk("name"); ok {
			request.ReservedInstanceName = v.(string)
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request.Description = v.(string)
		}
	}
	if d.HasChange("reserved_instance_name") {
		update = true
		if v, ok := d.GetOk("reserved_instance_name"); ok {
			request.ReservedInstanceName = v.(string)
		}
	}

	if update {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyReservedInstanceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	}

	if d.HasChanges("auto_renew", "auto_renew_period", "renewal_status") {
		request := map[string]interface{}{
			"ReservedInstanceId": []string{d.Id()},
			"RegionId":           client.RegionId,
		}
		if v, ok := d.GetOk("auto_renew_period"); ok {
			if formatInt(v) == 1 {
				request["Period"] = v
				request["PeriodUnit"] = "Month"
			} else {
				request["Period"] = formatInt(v) / 12
				request["PeriodUnit"] = "Year"
			}
		}
		if v, explicit := autoRenewFromConfig(d); explicit {
			if v {
				request["RenewalStatus"] = string(RenewAutoRenewal)
			} else {
				request["RenewalStatus"] = string(RenewNormal)
			}
		} else if v, ok := renewalStatusFromConfig(d); ok {
			request["RenewalStatus"] = v
		}

		action := "ModifyReservedInstanceAutoRenewAttribute"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			resp, err := client.RpcPost("Ecs", "2014-05-26", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	stateConf := BuildStateConf([]string{}, []string{"Active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsReservedInstanceStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	d.Partial(false)

	return resourceAliCloudReservedInstanceRead(d, meta)
}
func resourceAliCloudReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	reservedInstances, err := ecsService.DescribeReservedInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", reservedInstances.InstanceType)
	d.Set("allocation_status", reservedInstances.AllocationStatus)
	d.Set("scope", reservedInstances.Scope)
	d.Set("zone_id", reservedInstances.ZoneId)
	d.Set("instance_amount", reservedInstances.InstanceAmount)
	d.Set("platform", strings.Title(reservedInstances.Platform))
	d.Set("offering_type", reservedInstances.OfferingType)
	d.Set("name", reservedInstances.ReservedInstanceName)
	d.Set("reserved_instance_name", reservedInstances.ReservedInstanceName)
	d.Set("description", reservedInstances.Description)
	d.Set("resource_group_id", reservedInstances.ResourceGroupId)
	d.Set("status", reservedInstances.Status)
	d.Set("create_time", reservedInstances.CreationTime)
	d.Set("expired_time", reservedInstances.ExpiredTime)
	d.Set("start_time", reservedInstances.StartTime)

	operationLocks := make([]map[string]interface{}, 0)
	for _, lock := range reservedInstances.OperationLocks.OperationLock {
		operationLocks = append(operationLocks, map[string]interface{}{
			"lock_reason": lock.LockReason,
		})
	}
	d.Set("operation_locks", operationLocks)
	tags := make(map[string]string)
	for _, t := range reservedInstances.Tags.Tag {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)

	object, err := ecsService.DescribeReservedInstanceAutoRenewAttribute(d.Id())
	if err != nil {
		return WrapError(err)
	}
	renewalStatus, _ := object["RenewalStatus"].(string)
	d.Set("renewal_status", renewalStatus)
	d.Set("auto_renew", renewalStatus == string(RenewAutoRenewal))

	if v, ok := object["Duration"]; ok && formatInt(v) != 0 {
		renewPeriod := formatInt(v)
		if object["PeriodUnit"] == string(Year) {
			renewPeriod = renewPeriod * 12
		}
		d.Set("auto_renew_period", renewPeriod)
	}

	return nil
}
func resourceAliCloudReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// PurchaseReservedInstancesOffering can not be release.
	return nil
}
