// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudEnsSecurityGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEnsSecurityGroupRead,
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
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
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
						"instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"source_port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"source_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"creation_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dest_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6_source_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ipv6_dest_cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"direction": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_name": {
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

func dataSourceAliCloudEnsSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
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
	action := "DescribeSecurityGroups"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})

	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		request["SecurityGroupId"] = v
	}
	if v, ok := d.GetOk("security_group_name"); ok {
		request["SecurityGroupName"] = v
	}
	request["MaxResults"] = PageSizeLarge
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
			response, err = client.RpcPost("Ens", "2017-11-10", action, query, request, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.SecurityGroups.SecurityGroup[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["SecurityGroupName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SecurityGroupId"])]; !ok {
					continue
				}
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
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["SecurityGroupId"]

		mapping["description"] = objectRaw["Description"]
		mapping["security_group_name"] = objectRaw["SecurityGroupName"]
		mapping["create_time"] = objectRaw["CreationTime"]
		mapping["instance_count"] = objectRaw["InstanceCount"]
		mapping["security_group_id"] = objectRaw["SecurityGroupId"]

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["SecurityGroupName"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["SecurityGroupId"])
		mapping, err = dataSourceAliCloudEnsSecurityGroupReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["SecurityGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudEnsSecurityGroupReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	ensServiceV2 := EnsServiceV2{client}
	getResp, err := ensServiceV2.DescribeEnsSecurityGroup(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["description"] = objectRaw["Description"]
	mapping["security_group_name"] = objectRaw["SecurityGroupName"]
	mapping["security_group_id"] = objectRaw["SecurityGroupId"]

	permissionRaw, _ := jsonpath.Get("$.Permissions.Permission", objectRaw)
	permissionsMaps := make([]map[string]interface{}, 0)
	if permissionRaw != nil {
		for _, permissionChildRaw := range convertToInterfaceArray(permissionRaw) {
			permissionsMap := make(map[string]interface{})
			permissionChildRaw := permissionChildRaw.(map[string]interface{})
			permissionsMap["creation_time"] = permissionChildRaw["CreationTime"]
			permissionsMap["description"] = permissionChildRaw["Description"]
			permissionsMap["dest_cidr_ip"] = permissionChildRaw["DestCidrIp"]
			permissionsMap["direction"] = permissionChildRaw["Direction"]
			permissionsMap["ip_protocol"] = permissionChildRaw["IpProtocol"]
			permissionsMap["ipv6_dest_cidr_ip"] = permissionChildRaw["Ipv6DestCidrIp"]
			permissionsMap["ipv6_source_cidr_ip"] = permissionChildRaw["Ipv6SourceCidrIp"]
			permissionsMap["policy"] = permissionChildRaw["Policy"]
			permissionsMap["port_range"] = permissionChildRaw["PortRange"]
			permissionsMap["priority"] = permissionChildRaw["Priority"]
			permissionsMap["source_cidr_ip"] = permissionChildRaw["SourceCidrIp"]
			permissionsMap["source_port_range"] = permissionChildRaw["SourcePortRange"]

			permissionsMaps = append(permissionsMaps, permissionsMap)
		}
	}
	mapping["permissions"] = permissionsMaps

	return mapping, nil
}
