// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCloudFirewallVpcCenTrFirewalls() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudFirewallVpcCenTrFirewallRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"current_page": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"firewall_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewall_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewall_switch_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"region_no": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"route_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"transit_router_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewalls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cen_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cen_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"firewall_switch_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ips_config": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_all_patch": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"basic_rules": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"run_mode": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"precheck_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_no": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"result_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transit_router_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudCloudFirewallVpcCenTrFirewallRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DescribeTrFirewallsV2List"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if v, ok := d.GetOk("firewall_id"); ok {
		request["FirewallId"] = v
	}
	request["CenId"] = d.Get("cen_id")
	if v, ok := d.GetOkExists("current_page"); ok {
		request["CurrentPage"] = v
	}
	if v, ok := d.GetOk("firewall_id"); ok {
		request["FirewallId"] = v
	}
	request["FirewallName"] = d.Get("firewall_name")
	if v, ok := d.GetOk("firewall_switch_status"); ok {
		request["FirewallSwitchStatus"] = v
	}
	request["RegionNo"] = d.Get("region_no")
	request["RouteMode"] = d.Get("route_mode")
	request["TransitRouterId"] = d.Get("transit_router_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cloudfw", "2017-12-07", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.VpcTrFirewalls[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["FirewallId"])]; !ok {
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
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["FirewallId"]

		mapping["cen_id"] = objectRaw["CenId"]
		mapping["firewall_name"] = objectRaw["VpcFirewallName"]
		mapping["firewall_switch_status"] = objectRaw["FirewallSwitchStatus"]
		mapping["region_no"] = objectRaw["RegionNo"]
		mapping["route_mode"] = objectRaw["RouteMode"]
		mapping["transit_router_id"] = objectRaw["TransitRouterId"]
		mapping["cen_name"] = objectRaw["CenName"]
		mapping["firewall_id"] = objectRaw["FirewallId"]
		mapping["precheck_status"] = objectRaw["PrecheckStatus"]
		mapping["region_status"] = objectRaw["RegionStatus"]
		mapping["result_code"] = objectRaw["ResultCode"]

		ipsConfigMaps := make([]map[string]interface{}, 0)
		ipsConfigMap := make(map[string]interface{})
		ipsConfigRaw := make(map[string]interface{})
		if objectRaw["IpsConfig"] != nil {
			ipsConfigRaw = objectRaw["IpsConfig"].(map[string]interface{})
		}
		if len(ipsConfigRaw) > 0 {
			ipsConfigMap["basic_rules"] = ipsConfigRaw["BasicRules"]
			ipsConfigMap["enable_all_patch"] = ipsConfigRaw["EnableAllPatch"]
			ipsConfigMap["run_mode"] = ipsConfigRaw["RunMode"]

			ipsConfigMaps = append(ipsConfigMaps, ipsConfigMap)
		}
		mapping["ips_config"] = ipsConfigMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw[""])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("firewalls", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
