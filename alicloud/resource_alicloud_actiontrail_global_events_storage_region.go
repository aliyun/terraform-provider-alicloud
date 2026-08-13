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
				ValidateFunc: validation.StringInSlice([]string{
					"ap-southeast-1",
					"cn-hangzhou",
				}, false),
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
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			resp, err := client.RpcGet("Actiontrail", "2020-07-06", action, request, nil)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		// GetGlobalEventsStorageRegion is eventually consistent: immediately after
		// a successful UpdateGlobalEventsStorageRegion it can still report the
		// previous StorageRegion, so poll until the API returns the newly
		// configured region before refreshing state.
		actiontrailService := ActiontrailService{client}
		stateConf := BuildStateConf([]string{}, []string{fmt.Sprint(d.Get("storage_region"))}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, func() (interface{}, string, error) {
			object, err := actiontrailService.DescribeActiontrailGlobalEventsStorageRegion(d.Id())
			if err != nil {
				if NotFoundError(err) {
					return object, "", nil
				}
				return nil, "", WrapError(err)
			}
			return object, fmt.Sprint(object["StorageRegion"]), nil
		})
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudActiontrailGlobalEventsStorageRegionRead(d, meta)
}

func resourceAlicloudActiontrailGlobalEventsStorageRegionDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
