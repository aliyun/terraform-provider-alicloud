// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudWafv3DefenseTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudWafv3DefenseTemplateCreate,
		Read:   resourceAliCloudWafv3DefenseTemplateRead,
		Update: resourceAliCloudWafv3DefenseTemplateUpdate,
		Delete: resourceAliCloudWafv3DefenseTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"defense_scene": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"defense_template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"defense_template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_manager_resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_origin": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudWafv3DefenseTemplateCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDefenseTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = d.Get("instance_id")
	request["RegionId"] = client.RegionId

	request["DefenseScene"] = d.Get("defense_scene")
	request["TemplateOrigin"] = d.Get("template_origin")
	request["TemplateType"] = d.Get("template_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["TemplateName"] = d.Get("defense_template_name")
	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	request["TemplateStatus"] = d.Get("status")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, false)

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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_wafv3_defense_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["InstanceId"], response["TemplateId"]))

	return resourceAliCloudWafv3DefenseTemplateRead(d, meta)
}

func resourceAliCloudWafv3DefenseTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	wafv3ServiceV2 := Wafv3ServiceV2{client}

	objectRaw, err := wafv3ServiceV2.DescribeWafv3DefenseTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_wafv3_defense_template DescribeWafv3DefenseTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("defense_scene", objectRaw["DefenseScene"])
	d.Set("defense_template_name", objectRaw["TemplateName"])
	d.Set("description", objectRaw["Description"])
	d.Set("status", objectRaw["TemplateStatus"])
	d.Set("template_origin", objectRaw["TemplateOrigin"])
	d.Set("template_type", objectRaw["TemplateType"])
	d.Set("defense_template_id", objectRaw["TemplateId"])

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])
	d.Set("defense_template_id", parts[1])

	return nil
}

func resourceAliCloudWafv3DefenseTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDefenseTemplate"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["TemplateId"] = parts[1]
	query["InstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("defense_template_name") {
		update = true
	}
	request["TemplateName"] = d.Get("defense_template_name")
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, false)

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
		d.SetPartial("defense_template_name")
		d.SetPartial("description")
	}
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyDefenseTemplateStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["InstanceId"] = parts[0]
	query["TemplateId"] = parts[1]
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	if d.HasChange("status") {
		update = true
	}
	request["TemplateStatus"] = d.Get("status")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, false)

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
		d.SetPartial("status")
	}

	d.Partial(false)
	return resourceAliCloudWafv3DefenseTemplateRead(d, meta)
}

func resourceAliCloudWafv3DefenseTemplateDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDefenseTemplate"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	query["InstanceId"] = parts[0]
	query["TemplateId"] = parts[1]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, false)

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
