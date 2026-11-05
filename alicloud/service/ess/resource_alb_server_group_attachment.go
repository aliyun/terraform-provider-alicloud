package ess

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
)

func resourceAlicloudEssAlbServerGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssAlbServerGroupAttachmentCreate,
		Read:   resourceAliyunEssAlbServerGroupAttachmentRead,
		Update: resourceAliyunEssAlbServerGroupAttachmentUpdate,
		Delete: resourceAliyunEssAlbServerGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"alb_server_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"force_attach": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAliyunEssAlbServerGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	scalingGroupId := d.Get("scaling_group_id").(string)
	albServerGroupId := d.Get("alb_server_group_id").(string)
	portNumber, err := convert.FormatInt(d.Get("port"))
	if err != nil {
		return err
	}
	port := strconv.Itoa(portNumber)

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateAttachAlbServerGroupsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.ForceAttach = requests.NewBoolean(d.Get("force_attach").(bool))
	attachScalingGroupAlbServerGroups := make([]ess.AttachAlbServerGroupsAlbServerGroup, 0)

	weightValue, err := convert.FormatInt(d.Get("weight"))
	if err != nil {
		return err
	}
	attachScalingGroupAlbServerGroups = append(attachScalingGroupAlbServerGroups, ess.AttachAlbServerGroupsAlbServerGroup{
		AlbServerGroupId: albServerGroupId,
		Port:             port,
		Weight:           strconv.Itoa(weightValue),
	})
	request.AlbServerGroup = &attachScalingGroupAlbServerGroups
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)

	var raw interface{}
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.AttachAlbServerGroups(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus", "InvalidOperation.Conflict", "ScalingActivityInProgress"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ess.AttachAlbServerGroupsResponse)

	d.SetId(fmt.Sprint(scalingGroupId, ":", albServerGroupId, ":", port))
	if len(response.ScalingActivityId) == 0 {
		return resourceAliyunEssAlbServerGroupAttachmentRead(d, meta)
	}
	essService := EssService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(response.ScalingActivityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceAliyunEssAlbServerGroupAttachmentRead(d, meta)
}

func resourceAliyunEssAlbServerGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return errs.WrapErrorf(errs.Error("alb_server_group_attachment not support modify operation"), errs.DefaultErrorMsg, "alicloud_ess_alb_server_groups", "Modify", errs.AlibabaCloudSdkGoERROR)
}

func resourceAliyunEssAlbServerGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := helper.ParseResourceId(d.Id(), 3)
	scalingGroupId, albServerGroupId, port := strs[0], strs[1], strs[2]

	object, err := essService.DescribeEssScalingGroup(scalingGroupId)
	if err != nil {
		return errs.WrapError(err)
	}

	for _, v := range object.AlbServerGroups.AlbServerGroup {
		if v.AlbServerGroupId == albServerGroupId {
			portVal, portErr := convert.FormatInt(port)
			if portErr != nil {
				return portErr
			}
			if v.Port == portVal {
				d.Set("scaling_group_id", object.ScalingGroupId)
				d.Set("alb_server_group_id", v.AlbServerGroupId)
				d.Set("weight", v.Weight)
				d.Set("port", v.Port)
				return nil
			}
		}
	}
	return errs.WrapErrorf(errs.NotFoundErr("AlbServerGroup", d.Id()), errs.NotFoundMsg, errs.ProviderERROR)
}

func resourceAliyunEssAlbServerGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDetachAlbServerGroupsRequest()
	request.RegionId = client.RegionId
	strs, _ := helper.ParseResourceId(d.Id(), 3)
	scalingGroupId, albServerGroupId, port := strs[0], strs[1], strs[2]

	request.ScalingGroupId = scalingGroupId
	request.ForceDetach = requests.NewBoolean(d.Get("force_attach").(bool))
	detachScalingGroupAlbServerGroups := make([]ess.DetachAlbServerGroupsAlbServerGroup, 0)
	detachScalingGroupAlbServerGroups = append(detachScalingGroupAlbServerGroups, ess.DetachAlbServerGroupsAlbServerGroup{
		AlbServerGroupId: albServerGroupId,
		Port:             port,
	})
	request.AlbServerGroup = &detachScalingGroupAlbServerGroups
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)

	activityId := ""
	err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DetachAlbServerGroups(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus", "InvalidOperation.Conflict", "ScalingActivityInProgress"}) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response, _ := raw.(*ess.DetachAlbServerGroupsResponse)
		activityId = response.ScalingActivityId
		if len(response.ScalingActivityId) == 0 {
			return nil
		}

		logs.AddDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "The specified value of parameter \"ScalingGroupId\" is not valid") {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR)
	}
	essService := EssService{client}
	if activityId == "" {
		return nil
	}
	stateConf := helper.BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(activityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		if strings.Contains(err.Error(), "activity not found") {
			return nil
		}
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
