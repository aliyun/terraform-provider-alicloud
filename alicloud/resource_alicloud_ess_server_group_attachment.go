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

func resourceAliCloudEssServerGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssServerGroupAttachmentCreate,
		Read:   resourceAliyunEssServerGroupAttachmentRead,
		Update: resourceAliyunEssServerGroupAttachmentUpdate,
		Delete: resourceAliyunEssServerGroupAttachmentDelete,
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
			"server_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ALB", "NLB"}, false),
				Required:     true,
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

func resourceAliyunEssServerGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	scalingGroupId := d.Get("scaling_group_id").(string)
	serverGroupId := d.Get("server_group_id").(string)
	typeAttribute := d.Get("type").(string)
	port := strconv.Itoa(formatInt(d.Get("port")))

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateAttachServerGroupsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.ForceAttach = requests.NewBoolean(d.Get("force_attach").(bool))
	attachScalingGroupServerGroups := make([]ess.AttachServerGroupsServerGroup, 0)
	attachScalingGroupServerGroups = append(attachScalingGroupServerGroups, ess.AttachServerGroupsServerGroup{
		ServerGroupId: serverGroupId,
		Port:          port,
		Weight:        strconv.Itoa(formatInt(d.Get("weight"))),
		Type:          typeAttribute,
	})
	request.ServerGroup = &attachScalingGroupServerGroups
	wait := incrementalWait(1*time.Second, 2*time.Second)

	var raw interface{}
	var err error
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.AttachServerGroups(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus", "InvalidOperation.Conflict", "ScalingActivityInProgress"}) || NeedRetry(err) {
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

	response, _ := raw.(*ess.AttachServerGroupsResponse)

	d.SetId(fmt.Sprint(scalingGroupId, ":", serverGroupId, ":", typeAttribute, ":", port))
	if len(response.ScalingActivityId) == 0 {
		return resourceAliyunEssServerGroupAttachmentRead(d, meta)
	}
	essService := EssService{client}
	stateConf := BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(response.ScalingActivityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliyunEssServerGroupAttachmentRead(d, meta)
}

func resourceAliyunEssServerGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return WrapErrorf(Error("server_group_attachment not support modify operation"), DefaultErrorMsg, "alicloud_ess_server_groups", "Modify", AlibabaCloudSdkGoERROR)
}

func resourceAliyunEssServerGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := ParseResourceId(d.Id(), 4)
	scalingGroupId, serverGroupId, typeAttribute, port := strs[0], strs[1], strs[2], strs[3]

	object, err := essService.DescribeEssScalingGroup(scalingGroupId)
	if err != nil {
		return WrapError(err)
	}

	for _, v := range object.ServerGroups.ServerGroup {
		if v.ServerGroupId == serverGroupId && v.Port == formatInt(port) && v.Type == typeAttribute {
			d.Set("scaling_group_id", object.ScalingGroupId)
			d.Set("type", v.Type)
			d.Set("server_group_id", v.ServerGroupId)
			d.Set("weight", v.Weight)
			d.Set("port", v.Port)
			return nil
		}
	}
	return WrapErrorf(NotFoundErr("ServerGroup", d.Id()), NotFoundMsg, ProviderERROR)
}

func resourceAliyunEssServerGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDetachServerGroupsRequest()
	request.RegionId = client.RegionId
	strs, _ := ParseResourceId(d.Id(), 4)
	scalingGroupId, serverGroupId, typeAttribute, port := strs[0], strs[1], strs[2], strs[3]

	request.ScalingGroupId = scalingGroupId
	request.ForceDetach = requests.NewBoolean(d.Get("force_attach").(bool))
	detachScalingGroupServerGroups := make([]ess.DetachServerGroupsServerGroup, 0)
	detachScalingGroupServerGroups = append(detachScalingGroupServerGroups, ess.DetachServerGroupsServerGroup{
		ServerGroupId: serverGroupId,
		Port:          port,
		Type:          typeAttribute,
	})
	request.ServerGroup = &detachScalingGroupServerGroups

	activityId := ""
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DetachServerGroups(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus", "InvalidOperation.Conflict", "ScalingActivityInProgress"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		response, _ := raw.(*ess.DetachServerGroupsResponse)
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
	stateConf := BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, essService.ActivityStateRefreshFunc(activityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		if strings.Contains(err.Error(), "activity not found") {
			return nil
		}
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
