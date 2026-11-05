// Package actiontrail. This file is generated automatically. Please do not modify it manually, thank you!
package actiontrail

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

func ResourceAdvancedQueryTemplate() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAdvancedQueryTemplateCreate,
		Read:   ResourceAdvancedQueryTemplateRead,
		Update: ResourceAdvancedQueryTemplateUpdate,
		Delete: ResourceAdvancedQueryTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"simple_query": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_sql": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func ResourceAdvancedQueryTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAdvancedQueryTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["TemplateSql"] = d.Get("template_sql")
	request["SimpleQuery"] = d.Get("simple_query")
	if v, ok := d.GetOk("template_name"); ok {
		request["TemplateName"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_actiontrail_advanced_query_template", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return ResourceAdvancedQueryTemplateRead(d, meta)
}

func ResourceAdvancedQueryTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailServiceV2 := NewActiontrailServiceV2(client)

	objectRaw, err := actiontrailServiceV2.DescribeActiontrailAdvancedQueryTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_advanced_query_template DescribeActiontrailAdvancedQueryTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("simple_query", objectRaw["SimpleQuery"])
	d.Set("template_name", objectRaw["TemplateName"])
	d.Set("template_sql", objectRaw["TemplateSql"])

	return nil
}

func ResourceAdvancedQueryTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	action := "UpdateAdvancedQueryTemplate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["TemplateId"] = d.Id()

	if d.HasChange("template_sql") {
		update = true
	}
	request["TemplateSql"] = d.Get("template_sql")
	if d.HasChange("simple_query") {
		update = true
	}
	request["SimpleQuery"] = d.Get("simple_query")
	if d.HasChange("template_name") {
		update = true
		request["TemplateName"] = d.Get("template_name")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
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

	return ResourceAdvancedQueryTemplateRead(d, meta)
}

func ResourceAdvancedQueryTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAdvancedQueryTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TemplateId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)

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
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
