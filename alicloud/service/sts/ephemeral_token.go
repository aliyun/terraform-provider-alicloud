package sts

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const stsTokenEphemeralTypeName = "alicloud_sts_token"

var (
	_ ephemeral.EphemeralResource              = &stsTokenEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &stsTokenEphemeral{}
)

func NewStsTokenEphemeralResource() ephemeral.EphemeralResource {
	return &stsTokenEphemeral{}
}

type stsTokenEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type stsTokenEphemeralModel struct {
	RoleArn         types.String `tfsdk:"role_arn"`
	SessionName     types.String `tfsdk:"session_name"`
	DurationSeconds types.Int64  `tfsdk:"duration_seconds"`
	Policy          types.String `tfsdk:"policy"`
	ExternalId      types.String `tfsdk:"external_id"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	AccessKeySecret types.String `tfsdk:"access_key_secret"`
	SecurityToken   types.String `tfsdk:"security_token"`
	Expiration      types.String `tfsdk:"expiration"`
	AssumedRoleArn  types.String `tfsdk:"assumed_role_arn"`
}

func (e *stsTokenEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = stsTokenEphemeralTypeName
}

func (e *stsTokenEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"role_arn": schema.StringAttribute{
				Required: true,
			},
			"session_name": schema.StringAttribute{
				Required: true,
			},
			"duration_seconds": schema.Int64Attribute{
				Optional: true,
				Validators: []validator.Int64{
					int64validator.Between(900, 43200),
				},
			},
			"policy": schema.StringAttribute{
				Optional: true,
			},
			"external_id": schema.StringAttribute{
				Optional: true,
			},
			"access_key_id": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"access_key_secret": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"security_token": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"expiration": schema.StringAttribute{
				Computed: true,
			},
			"assumed_role_arn": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *stsTokenEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "AssumeRole"

	var config stsTokenEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"RoleArn":         config.RoleArn.ValueString(),
		"RoleSessionName": config.SessionName.ValueString(),
	}
	if !config.DurationSeconds.IsNull() {
		request["DurationSeconds"] = config.DurationSeconds.ValueInt64()
	}
	if !config.Policy.IsNull() {
		request["Policy"] = config.Policy.ValueString()
	}
	if !config.ExternalId.IsNull() {
		request["ExternalId"] = config.ExternalId.ValueString()
	}

	response, err := e.Meta().RpcPost("Sts", "2015-04-01", action, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", stsTokenEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", stsTokenEphemeralTypeName, action).Error(),
		)
		return
	}

	state := config

	credentials, ok := response["Credentials"].(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", stsTokenEphemeralTypeName, action),
			fmt.Sprintf("The response carried no Credentials: %v", response),
		)
		return
	}
	accessKeyId, _ := credentials["AccessKeyId"].(string)
	accessKeySecret, _ := credentials["AccessKeySecret"].(string)
	securityToken, _ := credentials["SecurityToken"].(string)
	if accessKeyId == "" || accessKeySecret == "" || securityToken == "" {
		// Never echo the response here: it may carry a partially issued credential.
		missing := ""
		if accessKeyId == "" {
			missing += " AccessKeyId"
		}
		if accessKeySecret == "" {
			missing += " AccessKeySecret"
		}
		if securityToken == "" {
			missing += " SecurityToken"
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", stsTokenEphemeralTypeName, action),
			fmt.Sprintf("The response carried an incomplete credential, missing:%s", missing),
		)
		return
	}
	state.AccessKeyId = types.StringValue(accessKeyId)
	state.AccessKeySecret = types.StringValue(accessKeySecret)
	state.SecurityToken = types.StringValue(securityToken)

	if v, ok := credentials["Expiration"].(string); ok {
		state.Expiration = types.StringValue(v)
	}

	if assumedRoleUser, ok := response["AssumedRoleUser"].(map[string]interface{}); ok {
		if v, ok := assumedRoleUser["Arn"].(string); ok {
			state.AssumedRoleArn = types.StringValue(v)
		}
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
}
