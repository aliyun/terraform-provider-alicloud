package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEssSuspendProcess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssSuspendCreate,
		Read:   resourceAliyunEssSuspendRead,
		Delete: resourceAliyunEssSuspendDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"process": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ScaleIn", "ScaleOut", "HealthCheck", "AlarmNotification", "ScheduledAction"}, false),
			},
		},
	}
}

func resourceAliyunEssSuspendCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var process []string
	process = append(process, d.Get("process").(string))
	request := ess.CreateSuspendProcessesRequest()
	request.RegionId = client.RegionId
	request.Process = &process
	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	if err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.SuspendProcesses(request)
		})
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapError(err)
	}
	d.SetId(fmt.Sprint(d.Get("scaling_group_id").(string), ":", d.Get("process").(string)))
	return resourceAliyunEssSuspendRead(d, meta)
}

func resourceAliyunEssSuspendRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	_, err := essService.DescribeEssScalingGroupSuspendProcess(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	d.Set("scaling_group_id", parts[0])
	d.Set("process", parts[1])
	return nil
}

func resourceAliyunEssSuspendDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := ParseResourceId(d.Id(), 2)
	scalingGroupId, processTemp := strs[0], strs[1]
	var process []string
	process = append(process, processTemp)
	request := ess.CreateResumeProcessesRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.Process = &process
	if err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := essService.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.ResumeProcesses(request)
		})
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapError(err)
	}
	return nil
}
