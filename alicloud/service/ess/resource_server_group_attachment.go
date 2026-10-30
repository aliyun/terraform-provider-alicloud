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
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
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
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"ALB", "NLB", "GWLB"}, false),
				Required:     true,
			},
			"weight": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
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
	portNumber, err := convert.FormatInt(d.Get("port"))
	if err != nil {
		return err
	}
	port := strconv.Itoa(portNumber)

	if typeAttribute == "GWLB" {
		if _, ok := d.GetOkExists("weight"); ok {
			return errs.WrapError(errs.Error("When the type is set to GWLB, both weight and port are prohibited. Please remove the configurations for weight and port!"))
		}
		if _, ok := d.GetOkExists("port"); ok {
			return errs.WrapError(errs.Error("When the type is set to GWLB, both weight and port are prohibited. Please remove the configurations for weight and port!"))
		}
	}

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateAttachServerGroupsRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = scalingGroupId
	request.ForceAttach = requests.NewBoolean(d.Get("force_attach").(bool))
	attachScalingGroupServerGroups := make([]ess.AttachServerGroupsServerGroup, 0)
	if typeAttribute == "GWLB" {
		attachScalingGroupServerGroups = append(attachScalingGroupServerGroups, ess.AttachServerGroupsServerGroup{
			ServerGroupId: serverGroupId,
			Type:          typeAttribute,
		})
	} else {

		weightValue, err := convert.FormatInt(d.Get("weight"))
		if err != nil {
			return err
		}
		attachScalingGroupServerGroups = append(attachScalingGroupServerGroups, ess.AttachServerGroupsServerGroup{
			ServerGroupId: serverGroupId,
			Port:          port,
			Weight:        strconv.Itoa(weightValue),
			Type:          typeAttribute,
		})
	}

	request.ServerGroup = &attachScalingGroupServerGroups
	wait := helper.IncrementalWait(1*time.Second, 2*time.Second)

	var raw interface{}
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		raw, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.AttachServerGroups(request)
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

	response, _ := raw.(*ess.AttachServerGroupsResponse)

	d.SetId(fmt.Sprint(scalingGroupId, ":", serverGroupId, ":", typeAttribute, ":", port))
	if len(response.ScalingActivityId) == 0 {
		return resourceAliyunEssServerGroupAttachmentRead(d, meta)
	}
	essService := EssService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, essService.ActivityStateRefreshFunc(response.ScalingActivityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return resourceAliyunEssServerGroupAttachmentRead(d, meta)
}

func resourceAliyunEssServerGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	return errs.WrapErrorf(errs.Error("server_group_attachment not support modify operation"), errs.DefaultErrorMsg, "alicloud_ess_server_groups", "Modify", errs.AlibabaCloudSdkGoERROR)
}

func resourceAliyunEssServerGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	strs, _ := helper.ParseResourceId(d.Id(), 4)
	scalingGroupId, serverGroupId, typeAttribute, port := strs[0], strs[1], strs[2], strs[3]

	object, err := essService.DescribeEssScalingGroup(scalingGroupId)
	if err != nil {
		return errs.WrapError(err)
	}

	for _, v := range object.ServerGroups.ServerGroup {
		if v.Type == "GWLB" {
			if v.ServerGroupId == serverGroupId && v.Type == typeAttribute {
				d.Set("scaling_group_id", object.ScalingGroupId)
				d.Set("type", v.Type)
				d.Set("server_group_id", v.ServerGroupId)
				return nil
			}
		} else {
			if v.ServerGroupId == serverGroupId {
				portVal, portErr := convert.FormatInt(port)
				if portErr != nil {
					return portErr
				}
				if v.Port == portVal && v.Type == typeAttribute {
					d.Set("scaling_group_id", object.ScalingGroupId)
					d.Set("type", v.Type)
					d.Set("server_group_id", v.ServerGroupId)
					d.Set("weight", v.Weight)
					d.Set("port", v.Port)
					return nil
				}
			}
		}

	}
	return errs.WrapErrorf(errs.NotFoundErr("ServerGroup", d.Id()), errs.NotFoundMsg, errs.ProviderERROR)
}

func resourceAliyunEssServerGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDetachServerGroupsRequest()
	request.RegionId = client.RegionId
	strs, _ := helper.ParseResourceId(d.Id(), 4)
	scalingGroupId, serverGroupId, typeAttribute, port := strs[0], strs[1], strs[2], strs[3]

	request.ScalingGroupId = scalingGroupId
	request.ForceDetach = requests.NewBoolean(d.Get("force_attach").(bool))
	detachScalingGroupServerGroups := make([]ess.DetachServerGroupsServerGroup, 0)
	if typeAttribute == "GWLB" {
		detachScalingGroupServerGroups = append(detachScalingGroupServerGroups, ess.DetachServerGroupsServerGroup{
			ServerGroupId: serverGroupId,
			Type:          typeAttribute,
		})
	} else {
		detachScalingGroupServerGroups = append(detachScalingGroupServerGroups, ess.DetachServerGroupsServerGroup{
			ServerGroupId: serverGroupId,
			Port:          port,
			Type:          typeAttribute,
		})
	}
	request.ServerGroup = &detachScalingGroupServerGroups

	activityId := ""
	err := retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DetachServerGroups(request)
		})
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus", "InvalidOperation.Conflict", "ScalingActivityInProgress"}) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		response, _ := raw.(*ess.DetachServerGroupsResponse)
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
	stateConf := helper.BuildStateConf([]string{}, []string{"Successful"}, d.Timeout(schema.TimeoutDelete), 1*time.Minute, essService.ActivityStateRefreshFunc(activityId, []string{"Failed", "Rejected"}))
	if _, err := stateConf.WaitForState(); err != nil {
		if strings.Contains(err.Error(), "activity not found") {
			return nil
		}
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}
	return nil
}
