package alicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAliCloudAckClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudAckClusterRead,
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
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_spec": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_mode": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"cluster_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deletion_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ip_stack": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"maintenance_window": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"recurrence": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"maintenance_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weekly_period": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"duration": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"node_cidr_mask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operation_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_auto_upgrade": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"channel": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"pod_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"profile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAliCloudAckClusterRead(d *schema.ResourceData, meta interface{}) error {
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
	var query map[string]*string
	// DescribeClustersForRegion
	action := fmt.Sprintf("/regions/%s/clusters", client.RegionId)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	query["region_id"] = StringPointer(client.RegionId)
	query["cluster_id"] = StringPointer(d.Get("cluster_id").(string))
	if v, ok := d.GetOk("cluster_id"); ok {
		query["cluster_id"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		query["name"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("cluster_spec"); ok {
		query["cluster_spec"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		query["cluster_type"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("profile"); ok {
		query["profile"] = StringPointer(v.(string))
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["page_size"] = StringPointer(strconv.Itoa(PageSizeLarge))
	query["page_number"] = StringPointer("1")
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaGet("CS", "2015-12-15", action, query, nil, nil)

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

		resp, _ := jsonpath.Get("$.clusters[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["cluster_id"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < PageSizeLarge {
			break
		}
		pageNum, _ := strconv.Atoi(*query["page_number"])
		query["page_number"] = StringPointer(strconv.Itoa(pageNum + 1))
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["cluster_id"]

		mapping["cluster_domain"] = objectRaw["cluster_domain"]
		mapping["cluster_name"] = objectRaw["name"]
		mapping["cluster_spec"] = objectRaw["cluster_spec"]
		mapping["cluster_type"] = objectRaw["cluster_type"]
		mapping["current_version"] = objectRaw["current_version"]
		mapping["deletion_protection"] = objectRaw["deletion_protection"]
		mapping["ip_stack"] = objectRaw["ip_stack"]
		mapping["pod_cidr"] = objectRaw["container_cidr"]
		mapping["profile"] = objectRaw["profile"]
		mapping["proxy_mode"] = objectRaw["proxy_mode"]
		mapping["region_id"] = objectRaw["region_id"]
		mapping["resource_group_id"] = objectRaw["resource_group_id"]
		mapping["security_group_id"] = objectRaw["security_group_id"]
		mapping["service_cidr"] = objectRaw["service_cidr"]
		mapping["state"] = objectRaw["state"]
		mapping["timezone"] = objectRaw["timezone"]
		mapping["vpc_id"] = objectRaw["vpc_id"]
		mapping["cluster_id"] = objectRaw["cluster_id"]

		tagsMaps := objectRaw["tags"]
		mapping["tags"] = tagsToMap(tagsMaps)
		vswitch_idsRaw := make([]interface{}, 0)
		if objectRaw["vswitch_ids"] != nil {
			vswitch_idsRaw = convertToInterfaceArray(objectRaw["vswitch_ids"])
		}

		mapping["vswitch_ids"] = vswitch_idsRaw

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, objectRaw["name"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(objectRaw["cluster_id"])
		mapping, err = dataSourceAliCloudAckClusterReadDescription(d, id, mapping, meta)
		if err != nil {
			return WrapError(err)
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("clusters", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func dataSourceAliCloudAckClusterReadDescription(d *schema.ResourceData, id string, object map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)

	ackServiceV2 := AckServiceV2{client}
	getResp, err := ackServiceV2.DescribeAckCluster(id)
	if err != nil {
		return nil, WrapError(err)
	}

	// Merge additional fields from Get API response to mapping
	// Reuse the response mapping template from Resource's read function
	mapping := object
	objectRaw := getResp

	mapping["cluster_domain"] = objectRaw["cluster_domain"]
	mapping["cluster_name"] = objectRaw["name"]
	mapping["cluster_spec"] = objectRaw["cluster_spec"]
	mapping["cluster_type"] = objectRaw["cluster_type"]
	mapping["current_version"] = objectRaw["current_version"]
	mapping["deletion_protection"] = objectRaw["deletion_protection"]
	mapping["ip_stack"] = objectRaw["ip_stack"]
	mapping["node_cidr_mask"] = objectRaw["node_cidr_mask"]
	mapping["pod_cidr"] = objectRaw["container_cidr"]
	mapping["profile"] = objectRaw["profile"]
	mapping["proxy_mode"] = objectRaw["proxy_mode"]
	mapping["region_id"] = objectRaw["region_id"]
	mapping["resource_group_id"] = objectRaw["resource_group_id"]
	mapping["security_group_id"] = objectRaw["security_group_id"]
	mapping["service_cidr"] = objectRaw["service_cidr"]
	mapping["state"] = objectRaw["state"]
	mapping["timezone"] = objectRaw["timezone"]
	mapping["vpc_id"] = objectRaw["vpc_id"]
	mapping["cluster_id"] = objectRaw["cluster_id"]

	autoModeMaps := make([]map[string]interface{}, 0)
	autoModeMap := make(map[string]interface{})
	auto_modeRaw := make(map[string]interface{})
	if objectRaw["auto_mode"] != nil {
		auto_modeRaw = objectRaw["auto_mode"].(map[string]interface{})
	}
	if len(auto_modeRaw) > 0 {
		autoModeMap["enabled"] = auto_modeRaw["enable"]

		autoModeMaps = append(autoModeMaps, autoModeMap)
	}
	mapping["auto_mode"] = autoModeMaps
	maintenanceWindowMaps := make([]map[string]interface{}, 0)
	maintenanceWindowMap := make(map[string]interface{})
	maintenance_windowRaw := make(map[string]interface{})
	if objectRaw["maintenance_window"] != nil {
		maintenance_windowRaw = objectRaw["maintenance_window"].(map[string]interface{})
	}
	if len(maintenance_windowRaw) > 0 {
		maintenanceWindowMap["duration"] = maintenance_windowRaw["duration"]
		maintenanceWindowMap["enable"] = maintenance_windowRaw["enable"]
		maintenanceWindowMap["maintenance_time"] = maintenance_windowRaw["maintenance_time"]
		maintenanceWindowMap["recurrence"] = maintenance_windowRaw["recurrence"]
		maintenanceWindowMap["weekly_period"] = maintenance_windowRaw["weekly_period"]

		maintenanceWindowMaps = append(maintenanceWindowMaps, maintenanceWindowMap)
	}
	mapping["maintenance_window"] = maintenanceWindowMaps
	operationPolicyMaps := make([]map[string]interface{}, 0)
	operationPolicyMap := make(map[string]interface{})
	cluster_auto_upgradeRawObj, _ := jsonpath.Get("$.operation_policy.cluster_auto_upgrade", objectRaw)
	cluster_auto_upgradeRaw := make(map[string]interface{})
	if cluster_auto_upgradeRawObj != nil {
		cluster_auto_upgradeRaw = cluster_auto_upgradeRawObj.(map[string]interface{})
	}
	if len(cluster_auto_upgradeRaw) > 0 {

		clusterAutoUpgradeMaps := make([]map[string]interface{}, 0)
		clusterAutoUpgradeMap := make(map[string]interface{})
		clusterAutoUpgradeMap["channel"] = cluster_auto_upgradeRaw["channel"]
		clusterAutoUpgradeMap["enabled"] = cluster_auto_upgradeRaw["enabled"]
		clusterAutoUpgradeMaps = append(clusterAutoUpgradeMaps, clusterAutoUpgradeMap)
		operationPolicyMap["cluster_auto_upgrade"] = clusterAutoUpgradeMaps
		operationPolicyMaps = append(operationPolicyMaps, operationPolicyMap)
	}
	mapping["operation_policy"] = operationPolicyMaps
	tagsMaps := objectRaw["tags"]
	mapping["tags"] = tagsToMap(tagsMaps)
	vswitch_idsRaw := make([]interface{}, 0)
	if objectRaw["vswitch_ids"] != nil {
		vswitch_idsRaw = convertToInterfaceArray(objectRaw["vswitch_ids"])
	}

	mapping["vswitch_ids"] = vswitch_idsRaw

	return mapping, nil
}
