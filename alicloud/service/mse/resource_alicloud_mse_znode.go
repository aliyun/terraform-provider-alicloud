package mse

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceMseZnode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMseZnodeCreate,
		Read:   resourceAlicloudMseZnodeRead,
		Update: resourceAlicloudMseZnodeUpdate,
		Delete: resourceAlicloudMseZnodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"accept_language": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"zh", "en"}, false),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^\/.+$`), "The value must start with a forward slash (/)."),
			},
		},
	}
}

func resourceAlicloudMseZnodeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateZnode"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	request["ClusterId"] = d.Get("cluster_id")
	if v, ok := d.GetOk("data"); ok {
		request["Data"] = v
	}
	request["Path"] = d.Get("path")
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_mse_znode", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["ClusterId"], ":", request["Path"]))

	return resourceAlicloudMseZnodeRead(d, meta)
}
func resourceAlicloudMseZnodeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mseService := NewService(client)
	object, err := mseService.DescribeMseZnode(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mse_znode mseService.DescribeMseZnode Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("cluster_id", parts[0])
	d.Set("path", object["Path"])
	d.Set("data", object["Data"])
	return nil
}
func resourceAlicloudMseZnodeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	request := map[string]interface{}{
		"ClusterId": parts[0],
		"Path":      parts[1],
	}
	if d.HasChange("data") {
		if v, ok := d.GetOk("data"); ok {
			request["Data"] = v
		}
	}
	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	action := "UpdateZnode"

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return resourceAlicloudMseZnodeRead(d, meta)
}
func resourceAlicloudMseZnodeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteZnode"
	var response map[string]interface{}
	request := map[string]interface{}{
		"ClusterId": parts[0],
		"Path":      parts[1],
	}

	if v, ok := d.GetOk("accept_language"); ok {
		request["AcceptLanguage"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("mse", "2019-05-31", action, nil, request, false)
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
	return nil
}
