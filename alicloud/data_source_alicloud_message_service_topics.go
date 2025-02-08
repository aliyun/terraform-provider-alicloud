package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudMessageServiceTopics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMessageServiceTopicsRead,
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
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topics": {
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
						"message_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_message_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"message_retention_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"logging_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"topic_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_inner_url": {
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

func dataSourceAlicloudMessageServiceTopicsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTopic"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("topic_name"); ok {
		request["TopicName"] = v
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
	var topicNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		topicNameRegex = r
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
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Mns-open", "2022-01-19", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_message_service_topics", action, AlibabaCloudSdkGoERROR)
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
			if topicNameRegex != nil && !topicNameRegex.MatchString(fmt.Sprint(item["TopicName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TopicName"])]; !ok {
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
			"id":                       fmt.Sprint(object["TopicName"]),
			"topic_name":               fmt.Sprint(object["TopicName"]),
			"message_count":            formatInt(object["MessageCount"]),
			"max_message_size":         formatInt(object["MaxMessageSize"]),
			"message_retention_period": formatInt(object["MessageRetentionPeriod"]),
			"logging_enabled":          object["LoggingEnabled"],
			"topic_url":                object["TopicUrl"],
			"topic_inner_url":          object["TopicInnerUrl"],
			"last_modify_time":         formatInt(object["LastModifyTime"]),
			"create_time":              formatInt(object["CreateTime"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["TopicName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("topics", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
