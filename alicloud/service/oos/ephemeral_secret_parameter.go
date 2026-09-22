package oos

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const secretParameterEphemeralTypeName = "alicloud_oos_secret_parameter"

var (
	_ ephemeral.EphemeralResource              = &secretParameterEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &secretParameterEphemeral{}
)

func NewSecretParameterEphemeralResource() ephemeral.EphemeralResource {
	return &secretParameterEphemeral{}
}

type secretParameterEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type secretParameterEphemeralModel struct {
	SecretParameterName types.String `tfsdk:"secret_parameter_name"`
	ParameterVersion    types.Int64  `tfsdk:"parameter_version"`
	Value               types.String `tfsdk:"value"`
}

func (e *secretParameterEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = secretParameterEphemeralTypeName
}

func (e *secretParameterEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"secret_parameter_name": schema.StringAttribute{
				Required: true,
			},
			"parameter_version": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"value": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (e *secretParameterEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "GetSecretParameter"

	var config secretParameterEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"Name":           config.SecretParameterName.ValueString(),
		"WithDecryption": true,
	}
	if !config.ParameterVersion.IsNull() {
		request["ParameterVersion"] = config.ParameterVersion.ValueInt64()
	}

	response, err := e.Client().RpcPost("oos", "2019-06-01", action, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", secretParameterEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", secretParameterEphemeralTypeName, action).Error(),
		)
		return
	}

	parameter, ok := response["Parameter"].(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", secretParameterEphemeralTypeName, action),
			fmt.Sprintf("The response carried no Parameter: %v", response),
		)
		return
	}

	state := config

	value, ok := parameter["Value"].(string)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", secretParameterEphemeralTypeName, action),
			fmt.Sprintf("The response carried no Value: %v", response),
		)
		return
	}
	state.Value = types.StringValue(value)

	if v, ok := parameter["ParameterVersion"].(json.Number); ok {
		if n, err := v.Int64(); err == nil {
			state.ParameterVersion = types.Int64Value(n)
		}
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
}
