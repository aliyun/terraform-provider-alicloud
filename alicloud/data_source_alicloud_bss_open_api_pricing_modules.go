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

func dataSourceAlicloudBssOpenApiPricingModules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudBssOpenApiPricingModulesRead,
		Schema: map[string]*schema.Schema{
			"product_code": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"product_type": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"subscription_type": {
				Required: true,
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
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"modules": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"pricing_module_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_code": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"product_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"subscription_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"unit": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"values": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"remark": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"value": {
										Computed: true,
										Type:     schema.TypeString,
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

func dataSourceAlicloudBssOpenApiPricingModulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})

	if v, ok := d.GetOk("product_code"); ok {
		request["ProductCode"] = v
	}
	if v, ok := d.GetOk("product_type"); ok {
		request["ProductType"] = v
	}
	if v, ok := d.GetOk("subscription_type"); ok {
		request["SubscriptionType"] = v
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

	var pricingModuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		pricingModuleNameRegex = r
	}

	var objects []interface{}
	var response map[string]interface{}
	var endpoint string
	var err error
	action := "DescribePricingModule"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPostWithEndpoint("BssOpenApi", "2017-12-14", action, nil, request, true, endpoint)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if !client.IsInternationalAccount() && IsExpectedErrors(err, []string{"NotApplicable"}) {
				endpoint = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_bss_open_api_pricing_modules", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Data.AttributeList.Attribute", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.AttributeList.Attribute", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["ProductCode"], ":", request["ProductType"], ":", request["SubscriptionType"], ":", item["Code"])]; !ok {
				continue
			}
		}

		if pricingModuleNameRegex != nil && !pricingModuleNameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                  fmt.Sprint(request["ProductCode"], ":", request["ProductType"], ":", request["SubscriptionType"], ":", object["Code"]),
			"code":                object["Code"],
			"pricing_module_name": object["Name"],
			"unit":                object["Unit"],
			"product_code":        request["ProductCode"],
			"product_type":        request["ProductType"],
			"subscription_type":   request["ProductType"],
		}

		values91Maps := make([]map[string]interface{}, 0)
		values91Raw, _ := jsonpath.Get("$.Values.AttributeValue", object)
		for _, value0 := range values91Raw.([]interface{}) {
			values91 := value0.(map[string]interface{})
			values91Map := make(map[string]interface{})
			values91Map["name"] = values91["Name"]
			values91Map["remark"] = values91["Remark"]
			values91Map["type"] = values91["Type"]
			values91Map["value"] = values91["Value"]
			values91Maps = append(values91Maps, values91Map)
		}
		mapping["values"] = values91Maps

		ids = append(ids, fmt.Sprint(request["ProductCode"], ":", request["ProductType"], ":", request["SubscriptionType"], ":", object["Code"]))
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

	if err := d.Set("modules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
