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

func resourceAliCloudCdnRefreshObjectCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCdnRefreshObjectCacheCreate,
		Read:   resourceAliCloudCdnRefreshObjectCacheRead,
		Delete: resourceAliCloudCdnRefreshObjectCacheDelete,
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
			"object_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "File",
				ValidateFunc: validation.StringInSlice([]string{"File", "Directory"}, false),
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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

func resourceAliCloudCdnRefreshObjectCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "RefreshObjectCaches"
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	request["ObjectPath"] = d.Get("object_path")
	if v, ok := d.GetOk("object_type"); ok {
		request["ObjectType"] = v
	}
	if v, ok := d.GetOkExists("force"); ok {
		request["Force"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cdn_refresh_object_cache", action, AlibabaCloudSdkGoERROR)
	}

	// M-RL-0043: create output RefreshTaskId is the TaskId
	d.SetId(fmt.Sprint(response["RefreshTaskId"]))

	return resourceAliCloudCdnRefreshObjectCacheRead(d, meta)
}

func resourceAliCloudCdnRefreshObjectCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cdnServiceV2 := CdnServiceV2{client}

	object, err := cdnServiceV2.DescribeCdnRefreshObjectCache(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cdn_refresh_object_cache DescribeCdnRefreshObjectCache Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("object_path", object["content"])
	d.Set("process", object["process"])
	d.Set("status", object["status"])
	d.Set("creation_time", object["gmtCreatedStr"])
	// M-RT-0044: object_type is create-only, not backfilled in Read

	return nil
}

func resourceAliCloudCdnRefreshObjectCacheDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource alicloud_cdn_refresh_object_cache. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
