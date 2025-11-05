package alicloud

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingGroupCreate,
		Read:   resourceAliyunEssScalingGroupRead,
		Update: resourceAliyunEssScalingGroupUpdate,
		Delete: resourceAliyunEssScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(0, 2000),
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: IntBetween(0, 2000),
			},
			"stop_instance_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(30, 240),
			},
			"desired_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 2000),
			},
			"scaling_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_check_type": {
				Type:          schema.TypeString,
				Computed:      true,
				ValidateFunc:  StringInSlice([]string{"ECS", "NONE", "LOAD_BALANCER"}, false),
				Optional:      true,
				ConflictsWith: []string{"health_check_types"},
			},
			"scaling_policy": {
				Type:         schema.TypeString,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"recycle", "release", "forceRecycle", "forceRelease"}, false),
				Optional:     true,
			},
			"max_instance_lifetime": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntAtLeast(86400),
			},
			"default_cooldown": {
				Type:         schema.TypeInt,
				Default:      300,
				Optional:     true,
				ValidateFunc: IntBetween(0, 86400),
			},
			"vswitch_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.7.1, and new field 'vswitch_ids' can replace it.",
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"container_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"health_check_types": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"health_check_type"},
			},
			"vswitch_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"removal_policies": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				MaxItems: 2,
				MinItems: 1,
			},
			"alb_server_group": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alb_server_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"db_instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MinItems: 0,
			},
			"loadbalancer_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MinItems: 0,
			},
			"multi_az_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PRIORITY",
				ValidateFunc: StringInSlice([]string{"PRIORITY", "BALANCE", "COST_OPTIMIZED", "COMPOSABLE"}, false),
				ForceNew:     true,
			},
			"az_balance": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allocation_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"priority", "lowestPrice"}, false),
				Computed:     true,
			},
			"spot_allocation_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"priority", "lowestPrice"}, false),
				Computed:     true,
			},
			"on_demand_base_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 1000),
			},
			"on_demand_percentage_above_base_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"spot_instance_pools": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(1, 10),
			},
			"spot_instance_remedy": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"compensate_with_on_demand": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"capacity_options_on_demand_base_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 1000),
			},
			"capacity_options_on_demand_percentage_above_base_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"capacity_options_compensate_with_on_demand": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"capacity_options_spot_auto_replace_on_demand": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"capacity_options_price_comparison_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"launch_template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"launch_template_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"group_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"ECS", "ECI"}, false),
			},
			"protected_instances": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"launch_template_override": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weighted_capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"spot_price_limit": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunEssScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	var response map[string]interface{}
	request, err := buildAlicloudEssScalingGroupArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.RpcPost("Ess", "2014-08-28", "CreateScalingGroup", nil, request, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{Throttling, "IncorrectLoadBalancerHealthCheck", "IncorrectLoadBalancerStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scaling_group", "CreateScalingGroup", AlibabaCloudSdkGoERROR)
	}
	d.SetId(fmt.Sprint(response["ScalingGroupId"]))
	d.Set("alb_server_group", request["AlbServerGroup"])

	if err := essService.WaitForEssScalingGroup(d.Id(), Inactive, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	// enable group if use launchTemplate
	if request["LaunchTemplateId"] != "" && request["LaunchTemplateId"] != nil {
		enableGroupRequest := ess.CreateEnableScalingGroupRequest()
		enableGroupRequest.ScalingGroupId = d.Id()

		err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.EnableScalingGroup(enableGroupRequest)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus"}) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(enableGroupRequest.GetActionName(), raw, enableGroupRequest.RpcRequest, enableGroupRequest)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scaling_group", enableGroupRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliyunEssScalingGroupUpdate(d, meta)
}

func resourceAliyunEssScalingGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScalingGroupById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("min_size", object["MinSize"])
	d.Set("max_size", object["MaxSize"])
	if object["StopInstanceTimeout"] != nil {
		d.Set("stop_instance_timeout", object["StopInstanceTimeout"])
	}
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("desired_capacity", object["DesiredCapacity"])
	d.Set("scaling_group_name", object["ScalingGroupName"])
	d.Set("default_cooldown", object["DefaultCooldown"])
	if object["MaxInstanceLifetime"] != nil {
		d.Set("max_instance_lifetime", object["MaxInstanceLifetime"])
	}
	d.Set("multi_az_policy", object["MultiAZPolicy"])
	d.Set("az_balance", object["AzBalance"])
	d.Set("allocation_strategy", object["AllocationStrategy"])
	d.Set("spot_allocation_strategy", object["SpotAllocationStrategy"])
	d.Set("on_demand_base_capacity", object["OnDemandBaseCapacity"])
	d.Set("on_demand_percentage_above_base_capacity", object["OnDemandPercentageAboveBaseCapacity"])
	if object["SpotInstancePools"] != nil {
		d.Set("spot_instance_pools", object["SpotInstancePools"])
	}

	d.Set("spot_instance_remedy", object["SpotInstanceRemedy"])
	d.Set("compensate_with_on_demand", object["CompensateWithOnDemand"])
	d.Set("group_deletion_protection", object["GroupDeletionProtection"])
	var polices []string
	if len(object["RemovalPolicies"].(map[string]interface{})["RemovalPolicy"].([]interface{})) > 0 {
		for _, v := range object["RemovalPolicies"].(map[string]interface{})["RemovalPolicy"].([]interface{}) {
			polices = append(polices, v.(string))
		}
	}
	d.Set("removal_policies", polices)
	var dbIds []string
	if len(object["DBInstanceIds"].(map[string]interface{})["DBInstanceId"].([]interface{})) > 0 {
		for _, v := range object["DBInstanceIds"].(map[string]interface{})["DBInstanceId"].([]interface{}) {
			dbIds = append(dbIds, v.(string))
		}
	}
	d.Set("db_instance_ids", dbIds)

	var slbIds []string
	if len(object["LoadBalancerIds"].(map[string]interface{})["LoadBalancerId"].([]interface{})) > 0 {
		for _, v := range object["LoadBalancerIds"].(map[string]interface{})["LoadBalancerId"].([]interface{}) {
			slbIds = append(slbIds, v.(string))
		}
	}
	d.Set("loadbalancer_ids", slbIds)

	var vswitchIds []string
	if object["VSwitchIds"] != nil && len(object["VSwitchIds"].(map[string]interface{})["VSwitchId"].([]interface{})) > 0 {
		for _, v := range object["VSwitchIds"].(map[string]interface{})["VSwitchId"].([]interface{}) {
			vswitchIds = append(vswitchIds, v.(string))
		}
	}

	var healthCheckTypes []string
	if object["HealthCheckTypes"] != nil && len(object["HealthCheckTypes"].(map[string]interface{})["HealthCheckType"].([]interface{})) > 0 {
		for _, v := range object["HealthCheckTypes"].(map[string]interface{})["HealthCheckType"].([]interface{}) {
			healthCheckTypes = append(healthCheckTypes, v.(string))
		}
	}

	if v := object["LaunchTemplateOverrides"]; v != nil {
		result := make([]map[string]interface{}, 0)
		for _, i := range v.(map[string]interface{})["LaunchTemplateOverride"].([]interface{}) {
			launchTemplateOverride := i.(map[string]interface{})
			l := map[string]interface{}{
				"instance_type": launchTemplateOverride["InstanceType"],
			}
			if launchTemplateOverride["SpotPriceLimit"] != nil {
				spotPriceLimitFloatformat, _ := launchTemplateOverride["SpotPriceLimit"].(json.Number).Float64()
				spotPriceLimit, _ := strconv.ParseFloat(strconv.FormatFloat(spotPriceLimitFloatformat, 'f', 2, 64), 64)
				l["spot_price_limit"] = spotPriceLimit
			}
			if launchTemplateOverride["WeightedCapacity"] != nil {
				l["weighted_capacity"] = launchTemplateOverride["WeightedCapacity"]
			}
			result = append(result, l)
		}
		err := d.Set("launch_template_override", result)
		if err != nil {
			return WrapError(err)
		}
	}

	if v := object["CapacityOptions"]; v != nil {
		m := v.(map[string]interface{})
		if m["OnDemandBaseCapacity"] != nil {
			d.Set("capacity_options_on_demand_base_capacity", m["OnDemandBaseCapacity"])
		}
		if m["OnDemandPercentageAboveBaseCapacity"] != nil {
			d.Set("capacity_options_on_demand_percentage_above_base_capacity", m["OnDemandPercentageAboveBaseCapacity"])
		}
		if m["CompensateWithOnDemand"] != nil {
			d.Set("capacity_options_compensate_with_on_demand", m["CompensateWithOnDemand"])
		}
		if m["SpotAutoReplaceOnDemand"] != nil {
			d.Set("capacity_options_spot_auto_replace_on_demand", m["SpotAutoReplaceOnDemand"])
		}
		if m["PriceComparisonMode"] != nil {
			d.Set("capacity_options_price_comparison_mode", m["PriceComparisonMode"])
		}
	}

	if v := object["AlbServerGroups"]; v != nil {
		result := make([]map[string]interface{}, 0)
		if w, ok := d.GetOk("alb_server_group"); ok {
			albServerGroups := w.(*schema.Set).List()
			for _, rew := range albServerGroups {
				item := rew.(map[string]interface{})
				for _, i := range v.(map[string]interface{})["AlbServerGroup"].([]interface{}) {
					r := i.(map[string]interface{})
					uu, _ := r["Port"].(json.Number).Int64()
					if albServerGroupId, ok := item["alb_server_group_id"].(string); ok && albServerGroupId != "" {
						if r["AlbServerGroupId"].(string) == albServerGroupId && int64(item["port"].(int)) == uu {
							l := map[string]interface{}{
								"alb_server_group_id": r["AlbServerGroupId"],
								"weight":              r["Weight"],
								"port":                r["Port"],
							}
							result = append(result, l)
						}
					}
				}
				err := d.Set("alb_server_group", result)
				if err != nil {
					return WrapError(err)
				}
			}
		}
	}

	d.Set("vswitch_ids", vswitchIds)
	d.Set("launch_template_id", object["LaunchTemplateId"])
	d.Set("launch_template_version", object["LaunchTemplateVersion"])
	d.Set("group_type", object["GroupType"])
	d.Set("health_check_type", object["HealthCheckType"])
	if object["HealthCheckType"] == nil {
		d.Set("health_check_types", healthCheckTypes)
	}
	d.Set("scaling_policy", object["ScalingPolicy"])

	listTagResourcesObject, err := essService.ListTagResources(d.Id(), client)
	if err != nil {
		return WrapError(err)
	}

	d.Set("tags", tagsToMap(listTagResourcesObject))

	instances, _ := essService.DescribeInstances(d.Id(), "Protected")
	var protectedInstances []string
	for _, v := range instances {
		protectedInstances = append(protectedInstances, v.InstanceId)
	}
	d.Set("protected_instances", protectedInstances)

	return nil
}

func resourceAliyunEssScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	action := "ModifyScalingGroup"
	request := map[string]interface{}{
		"ScalingGroupId": d.Id(),
	}

	essService := EssService{client}
	//开启 允许部分属性修改
	d.Partial(true)

	if d.HasChange("tags") {
		if err := essService.SetResourceTags(d, d.Id(), client); err != nil {
			return WrapError(err)
		}
	}
	if d.HasChange("resource_group_id") {
		if err := essService.ChangeResourceGroup(d, d.Id(), client); err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("scaling_group_name") {
		request["ScalingGroupName"] = d.Get("scaling_group_name").(string)
	}

	if d.HasChange("health_check_type") {
		request["HealthCheckType"] = d.Get("health_check_type").(string)
	}

	if d.HasChange("scaling_policy") {
		request["ScalingPolicy"] = d.Get("scaling_policy").(string)
	}

	if d.HasChange("min_size") {
		request["MinSize"] = requests.NewInteger(d.Get("min_size").(int))
	}

	if d.HasChange("max_size") {
		request["MaxSize"] = requests.NewInteger(d.Get("max_size").(int))
	}

	if d.HasChange("stop_instance_timeout") {
		if v, ok := d.GetOkExists("stop_instance_timeout"); ok {
			request["StopInstanceTimeout"] = requests.NewInteger(v.(int))
		}
	}

	if d.HasChange("desired_capacity") {
		if v, ok := d.GetOkExists("desired_capacity"); ok {
			request["DesiredCapacity"] = requests.NewInteger(v.(int))
		}
	}
	if d.HasChange("max_instance_lifetime") {
		if v, ok := d.GetOkExists("max_instance_lifetime"); ok {
			request["MaxInstanceLifetime"] = requests.NewInteger(v.(int))
		}
	}
	if d.HasChange("default_cooldown") {
		request["DefaultCooldown"] = requests.NewInteger(d.Get("default_cooldown").(int))
	}

	if d.HasChange("vswitch_ids") {
		vSwitchIds := expandStringList(d.Get("vswitch_ids").(*schema.Set).List())
		request["VSwitchIds"] = &vSwitchIds
	}

	if d.HasChange("health_check_types") {
		healthCheckTypes := expandStringList(d.Get("health_check_types").([]interface{}))
		request["HealthCheckTypes"] = &healthCheckTypes
	}

	if d.HasChange("removal_policies") {
		policyies := expandStringList(d.Get("removal_policies").([]interface{}))
		for i, p := range policyies {
			request[fmt.Sprintf("RemovalPolicy.%d", i+1)] = p
		}
	}

	if d.HasChange("on_demand_base_capacity") {
		request["OnDemandBaseCapacity"] = requests.NewInteger(d.Get("on_demand_base_capacity").(int))
	}

	if d.HasChange("on_demand_percentage_above_base_capacity") {
		request["OnDemandPercentageAboveBaseCapacity"] = requests.NewInteger(d.Get("on_demand_percentage_above_base_capacity").(int))
	}

	if d.HasChange("spot_instance_pools") {
		request["SpotInstancePools"] = requests.NewInteger(d.Get("spot_instance_pools").(int))
	}

	if d.HasChange("capacity_options_on_demand_base_capacity") {
		if v, ok := d.GetOkExists("capacity_options_on_demand_base_capacity"); ok {
			request["CapacityOptions.OnDemandBaseCapacity"] = requests.NewInteger(v.(int))
		}
	}

	if d.HasChange("capacity_options_on_demand_percentage_above_base_capacity") {
		if v, ok := d.GetOkExists("capacity_options_on_demand_percentage_above_base_capacity"); ok {
			request["CapacityOptions.OnDemandPercentageAboveBaseCapacity"] = requests.NewInteger(v.(int))
		}
	}

	if d.HasChange("capacity_options_compensate_with_on_demand") {
		if v, ok := d.GetOkExists("capacity_options_compensate_with_on_demand"); ok {
			request["CapacityOptions.CompensateWithOnDemand"] = requests.NewBoolean(v.(bool))
		}
	}

	if d.HasChange("capacity_options_spot_auto_replace_on_demand") {
		if v, ok := d.GetOkExists("capacity_options_spot_auto_replace_on_demand"); ok {
			request["CapacityOptions.SpotAutoReplaceOnDemand"] = requests.NewBoolean(v.(bool))
		}
	}

	if d.HasChange("capacity_options_price_comparison_mode") {
		if v, ok := d.GetOk("capacity_options_price_comparison_mode"); ok && v.(string) != "" {
			request["CapacityOptions.PriceComparisonMode"] = d.Get("capacity_options_price_comparison_mode").(string)
		}
	}

	if d.HasChange("spot_instance_remedy") {
		request["SpotInstanceRemedy"] = requests.NewBoolean(d.Get("spot_instance_remedy").(bool))
	}

	if d.HasChange("compensate_with_on_demand") {
		if v, ok := d.GetOkExists("compensate_with_on_demand"); ok {
			request["CompensateWithOnDemand"] = requests.NewBoolean(v.(bool))
		}
	}

	if d.HasChange("az_balance") {
		request["AzBalance"] = requests.NewBoolean(d.Get("az_balance").(bool))
	}

	if d.HasChange("allocation_strategy") {
		request["AllocationStrategy"] = d.Get("allocation_strategy").(string)
	}

	if d.HasChange("spot_allocation_strategy") {
		request["SpotAllocationStrategy"] = d.Get("spot_allocation_strategy").(string)
	}

	if d.HasChange("group_deletion_protection") {
		request["GroupDeletionProtection"] = requests.NewBoolean(d.Get("group_deletion_protection").(bool))
	}

	if d.HasChange("launch_template_id") || d.HasChange("launch_template_version") {
		request["LaunchTemplateId"] = d.Get("launch_template_id").(string)
		request["LaunchTemplateVersion"] = d.Get("launch_template_version").(string)

	}

	if d.HasChange("launch_template_override") {
		v, ok := d.GetOk("launch_template_override")
		if ok {
			launchTemplateOverrides := make([]map[string]interface{}, 0)
			for _, rew := range v.(*schema.Set).List() {
				item := rew.(map[string]interface{})
				l := map[string]interface{}{
					"InstanceType": item["instance_type"].(string),
				}
				if item["spot_price_limit"].(float64) != 0 {
					l["SpotPriceLimit"] = strconv.FormatFloat(item["spot_price_limit"].(float64), 'f', 2, 64)
				}
				if item["weighted_capacity"].(int) != 0 {
					l["WeightedCapacity"] = strconv.Itoa(item["weighted_capacity"].(int))
				}
				launchTemplateOverrides = append(launchTemplateOverrides, l)
			}
			request["LaunchTemplateVersion"] = d.Get("launch_template_version").(string)
			request["LaunchTemplateId"] = d.Get("launch_template_id").(string)
			request["LaunchTemplateOverride"] = &launchTemplateOverrides
		}
	}

	_, err = client.RpcPost("Ess", "2014-08-28", action, nil, request, false)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scaling_group", "ModifyScalingGroup", AlibabaCloudSdkGoERROR)
	}

	if d.HasChange("loadbalancer_ids") {
		oldLoadbalancers, newLoadbalancers := d.GetChange("loadbalancer_ids")
		err = attachOrDetachLoadbalancers(d, client, oldLoadbalancers.(*schema.Set), newLoadbalancers.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("alb_server_group") {
		oldAlbServerGroups, newAlbServerGroups := d.GetChange("alb_server_group")
		err = attachOrDetachAlbServerGroups(d, client, oldAlbServerGroups.(*schema.Set), newAlbServerGroups.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("db_instance_ids") {
		oldDbInstanceIds, newDbInstanceIds := d.GetChange("db_instance_ids")
		err = attachOrDetachDbInstances(d, client, oldDbInstanceIds.(*schema.Set), newDbInstanceIds.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
	}

	if d.HasChange("protected_instances") {
		oldProtectedInstances, newProtectedInstances := d.GetChange("protected_instances")
		err = setProtectedInstances(d, client, oldProtectedInstances.(*schema.Set), newProtectedInstances.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
	}

	d.Partial(false)
	return resourceAliyunEssScalingGroupRead(d, meta)
}

func resourceAliyunEssScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	request := ess.CreateDeleteScalingGroupRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = d.Id()
	request.ForceDelete = requests.NewBoolean(true)
	err := resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			if IsExpectedErrors(err, []string{"InternalError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultLongTimeout))
}

func buildAlicloudEssScalingGroupArgs(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"RegionId":        client.RegionId,
		"MinSize":         d.Get("min_size"),
		"MaxSize":         d.Get("max_size"),
		"DefaultCooldown": d.Get("default_cooldown"),
		"MultiAZPolicy":   d.Get("multi_az_policy"),
		"GroupType":       d.Get("group_type"),
	}

	slbService := SlbService{client}

	if v, ok := d.GetOk("scaling_group_name"); ok && v.(string) != "" {
		request["ScalingGroupName"] = v
	}

	if v, ok := d.GetOk("allocation_strategy"); ok && v.(string) != "" {
		request["AllocationStrategy"] = v
	}
	if v, ok := d.GetOk("spot_allocation_strategy"); ok && v.(string) != "" {
		request["SpotAllocationStrategy"] = v
	}

	if v, ok := d.GetOk("az_balance"); ok {
		request["AzBalance"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}

	if v, ok := d.GetOk("container_group_id"); ok {
		request["ContainerGroupId"] = v
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		count := 1
		for _, value := range v.(*schema.Set).List() {
			request[fmt.Sprintf("VSwitchIds.%d", count)] = value
			count++
		}
	}

	if v, ok := d.GetOk("health_check_types"); ok {
		count := 1
		for _, value := range convertToInterfaceArray(v) {
			request[fmt.Sprintf("HealthCheckTypes.%d", count)] = value
			count++
		}
	}

	if v, ok := d.GetOk("alb_server_group"); ok {
		albServerGroupsMaps := make([]map[string]interface{}, 0)
		albServerGroups := v.(*schema.Set).List()
		for _, rew := range albServerGroups {
			albServerGroupsMap := make(map[string]interface{})
			item := rew.(map[string]interface{})
			if albServerGroupId, ok := item["alb_server_group_id"].(string); ok && albServerGroupId != "" {
				albServerGroupsMap["AlbServerGroupId"] = albServerGroupId
			}
			albServerGroupsMap["Weight"] = item["weight"]
			albServerGroupsMap["Port"] = item["port"]
			albServerGroupsMaps = append(albServerGroupsMaps, albServerGroupsMap)
		}
		request["AlbServerGroup"] = albServerGroupsMaps
	}

	if dbs, ok := d.GetOk("db_instance_ids"); ok {
		request["DBInstanceIds"] = convertListToJsonString(dbs.(*schema.Set).List())
	}

	if lbs, ok := d.GetOk("loadbalancer_ids"); ok {
		for _, lb := range lbs.(*schema.Set).List() {
			if err := slbService.WaitForSlb(lb.(string), Active, DefaultTimeout); err != nil {
				return nil, WrapError(err)
			}
		}
		request["LoadBalancerIds"] = convertListToJsonString(lbs.(*schema.Set).List())
	}

	if v, ok := d.GetOk("desired_capacity"); ok {
		request["DesiredCapacity"] = v
	}

	if v, ok := d.GetOk("stop_instance_timeout"); ok {
		request["StopInstanceTimeout"] = v
	}

	if v, ok := d.GetOk("max_instance_lifetime"); ok {
		request["MaxInstanceLifetime"] = v
	}

	request["OnDemandBaseCapacity"] = d.Get("on_demand_base_capacity")

	request["OnDemandPercentageAboveBaseCapacity"] = d.Get("on_demand_percentage_above_base_capacity")

	if v, ok := d.GetOk("spot_instance_pools"); ok {
		request["SpotInstancePools"] = v
	}

	if v, ok := d.GetOkExists("capacity_options_on_demand_base_capacity"); ok {
		request["CapacityOptions.OnDemandBaseCapacity"] = v
	}
	if v, ok := d.GetOkExists("capacity_options_on_demand_percentage_above_base_capacity"); ok {
		request["CapacityOptions.OnDemandPercentageAboveBaseCapacity"] = v
	}

	if v, ok := d.GetOkExists("capacity_options_compensate_with_on_demand"); ok {
		request["CapacityOptions.CompensateWithOnDemand"] = v
	}

	if v, ok := d.GetOkExists("capacity_options_spot_auto_replace_on_demand"); ok {
		request["CapacityOptions.SpotAutoReplaceOnDemand"] = v
	}

	if v, ok := d.GetOk("capacity_options_price_comparison_mode"); ok && v.(string) != "" {
		request["CapacityOptions.PriceComparisonMode"] = v
	}

	if v, ok := d.GetOk("spot_instance_remedy"); ok {
		request["SpotInstanceRemedy"] = v
	}

	if v, ok := d.GetOkExists("compensate_with_on_demand"); ok {
		request["CompensateWithOnDemand"] = v
	}

	if v, ok := d.GetOk("health_check_type"); ok {
		request["HealthCheckType"] = v
	}

	if v, ok := d.GetOk("scaling_policy"); ok {
		request["ScalingPolicy"] = v
	}

	if v, ok := d.GetOk("group_deletion_protection"); ok {
		request["GroupDeletionProtection"] = v
	}

	if v, ok := d.GetOk("launch_template_id"); ok {
		request["LaunchTemplateId"] = v
	}

	if v, ok := d.GetOk("launch_template_override"); ok {
		launchTemplateOverridesMaps := make([]map[string]interface{}, 0)
		launchTemplateOverrides := v.(*schema.Set).List()
		for _, rew := range launchTemplateOverrides {
			launchTemplateOverridesMap := make(map[string]interface{})
			item := rew.(map[string]interface{})

			if instanceType, ok := item["instance_type"].(string); ok && instanceType != "" {
				launchTemplateOverridesMap["InstanceType"] = instanceType
			}
			if item["spot_price_limit"].(float64) != 0 {
				launchTemplateOverridesMap["SpotPriceLimit"] = strconv.FormatFloat(item["spot_price_limit"].(float64), 'f', 2, 64)
			}
			if item["weighted_capacity"].(int) != 0 {
				launchTemplateOverridesMap["WeightedCapacity"] = item["weighted_capacity"].(int)
			}
			launchTemplateOverridesMaps = append(launchTemplateOverridesMaps, launchTemplateOverridesMap)
		}
		request["LaunchTemplateOverride"] = launchTemplateOverridesMaps
	}

	if v, ok := d.GetOk("launch_template_version"); ok {
		request["LaunchTemplateVersion"] = v
	}

	return request, nil
}

func attachOrDetachLoadbalancers(d *schema.ResourceData, client *connectivity.AliyunClient, oldLoadbalancerSet *schema.Set, newLoadbalancerSet *schema.Set) error {
	detachLoadbalancerSet := oldLoadbalancerSet.Difference(newLoadbalancerSet)
	attachLoadbalancerSet := newLoadbalancerSet.Difference(oldLoadbalancerSet)
	if attachLoadbalancerSet.Len() > 0 {
		var subLists = partition(attachLoadbalancerSet, int(AttachDetachLoadbalancersBatchsize))
		for _, subList := range subLists {
			attachLoadbalancersRequest := ess.CreateAttachLoadBalancersRequest()
			attachLoadbalancersRequest.RegionId = client.RegionId
			attachLoadbalancersRequest.ScalingGroupId = d.Id()
			attachLoadbalancersRequest.ForceAttach = requests.NewBoolean(true)
			attachLoadbalancersRequest.LoadBalancer = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.AttachLoadBalancers(attachLoadbalancersRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), attachLoadbalancersRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(attachLoadbalancersRequest.GetActionName(), raw, attachLoadbalancersRequest.RpcRequest, attachLoadbalancersRequest)
		}
	}
	if detachLoadbalancerSet.Len() > 0 {
		var subLists = partition(detachLoadbalancerSet, int(AttachDetachLoadbalancersBatchsize))
		for _, subList := range subLists {
			detachLoadbalancersRequest := ess.CreateDetachLoadBalancersRequest()
			detachLoadbalancersRequest.RegionId = client.RegionId
			detachLoadbalancersRequest.ScalingGroupId = d.Id()
			detachLoadbalancersRequest.ForceDetach = requests.NewBoolean(false)
			detachLoadbalancersRequest.LoadBalancer = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DetachLoadBalancers(detachLoadbalancersRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), detachLoadbalancersRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(detachLoadbalancersRequest.GetActionName(), raw, detachLoadbalancersRequest.RpcRequest, detachLoadbalancersRequest)
		}
	}
	return nil
}

func setProtectedInstances(d *schema.ResourceData, client *connectivity.AliyunClient, oldInstances *schema.Set, newInstances *schema.Set) error {
	unprotected := oldInstances.Difference(newInstances)
	protected := newInstances.Difference(oldInstances)

	request := ess.CreateSetInstancesProtectionRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = d.Id()

	if protected.Len() > 0 {
		var subLists = partition(protected, 20)
		for _, subList := range subLists {
			request.InstanceId = &subList
			request.ProtectedFromScaleIn = requests.Boolean(strconv.FormatBool(true))
			_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.SetInstancesProtection(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}
	}

	if unprotected.Len() > 0 {
		var subLists = partition(unprotected, 20)
		for _, subList := range subLists {
			request.InstanceId = &subList
			request.ProtectedFromScaleIn = requests.Boolean(strconv.FormatBool(false))
			_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.SetInstancesProtection(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
		}
	}

	return nil
}

func attachOrDetachAlbServerGroups(d *schema.ResourceData, client *connectivity.AliyunClient, oldAlbServerGroupSet *schema.Set, newAlbServerGroupSet *schema.Set) error {
	detachAlbServerGroupSet := oldAlbServerGroupSet.Difference(newAlbServerGroupSet).List()
	attachAlbServerGroupSet := newAlbServerGroupSet.Difference(oldAlbServerGroupSet).List()
	var response map[string]interface{}
	var err error
	if len(detachAlbServerGroupSet) > 0 {
		var subLists = SplitSlice(detachAlbServerGroupSet, int(AttachDetachAlbServerGroupBatchsize))
		for _, subList := range subLists {
			action := "DetachAlbServerGroups"
			albRequest := map[string]interface{}{
				"ScalingGroupId": d.Id(),
				"ForceDetach":    true,
				"RegionId":       client.RegionId,
			}
			albServerGroupsMaps := make([]map[string]interface{}, 0)
			for _, rew := range subList {
				albServerGroupsMap := make(map[string]interface{})
				item := rew.(map[string]interface{})
				if albServerGroupId, ok := item["alb_server_group_id"].(string); ok && albServerGroupId != "" {
					albServerGroupsMap["AlbServerGroupId"] = albServerGroupId
				}
				albServerGroupsMap["Port"] = item["port"]
				albServerGroupsMaps = append(albServerGroupsMaps, albServerGroupsMap)
			}
			albRequest["AlbServerGroup"] = albServerGroupsMaps
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ess", "2014-08-28", action, nil, albRequest, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus"}) || NeedRetry(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			addDebug(action, response, albRequest)
		}
	}

	if len(attachAlbServerGroupSet) > 0 {
		var subLists = SplitSlice(attachAlbServerGroupSet, int(AttachDetachAlbServerGroupBatchsize))
		for _, subList := range subLists {
			action := "AttachAlbServerGroups"
			albRequest := map[string]interface{}{
				"ScalingGroupId": d.Id(),
				"ForceAttach":    true,
				"RegionId":       client.RegionId,
			}
			albServerGroupsMaps := make([]map[string]interface{}, 0)
			for _, rew := range subList {
				albServerGroupsMap := make(map[string]interface{})
				item := rew.(map[string]interface{})
				if albServerGroupId, ok := item["alb_server_group_id"].(string); ok && albServerGroupId != "" {
					albServerGroupsMap["AlbServerGroupId"] = albServerGroupId
				}
				albServerGroupsMap["Port"] = item["port"]
				albServerGroupsMap["Weight"] = item["weight"]
				albServerGroupsMaps = append(albServerGroupsMaps, albServerGroupsMap)
			}
			albRequest["AlbServerGroup"] = albServerGroupsMaps
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = client.RpcPost("Ess", "2014-08-28", action, nil, albRequest, true)
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectScalingGroupStatus"}) || NeedRetry(err) {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			addDebug(action, response, albRequest)
		}
	}
	return nil
}

func attachOrDetachDbInstances(d *schema.ResourceData, client *connectivity.AliyunClient, oldDbInstanceIdSet *schema.Set, newDbInstanceIdSet *schema.Set) error {
	detachDbInstanceSet := oldDbInstanceIdSet.Difference(newDbInstanceIdSet)
	attachDbInstanceSet := newDbInstanceIdSet.Difference(oldDbInstanceIdSet)
	if attachDbInstanceSet.Len() > 0 {
		var subLists = partition(attachDbInstanceSet, int(AttachDetachDbinstancesBatchsize))
		for _, subList := range subLists {
			attachDbInstancesRequest := ess.CreateAttachDBInstancesRequest()
			attachDbInstancesRequest.RegionId = client.RegionId
			attachDbInstancesRequest.ScalingGroupId = d.Id()
			attachDbInstancesRequest.ForceAttach = requests.NewBoolean(true)
			attachDbInstancesRequest.DBInstance = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.AttachDBInstances(attachDbInstancesRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), attachDbInstancesRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(attachDbInstancesRequest.GetActionName(), raw, attachDbInstancesRequest.RpcRequest, attachDbInstancesRequest)
		}
	}
	if detachDbInstanceSet.Len() > 0 {
		var subLists = partition(detachDbInstanceSet, int(AttachDetachDbinstancesBatchsize))
		for _, subList := range subLists {
			detachDbInstancesRequest := ess.CreateDetachDBInstancesRequest()
			detachDbInstancesRequest.RegionId = client.RegionId
			detachDbInstancesRequest.ScalingGroupId = d.Id()
			detachDbInstancesRequest.ForceDetach = requests.NewBoolean(true)
			detachDbInstancesRequest.DBInstance = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DetachDBInstances(detachDbInstancesRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), detachDbInstancesRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(detachDbInstancesRequest.GetActionName(), raw, detachDbInstancesRequest.RpcRequest, detachDbInstancesRequest)
		}
	}
	return nil
}

func partition(instanceIds *schema.Set, batchSize int) [][]string {
	var res [][]string
	size := instanceIds.Len()
	batchCount := int(math.Ceil(float64(size) / float64(batchSize)))
	idList := expandStringList(instanceIds.List())
	for i := 1; i <= batchCount; i++ {
		fromIndex := batchSize * (i - 1)
		toIndex := int(math.Min(float64(batchSize*i), float64(size)))
		subList := idList[fromIndex:toIndex]
		res = append(res, subList)
	}
	return res
}
