package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunLaunchTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunLaunchTemplateCreate,
		Read:   resourceAliyunLaunchTemplateRead,
		Delete: resourceAliyunLaunchTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateLaunchTemplateName,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateLaunchTemplateDescription,
			},

			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"image_owner_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateImageOwners,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceChargeType,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceName,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceType,
			},

			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInternetChargeType,
			},

			"internet_max_bandwidth_in": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(1, 200),
			},

			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(0, 100),
			},

			"io_optimized": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIoOptimized,
			},

			"key_pair_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateKeyPairName,
			},

			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceNetworkType,
			},

			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ActiveSecurityEnhancementStrategy),
					string(DeactiveSecurityEnhancementStrategy),
				}),
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"spot_price_limit": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},

			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceSpotStrategy,
			},

			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
			},
			"system_disk_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskDescription,
			},
			"system_disk_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskName,
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(20, 500),
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"userdata": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"network_interfaces": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"primary_ip": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateIpAddress,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"vswitch_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunLaunchTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	args := ecs.CreateCreateLaunchTemplateRequest()
	args.LaunchTemplateName = d.Get("name").(string)
	args.Description = d.Get("description").(string)
	args.HostName = d.Get("host_name").(string)
	args.ImageId = d.Get("image_id").(string)
	args.ImageOwnerAlias = d.Get("image_owner_alias").(string)
	args.InstanceChargeType = d.Get("instance_charge_type").(string)
	args.InstanceName = d.Get("instance_name").(string)
	args.InstanceType = d.Get("instance_type").(string)
	args.AutoReleaseTime = d.Get("auto_release_time").(string)
	args.InternetChargeType = d.Get("internet_charge_type").(string)
	args.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
	args.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))
	args.IoOptimized = d.Get("io_optimized").(string)
	args.KeyPairName = d.Get("key_pair_name").(string)
	args.NetworkType = d.Get("network_type").(string)

	args.RamRoleName = d.Get("ram_role_name").(string)
	args.ResourceGroupId = d.Get("resource_group_id").(string)
	args.SecurityEnhancementStrategy = d.Get("security_enhancement_strategy").(string)
	args.SecurityGroupId = d.Get("security_group_id").(string)
	args.SpotPriceLimit = requests.NewFloat(d.Get("spot_price_limit").(float64))
	args.SpotStrategy = d.Get("spot_strategy").(string)
	args.SystemDiskDiskName = d.Get("system_disk_name").(string)
	args.SystemDiskCategory = d.Get("system_disk_category").(string)
	args.SystemDiskDescription = d.Get("system_disk_description").(string)
	args.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
	args.UserData = d.Get("userdata").(string)
	args.VSwitchId = d.Get("vswitch_id").(string)
	args.VpcId = d.Get("vpc_id").(string)
	args.ZoneId = d.Get("zone_id").(string)
	netsRaw := d.Get("network_interfaces").([]interface{})
	if netsRaw != nil {
		var nets []ecs.CreateLaunchTemplateNetworkInterface
		for _, raw := range netsRaw {
			netRaw := raw.(map[string]interface{})
			net := ecs.CreateLaunchTemplateNetworkInterface{
				NetworkInterfaceName: netRaw["name"].(string),
				VSwitchId:            netRaw["vswitch_id"].(string),
				SecurityGroupId:      netRaw["security_group_id"].(string),
				Description:          netRaw["description"].(string),
				PrimaryIpAddress:     netRaw["primary_ip"].(string),
			}
			nets = append(nets, net)
		}
		args.NetworkInterface = &nets
	}

	disksRaw := d.Get("data_disks").([]interface{})
	if disksRaw != nil {
		var disks []ecs.CreateLaunchTemplateDataDisk
		for _, raw := range disksRaw {
			diskRaw := raw.(map[string]interface{})
			disk := ecs.CreateLaunchTemplateDataDisk{
				Size:               fmt.Sprintf("%d", diskRaw["size"].(int)),
				SnapshotId:         diskRaw["snapshot_id"].(string),
				Category:           diskRaw["category"].(string),
				Encrypted:          fmt.Sprintf("%v", diskRaw["encrypted"].(bool)),
				DiskName:           diskRaw["name"].(string),
				Description:        diskRaw["description"].(string),
				DeleteWithInstance: fmt.Sprintf("%v", diskRaw["delete_with_instance"].(bool)),
			}
			disks = append(disks, disk)
		}

		args.DataDisk = &disks
	}
	tagsRaw := d.Get("tags").(map[string]interface{})
	var tags []ecs.CreateLaunchTemplateTag
	for key, value := range tagsRaw {
		tags = append(tags, ecs.CreateLaunchTemplateTag{
			Key:   key,
			Value: value.(string),
		})
	}
	args.Tag = &tags

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateLaunchTemplate(args)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_launch_template", args.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(args.GetActionName(), raw)
	resp, _ := raw.(*ecs.CreateLaunchTemplateResponse)

	d.SetId(resp.LaunchTemplateId)

	return resourceAliyunLaunchTemplateRead(d, meta)
}

func resourceAliyunLaunchTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	template, err := ecsService.DescribeLaunchTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", template.LaunchTemplateName)
	d.Set("description", template.LaunchTemplateData.Description)
	d.Set("host_name", template.LaunchTemplateData.HostName)
	d.Set("image_id", template.LaunchTemplateData.ImageId)
	d.Set("image_owner_alias", template.LaunchTemplateData.ImageOwnerAlias)
	d.Set("instance_charge_type", template.LaunchTemplateData.InstanceChargeType)
	d.Set("instance_name", template.LaunchTemplateData.InstanceName)
	d.Set("instance_type", template.LaunchTemplateData.InstanceType)
	d.Set("auto_release_time", template.LaunchTemplateData.AutoReleaseTime)
	d.Set("internet_charge_type", template.LaunchTemplateData.InternetChargeType)
	d.Set("internet_max_bandwidth_in", template.LaunchTemplateData.InternetMaxBandwidthIn)
	d.Set("internet_max_bandwidth_out", template.LaunchTemplateData.InternetMaxBandwidthOut)
	d.Set("io_optimized", template.LaunchTemplateData.IoOptimized)
	d.Set("key_pair_name", template.LaunchTemplateData.KeyPairName)
	d.Set("network_type", template.LaunchTemplateData.NetworkType)
	d.Set("ram_role_name", template.LaunchTemplateData.RamRoleName)
	d.Set("resource_group_id", template.LaunchTemplateData.ResourceGroupId)
	d.Set("security_enhancement_strategy", template.LaunchTemplateData.SecurityEnhancementStrategy)
	d.Set("security_group_id", template.LaunchTemplateData.SecurityGroupId)
	d.Set("spot_price_limit", template.LaunchTemplateData.SpotPriceLimit)
	d.Set("spot_strategy", template.LaunchTemplateData.SpotStrategy)
	d.Set("system_disk_name", template.LaunchTemplateData.SystemDiskDiskName)
	d.Set("system_disk_category", template.LaunchTemplateData.SystemDiskCategory)
	d.Set("system_disk_description", template.LaunchTemplateData.SystemDiskDescription)
	d.Set("system_disk_size", template.LaunchTemplateData.SystemDiskSize)
	d.Set("resource_group_id", template.LaunchTemplateData.ResourceGroupId)
	d.Set("userdata", template.LaunchTemplateData.UserData)
	d.Set("vswitch_id", template.LaunchTemplateData.VSwitchId)
	d.Set("vpc_id", template.LaunchTemplateData.VpcId)
	d.Set("zone_id", template.LaunchTemplateData.ZoneId)

	var interfaces []map[string]interface{}
	for _, net := range template.LaunchTemplateData.NetworkInterfaces.NetworkInterface {
		ds := make(map[string]interface{})
		ds["vswitch_id"] = net.VSwitchId
		ds["security_group_id"] = net.SecurityGroupId
		ds["name"] = net.NetworkInterfaceName
		ds["description"] = net.Description
		ds["primary_ip"] = net.PrimaryIpAddress
		interfaces = append(interfaces, ds)
	}
	if err := d.Set("network_interfaces", interfaces); err != nil {
		return WrapError(err)
	}

	var disks []map[string]interface{}
	for _, disk := range template.LaunchTemplateData.DataDisks.DataDisk {
		ds := make(map[string]interface{})
		ds["size"] = disk.Size
		ds["snapshot_id"] = disk.SnapshotId
		ds["category"] = disk.Category
		ds["encrypted"] = (disk.Encrypted == "true")
		ds["name"] = disk.DiskName
		ds["description"] = disk.Description
		ds["delete_with_instance"] = disk.DeleteWithInstance
		disks = append(disks, ds)
	}
	if err := d.Set("data_disks", disks); err != nil {
		return WrapError(err)
	}

	tags := make(map[string]interface{})
	for _, tag := range template.LaunchTemplateData.Tags.InstanceTag {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func resourceAliyunLaunchTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	req := ecs.CreateDeleteLaunchTemplateRequest()
	req.LaunchTemplateId = d.Id()

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteLaunchTemplate(req)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(req.GetAcceptFormat(), raw)

	return nil
}
