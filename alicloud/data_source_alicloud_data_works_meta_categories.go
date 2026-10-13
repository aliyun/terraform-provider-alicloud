package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudDataWorksMetaCategories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDataWorksMetaCategoriesRead,
		Schema: map[string]*schema.Schema{
			"parent_category_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"categories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"parent_category_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDataWorksMetaCategoriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "GetMetaCategory"
	request := make(map[string]interface{})
	parentCategoryId := d.Get("parent_category_id").(int)
	request["ParentCategoryId"] = parentCategoryId
	request["PageNum"] = 1
	request["PageSize"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var objects []map[string]interface{}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("dataworks-public", "2020-05-18", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_data_works_meta_categories", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Data.DataEntityList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.DataEntityList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			// GetMetaCategory returns CategoryId as json.Number (RpcPost UseNumber
			// decoding); a .(float64) assertion panics. toInt tolerates json.Number.
			categoryId, err := toInt(item["CategoryId"])
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.DataEntityList[].CategoryId", item)
			}
			id := fmt.Sprintf("%d:%d", categoryId, parentCategoryId)
			if len(idsMap) > 0 {
				// Match against the full composite id ("categoryId:parentCategoryId"),
				// the same format the resource uses for SetId and that callers pass
				// in the ids list. The previous code matched on the bare category id
				// ("categoryId"), which never matched ids like "56855:56854" and
				// caused categories.# to be 0.
				if _, ok := idsMap[id]; !ok {
					continue
				}
			}
			createTime, err := toInt(item["CreateTime"])
			if err != nil {
				return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.DataEntityList[].CreateTime", item)
			}
			objects = append(objects, map[string]interface{}{
				"id":                 id,
				"category_id":        categoryId,
				"parent_category_id": parentCategoryId,
				"name":               item["Name"],
				"comment":            item["Comment"],
				"create_time":        createTime,
			})
		}
		// TotalCount is a json.Number from RpcPost; the previous float64 assertion
		// silently failed and left totalCount at 0, which broke pagination early.
		totalCount := 0
		if v, err := jsonpath.Get("$.Data.TotalCount", response); err == nil {
			if tc, err := toInt(v); err == nil {
				totalCount = tc
			}
		}
		if len(objects) >= totalCount || len(result) < PageSizeLarge {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		ids = append(ids, object["id"].(string))
		s = append(s, object)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("categories", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
