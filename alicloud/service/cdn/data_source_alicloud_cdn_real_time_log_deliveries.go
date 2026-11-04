package cdn

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlicloudCdnRealTimeLogDeliveries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCdnRealTimeLogDeliveriesRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"offline", "online"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deliveries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logstore": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCdnRealTimeLogDeliveriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDomainRealtimeLogDelivery"
	request := map[string]interface{}{}
	if v, ok := d.GetOk("domain"); ok {
		request["Domain"] = v
	}

	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcGet("Cdn", "2018-05-10", action, request, nil)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_cdn_real_time_log_deliveries", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$", response)
	}

	s := make([]map[string]interface{}, 0)
	if object, ok := resp.(map[string]interface{}); ok {
		if statusOk && status.(string) != "" && status.(string) == object["Status"].(string) {
			mapping := map[string]interface{}{
				"id":         fmt.Sprint(request["Domain"]),
				"domain":     fmt.Sprint(request["Domain"]),
				"logstore":   object["Logstore"],
				"project":    object["Project"],
				"sls_region": object["Region"],
				"status":     object["Status"],
			}
			s = append(s, mapping)
		}
	}

	d.SetId(fmt.Sprint(request["Domain"]))

	if err := d.Set("deliveries", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
