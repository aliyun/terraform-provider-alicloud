// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudEsaTransportLayerApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaTransportLayerApplicationRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"record_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cross_border_optimization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_access_rule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keep_alive_protection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"comment": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"edge_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"source_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_port": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_ip_pass_through_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"rules_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"site_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"static_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudEsaTransportLayerApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "ListTransportLayerApplications"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("site_id"); ok {
		query["SiteId"] = v
	}
	if v, ok := d.GetOk("match_type"); ok {
		query["MatchType"] = v.(string)
	}

	if v, ok := d.GetOk("record_name"); ok {
		query["RecordName"] = v.(string)
	}

	if v, ok := d.GetOk("site_id"); ok {
		query["SiteId"] = v.(string)
	}

	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcGet("ESA", "2024-09-10", action, query, request)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.Applications[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SiteId"], ":", item["ApplicationId"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(objectRaw["SiteId"], ":", objectRaw["ApplicationId"])

		mapping["cname"] = objectRaw["Cname"]
		mapping["cross_border_optimization"] = objectRaw["CrossBorderOptimization"]
		mapping["ip_access_rule"] = objectRaw["IpAccessRule"]
		mapping["ipv6"] = objectRaw["Ipv6"]
		mapping["keep_alive_protection"] = objectRaw["KeepAliveProtection"]
		mapping["record_name"] = objectRaw["RecordName"]
		mapping["rules_count"] = objectRaw["RulesCount"]
		mapping["static_ip"] = objectRaw["StaticIp"]
		mapping["application_id"] = objectRaw["ApplicationId"]
		if v, ok := objectRaw["SiteId"]; ok {
			mapping["site_id"] = v
		}

		rulesRaw := objectRaw["Rules"]
		rulesMaps := make([]map[string]interface{}, 0)
		if rulesRaw != nil {
			for _, rulesChildRaw := range convertToInterfaceArray(rulesRaw) {
				rulesMap := make(map[string]interface{})
				rulesChildRaw := rulesChildRaw.(map[string]interface{})
				rulesMap["client_ip_pass_through_mode"] = rulesChildRaw["ClientIPPassThroughMode"]
				rulesMap["comment"] = rulesChildRaw["Comment"]
				rulesMap["edge_port"] = rulesChildRaw["EdgePort"]
				rulesMap["protocol"] = rulesChildRaw["Protocol"]
				rulesMap["rule_id"] = rulesChildRaw["RuleId"]
				rulesMap["source"] = rulesChildRaw["Source"]
				rulesMap["source_port"] = rulesChildRaw["SourcePort"]
				rulesMap["source_type"] = rulesChildRaw["SourceType"]

				rulesMaps = append(rulesMaps, rulesMap)
			}
		}
		mapping["rules"] = rulesMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["SiteId"], ":", objectRaw["ApplicationId"])
		mapping, err = dataSourceAliCloudEsaTransportLayerApplicationReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("applications", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudEsaTransportLayerApplicationReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	esaServiceV2 := EsaServiceV2{client}
	getResp, err := esaServiceV2.DescribeEsaTransportLayerApplication(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["cname"] = objectRaw["Cname"]
	mapping["cross_border_optimization"] = objectRaw["CrossBorderOptimization"]
	mapping["ip_access_rule"] = objectRaw["IpAccessRule"]
	mapping["ipv6"] = objectRaw["Ipv6"]
	mapping["keep_alive_protection"] = objectRaw["KeepAliveProtection"]
	mapping["record_name"] = objectRaw["RecordName"]
	mapping["rules_count"] = objectRaw["RulesCount"]
	mapping["static_ip"] = objectRaw["StaticIp"]
	mapping["status"] = objectRaw["Status"]
	mapping["application_id"] = objectRaw["ApplicationId"]
	if v, ok := objectRaw["SiteId"]; ok {
		mapping["site_id"] = v
	}

	rulesRaw := objectRaw["Rules"]
	rulesMaps := make([]map[string]interface{}, 0)
	if rulesRaw != nil {
		for _, rulesChildRaw := range convertToInterfaceArray(rulesRaw) {
			rulesMap := make(map[string]interface{})
			rulesChildRaw := rulesChildRaw.(map[string]interface{})
			rulesMap["client_ip_pass_through_mode"] = rulesChildRaw["ClientIPPassThroughMode"]
			rulesMap["comment"] = rulesChildRaw["Comment"]
			rulesMap["edge_port"] = rulesChildRaw["EdgePort"]
			rulesMap["protocol"] = rulesChildRaw["Protocol"]
			rulesMap["rule_id"] = rulesChildRaw["RuleId"]
			rulesMap["source"] = rulesChildRaw["Source"]
			rulesMap["source_port"] = rulesChildRaw["SourcePort"]
			rulesMap["source_type"] = rulesChildRaw["SourceType"]

			rulesMaps = append(rulesMaps, rulesMap)
		}
	}
	mapping["rules"] = rulesMaps

	return mapping, nil
}
