// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec Umodel definition.
package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsUmodel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsUmodelCreate,
		Read:   resourceAliCloudCmsUmodelRead,
		Update: resourceAliCloudCmsUmodelUpdate,
		Delete: resourceAliCloudCmsUmodelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// common_schema_ref is system-managed: the CMS Umodel Create/Update
			// API ignores any user-supplied value (the create body only accepts
			// description) and the server populates the default itself. It is
			// therefore Computed-only so users cannot set it, and Read reflects
			// the server-managed value returned by Get.
			"common_schema_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"items": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCmsUmodelRequestBody(d *schema.ResourceData) map[string]interface{} {
	body := make(map[string]interface{})
	if v, ok := d.GetOk("description"); ok {
		body["description"] = v
	}
	// common_schema_ref is system-managed and is not accepted by the
	// Create/Update API, so it is intentionally not sent in the request body.
	return body
}

func resourceAliCloudCmsUmodelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	workspace := d.Get("workspace").(string)
	action := fmt.Sprintf("/workspace/%s/umodel", workspace)
	query := make(map[string]*string)
	body := buildCmsUmodelRequestBody(d)
	var response map[string]interface{}
	var err error
	request := body

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_umodel", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(workspace)
	return resourceAliCloudCmsUmodelRead(d, meta)
}

func resourceAliCloudCmsUmodelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsUmodel(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_umodel DescribeCmsUmodel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// description is a write-only field: the CMS Umodel Create/Update API
	// accepts it but the Get API never returns it, so Read must not overwrite
	// the value applied from config (otherwise a perpetual non-empty plan
	// would be triggered on every refresh).
	d.Set("region_id", objectRaw["regionId"])
	d.Set("workspace", objectRaw["workspace"])

	if v, ok := objectRaw["commonSchemaRef"]; ok && v != nil {
		if refs, ok := v.([]interface{}); ok {
			refList := make([]map[string]interface{}, 0, len(refs))
			for _, ref := range refs {
				refMap, ok := ref.(map[string]interface{})
				if !ok {
					continue
				}
				m := make(map[string]interface{})
				if g, ok := refMap["group"]; ok {
					m["group"] = g
				}
				if rawItems, ok := refMap["items"].([]interface{}); ok {
					items := make([]string, 0, len(rawItems))
					for _, it := range rawItems {
						items = append(items, fmt.Sprint(it))
					}
					m["items"] = items
				}
				if ver, ok := refMap["version"]; ok {
					m["version"] = ver
				}
				refList = append(refList, m)
			}
			d.Set("common_schema_ref", refList)
		}
	}

	return nil
}

func resourceAliCloudCmsUmodelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	workspace := d.Id()
	action := fmt.Sprintf("/workspace/%s/umodel", workspace)
	query := make(map[string]*string)
	body := buildCmsUmodelRequestBody(d)
	var response map[string]interface{}
	var err error
	request := body

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.RoaPut("Cms", "2024-03-30", action, query, nil, body, true)
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

	return resourceAliCloudCmsUmodelRead(d, meta)
}

func resourceAliCloudCmsUmodelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	workspace := d.Id()
	action := fmt.Sprintf("/workspace/%s/umodel", workspace)
	query := make(map[string]*string)
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
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
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"UmodelNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
