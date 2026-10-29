package imm

import (
	"fmt"
	"log"
	"regexp"
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

func resourceAlicloudImmProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudImmProjectCreate,
		Read:   resourceAlicloudImmProjectRead,
		Update: resourceAlicloudImmProjectUpdate,
		Delete: resourceAlicloudImmProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-z][a-zA-Z0-9\\-]{1,49}$`), "The length of the project name is limited to 1-50 characters, can only contain English letters, numbers and dashes (-), must start with an English letter, and cannot be ALL (reserved words)"),
			},
			"service_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudImmProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := make(map[string]interface{})
	var err error
	request["Project"] = d.Get("project")
	if v, ok := d.GetOk("service_role"); ok {
		request["ServiceRole"] = v
	}
	action := "PutProject"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		resp, err := client.RpcPost("imm", "2017-09-06", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_imm_project", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Project"]))

	return resourceAlicloudImmProjectUpdate(d, meta)
}
func resourceAlicloudImmProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	immService := NewService(client)
	object, err := immService.DescribeImmProject(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_imm_project immService.DescribeImmProject Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("project", d.Id())
	if v, ok := object["CU"]; ok && fmt.Sprint(v) != "0" {
		compute_unitVal, err := convert.FormatInt(v)
		if err != nil {
			return err
		}
		d.Set("compute_unit", compute_unitVal)
	}
	d.Set("service_role", object["ServiceRole"])

	return nil
}
func resourceAlicloudImmProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error

	update := false
	request := map[string]interface{}{
		"Project": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("service_role") {
		update = true
		if v, ok := d.GetOk("service_role"); ok {
			request["NewServiceRole"] = v
		}
	}
	if update {
		action := "UpdateProject"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			resp, err := client.RpcPost("imm", "2017-09-06", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudImmProjectRead(d, meta)
}
func resourceAlicloudImmProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	request := map[string]interface{}{
		"Project": d.Id(),
	}

	action := "DeleteProject"
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		resp, err := client.RpcPost("imm", "2017-09-06", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, resp, request)
		return nil
	})
	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
