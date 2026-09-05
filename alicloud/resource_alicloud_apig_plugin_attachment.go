// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
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

func resourceAliCloudApigPluginAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudApigPluginAttachmentCreate,
		Read:   resourceAliCloudApigPluginAttachmentRead,
		Update: resourceAliCloudApigPluginAttachmentUpdate,
		Delete: resourceAliCloudApigPluginAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"attach_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"attach_resource_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attach_resource_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attach_resource_parent_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attach_resource_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"plugin_class_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execute_priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execute_stage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"plugin_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"plugin_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudApigPluginAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/v1/plugin-attachments")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	if v, ok := d.GetOk("plugin_info"); ok {
		pluginInfoPluginConfigJsonPath, err := jsonpath.Get("$[0].plugin_config", v)
		if err == nil && pluginInfoPluginConfigJsonPath != "" {
			request["pluginConfig"] = pluginInfoPluginConfigJsonPath
		}
	}
	if v, ok := d.GetOk("environment_id"); ok {
		request["environmentId"] = v
	}
	if v, ok := d.GetOk("plugin_info"); ok {
		pluginInfoPluginIdJsonPath, err := jsonpath.Get("$[0].plugin_id", v)
		if err == nil && pluginInfoPluginIdJsonPath != "" {
			request["pluginId"] = pluginInfoPluginIdJsonPath
		}
	}
	if v, ok := d.GetOk("attach_resource_type"); ok {
		request["attachResourceType"] = v
	}
	if v, ok := d.GetOk("attach_resource_ids"); ok {
		attachResourceIdsMapsArray := convertToInterfaceArray(v)

		request["attachResourceIds"] = attachResourceIdsMapsArray
	}

	if v, ok := d.GetOk("plugin_info"); ok {
		pluginInfoGatewayIdJsonPath, err := jsonpath.Get("$[0].gateway_id", v)
		if err == nil && pluginInfoGatewayIdJsonPath != "" {
			request["gatewayId"] = pluginInfoGatewayIdJsonPath
		}
	}
	if v, ok := d.GetOkExists("enable"); ok {
		request["enable"] = v
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("APIG", "2024-03-27", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_apig_plugin_attachment", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.pluginAttachmentId", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudApigPluginAttachmentUpdate(d, meta)
}

func resourceAliCloudApigPluginAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	apigServiceV2 := ApigServiceV2{client}

	var objectRaw map[string]interface{}
	var err error
	if v := d.Get("plugin_info.0.gateway_id"); v != nil && fmt.Sprint(v) != "" {
		objectRaw, err = apigServiceV2.DescribeApigPluginAttachmentWithGateway(d.Id(), fmt.Sprint(v))
	} else {
		objectRaw, err = apigServiceV2.DescribeApigPluginAttachment(d.Id())
	}
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_apig_plugin_attachment DescribeApigPluginAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("attach_resource_type", objectRaw["attachResourceType"])
	d.Set("enable", objectRaw["enable"])

	environmentInfoRawObj, _ := jsonpath.Get("$.environmentInfo", objectRaw)
	environmentInfoRaw := make(map[string]interface{})
	if environmentInfoRawObj != nil {
		environmentInfoRaw = environmentInfoRawObj.(map[string]interface{})
	}
	d.Set("environment_id", environmentInfoRaw["environmentId"])

	resourceInfosRawObj, _ := jsonpath.Get("$.resourceInfos[*]", objectRaw)
	resourceInfosRaw := make([]interface{}, 0)
	if resourceInfosRawObj != nil {
		resourceInfosRaw = convertToInterfaceArray(resourceInfosRawObj)
	}

	attachResourceIds := make([]string, 0)
	attachResourceNames := make([]string, 0)
	for _, item := range resourceInfosRaw {
		if info, ok := item.(map[string]interface{}); ok {
			if v, ok := info["resourceId"]; ok && fmt.Sprint(v) != "" {
				attachResourceIds = append(attachResourceIds, fmt.Sprint(v))
			}
			if v, ok := info["resourceName"]; ok && fmt.Sprint(v) != "" {
				attachResourceNames = append(attachResourceNames, fmt.Sprint(v))
			}
		}
	}
	d.Set("attach_resource_ids", attachResourceIds)
	d.Set("attach_resource_names", attachResourceNames)
	if len(attachResourceIds) > 0 {
		d.Set("attach_resource_id", attachResourceIds[0])
	}

	parentResourceInfoRawObj, _ := jsonpath.Get("$.parentResourceInfo", objectRaw)
	parentResourceInfoRaw := make(map[string]interface{})
	if parentResourceInfoRawObj != nil {
		parentResourceInfoRaw = parentResourceInfoRawObj.(map[string]interface{})
	}
	attachResourceParentIds := make([]string, 0)
	if apiInfoRaw, ok := parentResourceInfoRaw["apiInfo"].(map[string]interface{}); ok {
		if v, ok := apiInfoRaw["httpApiId"]; ok && fmt.Sprint(v) != "" {
			attachResourceParentIds = append(attachResourceParentIds, fmt.Sprint(v))
		}
	}
	d.Set("attach_resource_parent_ids", attachResourceParentIds)

	pluginClassInfoMaps := make([]map[string]interface{}, 0)
	pluginClassInfoMap := make(map[string]interface{})
	pluginClassInfoRaw := make(map[string]interface{})
	if objectRaw["pluginClassInfo"] != nil {
		pluginClassInfoRaw = objectRaw["pluginClassInfo"].(map[string]interface{})
	}
	if len(pluginClassInfoRaw) > 0 {
		pluginClassInfoMap["direction"] = pluginClassInfoRaw["mode"]
		pluginClassInfoMap["execute_priority"] = pluginClassInfoRaw["executePriority"]
		pluginClassInfoMap["execute_stage"] = pluginClassInfoRaw["executeStage"]
		pluginClassInfoMap["name"] = pluginClassInfoRaw["name"]
		pluginClassInfoMap["type"] = pluginClassInfoRaw["type"]

		pluginClassInfoMaps = append(pluginClassInfoMaps, pluginClassInfoMap)
	}
	if err := d.Set("plugin_class_info", pluginClassInfoMaps); err != nil {
		return err
	}

	pluginInfoMaps := make([]map[string]interface{}, 0)
	pluginInfoMap := make(map[string]interface{})

	pluginInfoMap["plugin_config"] = objectRaw["pluginConfig"]
	pluginInfoMap["plugin_id"] = objectRaw["pluginId"]

	// ListPluginAttachments does not return gatewayId for the attachment.
	// Fall back to the value already in state to avoid a perpetual ForceNew diff.
	gatewayId := ""
	gatewayInfoRawObj, _ := jsonpath.Get("$.environmentInfo.gatewayInfo", objectRaw)
	gatewayInfoRaw := make(map[string]interface{})
	if gatewayInfoRawObj != nil {
		gatewayInfoRaw = gatewayInfoRawObj.(map[string]interface{})
	}
	if v, ok := gatewayInfoRaw["gatewayId"]; ok && fmt.Sprint(v) != "" {
		gatewayId = fmt.Sprint(v)
	} else if v := d.Get("plugin_info.0.gateway_id"); v != nil {
		gatewayId = fmt.Sprint(v)
	}
	pluginInfoMap["gateway_id"] = gatewayId

	pluginInfoMaps = append(pluginInfoMaps, pluginInfoMap)
	if err := d.Set("plugin_info", pluginInfoMaps); err != nil {
		return err
	}

	return nil
}

func resourceAliCloudApigPluginAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	pluginAttachmentId := d.Id()
	action := fmt.Sprintf("/v1/plugin-attachments/%s", pluginAttachmentId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if !d.IsNewResource() && d.HasChange("plugin_info.0.plugin_config") {
		update = true
	}
	if v, ok := d.GetOk("plugin_info"); ok || d.HasChange("plugin_info") {
		pluginInfoPluginConfigJsonPath, err := jsonpath.Get("$[0].plugin_config", v)
		if err == nil && pluginInfoPluginConfigJsonPath != "" {
			request["pluginConfig"] = pluginInfoPluginConfigJsonPath
		}
	}
	if !d.IsNewResource() && d.HasChange("attach_resource_ids") {
		update = true
	}
	if v, ok := d.GetOk("attach_resource_ids"); ok || d.HasChange("attach_resource_ids") {
		attachResourceIdsMapsArray := convertToInterfaceArray(v)

		request["attachResourceIds"] = attachResourceIdsMapsArray
	}

	if !d.IsNewResource() && d.HasChange("enable") {
		update = true
	}
	if v, ok := d.GetOkExists("enable"); ok || d.HasChange("enable") {
		request["enable"] = v
	}
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("APIG", "2024-03-27", action, query, nil, body, true)
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

	return resourceAliCloudApigPluginAttachmentRead(d, meta)
}

func resourceAliCloudApigPluginAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	pluginAttachmentId := d.Id()
	action := fmt.Sprintf("/v1/plugin-attachments/%s", pluginAttachmentId)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("APIG", "2024-03-27", action, query, nil, nil, true)
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
