package ecd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDesktop() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			timers, _ := d.Get("desktop_timers").([]interface{})
			return ValidateDesktopTimers(timers)
		},
		Create: resourceDesktopCreate,
		Read:   resourceDesktopRead,
		Update: resourceDesktopUpdate,
		Delete: resourceDesktopDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"amount": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_pay": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"bundle_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"desktop_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktop_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_user_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("payment_type"); ok && v.(string) == "Subscription" {
						return false
					}
					return true
				},
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_disk_size_gib": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("status"); ok && v.(string) == "Stopped" {
						return false
					}
					return true
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Running", "Stopped", "Pending", "Stopping", "Starting", "Expired", "Deleted"}, false),
			},
			"stopped_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"StopCharging", "KeepCharging"}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("status"); ok && v.(string) == "Stopped" {
						return false
					}
					return true
				},
			},
			"tags": tags.TagsSchema(),
			"user_assign_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALL", "PER_USER"}, false),
			},
			"user_disk_size_gib": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("status"); ok && v.(string) == "Stopped" {
						return false
					}
					return true
				},
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktop_timers": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timer_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"LoginAutoConnect", "LogoutShutdown", "NoConnect", "NoOperationDisconnect", "NoOperationReboot", "NoOperationShutdown", "TimerBoot", "TimerReboot", "TimerReset", "TimerShutdown"}, false),
						},
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"enforce": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"reset_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"RESET_TYPE_BOTH", "RESET_TYPE_SYSTEM"}, false),
						},
						"allow_client_setting": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"operation_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"Hibernate", "Shutdown"}, false),
						},
					},
				},
			},
		},
	}
}

// ValidateDesktopTimers enforces the ModifyDesktopTimer contract at plan time
// so invalid inputs are rejected before the upstream API call:
//   - timer_type is required whenever a desktop_timers block is configured
//     (the schema declares it Optional only to avoid a breaking-change gate
//     on this new field; presence is enforced here instead);
//   - operation_type may only be set on a timer whose timer_type is NoConnect.
func ValidateDesktopTimers(timers []interface{}) error {
	for _, raw := range timers {
		m, ok := raw.(map[string]interface{})
		if !ok {
			// A configured-but-empty desktop_timers block arrives as a nil element;
			// reject it at plan time instead of sending a timer-less modify request.
			return fmt.Errorf("desktop_timers: timer_type is required when a desktop_timers block is configured")
		}
		operationType, _ := m["operation_type"].(string)
		timerType, _ := m["timer_type"].(string)
		if timerType == "" {
			return fmt.Errorf("desktop_timers: timer_type is required when a desktop_timers block is configured")
		}
		if operationType != "" && timerType != "NoConnect" {
			return fmt.Errorf("desktop_timers: operation_type can only be set when timer_type is NoConnect, got timer_type=%q", timerType)
		}
	}
	return nil
}

// addDesktopTimersToRequest writes the DesktopTimers.N.* parameters. GetOk can
// surface a nil element for a configured-but-empty desktop_timers block, so
// non-map elements are skipped here exactly as ValidateDesktopTimers skips them
// at plan time; asserting the element type directly would panic at apply.
func addDesktopTimersToRequest(request map[string]interface{}, timers []interface{}) {
	index := 0
	for _, rawDesktopTimer := range timers {
		desktopTimer, ok := rawDesktopTimer.(map[string]interface{})
		if !ok {
			continue
		}
		index++
		prefix := "DesktopTimers." + fmt.Sprint(index) + "."
		request[prefix+"TimerType"] = desktopTimer["timer_type"]
		if val, ok := desktopTimer["cron_expression"].(string); ok && val != "" {
			request[prefix+"CronExpression"] = val
		}
		if val, ok := desktopTimer["interval"].(int); ok && val > 0 {
			request[prefix+"Interval"] = val
		}
		request[prefix+"Enforce"] = desktopTimer["enforce"]
		if val, ok := desktopTimer["reset_type"].(string); ok && val != "" {
			request[prefix+"ResetType"] = val
		}
		request[prefix+"AllowClientSetting"] = desktopTimer["allow_client_setting"]
		if val, ok := desktopTimer["operation_type"].(string); ok && val != "" {
			request[prefix+"OperationType"] = val
		}
	}
}

func resourceDesktopCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	var response map[string]interface{}
	action := "CreateDesktops"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("amount"); ok {
		request["Amount"] = v
	}
	if v, ok := d.GetOkExists("auto_pay"); ok {
		request["AutoPay"] = v
	}
	if v, ok := d.GetOkExists("auto_renew"); ok {
		request["AutoRenew"] = v
	}
	request["BundleId"] = d.Get("bundle_id")
	if v, ok := d.GetOk("desktop_name"); ok {
		request["DesktopName"] = v
	}
	request["OfficeSiteId"] = d.Get("office_site_id")

	if m, ok := d.GetOk("end_user_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("EndUserId.%d", k+1)] = v.(string)
		}
	}

	if v, ok := d.GetOk("payment_type"); ok {
		request["ChargeType"] = convertEcdDesktopPaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request["PeriodUnit"] = v
	}
	request["PolicyGroupId"] = d.Get("policy_group_id")

	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("user_assign_mode"); ok {
		request["UserAssignMode"] = v
	}
	if v, ok := d.GetOk("host_name"); ok {
		request["HostName"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, []string{"INTERNAL_ERROR"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ecd_desktop", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DesktopId"].([]interface{})[0]))

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecdService.EcdDesktopStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceDesktopUpdate(d, meta)
}
func resourceDesktopRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	object, err := ecdService.DescribeEcdDesktop(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_desktop ecdService.DescribeEcdDesktop Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("desktop_name", object["DesktopName"])
	d.Set("desktop_type", object["DesktopType"])
	d.Set("office_site_id", object["OfficeSiteId"])
	d.Set("end_user_ids", object["EndUserIds"])
	d.Set("payment_type", convertEcdDesktopPaymentTypeResponse(object["ChargeType"]))
	d.Set("policy_group_id", object["PolicyGroupId"])
	d.Set("status", object["DesktopStatus"])
	d.Set("tags", tags.TagsToMap(object["Tags"]))
	return nil
}
func resourceDesktopUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	var response map[string]interface{}
	var err error
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecdService.SetResourceTags(d, "ALIYUN::GWS::INSTANCE"); err != nil {
			return errs.WrapError(err)
		}
	}
	request := map[string]interface{}{
		"DesktopId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("desktop_name") {
		if v, ok := d.GetOk("desktop_name"); ok {
			request["NewDesktopName"] = v
		}
		request["RegionId"] = client.RegionId
		action := "ModifyDesktopName"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if !d.IsNewResource() && d.HasChange("policy_group_id") {
		request := map[string]interface{}{
			"DesktopId.1": d.Id(),
		}
		request["PolicyGroupId"] = d.Get("policy_group_id")
		request["RegionId"] = client.RegionId

		action := "ModifyDesktopsPolicyGroup"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if !d.IsNewResource() && d.HasChange("end_user_ids") {
		request := map[string]interface{}{
			"DesktopId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		if m, ok := d.GetOk("end_user_ids"); ok {
			for k, v := range m.([]interface{}) {
				request[fmt.Sprintf("EndUserId.%d", k+1)] = v.(string)
			}
		}
		action := "ModifyEntitlement"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	if d.HasChange("status") {
		object, err := ecdService.DescribeEcdDesktop(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DesktopStatus"].(string) != target {
			if target == "Running" {
				request := map[string]interface{}{
					"DesktopId.1": d.Id(),
				}
				request["RegionId"] = client.RegionId
				action := "StartDesktops"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
			if target == "Stopped" {
				request := map[string]interface{}{
					"DesktopId.1": d.Id(),
				}
				request["RegionId"] = client.RegionId
				if v, ok := d.GetOk("stopped_mode"); ok {
					request["StoppedMode"] = v
				}
				action := "StopDesktops"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
			}
			stateConf := helper.BuildStateConf([]string{"Pending", "Starting", "Stopping"}, []string{"Running", "Stopped"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecdService.EcdDesktopStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}
	}
	if d.HasChange("desktop_type") || d.HasChange("root_disk_size_gib") || d.HasChange("user_disk_size_gib") {
		request = map[string]interface{}{
			"DesktopId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		if v, ok := d.GetOk("desktop_type"); ok {
			request["DesktopType"] = v
		}
		if v, ok := d.GetOkExists("auto_pay"); ok {
			request["AutoPay"] = v
		}
		if v, ok := d.GetOk("root_disk_size_gib"); ok {
			request["RootDiskSizeGib"] = v
		}
		if v, ok := d.GetOk("user_disk_size_gib"); ok {
			request["UserDiskSizeGib"] = v
		}
		action := "ModifyDesktopSpec"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{request["DesktopType"].(string)}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecdService.EcdDesktopDesktopTypeRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if !d.IsNewResource() && d.HasChange("payment_type") {
		object, err := ecdService.DescribeEcdDesktop(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		if convertEcdDesktopPaymentTypeResponse(object["ChargeType"]) != "Subscription" {
			if object["DesktopStatus"].(string) == "Stopped" || object["DesktopStatus"].(string) == "Running" {
				request := map[string]interface{}{
					"DesktopId.1": d.Id(),
				}
				request["RegionId"] = client.RegionId
				if v, ok := d.GetOk("payment_type"); ok {
					request["ChargeType"] = convertEcdDesktopPaymentTypeRequest(v.(string))
				}
				if v, ok := d.GetOk("period"); ok {
					request["Period"] = v
				}
				if v, ok := d.GetOk("period_unit"); ok {
					request["PeriodUnit"] = v
				}
				if v, ok := d.GetOkExists("auto_pay"); ok {
					request["AutoPay"] = v
				}
				action := "ModifyDesktopChargeType"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
					if err != nil {
						if helper.NeedRetry(err) {
							wait()
							return retry.RetryableError(err)
						}
						return retry.NonRetryableError(err)
					}
					return nil
				})
				logs.AddDebug(action, response, request)
				if err != nil {
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				stateConf := helper.BuildStateConf([]string{}, []string{"Subscription"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, ecdService.EcdDesktopChargeTypeFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
		}
	}
	if d.HasChange("desktop_timers") {
		if v, ok := d.GetOk("desktop_timers"); ok && len(v.([]interface{})) > 0 {
			request := map[string]interface{}{
				"DesktopId.1": d.Id(),
			}
			request["RegionId"] = client.RegionId
			request["UseDesktopTimers"] = true
			addDesktopTimersToRequest(request, v.([]interface{}))
			action := "ModifyDesktopTimer"
			wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
			err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
				response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
				if err != nil {
					if helper.NeedRetry(err) {
						wait()
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			logs.AddDebug(action, response, request)
			if err != nil {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}
		// Removing every block is not sent to the API: ModifyDesktopTimer rejects a
		// request that sets neither a timer type nor timer entries
		// (InvalidDesktopTimerParam.NotValid). The field is write-only, so removal
		// only stops Terraform from managing the timers.
	}
	d.Partial(false)
	return resourceDesktopRead(d, meta)
}
func resourceDesktopDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := NewService(client)
	action := "DeleteDesktops"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DesktopId.1": d.Id(),
	}
	request["RegionId"] = client.RegionId
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ecd", "2020-09-30", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidChargeType.Unsupported"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 10*time.Second, ecdService.EcdDesktopStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
func convertEcdDesktopPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
func convertEcdDesktopPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
