package kms

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const secretEphemeralTypeName = "alicloud_kms_secret"

var (
	_ ephemeral.EphemeralResource              = &secretEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &secretEphemeral{}
)

func NewSecretEphemeralResource() ephemeral.EphemeralResource {
	return &secretEphemeral{}
}

type secretEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type secretEphemeralModel struct {
	SecretName     types.String   `tfsdk:"secret_name"`
	VersionId      types.String   `tfsdk:"version_id"`
	VersionStage   types.String   `tfsdk:"version_stage"`
	SecretData     types.String   `tfsdk:"secret_data"`
	VersionStages  []types.String `tfsdk:"version_stages"`
	SecretDataType types.String   `tfsdk:"secret_data_type"`
	SecretType     types.String   `tfsdk:"secret_type"`
}

func (e *secretEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = secretEphemeralTypeName
}

func (e *secretEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"secret_name": schema.StringAttribute{
				Required: true,
			},
			"version_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"version_stage": schema.StringAttribute{
				Optional: true,
			},
			"secret_data": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"version_stages": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"secret_data_type": schema.StringAttribute{
				Computed: true,
			},
			"secret_type": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *secretEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "GetSecretValue"

	var config secretEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"SecretName": config.SecretName.ValueString(),
	}
	if !config.VersionId.IsNull() {
		request["VersionId"] = config.VersionId.ValueString()
	}
	if !config.VersionStage.IsNull() {
		request["VersionStage"] = config.VersionStage.ValueString()
	}

	response, err := e.Client().RpcPost("Kms", "2016-01-20", action, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", secretEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", secretEphemeralTypeName, action).Error(),
		)
		return
	}

	state := config

	secretData, ok := response["SecretData"].(string)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", secretEphemeralTypeName, action),
			fmt.Sprintf("The response carried no SecretData: %v", response),
		)
		return
	}
	state.SecretData = types.StringValue(secretData)

	if v, ok := response["VersionId"].(string); ok {
		state.VersionId = types.StringValue(v)
	}

	if v, ok := response["SecretDataType"].(string); ok {
		state.SecretDataType = types.StringValue(v)
	}

	if v, ok := response["SecretType"].(string); ok {
		state.SecretType = types.StringValue(v)
	}

	if versionStages, ok := response["VersionStages"].(map[string]interface{}); ok {
		if stages, ok := versionStages["VersionStage"].([]interface{}); ok {
			state.VersionStages = make([]types.String, 0, len(stages))
			for _, stage := range stages {
				state.VersionStages = append(state.VersionStages, types.StringValue(fmt.Sprintf("%v", stage)))
			}
		}
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
}
