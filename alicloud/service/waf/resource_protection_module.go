package waf

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceProtectionModule() *schema.Resource {
	return &schema.Resource{
		Create: resourceProtectionModuleCreate,
		Read:   resourceProtectionModuleRead,
		Update: resourceProtectionModuleUpdate,
		Delete: resourceProtectionModuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"defense_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ac_cc", "antifraud", "dld", "normalized", "waf"}, false),
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func resourceProtectionModuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ModifyProtectionModuleMode"
	request := make(map[string]interface{})
	var err error
	request["DefenseType"] = d.Get("defense_type")
	request["Domain"] = d.Get("domain")
	request["InstanceId"] = d.Get("instance_id")
	request["Mode"] = d.Get("mode")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {

		response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_waf_protection_module", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["InstanceId"], ":", request["Domain"], ":", request["DefenseType"]))

	return resourceProtectionModuleUpdate(d, meta)
}
func resourceProtectionModuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafOpenapiService := NewService(client)
	object, err := wafOpenapiService.DescribeWafProtectionModule(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_waf_protection_module wafOpenapiService.DescribeWafProtectionModule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("defense_type", parts[2])
	d.Set("domain", parts[1])
	d.Set("instance_id", parts[0])
	modeVal, err := convert.FormatInt(object["Mode"])
	if err != nil {
		return err
	}
	d.Set("mode", modeVal)
	describeProtectionModuleStatusObject, err := wafOpenapiService.DescribeProtectionModuleStatus(d.Id())
	if err != nil {
		return errs.WrapError(err)
	}
	statusVal, err := convert.FormatInt(describeProtectionModuleStatusObject["ModuleStatus"])
	if err != nil {
		return err
	}
	d.Set("status", statusVal)
	return nil
}
func resourceProtectionModuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 3)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Partial(true)

	request := map[string]interface{}{
		"DefenseType": parts[2],
		"Domain":      parts[1],
		"InstanceId":  parts[0],
	}

	update := false
	if !d.IsNewResource() && d.HasChange("mode") {
		update = true
		if v, ok := d.GetOkExists("mode"); ok {
			request["Mode"] = v
		}
	}

	if update {
		action := "ModifyProtectionModuleMode"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
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
	}

	update = false
	if d.IsNewResource() || d.HasChange("status") {
		update = true
		if v, ok := d.GetOkExists("status"); ok {
			request["ModuleStatus"] = v
		}
	}

	if update {
		action := "ModifyProtectionModuleStatus"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("waf-openapi", "2019-09-10", action, nil, request, false)
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
	}

	d.Partial(false)
	return resourceProtectionModuleRead(d, meta)
}
func resourceProtectionModuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy ResourceProtectionModule. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
