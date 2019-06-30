package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudCmsAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAlarmUpdate,
		Read:   resourceAlicloudCmsAlarmRead,
		Update: resourceAlicloudCmsAlarmUpdate,
		Delete: resourceAlicloudCmsAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem:     schema.TypeString,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"statistics": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Average,
				ValidateFunc: validateAllowedStringValue([]string{
					string(Average), string(Minimum), string(Maximum),
				}),
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Equal,
				ValidateFunc: validateAllowedStringValue([]string{
					MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, Equal, NotEqual,
				}),
			},
			"threshold": {
				Type:     schema.TypeString,
				Required: true,
			},
			"triggered_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"contact_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 24),
				Deprecated:   "Field 'start_time' has been deprecated from provider version 1.50.0. New field 'effective_interval' instead.",
			},
			"end_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ValidateFunc: validateIntegerInRange(0, 24),
				Deprecated:   "Field 'end_time' has been deprecated from provider version 1.50.0. New field 'effective_interval' instead.",
			},
			"effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: validateIntegerInRange(300, 86400),
			},

			"notify_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
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

	d.Set("name", alarm.RuleName)
	d.Set("project", alarm.Namespace)
	d.Set("metric", alarm.MetricName)
	if period, err := strconv.Atoi(alarm.Period); err != nil {
		return fmt.Errorf("Atoi Period got an error: %#v.", err)
	} else {
		d.Set("period", period)
	}
	d.Set("statistics", alarm.Escalations.Critical.Statistics)
	oper := convertOperator(alarm.Escalations.Critical.ComparisonOperator)
	if oper == MoreThan && d.Get("operator").(string) == Equal {
		oper = Equal
	}
	d.Set("operator", oper)
	d.Set("threshold", alarm.Escalations.Critical.Threshold)
	if count, err := strconv.Atoi(alarm.Escalations.Critical.Times); err != nil {
		return fmt.Errorf("Atoi Escalations.Critical.Times got an error: %#v.", err)
	} else {
		d.Set("triggered_count", count)
	}
	d.Set("effective_interval", alarm.EffectiveInterval)
	//d.Set("start_time", parts[0])
	//d.Set("end_time", parts[1])
	if silence, err := strconv.Atoi(alarm.SilenceTime); err != nil {
		return fmt.Errorf("Atoi SilenceTime got an error: %#v.", err)
	} else {
		d.Set("silence_time", silence)
	}
	d.Set("status", alarm.AlertState)
	d.Set("enabled", alarm.EnableState)
	d.Set("webhook", alarm.Webhook)
	d.Set("contact_groups", strings.Split(alarm.ContactGroups, ","))

	var dims []string
	if alarm.Dimensions != "" {
		if err := json.Unmarshal([]byte(alarm.Dimensions), &dims); err != nil {
			return fmt.Errorf("Unmarshaling Dimensions got an error: %#v.", err)
		}
	}
	d.Set("dimensions", dims)

	return nil
}

func resourceAlicloudCmsAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	d.SetId(resource.UniqueId())
	d.Partial(true)

	request := cms.CreatePutResourceMetricRuleRequest()
	request.RuleId = d.Id()
	request.RuleName = d.Get("name").(string)
	request.Namespace = d.Get("project").(string)
	request.MetricName = d.Get("metric").(string)
	request.Period = strconv.Itoa(d.Get("period").(int))
	request.EscalationsCriticalStatistics = d.Get("statistics").(string)
	request.EscalationsCriticalComparisonOperator = convertOperator(d.Get("operator").(string))
	request.EscalationsCriticalThreshold = d.Get("threshold").(string)
	request.EscalationsCriticalTimes = requests.NewInteger(d.Get("triggered_count").(int))
	request.ContactGroups = strings.Join(expandStringList(d.Get("contact_groups").([]interface{})), ",")
	if v, ok := d.GetOk("effective_interval"); ok && v.(string) != "" {
		request.EffectiveInterval = v.(string)
	} else {
		request.EffectiveInterval = fmt.Sprintf("%d:00-%d:00", d.Get("start_time").(int), d.Get("end_time").(int))
	}
	request.SilenceTime = requests.NewInteger(d.Get("silence_time").(int))

	if webhook, ok := d.GetOk("webhook"); ok && webhook.(string) != "" {
		request.Webhook = webhook.(string)
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
			request.Resources = string(bytes[:])
		}
	}
	_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.PutResourceMetricRule(request)
	})
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
		request := cms.CreateEnableMetricRulesRequest()
		request.RuleId = &[]string{d.Id()}

		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.EnableMetricRules(request)
		})
		if err != nil {
			return fmt.Errorf("Enabling alarm got an error: %#v", err)
		}
	} else {
		request := cms.CreateDisableMetricRulesRequest()
		request.RuleId = &[]string{d.Id()}

		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DisableMetricRules(request)
		})
		if err != nil {
			return fmt.Errorf("Disableing alarm got an error: %#v", err)
		}
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
	request := cms.CreateDeleteMetricRulesRequest()

	request.Id = &[]string{d.Id()}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.DeleteMetricRules(request)
		})

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
		}

		_, err = cmsService.DescribeAlarm(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe alarm rule got an error: %#v", err))
		}

		return resource.RetryableError(fmt.Errorf("Deleting alarm rule got an error: %#v", err))
	})
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
