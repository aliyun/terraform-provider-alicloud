// Package alicloud. This file was originally generated and is now maintained by hand.
// See the file history for its generation provenance.
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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Update: resourceDatabaseUpdate,
		Delete: resourceDatabaseDelete,
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			// CreateDatabase applies AccountPrivilege only when AccountName is set.
			// Reject the silently ineffective combination at plan time instead of
			// accepting it. Unknown values are left to apply-time evaluation.
			raw := d.GetRawConfig()
			if raw == cty.NilVal || raw.IsNull() || !raw.IsKnown() {
				return nil
			}
			privilege := raw.GetAttr("account_privilege")
			accountName := raw.GetAttr("account_name")
			if privilege.IsNull() || !privilege.IsKnown() || privilege.AsString() == "" {
				return nil
			}
			if accountName.IsNull() || (accountName.IsKnown() && accountName.AsString() == "") {
				return fmt.Errorf("\"account_privilege\" requires \"account_name\" to be set: CreateDatabase only applies AccountPrivilege together with AccountName")
			}
			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_privilege": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"character_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"collate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ctype": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	// CustomizeDiff rejects the known combinations; an account_name that is
	// unknown at plan time can still resolve to an empty string at apply time.
	// CreateDatabase silently ignores AccountPrivilege without AccountName, so
	// repeat the check before the request is sent.
	if privilege, ok := d.GetOk("account_privilege"); ok && strings.TrimSpace(fmt.Sprint(privilege)) != "" {
		if accountName, hasAccountName := d.GetOk("account_name"); !hasAccountName || strings.TrimSpace(fmt.Sprint(accountName)) == "" {
			return fmt.Errorf("\"account_privilege\" requires \"account_name\": CreateDatabase only applies AccountPrivilege together with AccountName")
		}
	}

	action := "CreateDatabase"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("db_cluster_id"); ok {
		request["DBClusterId"] = v
	}
	if v, ok := d.GetOk("db_name"); ok {
		request["DBName"] = v
	}

	if v, ok := d.GetOk("account_name"); ok {
		request["AccountName"] = v
	}
	if v, ok := d.GetOk("account_privilege"); ok {
		request["AccountPrivilege"] = v
	}
	if v, ok := d.GetOk("db_description"); ok {
		request["DBDescription"] = v
	}
	if v, ok := d.GetOk("character_set_name"); ok {
		request["CharacterSetName"] = v
	} else {
		request["CharacterSetName"] = "utf8"
	}

	polarDBService := Service{client}
	cluster, err := polarDBService.DescribePolarDBCluster(fmt.Sprint(request["DBClusterId"]))
	if err != nil {
		return errs.WrapError(err)
	}

	if v, ok := d.GetOk("collate"); ok {
		request["Collate"] = v
	} else if cluster.DBType == "PostgreSQL" || cluster.DBType == "Oracle" {
		request["Collate"] = "C"
	}

	if v, ok := d.GetOk("ctype"); ok {
		request["Ctype"] = v
	} else if cluster.DBType == "PostgreSQL" || cluster.DBType == "Oracle" {
		request["Ctype"] = "C"
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "ServiceUnavailable", "InternalError", "LockTimeout"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_polardb_database", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["DBName"]))

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polarDbServiceV2.PolarDbDatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceDatabaseRead(d, meta)
}

func resourceDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	objectRaw, err := polarDbServiceV2.DescribePolarDbDatabase(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_database DescribePolarDbDatabase Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("character_set_name", objectRaw["CharacterSetName"])
	d.Set("db_description", objectRaw["DBDescription"])
	d.Set("engine", objectRaw["Engine"])
	d.Set("status", objectRaw["DBStatus"])
	d.Set("db_name", objectRaw["DBName"])

	refreshDatabaseAccountState(d, objectRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])

	return nil
}

// refreshDatabaseAccountState applies the configured account's identity and
// privilege to the state from a DescribeDatabases database object.
//
// DescribeDatabases returns one entry per authorized account and does not
// guarantee their order, so the account is matched by name. Reading the first
// entry could bind another account's privilege to the ForceNew
// account_privilege and trigger a spurious database replacement.
//
// An empty account_name (for example after an import) carries no identity, so
// no account is inferred from list position: importing a database cannot tell
// which account the configuration manages.
func refreshDatabaseAccountState(d *schema.ResourceData, objectRaw map[string]interface{}) {
	configuredAccountName := strings.TrimSpace(fmt.Sprint(d.Get("account_name")))
	if configuredAccountName == "" || configuredAccountName == "<nil>" {
		return
	}

	accounts, authoritative := databaseAccounts(objectRaw)
	if !authoritative {
		// A missing or malformed account block cannot prove that the configured
		// account lost its grant, so keep the recorded state instead of clearing it
		// on a transient response.
		return
	}

	account, found := selectDatabaseAccount(accounts, configuredAccountName)
	if !found {
		if d.IsNewResource() {
			// A create-time read can observe an eventually consistent empty list.
			// Keep the configured value here; the next refresh still exposes a
			// genuine removal.
			return
		}
		// The account list is authoritative and no longer contains the configured
		// account. This includes a legal empty list after the last grant is
		// revoked: clearing the privilege surfaces the ForceNew drift instead of
		// keeping a stale value.
		d.Set("account_privilege", "")
		return
	}

	d.Set("account_name", account["AccountName"])
	// Some engines (for example the PostgreSQL owner account) do not report an
	// explicit privilege in DescribeDatabases. In that case keep the configured
	// value instead of overwriting the state with an empty string.
	if privilege := strings.TrimSpace(fmt.Sprint(account["AccountPrivilege"])); privilege != "" && privilege != "<nil>" {
		d.Set("account_privilege", privilege)
	}
}

// selectDatabaseAccount returns the account entry whose AccountName equals
// accountName, together with a found flag. An empty name has no identity to
// match and never falls back to the first entry.
func selectDatabaseAccount(accounts []map[string]interface{}, accountName string) (map[string]interface{}, bool) {
	if accountName == "" {
		return nil, false
	}
	for _, record := range accounts {
		if fmt.Sprint(record["AccountName"]) == accountName {
			return record, true
		}
	}
	return nil, false
}

// databaseAccounts returns the Accounts.Account entries of a DescribeDatabases
// database object. The boolean reports whether the response is authoritative:
// the Accounts block exists and the Account entry is either absent (an empty
// container) or a list of objects. A missing or malformed block cannot prove
// that a grant was removed.
func databaseAccounts(objectRaw map[string]interface{}) ([]map[string]interface{}, bool) {
	accounts, ok := objectRaw["Accounts"].(map[string]interface{})
	if !ok {
		return nil, false
	}
	rawAccounts, exists := accounts["Account"]
	if !exists || rawAccounts == nil {
		return nil, true
	}
	accountList, ok := rawAccounts.([]interface{})
	if !ok {
		return nil, false
	}
	records := make([]map[string]interface{}, 0, len(accountList))
	for _, item := range accountList {
		record, ok := item.(map[string]interface{})
		if !ok {
			return nil, false
		}
		records = append(records, record)
	}
	return records, true
}

func resourceDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	polarDbServiceV2 := PolarDbServiceV2{client}
	objectRaw, _ := polarDbServiceV2.DescribePolarDbDatabase(d.Id())

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyDBDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]

	if d.HasChange("db_description") {
		update = true
	}
	request["DBDescription"] = d.Get("db_description")
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"Connect.Timeout"}) || helper.NeedRetry(err) {
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
	objectRaw, _ = polarDbServiceV2.DescribePolarDbDatabase(d.Id())
	enableGrantAccountPrivilege1 := false
	checkValue00 := objectRaw["Engine"]

	if helper.InArray(fmt.Sprint(checkValue00), []string{"PostgreSQL", "Oracle"}) {
		enableGrantAccountPrivilege1 = true
	}
	parts = strings.Split(d.Id(), ":")
	action = "GrantAccountPrivilege"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]

	if d.HasChange("account_name") {
		update = true
	}
	request["AccountName"] = d.Get("account_name")
	// The CloudSpec update mapping pins AccountPrivilege to the DBOwner constant
	// (requestConstValue), so GrantAccountPrivilege keeps the shipped owner-grant
	// semantics; user-configured account_privilege is create-only (ForceNew).
	request["AccountPrivilege"] = "DBOwner"
	if update && enableGrantAccountPrivilege1 {
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

	d.Partial(false)
	return resourceDatabaseRead(d, meta)
}

func resourceDatabaseDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDatabase"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBClusterId"] = parts[0]
	request["DBName"] = parts[1]

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "ServiceUnavailable", "InternalError", "LockTimeout"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidDBName.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := helper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, polarDbServiceV2.PolarDbDatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
