package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudResourceManagerServiceLinkedRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudResourceManagerServiceLinkedRoleCreate,
		Read:   resourceAliCloudResourceManagerServiceLinkedRoleRead,
		Delete: resourceAliCloudResourceManagerServiceLinkedRoleDelete,
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

func resourceAliCloudResourceManagerServiceLinkedRoleCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(d.Get("service_name"), ":", response["Role"].(map[string]interface{})["RoleName"]))

	return resourceAliCloudResourceManagerServiceLinkedRoleRead(d, meta)
}
func resourceAliCloudResourceManagerServiceLinkedRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	parts, _ := ParseResourceId(d.Id(), 2)

	object, err := ramService.DescribeRamServiceLinkedRole(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("role_id", object.Role.RoleId)
	d.Set("role_name", object.Role.RoleName)
	d.Set("arn", object.Role.Arn)
	d.Set("service_name", parts[0])
	return nil

}

func resourceAliCloudResourceManagerServiceLinkedRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	resourceManagerService := ResourceManagerService{client}
	var response map[string]interface{}
	action := "DeleteServiceLinkedRole"
	request := make(map[string]interface{})
	var err error
	parts, _ := ParseResourceId(d.Id(), 2)
	request["RoleName"] = parts[1]

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("ResourceManager", "2020-03-31", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_service_linked_role", action, AlibabaCloudSdkGoERROR)
	}

	taskId := fmt.Sprint(response["DeletionTaskId"])
	stateConf := BuildStateConf([]string{}, []string{"SUCCEEDED"}, d.Timeout(schema.TimeoutDelete), 0*time.Second, resourceManagerService.ResourceManagerServiceLinkedRoleStateRefreshFunc(taskId, []string{"FAILED", "INTERNAL_ERROR"}))
	if _, err := stateConf.WaitForState(); err != nil {
		object, e := resourceManagerService.GetServiceLinkedRoleDeletionStatus(taskId)
		if e != nil {
			return WrapErrorf(err, FailedToReachTargetStatusWithError, d.Id(), e)
		}
		return WrapErrorf(err, FailedToReachTargetStatusWithResponse, d.Id(), object)
	}
	return nil
}
