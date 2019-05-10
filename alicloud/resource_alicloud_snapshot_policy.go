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

	args := ecs.CreateCreateAutoSnapshotPolicyRequest()
	args.AutoSnapshotPolicyName = d.Get("name").(string)
	args.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	args.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	args.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateAutoSnapshotPolicy(args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_snapshot_policy", args.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	resp := raw.(*ecs.CreateAutoSnapshotPolicyResponse)
	d.SetId(resp.AutoSnapshotPolicyId)

	ecsService := EcsService{client}
	if err := ecsService.WaitForSnapshotPolicy(d.Id(), SnapshotPolicyNormal, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunSnapshotPolicyRead(d, meta)
}

func resourceAliyunSnapshotPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	snapshotPolicy, err := ecsService.DescribeSnapshotPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", snapshotPolicy.AutoSnapshotPolicyName)
	weekdays, err := convertJsonStringToList(snapshotPolicy.RepeatWeekdays)
	if err != nil {
		return WrapError(err)
	}
	d.Set("repeat_weekdays", weekdays)
	d.Set("retention_days", snapshotPolicy.RetentionDays)
	timePoints, err := convertJsonStringToList(snapshotPolicy.TimePoints)
	if err != nil {
		return WrapError(err)
	}
	d.Set("time_points", timePoints)

	return nil
}

func resourceAliyunSnapshotPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ecs.CreateModifyAutoSnapshotPolicyExRequest()
	args.AutoSnapshotPolicyId = d.Id()
	if d.HasChange("name") {
		args.AutoSnapshotPolicyName = d.Get("name").(string)
	}
	if d.HasChange("repeat_weekdays") {
		args.RepeatWeekdays = convertListToJsonString(d.Get("repeat_weekdays").(*schema.Set).List())
	}
	if d.HasChange("retention_dasy") {
		args.RetentionDays = requests.NewInteger(d.Get("retention_days").(int))
	}
	if d.HasChange("time_points") {
		args.TimePoints = convertListToJsonString(d.Get("time_points").(*schema.Set).List())
	}
	_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyAutoSnapshotPolicyEx(args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunSnapshotPolicyRead(d, meta)
}

func resourceAliyunSnapshotPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	args := ecs.CreateDeleteAutoSnapshotPolicyRequest()
	args.AutoSnapshotPolicyId = d.Id()

	err := resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteAutoSnapshotPolicy(args)
		})
		if err != nil {
			if IsExceptedErrors(err, SnapshotPolicyInvalidOperations) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), args.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(ecsService.WaitForSnapshotPolicy(d.Id(), Deleted, DefaultTimeout))
}
