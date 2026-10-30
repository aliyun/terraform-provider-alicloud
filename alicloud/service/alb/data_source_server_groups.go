package alb

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServerGroupsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"server_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"server_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Available", "Configuring", "Provisioning"}, false),
			},
			"tags": tags.TagsSchema(),
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health_check_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"health_check_codes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"health_check_connect_port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"health_check_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_http_version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"health_check_method": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"health_check_timeout": {
										Type:     schema.TypeInt,
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
								},
							},
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"server_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"sticky_session_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cookie": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cookie_timeout": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"sticky_session_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"sticky_session_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListServerGroups"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if m, ok := d.GetOk("server_group_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("ServerGroupIds.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("server_group_name"); ok {
		request["ServerGroupNames.1"] = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		request["Tag"] = tags.TagsFromMap(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VpcId"] = v
	}
	request["MaxResults"] = helper.PageSizeLarge
	var objects []map[string]interface{}
	var serverGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errs.WrapError(err)
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			response, err = client.RpcPost("Alb", "2020-06-16", action, nil, request, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DataDefaultErrorMsg, "alicloud_alb_server_groups", action, errs.AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ServerGroups", response)
		totalCount, formatErr := convert.FormatInt(response["TotalCount"])
		if formatErr != nil {
			return formatErr
		}
		if totalCount != 0 && err != nil {
			return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.ServerGroups", response)
		}
		result, _ := resp.([]interface{})
		if resp != nil {
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
			"protocol":              object["Protocol"],
			"scheduler":             object["Scheduler"],
			"id":                    fmt.Sprint(object["ServerGroupId"]),
			"server_group_id":       fmt.Sprint(object["ServerGroupId"]),
			"server_group_name":     object["ServerGroupName"],
			"status":                object["ServerGroupStatus"],
			"sticky_session_config": object["StickySessionConfig"],
			"vpc_id":                object["VpcId"],
		}

		healthCheckConfigSli := make([]map[string]interface{}, 0)
		if len(object["HealthCheckConfig"].(map[string]interface{})) > 0 {
			healthCheckConfig := object["HealthCheckConfig"]
			healthCheckConfigMap := make(map[string]interface{})
			healthCheckConfigMap["health_check_codes"] = healthCheckConfig.(map[string]interface{})["HealthCheckCodes"]
			healthCheckConnectPort, err := convert.FormatInt(healthCheckConfig.(map[string]interface{})["HealthCheckConnectPort"])
			if err != nil {
				return err
			}
			healthCheckConfigMap["health_check_connect_port"] = healthCheckConnectPort
			healthCheckConfigMap["health_check_enabled"] = healthCheckConfig.(map[string]interface{})["HealthCheckEnabled"]
			healthCheckConfigMap["health_check_host"] = healthCheckConfig.(map[string]interface{})["HealthCheckHost"]
			healthCheckConfigMap["health_check_http_version"] = healthCheckConfig.(map[string]interface{})["HealthCheckHttpVersion"]
			healthCheckInterval, err := convert.FormatInt(healthCheckConfig.(map[string]interface{})["HealthCheckInterval"])
			if err != nil {
				return err
			}
			healthCheckConfigMap["health_check_interval"] = healthCheckInterval
			healthCheckConfigMap["health_check_method"] = healthCheckConfig.(map[string]interface{})["HealthCheckMethod"]
			healthCheckConfigMap["health_check_path"] = healthCheckConfig.(map[string]interface{})["HealthCheckPath"]
			healthCheckConfigMap["health_check_protocol"] = healthCheckConfig.(map[string]interface{})["HealthCheckProtocol"]
			healthCheckTimeout, err := convert.FormatInt(healthCheckConfig.(map[string]interface{})["HealthCheckTimeout"])
			if err != nil {
				return err
			}
			healthCheckConfigMap["health_check_timeout"] = healthCheckTimeout
			healthyThreshold, err := convert.FormatInt(healthCheckConfig.(map[string]interface{})["HealthyThreshold"])
			if err != nil {
				return err
			}
			healthCheckConfigMap["healthy_threshold"] = healthyThreshold
			unhealthyThreshold, err := convert.FormatInt(healthCheckConfig.(map[string]interface{})["UnhealthyThreshold"])
			if err != nil {
				return err
			}
			healthCheckConfigMap["unhealthy_threshold"] = unhealthyThreshold
			healthCheckConfigSli = append(healthCheckConfigSli, healthCheckConfigMap)
		}
		mapping["health_check_config"] = healthCheckConfigSli

		stickySessionConfigSli := make([]map[string]interface{}, 0)
		if len(object["StickySessionConfig"].(map[string]interface{})) > 0 {
			stickySessionConfig := object["StickySessionConfig"]
			stickySessionConfigMap := make(map[string]interface{})
			stickySessionConfigMap["cookie"] = stickySessionConfig.(map[string]interface{})["Cookie"]
			stickySessionConfigMap["cookie_timeout"] = stickySessionConfig.(map[string]interface{})["CookieTimeout"]
			stickySessionConfigMap["sticky_session_enabled"] = stickySessionConfig.(map[string]interface{})["StickySessionEnabled"]
			stickySessionConfigMap["sticky_session_type"] = stickySessionConfig.(map[string]interface{})["StickySessionType"]
			stickySessionConfigSli = append(stickySessionConfigSli, stickySessionConfigMap)
		}
		mapping["sticky_session_config"] = stickySessionConfigSli
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ServerGroupName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["ServerGroupId"])
		albService := AlbService{client}
		serversList, err := albService.ListServerGroupServers(id)
		if err != nil {
			return errs.WrapError(err)
		}

		servers := make([]map[string]interface{}, 0)
		for _, v := range serversList {
			if m1, ok := v.(map[string]interface{}); ok {

				portVal, err := convert.FormatInt(m1["Port"])
				if err != nil {
					return err
				}

				weightVal, err := convert.FormatInt(m1["Weight"])
				if err != nil {
					return err
				}
				temp1 := map[string]interface{}{
					"description": m1["Description"],
					"port":        portVal,
					"server_id":   m1["ServerId"],
					"server_ip":   m1["ServerIp"],
					"server_type": m1["ServerType"],
					"status":      m1["Status"],
					"weight":      weightVal,
				}
				servers = append(servers, temp1)
			}
		}
		mapping["servers"] = servers

		tagResp, _ := albService.ListTagResources(id, "servergroup")
		mapping["tags"] = tags.TagsToMap(tagResp)

		s = append(s, mapping)
	}

	d.SetId(helper.DataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errs.WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return errs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		helper.WriteToFile(output.(string), s)
	}

	return nil
}
