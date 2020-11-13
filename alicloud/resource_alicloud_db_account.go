package alicloud

import (
	"fmt"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDBAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBAccountCreate,
		Read:   resourceAlicloudDBAccountRead,
		Update: resourceAlicloudDBAccountUpdate,
		Delete: resourceAlicloudDBAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"password": {
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
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Normal", "Super"}, false),
				ForceNew:     true,
				Computed:     true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDBAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateAccount"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": d.Get("instance_id").(string),
		"AccountName":  d.Get("name").(string),
		"SourceIp":     client.SourceIp,
	}

	password := d.Get("password").(string)
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
	}

	if password != "" {
		request["AccountPassword"] = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["AccountPassword"] = decryptResp.Plaintext
	}
	if v, ok := d.GetOk("AccountType"); ok {
		request["AccountType"] = v
	}
	// Description will not be set when account type is normal and it is a API bug
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request["AccountDescription"] = v
	}
	// wait instance running before modifying
	if err := rdsService.WaitForDBInstance(request["DBInstanceId"].(string), Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	conn, err := meta.(*connectivity.AliyunClient).NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, OperationDeniedDBStatus) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_db_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v%s%v", request["DBInstanceId"], COLON_SEPARATED, request["AccountName"]))

	if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudDBAccountRead(d, meta)
}

func resourceAlicloudDBAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object["DBInstanceId"])
	d.Set("name", object["AccountName"])
	d.Set("type", object["AccountType"])
	d.Set("description", object["AccountDescription"])

	return nil
}

func resourceAlicloudDBAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	d.Partial(true)
	parts := strings.Split(d.Id(), COLON_SEPARATED)
	instanceId := parts[0]
	accountName := parts[1]

	if d.HasChange("description") {
		if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		action := "ModifyAccountDescription"
		request := map[string]interface{}{
			"RegionId":           client.RegionId,
			"DBInstanceId":       instanceId,
			"AccountName":        accountName,
			"AccountDescription": d.Get("description"),
			"SourceIp":           client.SourceIp,
		}
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("description")
	}

	if d.HasChange("password") || d.HasChange("kms_encrypted_password") {
		if err := rdsService.WaitForAccount(d.Id(), Available, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		action := "ResetAccountPassword"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": instanceId,
			"AccountName":  accountName,
			"SourceIp":     client.SourceIp,
		}
		password := d.Get("password").(string)
		kmsPassword := d.Get("kms_encrypted_password").(string)

		if password == "" && kmsPassword == "" {
			return WrapError(Error("One of the 'password' and 'kms_encrypted_password' should be set."))
		}

		if password != "" {
			d.SetPartial("password")
			request["AccountPassword"] = password
		} else {
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["AccountPassword"] = decryptResp.Plaintext
			d.SetPartial("kms_encrypted_password")
			d.SetPartial("kms_encryption_context")
		}
		conn, err := client.NewRdsClient()
		if err != nil {
			return WrapError(err)
		}
		response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		d.SetPartial("password")
	}

	d.Partial(false)
	return resourceAlicloudDBAccountRead(d, meta)
}

func resourceAlicloudDBAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteAccount"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": parts[0],
		"AccountName":  parts[1],
		"SourceIp":     client.SourceIp,
	}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	response, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	return rdsService.WaitForAccount(d.Id(), Deleted, DefaultTimeoutMedium)
}
