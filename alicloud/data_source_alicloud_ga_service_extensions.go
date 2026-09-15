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

func dataSourceAlicloudGaServiceExtensions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaServiceExtensionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"service_extension_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
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
			"service_extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_extension_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"components": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_component_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"config": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fail_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"component_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accelerator_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"associate_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"create_time": {
							Type:     schema.TypeString,
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

func dataSourceAlicloudGaServiceExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListServiceExtensions"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("service_extension_name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	var err error
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = client.RpcPost("Ga", "2019-11-20", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_service_extensions", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ServiceExtensions", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ServiceExtensions", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ServiceExtensionId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != fmt.Sprint(item["State"]) {
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
	gaService := GaService{client}
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                   fmt.Sprint(object["ServiceExtensionId"]),
			"service_extension_id": fmt.Sprint(object["ServiceExtensionId"]),
			"name":                 object["Name"],
			"description":          object["Description"],
			"type":                 object["Type"],
			"state":                object["State"],
			"resource_group_id":    object["ResourceGroupId"],
			"create_time":          object["CreateTime"],
			"update_time":          object["UpdateTime"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["ServiceExtensionId"])
		getResp, err := gaService.DescribeGaServiceExtension(id)
		if err != nil {
			return WrapError(err)
		}
		if v, ok := getResp["Components"].([]interface{}); ok {
			components := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				temp := map[string]interface{}{
					"service_component_id": item["ServiceComponentId"],
					"priority":             item["Priority"],
					"timeout":              item["Timeout"],
					"config":               item["Config"],
					"fail_policy":          item["FailPolicy"],
					"component_name":       item["ComponentName"],
				}
				components = append(components, temp)
			}
			mapping["components"] = components
		}
		if v, ok := getResp["Resources"].([]interface{}); ok {
			resources := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				temp := map[string]interface{}{
					"accelerator_id": item["AcceleratorId"],
					"resource_type":  item["ResourceType"],
					"resource_id":    item["ResourceId"],
					"associate_id":   item["AssociateId"],
				}
				resources = append(resources, temp)
			}
			mapping["resources"] = resources
		}
		listTagResourcesObject, err := gaService.ListTagResources(id, "serviceextension")
		if err != nil {
			return WrapError(err)
		}
		mapping["tags"] = tagsToMap(listTagResourcesObject)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("service_extensions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
