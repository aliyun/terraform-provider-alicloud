package conns

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
)

type Resource struct {
	TypeName string
	// Called once per RPC; contrast SDKResource.Factory, called once.
	Factory      func() resource.Resource
	Interceptors []intercept.Interceptor
}

type DataSource struct {
	TypeName     string
	Factory      func() datasource.DataSource
	Interceptors []intercept.Interceptor
}

type Function struct {
	TypeName string // short name: "arn_build" → provider::alicloud::arn_build
	Factory  func() function.Function
}

// No Interceptors field, here or on Function, ListResource and Action: these have
// no wrapper yet, so the field would accept interceptors that never run.
type EphemeralResource struct {
	TypeName string
	Factory  func() ephemeral.EphemeralResource
}

type ListResource struct {
	TypeName string // matches the managed resource type name it lists
	Factory  func() list.ListResource
}

type Action struct {
	TypeName string
	Factory  func() action.Action
}
