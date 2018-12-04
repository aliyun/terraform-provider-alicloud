package alicloud

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEssScalingConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingConfigurationCreate,
		Read:   resourceAliyunEssScalingConfigurationRead,
		Update: resourceAliyunEssScalingConfigurationUpdate,
		Delete: resourceAliyunEssScalingConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scaling_group_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validateInstanceType,
			},
			"io_optimized": &schema.Schema{
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Attribute io_optimized has been deprecated on instance resource. All the launched alicloud instances will be IO optimized. Suggest to remove it from your template.",
			},
			"is_outdated": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"scaling_configuration_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internet_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      PayByBandwidth,
				ValidateFunc: validateInternetChargeType,
			},
			"internet_max_bandwidth_in": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"internet_max_bandwidth_out": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInternetMaxBandWidthOut,
			},
			"system_disk_category": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      DiskCloudEfficiency,
				ValidateFunc: validateDiskCategory,
			},
			"system_disk_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"data_disk": &schema.Schema{
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"category": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateDiskCategory,
						},
						"snapshot_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"device": &schema.Schema{
							Type:       schema.TypeString,
							Optional:   true,
							Deprecated: "Attribute device has been deprecated on disk attachment resource. Suggest to remove it from your template.",
						},
						"delete_with_instance": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"instance_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instance_ids' has been deprecated from provider version 1.6.0. New resource 'alicloud_ess_attachment' replaces it.",
			},

			"substitute": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"key_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"force_delete": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"instance_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ESS-Instance",
				ValidateFunc: validateInstanceName,
			},
		},
	}
}

func resourceAliyunEssScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {

	// Ensure instance_type is generation three
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}
	if err := ecsService.InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
		return err
	}

	args, err := buildAlicloudEssScalingConfigurationArgs(d, meta)
	if err != nil {
		return err
	}

	args.IoOptimized = string(IOOptimized)
	if d.Get("is_outdated").(bool) == true {
		args.IoOptimized = string(NoneOptimized)
	}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateScalingConfiguration(args)
		})
		if err != nil {
			if IsExceptedError(err, EssThrottling) || IsExceptedError(err, IncorrectScalingGroupStatus) {
				return resource.RetryableError(fmt.Errorf("Error Create Scaling Configuration: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error Create Scaling Configuration: %#v.", err))
		}
		scaling, _ := raw.(*ess.CreateScalingConfigurationResponse)
		d.SetId(scaling.ScalingConfigurationId)
		return nil
	}); err != nil {
		return err
	}

	return resourceAliyunEssScalingConfigurationUpdate(d, meta)
}

func resourceAliyunEssScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	d.Partial(true)
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	if d.HasChange("active") {
		c, err := essService.DescribeScalingConfigurationById(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error Describe ESS scaling configuration Attribute: %#v", err)
		}

		active := d.Get("active").(bool)

		if active {
			if c.LifecycleState == string(Inactive) {

				err := essService.ActiveScalingConfigurationById(c.ScalingGroupId, d.Id())
				if err != nil {
					return fmt.Errorf("Active scaling configuration %s err: %#v", d.Id(), err)
				}
			}
		} else {
			if c.LifecycleState == string(Active) {
				_, err := activeSubstituteScalingConfiguration(d, meta)
				if err != nil {
					return err
				}
			}
		}
		d.SetPartial("active")
	}

	if err := enableEssScalingConfiguration(d, meta); err != nil {
		return err
	}

	d.Partial(false)

	return resourceAliyunEssScalingConfigurationRead(d, meta)
}

func enableEssScalingConfiguration(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if d.HasChange("enable") {
		sgId := d.Get("scaling_group_id").(string)
		group, err := essService.DescribeScalingGroupById(sgId)
		if err != nil {
			return fmt.Errorf("DescribeScalingGroupById %s error: %#v", sgId, err)
		}
		enable := d.Get("enable").(bool)

		if enable {
			if group.LifecycleState == string(Inactive) {

				cs, err := essService.DescribeScalingConfifurations(sgId)

				if err != nil {
					return fmt.Errorf("Describe ScalingConfigurations by scaling group %s got an error: %#v", sgId, err)
				}
				activeConfig := ""
				var csIds []string
				for _, c := range cs {
					csIds = append(csIds, c.ScalingConfigurationId)
					if c.LifecycleState == string(Active) {
						activeConfig = c.ScalingConfigurationId
					}
				}

				if activeConfig == "" {
					return fmt.Errorf("Please active a scaling configuration before enabling scaling group %s. "+
						"Its all scaling configuration are %s.", sgId, strings.Join(csIds, ","))
				}

				req := ess.CreateEnableScalingGroupRequest()
				req.ScalingGroupId = sgId
				req.ActiveScalingConfigurationId = activeConfig

				_, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.EnableScalingGroup(req)
				})
				if err != nil {
					return fmt.Errorf("EnableScalingGroup %s got an error: %#v", sgId, err)
				}
				if err := essService.WaitForScalingGroup(sgId, Active, DefaultTimeout); err != nil {
					return fmt.Errorf("WaitForScalingGroup is %#v got an error: %#v.", Active, err)
				}

				d.SetPartial("scaling_configuration_id")
			}
		} else {
			if group.LifecycleState == string(Active) {
				req := ess.CreateDisableScalingGroupRequest()
				req.ScalingGroupId = sgId
				_, err = client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.DisableScalingGroup(req)
				})
				if err != nil {
					return fmt.Errorf("DisableScalingGroup %s got an error: %#v", sgId, err)
				}
				if err := essService.WaitForScalingGroup(sgId, Inactive, DefaultTimeout); err != nil {
					return fmt.Errorf("WaitForScalingGroup is %#v got an error: %#v.", Inactive, err)
				}
			}
		}
		d.SetPartial("enable")
	}

	return nil
}

func resourceAliyunEssScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}
	c, err := essService.DescribeScalingConfigurationById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS scaling configuration Attribute: %#v", err)
	}

	d.Set("scaling_group_id", c.ScalingGroupId)
	d.Set("active", c.LifecycleState == string(Active))
	d.Set("image_id", c.ImageId)
	d.Set("instance_type", c.InstanceType)
	d.Set("security_group_id", c.SecurityGroupId)
	d.Set("scaling_configuration_name", c.ScalingConfigurationName)
	d.Set("internet_charge_type", c.InternetChargeType)
	d.Set("internet_max_bandwidth_in", c.InternetMaxBandwidthIn)
	d.Set("internet_max_bandwidth_out", c.InternetMaxBandwidthOut)
	d.Set("system_disk_category", c.SystemDiskCategory)
	d.Set("system_disk_size", c.SystemDiskSize)
	d.Set("data_disk", essService.flattenDataDiskMappings(c.DataDisks.DataDisk))
	d.Set("role_name", c.RamRoleName)
	d.Set("key_name", c.KeyPairName)
	d.Set("user_data", userDataHashSum(c.UserData))
	d.Set("force_delete", d.Get("force_delete").(bool))
	d.Set("tags", essTagsToMap(c.Tags.Tag))
	d.Set("instance_name", c.InstanceName)

	return nil
}

func resourceAliyunEssScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	c, err := essService.DescribeScalingConfigurationById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return fmt.Errorf("Describing scaling configuration got an error: %#v.", err)
	}

	req := ess.CreateDescribeScalingConfigurationsRequest()
	req.ScalingGroupId = c.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(req)
	})
	resp, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if resp == nil || len(resp.ScalingConfigurations.ScalingConfiguration) < 1 {
		return nil
	}
	if len(resp.ScalingConfigurations.ScalingConfiguration) == 1 {
		if d.Get("force_delete").(bool) {
			return essService.DeleteScalingGroupById(c.ScalingGroupId)
		}
		return fmt.Errorf("Current scaling configuration %s is the last configuration for the scaling group %s. Please launch a new "+
			"active scaling configuration or set 'force_delete' to 'true' to delete it with deleting its scaling group.", d.Id(), c.ScalingGroupId)
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		req := ess.CreateDeleteScalingConfigurationRequest()
		req.ScalingConfigurationId = d.Id()

		_, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DeleteScalingConfiguration(req)
		})

		if err != nil {
			if IsExceptedErrors(err, []string{IncorrectScalingConfigurationLifecycleState}) {
				return resource.NonRetryableError(
					fmt.Errorf("Scaling configuration is active. Please active another one before deleting it and trying again."))
			}
			if IsExceptedErrors(err, []string{InvalidScalingGroupIdNotFound}) {
				return resource.RetryableError(
					fmt.Errorf("Delete scaling configuration timeout and got an error:%#v.", err))
			}
		}

		c, err := essService.DescribeScalingConfigurationById(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		instances, err := essService.DescribeScalingInstances(c.ScalingGroupId, d.Id(), make([]string, 0), "")
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if len(instances) > 0 {
			return resource.NonRetryableError(fmt.Errorf("There are still ECS instances in the scaling configuration - please remove them and try again."))
		}

		return resource.RetryableError(
			fmt.Errorf("Delete scaling configuration timeout and got an error:%#v.", err))
	})
}

func buildAlicloudEssScalingConfigurationArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingConfigurationRequest, error) {
	args := ess.CreateCreateScalingConfigurationRequest()
	args.ScalingGroupId = d.Get("scaling_group_id").(string)
	args.ImageId = d.Get("image_id").(string)
	args.InstanceType = d.Get("instance_type").(string)
	args.SecurityGroupId = d.Get("security_group_id").(string)

	if v := d.Get("scaling_configuration_name").(string); v != "" {
		args.ScalingConfigurationName = v
	}

	if v := d.Get("internet_charge_type").(string); v != "" {
		args.InternetChargeType = v
	}

	if v := d.Get("internet_max_bandwidth_in").(int); v != 0 {
		args.InternetMaxBandwidthIn = requests.NewInteger(v)
	}

	args.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))

	if v := d.Get("system_disk_category").(string); v != "" {
		args.SystemDiskCategory = v
	}

	if v := d.Get("system_disk_size").(int); v != 0 {
		args.SystemDiskSize = requests.NewInteger(v)
	}

	dds, ok := d.GetOk("data_disk")
	if ok {
		disks := dds.([]interface{})
		createDataDisks := make([]ess.CreateScalingConfigurationDataDisk, 0, len(disks))
		for _, e := range disks {
			pack := e.(map[string]interface{})
			dataDisk := ess.CreateScalingConfigurationDataDisk{
				Size:               strconv.Itoa(pack["size"].(int)),
				Category:           pack["category"].(string),
				SnapshotId:         pack["snapshot_id"].(string),
				DeleteWithInstance: strconv.FormatBool(pack["delete_with_instance"].(bool)),
			}
			createDataDisks = append(createDataDisks, dataDisk)
		}
		args.DataDisk = &createDataDisks
	}

	if v, ok := d.GetOk("role_name"); ok && v.(string) != "" {
		args.RamRoleName = v.(string)
	}

	if v, ok := d.GetOk("key_name"); ok && v.(string) != "" {
		args.KeyPairName = v.(string)
	}

	if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
		_, base64DecodeError := base64.StdEncoding.DecodeString(v.(string))
		if base64DecodeError == nil {
			args.UserData = v.(string)
		} else {
			args.UserData = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := "{"
		for key, value := range v.(map[string]interface{}) {
			tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
		}
		args.Tags = strings.TrimSuffix(tags, ",") + "}"
	}

	if v, ok := d.GetOk("instance_name"); ok && v.(string) != "" {
		args.InstanceName = v.(string)
	}

	return args, nil
}

func activeSubstituteScalingConfiguration(d *schema.ResourceData, meta interface{}) (configures []ess.ScalingConfiguration, err error) {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	substitute_id, ok := d.GetOk("substitute")

	c, err := essService.DescribeScalingConfigurationById(d.Id())
	if err != nil {
		return
	}

	req := ess.CreateDescribeScalingConfigurationsRequest()
	req.ScalingGroupId = c.ScalingGroupId

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DescribeScalingConfigurations(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*ess.DescribeScalingConfigurationsResponse)
	if resp == nil || len(resp.ScalingConfigurations.ScalingConfiguration) < 1 {
		return
	}

	if !ok || substitute_id.(string) == "" {

		if len(resp.ScalingConfigurations.ScalingConfiguration) == 1 {
			return configures, fmt.Errorf("Current scaling configuration %s is the last configuration for the scaling group %s, and it can't be inactive.", d.Id(), c.ScalingGroupId)
		}

		var configs []string
		for _, cc := range resp.ScalingConfigurations.ScalingConfiguration {
			if cc.ScalingConfigurationId != d.Id() {
				configs = append(configs, cc.ScalingConfigurationId)
			}
		}

		return configures, fmt.Errorf("Before inactivating current scaling configuration, you must select a substitute for scaling group from: %s.", strings.Join(configs, ","))

	}

	err = essService.ActiveScalingConfigurationById(c.ScalingGroupId, substitute_id.(string))
	if err != nil {
		return configures, fmt.Errorf("Inactive scaling configuration %s err: %#v. Substitute scaling configuration ID: %s",
			d.Id(), err, substitute_id.(string))
	}

	return resp.ScalingConfigurations.ScalingConfiguration, nil
}
