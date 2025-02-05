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

func dataSourceAlicloudAlbAscripts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlbAscriptsRead,
		Schema: map[string]*schema.Schema{
			"ascript_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"listener_id": {
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
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enable_details": {
				Optional: true,
				Type:     schema.TypeBool,
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"ascripts": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ascript_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ascript_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"ext_attribute_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"ext_attributes": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute_key": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"attribute_value": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"listener_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"load_balancer_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"position": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"script_content": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlbAscriptsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("ascript_name"); ok {
		request["AScriptNames.1"] = v
	}
	if v, ok := d.GetOk("listener_id"); ok {
		request["ListenerIds.1"] = v
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var ascriptNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		ascriptNameRegex = r
	}
	var objects []interface{}
	var response map[string]interface{}
	var err error
	for {
		action := "ListAScripts"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alb_ascripts", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AScripts", response)
		if err != nil {
			resp = make([]interface{}, 0)
			// return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AScripts", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AScriptId"])]; !ok {
					continue
				}
			}

			if ascriptNameRegex != nil && !ascriptNameRegex.MatchString(fmt.Sprint(item["AScriptName"])) {
				continue
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                    fmt.Sprint(object["AScriptId"]),
			"ascript_id":            object["AScriptId"],
			"ascript_name":          object["AScriptName"],
			"enabled":               object["Enabled"],
			"ext_attribute_enabled": object["ExtAttributeEnabled"],
			"listener_id":           object["ListenerId"],
			"position":              object["Position"],
			"script_content":        object["ScriptContent"],
			"status":                object["AScriptStatus"],
		}

		extAttributes81Maps := make([]map[string]interface{}, 0)
		extAttributes81Raw := object["ExtAttributes"]
		for _, value0 := range extAttributes81Raw.([]interface{}) {
			extAttributes81 := value0.(map[string]interface{})
			extAttributes81Map := make(map[string]interface{})
			extAttributes81Map["attribute_key"] = extAttributes81["AttributeKey"]
			extAttributes81Map["attribute_value"] = extAttributes81["AttributeValue"]
			extAttributes81Maps = append(extAttributes81Maps, extAttributes81Map)
		}
		mapping["ext_attributes"] = extAttributes81Maps

		ids = append(ids, fmt.Sprint(object["AScriptId"]))
		names = append(names, object["AScriptName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ascripts", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
