package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudOtsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsInstanceCreate,
		Read:   resourceAliyunOtsInstanceRead,
		Update: resourceAliyunOtsInstanceUpdate,
		Delete: resourceAliyunOtsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(3, 16),
			},

			"accessed_by": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  AnyNetwork,
				ValidateFunc: validateAllowedStringValue([]string{
					string(AnyNetwork), string(VpcOnly), string(VpcOrConsole),
				}),
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  OtsHighPerformance,
				ValidateFunc: validateAllowedStringValue([]string{
					string(OtsCapacity), string(OtsHighPerformance),
				}),
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return !d.IsNewResource()
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunOtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	req := ots.CreateInsertInstanceRequest()
	req.ClusterType = convertInstanceType(OtsInstanceType(d.Get("instance_type").(string)))
	req.InstanceName = d.Get("name").(string)
	req.Description = d.Get("description").(string)
	req.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))

	if _, err := client.otsconn.InsertInstance(req); err != nil {
		return fmt.Errorf("failed to create instance with error: %s", err)
	}

	d.SetId(req.InstanceName)
	if err := client.WaitForOtsInstance(req.InstanceName, Running, DefaultTimeout); err != nil {
		return err
	}
	return resourceAliyunOtsInstanceUpdate(d, meta)
}

func resourceAliyunOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	inst, err := meta.(*AliyunClient).DescribeOtsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to describe instance with error: %s", err)
	}

	d.Set("name", inst.InstanceName)
	d.Set("accessed_by", convertInstanceAccessedByRevert(inst.Network))
	d.Set("instance_type", convertInstanceTypeRevert(inst.ClusterType))
	d.Set("description", inst.Description)
	d.Set("tags", otsTagsToMap(inst.TagInfos.TagInfo))
	return nil
}

func resourceAliyunOtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("accessed_by") {
		req := ots.CreateUpdateInstanceRequest()
		req.InstanceName = d.Id()
		req.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))
		if _, err := client.otsconn.UpdateInstance(req); err != nil {
			return fmt.Errorf("UpdateInstance %s got an error: %#v.", d.Id(), err)
		}
		d.SetPartial("accessed_by")
	}

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		if len(remove) > 0 {
			args := ots.CreateDeleteTagsRequest()
			args.InstanceName = d.Id()
			var tags []ots.DeleteTagsTagInfo
			for _, t := range remove {
				tags = append(tags, ots.DeleteTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			args.TagInfo = &tags
			if _, err := client.otsconn.DeleteTags(args); err != nil {
				return fmt.Errorf("Remove tags got error: %s", err)
			}
		}

		if len(create) > 0 {
			args := ots.CreateInsertTagsRequest()
			args.InstanceName = d.Id()
			var tags []ots.InsertTagsTagInfo
			for _, t := range create {
				tags = append(tags, ots.InsertTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			args.TagInfo = &tags
			if _, err := client.otsconn.InsertTags(args); err != nil {
				return fmt.Errorf("Insertting tags got error: %s", err)
			}
		}
		d.SetPartial("tags")
	}
	if err := client.WaitForOtsInstance(d.Id(), Running, DefaultTimeout); err != nil {
		return err
	}
	d.Partial(false)
	return resourceAliyunOtsInstanceRead(d, meta)
}

func resourceAliyunOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	req := ots.CreateDeleteInstanceRequest()
	req.InstanceName = d.Id()
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		if _, err := meta.(*AliyunClient).DescribeOtsInstance(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("When deleting instance, failed to describe instance with error: %s", err))
		}

		if _, err := client.otsconn.DeleteInstance(req); err != nil {
			if NotFoundError(err) {
				return nil
			}
			if IsExceptedErrors(err, []string{"AuthFailed", "InvalidStatus", "ValidationFailed"}) {
				return resource.RetryableError(fmt.Errorf("Deleting instance %s timeout and got an error: %#v.", d.Id(), err))
			}
			return resource.NonRetryableError(err)
		}

		return resource.RetryableError(fmt.Errorf("Deleting instance %s timeout.", d.Id()))
	})
}
