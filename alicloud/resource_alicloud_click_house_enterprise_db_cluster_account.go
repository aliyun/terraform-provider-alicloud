// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudClickHouseEnterpriseDBClusterAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudClickHouseEnterpriseDBClusterAccountCreate,
		Read:   resourceAliCloudClickHouseEnterpriseDBClusterAccountRead,
		Update: resourceAliCloudClickHouseEnterpriseDBClusterAccountUpdate,
		Delete: resourceAliCloudClickHouseEnterpriseDBClusterAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dml_auth_setting": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ddl_authority": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"dml_authority": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"allow_dictionaries": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allow_databases": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceAliCloudClickHouseEnterpriseDBClusterAccountCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("account"); ok {
		request["Account"] = v
	}
	if v, ok := d.GetOk("db_instance_id"); ok {
		request["DBInstanceId"] = v
	}
	request["RegionId"] = client.RegionId

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["Password"] = d.Get("password")
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("dml_auth_setting"); !IsNil(v) {
		allowDatabases1, _ := jsonpath.Get("$[0].allow_databases", v)
		if allowDatabases1 != nil && allowDatabases1 != "" {
			objectDataLocalMap["AllowDatabases"] = allowDatabases1
		}
		dmlAuthority1, _ := jsonpath.Get("$[0].dml_authority", v)
		if dmlAuthority1 != nil && dmlAuthority1 != "" {
			objectDataLocalMap["DmlAuthority"] = dmlAuthority1
		}
		ddlAuthority1, _ := jsonpath.Get("$[0].ddl_authority", v)
		if ddlAuthority1 != nil && ddlAuthority1 != "" {
			objectDataLocalMap["DdlAuthority"] = ddlAuthority1
		}
		allowDictionaries1, _ := jsonpath.Get("$[0].allow_dictionaries", v)
		if allowDictionaries1 != nil && allowDictionaries1 != "" {
			objectDataLocalMap["AllowDictionaries"] = allowDictionaries1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["DmlAuthSetting"] = string(objectDataLocalMapJson)
	}

	request["AccountType"] = d.Get("account_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_enterprise_db_cluster_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBInstanceId"], request["Account"]))

	return resourceAliCloudClickHouseEnterpriseDBClusterAccountRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDBClusterAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickHouseServiceV2 := ClickHouseServiceV2{client}

	objectRaw, err := clickHouseServiceV2.DescribeClickHouseEnterpriseDBClusterAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_enterprise_db_cluster_account DescribeClickHouseEnterpriseDBClusterAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	dmlAuthSettingMaps := make([]map[string]interface{}, 0)
	dmlAuthSettingMap := make(map[string]interface{})

	dmlAuthSettingMap["ddl_authority"] = objectRaw["DdlAuthority"]
	dmlAuthSettingMap["dml_authority"] = objectRaw["DmlAuthority"]

	allowDatabasesRaw, _ := jsonpath.Get("$.AllowDatabases", objectRaw)
	if allowDatabasesRaw != nil {
		allowDatabasesFiltered := make([]interface{}, 0)
		for _, allowDatabase := range allowDatabasesRaw.([]interface{}) {
			if fmt.Sprint(allowDatabase) != "system" {
				allowDatabasesFiltered = append(allowDatabasesFiltered, allowDatabase)
			}
		}
		dmlAuthSettingMap["allow_databases"] = allowDatabasesFiltered
	}
	allowDictionariesRaw, _ := jsonpath.Get("$.AllowDictionaries", objectRaw)
	dmlAuthSettingMap["allow_dictionaries"] = allowDictionariesRaw
	dmlAuthSettingMaps = append(dmlAuthSettingMaps, dmlAuthSettingMap)
	if err := d.Set("dml_auth_setting", dmlAuthSettingMaps); err != nil {
		return err
	}

	objectRaw, err = clickHouseServiceV2.DescribeEnterpriseDBClusterAccountDescribeAccounts(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("account_type", convertClickHouseEnterpriseDBClusterAccountDataAccountsAccountTypeResponse(objectRaw["AccountType"]))
	d.Set("description", objectRaw["Description"])
	d.Set("account", objectRaw["Account"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_instance_id", parts[0])

	return nil
}

func resourceAliCloudClickHouseEnterpriseDBClusterAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyAccountAuthority"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Account"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("dml_auth_setting") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("dml_auth_setting"); v != nil {
		allowDatabases1, _ := jsonpath.Get("$[0].allow_databases", d.Get("dml_auth_setting"))
		if allowDatabases1 != nil && (d.HasChange("dml_auth_setting.0.allow_databases") || allowDatabases1 != "") {
			objectDataLocalMap["AllowDatabases"] = allowDatabases1
		}
		dmlAuthority1, _ := jsonpath.Get("$[0].dml_authority", v)
		if dmlAuthority1 != nil && (d.HasChange("dml_auth_setting.0.dml_authority") || dmlAuthority1 != "") {
			objectDataLocalMap["DmlAuthority"] = dmlAuthority1
		}
		ddlAuthority1, _ := jsonpath.Get("$[0].ddl_authority", v)
		if ddlAuthority1 != nil && (d.HasChange("dml_auth_setting.0.ddl_authority") || ddlAuthority1 != "") {
			objectDataLocalMap["DdlAuthority"] = ddlAuthority1
		}
		allowDictionaries1, _ := jsonpath.Get("$[0].allow_dictionaries", d.Get("dml_auth_setting"))
		if allowDictionaries1 != nil && (d.HasChange("dml_auth_setting.0.allow_dictionaries") || allowDictionaries1 != "") {
			objectDataLocalMap["AllowDictionaries"] = allowDictionaries1
		}

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		request["DmlAuthSetting"] = string(objectDataLocalMapJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ModifyAccountDescription"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Account"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
	}
	request["Description"] = d.Get("description")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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
	update = false
	parts = strings.Split(d.Id(), ":")
	action = "ResetAccountPassword"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["Account"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("password") {
		update = true
	}
	request["Password"] = d.Get("password")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)
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

	d.Partial(false)
	return resourceAliCloudClickHouseEnterpriseDBClusterAccountRead(d, meta)
}

func resourceAliCloudClickHouseEnterpriseDBClusterAccountDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteAccount"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Account"] = parts[1]
	request["DBInstanceId"] = parts[0]
	request["RegionId"] = client.RegionId

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2023-05-22", action, query, request, true)

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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertClickHouseEnterpriseDBClusterAccountDataAccountsAccountTypeResponse(source interface{}) interface{} {
	source = fmt.Sprint(source)
	switch source {
	case "1":
		return "NormalAccount"
	case "6":
		return "SuperAccount"
	}
	return source
}
