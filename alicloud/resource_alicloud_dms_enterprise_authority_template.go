// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDMSEnterpriseAuthorityTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDMSEnterpriseAuthorityTemplateCreate,
		Read:   resourceAliCloudDMSEnterpriseAuthorityTemplateRead,
		Update: resourceAliCloudDMSEnterpriseAuthorityTemplateUpdate,
		Delete: resourceAliCloudDMSEnterpriseAuthorityTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"authority_template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"authority_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudDMSEnterpriseAuthorityTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAuthorityTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Tid"] = d.Get("tid")

	request["Name"] = d.Get("authority_template_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dms_enterprise_authority_template", action, AlibabaCloudSdkGoERROR)
	}

	templateId, _ := jsonpath.Get("$.AuthorityTemplateView.TemplateId", response)
	d.SetId(fmt.Sprintf("%v:%v", query["Tid"], templateId))

	return resourceAliCloudDMSEnterpriseAuthorityTemplateRead(d, meta)
}

func resourceAliCloudDMSEnterpriseAuthorityTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dMSEnterpriseServiceV2 := DMSEnterpriseServiceV2{client}

	objectRaw, err := dMSEnterpriseServiceV2.DescribeDMSEnterpriseAuthorityTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dms_enterprise_authority_template DescribeDMSEnterpriseAuthorityTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	authorityTemplateView1RawObj, _ := jsonpath.Get("$.AuthorityTemplateView", objectRaw)
	authorityTemplateView1Raw := make(map[string]interface{})
	if authorityTemplateView1RawObj != nil {
		authorityTemplateView1Raw = authorityTemplateView1RawObj.(map[string]interface{})
	}
	d.Set("authority_template_name", authorityTemplateView1Raw["Name"])
	d.Set("create_time", authorityTemplateView1Raw["CreateTime"])
	d.Set("description", authorityTemplateView1Raw["Description"])

	parts := strings.Split(d.Id(), ":")
	d.Set("tid", formatInt(parts[0]))
	d.Set("authority_template_id", formatInt(parts[1]))
	return nil
}

func resourceAliCloudDMSEnterpriseAuthorityTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	action := "UpdateAuthorityTemplate"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TemplateId"] = parts[1]
	query["Tid"] = parts[0]
	if d.HasChange("authority_template_name") {
		update = true
	}
	request["Name"] = d.Get("authority_template_name")
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudDMSEnterpriseAuthorityTemplateRead(d, meta)
}

func resourceAliCloudDMSEnterpriseAuthorityTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAuthorityTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["Tid"] = parts[0]
	query["TemplateId"] = parts[1]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("dms-enterprise", "2018-11-01", action, query, request, true)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
