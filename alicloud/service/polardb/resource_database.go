// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package polardb

import (
	tferrs "github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	tflogs "github.com/aliyun/terraform-provider-alicloud/alicloud/logs"

	tfhelper "github.com/aliyun/terraform-provider-alicloud/alicloud/helper"

	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfretry "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	jsonata "github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseCreate,
		Read:   resourceDatabaseRead,
		Update: resourceDatabaseUpdate,
		Delete: resourceDatabaseDelete,
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
	// AccountPrivilege is intentionally NOT set on CreateDatabase. Per OpenAPI,
	// AccountPrivilege accepts ReadWrite/ReadOnly/DMLOnly/DDLOnly/ReadIndex (no
	// DBOwner) and defaults to ReadWrite; the DB Owner is determined by AccountName
	// ("若指定了AccountName，则创建的DB的Owner为该账号"), not by AccountPrivilege=DBOwner.
	// The CloudSpec model const=DBOwner for CreateDatabase.AccountPrivilege conflicts
	// with the API contract (ACC verified MySQL rejects DBOwner with
	// InvalidAccountPrivilege.Malformed); this gap is surfaced to model maintainers.
	// Baseline GrantAccountPrivilege DBOwner for PG/Oracle (Update) is preserved.
	if v, ok := d.GetOk("character_set_name"); ok {
		request["CharacterSetName"] = v
	} else {
		request["CharacterSetName"] = "utf8"
	}
	cluster, err := (&PolardbServiceV2{client}).DescribePolardbCluster(fmt.Sprint(request["DBClusterId"]))
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
	err = tfretry.Retry(d.Timeout(schema.TimeoutCreate), func() *tfretry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"ConcurrentTaskExceeded", "Database.ConnectError", "IncorrectDBInstanceState", "InstanceConnectTimeoutFault", "InternalError", "LockTimeout", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage", "ServiceUnavailable"}) || tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
		}
		return nil
	})
	tflogs.AddDebug(action, response, request)

	if err != nil {
		return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, "alicloud_polardb_database", action, tferrs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["DBName"]))

	polardbServiceV2 := PolardbServiceV2{client}
	stateConf := tfhelper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polardbServiceV2.PolardbDatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return tferrs.WrapErrorf(err, tferrs.IdMsg, d.Id())
	}

	return resourceDatabaseRead(d, meta)
}

func resourceDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polardbServiceV2 := PolardbServiceV2{client}

	objectRaw, err := polardbServiceV2.DescribePolardbDatabase(d.Id())
	if err != nil {
		if !d.IsNewResource() && tferrs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_database DescribePolardbDatabase Failed!!! %s", err)
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

	e := jsonata.MustCompile("$.Accounts.Account[0].AccountName")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("account_name", evaluation)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])

	return nil
}

func resourceDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	polardbServiceV2 := PolardbServiceV2{client}
	objectRaw, _ := polardbServiceV2.DescribePolardbDatabase(d.Id())

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
		err = tfretry.Retry(d.Timeout(schema.TimeoutUpdate), func() *tfretry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if tferrs.IsExpectedErrors(err, []string{"Connect.Timeout"}) || tfhelper.NeedRetry(err) {
					wait()
					return tfretry.RetryableError(err)
				}
				return tfretry.NonRetryableError(err)
			}
			return nil
		})
		tflogs.AddDebug(action, response, request)
		if err != nil {
			return tferrs.WrapErrorf(err, tferrs.DefaultErrorMsg, d.Id(), action, tferrs.AlibabaCloudSdkGoERROR)
		}
	}
	update = false
	objectRaw, _ = polardbServiceV2.DescribePolardbDatabase(d.Id())
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
	request["AccountPrivilege"] = "DBOwner"
	if update && enableGrantAccountPrivilege1 {
		wait := tfhelper.IncrementalWait(10*time.Second, 10*time.Second)
		err = tfretry.Retry(d.Timeout(schema.TimeoutUpdate), func() *tfretry.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if tferrs.IsExpectedErrors(err, []string{"ConcurrentTaskExceeded", "Database.ConnectError", "IncorrectDBInstanceState", "InstanceConnectTimeoutFault", "LockTimeout", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage"}) || tfhelper.NeedRetry(err) {
					wait()
					return tfretry.RetryableError(err)
				}
				return tfretry.NonRetryableError(err)
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
	err = tfretry.Retry(d.Timeout(schema.TimeoutDelete), func() *tfretry.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if tferrs.IsExpectedErrors(err, []string{"ConcurrentTaskExceeded", "Database.ConnectError", "IncorrectDBInstanceState", "InstanceConnectTimeoutFault", "InternalError", "LockTimeout", "OperationDenied.DBClusterStatus", "OperationDenied.DBInstanceStatus", "OperationDenied.DBStatus", "OperationDenied.OutofUsage", "ServiceUnavailable"}) || tfhelper.NeedRetry(err) {
				wait()
				return tfretry.RetryableError(err)
			}
			return tfretry.NonRetryableError(err)
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

	polardbServiceV2 := PolardbServiceV2{client}
	stateConf := tfhelper.BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, polardbServiceV2.PolardbDatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return tferrs.WrapErrorf(err, tferrs.IdMsg, d.Id())
	}

	return nil
}
