package fnf

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceFnfFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceFnfFlowCreate,
		Read:   resourceFnfFlowRead,
		Update: resourceFnfFlowUpdate,
		Delete: resourceFnfFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"definition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flow_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "FDL"}, false),
			},
		},
	}
}

func resourceFnfFlowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateFlow"
	request := make(map[string]interface{})
	var err error
	request["Definition"] = d.Get("definition")
	request["Description"] = d.Get("description")
	request["Name"] = d.Get("name")
	if v, ok := d.GetOk("role_arn"); ok {
		request["RoleArn"] = v
	}

	request["Type"] = d.Get("type")
	wait := helper.IncrementalWait(3*time.Second, 1*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("fnf", "2019-03-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_fnf_flow", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Name"]))

	return resourceFnfFlowRead(d, meta)
}
func resourceFnfFlowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fnfService := NewService(client)
	object, err := fnfService.DescribeFnfFlow(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fnf_flow fnfService.DescribeFnfFlow Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("name", d.Id())
	d.Set("definition", object["Definition"])
	d.Set("description", object["Description"])
	d.Set("flow_id", object["Id"])
	d.Set("last_modified_time", object["LastModifiedTime"])
	d.Set("role_arn", object["RoleArn"])
	d.Set("type", object["Type"])
	return nil
}
func resourceFnfFlowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	if d.HasChange("definition") {
		update = true
		request["Definition"] = d.Get("definition")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("role_arn") {
		update = true
		request["RoleArn"] = d.Get("role_arn")
	}
	if d.HasChange("type") {
		update = true
		request["Type"] = d.Get("type")
	}
	if update {
		action := "UpdateFlow"
		wait := helper.IncrementalWait(3*time.Second, 1*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("fnf", "2019-03-15", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || helper.NeedRetry(err) {
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
	return resourceFnfFlowRead(d, meta)
}
func resourceFnfFlowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFlow"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"Name": d.Id(),
	}

	wait := helper.IncrementalWait(3*time.Second, 1*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcGet("fnf", "2019-03-15", action, request, nil)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"InternalServerError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"FlowNotExists"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
