package cloudsso

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	errs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	helper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	logs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCloudSsoScimServerCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudSsoScimServerCredentialCreate,
		Read:   resourceAliCloudCloudSsoScimServerCredentialRead,
		Update: resourceAliCloudCloudSsoScimServerCredentialUpdate,
		Delete: resourceAliCloudCloudSsoScimServerCredentialDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"Enabled", "Disabled"}, false),
			},
			"credential_secret_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"credential_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"credential_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCloudSsoScimServerCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSCIMServerCredential"
	request := make(map[string]interface{})
	var err error

	request["DirectoryId"] = d.Get("directory_id")

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_sso_scim_server_credential", action, errs.AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.SCIMServerCredential", response); err != nil || resp == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_cloud_sso_scim_server_credential")
	} else {
		credentialId := resp.(map[string]interface{})["CredentialId"]
		credentialSecret := resp.(map[string]interface{})["CredentialSecret"]
		d.SetId(fmt.Sprintf("%v:%v", request["DirectoryId"], credentialId))

		credentialSecretData := map[string]interface{}{
			"CredentialId":     credentialId,
			"CredentialSecret": credentialSecret,
		}

		jsonData, err := json.Marshal(credentialSecretData)
		if err != nil {
			return errs.WrapError(err)
		}

		if output, ok := d.GetOk("credential_secret_file"); ok && output != nil {
			// create a private_key_data_file and write private key to it.
			helper.WriteToFile(output.(string), string(jsonData))
		}
	}

	return resourceAliCloudCloudSsoScimServerCredentialUpdate(d, meta)
}

func resourceAliCloudCloudSsoScimServerCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}

	object, err := cloudssoService.DescribeCloudSsoScimServerCredential(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_scim_server_credential cloudssoService.DescribeCloudSsoScimServerCredential Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("directory_id", object["DirectoryId"])
	d.Set("status", object["Status"])
	d.Set("credential_id", object["CredentialId"])
	d.Set("credential_type", object["CredentialType"])
	d.Set("create_time", object["CreateTime"])
	d.Set("expire_time", object["ExpireTime"])

	return nil
}

func resourceAliCloudCloudSsoScimServerCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DirectoryId":  parts[0],
		"CredentialId": parts[1],
	}

	if d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok {
		request["NewStatus"] = v
	}

	if update {
		action := "UpdateSCIMServerCredentialStatus"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCloudSsoScimServerCredentialRead(d, meta)
}

func resourceAliCloudCloudSsoScimServerCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSCIMServerCredential"
	var response map[string]interface{}

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"DirectoryId":  parts[0],
		"CredentialId": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
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
		if errs.IsExpectedErrors(err, []string{"EntityNotExists.SCIMCredential"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
