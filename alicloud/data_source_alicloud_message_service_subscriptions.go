package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudMessageServiceSubscriptions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMessageServiceSubscriptionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscription_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PageSizeLarge,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_tag": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_content_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_modify_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMessageServiceSubscriptionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSubscriptionByTopic"
	request := make(map[string]interface{})
	request["TopicName"] = d.Get("topic_name")

	if v, ok := d.GetOk("subscription_name"); ok {
		request["SubscriptionName"] = v
	}

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNum"] = v.(int)
	} else {
		request["PageNum"] = 1
	}

	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}

	var objects []map[string]interface{}
	var subscriptionNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		subscriptionNameRegex = r
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

	var response map[string]interface{}
	conn, err := client.NewMnsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-01-19"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_message_service_subscriptions", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Data.PageData", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.PageData", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if subscriptionNameRegex != nil && !subscriptionNameRegex.MatchString(fmt.Sprint(item["SubscriptionName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TopicName"], ":", item["SubscriptionName"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                    fmt.Sprint(object["TopicName"], ":", object["SubscriptionName"]),
			"topic_name":            object["TopicName"],
			"subscription_name":     object["SubscriptionName"],
			"endpoint":              object["Endpoint"],
			"filter_tag":            object["FilterTag"],
			"notify_content_format": object["NotifyContentFormat"],
			"notify_strategy":       object["NotifyStrategy"],
			"topic_owner":           object["TopicOwner"],
			"subscription_url":      object["SubscriptionURL"],
			"last_modify_time":      formatInt(object["LastModifyTime"]),
			"create_time":           formatInt(object["CreateTime"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["SubscriptionName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("subscriptions", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
