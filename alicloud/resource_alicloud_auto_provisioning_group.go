package alicloud

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"
)

func resourceAlicloudAutoProvisioningGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAutoProvisioningGroupCreate,
		Read:   resourceAlicloudAutoProvisioningGroupRead,
		Update: resourceAlicloudAutoProvisioningGroupUpdate,
		Delete: resourceAlicloudAutoProvisioningGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"total_target_capacity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"launch_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"auto_provisioning_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_provisioning_group_type": {
				Type:         schema.TypeString,
				Default:      "maintain",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"request", "maintain"}, false),
			},
			"spot_allocation_strategy": {
				Type:         schema.TypeString,
				Default:      "lowest-price",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"lowest-price", "diversified"}, false),
			},
			"spot_instance_interruption_behavior": {
				Type:         schema.TypeString,
				Default:      "stop",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"stop", "terminate"}, false),
			},
			"spot_instance_pools_to_use_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"pay_as_you_go_allocation_strategy": {
				Type:         schema.TypeString,
				Default:      "lowest-price",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"prioritized", "lowest-price"}, false),
			},
			"excess_capacity_termination_policy": {
				Type:         schema.TypeString,
				Default:      "no-termination",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"no-termination", "termination"}, false),
			},
			"valid_from": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"valid_until": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"terminate_instances_with_expiration": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"terminate_instances": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},
			"max_spot_price": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"pay_as_you_go_target_capacity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spot_target_capacity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_target_capacity_type": {
				Type:         schema.TypeString,
				Default:      "Spot",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Spot", "PayAsYouGo"}, false),
			},
			"launch_template_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"launch_template_config": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_price": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weighted_capacity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"launch_configuration": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"image_family": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"host_names": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"credit_specification": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"deployment_set_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_release_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"io_optimized": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"security_enhancement_strategy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"internet_max_bandwidth_in": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"password_inherit": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"key_pair_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ram_role_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"system_disk_category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"system_disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"system_disk_performance_level": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"system_disk_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"system_disk_description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"system_disk": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"encrypted": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"kms_key_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"encrypt_algorithm": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"provisioned_iops": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"bursting_enabled": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"data_disk": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"disk_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"size": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"device": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"delete_with_instance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"encrypted": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"kms_key_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"encrypt_algorithm": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"performance_level": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"provisioned_iops": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"bursting_enabled": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"network_interface": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"security_group_ids": {
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"instance_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"tag": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"arn": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rolearn": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"role_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"assume_role_for": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"period": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"period_unit": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_renew": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"auto_renew_period": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"additional_info": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pvd_config": {
										Type:     schema.TypeString,
										Optional: true,
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
func resourceAlicloudAutoProvisioningGroupCreate(d *schema.ResourceData, meta interface{}) error {
	request, err := buildAlicloudAutoProvisioningGroupArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)
	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateAutoProvisioningGroup(request)
		})
		if NeedRetry(err) {
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.CreateAutoProvisioningGroupResponse)
		d.SetId(response.AutoProvisioningGroupId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_auto_provisioning_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return resourceAlicloudAutoProvisioningGroupRead(d, meta)
}
func resourceAlicloudAutoProvisioningGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeAutoProvisioningGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_auto_provisioning_group ecsService.DescribeAutoProvisioningGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("auto_provisioning_group_name", object.AutoProvisioningGroupName)
	d.Set("launch_template_id", object.LaunchTemplateId)
	d.Set("auto_provisioning_group_type", object.AutoProvisioningGroupType)
	d.Set("auto_provisioning_group_id", object.AutoProvisioningGroupId)
	d.Set("create_time", object.CreationTime)
	d.Set("excess_capacity_termination_policy", object.ExcessCapacityTerminationPolicy)
	d.Set("terminate_instances_with_expiration", object.TerminateInstancesWithExpiration)
	d.Set("max_spot_price", object.MaxSpotPrice)
	d.Set("terminate_instances", object.TerminateInstances)
	d.Set("valid_from", object.ValidFrom)
	d.Set("valid_until", object.ValidUntil)
	d.Set("status", object.Status)
	d.Set("state", object.State)
	d.Set("pay_as_you_go_target_capacity", fmt.Sprintf("%v", object.TargetCapacitySpecification.PayAsYouGoTargetCapacity))
	d.Set("pay_as_you_go_allocation_strategy", object.PayAsYouGoOptions.AllocationStrategy)
	d.Set("spot_target_capacity", fmt.Sprintf("%v", object.TargetCapacitySpecification.SpotTargetCapacity))
	d.Set("spot_allocation_strategy", object.SpotOptions.AllocationStrategy)
	d.Set("spot_instance_interruption_behavior", object.SpotOptions.InstanceInterruptionBehavior)
	d.Set("spot_instance_pools_to_use_count", object.SpotOptions.InstancePoolsToUseCount)
	d.Set("total_target_capacity", fmt.Sprintf("%v", object.TargetCapacitySpecification.TotalTargetCapacity))
	d.Set("default_target_capacity_type", object.TargetCapacitySpecification.DefaultTargetCapacityType)
	d.Set("launch_template_version", object.LaunchTemplateVersion)
	launch_template_config := []map[string]interface{}{}
	para := map[string]interface{}{}
	for _, mappara := range object.LaunchTemplateConfigs.LaunchTemplateConfig {
		para["instance_type"] = mappara.InstanceType
		para["vswitch_id"] = mappara.VSwitchId
		para["weighted_capacity"] = fmt.Sprintf("%v", mappara.WeightedCapacity)
		para["max_price"] = fmt.Sprintf("%v", mappara.MaxPrice)
		para["priority"] = fmt.Sprintf("%v", mappara.Priority)
	}
	launch_template_config = append(launch_template_config, para)
	d.Set("launch_template_config", launch_template_config)
	return nil
}
func resourceAlicloudAutoProvisioningGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateModifyAutoProvisioningGroupRequest()
	request.RegionId = client.RegionId
	request.AutoProvisioningGroupId = d.Id()
	if d.HasChange("excess_capacity_termination_policy") {
		if v, ok := d.GetOk("excess_capacity_termination_policy"); ok {
			request.ExcessCapacityTerminationPolicy = v.(string)
		}
	}
	if d.HasChange("default_target_capacity_type") {
		if v, ok := d.GetOk("default_target_capacity_type"); ok {
			request.DefaultTargetCapacityType = v.(string)
		}
	}
	if d.HasChange("terminate_instances_with_expiration") {
		request.TerminateInstancesWithExpiration = requests.NewBoolean(d.Get("terminate_instances_with_expiration").(bool))
	}
	if d.HasChange("max_spot_price") {
		request.MaxSpotPrice = requests.NewFloat(d.Get("max_spot_price").(float64))
	}
	if d.HasChange("total_target_capacity") {
		if v, ok := d.GetOk("total_target_capacity"); ok {
			request.TotalTargetCapacity = v.(string)
		}
	}
	if d.HasChange("pay_as_you_go_target_capacity") {
		if v, ok := d.GetOk("pay_as_you_go_target_capacity"); ok {
			request.PayAsYouGoTargetCapacity = v.(string)
		}
	}
	if d.HasChange("spot_target_capacity") {
		if v, ok := d.GetOk("spot_target_capacity"); ok {
			request.SpotTargetCapacity = v.(string)
		}
	}
	if d.HasChange("auto_provisioning_group_name") {
		if v, ok := d.GetOk("auto_provisioning_group_name"); ok {
			request.AutoProvisioningGroupName = v.(string)
		}
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyAutoProvisioningGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAlicloudAutoProvisioningGroupRead(d, meta)
}
func resourceAlicloudAutoProvisioningGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := ecs.CreateDeleteAutoProvisioningGroupRequest()
	request.RegionId = client.RegionId
	request.AutoProvisioningGroupId = d.Id()
	request.TerminateInstances = requests.NewBoolean(d.Get("terminate_instances").(bool))
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteAutoProvisioningGroup(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidAutoProvisioningGroupId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(ecsService.WaitForAutoProvisioningGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}
func buildAlicloudAutoProvisioningGroupArgs(d *schema.ResourceData, meta interface{}) (*ecs.CreateAutoProvisioningGroupRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateCreateAutoProvisioningGroupRequest()
	request.RegionId = client.RegionId
	request.LaunchTemplateId = d.Get("launch_template_id").(string)
	request.TotalTargetCapacity = d.Get("total_target_capacity").(string)
	// At least one of launch_template_id or launch_configuration must be specified
	if d.Get("launch_template_id").(string) == "" {
		if _, ok := d.GetOk("launch_configuration"); !ok {
			return nil, WrapError(fmt.Errorf("either launch_template_id or launch_configuration must be specified"))
		}
	}
	if v, ok := d.GetOk("auto_provisioning_group_name"); ok && v.(string) != "" {
		request.AutoProvisioningGroupName = v.(string)
	}
	if v, ok := d.GetOk("auto_provisioning_group_type"); ok && v.(string) != "" {
		request.AutoProvisioningGroupType = v.(string)
	}
	if v, ok := d.GetOk("spot_allocation_strategy"); ok && v.(string) != "" {
		request.SpotAllocationStrategy = v.(string)
	}
	if v, ok := d.GetOk("spot_instance_interruption_behavior"); ok && v.(string) != "" {
		request.SpotInstanceInterruptionBehavior = v.(string)
	}
	if v, ok := d.GetOk("spot_instancePools_to_use_count"); ok && v.(string) != "" {
		request.SpotInstancePoolsToUseCount = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("pay_as_you_go_allocation_strategy"); ok && v.(string) != "" {
		request.PayAsYouGoAllocationStrategy = v.(string)
	}
	if v, ok := d.GetOk("excess_capacity_termination_policy"); ok && v.(string) != "" {
		request.ExcessCapacityTerminationPolicy = v.(string)
	}
	if v, ok := d.GetOk("valid_from"); ok && v.(string) != "" {
		request.ValidFrom = v.(string)
	}
	if v, ok := d.GetOk("valid_until"); ok && v.(string) != "" {
		request.ValidUntil = v.(string)
	}
	if v, ok := d.GetOk("terminate_instances_with_expiration"); ok {
		request.TerminateInstancesWithExpiration = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("terminate_instances"); ok {
		request.TerminateInstances = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("max_spot_price"); ok && v.(float64) != 0.0 {
		request.MaxSpotPrice = requests.NewFloat(v.(float64))
	}
	if v, ok := d.GetOk("pay_as_you_go_target_capacity"); ok && v.(string) != "" {
		request.PayAsYouGoTargetCapacity = v.(string)
	}
	if v, ok := d.GetOk("spot_target_capacity"); ok && v.(string) != "" {
		request.SpotTargetCapacity = v.(string)
	}
	if v, ok := d.GetOk("default_target_capacity_type"); ok && v.(string) != "" {
		request.DefaultTargetCapacityType = v.(string)
	}
	if v, ok := d.GetOk("launch_template_version"); ok && v.(string) != "" {
		request.LaunchTemplateVersion = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	configs := d.Get("launch_template_config")
	confs := configs.([]interface{})
	createConfigs := make([]ecs.CreateAutoProvisioningGroupLaunchTemplateConfig, 0, len(confs))
	for _, c := range confs {
		cc := c.(map[string]interface{})
		conf := ecs.CreateAutoProvisioningGroupLaunchTemplateConfig{
			MaxPrice:         cc["max_price"].(string),
			VSwitchId:        cc["vswitch_id"].(string),
			WeightedCapacity: cc["weighted_capacity"].(string),
		}
		if v, ok := cc["instance_type"]; ok && v.(string) != "" {
			conf.InstanceType = cc["instance_type"].(string)
		}
		if v, ok := cc["priority"]; ok && v.(string) != "" {
			conf.Priority = cc["priority"].(string)
		}
		createConfigs = append(createConfigs, conf)
	}
	request.LaunchTemplateConfig = &createConfigs
	// Set LaunchConfiguration.* fields when launch_configuration block is provided.
	// These fields are create-only: DescribeAutoProvisioningGroups does not return
	// them and ModifyAutoProvisioningGroup does not accept them, so the entire
	// block is ForceNew and Read does not populate these fields.
	if lcRaw, ok := d.GetOk("launch_configuration"); ok {
		lcList := lcRaw.([]interface{})
		if len(lcList) > 0 {
			lc := lcList[0].(map[string]interface{})
			if v, ok := lc["image_id"].(string); ok && v != "" {
				request.LaunchConfigurationImageId = v
			}
			if v, ok := lc["image_family"].(string); ok && v != "" {
				request.LaunchConfigurationImageFamily = v
			}
			if v, ok := lc["instance_name"].(string); ok && v != "" {
				request.LaunchConfigurationInstanceName = v
			}
			if v, ok := lc["instance_description"].(string); ok && v != "" {
				request.LaunchConfigurationInstanceDescription = v
			}
			if v, ok := lc["host_name"].(string); ok && v != "" {
				request.LaunchConfigurationHostName = v
			}
			if v, ok := lc["host_names"].([]interface{}); ok && len(v) > 0 {
				hostNames := make([]string, 0, len(v))
				for _, hn := range v {
					hostNames = append(hostNames, hn.(string))
				}
				request.LaunchConfigurationHostNames = &hostNames
			}
			if v, ok := lc["credit_specification"].(string); ok && v != "" {
				request.LaunchConfigurationCreditSpecification = v
			}
			if v, ok := lc["deployment_set_id"].(string); ok && v != "" {
				request.LaunchConfigurationDeploymentSetId = v
			}
			if v, ok := lc["auto_release_time"].(string); ok && v != "" {
				request.LaunchConfigurationAutoReleaseTime = v
			}
			if v, ok := lc["io_optimized"].(string); ok && v != "" {
				request.LaunchConfigurationIoOptimized = v
			}
			if v, ok := lc["security_group_id"].(string); ok && v != "" {
				request.LaunchConfigurationSecurityGroupId = v
			}
			if v, ok := lc["security_group_ids"].([]interface{}); ok && len(v) > 0 {
				sgIds := make([]string, 0, len(v))
				for _, sg := range v {
					sgIds = append(sgIds, sg.(string))
				}
				request.LaunchConfigurationSecurityGroupIds = &sgIds
			}
			if v, ok := lc["security_enhancement_strategy"].(string); ok && v != "" {
				request.LaunchConfigurationSecurityEnhancementStrategy = v
			}
			if v, ok := lc["internet_max_bandwidth_in"].(int); ok && v != 0 {
				request.LaunchConfigurationInternetMaxBandwidthIn = requests.NewInteger(v)
			}
			if v, ok := lc["internet_max_bandwidth_out"].(int); ok && v != 0 {
				request.LaunchConfigurationInternetMaxBandwidthOut = requests.NewInteger(v)
			}
			if v, ok := lc["internet_charge_type"].(string); ok && v != "" {
				request.LaunchConfigurationInternetChargeType = v
			}
			if v, ok := lc["password"].(string); ok && v != "" {
				request.LaunchConfigurationPassword = v
			}
			if v, ok := lc["password_inherit"].(bool); ok {
				request.LaunchConfigurationPasswordInherit = requests.NewBoolean(v)
			}
			if v, ok := lc["key_pair_name"].(string); ok && v != "" {
				request.LaunchConfigurationKeyPairName = v
			}
			if v, ok := lc["ram_role_name"].(string); ok && v != "" {
				request.LaunchConfigurationRamRoleName = v
			}
			if v, ok := lc["user_data"].(string); ok && v != "" {
				request.LaunchConfigurationUserData = v
			}
			if v, ok := lc["resource_group_id"].(string); ok && v != "" {
				request.LaunchConfigurationResourceGroupId = v
			}
			if v, ok := lc["system_disk_category"].(string); ok && v != "" {
				request.LaunchConfigurationSystemDiskCategory = v
			}
			if v, ok := lc["system_disk_size"].(int); ok && v != 0 {
				request.LaunchConfigurationSystemDiskSize = requests.NewInteger(v)
			}
			if v, ok := lc["system_disk_performance_level"].(string); ok && v != "" {
				request.LaunchConfigurationSystemDiskPerformanceLevel = v
			}
			if v, ok := lc["system_disk_name"].(string); ok && v != "" {
				request.LaunchConfigurationSystemDiskName = v
			}
			if v, ok := lc["system_disk_description"].(string); ok && v != "" {
				request.LaunchConfigurationSystemDiskDescription = v
			}
			if sdRaw, ok := lc["system_disk"].([]interface{}); ok && len(sdRaw) > 0 {
				sd := sdRaw[0].(map[string]interface{})
				systemDisk := ecs.CreateAutoProvisioningGroupLaunchConfigurationSystemDisk{}
				if v, ok := sd["encrypted"].(string); ok && v != "" {
					systemDisk.Encrypted = v
				}
				if v, ok := sd["kms_key_id"].(string); ok && v != "" {
					systemDisk.KMSKeyId = v
				}
				if v, ok := sd["encrypt_algorithm"].(string); ok && v != "" {
					systemDisk.EncryptAlgorithm = v
				}
				if v, ok := sd["provisioned_iops"].(string); ok && v != "" {
					systemDisk.ProvisionedIops = v
				}
				if v, ok := sd["bursting_enabled"].(string); ok && v != "" {
					systemDisk.BurstingEnabled = v
				}
				request.LaunchConfigurationSystemDisk = systemDisk
			}
			if ddRaw, ok := lc["data_disk"].([]interface{}); ok && len(ddRaw) > 0 {
				dataDisks := make([]ecs.CreateAutoProvisioningGroupLaunchConfigurationDataDisk, 0, len(ddRaw))
				for _, ddItem := range ddRaw {
					dd := ddItem.(map[string]interface{})
					dataDisk := ecs.CreateAutoProvisioningGroupLaunchConfigurationDataDisk{}
					if v, ok := dd["category"].(string); ok && v != "" {
						dataDisk.Category = v
					}
					if v, ok := dd["disk_name"].(string); ok && v != "" {
						dataDisk.DiskName = v
					}
					if v, ok := dd["size"].(string); ok && v != "" {
						dataDisk.Size = v
					}
					if v, ok := dd["device"].(string); ok && v != "" {
						dataDisk.Device = v
					}
					if v, ok := dd["snapshot_id"].(string); ok && v != "" {
						dataDisk.SnapshotId = v
					}
					if v, ok := dd["description"].(string); ok && v != "" {
						dataDisk.Description = v
					}
					if v, ok := dd["delete_with_instance"].(string); ok && v != "" {
						dataDisk.DeleteWithInstance = v
					}
					if v, ok := dd["encrypted"].(string); ok && v != "" {
						dataDisk.Encrypted = v
					}
					if v, ok := dd["kms_key_id"].(string); ok && v != "" {
						dataDisk.KmsKeyId = v
					}
					if v, ok := dd["encrypt_algorithm"].(string); ok && v != "" {
						dataDisk.EncryptAlgorithm = v
					}
					if v, ok := dd["performance_level"].(string); ok && v != "" {
						dataDisk.PerformanceLevel = v
					}
					if v, ok := dd["provisioned_iops"].(string); ok && v != "" {
						dataDisk.ProvisionedIops = v
					}
					if v, ok := dd["bursting_enabled"].(string); ok && v != "" {
						dataDisk.BurstingEnabled = v
					}
					dataDisks = append(dataDisks, dataDisk)
				}
				request.LaunchConfigurationDataDisk = &dataDisks
			}
			if niRaw, ok := lc["network_interface"].([]interface{}); ok && len(niRaw) > 0 {
				nis := make([]ecs.CreateAutoProvisioningGroupLaunchConfigurationNetworkInterface, 0, len(niRaw))
				for _, niItem := range niRaw {
					ni := niItem.(map[string]interface{})
					networkInterface := ecs.CreateAutoProvisioningGroupLaunchConfigurationNetworkInterface{}
					if v, ok := ni["security_group_id"].(string); ok && v != "" {
						networkInterface.SecurityGroupId = v
					}
					if sgRaw, ok := ni["security_group_ids"].([]interface{}); ok && len(sgRaw) > 0 {
						sgIds := make([]string, 0, len(sgRaw))
						for _, sg := range sgRaw {
							sgIds = append(sgIds, sg.(string))
						}
						networkInterface.SecurityGroupIds = &sgIds
					}
					if v, ok := ni["instance_type"].(string); ok && v != "" {
						networkInterface.InstanceType = v
					}
					nis = append(nis, networkInterface)
				}
				request.LaunchConfigurationNetworkInterface = &nis
			}
			if tagRaw, ok := lc["tag"].([]interface{}); ok && len(tagRaw) > 0 {
				tags := make([]ecs.CreateAutoProvisioningGroupLaunchConfigurationTag, 0, len(tagRaw))
				for _, tagItem := range tagRaw {
					tag := tagItem.(map[string]interface{})
					lcTag := ecs.CreateAutoProvisioningGroupLaunchConfigurationTag{}
					if v, ok := tag["key"].(string); ok {
						lcTag.Key = v
					}
					if v, ok := tag["value"].(string); ok {
						lcTag.Value = v
					}
					tags = append(tags, lcTag)
				}
				request.LaunchConfigurationTag = &tags
			}
			if arnRaw, ok := lc["arn"].([]interface{}); ok && len(arnRaw) > 0 {
				arns := make([]ecs.CreateAutoProvisioningGroupLaunchConfigurationArn, 0, len(arnRaw))
				for _, arnItem := range arnRaw {
					arnMap := arnItem.(map[string]interface{})
					lcArn := ecs.CreateAutoProvisioningGroupLaunchConfigurationArn{}
					if v, ok := arnMap["rolearn"].(string); ok && v != "" {
						lcArn.Rolearn = v
					}
					if v, ok := arnMap["role_type"].(string); ok && v != "" {
						lcArn.RoleType = v
					}
					if v, ok := arnMap["assume_role_for"].(string); ok && v != "" {
						lcArn.AssumeRoleFor = v
					}
					arns = append(arns, lcArn)
				}
				request.LaunchConfigurationArn = &arns
			}
			lcStruct := ecs.CreateAutoProvisioningGroupLaunchConfiguration{}
			needLcStruct := false
			if v, ok := lc["period"].(string); ok && v != "" {
				lcStruct.Period = v
				needLcStruct = true
			}
			if v, ok := lc["period_unit"].(string); ok && v != "" {
				lcStruct.PeriodUnit = v
				needLcStruct = true
			}
			if v, ok := lc["auto_renew"].(string); ok && v != "" {
				lcStruct.AutoRenew = v
				needLcStruct = true
			}
			if v, ok := lc["auto_renew_period"].(string); ok && v != "" {
				lcStruct.AutoRenewPeriod = v
				needLcStruct = true
			}
			if needLcStruct {
				request.LaunchConfiguration = lcStruct
			}
			if aiRaw, ok := lc["additional_info"].([]interface{}); ok && len(aiRaw) > 0 {
				ai := aiRaw[0].(map[string]interface{})
				additionalInfo := ecs.CreateAutoProvisioningGroupLaunchConfigurationAdditionalInfo{}
				if v, ok := ai["pvd_config"].(string); ok && v != "" {
					additionalInfo.PvdConfig = v
				}
				request.LaunchConfigurationAdditionalInfo = additionalInfo
			}
		}
	}
	return request, nil
}
