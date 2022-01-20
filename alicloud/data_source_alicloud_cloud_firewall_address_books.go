package alicloud

import (
	"fmt"
	"github.com/alibabacloud-go/tea-rpc/client"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudCloudFirewallAddressBooks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudFirewallAddressBooksRead,
		Schema: map[string]*schema.Schema{
			"contain_port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(AddressBookGroupTypes, true),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"books": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"auto_add_tag_ecs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tag_relation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_tags": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tag_value": {
										Type:     schema.TypeString,
										Optional: true,
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

func dataSourceAlicloudCloudFirewallAddressBooksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := make(map[string]interface{})
	if v, ok := d.GetOk("contain_port"); ok {
		request["ContainPort"] = v
	}

	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1

	idsMapExisted := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMapExisted[vv.(string)] = false
		}
	}

	var response map[string]interface{}
	var objects []map[string]interface{}
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("group_type"); ok {
		request["GroupType"] = v
		objects, err = describeAddressBook(response, conn, request, idsMapExisted, objects)
	} else {
		for _, v := range AddressBookGroupTypes {
			if len(idsMapExisted) > 0 {
				request["GroupType"] = v
				objects, err = describeAddressBook(response, conn, request, idsMapExisted, objects)
			}
		}
	}

	if err != nil {
		return WrapError(err)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {

		mapping := map[string]interface{}{
			"address_list":     object["AddressList"],
			"auto_add_tag_ecs": formatInt(object["AutoAddTagEcs"]),
			"description":      object["Description"],
			"group_name":       object["GroupName"],
			"group_type":       object["GroupType"],
			"id":               fmt.Sprint(object["GroupUuid"]),
			"group_uuid":       fmt.Sprint(object["GroupUuid"]),
			"tag_relation":     object["TagRelation"],
		}

		tags := make([]map[string]interface{}, 0)
		t, _ := jsonpath.Get("$.TagList", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				ecsTagItem := make(map[string]interface{})
				ecsTagItem["tag_value"] = t.(map[string]interface{})["TagValue"]
				ecsTagItem["tag_key"] = t.(map[string]interface{})["TagKey"]
				tags = append(tags, ecsTagItem)
			}
		}
		mapping["ecs_tags"] = tags
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
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

func describeAddressBook(response map[string]interface{}, conn *client.Client, request map[string]interface{}, idsMap map[string]bool, objects []map[string]interface{}) ([]map[string]interface{}, error) {
	action := "DescribeAddressBook"
	var err error
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
			return nil, WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_address_books", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Acls", response)
		if err != nil {
			return nil, WrapErrorf(err, FailedGetAttributeMsg, action, "$.Acls", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				gruopuuid := fmt.Sprint(item["GroupUuid"])
				if existed, ok := idsMap[gruopuuid]; ok {
					if !existed {
						delete(idsMap, gruopuuid)
						objects = append(objects, item)
					}
				}
			}
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	return objects, nil
}
