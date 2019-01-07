package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliyunDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDiskCreate,
		Read:   resourceAliyunDiskRead,
		Update: resourceAliyunDiskUpdate,
		Delete: resourceAliyunDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskName,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDiskDescription,
			},

			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
				Default:      DiskCloudEfficiency,
			},

			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	availabilityZone, err := ecsService.DescribeZone(d.Get("availability_zone").(string))
	if err != nil {
		return err
	}

	args := ecs.CreateCreateDiskRequest()
	args.ZoneId = availabilityZone.ZoneId

	if v, ok := d.GetOk("category"); ok && v.(string) != "" {
		category := DiskCategory(v.(string))
		if err := ecsService.DiskAvailable(availabilityZone, category); err != nil {
			return err
		}
		args.DiskCategory = v.(string)
	}

	var size int
	if v, ok := d.GetOk("size"); ok {
		size = v.(int)
		if args.DiskCategory == string(DiskCloud) && (size < 5 || size > 2000) {
			return fmt.Errorf("the size of cloud disk must between 5 to 2000")
		}

		if (args.DiskCategory == string(DiskCloudEfficiency) || args.DiskCategory == string(DiskCloudSSD)) &&
			(size < 20 || size > 32768) {
			return fmt.Errorf("the size of %s disk must between 20 to 32768", args.DiskCategory)
		}
		args.Size = requests.NewInteger(size)

		d.Set("size", args.Size)
	}

	if v, ok := d.GetOk("snapshot_id"); ok && v.(string) != "" {
		args.SnapshotId = v.(string)
	}

	if size <= 0 && args.SnapshotId == "" {
		return fmt.Errorf("One of size or snapshot_id is required when specifying an ECS disk.")
	}

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		args.DiskName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		args.Description = v.(string)
	}

	if v, ok := d.GetOk("encrypted"); ok {
		args.Encrypted = requests.NewBoolean(v.(bool))
	}
	args.ClientToken = buildClientToken("TF-CreateDisk")
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateDisk(args)
	})
	if err != nil {
		return fmt.Errorf("CreateDisk got a error: %#v", err)
	}
	resp, _ := raw.(*ecs.CreateDiskResponse)
	if resp == nil {
		return fmt.Errorf("CreateDisk got a nil response: %#v", resp)
	}

	d.SetId(resp.DiskId)

	if err := ecsService.WaitForEcsDisk(d.Id(), Available, DefaultTimeout); err != nil {
		return fmt.Errorf("Waitting for disk %s got an error: %#v.", Available, err)
	}

	return resourceAliyunDiskUpdate(d, meta)
}

func resourceAliyunDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	disk, err := ecsService.DescribeDiskById("", d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error DescribeDisk: %#v", err)
	}

	d.Set("availability_zone", disk.ZoneId)
	d.Set("category", disk.Category)
	d.Set("size", disk.Size)
	d.Set("status", disk.Status)
	d.Set("name", disk.DiskName)
	d.Set("description", disk.Description)
	d.Set("snapshot_id", disk.SourceSnapshotId)
	d.Set("encrypted", disk.Encrypted)

	tags, err := ecsService.DescribeTags(d.Id(), TagResourceDisk)
	if err != nil && !NotFoundError(err) {
		return fmt.Errorf("[ERROR] DescribeTags for disk got error: %#v", err)
	}
	if len(tags) > 0 {
		d.Set("tags", tagsToMap(tags))
	}

	return nil
}

func resourceAliyunDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if err := setTags(client, TagResourceDisk, d); err != nil {
		return fmt.Errorf("Set tags for instance got error: %#v", err)
	} else {
		d.SetPartial("tags")
	}

	if d.HasChange("size") && !d.IsNewResource() {
		size := d.Get("size").(int)
		args := ecs.CreateResizeDiskRequest()
		args.DiskId = d.Id()
		args.NewSize = requests.NewInteger(size)
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ResizeDisk(args)
		})
		if err != nil {
			return fmt.Errorf("Resize disk failed: %#v", err)
		}
		d.SetPartial("size")
	}

	attributeUpdate := false
	args := ecs.CreateModifyDiskAttributeRequest()
	args.DiskId = d.Id()

	if d.HasChange("name") {
		d.SetPartial("name")
		val := d.Get("name").(string)
		args.DiskName = val

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		val := d.Get("description").(string)
		args.Description = val

		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyDiskAttribute(args)
		})
		if err != nil {
			return err
		}
	}

	d.Partial(false)

	return resourceAliyunDiskRead(d, meta)
}

func resourceAliyunDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	req := ecs.CreateDeleteDiskRequest()
	req.DiskId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteDisk(req)
		})
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			if IsExceptedErrors(err, DiskInvalidOperation) {
				return resource.RetryableError(fmt.Errorf("Deleting Disk %s timeout and got an error: %#v.", d.Id(), err))
			}
			return resource.NonRetryableError(err)
		}

		disk, descErr := ecsService.DescribeDiskById("", d.Id())

		if descErr != nil {
			if NotFoundError(descErr) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("While deleting disk %s, describing disk got an error: %#v.", d.Id(), descErr))
		}
		if disk.DiskId == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting Disk %s timeout.", d.Id()))
	})
}
