package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/config"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudConfigDeliveryChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudConfigDeliveryChannelCreate,
		Read:   resourceAlicloudConfigDeliveryChannelRead,
		Update: resourceAlicloudConfigDeliveryChannelUpdate,
		Delete: resourceAlicloudConfigDeliveryChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"delivery_channel_assume_role_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delivery_channel_condition": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delivery_channel_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"delivery_channel_target_arn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delivery_channel_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MNS", "OSS", "SLS"}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func resourceAlicloudConfigDeliveryChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := config.CreatePutDeliveryChannelRequest()
	request.DeliveryChannelAssumeRoleArn = d.Get("delivery_channel_assume_role_arn").(string)
	if v, ok := d.GetOk("delivery_channel_condition"); ok {
		request.DeliveryChannelCondition = v.(string)
	}

	if v, ok := d.GetOk("delivery_channel_name"); ok {
		request.DeliveryChannelName = v.(string)
	}

	request.DeliveryChannelTargetArn = d.Get("delivery_channel_target_arn").(string)
	request.DeliveryChannelType = d.Get("delivery_channel_type").(string)
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		request.Status = requests.NewInteger(v.(int))
	}

	raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
		return configClient.PutDeliveryChannel(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_config_delivery_channel", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*config.PutDeliveryChannelResponse)
	d.SetId(fmt.Sprintf("%v", response.DeliveryChannelId))

	return resourceAlicloudConfigDeliveryChannelRead(d, meta)
}
func resourceAlicloudConfigDeliveryChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	configService := ConfigService{client}
	object, err := configService.DescribeConfigDeliveryChannel(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_config_delivery_channel configService.DescribeConfigDeliveryChannel Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("delivery_channel_assume_role_arn", object.DeliveryChannelAssumeRoleArn)
	d.Set("delivery_channel_condition", object.DeliveryChannelCondition)
	d.Set("delivery_channel_name", object.DeliveryChannelName)
	d.Set("delivery_channel_target_arn", object.DeliveryChannelTargetArn)
	d.Set("delivery_channel_type", object.DeliveryChannelType)
	d.Set("description", object.Description)
	d.Set("status", object.Status)
	return nil
}
func resourceAlicloudConfigDeliveryChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := config.CreatePutDeliveryChannelRequest()
	request.DeliveryChannelId = d.Id()
	if d.HasChange("delivery_channel_assume_role_arn") {
		update = true
	}
	request.DeliveryChannelAssumeRoleArn = d.Get("delivery_channel_assume_role_arn").(string)
	if d.HasChange("delivery_channel_target_arn") {
		update = true
	}
	request.DeliveryChannelTargetArn = d.Get("delivery_channel_target_arn").(string)
	request.DeliveryChannelType = d.Get("delivery_channel_type").(string)
	if d.HasChange("delivery_channel_condition") {
		update = true
		request.DeliveryChannelCondition = d.Get("delivery_channel_condition").(string)
	}
	if d.HasChange("delivery_channel_name") {
		update = true
		request.DeliveryChannelName = d.Get("delivery_channel_name").(string)
	}
	if d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if d.HasChange("status") {
		update = true
		request.Status = requests.NewInteger(d.Get("status").(int))
	}
	if update {
		raw, err := client.WithConfigClient(func(configClient *config.Client) (interface{}, error) {
			return configClient.PutDeliveryChannel(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudConfigDeliveryChannelRead(d, meta)
}
func resourceAlicloudConfigDeliveryChannelDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudConfigDeliveryChannel. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
