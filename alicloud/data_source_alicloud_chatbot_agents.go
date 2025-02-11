package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudChatbotAgents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudChatbotAgentsRead,

		Schema: map[string]*schema.Schema{
			"agent_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"agents": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudChatbotAgentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "ListAgent"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("agent_name"); ok {
		request["AgentName"] = v
	}
	setPagingRequest(d, request, PageSizeLarge)
	var objects []map[string]interface{}
	var accountNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		accountNameRegex = r
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
			response, err = client.RpcPost("Chatbot", "2022-04-08", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_chatbot_agents", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if accountNameRegex != nil && !accountNameRegex.MatchString(fmt.Sprint(item["AgentName"])) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[request["AgentId"].(string)]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":         fmt.Sprint(object["AgentId"]),
			"agent_id":   fmt.Sprint(object["AgentId"]),
			"agent_key":  fmt.Sprint(object["AgentKey"]),
			"agent_name": object["AgentName"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["AgentName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("agents", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
