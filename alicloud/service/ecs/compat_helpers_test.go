package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	ecs "github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"
)

// Root-layout test helpers that migrated acceptance helpers still call by
// bare name. Under the split package they are aliases in package ecs
// (aliases.go / errors.go): the exported ones (NotFoundError, WrapError,
// WrapErrorf, NeedRetry, PageSizeLarge, DefaultErrorMsg, IsExpectedErrors,
// AlibabaCloudSdkGoERROR, FailedGetAttributeMsg) are reached via ecs; the
// unexported ones (addDebug = logs.AddDebug, incrementalWait =
// helper.IncrementalWait) cannot be reached cross-package, so they alias the
// leaf logs/helper packages directly, exactly as package ecs does. These
// restore bare-name visibility for ecs_test so the migrated sources
// (disk_default_kms_key, image_copy, security_group_rule) read verbatim,
// mirroring compat_enums_test.go. Values are the identical ecs declarations,
// not redefinitions.
var (
	NotFoundError          = ecs.NotFoundError
	WrapError              = ecs.WrapError
	WrapErrorf             = ecs.WrapErrorf
	NeedRetry              = ecs.NeedRetry
	PageSizeLarge          = ecs.PageSizeLarge
	DefaultErrorMsg        = ecs.DefaultErrorMsg
	IsExpectedErrors       = ecs.IsExpectedErrors
	AlibabaCloudSdkGoERROR = ecs.AlibabaCloudSdkGoERROR
	FailedGetAttributeMsg  = ecs.FailedGetAttributeMsg

	addDebug        = logs.AddDebug
	incrementalWait = helper.IncrementalWait

	// formatInt (aliases.go = convert.FormatInt) and buildClientToken
	// (aliases.go = helper.BuildClientToken) are unexported package-ecs aliases
	// used bare by migrated ecs_test helpers (selectInstancePrivatePool,
	// testAccPrepareSecurityGroupRuleCrossAccountPeer). They alias the leaf
	// convert / helper packages directly, exactly as addDebug /
	// incrementalWait do, restoring bare-name visibility for ecs_test.
	formatInt        = convert.FormatInt
	buildClientToken = helper.BuildClientToken
)
