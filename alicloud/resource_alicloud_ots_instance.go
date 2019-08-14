package alicloud

import (
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
	request := ots.CreateInsertInstanceRequest()
	request.ClusterType = convertInstanceType(OtsInstanceType(instanceType))
	types, err := otsService.DescribeOtsInstanceTypes()
	if err != nil {
		return WrapError(err)
	}
	valid := false
	for _, t := range types {
		if request.ClusterType == t {
			valid = true
			break
		}
	}
	if !valid {
		return WrapError(Error("The instance type %s is not available in the region %s.", instanceType, client.RegionId))
	}
	request.InstanceName = d.Get("name").(string)
	request.Description = d.Get("description").(string)
	request.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))

	raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.InsertInstance(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	d.SetId(request.InstanceName)
	if err := otsService.WaitForOtsInstance(request.InstanceName, Running, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunOtsInstanceUpdate(d, meta)
}

func resourceAliyunOtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.InstanceName)
	d.Set("accessed_by", convertInstanceAccessedByRevert(object.Network))
	d.Set("instance_type", convertInstanceTypeRevert(object.ClusterType))
	d.Set("description", object.Description)
	d.Set("tags", otsTagsToMap(object.TagInfos.TagInfo))
	return nil
}

func resourceAliyunOtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("accessed_by") {
		request := ots.CreateUpdateInstanceRequest()
		request.InstanceName = d.Id()
		request.Network = convertInstanceAccessedBy(InstanceAccessedByType(d.Get("accessed_by").(string)))
		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.UpdateInstance(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw)
		d.SetPartial("accessed_by")
	}

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		if len(remove) > 0 {
			request := ots.CreateDeleteTagsRequest()
			request.InstanceName = d.Id()
			var tags []ots.DeleteTagsTagInfo
			for _, t := range remove {
				tags = append(tags, ots.DeleteTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			request.TagInfo = &tags
			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.DeleteTags(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
		}

		if len(create) > 0 {
			request := ots.CreateInsertTagsRequest()
			request.InstanceName = d.Id()
			var tags []ots.InsertTagsTagInfo
			for _, t := range create {
				tags = append(tags, ots.InsertTagsTagInfo{
					TagKey:   t.Key,
					TagValue: t.Value,
				})
			}
			request.TagInfo = &tags
			raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
				return otsClient.InsertTags(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
		}
		d.SetPartial("tags")
	}
	if err := otsService.WaitForOtsInstance(d.Id(), Running, DefaultTimeout); err != nil {
		return WrapError(err)
	}
	d.Partial(false)
	return resourceAliyunOtsInstanceRead(d, meta)
}

func resourceAliyunOtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	request := ots.CreateDeleteInstanceRequest()
	request.InstanceName = d.Id()
	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
			return otsClient.DeleteInstance(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{"AuthFailed", "InvalidStatus", "ValidationFailed"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(otsService.WaitForOtsInstance(d.Id(), Deleted, DefaultLongTimeout))
}
