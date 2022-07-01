package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAlarmCreate,
		Read:   resourceAlicloudCmsAlarmRead,
		Update: resourceAlicloudCmsAlarmUpdate,
		Delete: resourceAlicloudCmsAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:          schema.TypeMap,
				Optional:      true,
				Computed:      true,
				Elem:          schema.TypeString,
				ConflictsWith: []string{"metric_dimensions"},
				Deprecated:    "Field 'dimensions' has been deprecated from version 1.173.0. Use 'metric_dimensions' instead.",
			},
			"metric_dimensions": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ConflictsWith:    []string{"dimensions"},
				DiffSuppressFunc: CmsAlarmDiffSuppressFunc,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"escalations_critical": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Availability, Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientCriticalSuppressFunc,
				MaxItems:         1,
			},
			"escalations_warn": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Availability, Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientWarnSuppressFunc,
				MaxItems:         1,
			},
			"escalations_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"statistics": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      Average,
							ValidateFunc: validation.StringInSlice([]string{Availability, Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Equal,
							ValidateFunc: validation.StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
							}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
				DiffSuppressFunc: cmsClientInfoSuppressFunc,
				MaxItems:         1,
			},
			"statistics": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{Average, Minimum, Maximum, ErrorCodeMaximum, Value, Sum, Count}, false),
				Deprecated:   "Field 'statistics' has been deprecated from provider version 1.94.0. New field 'escalations_critical.statistics' instead.",
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, Equal, NotEqual,
				}, false),
				Deprecated: "Field 'operator' has been deprecated from provider version 1.94.0. New field 'escalations_critical.comparison_operator' instead.",
			},
			"threshold": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'threshold' has been deprecated from provider version 1.94.0. New field 'escalations_critical.threshold' instead.",
			},
			"triggered_count": {
				Type:       schema.TypeInt,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'triggered_count' has been deprecated from provider version 1.94.0. New field 'escalations_critical.times' instead.",
			},
			"contact_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:      0,
				//ValidateFunc: validation.IntBetween(0, 24),
				Deprecated: "Field 'start_time' has been deprecated from provider version 1.50.0. New field 'effective_interval' instead.",
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:      24,
				//ValidateFunc: validation.IntBetween(0, 24),
				Deprecated: "Field 'end_time' has been deprecated from provider version 1.50.0. New field 'effective_interval' instead.",
			},
			"effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "00:00-23:59",
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validation.IntBetween(300, 86400),
			},

			"notify_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Removed:      "Field 'notify_type' has been removed from provider version 1.50.0.",
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudCmsAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(resource.UniqueId())
	return resourceAlicloudCmsAlarmUpdate(d, meta)
}

func resourceAlicloudCmsAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	alarm, err := cmsService.DescribeAlarm(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", alarm["RuleName"])
	d.Set("project", alarm["Namespace"])
	d.Set("metric", alarm["MetricName"])
	d.Set("period", alarm["Period"])

	if escalations, ok := alarm["Escalations"].(map[string]interface{}); ok {
		if critical, ok := escalations["Critical"].(map[string]interface{}); ok {
			d.Set("statistics", critical["Statistics"])
			oper := convertOperator(fmt.Sprint(critical["ComparisonOperator"]))
			if oper == MoreThan && d.Get("operator").(string) == Equal {
				oper = Equal
			}
			d.Set("operator", oper)

			d.Set("threshold", critical["Threshold"])
			d.Set("triggered_count", critical["Times"])

			mapping := map[string]interface{}{
				"statistics":          critical["Statistics"],
				"comparison_operator": convertOperator(fmt.Sprint(critical["ComparisonOperator"])),
				"threshold":           critical["Threshold"],
				"times":               formatInt(critical["Times"]),
			}
			d.Set("escalations_critical", []map[string]interface{}{mapping})
		}
		if warn, ok := escalations["Warn"].(map[string]interface{}); ok {
			if warn["Times"] != "" {
				mappingWarn := map[string]interface{}{
					"statistics":          warn["Statistics"],
					"comparison_operator": convertOperator(fmt.Sprint(warn["ComparisonOperator"])),
					"threshold":           warn["Threshold"],
					"times":               formatInt(warn["Times"]),
				}
				d.Set("escalations_warn", []map[string]interface{}{mappingWarn})
			}
		}
		if info, ok := escalations["Info"].(map[string]interface{}); ok {
			if fmt.Sprint(info["Times"]) != "" {
				mappingInfo := map[string]interface{}{
					"statistics":          info["Statistics"],
					"comparison_operator": convertOperator(fmt.Sprint(info["ComparisonOperator"])),
					"threshold":           info["Threshold"],
					"times":               formatInt(info["Times"]),
				}
				d.Set("escalations_info", []map[string]interface{}{mappingInfo})
			}
		}
	}

	d.Set("effective_interval", alarm["EffectiveInterval"])

	d.Set("silence_time", alarm["SilenceTime"])

	d.Set("status", alarm["AlertState"])
	d.Set("enabled", alarm["EnableState"])
	d.Set("webhook", alarm["Webhook"])
	d.Set("contact_groups", strings.Split(alarm["ContactGroups"].(string), ","))

	dims := make([]map[string]interface{}, 0)
	if fmt.Sprint(alarm["Resources"]) != "" {
		if err := json.Unmarshal([]byte(alarm["Resources"].(string)), &dims); err != nil {
			return fmt.Errorf("Unmarshaling Dimensions got an error: %#v.", err)
		}
	}

	dimensionList := make(map[string]interface{}, 0)
	for _, raw := range dims {
		for k, v := range raw {
			if dimensionListValue, ok := dimensionList[k]; ok {
				dimensionList[k] = fmt.Sprint(dimensionListValue.(string), ",", v.(string))
			} else {
				dimensionList[k] = v
			}
		}
	}

	if err := d.Set("dimensions", dimensionList); err != nil {
		return WrapError(err)
	}
	if err := d.Set("metric_dimensions", alarm["Resources"]); err != nil {
		return WrapError(err)
	}
	return nil
}

func resourceAlicloudCmsAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutResourceMetricRule"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	cmsService := CmsService{client}
	d.Partial(true)

	request["RuleId"] = d.Id()
	request["RuleName"] = d.Get("name").(string)
	request["Namespace"] = d.Get("project").(string)
	request["MetricName"] = d.Get("metric").(string)
	request["Period"] = strconv.Itoa(d.Get("period").(int))
	request["ContactGroups"] = strings.Join(expandStringList(d.Get("contact_groups").([]interface{})), ",")

	// 兼容弃用参数
	request["Escalations.Critical.Statistics"] = d.Get("statistics").(string)
	request["Escalations.Critical.ComparisonOperator"] = convertOperator(d.Get("operator").(string))
	if v, ok := d.GetOk("threshold"); ok && v.(string) != "" {
		request["Escalations.Critical.Threshold"] = v.(string)
	}
	request["Escalations.Critical.Threshold"] = d.Get("threshold").(string)
	request["Escalations.Critical.Times"] = requests.NewInteger(d.Get("triggered_count").(int))

	// Critical
	if v, ok := d.GetOk("escalations_critical"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request["Escalations.Critical.Statistics"] = val["statistics"].(string)
			request["Escalations.Critical.ComparisonOperator"] = convertOperator(fmt.Sprint(val["comparison_operator"]))
			request["Escalations.Critical.Threshold"] = val["threshold"].(string)
			request["Escalations.Critical.Times"] = requests.NewInteger(val["times"].(int))
		}
	}
	// Warn
	if v, ok := d.GetOk("escalations_warn"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request["Escalations.Warn.Statistics"] = val["statistics"].(string)
			request["Escalations.Warn.ComparisonOperator"] = convertOperator(fmt.Sprint(val["comparison_operator"]))
			request["Escalations.Warn.Threshold"] = val["threshold"].(string)
			request["Escalations.Warn.Times"] = requests.NewInteger(val["times"].(int))
		}
	}
	// Info
	if v, ok := d.GetOk("escalations_info"); ok && len(v.([]interface{})) != 0 {
		for _, val := range v.([]interface{}) {
			val := val.(map[string]interface{})
			request["Escalations.Info.Statistics"] = val["statistics"].(string)
			request["Escalations.Info.ComparisonOperator"] = convertOperator(fmt.Sprint(val["comparison_operator"]))
			request["Escalations.Info.Threshold"] = val["threshold"].(string)
			request["Escalations.Info.Times"] = requests.NewInteger(val["times"].(int))

		}
	}

	if v, ok := d.GetOk("effective_interval"); ok && v.(string) != "" {
		request["EffectiveInterval"] = v.(string)
	} else {
		start, startOk := d.GetOk("start_time")
		end, endOk := d.GetOk("end_time")
		if startOk && endOk && end.(int) > 0 {
			// The EffectiveInterval valid value between 00:00 and 23:59
			request["EffectiveInterval"] = fmt.Sprintf("%d:00-%d:59", start.(int), end.(int)-1)
		}
	}
	request["SilenceTime"] = requests.NewInteger(d.Get("silence_time").(int))

	if webhook, ok := d.GetOk("webhook"); ok && webhook.(string) != "" {
		request["Webhook"] = webhook.(string)
	}

	var dimList []map[string]string
	if dimensions, ok := d.GetOk("dimensions"); ok {
		for k, v := range dimensions.(map[string]interface{}) {
			values := strings.Split(v.(string), COMMA_SEPARATED)
			if len(values) > 0 {
				for _, vv := range values {
					dimList = append(dimList, map[string]string{k: Trim(vv)})
				}
			} else {
				dimList = append(dimList, map[string]string{k: Trim(v.(string))})
			}

		}
	}
	if len(dimList) > 0 {
		if bytes, err := json.Marshal(dimList); err != nil {
			return fmt.Errorf("Marshaling dimensions to json string got an error: %#v.", err)
		} else {
			request["Resources"] = string(bytes[:])
		}
	}

	if v, ok := d.GetOk("metric_dimensions"); ok && v.(string) != "" {
		request["Resources"] = v.(string)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if err != nil {
		return fmt.Errorf("Putting alarm got an error: %#v", err)
	}
	d.SetPartial("name")
	d.SetPartial("period")
	d.SetPartial("statistics")
	d.SetPartial("operator")
	d.SetPartial("threshold")
	d.SetPartial("triggered_count")
	d.SetPartial("contact_groups")
	d.SetPartial("effective_interval")
	d.SetPartial("start_time")
	d.SetPartial("end_time")
	d.SetPartial("silence_time")
	d.SetPartial("notify_type")
	d.SetPartial("webhook")

	if d.Get("enabled").(bool) {
		action := "EnableMetricRules"
		enableMetricRequest := make(map[string]interface{})
		enableMetricRequest["RuleId"] = []string{d.Id()}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, enableMetricRequest, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, enableMetricRequest)
	} else {
		action := "DisableMetricRules"
		disableMetricRequest := make(map[string]interface{})
		disableMetricRequest["RuleId"] = []string{d.Id()}

		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, disableMetricRequest, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{ThrottlingUser}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, disableMetricRequest)
	}
	if err := cmsService.WaitForCmsAlarm(d.Id(), d.Get("enabled").(bool), 102); err != nil {
		return err
	}

	d.Partial(false)

	return resourceAlicloudCmsAlarmRead(d, meta)
}

func resourceAlicloudCmsAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	action := "DeleteMetricRules"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Id": []string{d.Id()},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{ThrottlingUser}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		_, err = cmsService.DescribeAlarm(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe alarm rule got an error: %#v", err))
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertOperator(operator string) string {
	switch operator {
	case MoreThan:
		return "GreaterThanThreshold"
	case MoreThanOrEqual:
		return "GreaterThanOrEqualToThreshold"
	case LessThan:
		return "LessThanThreshold"
	case LessThanOrEqual:
		return "LessThanOrEqualToThreshold"
	case NotEqual:
		return "NotEqualToThreshold"
	case Equal:
		return "GreaterThanThreshold"
	case "GreaterThanThreshold":
		return MoreThan
	case "GreaterThanOrEqualToThreshold":
		return MoreThanOrEqual
	case "LessThanThreshold":
		return LessThan
	case "LessThanOrEqualToThreshold":
		return LessThanOrEqual
	case "NotEqualToThreshold":
		return NotEqual
	default:
		return ""
	}
}
