package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudRamAccessKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccessKeyCreate,
		Read:   resourceAlicloudRamAccessKeyRead,
		Update: resourceAlicloudRamAccessKeyUpdate,
		Delete: resourceAlicloudRamAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      Active,
				ValidateFunc: validateRamAKStatus,
			},
		},
	}
}

func resourceAlicloudRamAccessKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateCreateAccessKeyRequest()
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.CreateAccessKey(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "ram_access_key", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.CreateAccessKeyResponse)

	// create a secret_file and write access key to it.
	if output, ok := d.GetOk("secret_file"); ok && output != nil {
		writeToFile(output.(string), response.AccessKey)
	}

	d.SetId(response.AccessKey.AccessKeyId)
	return resourceAlicloudRamAccessKeyUpdate(d, meta)
}

func resourceAlicloudRamAccessKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	request := ram.CreateUpdateAccessKeyRequest()
	request.UserAccessKeyId = d.Id()
	request.Status = d.Get("status").(string)

	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	if d.HasChange("status") {
		d.SetPartial("status")
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.UpdateAccessKey(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamAccessKeyRead(d, meta)
}

func resourceAlicloudRamAccessKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateListAccessKeysRequest()
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
	}

	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListAccessKeys(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*ram.ListAccessKeysResponse)
	accessKeys := response.AccessKeys.AccessKey
	if len(accessKeys) < 1 {
		return WrapError(Error("No access keys found."))
	}

	for _, v := range accessKeys {
		if v.AccessKeyId == d.Id() {
			d.Set("status", v.Status)
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceAlicloudRamAccessKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateDeleteAccessKeyRequest()
	request.UserAccessKeyId = d.Id()

	request1 := ram.CreateListAccessKeysRequest()

	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		request.UserName = v.(string)
		request1.UserName = v.(string)
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.DeleteAccessKey(request)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error deleting access key: %#v", err))
		}

		raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListAccessKeys(request1)
		})
		if err != nil {
			if RamEntityNotExist(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		response, _ := raw.(*ram.ListAccessKeysResponse)

		if len(response.AccessKeys.AccessKey) < 1 {
			return nil
		}

		return resource.RetryableError(WrapErrorf(err, DeleteTimeoutMsg, d.Id(), request1.GetActionName(), ProviderERROR))
	})
}
