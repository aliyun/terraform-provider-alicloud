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

func dataSourceAlicloudNlbServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNlbServerGroupsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
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
			"server_group_names": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"server_group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Instance", "Ip"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Configuring", "Creating"}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_ip_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_drain": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"connection_drain_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"health_check_connect_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_connect_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"health_check_http_code": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"health_check_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"healthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"unhealthy_threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"http_check_method": {
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
						"related_load_balancer_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"server_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"preserve_client_ip_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudNlbServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListServerGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("server_group_names"); ok {
		request["ServerGroupNames"] = v
	}
	if v, ok := d.GetOk("server_group_type"); ok {
		request["ServerGroupType"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		request["Tag"] = tagsFromMap(v.(map[string]interface{}))
	}
	var objects []map[string]interface{}
	var serverGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		serverGroupNameRegex = r
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
			response, err = client.RpcPost("Nlb", "2022-04-30", action, nil, request, true)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nlb_server_groups", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.ServerGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ServerGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if serverGroupNameRegex != nil && !serverGroupNameRegex.MatchString(fmt.Sprint(item["ServerGroupName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ServerGroupId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["ServerGroupStatus"].(string) {
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
			"address_ip_version":         object["AddressIPVersion"],
			"connection_drain_timeout":   formatInt(object["ConnectionDrainTimeout"]),
			"protocol":                   object["Protocol"],
			"related_load_balancer_ids":  object["RelatedLoadBalancerIds"],
			"scheduler":                  object["Scheduler"],
			"server_count":               formatInt(object["ServerCount"]),
			"server_group_name":          object["ServerGroupName"],
			"status":                     object["ServerGroupStatus"],
			"id":                         fmt.Sprint(object["ServerGroupId"]),
			"vpc_id":                     object["VpcId"],
			"server_group_type":          object["ServerGroupType"],
			"connection_drain":           object["ConnectionDrainEnabled"],
			"resource_group_id":          object["ResourceGroupId"],
			"preserve_client_ip_enabled": object["PreserveClientIpEnabled"],
		}
		if v, ok := object["HealthCheck"]; ok {
			healthCheckSli := make([]map[string]interface{}, 0)
			if len(v.(map[string]interface{})) > 0 {
				healthCheck := v.(map[string]interface{})
				healthCheckMap := make(map[string]interface{})
				healthCheckMap["health_check_connect_port"] = healthCheck["HealthCheckConnectPort"]
				healthCheckMap["health_check_connect_timeout"] = healthCheck["HealthCheckConnectTimeout"]
				healthCheckMap["health_check_domain"] = healthCheck["HealthCheckDomain"]
				healthCheckMap["health_check_enabled"] = healthCheck["HealthCheckEnabled"]
				healthCheckMap["health_check_interval"] = healthCheck["HealthCheckInterval"]
				healthCheckMap["health_check_type"] = healthCheck["HealthCheckType"]
				healthCheckMap["health_check_url"] = healthCheck["HealthCheckUrl"]
				healthCheckMap["healthy_threshold"] = healthCheck["HealthyThreshold"]
				healthCheckMap["unhealthy_threshold"] = healthCheck["UnhealthyThreshold"]
				healthCheckMap["http_check_method"] = healthCheck["HttpCheckMethod"]
				healthCheckMap["health_check_http_code"] = healthCheck["HealthCheckHttpCode"]
				healthCheckSli = append(healthCheckSli, healthCheckMap)
			}
			mapping["health_check"] = healthCheckSli
		}

		mapping["tags"] = tagsToMap(object["Tags"])
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ServerGroupName"])

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
