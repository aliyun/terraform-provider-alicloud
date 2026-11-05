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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlicloudKmsCiphertext() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsCiphertextCreate,
		Read:   schema.Noop,
		Delete: resourceAlicloudKmsCiphertextDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("plaintext"), cty.GetAttrPath("plaintext_wo")),
		},

		Schema: map[string]*schema.Schema{
			"plaintext": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"plaintext", "plaintext_wo"},
			},
			"plaintext_wo": {
				Type:         schema.TypeString,
				Optional:     true,
				WriteOnly:    true,
				ExactlyOneOf: []string{"plaintext", "plaintext_wo"},
				RequiredWith: []string{"plaintext_wo_version"},
			},
			"plaintext_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"plaintext_wo"},
			},
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
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

func resourceAlicloudKmsCiphertextCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// Since a ciphertext has no ID, we create an ID based on
	// current unix time.
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	action := "Encrypt"
	request := make(map[string]interface{})
	request["Plaintext"] = d.Get("plaintext")
	if woValue, err := helper.GetWriteOnlyValue(d, cty.GetAttrPath("plaintext_wo"), cty.String); err != nil {
		return errs.WrapError(err)
	} else if !woValue.IsNull() {
		request["Plaintext"] = woValue.AsString()
	}
	request["KeyId"] = d.Get("key_id")
	request["RegionId"] = client.RegionId

	if context := d.Get("encryption_context"); context != nil {
		cm := context.(map[string]interface{})
		contextJson, err := convert.ConvertMaptoJsonString(cm)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_kms_ciphertext", action, errs.AlibabaCloudSdkGoERROR)
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
	resp, err := jsonpath.Get("$.CiphertextBlob", response)
	if err != nil {
		return errs.WrapErrorf(err, errs.FailedGetAttributeMsg, action, "$.CiphertextBlob", response)
	}
	d.Set("ciphertext_blob", resp)

	return nil
}

func resourceAlicloudKmsCiphertextDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
