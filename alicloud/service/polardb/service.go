// Package polardb. This file is generated automatically. Please do not modify it manually, thank you!
// Self-contained service helpers for the framework-native alicloud_polardb_backup resource.
// fwadapt.ResourceBase only exposes Client()/Configure(); the retry/token/tags/state-wait
// helpers live in package alicloud (circular from here) or are unexported in connectivity,
// so this file mirrors the v1 GwlbServiceV2 surface independently inside package polardb.
package polardb

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

type fwTimeoutDurationValidator struct{}

type fwWireField struct {
	Name   string
	Single bool
	Read   bool
	Input  bool
	Fields map[string]fwWireField
}

type fwMapping struct {
	Source, Target, Encoding, Constant string
	Group, GroupKey, GroupValue        string
	Repeat                             bool
	IndexSources                       []int
	Values                             map[string]string
}

type fwMappedValue struct {
	Value   interface{}
	Indexes []int
}
type fwDeleteCondition struct {
	Property, Action string
	Values           []string
}
type fwScopedReadIdentity struct{ Source, Target, Value string }

func fwTimeoutTypes() map[string]attr.Type {
	return map[string]attr.Type{"create": types.StringType, "update": types.StringType, "delete": types.StringType}
}

func fwTimeoutBlock() schema.SingleNestedBlock {
	attributes := map[string]schema.Attribute{}
	for _, operation := range []string{"create", "update", "delete"} {
		attributes[operation] = schema.StringAttribute{
			Optional:    true,
			Description: "Maximum retry and wait duration for this operation, for example 10m.",
			Validators:  []validator.String{fwTimeoutDurationValidator{}},
		}
	}
	return schema.SingleNestedBlock{Attributes: attributes}
}

func (fwTimeoutDurationValidator) Description(context.Context) string {
	return "A positive Go duration, for example 10m or 1h."
}
func (v fwTimeoutDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}
func (fwTimeoutDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	duration, err := time.ParseDuration(req.ConfigValue.ValueString())
	if err != nil || duration <= 0 {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
	}
}

func fwTimeoutDuration(value types.Object, operation string, fallback time.Duration) (time.Duration, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	if value.IsNull() {
		return fallback, diagnostics
	}
	if value.IsUnknown() {
		diagnostics.AddError("Unknown operation timeout", "Resolve the timeouts block before applying.")
		return 0, diagnostics
	}
	raw, ok := value.Attributes()[operation]
	if !ok || raw.IsNull() {
		return fallback, diagnostics
	}
	configured, ok := raw.(types.String)
	if !ok || configured.IsUnknown() {
		diagnostics.AddError("Invalid operation timeout", "Resolve the operation timeout to a duration string before applying.")
		return 0, diagnostics
	}
	duration, err := time.ParseDuration(configured.ValueString())
	if err != nil || duration <= 0 {
		diagnostics.AddError("Invalid operation timeout", "Use a positive Go duration, for example 10m or 1h.")
		return 0, diagnostics
	}
	return duration, diagnostics
}

// retryRpcPost mirrors alicloud's retry.Retry+NeedRetry wrapper around RpcPost.
// RpcPost returns *tea.SDKError, so the classifier only handles that path.
func retryRpcPost(client *connectivity.AliyunClient, product, version, action string, request map[string]interface{}, timeout time.Duration) (map[string]interface{}, error) {
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var response map[string]interface{}
	err := retry.Retry(timeout, func() *retry.RetryError {
		resp, e := client.RpcPost(product, version, action, nil, request, true)
		response = resp
		if e != nil {
			if needRetry(e) {
				wait()
				return retry.RetryableError(e)
			}
			return retry.NonRetryableError(e)
		}
		return nil
	})
	if err != nil {
		return response, errs.WrapErrorf(err, "calling %s failed", action)
	}
	return response, nil
}

// buildClientToken mirrors alicloud.buildClientToken (TF-<action>-<unix>-<uuid>, <=64 chars).
func buildClientToken(action string) string {
	token := strings.TrimSpace(fmt.Sprintf("TF-%s-%d-%s", action, time.Now().Unix(), strings.Trim(uuid.New().String(), "-")))
	if len(token) > 64 {
		token = token[0:64]
	}
	return token
}

// incrementalWait mirrors alicloud.incrementalWait: exponential backoff capped at 20s with jitter.
func incrementalWait(first, incr time.Duration) func() {
	retryCount := 0
	return func() {
		upper := float64(first) + (math.Pow(2, float64(retryCount))-1)*float64(incr)
		if upper > float64(20*time.Second) {
			upper = float64(20 * time.Second)
		}
		waitTime := time.Duration(upper * (0.8 + 0.4*rand.Float64()))
		time.Sleep(waitTime)
		retryCount++
	}
}

// needRetry mirrors alicloud.NeedRetry for the *tea.SDKError path that RpcPost produces.
func needRetry(err error) bool {
	if err == nil {
		return false
	}
	if strings.HasPrefix(err.Error(), "Post https://") {
		return true
	}
	if e, ok := err.(*tea.SDKError); ok {
		if e.Code == nil || e.Message == nil {
			return false
		}
		code := *e.Code
		msg := *e.Message
		if strings.Contains(msg, "code: 500, 您已开通过") {
			return false
		}
		if strings.Contains(msg, "Client.Timeout") {
			return true
		}
		if code == "ServiceUnavailable" || code == "Rejected.Throttling" || strings.Contains(code, "Throttling") || strings.HasPrefix(msg, "code: 5") {
			return true
		}
		// gap#5: v1 retries delete/read per-call on conflict/status/busy codes
		// (OperationConflict, LastTokenProcessing, IncorrectStatus, SystemBusy) and on
		// DependencyViolation (deps clear on retry, e.g. DeleteHaVip returning
		// 400 DependencyViolation.HaVip). The v2 generator uses a single global needRetry for
		// all RpcPost calls, so these codes never retried and Delete failed fast on
		// DependencyViolation. Including them here mirrors v1; all are transient and bounded
		// by the call's retry timeout. DependencyViolation carries a product suffix
		// (.HaVip/.GatewayRouteTable/...), so match by prefix.
		if strings.HasPrefix(code, "DependencyViolation") || code == "OperationConflict" ||
			code == "LastTokenProcessing" || code == "IncorrectStatus" || code == "SystemBusy" {
			return true
		}
	}
	return false
}

// isExpectedErrors matches structured API codes, including wrapped RPC failures.
// Error text, HTTP status alone and diagnostic Data are not resource identity evidence.
func isExpectedErrors(err error, codes []string) bool {
	for err != nil {
		if sdkError, ok := err.(*tea.SDKError); ok {
			if sdkError.Code == nil {
				return false
			}
			for _, code := range codes {
				if *sdkError.Code == code {
					return true
				}
			}
			return false
		}
		if wrapped, ok := err.(*errs.ComplexError); ok {
			err = wrapped.Cause
		} else if wrapped, ok := err.(interface{ Unwrap() error }); ok {
			err = wrapped.Unwrap()
		} else {
			return false
		}
	}
	return false
}

// backupReadNotFoundCodes are the NotFound error codes returned by DescribeBackupPolicy when
// the cluster (and therefore its backup) no longer exists. F5: the empty codes list previously
// never recognized NotFound, so an externally deleted cluster surfaced as a diagnostic instead
// of clearing state.
var backupReadNotFoundCodes = []string{"InvalidDBCluster.NotFound", "InvalidDBClusterId.NotFound"}

// backupDeleteNotFoundCodes are the NotFound error codes returned by DeleteBackup when the
// backup set or cluster is already gone. F5: Delete only accepted generic ResourceNotFound
// codes and rejected the official InvalidBackupSetID.NotFound / InvalidDBClusterId.NotFound.
var backupDeleteNotFoundCodes = []string{"InvalidBackupSetID.NotFound", "InvalidDBClusterId.NotFound", "ResourceNotFound.Backup", "ResourceNotFound"}

// backupParseID splits the composite resource id "DBClusterId:BackupId" into its parts.
// F1: the stable identity of a backup set is the BackupId returned by CreateBackup (the
// BackupJobId, which is the backup set id). Pairing it with the owning DBClusterId lets Read
// find the exact backup and Delete target only the owned backup instead of "last wins".
// A single-segment (legacy) id is treated as DBClusterId with an empty BackupId.
func backupParseID(id string) (dbClusterId, backupId string) {
	if idx := strings.Index(id, ":"); idx >= 0 {
		return id[:idx], id[idx+1:]
	}
	return id, ""
}

// backupComposeID builds the composite resource id "DBClusterId:BackupId".
func backupComposeID(dbClusterId, backupId string) string {
	return dbClusterId + ":" + backupId
}

// backupGet describes the resource by id. Returns (nil, nil) when it is gone, mirroring
// v1 GwlbServiceV2.DescribeGwlbListener (Get RPC + retry + ResourceNotFound -> not found).
func backupGet(client *connectivity.AliyunClient, id string, includeSecondary ...bool) (map[string]interface{}, error) {
	completeRead := len(includeSecondary) == 0 || includeSecondary[0]
	// F1: the resource id is the composite "DBClusterId:BackupId". DescribeBackupPolicy is
	// keyed on DBClusterId, so resolve the cluster part; the stable BackupId is injected into
	// the read object below so DescribeBackups/DescribeBackupTasks filter on the exact backup
	// instead of merging "last wins" over every backup set in the cluster.
	dbClusterId, backupId := backupParseID(id)
	getValues := map[string]interface{}{"RegionId": client.RegionId, "DBClusterId": dbClusterId}
	request := map[string]interface{}{}
	request["RegionId"] = client.RegionId
	if err := fwApplyMappings(request, getValues, []fwMapping{
		{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	}, false); err != nil {
		return nil, fmt.Errorf("building DescribeBackupPolicy request: %w", err)
	}
	if value, ok := request["DBClusterId"].(string); !ok || value == "" || value != dbClusterId {
		return nil, fmt.Errorf("building DescribeBackupPolicy request: required identity does not match")
	}
	response, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeBackupPolicy", request, 1*time.Minute)
	if err != nil {
		// F5: recognize the official NotFound codes so an externally deleted cluster clears
		// state instead of surfacing as an unresolved diagnostic.
		if isExpectedErrors(err, backupReadNotFoundCodes) {
			return nil, nil
		}
		return nil, err
	}
	response, err = fwSelectScopedReadObject(response, request, []fwScopedReadIdentity{
		{Source: "$.DBClusterId", Target: "DBClusterId", Value: dbClusterId},
	}, []string{"$.AdvancedDataPolicies.AdvancedDataPolicy[*].AutoCreated", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].BakType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestRegion", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DumpAction", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterKey", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterValue", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachDay", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachHour", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].PolicyId", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionValue", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcRegion", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcType", "$.AdvancedDataPolicies.AdvancedDataPolicy[*].StorageClass", "$.BackupFrequency", "$.BackupPolicyLevel", "$.BackupRetentionPolicyOnClusterDeletion", "$.DataLevel1BackupFrequency", "$.DataLevel1BackupPeriod", "$.DataLevel1BackupRetentionPeriod", "$.DataLevel1BackupTime", "$.DataLevel2BackupAnotherRegionRegion", "$.DataLevel2BackupAnotherRegionRetentionPeriod", "$.DataLevel2BackupPeriod", "$.DataLevel2BackupRetentionPeriod", "$.PreferredBackupPeriod", "$.PreferredBackupTime", "$.PreferredNextBackupTime"})
	if err != nil {
		return nil, fmt.Errorf("selecting DescribeBackupPolicy resource: %w", err)
	}
	if response == nil {
		return nil, nil
	}

	// F1: inject the stable BackupId (and its task alias BackupJobId, which is the same backup
	// set id) so the secondary DescribeBackups/DescribeBackupTasks reads filter on the exact
	// owned backup instead of merging every backup set in the cluster with "last wins".
	if backupId != "" {
		response["BackupId"] = backupId
		response["BackupJobId"] = backupId
	}
	return backupCompleteRead(client, response, dbClusterId, backupId, completeRead)
}

// backupQueryWindow returns a default UTC time window [StartTime, EndTime] for
// the DescribeBackups / DescribeDetachedBackups list APIs, which mandate
// StartTime in YYYY-MM-DDThh:mmZ (UTC, minute precision). DescribeBackupPolicy
// (the primary Get) does not return BackupStartTime/BackupEndTime, and response
// timestamps may use a different precision than the request format, so callers
// always supply this 7-day window ending now (the standard PolarDB backup
// retention period) to fetch the most recent backup sets; the response mappings
// then populate the real values into state.
func backupQueryWindow() (startTime, endTime string) {
	now := time.Now().UTC()
	return now.Add(-7 * 24 * time.Hour).Format("2006-01-02T15:04Z"), now.Format("2006-01-02T15:04Z")
}

func backupCompleteRead(client *connectivity.AliyunClient, object map[string]interface{}, dbClusterId, backupId string, includeSecondary bool) (map[string]interface{}, error) {
	object = fwNormalizeResponse(object, []fwMapping{
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].AutoCreated", Target: "$.AdvancedDataPolicies[*].AutoCreated", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].BakType", Target: "$.AdvancedDataPolicies[*].BakType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestRegion", Target: "$.AdvancedDataPolicies[*].DestRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DestType", Target: "$.AdvancedDataPolicies[*].DestType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].DumpAction", Target: "$.AdvancedDataPolicies[*].DumpAction", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterKey", Target: "$.AdvancedDataPolicies[*].FilterKey", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterType", Target: "$.AdvancedDataPolicies[*].FilterType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].FilterValue", Target: "$.AdvancedDataPolicies[*].FilterValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachDay", Target: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachDay", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].OnlyPreserveOneEachHour", Target: "$.AdvancedDataPolicies[*].OnlyPreserveOneEachHour", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].PolicyId", Target: "$.AdvancedDataPolicies[*].PolicyId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionType", Target: "$.AdvancedDataPolicies[*].RetentionType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].RetentionValue", Target: "$.AdvancedDataPolicies[*].RetentionValue", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcRegion", Target: "$.AdvancedDataPolicies[*].SrcRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].SrcType", Target: "$.AdvancedDataPolicies[*].SrcType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.AdvancedDataPolicies.AdvancedDataPolicy[*].StorageClass", Target: "$.AdvancedDataPolicies[*].StorageClass", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupFrequency", Target: "$.BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupPolicyLevel", Target: "$.BackupPolicyLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupRetentionPolicyOnClusterDeletion", Target: "$.BackupRetentionPolicyOnClusterDeletion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupFrequency", Target: "$.DataLevel1BackupFrequency", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupPeriod", Target: "$.DataLevel1BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupRetentionPeriod", Target: "$.DataLevel1BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel1BackupTime", Target: "$.DataLevel1BackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupAnotherRegionRegion", Target: "$.DataLevel2BackupAnotherRegionRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupAnotherRegionRetentionPeriod", Target: "$.DataLevel2BackupAnotherRegionRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupPeriod", Target: "$.DataLevel2BackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DataLevel2BackupRetentionPeriod", Target: "$.DataLevel2BackupRetentionPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PreferredBackupPeriod", Target: "$.PreferredBackupPeriod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PreferredBackupTime", Target: "$.PreferredBackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.PreferredNextBackupTime", Target: "$.PreferredNextBackupTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupEndTime", Target: "$.BackupEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupId", Target: "$.BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMethod", Target: "$.BackupMethod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMode", Target: "$.BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupSetSize", Target: "$.BackupSetSize", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupStartTime", Target: "$.BackupStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupType", Target: "$.BackupType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupsLevel", Target: "$.BackupsLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Comment", Target: "$.Comment", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ConsistentTime", Target: "$.ConsistentTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.IsAvail", Target: "$.IsAvail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Status", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupJobId", Target: "$.BackupJobId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupProgressStatus", Target: "$.BackupProgressStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.CreateTime", Target: "$.CreateTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.JobMode", Target: "$.JobMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Process", Target: "$.Process", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.TaskAction", Target: "$.TaskAction", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupEndTime", Target: "$.BackupEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupId", Target: "$.BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMethod", Target: "$.BackupMethod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupMode", Target: "$.BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupSetSize", Target: "$.BackupSetSize", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupStartTime", Target: "$.BackupStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupType", Target: "$.BackupType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.BackupsLevel", Target: "$.BackupsLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.ConsistentTime", Target: "$.ConsistentTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.IsAvail", Target: "$.IsAvail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.Status", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.StoreStatus", Target: "$.StoreStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		{Source: "$.DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
	})
	object["DBClusterId"] = dbClusterId
	// F1: keep the injected stable BackupId/BackupJobId so the secondary reads below filter on
	// the exact owned backup. They are overwritten with the materialized values when the backup
	// set is found; if it is never found the gate at the end reports the resource as gone.
	if backupId != "" {
		if _, ok := object["BackupId"]; !ok {
			object["BackupId"] = backupId
		}
		if _, ok := object["BackupJobId"]; !ok {
			object["BackupJobId"] = backupId
		}
	}
	if !includeSecondary {
		return object, nil
	}
	// F1: track whether the exact owned backup set is still present. DescribeBackupPolicy (the
	// primary Get) reflects the cluster policy and exists as long as the cluster does; the
	// backup-set identity is confirmed by DescribeBackups/DescribeDetachedBackups returning the
	// matching BackupId. If neither finds it, the resource is gone and Read must clear state.
	backupFound := false
	readValues0 := map[string]interface{}{}
	for key, value := range object {
		readValues0[key] = value
	}
	readValues0["RegionId"] = client.RegionId
	readEligible0, err := fwOperationEligible(readValues0, []fwDeleteCondition{})
	if err != nil {
		return nil, fmt.Errorf("selecting DescribeBackups read: %w", err)
	}
	if readEligible0 {
		// DescribeBackupPolicy does not return BackupStartTime/BackupEndTime,
		// and DescribeBackups mandates StartTime in YYYY-MM-DDThh:mmZ (UTC).
		// Always supply the fallback 7-day window so the most recent backup sets
		// are returned; the response mappings below populate the real values.
		start, end := backupQueryWindow()
		readValues0["BackupStartTime"] = start
		readValues0["BackupEndTime"] = end
		readRequest0 := map[string]interface{}{}
		if err := fwApplyMappings(readRequest0, readValues0, []fwMapping{
			{Source: "$.BackupEndTime", Target: "EndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupId", Target: "BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupMode", Target: "BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupStartTime", Target: "StartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "BackupRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Status", Target: "BackupStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			return nil, err
		}
		readResponse0, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeBackups", readRequest0, 1*time.Minute)
		if err != nil {
			return nil, err
		}
		readObject0 := fwNormalizeResponse(readResponse0, []fwMapping{
			{Source: "$.Items.Backup[*].BackupEndTime", Target: "$.BackupEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupId", Target: "$.BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupMethod", Target: "$.BackupMethod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupMode", Target: "$.BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupSetSize", Target: "$.BackupSetSize", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupStartTime", Target: "$.BackupStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupType", Target: "$.BackupType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupsLevel", Target: "$.BackupsLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].Comment", Target: "$.Comment", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].ConsistentTime", Target: "$.ConsistentTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].IsAvail", Target: "$.IsAvail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupStatus", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		})
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupEndTime"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupEndTime"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupId"), nil); len(entries) > 0 {
			backupFound = true
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupId"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupMethod"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupMethod"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupMode"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupMode"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupSetSize"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupSetSize"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupStartTime"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupStartTime"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupType"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupType"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.BackupsLevel"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupsLevel"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.Comment"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.Comment"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.ConsistentTime"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.ConsistentTime"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.DBClusterId"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.DBClusterId"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.IsAvail"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.IsAvail"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject0, fwMappingPath("$.Status"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.Status"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
	}
	readValues1 := map[string]interface{}{}
	for key, value := range object {
		readValues1[key] = value
	}
	readValues1["RegionId"] = client.RegionId
	readEligible1, err := fwOperationEligible(readValues1, []fwDeleteCondition{})
	if err != nil {
		return nil, fmt.Errorf("selecting DescribeBackupTasks read: %w", err)
	}
	if readEligible1 {
		readRequest1 := map[string]interface{}{}
		if err := fwApplyMappings(readRequest1, readValues1, []fwMapping{
			{Source: "$.BackupJobId", Target: "BackupJobId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupMode", Target: "BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "RegionId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			return nil, err
		}
		readResponse1, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeBackupTasks", readRequest1, 1*time.Minute)
		if err != nil {
			return nil, err
		}
		readObject1 := fwNormalizeResponse(readResponse1, []fwMapping{
			{Source: "$.Items.BackupJob[*].BackupJobId", Target: "$.BackupJobId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.BackupJob[*].BackupProgressStatus", Target: "$.BackupProgressStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.BackupJob[*].StartTime", Target: "$.CreateTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.BackupJob[*].JobMode", Target: "$.JobMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.BackupJob[*].Process", Target: "$.Process", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.BackupJob[*].TaskAction", Target: "$.TaskAction", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		})
		if entries := fwWalkMappedPath(readObject1, fwMappingPath("$.BackupJobId"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupJobId"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject1, fwMappingPath("$.BackupProgressStatus"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupProgressStatus"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject1, fwMappingPath("$.CreateTime"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.CreateTime"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject1, fwMappingPath("$.JobMode"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.JobMode"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject1, fwMappingPath("$.Process"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.Process"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject1, fwMappingPath("$.TaskAction"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.TaskAction"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
	}
	readValues2 := map[string]interface{}{}
	for key, value := range object {
		readValues2[key] = value
	}
	readValues2["RegionId"] = client.RegionId
	readEligible2, err := fwOperationEligible(readValues2, []fwDeleteCondition{})
	if err != nil {
		return nil, fmt.Errorf("selecting DescribeDetachedBackups read: %w", err)
	}
	if readEligible2 {
		// DescribeDetachedBackups also mandates StartTime in YYYY-MM-DDThh:mmZ.
		// Use the same fallback 7-day window regardless of any value read above,
		// since response timestamps may use a different precision than the
		// request format.
		start, end := backupQueryWindow()
		readValues2["BackupStartTime"] = start
		readValues2["BackupEndTime"] = end
		readRequest2 := map[string]interface{}{}
		if err := fwApplyMappings(readRequest2, readValues2, []fwMapping{
			{Source: "$.BackupEndTime", Target: "EndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupId", Target: "BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupMode", Target: "BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.BackupStartTime", Target: "StartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.DBClusterId", Target: "DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.RegionId", Target: "BackupRegion", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Status", Target: "BackupStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		}, false); err != nil {
			return nil, err
		}
		readResponse2, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeDetachedBackups", readRequest2, 1*time.Minute)
		if err != nil {
			// A live cluster has no detached backups; DescribeDetachedBackups
			// reports InvalidDBClusterId.NotFound in that case (its generic
			// message references DBInstanceId, but the request is keyed on
			// DBClusterId). Treat it as an empty result so the read continues
			// with the values populated by DescribeBackups above.
			if !isExpectedErrors(err, []string{"InvalidDBClusterId.NotFound"}) {
				return nil, err
			}
			readResponse2 = map[string]interface{}{}
		}
		readObject2 := fwNormalizeResponse(readResponse2, []fwMapping{
			{Source: "$.Items.Backup[*].BackupEndTime", Target: "$.BackupEndTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupId", Target: "$.BackupId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupMethod", Target: "$.BackupMethod", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupMode", Target: "$.BackupMode", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupSetSize", Target: "$.BackupSetSize", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupStartTime", Target: "$.BackupStartTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupType", Target: "$.BackupType", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupsLevel", Target: "$.BackupsLevel", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].ConsistentTime", Target: "$.ConsistentTime", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].DBClusterId", Target: "$.DBClusterId", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].IsAvail", Target: "$.IsAvail", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].BackupStatus", Target: "$.Status", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
			{Source: "$.Items.Backup[*].StoreStatus", Target: "$.StoreStatus", Encoding: "", Constant: "", Repeat: false, Values: map[string]string{}},
		})
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupEndTime"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupEndTime"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupId"), nil); len(entries) > 0 {
			backupFound = true
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupId"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupMethod"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupMethod"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupMode"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupMode"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupSetSize"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupSetSize"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupStartTime"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupStartTime"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupType"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupType"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.BackupsLevel"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.BackupsLevel"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.ConsistentTime"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.ConsistentTime"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.DBClusterId"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.DBClusterId"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.IsAvail"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.IsAvail"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.Status"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.Status"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
		if entries := fwWalkMappedPath(readObject2, fwMappingPath("$.StoreStatus"), nil); len(entries) > 0 {
			for _, entry := range entries {
				object = fwSetMappedPath(object, fwMappingPath("$.StoreStatus"), entry.Indexes, entry.Value).(map[string]interface{})
			}
		}
	}
	// F1: a composite identity (DBClusterId:BackupId) means the resource tracks one specific
	// backup set. The cluster policy can still exist while the owned backup was externally
	// deleted; in that case neither DescribeBackups nor DescribeDetachedBackups returned the
	// matching BackupId, so report the resource as gone so Read clears state and Delete is
	// idempotent instead of operating on a stale/foreign backup.
	if backupId != "" && !backupFound {
		return nil, nil
	}
	return object, nil
}

// backupTaskList extracts the BackupJob task entries from a DescribeBackupTasks response.
func backupTaskList(resp map[string]interface{}) []map[string]interface{} {
	out := []map[string]interface{}{}
	items, _ := resp["Items"].(map[string]interface{})
	jobs, _ := items["BackupJob"].([]interface{})
	for _, j := range jobs {
		if m, ok := j.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// backupBackupsList extracts the backup-set entries from a DescribeBackups response.
func backupBackupsList(resp map[string]interface{}) []map[string]interface{} {
	out := []map[string]interface{}{}
	items, _ := resp["Items"].(map[string]interface{})
	backs, _ := items["Backup"].([]interface{})
	for _, b := range backs {
		if m, ok := b.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

// backupIdGreater reports whether backup id a is greater than b, comparing as integers
// (PolarDB BackupIds are sequential); falls back to lexical comparison on parse failure.
func backupIdGreater(a, b string) bool {
	ai, ea := strconv.ParseInt(a, 10, 64)
	bi, eb := strconv.ParseInt(b, 10, 64)
	if ea == nil && eb == nil {
		return ai > bi
	}
	return a > b
}

// backupTaskRefreshFunc polls the CreateBackup task via DescribeBackupTasks(BackupJobId) until
// BackupProgressStatus reaches Finished (the terminal success state). CreateBackup only submits
// the background task; the backup set (BackupId) is not stable until the task finishes. Status
// values: NoStart/Preparing/Waiting/Uploading/Checking/Finished.
func backupTaskRefreshFunc(client *connectivity.AliyunClient, dbClusterId, backupJobId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{
			"DBClusterId": dbClusterId,
			"BackupJobId": backupJobId,
			"RegionId":    client.RegionId,
		}
		resp, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeBackupTasks", req, 1*time.Minute)
		if err != nil {
			return nil, "", err
		}
		for _, j := range backupTaskList(resp) {
			if fmt.Sprint(j["BackupJobId"]) == backupJobId {
				return j, fmt.Sprint(j["BackupProgressStatus"]), nil
			}
		}
		// Task not listed yet (not registered or between states): keep pending.
		return nil, "", nil
	}
}

// backupWaitForBackupTask waits for the CreateBackup task (BackupJobId) to reach Finished.
// F2: CreateBackup only submits the task and returns BackupJobId; the backup SET (BackupId)
// is not queryable until the task finishes. This bounds the create wait on the real task
// state so a still-running task cannot let Create return with an unfinished backup, and a
// timeout surfaces a diagnostic. Timeout keeps the resource recoverable by its composite id.
func backupWaitForBackupTask(client *connectivity.AliyunClient, dbClusterId, backupJobId string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"", "NoStart", "Preparing", "Waiting", "Uploading", "Checking"},
		Target:     []string{"Finished"},
		Refresh:    backupTaskRefreshFunc(client, dbClusterId, backupJobId),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err := stateConf.WaitForState()
	return err
}

// backupResolveBackupId finds the real backup SET BackupId produced by a finished backup task.
// DescribeBackups/DeleteBackup key on BackupId (the backup set id), which is NOT the
// CreateBackup BackupJobId; DescribeBackups entries do not expose BackupJobId, so the set is
// correlated by BackupJobId when the raw response includes it, else by Comment (the CreateBackup
// Comment is stored on the backup set), else by the most recent backup set (largest BackupId).
// DescribeBackups may lag slightly behind the task's Finished state, so the lookup retries
// briefly before giving up.
func backupResolveBackupId(client *connectivity.AliyunClient, dbClusterId, backupJobId, comment string) (string, error) {
	// Up to ~2 min of retries for DescribeBackups to list the newly completed backup set.
	for i := 0; i < 24; i++ {
		bid, found, err := backupLookupBackupId(client, dbClusterId, backupJobId, comment)
		if err != nil {
			return "", err
		}
		if found && bid != "" {
			return bid, nil
		}
		time.Sleep(5 * time.Second)
	}
	return "", nil
}

// backupLookupBackupId performs a single DescribeBackups list and correlates the owned backup
// set. Returns (backupId, found, err); found is true only when a reliable match (BackupJobId or
// Comment) is found, or - when no Comment is configured - when any backup set exists (the most
// recent one, since the task just Finished).
func backupLookupBackupId(client *connectivity.AliyunClient, dbClusterId, backupJobId, comment string) (string, bool, error) {
	start, end := backupQueryWindow()
	req := map[string]interface{}{
		"DBClusterId":  dbClusterId,
		"StartTime":    start,
		"EndTime":      end,
		"BackupRegion": client.RegionId,
	}
	resp, err := retryRpcPost(client, "polardb", "2017-08-01", "DescribeBackups", req, 1*time.Minute)
	if err != nil {
		return "", false, err
	}
	entries := backupBackupsList(resp)
	// 1. correlate by BackupJobId when the raw response carries it (undocumented but often present)
	for _, e := range entries {
		jid := fmt.Sprint(e["BackupJobId"])
		if jid != "" && jid != "<nil>" && jid == backupJobId {
			if bid := fmt.Sprint(e["BackupId"]); bid != "" && bid != "<nil>" {
				return bid, true, nil
			}
		}
	}
	// 2. correlate by Comment (the CreateBackup Comment is stored on the backup set)
	if comment != "" {
		for _, e := range entries {
			if fmt.Sprint(e["Comment"]) == comment {
				if bid := fmt.Sprint(e["BackupId"]); bid != "" && bid != "<nil>" {
					return bid, true, nil
				}
			}
		}
		return "", false, nil // Comment configured but not yet visible: keep retrying.
	}
	// 3. no Comment to correlate on: the most recent backup set (largest BackupId). Only
	// reliable as "ours" once the task is Finished; any non-empty result is accepted.
	var best string
	for _, e := range entries {
		bid := fmt.Sprint(e["BackupId"])
		if bid == "" || bid == "<nil>" {
			continue
		}
		if best == "" || backupIdGreater(bid, best) {
			best = bid
		}
	}
	return best, best != "", nil
}

// A required exact request can identify a singleton configuration whose API omits its ID.
func fwSelectScopedReadObject(raw interface{}, request map[string]interface{}, identities []fwScopedReadIdentity, payloadPaths []string) (map[string]interface{}, error) {
	if len(identities) == 0 {
		return nil, fmt.Errorf("no exact request identity was declared")
	}
	for _, identity := range identities {
		value, exists := request[identity.Target].(string)
		if !exists || identity.Value == "" || value != identity.Value {
			return nil, fmt.Errorf("singleton request does not match the resource identity")
		}
	}
	object, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("singleton response is not an object")
	}
	for _, identity := range identities {
		values := fwWalkMappedPath(object, fwMappingPath(identity.Source), nil)
		if len(values) > 1 {
			return nil, fmt.Errorf("singleton response has ambiguous identity")
		}
		if len(values) == 1 && (values[0].Value == nil || fmt.Sprint(values[0].Value) != identity.Value) {
			return nil, fmt.Errorf("singleton response identity does not match the request")
		}
	}
	for _, path := range payloadPaths {
		values := fwWalkMappedPath(object, fwMappingPath(path), nil)
		for _, value := range values {
			if value.Value != nil {
				return object, nil
			}
		}
	}
	return nil, fmt.Errorf("singleton response has no declared configuration payload")
}

// Optional update parameters only participate when their mapped value changed.
func fwChangedUpdateMappings(bindings []fwMapping, desired, prior map[string]interface{}, requiredSources []string) []fwMapping {
	selected := make([]fwMapping, 0, len(bindings))
	for _, binding := range bindings {
		required := false
		for _, source := range requiredSources {
			if binding.Source == source || strings.HasPrefix(binding.Source, source+".") || strings.HasPrefix(binding.Source, source+"[*]") {
				required = true
				break
			}
		}
		if required || !fwValuesEqual(fwWalkMappedPath(desired, fwMappingPath(binding.Source), nil), fwWalkMappedPath(prior, fwMappingPath(binding.Source), nil)) {
			selected = append(selected, binding)
		}
	}
	return selected
}

// fwStr navigates a dotted key in m and returns a types.String (missing -> StringNull).
// Used by nested-list ObjectValues helpers to read scalar sub-fields from each list item;
// the top-level State closures in resource.go capture the Get response, so they cannot read
// per-item sub-maps.
func fwStr(m map[string]interface{}, key string) types.String {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.StringNull()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.StringNull()
		}
	}
	if v, ok := cur.(string); ok {
		return types.StringValue(v)
	}
	return types.StringNull()
}

func fwBool(m map[string]interface{}, key string) types.Bool {
	cur := interface{}(m)
	for _, seg := range strings.Split(key, ".") {
		mm, ok := cur.(map[string]interface{})
		if !ok {
			return types.BoolNull()
		}
		cur, ok = mm[seg]
		if !ok {
			return types.BoolNull()
		}
	}
	if v, ok := cur.(bool); ok {
		return types.BoolValue(v)
	}
	return types.BoolNull()
}

// backupAdvancedDataPoliciesToTypeList wraps the API object array into a framework
// types.List of nested objects (types.Object). nil/empty/non-array -> empty list
// (.#=0) to match v1 d.Set(..., []), so TestCheckResourceAttrSet("x.#") passes and
// UseStateForUnknown unknown->0 stays consistent (ListNull leaves .# unset, failing
// the .# check).
func backupAdvancedDataPoliciesToTypeList(raw interface{}) types.List {
	objType := backupAdvancedDataPoliciesObjType()
	elems := make([]attr.Value, 0)
	// Handle both array (nested list) and single object (MaxItems:1 object type, gap#7).
	// F6: when the API returns AdvancedDataPolicies:{AdvancedDataPolicy:[]}, the leaf mappings
	// produce zero entries and raw arrives as the wrapper map. Unwrap the inner
	// "AdvancedDataPolicy" container so an empty array stays an empty list instead of being
	// wrapped into a single all-null block (state item count 0 -> 1).
	if raw != nil {
		var arr []interface{}
		switch v := raw.(type) {
		case []interface{}:
			arr = v
		case map[string]interface{}:
			if inner, ok := v["AdvancedDataPolicy"].([]interface{}); ok {
				arr = inner
			} else {
				arr = []interface{}{v}
			}
		}
		for _, item := range arr {
			m, _ := item.(map[string]interface{})
			obj, _ := types.ObjectValue(backupAdvancedDataPoliciesAttrTypes(), backupAdvancedDataPoliciesObjectValues(m))
			elems = append(elems, obj)
		}
	}
	l, _ := types.ListValue(objType, elems)
	return l
}

// backupAdvancedDataPoliciesAttrTypes returns the framework attr types for each nested field.
func backupAdvancedDataPoliciesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"action_type":                 types.StringType,
		"only_preserve_one_each_hour": types.BoolType,
		"storage_class":               types.StringType,
		"retention_value":             types.StringType,
		"bak_type":                    types.StringType,
		"filter_key":                  types.StringType,
		"filter_value":                types.StringType,
		"dest_type":                   types.StringType,
		"filter_type":                 types.StringType,
		"src_region":                  types.StringType,
		"auto_created":                types.BoolType,
		"only_preserve_one_each_day":  types.BoolType,
		"src_type":                    types.StringType,
		"dump_action":                 types.StringType,
		"dest_region":                 types.StringType,
		"policy_id":                   types.StringType,
		"retention_type":              types.StringType,
	}
}

// backupAdvancedDataPoliciesObjType is the element ObjectType for the nested list.
func backupAdvancedDataPoliciesObjType() types.ObjectType {
	return types.ObjectType{AttrTypes: backupAdvancedDataPoliciesAttrTypes()}
}

// backupAdvancedDataPoliciesObjectValues maps a single API object (map) to framework attr
// values. Nested sub-lists recurse via their own ToTypeList helpers; scalars use fwStr/fwBool/fwInt.
func backupAdvancedDataPoliciesObjectValues(m map[string]interface{}) map[string]attr.Value {
	out := map[string]attr.Value{}
	out["action_type"] = types.StringNull()
	out["only_preserve_one_each_hour"] = fwBool(m, "OnlyPreserveOneEachHour")
	out["storage_class"] = fwStr(m, "StorageClass")
	out["retention_value"] = fwStr(m, "RetentionValue")
	out["bak_type"] = fwStr(m, "BakType")
	out["filter_key"] = fwStr(m, "FilterKey")
	out["filter_value"] = fwStr(m, "FilterValue")
	out["dest_type"] = fwStr(m, "DestType")
	out["filter_type"] = fwStr(m, "FilterType")
	out["src_region"] = fwStr(m, "SrcRegion")
	out["auto_created"] = fwBool(m, "AutoCreated")
	out["only_preserve_one_each_day"] = fwBool(m, "OnlyPreserveOneEachDay")
	out["src_type"] = fwStr(m, "SrcType")
	out["dump_action"] = fwStr(m, "DumpAction")
	out["dest_region"] = fwStr(m, "DestRegion")
	out["policy_id"] = fwStr(m, "PolicyId")
	out["retention_type"] = fwStr(m, "RetentionType")
	return out
}

// fwWireField preserves API/source names independently of Terraform schema names.

// fwPlannedChanged ignores computed unknown leaves while comparing configured nested values.
func fwPlannedChanged(plan, state attr.Value) bool {
	if plan == nil || plan.IsUnknown() {
		return false
	}
	if state == nil || state.IsNull() || state.IsUnknown() || plan.IsNull() {
		return state == nil || !plan.Equal(state)
	}
	switch desired := plan.(type) {
	case types.Object:
		previous, ok := state.(types.Object)
		if !ok {
			return true
		}
		oldItems := previous.Attributes()
		for key, value := range desired.Attributes() {
			if fwPlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	case types.List:
		previous, ok := state.(types.List)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for index, value := range items {
			if fwPlannedChanged(value, oldItems[index]) {
				return true
			}
		}
		return false
	case types.Set:
		previous, ok := state.(types.Set)
		if !ok {
			return true
		}
		return !fwSetMatches(desired.Elements(), previous.Elements(), func(a, b attr.Value) bool { return !fwPlannedChanged(a, b) })
	case types.Map:
		previous, ok := state.(types.Map)
		if !ok {
			return true
		}
		items, oldItems := desired.Elements(), previous.Elements()
		if len(items) != len(oldItems) {
			return true
		}
		for key, value := range items {
			if fwPlannedChanged(value, oldItems[key]) {
				return true
			}
		}
		return false
	}
	return !plan.Equal(state)
}

// fwSetMatches pairs unordered members without depending on response or plan order.
func fwSetMatches(left, right []attr.Value, compatible func(attr.Value, attr.Value) bool) bool {
	if len(left) != len(right) {
		return false
	}
	owners := make([]int, len(right))
	for index := range owners {
		owners[index] = -1
	}
	var assign func(int, []bool) bool
	assign = func(index int, seen []bool) bool {
		for candidate := range right {
			if seen[candidate] || !compatible(left[index], right[candidate]) {
				continue
			}
			seen[candidate] = true
			if owners[candidate] == -1 || assign(owners[candidate], seen) {
				owners[candidate] = index
				return true
			}
		}
		return false
	}
	for index := range left {
		if !assign(index, make([]bool, len(right))) {
			return false
		}
	}
	return true
}

func fwDeleteEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	return fwOperationEligible(values, conditions)
}

func fwOperationEligible(values map[string]interface{}, conditions []fwDeleteCondition) (bool, error) {
	eligible := true
	for _, condition := range conditions {
		found := fwWalkMappedPath(values, fwMappingPath(condition.Property), nil)
		if len(found) == 0 || found[0].Value == nil {
			return false, fmt.Errorf("cannot select an operation: %s is unknown", condition.Property)
		}
		matches := false
		for _, expected := range condition.Values {
			if fmt.Sprint(found[0].Value) == expected {
				matches = true
				break
			}
		}
		switch condition.Action {
		case "valid", "required":
			eligible = eligible && matches
		case "invalid":
			eligible = eligible && !matches
		default:
			return false, fmt.Errorf("unsupported operation dependency action: %s", condition.Action)
		}
	}
	return eligible, nil
}

func fwKnownValue(value attr.Value) attr.Value {
	switch v := value.(type) {
	case types.String:
		if v.IsUnknown() {
			return types.StringNull()
		}
	case types.Bool:
		if v.IsUnknown() {
			return types.BoolNull()
		}
	case types.Int64:
		if v.IsUnknown() {
			return types.Int64Null()
		}
	case types.Float64:
		if v.IsUnknown() {
			return types.Float64Null()
		}
	case types.List:
		if v.IsUnknown() {
			return types.ListNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, fwKnownValue(item))
			}
			return types.ListValueMust(v.ElementType(context.Background()), items)
		}
	case types.Set:
		if v.IsUnknown() {
			return types.SetNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := []attr.Value{}
			for _, item := range v.Elements() {
				items = append(items, fwKnownValue(item))
			}
			return types.SetValueMust(v.ElementType(context.Background()), items)
		}
	case types.Map:
		if v.IsUnknown() {
			return types.MapNull(v.ElementType(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Elements() {
				items[key] = fwKnownValue(item)
			}
			return types.MapValueMust(v.ElementType(context.Background()), items)
		}
	case types.Object:
		if v.IsUnknown() {
			return types.ObjectNull(v.AttributeTypes(context.Background()))
		}
		if !v.IsNull() {
			items := map[string]attr.Value{}
			for key, item := range v.Attributes() {
				items[key] = fwKnownValue(item)
			}
			return types.ObjectValueMust(v.AttributeTypes(context.Background()), items)
		}
	}
	return value
}

func fwValueToInterface(value attr.Value, fields map[string]fwWireField, single bool) interface{} {
	if value == nil || value.IsNull() || value.IsUnknown() {
		return nil
	}
	switch v := value.(type) {
	case interface{ ValueString() string }:
		return v.ValueString()
	case types.Bool:
		return v.ValueBool()
	case types.Int64:
		return v.ValueInt64()
	case types.Float64:
		return v.ValueFloat64()
	case types.Object:
		out := map[string]interface{}{}
		for key, child := range v.Attributes() {
			field, mapped := fields[key]
			if !mapped {
				field.Name = key
			}
			if raw := fwValueToInterface(child, field.Fields, field.Single); raw != nil {
				out[field.Name] = raw
			}
		}
		return out
	case types.Map:
		out := map[string]interface{}{}
		for key, child := range v.Elements() {
			if raw := fwValueToInterface(child, nil, false); raw != nil {
				out[key] = raw
			}
		}
		return out
	case types.List:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := fwValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		if single {
			if len(out) == 0 {
				return nil
			}
			return out[0]
		}
		return out
	case types.Set:
		out := make([]interface{}, 0, len(v.Elements()))
		for _, child := range v.Elements() {
			if raw := fwValueToInterface(child, fields, false); raw != nil {
				out = append(out, raw)
			}
		}
		return out
	}
	return nil
}

func fwMappingPath(path string) []string {
	path = strings.TrimPrefix(path, "$.")
	path = strings.TrimPrefix(path, "$")
	path = strings.ReplaceAll(path, "[*]", ".*")
	if path == "" {
		return nil
	}
	return strings.Split(strings.TrimPrefix(path, "."), ".")
}

func fwWalkMappedPath(value interface{}, path []string, indexes []int) []fwMappedValue {
	if len(path) == 0 {
		return []fwMappedValue{{Value: value, Indexes: indexes}}
	}
	if path[0] == "*" {
		var out []fwMappedValue
		if items, ok := value.([]interface{}); ok {
			for i, item := range items {
				next := append(append([]int{}, indexes...), i)
				out = append(out, fwWalkMappedPath(item, path[1:], next)...)
			}
		}
		return out
	}
	if m, ok := value.(map[string]interface{}); ok {
		if child, exists := m[path[0]]; exists {
			return fwWalkMappedPath(child, path[1:], indexes)
		}
	}
	return nil
}

func fwSetMappedPath(value interface{}, path []string, indexes []int, raw interface{}) interface{} {
	if len(path) == 0 {
		return raw
	}
	if path[0] == "*" {
		if len(indexes) == 0 {
			return value
		}
		items, _ := value.([]interface{})
		for len(items) <= indexes[0] {
			items = append(items, nil)
		}
		items[indexes[0]] = fwSetMappedPath(items[indexes[0]], path[1:], indexes[1:], raw)
		return items
	}
	m, _ := value.(map[string]interface{})
	if m == nil {
		m = map[string]interface{}{}
	}
	m[path[0]] = fwSetMappedPath(m[path[0]], path[1:], indexes, raw)
	return m
}

func fwApplyMappings(request, values map[string]interface{}, bindings []fwMapping, roa bool) error {
	mapped := map[string]interface{}{}
	repeated := map[string]bool{}
	for _, binding := range bindings {
		source, target := binding.Source, binding.Target
		if binding.Encoding == "comma" || binding.Encoding == "first" || binding.Encoding == "wrap-json" || binding.Encoding == "wrap-array" {
			source = strings.TrimSuffix(source, "[*]")
			target = strings.TrimSuffix(target, "[*]")
		}
		entries := fwWalkMappedPath(values, fwMappingPath(source), nil)
		if binding.Encoding == "scoped-const" {
			parent := source[:strings.LastIndex(source, ".")]
			entries = fwWalkMappedPath(values, fwMappingPath(parent), nil)
		}
		if strings.HasPrefix(binding.Encoding, "concat|") {
			parts := strings.Split(binding.Encoding, "|")
			fields := []string{}
			complete := true
			for _, propertyPath := range parts[2:] {
				found := fwWalkMappedPath(values, fwMappingPath(propertyPath), nil)
				if len(found) == 0 || found[0].Value == nil {
					complete = false
					break
				}
				fields = append(fields, fmt.Sprint(found[0].Value))
			}
			entries = nil
			if complete {
				entries = []fwMappedValue{{Value: strings.Join(fields, parts[1])}}
			}
		}
		if strings.Contains(source, "[*]") && !strings.Contains(target, "[*]") && binding.Encoding == "" {
			items := []interface{}{}
			for _, entry := range entries {
				if entry.Value != nil {
					items = append(items, entry.Value)
				}
			}
			if len(entries) > 0 {
				entries = []fwMappedValue{{Value: items}}
			}
		}
		concreteEntries := entries[:0]
		for _, entry := range entries {
			if entry.Value != nil {
				concreteEntries = append(concreteEntries, entry)
			}
		}
		entries = concreteEntries
		if binding.Constant != "" && !(binding.Encoding == "default" && len(entries) > 0) {
			if (binding.Encoding == "key" || binding.Encoding == "scoped-const") && len(entries) == 0 {
				continue
			}
			var raw interface{}
			if err := json.Unmarshal([]byte(binding.Constant), &raw); err != nil {
				return err
			}
			if len(entries) == 0 {
				entries = []fwMappedValue{{Value: raw}}
			} else {
				for index := range entries {
					entries[index].Value = raw
				}
			}
		}
		if strings.HasPrefix(binding.Encoding, "uniform-scalar|") && len(entries) > 0 {
			scalarType := strings.TrimPrefix(binding.Encoding, "uniform-scalar|")
			for _, entry := range entries {
				valid := false
				switch scalarType {
				case "boolean":
					_, valid = entry.Value.(bool)
				case "integer":
					switch value := entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32:
						valid = true
					case float64:
						valid = value == float64(int64(value))
					case json.Number:
						_, err := value.Int64()
						valid = err == nil
					}
				case "number":
					switch entry.Value.(type) {
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, json.Number:
						valid = true
					}
				}
				if !valid {
					return fmt.Errorf("configured value does not match scalar API parameter %s type %s", target, scalarType)
				}
			}
			first, err := json.Marshal(entries[0].Value)
			if err != nil {
				return err
			}
			for _, entry := range entries[1:] {
				next, err := json.Marshal(entry.Value)
				if err != nil {
					return err
				}
				if string(next) != string(first) {
					return fmt.Errorf("conflicting configured values for scalar API parameter %s", target)
				}
			}
			entries = entries[:1]
		}
		for _, entry := range entries {
			raw := entry.Value
			if raw == nil {
				continue
			}
			if mappedValue, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mappedValue
			}
			if strings.HasPrefix(binding.Encoding, "suffix|") {
				raw = fmt.Sprint(raw) + strings.TrimPrefix(binding.Encoding, "suffix|")
			}
			switch binding.Encoding {
			case "first":
				if list, ok := raw.([]interface{}); ok {
					if len(list) == 0 {
						continue
					}
					raw = list[0]
				}
			case "comma":
				if list, ok := raw.([]interface{}); ok {
					parts := make([]string, len(list))
					for i, item := range list {
						parts[i] = fmt.Sprint(item)
					}
					raw = strings.Join(parts, ",")
				}
			case "wrap-array":
				raw = []interface{}{raw}
			case "wrap-json":
				encoded, err := json.Marshal([]interface{}{raw})
				if err != nil {
					return err
				}
				raw = string(encoded)
			}
			indexes := append([]int(nil), entry.Indexes...)
			if len(binding.IndexSources) > 0 {
				if len(binding.IndexSources) != strings.Count(target, "[*]") {
					return fmt.Errorf("request wildcard alignment has an unexpected target dimension count")
				}
				indexes = make([]int, len(binding.IndexSources))
				for position, sourceIndex := range binding.IndexSources {
					if sourceIndex == -1 {
						continue
					}
					if sourceIndex < 0 || sourceIndex >= len(entry.Indexes) {
						return fmt.Errorf("request wildcard alignment is missing a source dimension")
					}
					indexes[position] = entry.Indexes[sourceIndex]
				}
			} else {
				for len(indexes) < strings.Count(target, "[*]") {
					indexes = append(indexes, 0)
				}
			}
			mapped = fwSetMappedPath(mapped, fwMappingPath(target), indexes, raw).(map[string]interface{})
		}
		if binding.Repeat {
			if parts := fwMappingPath(target); len(parts) > 0 {
				repeated[parts[0]] = true
			}
		}
	}
	if err := fwCompactKvPairs(mapped, bindings); err != nil {
		return err
	}
	for key, raw := range mapped {
		if roa {
			request[key] = raw
			continue
		}
		switch raw.(type) {
		case map[string]interface{}, []interface{}:
			if repeated[key] {
				fwFlattenRequest(request, key, raw)
			} else {
				encoded, err := json.Marshal(raw)
				if err != nil {
					return err
				}
				request[key] = string(encoded)
			}
		default:
			request[key] = raw
		}
	}
	return nil
}

// Only explicitly declared kvPairs groups are compacted; ordinary indexed lists retain shape.
func fwCompactKvPairs(mapped map[string]interface{}, bindings []fwMapping) error {
	groups := map[string][2]string{}
	for _, binding := range bindings {
		if binding.Group == "" {
			continue
		}
		fields := [2]string{binding.GroupKey, binding.GroupValue}
		if fields[0] == "" || fields[1] == "" || fields[0] == fields[1] {
			return fmt.Errorf("Invalid kvPairs member declaration")
		}
		if previous, exists := groups[binding.Group]; exists && previous != fields {
			return fmt.Errorf("Conflicting kvPairs member declaration")
		}
		groups[binding.Group] = fields
	}
	for group, fields := range groups {
		raw, exists := mapped[group]
		if !exists {
			continue
		}
		members, ok := raw.(map[string]interface{})
		if !ok {
			return fmt.Errorf("kvPairs group is not an indexed object")
		}
		indexes := []int{}
		for key := range members {
			index, err := strconv.Atoi(key)
			if err != nil || index < 1 || strconv.Itoa(index) != key {
				return fmt.Errorf("Invalid kvPairs member index")
			}
			indexes = append(indexes, index)
		}
		sort.Ints(indexes)
		compact := map[string]interface{}{}
		codes := map[string]bool{}
		for position, index := range indexes {
			pair, valid := members[strconv.Itoa(index)].(map[string]interface{})
			if !valid || len(pair) != 2 || pair[fields[0]] == nil || pair[fields[1]] == nil {
				return fmt.Errorf("Incomplete kvPairs member")
			}
			code, valid := pair[fields[0]].(string)
			if !valid || code == "" || codes[code] {
				return fmt.Errorf("Invalid or duplicate kvPairs code")
			}
			codes[code] = true
			compact[strconv.Itoa(position+1)] = pair
		}
		mapped[group] = compact
	}
	return nil
}

func fwFlattenRequest(request map[string]interface{}, key string, raw interface{}) {
	switch v := raw.(type) {
	case []interface{}:
		for i, item := range v {
			fwFlattenRequest(request, fmt.Sprintf("%s.%d", key, i+1), item)
		}
	case map[string]interface{}:
		for name, item := range v {
			fwFlattenRequest(request, key+"."+name, item)
		}
	default:
		if v != nil {
			request[key] = v
		}
	}
}

func fwNormalizeResponse(object map[string]interface{}, bindings []fwMapping) map[string]interface{} {
	// Whole-object mappings must precede their leaf aliases, regardless of AM order.
	// Stable insertion order preserves same-depth precedence and source array indexes.
	bindings = append([]fwMapping(nil), bindings...)
	for index := 1; index < len(bindings); index++ {
		for current := index; current > 0 && len(fwMappingPath(bindings[current].Target)) < len(fwMappingPath(bindings[current-1].Target)); current-- {
			bindings[current], bindings[current-1] = bindings[current-1], bindings[current]
		}
	}
	out := map[string]interface{}{}
	for key, value := range object {
		out[key] = value
	}
	for _, binding := range bindings {
		if strings.HasSuffix(binding.Source, "[*]") && strings.HasSuffix(binding.Target, "[*]") {
			source := strings.TrimSuffix(binding.Source, "[*]")
			target := strings.TrimSuffix(binding.Target, "[*]")
			for _, container := range fwWalkMappedPath(object, fwMappingPath(source), nil) {
				if values, ok := container.Value.([]interface{}); ok && len(values) == 0 {
					out = fwSetMappedPath(out, fwMappingPath(target), container.Indexes, []interface{}{}).(map[string]interface{})
				}
			}
		}
		for _, entry := range fwWalkMappedPath(object, fwMappingPath(binding.Source), nil) {
			raw := entry.Value
			target := binding.Target
			if binding.Encoding == "split-comma" {
				target = strings.TrimSuffix(target, "[*]")
				items := []interface{}{}
				if text, ok := raw.(string); ok && text != "" {
					for _, part := range strings.Split(text, ",") {
						items = append(items, strings.TrimSpace(part))
					}
				}
				raw = items
			}
			if binding.Encoding == "scalar-string" && raw != nil {
				switch raw.(type) {
				case string:
				case int, int32, int64, float32, float64, json.Number:
					raw = fmt.Sprint(raw)
				}
			}
			if mapped, ok := binding.Values[fmt.Sprint(raw)]; ok {
				raw = mapped
			}
			out = fwSetMappedPath(out, fwMappingPath(target), entry.Indexes, raw).(map[string]interface{})
		}
	}
	return out
}

func fwValuesEqual(a, b interface{}) bool {
	aa, errA := json.Marshal(a)
	bb, errB := json.Marshal(b)
	return errA == nil && errB == nil && string(aa) == string(bb)
}
