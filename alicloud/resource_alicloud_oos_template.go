package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOosTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosTemplateCreate,
		Read:   resourceAlicloudOosTemplateRead,
		Update: resourceAlicloudOosTemplateUpdate,
		Delete: resourceAlicloudOosTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"auto_delete_executions": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"has_trigger": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"template_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlicloudOosTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateTemplate"
	request := make(map[string]interface{})
	var err error
	request["Content"] = d.Get("content")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	request["TemplateName"] = d.Get("template_name")
	if v, ok := d.GetOk("version_name"); ok {
		request["VersionName"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_template", action, AlibabaCloudSdkGoERROR)
	}
	responseTemplate := response["Template"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseTemplate["TemplateName"]))

	return resourceAlicloudOosTemplateRead(d, meta)
}

func resourceAlicloudOosTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_template oosService.DescribeOosTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("template_name", d.Id())
	d.Set("created_by", object["CreatedBy"])
	d.Set("created_date", object["CreatedDate"])
	d.Set("description", object["Description"])
	d.Set("has_trigger", object["HasTrigger"])
	d.Set("share_type", object["ShareType"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v))
	}
	d.Set("template_format", object["TemplateFormat"])
	d.Set("template_id", object["TemplateId"])
	d.Set("template_type", object["TemplateType"])
	d.Set("template_version", object["TemplateVersion"])
	d.Set("updated_by", object["UpdatedBy"])
	d.Set("updated_date", object["UpdatedDate"])
	d.Set("resource_group_id", object["ResourceGroupId"])

	// $.ResourceIds is exposed as a computed attribute and sourced from
	// ListTagResources ($.TagResources.TagResource[*].ResourceId). The query
	// is scoped to the current template, so de-duplicate the echoed resource
	// ids that ListTagResources returns once per bound tag. Only non-empty
	// string ResourceId values are accepted; missing, null or non-string
	// entries are ignored so that "<nil>" never leaks into state.
	tagResp, err := oosService.ListOosTemplateTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}
	tagResources, _ := jsonpath.Get("$.TagResources.TagResource", tagResp)
	resourceIds := make([]string, 0)
	seen := make(map[string]bool)
	for _, tagResource := range convertToInterfaceArray(tagResources) {
		tagResourceMap, ok := tagResource.(map[string]interface{})
		if !ok {
			continue
		}
		rid, ok := tagResourceMap["ResourceId"].(string)
		if !ok || rid == "" {
			continue
		}
		if seen[rid] {
			continue
		}
		seen[rid] = true
		resourceIds = append(resourceIds, rid)
	}
	d.Set("resource_ids", resourceIds)
	return nil
}

func resourceAlicloudOosTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"TemplateName": d.Id(),
	}
	if d.HasChange("content") {
		update = true
	}
	request["Content"] = d.Get("content")
	request["RegionId"] = client.RegionId
	if d.HasChange("tags") {
		update = true
		respJson, err := convertMaptoJsonString(d.Get("tags").(map[string]interface{}))
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_template", "UpdateTemplate", AlibabaCloudSdkGoERROR)
		}
		request["Tags"] = respJson
	}
	if d.HasChange("version_name") {
		update = true
		request["VersionName"] = d.Get("version_name")
	}
	if d.HasChange("resource_group_id") {
		update = true
		request["ResourceGroupId"] = d.Get("resource_group_id")
	}
	if update {
		action := "UpdateTemplate"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
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
	return resourceAlicloudOosTemplateRead(d, meta)
}

func resourceAlicloudOosTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteTemplate"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"TemplateName": d.Id(),
	}

	if v, ok := d.GetOkExists("auto_delete_executions"); ok {
		request["AutoDeleteExecutions"] = v
	}
	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
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
		if IsExpectedErrors(err, []string{"EntityNotExists.Template"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
