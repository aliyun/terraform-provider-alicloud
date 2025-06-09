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

func resourceAliCloudAdbLakeAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAdbLakeAccountCreate,
		Read:   resourceAliCloudAdbLakeAccountRead,
		Update: resourceAliCloudAdbLakeAccountUpdate,
		Delete: resourceAliCloudAdbLakeAccountDelete,
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
			"account_privileges": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"privilege_object": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"column": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"database": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"privileges": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"privilege_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"account_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
		},
	}
}

func resourceAliCloudAdbLakeAccountCreate(d *schema.ResourceData, meta interface{}) error {

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

	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	request["AccountPassword"] = d.Get("account_password")
	request["AccountType"] = d.Get("account_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_lake_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["AccountName"]))

	adbServiceV2 := AdbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, adbServiceV2.AdbLakeAccountStateRefreshFunc(d.Id(), "AccountStatus", []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAdbLakeAccountUpdate(d, meta)
}

func resourceAliCloudAdbLakeAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbServiceV2 := AdbServiceV2{client}

	objectRaw, err := adbServiceV2.DescribeAdbLakeAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_lake_account DescribeAdbLakeAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	resultRaw, _ := jsonpath.Get("$.Data.Result", objectRaw)

	accountPrivilegesMaps := make([]map[string]interface{}, 0)
	if resultRaw != nil {
		for _, resultChildRaw := range resultRaw.([]interface{}) {
			accountPrivilegesMap := make(map[string]interface{})
			resultChildRaw := resultChildRaw.(map[string]interface{})
			accountPrivilegesMap["privilege_type"] = resultChildRaw["PrivilegeType"]

			privilegeObjectMaps := make([]map[string]interface{}, 0)
			privilegeObjectMap := make(map[string]interface{})
			privilegeObjectRawObj, _ := jsonpath.Get("$.PrivilegeObject", resultChildRaw)
			privilegeObjectRaw := make(map[string]interface{})
			if privilegeObjectRawObj != nil {
				privilegeObjectRaw = privilegeObjectRawObj.(map[string]interface{})
			}
			if len(privilegeObjectRaw) > 0 {
				privilegeObjectMap["column"] = privilegeObjectRaw["Column"]
				privilegeObjectMap["database"] = privilegeObjectRaw["Database"]
				privilegeObjectMap["table"] = privilegeObjectRaw["Table"]

				privilegeObjectMaps = append(privilegeObjectMaps, privilegeObjectMap)
			}
			accountPrivilegesMap["privilege_object"] = privilegeObjectMaps

			privileges1Raw, _ := jsonpath.Get("$.Privileges", resultChildRaw)
			accountPrivilegesMap["privileges"] = privileges1Raw
			accountPrivilegesMaps = append(accountPrivilegesMaps, accountPrivilegesMap)
		}
	}
	if err := d.Set("account_privileges", accountPrivilegesMaps); err != nil {
		return err
	}

	objectRaw, err = adbServiceV2.DescribeLakeAccountDescribeAccounts(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	d.Set("account_description", objectRaw["AccountDescription"])
	d.Set("account_type", objectRaw["AccountType"])
	d.Set("status", objectRaw["AccountStatus"])
	d.Set("account_name", objectRaw["AccountName"])

	parts := strings.Split(d.Id(), ":")
	d.Set("db_cluster_id", parts[0])

	return nil
}

func resourceAliCloudAdbLakeAccountUpdate(d *schema.ResourceData, meta interface{}) error {
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
			response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
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
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("account_description") {
		update = true
		request["AccountDescription"] = d.Get("account_description")
	}

	if !d.IsNewResource() && d.HasChange("account_password") {
		update = true
	}
	request["AccountPassword"] = d.Get("account_password")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
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
	action = "ModifyAccountPrivileges"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DBClusterId"] = parts[0]
	query["AccountName"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("account_privileges") {
		update = true
		if v, ok := d.GetOk("account_privileges"); ok {
			accountPrivilegesMaps := make([]map[string]interface{}, 0)
			for _, dataLoop := range v.(*schema.Set).List() {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["PrivilegeType"] = dataLoopTmp["privilege_type"]
				dataLoopMap["Privileges"] = dataLoopTmp["privileges"]
				privilegeObjectMap := make(map[string]interface{})
				database, _ := jsonpath.Get("$[0].database", dataLoopTmp["privilege_object"])
				privilegeObjectMap["Database"] = database
				table, _ := jsonpath.Get("$[0].table", dataLoopTmp["privilege_object"])
				privilegeObjectMap["Table"] = table
				column, _ := jsonpath.Get("$[0].column", dataLoopTmp["privilege_object"])
				privilegeObjectMap["Column"] = column
				dataLoopMap["PrivilegeObject"] = privilegeObjectMap
				accountPrivilegesMaps = append(accountPrivilegesMaps, dataLoopMap)
			}
			accountPrivilegesMapsJson, err := json.Marshal(accountPrivilegesMaps)
			if err != nil {
				return WrapError(err)
			}
			query["AccountPrivileges"] = string(accountPrivilegesMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
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
	return resourceAliCloudAdbLakeAccountRead(d, meta)
}

func resourceAliCloudAdbLakeAccountDelete(d *schema.ResourceData, meta interface{}) error {

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
		response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)

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
