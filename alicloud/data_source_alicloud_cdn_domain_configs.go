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

func dataSourceAliCloudCdnDomainConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCdnDomainConfigRead,
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
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"function_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"success", "testing", "failed", "configuring"}, false),
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
			"configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"function_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"function_args": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"arg_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"arg_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCdnDomainConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeCdnDomainConfigs"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId

	request["DomainName"] = d.Get("domain_name")

	if v, ok := d.GetOk("function_name"); ok {
		request["FunctionNames"] = v
	}

	if v, ok := d.GetOk("config_id"); ok {
		request["ConfigId"] = v
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

	var domainConfigNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		domainConfigNameRegex = r
	}

	var response map[string]interface{}
	var err error

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Cdn", "2018-05-10", action, nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"FlowControlError"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cdn_domain_configs", action, AlibabaCloudSdkGoERROR)
	}

	resp, err := jsonpath.Get("$.DomainConfigs.DomainConfig", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DomainConfigs.DomainConfig", response)
	}

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprintf("%v:%v:%v", request["DomainName"], item["FunctionName"], item["ConfigId"])]; !ok {
				continue
			}
		}

		if domainConfigNameRegex != nil {
			if !domainConfigNameRegex.MatchString(fmt.Sprint(item["FunctionName"])) {
				continue
			}
		}

		if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
			continue
		}

		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":            fmt.Sprintf("%v:%v:%v", request["DomainName"], object["FunctionName"], object["ConfigId"]),
			"function_name": fmt.Sprint(object["FunctionName"]),
			"config_id":     fmt.Sprint(object["ConfigId"]),
			"parent_id":     object["ParentId"],
			"status":        object["Status"],
		}

		if functionArgsMap, ok := object["FunctionArgs"].(map[string]interface{}); ok && functionArgsMap != nil {
			if functionArgList, ok := functionArgsMap["FunctionArg"]; ok && functionArgList != nil {
				functionArgMaps := make([]map[string]interface{}, 0)
				for _, functionArgLists := range functionArgList.([]interface{}) {
					functionArgMap := make(map[string]interface{})
					functionArg := functionArgLists.(map[string]interface{})

					if argName, ok := functionArg["ArgName"]; ok {
						functionArgMap["arg_name"] = argName
					}

					if argValue, ok := functionArg["ArgValue"]; ok {
						functionArgMap["arg_value"] = argValue
					}

					functionArgMaps = append(functionArgMaps, functionArgMap)
				}

				mapping["function_args"] = functionArgMaps
			}
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["FunctionName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("configs", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
