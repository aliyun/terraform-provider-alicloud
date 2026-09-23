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

func dataSourceAliCloudThreatDetectionAttackPathWhitelists() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudThreatDetectionAttackPathWhitelistRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path_name_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"whitelist_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"whitelists": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attack_path_asset_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"asset_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"node_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vendor": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"asset_sub_type": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"attack_path_whitelist_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"whitelist_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"whitelist_type": {
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

func dataSourceAliCloudThreatDetectionAttackPathWhitelistRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "ListAttackPathWhitelist"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("path_name_desc"); ok {
		request["PathNameDesc"] = v
	}
	if v, ok := d.GetOk("path_type"); ok {
		request["PathType"] = v
	}
	if v, ok := d.GetOk("whitelist_name"); ok {
		request["WhitelistName"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Sas", "2018-12-03", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.AttackPathWhitelistList[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AttackPathWhitelistId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["AttackPathWhitelistId"]

		mapping["path_name"] = objectRaw["PathName"]
		mapping["path_type"] = objectRaw["PathType"]
		mapping["remark"] = objectRaw["Remark"]
		mapping["whitelist_name"] = objectRaw["WhitelistName"]
		mapping["whitelist_type"] = objectRaw["WhitelistType"]
		mapping["attack_path_whitelist_id"] = objectRaw["AttackPathWhitelistId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["AttackPathWhitelistId"])
		mapping, err = dataSourceAliCloudThreatDetectionAttackPathWhitelistReadDescription(d, id, mapping, meta)
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

	if err := d.Set("whitelists", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudThreatDetectionAttackPathWhitelistReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	threatDetectionServiceV2 := ThreatDetectionServiceV2{client}
	getResp, err := threatDetectionServiceV2.DescribeThreatDetectionAttackPathWhitelist(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["path_name"] = objectRaw["PathName"]
	mapping["path_type"] = objectRaw["PathType"]
	mapping["remark"] = objectRaw["Remark"]
	mapping["whitelist_name"] = objectRaw["WhitelistName"]
	mapping["whitelist_type"] = objectRaw["WhitelistType"]
	mapping["attack_path_whitelist_id"] = objectRaw["AttackPathWhitelistId"]

	attackPathAssetListRaw := objectRaw["AttackPathAssetList"]
	attackPathAssetListMaps := make([]map[string]interface{}, 0)
	if attackPathAssetListRaw != nil {
		for _, attackPathAssetListChildRaw := range convertToInterfaceArray(attackPathAssetListRaw) {
			attackPathAssetListMap := make(map[string]interface{})
			attackPathAssetListChildRaw := attackPathAssetListChildRaw.(map[string]interface{})
			attackPathAssetListMap["asset_sub_type"] = attackPathAssetListChildRaw["AssetSubType"]
			attackPathAssetListMap["asset_type"] = attackPathAssetListChildRaw["AssetType"]
			attackPathAssetListMap["instance_id"] = attackPathAssetListChildRaw["InstanceId"]
			attackPathAssetListMap["node_type"] = attackPathAssetListChildRaw["NodeType"]
			attackPathAssetListMap["region_id"] = attackPathAssetListChildRaw["RegionId"]
			attackPathAssetListMap["vendor"] = attackPathAssetListChildRaw["Vendor"]

			attackPathAssetListMaps = append(attackPathAssetListMaps, attackPathAssetListMap)
		}
	}
	mapping["attack_path_asset_list"] = attackPathAssetListMaps

	return mapping, nil
}
