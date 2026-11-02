package adb

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// This file holds the product-local copies of symbols the migrated adb source
// referenced from the root package alicloud. Product subpackages cannot import
// the root package, the rds-owned DBInstanceStatusCatcher, or the cross-product
// VpcService/KmsService, so faithful 1:1 copies are kept here. Behavior is
// preserved; see the migration report for the per-symbol provenance.

// Status mirrors the root alicloud.Status type so the migrated AdbService wait
// methods keep their original `status Status` signatures.
type Status string

const (
	Pending     = Status("Pending")
	Creating    = Status("Creating")
	Running     = Status("Running")
	Available   = Status("Available")
	Unavailable = Status("Unavailable")
	Modifying   = Status("Modifying")
	Deleting    = Status("Deleting")
	Starting    = Status("Starting")
	Stopping    = Status("Stopping")
	Deleted     = Status("Deleted")
)

type PayType string

const (
	PrePaid  = PayType("PrePaid")
	PostPaid = PayType("PostPaid")
	Prepaid  = PayType("Prepaid")
	Postpaid = PayType("Postpaid")
)

type TimeType string

const (
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

type RenewalStatus string

const (
	RenewAutoRenewal = RenewalStatus("AutoRenewal")
	RenewNormal      = RenewalStatus("Normal")
	RenewNotRenewal  = RenewalStatus("NotRenewal")
)

const (
	LOCAL_HOST_IP        = "127.0.0.1"
	DefaultTimeout       = 120
	DefaultTimeoutMedium = 500
)

// dbInstanceStatusCatcher mirrors the rds-owned root DBInstanceStatusCatcher used
// by AdbService.DescribeAdbAccount; the original is not importable from a product
// subpackage.
var dbInstanceStatusCatcher = helper.Catcher{Reason: "OperationDenied.DBInstanceStatus", RetryCount: 60, RetryWaitSeconds: 5}

// kmsDecrypt replicates the root KmsService.Decrypt behavior for the migrated adb
// account resource, calling the KMS Decrypt RPC directly via the shared client.
func kmsDecrypt(client *connectivity.AliyunClient, ciphertextBlob string, encryptionContext map[string]interface{}) (plaintext string, err error) {
	context, err := json.Marshal(encryptionContext)
	if err != nil {
		return plaintext, errs.WrapError(err)
	}

	var response map[string]interface{}
	action := "Decrypt"
	request := map[string]interface{}{
		"RegionId":          client.RegionId,
		"CiphertextBlob":    ciphertextBlob,
		"EncryptionContext": string(context[:]),
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		response, err = client.RpcPost("Kms", "2016-01-20", action, nil, request, false)
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
		return plaintext, errs.WrapErrorf(err, errs.DefaultErrorMsg, context, action, errs.AlibabaCloudSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Plaintext", response)
	if err != nil {
		return plaintext, errs.WrapErrorf(err, errs.FailedGetAttributeMsg, context, "$.Plaintext", response)
	}

	return fmt.Sprint(v), err
}

// describeVSwitch replicates the root VpcService.DescribeVSwitchWithTeadsl
// behavior for the migrated adb db_cluster resource's create path, calling the
// VPC DescribeVSwitchAttributes RPC directly via the shared client.
func describeVSwitch(client *connectivity.AliyunClient, id string) (object map[string]interface{}, err error) {
	action := "DescribeVSwitchAttributes"
	request := map[string]interface{}{
		"RegionId":  client.RegionId,
		"VSwitchId": id,
	}

	response, err := client.RpcPost("Vpc", "2016-04-28", action, nil, request, true)
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"InvalidVswitchID.NotFound"}) {
			return nil, errs.WrapErrorf(err, errs.NotFoundMsg, errs.AlibabaCloudSdkGoERROR)
		}
		return nil, errs.WrapErrorf(err, errs.DefaultErrorMsg, id, action, errs.AlibabaCloudSdkGoERROR)
	}
	logs.AddDebug(action, response, request)
	if v, ok := response["VSwitchId"].(string); ok && v != id {
		return nil, errs.WrapErrorf(errs.NotFoundErr("vswitch", id), errs.NotFoundMsg, errs.ProviderERROR)
	}
	return response, nil
}

// getChargeType is the product-local copy of the root alicloud.getChargeType
// helper (defined in service_alicloud_elasticsearch.go) that the migrated adb
// data_source_db_clusters reuses for the "charge_type" attribute. Product
// subpackages cannot reach the root package, so a faithful 1:1 copy is kept
// here. Behavior is preserved.
func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == strings.ToLower(string(PostPaid)) {
		return string(PostPaid)
	} else {
		return string(PrePaid)
	}
}

// BACKUP_TIME is the product-local copy of the root alicloud.BACKUP_TIME var
// (defined in extension_rds.go, rds-owned) that the migrated adb
// resource_backup_policy reuses for its preferred_backup_period validator.
// Product subpackages cannot reach the root package, so a faithful 1:1 copy is
// kept here. Behavior is preserved.
var BACKUP_TIME = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}

// COLON_SEPARATED is the product-local copy of the root common.go const
// COLON_SEPARATED = ":" that the migrated adb resource_connection reuses for
// ID composition/splitting. Product subpackages cannot reach the root package,
// so a faithful 1:1 copy is kept here. Behavior is preserved.
const COLON_SEPARATED = ":"
