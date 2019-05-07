package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEssAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssAlarmCreate,
		Read:   resourceAliyunEssAlarmRead,
		Update: resourceAliyunEssAlarmUpdate,
		Delete: resourceAliyunEssAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_actions": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 5,
				MinItems: 1,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  System,
				ValidateFunc: validateAllowedStringValue([]string{
					string(System),
					string(Custom),
				}),
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  FiveMinite,
				ValidateFunc: validateAllowedIntValue([]int{
					int(OneMinite),
					int(TwoMinite),
					int(FiveMinite),
					int(FifteenMinite),
				}),
			},
			"statistics": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Average,
				ValidateFunc: validateAllowedStringValue([]string{
					string(Average),
					string(Minimum),
					string(Maximum),
				}),
			},
			"threshold": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comparison_operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  Gte,
				ValidateFunc: validateAllowedStringValue([]string{
					string(Gt),
					string(Gte),
					string(Lt),
					string(Lte),
				}),
			},
			"evaluation_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validateEvaluationCount,
			},
			"cloud_monitor_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"dimensions": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunEssAlarmCreate(d *schema.ResourceData, meta interface{}) error {

	args, err := buildAlicloudEssAlarmArgs(d)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateAlarm(args)
		})
		if err != nil {
			if IsExceptedError(err, EssThrottling) {
				return resource.RetryableError(fmt.Errorf("CreateAlarm timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateAlarm got an error: %#v.", err))
		}
		alarm, _ := raw.(*ess.CreateAlarmResponse)
		d.SetId(alarm.AlarmTaskId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunEssAlarmRead(d, meta)
}

func resourceAliyunEssAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	alarm, err := essService.DescribeEssAlarmById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS Alarm Attribute: %#v", err)
	}

	d.Set("name", alarm.Name)
	d.Set("description", alarm.Description)
	d.Set("alarm_actions", alarm.AlarmActions.AlarmAction)
	d.Set("scaling_group_id", alarm.ScalingGroupId)
	d.Set("metric_type", alarm.MetricType)
	d.Set("metric_name", alarm.MetricName)
	d.Set("period", alarm.Period)
	d.Set("statistics", alarm.Statistics)
	d.Set("threshold", strconv.FormatFloat(alarm.Threshold, 'f', -1, 32))
	d.Set("comparison_operator", alarm.ComparisonOperator)
	d.Set("evaluation_count", alarm.EvaluationCount)
	d.Set("state", alarm.State)

	dims := make([]ess.Dimension, 0, len(alarm.Dimensions.Dimension))
	for _, dimension := range alarm.Dimensions.Dimension {
		if dimension.DimensionKey == GroupId {
			d.Set("cloud_monitor_group_id", dimension.DimensionValue)
		} else {
			dims = append(dims, dimension)
		}
	}

	if err := d.Set("dimensions", essService.flattenDimensionsToMap(dims)); err != nil {
		return err
	}

	return nil
}

func resourceAliyunEssAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := ess.CreateModifyAlarmRequest()
	args.AlarmTaskId = d.Id()

	if d.HasChange("name") {
		args.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		args.Description = d.Get("description").(string)
	}

	if d.HasChange("alarm_actions") {
		if v, ok := d.GetOk("alarm_actions"); ok {
			alarmActions := expandStringList(v.(*schema.Set).List())
			if len(alarmActions) > 0 {
				args.AlarmAction = &alarmActions
			}
		}
	}

	_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyAlarm(args)
	})
	if err != nil {
		return err
	}

	return resourceAliyunEssAlarmRead(d, meta)
}

func resourceAliyunEssAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	id := d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := ess.CreateDeleteAlarmRequest()
		req.AlarmTaskId = id

		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteAlarm(req)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidEssAlarmTaskNotFound}) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete ess alarm timeout and got an error:%#v.", err))
		}
		_, err = essService.DescribeEssAlarmById(id)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(fmt.Errorf("Delete ess alarm timeout and got an error:%#v.", err))
	})
}

func buildAlicloudEssAlarmArgs(d *schema.ResourceData) (*ess.CreateAlarmRequest, error) {
	args := ess.CreateCreateAlarmRequest()

	if name := d.Get("name").(string); name != "" {
		args.Name = name
	}

	if description := d.Get("description").(string); description != "" {
		args.Description = description
	}

	if v, ok := d.GetOk("alarm_actions"); ok {
		alarmActions := expandStringList(v.(*schema.Set).List())
		args.AlarmAction = &alarmActions
	}

	if scalingGroupId := d.Get("scaling_group_id").(string); scalingGroupId != "" {
		args.ScalingGroupId = scalingGroupId
	}

	if metricType := d.Get("metric_type").(string); metricType != "" {
		args.MetricType = metricType
	}

	if metricName := d.Get("metric_name").(string); metricName != "" {
		args.MetricName = metricName
	}

	if period, ok := d.GetOk("period"); ok && period.(int) > 0 {
		args.Period = requests.NewInteger(period.(int))
	}

	if statistics := d.Get("statistics").(string); statistics != "" {
		args.Statistics = statistics
	}

	if v, ok := d.GetOk("threshold"); ok {
		threshold, err := strconv.ParseFloat(v.(string), 32)
		if err != nil {
			return nil, err
		}
		args.Threshold = requests.NewFloat(threshold)
	}

	if comparisonOperator := d.Get("comparison_operator").(string); comparisonOperator != "" {
		args.ComparisonOperator = comparisonOperator
	}

	if evaluationCount, ok := d.GetOk("evaluation_count"); ok && evaluationCount.(int) > 0 {
		args.EvaluationCount = requests.NewInteger(evaluationCount.(int))
	}

	if groupId, ok := d.GetOk("cloud_monitor_group_id"); ok {
		args.GroupId = requests.NewInteger(groupId.(int))
	}

	if v, ok := d.GetOk("dimensions"); ok {
		dimensions := v.(map[string]interface{})
		createAlarmDimensions := make([]ess.CreateAlarmDimension, 0, len(dimensions))
		for k, v := range dimensions {
			if k == UserId || k == ScalingGroup {
				return nil, fmt.Errorf("Invalide dimension keys, %s", k)
			}
			if k != "" {
				dimension := ess.CreateAlarmDimension{
					DimensionKey:   k,
					DimensionValue: v.(string),
				}
				createAlarmDimensions = append(createAlarmDimensions, dimension)
			}
		}
		args.Dimension = &createAlarmDimensions
	}

	return args, nil
}
