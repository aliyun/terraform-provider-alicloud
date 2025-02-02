package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpcPrefixLists() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcPrefixListsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"prefix_list_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lists": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entrys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_entries": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_list_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix_list_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func dataSourceAlicloudVpcPrefixListsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPrefixLists"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("prefix_list_name"); ok {
		request["PrefixListName"] = v
	}
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var prefixListNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		prefixListNameRegex = r
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
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_prefix_lists", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.PrefixLists", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PrefixLists", response)
		}
		if result, ok := resp.([]interface{}); ok && len(result) > 0 {
			for _, v := range result {
				item := v.(map[string]interface{})
				if prefixListNameRegex != nil && !prefixListNameRegex.MatchString(fmt.Sprint(item["PrefixListName"])) {
					continue
				}
				if len(idsMap) > 0 {
					if _, ok := idsMap[fmt.Sprint(item["PrefixListId"])]; !ok {
						continue
					}
				}
				objects = append(objects, item)
			}
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"create_time":             object["CreationTime"],
			"ip_version":              object["IpVersion"],
			"max_entries":             formatInt(object["MaxEntries"]),
			"id":                      fmt.Sprint(object["PrefixListId"]),
			"prefix_list_id":          fmt.Sprint(object["PrefixListId"]),
			"prefix_list_name":        object["PrefixListName"],
			"prefix_list_description": object["PrefixListDescription"],
			"share_type":              object["ShareType"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["PrefixListName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["PrefixListId"])
		vpcService := VpcService{client}
		getResp, err := vpcService.GetVpcPrefixListEntries(id)
		if err != nil {
			return WrapError(err)
		}

		prefixListEntry := make([]map[string]interface{}, 0)
		for _, v := range getResp {
			temp1 := map[string]interface{}{
				"cidr":        v["Cidr"],
				"description": v["Description"],
			}
			prefixListEntry = append(prefixListEntry, temp1)
		}
		mapping["entrys"] = prefixListEntry
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("lists", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
