package actiontrail

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGlobalEventsStorageRegion() *schema.Resource {
	return &schema.Resource{
		Create: ResourceGlobalEventsStorageRegionCreate,
		Read:   ResourceGlobalEventsStorageRegionRead,
		Update: ResourceGlobalEventsStorageRegionUpdate,
		Delete: ResourceGlobalEventsStorageRegionDelete,
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

func ResourceGlobalEventsStorageRegionCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(fmt.Sprint("GlobalEventsStorageRegion"))

	return ResourceGlobalEventsStorageRegionUpdate(d, meta)
}

func ResourceGlobalEventsStorageRegionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	actiontrailService := NewActiontrailService(client)

	object, err := actiontrailService.DescribeActiontrailGlobalEventsStorageRegion(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_actiontrail_global_events_storage_region actiontrailService.DescribeActiontrailGlobalEventsStorageRegion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	d.Set("storage_region", object["StorageRegion"])

	return nil
}

func ResourceGlobalEventsStorageRegionUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			resp, err := client.RpcGet("Actiontrail", "2020-07-06", action, request, nil)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, resp, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return ResourceGlobalEventsStorageRegionRead(d, meta)
}

func ResourceGlobalEventsStorageRegionDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}
