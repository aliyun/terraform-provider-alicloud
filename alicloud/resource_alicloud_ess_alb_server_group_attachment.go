package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	port := strconv.Itoa(formatInt(d.Get("port")))

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateAttachAlbServerGroupsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.ForceAttach = requests.NewBoolean(d.Get("force_attach").(bool))
	attachScalingGroupAlbServerGroups := make([]ess.AttachAlbServerGroupsAlbServerGroup, 0)
	attachScalingGroupAlbServerGroups = append(attachScalingGroupAlbServerGroups, ess.AttachAlbServerGroupsAlbServerGroup{
		AlbServerGroupId: albServerGroupId,
		Port:             port,
		Weight:           strconv.Itoa(formatInt(d.Get("weight"))),
	})
	request.AlbServerGroup = &attachScalingGroupAlbServerGroups
	wait := incrementalWait(1*time.Second, 2*time.Second)

	var raw interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.AttachAlbServerGroups(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ess.AttachAlbServerGroupsResponse)

	d.SetId(fmt.Sprint(scalingGroupId, ":", albServerGroupId, ":", port))
	if len(response.ScalingActivityId) == 0 {
		return resourceAliyunEssAlbServerGroupAttachmentRead(d, meta)
	}
	essService := EssService{client}
	stateConf := BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(response.ScalingActivityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliyunEssAlbServerGroupAttachmentRead(d, meta)
}

func resourceAliyunEssAlbServerGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return WrapErrorf(Error("alb_server_group_attachment not support modify operation"), DefaultErrorMsg, "alicloud_ess_alb_server_groups", "Modify", AlibabaCloudSdkGoERROR)
}

func resourceAliyunEssAlbServerGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := ParseResourceId(d.Id(), 3)
	scalingGroupId, albServerGroupId, port := strs[0], strs[1], strs[2]

	object, err := essService.DescribeEssScalingGroup(scalingGroupId)
	if err != nil {
		return WrapError(err)
	}

	for _, v := range object.AlbServerGroups.AlbServerGroup {
		if v.AlbServerGroupId == albServerGroupId && v.Port == formatInt(port) {
			d.Set("scaling_group_id", object.ScalingGroupId)
			d.Set("alb_server_group_id", v.AlbServerGroupId)
			d.Set("weight", v.Weight)
			d.Set("port", v.Port)
			return nil
		}
	}
	return WrapErrorf(NotFoundErr("AlbServerGroup", d.Id()), NotFoundMsg, ProviderERROR)
}

func resourceAliyunEssAlbServerGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDetachAlbServerGroupsRequest()
	request.RegionId = client.RegionId
	strs, _ := ParseResourceId(d.Id(), 3)
	scalingGroupId, albServerGroupId, port := strs[0], strs[1], strs[2]

	request.ScalingGroupId = scalingGroupId
	request.ForceDetach = requests.NewBoolean(d.Get("force_attach").(bool))
	detachScalingGroupAlbServerGroups := make([]ess.DetachAlbServerGroupsAlbServerGroup, 0)
	detachScalingGroupAlbServerGroups = append(detachScalingGroupAlbServerGroups, ess.DetachAlbServerGroupsAlbServerGroup{
		AlbServerGroupId: albServerGroupId,
		Port:             port,
	})
	request.AlbServerGroup = &detachScalingGroupAlbServerGroups

	activityId := ""
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DetachAlbServerGroups(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, _ := raw.(*ess.DetachAlbServerGroupsResponse)
		activityId = response.ScalingActivityId
		if len(response.ScalingActivityId) == 0 {
			return nil
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "The specified value of parameter \"ScalingGroupId\" is not valid") {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	essService := EssService{client}
	if activityId == "" {
		return nil
	}
	stateConf := BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(activityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		if strings.Contains(err.Error(), "activity not found") {
			return nil
		}
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
