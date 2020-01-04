package alicloud

import (
	"fmt"
	"strings"
	"time"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r_kvstore"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudKVstoreAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKVStoreAccountCreate,
		Read:   resourceAlicloudKVStoreAccountRead,
		Update: resourceAlicloudKVStoreAccountUpdate,
		Delete: resourceAlicloudKVStoreAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal"}, false),
				ForceNew:     true,
				Default:      "Normal",
			},
			"account_privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"RoleReadOnly", "RoleReadWrite", "RoleRepl"}, false),
				Default:      "RoleReadWrite",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudKVStoreAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	request := r_kvstore.CreateCreateAccountRequest()
	request.RegionId = client.RegionId
	request.InstanceId = d.Get("instance_id").(string)
	request.AccountName = d.Get("account_name").(string)

	password := d.Get("account_password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		request.AccountPassword = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.AccountPassword = decryptResp.Plaintext
	}
	request.AccountType = d.Get("account_type").(string)

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.AccountDescription = v.(string)
	}
	// wait instance running before modifying
	if err := kvstoreService.WaitForKVstoreInstance(request.InstanceId, Normal, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.CreateAccount(request)
		})
		if err != nil {
			if IsExceptedErrors(err, []string{"IncorrectDBInstanceState"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_kvstore_account", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", request.InstanceId, COLON_SEPARATED, request.AccountName))

	if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudKVStoreAccountRead(d, meta)
}

func resourceAlicloudKVStoreAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	object, err := kvstoreService.DescribeKVstoreAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.InstanceId)
	d.Set("account_name", object.AccountName)
	d.Set("account_type", object.AccountType)
	d.Set("description", object.AccountDescription)
	d.Set("account_privilege", object.DatabasePrivileges.DatabasePrivilege[0].AccountPrivilege)

	return nil
}

func resourceAlicloudKVStoreAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("description") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateModifyAccountDescriptionRequest()
		request.RegionId = client.RegionId
		request.InstanceId = instanceId
		request.AccountName = accountName
		request.AccountDescription = d.Get("description").(string)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ModifyAccountDescription(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("description")
	}

	if d.HasChange("account_privilege") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateGrantAccountPrivilegeRequest()
		request.RegionId = client.RegionId
		request.InstanceId = instanceId
		request.AccountName = accountName
		request.AccountPrivilege = d.Get("account_privilege").(string)

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.GrantAccountPrivilege(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("account_privilege")
	}

	if d.HasChange("account_password") || d.HasChange("kms_encrypted_password") {
		if err := kvstoreService.WaitForKVstoreAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		request := r_kvstore.CreateResetAccountPasswordRequest()
		request.RegionId = client.RegionId
		request.InstanceId = instanceId
		request.AccountName = accountName

		password := d.Get("account_password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'account_password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			request.AccountPassword = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request.AccountPassword = decryptResp.Plaintext

		}

		raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.ResetAccountPassword(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}

		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
		d.SetPartial("account_password")
	}

	d.Partial(false)
	return resourceAlicloudKVStoreAccountRead(d, meta)
}

func resourceAlicloudKVStoreAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	kvstoreService := KvstoreService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := r_kvstore.CreateDeleteAccountRequest()
	request.RegionId = client.RegionId
	request.InstanceId = parts[0]
	request.AccountName = parts[1]

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DeleteAccount(request)
	})
	if err != nil {
		if IsExceptedError(err, InvalidAccountNameNotFound) {
			return nil
		} else {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return kvstoreService.WaitForKVstoreAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
