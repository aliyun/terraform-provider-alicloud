package brainindustrial

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourcePidProject() *schema.Resource {
	return &schema.Resource{
		Create: resourcePidProjectCreate,
		Read:   resourcePidProjectRead,
		Update: resourcePidProjectUpdate,
		Delete: resourcePidProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.222.0 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"pid_organization_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pid_project_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pid_project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePidProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePidProject"
	request := make(map[string]interface{})
	var err error
	request["PidOrganisationId"] = d.Get("pid_organization_id")
	if v, ok := d.GetOk("pid_project_desc"); ok {
		request["PidProjectDesc"] = v
	}

	request["PidProjectName"] = d.Get("pid_project_name")
	request["ClientToken"] = helper.BuildClientToken("CreatePidProject")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_brain_industrial_pid_project", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	d.SetId(fmt.Sprint(response["PidProjectId"]))

	return resourcePidProjectRead(d, meta)
}
func resourcePidProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	object, err := service.DescribePidProject(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_brain_industrial_pid_project brainindustrialService.DescribePidProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("pid_organization_id", object["PidOrganisationId"])
	d.Set("pid_project_desc", object["PidProjectDesc"])
	d.Set("pid_project_name", object["PidProjectName"])
	return nil
}
func resourcePidProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PidProjectId": d.Id(),
	}
	if d.HasChange("pid_organization_id") {
		update = true
	}
	request["PidOrganisationId"] = d.Get("pid_organization_id")
	if d.HasChange("pid_project_desc") {
		update = true
		request["PidProjectDesc"] = d.Get("pid_project_desc")
	}
	if d.HasChange("pid_project_name") {
		update = true
		request["PidProjectName"] = d.Get("pid_project_name")
	}
	if update {
		action := "UpdatePidProject"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourcePidProjectRead(d, meta)
}
func resourcePidProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePidProject"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PidProjectId": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("brain-industrial", "2020-09-20", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
