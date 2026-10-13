package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudCmsTransformers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsTransformersRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"transformer_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transformer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"transformers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transformer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transformer_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quit_after_match": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sort_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
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
		},
	}
}

func dataSourceAliCloudCmsTransformersRead(d *schema.ResourceData, meta interface{}) error {
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

	action := "/transformers"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	query["regionId"] = StringPointer(client.RegionId)
	query["workspace"] = StringPointer(d.Get("workspace").(string))
	if v, ok := d.GetOk("transformer_name"); ok {
		query["transformerName"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("transformer_id"); ok {
		query["transformerId"] = StringPointer(v.(string))
	}
	if v, ok := d.GetOk("enable"); ok {
		query["enable"] = StringPointer(strconv.FormatBool(v.(bool)))
	}

	query["MaxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("Cms", "2024-03-30", action, query, nil, nil)
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

		resp, _ := jsonpath.Get("$.transformerList[*]", response)
		result, _ := resp.([]interface{})
		for _, v := range result {
			item, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["transformerId"], ":", item["workspace"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			query["NextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}
		mapping["id"] = fmt.Sprint(objectRaw["transformerId"], ":", objectRaw["workspace"])
		mapping["transformer_id"] = objectRaw["transformerId"]
		mapping["transformer_name"] = objectRaw["transformerName"]
		mapping["workspace"] = objectRaw["workspace"]
		mapping["description"] = objectRaw["description"]
		mapping["quit_after_match"] = objectRaw["quitAfterMatch"]
		mapping["sort_id"] = objectRaw["sortId"]
		mapping["enable"] = objectRaw["enable"]
		mapping["create_time"] = objectRaw["createTime"]
		mapping["update_time"] = objectRaw["updateTime"]
		mapping["user_id"] = objectRaw["userId"]
		mapping["region_id"] = objectRaw["regionId"]

		ids = append(ids, mapping["id"].(string))
		s = append(s, mapping)
	}

	d.SetId(fmt.Sprintf("%v", request))
	d.Set("transformers", s)
	d.Set("ids", ids)

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			log.Printf("[WARN] writeToFile for alicloud_cms_transformers datasource failed: %s", err)
		}
	}

	return nil
}
