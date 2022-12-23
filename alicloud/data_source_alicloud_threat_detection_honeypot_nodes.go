package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudThreatDetectionHoneypotNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudThreatDetectionHoneypotNodesRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"node_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"node_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"page_number": {
				Optional: true,
				Type:     schema.TypeInt,
			},
			"page_size": {
				Optional: true,
				Type:     schema.TypeInt,
				Default:  20,
			},
			"nodes": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"allow_honeypot_access_internet": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"available_probe_num": {
							Computed: true,
							Type:     schema.TypeInt,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"node_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"node_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"security_group_probe_ip_list": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Computed: true,
							Type:     schema.TypeInt,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudThreatDetectionHoneypotNodesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("node_name"); ok {
		request["NodeName"] = v
	}
	if v, ok := d.GetOk("page_number"); ok {
		request["CurrentPage"] = v
	} else {
		request["CurrentPage"] = 1
	}
	if v, ok := d.GetOk("node_id"); ok {
		request["NodeId"] = v
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}

	var honeypotNodesNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		honeypotNodesNameRegex = r
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

	conn, err := client.NewThreatdetectionClient()
	if err != nil {
		return WrapError(err)
	}
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListHoneypotNode"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-12-03"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_threat_detection_honeypot_nodes", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.HoneypotNodeList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.HoneypotNodeList", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if honeypotNodesNameRegex != nil && !honeypotNodesNameRegex.MatchString(fmt.Sprint(item["NodeName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["NodeId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                             fmt.Sprint(object["NodeId"]),
			"allow_honeypot_access_internet": object["AllowHoneypotAccessInternet"],
			"available_probe_num":            object["ProbeTotalCount"],
			"node_id":                        object["NodeId"],
			"node_name":                      object["NodeName"],
			"security_group_probe_ip_list":   object["SecurityGroupProbeIpList"].([]interface{}),
			"status":                         object["TotalStatus"],
			"create_time":                    object["CreateTime"],
		}

		ids = append(ids, fmt.Sprint(object["NodeId"]))
		names = append(names, object["NodeName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("nodes", s); err != nil {
		return WrapError(err)
	}
	return nil
}
