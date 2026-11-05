package dts

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceConsumerChannels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConsumerChannelsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dts_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumer_group_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumption_checkpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_delay": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unconsumed_data": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceConsumerChannelsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeConsumerChannel"
	request := make(map[string]interface{})
	request["DtsInstanceId"] = d.Get("dts_instance_id")
	request["RegionId"] = client.RegionId
	request["PageSize"] = helper.PageSizeMedium
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	var err error
	for {
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {

			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, true)
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
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_dts_consumer_channels", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ConsumerChannels", response)
		if err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ConsumerChannels", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["DtsInstanceId"], ":", item["ConsumerGroupId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < helper.PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		message_delayVal, err := convert.FormatInt(object["MessageDelay"])
		if err != nil {
			return err
		}
		unconsumed_dataVal, err := convert.FormatInt(object["UnconsumedData"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"id":                       fmt.Sprint(request["DtsInstanceId"], ":", object["ConsumerGroupId"]),
			"consumer_group_id":        fmt.Sprint(object["ConsumerGroupId"]),
			"consumer_group_name":      object["ConsumerGroupName"],
			"consumer_group_user_name": object["ConsumerGroupUserName"],
			"consumption_checkpoint":   fmt.Sprint(object["ConsumptionCheckpoint"]),
			"message_delay":            message_delayVal,
			"unconsumed_data":          unconsumed_dataVal,
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("channels", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
