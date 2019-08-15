package alicloud

import (
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunLaunchTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunLaunchTemplateCreate,
		Read:   resourceAliyunLaunchTemplateRead,
		Update: resourceAliyunLaunchTemplateUpdate,
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
				ValidateFunc: validateLaunchTemplateDescription,
			},

			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"image_owner_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateImageOwners,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceChargeType,
			},

			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceName,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceType,
			},

			"auto_release_time": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetChargeType,
			},

			"internet_max_bandwidth_in": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(1, 200),
			},

			"internet_max_bandwidth_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 100),
			},

			"io_optimized": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIoOptimized,
			},

			"key_pair_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateKeyPairName,
			},

			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceNetworkType,
			},

			"ram_role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"security_enhancement_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ActiveSecurityEnhancementStrategy),
					string(DeactiveSecurityEnhancementStrategy),
				}),
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"spot_price_limit": {
				Type:     schema.TypeFloat,
				Optional: true,
			},

			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInstanceSpotStrategy,
			},

			"system_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskCategory,
			},
			"system_disk_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskDescription,
			},
			"system_disk_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskName,
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(20, 500),
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"userdata": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"network_interfaces": {
				Type:     schema.TypeList,
				Optional: true,
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

	request := ecs.CreateCreateLaunchTemplateRequest()
	request.RegionId = client.RegionId
	request.LaunchTemplateName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.HostName = d.Get("host_name").(string)
	request.ImageId = d.Get("image_id").(string)
	request.ImageOwnerAlias = d.Get("image_owner_alias").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	request.InstanceName = d.Get("instance_name").(string)
	request.InstanceType = d.Get("instance_type").(string)
	request.AutoReleaseTime = d.Get("auto_release_time").(string)
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
	request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))
	request.IoOptimized = d.Get("io_optimized").(string)
	request.KeyPairName = d.Get("key_pair_name").(string)
	request.NetworkType = d.Get("network_type").(string)

	request.RamRoleName = d.Get("ram_role_name").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.SecurityEnhancementStrategy = d.Get("security_enhancement_strategy").(string)
	request.SecurityGroupId = d.Get("security_group_id").(string)
	request.SpotPriceLimit = requests.NewFloat(d.Get("spot_price_limit").(float64))
	request.SpotStrategy = d.Get("spot_strategy").(string)
	request.SystemDiskDiskName = d.Get("system_disk_name").(string)
	request.SystemDiskCategory = d.Get("system_disk_category").(string)
	request.SystemDiskDescription = d.Get("system_disk_description").(string)
	request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
	request.UserData = d.Get("userdata").(string)
	request.VSwitchId = d.Get("vswitch_id").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.ZoneId = d.Get("zone_id").(string)
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
		request.NetworkInterface = &nets
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

		request.DataDisk = &disks
	}
	tagsRaw := d.Get("tags").(map[string]interface{})
	var tags []ecs.CreateLaunchTemplateTag
	for key, value := range tagsRaw {
		tags = append(tags, ecs.CreateLaunchTemplateTag{
			Key:   key,
			Value: value.(string),
		})
	}
	request.Tag = &tags

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateLaunchTemplate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_launch_template", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.CreateLaunchTemplateResponse)

	d.SetId(response.LaunchTemplateId)

	return resourceAliyunLaunchTemplateRead(d, meta)
}

func resourceAliyunLaunchTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeLaunchTemplate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	latestVersion, err := ecsService.DescribeLaunchTemplateVersion(d.Id(), int(object.LatestVersionNumber))
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", latestVersion.LaunchTemplateName)
	d.Set("description", latestVersion.LaunchTemplateData.Description)
	d.Set("host_name", latestVersion.LaunchTemplateData.HostName)
	d.Set("image_id", latestVersion.LaunchTemplateData.ImageId)
	d.Set("image_owner_alias", latestVersion.LaunchTemplateData.ImageOwnerAlias)
	d.Set("instance_charge_type", latestVersion.LaunchTemplateData.InstanceChargeType)
	d.Set("instance_name", latestVersion.LaunchTemplateData.InstanceName)
	d.Set("instance_type", latestVersion.LaunchTemplateData.InstanceType)
	d.Set("auto_release_time", latestVersion.LaunchTemplateData.AutoReleaseTime)
	d.Set("internet_charge_type", latestVersion.LaunchTemplateData.InternetChargeType)
	d.Set("internet_max_bandwidth_in", latestVersion.LaunchTemplateData.InternetMaxBandwidthIn)
	d.Set("internet_max_bandwidth_out", latestVersion.LaunchTemplateData.InternetMaxBandwidthOut)
	d.Set("io_optimized", latestVersion.LaunchTemplateData.IoOptimized)
	d.Set("key_pair_name", latestVersion.LaunchTemplateData.KeyPairName)
	d.Set("network_type", latestVersion.LaunchTemplateData.NetworkType)
	d.Set("ram_role_name", latestVersion.LaunchTemplateData.RamRoleName)
	d.Set("resource_group_id", latestVersion.LaunchTemplateData.ResourceGroupId)
	d.Set("security_enhancement_strategy", latestVersion.LaunchTemplateData.SecurityEnhancementStrategy)
	d.Set("security_group_id", latestVersion.LaunchTemplateData.SecurityGroupId)
	d.Set("spot_price_limit", latestVersion.LaunchTemplateData.SpotPriceLimit)
	d.Set("spot_strategy", latestVersion.LaunchTemplateData.SpotStrategy)
	d.Set("system_disk_name", latestVersion.LaunchTemplateData.SystemDiskDiskName)
	d.Set("system_disk_category", latestVersion.LaunchTemplateData.SystemDiskCategory)
	d.Set("system_disk_description", latestVersion.LaunchTemplateData.SystemDiskDescription)
	d.Set("system_disk_size", latestVersion.LaunchTemplateData.SystemDiskSize)
	d.Set("resource_group_id", latestVersion.LaunchTemplateData.ResourceGroupId)
	d.Set("userdata", latestVersion.LaunchTemplateData.UserData)
	d.Set("vswitch_id", latestVersion.LaunchTemplateData.VSwitchId)
	d.Set("vpc_id", latestVersion.LaunchTemplateData.VpcId)
	d.Set("zone_id", latestVersion.LaunchTemplateData.ZoneId)
	var interfaces []map[string]interface{}
	for _, net := range latestVersion.LaunchTemplateData.NetworkInterfaces.NetworkInterface {
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
	for _, disk := range latestVersion.LaunchTemplateData.DataDisks.DataDisk {
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
	for _, tag := range latestVersion.LaunchTemplateData.Tags.InstanceTag {
		tags[tag.Key] = tag.Value
	}
	d.Set("tags", tags)

	return nil
}

func resourceAliyunLaunchTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	versions, err := getLaunchTemplateVersions(d.Id(), meta)
	if err != nil {
		return WrapError(err)
	}
	// Remove one of the oldest and non-default version when the total number reach 30
	if len(versions) > 29 {
		var oldestVersion int64
		for _, version := range versions {
			if !version.DefaultVersion && (oldestVersion == 0 || version.VersionNumber < oldestVersion) {
				oldestVersion = version.VersionNumber
			}
		}

		err = deleteLaunchTemplateVersion(d.Id(), int(oldestVersion), meta)
		if err != nil {
			return WrapError(err)
		}
	}
	return WrapError(createLaunchTemplateVersion(d, meta))

}

func resourceAliyunLaunchTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateDeleteLaunchTemplateRequest()
	request.RegionId = client.RegionId
	request.LaunchTemplateId = d.Id()

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteLaunchTemplate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetAcceptFormat(), raw, request.RpcRequest, request)
	ecsService := EcsService{client}
	if err := ecsService.WaitForLaunchTemplate(d.Id(), Deleted, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunLaunchTemplateRead(d, meta)
}

func getLaunchTemplateVersions(id string, meta interface{}) ([]ecs.LaunchTemplateVersionSet, error) {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateDescribeLaunchTemplateVersionsRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(50)
	request.LaunchTemplateId = id
	raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
		return client.DescribeLaunchTemplateVersions(request)
	})
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.DescribeLaunchTemplateVersionsResponse)
	if len(response.LaunchTemplateVersionSets.LaunchTemplateVersionSet) > 0 {
		return response.LaunchTemplateVersionSets.LaunchTemplateVersionSet, nil
	} else {
		return nil, WrapErrorf(Error(GetNotFoundMessage("LaunchTemplate", id)), NotFoundMsg, ProviderERROR)
	}
}

func deleteLaunchTemplateVersion(id string, version int, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateDeleteLaunchTemplateVersionRequest()
	request.RegionId = client.RegionId
	request.LaunchTemplateId = id
	request.DeleteVersion = &[]string{strconv.FormatInt(int64(version), 10)}
	raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
		return client.DeleteLaunchTemplateVersion(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), ProviderERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}

func createLaunchTemplateVersion(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateCreateLaunchTemplateVersionRequest()
	request.RegionId = client.RegionId
	request.LaunchTemplateId = d.Id()
	request.Description = d.Get("description").(string)
	request.HostName = d.Get("host_name").(string)
	request.ImageId = d.Get("image_id").(string)
	request.ImageOwnerAlias = d.Get("image_owner_alias").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	request.InstanceName = d.Get("instance_name").(string)
	request.InstanceType = d.Get("instance_type").(string)
	request.AutoReleaseTime = d.Get("auto_release_time").(string)
	request.InternetChargeType = d.Get("internet_charge_type").(string)
	request.InternetMaxBandwidthIn = requests.NewInteger(d.Get("internet_max_bandwidth_in").(int))
	request.InternetMaxBandwidthOut = requests.NewInteger(d.Get("internet_max_bandwidth_out").(int))
	request.IoOptimized = d.Get("io_optimized").(string)
	request.KeyPairName = d.Get("key_pair_name").(string)
	request.NetworkType = d.Get("network_type").(string)

	request.RamRoleName = d.Get("ram_role_name").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)
	request.SecurityEnhancementStrategy = d.Get("security_enhancement_strategy").(string)
	request.SecurityGroupId = d.Get("security_group_id").(string)
	request.SpotPriceLimit = requests.NewFloat(d.Get("spot_price_limit").(float64))
	request.SpotStrategy = d.Get("spot_strategy").(string)
	request.SystemDiskDiskName = d.Get("system_disk_name").(string)
	request.SystemDiskCategory = d.Get("system_disk_category").(string)
	request.SystemDiskDescription = d.Get("system_disk_description").(string)
	request.SystemDiskSize = requests.NewInteger(d.Get("system_disk_size").(int))
	request.UserData = d.Get("userdata").(string)
	request.VSwitchId = d.Get("vswitch_id").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.ZoneId = d.Get("zone_id").(string)
	netsRaw := d.Get("network_interfaces").([]interface{})
	if netsRaw != nil {
		var nets []ecs.CreateLaunchTemplateVersionNetworkInterface
		for _, raw := range netsRaw {
			netRaw := raw.(map[string]interface{})
			net := ecs.CreateLaunchTemplateVersionNetworkInterface{
				NetworkInterfaceName: netRaw["name"].(string),
				VSwitchId:            netRaw["vswitch_id"].(string),
				SecurityGroupId:      netRaw["security_group_id"].(string),
				Description:          netRaw["description"].(string),
				PrimaryIpAddress:     netRaw["primary_ip"].(string),
			}
			nets = append(nets, net)
		}
		request.NetworkInterface = &nets
	}

	disksRaw := d.Get("data_disks").([]interface{})
	if disksRaw != nil {
		var disks []ecs.CreateLaunchTemplateVersionDataDisk
		for _, raw := range disksRaw {
			diskRaw := raw.(map[string]interface{})
			disk := ecs.CreateLaunchTemplateVersionDataDisk{
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

		request.DataDisk = &disks
	}
	tagsRaw := d.Get("tags").(map[string]interface{})
	var tags []ecs.CreateLaunchTemplateVersionTag
	for key, value := range tagsRaw {
		tags = append(tags, ecs.CreateLaunchTemplateVersionTag{
			Key:   key,
			Value: value.(string),
		})
	}
	request.Tag = &tags

	raw, err := client.WithEcsClient(func(client *ecs.Client) (interface{}, error) {
		return client.CreateLaunchTemplateVersion(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil

}
