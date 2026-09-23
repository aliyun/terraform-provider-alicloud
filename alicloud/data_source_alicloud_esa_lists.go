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

func dataSourceAliCloudEsaLists() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEsaListRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"query_args": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id_like": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name_like": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description_like": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name_item_like": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"item_like": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"order_by": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"desc": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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
			"lists": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"list_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEsaListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListLists"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("query_args"); ok {
		queryArgsMap := map[string]interface{}{}
		for _, queryArgsList := range v.([]interface{}) {
			queryArgsArg := queryArgsList.(map[string]interface{})

			if idLike, ok := queryArgsArg["id_like"]; ok {
				queryArgsMap["IdLike"] = idLike
			}

			if nameLike, ok := queryArgsArg["name_like"]; ok {
				queryArgsMap["NameLike"] = nameLike
			}

			if descriptionLike, ok := queryArgsArg["description_like"]; ok {
				queryArgsMap["DescriptionLike"] = descriptionLike
			}

			if nameItemLike, ok := queryArgsArg["name_item_like"]; ok {
				queryArgsMap["NameItemLike"] = nameItemLike
			}

			if itemLike, ok := queryArgsArg["item_like"]; ok {
				queryArgsMap["ItemLike"] = itemLike
			}

			if kind, ok := queryArgsArg["kind"]; ok {
				queryArgsMap["Kind"] = kind
			}

			if orderBy, ok := queryArgsArg["order_by"]; ok {
				queryArgsMap["OrderBy"] = orderBy
			}

			if desc, ok := d.GetOkExists("query_args.0.desc"); ok {
				queryArgsMap["Desc"] = desc
			}
		}

		request["QueryArgs"] = convertObjectToJsonString(queryArgsMap)
	}

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

	var listNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}

		listNameRegex = r
	}

	var response map[string]interface{}
	var err error

	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("ESA", "2024-09-10", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_esa_lists", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Lists", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Lists", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
					continue
				}
			}

			if listNameRegex != nil {
				if !listNameRegex.MatchString(fmt.Sprint(item["Name"])) {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":          fmt.Sprint(object["Id"]),
			"list_id":     fmt.Sprint(object["Id"]),
			"name":        object["Name"],
			"kind":        object["Kind"],
			"description": object["Description"],
			"length":      object["Length"],
			"update_time": object["UpdateTime"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
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
