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

func dataSourceAliCloudEnsNatGatewayForwardEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsNatGatewayForwardEntryRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forward_entry_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forward_entry_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_entry_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internal_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
		},
	}
}

func dataSourceAliCloudEnsNatGatewayForwardEntryRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DescribeForwardTableEntries"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("forward_entry_id"); ok {
		request["ForwardEntryId"] = v
	}
	request["ExternalIp"] = d.Get("external_ip")
	if v, ok := d.GetOk("forward_entry_id"); ok {
		request["ForwardEntryId"] = v
	}
	if v, ok := d.GetOk("forward_entry_name"); ok {
		request["ForwardEntryName"] = v
	}
	request["InternalIp"] = d.Get("internal_ip")
	if v, ok := d.GetOk("ip_protocol"); ok {
		request["IpProtocol"] = v
	}
	request["NatGatewayId"] = d.Get("nat_gateway_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.ForwardTableEntries[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ForwardEntryId"])]; !ok {
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

		mapping["id"] = objectRaw["ForwardEntryId"]

		mapping["external_ip"] = objectRaw["ExternalIp"]
		mapping["external_port"] = objectRaw["ExternalPort"]
		mapping["forward_entry_name"] = objectRaw["ForwardEntryName"]
		mapping["health_check_port"] = formatInt(objectRaw["HealthCheckPort"])
		mapping["internal_ip"] = objectRaw["InternalIp"]
		mapping["internal_port"] = objectRaw["InternalPort"]
		mapping["ip_protocol"] = objectRaw["IpProtocol"]
		mapping["nat_gateway_id"] = objectRaw["NatGatewayId"]
		mapping["status"] = objectRaw["Status"]
		mapping["forward_entry_id"] = objectRaw["ForwardEntryId"]

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
