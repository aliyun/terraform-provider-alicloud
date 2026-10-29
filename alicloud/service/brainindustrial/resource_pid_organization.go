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

func ResourcePidOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourcePidOrganizationCreate,
		Read:   resourcePidOrganizationRead,
		Update: resourcePidOrganizationUpdate,
		Delete: resourcePidOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.222.0 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"parent_pid_organization_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"pid_organization_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePidOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePidOrganization"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("parent_pid_organization_id"); ok {
		request["ParentOrganizationId"] = v
	}

	request["OrganizationName"] = d.Get("pid_organization_name")
	request["ClientToken"] = helper.BuildClientToken("CreatePidOrganization")
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_brain_industrial_pid_organization", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	d.SetId(fmt.Sprint(response["OrganizationId"]))

	return resourcePidOrganizationRead(d, meta)
}
func resourcePidOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	object, err := service.DescribePidOrganization(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_brain_industrial_pid_organization brainindustrialService.DescribePidOrganization Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("pid_organization_name", object["OrganizationName"])
	return nil
}
func resourcePidOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	if d.HasChange("pid_organization_name") {
		request := map[string]interface{}{
			"OrganizationId": d.Id(),
		}
		request["OrganizationName"] = d.Get("pid_organization_name")
		action := "UpdatePidOrganization"
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
	return resourcePidOrganizationRead(d, meta)
}
func resourcePidOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePidOrganization"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"OrganizationId": d.Id(),
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
