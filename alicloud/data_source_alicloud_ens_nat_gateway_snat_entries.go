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

func dataSourceAliCloudEnsNatGatewaySnatEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsNatGatewaySnatEntryRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snat_entry_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snat_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snat_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_affinity": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"idle_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"isp_affinity": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_entry_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standby_snat_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standby_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
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
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudEnsNatGatewaySnatEntryRead(d *schema.ResourceData, meta interface{}) error {
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

	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeSnatTableEntries"
	var err error
	query = make(map[string]interface{})

	if v, ok := d.GetOk("snat_entry_id"); ok {
		query["SnatEntryId"] = v
	}
	if v, ok := d.GetOk("nat_gateway_id"); ok {
		query["NatGatewayId"] = v.(string)
	}

	if v, ok := d.GetOk("snat_entry_id"); ok {
		query["SnatEntryId"] = v.(string)
	}

	if v, ok := d.GetOk("snat_entry_name"); ok {
		query["SnatEntryName"] = v.(string)
	}

	if v, ok := d.GetOk("snat_ip"); ok {
		query["SnatIp"] = v.(string)
	}

	if v, ok := d.GetOk("source_cidr"); ok {
		query["SourceCIDR"] = v.(string)
	}

	query["PageSize"] = PageSizeLarge
	query["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcGet("Ens", "2017-11-10", action, query, nil)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, query)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.SnatTableEntries[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SnatEntryId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		query["PageNumber"] = query["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["SnatEntryId"]

		mapping["eip_affinity"] = objectRaw["EipAffinity"]
		mapping["idle_timeout"] = objectRaw["IdleTimeout"]
		mapping["isp_affinity"] = objectRaw["IspAffinity"]
		mapping["nat_gateway_id"] = objectRaw["NatGatewayId"]
		mapping["snat_entry_name"] = objectRaw["SnatEntryName"]
		mapping["snat_ip"] = objectRaw["SnatIp"]
		mapping["source_cidr"] = objectRaw["SourceCIDR"]
		mapping["standby_snat_ip"] = objectRaw["StandbySnatIp"]
		mapping["standby_status"] = objectRaw["StandbyStatus"]
		mapping["status"] = objectRaw["Status"]
		mapping["snat_entry_id"] = objectRaw["SnatEntryId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["SnatEntryId"])
		mapping, err = dataSourceAliCloudEnsNatGatewaySnatEntryReadDescription(d, id, mapping, meta)
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

	if err := d.Set("entries", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudEnsNatGatewaySnatEntryReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	ensServiceV2 := EnsServiceV2{client}
	getResp, err := ensServiceV2.DescribeEnsNatGatewaySnatEntry(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["creation_time"] = objectRaw["CreationTime"]
	mapping["dest_cidr"] = objectRaw["DestCIDR"]
	mapping["eip_affinity"] = objectRaw["EipAffinity"]
	mapping["idle_timeout"] = objectRaw["IdleTimeout"]
	mapping["isp_affinity"] = objectRaw["IspAffinity"]
	mapping["nat_gateway_id"] = objectRaw["NatGatewayId"]
	mapping["snat_entry_name"] = objectRaw["SnatEntryName"]
	mapping["snat_ip"] = objectRaw["SnatIp"]
	mapping["source_cidr"] = objectRaw["SourceCIDR"]
	mapping["standby_snat_ip"] = objectRaw["StandbySnatIp"]
	mapping["standby_status"] = objectRaw["StandbyStatus"]
	mapping["status"] = objectRaw["Status"]
	mapping["type"] = objectRaw["Type"]
	mapping["snat_entry_id"] = objectRaw["SnatEntryId"]

	return mapping, nil
}
