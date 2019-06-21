package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSnapshotPolicyCreate,
		Read:   resourceAliyunSnapshotPolicyRead,
		Update: resourceAliyunSnapshotPolicyUpdate,
		Delete: resourceAliyunSnapshotPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repeat_weekdays": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"time_points": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAliyunSnapshotPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateCreateAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyName = d.Get("name").(string)
	request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateAutoSnapshotPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_snapshot_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response := raw.(*ecs.CreateAutoSnapshotPolicyResponse)
	d.SetId(response.AutoSnapshotPolicyId)

	ecsService := EcsService{client}
	if err := ecsService.WaitForSnapshotPolicy(d.Id(), SnapshotPolicyNormal, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunSnapshotPolicyRead(d, meta)
}

func resourceAliyunSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeSnapshotPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.AutoSnapshotPolicyName)
	weekdays, err := convertJsonStringToList(object.RepeatWeekdays)
	if err != nil {
		return WrapError(err)
	}
	d.Set("repeat_weekdays", weekdays)
	d.Set("retention_days", object.RetentionDays)
	timePoints, err := convertJsonStringToList(object.TimePoints)
	if err != nil {
		return WrapError(err)
	}
	d.Set("time_points", timePoints)

	return nil
}

func resourceAliyunSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateModifyAutoSnapshotPolicyExRequest()
	request.AutoSnapshotPolicyId = d.Id()
	if d.HasChange("name") {
		request.AutoSnapshotPolicyName = d.Get("name").(string)
	}
	if d.HasChange("repeat_weekdays") {
		request.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	}
	if d.HasChange("retention_days") {
		request.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	}
	if d.HasChange("time_points") {
		request.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyAutoSnapshotPolicyEx(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	return resourceAliyunSnapshotPolicyRead(d, meta)
}

func resourceAliyunSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteAutoSnapshotPolicyRequest()
	request.AutoSnapshotPolicyId = d.Id()
	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteAutoSnapshotPolicy(request)
		})
		if err != nil {
			if IsExceptedErrors(err, SnapshotPolicyInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(ecsService.WaitForSnapshotPolicy(d.Id(), Deleted, DefaultTimeout))
}
