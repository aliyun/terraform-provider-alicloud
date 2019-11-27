package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDisksRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "data"}, false),
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "cloud_essd", "cloud_ssd", "ephemeral_ssd", "cloud_efficiency"}, false),
			},
			"encrypted": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expiration_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsSchema(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDisksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ecs.CreateDescribeDisksRequest()
	request.RegionId = client.RegionId

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		request.DiskIds = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		request.DiskType = v.(string)
	}
	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		request.Category = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok && v.(string) != "" {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("encrypted"); ok && v.(string) != "" {
		if v == string(OnFlag) {
			request.Encrypted = requests.NewBoolean(true)
		} else {
			request.Encrypted = requests.NewBoolean(false)
		}
	}
	if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" {
		request.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeDisksTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeDisksTag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.Tag = &tags
	}

	var allDisks []ecs.Disk
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_disks", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.DescribeDisksResponse)

		if response == nil || len(response.Disks.Disk) < 1 {
			break
		}

		allDisks = append(allDisks, response.Disks.Disk...)

		if len(response.Disks.Disk) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(request.PageNumber)
		if err != nil {
			return WrapError(err)
		}
		request.PageNumber = page
	}

	var filteredDisksTemp []ecs.Disk

	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, disk := range allDisks {
			if r != nil && !r.MatchString(disk.DiskName) {
				continue
			}

			filteredDisksTemp = append(filteredDisksTemp, disk)
		}
	} else {
		filteredDisksTemp = allDisks
	}
	return disksDescriptionAttributes(d, filteredDisksTemp)
}

func disksDescriptionAttributes(d *schema.ResourceData, disks []ecs.Disk) error {
	var ids []string
	var s []map[string]interface{}
	for _, disk := range disks {
		mapping := map[string]interface{}{
			"id":                disk.DiskId,
			"name":              disk.DiskName,
			"resource_group_id": disk.ResourceGroupId,
			"description":       disk.Description,
			"region_id":         disk.RegionId,
			"availability_zone": disk.ZoneId,
			"status":            disk.Status,
			"type":              disk.Type,
			"category":          disk.Category,
			"encrypted":         string(OnFlag),
			"size":              disk.Size,
			"image_id":          disk.ImageId,
			"snapshot_id":       disk.SourceSnapshotId,
			"instance_id":       disk.InstanceId,
			"creation_time":     disk.CreationTime,
			"attached_time":     disk.AttachedTime,
			"detached_time":     disk.DetachedTime,
			"expiration_time":   disk.ExpiredTime,
			"tags":              tagsToMap(disk.Tags.Tag),
		}
		if !disk.Encrypted {
			mapping["encrypted"] = string(OffFlag)
		}

		ids = append(ids, disk.DiskId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("disks", s); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
