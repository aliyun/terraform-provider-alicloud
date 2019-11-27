package alicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsKeyCreate,
		Read:   resourceAlicloudKmsKeyRead,
		Update: resourceAlicloudKmsKeyUpdate,
		Delete: resourceAlicloudKmsKeyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "From Terraform",
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 8192),
			},
			"key_usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENCRYPT/DECRYPT"}, false),
				Default:      "ENCRYPT/DECRYPT",
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion_window_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(7, 30),
				Default:      30,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := kms.CreateCreateKeyRequest()
	request.RegionId = client.RegionId
	request.KeyUsage = d.Get("key_usage").(string)

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateKey(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_key", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.CreateKeyResponse)
	d.SetId(response.KeyMetadata.KeyId)

	return resourceAlicloudKmsKeyUpdate(d, meta)
}

func resourceAlicloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	kmsService := &KmsService{client: client}
	object, err := kmsService.DescribeKmsKey(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", object.KeyMetadata.Description)
	d.Set("key_usage", object.KeyMetadata.KeyUsage)
	d.Set("is_enabled", KeyState(object.KeyMetadata.KeyState) == Enabled)
	d.Set("deletion_window_in_days", d.Get("deletion_window_in_days").(int))
	d.Set("arn", object.KeyMetadata.Arn)

	return nil
}

func resourceAlicloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	if d.HasChange("is_enabled") {
		kmsService := &KmsService{client: client}
		key, err := kmsService.DescribeKmsKey(d.Id())
		if err != nil {
			return WrapError(err)
		}
		if d.Get("is_enabled").(bool) && KeyState(key.KeyMetadata.KeyState) == Disabled {
			request := kms.CreateEnableKeyRequest()
			request.RegionId = client.RegionId
			request.KeyId = d.Id()
			raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
				return kmsClient.EnableKey(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}

		if !d.Get("is_enabled").(bool) && KeyState(key.KeyMetadata.KeyState) == Enabled {
			request := kms.CreateDisableKeyRequest()
			request.RegionId = client.RegionId
			request.KeyId = d.Id()
			raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
				return kmsClient.DisableKey(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		}
		d.SetPartial("is_enabled")
	}

	d.Partial(false)

	return resourceAlicloudKmsKeyRead(d, meta)
}

func resourceAlicloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kmsService := KmsService{client}
	request := kms.CreateScheduleKeyDeletionRequest()
	request.RegionId = client.RegionId
	request.KeyId = d.Id()
	request.PendingWindowInDays = requests.NewInteger((d.Get("deletion_window_in_days").(int)))
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.ScheduleKeyDeletion(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(kmsService.WaitForKmsKey(d.Id(), Deleted, DefaultTimeoutMedium))
}
