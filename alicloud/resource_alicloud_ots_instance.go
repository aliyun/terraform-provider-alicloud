package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateStringLengthInRange(3, 16),
			},

			"accessed_by": {
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAliyunOtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	instanceType := d.Get("instance_type").(string)
	req := ots.CreateInsertInstanceRequest()
	req.ClusterType = convertInstanceType(OtsInstanceType(instanceType))
	types, err := otsService.DescribeOtsInstanceTypes()
	if err != nil {
		return err
	}
	valid := false
	for _, t := range types {
		if req.ClusterType == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("The instance type %s is not available in the region %s.", instanceType, client.RegionId)
	}
	req.InstanceName = d.Get("name").(string)
	req.Description = d.Get("description").(string)
	req.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))

	_, err = client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.InsertInstance(req)
	})
	if err != nil {
		return fmt.Errorf("failed to create instance with error: %s", err)
	}

	d.SetId(req.InstanceName)
	if err := otsService.WaitForOtsInstance(req.InstanceName, Running, DefaultTimeout); err != nil {
		return err
	}
	return resourceAliyunOtsInstanceUpdate(d, meta)
}

func resourceAliyunOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	inst, err := otsService.DescribeOtsInstance(d.Id())
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
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("accessed_by") {
		req := ots.CreateUpdateInstanceRequest()
		req.InstanceName = d.Id()
		req.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))
		_, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(req)
		})
		if err != nil {
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
			_, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.DeleteTags(args)
			})
			if err != nil {
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
			_, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.InsertTags(args)
			})
			if err != nil {
				return fmt.Errorf("Insertting tags got error: %s", err)
			}
		}
		d.SetPartial("tags")
	}
	if err := otsService.WaitForOtsInstance(d.Id(), Running, DefaultTimeout); err != nil {
		return err
	}
	d.Partial(false)
	return resourceAliyunOtsInstanceRead(d, meta)
}

func resourceAliyunOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	req := ots.CreateDeleteInstanceRequest()
	req.InstanceName = d.Id()
	return resource.Retry(10*time.Minute, func() *resource.RetryError {
		if _, err := otsService.DescribeOtsInstance(d.Id()); err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("When deleting instance, failed to describe instance with error: %s", err))
		}

		_, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.DeleteInstance(req)
		})
		if err != nil {
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
