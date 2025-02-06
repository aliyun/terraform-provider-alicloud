package alicloud

import (
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCloudSsoService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudSsoServiceRead,
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudCloudSsoServiceRead(d *schema.ResourceData, meta interface{}) error {
	var response map[string]interface{}
	request := map[string]interface{}{}
	client := meta.(*connectivity.AliyunClient)
	var err error
	enable := d.Get("enable").(string)
	if enable == "On" {
		action := "EnableService"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_sso_service", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId("CloudSsoServiceHasBeenOpened")
		d.Set("status", "Opened")
		return nil
	}

	if enable == "Off" {
		action := "DisableService"
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_sso_service", action, AlibabaCloudSdkGoERROR)
		}

		d.SetId("CloudSsoServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	return nil
}
