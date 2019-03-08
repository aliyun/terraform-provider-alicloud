package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateStringLengthInRange(0, 8192),
			},
			"key_usage": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !(value == "ENCRYPT/DECRYPT") {
						es = append(es, fmt.Errorf(
							"%q must be %s", k, "ENCRYPT/DECRYPT"))
					}
					return
				},
				Default: "ENCRYPT/DECRYPT",
			},
			"is_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion_window_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(7, 30),
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
	request.KeyUsage = d.Get("key_usage").(string)

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateKey(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "kms_key", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	resp, _ := raw.(*kms.CreateKeyResponse)
	d.SetId(resp.KeyMetadata.KeyId)

	return resourceAlicloudKmsKeyUpdate(d, meta)
}

func resourceAlicloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	kmsService := &KmsService{client: client}
	key, err := kmsService.DescribeKmsKey(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("description", key.KeyMetadata.Description)
	d.Set("key_usage", key.KeyMetadata.KeyUsage)
	d.Set("is_enabled", KeyState(key.KeyMetadata.KeyState) == Enabled)
	d.Set("deletion_window_in_days", d.Get("deletion_window_in_days").(int))
	d.Set("arn", key.KeyMetadata.Arn)

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
			request.KeyId = d.Id()
			raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
				return kmsClient.EnableKey(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
		}

		if !d.Get("is_enabled").(bool) && KeyState(key.KeyMetadata.KeyState) == Enabled {
			request := kms.CreateDisableKeyRequest()
			request.KeyId = d.Id()
			raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
				return kmsClient.DisableKey(request)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(request.GetActionName(), raw)
		}
		d.SetPartial("is_enabled")
	}

	d.Partial(false)

	return resourceAlicloudKmsKeyRead(d, meta)
}

func resourceAlicloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := kms.CreateScheduleKeyDeletionRequest()
	request.KeyId = d.Id()
	request.PendingWindowInDays = requests.NewInteger((d.Get("deletion_window_in_days").(int)))
	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.ScheduleKeyDeletion(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		kmsService := &KmsService{client: client}
		_, err := kmsService.DescribeKmsKey(d.Id())
		if err != nil {
			if NotFoundError(err) {
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request.GetActionName(), ProviderERROR))
	})
}
