package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudMongodbShardingAuditFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudMongodbShardingAuditFilterCreate,
		Read:   resourceAliCloudMongodbShardingAuditFilterRead,
		Update: resourceAliCloudMongodbShardingAuditFilterUpdate,
		Delete: resourceAliCloudMongodbShardingAuditFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"audit_status": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"enable", "disabled"}, false),
			},
			// Write-only API parameter: ModifyAuditPolicy accepts it as an optional request source and
			// none of the describe APIs return it, so Read must neither backfill nor clear an explicitly
			// configured value. No schema Default supplies it: when it is not configured the request
			// parameter is left out instead of being sent as an empty or provider-chosen value.
			"audit_log_switch_source": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"Console"}, false),
			},
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"hot_storage_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"service_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     StringInSlice([]string{"Standard", "V2_Standard"}, false),
				DiffSuppressFunc: mongodbShardingAuditFilterServiceTypeDiffSuppress,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// mongodbShardingAuditFilterServiceTypeDiffSuppress ignores service_type drift while the audit log is
// disabled: disabling the audit log resets the server-side type, which would otherwise diff against
// the declared value. Mirrors the audit_policy behaviour for the same server-side quirk.
func mongodbShardingAuditFilterServiceTypeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return d.Get("audit_status").(string) == "disabled"
}

func resourceAliCloudMongodbShardingAuditFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	// 1. ModifyAuditPolicy — audit_status / audit_log_switch_source / storage_period / hot_storage_period / service_type.
	action := "ModifyAuditPolicy"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["DBInstanceId"] = d.Get("db_instance_id")
	request["RegionId"] = client.RegionId
	// Optional request attribution: omit the parameter entirely when it is not configured rather than
	// sending an empty value.
	if v, ok := d.GetOk("audit_log_switch_source"); ok {
		request["AuditLogSwitchSource"] = v
	}
	if v, ok := d.GetOkExists("storage_period"); ok {
		request["StoragePeriod"] = v
	}
	if v, ok := d.GetOkExists("hot_storage_period"); ok {
		request["HotStoragePeriod"] = v
	}
	// service_type is Optional+Computed with no schema Default; fall back to Standard so HCL that
	// omits the field keeps the historical behaviour of enabling the Standard edition.
	serviceType := d.Get("service_type").(string)
	if serviceType == "" {
		serviceType = "Standard"
	}
	request["ServiceType"] = serviceType
	request["AuditStatus"] = convertMongodbAuditPolicyAuditStatusRequest(d.Get("audit_status").(string))
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mongodb_sharding_audit_filter", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DBInstanceId"]))

	mongodbServiceV2 := MongodbServiceV2{client}
	stateConf := BuildStateConf([]string{}, []string{"[Running]"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, mongodbServiceV2.DescribeAsyncMongodbAuditPolicyStateRefreshFunc(d, response, "$.DBInstances.DBInstance[*].DBInstanceStatus", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	// 2. ModifyAuditLogFilter — filter (required) / role_type (optional).
	action = "ModifyAuditLogFilter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	request["Filter"] = d.Get("filter")
	if v, ok := d.GetOk("role_type"); ok && v.(string) != "" {
		request["RoleType"] = v
	}
	wait = incrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
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
	stateConf = BuildStateConf([]string{}, []string{"[Running]"}, d.Timeout(schema.TimeoutCreate), 5*time.Minute, mongodbServiceV2.DescribeAsyncMongodbAuditPolicyStateRefreshFunc(d, response, "$.DBInstances.DBInstance[*].DBInstanceStatus", []string{}))
	if jobDetail, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
	}

	return resourceAliCloudMongodbShardingAuditFilterRead(d, meta)
}

func resourceAliCloudMongodbShardingAuditFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}

	objectRaw, err := mongodbServiceV2.DescribeMongodbAuditPolicy(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mongodb_sharding_audit_filter DescribeMongodbAuditPolicy Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	// DescribeMongoDBLogConfig reports the retention in edition-specific fields: the Standard edition
	// keeps it in TtlForStandard, while the V2_Standard edition reports the cold storage duration
	// (ModifyAuditPolicy's StoragePeriod for that edition) in TtlForV2Standard together with the hot
	// storage duration in HotTtlForV2Standard. Reading TtlForStandard for a V2_Standard audit log
	// returns the untouched Standard-edition default instead of the configured value.
	serviceType := objectRaw["ServiceType"]
	storagePeriod, hotStoragePeriod := mongodbShardingAuditFilterStorageValues(objectRaw, fmt.Sprint(serviceType), d.Get("service_type").(string))
	// While the audit log is disabled the server reports the Trial edition and may omit the retention
	// fields; unconditionally Set-ing them would write 0 over the configured value and produce a
	// perpetual diff on every Refresh. Keep the state value when the field is absent so re-enabling
	// can restore it (mirrors the service_type retention below).
	if storagePeriod != nil {
		d.Set("storage_period", storagePeriod)
	}
	if hotStoragePeriod != nil {
		d.Set("hot_storage_period", hotStoragePeriod)
	}

	objectRaw, err = mongodbServiceV2.DescribeAuditPolicyDescribeAuditPolicy(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	auditStatus := convertMongodbAuditPolicyLogAuditStatusResponse(objectRaw["LogAuditStatus"])
	d.Set("audit_status", auditStatus)
	// Disabling the audit log flips the server-side service_type to Trial; keep the declared value so it
	// can be restored on re-enable without a perpetual diff.
	if auditStatus != "disabled" {
		d.Set("service_type", serviceType)
	}

	d.Set("db_instance_id", d.Id())
	d.Set("region_id", client.RegionId)

	objectRaw, err = mongodbServiceV2.DescribeAuditPolicyDescribeAuditLogFilter(d.Id())
	if err != nil && !NotFoundError(err) {
		return WrapError(err)
	}

	// DescribeAuditLogFilter returns a merged view like "mongos@admin,slow-db@admin,slow" with
	// RoleType=logic when RoleType is omitted on read. Parse out the segment for the configured
	// role_type so state matches the config; on import (role_type empty) fall back to the common
	// value across roles.
	configRoleType := d.Get("role_type").(string)
	d.Set("filter", parseMongodbShardingAuditFilter(fmt.Sprint(objectRaw["Filter"]), configRoleType))
	if configRoleType != "" {
		d.Set("role_type", configRoleType)
	} else if v, ok := objectRaw["RoleType"]; ok && v != nil {
		d.Set("role_type", v)
	}

	return nil
}

// parseMongodbShardingAuditFilter extracts the filter value for the configured role_type from the
// DescribeAuditLogFilter response. The API returns either a single value ("admin,slow") when a
// RoleType is honoured on read, or a merged view ("mongos@admin,slow-db@admin,slow") with
// RoleType=logic when RoleType is omitted. This function picks the segment matching the configured
// role_type, or the common value across segments when role_type is unset (e.g. on import).
func parseMongodbShardingAuditFilter(filter, roleType string) string {
	// "logic" is the API's all-role merged sentinel: DescribeAuditLogFilter returns it (with a
	// merged view like "mongos@admin,slow-db@admin,slow") when RoleType is omitted on read. A
	// resource whose config does not set role_type ends up with state.role_type="logic" after the
	// first Read, so on the next Refresh configRoleType="logic" would miss every segment and fall
	// back to the raw merged string, causing a perpetual diff. Treat "logic" like an unset
	// role_type so the common value across segments is extracted, matching the datasource path.
	if roleType == "logic" {
		roleType = ""
	}
	if filter == "" || !strings.Contains(filter, "@") {
		return filter
	}
	var values []string
	for _, seg := range strings.Split(filter, "-") {
		idx := strings.Index(seg, "@")
		if idx < 0 {
			continue
		}
		role, val := seg[:idx], seg[idx+1:]
		if roleType != "" {
			if role == roleType {
				return val
			}
			continue
		}
		values = append(values, val)
	}
	if roleType == "" && len(values) > 0 {
		first := values[0]
		allSame := true
		for _, v := range values[1:] {
			if v != first {
				allSame = false
				break
			}
		}
		if allSame {
			return first
		}
	}
	return filter
}

// mongodbShardingAuditFilterStorageValues picks the retention values that match the audit log edition
// in use. The reported edition wins; when the server reports a different one (it reports Trial while
// the audit log is disabled), the declared edition decides. A missing field yields nil so the caller
// keeps the state value.
func mongodbShardingAuditFilterStorageValues(objectRaw map[string]interface{}, reportedServiceType, declaredServiceType string) (storagePeriod, hotStoragePeriod interface{}) {
	edition := reportedServiceType
	if edition != "Standard" && edition != "V2_Standard" {
		edition = declaredServiceType
	}
	if edition == "V2_Standard" {
		return objectRaw["TtlForV2Standard"], objectRaw["HotTtlForV2Standard"]
	}
	return objectRaw["TtlForStandard"], nil
}

func resourceAliCloudMongodbShardingAuditFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mongodbServiceV2 := MongodbServiceV2{client}
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	// 1. ModifyAuditPolicy — audit_status / audit_log_switch_source / storage_period / hot_storage_period / service_type.
	action := "ModifyAuditPolicy"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	// audit_log_switch_source is request attribution, not a policy value: send it with a real policy
	// change only when it is configured, and never rewrite the whole policy because of it alone. An
	// imported state has no service_type while the audit log is disabled (Read keeps the declared
	// value only when enabled) and the disabled diff suppression hides the declared edition, so such
	// a write would fall back to ServiceType=Standard and could switch editions.
	if v, ok := d.GetOk("audit_log_switch_source"); ok {
		request["AuditLogSwitchSource"] = v
	}
	if !d.IsNewResource() && d.HasChange("storage_period") {
		update = true
	}
	if v, ok := d.GetOkExists("storage_period"); ok {
		request["StoragePeriod"] = v
	}
	if !d.IsNewResource() && d.HasChange("hot_storage_period") {
		update = true
	}
	if v, ok := d.GetOkExists("hot_storage_period"); ok {
		request["HotStoragePeriod"] = v
	}
	if !d.IsNewResource() && d.HasChange("service_type") {
		update = true
	}
	if !d.IsNewResource() && d.HasChange("audit_status") {
		update = true
	}
	request["AuditStatus"] = convertMongodbAuditPolicyAuditStatusRequest(d.Get("audit_status").(string))
	serviceType := d.Get("service_type").(string)
	if serviceType == "" {
		serviceType = "Standard"
	}
	request["ServiceType"] = serviceType
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"[Running]"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, mongodbServiceV2.DescribeAsyncMongodbAuditPolicyStateRefreshFunc(d, response, "$.DBInstances.DBInstance[*].DBInstanceStatus", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	// 2. ModifyAuditLogFilter — filter / role_type.
	update = false
	action = "ModifyAuditLogFilter"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["DBInstanceId"] = d.Id()
	request["RegionId"] = client.RegionId
	if !d.IsNewResource() && d.HasChange("filter") {
		update = true
	}
	// role_type is Optional+Computed; only forward a non-empty value so the drift between an empty
	// config and the API-read state does not push an unintended value to the server.
	if v, ok := d.GetOk("role_type"); ok && v.(string) != "" {
		if !d.IsNewResource() && d.HasChange("role_type") {
			update = true
		}
		request["RoleType"] = v
	}
	request["Filter"] = d.Get("filter")
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dds", "2015-12-01", action, query, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
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
		stateConf := BuildStateConf([]string{}, []string{"[Running]"}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, mongodbServiceV2.DescribeAsyncMongodbAuditPolicyStateRefreshFunc(d, response, "$.DBInstances.DBInstance[*].DBInstanceStatus", []string{}))
		if jobDetail, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id(), jobDetail)
		}
	}

	d.Partial(false)
	return resourceAliCloudMongodbShardingAuditFilterRead(d, meta)
}

func resourceAliCloudMongodbShardingAuditFilterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resource AliCloud Resource Mongodb Sharding Audit Filter. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
