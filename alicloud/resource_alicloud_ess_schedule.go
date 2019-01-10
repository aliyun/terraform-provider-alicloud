package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEssSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScheduleCreate,
		Read:   resourceAliyunEssScheduleRead,
		Update: resourceAliyunEssScheduleUpdate,
		Delete: resourceAliyunEssScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scheduled_action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"launch_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scheduled_task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"launch_expiration_time": {
				Type:         schema.TypeInt,
				Default:      600,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 21600),
			},
			"recurrence_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Daily), string(Weekly), string(Monthly)}),
			},
			"recurrence_value": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"recurrence_end_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"task_enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceAliyunEssScheduleCreate(d *schema.ResourceData, meta interface{}) error {

	args, err := buildAlicloudEssScheduleArgs(d, meta)
	if err != nil {
		return err
	}

	client := meta.(*connectivity.AliyunClient)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScheduledTask(args)
	})
	if err != nil {
		return err
	}
	rule, _ := raw.(*ess.CreateScheduledTaskResponse)
	d.SetId(rule.ScheduledTaskId)

	return resourceAliyunEssScheduleUpdate(d, meta)
}

func resourceAliyunEssScheduleRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	rule, err := essService.DescribeScheduleById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS schedule Attribute: %#v", err)
	}

	d.Set("scheduled_action", rule.ScheduledAction)
	d.Set("launch_time", rule.LaunchTime)
	d.Set("scheduled_task_name", rule.ScheduledTaskName)
	d.Set("description", rule.Description)
	d.Set("launch_expiration_time", rule.LaunchExpirationTime)
	d.Set("recurrence_type", rule.RecurrenceType)
	d.Set("recurrence_value", rule.RecurrenceValue)
	d.Set("recurrence_end_time", rule.RecurrenceEndTime)
	d.Set("task_enabled", rule.TaskEnabled)

	return nil
}

func resourceAliyunEssScheduleUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	args := ess.CreateModifyScheduledTaskRequest()
	args.ScheduledTaskId = d.Id()

	if d.HasChange("scheduled_task_name") {
		args.ScheduledTaskName = d.Get("scheduled_task_name").(string)
	}

	if d.HasChange("description") {
		args.Description = d.Get("description").(string)
	}

	if d.HasChange("scheduled_action") {
		args.ScheduledAction = d.Get("scheduled_action").(string)
	}

	if d.HasChange("launch_time") {
		args.LaunchTime = d.Get("launch_time").(string)
	}

	if d.HasChange("launch_expiration_time") {
		args.LaunchExpirationTime = requests.NewInteger(d.Get("launch_expiration_time").(int))
	}

	if d.HasChange("recurrence_type") {
		args.RecurrenceType = d.Get("recurrence_type").(string)
	}

	if d.HasChange("recurrence_value") {
		args.RecurrenceValue = d.Get("recurrence_value").(string)
	}

	if d.HasChange("recurrence_end_time") {
		args.RecurrenceEndTime = d.Get("recurrence_end_time").(string)
	}

	if d.HasChange("task_enabled") {
		args.TaskEnabled = requests.NewBoolean(d.Get("task_enabled").(bool))
	}

	_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScheduledTask(args)
	})
	if err != nil {
		return err
	}

	return resourceAliyunEssScheduleRead(d, meta)
}

func resourceAliyunEssScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	return resource.Retry(2*time.Minute, func() *resource.RetryError {
		err := essService.DeleteScheduleById(d.Id())

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Delete scaling schedule timeout and got an error:%#v.", err))
		}

		_, err = essService.DescribeScheduleById(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Delete scaling schedule timeout and got an error:%#v.", err))
	})
}

func buildAlicloudEssScheduleArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScheduledTaskRequest, error) {
	args := ess.CreateCreateScheduledTaskRequest()
	args.ScheduledAction = d.Get("scheduled_action").(string)
	args.LaunchTime = d.Get("launch_time").(string)
	args.TaskEnabled = requests.NewBoolean(d.Get("task_enabled").(bool))

	if v := d.Get("scheduled_task_name").(string); v != "" {
		args.ScheduledTaskName = v
	}

	if v := d.Get("description").(string); v != "" {
		args.Description = v
	}

	if v := d.Get("recurrence_type").(string); v != "" {
		args.RecurrenceType = v
	}

	if v := d.Get("recurrence_value").(string); v != "" {
		args.RecurrenceValue = v
	}

	if v := d.Get("recurrence_end_time").(string); v != "" {
		args.RecurrenceEndTime = v
	}

	if v := d.Get("launch_expiration_time").(int); v != 0 {
		args.LaunchExpirationTime = requests.NewInteger(v)
	}

	return args, nil
}
