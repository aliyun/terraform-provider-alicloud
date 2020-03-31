package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/eci"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudEciopenapiImageCache() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEciopenapiImageCacheCreate,
		Read:   resourceAlicloudEciopenapiImageCacheRead,
		Delete: resourceAlicloudEciopenapiImageCacheDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"eip_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"image_cache_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_cache_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_cache_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"images": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEciopenapiImageCacheCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	//eciService := EciService{client}

	request := eci.CreateCreateImageCacheRequest()
	image := []string{}
	v := d.Get("images").([]interface{})
	for _, m := range v {
		image = append(image, m.(string))
	}
	request.Image = image

	if v, ok := d.GetOk("eip_instance_id"); ok {
		request.EipInstanceId = v.(string)
	}
	request.ImageCacheName = d.Get("image_cache_name").(string)
	if v, ok := d.GetOk("image_cache_size"); ok {
		request.ImageCacheSize = requests.NewInteger(v.(int))
	}
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	if v, ok := d.GetOk("retention_days"); ok {
		request.RetentionDays = requests.NewInteger(v.(int))
	}
	request.SecurityGroupId = d.Get("security_group_id").(string)
	request.VSwitchId = d.Get("vswitch_id").(string)
	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = v.(string)
	}
	raw, err := client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.CreateImageCache(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_eciopenapi_image_cache", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*eci.CreateImageCacheResponse)
	d.SetId(response.ImageCacheId)

	return resourceAlicloudEciopenapiImageCacheRead(d, meta)
}
func resourceAlicloudEciopenapiImageCacheRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	eciService := &EciService{client: client}
	object, err := eciService.DescribeEciopenapiImageCache(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	err = d.Set("images", object.Images)
	if err != nil {
		return err
	}
	d.Set("container_groupId", object.ContainerGroupId)
	d.Set("creation_time", object.CreationTime)
	d.Set("events", object.Events)
	d.Set("expire_date_time", object.ExpireDateTime)
	d.Set("image_cache_id", object.ImageCacheId)
	d.Set("image_cache_name", object.ImageCacheName)
	d.Set("progress", object.Progress)
	d.Set("snapshot_id", object.SnapshotId)
	d.Set("status", object.Status)
	return nil

}
func resourceAlicloudEciopenapiImageCacheDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := eci.CreateDeleteImageCacheRequest()
	request.ImageCacheId = d.Id()
	request.RegionId = client.RegionId
	raw, err := client.WithEciClient(func(eciClient *eci.Client) (interface{}, error) {
		return eciClient.DeleteImageCache(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
