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

func dataSourceAliCloudDirectMailConfigSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudDirectMailConfigSetRead,
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

			"sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_pool_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_public_channel_backoff": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"validation_option": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"forbidden_sub_status_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"forbidden_status_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudDirectMailConfigSetRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "ConfigSetList"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	if v, ok := d.GetOkExists("all"); ok {
		request["All"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageIndex"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Dm", "2015-11-23", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.ConfigSets[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
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

		mapping["description"] = objectRaw["Description"]
		mapping["is_public_channel_backoff"] = objectRaw["IsPublicChannelBackoff"]
		mapping["name"] = objectRaw["Name"]

		ipPoolRawObj, _ := jsonpath.Get("$.IpPool", objectRaw)
		ipPoolRaw := make(map[string]interface{})
		if ipPoolRawObj != nil {
			ipPoolRaw = ipPoolRawObj.(map[string]interface{})
		}
		mapping["ip_pool_id"] = ipPoolRaw["IpPoolId"]
		mapping["ip_pool_name"] = ipPoolRaw["IpPoolName"]

		validationOptionMaps := make([]map[string]interface{}, 0)
		validationOptionMap := make(map[string]interface{})
		validationOptionRaw := make(map[string]interface{})
		if objectRaw["ValidationOption"] != nil {
			validationOptionRaw = objectRaw["ValidationOption"].(map[string]interface{})
		}
		if len(validationOptionRaw) > 0 {
			validationOptionMap["enabled"] = validationOptionRaw["Enabled"]

			forbiddenStatusListRaw := make([]interface{}, 0)
			if validationOptionRaw["ForbiddenStatusList"] != nil {
				forbiddenStatusListRaw = convertToInterfaceArray(validationOptionRaw["ForbiddenStatusList"])
			}
			validationOptionMap["forbidden_status_list"] = forbiddenStatusListRaw

			forbiddenSubStatusListRaw := make([]interface{}, 0)
			if validationOptionRaw["ForbiddenSubStatusList"] != nil {
				forbiddenSubStatusListRaw = convertToInterfaceArray(validationOptionRaw["ForbiddenSubStatusList"])
			}
			validationOptionMap["forbidden_sub_status_list"] = forbiddenSubStatusListRaw

			validationOptionMaps = append(validationOptionMaps, validationOptionMap)
		}
		mapping["validation_option"] = validationOptionMaps

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["Id"])
		mapping, err = dataSourceAliCloudDirectMailConfigSetReadDescription(d, id, mapping, meta)
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

	if err := d.Set("sets", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudDirectMailConfigSetReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	directMailServiceV2 := DirectMailServiceV2{client}
	getResp, err := directMailServiceV2.DescribeDirectMailConfigSet(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	detailRawObj, _ := jsonpath.Get("$.Detail", objectRaw)
	detailRaw := make(map[string]interface{})
	if detailRawObj != nil {
		detailRaw = detailRawObj.(map[string]interface{})
	}
	mapping["description"] = detailRaw["Description"]
	mapping["is_public_channel_backoff"] = detailRaw["IsPublicChannelBackoff"]
	mapping["name"] = detailRaw["Name"]
	mapping["id"] = detailRaw["Id"]

	ipPoolRawObj, _ := jsonpath.Get("$.Detail.IpPool", objectRaw)
	ipPoolRaw := make(map[string]interface{})
	if ipPoolRawObj != nil {
		ipPoolRaw = ipPoolRawObj.(map[string]interface{})
	}
	mapping["ip_pool_id"] = ipPoolRaw["IpPoolId"]
	mapping["ip_pool_name"] = ipPoolRaw["IpPoolName"]

	validationOptionMaps := make([]map[string]interface{}, 0)
	validationOptionMap := make(map[string]interface{})
	validationOptionRaw := make(map[string]interface{})
	if detailRaw["ValidationOption"] != nil {
		validationOptionRaw = detailRaw["ValidationOption"].(map[string]interface{})
	}
	if len(validationOptionRaw) > 0 {
		validationOptionMap["enabled"] = validationOptionRaw["Enabled"]

		forbiddenStatusListRaw := make([]interface{}, 0)
		if validationOptionRaw["ForbiddenStatusList"] != nil {
			forbiddenStatusListRaw = convertToInterfaceArray(validationOptionRaw["ForbiddenStatusList"])
		}

		validationOptionMap["forbidden_status_list"] = forbiddenStatusListRaw
		forbiddenSubStatusListRaw := make([]interface{}, 0)
		if validationOptionRaw["ForbiddenSubStatusList"] != nil {
			forbiddenSubStatusListRaw = convertToInterfaceArray(validationOptionRaw["ForbiddenSubStatusList"])
		}

		validationOptionMap["forbidden_sub_status_list"] = forbiddenSubStatusListRaw
		validationOptionMaps = append(validationOptionMaps, validationOptionMap)
	}
	mapping["validation_option"] = validationOptionMaps

	return mapping, nil
}
