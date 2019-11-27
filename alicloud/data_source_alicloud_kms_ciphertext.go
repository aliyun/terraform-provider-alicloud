package alicloud

import (
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func dataSourceAlicloudKmsCiphertext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsCiphertextRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ciphertext_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudKmsCiphertextRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// Since a ciphertext has no ID, we create an ID based on
	// current unix time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	request := kms.CreateEncryptRequest()
	request.Plaintext = d.Get("plaintext").(string)
	request.KeyId = d.Get("key_id").(string)
	request.RegionId = client.RegionId

	if context := d.Get("encryption_context"); context != nil {
		cm := context.(map[string]interface{})
		contextJson, err := convertMaptoJsonString(cm)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_ciphertext", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		request.EncryptionContext = string(contextJson)
	}

	raw, err := client.WithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.Encrypt(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_ciphertext", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*kms.EncryptResponse)
	d.Set("ciphertext_blob", response.CiphertextBlob)

	return nil
}
