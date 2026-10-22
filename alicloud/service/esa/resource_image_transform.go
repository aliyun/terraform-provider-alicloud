// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package esa

import (
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	tfvalidators "github.com/aliyun/terraform-provider-alicloud/alicloud/validators"

	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceImageTransform() *schema.Resource {
	return &schema.Resource{
		Create: resourceImageTransformCreate,
		Read:   resourceImageTransformRead,
		Update: resourceImageTransformUpdate,
		Delete: resourceImageTransformDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auto_avif": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tfvalidators.StringInSlice([]string{"on", "off"}, false),
			},
			"auto_webp": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tfvalidators.StringInSlice([]string{"on", "off"}, false),
			},
			"config_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"config_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tfvalidators.StringInSlice([]string{"on", "off"}, false),
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_enable": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tfvalidators.StringInSlice([]string{"on", "off"}, false),
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceImageTransformCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateImageTransform"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("site_id"); ok {
		request["SiteId"] = v
	}

	if v, ok := d.GetOk("auto_avif"); ok {
		request["AutoAvif"] = v
	}
	if v, ok := d.GetOk("auto_webp"); ok {
		request["AutoWebp"] = v
	}
	if v, ok := d.GetOk("enable"); ok {
		request["Enable"] = v
	}
	if v, ok := d.GetOk("rule"); ok {
		request["Rule"] = v
	}
	if v, ok := d.GetOk("rule_enable"); ok {
		request["RuleEnable"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	if v, ok := d.GetOkExists("sequence"); ok {
		request["Sequence"] = v
	}
	if v, ok := d.GetOkExists("site_version"); ok {
		request["SiteVersion"] = v
	}
	wait := tfhelper.IncrementalWait(5*time.Second, 5*time.Second)
	err = tfretry.Retry(d.Timeout(schema.TimeoutCreate), func() *tfretry.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"LockFailed", "Site.ServiceBusy", "TooManyRequests"}) || tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_esa_image_transform", action, tferrs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["SiteId"], response["ConfigId"]))

	return resourceImageTransformUpdate(d, meta)
}

func resourceImageTransformRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	esaServiceV2 := EsaServiceV2{client}

	objectRaw, err := esaServiceV2.DescribeEsaImageTransform(d.Id())
	if err != nil {
		if !d.IsNewResource() && tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_esa_image_transform DescribeEsaImageTransform Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return tferrs.WrapError(err)
	}

	d.Set("auto_avif", objectRaw["AutoAvif"])
	d.Set("auto_webp", objectRaw["AutoWebp"])
	d.Set("config_type", objectRaw["ConfigType"])
	d.Set("enable", objectRaw["Enable"])
	d.Set("rule", objectRaw["Rule"])
	d.Set("rule_enable", objectRaw["RuleEnable"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("sequence", objectRaw["Sequence"])
	d.Set("site_version", objectRaw["SiteVersion"])
	d.Set("config_id", objectRaw["ConfigId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("site_id", fmt.Sprint(parts[0]))

	return nil
}

func resourceImageTransformUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "UpdateImageTransform"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("auto_avif") {
		update = true
		request["AutoAvif"] = d.Get("auto_avif")
	}

	if !d.IsNewResource() && d.HasChange("auto_webp") {
		update = true
		request["AutoWebp"] = d.Get("auto_webp")
	}

	if !d.IsNewResource() && d.HasChange("enable") {
		update = true
		request["Enable"] = d.Get("enable")
	}

	if !d.IsNewResource() && d.HasChange("rule") {
		update = true
		request["Rule"] = d.Get("rule")
	}

	if !d.IsNewResource() && d.HasChange("rule_enable") {
		update = true
		request["RuleEnable"] = d.Get("rule_enable")
	}

	if !d.IsNewResource() && d.HasChange("rule_name") {
		update = true
		request["RuleName"] = d.Get("rule_name")
	}

	if !d.IsNewResource() && d.HasChange("sequence") {
		update = true
		request["Sequence"] = d.Get("sequence")
	}

	if update {
		wait := tfhelper.IncrementalWait(5*time.Second, 5*time.Second)
		err = tfretry.Retry(d.Timeout(schema.TimeoutUpdate), func() *tfretry.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
			if err != nil {
				if tferrs.IsExpectedErrors(err, []string{"LockFailed", "Site.ServiceBusy", "TooManyRequests"}) || tfhelper.NeedRetry(err) {
					wait()
					return tfretry.RetryableError(err)
				}
				return tfretry.NonRetryableError(err)
			}
			return nil
		})
		tflogs.AddDebug(action, response, request)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceImageTransformRead(d, meta)
}

func resourceImageTransformDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteImageTransform"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["ConfigId"] = parts[1]
	request["SiteId"] = parts[0]

	wait := tfhelper.IncrementalWait(5*time.Second, 5*time.Second)
	err = tfretry.Retry(d.Timeout(schema.TimeoutDelete), func() *tfretry.RetryError {
		response, err = client.RpcPost("ESA", "2024-09-10", action, query, request, true)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		if tferrs.NotFoundError(err) {
			return nil
		}
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
