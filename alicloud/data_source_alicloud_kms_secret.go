package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudKmsSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsSecretRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"context": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ciphertext_blob": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAlicloudKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// Since a secret has no ID, we create an ID based on
	// current unix time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	request := kms.CreateDecryptRequest()
	request.CiphertextBlob = d.Get("ciphertext_blob").(string)
	request.RegionId = client.RegionId

	if context := d.Get("context"); context != nil {
		request.EncryptionContext = context.(string)
	}

	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.Decrypt(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_secret", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.DecryptResponse)
	d.Set("plaintext", response.Plaintext)
	d.Set("key_id", response.KeyId)

	return nil
}
