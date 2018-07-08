package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"encoding/base64"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      40,
				ValidateFunc: validateIntegerInRange(40, 500),
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

			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn

	// Ensure instance_type is valid
	zoneId, validZones, err := meta.(*AliyunClient).DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}
	if err := meta.(*AliyunClient).InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
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

	resp, err := conn.CreateInstance(args)
	if err != nil {
		return fmt.Errorf("Error creating Aliyun ecs instance: %#v", err)
	}
	if resp == nil {
		return fmt.Errorf("Creating Ecs instance got a response: %#v.", resp)
	}

	d.SetId(resp.InstanceId)

	// after instance created, its status is pending,
	// so we need to wait it become to stopped and then start it
	if err := client.WaitForEcsInstance(d.Id(), Stopped, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Stopped, err)
	}

	out, err := ConvertIntegerToInt(args.InternetMaxBandwidthOut)
	if err != nil {
		return err
	}
	if out > 0 {
		req := ecs.CreateAllocatePublicIpAddressRequest()
		req.InstanceId = d.Id()
		if _, err := conn.AllocatePublicIpAddress(req); err != nil {
			return fmt.Errorf("[DEBUG] AllocatePublicIpAddress for instance got error: %#v", err)
		}
	}

	startArgs := ecs.CreateStartInstanceRequest()
	startArgs.InstanceId = d.Id()
	if _, err := conn.StartInstance(startArgs); err != nil {
		return fmt.Errorf("Start instance got error: %#v", err)
	}

	if err := client.WaitForEcsInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	return resourceAliyunInstanceUpdate(d, meta)
}

func resourceAliyunInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn

	instance, err := client.DescribeInstanceById(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error DescribeInstanceAttribute: %#v", err)
	}

	disk, diskErr := client.QueryInstanceSystemDisk(d.Id())

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
		resp, err := conn.DescribeUserData(args)

		if err != nil {
			log.Printf("[ERROR] DescribeUserData for instance got error: %#v", err)
		}
		if resp != nil {
			d.Set("user_data", userDataHashSum(resp.UserData))
		}
	}

	if len(instance.VpcAttributes.VSwitchId) > 0 {
		args := ecs.CreateDescribeInstanceRamRoleRequest()
		args.InstanceIds = convertListToJsonString([]interface{}{d.Id()})
		response, err := conn.DescribeInstanceRamRole(args)
		if err != nil {
			return fmt.Errorf("[ERROR] DescribeInstanceRamRole for instance got error: %#v", err)
		}

		if response != nil && len(response.InstanceRamRoleSets.InstanceRamRoleSet) > 1 {
			d.Set("role_name", response.InstanceRamRoleSets.InstanceRamRoleSet[0].RamRoleName)
		}
	}

	if instance.InstanceChargeType == string(PrePaid) {
		args := ecs.CreateDescribeInstanceAutoRenewAttributeRequest()
		args.InstanceId = d.Id()
		resp, err := conn.DescribeInstanceAutoRenewAttribute(args)
		if err != nil {
			return fmt.Errorf("DescribeInstanceAutoRenewAttribute got an error: %#v.", err)
		}
		if resp != nil && len(resp.InstanceRenewAttributes.InstanceRenewAttribute) > 0 {
			renew := resp.InstanceRenewAttributes.InstanceRenewAttribute[0]
			d.Set("renewal_status", renew.RenewalStatus)
			d.Set("auto_renew_period", renew.Duration)
		}

	}
	tags, err := client.DescribeTags(d.Id(), TagResourceInstance)
	if err != nil && !NotFoundError(err) {
		return fmt.Errorf("[ERROR] DescribeTags for instance got error: %#v", err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn

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
			err := client.JoinSecurityGroups(d.Id(), al)
			if err != nil {
				return err
			}
		}
		if len(rl) > 0 {
			err := client.LeaveSecurityGroups(d.Id(), rl)
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

		if _, err := client.ecsconn.ModifyInstanceAutoRenewAttribute(args); err != nil {
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
		instance, errDesc := client.DescribeInstanceById(d.Id())
		if errDesc != nil {
			return fmt.Errorf("Describe instance got an error: %#v", errDesc)
		}
		if instance.Status == string(Running) {
			log.Printf("[DEBUG] Stop instance when changing image or password or vpc attribute")
			stop := ecs.CreateStopInstanceRequest()
			stop.InstanceId = d.Id()
			stop.ForceStop = requests.NewBoolean(false)
			if _, err := conn.StopInstance(stop); err != nil {
				return fmt.Errorf("StopInstance got error: %#v", err)
			}
		}

		if err := client.WaitForEcsInstance(d.Id(), Stopped, DefaultTimeout); err != nil {
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

		if _, err := conn.StartInstance(start); err != nil {
			return fmt.Errorf("StartInstance got error: %#v", err)
		}

		// Start instance sometimes costs more than 8 minutes when os type is centos.
		if err := client.WaitForEcsInstance(d.Id(), Running, DefaultTimeoutMedium); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
	}

	if err := modifyInstanceNetworkSpec(d, meta); err != nil {
		return err
	}

	if err := modifyInstanceChargeType(d, meta); err != nil {
		return err
	}

	d.Partial(false)
	return resourceAliyunInstanceRead(d, meta)
}

func resourceAliyunInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.ecsconn
	if d.Get("instance_charge_type").(string) == string(PrePaid) {
		return fmt.Errorf("At present, 'PrePaid' instance cannot be deleted and must wait it to be expired and release it automatically.")
	}
	stop := ecs.CreateStopInstanceRequest()
	stop.InstanceId = d.Id()
	stop.ForceStop = requests.NewBoolean(true)

	deld := ecs.CreateDeleteInstanceRequest()
	deld.InstanceId = d.Id()
	deld.Force = requests.NewBoolean(true)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		instance, err := client.DescribeInstanceById(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
		}

		if instance.Status != string(Stopped) {
			if _, err := conn.StopInstance(stop); err != nil {
				return resource.RetryableError(fmt.Errorf("Stop instance timeout and got an error: %#v.", err))
			}

			if err := client.WaitForEcsInstance(d.Id(), Stopped, DefaultTimeout); err != nil {
				return resource.RetryableError(fmt.Errorf("Waiting for ecs stopped timeout and got an error: %#v.", err))
			}
		}

		if _, err := conn.DeleteInstance(deld); err != nil {
			if NotFoundError(err) || IsExceptedErrors(err, EcsNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete instance timeout and got an error: %#v.", err))
		}

		return nil
	})

}

func buildAliyunInstanceArgs(d *schema.ResourceData, meta interface{}) (*ecs.CreateInstanceRequest, error) {
	client := meta.(*AliyunClient)

	args := ecs.CreateCreateInstanceRequest()
	args.InstanceType = d.Get("instance_type").(string)

	imageID := d.Get("image_id").(string)

	args.ImageId = imageID

	systemDiskCategory := DiskCategory(d.Get("system_disk_category").(string))

	zoneID := d.Get("availability_zone").(string)
	// check instanceType and systemDiskCategory, when zoneID is not empty
	if zoneID != "" {
		zone, err := client.DescribeZone(zoneID)
		if err != nil {
			return nil, err
		}

		if err := client.ResourceAvailable(zone, ResourceTypeInstance); err != nil {
			return nil, err
		}

		if err := client.DiskAvailable(zone, systemDiskCategory); err != nil {
			return nil, err
		}

		args.ZoneId = zoneID

	}

	args.SystemDiskCategory = string(systemDiskCategory)
	args.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))

	sgs, ok := d.GetOk("security_groups")

	if ok {
		sgList := expandStringList(sgs.(*schema.Set).List())
		sg0 := sgList[0]
		// check security group instance exist
		_, err := client.DescribeSecurityGroupAttribute(sg0)
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

	return args, nil
}

func modifyInstanceChargeType(d *schema.ResourceData, meta interface{}) error {
	if d.IsNewResource() {
		return nil
	}

	conn := meta.(*AliyunClient).ecsconn

	if d.HasChange("instance_charge_type") {
		chargeType := d.Get("instance_charge_type").(string)
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
			if _, err := conn.ModifyInstanceChargeType(args); err != nil {
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

		d.SetPartial("instance_charge_type")
		return nil
	}

	return nil
}

func modifyInstanceImage(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		return false, nil
	}
	client := meta.(*AliyunClient)
	update := false
	if d.HasChange("image_id") {
		update = true
		if !run {
			return update, nil
		}
		args := ecs.CreateReplaceSystemDiskRequest()
		args.InstanceId = d.Id()
		args.ImageId = d.Get("image_id").(string)
		args.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))

		_, err := client.ecsconn.ReplaceSystemDisk(args)
		if err != nil {
			return update, fmt.Errorf("Replace system disk got an error: %#v", err)
		}

		// Ensure instance's image has been replaced successfully.
		timeout := DefaultTimeoutMedium
		for {
			instance, errDesc := client.DescribeInstanceById(d.Id())
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
		if _, err := meta.(*AliyunClient).ecsconn.ModifyInstanceAttribute(args); err != nil {
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
		if _, err := meta.(*AliyunClient).ecsconn.ModifyInstanceVpcAttribute(vpcArgs); err != nil {
			return update, fmt.Errorf("ModifyInstanceVPCAttribute got an error: %#v.", err)
		}
	}
	return update, nil
}

func modifyInstanceType(d *schema.ResourceData, meta interface{}, run bool) (bool, error) {
	if d.IsNewResource() {
		return false, nil
	}
	client := meta.(*AliyunClient)
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
		zoneId, validZones, err := meta.(*AliyunClient).DescribeAvailableResources(d, meta, InstanceTypeResource)
		if err != nil {
			return update, err
		}
		if err := meta.(*AliyunClient).InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
			return update, err
		}

		d.SetPartial("instance_type")

		//An instance that was successfully modified once cannot be modified again within 5 minutes.
		args := ecs.CreateModifyInstanceSpecRequest()
		args.InstanceId = d.Id()
		args.InstanceType = d.Get("instance_type").(string)

		err = resource.Retry(6*time.Minute, func() *resource.RetryError {
			if _, err := client.ecsconn.ModifyInstanceSpec(args); err != nil {
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
	if update {
		if err := resource.Retry(6*time.Minute, func() *resource.RetryError {
			if _, err := meta.(*AliyunClient).ecsconn.ModifyInstanceNetworkSpec(args); err != nil {
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
			if _, err := meta.(*AliyunClient).ecsconn.AllocatePublicIpAddress(req); err != nil {
				return fmt.Errorf("[DEBUG] AllocatePublicIpAddress for instance got error: %#v", err)
			}
		}
	}
	return nil
}
