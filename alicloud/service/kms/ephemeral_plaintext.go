package kms

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

const plaintextEphemeralTypeName = "alicloud_kms_plaintext"

var (
	_ ephemeral.EphemeralResource              = &plaintextEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &plaintextEphemeral{}
)

func NewPlaintextEphemeralResource() ephemeral.EphemeralResource {
	return &plaintextEphemeral{}
}

type plaintextEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type plaintextEphemeralModel struct {
	CiphertextBlob    types.String `tfsdk:"ciphertext_blob"`
	EncryptionContext types.Map    `tfsdk:"encryption_context"`
	Plaintext         types.String `tfsdk:"plaintext"`
	KeyId             types.String `tfsdk:"key_id"`
}

func (e *plaintextEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = plaintextEphemeralTypeName
}

func (e *plaintextEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ciphertext_blob": schema.StringAttribute{
				Required: true,
			},
			"encryption_context": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"plaintext": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"key_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *plaintextEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "Decrypt"

	var config plaintextEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{
		"CiphertextBlob": config.CiphertextBlob.ValueString(),
	}
	if !config.EncryptionContext.IsNull() {
		var encryptionContext map[string]string
		resp.Diagnostics.Append(config.EncryptionContext.ElementsAs(ctx, &encryptionContext, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		contextJson, err := json.Marshal(encryptionContext)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Opening %s: calling %s", plaintextEphemeralTypeName, action),
				fmt.Sprintf("Marshaling the encryption context: %s", err),
			)
			return
		}
		request["EncryptionContext"] = string(contextJson)
	}

	response, err := e.Client().RpcPost("Kms", "2016-01-20", action, nil, request, true)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", plaintextEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", plaintextEphemeralTypeName, action).Error(),
		)
		return
	}

	state := config

	plaintext, ok := response["Plaintext"].(string)
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", plaintextEphemeralTypeName, action),
			fmt.Sprintf("The response carried no Plaintext: %v", response),
		)
		return
	}
	state.Plaintext = types.StringValue(plaintext)

	if v, ok := response["KeyId"].(string); ok {
		state.KeyId = types.StringValue(v)
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
}
