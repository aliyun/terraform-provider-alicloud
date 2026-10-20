package ram

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

const (
	ramAccessKeyEphemeralTypeName = "alicloud_ram_access_key"
	// Private state key under which Open stashes the data Close needs to revoke
	// the access key. CloseRequest carries Private but not Result, so the data
	// Open writes to Result for the Terraform configuration cannot be read back
	// in Close; it must round-trip through private state as JSON.
	ramAccessKeyEphemeralPrivateCleanupKey = "cleanup"
)

var (
	_ ephemeral.EphemeralResource              = &ramAccessKeyEphemeral{}
	_ ephemeral.EphemeralResourceWithConfigure = &ramAccessKeyEphemeral{}
)

// NewRamAccessKeyEphemeralResource returns an ephemeral resource that issues a
// RAM access key for the duration of its reference and revokes it when the
// reference closes. The secret never lands in the Terraform state file: callers
// consume access_key_id / access_key_secret from the ephemeral output within
// the same Terraform operation, for example by configuring an aliased provider
// with them. Once the reference ends, Close deletes the access key so no
// residual credential remains; credentials that must outlive the operation
// belong to the managed alicloud_ram_access_key resource instead.
func NewRamAccessKeyEphemeralResource() ephemeral.EphemeralResource {
	return &ramAccessKeyEphemeral{}
}

type ramAccessKeyEphemeral struct {
	fwadapt.EphemeralResourceBase
}

type ramAccessKeyEphemeralModel struct {
	UserName        types.String `tfsdk:"user_name"`
	AccessKeyId     types.String `tfsdk:"access_key_id"`
	AccessKeySecret types.String `tfsdk:"access_key_secret"`
	Status          types.String `tfsdk:"status"`
	CreateDate      types.String `tfsdk:"create_date"`
}

// ramAccessKeyCleanup carries the minimum data Close needs to revoke the access
// key Open created. Result is only visible to the Terraform configuration, not
// to Close, so this is stashed in ephemeral private state during Open.
type ramAccessKeyCleanup struct {
	AccessKeyId string `json:"access_key_id"`
	UserName    string `json:"user_name,omitempty"`
}

func (e *ramAccessKeyEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = ramAccessKeyEphemeralTypeName
}

func (e *ramAccessKeyEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user_name": schema.StringAttribute{
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
			"status": schema.StringAttribute{
				Computed: true,
			},
			"create_date": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (e *ramAccessKeyEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	const action = "CreateAccessKey"

	var config ramAccessKeyEphemeralModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := map[string]interface{}{}
	if !config.UserName.IsNull() {
		request["UserName"] = config.UserName.ValueString()
	}

	response, err := createRamAccessKey(e.Client(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", ramAccessKeyEphemeralTypeName, action),
			errs.WrapErrorf(err, "Opening %s: calling %s", ramAccessKeyEphemeralTypeName, action).Error()+
				". CreateAccessKey has no idempotency token and is not retried automatically: if the request reached RAM, an access key may have been created without being tracked for revocation. List the user's access keys with ListAccessKeys and delete any unexpected key before retrying.",
		)
		return
	}

	state := config

	accessKey, ok := response["AccessKey"].(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", ramAccessKeyEphemeralTypeName, action),
			fmt.Sprintf("The %s response carried no AccessKey", action),
		)
		return
	}
	accessKeyId, _ := accessKey["AccessKeyId"].(string)
	accessKeySecret, _ := accessKey["AccessKeySecret"].(string)
	if accessKeyId == "" || accessKeySecret == "" {
		// Never echo the response here: it may carry a partially issued credential.
		missing := ""
		if accessKeyId == "" {
			missing += " AccessKeyId"
		}
		if accessKeySecret == "" {
			missing += " AccessKeySecret"
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", ramAccessKeyEphemeralTypeName, action),
			fmt.Sprintf("The response carried an incomplete access key, missing:%s", missing),
		)
		return
	}
	state.AccessKeyId = types.StringValue(accessKeyId)
	state.AccessKeySecret = types.StringValue(accessKeySecret)

	if v, ok := accessKey["Status"].(string); ok && v != "" {
		state.Status = types.StringValue(v)
	}
	if v, ok := accessKey["CreateDate"].(string); ok && v != "" {
		state.CreateDate = types.StringValue(v)
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Stash the cleanup data for Close. Result is visible to the Terraform
	// configuration, not to Close; Close reads private state only.
	cleanup := ramAccessKeyCleanup{AccessKeyId: accessKeyId}
	if !config.UserName.IsNull() {
		cleanup.UserName = config.UserName.ValueString()
	}
	cleanupJSON, err := json.Marshal(cleanup)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Opening %s: calling %s", ramAccessKeyEphemeralTypeName, action),
			fmt.Sprintf("marshalling cleanup data: %s", err.Error()),
		)
		return
	}
	resp.Diagnostics.Append(resp.Private.SetKey(ctx, ramAccessKeyEphemeralPrivateCleanupKey, cleanupJSON)...)
}

func (e *ramAccessKeyEphemeral) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	const action = "DeleteAccessKey"

	raw, diags := req.Private.GetKey(ctx, ramAccessKeyEphemeralPrivateCleanupKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(raw) == 0 {
		// Open never stashed cleanup data (failed early, or nothing opened).
		// Nothing to revoke; the goal of Close — no residual credential — is met.
		return
	}

	var cleanup ramAccessKeyCleanup
	if err := json.Unmarshal(raw, &cleanup); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Closing %s: calling %s", ramAccessKeyEphemeralTypeName, action),
			fmt.Sprintf("unmarshalling cleanup data: %s", err.Error()),
		)
		return
	}

	if err := deleteRamAccessKey(e.Client(), cleanup); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Closing %s: calling %s", ramAccessKeyEphemeralTypeName, action),
			errs.WrapErrorf(err, "Closing %s: calling %s", ramAccessKeyEphemeralTypeName, action).Error(),
		)
	}
}

// ramAccessKeyClient is the subset of the Aliyun client the ephemeral resource
// depends on. The narrow interface keeps creation, revocation and retry
// behavior unit-testable without cloud calls.
type ramAccessKeyClient interface {
	RpcPost(apiProductCode string, apiVersion string, apiName string, query map[string]interface{}, body map[string]interface{}, autoRetry bool) (map[string]interface{}, error)
}

func createRamAccessKey(client ramAccessKeyClient, request map[string]interface{}) (map[string]interface{}, error) {
	// CreateAccessKey has no idempotency token, so a retried request can issue a
	// second access key that nothing revokes. Never let the runtime retry it;
	// the caller surfaces the ambiguous outcome to the operator instead.
	return client.RpcPost("Ram", "2015-05-01", "CreateAccessKey", nil, request, false)
}

func deleteRamAccessKey(client ramAccessKeyClient, cleanup ramAccessKeyCleanup) error {
	if cleanup.AccessKeyId == "" {
		return nil
	}

	request := map[string]interface{}{
		"UserAccessKeyId": cleanup.AccessKeyId,
	}
	if cleanup.UserName != "" {
		request["UserName"] = cleanup.UserName
	}

	if _, err := client.RpcPost("Ram", "2015-05-01", "DeleteAccessKey", nil, request, true); err != nil {
		// The access key may already be gone (deleted out-of-band, or by a prior
		// Close). Such cases are not failures: the goal of Close — no residual
		// credential — is met.
		if errs.NotFoundError(err) || errs.IsExpectedErrors(err, []string{"EntityNotExist.User.AccessKey", "EntityNotExist.User"}) {
			return nil
		}
		return err
	}
	return nil
}
