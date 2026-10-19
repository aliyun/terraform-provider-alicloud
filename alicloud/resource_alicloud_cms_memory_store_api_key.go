// Package alicloud. Hand-written implementation for the CMS MemoryStoreAPIKey resource.
// Backing OpenAPI: Cms 2024-03-30 (RESTful/ROA, @visibility Private).
// CRUD: CreateMemoryStoreAPIKey (POST), ListMemoryStoreAPIKeys (GET, single-resource
// read via list filter), DeleteMemoryStoreAPIKey (DELETE). No Update API exists by
// design — the API Key name is immutable and the key value is server-generated, so any
// user-settable field change forces recreation.
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

func resourceAliCloudCmsMemoryStoreAPIKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsMemoryStoreAPIKeyCreate,
		Read:   resourceAliCloudCmsMemoryStoreAPIKeyRead,
		Delete: resourceAliCloudCmsMemoryStoreAPIKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory_store_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsMemoryStoreAPIKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	workspace := d.Get("workspace")
	memoryStoreName := d.Get("memory_store_name")
	name := d.Get("name")
	action := fmt.Sprintf("/workspace/%s/memorystore/%s/apikey", workspace, memoryStoreName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["name"] = name
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_memory_store_api_key", action, AlibabaCloudSdkGoERROR)
	}

	responseName, ok := response["name"]
	if !ok || fmt.Sprint(responseName) == "" {
		responseName = name
	}
	d.SetId(fmt.Sprintf("%s:%s:%s", workspace, memoryStoreName, responseName))

	return resourceAliCloudCmsMemoryStoreAPIKeyRead(d, meta)
}

func resourceAliCloudCmsMemoryStoreAPIKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsMemoryStoreAPIKey(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_memory_store_api_key DescribeCmsMemoryStoreAPIKey Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("api_key", objectRaw["apiKey"])
	d.Set("create_time", objectRaw["createTime"])
	d.Set("memory_store_name", objectRaw["memoryStoreName"])
	d.Set("name", objectRaw["name"])
	d.Set("region_id", client.RegionId)
	d.Set("workspace", objectRaw["workspace"])

	return nil
}

func resourceAliCloudCmsMemoryStoreAPIKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", d.Id(), 3, len(parts)))
	}
	workspace := parts[0]
	memoryStoreName := parts[1]
	name := parts[2]
	action := fmt.Sprintf("/workspace/%s/memorystore/%s/apikey/%s", workspace, memoryStoreName, name)
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
		if IsExpectedErrors(err, []string{"MemoryStoreAPIKeyNotExist", "NotFound", "WorkspaceNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
