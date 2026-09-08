// Package alicloud
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

func resourceAliCloudCmsContextStore() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsContextStoreCreate,
		Read:   resourceAliCloudCmsContextStoreRead,
		Update: resourceAliCloudCmsContextStoreUpdate,
		Delete: resourceAliCloudCmsContextStoreDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"context_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"context_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata_field": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"logstore": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataset": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsContextStoreCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	workspace := d.Get("workspace").(string)
	contextStoreName := d.Get("context_store_name").(string)
	action := fmt.Sprintf("/workspace/%s/contextstore", workspace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["contextStoreName"] = contextStoreName
	request["contextType"] = d.Get("context_type")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("config"); ok && len(v.([]interface{})) > 0 {
		request["config"] = expandCmsContextStoreConfig(v.([]interface{}))
	}
	body = request
	wait := incrementalWait(3*time.Second, 0*time.Second)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_context_store", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s:%s", workspace, contextStoreName))

	return resourceAliCloudCmsContextStoreRead(d, meta)
}

func resourceAliCloudCmsContextStoreRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsContextStore(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_context_store DescribeCmsContextStore Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("workspace", objectRaw["workspace"])
	d.Set("context_store_name", objectRaw["contextStoreName"])
	d.Set("context_type", objectRaw["contextType"])
	d.Set("description", objectRaw["description"])
	d.Set("status", objectRaw["status"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("update_time", objectRaw["updateTime"])
	d.Set("config", flattenCmsContextStoreConfig(objectRaw))
	d.Set("dataset", flattenCmsContextStoreDataset(objectRaw))

	return nil
}

func resourceAliCloudCmsContextStoreUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("Invalid resource id: %s, expected format <workspace>:<context_store_name>", d.Id()))
	}
	workspace, contextStoreName := parts[0], parts[1]
	action := fmt.Sprintf("/workspace/%s/contextstore/%s", workspace, contextStoreName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if d.HasChange("context_type") || d.HasChange("description") || d.HasChange("config") {
		request["contextType"] = d.Get("context_type")
		request["description"] = d.Get("description")
		if v, ok := d.GetOk("config"); ok && len(v.([]interface{})) > 0 {
			request["config"] = expandCmsContextStoreConfig(v.([]interface{}))
		}
		body = request
		wait := incrementalWait(3*time.Second, 0*time.Second)
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
	}

	return resourceAliCloudCmsContextStoreRead(d, meta)
}

func resourceAliCloudCmsContextStoreDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("Invalid resource id: %s, expected format <workspace>:<context_store_name>", d.Id()))
	}
	workspace, contextStoreName := parts[0], parts[1]
	action := fmt.Sprintf("/workspace/%s/contextstore/%s", workspace, contextStoreName)
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error

	wait := incrementalWait(3*time.Second, 0*time.Second)
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
	addDebug(action, response, nil)

	if err != nil {
		if IsExpectedErrors(err, []string{"ContextStoreNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func expandCmsContextStoreConfig(configs []interface{}) map[string]interface{} {
	if len(configs) == 0 {
		return nil
	}
	config := configs[0].(map[string]interface{})
	configReq := map[string]interface{}{}
	if sourceRaw, ok := config["source"]; ok && len(sourceRaw.([]interface{})) > 0 {
		sourceMap := sourceRaw.([]interface{})[0].(map[string]interface{})
		sourceReq := map[string]interface{}{}
		if v, ok := sourceMap["project"].(string); ok && v != "" {
			sourceReq["project"] = v
		}
		if v, ok := sourceMap["logstore"].(string); ok && v != "" {
			sourceReq["logstore"] = v
		}
		if v, ok := sourceMap["start_time"].(string); ok && v != "" {
			sourceReq["startTime"] = v
		}
		configReq["source"] = sourceReq
	}
	if mf, ok := config["metadata_field"]; ok && mf != nil {
		if mfMap, ok := mf.(map[string]interface{}); ok && len(mfMap) > 0 {
			configReq["metadataField"] = mfMap
		}
	}
	return configReq
}

func flattenCmsContextStoreConfig(response map[string]interface{}) []interface{} {
	configList := make([]interface{}, 0)
	if configRaw, ok := response["config"]; ok && configRaw != nil {
		if configMap, ok := configRaw.(map[string]interface{}); ok {
			item := map[string]interface{}{}
			if sourceRaw, ok := configMap["source"]; ok && sourceRaw != nil {
				if sourceMap, ok := sourceRaw.(map[string]interface{}); ok {
					item["source"] = []interface{}{map[string]interface{}{
						"project":    sourceMap["project"],
						"logstore":   sourceMap["logstore"],
						"start_time": sourceMap["startTime"],
					}}
				}
			}
			if mf, ok := configMap["metadataField"]; ok && mf != nil {
				item["metadata_field"] = mf
			}
			configList = append(configList, item)
		}
	}
	return configList
}

func flattenCmsContextStoreDataset(response map[string]interface{}) []interface{} {
	if datasetRaw, ok := response["dataset"]; ok && datasetRaw != nil {
		if datasetMap, ok := datasetRaw.(map[string]interface{}); ok {
			return []interface{}{map[string]interface{}{
				"name": datasetMap["name"],
			}}
		}
	}
	return []interface{}{}
}
