package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudRamAccountAlias() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamAccountAliasRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudRamAccountAliasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	action := "GetAccountAlias"
	var response map[string]interface{}
	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutRead)), func() *retry.RetryError {
		response, err = client.RpcPost("Ram", "2015-05-01", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_account_alias", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AccountAlias"]))
	d.Set("account_alias", response["AccountAlias"])
	if output, ok := d.GetOk("output_file"); ok && output != nil {
		s := map[string]interface{}{"account_alias": response["AccountAlias"]}
		writeToFile(output.(string), s)
	}
	return nil
}
