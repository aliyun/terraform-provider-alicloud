package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceAlicloudEssSuspendProcess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssSuspendCreate,
		Read:   resourceAliyunEssSuspendRead,
		Delete: resourceAliyunEssSuspendDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"process": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	d.SetId(fmt.Sprint(d.Get("scaling_group_id").(string), ":", d.Get("process").(string)))
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.SuspendProcesses(request)
		})
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		time.Sleep(3 * time.Second)
		return nil
	}); err != nil {
		return WrapError(err)
	}
	return resourceAliyunEssSuspendRead(d, meta)
}

func resourceAliyunEssSuspendRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssScalingGroup(d.Get("scaling_group_id").(string))
	if err != nil {
		return WrapError(err)
	}
	var process []string
	process = object.SuspendedProcesses.SuspendedProcess
	if len(process) < 1 {
		return nil
	}
	strs, _ := ParseResourceId(d.Id(), 2)
	processTemp := strs[1]
	d.Set("process", processTemp)
	return nil
}

func resourceAliyunEssSuspendDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := ParseResourceId(d.Id(), 2)
	scalingGroupId, processTemp := strs[0], strs[1]
	var process []string
	process = append(process, processTemp)
	if len(process) < 1 {
		return nil
	}

	request := ess.CreateResumeProcessesRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.Process = &process
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := essService.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.ResumeProcesses(request)
		})
		if err != nil {
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		time.Sleep(3 * time.Second)
		return nil
	}); err != nil {
		return WrapError(err)
	}
	return nil
}
