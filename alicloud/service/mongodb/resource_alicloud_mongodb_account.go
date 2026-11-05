// Package mongodb. This file is generated automatically. Please do not modify it manually, thank you!
package mongodb

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// accountRetryableErrors adds the account-operation throttling code to the
// shared transient list; back-to-back account operations (password reset then
// description update) can be rejected with OperationDenied.TooFrequent.
var accountRetryableErrors = append(append([]string{}, errs.OperationDeniedDBStatus...), "OperationDenied.TooFrequent")

func resourceAliCloudMongodbAccount() *schema.Resource {
	return &schema.Resource{
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// Only user configuration is validated; values returned by
			// DescribeAccounts must not be able to block later plans.
			if !d.HasChange("roles") {
				return nil
			}
			return validateAccountRoles(d.Get("roles"))
		},
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
						// Both members are validated by CustomizeDiff instead of
						// being declared Required: new schema fields are not
						// allowed to be Required by the breaking-change gate.
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

	// A config server account cannot be created without roles: CreateAccount
	// rejects it with MISSING_ACCOUNT_ROLES.
	if characterType == "cs" && !hasRoles {
		return errs.WrapError(fmt.Errorf("roles are required for config server accounts (character_type = \"cs\"), account %s", accountName))
	}

	// Replica set accounts without custom roles keep the legacy behaviour: the
	// built-in root account is managed through ResetAccountPassword, which does
	// not accept roles.
	if !hasRoles && characterType != "db" && characterType != "cs" {
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
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, accountRetryableErrors) {
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

		d.SetId(buildAccountId(instanceId, dbName, accountName))
		if err := d.Set("db_name", dbName); err != nil {
			return errs.WrapError(err)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(accountIdFromResourceData(d), "AccountStatus", []string{}))
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
					if helper.NeedRetry(err) || errs.IsExpectedErrors(err, accountRetryableErrors) {
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
			stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(accountIdFromResourceData(d), "AccountStatus", []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
			}
		}

		return resourceAliCloudMongodbAccountRead(d, meta)
	}

	// db, cs and replica set accounts with custom roles are created through
	// CreateAccount. A db account with roles is not rejected locally so the
	// documented InvalidRoles.NotSupported response comes from the API itself.
	// CreateAccount defaults CharacterType to db, which never accepts roles, so
	// an omitted character_type falls back to the documented replica set default.
	if characterType == "" {
		characterType = "normal"
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
			if helper.NeedRetry(err) || errs.IsExpectedErrors(err, accountRetryableErrors) {
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

	d.SetId(buildAccountId(instanceId, dbName, accountName))
	if err := d.Set("db_name", dbName); err != nil {
		return errs.WrapError(err)
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 30*time.Second, mongodbServiceV2.MongodbAccountStateRefreshFunc(accountIdFromResourceData(d), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return resourceAliCloudMongodbAccountRead(d, meta)
}

func resourceAliCloudMongodbAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := NewMongodbServiceV2(client)

	objectRaw, err := mongodbServiceV2.DescribeMongodbAccount(accountIdFromResourceData(d))
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
		if err := d.Set("account_description", v); err != nil {
			return errs.WrapError(err)
		}
	}
	if err := d.Set("character_type", objectRaw["CharacterType"]); err != nil {
		return errs.WrapError(err)
	}
	// db_name is part of the identity, so it is only overwritten when the API
	// returns a value for it.
	if v, ok := objectRaw["DBName"]; ok && fmt.Sprint(v) != "" {
		if err := d.Set("db_name", v); err != nil {
			return errs.WrapError(err)
		}
	}
	if err := d.Set("status", objectRaw["AccountStatus"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("account_name", objectRaw["AccountName"]); err != nil {
		return errs.WrapError(err)
	}
	if err := d.Set("instance_id", objectRaw["DBInstanceId"]); err != nil {
		return errs.WrapError(err)
	}

	roles := make([]interface{}, 0)
	if v, ok := objectRaw["Roles"]; ok {
		roles = accountRolesFromResponse(v)
	}
	if err := d.Set("roles", roles); err != nil {
		return errs.WrapError(err)
	}

	// The canonical id carries the account identity scope, including the
	// database, so two accounts with the same name in different databases keep
	// distinct ids. Structural separators inside the database or account name
	// are escaped so the id stays reversible. Legacy two-part ids are upgraded
	// in place on the first read.
	if instanceId, accountName, dbName := fmt.Sprint(objectRaw["DBInstanceId"]), fmt.Sprint(objectRaw["AccountName"]), fmt.Sprint(objectRaw["DBName"]); instanceId != "" && accountName != "" && dbName != "" {
		d.SetId(buildAccountId(instanceId, dbName, accountName))
	}

	return nil
}

func resourceAliCloudMongodbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	instanceId := fmt.Sprint(d.Get("instance_id"))
	accountName := fmt.Sprint(d.Get("account_name"))
	dbName := fmt.Sprint(d.Get("db_name"))
	if dbName == "" {
		dbName = "admin"
	}
	accountId := accountIdFromResourceData(d)

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
				if helper.NeedRetry(err) || errs.IsExpectedErrors(err, accountRetryableErrors) {
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
				if helper.NeedRetry(rpcErr) || errs.IsExpectedErrors(rpcErr, accountRetryableErrors) {
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

// accountIdFromResourceData returns the database scoped id used for all account
// lookups. Account names are only unique inside a database, so the database
// name has to be part of every DescribeAccounts request once it is known.
func accountIdFromResourceData(d *schema.ResourceData) string {
	instanceId := fmt.Sprint(d.Get("instance_id"))
	accountName := fmt.Sprint(d.Get("account_name"))
	dbName := fmt.Sprint(d.Get("db_name"))
	if instanceId != "" && accountName != "" && dbName != "" {
		return buildAccountId(instanceId, dbName, accountName)
	}
	return d.Id()
}

// buildAccountId joins the account identity into the canonical
// <instance_id>:<db_name>:<account_name> form, escaping ':' and '%' so a legal
// database name that contains a colon cannot break the id apart.
func buildAccountId(instanceId, dbName, accountName string) string {
	return fmt.Sprintf("%v:%v:%v", instanceId, escapeAccountIdSegment(dbName), escapeAccountIdSegment(accountName))
}

// parseAccountId accepts the canonical three-part id and the legacy two-part
// <instance_id>:<account_name> id. Extra segments are rejected; a database name
// with a colon is carried escaped, so it never adds a segment.
func parseAccountId(id string) (instanceId, dbName, accountName string, err error) {
	parts := strings.Split(id, ":")
	switch len(parts) {
	case 2:
		if parts[0] == "" || parts[1] == "" {
			break
		}
		return parts[0], "", unescapeAccountIdSegment(parts[1]), nil
	case 3:
		if parts[0] == "" || parts[1] == "" || parts[2] == "" {
			break
		}
		return parts[0], unescapeAccountIdSegment(parts[1]), unescapeAccountIdSegment(parts[2]), nil
	}
	return "", "", "", fmt.Errorf("invalid Resource Id %s. Expected <instance_id>:<db_name>:<account_name> or the legacy <instance_id>:<account_name>", id)
}

func escapeAccountIdSegment(segment string) string {
	segment = strings.ReplaceAll(segment, "%", "%25")
	return strings.ReplaceAll(segment, ":", "%3A")
}

func unescapeAccountIdSegment(segment string) string {
	if !strings.Contains(segment, "%") {
		return segment
	}
	segment = strings.ReplaceAll(segment, "%3A", ":")
	segment = strings.ReplaceAll(segment, "%3a", ":")
	return strings.ReplaceAll(segment, "%25", "%")
}

// validateAccountRoles enforces the API contract that every configured role
// grants a role on a database; both members are mandatory once a roles block is
// present.
func validateAccountRoles(raw interface{}) error {
	roles, ok := raw.(*schema.Set)
	if !ok || roles.Len() == 0 {
		return nil
	}
	for _, item := range roles.List() {
		role, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if fmt.Sprint(role["role"]) == "" {
			return fmt.Errorf("every roles block requires a non-empty \"role\"")
		}
		if fmt.Sprint(role["db_name"]) == "" {
			return fmt.Errorf("every roles block requires a non-empty \"db_name\"")
		}
	}
	return nil
}

func accountRolesFromResponse(raw interface{}) []interface{} {
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
