package alicloud

import (
	"fmt"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudEssLifecycleHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssLifeCycleHookCreate,
		Read:   resourceAliyunEssLifeCycleHookRead,
		Update: resourceAliyunEssLifeCycleHookUpdate,
		Delete: resourceAliyunEssLifeCycleHookDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lifecycle_transition": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateLifecycleTransaction,
			},
			"heartbeat_timeout": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      600,
				ValidateFunc: validateIntegerInRange(30, 21600),
			},
			"default_result": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      Continue,
				ValidateFunc: validateActionResult,
			},
			"notification_arn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"notification_metadata": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunEssLifeCycleHookCreate(d *schema.ResourceData, meta interface{}) error {

	args := buildAlicloudEssLifeCycleHookArgs(d, meta)
	essconn := meta.(*AliyunClient).essconn

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		hook, err := essconn.CreateLifecycleHook(args)
		if err != nil {
			if IsExceptedError(err, EssThrottling) {
				return resource.RetryableError(fmt.Errorf("CreateLifecycleHook timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("CreateLifecycleHook got an error: %#v.", err))
		}
		d.SetId(hook.LifecycleHookId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunEssLifeCycleHookRead(d, meta)
}

func resourceAliyunEssLifeCycleHookRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	hook, err := client.DescribeLifecycleHookById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS Lifecycle Hook Attribute: %#v", err)
	}

	d.Set("scaling_group_id", hook.ScalingGroupId)
	d.Set("name", hook.LifecycleHookName)
	d.Set("lifecycle_transition", hook.LifecycleTransition)
	d.Set("heartbeat_timeout", hook.HeartbeatTimeout)
	d.Set("default_result", hook.DefaultResult)
	d.Set("notification_arn", hook.NotificationArn)
	d.Set("notification_metadata", hook.NotificationMetadata)

	return nil
}

func resourceAliyunEssLifeCycleHookUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).essconn
	args := ess.CreateModifyLifecycleHookRequest()
	args.LifecycleHookId = d.Id()

	if d.HasChange("lifecycle_transition") {
		args.LifecycleTransition = d.Get("lifecycle_transition").(string)
	}

	if d.HasChange("heartbeat_timeout") {
		args.HeartbeatTimeout = requests.NewInteger(d.Get("heartbeat_timeout").(int))
	}

	if d.HasChange("default_result") {
		args.DefaultResult = d.Get("default_result").(string)
	}

	if d.HasChange("notification_arn") {
		args.NotificationArn = d.Get("notification_arn").(string)
	}

	if d.HasChange("notification_metadata") {
		args.NotificationMetadata = d.Get("notification_metadata").(string)
	}

	if _, err := conn.ModifyLifecycleHook(args); err != nil {
		return err
	}

	return resourceAliyunEssLifeCycleHookRead(d, meta)
}

func resourceAliyunEssLifeCycleHookDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	id := d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		req := ess.CreateDeleteLifecycleHookRequest()
		req.LifecycleHookId = id

		_, err := client.essconn.DeleteLifecycleHook(req)
		if err != nil {
			if IsExceptedErrors(err, []string{InvalidLifecycleHookIdNotFound}) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete lifecycle hook  timeout and got an error:%#v.", err))
		}
		_, err = client.DescribeLifecycleHookById(id)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Delete lifecycle hook timeout and got an error:%#v.", err))
	})
}

func buildAlicloudEssLifeCycleHookArgs(d *schema.ResourceData, meta interface{}) *ess.CreateLifecycleHookRequest {
	args := ess.CreateCreateLifecycleHookRequest()

	args.ScalingGroupId = d.Get("scaling_group_id").(string)

	if name := d.Get("name").(string); name != "" {
		args.LifecycleHookName = name
	}

	if transition := d.Get("lifecycle_transition").(string); transition != "" {
		args.LifecycleTransition = transition
	}

	if timeout, ok := d.GetOk("heartbeat_timeout"); ok && timeout.(int) > 0 {
		args.HeartbeatTimeout = requests.NewInteger(timeout.(int))
	}

	if result := d.Get("default_result").(string); result != "" {
		args.DefaultResult = result
	}

	if arn := d.Get("notification_arn").(string); arn != "" {
		args.NotificationArn = arn
	}

	if metadata := d.Get("notification_metadata").(string); meta != "" {
		args.NotificationMetadata = metadata
	}

	return args
}
