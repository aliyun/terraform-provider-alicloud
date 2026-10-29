// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Update: resourceDatabaseUpdate,
		Delete: resourceDatabaseDelete,
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			// CreateDatabase applies AccountPrivilege only when AccountName is set
			// (the account is created as the database owner). Reject the silently
			// ineffective combination at plan time instead of accepting it.
			// Unknown values are left to apply-time evaluation.
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
	if v, ok := d.GetOk("character_set_name"); ok {
		request["CharacterSetName"] = v
	} else {
		request["CharacterSetName"] = "utf8"
	}
	cluster, err := (&Service{client}).DescribeCluster(fmt.Sprint(request["DBClusterId"]))
	if err != nil {
		return tferrs.WrapError(err)
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
	if v, ok := d.GetOk("db_description"); ok {
		request["DBDescription"] = v
	}
	wait := tfhelper.IncrementalWait(10*time.Second, 10*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"ConcurrentTaskExceeded", "Database.ConnectError", "IncorrectDBInstanceState", "InstanceConnectTimeoutFault", "InternalError", "LockTimeout", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage", "ServiceUnavailable"}) || tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_polardb_database", action, tferrs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["DBName"]))

	polardbServiceV2 := Service{client}
	stateConf := tfhelper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polardbServiceV2.DatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return tferrs.WrapErrorf(err, tferrs.IdMsg, d.Id())
	}

	return resourceDatabaseRead(d, meta)
}

func resourceDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polardbServiceV2 := Service{client}

	objectRaw, err := polardbServiceV2.DescribeDatabase(d.Id())
	if err != nil {
		if !d.IsNewResource() && tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_database DescribeDatabase Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return tferrs.WrapError(err)
	}

	d.Set("character_set_name", objectRaw["CharacterSetName"])
	d.Set("db_description", objectRaw["DBDescription"])
	d.Set("engine", objectRaw["Engine"])
	d.Set("status", objectRaw["DBStatus"])
	d.Set("db_name", objectRaw["DBName"])

	// DescribeDatabases returns one entry per authorized account and does not
	// guarantee their order, so the account is matched by name. Reading the first
	// entry could bind another account's privilege to the ForceNew
	// account_privilege and trigger a spurious database replacement. When no
	// account is configured (for example after an import) the first entry is used,
	// which is the only account information available.
	matchedAccount := selectDatabaseAccount(objectRaw, fmt.Sprint(d.Get("account_name")))
	if matchedAccount != nil {
		d.Set("account_name", matchedAccount["AccountName"])
		// Some engines (for example the PostgreSQL owner account) do not report an
		// explicit privilege in DescribeDatabases. In that case keep the configured
		// value instead of overwriting the state with an empty string.
		if privilege := strings.TrimSpace(fmt.Sprint(matchedAccount["AccountPrivilege"])); privilege != "" && privilege != "<nil>" {
			d.Set("account_privilege", privilege)
		}
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])

	return nil
}

// selectDatabaseAccount returns the Accounts.Account entry that belongs to the
// given account name. An empty accountName (no account configured, for example
// after an import) falls back to the first entry, which is the only account
// information the API response offers in that case.
func selectDatabaseAccount(objectRaw map[string]interface{}, accountName string) map[string]interface{} {
	accounts, ok := objectRaw["Accounts"].(map[string]interface{})
	if !ok {
		return nil
	}
	accountList, ok := accounts["Account"].([]interface{})
	if !ok {
		return nil
	}
	for _, item := range accountList {
		record, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if accountName == "" || fmt.Sprint(record["AccountName"]) == accountName {
			return record
		}
	}
	return nil
}

func resourceDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	polardbServiceV2 := Service{client}
	objectRaw, _ := polardbServiceV2.DescribeDatabase(d.Id())

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
		wait := tfhelper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if tferrs.IsExpectedErrors(err, []string{"Connect.Timeout"}) || tfhelper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		tflogs.AddDebug(action, response, request)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	objectRaw, _ = polardbServiceV2.DescribeDatabase(d.Id())
	enableGrantAccountPrivilege1 := false
	checkValue00 := objectRaw["Engine"]
	if tfhelper.InArray(fmt.Sprint(checkValue00), []string{"PostgreSQL", "Oracle"}) {
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
		wait := tfhelper.IncrementalWait(10*time.Second, 10*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if tferrs.IsExpectedErrors(err, []string{"ConcurrentTaskExceeded", "Database.ConnectError", "IncorrectDBInstanceState", "InstanceConnectTimeoutFault", "LockTimeout", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage"}) || tfhelper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		tflogs.AddDebug(action, response, request)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
		}
	}

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

	wait := tfhelper.IncrementalWait(10*time.Second, 10*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"ConcurrentTaskExceeded", "Database.ConnectError", "IncorrectDBInstanceState", "InstanceConnectTimeoutFault", "InternalError", "LockTimeout", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage", "ServiceUnavailable"}) || tfhelper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		if tferrs.IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidDBName.NotFound"}) || tferrs.NotFoundError(err) {
			return nil
		}
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
	}

	polardbServiceV2 := Service{client}
	stateConf := tfhelper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, polardbServiceV2.DatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return tferrs.WrapErrorf(err, tferrs.IdMsg, d.Id())
	}

	return nil
}
