// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudWafv3AddressBooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudWafv3AddressBookRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"books": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_book_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_book_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_book_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_list": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"description": {
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

func dataSourceAliCloudWafv3AddressBookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
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

	action := "DescribeDefenseRules"
	queryFilter, _ := json.Marshal(map[string]interface{}{"scene": "address_book"})
	pageNumber := 1
	pageSize := PageSizeLarge
	instanceId := d.Get("instance_id")

	type bookEntry struct {
		composedId string
		ruleId     string
		config     map[string]interface{}
	}
	entries := make([]bookEntry, 0)

	for {
		request := map[string]interface{}{
			"RegionId":    client.RegionId,
			"InstanceId":  instanceId,
			"DefenseType": "global",
			"Query":       string(queryFilter),
			"PageNumber":  pageNumber,
			"PageSize":    pageSize,
		}

		var response map[string]interface{}
		var err error
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("waf-openapi", "2021-10-01", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_wafv3_address_books", action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.Rules", response)
		result, _ := resp.([]interface{})

		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			configRaw := decodeWafv3AddressBookConfig(item["Config"])
			bookName := ""
			if configRaw != nil {
				if n, ok := configRaw["name"].(string); ok {
					bookName = n
				}
			}
			if nameRegex != nil && !nameRegex.MatchString(bookName) {
				continue
			}

			composedId := fmt.Sprintf("%v:%v", instanceId, item["RuleId"])
			if len(idsMap) > 0 {
				if _, ok := idsMap[composedId]; !ok {
					continue
				}
			}

			entries = append(entries, bookEntry{
				composedId: composedId,
				ruleId:     fmt.Sprint(item["RuleId"]),
				config:     configRaw,
			})
		}

		if len(result) < pageSize {
			break
		}
		pageNumber++
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	enableDetails := d.Get("enable_details").(bool)
	wafv3ServiceV2 := Wafv3ServiceV2{client}

	for _, entry := range entries {
		mapping := map[string]interface{}{
			"id":              entry.composedId,
			"address_book_id": entry.ruleId,
		}

		bookName := ""
		bookDesc := ""
		bookType := ""
		if entry.config != nil {
			if n, ok := entry.config["name"].(string); ok {
				bookName = n
			}
			if dRaw, ok := entry.config["description"].(string); ok {
				bookDesc = dRaw
			}
			if t, ok := entry.config["valueType"].(string); ok {
				bookType = t
			}
		}
		mapping["address_book_name"] = bookName
		mapping["address_book_type"] = bookType
		mapping["description"] = bookDesc

		if enableDetails {
			addressesRaw, err := wafv3ServiceV2.DescribeWafv3AddressBookAddresses(entry.composedId)
			if err != nil && !NotFoundError(err) {
				return WrapError(err)
			}
			addressList := make([]interface{}, 0)
			if addrListRawObj, _ := jsonpath.Get("$.AddressList", addressesRaw); addrListRawObj != nil {
				if addrItems, ok := addrListRawObj.([]interface{}); ok {
					for _, addrItem := range addrItems {
						if addrEntry, ok := addrItem.(map[string]interface{}); ok {
							if addr, ok := addrEntry["Address"]; ok && addr != nil {
								addressList = append(addressList, addr)
							}
						}
					}
				}
			}
			mapping["address_list"] = addressList
		}

		ids = append(ids, entry.composedId)
		names = append(names, bookName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("books", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func decodeWafv3AddressBookConfig(raw interface{}) map[string]interface{} {
	switch v := raw.(type) {
	case map[string]interface{}:
		return v
	case string:
		if v == "" {
			return nil
		}
		out := make(map[string]interface{})
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return nil
		}
		return out
	}
	return nil
}
