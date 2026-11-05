package kms

import (
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlicloudKmsPlaintext() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsPlaintextRead,

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_id": {
				Type:     schema.TypeString,
				Computed: true,
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
				Required: true,
			},
		},
	}
}

func dataSourceAlicloudKmsPlaintextRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// Since a plaintext has no ID, we create an ID based on
	// current unix time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	action := "Decrypt"
	request := make(map[string]interface{})
	request["CiphertextBlob"] = d.Get("ciphertext_blob")
	request["RegionId"] = client.RegionId

	if context := d.Get("encryption_context"); context != nil {
		cm := context.(map[string]interface{})
		contextJson, err := convert.ConvertMaptoJsonString(cm)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kms_plaintext", action, errs.AlibabaCloudSdkGoERROR)
		}
		request["EncryptionContext"] = contextJson
	}

	var response map[string]interface{}
	var err error
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kms_ciphertext", action, errs.AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.Plaintext", response)
	}
	d.Set("plaintext", resp)
	resp, err = jsonpath.Get("$.KeyId", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.KeyId", response)
	}
	d.Set("key_id", resp)

	return nil
}
