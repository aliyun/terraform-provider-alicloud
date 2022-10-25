package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudConfigDeliveries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigDeliveriesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
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
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"configuration_item_change_notification": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"configuration_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
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
						"non_compliant_notification": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oversized_data_oss_target_arn": {
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

func dataSourceAlicloudConfigDeliveriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListConfigDeliveryChannels"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var channelNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		channelNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-07"), StringPointer("AK"), nil, request, &runtime)
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
		if IsExpectedErrors(err, []string{"DeliveryChannelNotExists"}) {
			d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
			return nil
		}
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_deliveries", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.DeliveryChannels", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DeliveryChannels", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if channelNameRegex != nil {
			if !channelNameRegex.MatchString(item["DeliveryChannelName"].(string)) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DeliveryChannelId"])]; !ok {
				continue
			}
		}
		if statusOk && status.(int) != formatInt(item["Status"]) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_id":                             fmt.Sprint(object["AccountId"]),
			"configuration_item_change_notification": object["ConfigurationItemChangeNotification"],
			"configuration_snapshot":                 object["ConfigurationSnapshot"],
			"delivery_channel_assume_role_arn":       object["DeliveryChannelAssumeRoleArn"],
			"delivery_channel_condition":             object["DeliveryChannelCondition"],
			"id":                                     fmt.Sprint(object["DeliveryChannelId"]),
			"delivery_channel_id":                    fmt.Sprint(object["DeliveryChannelId"]),
			"delivery_channel_name":                  object["DeliveryChannelName"],
			"delivery_channel_target_arn":            object["DeliveryChannelTargetArn"],
			"delivery_channel_type":                  object["DeliveryChannelType"],
			"description":                            object["Description"],
			"non_compliant_notification":             object["NonCompliantNotification"],
			"oversized_data_oss_target_arn":          object["OversizedDataOSSTargetArn"],
			"status":                                 formatInt(object["Status"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DeliveryChannelName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("deliveries", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
