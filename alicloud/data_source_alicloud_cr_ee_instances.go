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

func dataSourceAliCloudCrInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCrInstanceRead,
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
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_endpoints": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domains": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"domain": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"endpoint_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_issue": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"renew_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"renewal_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAliCloudCrInstanceRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "ListInstance"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["InstanceName"] = d.Get("instance_name")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNo"] = 1
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RpcPost("cr", "2018-12-01", action, query, request, true)

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

		resp, _ := jsonpath.Get("$.Instances[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["InstanceName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["InstanceId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		request["PageNo"] = request["PageNo"].(int) + 1
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["InstanceId"]

		mapping["create_time"] = objectRaw["CreateTime"]
		mapping["instance_issue"] = objectRaw["InstanceIssue"]
		mapping["instance_name"] = objectRaw["InstanceName"]
		mapping["modified_time"] = formatInt(objectRaw["ModifiedTime"])
		mapping["region_id"] = objectRaw["RegionId"]
		mapping["resource_group_id"] = objectRaw["ResourceGroupId"]
		mapping["instance_id"] = objectRaw["InstanceId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["InstanceName"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["InstanceId"])
		mapping, err = dataSourceAliCloudCrInstanceReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["InstanceName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudCrInstanceReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	crServiceV2 := CrServiceV2{client}
	getResp, err := crServiceV2.DescribeCrInstance(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["create_time"] = objectRaw["CreateTime"]
	mapping["instance_issue"] = objectRaw["InstanceIssue"]
	mapping["instance_name"] = objectRaw["InstanceName"]
	mapping["modified_time"] = objectRaw["ModifiedTime"]
	mapping["resource_group_id"] = objectRaw["ResourceGroupId"]
	mapping["status"] = objectRaw["InstanceStatus"]
	mapping["instance_id"] = objectRaw["InstanceId"]

	objectRaw = getResp

	mapping["create_time"] = objectRaw["CreateTime"]
	mapping["end_time"] = objectRaw["EndTime"]
	mapping["payment_type"] = objectRaw["SubscriptionType"]
	mapping["region_id"] = objectRaw["Region"]
	mapping["renew_period"] = objectRaw["RenewalDuration"]
	mapping["renewal_status"] = objectRaw["RenewStatus"]
	mapping["instance_id"] = objectRaw["InstanceID"]

	objectRaw = getResp

	endpointsRaw := objectRaw["Endpoints"]
	instanceEndpointsMaps := make([]map[string]interface{}, 0)
	if endpointsRaw != nil {
		for _, endpointsChildRaw := range convertToInterfaceArray(endpointsRaw) {
			instanceEndpointsMap := make(map[string]interface{})
			endpointsChildRaw := endpointsChildRaw.(map[string]interface{})
			instanceEndpointsMap["enable"] = endpointsChildRaw["Enable"]
			instanceEndpointsMap["endpoint_type"] = endpointsChildRaw["EndpointType"]

			domainsRaw := endpointsChildRaw["Domains"]
			domainsMaps := make([]map[string]interface{}, 0)
			if domainsRaw != nil {
				for _, domainsChildRaw := range convertToInterfaceArray(domainsRaw) {
					domainsMap := make(map[string]interface{})
					domainsChildRaw := domainsChildRaw.(map[string]interface{})
					domainsMap["domain"] = domainsChildRaw["Domain"]
					domainsMap["type"] = domainsChildRaw["Type"]

					domainsMaps = append(domainsMaps, domainsMap)
				}
			}
			instanceEndpointsMap["domains"] = domainsMaps
			instanceEndpointsMaps = append(instanceEndpointsMaps, instanceEndpointsMap)
		}
	}
	mapping["instance_endpoints"] = instanceEndpointsMaps

	return mapping, nil
}
