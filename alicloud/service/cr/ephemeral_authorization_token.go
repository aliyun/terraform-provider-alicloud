package cr

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const authorizationTokenEphemeralTypeName = "alicloud_cr_ee_authorization_token"

var (
	_ ephemeral.EphemeralResource              = &authorizationTokenEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &authorizationTokenEphemeral{}
)

func NewAuthorizationTokenEphemeralResource() ephemeral.EphemeralResource {
	return &authorizationTokenEphemeral{}
}

type authorizationTokenEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type authorizationTokenEphemeralModel struct {
	InstanceId         types.String `tfsdk:"instance_id"`
	ExpiresInHours     types.Int64  `tfsdk:"expires_in_hours"`
	TempUsername       types.String `tfsdk:"temp_username"`
	AuthorizationToken types.String `tfsdk:"authorization_token"`
	ExpireTime         types.Int64  `tfsdk:"expire_time"`
}

func (e *authorizationTokenEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = authorizationTokenEphemeralTypeName
}

func (e *authorizationTokenEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Required: true,
			},
			"expires_in_hours": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(1, 24),
				},
			},
			"temp_username": schema.StringAttribute{
				Computed: true,
			},
			"authorization_token": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"expire_time": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (e *authorizationTokenEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "GetAuthorizationToken"

	var config authorizationTokenEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"InstanceId": config.InstanceId.ValueString(),
	}
	if !config.ExpiresInHours.IsNull() {
		request["ExpiresInHours"] = config.ExpiresInHours.ValueInt64()
	}

	response, err := e.Client().RpcPost("cr", "2018-12-01", action, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", authorizationTokenEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", authorizationTokenEphemeralTypeName, action).Error(),
		)
		return
	}

	state := config

	if success, ok := response["IsSuccess"].(bool); ok && !success {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", authorizationTokenEphemeralTypeName, action),
			fmt.Sprintf("The API reported failure: %v", response),
		)
		return
	}

	token, ok := response["AuthorizationToken"].(string)
	if !ok || token == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", authorizationTokenEphemeralTypeName, action),
			fmt.Sprintf("The response carried no AuthorizationToken: %v", response),
		)
		return
	}
	state.AuthorizationToken = types.StringValue(token)

	if v, ok := response["TempUsername"].(string); ok {
		state.TempUsername = types.StringValue(v)
	}

	// RpcPost parses the response body with json.Decoder.UseNumber(), so every
	// numeric field arrives as json.Number rather than float64.
	if v, ok := response["ExpireTime"].(json.Number); ok {
		if ms, err := v.Int64(); err == nil {
			state.ExpireTime = types.Int64Value(ms)
		}
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
}
