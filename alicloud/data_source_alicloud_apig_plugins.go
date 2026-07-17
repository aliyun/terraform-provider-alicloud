// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudApigPlugins() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudApigPluginRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_class_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugin_class_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_class_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_class_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_id": {
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
		},
	}
}

func dataSourceAliCloudApigPluginRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]*string
	// ListPlugins
	action := fmt.Sprintf("/v1/plugins")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("gateway_id"); ok {
		query["gatewayId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("plugin_class_id"); ok {
		query["pluginClassId"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("plugin_class_name"); ok {
		query["pluginClassName"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["pageSize"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["pageNumber"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.data.items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["pluginId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["pluginId"]

		mapping["plugin_id"] = objectRaw["pluginId"]

		gatewayInfoRawObj, _ := jsonpath.Get("$.gatewayInfo", objectRaw)
		gatewayInfoRaw := make(map[string]interface{})
		if gatewayInfoRawObj != nil {
			gatewayInfoRaw = gatewayInfoRawObj.(map[string]interface{})
		}
		mapping["gateway_name"] = gatewayInfoRaw["name"]
		mapping["gateway_id"] = gatewayInfoRaw["gatewayId"]

		pluginClassInfoRawObj, _ := jsonpath.Get("$.pluginClassInfo", objectRaw)
		pluginClassInfoRaw := make(map[string]interface{})
		if pluginClassInfoRawObj != nil {
			pluginClassInfoRaw = pluginClassInfoRawObj.(map[string]interface{})
		}
		mapping["plugin_class_name"] = pluginClassInfoRaw["name"]
		mapping["plugin_class_id"] = pluginClassInfoRaw["pluginClassId"]

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("plugins", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
