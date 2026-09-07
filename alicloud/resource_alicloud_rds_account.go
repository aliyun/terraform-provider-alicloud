package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAliCloudRdsAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudRdsAccountCreate,
		Read:   resourceAliCloudRdsAccountRead,
		Update: resourceAliCloudRdsAccountUpdate,
		Delete: resourceAliCloudRdsAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"description"},
			},
			"description": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'description' has been deprecated from provider version 1.120.0. New field 'account_description' instead.",
				ConflictsWith: []string{"account_description"},
			},
			"account_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{0,61}[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and must begin with letters and end with letters or numbers"),
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{0,61}[a-z0-9]$`), "The name can consist of lowercase letters, numbers, underscores, and must begin with letters and end with letters or numbers"),
				Deprecated:    "Field 'name' has been deprecated from provider version 1.120.0. New field 'account_name' instead.",
				ConflictsWith: []string{"account_name"},
			},
			"account_password": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"password"},
			},
			"password": {
				Type:          schema.TypeString,
				Sensitive:     true,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'password' has been deprecated from provider version 1.120.0. New field 'account_password' instead.",
				ConflictsWith: []string{"account_password"},
			},
			"account_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Normal", "Super", "Sysadmin"}, false),
				ConflictsWith: []string{"type"},
			},
			"check_policy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Normal", "Super"}, false),
				Deprecated:    "Field 'type' has been deprecated from provider version 1.120.0. New field 'account_type' instead.",
				ConflictsWith: []string{"account_type"},
			},
			"db_instance_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"db_instance_id", "instance_id"},
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'instance_id' has been deprecated from provider version 1.120.0. New field 'db_instance_id' instead.",
				ConflictsWith: []string{"db_instance_id"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
			"reset_permission_flag": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliCloudRdsAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var response map[string]interface{}
	action := "CreateAccount"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	var err error
	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	} else if v, ok := d.GetOk("description"); ok {
		request["AccountDescription"] = v
	}

	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["AccountName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "account_name" must be set one!`))
	}

	request["AccountPassword"] = d.Get("account_password")
	if v, ok := d.GetOk("account_password"); ok {
		request["AccountPassword"] = v
	} else if v, ok := d.GetOk("password"); ok {
		request["AccountPassword"] = v
	} else if v, ok := d.GetOk("kms_encrypted_password"); ok {
		kmsService := KmsService{client}
		decryptResp, err := kmsService.Decrypt(v.(string), d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["AccountPassword"] = decryptResp
	} else {
		return WrapError(Error("One of the 'account_password' and 'password' and 'kms_encrypted_password' should be set."))
	}
	if v, ok := d.GetOk("account_type"); ok {
		request["AccountType"] = v
	} else if v, ok := d.GetOk("type"); ok {
		request["AccountType"] = v
	}

	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	} else if v, ok := d.GetOk("instance_id"); ok {
		request["DBInstanceId"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "instance_id" or "db_instance_id" must be set one!`))
	}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(request["DBInstanceId"].(string), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_rds_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"], ":", request["AccountName"]))
	stateConf = BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	stateConf = BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, rdsService.RdsDBInstanceStateRefreshFunc(request["DBInstanceId"].(string), []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudRdsAccountRead(d, meta)
}
func resourceAliCloudRdsAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeRdsAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_rds_account rdsService.DescribeRdsAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("account_name", parts[1])
	d.Set("name", parts[1])
	d.Set("db_instance_id", parts[0])
	d.Set("instance_id", parts[0])
	d.Set("check_policy", object["CheckPolicy"])
	d.Set("account_description", object["AccountDescription"])
	d.Set("description", object["AccountDescription"])
	d.Set("account_type", object["AccountType"])
	d.Set("type", object["AccountType"])
	d.Set("status", object["AccountStatus"])
	return nil
}

func resourceAliCloudRdsAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	var err error
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Partial(true)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	update := false
	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
		"RegionId":     client.RegionId,
	}

	rdsServiceV2 := RdsServiceV2{client}
	objectRaw, _ := rdsServiceV2.DescribeRdsAccount(d.Id())

	if d.HasChange("status") {
		var err error
		target := d.Get("status").(string)

		currentStatus, err := jsonpath.Get("AccountStatus", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "AccountStatus", objectRaw)
		}
		if fmt.Sprint(currentStatus) != target {
			if target == "Lock" {
				action := "LockAccount"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
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
			if target == "Available" {
				action := "UnlockAccount"
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
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
		}
	}

	if d.HasChange("account_description") {
		update = true
		request["AccountDescription"] = d.Get("account_description")
	} else if d.HasChange("description") {
		update = true
		request["AccountDescription"] = d.Get("description")
	}
	if update {
		action := "ModifyAccountDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("description")
		d.SetPartial("account_description")
	}
	update = false
	resetAccountPasswordReq := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
	}
	if d.HasChange("account_password") {
		update = true
		resetAccountPasswordReq["AccountPassword"] = d.Get("account_password").(string)
	} else if d.HasChange("password") {
		update = true
		resetAccountPasswordReq["AccountPassword"] = d.Get("password").(string)
	} else if d.HasChange("kms_encrypted_password") {
		update = true
		kmsPassword := d.Get("kms_encrypted_password").(string)
		kmsService := KmsService{meta.(*connectivity.AliyunClient)}
		decryptResp, err := kmsService.Decrypt(kmsPassword, d.Get("kms_encryption_context").(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		resetAccountPasswordReq["AccountPassword"] = decryptResp
	}
	if update {
		action := "ResetAccountPassword"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, resetAccountPasswordReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, resetAccountPasswordReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("password")
		d.SetPartial("kms_encrypted_password")
		d.SetPartial("kms_encryption_context")
		d.SetPartial("account_password")
	}

	resetPermission := false
	resetAccountReq := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
	}
	if v, _ := d.GetOkExists("reset_permission_flag"); v.(bool) && !d.IsNewResource() && d.HasChange("reset_permission_flag") {
		if accountType, ok := d.GetOk("account_type"); ok {
			if accountType.(string) == "Super" {
				resetPermission = true
			}
		}
		if accountType, ok := d.GetOk("type"); ok {
			if accountType.(string) == "Super" {
				resetPermission = true
			}
		}
	}

	if d.HasChange("check_policy") {
		action := "ModifyAccountCheckPolicy"
		request["CheckPolicy"] = d.Get("check_policy")
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, true)
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
	if resetPermission {
		if v, ok := d.GetOk("account_password"); ok {
			resetAccountReq["AccountPassword"] = v.(string)
		}
		if v, ok := d.GetOk("password"); ok {
			resetAccountReq["AccountPassword"] = v.(string)
		}
		// ResetAccount interface can also reset the database account password
		action := "ResetAccount"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Rds", "2014-08-15", action, nil, resetAccountReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, resetAccountReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	d.Partial(false)
	return resourceAliCloudRdsAccountRead(d, meta)
}

func resourceAliCloudRdsAccountDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	rdsService := RdsService{client}
	action := "DeleteAccount"
	var response map[string]interface{}
	request := map[string]interface{}{
		"AccountName":  parts[1],
		"DBInstanceId": parts[0],
		"SourceIp":     client.SourceIp,
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"InternalError", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "AccountActionForbidden", "IncorrectDBInstanceState"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		object, err := rdsService.DescribeRdsAccount(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(err)
		}
		if fmt.Sprint(object["AccountStatus"]) == "Lock" {
			action = "UnlockAccount"
			wait = incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
				response, err = client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				addDebug(action, response, request)
				return nil
			})
			if err != nil {
				return resource.NonRetryableError(err)
			}
			action = "DeleteAccount"
			return resource.RetryableError(fmt.Errorf("there need to delete account %s again after unlock it", d.Id()))
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, rdsService.RdsAccountStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
