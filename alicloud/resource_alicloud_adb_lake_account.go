package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudAdbLakeAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAdbLakeAccountCreate,
		Read:   resourceAliCloudAdbLakeAccountRead,
		Update: resourceAliCloudAdbLakeAccountUpdate,
		Delete: resourceAliCloudAdbLakeAccountDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAdbLakeAccountImportState,
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
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"AnalyticDB", "Clickhouse"}, false),
			},
			"ram_user_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// resourceAdbLakeAccountImportState parses the import id and, when an engine is
// supplied, records it so that Read/Delete target the correct engine.
//
// DescribeAccounts accepts an Engine filter whose API default is AnalyticDB.
// An account created on the Clickhouse engine cannot be located (and therefore
// cannot be imported, refreshed or deleted) without passing Engine=Clickhouse,
// because the default-engine query does not return it. The legacy two-segment id
// `<db_cluster_id>:<account_name>` is retained for backwards compatibility and
// resolves to the default engine (AnalyticDB); Clickhouse accounts must be
// imported with the three-segment form `<db_cluster_id>:<account_name>:<engine>`,
// which also disambiguates same-named accounts across engines.
func resourceAdbLakeAccountImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 3)
	switch len(parts) {
	case 3:
		// <db_cluster_id>:<account_name>:<engine>
		engine := parts[2]
		if engine != "AnalyticDB" && engine != "Clickhouse" {
			return nil, WrapError(fmt.Errorf("invalid import id %q: engine segment must be AnalyticDB or Clickhouse", d.Id()))
		}
		d.Set("engine", engine)
		d.SetId(fmt.Sprintf("%s:%s", parts[0], parts[1]))
	case 2:
		// legacy <db_cluster_id>:<account_name>; engine is unknown so Read queries
		// DescribeAccounts without an Engine filter (API default AnalyticDB).
		if parts[0] == "" || parts[1] == "" {
			return nil, WrapError(fmt.Errorf("invalid import id %q: expected <db_cluster_id>:<account_name>[:<engine>]", d.Id()))
		}
	default:
		return nil, WrapError(fmt.Errorf("invalid import id %q: expected <db_cluster_id>:<account_name>[:<engine>]", d.Id()))
	}
	return []*schema.ResourceData{d}, nil
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

	if v, ok := d.GetOk("ram_user_list"); ok {
		ramUserListMapsArray := convertToInterfaceArray(v)

		ramUserListMapsJson, err := json.Marshal(ramUserListMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["RamUserList"] = string(ramUserListMapsJson)
	}

	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}
	request["AccountPassword"] = d.Get("account_password")
	request["AccountType"] = d.Get("account_type")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_lake_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["DBClusterId"], request["AccountName"]))

	adbServiceV2 := AdbServiceV2{client}
	// The create waiter queries DescribeAccounts, which accepts an Engine
	// filter. Pass the engine selected at create time so Clickhouse accounts are
	// polled on the correct engine instead of the API default (AnalyticDB).
	engine := d.Get("engine").(string)
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, adbServiceV2.AdbLakeAccountStateRefreshFuncWithApi(d.Id(), "AccountStatus", []string{}, func(id string) (map[string]interface{}, error) {
		return adbServiceV2.DescribeLakeAccountDescribeAccounts(id, engine)
	}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAdbLakeAccountUpdate(d, meta)
}

func resourceAliCloudAdbLakeAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbServiceV2 := AdbServiceV2{client}

	// DescribeAccounts accepts an Engine filter and is the authoritative
	// identity check for the (db_cluster_id, account_name, engine) tuple.
	// Query it first so that a NotFound on the specified engine removes the
	// resource from state (refresh) or fails the import, instead of silently
	// writing empty attributes. A three-segment import
	// (<db_cluster_id>:<account_name>:<engine>) populates engine before Read,
	// so Clickhouse accounts are queried on the correct engine. A legacy
	// two-segment id leaves engine empty, so DescribeAccounts uses the API
	// default (AnalyticDB) and locates the account by AccountName.
	objectRaw, err := adbServiceV2.DescribeLakeAccountDescribeAccounts(d.Id(), d.Get("engine").(string))
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_lake_account DescribeAccounts NotFound!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// DescribeAccountAllPrivileges does not accept an Engine parameter; it is
	// queried only after the engine-qualified account has been confirmed to
	// exist, so privileges are never read for an identity that has been denied.
	accountAllPrivilegesRaw, err := adbServiceV2.DescribeAdbLakeAccount(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_lake_account DescribeAdbLakeAccount Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	resultRaw, _ := jsonpath.Get("$.Data.Result", accountAllPrivilegesRaw)

	accountPrivilegesMaps := make([]map[string]interface{}, 0)
	if resultRaw != nil {
		for _, resultChildRaw := range convertToInterfaceArray(resultRaw) {
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
			privilegesRaw, _ := jsonpath.Get("$.Privileges", resultChildRaw)
			accountPrivilegesMap["privileges"] = privilegesRaw
			accountPrivilegesMaps = append(accountPrivilegesMaps, accountPrivilegesMap)
		}
	}
	if err := d.Set("account_privileges", accountPrivilegesMaps); err != nil {
		return err
	}

	d.Set("account_description", objectRaw["AccountDescription"])
	d.Set("account_type", objectRaw["AccountType"])
	d.Set("status", objectRaw["AccountStatus"])
	d.Set("account_name", objectRaw["AccountName"])
	// DescribeAccounts does not return the Engine field for accounts on the
	// default engine (AnalyticDB). Keep the value dispatched at create time
	// when the API omits it, so refresh does not produce an empty diff that
	// would force a replacement of the resource; never overwrite a known
	// engine with an empty value to avoid historical state drift.
	if v, ok := objectRaw["Engine"].(string); ok && v != "" {
		d.Set("engine", v)
	}

	ramUserListRaw, _ := jsonpath.Get("$.RamUserList.RamUserList", objectRaw)
	d.Set("ram_user_list", ramUserListRaw)

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
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}

	if !d.IsNewResource() && d.HasChange("account_description") {
		update = true
	}
	request["AccountDescription"] = d.Get("account_description")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}

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
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]
	request["RegionId"] = client.RegionId
	if d.HasChange("account_privileges") {
		update = true
	}
	if v, ok := d.GetOk("account_privileges"); ok || d.HasChange("account_privileges") {
		accountPrivilegesMapsArray := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(v) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["PrivilegeType"] = dataLoopTmp["privilege_type"]
			dataLoopMap["Privileges"] = dataLoopTmp["privileges"]
			if !IsNil(dataLoopTmp["privilege_object"]) {
				localData1 := make(map[string]interface{})
				column1, _ := jsonpath.Get("$[0].column", dataLoopTmp["privilege_object"])
				if column1 != nil && column1 != "" {
					localData1["Column"] = column1
				}
				database1, _ := jsonpath.Get("$[0].database", dataLoopTmp["privilege_object"])
				if database1 != nil && database1 != "" {
					localData1["Database"] = database1
				}
				table1, _ := jsonpath.Get("$[0].table", dataLoopTmp["privilege_object"])
				if table1 != nil && table1 != "" {
					localData1["Table"] = table1
				}
				if len(localData1) > 0 {
					dataLoopMap["PrivilegeObject"] = localData1
				}
			}
			accountPrivilegesMapsArray = append(accountPrivilegesMapsArray, dataLoopMap)
		}
		accountPrivilegesMapsJson, err := json.Marshal(accountPrivilegesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		request["AccountPrivileges"] = string(accountPrivilegesMapsJson)
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	action = "BindAccount"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["AccountName"] = parts[1]
	request["DBClusterId"] = parts[0]

	if !d.IsNewResource() && d.HasChange("ram_user_list") {
		update = true
		if v, ok := d.GetOk("ram_user_list"); ok || d.HasChange("ram_user_list") {
			ramUserListMapsArray := convertToInterfaceArray(v)

			ramUserListMapsJson, err := json.Marshal(ramUserListMapsArray)
			if err != nil {
				return WrapError(err)
			}
			request["RamUserList"] = string(ramUserListMapsJson)
		}
	}

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
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
	if v, ok := d.GetOk("engine"); ok {
		request["Engine"] = v
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, query, request, true)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
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
