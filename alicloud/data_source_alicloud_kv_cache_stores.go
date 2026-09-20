// Package alicloud. Hand-written data source for KVCacheStore List operation
// (Kvcachestore 2026-06-17 ListKVCacheStores).
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudKvCacheStores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKvCacheStoresRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"output_file": {
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
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKvCacheStoresRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{}
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

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []interface{}
	var response map[string]interface{}
	pageNumber := 1
	pageSize := 50
	request["PageSize"] = pageSize

	for {
		request["PageNumber"] = pageNumber
		action := "ListKVCacheStores"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Kvcachestore", "2026-06-17", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kv_cache_stores", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.KVCacheStores", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.KVCacheStores", response)
		}
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

		pageTotal := 0
		if v, ok := response["PageTotal"]; ok {
			pageTotal, _ = strconv.Atoi(fmt.Sprint(v))
		}
		if pageNumber >= pageTotal || pageTotal == 0 {
			break
		}
		pageNumber++
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"capacity":          object["Capacity"],
			"create_time":       object["CreateTime"],
			"description":       object["Description"],
			"hpn_zone":          object["HpnZone"],
			"kvcs_id":           fmt.Sprint(object["KvcsId"]),
			"name":              object["Name"],
			"payment_type":      object["PaymentType"],
			"region_id":         object["RegionId"],
			"resource_group_id": object["ResourceGroupId"],
			"status":            object["Status"],
			"tags":              tagsToMap(object["Tags"]),
			"zone_id":           object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(object["KvcsId"]))
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
