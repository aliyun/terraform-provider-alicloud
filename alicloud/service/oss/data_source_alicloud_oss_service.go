package oss

import (
	"fmt"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceAlicloudOssService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOssServiceRead,

		Schema: map[string]*schema.Schema{
			"enable": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"On", "Off"}, false),
				Optional:     true,
				Default:      "Off",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceAlicloudOssServiceRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("enable"); !ok || v.(string) != "On" {
		d.SetId("OssServiceHasNotBeenOpened")
		d.Set("status", "")
		return nil
	}

	conn, err := meta.(*connectivity.AliyunClient).NewTeaCommonClient(connectivity.OpenOssService)
	if err != nil {
		return errs.WrapError(err)
	}
	response, err := conn.DoRequest(helper.StringPointer("OpenOssService"), nil, helper.StringPointer("POST"), helper.StringPointer("2019-04-22"), helper.StringPointer("AK"), nil, nil, &util.RuntimeOptions{})

	logs.AddDebug("OpenOssService", response, nil)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"SYSTEM.SALE_VALIDATE_NO_SPECIFIC_CODE_FAILEDError", "ORDER.OPEND"}) {
			d.SetId("OssServicHasBeenOpened")
			d.Set("status", "Opened")
			return nil
		}
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_oss_service", "OpenOssService", errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v", response["OrderId"]))
	d.Set("status", "Opened")

	return nil
}
