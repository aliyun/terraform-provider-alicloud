// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudPolarDbDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudPolarDbDatabaseCreate,
		Read:   resourceAliCloudPolarDbDatabaseRead,
		Update: resourceAliCloudPolarDbDatabaseUpdate,
		Delete: resourceAliCloudPolarDbDatabaseDelete,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudPolarDbDatabaseCreate(d *schema.ResourceData, meta interface{}) error {

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
	if v, ok := d.GetOk("db_description"); ok {
		request["DBDescription"] = v
	}
	if v, ok := d.GetOk("character_set_name"); ok {
		request["CharacterSetName"] = v
	} else {
		request["CharacterSetName"] = "utf8"
	}

	polarDBService := PolarDBService{client}
	cluster, err := polarDBService.DescribePolarDBCluster(fmt.Sprint(request["DBClusterId"]))
	if err != nil {
		return WrapError(err)
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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "ServiceUnavailable", "InternalError", "LockTimeout"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_polardb_database", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["DBName"]))

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, polarDbServiceV2.PolarDbDatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudPolarDbDatabaseRead(d, meta)
}

func resourceAliCloudPolarDbDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	polarDbServiceV2 := PolarDbServiceV2{client}

	objectRaw, err := polarDbServiceV2.DescribePolarDbDatabase(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_polardb_database DescribePolarDbDatabase Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("character_set_name", objectRaw["CharacterSetName"])
	d.Set("db_description", objectRaw["DBDescription"])
	d.Set("status", objectRaw["DBStatus"])
	d.Set("db_name", objectRaw["DBName"])

	e := jsonata.MustCompile("$.Accounts.Account[0].AccountName")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("account_name", evaluation)

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])

	return nil
}

func resourceAliCloudPolarDbDatabaseUpdate(d *schema.ResourceData, meta interface{}) error {
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
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"Connect.Timeout"}) || NeedRetry(err) {
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
	objectRaw, _ = polarDbServiceV2.DescribePolarDbDatabase(d.Id())
	enableGrantAccountPrivilege1 := false
	checkValue00 := objectRaw["Engine"]

	if InArray(fmt.Sprint(checkValue00), []string{"PostgreSQL", "Oracle"}) {
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
	return resourceAliCloudPolarDbDatabaseRead(d, meta)
}

func resourceAliCloudPolarDbDatabaseDelete(d *schema.ResourceData, meta interface{}) error {

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

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("polardb", "2017-08-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDBInstanceState", "OperationDenied.OutofUsage", "InstanceConnectTimeoutFault", "OperationDenied.DBInstanceStatus", "ConcurrentTaskExceeded", "OperationDenied.DBClusterStatus", "OperationDenied.DBStatus", "Database.ConnectError", "ServiceUnavailable", "InternalError", "LockTimeout"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBClusterId.NotFound", "InvalidDBName.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	polarDbServiceV2 := PolarDbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{""}, d.Timeout(schema.TimeoutDelete), 5*time.Second, polarDbServiceV2.PolarDbDatabaseStateRefreshFunc(d.Id(), "DBStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
