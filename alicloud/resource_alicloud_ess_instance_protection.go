package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceAlicloudEssInstanceProtection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssInstanceProtectionCreate,
		Read:   resourceAliyunEssInstanceProtectionRead,
		Delete: resourceAliyunEssInstanceProtectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceAliyunEssInstanceProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	scalingGroupId := d.Get("scaling_group_id").(string)
	instanceId := d.Get("instance_id").(string)
	instanceIds := []string{instanceId}
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssScalingGroup(scalingGroupId)
	if err != nil {
		return WrapError(err)
	}
	if object.LifecycleState == string(Inactive) {
		return WrapError(Error("Scaling group current status is %s, please active it before setting instance protection.", object.LifecycleState))
	} else {
		if err := essService.WaitForEssScalingGroup(object.ScalingGroupId, Active, DefaultTimeout); err != nil {
			return WrapError(err)
		}
	}
	request := ess.CreateSetInstancesProtectionRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.InstanceId = &instanceIds
	request.ProtectedFromScaleIn = "true"

	d.SetId(fmt.Sprint(scalingGroupId, ":", instanceId))
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.SetInstancesProtection(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus"}) {
				time.Sleep(5)
				return resource.RetryableError(WrapError(err))
			}
			return resource.NonRetryableError(WrapError(err))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapError(err)
	}
	return resourceAliyunEssInstanceProtectionRead(d, meta)
}

func resourceAliyunEssInstanceProtectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := ParseResourceId(d.Id(), 2)
	object, err := essService.DescribeInstanceByStatus(strs[0], string(Protected))
	if err != nil {
		return WrapError(err)
	}
	var instanceId string
	for _, inst := range object {
		if inst.InstanceId == strs[1] {
			instanceId = inst.InstanceId
		}
	}
	d.Set("scaling_group_id", object[0].ScalingGroupId)
	d.Set("instance_id", instanceId)

	return nil
}

func resourceAliyunEssInstanceProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := ParseResourceId(d.Id(), 2)
	scalingGroupId, _ := strs[0], strs[1]

	object, err := essService.DescribeEssScalingGroup(scalingGroupId)
	if err != nil {
		return WrapError(err)
	}

	if err := essService.WaitForEssScalingGroup(object.ScalingGroupId, Active, DefaultTimeout); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	if err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request := ess.CreateSetInstancesProtectionRequest()
		strs, _ := ParseResourceId(d.Id(), 2)
		request.RegionId = client.RegionId
		request.ScalingGroupId = strs[0]
		request.ProtectedFromScaleIn = "false"
		instanceIds := []string{strs[1]}
		request.InstanceId = &instanceIds
		raw, err := essService.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.SetInstancesProtection(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"ScalingActivityInProgress", "IncorrectScalingGroupStatus"}) {
				time.Sleep(5)
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		time.Sleep(3 * time.Second)
		object, err := essService.DescribeInstanceByStatus(strs[0], string(Protected))
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		if len(object) > 0 {
			instanceIds = make([]string, 0)
			for _, inst := range object {
				if inst.InstanceId == strs[1] {
					return resource.RetryableError(WrapError(Error("ECS instance is still protection status in the scaling group.")))
				}
			}
		}
		return nil
	}); err != nil {
		return WrapError(err)
	}
	return nil
}
