package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"encoding/base64"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunInstanceCreate,
		Read:   resourceAliyunInstanceRead,
		Update: resourceAliyunInstanceUpdate,
		Delete: resourceAliyunInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateInstanceType,
			},

			"security_groups": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"allocate_public_ip": &schema.Schema{
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Field 'allocate_public_ip' has been deprecated from provider version 1.6.1. Setting 'internet_max_bandwidth_out' larger than 0 will allocate public ip for instance.",
			},

			"instance_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ECS-Instance",
				ValidateFunc: validateInstanceName,
			},

			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceDescription,
			},

			"internet_charge_type": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validateInternetChargeType,
				Default:          PayByTraffic,
				DiffSuppressFunc: ecsInternetDiffSuppressFunc,
			},
			"internet_max_bandwidth_in": &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateFunc:     validateIntegerInRange(1, 200),
				Computed:         true,
				DiffSuppressFunc: ecsInternetDiffSuppressFunc,
			},
			"internet_max_bandwidth_out": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 100),
			},
			"host_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
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
			"system_disk_category": &schema.Schema{
				Type:         schema.TypeString,
				Default:      DiskCloudEfficiency,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
			},
			"system_disk_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40,
			},
			"data_disks": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 15,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validateDiskName,
						},
						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"category": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateDiskCategory,
							Default:      DiskCloudEfficiency,
							ForceNew:     true,
						},
						"encrypted": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"snapshot_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"delete_with_instance": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"description": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validateDiskDescription,
						},
					},
				},
			},

			//subnet_id and vswitch_id both exists, cause compatible old version, and aws habit.
			"subnet_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true, //add this schema cause subnet_id not used enter parameter, will different, so will be ForceNew
				ConflictsWith: []string{"vswitch_id"},
			},

			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instance_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceChargeType,
				Default:      PostPaid,
			},
			"period": &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validateInstanceChargeTypePeriod,
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
			},
			"period_unit": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				Default:          Month,
				ValidateFunc:     validateInstanceChargeTypePeriodUnit,
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
			},
			"renewal_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  RenewNormal,
				ValidateFunc: validateAllowedStringValue([]string{
					string(RenewAutoRenewal),
					string(RenewNormal),
					string(RenewNotRenewal)}),
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
			},
			"auto_renew_period": &schema.Schema{
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: ecsNotAutoRenewDiffSuppressFunc,
			},
			"include_data_disks": &schema.Schema{
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          true,
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
			},
			"dry_run": &schema.Schema{
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
			},

			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"role_name": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: vpcTypeResourceDiffSuppressFunc,
			},

			"key_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"spot_strategy": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          NoSpot,
				ValidateFunc:     validateInstanceSpotStrategy,
				DiffSuppressFunc: ecsSpotStrategyDiffSuppressFunc,
			},

			"spot_price_limit": &schema.Schema{
				Type:             schema.TypeFloat,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: ecsSpotPriceLimitDiffSuppressFunc,
			},

			"deletion_protection": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"force_delete": &schema.Schema{
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				Description:      descriptions["A behavior mark used to delete 'PrePaid' ECS instance forcibly."],
				DiffSuppressFunc: ecsPostPaidDiffSuppressFunc,
			},

			"security_enhancement_strategy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ActiveSecurityEnhancementStrategy),
					string(DeactiveSecurityEnhancementStrategy),
				}),
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	// Ensure instance_type is valid
	zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}
	if err := ecsService.InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
		return err
	}

	args, err := buildAliyunInstanceArgs(d, meta)
	if err != nil {
		return err
	}
	args.IoOptimized = "optimized"
	if d.Get("is_outdated").(bool) == true {
		args.IoOptimized = "none"
	}

	err = resource.Retry(DefaultTimeout*time.Second, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.RunInstances(args)
		})
		if err != nil {
			if IsExceptedError(err, InvalidPrivateIpAddressDuplicated) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(fmt.Errorf("Error creating Aliyun ecs instance: %#v", err))
		}
		resp, _ := raw.(*ecs.RunInstancesResponse)
		if resp == nil {
			return resource.NonRetryableError(fmt.Errorf("Creating Ecs instance got a response: %#v.", resp))
		}

		if len(resp.InstanceIdSets.InstanceIdSet) != 1 {
			return resource.NonRetryableError(fmt.Errorf("run instance failed, invalid instance ID list: %#v", resp.InstanceIdSets.InstanceIdSet))
		}

		d.SetId(resp.InstanceIdSets.InstanceIdSet[0])

		return nil
	})
	if err != nil {
		return err
	}

	if err := ecsService.WaitForEcsInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	return resourceAliyunInstanceUpdate(d, meta)
}

func resourceAliyunInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	instance, err := ecsService.DescribeInstanceById(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error DescribeInstanceAttribute: %#v", err)
	}

	disk, diskErr := ecsService.QueryInstanceSystemDisk(d.Id())

	if diskErr != nil {
		if NotFoundError(diskErr) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error DescribeSystemDisk: %#v", err)
	}

	d.Set("instance_name", instance.InstanceName)
	d.Set("description", instance.Description)
	d.Set("status", instance.Status)
	d.Set("availability_zone", instance.ZoneId)
	d.Set("host_name", instance.HostName)
	d.Set("image_id", instance.ImageId)
	d.Set("instance_type", instance.InstanceType)
	d.Set("system_disk_category", disk.Category)
	d.Set("system_disk_size", disk.Size)
	d.Set("password", d.Get("password"))
	d.Set("internet_max_bandwidth_out", instance.InternetMaxBandwidthOut)
	d.Set("internet_max_bandwidth_in", instance.InternetMaxBandwidthIn)
	d.Set("instance_charge_type", instance.InstanceChargeType)
	d.Set("key_name", instance.KeyPairName)
	d.Set("spot_strategy", instance.SpotStrategy)
	d.Set("spot_price_limit", instance.SpotPriceLimit)
	d.Set("internet_charge_type", instance.InternetChargeType)
	d.Set("deletion_protection", instance.DeletionProtection)

	if len(instance.PublicIpAddress.IpAddress) > 0 {
		d.Set("public_ip", instance.PublicIpAddress.IpAddress[0])
	} else {
		d.Set("public_ip", "")
	}

	d.Set("subnet_id", instance.VpcAttributes.VSwitchId)
	d.Set("vswitch_id", instance.VpcAttributes.VSwitchId)

	if len(instance.VpcAttributes.PrivateIpAddress.IpAddress) > 0 {
		d.Set("private_ip", instance.VpcAttributes.PrivateIpAddress.IpAddress[0])
	} else {
		d.Set("private_ip", strings.Join(instance.InnerIpAddress.IpAddress, ","))
	}

	sgs := make([]string, 0, len(instance.SecurityGroupIds.SecurityGroupId))
	for _, sg := range instance.SecurityGroupIds.SecurityGroupId {
		sgs = append(sgs, sg)
	}
	log.Printf("[DEBUG] Setting Security Group Ids: %#v", sgs)
	if err := d.Set("security_groups", sgs); err != nil {
		return err
	}

	if d.Get("user_data").(string) != "" {
		args := ecs.CreateDescribeUserDataRequest()
		args.InstanceId = d.Id()
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeUserData(args)
		})

		if err != nil {
			log.Printf("[ERROR] DescribeUserData for instance got error: %#v", err)
		}
		resp, _ := raw.(*ecs.DescribeUserDataResponse)
		if resp != nil {
			d.Set("user_data", userDataHashSum(resp.UserData))
		}
	}

	if len(instance.VpcAttributes.VSwitchId) > 0 {
		args := ecs.CreateDescribeInstanceRamRoleRequest()
		args.InstanceIds = convertListToJsonString([]interface{}{d.Id()})
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(args)
		})
		if err != nil {
			return fmt.Errorf("[ERROR] DescribeInstanceRamRole for instance got error: %#v", err)
		}
		response, _ := raw.(*ecs.DescribeInstanceRamRoleResponse)
		if response != nil && len(response.InstanceRamRoleSets.InstanceRamRoleSet) > 1 {
			d.Set("role_name", response.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName)
		}
	}

	if instance.InstanceChargeType == string(PrePaid) {
		args := ecs.CreateDescribeInstanceAutoRenewAttributeRequest()
		args.InstanceId = d.Id()
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceAutoRenewAttribute(args)
		})
		if err != nil {
			return fmt.Errorf("DescribeInstanceAutoRenewAttribute got an error: %#v.", err)
		}
		resp, _ := raw.(*ecs.DescribeInstanceAutoRenewAttributeResponse)
		if resp != nil && len(resp.InstanceRenewAttributes.InstanceRenewAttribute) > 0 {
			renew := resp.InstanceRenewAttributes.InstanceRenewAttribute[0]
			d.Set("renewal_status", renew.RenewalStatus)
			d.Set("auto_renew_period", renew.Duration)
		}

	}
	tags, err := ecsService.DescribeTags(d.Id(), TagResourceInstance)
	if err != nil && !NotFoundError(err) {
		return fmt.Errorf("[ERROR] DescribeTags for instance got error: %#v", err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	d.Partial(true)

	if err := setTags(client, TagResourceInstance, d); err != nil {
		log.Printf("[DEBUG] Set tags for instance got error: %#v", err)
		return fmt.Errorf("Set tags for instance got error: %#v", err)
	} else {
		d.SetPartial("tags")
	}

	if d.HasChange("security_groups") {
		o, n := d.GetChange("security_groups")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		rl := expandStringList(os.Difference(ns).List())
		al := expandStringList(ns.Difference(os).List())

		if len(al) > 0 {
			err := ecsService.JoinSecurityGroups(d.Id(), al)
			if err != nil {
				return err
			}
		}
		if len(rl) > 0 {
			err := ecsService.LeaveSecurityGroups(d.Id(), rl)
			if err != nil {
				return err
			}
		}

		d.SetPartial("security_groups")
	}
	if d.HasChange("renewal_status") || d.HasChange("auto_renew_period") {
		status := d.Get("renewal_status").(string)
		args := ecs.CreateModifyInstanceAutoRenewAttributeRequest()
		args.InstanceId = d.Id()
		args.RenewalStatus = status

		if status == string(RenewAutoRenewal) {
			args.Duration = requests.NewInteger(d.Get("auto_renew_period").(int))
		}

		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceAutoRenewAttribute(args)
		})
		if err != nil {
			return fmt.Errorf("ModifyInstanceAutoRenewAttribute got an error: %#v", err)
		}
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	run := false
	imageUpdate, err := modifyInstanceImage(d, meta, run)
	if err != nil {
		return err
	}

	vpcUpdate, err := modifyVpcAttribute(d, meta, run)
	if err != nil {
		return err
	}

	passwordUpdate, err := modifyInstanceAttribute(d, meta)
	if err != nil {
		return err
	}

	typeUpdate, err := modifyInstanceType(d, meta, run)
	if err != nil {
		return err
	}
	if imageUpdate || vpcUpdate || passwordUpdate || typeUpdate {
		run = true
		log.Printf("[INFO] Need rebooting to make all changes valid.")
		instance, errDesc := ecsService.DescribeInstanceById(d.Id())
		if errDesc != nil {
			return fmt.Errorf("Describe instance got an error: %#v", errDesc)
		}
		if instance.Status == string(Running) {
			log.Printf("[DEBUG] Stop instance when changing image or password or vpc attribute")
			stop := ecs.CreateStopInstanceRequest()
			stop.InstanceId = d.Id()
			stop.ForceStop = requests.NewBoolean(false)
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.StopInstance(stop)
			})
			if err != nil {
				return fmt.Errorf("StopInstance got error: %#v", err)
			}
		}

		if err := ecsService.WaitForEcsInstance(d.Id(), Stopped, DefaultTimeout); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Stopped, err)
		}

		if _, err := modifyInstanceImage(d, meta, run); err != nil {
			return err
		}

		if _, err := modifyVpcAttribute(d, meta, run); err != nil {
			return err
		}

		if _, err := modifyInstanceType(d, meta, run); err != nil {
			return err
		}

		log.Printf("[DEBUG] Start instance after changing image or password or vpc attribute")
		start := ecs.CreateStartInstanceRequest()
		start.InstanceId = d.Id()

		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.StartInstance(start)
		})
		if err != nil {
			return fmt.Errorf("StartInstance got error: %#v", err)
		}

		// Start instance sometimes costs more than 8 minutes when os type is centos.
		if err := ecsService.WaitForEcsInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
	}

	if err := modifyInstanceNetworkSpec(d, meta); err != nil {
		return err
	}

	updateRenewal := false
	if d.HasChange("instance_charge_type") {
		if _, n := d.GetChange("instance_charge_type"); n.(string) == string(PrePaid) {
			updateRenewal = true
		}
	}
	if err := modifyInstanceChargeType(d, meta, false); err != nil {
		return err
	}

	// Only PrePaid instance can support modifying renewal attribute
	if updateRenewal && (d.HasChange("renewal_status") || d.HasChange("auto_renew_period")) {
		status := d.Get("renewal_status").(string)
		args := ecs.CreateModifyInstanceAutoRenewAttributeRequest()
		args.InstanceId = d.Id()
		args.RenewalStatus = status

		if status == string(RenewAutoRenewal) {
			args.Duration = requests.NewInteger(d.Get("auto_renew_period").(int))
		}

		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceAutoRenewAttribute(args)
		})
		if err != nil {
			return fmt.Errorf("ModifyInstanceAutoRenewAttribute got an error: %#v", err)
		}
		d.SetPartial("renewal_status")
		d.SetPartial("auto_renew_period")
	}

	if d.HasChange("force_delete") {
		d.SetPartial("force_delete")
	}

	d.Partial(false)
	return resourceAliyunInstanceRead(d, meta)
}

func resourceAliyunInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		force := d.Get("force_delete").(bool)
		if !force {
			return fmt.Errorf("Please convert 'PrePaid' instance to 'PostPaid' or set 'force_delete' as true before deleting 'PrePaid' instance.")
		} else if err := modifyInstanceChargeType(d, meta, force); err != nil {
			return fmt.Errorf("Before deleteing ECS instance forcibly, converting instance charge type got an error: %#v.", err)
		}
	}
	stop := ecs.CreateStopInstanceRequest()
	stop.InstanceId = d.Id()
	stop.ForceStop = requests.NewBoolean(true)

	deld := ecs.CreateDeleteInstanceRequest()
	deld.InstanceId = d.Id()
	deld.Force = requests.NewBoolean(true)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		instance, err := ecsService.DescribeInstanceById(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
		}

		if instance.Status != string(Stopped) {
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.StopInstance(stop)
			})
			if err != nil {
				return resource.RetryableError(fmt.Errorf("Stop instance timeout and got an error: %#v.", err))
			}

			if err := ecsService.WaitForEcsInstance(d.Id(), Stopped, DefaultTimeout); err != nil {
				return resource.RetryableError(fmt.Errorf("Waiting for ecs stopped timeout and got an error: %#v.", err))
			}
		}

		_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteInstance(deld)
		})
		if err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, EcsNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete instance timeout and got an error: %#v.", err))
		}

		return nil
	})

}

func buildAliyunInstanceArgs(d *schema.ResourceData, meta interface{}) (*ecs.RunInstancesRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	args := ecs.CreateRunInstancesRequest()
	args.InstanceType = d.Get("instance_type").(string)

	imageID := d.Get("image_id").(string)

	args.ImageId = imageID

	systemDiskCategory := DiskCategory(d.Get("system_disk_category").(string))

	zoneID := d.Get("availability_zone").(string)
	// check instanceType and systemDiskCategory, when zoneID is not empty
	if zoneID != "" {
		zone, err := ecsService.DescribeZone(zoneID)
		if err != nil {
			return nil, err
		}

		if err := ecsService.ResourceAvailable(zone, ResourceTypeInstance); err != nil {
			return nil, err
		}

		if err := ecsService.DiskAvailable(zone, systemDiskCategory); err != nil {
			return nil, err
		}

		args.ZoneId = zoneID

	}

	args.SystemDiskCategory = string(systemDiskCategory)
	args.SystemDiskSize = strconv.Itoa(d.Get("system_disk_size").(int))

	sgs, ok := d.GetOk("security_groups")

	if ok {
		sgList := expandStringList(sgs.(*schema.Set).List())
		sg0 := sgList[0]
		// check security group instance exist
		_, err := ecsService.DescribeSecurityGroupAttribute(sg0)
		if err == nil {
			args.SecurityGroupId = sg0
		}
	}

	if v := d.Get("instance_name").(string); v != "" {
		args.InstanceName = v
	}

	if v := d.Get("description").(string); v != "" {
		args.Description = v
	}

	if v := d.Get("internet_charge_type").(string); v != "" {
		args.InternetChargeType = v
	}

	args.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))

	if v := d.Get("host_name").(string); v != "" {
		args.HostName = v
	}

	if v := d.Get("password").(string); v != "" {
		args.Password = v
	}

	vswitchValue := d.Get("subnet_id").(string)
	if vswitchValue == "" {
		vswitchValue = d.Get("vswitch_id").(string)
	}
	if vswitchValue != "" {
		args.VSwitchId = vswitchValue
		if v, ok := d.GetOk("private_ip"); ok && v.(string) != "" {
			args.PrivateIpAddress = v.(string)
		}
	}

	if v := d.Get("instance_charge_type").(string); v != "" {
		args.InstanceChargeType = v
	}

	if args.InstanceChargeType == string(PrePaid) {
		args.Period = requests.NewInteger(d.Get("period").(int))
		args.PeriodUnit = d.Get("period_unit").(string)
	} else {
		if v := d.Get("spot_strategy").(string); v != "" {
			args.SpotStrategy = v
		}
		if v := d.Get("spot_price_limit").(float64); v > 0 {
			args.SpotPriceLimit = requests.NewFloat(v)
		}
	}

	if v := d.Get("user_data").(string); v != "" {
		args.UserData = base64.StdEncoding.EncodeToString([]byte(v))
	}

	if v := d.Get("role_name").(string); v != "" {
		if vswitchValue == "" {
			return nil, fmt.Errorf("Role name only supported for VPC instance.")
		}
		args.RamRoleName = v
	}

	if v := d.Get("key_name").(string); v != "" {
		args.KeyPairName = v
	}

	if v, ok := d.GetOk("security_enhancement_strategy"); ok {
		value := v.(string)
		args.SecurityEnhancementStrategy = value
	}

	args.DeletionProtection = requests.NewBoolean(d.Get("deletion_protection").(bool))
	args.ClientToken = buildClientToken("TF-CreateInstance")

	if v, ok := d.GetOk("data_disks"); ok {
		disks := v.([]interface{})
		var dataDiskRequests []ecs.RunInstancesDataDisk
		for i := range disks {
			disk := disks[i].(map[string]interface{})

			req := ecs.RunInstancesDataDisk{
				Category:           disk["category"].(string),
				DeleteWithInstance: strconv.FormatBool(disk["delete_with_instance"].(bool)),
				Encrypted:          strconv.FormatBool(disk["encrypted"].(bool)),
			}

			if name, ok := disk["name"]; ok {
				req.DiskName = name.(string)
			}
			if snapshotId, ok := disk["snapshot_id"]; ok {
				req.SnapshotId = snapshotId.(string)
			}
			if description, ok := disk["description"]; ok {
				req.Description = description.(string)
			}
			req.Size = fmt.Sprintf("%d", disk["size"].(int))
			req.Category = disk["category"].(string)
			if req.Category == string(DiskEphemeralSSD) {
				req.DeleteWithInstance = ""
			}

			dataDiskRequests = append(dataDiskRequests, req)
		}
		args.DataDisk = &dataDiskRequests
	}
	return args, nil
}

func modifyInstanceChargeType(d *schema.ResourceData, meta interface{}, forceDelete bool) error {
	if d.IsNewResource() {
		return nil
	}

	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	chargeType := d.Get("instance_charge_type").(string)
	if d.HasChange("instance_charge_type") || forceDelete {
		if forceDelete {
			chargeType = string(PostPaid)
		}
		args := ecs.CreateModifyInstanceChargeTypeRequest()
		args.InstanceIds = convertListToJsonString(append(make([]interface{}, 0, 1), d.Id()))
		args.IncludeDataDisks = requests.NewBoolean(d.Get("include_data_disks").(bool))
		args.AutoPay = requests.NewBoolean(true)
		args.DryRun = requests.NewBoolean(d.Get("dry_run").(bool))
		args.ClientToken = fmt.Sprintf("terraform-modify-instance-charge-type-%s", d.Id())
		if chargeType == string(PrePaid) {
			args.Period = requests.NewInteger(d.Get("period").(int))
			args.PeriodUnit = d.Get("period_unit").(string)
		}
		args.InstanceChargeType = chargeType
		if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceChargeType(args)
			})
			if err != nil {
				if IsExceptedErrors(err, []string{Throttling}) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(fmt.Errorf("Modifying instance %s chareType timeout and got an error:%#v.", d.Id(), err))
				}
				return resource.NonRetryableError(fmt.Errorf("Modifying instance %s chareType timeout and got an error: %#v.", d.Id(), err))
			}
			return nil
		}); err != nil {
			return err
		}
		// Wait for instance charge type has been changed
		if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
			if instance, err := ecsService.DescribeInstanceById(d.Id()); err != nil {
				return resource.NonRetryableError(fmt.Errorf("Describing instance %s got an error: %#v.", d.Id(), err))
			} else if instance.InstanceChargeType == chargeType {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Waitting for instance %s to be %s timeout.", d.Id(), chargeType))
		}); err != nil {
			return err
		}

		d.SetPartial("instance_charge_type")
		return nil
	}

	return nil
}

func modifyInstanceImage(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		return false, nil
	}
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChange("image_id") {
		update = true
		if !run {
			return update, nil
		}
		instance, e := ecsService.DescribeInstanceById(d.Id())
		if e != nil {
			return update, e
		}
		keyPairName := instance.KeyPairName
		args := ecs.CreateReplaceSystemDiskRequest()
		args.InstanceId = d.Id()
		args.ImageId = d.Get("image_id").(string)
		args.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
		args.ClientToken = buildClientToken("TF-ReplaceSystemDisk")
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ReplaceSystemDisk(args)
		})
		if err != nil {
			return update, fmt.Errorf("Replace system disk got an error: %#v", err)
		}

		// Ensure instance's image has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, errDesc := ecsService.DescribeInstanceById(d.Id())
			if errDesc != nil {
				return update, fmt.Errorf("Describe instance got an error: %#v", errDesc)
			}

			if instance.ImageId == d.Get("image_id") {
				break
			}
			time.Sleep(DefaultIntervalShort * time.Second)

			timeout = timeout - DefaultIntervalShort
			if timeout <= 0 {
				return update, GetTimeErrorFromString(fmt.Sprintf("Replacing instance %s system disk timeout.", d.Id()))
			}
		}

		d.SetPartial("system_disk_size")
		d.SetPartial("image_id")

		// After updating image, it need to re-attach key pair
		if keyPairName != "" {
			if err := ecsService.AttachKeyPair(keyPairName, []interface{}{d.Id()}); err != nil {
				return update, fmt.Errorf("After updating image, attaching key pair %s got an error: %#v.", keyPairName, err)
			}
		}
	}
	// Provider doesn't support change 'system_disk_size'separately.
	if d.HasChange("system_disk_size") && !d.HasChange("image_id") {
		return update, fmt.Errorf("Update resource failed. 'system_disk_size' isn't allowed to change separately. You can update it via renewing instance or replacing system disk.")
	}
	return update, nil
}

func modifyInstanceAttribute(d *schema.ResourceData, meta interface{}) (bool, error) {
	if d.IsNewResource() {
		return false, nil
	}

	update := false
	reboot := false
	args := ecs.CreateModifyInstanceAttributeRequest()
	args.InstanceId = d.Id()

	if d.HasChange("instance_name") {
		log.Printf("[DEBUG] ModifyInstanceAttribute instance_name")
		d.SetPartial("instance_name")
		args.InstanceName = d.Get("instance_name").(string)
		update = true
	}

	if d.HasChange("description") {
		log.Printf("[DEBUG] ModifyInstanceAttribute description")
		d.SetPartial("description")
		args.Description = d.Get("description").(string)
		update = true
	}

	if d.HasChange("host_name") {
		log.Printf("[DEBUG] ModifyInstanceAttribute host_name")
		d.SetPartial("host_name")
		args.HostName = d.Get("host_name").(string)
		update = true
		reboot = true
	}

	if d.HasChange("password") {
		log.Printf("[DEBUG] ModifyInstanceAttribute password")
		d.SetPartial("password")
		args.Password = d.Get("password").(string)
		update = true
		reboot = true
	}

	if update {
		client := meta.(*connectivity.AliyunClient)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceAttribute(args)
		})
		if err != nil {
			return reboot, fmt.Errorf("Modify instance attribute got error: %#v", err)
		}
	}
	return reboot, nil
}

func modifyVpcAttribute(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		return false, nil
	}

	update := false
	vpcArgs := ecs.CreateModifyInstanceVpcAttributeRequest()
	vpcArgs.InstanceId = d.Id()
	vpcArgs.VSwitchId = d.Get("vswitch_id").(string)

	if d.HasChange("vswitch_id") {
		update = true
		if d.Get("vswitch_id").(string) == "" {
			return update, fmt.Errorf("Field 'vswitch_id' is required when modifying the instance VPC attribute.")
		}
		d.SetPartial("vswitch_id")
	}

	if d.HasChange("subnet_id") {
		update = true
		if d.Get("subnet_id").(string) == "" {
			return update, fmt.Errorf("Field 'subnet_id' is required when modifying the instance VPC attribute.")
		}
		vpcArgs.VSwitchId = d.Get("subnet_id").(string)
		d.SetPartial("subnet_id")
	}

	if vpcArgs.VSwitchId != "" && d.HasChange("private_ip") {
		vpcArgs.PrivateIpAddress = d.Get("private_ip").(string)
		update = true
		d.SetPartial("private_ip")
	}

	if !run {
		return update, nil
	}

	if update {
		client := meta.(*connectivity.AliyunClient)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyInstanceVpcAttribute(vpcArgs)
		})
		if err != nil {
			return update, fmt.Errorf("ModifyInstanceVPCAttribute got an error: %#v.", err)
		}
	}
	return update, nil
}

func modifyInstanceType(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		return false, nil
	}
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	update := false
	if d.HasChange("instance_type") {
		update = true
		if !run {
			return update, nil
		}
		if d.Get("instance_charge_type").(string) == string(PrePaid) {
			return update, fmt.Errorf("At present, 'PrePaid' instance type cannot be modified.")
		}
		// Ensure instance_type is valid
		zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
		if err != nil {
			return update, err
		}
		if err := ecsService.InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
			return update, err
		}

		d.SetPartial("instance_type")

		//An instance that was successfully modified once cannot be modified again within 5 minutes.
		args := ecs.CreateModifyInstanceSpecRequest()
		args.InstanceId = d.Id()
		args.InstanceType = d.Get("instance_type").(string)
		args.ClientToken = buildClientToken("TF-ModifyInstanceSpec")

		err = resource.Retry(6*time.Minute, func() *resource.RetryError {
			_, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceSpec(args)
			})
			if err != nil {
				if IsExceptedError(err, EcsThrottling) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(fmt.Errorf("Modify instance type timeout and got an error; %#v", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Modify instance type got an error: %#v", err))
			}
			return nil
		})
		return update, err
	}
	return update, nil
}

func modifyInstanceNetworkSpec(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}

	allocate := false
	update := false
	args := ecs.CreateModifyInstanceNetworkSpecRequest()
	args.InstanceId = d.Id()
	args.ClientToken = buildClientToken("TF-ModifyInstanceNetworkSpec")

	if d.HasChange("internet_charge_type") {
		args.NetworkChargeType = d.Get("internet_charge_type").(string)
		update = true
		d.SetPartial("internet_charge_type")
	}

	if d.HasChange("internet_max_bandwidth_out") {
		o, n := d.GetChange("internet_max_bandwidth_out")
		if o.(int) <= 0 && n.(int) > 0 {
			allocate = true
		}
		args.InternetMaxBandwidthOut = requests.NewInteger(n.(int))
		update = true
		d.SetPartial("internet_max_bandwidth_out")
	}

	if d.HasChange("internet_max_bandwidth_in") {
		args.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
		update = true
		d.SetPartial("internet_max_bandwidth_in")
	}

	//An instance that was successfully modified once cannot be modified again within 5 minutes.
	client := meta.(*connectivity.AliyunClient)
	if update {
		if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ModifyInstanceNetworkSpec(args)
			})
			if err != nil {
				if IsExceptedError(err, EcsThrottling) {
					time.Sleep(10 * time.Second)
					return resource.RetryableError(fmt.Errorf("Modify instance network bandwidth timeout and got an error; %#v", err))
				}
				if IsExceptedError(err, EcsInternalError) {
					return resource.RetryableError(fmt.Errorf("Modify instance network bandwidth timeout and got an error; %#v", err))
				}
				return resource.NonRetryableError(fmt.Errorf("Modify instance network bandwidth got an error: %#v", err))
			}
			return nil
		}); err != nil {
			return err
		}
		if allocate {
			req := ecs.CreateAllocatePublicIpAddressRequest()
			req.InstanceId = d.Id()
			_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.AllocatePublicIpAddress(req)
			})
			if err != nil {
				return fmt.Errorf("[DEBUG] AllocatePublicIpAddress for instance got error: %#v", err)
			}
		}
	}
	return nil
}
