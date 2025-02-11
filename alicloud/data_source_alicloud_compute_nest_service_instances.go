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

func dataSourceAlicloudComputeNestServiceInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudComputeNestServiceInstancesRead,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Created", "Deploying", "DeployedFailed", "Deployed", "Upgrading", "Deleting", "Deleted", "DeletedFailed"}, false),
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 10,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"tags": tagsSchemaForceNew(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"service_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_instance_ops": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resources": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operated_service_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"deploy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"supplier_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"supplier_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"publish_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_infos": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"short_description": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"image": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"locale": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudComputeNestServiceInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListServiceInstances"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge

	if v, ok := d.GetOk("filter"); ok {
		filterMaps := make([]map[string]interface{}, 0)
		for _, filterList := range v.([]interface{}) {
			filterArg := filterList.(map[string]interface{})
			filterMap := map[string]interface{}{}

			if name, ok := filterArg["name"]; ok {
				filterMap["Name"] = name
			}

			if value, ok := filterArg["value"]; ok {
				filterMap["Value"] = value
			}

			filterMaps = append(filterMaps, filterMap)
		}

		request["Filter"] = filterMaps
	}

	if v, ok := d.GetOk("tags"); ok {
		request["Tag"] = tagsFromMap(v.(map[string]interface{}))
	}

	status, statusOk := d.GetOk("status")

	var objects []map[string]interface{}
	var serviceInstanceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		serviceInstanceNameRegex = r
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
			response, err = client.RpcPost("ComputeNest", "2021-06-01", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_compute_nest_service_instances", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.ServiceInstances", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ServiceInstances", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if serviceInstanceNameRegex != nil && !serviceInstanceNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}

			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ServiceInstanceId"])]; !ok {
					continue
				}
			}

			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
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
			"id":                           fmt.Sprint(object["ServiceInstanceId"]),
			"service_instance_id":          fmt.Sprint(object["ServiceInstanceId"]),
			"service_instance_name":        object["Name"],
			"parameters":                   object["Parameters"],
			"enable_instance_ops":          object["EnableInstanceOps"],
			"template_name":                object["TemplateName"],
			"operation_start_time":         object["OperationStartTime"],
			"operation_end_time":           object["OperationEndTime"],
			"resources":                    object["Resources"],
			"operated_service_instance_id": object["OperatedServiceInstanceId"],
			"source":                       object["Source"],
			"tags":                         tagsToMap(object["Tags"]),
			"status":                       object["Status"],
		}

		if v, ok := object["Service"]; ok {
			serviceMaps := make([]map[string]interface{}, 0)
			serviceArg := v.(map[string]interface{})
			serviceMap := map[string]interface{}{}

			if v, ok := serviceArg["ServiceId"]; ok {
				serviceMap["service_id"] = v
			}

			if v, ok := serviceArg["ServiceType"]; ok {
				serviceMap["service_type"] = v
			}

			if v, ok := serviceArg["DeployType"]; ok {
				serviceMap["deploy_type"] = v
			}

			if v, ok := serviceArg["SupplierName"]; ok {
				serviceMap["supplier_name"] = v
			}

			if v, ok := serviceArg["SupplierUrl"]; ok {
				serviceMap["supplier_url"] = v
			}

			if v, ok := serviceArg["PublishTime"]; ok {
				serviceMap["publish_time"] = v
			}

			if v, ok := serviceArg["Version"]; ok {
				serviceMap["version"] = v
			}

			if v, ok := serviceArg["VersionName"]; ok {
				serviceMap["version_name"] = v
			}

			if v, ok := serviceArg["ServiceInfos"]; ok {
				serviceInfosMaps := make([]map[string]interface{}, 0)
				for _, serviceInfosList := range v.([]interface{}) {
					serviceInfosArg := serviceInfosList.(map[string]interface{})
					serviceInfosMap := map[string]interface{}{}

					if v, ok := serviceInfosArg["Name"]; ok {
						serviceInfosMap["name"] = v
					}

					if v, ok := serviceInfosArg["ShortDescription"]; ok {
						serviceInfosMap["short_description"] = v
					}

					if v, ok := serviceInfosArg["Image"]; ok {
						serviceInfosMap["image"] = v
					}

					if v, ok := serviceInfosArg["Locale"]; ok {
						serviceInfosMap["locale"] = v
					}

					serviceInfosMaps = append(serviceInfosMaps, serviceInfosMap)
				}

				serviceMap["service_infos"] = serviceInfosMaps
			}

			if v, ok := serviceArg["Status"]; ok {
				serviceMap["status"] = v
			}

			serviceMaps = append(serviceMaps, serviceMap)
			mapping["service"] = serviceMaps
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

	if err := d.Set("service_instances", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
