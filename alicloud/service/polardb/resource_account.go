// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccountCreate,
		Read:   resourceAccountRead,
		Update: resourceAccountUpdate,
		Delete: resourceAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(8 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(8 * time.Minute),
		},
		ValidateRawResourceConfigFuncs: []schema.ValidateRawResourceConfigFunc{
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("account_password"), cty.GetAttrPath("account_password_wo")),
			validation.PreferWriteOnlyAttribute(cty.GetAttrPath("kms_encrypted_password"), cty.GetAttrPath("account_password_wo")),
		},
		// CreateAccount applies AccountPrivilege only when DBName is sent in the same request, and
		// applies PrivForAllDB only when AccountPrivilege is sent. Reject an incomplete pairing at
		// plan time instead of silently creating an account without the declared grants. Values that
		// are still unknown at plan time are validated again in Create once they resolve.
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
			return validateAccountCreatePairing(
				diff.Get("db_name").(string), diff.NewValueKnown("db_name"),
				diff.Get("account_privilege").(string), diff.NewValueKnown("account_privilege"),
				diff.Get("priv_for_all_db").(string), diff.NewValueKnown("priv_for_all_db"))
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
				ValidateFunc: validators.StringInSlice([]string{"UnLock", "Lock"}, false),
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"kms_encrypted_password", "account_password_wo"},
				AtLeastOneOf:  []string{"account_password", "kms_encrypted_password", "account_password_wo"},
			},
			"account_password_wo": {
				Type:          schema.TypeString,
				Optional:      true,
				WriteOnly:     true,
				ConflictsWith: []string{"account_password", "kms_encrypted_password"},
				RequiredWith:  []string{"account_password_wo_version"},
				AtLeastOneOf:  []string{"account_password", "kms_encrypted_password", "account_password_wo"},
			},
			"account_password_wo_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"account_password_wo"},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"account_password_valid_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_privilege": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"db_name"},
				ValidateFunc: validatePolarDbAccountPrivilege,
			},
			"account_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateFunc:     validators.StringInSlice([]string{"Normal", "Super", "DynamoDB"}, false),
				DiffSuppressFunc: polarDbAccountTypeDiffSuppressFunc,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"Search"}, false),
			},
			"password_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priv_for_all_db": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"0", "1"}, false),
			},
			"reset_permission_flag": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_encrypted_password": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: kmsDiffSuppressFunc,
				ConflictsWith:    []string{"account_password", "account_password_wo"},
				AtLeastOneOf:     []string{"account_password", "kms_encrypted_password", "account_password_wo"},
			},
			"kms_encryption_context": {
				Type:     schema.TypeMap,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("kms_encrypted_password").(string) == ""
				},
				Elem: schema.TypeString,
			},
			"dynamodb_auth_password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The DynamoDB authentication password. Only available for DynamoDB account type.",
			},
		},
	}
}

func resourceAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	// Create-time values are known even when the plan could not resolve them; re-check the pairing
	// that the API silently ignores instead of sending a request without the declared grants.
	if err := validateAccountCreatePairing(d.Get("db_name").(string), true, d.Get("account_privilege").(string), true, d.Get("priv_for_all_db").(string), true); err != nil {
		return errs.WrapError(err)
	}

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

	request["ClientToken"] = helper.BuildClientToken(action)

	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}

	password := d.Get("account_password").(string)
	if woValue, err := helper.GetWriteOnlyValue(d, cty.GetAttrPath("account_password_wo"), cty.String); err != nil {
		return errs.WrapError(err)
	} else if !woValue.IsNull() {
		password = woValue.AsString()
	}
	kmsPassword := d.Get("kms_encrypted_password").(string)

	if password == "" && kmsPassword == "" {
		return errs.WrapError(errs.Error("One of the 'account_password', 'kms_encrypted_password' and 'account_password_wo' should be set."))
	}
	if password != "" {
		request["AccountPassword"] = password
	} else {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return errs.WrapError(err)
		}
		request["AccountPassword"] = decryptResp
	}
	if v, ok := d.GetOk("account_type"); ok {
		request["AccountType"] = v
	} else {
		request["AccountType"] = "Normal"
	}
	if v, ok := d.GetOk("db_name"); ok {
		request["DBName"] = v
	}
	if v, ok := d.GetOk("node_type"); ok {
		request["NodeType"] = v
	}
	if v, ok := d.GetOk("priv_for_all_db"); ok {
		request["PrivForAllDB"] = v
	}
	if v, ok := d.GetOk("account_privilege"); ok {
		request["AccountPrivilege"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout", "EngineMigration.ActionDisabled"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["AccountName"]))
	if v, ok := d.GetOk("node_type"); ok {
		d.SetId(fmt.Sprintf("%v:%v", d.Id(), v))
	}

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, polarDbServiceV2.PolarDbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAccountUpdate(d, meta)
}

func resourceAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	objectRaw, err := polarDbServiceV2.DescribePolarDbAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_account DescribePolarDbAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("account_description", objectRaw["AccountDescription"])
	d.Set("account_lock_state", objectRaw["AccountLockState"])
	d.Set("account_password_valid_time", objectRaw["AccountPasswordValidTime"])
	if v, ok := objectRaw["AccountType"].(string); ok && v != "" {
		d.Set("account_type", v)
	}
	d.Set("status", objectRaw["AccountStatus"])
	d.Set("account_name", objectRaw["AccountName"])
	if v, ok := objectRaw["DynamoDBAuthPassword"]; ok && v != nil {
		d.Set("dynamodb_auth_password", v)
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])
	if len(parts) == 3 && parts[2] != "" {
		d.Set("node_type", parts[2])
	}

	return nil
}

func resourceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || helper.NeedRetry(err) {
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
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || helper.NeedRetry(err) {
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ResetAccountPassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	password := ""
	woValue, err := helper.GetWriteOnlyValue(d, cty.GetAttrPath("account_password_wo"), cty.String)
	if err != nil {
		return errs.WrapError(err)
	}
	woPassword := ""
	if !woValue.IsNull() {
		woPassword = woValue.AsString()
	}
	passwordChanged := !d.IsNewResource() && d.HasChange("account_password")
	kmsChanged := !d.IsNewResource() && d.HasChange("kms_encrypted_password")
	woVersionChanged := !d.IsNewResource() && d.HasChange("account_password_wo_version")
	passwordType := d.Get("password_type").(string)
	// password_type only changes how a password is written: setting it without a password source
	// change must not issue a password write.
	if accountPasswordSourceChanged(passwordChanged, d.Get("account_password").(string), kmsChanged, d.Get("kms_encrypted_password").(string), woVersionChanged, woPassword) {
		update = true
		switch {
		case passwordChanged && d.Get("account_password").(string) != "":
			password = d.Get("account_password").(string)
		case kmsChanged && d.Get("kms_encrypted_password").(string) != "":
			kmsService := KmsService{meta.(*connectivity.AliyunClient)}
			decryptResp, err := kmsService.Decrypt(d.Get("kms_encrypted_password").(string), d.Get("kms_encryption_context").(map[string]interface{}))
			if err != nil {
				return errs.WrapError(err)
			}
			password = decryptResp
		case woVersionChanged && woPassword != "":
			password = woPassword
		}
	}
	if update {
		// ModifyAccountPassword is the only API that accepts PasswordType; the regular
		// password update goes through ResetAccountPassword.
		if passwordType != "" {
			action = "ModifyAccountPassword"
			request["NewAccountPassword"] = password
			request["PasswordType"] = passwordType
		} else {
			request["AccountPassword"] = password
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout", "EngineMigration.ActionDisabled"}) || helper.NeedRetry(err) {
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

	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ResetAccount"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("reset_permission_flag") && d.Get("reset_permission_flag").(bool) {
		if v, ok := d.GetOk("account_type"); !ok || v.(string) != "Super" {
			return errs.WrapError(errs.Error("The 'reset_permission_flag' can only be used to reset the permissions of a Super account."))
		}
		update = true
		if password == "" {
			password = d.Get("account_password").(string)
		}
		if password != "" {
			request["AccountPassword"] = password
		}
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || helper.NeedRetry(err) {
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

		polarDbServiceV2 := PolarDbServiceV2{client}
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, polarDbServiceV2.PolarDbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, IdMsg, d.Id())
		}
	}

	d.Partial(false)
	return resourceAccountRead(d, meta)
}

func resourceAccountDelete(d *schema.ResourceData, meta interface{}) error {

	// DynamoDB type account does not support deletion
	if v, ok := d.GetOk("account_type"); ok && v.(string) == "DynamoDB" {
		return nil
	}

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

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "LockTimeout"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidAccountName.NotFound", "InvalidDBClusterId.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, polarDbServiceV2.PolarDbAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

// DescribeAccounts reports the privilege-scoped account type (for example ReadOnly for an account
// created with a cluster-wide read-only privilege) even though CreateAccount was called with Normal.
// Suppressing that one transition keeps the configured type as the plan intent; any change the user
// actually makes to account_type still produces a replacement diff.
func polarDbAccountTypeDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return new == "Normal" && old == "ReadOnly"
}

func validatePolarDbAccountPrivilege(v interface{}, k string) (ws []string, es []error) {
	value, ok := v.(string)
	if !ok {
		return nil, []error{fmt.Errorf("%q must be a string", k)}
	}
	if len(value) > 900 {
		return nil, []error{fmt.Errorf("%q cannot exceed 900 characters", k)}
	}
	allowed := map[string]bool{
		"ReadWrite": true,
		"ReadOnly":  true,
		"DMLOnly":   true,
		"DDLOnly":   true,
		"ReadIndex": true,
	}
	for _, privilege := range strings.Split(value, ",") {
		if privilege == "" {
			es = append(es, fmt.Errorf("%q contains an empty value, use a comma-separated list such as \"ReadWrite,ReadOnly\"", k))
			continue
		}
		if !allowed[privilege] {
			es = append(es, fmt.Errorf("%q contains invalid value %q, valid values: ReadWrite, ReadOnly, DMLOnly, DDLOnly, ReadIndex", k, privilege))
		}
	}
	return ws, es
}

// validateAccountCreatePairing mirrors the CreateAccount contract: AccountPrivilege takes effect
// only when DBName is sent in the same request, and PrivForAllDB="1" takes effect only when
// AccountPrivilege is sent. Unknown plan values are marked not known and are re-validated in Create.
func validateAccountCreatePairing(dbName string, dbNameKnown bool, accountPrivilege string, privilegeKnown bool, privForAllDB string, privKnown bool) error {
	if privilegeKnown && accountPrivilege != "" && dbNameKnown && dbName == "" {
		return fmt.Errorf(`"account_privilege" requires "db_name" to contain at least one database name`)
	}
	if privKnown && privForAllDB == "1" && privilegeKnown && accountPrivilege == "" {
		return fmt.Errorf(`"priv_for_all_db" set to "1" requires "account_privilege" to be set`)
	}
	return nil
}

// accountPasswordSourceChanged reports whether the configuration changes the password itself.
// password_type only changes how a password is written, so it must not trigger a write on its own,
// and removing a configured password source is not a password write either.
func accountPasswordSourceChanged(passwordChanged bool, password string, kmsPasswordChanged bool, kmsPassword string, woVersionChanged bool, woPassword string) bool {
	return passwordChanged && password != "" || kmsPasswordChanged && kmsPassword != "" || woVersionChanged && woPassword != ""
}
