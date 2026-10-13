// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec EndpointConnector definition.
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCmsEndpointConnector() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsEndpointConnectorCreate,
		Read:   resourceAliCloudCmsEndpointConnectorRead,
		Update: resourceAliCloudCmsEndpointConnectorUpdate,
		Delete: resourceAliCloudCmsEndpointConnectorDelete,
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
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"model_service", "agent_app"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"credential": {
				Type:      schema.TypeMap,
				Required:  true,
				Sensitive: true,
				Elem:      &schema.Schema{Type: schema.TypeString},
			},
			"headers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connector_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsEndpointConnectorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	workspace := d.Get("workspace")
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s", workspace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["type"] = d.Get("type")
	request["name"] = d.Get("name")
	if v, ok := d.GetOk("alias"); ok {
		request["alias"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	request["endpoint"] = d.Get("endpoint")
	request["credential"] = d.Get("credential")
	if v, ok := d.GetOk("headers"); ok {
		request["headers"] = buildCmsEndpointConnectorHeaders(v.([]interface{}))
	}
	if v, ok := d.GetOk("properties"); ok {
		request["properties"] = v
	}
	body = request
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_endpoint_connector", action, AlibabaCloudSdkGoERROR)
	}

	connectorId, ok := response["connectorId"]
	if !ok || connectorId == nil || fmt.Sprint(connectorId) == "" {
		return WrapErrorf(fmt.Errorf("connector_id is empty in the CreateEndpointConnector response"), IdMsg, "alicloud_cms_endpoint_connector")
	}
	d.SetId(fmt.Sprintf("%v:%v", workspace, connectorId))

	return resourceAliCloudCmsEndpointConnectorRead(d, meta)
}

func resourceAliCloudCmsEndpointConnectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsEndpointConnector(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_endpoint_connector DescribeCmsEndpointConnector Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("connector_id", objectRaw["connectorId"])
	d.Set("workspace", objectRaw["workspace"])
	d.Set("type", objectRaw["type"])
	d.Set("name", objectRaw["name"])
	d.Set("alias", objectRaw["alias"])
	d.Set("endpoint", objectRaw["endpoint"])
	d.Set("description", objectRaw["description"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("created_at", objectRaw["createdAt"])
	d.Set("updated_at", objectRaw["updatedAt"])

	if v, ok := objectRaw["credential"]; ok && v != nil {
		d.Set("credential", v)
	}
	if v, ok := objectRaw["properties"]; ok && v != nil {
		d.Set("properties", v)
	}
	if v, ok := objectRaw["headers"]; ok && v != nil {
		d.Set("headers", flattenCmsEndpointConnectorHeaders(v))
	}

	return nil
}

func resourceAliCloudCmsEndpointConnectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	var err error

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 2, len(parts)))
	}
	workspace := parts[0]
	connectorId := parts[1]
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", workspace, connectorId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	request["name"] = d.Get("name")
	if v, ok := d.GetOk("alias"); ok || d.HasChange("alias") {
		request["alias"] = v
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	request["endpoint"] = d.Get("endpoint")
	request["credential"] = d.Get("credential")
	if v, ok := d.GetOk("headers"); ok || d.HasChange("headers") {
		request["headers"] = buildCmsEndpointConnectorHeaders(v.([]interface{}))
	}
	if v, ok := d.GetOk("properties"); ok || d.HasChange("properties") {
		request["properties"] = v
	}
	body = request
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

	return resourceAliCloudCmsEndpointConnectorRead(d, meta)
}

func resourceAliCloudCmsEndpointConnectorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 2, len(parts)))
	}
	workspace := parts[0]
	connectorId := parts[1]
	action := fmt.Sprintf("/api/v1/endpoint-connectors/%s/%s", workspace, connectorId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
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
		if IsExpectedErrors(err, []string{"EndpointConnectorNotExist", "WorkspaceNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func buildCmsEndpointConnectorHeaders(headers []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(headers))
	for _, h := range headers {
		if h == nil {
			continue
		}
		hm := h.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"key":   hm["key"],
			"value": hm["value"],
		})
	}
	return result
}

func flattenCmsEndpointConnectorHeaders(v interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	if list, ok := v.([]interface{}); ok {
		for _, h := range list {
			if h == nil {
				continue
			}
			hm := h.(map[string]interface{})
			result = append(result, map[string]interface{}{
				"key":   hm["key"],
				"value": hm["value"],
			})
		}
	}
	return result
}
