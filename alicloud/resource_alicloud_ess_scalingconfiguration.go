package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				Computed:     true,
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
				Default:      ecs.DiskCategoryCloudEfficiency,
				ValidateFunc: validateDiskCategory,
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
		},
	}
}

func resourceAliyunEssScalingConfigurationCreate(d *schema.ResourceData, meta interface{}) error {

	// Ensure instance_type is generation three
	validData, err := meta.(*AliyunClient).CheckParameterValidity(d, meta)
	if err != nil {
		return err
	}

	args, err := buildAlicloudEssScalingConfigurationArgs(d, meta)
	if err != nil {
		return err
	}

	if validData[IoOptimizedKey].(ecs.IoOptimized) == ecs.IoOptimizedOptimized {
		args.IoOptimized = ecs.IoOptimizedOptimized
	}

	essconn := meta.(*AliyunClient).essconn

	scaling, err := essconn.CreateScalingConfiguration(args)
	if err != nil && !IsExceptedError(err, IncorrectScalingGroupStatus) {
		return fmt.Errorf("Error Create Scaling Configuration: %#v", err)
	}

	d.SetId(scaling.ScalingConfigurationId)

	return resourceAliyunEssScalingConfigurationUpdate(d, meta)
}

func resourceAliyunEssScalingConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	d.Partial(true)
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	if d.HasChange("active") {
		c, err := client.DescribeScalingConfigurationById(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("Error Describe ESS scaling configuration Attribute: %#v", err)
		}

		active := d.Get("active").(bool)

		if active {
			if c.LifecycleState == ess.Inacitve {

				err := client.ActiveScalingConfigurationById(c.ScalingGroupId, d.Id())
				if err != nil {
					return fmt.Errorf("Active scaling configuration %s err: %#v", d.Id(), err)
				}
			}
		} else {
			if c.LifecycleState == ess.Active {
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
	client := meta.(*AliyunClient)

	if d.HasChange("enable") {
		sgId := d.Get("scaling_group_id").(string)
		group, err := client.DescribeScalingGroupById(sgId)
		if err != nil {
			return fmt.Errorf("DescribeScalingGroupById %s error: %#v", sgId, err)
		}
		enable := d.Get("enable").(bool)

		if enable {
			if group.LifecycleState == ess.Inacitve {

				cs, _, err := client.essconn.DescribeScalingConfigurations(&ess.DescribeScalingConfigurationsArgs{
					RegionId:       getRegion(d, meta),
					ScalingGroupId: sgId,
					Pagination:     getPagination(1, 50),
				})

				if err != nil {
					return fmt.Errorf("Describe ScalingConfigurations by scaling group %s got an error: %#v", sgId, err)
				}
				activeConfig := ""
				var csIds []string
				for _, c := range cs {
					csIds = append(csIds, c.ScalingConfigurationId)
					if c.LifecycleState == ess.Active {
						activeConfig = c.ScalingConfigurationId
					}
				}

				if activeConfig == "" {
					return fmt.Errorf("Please active a scaling configuration before enabling scaling group %s. "+
						"Its all scaling configuration are %s.", sgId, strings.Join(csIds, ","))
				}

				if _, err := client.essconn.EnableScalingGroup(&ess.EnableScalingGroupArgs{
					ScalingGroupId:               sgId,
					ActiveScalingConfigurationId: activeConfig,
				}); err != nil {
					return fmt.Errorf("EnableScalingGroup %s got an error: %#v", sgId, err)
				}
				if err := client.essconn.WaitForScalingGroup(getRegion(d, meta), sgId, ess.Active, defaultTimeout); err != nil {
					return fmt.Errorf("WaitForScalingGroup is %#v got an error: %#v.", ess.Active, err)
				}

				d.SetPartial("scaling_configuration_id")
			}
		} else {
			if group.LifecycleState == ess.Active {
				if _, err := client.essconn.DisableScalingGroup(&ess.DisableScalingGroupArgs{
					ScalingGroupId: sgId,
				}); err != nil {
					return fmt.Errorf("DisableScalingGroup %s got an error: %#v", sgId, err)
				}
				if err := client.essconn.WaitForScalingGroup(getRegion(d, meta), sgId, ess.Inacitve, defaultTimeout); err != nil {
					return fmt.Errorf("WaitForScalingGroup is %#v got an error: %#v.", ess.Inacitve, err)
				}
			}
		}
		d.SetPartial("enable")
	}

	return nil
}

func resourceAliyunEssScalingConfigurationRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}
	c, err := client.DescribeScalingConfigurationById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS scaling configuration Attribute: %#v", err)
	}

	d.Set("scaling_group_id", c.ScalingGroupId)
	d.Set("active", c.LifecycleState == ess.Active)
	d.Set("image_id", c.ImageId)
	d.Set("instance_type", c.InstanceType)
	d.Set("security_group_id", c.SecurityGroupId)
	d.Set("scaling_configuration_name", c.ScalingConfigurationName)
	d.Set("internet_charge_type", c.InternetChargeType)
	d.Set("internet_max_bandwidth_in", c.InternetMaxBandwidthIn)
	d.Set("internet_max_bandwidth_out", c.InternetMaxBandwidthOut)
	d.Set("system_disk_category", c.SystemDiskCategory)
	d.Set("data_disk", flattenDataDiskMappings(c.DataDisks.DataDisk))
	d.Set("role_name", c.RamRoleName)
	d.Set("key_name", c.KeyPairName)
	d.Set("user_data", userDataHashSum(c.UserData))
	d.Set("force_delete", d.Get("force_delete").(bool))
	d.Set("tags", essTagsToMap(c.Tags.Tag))

	return nil
}

func resourceAliyunEssScalingConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	if strings.Contains(d.Id(), COLON_SEPARATED) {
		d.SetId(strings.Split(d.Id(), COLON_SEPARATED)[1])
	}

	configs, _ := activeSubstituteScalingConfiguration(d, meta)
	if len(configs) <= 1 {
		if len(configs) == 0 {
			return nil
		}
		if d.Get("force_delete").(bool) {
			return client.DeleteScalingGroupById(configs[0].ScalingGroupId)
		}
		return fmt.Errorf("Current scaling configuration %s is the last configuration for the scaling group %s. Please launch a new "+
			"active scaling configuration or set 'force_delete' to 'true' to delete it with deleting its scaling group.", d.Id(), configs[0].ScalingGroupId)
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err := client.essconn.DeleteScalingConfiguration(&ess.DeleteScalingConfigurationArgs{
			ScalingConfigurationId: d.Id(),
		})

		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == IncorrectScalingConfigurationLifecycleState {
				return resource.NonRetryableError(
					fmt.Errorf("Scaling configuration is active. Please active another one before deleting it and trying again."))
			}
			if e.ErrorResponse.Code != InvalidScalingGroupIdNotFound {
				return resource.RetryableError(
					fmt.Errorf("Delete scaling configuration timeout and got an error:%#v.", err))
			}
		}

		c, err := client.DescribeScalingConfigurationById(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		instances, _, err := client.essconn.DescribeScalingInstances(&ess.DescribeScalingInstancesArgs{
			RegionId:               getRegion(d, meta),
			ScalingGroupId:         c.ScalingGroupId,
			ScalingConfigurationId: d.Id(),
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(instances) > 0 {
			return resource.NonRetryableError(fmt.Errorf("There are still ECS instances in the scaling configuration - please remove them and try again."))
		}

		return resource.RetryableError(
			fmt.Errorf("Delete scaling configuration timeout and got an error:%#v.", err))
	})
}

func buildAlicloudEssScalingConfigurationArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingConfigurationArgs, error) {
	args := &ess.CreateScalingConfigurationArgs{
		ScalingGroupId:  d.Get("scaling_group_id").(string),
		ImageId:         d.Get("image_id").(string),
		InstanceType:    d.Get("instance_type").(string),
		SecurityGroupId: d.Get("security_group_id").(string),
	}

	if v := d.Get("scaling_configuration_name").(string); v != "" {
		args.ScalingConfigurationName = v
	}

	if v := d.Get("internet_charge_type").(string); v != "" {
		args.InternetChargeType = common.InternetChargeType(v)
	}

	if v := d.Get("internet_max_bandwidth_in").(int); v != 0 {
		args.InternetMaxBandwidthIn = v
	}

	if v := d.Get("internet_max_bandwidth_out").(int); v != 0 {
		args.InternetMaxBandwidthOut = v
	}

	if v := d.Get("system_disk_category").(string); v != "" {
		args.SystemDisk_Category = common.UnderlineString(v)
	}

	dds, ok := d.GetOk("data_disk")
	if ok {
		disks := dds.([]interface{})
		diskTypes := []ess.DataDiskType{}

		for _, e := range disks {
			pack := e.(map[string]interface{})
			disk := ess.DataDiskType{
				Size:       pack["size"].(int),
				Category:   pack["category"].(string),
				SnapshotId: pack["snapshot_id"].(string),
			}
			if v := pack["size"].(int); v != 0 {
				disk.Size = v
			}
			if v := pack["category"].(string); v != "" {
				disk.Category = v
			}
			if v := pack["snapshot_id"].(string); v != "" {
				disk.SnapshotId = v
			}
			diskTypes = append(diskTypes, disk)
		}
		args.DataDisk = diskTypes
	}

	if v, ok := d.GetOk("role_name"); ok && v.(string) != "" {
		args.RamRoleName = v.(string)
	}

	if v, ok := d.GetOk("key_name"); ok && v.(string) != "" {
		args.KeyPairName = v.(string)
	}

	if v, ok := d.GetOk("user_data"); ok && v.(string) != "" {
		args.UserData = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := "{"
		for key, value := range v.(map[string]interface{}) {
			tags += "\"" + key + "\"" + ":" + "\"" + value.(string) + "\"" + ","
		}
		args.Tags = strings.TrimSuffix(tags, ",") + "}"
	}

	return args, nil
}

func activeSubstituteScalingConfiguration(d *schema.ResourceData, meta interface{}) ([]ess.ScalingConfigurationItemType, error) {
	client := meta.(*AliyunClient)
	substitute_id, ok := d.GetOk("substitute")

	c, err := client.DescribeScalingConfigurationById(d.Id())
	if err != nil {
		return nil, fmt.Errorf("DescribeScalingConfigurationById error: %#v", err)
	}

	cs, _, err := client.essconn.DescribeScalingConfigurations(&ess.DescribeScalingConfigurationsArgs{
		RegionId:       getRegion(d, meta),
		ScalingGroupId: c.ScalingGroupId,
	})
	if err != nil {
		return nil, fmt.Errorf("DescribeScalingConfigurations error: %#v", err)
	}

	if !ok || substitute_id.(string) == "" {

		if len(cs) <= 1 {
			return cs, fmt.Errorf("Current scaling configuration %s is the last configuration for the scaling group %s, and it can't be inactive.", d.Id(), c.ScalingGroupId)
		}

		var configs []string
		for _, cc := range cs {
			if cc.ScalingConfigurationId != d.Id() {
				configs = append(configs, cc.ScalingConfigurationId)
			}
		}

		return cs, fmt.Errorf("Before inactivating current scaling configuration, you must select a substitute for scaling group from: %s.", strings.Join(configs, ","))
	}

	err = client.ActiveScalingConfigurationById(c.ScalingGroupId, substitute_id.(string))
	if err != nil {
		return cs, fmt.Errorf("Inactive scaling configuration %s err: %#v. Substitute scaling configuration ID: %s",
			d.Id(), err, substitute_id.(string))
	}

	return cs, nil
}
