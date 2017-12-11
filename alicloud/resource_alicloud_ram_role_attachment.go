package alicloud

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamRoleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudInstanceRoleAttachmentCreate,
		Read:   resourceAlicloudInstanceRoleAttachmentRead,
		Delete: resourceAlicloudInstanceRoleAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"role_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRamName,
				ForceNew:     true,
			},
			"instance_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudInstanceRoleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn

	instanceIds := convertListToJsonString(d.Get("instance_ids").(*schema.Set).List())

	args := ecs.AttachInstancesArgs{
		RegionId:    getRegion(d, meta),
		InstanceIds: instanceIds,
		RamRoleName: d.Get("role_name").(string),
	}

	err := client.JudgeRolePolicyPrincipal(args.RamRoleName)
	if err != nil {
		return err
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		if err := conn.AttachInstanceRamRole(&args); err != nil {
			if IsExceptedError(err, RoleAttachmentUnExpectedJson) {
				return resource.RetryableError(fmt.Errorf("Please trying again."))
			}
			return resource.NonRetryableError(fmt.Errorf("AttachInstanceRamRole got an error: %#v", err))
		}
		d.SetId(d.Get("role_name").(string) + ":" + instanceIds)
		return resource.NonRetryableError(resourceAlicloudInstanceRoleAttachmentRead(d, meta))
	})
}

func resourceAlicloudInstanceRoleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	roleName := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	args := ecs.AttachInstancesArgs{
		RegionId:    getRegion(d, meta),
		InstanceIds: instanceIds,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err := conn.DescribeInstanceRamRole(&args)
		if err != nil {
			if IsExceptedError(err, RoleAttachmentUnExpectedJson) {
				return resource.RetryableError(fmt.Errorf("Please trying again."))
			}
			if IsExceptedError(err, InvalidRamRoleNotFound) {
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("DescribeInstanceRamRole got an error: %#v", err))
		}

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
		return resource.NonRetryableError(fmt.Errorf("No ram role for instances found."))
	})
}

func resourceAlicloudInstanceRoleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	roleName := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DetachInstanceRamRole(&ecs.AttachInstancesArgs{
			RegionId:    getRegion(d, meta),
			RamRoleName: roleName,
			InstanceIds: instanceIds,
		})

		if err != nil {
			if IsExceptedError(err, RoleAttachmentUnExpectedJson) {
				return resource.RetryableError(fmt.Errorf("Please trying again."))
			}
			return resource.NonRetryableError(fmt.Errorf("Error DetachInstanceRamRole:%#v", err))
		}
		return nil
	})
}
