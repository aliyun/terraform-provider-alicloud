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

func dataSourceAliCloudEsaHttpsBasicConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaHttpsBasicConfigurationRead,
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"configurations": {
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
						"ciphersuite": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ciphersuite_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ocsp_stapling": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http2": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"http3": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"https": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls10": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls11": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls12": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls13": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaHttpsBasicConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListHttpsBasicConfigurations"
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

	var httpsBasicConfigurationNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		httpsBasicConfigurationNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_https_basic_configurations", action, AlibabaCloudSdkGoERROR)
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

			if httpsBasicConfigurationNameRegex != nil {
				if !httpsBasicConfigurationNameRegex.MatchString(fmt.Sprint(item["RuleName"])) {
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
			"id":                fmt.Sprintf("%v:%v", request["SiteId"], object["ConfigId"]),
			"config_id":         fmt.Sprint(object["ConfigId"]),
			"config_type":       object["ConfigType"],
			"rule_enable":       object["RuleEnable"],
			"rule_name":         object["RuleName"],
			"rule":              object["Rule"],
			"sequence":          object["Sequence"],
			"ciphersuite":       object["Ciphersuite"],
			"ciphersuite_group": object["CiphersuiteGroup"],
			"ocsp_stapling":     object["OcspStapling"],
			"http2":             object["Http2"],
			"http3":             object["Http3"],
			"https":             object["Https"],
			"tls10":             object["Tls10"],
			"tls11":             object["Tls11"],
			"tls12":             object["Tls12"],
			"tls13":             object["Tls13"],
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

	if err := d.Set("configurations", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
