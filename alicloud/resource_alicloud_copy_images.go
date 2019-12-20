package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAliCloudCopyImage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCopyImageCreate,
		Read:   resourceAliCloudCopyImageRead,
		Update: resourceAliCloudCopyImageUpdate,
		Delete: resourceAliCloudCopyImageDelete,
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
			"source_region_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}
func resourceAliCloudCopyImageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateCopyImageRequest()
	request.RegionId = d.Get("source_region_id").(string)
	request.ImageId = d.Get("source_image_id").(string)
	request.DestinationRegionId = client.RegionId
	request.DestinationImageName = d.Get("name").(string)
	request.DestinationDescription = d.Get("description").(string)
	request.Encrypted = requests.NewBoolean(d.Get("encrypted").(bool))
	request.KMSKeyId = d.Get("kms_key_id").(string)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_copy_image", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response, _ := raw.(*ecs.CopyImageResponse)
	d.SetId(response.ImageId)
	stateConf := BuildStateConf([]string{"Creating"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, ecsService.ImageStateRefreshFunc(d.Id(), []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudCopyImageRead(d, meta)
}
func resourceAliCloudCopyImageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	err := ecsService.updateImage(d)
	if err != nil {
		return WrapError(err)
	}
	return resourceAliCloudImageRead(d, meta)
}
func resourceAliCloudCopyImageRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("description", object.Description)

	tags := object.Tags.Tag
	if len(tags) > 0 {
		err = d.Set("tags", tagsToMap(tags))
	}
	return WrapError(err)
}

func resourceAliCloudCopyImageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	return ecsService.deleteImage(d)
}
