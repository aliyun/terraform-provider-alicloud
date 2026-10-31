package config

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDeliveryChannels() *schema.Resource {
	return &schema.Resource{
		Read:               DataSourceDeliveryChannelsRead,
		DeprecationMessage: "This data source has been deprecated since v1.173.0 and will be removed in the future. Please use 'alicloud_config_deliveries' instead.",
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
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
						"delivery_channel_assume_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_target_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delivery_channel_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceDeliveryChannelsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDeliveryChannels"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var deliveryChannelNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
		}
		deliveryChannelNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOkExists("status")
	response, err := client.RpcGet("Config", "2019-01-08", action, request, nil)
	if err != nil {
		return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_config_delivery_channels", action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)

	resp, err := jsonpath.Get("$.DeliveryChannels", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.DeliveryChannels", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if deliveryChannelNameRegex != nil {
			if !deliveryChannelNameRegex.MatchString(item["DeliveryChannelName"].(string)) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DeliveryChannelId"])]; !ok {
				continue
			}
		}
		if statusOk {
			statusVal, statusErr := convert.FormatInt(item["Status"])
			if statusErr != nil {
				return statusErr
			}
			if status.(int) != statusVal {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		statusVal, err := convert.FormatInt(object["Status"])
		if err != nil {
			return err
		}
		mapping := map[string]interface{}{
			"delivery_channel_assume_role_arn": object["DeliveryChannelAssumeRoleArn"],
			"delivery_channel_condition":       object["DeliveryChannelCondition"],
			"id":                               fmt.Sprint(object["DeliveryChannelId"]),
			"delivery_channel_id":              fmt.Sprint(object["DeliveryChannelId"]),
			"delivery_channel_name":            object["DeliveryChannelName"],
			"delivery_channel_target_arn":      object["DeliveryChannelTargetArn"],
			"delivery_channel_type":            object["DeliveryChannelType"],
			"description":                      object["Description"],
			"status":                           statusVal,
		}
		ids = append(ids, fmt.Sprint(object["DeliveryChannelId"]))
		names = append(names, object["DeliveryChannelName"].(string))
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
