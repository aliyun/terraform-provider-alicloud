package alicloud

import (
	"fmt"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDcdnKvAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDcdnKvAccountRead,

		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Computed:     true,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"online", "offline"}, false),
			},
		},
	}
}
func dataSourceAlicloudDcdnKvAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	object, err := dcdnService.DescribeDcdnKvAccountStatus()
	d.SetId("DcdnKvAccount")
	if err != nil {
		d.Set("status", "")
		return WrapError(err)
	}

	if v, ok := d.GetOk("status"); !ok || v.(string) == fmt.Sprint(object["Status"]) {
		d.Set("status", object["Status"])

		return nil

	}

	action := "PutDcdnKvAccount"
	request := map[string]interface{}{
		"AccountType": "prod",
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenDcdnService)
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				return resource.RetryableError(err)
			}
			addDebug(action, response, nil)
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, nil)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dcdn_kv_account", action, AlibabaCloudSdkGoERROR)
	}

	d.Set("status", request["Status"])

	return nil
}
