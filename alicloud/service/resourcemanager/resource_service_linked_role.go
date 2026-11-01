package resourcemanager

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: ResourceServiceLinkedRoleCreate,
		Read:   ResourceServiceLinkedRoleRead,
		Delete: ResourceServiceLinkedRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"custom_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateServiceLinkedRole"
	request := make(map[string]interface{})
	var err error

	request["ServiceName"] = d.Get("service_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("custom_suffix"); ok {
		request["CustomSuffix"] = v
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ram_service_linked_role", action, errs.AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(d.Get("service_name"), ":", response["Role"].(map[string]interface{})["RoleName"]))

	return ResourceServiceLinkedRoleRead(d, meta)
}
func ResourceServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := NewResourceManagerService(client)
	parts, _ := helper.ParseResourceId(d.Id(), 2)

	object, err := resourceManagerService.DescribeRamServiceLinkedRole(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("role_id", object.Role.RoleId)
	d.Set("role_name", object.Role.RoleName)
	d.Set("arn", object.Role.Arn)
	d.Set("service_name", parts[0])
	return nil

}

func ResourceServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourceManagerService{client}
	var response map[string]interface{}
	action := "DeleteServiceLinkedRole"
	request := make(map[string]interface{})
	var err error
	parts, _ := helper.ParseResourceId(d.Id(), 2)
	request["RoleName"] = parts[1]

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_ram_service_linked_role", action, errs.AlibabaCloudSdkGoERROR)
	}

	taskId := fmt.Sprint(response["DeletionTaskId"])
	stateConf := helper.BuildStateConf([]string{}, []string{"SUCCEEDED"}, d.Timeout(schema.TimeoutDelete), 0*time.Second, resourceManagerService.ResourceManagerServiceLinkedRoleStateRefreshFunc(taskId, []string{"FAILED", "INTERNAL_ERROR"}))
	if _, err := stateConf.WaitForState(); err != nil {
		object, e := resourceManagerService.GetServiceLinkedRoleDeletionStatus(taskId)
		if e != nil {
			return errs.WrapErrorf(err, errs.FailedToReachTargetStatusWithError, d.Id(), e)
		}
		return errs.WrapErrorf(err, errs.FailedToReachTargetStatusWithResponse, d.Id(), object)
	}
	return nil
}
