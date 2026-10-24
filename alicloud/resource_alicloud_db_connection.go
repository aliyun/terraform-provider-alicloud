package alicloud

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const dbConnectionSuffixRegex = "\\.mysql\\.([a-zA-Z0-9\\-]+\\.){0,1}rds\\.aliyuncs\\.com"
const dbConnectionIdWithSuffixRegex = "^([a-zA-Z0-9\\-_]+:[a-zA-Z0-9\\-_]+)" + dbConnectionSuffixRegex + "$"

var dbConnectionIdWithSuffixRegexp = regexp.MustCompile(dbConnectionIdWithSuffixRegex)

func resourceAlicloudDBConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDBConnectionCreate,
		Read:   resourceAlicloudDBConnectionRead,
		Update: resourceAlicloudDBConnectionUpdate,
		Delete: resourceAlicloudDBConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateDBConnectionPort,
				Default:      "3306",
			},
			"babelfish_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pgbouncer_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"general_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDBConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	instanceId := d.Get("instance_id").(string)
	prefix := d.Get("connection_prefix").(string)
	if prefix == "" {
		prefix = fmt.Sprintf("%stf", instanceId)
	}
	action := "AllocateInstancePublicConnection"
	request := map[string]interface{}{
		"RegionId":               client.RegionId,
		"DBInstanceId":           instanceId,
		"ConnectionStringPrefix": prefix,
		"Port":                   d.Get("port"),
		"SourceIp":               client.SourceIp,
	}

	if v, ok := d.GetOk("babelfish_port"); ok {
		request["BabelfishPort"] = v
	}
	setDBConnectionOptionalParams(request, d)

	var err error
	err = retry.Retry(8*time.Minute, func() *retry.RetryError {
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if isDBConnectionAllocateTransientError(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_db_connection", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%v", instanceId, COLON_SEPARATED, request["ConnectionStringPrefix"]))

	if err := rdsService.WaitForDBConnection(d.Id(), Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// wait instance running after allocating
	if err := rdsService.WaitForDBInstance(instanceId, Running, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	// A configured PgBouncer port is only converged once the API returns it for
	// this connection; otherwise the state would claim a port that was never
	// applied. The waiter only applies when the parent instance reports
	// PgBouncer as enabled: on other instances the parameter is not applicable
	// and is not applied.
	if v, ok := d.GetOk("pgbouncer_port"); ok && v.(string) != "" {
		if err := waitDBConnectionPgBouncerPortIfEnabled(rdsService, instanceId, d.Id(), v.(string), d.Timeout(schema.TimeoutCreate)); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudDBConnectionRead(d, meta)
}

func resourceAlicloudDBConnectionRead(d *schema.ResourceData, meta interface{}) error {
	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}
	object, err := rdsService.DescribeDBConnection(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	// The parent instance is the identity scope of this connection: once it is
	// confirmed gone the connection can neither be read nor addressed, so it
	// must not stay in state. A failed parent lookup is propagated instead of
	// being read as absence.
	instance, err := rdsService.describeRdsParentInstance(parts[0], 5*time.Minute)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_id", parts[0])
	d.Set("connection_prefix", parts[1])
	d.Set("port", object["Port"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("ip_address", object["IPAddress"])
	if object["BabelfishPort"] != nil {
		d.Set("babelfish_port", object["BabelfishPort"])
	}
	// PGBouncerPort is only returned for PostgreSQL instances with PgBouncer
	// enabled. On PostgreSQL the configured value is unconfirmed whenever the
	// API omits it, so it must not stay as converged state; on other engines or
	// when PgBouncer is known to be disabled the parameter is not applicable and
	// the configured value is kept instead of producing a permanent diff.
	if object["PGBouncerPort"] != nil {
		d.Set("pgbouncer_port", object["PGBouncerPort"])
	} else if dbConnectionPgBouncerPortUnconfirmed(instance) {
		d.Set("pgbouncer_port", "")
	}
	// GeneralGroupName is an instance attribute and is only returned for
	// MyBase MySQL Standard Edition instances.
	if v, ok := instance["GeneralGroupName"].(string); ok && v != "" {
		d.Set("general_group_name", v)
	}
	return nil
}

func resourceAlicloudDBConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	if d.HasChange("port") || d.HasChange("babelfish_port") || d.HasChange("pgbouncer_port") || d.HasChange("general_group_name") {
		action := "ModifyDBInstanceConnectionString"
		request := map[string]interface{}{
			"RegionId":     client.RegionId,
			"DBInstanceId": parts[0],
			"SourceIp":     client.SourceIp,
		}
		object, err := rdsService.DescribeDBConnection(d.Id())
		if err != nil {
			return WrapError(err)
		}
		request["CurrentConnectionString"] = object["ConnectionString"]
		request["ConnectionStringPrefix"] = parts[1]
		request["Port"] = d.Get("port")
		if v, ok := d.GetOk("babelfish_port"); ok {
			request["BabelfishPort"] = v
		}
		setDBConnectionOptionalParams(request, d)
		if err := retry.Retry(8*time.Minute, func() *retry.RetryError {
			response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, OperationDeniedDBStatus) || NeedRetry(err) {
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		// wait instance running after modifying
		if err := rdsService.WaitForDBInstance(request["DBInstanceId"].(string), Running, DefaultTimeoutMedium); err != nil {
			return WrapError(err)
		}
		if d.HasChange("pgbouncer_port") {
			if v, ok := d.GetOk("pgbouncer_port"); ok && v.(string) != "" {
				if err := waitDBConnectionPgBouncerPortIfEnabled(rdsService, parts[0], d.Id(), v.(string), d.Timeout(schema.TimeoutUpdate)); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}
	}
	return resourceAlicloudDBConnectionRead(d, meta)
}

// dbInstancePgBouncerState reports the PgBouncer flag of the parent instance.
// The second return value is false when the API omits the field or returns a
// value that is not a boolean, so callers can tell "disabled" from "unknown".
func dbInstancePgBouncerState(instance map[string]interface{}) (enabled bool, known bool) {
	raw, ok := instance["PGBouncerEnabled"]
	if !ok || raw == nil {
		return false, false
	}
	parsed, err := strconv.ParseBool(fmt.Sprint(raw))
	if err != nil {
		return false, false
	}
	return parsed, true
}

// dbConnectionPgBouncerPortUnconfirmed reports whether a missing PGBouncerPort
// means the configured value cannot be confirmed: the parameter only applies to
// PostgreSQL, and there it is unconfirmed unless the instance explicitly reports
// PgBouncer as disabled.
func dbConnectionPgBouncerPortUnconfirmed(instance map[string]interface{}) bool {
	engine, _ := instance["Engine"].(string)
	if engine == "" || engine != "PostgreSQL" {
		return false
	}
	enabled, known := dbInstancePgBouncerState(instance)
	return !known || enabled
}

// waitDBConnectionPgBouncerPortIfEnabled waits for the requested PgBouncer port
// only when the parent instance reports PgBouncer as enabled. The PgBouncer flag
// is polled until it is known: an omitted flag is not evidence that PgBouncer is
// disabled, so the apply must not record an unconfirmed port. Instances of other
// engines never apply this parameter.
func waitDBConnectionPgBouncerPortIfEnabled(rdsService RdsService, instanceId, connectionId, expected string, timeout time.Duration) error {
	// The flag wait is bounded by the caller's timeout, capped so a long default
	// resource timeout cannot stretch it to an unreasonable window. The port wait
	// below shares the same deadline, so the two stages never consume two full
	// budgets.
	limit := timeout
	if limit <= 0 || limit > 15*time.Minute {
		limit = 15 * time.Minute
	}
	deadline := time.Now().Add(limit)
	wait := incrementalWait(2*time.Second, 2*time.Second)
	var enabled bool
	err := retry.RetryContext(context.Background(), limit, func() *retry.RetryError {
		instance, err := rdsService.describeRdsParentInstance(instanceId, 0)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		if engine, _ := instance["Engine"].(string); engine != "" && engine != "PostgreSQL" {
			enabled = false
			return nil
		}
		value, known := dbInstancePgBouncerState(instance)
		if !known {
			wait()
			return retry.RetryableError(fmt.Errorf("pgbouncer state not confirmed yet: instance %s reports status %v without PGBouncerEnabled", instanceId, instance["DBInstanceStatus"]))
		}
		enabled = value
		return nil
	})
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	remaining, err := pgBouncerPortRemainingBudget(deadline, limit)
	if err != nil {
		return err
	}
	return waitDBConnectionPgBouncerPortConverged(rdsService, connectionId, expected, remaining)
}

// pgBouncerPortRemainingBudget returns the part of the shared budget that is
// still available for the port wait, or an error when the PgBouncer state wait
// consumed all of it.
func pgBouncerPortRemainingBudget(deadline time.Time, limit time.Duration) (time.Duration, error) {
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return 0, fmt.Errorf("pgbouncer state wait used the whole %s budget; the port was not read", limit)
	}
	return remaining, nil
}

// waitDBConnectionPgBouncerPortConverged waits until DescribeDBInstanceNetInfo
// returns the requested PgBouncer port for this connection. A successful
// Allocate/Modify only means the request was accepted; an omitted port is not
// confirmation that the value was applied.
func waitDBConnectionPgBouncerPortConverged(rdsService RdsService, id, expected string, timeout time.Duration) error {
	wait := incrementalWait(3*time.Second, 3*time.Second)
	return retry.RetryContext(context.Background(), timeout, func() *retry.RetryError {
		object, err := rdsService.DescribeDBConnection(id)
		if err != nil {
			if NotFoundError(err) || NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		raw, ok := object["PGBouncerPort"]
		if !ok || raw == nil || fmt.Sprint(raw) == "" {
			wait()
			return retry.RetryableError(fmt.Errorf("pgbouncer_port not confirmed yet: DescribeDBInstanceNetInfo omits PGBouncerPort"))
		}
		if fmt.Sprint(raw) != expected {
			wait()
			return retry.RetryableError(fmt.Errorf("pgbouncer_port not converged: expected %s, got %v", expected, raw))
		}
		return nil
	})
}

// isDBConnectionAllocateTransientError reports whether
// AllocateInstancePublicConnection failed for a transient reason: the instance
// is busy with another task, or the connection is still being created
// (NetworkConnectionCreating).
func isDBConnectionAllocateTransientError(err error) bool {
	return IsExpectedErrors(err, OperationDeniedDBStatus) ||
		IsExpectedErrors(err, []string{"NetworkConnectionCreating"}) ||
		NeedRetry(err)
}

// setDBConnectionOptionalParams adds the connection-string parameters that are
// managed by this resource but are optional: the PgBouncer port (PostgreSQL
// only, requires PgBouncer enabled) and the general group name (MyBase MySQL
// Standard Edition only).
func setDBConnectionOptionalParams(request map[string]interface{}, d *schema.ResourceData) {
	if v, ok := d.GetOk("pgbouncer_port"); ok {
		request["PGBouncerPort"] = v
	}
	if v, ok := d.GetOk("general_group_name"); ok {
		request["GeneralGroupName"] = v
	}
}

func resourceAlicloudDBConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	rdsService := RdsService{client}

	submatch := dbConnectionIdWithSuffixRegexp.FindStringSubmatch(d.Id())
	if len(submatch) > 1 {
		d.SetId(submatch[1])
	}

	split := strings.Split(d.Id(), COLON_SEPARATED)
	action := "ReleaseInstancePublicConnection"
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"DBInstanceId": split[0],
		"SourceIp":     client.SourceIp,
	}
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		object, err := rdsService.DescribeDBConnection(d.Id())
		if err != nil {
			return retry.NonRetryableError(WrapError(err))
		}
		request["CurrentConnectionString"] = object["ConnectionString"]
		response, err := client.RpcPost("Rds", "2014-08-15", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationDenied.DBInstanceStatus"}) || NeedRetry(err) {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})

	if err != nil {
		if NotFoundError(err) || IsExpectedErrors(err, []string{"InvalidCurrentConnectionString.NotFound", "AtLeastOneNetTypeExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return rdsService.WaitForDBConnection(d.Id(), Deleted, DefaultTimeoutMedium)
}
