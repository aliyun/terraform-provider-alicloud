package conns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
)

type SDKResource struct {
	TypeName string // e.g. "alicloud_vpc"
	// Called once, at build time, and the result reused for every operation.
	// Contrast Resource.Factory, called once per RPC.
	Factory      func() *schema.Resource
	Interceptors []intercept.Interceptor
}

type SDKDataSource struct {
	TypeName     string
	Factory      func() *schema.Resource
	Interceptors []intercept.Interceptor
}
