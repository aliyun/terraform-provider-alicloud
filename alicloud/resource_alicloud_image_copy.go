package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudImageCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudImageCopyCreate,
		Read:   resourceAliCloudImageCopyRead,
		Update: resourceAliCloudImageCopyUpdate,
		Delete: resourceAliCloudImageCopyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"source_image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_region_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"name": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Attribute 'name' has been deprecated from version 1.69.0. Use `image_name` instead.",
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": tagsSchema(),
			// Not the public attribute and it used to automatically delete dependence snapshots while deleting the image.
			// Available in 1.136.0
			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}
func resourceAliCloudImageCopyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateCopyImageRequest()
	request.RegionId = client.RegionId
	request.ImageId = d.Get("source_image_id").(string)
	if v, ok := d.GetOk("destination_region_id"); ok && v.(string) != "" {
		request.DestinationRegionId = v.(string)
	}
	if v, ok := d.GetOk("image_name"); ok && v.(string) != "" {
		request.DestinationImageName = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.DestinationDescription = v.(string)
	}
	if v, ok := d.GetOk("kms_key_id"); ok && v.(string) != "" {
		request.KMSKeyId = v.(string)
	}
	if v, ok := d.GetOk("encrypted"); ok {
		request.Encrypted = requests.NewBoolean(v.(bool))
	}
	tags := d.Get("tags").(map[string]interface{})
	if tags != nil && len(tags) > 0 {
		imageTags := make([]ecs.CopyImageTag, 0, len(tags))
		for k, v := range tags {
			imageTag := ecs.CopyImageTag{
				Key:   k,
				Value: v.(string),
			}
			imageTags = append(imageTags, imageTag)
		}
		request.Tag = &imageTags
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CopyImage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_copy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ecs.CopyImageResponse)
	d.SetId(response.ImageId)
	stateConf := BuildStateConf([]string{"Creating", "Waiting"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFunc(d.Id(), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		// If the copying is timeout, the progress should be cancelled
		// Currently the product does not support cancel some specify images
		//cancelCopyImageRequest := ecs.CreateCancelCopyImageRequest()
		//cancelCopyImageRequest.ImageId = d.Id()
		//cancelCopyImageRequest.RegionId = client.RegionId
		//if _, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		//	return ecsClient.CancelCopyImage(cancelCopyImageRequest)
		//}); err != nil {
		//	return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_copy", cancelCopyImageRequest.GetActionName(), AlibabaCloudSdkGoERROR)
		//}
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudImageCopyRead(d, meta)
}
func resourceAliCloudImageCopyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return WrapError(err)
	}
	return resourceAliCloudImageRead(d, meta)
}
func resourceAliCloudImageCopyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ecsService := EcsService{client}
	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.ImageName)
	d.Set("image_name", object.ImageName)
	d.Set("description", object.Description)

	tags := object.Tags.Tag
	if len(tags) > 0 {
		err = d.Set("tags", ecsService.tagsToMap(tags))
	}
	return WrapError(err)
}

func resourceAliCloudImageCopyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	return ecsService.deleteImage(d)
}
