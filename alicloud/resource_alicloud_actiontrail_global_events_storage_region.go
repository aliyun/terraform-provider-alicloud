package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlicloudActiontrailGlobalEventsStorageRegion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudActiontrailGlobalEventsStorageRegionCreate,
		Read:   resourceAlicloudActiontrailGlobalEventsStorageRegionRead,
		Update: resourceAlicloudActiontrailGlobalEventsStorageRegionUpdate,
		Delete: resourceAlicloudActiontrailGlobalEventsStorageRegionDelete,
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"storage_region": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func resourceAlicloudActiontrailGlobalEventsStorageRegionCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(fmt.Sprint("GlobalEventsStorageRegion"))

	return resourceAlicloudActiontrailGlobalEventsStorageRegionUpdate(d, meta)
}

func resourceAlicloudActiontrailGlobalEventsStorageRegionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := ActiontrailService{client}

	object, err := actiontrailService.DescribeActiontrailGlobalEventsStorageRegion(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_global_events_storage_region actiontrailService.DescribeActiontrailGlobalEventsStorageRegion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("storage_region", object["StorageRegion"])

	return nil
}

func resourceAlicloudActiontrailGlobalEventsStorageRegionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var err error
	update := false
	request := map[string]interface{}{}
	if d.HasChange("storage_region") {
		update = true
		if v, ok := d.GetOk("storage_region"); ok {
			request["StorageRegion"] = v
		}
	}

	if update {
		action := "UpdateGlobalEventsStorageRegion"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			resp, err := client.RpcGet("Actiontrail", "2020-07-06", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAlicloudActiontrailGlobalEventsStorageRegionRead(d, meta)
}

func resourceAlicloudActiontrailGlobalEventsStorageRegionDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
