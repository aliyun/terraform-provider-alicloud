package alicloud

import (
	"log"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				//must contain a valid status, expected Creating, Starting, Running, Stopping, Stopped
				ValidateFunc: validation.StringInSlice([]string{
					string(Running),
					string(Stopped),
					string(Creating),
					string(Starting),
					string(Stopping),
				}, false),
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
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
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spot_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_device_mappings": {
							Type:     schema.TypeList,
							Computed: true,
							//Set:      imageDiskDeviceMappingHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}
func dataSourceAlicloudInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateDescribeInstancesRequest()
	request.RegionId = client.RegionId
	request.Status = d.Get("status").(string)

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		request.InstanceIds = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.VpcId = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("availability_zone"); ok && v.(string) != "" {
		request.ZoneId = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeInstancesTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeInstancesTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	var allInstances []ecs.Instance
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instances", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeInstancesResponse)
		if len(response.Instances.Instance) < 1 {
			break
		}

		allInstances = append(allInstances, response.Instances.Instance...)

		if len(response.Instances.Instance) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredInstancesTemp []ecs.Instance

	nameRegex, ok := d.GetOk("name_regex")
	imageId, okImg := d.GetOk("image_id")
	if (ok && nameRegex.(string) != "") || (okImg && imageId.(string) != "") {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, inst := range allInstances {
			if r != nil && !r.MatchString(inst.InstanceName) {
				continue
			}
			if imageId.(string) != "" && inst.ImageId != imageId.(string) {
				continue
			}
			filteredInstancesTemp = append(filteredInstancesTemp, inst)
		}
	} else {
		filteredInstancesTemp = allInstances
	}

	return instancessDescriptionAttributes(d, filteredInstancesTemp, meta)
}

// populate the numerous fields that the instance description returns.
func instancessDescriptionAttributes(d *schema.ResourceData, instances []ecs.Instance, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, inst := range instances {
		mapping := map[string]interface{}{
			"id":                         inst.InstanceId,
			"region_id":                  inst.RegionId,
			"availability_zone":          inst.ZoneId,
			"status":                     inst.Status,
			"name":                       inst.InstanceName,
			"instance_type":              inst.InstanceType,
			"vpc_id":                     inst.VpcAttributes.VpcId,
			"vswitch_id":                 inst.VpcAttributes.VSwitchId,
			"image_id":                   inst.ImageId,
			"description":                inst.Description,
			"security_groups":            inst.SecurityGroupIds.SecurityGroupId,
			"resource_group_id":          inst.ResourceGroupId,
			"eip":                        inst.EipAddress.IpAddress,
			"key_name":                   inst.KeyPairName,
			"spot_strategy":              inst.SpotStrategy,
			"creation_time":              inst.CreationTime,
			"instance_charge_type":       inst.InstanceChargeType,
			"internet_charge_type":       inst.InternetChargeType,
			"internet_max_bandwidth_out": inst.InternetMaxBandwidthOut,
			// Complex types get their own functions
			"disk_device_mappings": instanceDisksMappings(d, inst.InstanceId, meta),
			"tags":                 tagsToMap(inst.Tags.Tag),
		}
		if len(inst.InnerIpAddress.IpAddress) > 0 {
			mapping["private_ip"] = inst.InnerIpAddress.IpAddress[0]
		} else {
			mapping["private_ip"] = inst.VpcAttributes.PrivateIpAddress.IpAddress[0]
		}
		if len(inst.PublicIpAddress.IpAddress) > 0 {
			mapping["public_ip"] = inst.PublicIpAddress.IpAddress[0]
		} else {
			mapping["public_ip"] = inst.VpcAttributes.NatIpAddress
		}

		ids = append(ids, inst.InstanceId)
		names = append(names, inst.InstanceName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	d.Set("ids", ids)
	d.Set("names", names)
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

//Returns a mapping of instance disks
func instanceDisksMappings(d *schema.ResourceData, instanceId string, meta interface{}) []map[string]interface{} {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateDescribeDisksRequest()
	request.InstanceId = instanceId

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(request)
	})

	if err != nil {
		log.Printf("[ERROR] DescribeDisks for instance got error: %#v", err)
		return nil
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ecs.DescribeDisksResponse)
	if len(response.Disks.Disk) < 1 {
		return nil
	}

	var s []map[string]interface{}

	for _, v := range response.Disks.Disk {
		mapping := map[string]interface{}{
			"device":   v.Device,
			"size":     v.Size,
			"category": v.Category,
			"type":     v.Type,
		}

		log.Printf("[DEBUG] alicloud_instances - adding disk device mapping: %v", mapping)
		s = append(s, mapping)
	}

	return s
}
