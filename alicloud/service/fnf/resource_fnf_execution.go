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

func ResourceFnFExecution() *schema.Resource {
	return &schema.Resource{
		Create: resourceFnFExecutionCreate,
		Read:   resourceFnFExecutionRead,
		Update: resourceFnFExecutionUpdate,
		Delete: resourceFnFExecutionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"execution_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flow_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"input": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Stopped"}, false),
			},
		},
	}
}

func resourceFnFExecutionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "StartExecution"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("execution_name"); ok {
		request["ExecutionName"] = v
	}
	request["FlowName"] = d.Get("flow_name")
	if v, ok := d.GetOk("input"); ok {
		request["Input"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("fnf", "2019-03-15", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_fnf_execution", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FlowName"], ":", request["ExecutionName"]))
	fnfService := NewService(client)
	stateConf := helper.BuildStateConf([]string{}, []string{"Running", "Succeeded"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, fnfService.FnFExecutionStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceFnFExecutionUpdate(d, meta)
}
func resourceFnFExecutionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fnfService := NewService(client)
	object, err := fnfService.DescribeFnFExecution(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fnf_execution fnfService.DescribeFnFExecution Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("execution_name", parts[1])
	d.Set("flow_name", parts[0])
	d.Set("input", object["Input"])
	d.Set("status", object["Status"])
	return nil
}
func resourceFnFExecutionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fnfService := NewService(client)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)

	if d.HasChange("status") {
		object, err := fnfService.DescribeFnFExecution(d.Id())
		if err != nil {
			return errs.WrapError(err)
		}
		target := d.Get("status").(string)
		if object["Status"].(string) != target {
			if target == "Stopped" {
				request := map[string]interface{}{
					"ExecutionName": parts[1],
					"FlowName":      parts[0],
				}
				action := "StopExecution"
				wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
				err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					response, err = client.RpcPost("fnf", "2019-03-15", action, nil, request, false)
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
					return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
				}
				stateConf := helper.BuildStateConf([]string{}, []string{"Stopped"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, fnfService.FnFExecutionStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return errs.WrapErrorf(err, errs.IdMsg, d.Id())
				}
			}
		}
	}
	d.Partial(false)
	return resourceFnFExecutionRead(d, meta)
}
func resourceFnFExecutionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceFnFExecution. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
