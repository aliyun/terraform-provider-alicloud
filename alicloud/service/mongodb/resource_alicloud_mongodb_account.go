// Package mongodb. This file is generated automatically. Please do not modify it manually, thank you!
package mongodb

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudMongodbAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbAccountCreate,
		Read:   resourceAliCloudMongodbAccountRead,
		Update: resourceAliCloudMongodbAccountUpdate,
		Delete: resourceAliCloudMongodbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"character_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validators.StringInSlice([]string{"db", "cs", "normal"}, false),
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudMongodbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	instanceId := fmt.Sprint(d.Get("instance_id"))
	accountName := fmt.Sprint(d.Get("account_name"))
	characterType := fmt.Sprint(d.Get("character_type"))
	dbName := fmt.Sprint(d.Get("db_name"))
	if dbName == "" {
		dbName = "admin"
	}
	_, hasRoles := d.GetOk("roles")
	mongodbServiceV2 := NewMongodbServiceV2(client)

	if characterType != "db" && characterType != "cs" {
		// Replica set instances expose the built-in root account, which is
		// managed through ResetAccountPassword; that API does not accept roles.
		if hasRoles {
			return errs.WrapError(fmt.Errorf("roles can only be set when creating a config server account (character_type = \"cs\"); account %s is a replica set account", accountName))
		}
		action := "ResetAccountPassword"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["AccountName"] = accountName
		request["DBInstanceId"] = instanceId
		request["RegionId"] = client.RegionId
		request["AccountPassword"] = d.Get("account_password")
		request["DBName"] = dbName
		if v, ok := d.GetOk("character_type"); ok {
			request["CharacterType"] = v
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_mongodb_account", action, errs.AlibabaCloudSdkGoERROR)
		}

		d.SetId(fmt.Sprintf("%v:%v:%v", instanceId, dbName, accountName))
		if err := d.Set("db_name", dbName); err != nil {
			return errs.WrapError(err)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(mongodbAccountId(d), "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

		if v, ok := d.GetOk("account_description"); ok {
			// ModifyAccountDescription only supports replica set accounts, so a
			// create-time description is applied here instead of reusing Update.
			action = "ModifyAccountDescription"
			request = make(map[string]interface{})
			query = make(map[string]interface{})
			request["AccountName"] = accountName
			request["DBInstanceId"] = instanceId
			request["RegionId"] = client.RegionId
			request["AccountDescription"] = v
			request["DBName"] = dbName
			err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
				response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
				if err != nil {
					if helper.NeedRetry(err) || errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
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
			stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(mongodbAccountId(d), "AccountStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

		return resourceAliCloudMongodbAccountRead(d, meta)
	}

	// Product behaviour: shard direct-connection accounts have fixed privileges,
	// so CreateAccount rejects Roles with InvalidRoles.NotSupported; config server
	// accounts require Roles.
	if characterType == "db" && hasRoles {
		return errs.WrapError(fmt.Errorf("roles are not supported for shard direct-connection accounts (character_type = \"db\"): CreateAccount rejects them with InvalidRoles.NotSupported, account %s", accountName))
	}
	if characterType == "cs" && !hasRoles {
		return errs.WrapError(fmt.Errorf("roles are required for config server accounts (character_type = \"cs\"), account %s", accountName))
	}

	action := "CreateAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["AccountName"] = accountName
	request["DBInstanceId"] = instanceId
	request["RegionId"] = client.RegionId
	request["AccountPassword"] = d.Get("account_password")
	request["DBName"] = dbName
	request["CharacterType"] = characterType
	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	if hasRoles {
		rolesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range d.Get("roles").(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["DBName"] = dataLoopTmp["db_name"]
			dataLoopMap["Role"] = dataLoopTmp["role"]
			rolesMapsArray = append(rolesMapsArray, dataLoopMap)
		}
		request["Roles"] = rolesMapsArray
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
		if err != nil {
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_mongodb_account", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v:%v", instanceId, dbName, accountName))
	if err := d.Set("db_name", dbName); err != nil {
		return errs.WrapError(err)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(mongodbAccountId(d), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudMongodbAccountRead(d, meta)
}

func resourceAliCloudMongodbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := NewMongodbServiceV2(client)

	objectRaw, err := mongodbServiceV2.DescribeMongodbAccount(mongodbAccountId(d))
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_account DescribeMongodbAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	// DescribeAccounts only returns AccountDescription when the description was
	// set through ModifyAccountDescription, which is not supported for shard
	// accounts, so a configured description is kept when the API omits it.
	if v, ok := objectRaw["AccountDescription"]; ok && fmt.Sprint(v) != "" {
		d.Set("account_description", v)
	}
	d.Set("character_type", objectRaw["CharacterType"])
	// db_name is part of the identity, so it is only overwritten when the API
	// returns a value for it.
	if v, ok := objectRaw["DBName"]; ok && fmt.Sprint(v) != "" {
		d.Set("db_name", v)
	}
	d.Set("status", objectRaw["AccountStatus"])
	d.Set("account_name", objectRaw["AccountName"])
	d.Set("instance_id", objectRaw["DBInstanceId"])

	roles := make([]interface{}, 0)
	if v, ok := objectRaw["Roles"]; ok {
		roles = mongodbAccountRoles(v)
	}
	if err := d.Set("roles", roles); err != nil {
		return errs.WrapError(err)
	}

	// The canonical id carries the account identity scope, including the
	// database, so two accounts with the same name in different databases keep
	// distinct ids. Legacy two-part ids are upgraded in place on the first read.
	if instanceId, accountName, dbName := fmt.Sprint(objectRaw["DBInstanceId"]), fmt.Sprint(objectRaw["AccountName"]), fmt.Sprint(objectRaw["DBName"]); instanceId != "" && accountName != "" && dbName != "" {
		d.SetId(fmt.Sprintf("%v:%v:%v", instanceId, dbName, accountName))
	}

	return nil
}

func resourceAliCloudMongodbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	instanceId := fmt.Sprint(d.Get("instance_id"))
	accountName := fmt.Sprint(d.Get("account_name"))
	dbName := fmt.Sprint(d.Get("db_name"))
	accountId := mongodbAccountId(d)

	mongodbServiceV2 := NewMongodbServiceV2(client)

	if d.HasChange("account_password") {
		action := "ResetAccountPassword"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var err error
		request = make(map[string]interface{})
		request["AccountName"] = accountName
		request["DBInstanceId"] = instanceId
		request["RegionId"] = client.RegionId
		request["AccountPassword"] = d.Get("account_password")
		if dbName != "" {
			request["DBName"] = dbName
		}
		if v, ok := d.GetOk("character_type"); ok {
			// CharacterType is mandatory for sharded instances.
			request["CharacterType"] = v
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
			if err != nil {
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, errs.OperationDeniedDBStatus) {
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
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(accountId, "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	if d.HasChange("account_description") {
		objectRaw, err := mongodbServiceV2.DescribeMongodbAccount(accountId)
		if err != nil {
			return errs.WrapError(err)
		}
		if fmt.Sprint(objectRaw["CharacterType"]) != "normal" {
			return errs.WrapError(fmt.Errorf("account_description of account %s cannot be modified: ModifyAccountDescription only supports replica set (normal) accounts, the account character_type is %v", d.Id(), objectRaw["CharacterType"]))
		}
		action := "ModifyAccountDescription"
		var request map[string]interface{}
		var response map[string]interface{}
		query := make(map[string]interface{})
		var rpcErr error
		request = make(map[string]interface{})
		request["AccountName"] = accountName
		request["DBInstanceId"] = instanceId
		request["RegionId"] = client.RegionId
		request["AccountDescription"] = d.Get("account_description")
		if dbName != "" {
			request["DBName"] = dbName
		}
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		rpcErr = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, rpcErr = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
			if rpcErr != nil {
				if helper.NeedRetry(rpcErr) || errs.IsExpectedErrors(rpcErr, errs.OperationDeniedDBStatus) {
					wait()
					return retry.RetryableError(rpcErr)
				}
				return retry.NonRetryableError(rpcErr)
			}
			return nil
		})
		logs.AddDebug(action, response, request)
		if rpcErr != nil {
			return errs.WrapErrorf(rpcErr, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}
		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(accountId, "AccountStatus", []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}
	}

	return resourceAliCloudMongodbAccountRead(d, meta)
}

func resourceAliCloudMongodbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Account. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

// mongodbAccountId returns the database scoped id used for all account lookups.
// Account names are only unique inside a database, so the database name has to
// be part of every DescribeAccounts request once it is known.
func mongodbAccountId(d *schema.ResourceData) string {
	instanceId := fmt.Sprint(d.Get("instance_id"))
	accountName := fmt.Sprint(d.Get("account_name"))
	dbName := fmt.Sprint(d.Get("db_name"))
	if instanceId != "" && accountName != "" && dbName != "" {
		return fmt.Sprintf("%v:%v:%v", instanceId, dbName, accountName)
	}
	return d.Id()
}

func mongodbAccountRoles(raw interface{}) []interface{} {
	roles := make([]interface{}, 0)
	roleMap, ok := raw.(map[string]interface{})
	if !ok {
		return roles
	}
	roleList, ok := roleMap["Role"].([]interface{})
	if !ok {
		return roles
	}
	for _, item := range roleList {
		role, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		roles = append(roles, map[string]interface{}{
			"db_name": role["DBName"],
			"role":    role["Role"],
		})
	}
	return roles
}
