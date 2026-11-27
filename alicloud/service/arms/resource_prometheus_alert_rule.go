package arms

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePrometheusAlertRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourcePrometheusAlertRuleCreate,
		Read:   ResourcePrometheusAlertRuleRead,
		Update: ResourcePrometheusAlertRuleUpdate,
		Delete: ResourcePrometheusAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"annotations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dispatch_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("notify_type"); ok && v.(string) == "DISPATCH_RULE" {
						return false
					}
					return true
				},
			},
			"duration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"notify_type": {
				Type:     schema.TypeString,
				Optional: true,
				// The server always materializes a notify type (default
				// ALERT_MANAGER), and UpdatePrometheusAlertRule is a
				// full-overwrite API: an update request that omits NotifyType
				// resets the server-side value to the default instead of
				// preserving it (verified against the live API, see the
				// unconfiguredNotifyType acceptance test). Marking the
				// attribute Computed keeps the state value authoritative when
				// the configuration leaves it unset, so Update carries it in
				// full and the plan stays converged; the same pattern is
				// already used for "type" above.
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_MANAGER", "DISPATCH_RULE"}, false),
			},
			"prometheus_alert_rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"prometheus_alert_rule_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// setAsJsonString serializes a TypeSet of {name, value} elements into the
// JSON string the ARMS alert rule APIs expect ("[]" for an empty set).
func setAsJsonString(d *schema.ResourceData, key string) (string, error) {
	items := make([]map[string]interface{}, 0)
	for _, item := range d.Get(key).(*schema.Set).List() {
		items = append(items, item.(map[string]interface{}))
	}
	return convert.ConvertArrayObjectToJsonString(items)
}

func ResourcePrometheusAlertRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePrometheusAlertRule"
	request := make(map[string]interface{})
	var err error
	if _, ok := d.GetOk("annotations"); ok {
		annotations, err := setAsJsonString(d, "annotations")
		if err != nil {
			return errs.WrapError(err)
		}
		request["Annotations"] = annotations
	}
	request["ClusterId"] = d.Get("cluster_id")
	if v, ok := d.GetOk("dispatch_rule_id"); ok {
		request["DispatchRuleId"] = v
	}
	request["Duration"] = d.Get("duration")
	request["Expression"] = d.Get("expression")
	if _, ok := d.GetOk("labels"); ok {
		labels, err := setAsJsonString(d, "labels")
		if err != nil {
			return errs.WrapError(err)
		}
		request["Labels"] = labels
	}
	request["Message"] = d.Get("message")
	if v, ok := d.GetOk("notify_type"); ok {
		request["NotifyType"] = v
	}
	request["AlertName"] = d.Get("prometheus_alert_rule_name")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_arms_prometheus_alert_rule", action, errs.AlibabaCloudSdkGoERROR)
	}
	responsePrometheusAlertRule := response["PrometheusAlertRule"].(map[string]interface{})
	d.SetId(fmt.Sprint(request["ClusterId"], ":", responsePrometheusAlertRule["AlertId"]))

	return ResourcePrometheusAlertRuleRead(d, meta)
}
func ResourcePrometheusAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	armsService := Service{client}
	object, err := armsService.DescribeArmsPrometheusAlertRule(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_arms_prometheus_alert_rule armsService.DescribeArmsPrometheusAlertRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("cluster_id", parts[0])
	// prometheus_alert_rule_id is declared as an int in the schema, so the
	// string alert id parsed from the resource id must be converted before Set,
	// otherwise a successful create/update panics in ResourceData.Set.
	prometheusAlertRuleId, err := convert.FormatInt(parts[1])
	if err != nil {
		return err
	}
	d.Set("prometheus_alert_rule_id", prometheusAlertRuleId)
	if v, ok := object["Annotations"].([]interface{}); ok {
		annotations := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			if item["Name"] == "message" {
				continue
			}

			temp := map[string]interface{}{
				"name":  item["Name"],
				"value": item["Value"],
			}

			annotations = append(annotations, temp)
		}
		if err := d.Set("annotations", annotations); err != nil {
			return errs.WrapError(err)
		}
	}
	dispatchRuleId, err := convert.FormatInt(object["DispatchRuleId"])
	if err != nil {
		return err
	}
	d.Set("dispatch_rule_id", fmt.Sprint(dispatchRuleId))
	d.Set("duration", object["Duration"])
	d.Set("expression", object["Expression"])
	if v, ok := object["Labels"].([]interface{}); ok {
		labels := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"name":  item["Name"],
				"value": item["Value"],
			}

			labels = append(labels, temp)
		}
		if err := d.Set("labels", labels); err != nil {
			return errs.WrapError(err)
		}
	}
	d.Set("message", object["Message"])
	d.Set("notify_type", object["NotifyType"])
	d.Set("prometheus_alert_rule_name", object["AlertName"])
	actualStatus, err := convert.FormatInt(object["Status"])
	if err != nil {
		return err
	}
	d.Set("status", actualStatus)
	d.Set("type", object["Type"])
	return nil
}
func ResourcePrometheusAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"ClusterId": parts[0],
		"AlertId":   parts[1],
	}

	request["Duration"] = d.Get("duration")
	if d.HasChange("duration") {
		update = true
	}
	request["Expression"] = d.Get("expression")
	if d.HasChange("expression") {
		update = true
	}
	request["Message"] = d.Get("message")
	if d.HasChange("message") {
		update = true
	}
	request["AlertName"] = d.Get("prometheus_alert_rule_name")
	if d.HasChange("prometheus_alert_rule_name") {
		update = true
	}
	request["RegionId"] = client.RegionId
	// UpdatePrometheusAlertRule is a full-overwrite API: fields omitted from
	// the request are cleared or reset on the server side (observed for
	// NotifyType against the live API, see the unconfiguredNotifyType
	// acceptance test), so labels, annotations, notify_type and type must
	// always be carried in full instead of only when they changed; otherwise
	// an update triggered by another field wipes them and the next plan
	// shows a permanent diff. DispatchRuleId is the exception the API
	// contract itself defines: the parameter is only meaningful for
	// notify_type = DISPATCH_RULE (the ListPrometheusAlertRules response
	// only reports it then), and Read stores the placeholder "0" for rules
	// without a dispatch rule, so it must not be sent otherwise.
	if d.HasChange("annotations") {
		update = true
	}
	annotations, err := setAsJsonString(d, "annotations")
	if err != nil {
		return errs.WrapError(err)
	}
	request["Annotations"] = annotations
	if d.HasChange("dispatch_rule_id") {
		update = true
	}
	if v, ok := d.GetOk("dispatch_rule_id"); ok && d.Get("notify_type").(string) == "DISPATCH_RULE" {
		request["DispatchRuleId"] = v
	}
	if d.HasChange("labels") {
		update = true
	}
	labels, err := setAsJsonString(d, "labels")
	if err != nil {
		return errs.WrapError(err)
	}
	request["Labels"] = labels
	if d.HasChange("notify_type") {
		update = true
	}
	// notify_type is Optional+Computed, so the state value (the server
	// default ALERT_MANAGER when the configuration leaves it unset) is
	// authoritative; it is always carried, never conditionally omitted.
	request["NotifyType"] = d.Get("notify_type")
	// type is ForceNew, so its value cannot change in place: carrying the
	// Computed state value (populated by Read from the server response)
	// keeps the server-side type intact across updates instead of relying
	// on the undocumented omit-preservation behavior.
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}

	if update {
		action := "UpdatePrometheusAlertRule"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return ResourcePrometheusAlertRuleRead(d, meta)
}
func ResourcePrometheusAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeletePrometheusAlertRule"
	var response map[string]interface{}
	request := map[string]interface{}{
		"AlertId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ARMS", "2019-08-08", action, nil, request, false)
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
	return nil
}
