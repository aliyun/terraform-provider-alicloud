// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudKvcachestoreKvCacheStores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudKvcachestoreKvCacheStoreRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"kvcs_ids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"stores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hpn_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kvcs_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
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

func dataSourceAliCloudKvcachestoreKvCacheStoreRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "ListKVCacheStores"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("kvcs_ids"); ok {
		request["KvcsIds"] = v
	}
	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Kvcachestore", "2026-06-17", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.KVCacheStores[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["KvcsId"])]; !ok {
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

		mapping["id"] = objectRaw["KvcsId"]

		mapping["capacity"] = objectRaw["Capacity"]
		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["description"] = objectRaw["Description"]
		mapping["hpn_zone"] = objectRaw["HpnZone"]
		mapping["name"] = objectRaw["Name"]
		mapping["payment_type"] = objectRaw["PaymentType"]
		mapping["region_id"] = objectRaw["RegionId"]
		mapping["zone_id"] = objectRaw["ZoneId"]
		mapping["kvcs_id"] = objectRaw["KvcsId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["KvcsId"])
		mapping, err = dataSourceAliCloudKvcachestoreKvCacheStoreReadDescription(d, id, mapping, meta)
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

	if err := d.Set("stores", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudKvcachestoreKvCacheStoreReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	kvcachestoreServiceV2 := KvcachestoreServiceV2{client}
	getResp, err := kvcachestoreServiceV2.DescribeKvcachestoreKvCacheStore(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["capacity"] = objectRaw["Capacity"]
	mapping["create_time"] = objectRaw["CreateTime"]
	mapping["description"] = objectRaw["Description"]
	mapping["hpn_zone"] = objectRaw["HpnZone"]
	mapping["name"] = objectRaw["Name"]
	mapping["payment_type"] = objectRaw["PaymentType"]
	mapping["region_id"] = objectRaw["RegionId"]
	mapping["status"] = objectRaw["Status"]
	mapping["zone_id"] = objectRaw["ZoneId"]
	mapping["kvcs_id"] = objectRaw["KvcsId"]

	return mapping, nil
}
