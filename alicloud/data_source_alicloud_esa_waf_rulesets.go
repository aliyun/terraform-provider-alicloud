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

func dataSourceAliCloudEsaWafRuleSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaWafRuleSetRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"phase": {
				Type:     schema.TypeString,
				Required: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"on", "off"}, true),
			},
			"query_args": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"any_like": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name_like": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"order_by": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"desc": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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
			"sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ruleset_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaWafRuleSetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListWafRulesets"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	request["SiteId"] = d.Get("site_id")
	request["Phase"] = d.Get("phase")
	request["SiteVersion"] = d.Get("site_version")

	if v, ok := d.GetOk("query_args"); ok {
		queryArgsMap := map[string]interface{}{}
		for _, queryArgsList := range v.([]interface{}) {
			queryArgsArg := queryArgsList.(map[string]interface{})

			if anyLike, ok := queryArgsArg["any_like"]; ok {
				queryArgsMap["AnyLike"] = anyLike
			}

			if nameLike, ok := queryArgsArg["name_like"]; ok {
				queryArgsMap["NameLike"] = nameLike
			}

			if orderBy, ok := queryArgsArg["order_by"]; ok {
				queryArgsMap["OrderBy"] = orderBy
			}

			if desc, ok := d.GetOkExists("query_args.0.desc"); ok {
				queryArgsMap["Desc"] = desc
			}
		}

		request["QueryArgs"] = convertObjectToJsonString(queryArgsMap)
	}

	status, statusOk := d.GetOk("status")

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

	var wafRuleSetNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		wafRuleSetNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Site.ServiceBusy", "TooManyRequests"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_waf_rulesets", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Rulesets", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Rulesets", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", item["Id"], request["SiteId"])]; !ok {
					continue
				}
			}

			if wafRuleSetNameRegex != nil {
				if !wafRuleSetNameRegex.MatchString(fmt.Sprint(item["Name"])) {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
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
			"id":          fmt.Sprintf("%v:%v", object["Id"], request["SiteId"]),
			"ruleset_id":  fmt.Sprint(object["Id"]),
			"phase":       object["Phase"],
			"name":        object["Name"],
			"target":      object["Target"],
			"status":      object["Status"],
			"update_time": object["UpdateTime"],
			"types":       object["Types"],
			"fields":      object["Fields"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("sets", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
