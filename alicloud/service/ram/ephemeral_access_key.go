package ram

import (
	"context"
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const accessKeyEphemeralTypeName = "alicloud_ram_access_key"

var (
	_ ephemeral.EphemeralResource              = &accessKeyEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &accessKeyEphemeral{}
)

func NewAccessKeyEphemeralResource() ephemeral.EphemeralResource {
	return &accessKeyEphemeral{}
}

type accessKeyEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type accessKeyEphemeralModel struct {
	UserName        types.String `tfsdk:"user_name"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	AccessKeySecret types.String `tfsdk:"access_key_secret"`
	Status          types.String `tfsdk:"status"`
	CreateTime      types.String `tfsdk:"create_time"`
}

func (e *accessKeyEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = accessKeyEphemeralTypeName
}

func (e *accessKeyEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_name": schema.StringAttribute{
				Optional: true,
			},
			"access_key_id": schema.StringAttribute{
				Computed: true,
			},
			"access_key_secret": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"create_time": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *accessKeyEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "CreateAccessKey"

	var config accessKeyEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}
	if !config.UserName.IsNull() && config.UserName.ValueString() != "" {
		request["UserName"] = config.UserName.ValueString()
	}

	response, err := e.Client().RpcPost("Ram", "2015-05-01", action, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", accessKeyEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", accessKeyEphemeralTypeName, action).Error(),
		)
		return
	}

	state := config

	accessKey, ok := response["AccessKey"].(map[string]interface{})
	if !ok {
		// Never echo the response here: it may carry a partially issued credential.
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", accessKeyEphemeralTypeName, action),
			"The response carried no AccessKey.",
		)
		return
	}
	accessKeyId, _ := accessKey["AccessKeyId"].(string)
	accessKeySecret, _ := accessKey["AccessKeySecret"].(string)
	if accessKeyId == "" || accessKeySecret == "" {
		missing := ""
		if accessKeyId == "" {
			missing += " AccessKeyId"
		}
		if accessKeySecret == "" {
			missing += " AccessKeySecret"
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", accessKeyEphemeralTypeName, action),
			fmt.Sprintf("The response carried an incomplete credential, missing:%s", missing),
		)
		return
	}

	state.AccessKeyId = types.StringValue(accessKeyId)
	state.AccessKeySecret = types.StringValue(accessKeySecret)

	if v, ok := accessKey["Status"].(string); ok {
		state.Status = types.StringValue(v)
	}
	if v, ok := accessKey["CreateDate"].(string); ok {
		state.CreateTime = types.StringValue(v)
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
}
