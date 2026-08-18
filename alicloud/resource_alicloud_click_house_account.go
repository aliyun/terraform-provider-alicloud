package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// clickHouseEngineVersionAfter218 reports whether the given ClickHouse cluster
// EngineVersion (e.g. "21.8.10.19", "22.8.5.29") is strictly greater than 21.8.
// ModifyAccountAuthority only applies to 21.8 and earlier cluster versions, so
// the authority fields (dml_authority, ddl_authority, allow_databases and
// allow_dictionaries) must not be sent to clusters newer than 21.8.
func clickHouseEngineVersionAfter218(engineVersion string) bool {
	parts := strings.Split(engineVersion, ".")
	if len(parts) < 2 {
		return false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return false
	}
	return major > 21 || (major == 21 && minor > 8)
}

func resourceAlicloudClickHouseAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudClickHouseAccountCreate,
		Read:   resourceAlicloudClickHouseAccountRead,
		Update: resourceAlicloudClickHouseAccountUpdate,
		Delete: resourceAlicloudClickHouseAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"account_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`^[a-z][a-z0-9_]{1,15}`), "The account_name most consist of lowercase letters, numbers, and underscores, starting with a lowercase letter"),
			},
			"account_password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringMatch(regexp.MustCompile(`[a-zA-Z!#$%^&*()_+-=]{8,32}`), "account_password must consist of uppercase letters, lowercase letters, numbers, and special characters"),
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
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"Normal", "Super"}, false),
			},
			"dml_authority": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"all", "readOnly,modify"}, false),
				Deprecated:   "Field 'dml_authority' has been deprecated from version 1.290.0. ClickHouse clusters newer than 21.8 no longer support ModifyAccountAuthority; use the SQL account authority model instead.",
			},
			"ddl_authority": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'ddl_authority' has been deprecated from version 1.290.0. ClickHouse clusters newer than 21.8 no longer support ModifyAccountAuthority; use the SQL account authority model instead.",
			},
			"allow_databases": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'allow_databases' has been deprecated from version 1.290.0. ClickHouse clusters newer than 21.8 no longer support ModifyAccountAuthority; use the SQL account authority model instead.",
			},
			"total_databases": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'total_databases' has been deprecated from version 1.223.1 and it will be removed in the future version.",
			},
			"allow_dictionaries": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'allow_dictionaries' has been deprecated from version 1.290.0. ClickHouse clusters newer than 21.8 no longer support ModifyAccountAuthority; use the SQL account authority model instead.",
			},
			"total_dictionaries": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'total_dictionaries' has been deprecated from version 1.223.1 and it will be removed in the future version.",
			},
		},
	}
}

func resourceAlicloudClickHouseAccountCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	var response map[string]interface{}
	var err error

	dbClusterId := d.Get("db_cluster_id").(string)
	// DescribeDBClusterAttribute 读取集群 EngineVersion，用于按版本兼容切换
	// CreateAccount/CreateSQLAccount 以及权限字段防写。
	cluster, err := clickhouseService.DescribeClickHouseDbCluster(dbClusterId)
	if err != nil {
		return WrapError(err)
	}
	engineVersion, _ := cluster["EngineVersion"].(string)
	after218 := clickHouseEngineVersionAfter218(engineVersion)

	// ModifyAccountAuthority 仅适用于 21.8 及以下集群；高版本集群设置权限字段前置报错防写。
	if after218 {
		for _, field := range []string{"dml_authority", "ddl_authority", "allow_databases", "allow_dictionaries"} {
			var set bool
			if field == "ddl_authority" {
				_, set = d.GetOkExists(field)
			} else {
				_, set = d.GetOk(field)
			}
			if set {
				return WrapError(fmt.Errorf("setting %s is not supported on ClickHouse clusters newer than 21.8 (EngineVersion %q); ModifyAccountAuthority only applies to 21.8 and earlier versions", field, engineVersion))
			}
		}
	}

	action := "CreateAccount"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	request["AccountName"] = d.Get("account_name")
	request["AccountPassword"] = d.Get("account_password")
	request["DBClusterId"] = dbClusterId
	// 21.8 以上集群用 CreateSQLAccount 创建账号（Normal 与 Super 均走该 API）；
	// 21.8 及以下维持原逻辑：Super 走 CreateSQLAccount，Normal 走 CreateAccount。
	if after218 || d.Get("type") == "Super" {
		action = "CreateSQLAccount"
		request["AccountType"] = d.Get("type")
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectAccountStatus", "IncorrectDBInstanceState"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_click_house_account", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBClusterId"], ":", request["AccountName"]))

	return resourceAlicloudClickHouseAccountUpdate(d, meta)
}
func resourceAlicloudClickHouseAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	object, err := clickhouseService.DescribeClickHouseAccount(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_click_house_account clickhouseService.DescribeClickHouseAccount Failed!!! %s", err)
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
	d.Set("db_cluster_id", parts[0])
	d.Set("account_description", object["AccountDescription"])
	d.Set("status", object["AccountStatus"])
	d.Set("type", object["AccountType"])

	authority, err := clickhouseService.DescribeClickHouseAccountAuthority(d.Id())
	d.Set("dml_authority", authority["DmlAuthority"])
	d.Set("ddl_authority", authority["DdlAuthority"])

	d.Set("allow_databases", convertArrayToString(authority["AllowDatabases"], ","))
	d.Set("allow_dictionaries", convertArrayToString(authority["AllowDictionaries"], ","))
	d.Set("total_databases", convertArrayToString(authority["TotalDatabases"], ","))
	d.Set("total_dictionaries", convertArrayToString(authority["TotalDictionaries"], ","))
	if err != nil {
		return WrapError(err)
	}

	return nil
}
func resourceAlicloudClickHouseAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	update := false
	d.Partial(true)
	request := map[string]interface{}{
		"AccountName": parts[1],
		"DBClusterId": parts[0],
	}
	if !d.IsNewResource() && d.HasChange("account_description") {
		update = true
	}
	if v, ok := d.GetOk("account_description"); ok {
		request["AccountDescription"] = v
	}
	if update {
		action := "ModifyAccountDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectAccountStatus", "IncorrectDBInstanceState"}) || NeedRetry(err) {
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
		d.SetPartial("account_description")
	}
	update = false
	request = map[string]interface{}{
		"AccountName": parts[1],
		"DBClusterId": parts[0],
	}
	if !d.IsNewResource() && d.HasChange("account_password") {
		update = true
	}
	request["AccountPassword"] = d.Get("account_password")
	if update {
		action := "ResetAccountPassword"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectAccountStatus", "IncorrectDBInstanceState"}) || NeedRetry(err) {
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
		d.SetPartial("account_password")
	}

	update = false
	request = map[string]interface{}{
		"AccountName": parts[1],
		"DBClusterId": parts[0],
	}
	request["RegionId"] = client.RegionId
	// 权限字段(dml_authority/ddl_authority/allow_databases/allow_dictionaries)变更前
	// 确认集群版本，ModifyAccountAuthority 仅适用于 21.8 及以下集群，>21.8 前置报错防写。
	if d.HasChange("dml_authority") || d.HasChange("ddl_authority") || d.HasChange("allow_databases") || d.HasChange("allow_dictionaries") {
		cluster, err := clickhouseService.DescribeClickHouseDbCluster(parts[0])
		if err != nil {
			return WrapError(err)
		}
		engineVersion, _ := cluster["EngineVersion"].(string)
		if clickHouseEngineVersionAfter218(engineVersion) {
			return WrapError(fmt.Errorf("modifying account authority is not supported on ClickHouse clusters newer than 21.8 (EngineVersion %q); ModifyAccountAuthority only applies to 21.8 and earlier versions", engineVersion))
		}
	}
	if d.HasChange("dml_authority") {
		update = true
	}
	if v, ok := d.GetOk("dml_authority"); ok {
		request["DmlAuthority"] = v
	}

	if d.HasChange("ddl_authority") {
		update = true
	}
	if v, ok := d.GetOkExists("ddl_authority"); ok {
		request["DdlAuthority"] = v
	}

	if d.HasChange("allow_databases") {
		update = true
	}
	if v, ok := d.GetOk("allow_databases"); ok {
		request["AllowDatabases"] = v
	}

	if d.HasChange("total_databases") {
		update = true
	}
	if v, ok := d.GetOk("total_databases"); ok {
		request["TotalDatabases"] = v
	}

	if d.HasChange("allow_dictionaries") {
		update = true
	}
	if v, ok := d.GetOk("allow_dictionaries"); ok {
		request["AllowDictionaries"] = v
	}

	if d.HasChange("total_dictionaries") {
		update = true
	}
	if v, ok := d.GetOk("total_dictionaries"); ok {
		request["TotalDictionaries"] = v
	}

	if update {
		action := "ModifyAccountAuthority"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"IncorrectAccountStatus", "IncorrectDBInstanceState"}) || NeedRetry(err) {
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
		d.SetPartial("account_password")
	}
	d.Partial(false)

	return resourceAlicloudClickHouseAccountRead(d, meta)
}
func resourceAlicloudClickHouseAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	clickhouseService := ClickhouseService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteAccount"
	var response map[string]interface{}
	request := map[string]interface{}{
		"AccountName": parts[1],
		"DBClusterId": parts[0],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("clickhouse", "2019-11-11", action, nil, request, false)
		if err != nil {
			// InstanceConnectFailed is a transient backend-side failure to reach
			// the ClickHouse node (HTTP 400 with a connectivity message). It is
			// not a validation error and typically clears on its own, so retry
			// it like the account/instance-state errors above instead of
			// surfacing it to the user on the first attempt.
			if IsExpectedErrors(err, []string{"IncorrectAccountStatus", "IncorrectDBInstanceState", "InstanceConnectFailed"}) || NeedRetry(err) {
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
	stateConf := BuildStateConf([]string{"Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, clickhouseService.ClickhouseStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
