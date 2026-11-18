package connectivity

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/client/config"
)

// Config 保留类型别名，兼容已有 provider 和资源代码；实现位于 client/config。
type Config = config.Config

type AssumeRoleWithOidc = config.AssumeRoleWithOidc

const DefaultClientRetryCountSmall = config.DefaultClientRetryCountSmall
const DefaultClientRetryCountMedium = config.DefaultClientRetryCountMedium
const DefaultClientRetryCountLarge = config.DefaultClientRetryCountLarge

const Terraform = config.Terraform
const Provider = config.Provider
const Module = config.Module
const TerraformTraceId = config.TerraformTraceId
