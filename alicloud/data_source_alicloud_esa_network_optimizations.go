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

func dataSourceAliCloudEsaNetworkOptimizations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaNetworkOptimizationRead,
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
			"optimizations": {
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
						"smart_routing": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grpc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http2_origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"websocket": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"upload_max_filesize": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaNetworkOptimizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListNetworkOptimizations"
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

	var networkOptimizationNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		networkOptimizationNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_network_optimizations", action, AlibabaCloudSdkGoERROR)
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

			if networkOptimizationNameRegex != nil {
				if !networkOptimizationNameRegex.MatchString(fmt.Sprint(item["RuleName"])) {
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
			"id":                  fmt.Sprintf("%v:%v", request["SiteId"], object["ConfigId"]),
			"config_id":           fmt.Sprint(object["ConfigId"]),
			"config_type":         object["ConfigType"],
			"site_version":        object["SiteVersion"],
			"rule_enable":         object["RuleEnable"],
			"rule_name":           object["RuleName"],
			"rule":                object["Rule"],
			"sequence":            object["Sequence"],
			"smart_routing":       object["SmartRouting"],
			"grpc":                object["Grpc"],
			"http2_origin":        object["Http2Origin"],
			"websocket":           object["Websocket"],
			"upload_max_filesize": object["UploadMaxFilesize"],
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

	if err := d.Set("optimizations", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
