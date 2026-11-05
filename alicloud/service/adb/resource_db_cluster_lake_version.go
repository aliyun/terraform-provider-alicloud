package adb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceAdbDbClusterLakeVersion() *schema.Resource {
	return &schema.Resource{
		Create: ResourceAdbDbClusterLakeVersionCreate,
		Read:   ResourceAdbDbClusterLakeVersionRead,
		Update: ResourceAdbDbClusterLakeVersionUpdate,
		Delete: ResourceAdbDbClusterLakeVersionDelete,
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
				ValidateFunc: validators.StringInSlice([]string{"5.0"}, false),
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
				ValidateFunc: validators.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
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
			// ProductForm<->ProductVersion form conversion is an in-place
			// ModifyDBCluster operation ("原地切换" per the published API):
			// LegacyForm -> IntegrationForm enterprise upgrade sends
			// ProductVersion=enterprise_version together with ProductForm,
			// ReservedNodeSize and ReservedNodeCount, without rebuilding the
			// cluster (DBClusterId is unchanged). ForceNew would rebuild and
			// break the customer's same-instance upgrade requirement, so it
			// must not be set here.
			"product_version": {
				Type:     schema.TypeString,
				Optional: true,
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
				ValidateFunc: validators.StringInSlice([]string{"Enable", "Disable", "Update"}, false),
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
				ValidateFunc: validators.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
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
				ValidateFunc: validators.StringInSlice([]string{"backup", "timepoint"}, false),
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
			return errs.WrapError(err)
		}
	case "Disable":
		if err := d.SetNew("enable_ssl", false); err != nil {
			return errs.WrapError(err)
		}
	case "Update":
		// Renewal: preserve the server's actual SSL state from the
		// prior state, not the user's new config value. This prevents
		// a conflicting enable_ssl change from driving the next
		// round's operation and unexpectedly disabling SSL.
		oldEnableSsl, _ := d.GetChange("enable_ssl")
		if err := d.SetNew("enable_ssl", oldEnableSsl); err != nil {
			return errs.WrapError(err)
		}
	}
	return nil
}

func ResourceAdbDbClusterLakeVersionCreate(d *schema.ResourceData, meta interface{}) error {
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

	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_adb_db_cluster_lake_version", action, errs.AlibabaCloudSdkGoERROR)
	}

	if v, err := jsonpath.Get("$.DBClusterId", response); err != nil || v == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_adb_db_cluster_lake_version")
	} else {
		d.SetId(fmt.Sprint(v))
	}

	stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{"Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	return ResourceAdbDbClusterLakeVersionUpdate(d, meta)
}

func ResourceAdbDbClusterLakeVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	adbService := AdbService{client}

	object, err := adbService.DescribeAdbDbClusterLakeVersion(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_adb_db_cluster_lake_version adbService.DescribeAdbDbClusterLakeVersion Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
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
		return errs.WrapError(err)
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

func ResourceAdbDbClusterLakeVersionUpdate(d *schema.ResourceData, meta interface{}) error {
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

	// ProductVersion is required by ModifyDBCluster for an in-place
	// ProductForm form conversion (e.g. LegacyForm -> IntegrationForm
	// enterprise upgrade: ProductVersion=enterprise_version together with
	// ProductForm, ReservedNodeSize and ReservedNodeCount). Without it the
	// API rejects the upgrade with InvalidParameter.ReservedNodeSize/
	// ReservedNodeCount ("...required for form conversion"). It is also
	// sent when the user changes product_version directly. Send the value
	// as-is from config/state; the stored/returned format is CamelCase
	// (e.g. EnterpriseVersion), matching CreateDBCluster.
	if !d.IsNewResource() && d.HasChange("product_version") {
		update = true

		if v, ok := d.GetOk("product_version"); ok {
			request["ProductVersion"] = v
		}
	}

	// Form conversion (ProductForm changed): force-send the target form's
	// required parameter set. These Computed attributes may be read back
	// from DescribeDBClusters with the same value the user configured,
	// making d.HasChange false and omitting them from the regular
	// per-parameter blocks above. Use d.Get (not d.GetOk) because d.GetOk
	// for Optional+Computed fields returns false when the config value
	// equals the state value (no diff entry), even though the value IS
	// present in config — this would silently drop ComputeResource/
	// StorageResource from the form-conversion request. d.Get returns the
	// config-or-state value regardless of diff state, so the target form's
	// parameter set is always sent.
	//
	// The API rejects a bare ProductForm change without its target parameter
	// set:
	// - LegacyForm->IntegrationForm upgrade: InvalidParameter.ReservedNodeSize/
	//   ReservedNodeCount ("...required for form conversion"). Per
	//   ModifyDBCluster docs, the upgrade must specify ProductForm,
	//   ReservedNodeSize and ReservedNodeCount together with
	//   ProductVersion=enterprise_version.
	// - IntegrationForm->LegacyForm downgrade: the request includes
	//   ProductForm=LegacyForm plus ComputeResource/StorageResource (via
	//   d.Get), yet the API still returns NoNeedToScale. This is handled
	//   as a non-fatal warning below — the conversion may be applied
	//   server-side; Read reconciles the actual state.
	if !d.IsNewResource() && d.HasChange("product_form") {
		targetForm := d.Get("product_form").(string)
		if targetForm == "IntegrationForm" {
			if v := d.Get("product_version").(string); v != "" {
				request["ProductVersion"] = v
			}
			if v := d.Get("reserved_node_size").(string); v != "" {
				request["ReservedNodeSize"] = v
			}
			if v := d.Get("reserved_node_count").(int); v != 0 {
				request["ReservedNodeCount"] = v
			}
		} else if targetForm == "LegacyForm" {
			if v := d.Get("compute_resource").(string); v != "" {
				request["ComputeResource"] = v
			}
			if v := d.Get("storage_resource").(string); v != "" {
				request["StorageResource"] = v
			}
		}
	}

	if update {
		action := "ModifyDBCluster"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"OperationDenied.OrderProcessing"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			// NoNeedToScale during a ProductForm conversion: the API
			// determined no scaling is required for the form conversion
			// (IntegrationForm→LegacyForm where reserved nodes map to
			// equivalent compute/storage). The conversion may be applied
			// server-side despite the 400 response; treat as non-fatal and
			// let Read reconcile state from DescribeDBClusters. Not silenced
			// for regular scaling requests (non-form-conversion), where
			// NoNeedToScale would surface as a genuine error to the user.
			if errs.IsExpectedErrors(err, []string{"NoNeedToScale"}) && d.HasChange("product_form") {
				logs.AddDebug(action, "NoNeedToScale during ProductForm conversion (non-fatal)", request)
				err = nil
			} else {
				return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
			}
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyClusterAccessWhiteListReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyClusterAccessWhiteListReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	update = false
	modifyDBClusterSSLReq := map[string]interface{}{
		"DBClusterId": d.Id(),
	}

	// ssl_enabled (Enable/Disable/Update) takes priority over enable_ssl
	// when both are set. CustomizeDiff resolves the priority via SetNew:
	// Enable/Disable override enable_ssl; Update preserves the old actual
	// enable_ssl state. Send SSLEnabled when ssl_enabled changed; otherwise
	// fall back to EnableSSL for a pure enable_ssl change. When ssl_enabled
	// is removed from config, CustomizeDiff SetNew removes it, making
	// d.HasChange("ssl_enabled") false and falling through to the EnableSSL
	// branch — a later enable_ssl-only change compares the preserved state
	// value to itself (equal-to-equal) and correctly falls through here.
	if d.HasChange("enable_ssl") || d.HasChange("ssl_enabled") {
		update = true
		if d.HasChange("ssl_enabled") {
			if v, ok := d.GetOk("ssl_enabled"); ok && v.(string) != "" {
				modifyDBClusterSSLReq["SSLEnabled"] = v
			}
		} else if v, ok := d.GetOkExists("enable_ssl"); ok {
			modifyDBClusterSSLReq["EnableSSL"] = v
		}
	}

	if update {
		action := "ModifyDBClusterSSL"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterSSLReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterSSLReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

		// Determine the target SSL state for the SSL state refresh.
		// SSLEnabled (Enable/Disable/Update) is the new API parameter;
		// EnableSSL (bool) is the legacy parameter. Map both to the
		// boolean string that
		// AdbDbClusterLakeVersionDBClusterSSLStateRefreshFunc returns
		// ("true"/"false"). For "Update" (certificate renewal), the SSL
		// state does not change — renewal is only valid when SSL is
		// enabled, so the target is "true"; the state refresh confirms
		// SSL is enabled but does NOT verify the renewal itself (that
		// requires checking ssl_expire_time changes independently).
		var sslTarget string
		if v, ok := modifyDBClusterSSLReq["SSLEnabled"]; ok {
			switch fmt.Sprint(v) {
			case "Enable":
				sslTarget = "true"
			case "Disable":
				sslTarget = "false"
			case "Update":
				sslTarget = "true"
			}
		} else if v, ok := modifyDBClusterSSLReq["EnableSSL"]; ok {
			sslTarget = fmt.Sprint(v)
		}

		if sslTarget != "" {
			sslEnabledStateConf := helper.BuildStateConf([]string{}, []string{sslTarget}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionDBClusterSSLStateRefreshFunc(d, []string{}))
			if _, err := sslEnabledStateConf.WaitForState(); err != nil {
				return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterDescriptionReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterDescriptionReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2021-12-01", action, nil, modifyDBClusterResourceGroupReq, true)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterResourceGroupReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
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
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("adb", "2019-03-15", action, nil, modifyDBClusterPayTypeReq, true)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"OperationDenied.DBClusterStatus"}) || helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, modifyDBClusterPayTypeReq)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		stateConf := helper.BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errs.WrapErrorf(err, errs.IdMsg, d.Id())
		}

	}

	d.Partial(false)

	return ResourceAdbDbClusterLakeVersionRead(d, meta)
}

func ResourceAdbDbClusterLakeVersionDelete(d *schema.ResourceData, meta interface{}) error {

	if v, ok := d.GetOk("payment_type"); ok && fmt.Sprint(v) == "Subscription" {
		return errs.WrapError(errs.Error("'Subscription' instance cannot be destroyed: please convert it to 'PayAsYouGo' (set payment_type = \"PayAsYouGo\" and re-apply) or unsubscribe it in the console. To remove it from the state only, run 'terraform state rm'."))
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
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("adb", "2021-12-01", action, nil, request, true)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	adbService := AdbService{client}
	stateConf := helper.BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutUpdate), 5*time.Minute, adbService.AdbDbClusterLakeVersionStateRefreshFunc(d, []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errs.WrapErrorf(err, errs.IdMsg, d.Id())
	}

	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidDBCluster.NotFound"}) || errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
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
