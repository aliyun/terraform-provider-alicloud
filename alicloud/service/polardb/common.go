package polardb

import (
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
)

// Root-package leaf aliases re-homed into the PolarDB product package so the
// migrated sources keep their bare names without reverse-importing the root
// package. Each entry mirrors the root alias verbatim (same kind, same target);
// the migration does not alter behavior.
type (
	Tag             = tags.Tag
	TagResourceType = tags.TagResourceType
	Protocol        = helper.Protocol
)

// Function aliases (var-form, identical to the root package).
var (
	dataResourceIdHash = helper.DataResourceIdHash
	String             = helper.String
	Trim               = helper.Trim
	mapMerge           = helper.MapMerge
	GetNotFoundMessage = errs.GetNotFoundMessage
)

// String-constant aliases re-exported from errs (same values as the root package).
const (
	Throttling                = errs.Throttling
	DeleteTimeoutMsg          = errs.DeleteTimeoutMsg
	WaitTimeoutMsg            = errs.WaitTimeoutMsg
	IdMsg                     = errs.IdMsg
	FailedToReachTargetStatus = errs.FailedToReachTargetStatus
	NotFoundMsg               = errs.NotFoundMsg
	RequiredWhenMsg           = errs.RequiredWhenMsg
	InvalidAttributeValue     = errs.InvalidAttributeValue
)

// Status mirrors the root-package status enum so migrated PolarDB helpers keep
// their typed lifecycle comparisons without reverse-importing the root package.
type Status string

const (
	Pending         = Status("Pending")
	Creating        = Status("Creating")
	Running         = Status("Running")
	Available       = Status("Available")
	Unavailable     = Status("Unavailable")
	Modifying       = Status("Modifying")
	Deleting        = Status("Deleting")
	Starting        = Status("Starting")
	Stopping        = Status("Stopping")
	Stopped         = Status("Stopped")
	Normal          = Status("Normal")
	Changing        = Status("Changing")
	Online          = Status("online")
	Configuring     = Status("configuring")
	Associating     = Status("Associating")
	Unassociating   = Status("Unassociating")
	InUse           = Status("InUse")
	DiskInUse       = Status("In_use")
	Active          = Status("Active")
	Inactive        = Status("Inactive")
	Idle            = Status("Idle")
	SoldOut         = Status("SoldOut")
	InService       = Status("InService")
	Removing        = Status("Removing")
	DisabledStatus  = Status("Disabled")
	Init            = Status("Init")
	Provisioning    = Status("Provisioning")
	Updating        = Status("Updating")
	FinancialLocked = Status("FinancialLocked")
	PUBLISHED       = Status("Published")
	NOPUBLISHED     = Status("NonPublished")
	Deleted         = Status("Deleted")
	Null            = Status("Null")
	Enable          = Status("Enable")
	BINDED          = Status("BINDED")
)

// timeout (minutes), preserved from the root package for PolarDB long waits.
const DefaultTimeout = 120
const Timeout5Minute = 300
const DefaultTimeoutMedium = 500
const DefaultLongTimeout = 1000

const DEFAULT_REGION = "cn-beijing"
const INT_MAX = 2147483647

// separators & literals preserved from the root package.
const (
	COLON_SEPARATED = ":"
	SLASH_SEPARATED = "/"
	LOCAL_HOST_IP   = "127.0.0.1"
	MULTI_IZ_SYMBOL = helper.MultiIzSymbol
	COMMA_SEPARATED = helper.CommaSeparated
)

// DBInstanceStatusCatcher retries OperationDenied.DBInstanceStatus; replicated
// here (same value as the root RDS/ADB service) so the product package does not
// reverse-import the root package. shared with unmigrated root consumers.
var DBInstanceStatusCatcher = helper.Catcher{Reason: "OperationDenied.DBInstanceStatus", RetryCount: 60, RetryWaitSeconds: 5}

// Product enums re-homed from the root package (extension_ecs.go / common.go /
// extension_rds.go) so migrated PolarDB resources keep their typed payment,
// renewal and network comparisons without reverse-importing the root package.
type PayType string

const (
	PrePaid    = PayType("PrePaid")
	PostPaid   = PayType("PostPaid")
	Prepaid    = PayType("Prepaid")
	Postpaid   = PayType("Postpaid")
	Serverless = PayType("Serverless")
)

type RenewalStatus string

const (
	RenewAutoRenewal = RenewalStatus("AutoRenewal")
	RenewNormal      = RenewalStatus("Normal")
	RenewNotRenewal  = RenewalStatus("NotRenewal")
)

type NetworkType string

const (
	Vpc = NetworkType("Vpc")
)

// TimeType mirrors the root-package time-unit enum (common.go) referenced by
// the PolarDB account/cluster renewal-period comparisons.
type TimeType string

const (
	Month = TimeType("Month")
	Year  = TimeType("Year")
)

// BACKUP_TIME is the catalog of hourly backup windows (extension_rds.go),
// referenced by the PolarDB backup policy resource.
var BACKUP_TIME = []string{
	"00:00Z-01:00Z", "01:00Z-02:00Z", "02:00Z-03:00Z", "03:00Z-04:00Z", "04:00Z-05:00Z",
	"05:00Z-06:00Z", "06:00Z-07:00Z", "07:00Z-08:00Z", "08:00Z-09:00Z", "09:00Z-10:00Z",
	"10:00Z-11:00Z", "11:00Z-12:00Z", "12:00Z-13:00Z", "13:00Z-14:00Z", "14:00Z-15:00Z",
	"15:00Z-16:00Z", "16:00Z-17:00Z", "17:00Z-18:00Z", "18:00Z-19:00Z", "19:00Z-20:00Z",
	"20:00Z-21:00Z", "21:00Z-22:00Z", "22:00Z-23:00Z", "23:00Z-24:00Z",
}

// Function aliases (var-form) re-exporting leaf-package helpers, identical to
// the root package.
var (
	incrementalWait = helper.IncrementalWait
	NewInvoker      = helper.NewInvoker
	BuildStateConf  = helper.BuildStateConf
	addDebug        = logs.AddDebug
	NotFoundErr     = errs.NotFoundErr
	NeedRetry       = helper.NeedRetry
	NotFoundError   = errs.NotFoundError
)

// String-constant aliases re-exporting errs string constants, identical to the
// root package. WrapErrorf and IsExpectedErrors are functions and MUST stay
// var-form (a function value cannot be assigned to a const).
const (
	DefaultErrorMsg        = errs.DefaultErrorMsg
	DataDefaultErrorMsg    = errs.DataDefaultErrorMsg
	AlibabaCloudSdkGoERROR = errs.AlibabaCloudSdkGoERROR
	ProviderERROR          = errs.ProviderERROR
	FailedGetAttributeMsg  = errs.FailedGetAttributeMsg
	NotFoundWithResponse   = errs.NotFoundWithResponse
	NotFoundWithError      = errs.NotFoundWithError
)

// Function aliases (var-form) re-exporting errs helpers, identical to the root
// package. Functions cannot live in a const block.
var (
	WrapError        = errs.WrapError
	WrapErrorf       = errs.WrapErrorf
	IsExpectedErrors = errs.IsExpectedErrors
)

// getChargeType normalizes a raw payment type to the canonical PrePaid/PostPaid
// display value. Replicated from the root package (service_alicloud_elasticsearch.go)
// with string literals — behavior-equivalent, no PayType type cascade — so the
// product package does not reverse-import the root package.
func getChargeType(paymentType string) string {
	if strings.ToLower(paymentType) == "postpaid" {
		return "PostPaid"
	}
	return "PrePaid"
}
