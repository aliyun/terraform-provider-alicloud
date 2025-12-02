// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPolarDbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPolarDbAccountCreate,
		Read:   resourceAliCloudPolarDbAccountRead,
		Update: resourceAliCloudPolarDbAccountUpdate,
		Delete: resourceAliCloudPolarDbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(8 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(8 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_lock_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"UnLock", "Lock"}, false),
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"account_password_valid_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "Super"}, false),
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
		},
	}
}

func resourceAliCloudPolarDbAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	}
	if v, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = v
	}

	request["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}

	password := d.Get("account_password").(string)
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
		request["AccountPassword"] = decryptResp
	}
	if v, ok := d.GetOk("account_type"); ok {
		request["AccountType"] = v
	} else {
		request["AccountType"] = "Normal"
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["AccountName"]))

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, polarDbServiceV2.PolarDbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPolarDbAccountUpdate(d, meta)
}

func resourceAliCloudPolarDbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	objectRaw, err := polarDbServiceV2.DescribePolarDbAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_account DescribePolarDbAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("account_description", objectRaw["AccountDescription"])
	d.Set("account_lock_state", objectRaw["AccountLockState"])
	d.Set("account_password_valid_time", objectRaw["AccountPasswordValidTime"])
	d.Set("account_type", objectRaw["AccountType"])
	d.Set("status", objectRaw["AccountStatus"])
	d.Set("account_name", objectRaw["AccountName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])

	return nil
}

func resourceAliCloudPolarDbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyAccountDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("account_description") {
		update = true
	}
	request["AccountDescription"] = d.Get("account_description")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || NeedRetry(err) {
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyAccountLockState"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	if d.HasChange("account_password_valid_time") {
		update = true
		request["AccountPasswordValidTime"] = d.Get("account_password_valid_time")
	}

	if d.HasChange("account_lock_state") {
		update = true
	}
	request["AccountLockState"] = d.Get("account_lock_state")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || NeedRetry(err) {
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyAccountPassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("account_password") {
		update = true

		request["NewAccountPassword"] = d.Get("account_password")
	}

	if !d.IsNewResource() && d.HasChange("kms_encrypted_password") {
		update = true

		kmsService := KmsService{meta.(*connectivity.AliyunClient)}
		decryptResp, err := kmsService.Decrypt(d.Get("kms_encrypted_password").(string), d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["NewAccountPassword"] = decryptResp
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || NeedRetry(err) {
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

	d.Partial(false)
	return resourceAliCloudPolarDbAccountRead(d, meta)
}

func resourceAliCloudPolarDbAccountDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidAccountName.NotFound", "InvalidDBClusterId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, polarDbServiceV2.PolarDbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
