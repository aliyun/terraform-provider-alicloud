package alicloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudAdbDbClusterLakeVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudAdbDbClusterLakeVersionCreate,
		Read:   resourceAliCloudAdbDbClusterLakeVersionRead,
		Update: resourceAliCloudAdbDbClusterLakeVersionUpdate,
		Delete: resourceAliCloudAdbDbClusterLakeVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(72 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"db_cluster_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: StringInSlice([]string{"5.0"}, false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"secondary_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secondary_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_form": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"product_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"compute_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_resource": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reserved_node_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"reserved_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"disk_encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"kms_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_ips": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ssl_enabled": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"Enable", "Disable", "Update"}, false),
			},
			"ssl_connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_cluster_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
			},
			"enable_default_resource_group": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source_db_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backup_set_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restore_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"backup", "timepoint"}, false),
			},
			"restore_to_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"commodity_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: resourceAdbDbClusterLakeVersionCustomizeDiff,
	}
}

// resourceAdbDbClusterLakeVersionCustomizeDiff implements the product
// contract: when ssl_enabled is explicitly set in HCL config, it takes
// priority over enable_ssl; otherwise enable_ssl (the legacy boolean) is
// used as-is.
//
// GetRawConfig returns the raw HCL config as a cty.Value. Unlike GetOk
// (which reads merged state/config/diff via ReadFieldMerge and cannot
// distinguish an explicitly configured ssl_enabled from a state-retained
// value), GetRawConfig only reflects what the user wrote in HCL. When
// ssl_enabled is omitted from config, GetRawConfig returns a typed null
// for that attribute, so CustomizeDiff correctly returns nil and lets
// the legacy enable_ssl boolean drive the operation. This prevents a
// removed ssl_enabled command from reading stale state and swallowing a
// legitimate enable_ssl true→false change.
//
// Unknown guard: cross-resource references may produce cty unknown
// values during planning. Skip processing until the value is known.
//
// When ssl_enabled is explicitly set and known:
//   - "Enable"  → plan enable_ssl=true  (priority)
//   - "Disable" → plan enable_ssl=false (priority)
//   - "Update"  → preserve the old actual SSL state (from prior state),
//     not the user's new enable_ssl config value. Without this, a
//     conflicting enable_ssl change would drive the next round's
//     operation and unexpectedly disable SSL. By forcing enable_ssl
//     back to its old state value, HasChange("enable_ssl") returns
//     false on the next plan, preventing an unintended Disable.
//
// When ssl_enabled is absent (not in config), enable_ssl is used as-is.
func resourceAdbDbClusterLakeVersionCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	rawConfig := d.GetRawConfig()
	// Defensive: in production Terraform always provides a valid object.
	// This guards against nil/non-object raw config in edge cases.
	if rawConfig.IsNull() || !rawConfig.Type().IsObjectType() {
		return nil
	}
	sslEnabledVal := rawConfig.GetAttr("ssl_enabled")
	// Unknown guard: skip when the value is unknown (cross-resource
	// references during planning).
	if !sslEnabledVal.IsKnown() {
		return nil
	}
	// Not configured: ssl_enabled is absent from HCL (typed null).
	// The legacy enable_ssl boolean is used as-is.
	if sslEnabledVal.IsNull() {
		return nil
	}
	sslEnabledStr := sslEnabledVal.AsString()
	if sslEnabledStr == "" {
		return nil
	}
	switch sslEnabledStr {
	case "Enable":
		if err := d.SetNew("enable_ssl", true); err != nil {
			return WrapError(err)
		}
	case "Disable":
		if err := d.SetNew("enable_ssl", false); err != nil {
			return WrapError(err)
		}
	case "Update":
		// Renewal: preserve the server's actual SSL state from the
		// prior state, not the user's new config value. This prevents
		// a conflicting enable_ssl change from driving the next
		// round's operation and unexpectedly disabling SSL.
		oldEnableSsl, _ := d.GetChange("enable_ssl")
		if err := d.SetNew("enable_ssl", oldEnableSsl); err != nil {
			return WrapError(err)
		}
	}
	return nil
}

func resourceAliCloudAdbDbClusterLakeVersionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	var response map[string]interface{}
	var err error
	action := "CreateDBCluster"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["DBClusterVersion"] = d.Get("db_cluster_version")
	request["VPCId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["ZoneId"] = d.Get("zone_id")
	request["PayType"] = convertAdbDbClusterLakeVersionPaymentTypeRequest(d.Get("payment_type"))

	if v, ok := d.GetOk("secondary_vswitch_id"); ok {
		request["SecondaryVSwitchId"] = v
	}

	if v, ok := d.GetOk("secondary_zone_id"); ok {
		request["SecondaryZoneId"] = v
	}

	if v, ok := d.GetOk("product_form"); ok {
		request["ProductForm"] = v
	}

	if v, ok := d.GetOk("product_version"); ok {
		request["ProductVersion"] = v
	}

	if v, ok := d.GetOk("compute_resource"); ok {
		request["ComputeResource"] = v
	}

	if v, ok := d.GetOk("storage_resource"); ok {
		request["StorageResource"] = v
	}

	if v, ok := d.GetOk("reserved_node_size"); ok {
		request["ReservedNodeSize"] = v
	}

	if v, ok := d.GetOkExists("reserved_node_count"); ok {
		request["ReservedNodeCount"] = v
	}

	if v, ok := d.GetOkExists("disk_encryption"); ok {
		request["DiskEncryption"] = v
	}

	if v, ok := d.GetOk("kms_id"); ok {
		request["KmsId"] = v
	}

	if v, ok := d.GetOk("db_cluster_description"); ok {
		request["DBClusterDescription"] = v
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if request["PayType"] == string(Prepaid) {
		if v, ok := d.GetOkExists("period"); ok {
			usedTime := v.(int)
			request["UsedTime"] = strconv.Itoa(usedTime)
			request["Period"] = string(Month)
			if usedTime > 9 {
				request["UsedTime"] = strconv.Itoa(usedTime / 12)
				request["Period"] = string(Year)
			}
		}
	}

	if v, ok := d.GetOkExists("enable_default_resource_group"); ok {
		request["EnableDefaultResourcePool"] = v
	}

	if v, ok := d.GetOk("source_db_cluster_id"); ok {
		request["SourceDbClusterId"] = v
	}

	if v, ok := d.GetOk("backup_set_id"); ok {
		request["BackupSetId"] = v
	}

	if v, ok := d.GetOk("restore_type"); ok {
		request["RestoreType"] = v
	}

	if v, ok := d.GetOk("restore_to_time"); ok {
		request["RestoreToTime"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_adb_db_cluster_lake_version", action, AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.DBClusterId", response); err != nil || v == nil {
		return WrapErrorf(err, IdMsg, "alicloud_adb_db_cluster_lake_version")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAliCloudAdbDbClusterLakeVersionUpdate(d, meta)
}

func resourceAliCloudAdbDbClusterLakeVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	object, err := adbService.DescribeAdbDbClusterLakeVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version adbService.DescribeAdbDbClusterLakeVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("db_cluster_version", object["DBVersion"])
	d.Set("vpc_id", object["VPCId"])
	d.Set("vswitch_id", object["VSwitchId"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("compute_resource", object["ComputeResource"])
	d.Set("storage_resource", object["StorageResource"])
	d.Set("payment_type", convertAdbDbClusterLakeVersionPaymentTypeResponse(object["PayType"]))
	d.Set("secondary_vswitch_id", object["SecondaryVSwitchId"])
	d.Set("secondary_zone_id", object["SecondaryZoneId"])
	d.Set("product_form", object["ProductForm"])
	d.Set("product_version", object["ProductVersion"])
	d.Set("reserved_node_size", object["ReservedNodeSize"])
	d.Set("reserved_node_count", object["ReservedNodeCount"])
	d.Set("disk_encryption", object["DiskEncryption"])
	d.Set("kms_id", object["KmsId"])
	d.Set("db_cluster_description", object["DBClusterDescription"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("engine", object["Engine"])
	d.Set("engine_version", object["EngineVersion"])
	d.Set("commodity_code", object["CommodityCode"])
	d.Set("connection_string", object["ConnectionString"])
	d.Set("port", object["Port"])
	d.Set("lock_mode", object["LockMode"])
	d.Set("lock_reason", object["LockReason"])
	d.Set("expired", object["Expired"])
	d.Set("create_time", object["CreationTime"])
	d.Set("expire_time", object["ExpireTime"])
	d.Set("status", object["DBClusterStatus"])

	SecurityIPsObject, err := adbService.DescribeClusterAccessWhiteList(d.Id())
	if err != nil {
		return WrapError(err)
	}

	d.Set("security_ips", SecurityIPsObject["SecurityIPList"])

	sslObject, err := adbService.DescribeAdbDbClusterLakeVersionDBClusterSSL(d.Id())
	if err != nil {
		log.Printf("[WARN] Resource alicloud_adb_db_cluster_lake_version DescribeAdbDbClusterLakeVersionDBClusterSSL Failed!!! %s", err)
	} else {
		d.Set("enable_ssl", sslObject["SSLEnabled"])
		d.Set("ssl_connection_string", sslObject["ConnectionString"])
		d.Set("ssl_expire_time", sslObject["ExpireTime"])
		// SslEnabled (Enable/Disable/Update) is a write-only SSL operation
		// command, so the API does not return it directly. When the state is
		// "Update", retain it to prevent a persistent diff that would
		// re-trigger renewal on every plan/apply — this is for Terraform
		// state convergence only and does NOT prove the renewal was executed
		// on the server; independent evidence (ssl_expire_time change) is
		// required to verify renewal. For "Enable"/"Disable", derive the
		// actual SSL state from the SSLEnabled boolean. This read is also
		// required by the Update trigger: a later enable_ssl-only change
		// (with ssl_enabled removed from config) compares the preserved
		// state value to itself (equal-to-equal) and falls through to the
		// EnableSSL branch. When both ssl_enabled and enable_ssl are set in
		// config, CustomizeDiff resolves the priority (ssl_enabled takes
		// priority) via SetNew: Enable/Disable override enable_ssl, and
		// Update preserves the old actual enable_ssl state to prevent a
		// conflicting bool from driving the next round's operation.
		if d.Get("ssl_enabled").(string) != "Update" {
			sslEnabledStr := fmt.Sprint(sslObject["SSLEnabled"])
			var sslEnabledValue string
			switch sslEnabledStr {
			case "true", "Enable", "1":
				sslEnabledValue = "Enable"
			default:
				sslEnabledValue = "Disable"
			}
			d.Set("ssl_enabled", sslEnabledValue)
		}
	}

	return nil
}

func resourceAliCloudAdbDbClusterLakeVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("compute_resource") {
		update = true
	}
	if v, ok := d.GetOk("compute_resource"); ok {
		request["ComputeResource"] = v
	}

	if !d.IsNewResource() && d.HasChange("storage_resource") {
		update = true
	}
	if v, ok := d.GetOk("storage_resource"); ok {
		request["StorageResource"] = v
	}

	if !d.IsNewResource() && d.HasChange("reserved_node_size") {
		update = true

		if v, ok := d.GetOk("reserved_node_size"); ok {
			request["ReservedNodeSize"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("reserved_node_count") {
		update = true

		if v, ok := d.GetOkExists("reserved_node_count"); ok {
			request["ReservedNodeCount"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("product_form") {
		update = true

		if v, ok := d.GetOk("product_form"); ok {
			request["ProductForm"] = v
		}
	}

	if update {
		action := "ModifyDBCluster"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || NeedRetry(err) {
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

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyClusterAccessWhiteListReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if d.HasChange("security_ips") {
		update = true
	}
	if v, ok := d.GetOk("security_ips"); ok {
		modifyClusterAccessWhiteListReq["SecurityIps"] = v
	}

	if update {
		action := "ModifyClusterAccessWhiteList"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyClusterAccessWhiteListReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyClusterAccessWhiteListReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterSSLReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if d.HasChange("enable_ssl") || d.HasChange("ssl_enabled") {
		update = true

		// SslEnabled (Enable/Disable/Update) is a write-only operation command;
		// EnableSSL (boolean) reports/sets the actual state. Detect an explicit
		// ssl_enabled change with HasChange rather than GetOk: ssl_enabled is
		// read back into state as the actual Enable/Disable value, so GetOk
		// would always be true once SSL is configured and would misroute a
		// later enable_ssl-only change (e.g. enable_ssl=false) into the
		// SSLEnabled branch, omitting EnableSSL and failing with
		// MissingEnableSSL. HasChange is false when ssl_enabled is removed
		// from config because Optional+Computed preserves the prior state
		// value (equal-to-equal), so the change routes to EnableSSL.
		if d.HasChange("ssl_enabled") {
			if v, ok := d.GetOk("ssl_enabled"); ok && v.(string) != "" {
				modifyDBClusterSSLReq["SSLEnabled"] = v
				// EnableSSL is mandatory for ModifyDBClusterSSL: the API
				// rejects the request with MissingEnableSSL when it is omitted.
				// Derive the required boolean from the SslEnabled command so
				// both the operation command and the mandatory flag are sent:
				// Enable turns SSL on, Disable turns it off, and Update
				// derives EnableSSL from the actual old SSL state (preserved
				// by CustomizeDiff via SetNew), not a hardcoded true.
				switch v.(string) {
				case "Enable":
					modifyDBClusterSSLReq["EnableSSL"] = true
				case "Disable":
					modifyDBClusterSSLReq["EnableSSL"] = false
				case "Update":
					oldEnableSsl, _ := d.GetChange("enable_ssl")
					modifyDBClusterSSLReq["EnableSSL"] = oldEnableSsl
				}
			}
		} else if v, ok := d.GetOkExists("enable_ssl"); ok {
			modifyDBClusterSSLReq["EnableSSL"] = v
		}
	}

	if update {
		action := "ModifyDBClusterSSL"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterSSLReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterSSLReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

		// Determine the target SSL state for state refresh
		var sslTarget string
		if v, ok := modifyDBClusterSSLReq["SSLEnabled"]; ok {
			switch fmt.Sprint(v) {
			case "Enable":
				sslTarget = "true"
			case "Disable":
				sslTarget = "false"
			case "Update":
				// Certificate renewal does not change the SSL enabled state —
				// SSLEnabled stays true after renewal. Always target "true"
				// regardless of enable_ssl to avoid waiting on a state (e.g.
				// false) that the renewal operation will never reach.
				sslTarget = "true"
			}
		} else if v, ok := modifyDBClusterSSLReq["EnableSSL"]; ok {
			sslTarget = fmt.Sprint(v)
		}

		if sslTarget != "" {
			// "provisioning" is a pending state returned by the SSL state
			// refresh when SSLEnabled is true but ConnectionString has not
			// been provisioned yet. Listing it as Pending keeps the state
			// refresh polling until the SSL certificate is fully configured.
			sslEnabledStateConf := BuildStateConf([]string{"provisioning"}, []string{sslTarget}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionDBClusterSSLStateRefreshFunc(d, []string{}))
			if _, err := sslEnabledStateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

	}

	update = false
	modifyDBClusterDescriptionReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("db_cluster_description") {
		update = true
	}
	if v, ok := d.GetOk("db_cluster_description"); ok {
		modifyDBClusterDescriptionReq["DBClusterDescription"] = v
	}

	if update {
		action := "ModifyDBClusterDescription"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterDescriptionReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterDescriptionReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterResourceGroupReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		modifyDBClusterResourceGroupReq["NewResourceGroupId"] = v
	}

	if update {
		action := "ModifyDBClusterResourceGroup"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterResourceGroupReq, true)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterResourceGroupReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterPayTypeReq := map[string]interface{}{
		"RegionId":    client.RegionId,
		"DbClusterId": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("payment_type") {
		update = true
	}
	modifyDBClusterPayTypeReq["PayType"] = convertAdbDbClusterLakeVersionPaymentTypeRequest(d.Get("payment_type"))

	if modifyDBClusterPayTypeReq["PayType"] == string(Prepaid) {
		if v, ok := d.GetOkExists("period"); ok {
			usedTime := v.(int)
			modifyDBClusterPayTypeReq["UsedTime"] = strconv.Itoa(usedTime)
			modifyDBClusterPayTypeReq["Period"] = string(Month)
			if usedTime > 9 {
				modifyDBClusterPayTypeReq["UsedTime"] = strconv.Itoa(usedTime / 12)
				modifyDBClusterPayTypeReq["Period"] = string(Year)
			}
		}
	}

	if update {
		action := "ModifyDBClusterPayType"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterPayTypeReq, true)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus"}) || NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyDBClusterPayTypeReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}

	}

	d.Partial(false)

	return resourceAliCloudAdbDbClusterLakeVersionRead(d, meta)
}

func resourceAliCloudAdbDbClusterLakeVersionDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok && fmt.Sprint(v) == "Subscription" {
		return WrapError(Error("'Subscription' instance cannot be destroyed: please convert it to 'PayAsYouGo' (set payment_type = \"PayAsYouGo\" and re-apply) or unsubscribe it in the console. To remove it from the state only, run 'terraform state rm'."))
	}

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDBCluster"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"DBClusterId": d.Id(),
		"RegionId":    client.RegionId,
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
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

	// If DeleteDBCluster returned an error, check for NotFound before
	// running WaitForState. InvalidDBClusterId.NotFound means the cluster
	// is already gone — return nil promptly without a meaningless wait.
	// InvalidDBCluster.NotFound is preserved for backward compatibility.
	// Other errors (permission, throttling, etc.) are propagated
	// transparently without an unnecessary state wait.
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	// DeleteDBCluster succeeded; wait for the cluster to be fully deleted.
	adbService := AdbService{client}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}

func convertAdbDbClusterLakeVersionPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "Postpaid"
	case "Subscription":
		return "Prepaid"
	}

	return source
}

func convertAdbDbClusterLakeVersionPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "Postpaid":
		return "PayAsYouGo"
	case "Prepaid":
		return "Subscription"
	}

	return source
}
