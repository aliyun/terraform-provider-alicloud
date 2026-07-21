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

func dataSourceAliCloudEsaOriginRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaOriginRuleRead,
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
			"config_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_version": {
				Type:     schema.TypeInt,
				Optional: true,
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
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"site_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_enable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sequence": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"origin_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_scheme": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_sni": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_http_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_https_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_read_timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_verify": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin_mtls": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"follow302_enable": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"follow302_max_tries": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"follow302_target_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"follow302_retain_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"follow302_retain_args": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"range_chunk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaOriginRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListOriginRules"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	request["SiteId"] = d.Get("site_id")

	if v, ok := d.GetOk("config_id"); ok {
		request["ConfigId"] = v
	}

	if v, ok := d.GetOk("config_type"); ok {
		request["ConfigType"] = v
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}

	if v, ok := d.GetOkExists("site_version"); ok {
		request["SiteVersion"] = v
	}

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

	var originRuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		originRuleNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcGet("ESA", "2024-09-10", action, request, nil)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_origin_rules", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Configs", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Configs", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%v:%v", request["SiteId"], item["ConfigId"])]; !ok {
					continue
				}
			}

			if originRuleNameRegex != nil {
				if !originRuleNameRegex.MatchString(fmt.Sprint(item["RuleName"])) {
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
			"id":                      fmt.Sprintf("%v:%v", request["SiteId"], object["ConfigId"]),
			"config_id":               fmt.Sprint(object["ConfigId"]),
			"config_type":             object["ConfigType"],
			"site_version":            object["SiteVersion"],
			"rule_enable":             object["RuleEnable"],
			"rule_name":               object["RuleName"],
			"rule":                    object["Rule"],
			"sequence":                object["Sequence"],
			"origin_host":             object["OriginHost"],
			"origin_scheme":           object["OriginScheme"],
			"origin_sni":              object["OriginSni"],
			"origin_http_port":        object["OriginHttpPort"],
			"origin_https_port":       object["OriginHttpsPort"],
			"origin_read_timeout":     object["OriginReadTimeout"],
			"dns_record":              object["DnsRecord"],
			"origin_verify":           object["OriginVerify"],
			"origin_mtls":             object["OriginMtls"],
			"follow302_enable":        object["Follow302Enable"],
			"follow302_max_tries":     object["Follow302MaxTries"],
			"follow302_target_host":   object["Follow302TargetHost"],
			"follow302_retain_header": object["Follow302RetainHeader"],
			"follow302_retain_args":   object["Follow302RetainArgs"],
			"range":                   object["Range"],
			"range_chunk_size":        object["RangeChunkSize"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["RuleName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
