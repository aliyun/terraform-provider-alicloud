// Package apig. This file is generated automatically. Please do not modify it manually, thank you!
package apig

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServiceRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"create_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"express_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"unhealthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"http_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"healthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"expected_statuses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"health_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"healthy_panic_threshold": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"outlier_detection_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"failure_percentage_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"base_ejection_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"failure_percentage_minimum_hosts": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"outlier_endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qualifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_detail_error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_detail_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"unhealthy_endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"update_timestamp": {
							Type:     schema.TypeInt,
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

func dataSourceServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	// ListServices
	action := fmt.Sprintf("/v1/services")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)

	if v, ok := d.GetOk("gateway_id"); ok {
		query["gatewayId"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		query["resourceGroupId"] = helper.StringPointer(v.(string))
	}

	if v, ok := d.GetOk("source_type"); ok {
		query["sourceType"] = helper.StringPointer(v.(string))
	}

	query["pageSize"] = helper.StringPointer(strconv.Itoa(helper.PageSizeLarge))
	query["pageNumber"] = helper.StringPointer("1")
	for {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RoaGet("APIG", "2024-03-27", action, query, nil, nil)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.data.items[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["serviceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < helper.PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["pageNumber"])
		query["pageNumber"] = helper.StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["serviceId"]

		mapping["create_timestamp"] = objectRaw["createTimestamp"]
		mapping["express_type"] = objectRaw["expressType"]
		mapping["gateway_id"] = objectRaw["gatewayId"]
		mapping["health_status"] = objectRaw["healthStatus"]
		mapping["healthy_panic_threshold"] = objectRaw["healthyPanicThreshold"]
		mapping["namespace"] = objectRaw["namespace"]
		mapping["qualifier"] = objectRaw["qualifier"]
		mapping["resource_group_id"] = objectRaw["resourceGroupId"]
		mapping["runtime_detail_error_code"] = objectRaw["runtimeDetailErrorCode"]
		mapping["runtime_detail_status"] = objectRaw["runtimeDetailStatus"]
		mapping["service_name"] = objectRaw["name"]
		mapping["source_type"] = objectRaw["sourceType"]
		mapping["update_timestamp"] = objectRaw["updateTimestamp"]
		mapping["protocol"] = objectRaw["protocol"]
		mapping["service_id"] = objectRaw["serviceId"]

		addressesRaw := make([]interface{}, 0)
		if objectRaw["addresses"] != nil {
			addressesRaw = convert.ConvertToInterfaceArray(objectRaw["addresses"])
		}

		mapping["addresses"] = addressesRaw
		dnsServersRaw := make([]interface{}, 0)
		if objectRaw["dnsServers"] != nil {
			dnsServersRaw = convert.ConvertToInterfaceArray(objectRaw["dnsServers"])
		}

		mapping["dns_servers"] = dnsServersRaw
		healthCheckConfigMaps := make([]map[string]interface{}, 0)
		healthCheckConfigMap := make(map[string]interface{})
		healthCheckRaw := make(map[string]interface{})
		if objectRaw["healthCheck"] != nil {
			healthCheckRaw = objectRaw["healthCheck"].(map[string]interface{})
		}
		if len(healthCheckRaw) > 0 {
			healthCheckConfigMap["enable"] = healthCheckRaw["enable"]
			healthCheckConfigMap["healthy_threshold"] = healthCheckRaw["healthyThreshold"]
			healthCheckConfigMap["http_host"] = healthCheckRaw["httpHost"]
			healthCheckConfigMap["http_path"] = healthCheckRaw["httpPath"]
			healthCheckConfigMap["interval"] = healthCheckRaw["interval"]
			healthCheckConfigMap["protocol"] = healthCheckRaw["protocol"]
			healthCheckConfigMap["timeout"] = healthCheckRaw["timeout"]
			healthCheckConfigMap["unhealthy_threshold"] = healthCheckRaw["unhealthyThreshold"]

			healthCheckConfigMaps = append(healthCheckConfigMaps, healthCheckConfigMap)
		}
		mapping["health_check_config"] = healthCheckConfigMaps
		outlierDetectionConfigMaps := make([]map[string]interface{}, 0)
		outlierDetectionConfigMap := make(map[string]interface{})
		outlierDetectionRaw := make(map[string]interface{})
		if objectRaw["outlierDetection"] != nil {
			outlierDetectionRaw = objectRaw["outlierDetection"].(map[string]interface{})
		}
		if len(outlierDetectionRaw) > 0 {
			outlierDetectionConfigMap["base_ejection_time"] = outlierDetectionRaw["baseEjectionTime"]
			outlierDetectionConfigMap["enable"] = outlierDetectionRaw["enable"]
			outlierDetectionConfigMap["failure_percentage_minimum_hosts"] = outlierDetectionRaw["failurePercentageMinimumHosts"]
			outlierDetectionConfigMap["failure_percentage_threshold"] = outlierDetectionRaw["failurePercentageThreshold"]
			outlierDetectionConfigMap["interval"] = outlierDetectionRaw["interval"]

			outlierDetectionConfigMaps = append(outlierDetectionConfigMaps, outlierDetectionConfigMap)
		}
		mapping["outlier_detection_config"] = outlierDetectionConfigMaps
		outlierEndpointsRaw := make([]interface{}, 0)
		if objectRaw["outlierEndpoints"] != nil {
			outlierEndpointsRaw = convert.ConvertToInterfaceArray(objectRaw["outlierEndpoints"])
		}

		mapping["outlier_endpoints"] = outlierEndpointsRaw
		unhealthyEndpointsRaw := make([]interface{}, 0)
		if objectRaw["unhealthyEndpoints"] != nil {
			unhealthyEndpointsRaw = convert.ConvertToInterfaceArray(objectRaw["unhealthyEndpoints"])
		}

		mapping["unhealthy_endpoints"] = unhealthyEndpointsRaw
		portsRaw := objectRaw["ports"]
		portsMaps := make([]map[string]interface{}, 0)
		if portsRaw != nil {
			for _, portsChildRaw := range convert.ConvertToInterfaceArray(portsRaw) {
				portsMap := make(map[string]interface{})
				portsChildRaw := portsChildRaw.(map[string]interface{})
				portsMap["name"] = portsChildRaw["name"]
				portsMap["port"] = portsChildRaw["port"]
				portsMap["protocol"] = portsChildRaw["protocol"]

				portsMaps = append(portsMaps, portsMap)
			}
		}
		mapping["ports"] = portsMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("services", s); err != nil {
		return errs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}
	return nil
}
