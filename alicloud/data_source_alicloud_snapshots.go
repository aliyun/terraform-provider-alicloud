package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encrypted": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MinItems: 1,
				MaxItems: 100,
			},
			"name_regex": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateNameRegex,
			},
			"status": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"progressing", "accomplished", "failed", "all"}),
				Default:      "all",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"auto", "user", "all"}),
				Default:      "all",
			},
			"source_disk_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{"System", "Data"}),
			},
			"usage": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{"image", "disk", "image_disk", "node"}),
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"output_file": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"snapshots": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"progress": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_size": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_days": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"remain_time": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_time": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"usage": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	args := ecs.CreateDescribeSnapshotsRequest()

	if instanceId, ok := d.GetOk("instance_id"); ok {
		args.InstanceId = instanceId.(string)
	}
	if diskId, ok := d.GetOk("disk_id"); ok {
		args.DiskId = diskId.(string)
	}
	if encrypted, ok := d.GetOk("encrypted"); ok {
		args.Encrypted = requests.NewBoolean(encrypted.(bool))
	}
	if ids, ok := d.GetOk("ids"); ok {
		args.SnapshotIds = convertListToJsonString(ids.(*schema.Set).List())
	}
	if status, ok := d.GetOk("status"); ok {
		args.Status = status.(string)
	}
	if typ, ok := d.GetOk("type"); ok {
		args.SnapshotType = typ.(string)
	}

	if diskType, ok := d.GetOk("source_disk_type"); ok {
		args.SourceDiskType = diskType.(string)
	}
	if usage, ok := d.GetOk("usage"); ok {
		args.Usage = usage.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []ecs.DescribeSnapshotsTag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, ecs.DescribeSnapshotsTag{
				Key:   key,
				Value: value.(string),
			})
		}
		args.Tag = &tags
	}

	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	var allSnapshots []ecs.Snapshot
	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSnapshots(args)
		})
		if err != nil {
			return WrapError(err)
		}
		resp := raw.(*ecs.DescribeSnapshotsResponse)
		allSnapshots = append(allSnapshots, resp.Snapshots.Snapshot...)

		if len(resp.Snapshots.Snapshot) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(args.PageNumber); err != nil {
			return WrapError(err)
		} else {
			args.PageNumber = page
		}
	}

	var filteredSnapshots []ecs.Snapshot
	nameRegex, ok := d.GetOk("name_regex")
	if ok && nameRegex.(string) != "" {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, snapshot := range allSnapshots {
			if r != nil && !r.MatchString(snapshot.SnapshotName) {
				continue
			}

			filteredSnapshots = append(filteredSnapshots, snapshot)
		}
	} else {
		filteredSnapshots = allSnapshots
	}

	return snapshotsDescriptionAttributes(d, filteredSnapshots)
}

func snapshotsDescriptionAttributes(d *schema.ResourceData, snapshots []ecs.Snapshot) error {
	var s []map[string]interface{}
	var ids []string
	var names []string
	for _, snapshot := range snapshots {
		mapping := map[string]interface{}{
			"id":               snapshot.SnapshotId,
			"name":             snapshot.SnapshotName,
			"description":      snapshot.Description,
			"encrypted":        snapshot.Encrypted,
			"progress":         snapshot.Progress,
			"source_disk_id":   snapshot.SourceDiskId,
			"source_disk_type": snapshot.SourceDiskType,
			"source_disk_size": snapshot.SourceDiskSize,
			"product_code":     snapshot.ProductCode,
			"retention_days":   snapshot.RetentionDays,
			"remain_time":      snapshot.RemainTime,
			"creation_time":    snapshot.CreationTime,
			"status":           snapshot.Status,
			"usage":            snapshot.Usage,
		}
		s = append(s, mapping)
		ids = append(ids, snapshot.SnapshotId)
		names = append(names, snapshot.SnapshotName)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("snapshots", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
