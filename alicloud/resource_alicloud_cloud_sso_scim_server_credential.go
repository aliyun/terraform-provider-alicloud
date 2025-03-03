package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				ValidateFunc: StringInSlice([]string{"Enabled", "Disabled"}, false),
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

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_sso_scim_server_credential", action, AlibabaCloudSdkGoERROR)
	}

	if resp, err := jsonpath.Get("$.SCIMServerCredential", response); err != nil || resp == nil {
		return WrapErrorf(err, IdMsg, "alicloud_cloud_sso_scim_server_credential")
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
			return WrapError(err)
		}

		if output, ok := d.GetOk("credential_secret_file"); ok && output != nil {
			// create a private_key_data_file and write private key to it.
			writeToFile(output.(string), string(jsonData))
		}
	}

	return resourceAliCloudCloudSsoScimServerCredentialUpdate(d, meta)
}

func resourceAliCloudCloudSsoScimServerCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudssoService := CloudssoService{client}

	object, err := cloudssoService.DescribeCloudSsoScimServerCredential(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_sso_scim_server_credential cloudssoService.DescribeCloudSsoScimServerCredential Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudCloudSsoScimServerCredentialRead(d, meta)
}

func resourceAliCloudCloudSsoScimServerCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteSCIMServerCredential"
	var response map[string]interface{}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := map[string]interface{}{
		"DirectoryId":  parts[0],
		"CredentialId": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("cloudsso", "2021-05-15", action, nil, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"EntityNotExists.SCIMCredential"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
