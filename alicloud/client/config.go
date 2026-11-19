package client

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
)

// Config 是客户端共享配置的类型别名，实现在 client/config 包中。
type Config = config.Config

type AssumeRoleWithOidc = config.AssumeRoleWithOidc

const DefaultClientRetryCountSmall = config.DefaultClientRetryCountSmall
const DefaultClientRetryCountMedium = config.DefaultClientRetryCountMedium
const DefaultClientRetryCountLarge = config.DefaultClientRetryCountLarge

const Terraform = config.Terraform
const Provider = config.Provider
const Module = config.Module
const TerraformTraceId = config.TerraformTraceId

func needRetry(err error) bool {
	return config.NeedRetry(err)
}

func isExpectedErrors(err error, expectCodes []string) bool {
	return config.IsExpectedErrors(err, expectCodes)
}
