package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateNameRegex,
				ForceNew:     true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(DiskTypeSystem), string(DiskTypeData)}),
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(DiskCloud), string(DiskEphemeralSSD), string(DiskCloudEfficiency), string(DiskCloudSSD)}),
			},
			"encrypted": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(OnFlag), string(OffFlag)}),
			},
			"instance_id": {
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

	args := ecs.CreateDescribeDisksRequest()

	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		args.DiskIds = convertListToJsonString(v.([]interface{}))
	}
	if v, ok := d.GetOk("type"); ok && v.(string) != "" {
		args.DiskType = v.(string)
	}
	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		args.Category = v.(string)
	}
	if v, ok := d.GetOk("encrypted"); ok && v.(string) != "" {
		if v == string(OnFlag) {
			args.Encrypted = requests.NewBoolean(true)
		} else {
			args.Encrypted = requests.NewBoolean(false)
		}
	}
	if v, ok := d.GetOk("instance_id"); ok && v.(string) != "" {
		args.InstanceId = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeDisksTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeDisksTag{
				Key:   key,
				Value: value.(string),
			})
		}
		args.Tag = &tags
	}

	var allDisks []ecs.Disk
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeDisks(args)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "disks", args.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		resp, _ := raw.(*ecs.DescribeDisksResponse)

		if resp == nil || len(resp.Disks.Disk) < 1 {
			break
		}

		allDisks = append(allDisks, resp.Disks.Disk...)

		if len(resp.Disks.Disk) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return WrapError(err)
		} else {
			args.PageNumber = page
		}
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
