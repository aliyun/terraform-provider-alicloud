// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudEnsGatewayQos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsGatewayQosRead,
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
			"ens_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_qos_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_qos_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_qos_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"qos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth_in": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ens_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_qos_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_qos_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_qos_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceAliCloudEnsGatewayQosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

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

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeEnsGatewayQoses"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("gateway_qos_id"); ok {
		request["GatewayQosId"] = v
	}
	if v, ok := d.GetOk("ens_region_id"); ok {
		request["EnsRegionId"] = v
	}
	if v, ok := d.GetOk("gateway_qos_id"); ok {
		request["GatewayQosId"] = v
	}
	if v, ok := d.GetOk("gateway_qos_name"); ok {
		request["GatewayQosName"] = v
	}
	request["GatewayQosType"] = d.Get("gateway_qos_type")
	if v, ok := d.GetOk("instances"); ok && v.(*schema.Set).Len() > 0 {
		for _, dataLoop := range v.(*schema.Set).List() {
			if dataLoopTmp, ok := dataLoop.(map[string]interface{}); ok {
				if instanceId, ok := dataLoopTmp["instance_id"].(string); ok && instanceId != "" {
					request["InstanceId"] = instanceId
					break
				}
			}
		}
	}
	request["NetworkId"] = d.Get("network_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.GatewayQoses.GatewayQos[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["GatewayQosName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["GatewayQosId"])]; !ok {
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
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["GatewayQosId"]

		mapping["bandwidth_in"] = objectRaw["BandwidthIn"]
		mapping["bandwidth_out"] = objectRaw["BandwidthOut"]
		mapping["creation_time"] = objectRaw["CreationTime"]
		mapping["ens_region_id"] = objectRaw["EnsRegionId"]
		mapping["gateway_qos_name"] = objectRaw["GatewayQosName"]
		mapping["gateway_qos_type"] = objectRaw["GatewayQosType"]
		mapping["network_id"] = objectRaw["NetworkId"]
		mapping["status"] = objectRaw["Status"]
		mapping["gateway_qos_id"] = objectRaw["GatewayQosId"]

		instanceRaw, _ := jsonpath.Get("$.Instances.Instance", objectRaw)
		instancesMaps := make([]map[string]interface{}, 0)
		if instanceRaw != nil {
			for _, instanceChildRaw := range convertToInterfaceArray(instanceRaw) {
				instancesMap := make(map[string]interface{})
				instanceChildRaw := instanceChildRaw.(map[string]interface{})
				instancesMap["instance_id"] = instanceChildRaw["InstanceId"]
				instancesMap["instance_type"] = instanceChildRaw["InstanceType"]
				instancesMap["status"] = instanceChildRaw["Status"]

				instancesMaps = append(instancesMaps, instancesMap)
			}
		}
		mapping["instances"] = instancesMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["GatewayQosName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("qos", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
