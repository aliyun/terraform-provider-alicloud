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

func resourceAliCloudCmsAlertNotifyTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAlertNotifyTemplateCreate,
		Read:   resourceAliCloudCmsAlertNotifyTemplateRead,
		Update: resourceAliCloudCmsAlertNotifyTemplateUpdate,
		Delete: resourceAliCloudCmsAlertNotifyTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_notify_template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alert_notify_template_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Type:     schema.TypeString,
							Required: true,
						},
						"title": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"program_lang": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsAlertNotifyTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/alertNotifyTemplate"
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	query["regionId"] = StringPointer(client.RegionId)

	body["alertNotifyTemplateId"] = d.Get("alert_notify_template_id").(string)
	if v, ok := d.GetOk("alert_notify_template_name"); ok {
		body["alertNotifyTemplateName"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		body["type"] = v
	}
	if v, ok := d.GetOk("program_lang"); ok {
		body["programLang"] = v
	}
	if v, ok := d.GetOk("templates"); ok {
		body["templates"] = buildAlertNotifyTemplateRequestMap(v)
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, body)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_notify_template", action, AlibabaCloudSdkGoERROR)
	}

	if id, e := jsonpath.Get("$.alertNotifyTemplateId", response); e == nil && id != nil && id != "" {
		d.SetId(fmt.Sprintf("%v", id))
	} else {
		d.SetId(d.Get("alert_notify_template_id").(string))
	}

	return resourceAliCloudCmsAlertNotifyTemplateRead(d, meta)
}

func resourceAliCloudCmsAlertNotifyTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAlertNotifyTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_alert_notify_template DescribeCmsAlertNotifyTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alert_notify_template_id", objectRaw["alertNotifyTemplateId"])
	d.Set("alert_notify_template_name", objectRaw["alertNotifyTemplateName"])
	d.Set("type", objectRaw["type"])
	d.Set("program_lang", objectRaw["programLang"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("templates", flattenAlertNotifyTemplates(objectRaw["templates"]))

	return nil
}

func resourceAliCloudCmsAlertNotifyTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	templateId := d.Id()
	action := fmt.Sprintf("/alertNotifyTemplate/%s", templateId)
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	query["regionId"] = StringPointer(client.RegionId)

	if d.HasChange("alert_notify_template_name") || d.HasChange("templates") {
		if v, ok := d.GetOk("alert_notify_template_name"); ok {
			body["alertNotifyTemplateName"] = v
		}
		if v, ok := d.GetOk("templates"); ok {
			body["templates"] = buildAlertNotifyTemplateRequestMap(v)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		var err error
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPatch("Cms", "2024-03-30", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, body)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCmsAlertNotifyTemplateRead(d, meta)
}

func resourceAliCloudCmsAlertNotifyTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/alertNotifyTemplates"
	var response map[string]interface{}
	query := make(map[string]*string)
	query["regionId"] = StringPointer(client.RegionId)
	query["alertNotifyTemplateIds"] = StringPointer(fmt.Sprintf(`["%s"]`, d.Id()))

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var err error
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("Cms", "2024-03-30", action, query, nil, nil, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, query)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

// buildAlertNotifyTemplateRequestMap converts the Terraform templates list
// (list of objects, each carrying a channel name and its fields) into the
// map[string]map[string]interface{} payload shape expected by the CMS
// AlertNotifyTemplate API, which keys templates by channel name.
func buildAlertNotifyTemplateRequestMap(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	if raw == nil {
		return result
	}
	list, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		channel, ok := m["channel"].(string)
		if !ok || channel == "" {
			continue
		}
		entry := make(map[string]interface{})
		if v, ok := m["title"]; ok && v != nil {
			entry["title"] = v
		}
		if v, ok := m["content"]; ok && v != nil {
			entry["content"] = v
		}
		result[channel] = entry
	}
	return result
}

// flattenAlertNotifyTemplates converts the API templates map (keyed by
// channel name) into the list-of-objects shape used by the Terraform
// TypeList schema. Each entry carries the channel name plus its title/content.
func flattenAlertNotifyTemplates(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if raw == nil {
		return result
	}
	m, ok := raw.(map[string]interface{})
	if !ok {
		return result
	}
	for channel, v := range m {
		entry := map[string]interface{}{
			"channel": channel,
		}
		if inner, ok := v.(map[string]interface{}); ok {
			if title, ok := inner["title"]; ok && title != nil {
				entry["title"] = fmt.Sprint(title)
			}
			if content, ok := inner["content"]; ok && content != nil {
				entry["content"] = fmt.Sprint(content)
			}
		}
		result = append(result, entry)
	}
	return result
}
