package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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
		// UpdateTemplate replaces the whole tag set when Tags is present, so the ignored tags
		// already on the template have to be sent back when ignore_tags is configured; without
		// the configuration the configured tags are sent as-is.
		tagsPayload := d.Get("tags").(map[string]interface{})
		if !client.IgnoreTags.Empty() {
			oosService := OosService{client}
			object, err := oosService.DescribeOosTemplate(d.Id())
			if err != nil {
				return WrapError(err)
			}
			current := make(map[string]string)
			for key, value := range tagsToMap(object["Tags"]) {
				current[key] = value.(string)
			}
			desired := make(map[string]string)
			for key, value := range tagsPayload {
				desired[key] = value.(string)
			}
			merged := mergeIgnoredTags(desired, current, client.IgnoreTags)
			tagsPayload = make(map[string]interface{}, len(merged))
			for key, value := range merged {
				tagsPayload[key] = value
			}
		}
		respJson, err := convertMaptoJsonString(tagsPayload)
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
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// GetTemplate lags behind UpdateTemplate; wait until the template reports the
		// exact tag pairs this update sent so the trailing Read does not observe the
		// previous set.
		if d.HasChange("tags") {
			if tagsJson, ok := request["Tags"].(string); ok {
				parsed := map[string]interface{}{}
				if json.Unmarshal([]byte(tagsJson), &parsed) == nil && len(parsed) > 0 {
					err = retry.Retry(2*time.Minute, func() *retry.RetryError {
						object, err := oosServiceDescribeTemplateTags(client, d.Id())
						if err != nil {
							return retry.NonRetryableError(err)
						}
						for key, value := range parsed {
							if object[key] != fmt.Sprint(value) {
								return retry.RetryableError(fmt.Errorf("template tags not yet updated"))
							}
						}
						return nil
					})
					if err != nil {
						return WrapErrorf(err, "waiting for oos template %s tags to converge", d.Id())
					}
				}
			}
		}
	}
	return resourceAlicloudOosTemplateRead(d, meta)
}

func oosServiceDescribeTemplateTags(client *connectivity.AliyunClient, name string) (map[string]string, error) {
	object, err := (&OosService{client}).DescribeOosTemplate(name)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		for key, value := range v {
			result[key] = fmt.Sprint(value)
		}
	}
	return result, nil
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
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("oos", "2019-06-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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
