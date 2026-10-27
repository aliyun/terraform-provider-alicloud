// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudDirectMailDedicatedIpPools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDirectMailDedicatedIpPoolRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudDirectMailDedicatedIpPoolRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DedicatedIpPoolList"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	// GetOkExists is unreliable for boolean zero values in plugin-sdk v2; use the
	// raw config to distinguish an explicit false from an omitted attribute.
	if rawConfig := d.GetRawConfig(); rawConfig != cty.NilVal && !rawConfig.IsNull() {
		if v, ok := rawConfig.AsValueMap()["all"]; ok && !v.IsNull() && v.IsKnown() {
			request["All"] = v.True()
		}
	}
	request["PageIndex"] = 1
	if v, ok := d.GetOk("pool_id"); ok {
		request["PoolId"] = v
	}
	request["PageSize"] = PageSizeLarge
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.IpPools[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if poolId, ok := d.GetOk("pool_id"); ok && fmt.Sprint(item["Id"]) != fmt.Sprint(poolId) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageIndex"] = request["PageIndex"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["Id"]

		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["ip_count"] = objectRaw["IpCount"]
		mapping["name"] = objectRaw["Name"]

		ipsRaw := objectRaw["Ips"]
		ipsMaps := make([]map[string]interface{}, 0)
		if ipsRaw != nil {
			for _, ipsChildRaw := range convertToInterfaceArray(ipsRaw) {
				ipsMap := make(map[string]interface{})
				ipsChildRaw := ipsChildRaw.(map[string]interface{})
				ipsMap["id"] = ipsChildRaw["Id"]
				ipsMap["ip"] = ipsChildRaw["Ip"]
				ipsMap["zone_id"] = ipsChildRaw["ZoneId"]

				ipsMaps = append(ipsMaps, ipsMap)
			}
		}
		mapping["ips"] = ipsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("pools", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
