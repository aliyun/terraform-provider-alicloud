// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudActiontrailAdvancedQueryTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudActiontrailAdvancedQueryTemplateCreate,
		Read:   resourceAliCloudActiontrailAdvancedQueryTemplateRead,
		Update: resourceAliCloudActiontrailAdvancedQueryTemplateUpdate,
		Delete: resourceAliCloudActiontrailAdvancedQueryTemplateDelete,
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

func resourceAliCloudActiontrailAdvancedQueryTemplateCreate(d *schema.ResourceData, meta interface{}) error {

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
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_actiontrail_advanced_query_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["TemplateId"]))

	return resourceAliCloudActiontrailAdvancedQueryTemplateRead(d, meta)
}

func resourceAliCloudActiontrailAdvancedQueryTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailServiceV2 := ActiontrailServiceV2{client}

	objectRaw, err := actiontrailServiceV2.DescribeActiontrailAdvancedQueryTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_advanced_query_template DescribeActiontrailAdvancedQueryTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("simple_query", objectRaw["SimpleQuery"])
	d.Set("template_name", objectRaw["TemplateName"])
	d.Set("template_sql", objectRaw["TemplateSql"])

	return nil
}

func resourceAliCloudActiontrailAdvancedQueryTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudActiontrailAdvancedQueryTemplateRead(d, meta)
}

func resourceAliCloudActiontrailAdvancedQueryTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAdvancedQueryTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["TemplateId"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Actiontrail", "2020-07-06", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
