package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudEssScalingGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEssScalingGroupsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"active_scaling_configuration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"launch_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"launch_template_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"desired_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_instance_lifetime": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"stop_instance_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"on_demand_base_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cooldown_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"removal_policies": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"load_balancer_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"db_instance_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"vswitch_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lifecycle_state": {
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
						"health_check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az_balance": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"system_suspended": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"monitor_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"suspended_processes": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"multi_az_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_allocation_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_deletion_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"on_demand_percentage_above_base_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spot_instance_remedy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_desired_capacity": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allocation_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modification_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spot_instance_pools": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"init_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pending_wait_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"removing_wait_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protected_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"standby_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spot_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"stopped_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_instance_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"active_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pending_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"removing_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEssScalingGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PageSize":   requests.NewInteger(PageSizeLarge),
		"PageNumber": requests.NewInteger(1),
		"RegionId":   client.RegionId,
	}

	var allScalingGroups []interface{}
	for {
		response, err = client.RpcPost("Ess", "2014-08-28", "DescribeScalingGroups", nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_scaling_groups", "DescribeScalingGroups", AlibabaCloudSdkGoERROR)
		}

		v, err := jsonpath.Get("$.ScalingGroups.ScalingGroup", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "$.ScalingGroups.ScalingGroup", response)
		}

		addDebug("DescribeScalingGroups", response, request, request)
		if len(v.([]interface{})) < 1 {
			break
		}

		allScalingGroups = append(allScalingGroups, v.([]interface{})...)

		if len(v.([]interface{})) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(requests.Integer(fmt.Sprint(request["PageNumber"]))); err != nil {
			return WrapError(err)
		} else {
			request["PageNumber"] = page
		}
	}

	var filteredScalingGroupsTemp []interface{}

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}
	if okNameRegex || okIds {
		for _, group := range allScalingGroups {
			var object map[string]interface{}
			object = group.(map[string]interface{})
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if r != nil && !r.MatchString(object["ScalingGroupName"].(string)) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[object["ScalingGroupId"].(string)]; !ok {
					continue
				}
			}
			filteredScalingGroupsTemp = append(filteredScalingGroupsTemp, group)
		}
	} else {
		filteredScalingGroupsTemp = allScalingGroups
	}
	return scalingGroupsDescriptionAttribute(d, filteredScalingGroupsTemp, meta, client)
}

func scalingGroupsDescriptionAttribute(d *schema.ResourceData, scalingGroups []interface{}, meta interface{}, client *connectivity.AliyunClient) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	essService := EssService{client}

	for _, scalingGroup := range scalingGroups {
		var object map[string]interface{}
		object = scalingGroup.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":                           object["ScalingGroupId"],
			"name":                         object["ScalingGroupName"],
			"active_scaling_configuration": object["ActiveScalingConfigurationId"],
			"launch_template_id":           object["LaunchTemplateId"],
			"launch_template_version":      object["LaunchTemplateVersion"],
			"region_id":                    object["RegionId"],
			"min_size":                     object["MinSize"],
			"max_size":                     object["MaxSize"],
			"group_type":                   object["GroupType"],
			"monitor_group_id":             object["MonitorGroupId"],
			"cooldown_time":                object["DefaultCooldown"],
			"stop_instance_timeout":        object["StopInstanceTimeout"],
			"lifecycle_state":              object["LifecycleState"],
			"total_capacity":               object["TotalCapacity"],
			"active_capacity":              object["ActiveCapacity"],
			"enable_desired_capacity":      object["EnableDesiredCapacity"],
			"pending_capacity":             object["PendingCapacity"],
			"removing_capacity":            object["RemovingCapacity"],
			"total_instance_count":         object["TotalInstanceCount"],
			"init_capacity":                object["InitCapacity"],
			"scaling_policy":               object["ScalingPolicy"],
			"pending_wait_capacity":        object["PendingWaitCapacity"],
			"removing_wait_capacity":       object["RemovingWaitCapacity"],
			"protected_capacity":           object["ProtectedCapacity"],
			"standby_capacity":             object["StandbyCapacity"],
			"spot_capacity":                object["SpotCapacity"],
			"stopped_capacity":             object["StoppedCapacity"],
			"system_suspended":             object["SystemSuspended"],
			"vpc_id":                       object["VpcId"],
			"vswitch_id":                   object["VSwitchId"],
			"health_check_type":            object["HealthCheckType"],
			"group_deletion_protection":    object["GroupDeletionProtection"],
			"spot_instance_remedy":         object["SpotInstanceRemedy"],
			"modification_time":            object["ModificationTime"],
			"creation_time":                object["CreationTime"],
			"multi_az_policy":              object["MultiAZPolicy"],
			"resource_group_id":            object["ResourceGroupId"],
		}
		if object["DesiredCapacity"] != nil {
			mapping["desired_capacity"] = object["DesiredCapacity"]
		}
		if object["AzBalance"] != nil {
			mapping["az_balance"] = object["AzBalance"]
		}
		if object["SpotAllocationStrategy"] != nil {
			mapping["spot_allocation_strategy"] = object["SpotAllocationStrategy"]
		}
		if object["AllocationStrategy"] != nil {
			mapping["allocation_strategy"] = object["AllocationStrategy"]
		}
		if object["OnDemandBaseCapacity"] != nil {
			mapping["on_demand_base_capacity"] = object["OnDemandBaseCapacity"]
		}
		if object["OnDemandPercentageAboveBaseCapacity"] != nil {
			mapping["on_demand_percentage_above_base_capacity"] = object["OnDemandPercentageAboveBaseCapacity"]
		}
		if object["SpotInstancePools"] != nil {
			mapping["spot_instance_pools"] = object["SpotInstancePools"]
		}
		if object["MaxInstanceLifetime"] != nil {
			mapping["max_instance_lifetime"] = object["MaxInstanceLifetime"]
		}
		var dbIds []string
		if len(object["DBInstanceIds"].(map[string]interface{})["DBInstanceId"].([]interface{})) > 0 {
			for _, v := range object["DBInstanceIds"].(map[string]interface{})["DBInstanceId"].([]interface{}) {
				dbIds = append(dbIds, v.(string))
			}
			mapping["db_instance_ids"] = dbIds
		}

		var slbIds []string
		if len(object["LoadBalancerIds"].(map[string]interface{})["LoadBalancerId"].([]interface{})) > 0 {
			for _, v := range object["LoadBalancerIds"].(map[string]interface{})["LoadBalancerId"].([]interface{}) {
				slbIds = append(slbIds, v.(string))
			}
			mapping["loadbalancer_ids"] = slbIds
		}

		var polices []string
		if len(object["RemovalPolicies"].(map[string]interface{})["RemovalPolicy"].([]interface{})) > 0 {
			for _, v := range object["RemovalPolicies"].(map[string]interface{})["RemovalPolicy"].([]interface{}) {
				polices = append(polices, v.(string))
			}
			mapping["removal_policies"] = polices
		}

		var vswitchIds []string
		if object["VSwitchIds"] != nil && len(object["VSwitchIds"].(map[string]interface{})["VSwitchId"].([]interface{})) > 0 {
			for _, v := range object["VSwitchIds"].(map[string]interface{})["VSwitchId"].([]interface{}) {
				vswitchIds = append(vswitchIds, v.(string))
			}
			mapping["vswitch_ids"] = vswitchIds
		}

		var suspendedProcesses []string
		if object["SuspendedProcesses"] != nil && len(object["SuspendedProcesses"].(map[string]interface{})["SuspendedProcess"].([]interface{})) > 0 {
			for _, v := range object["SuspendedProcesses"].(map[string]interface{})["SuspendedProcess"].([]interface{}) {
				suspendedProcesses = append(suspendedProcesses, v.(string))
			}
			mapping["suspended_processes"] = suspendedProcesses
		}

		listTagResourcesObject, err := essService.ListTagResources(d.Id(), client)
		if err != nil {
			return WrapError(err)
		}
		mapping["tags"] = tagsToMap(listTagResourcesObject)

		ids = append(ids, object["ScalingGroupId"].(string))
		names = append(names, object["ScalingGroupName"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
