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

func dataSourceAlicloudApiGatewayPlugins() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudApiGatewayPluginsRead,
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
			"plugin_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"plugin_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"backendSignature", "caching", "cors", "ipControl", "jwtAuth", "trafficControl"}, false),
			},
			"tags": tagsSchema(),
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  50,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plugins": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plugin_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudApiGatewayPluginsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribePlugins"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("plugin_name"); ok {
		request["PluginName"] = v
	}
	if v, ok := d.GetOk("plugin_type"); ok {
		request["PluginType"] = v
	}
	setPagingRequest(d, request, PageSizeLarge)
	var objects []map[string]interface{}
	var pluginNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		pluginNameRegex = r
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
	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
	}
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("CloudAPI", "2016-07-14", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_api_gateway_plugins", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Plugins.PluginAttribute", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Plugins.PluginAttribute", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if pluginNameRegex != nil && !pluginNameRegex.MatchString(fmt.Sprint(item["PluginName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PluginId"])]; !ok {
					continue
				}
			}
			if len(tagsMap) > 0 {
				if _, ok := item["Tags"]; !ok {
					continue
				}
				if len(item["Tags"].(map[string]interface{})["TagInfo"].([]interface{})) != len(tagsMap) {
					continue
				}
				match := true
				for _, tag := range item["Tags"].(map[string]interface{})["TagInfo"].([]interface{}) {
					if v, ok := tagsMap[tag.(map[string]interface{})["Key"].(string)]; !ok || v.(string) != tag.(map[string]interface{})["Value"].(string) {
						match = false
						break
					}
				}
				if !match {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"create_time":   object["CreatedTime"],
			"description":   object["Description"],
			"modified_time": object["ModifiedTime"],
			"plugin_data":   object["PluginData"],
			"id":            fmt.Sprint(object["PluginId"]),
			"plugin_id":     fmt.Sprint(object["PluginId"]),
			"plugin_name":   object["PluginName"],
			"plugin_type":   object["PluginType"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.TagInfo", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["Key"].(string)
				value := t.(map[string]interface{})["Value"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["PluginName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
