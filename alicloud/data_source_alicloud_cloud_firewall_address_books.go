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

func dataSourceAliCloudCloudFirewallAddressBooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudFirewallAddressBooksRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ip", "ipv6", "domain", "port", "tag"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"books": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_add_tag_ecs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tag_relation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ecs_tags": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudCloudFirewallAddressBooksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeAddressBook"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1

	if v, ok := d.GetOk("group_type"); ok {
		request["GroupType"] = v
	}

	var objects []map[string]interface{}
	var addressBookNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		addressBookNameRegex = r
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
	var endpoint string

	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPostWithEndpoint("Cloudfw", "2017-12-07", action, nil, request, true, endpoint)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				} else if IsExpectedErrors(err, []string{"not buy user"}) {
					endpoint = connectivity.CloudFirewallOpenAPIEndpointControlPolicy
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}

			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_address_books", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Acls", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Acls", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if addressBookNameRegex != nil && !addressBookNameRegex.MatchString(fmt.Sprint(item["GroupName"])) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GroupUuid"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprint(object["GroupUuid"]),
			"group_uuid":       fmt.Sprint(object["GroupUuid"]),
			"group_name":       object["GroupName"],
			"group_type":       object["GroupType"],
			"description":      object["Description"],
			"auto_add_tag_ecs": formatInt(object["AutoAddTagEcs"]),
			"tag_relation":     object["TagRelation"],
			"address_list":     object["AddressList"],
		}

		ecsTags := make([]map[string]interface{}, 0)
		for _, tagListItem := range object["TagList"].([]interface{}) {
			ecsTagItem := make(map[string]interface{})
			ecsTagItem["tag_value"] = tagListItem.(map[string]interface{})["TagValue"]
			ecsTagItem["tag_key"] = tagListItem.(map[string]interface{})["TagKey"]
			ecsTags = append(ecsTags, ecsTagItem)
		}

		mapping["ecs_tags"] = ecsTags

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["GroupName"])
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
