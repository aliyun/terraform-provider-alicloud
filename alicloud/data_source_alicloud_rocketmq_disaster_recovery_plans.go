package alicloud

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRocketmqDisasterRecoveryPlans() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRocketmqDisasterRecoveryPlansRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plans": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_sync_checkpoint": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auth_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"consumer_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"endpoint_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_role": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message_property": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"property_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"property_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"network_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"password": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vswitch_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"plan_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"plan_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sync_checkpoint_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRocketmqDisasterRecoveryPlansRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rocketmqServiceV2 := RocketmqServiceV2{client}

	filters := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		filters["instance_id"] = v.(string)
	}

	objects, err := rocketmqServiceV2.ListRocketmqDisasterRecoveryPlans(filters)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rocketmq_disaster_recovery_plans", "ListDisasterRecoveryPlans", AlibabaCloudSdkGoERROR)
	}

	var planNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		planNameRegex = r
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

	s := make([]map[string]interface{}, 0)
	ids := make([]string, 0)
	for _, object := range objects {
		item := object.(map[string]interface{})
		planId := fmt.Sprint(item["planId"])

		if planNameRegex != nil {
			if !planNameRegex.MatchString(fmt.Sprint(item["planName"])) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[planId]; !ok {
				continue
			}
		}

		instancesMaps := flattenRocketmqDisasterRecoveryPlanInstances(item["instances"])
		mapping := map[string]interface{}{
			"auto_sync_checkpoint":    item["autoSyncCheckpoint"],
			"create_time":             item["createTime"],
			"instances":               instancesMaps,
			"plan_desc":               item["planDesc"],
			"plan_id":                 planId,
			"plan_name":               item["planName"],
			"plan_type":               item["planType"],
			"region_id":               client.RegionId,
			"status":                  item["planStatus"],
			"sync_checkpoint_enabled": item["syncCheckpointEnabled"],
			"update_time":             item["updateTime"],
		}
		ids = append(ids, planId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("plans", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

func flattenRocketmqDisasterRecoveryPlanInstances(instancesRaw interface{}) []map[string]interface{} {
	instancesMaps := make([]map[string]interface{}, 0)
	if instancesRaw == nil {
		return instancesMaps
	}
	instances, ok := instancesRaw.([]interface{})
	if !ok {
		return instancesMaps
	}
	for _, instanceRaw := range instances {
		instanceMap := instanceRaw.(map[string]interface{})
		item := make(map[string]interface{})
		item["auth_type"] = instanceMap["authType"]
		item["consumer_group_id"] = instanceMap["consumerGroupId"]
		item["endpoint_url"] = instanceMap["endpointUrl"]
		item["instance_id"] = instanceMap["instanceId"]
		item["instance_role"] = instanceMap["instanceRole"]
		item["instance_type"] = instanceMap["instanceType"]
		item["network_type"] = instanceMap["networkType"]
		item["password"] = instanceMap["password"]
		item["region_id"] = instanceMap["regionId"]
		item["security_group_id"] = instanceMap["securityGroupId"]
		item["username"] = instanceMap["username"]
		item["vpc_id"] = instanceMap["vpcId"]
		item["vswitch_id"] = instanceMap["vSwitchId"]

		messagePropertyMaps := make([]map[string]interface{}, 0)
		if instanceMap["messageProperty"] != nil {
			mpRaw := instanceMap["messageProperty"].(map[string]interface{})
			if len(mpRaw) > 0 {
				mpMap := make(map[string]interface{})
				mpMap["property_key"] = mpRaw["propertyKey"]
				mpMap["property_value"] = mpRaw["propertyValue"]
				messagePropertyMaps = append(messagePropertyMaps, mpMap)
			}
		}
		item["message_property"] = messagePropertyMaps
		instancesMaps = append(instancesMaps, item)
	}
	return instancesMaps
}
