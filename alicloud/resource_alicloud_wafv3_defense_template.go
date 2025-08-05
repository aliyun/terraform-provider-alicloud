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
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"user_default", "user_custom"}, false),
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
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["TemplateName"] = d.Get("defense_template_name")
	request["TemplateType"] = d.Get("template_type")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["DefenseScene"] = d.Get("defense_scene")
	request["TemplateStatus"] = d.Get("status")
	request["TemplateOrigin"] = d.Get("template_origin")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_wafv3_defense_template", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], response["TemplateId"]))

	return resourceAliCloudWafv3DefenseTemplateUpdate(d, meta)
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

	objectRaw, err = wafv3ServiceV2.DescribeDefenseTemplateDescribeTemplateResources(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("defense_template_id", objectRaw["TemplateId"])

	resourcesRaw := make([]interface{}, 0)
	if objectRaw["Resources"] != nil {
		resourcesRaw = objectRaw["Resources"].([]interface{})
	}

	d.Set("resources", resourcesRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceAliCloudWafv3DefenseTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDefenseTemplate"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["TemplateId"] = parts[1]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("defense_template_name") {
		update = true
	}
	request["TemplateName"] = d.Get("defense_template_name")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}

	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyDefenseTemplateStatus"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["TemplateId"] = parts[1]
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("status") {
		update = true
	}
	request["TemplateStatus"] = d.Get("status")
	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

	if d.HasChange("resource_groups") {
		oldEntry, newEntry := d.GetChange("resource_groups")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "ModifyTemplateResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = parts[0]
			request["TemplateId"] = parts[1]
			request["RegionId"] = client.RegionId
			localData := removed.List()
			unbindResourceGroupsMapsArray := localData
			request["UnbindResourceGroups"] = unbindResourceGroupsMapsArray

			if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
				request["ResourceManagerResourceGroupId"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "ModifyTemplateResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = parts[0]
			request["TemplateId"] = parts[1]
			request["RegionId"] = client.RegionId
			localData := added.List()
			bindResourceGroupsMapsArray := localData
			request["BindResourceGroups"] = bindResourceGroupsMapsArray

			if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
				request["ResourceManagerResourceGroupId"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

	}
	if d.HasChange("resources") {
		oldEntry, newEntry := d.GetChange("resources")
		oldEntrySet := oldEntry.(*schema.Set)
		newEntrySet := newEntry.(*schema.Set)
		removed := oldEntrySet.Difference(newEntrySet)
		added := newEntrySet.Difference(oldEntrySet)

		if removed.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "ModifyTemplateResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = parts[0]
			request["TemplateId"] = parts[1]
			request["RegionId"] = client.RegionId
			localData := removed.List()
			unbindResourcesMapsArray := localData
			request["UnbindResources"] = unbindResourcesMapsArray

			if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
				request["ResourceManagerResourceGroupId"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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

		if added.Len() > 0 {
			parts := strings.Split(d.Id(), ":")
			action := "ModifyTemplateResources"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["InstanceId"] = parts[0]
			request["TemplateId"] = parts[1]
			request["RegionId"] = client.RegionId
			localData := added.List()
			bindResourcesMapsArray := localData
			request["BindResources"] = bindResourcesMapsArray

			if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
				request["ResourceManagerResourceGroupId"] = v
			}
			wait := incrementalWait(3*time.Second, 5*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)
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
	request["InstanceId"] = parts[0]
	request["TemplateId"] = parts[1]
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("resource_manager_resource_group_id"); ok {
		request["ResourceManagerResourceGroupId"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("waf-openapi", "2021-10-01", action, query, request, true)

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
		if IsExpectedErrors(err, []string{"Defense.Control.DefenseTemplateNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
