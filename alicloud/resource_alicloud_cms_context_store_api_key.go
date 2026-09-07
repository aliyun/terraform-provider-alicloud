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

func resourceAliCloudCmsContextStoreApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsContextStoreApiKeyCreate,
		Read:   resourceAliCloudCmsContextStoreApiKeyRead,
		Delete: resourceAliCloudCmsContextStoreApiKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"context_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsContextStoreApiKeyCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/workspace/%v/contextstore/%v/apikey", d.Get("workspace"), d.Get("context_store_name"))
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})

	request["name"] = d.Get("name")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_context_store_api_key", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", d.Get("workspace"), d.Get("context_store_name"), d.Get("name")))

	if v, ok := response["apiKey"]; ok && fmt.Sprint(v) != "" {
		d.Set("api_key", v)
	}

	return resourceAliCloudCmsContextStoreApiKeyRead(d, meta)
}

func resourceAliCloudCmsContextStoreApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsContextStoreApiKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_context_store_api_key DescribeCmsContextStoreApiKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	parts := strings.SplitN(d.Id(), ":", 3)
	if v, ok := objectRaw["apiKey"]; ok && fmt.Sprint(v) != "" {
		d.Set("api_key", v)
	}
	d.Set("create_time", objectRaw["createTime"])
	d.Set("region_id", client.RegionId)
	if v, ok := objectRaw["name"]; ok && fmt.Sprint(v) != "" {
		d.Set("name", v)
	} else {
		d.Set("name", parts[2])
	}
	if v, ok := objectRaw["workspace"]; ok && fmt.Sprint(v) != "" {
		d.Set("workspace", v)
	} else {
		d.Set("workspace", parts[0])
	}
	if v, ok := objectRaw["contextStoreName"]; ok && fmt.Sprint(v) != "" {
		d.Set("context_store_name", v)
	} else {
		d.Set("context_store_name", parts[1])
	}

	return nil
}

func resourceAliCloudCmsContextStoreApiKeyDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.SplitN(d.Id(), ":", 3)
	action := fmt.Sprintf("/workspace/%v/contextstore/%v/apikey/%v", parts[0], parts[1], parts[2])
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

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
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"WorkspaceNotExist", "ContextStoreNotExist", "NotFound", "EntityNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
