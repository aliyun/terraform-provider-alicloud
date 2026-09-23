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

func resourceAliCloudGpdbApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudGpdbApiKeyCreate,
		Read:   resourceAliCloudGpdbApiKeyRead,
		Delete: resourceAliCloudGpdbApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudGpdbApiKeyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateApiKey"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("workspace_id"); ok {
		request["WorkspaceId"] = v
	}
	request["RegionId"] = client.RegionId

	request["KeyName"] = d.Get("key_name")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("service_ids"); ok {
		serviceIdsJsonPath, err := jsonpath.Get("$", v)
		if err == nil && serviceIdsJsonPath != "" {
			request["ServiceIds"] = convertListToJsonString(convertToInterfaceArray(serviceIdsJsonPath))
		}
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_gpdb_api_key", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["WorkspaceId"], response["KeyId"]))

	return resourceAliCloudGpdbApiKeyRead(d, meta)
}

func resourceAliCloudGpdbApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gpdbServiceV2 := GpdbServiceV2{client}

	objectRaw, err := gpdbServiceV2.DescribeGpdbApiKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_gpdb_api_key DescribeGpdbApiKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("key_id", objectRaw["KeyId"])
	d.Set("key_name", objectRaw["KeyName"])
	d.Set("description", objectRaw["Description"])

	if v, ok := objectRaw["AuthServices"]; ok && v != nil {
		authServices, _ := v.([]interface{})
		serviceIds := make([]string, 0, len(authServices))
		for _, item := range authServices {
			if m, ok := item.(map[string]interface{}); ok {
				if sid, ok := m["ServiceId"]; ok {
					serviceIds = append(serviceIds, fmt.Sprint(sid))
				}
			}
		}
		d.Set("service_ids", serviceIds)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("workspace_id", parts[0])

	return nil
}

func resourceAliCloudGpdbApiKeyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteApiKey"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["KeyId"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("gpdb", "2016-05-03", action, query, request, true)
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
