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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourcePidLoop() *schema.Resource {
	return &schema.Resource{
		Create: resourcePidLoopCreate,
		Read:   resourcePidLoopRead,
		Update: resourcePidLoopUpdate,
		Delete: resourcePidLoopDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "This resource has been deprecated since v1.229.1 and will be removed in the future.",
		Schema: map[string]*schema.Schema{
			"pid_loop_configuration": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := helper.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"pid_loop_dcs_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard"}, false),
			},
			"pid_loop_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pid_loop_is_crucial": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"pid_loop_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pid_loop_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"0", "1", "2", "3", "4", "5"}, false),
			},
			"pid_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePidLoopCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePidLoop"
	request := make(map[string]interface{})
	var err error
	request["PidLoopConfiguration"] = d.Get("pid_loop_configuration")
	request["PidLoopDcsType"] = d.Get("pid_loop_dcs_type")
	if v, ok := d.GetOk("pid_loop_desc"); ok {
		request["PidLoopDesc"] = v
	}

	request["IsCrucialPidLoop"] = d.Get("pid_loop_is_crucial")
	request["PidLoopName"] = d.Get("pid_loop_name")
	request["PidLoopType"] = d.Get("pid_loop_type")
	request["PidProjectId"] = d.Get("pid_project_id")
	request["ClientToken"] = helper.BuildClientToken("CreatePidLoop")
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_brain_industrial_pid_loop", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	d.SetId(fmt.Sprint(response["PidLoopId"]))

	return resourcePidLoopRead(d, meta)
}
func resourcePidLoopRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	service := NewService(client)
	object, err := service.DescribePidLoop(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_brain_industrial_pid_loop brainindustrialService.DescribePidLoop Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("pid_loop_configuration", object["PidLoopConfiguration"])
	d.Set("pid_loop_dcs_type", object["PidLoopDcsType"])
	d.Set("pid_loop_desc", object["PidLoopDesc"])
	d.Set("pid_loop_is_crucial", object["PidLoopIsCrucial"])
	d.Set("pid_loop_name", object["PidLoopName"])
	d.Set("pid_loop_type", object["PidLoopType"])
	d.Set("pid_project_id", object["PidProjectId"])
	d.Set("status", object["Status"])
	return nil
}
func resourcePidLoopUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"PidLoopId": d.Id(),
	}
	if d.HasChange("pid_loop_configuration") {
		update = true
	}
	request["PidLoopConfiguration"] = d.Get("pid_loop_configuration")
	if d.HasChange("pid_loop_dcs_type") {
		update = true
	}
	request["PidLoopDcsType"] = d.Get("pid_loop_dcs_type")
	if d.HasChange("pid_loop_is_crucial") {
		update = true
	}
	request["Crucial"] = d.Get("pid_loop_is_crucial")
	if d.HasChange("pid_loop_name") {
		update = true
	}
	request["PidLoopName"] = d.Get("pid_loop_name")
	if d.HasChange("pid_loop_type") {
		update = true
	}
	request["PidLoopType"] = d.Get("pid_loop_type")
	if d.HasChange("pid_project_id") {
		update = true
	}
	request["PidProjectId"] = d.Get("pid_project_id")
	if d.HasChange("pid_loop_desc") {
		update = true
		request["PidLoopDesc"] = d.Get("pid_loop_desc")
	}
	if update {
		action := "UpdatePidLoop"
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
	return resourcePidLoopRead(d, meta)
}
func resourcePidLoopDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePidLoop"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PidLoopId": d.Id(),
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
