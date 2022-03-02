package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcsImageComponents() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsImageComponentsRead,
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
			"image_component_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"owner": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALIYUN", "SELF"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_component_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_component_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_type": {
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

func dataSourceAlicloudEcsImageComponentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeImageComponents"
	request := make(map[string]interface{})
	if m, ok := d.GetOk("ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("ImageComponentId.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("image_component_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("owner"); ok {
		request["Owner"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var imageComponentNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		imageComponentNameRegex = r
	}

	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_image_components", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ImageComponent.ImageComponentSet", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ImageComponent.ImageComponentSet", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if imageComponentNameRegex != nil && !imageComponentNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			objects = append(objects, item)
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
			"component_type":       object["ComponentType"],
			"content":              object["Content"],
			"create_time":          object["CreationTime"],
			"description":          object["Description"],
			"id":                   fmt.Sprint(object["ImageComponentId"]),
			"image_component_id":   fmt.Sprint(object["ImageComponentId"]),
			"image_component_name": object["Name"],
			"owner":                object["Owner"],
			"resource_group_id":    object["ResourceGroupId"],
			"system_type":          object["SystemType"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
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

	if err := d.Set("components", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
