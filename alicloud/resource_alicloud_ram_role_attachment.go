package alicloud

import (
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudInstanceRoleAttachmentCreate,
		Read:   resourceAlicloudInstanceRoleAttachmentRead,
		Delete: resourceAlicloudInstanceRoleAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudInstanceRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}

	instanceIds := convertListToJsonString(d.Get("instance_ids").(*schema.Set).List())

	request := ecs.CreateAttachInstanceRamRoleRequest()
	request.InstanceIds = instanceIds
	request.RamRoleName = d.Get("role_name").(string)

	err := ramService.JudgeRolePolicyPrincipal(request.RamRoleName)
	if err != nil {
		return WrapError(err)
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachInstanceRamRole(request)
		})
		if err != nil {
			if IsExceptedError(err, RoleAttachmentUnExpectedJson) {
				return resource.RetryableError(WrapError(Error("Please trying again.")))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, "ram_role_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		d.SetId(d.Get("role_name").(string) + ":" + instanceIds)
		return resource.NonRetryableError(WrapError(resourceAlicloudInstanceRoleAttachmentRead(d, meta)))
	})
}

func resourceAlicloudInstanceRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	roleName := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	request := ecs.CreateDescribeInstanceRamRoleRequest()
	request.InstanceIds = instanceIds

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{RoleAttachmentUnExpectedJson}) {
				return resource.RetryableError(WrapError(Error("Please trying again.")))
			}
			if IsExceptedErrors(err, []string{InvalidRamRoleNotFound}) {
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		resp, _ := raw.(*ecs.DescribeInstanceRamRoleResponse)
		instRoleSets := resp.InstanceRamRoleSets.InstanceRamRoleSet
		if len(instRoleSets) > 0 {
			var instIds []string
			for _, item := range instRoleSets {
				if item.RamRoleName == roleName {
					instIds = append(instIds, item.InstanceId)
				}
			}
			ids := strings.Split(strings.TrimRight(strings.TrimLeft(strings.Replace(instanceIds, "\"", "", -1), "["), "]"), ",")
			sort.Strings(instIds)
			sort.Strings(ids)
			if reflect.DeepEqual(instIds, ids) {
				d.Set("role_name", resp.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName)
				d.Set("instance_ids", instIds)
				return nil
			}
		}
		return resource.NonRetryableError(WrapError(Error("No ram role for instances found.")))
	})
}

func resourceAlicloudInstanceRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	roleName := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	request := ecs.CreateDetachInstanceRamRoleRequest()
	request.RamRoleName = roleName
	request.InstanceIds = instanceIds

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DetachInstanceRamRole(request)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{RoleAttachmentUnExpectedJson}) {
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		return nil
	})
}
