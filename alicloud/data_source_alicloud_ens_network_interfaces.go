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

func dataSourceAliCloudEnsNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsNetworkInterfaceRead,
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
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_interface_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interfaces": {
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
						"ens_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_ip_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vmnc_learn": {
							Type:     schema.TypeBool,
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

func dataSourceAliCloudEnsNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DescribeNetworkInterfaces"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("network_interface_id"); ok {
		request["NetworkInterfaceId"] = v
	}
	if v, ok := d.GetOk("ens_region_id"); ok {
		request["EnsRegionId"] = v
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("network_id"); ok {
		request["NetworkId"] = v
	}
	if v, ok := d.GetOk("network_interface_id"); ok {
		request["NetworkInterfaceId"] = v
	}
	if v, ok := d.GetOk("network_interface_name"); ok {
		request["NetworkInterfaceName"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
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

		resp, _ := jsonpath.Get("$.NetworkInterfaceSets.NetworkInterfaceSet[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["NetworkInterfaceName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["NetworkInterfaceId"])]; !ok {
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

		mapping["id"] = objectRaw["NetworkInterfaceId"]

		mapping["create_time"] = objectRaw["CreationTime"]
		mapping["description"] = objectRaw["Description"]
		mapping["ens_region_id"] = objectRaw["EnsRegionId"]
		mapping["instance_id"] = objectRaw["InstanceId"]
		mapping["mac_address"] = objectRaw["MacAddress"]
		mapping["network_id"] = objectRaw["NetworkId"]
		mapping["network_interface_name"] = objectRaw["NetworkInterfaceName"]
		mapping["primary_ip"] = objectRaw["PrimaryIp"]
		mapping["primary_ip_type"] = objectRaw["PrimaryIpType"]
		mapping["status"] = objectRaw["Status"]
		mapping["vswitch_id"] = objectRaw["VSwitchId"]
		mapping["vmnc_learn"] = objectRaw["VmncLearn"]
		mapping["network_interface_id"] = objectRaw["NetworkInterfaceId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["NetworkInterfaceName"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["NetworkInterfaceId"])
		mapping, err = dataSourceAliCloudEnsNetworkInterfaceReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["NetworkInterfaceName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("interfaces", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudEnsNetworkInterfaceReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	ensServiceV2 := EnsServiceV2{client}
	getResp, err := ensServiceV2.DescribeEnsNetworkInterface(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["create_time"] = objectRaw["CreationTime"]
	mapping["description"] = objectRaw["Description"]
	mapping["ens_region_id"] = objectRaw["EnsRegionId"]
	mapping["instance_id"] = objectRaw["InstanceId"]
	mapping["mac_address"] = objectRaw["MacAddress"]
	mapping["network_id"] = objectRaw["NetworkId"]
	mapping["network_interface_name"] = objectRaw["NetworkInterfaceName"]
	mapping["primary_ip"] = objectRaw["PrimaryIp"]
	mapping["primary_ip_type"] = objectRaw["PrimaryIpType"]
	mapping["status"] = objectRaw["Status"]
	mapping["vswitch_id"] = objectRaw["VSwitchId"]
	mapping["vmnc_learn"] = objectRaw["VmncLearn"]
	mapping["network_interface_id"] = objectRaw["NetworkInterfaceId"]

	securityGroupRaw, _ := jsonpath.Get("$.SecurityGroupIds.SecurityGroup", objectRaw)
	mapping["security_group_ids"] = securityGroupRaw

	return mapping, nil
}
