package alicloud

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func resourceAlicloudEssInstanceRefresh() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssInstanceRefreshCreate,
		Read:   resourceAliyunEssInstanceRefreshRead,
		Update: resourceAliyunEssInstanceRefreshUpdate,
		Delete: resourceAliyunEssInstanceRefreshDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"min_healthy_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(0, 100),
				Computed:     true,
			},
			"max_healthy_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: IntBetween(100, 200),
				Computed:     true,
			},
			"desired_configuration_image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"desired_configuration_launch_template_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"desired_configuration_launch_template_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"desired_configuration_launch_template_overrides": {
				Optional: true,
				Type:     schema.TypeList,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"desired_configuration_containers": {
				Optional: true,
				Type:     schema.TypeList,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"image": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"environment_vars": {
							Type:     schema.TypeList,
							ForceNew: true,
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
										Computed: true,
									},
									"field_ref_field_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"commands": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"args": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"skip_matching": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"checkpoints": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"percentage": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"checkpoint_pause_time": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunEssInstanceRefreshCreate(d *schema.ResourceData, meta interface{}) error {
	request, errRequest := buildAlicloudEssInstanceRefreshArgs(d, meta)
	if errRequest != nil {
		return WrapError(errRequest)
	}
	client := meta.(*connectivity.AliyunClient)
	wait := incrementalWait(1*time.Second, 2*time.Second)

	var raw map[string]interface{}
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = client.RpcPost("Ess", "2014-08-28", "StartInstanceRefresh", nil, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidInstance.NoInstanceNeedBeRefreshed"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_instance_refresh", "StartInstanceRefresh", AlibabaCloudSdkGoERROR)
	}
	addDebug("StartInstanceRefresh", raw, request, request)
	d.SetId(fmt.Sprint(request["ScalingGroupId"], ":", raw["InstanceRefreshTaskId"]))
	return resourceAliyunEssInstanceRefreshRead(d, meta)
}

func resourceAliyunEssInstanceRefreshRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	var object map[string]interface{}
	var err error
	object, err = essService.DescribeEssInstanceRefresh(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("scaling_group_id", object["ScalingGroupId"])

	d.Set("status", object["Status"])
	if object["MinHealthyPercentage"] != nil {
		d.Set("min_healthy_percentage", object["MinHealthyPercentage"])
	}
	if object["MaxHealthyPercentage"] != nil {
		d.Set("max_healthy_percentage", object["MaxHealthyPercentage"])
	}
	if v := object["DesiredConfiguration"]; v != nil {
		m := v.(map[string]interface{})
		if m["ImageId"] != nil {
			d.Set("desired_configuration_image_id", m["ImageId"])

		}
		if m["LaunchTemplateId"] != nil {
			d.Set("desired_configuration_launch_template_id", m["LaunchTemplateId"])
		}
		if m["LaunchTemplateVersion"] != nil {
			d.Set("desired_configuration_launch_template_version", m["LaunchTemplateVersion"])
		}

		if v := m["LaunchTemplateOverrides"]; v != nil {
			result := make([]map[string]interface{}, 0)
			for _, i := range v.(map[string]interface{})["LaunchTemplateOverride"].([]interface{}) {
				launchTemplateOverride := i.(map[string]interface{})
				l := map[string]interface{}{
					"instance_type": launchTemplateOverride["InstanceType"],
				}
				result = append(result, l)
			}
			err := d.Set("desired_configuration_launch_template_overrides", result)
			if err != nil {
				return WrapError(err)
			}
		}
		if v := m["Containers"]; v != nil {
			result := make([]map[string]interface{}, 0)
			for _, i := range v.(map[string]interface{})["Container"].([]interface{}) {
				container := i.(map[string]interface{})
				l := map[string]interface{}{}
				if container["Name"] != nil {
					l["name"] = container["Name"]
				}
				if container["Image"] != nil {
					l["image"] = container["Image"]
				}
				if container["Commands"] != nil {
					var commands []string
					commandItem := container["Commands"].(map[string]interface{})
					for _, command := range commandItem["Command"].([]interface{}) {
						commands = append(commands, command.(string))
					}
					l["commands"] = commands
				}
				if container["Args"] != nil {
					var args []string
					argItem := container["Args"].(map[string]interface{})
					for _, arg := range argItem["Arg"].([]interface{}) {
						args = append(args, arg.(string))
					}
					l["args"] = args
				}
				if container["EnvironmentVars"] != nil {
					envVars := make([]map[string]interface{}, 0)
					for _, j := range container["EnvironmentVars"].(map[string]interface{})["EnvironmentVar"].([]interface{}) {
						envVar := j.(map[string]interface{})
						w := map[string]interface{}{}
						if envVar["Key"] != nil {
							w["key"] = envVar["Key"]
						}
						if envVar["Value"] != nil {
							w["value"] = envVar["Value"]
						}
						if envVar["FieldRefFieldPath"] != nil {
							w["field_ref_field_path"] = envVar["FieldRefFieldPath"]
						}
						envVars = append(envVars, w)
					}
					l["environment_vars"] = envVars
				}

				result = append(result, l)
			}
			err := d.Set("desired_configuration_containers", result)
			if err != nil {
				return WrapError(err)
			}
		}
	}

	if object["SkipMatching"] != nil {
		d.Set("skip_matching", object["SkipMatching"])
	}
	if object["CheckpointPauseTime"] != nil {
		d.Set("checkpoint_pause_time", object["CheckpointPauseTime"])
	}

	if v := object["Checkpoints"]; v != nil {
		result := make([]map[string]interface{}, 0)
		for _, i := range v.(map[string]interface{})["Checkpoint"].([]interface{}) {
			checkpoint := i.(map[string]interface{})
			l := map[string]interface{}{}
			if checkpoint["Percentage"] != nil {
				l["percentage"] = checkpoint["Percentage"]
			}
			result = append(result, l)
		}
		err := d.Set("checkpoints", result)
		if err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAliyunEssInstanceRefreshDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceAliyunEssInstanceRefreshUpdate(d *schema.ResourceData, meta interface{}) error {
	strs, _ := ParseResourceId(d.Id(), 2)
	scalingGroupId, instanceRefreshTaskId := strs[0], strs[1]
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"ScalingGroupId":        scalingGroupId,
		"RegionId":              client.RegionId,
		"InstanceRefreshTaskId": instanceRefreshTaskId,
	}
	update := false
	var status string
	if d.HasChange("status") && (d.Get("status").(string) == "Cancelled" || d.Get("status").(string) == "RollbackSuccessful") {
		status = d.Get("status").(string)
		update = true
	}
	if update {
		var raw map[string]interface{}
		var err error
		var operatorErr error
		var action string
		essService := EssService{client}
		wait := incrementalWait(1*time.Second, 2*time.Second)
		switch status {
		case "Cancelled":
			{
				action = "CancelInstanceRefresh"
				raw, operatorErr = client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
				err = resource.Retry(10*time.Minute, func() *resource.RetryError {
					object, errQuery := essService.DescribeEssInstanceRefresh(d.Id())
					if errQuery != nil {
						return resource.NonRetryableError(err)
					}
					if object["Status"].(string) == "Cancelling" {
						cancellingErr := errors.New(d.Id() + " is " + object["Status"].(string) + ",Please wait it")
						wait()
						return resource.RetryableError(cancellingErr)
					}
					return nil
				})
			}
		case "RollbackSuccessful":
			{
				action = "RollbackInstanceRefresh"
				raw, operatorErr = client.RpcPost("Ess", "2014-08-28", action, nil, request, true)
				err = resource.Retry(10*time.Minute, func() *resource.RetryError {
					object, errQuery := essService.DescribeEssInstanceRefresh(d.Id())
					if errQuery != nil {
						return resource.NonRetryableError(err)
					}
					if object["Status"].(string) == "RollbackInProgress" {
						rollbackInProgressErr := errors.New(d.Id() + " is " + object["Status"].(string) + ",Please wait it")
						wait()
						return resource.RetryableError(rollbackInProgressErr)
					}
					if object["Status"].(string) == "RollbackFailed" {
						rollbackInProgressErr := errors.New(d.Id() + " is RollbackFailed")
						return resource.NonRetryableError(rollbackInProgressErr)
					}
					return nil
				})
			}
		}
		if operatorErr != nil {
			return WrapErrorf(operatorErr, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, raw, request, request)
	}

	return resourceAliyunEssInstanceRefreshRead(d, meta)
}

func buildAlicloudEssInstanceRefreshArgs(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"ScalingGroupId": d.Get("scaling_group_id"),
		"RegionId":       client.RegionId,
	}
	if v, ok := d.GetOkExists("min_healthy_percentage"); ok {
		request["MinHealthyPercentage"] = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("max_healthy_percentage"); ok {
		request["MaxHealthyPercentage"] = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("desired_configuration_image_id"); ok && v.(string) != "" {
		request["DesiredConfiguration.ImageId"] = v
	}
	if v, ok := d.GetOk("desired_configuration_launch_template_id"); ok && v.(string) != "" {
		request["DesiredConfiguration.LaunchTemplateId"] = v
	}
	if v, ok := d.GetOk("desired_configuration_scaling_configuration_id"); ok && v.(string) != "" {
		request["DesiredConfiguration.ScalingConfigurationId"] = v
	}
	if v, ok := d.GetOk("desired_configuration_launch_templateId"); ok && v.(string) != "" {
		request["DesiredConfiguration.LaunchTemplateId"] = v
	}
	if v, ok := d.GetOk("desired_configuration_launch_template_version"); ok && v.(string) != "" {
		request["DesiredConfiguration.LaunchTemplateVersion"] = v
	}
	if v, ok := d.GetOk("desired_configuration_launch_template_overrides"); ok {
		desiredConfigurationLaunchTemplateOverrides := v.([]interface{})
		for i, rew := range desiredConfigurationLaunchTemplateOverrides {
			item := rew.(map[string]interface{})
			if instanceType, ok := item["instance_type"].(string); ok && instanceType != "" {
				request[fmt.Sprintf("DesiredConfiguration.LaunchTemplateOverrides.%d.InstanceType", i+1)] = instanceType
			}
		}
	}
	if v, ok := d.GetOk("checkpoints"); ok {
		checkpoints := v.([]interface{})
		for i, rew := range checkpoints {
			item := rew.(map[string]interface{})
			if percentage, ok := item["percentage"].(int); ok {
				request[fmt.Sprintf("Checkpoints.%d.Percentage", i+1)] = percentage
			}
		}
	}

	request["SkipMatching"] = d.Get("skip_matching")
	if v, ok := d.GetOk("checkpoint_pause_time"); ok {
		request["CheckpointPauseTime"] = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("desired_configuration_containers"); ok {
		desiredConfigurationContainers := v.([]interface{})
		for i, rew := range desiredConfigurationContainers {
			item := rew.(map[string]interface{})
			i = i + 1
			if name, ok := item["name"].(string); ok && name != "" {
				request[fmt.Sprintf("DesiredConfiguration.Containers.%d.Name", i)] = name
			}
			if image, ok := item["image"].(string); ok && image != "" {
				request[fmt.Sprintf("DesiredConfiguration.Containers.%d.Image", i)] = image
			}
			if item["environment_vars"] != nil {
				environmentVars := item["environment_vars"].([]interface{})
				for j, envVar := range environmentVars {
					environmentVar := envVar.(map[string]interface{})
					j = j + 1
					if environmentVar["key"] != nil {
						request[fmt.Sprintf("DesiredConfiguration.Containers.%d.EnvironmentVars.%d.Key", i, j)] = environmentVar["key"].(string)
					}
					if environmentVar["value"] != nil {
						request[fmt.Sprintf("DesiredConfiguration.Containers.%d.EnvironmentVars.%d.Value", i, j)] = environmentVar["value"].(string)
					}
					if environmentVar["field_ref_field_path"] != nil {
						request[fmt.Sprintf("DesiredConfiguration.Containers.%d.EnvironmentVars.%d.FieldRefFieldPath", i, j)] = environmentVar["field_ref_field_path"].(string)
					}
				}
			}
			if item["commands"] != nil {
				commands := item["commands"].([]interface{})
				for commandIndex, command := range commands {
					commandIndex = commandIndex + 1
					if command != nil {
						request[fmt.Sprintf("DesiredConfiguration.Containers.%d.Commands.%d", i, commandIndex)] = command.(string)
					}
				}
			}
			if item["args"] != nil {
				args := item["args"].([]interface{})
				for argIndex, arg := range args {
					argIndex = argIndex + 1
					if arg != nil {
						request[fmt.Sprintf("DesiredConfiguration.Containers.%d.Args.%d", i, argIndex)] = arg.(string)
					}
				}
			}
		}
	}

	return request, nil
}
