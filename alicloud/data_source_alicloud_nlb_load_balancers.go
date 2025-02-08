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

func dataSourceAlicloudNlbLoadBalancers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNlbLoadBalancersRead,
		Schema: map[string]*schema.Schema{
			"address_ip_version": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Ipv4", "DualStack"}, false),
			},
			"address_type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Internet", "Intranet"}, false),
			},
			"dns_name": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"ipv6_address_type": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Internet", "Intranet"}, false),
			},
			"load_balancer_business_status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Abnormal", "Normal"}, false),
			},
			"resource_group_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"zone_id": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeString,
			},
			"vpc_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"load_balancer_names": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"Inactive", "Active", "Provisioning", "Configuring", "Deleting", "Deleted"}, false),
			},
			"tags": tagsSchema(),
			"ids": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Optional:     true,
				ForceNew:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_file": {
				Optional: true,
				Type:     schema.TypeString,
			},
			"balancers": {
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"address_ip_version": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"address_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"bandwidth_package_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"create_time": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"cross_zone_enabled": {
							Computed: true,
							Type:     schema.TypeBool,
						},
						"dns_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"ipv6_address_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"load_balancer_business_status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"load_balancer_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"load_balancer_name": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"load_balancer_type": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"resource_group_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"security_group_ids": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"tags": tagsSchema(),
						"vpc_id": {
							Computed: true,
							Type:     schema.TypeString,
						},
						"zone_mappings": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allocation_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"eni_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"ipv6_address": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"private_ipv4_address": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"public_ipv4_address": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"vswitch_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"zone_id": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
						"operation_locks": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lock_type": {
										Computed: true,
										Type:     schema.TypeString,
									},
									"lock_reason": {
										Computed: true,
										Type:     schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNlbLoadBalancersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIpVersion"] = v
	}
	if v, ok := d.GetOk("address_type"); ok {
		request["AddressType"] = v
	}
	if v, ok := d.GetOk("dns_name"); ok {
		request["DNSName"] = v
	}
	if v, ok := d.GetOk("vpc_ids"); ok {
		request["VpcIds"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("load_balancer_names"); ok {
		request["LoadBalancerNames"] = v
	}
	if v, ok := d.GetOk("ipv6_address_type"); ok {
		request["Ipv6AddressType"] = v
	}
	if v, ok := d.GetOk("load_balancer_business_status"); ok {
		request["LoadBalancerBusinessStatus"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["LoadBalancerStatus"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		request["Tag"] = tagsFromMap(v.(map[string]interface{}))
	}
	request["MaxResults"] = PageSizeLarge

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var loadBalancerNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		loadBalancerNameRegex = r
	}

	var err error
	var objects []interface{}
	var response map[string]interface{}

	for {
		action := "ListLoadBalancers"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err := client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			response = resp
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nlb_load_balancers", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.LoadBalancers", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.LoadBalancers", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["LoadBalancerId"])]; !ok {
					continue
				}
			}

			if loadBalancerNameRegex != nil && !loadBalancerNameRegex.MatchString(fmt.Sprint(item["LoadBalancerName"])) {
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
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                            fmt.Sprint(object["LoadBalancerId"]),
			"address_ip_version":            object["AddressIpVersion"],
			"address_type":                  object["AddressType"],
			"bandwidth_package_id":          object["BandwidthPackageId"],
			"create_time":                   object["CreateTime"],
			"cross_zone_enabled":            object["CrossZoneEnabled"],
			"dns_name":                      object["DNSName"],
			"ipv6_address_type":             object["Ipv6AddressType"],
			"load_balancer_business_status": object["LoadBalancerBusinessStatus"],
			"load_balancer_id":              object["LoadBalancerId"],
			"load_balancer_name":            object["LoadBalancerName"],
			"load_balancer_type":            object["LoadBalancerType"],
			"resource_group_id":             object["ResourceGroupId"],
			"status":                        object["LoadBalancerStatus"],
			"vpc_id":                        object["VpcId"],
		}
		if v, ok := object["SecurityGroupIds"]; ok {
			mapping["security_group_ids"] = v.([]interface{})
		}

		tagsMap := make(map[string]interface{})
		tagsRaw, _ := jsonpath.Get("$.Tags", object)
		if tagsRaw != nil {
			for _, value0 := range tagsRaw.([]interface{}) {
				tags := value0.(map[string]interface{})
				key := tags["Key"].(string)
				value := tags["Value"]
				if !ignoredTags(key, value) {
					tagsMap[key] = value
				}
			}
		}
		if len(tagsMap) > 0 {
			mapping["tags"] = tagsMap
		}
		zoneMappingsMaps := make([]map[string]interface{}, 0)
		if zoneMappingsRaw, ok := object["ZoneMappings"]; ok {
			for _, value0 := range zoneMappingsRaw.([]interface{}) {
				zoneMappings := value0.(map[string]interface{})
				zoneMappingsMap := make(map[string]interface{})
				zoneMappingsMap["vswitch_id"] = zoneMappings["VSwitchId"]
				zoneMappingsMap["zone_id"] = zoneMappings["ZoneId"]
				if v, ok := zoneMappings["LoadBalancerAddresses"]; ok && len(v.([]interface{})) > 0 {
					LoadBalancerAddressesMap := v.([]interface{})[0].(map[string]interface{})
					zoneMappingsMap["allocation_id"] = LoadBalancerAddressesMap["AllocationId"]
					zoneMappingsMap["eni_id"] = LoadBalancerAddressesMap["EniId"]
					zoneMappingsMap["ipv6_address"] = LoadBalancerAddressesMap["Ipv6Address"]
					zoneMappingsMap["private_ipv4_address"] = LoadBalancerAddressesMap["PrivateIPv4Address"]
					zoneMappingsMap["public_ipv4_address"] = LoadBalancerAddressesMap["PublicIPv4Address"]
				}
				zoneMappingsMaps = append(zoneMappingsMaps, zoneMappingsMap)
			}
		}
		mapping["zone_mappings"] = zoneMappingsMaps

		operationLocksMaps := make([]map[string]interface{}, 0)
		if operationLocksRaw, ok := object["OperationLocks"]; ok {
			for _, value0 := range operationLocksRaw.([]interface{}) {
				operationLocks := value0.(map[string]interface{})
				operationLocksMap := make(map[string]interface{})
				operationLocksMap["lock_type"] = operationLocks["LockType"]
				operationLocksMap["lock_reason"] = operationLocks["LockReason"]
				operationLocksMaps = append(operationLocksMaps, operationLocksMap)
			}
		}
		mapping["operation_locks"] = operationLocksMaps

		ids = append(ids, fmt.Sprint(object["LoadBalancerId"]))
		names = append(names, object["LoadBalancerName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("balancers", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
