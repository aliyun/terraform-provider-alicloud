package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudCdnPreloadObjectCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCdnPreloadObjectCacheCreate,
		Read:   resourceAliCloudCdnPreloadObjectCacheRead,
		Delete: resourceAliCloudCdnPreloadObjectCacheDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"object_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"area": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"domestic", "overseas"}, false),
			},
			"l2_preload": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCdnPreloadObjectCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "PushObjectCache"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["ObjectPath"] = d.Get("object_path")
	if v, ok := d.GetOk("area"); ok {
		request["Area"] = v
	}
	if v, ok := d.GetOkExists("l2_preload"); ok {
		request["L2Preload"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, query, request, true)
		if err != nil {
			// PushObjectCache operates on edge nodes; after the CDN domain
			// control plane reports online, the edge data plane may lag briefly,
			// returning InvalidDomain.Offline. Retry within the create timeout
			// until the edge is ready.
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InvalidDomain.Offline"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cdn_preload_object_cache", action, AlibabaCloudSdkGoERROR)
	}

	// M-RL-0043: create output PushTaskId is the TaskId
	d.SetId(fmt.Sprint(response["PushTaskId"]))

	return resourceAliCloudCdnPreloadObjectCacheRead(d, meta)
}

func resourceAliCloudCdnPreloadObjectCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnServiceV2 := CdnServiceV2{client}

	object, err := cdnServiceV2.DescribeCdnPreloadObjectCache(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cdn_preload_object_cache DescribeCdnPreloadObjectCache Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("object_path", object["content"])
	d.Set("process", object["process"])
	d.Set("status", object["status"])
	d.Set("creation_time", object["gmtCreatedStr"])
	// M-RT-0044: area and l2_preload are create-only, not backfilled in Read

	return nil
}

func resourceAliCloudCdnPreloadObjectCacheDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource alicloud_cdn_preload_object_cache. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
