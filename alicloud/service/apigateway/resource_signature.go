// Package apigateway implements the Alibaba Cloud API Gateway (CloudAPI) resources
// and data sources on top of terraform-plugin-framework. The CRUD calls reuse the
// alibaba-cloud-sdk-go cloudapi client through connectivity.AliyunClient, matching
// the request/response shapes the service exposes; the resource/datasource plumbing
// is Framework-native and registered in the framework provider server.
package apigateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const signatureTypeName = "alicloud_api_gateway_signature"

var (
	_ resource.Resource                = &signatureResource{}
	_ resource.ResourceWithImportState = &signatureResource{}
)

func NewSignatureResource() resource.Resource {
	return &signatureResource{}
}

type signatureResource struct {
	fwadapt.ResourceBase
}

type signatureResourceModel struct {
	Id              types.String `tfsdk:"id"`
	SignatureName   types.String `tfsdk:"signature_name"`
	SignatureKey    types.String `tfsdk:"signature_key"`
	SignatureSecret types.String `tfsdk:"signature_secret"`
	CreatedTime     types.String `tfsdk:"created_time"`
	ModifiedTime    types.String `tfsdk:"modified_time"`
}

func (r *signatureResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = signatureTypeName
}

func (r *signatureResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides an API Gateway Signature resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of the signature allocated by the system.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"signature_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the signature. It must be 4 to 15 characters in length and can contain Chinese characters, English letters, numbers, and underscores (_). It must start with a letter or underscore.",
			},
			"signature_key": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The key of the signature. It is used together with the signature secret to verify the identity of the API caller.",
			},
			"signature_secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The secret of the signature. It is used together with the signature key to verify the identity of the API caller.",
			},
			"created_time": schema.StringAttribute{
				Computed:    true,
				Description: "The creation time of the signature.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"modified_time": schema.StringAttribute{
				Computed:    true,
				Description: "The last modification time of the signature.",
				// No UseStateForUnknown: ModifySignature bumps the timestamp, so the
				// planned state value can never survive an apply.
			},
		},
	}
}

func (r *signatureResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan signatureResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	request := cloudapi.CreateCreateSignatureRequest()
	request.RegionId = client.RegionId
	request.SignatureName = plan.SignatureName.ValueString()
	request.SignatureKey = plan.SignatureKey.ValueString()
	request.SignatureSecret = plan.SignatureSecret.ValueString()

	// RepeatedCommit surfaces when the service is still settling a prior create; retry
	// for the same bounded window the SDKv2 resource did, instead of failing the apply.
	var signatureId string
	if err := retry(5*time.Minute, func() (bool, error) {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.CreateSignature(request)
		})
		if err != nil {
			if isExpectedError(err, "RepeatedCommit") {
				return true, err
			}
			return false, err
		}
		response, _ := raw.(*cloudapi.CreateSignatureResponse)
		signatureId = response.SignatureId
		return false, nil
	}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSignature", signatureTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, signatureTypeName, request.GetActionName(), errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if signatureId == "" {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSignature", signatureTypeName),
			fmt.Sprintf("The CreateSignature response carried no SignatureId: %v", request),
		)
		return
	}

	object, err := describeSignature(client, signatureId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeSignatures", signatureTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, signatureId, "DescribeSignatures", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling CreateSignature", signatureTypeName),
			fmt.Sprintf("The signature %s is not visible via DescribeSignatures right after creation", signatureId),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, signatureState(object))...)
}

func (r *signatureResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state signatureResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := describeSignature(r.Client(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeSignatures", signatureTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeSignatures", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, signatureState(object))...)
}

func (r *signatureResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state signatureResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	request := cloudapi.CreateModifySignatureRequest()
	request.RegionId = client.RegionId
	request.SignatureId = state.Id.ValueString()
	request.SignatureName = plan.SignatureName.ValueString()
	request.SignatureKey = plan.SignatureKey.ValueString()
	request.SignatureSecret = plan.SignatureSecret.ValueString()

	if _, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.ModifySignature(request)
	}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling ModifySignature", signatureTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}

	// ModifySignature is eventually consistent for the read-back; let the describe retry
	// bridge the window rather than trusting a fresh Read to see the new values.
	object, err := describeSignature(client, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling DescribeSignatures", signatureTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), "DescribeSignatures", errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
	if object == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, signatureState(object))...)
}

func (r *signatureResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state signatureResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.Client()
	request := cloudapi.CreateDeleteSignatureRequest()
	request.RegionId = client.RegionId
	request.SignatureId = state.Id.ValueString()

	if err := retry(5*time.Minute, func() (bool, error) {
		_, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DeleteSignature(request)
		})
		if err != nil {
			// NotFoundSignature means the signature is already gone; treat as success.
			if isExpectedError(err, "NotFoundSignature") {
				return false, nil
			}
			return true, err
		}
		return false, nil
	}); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling DeleteSignature", signatureTypeName),
			errs.WrapErrorf(err, errs.DefaultErrorMsg, state.Id.ValueString(), request.GetActionName(), errs.AlibabaCloudSdkGoERROR).Error(),
		)
		return
	}
}

func (r *signatureResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// describeSignature queries DescribeSignatures by signature ID and returns the matching
// SignatureInfo, or nil when the signature no longer exists (the service returns an
// empty list rather than a not-found error).
func describeSignature(client *connectivity.AliyunClient, id string) (*cloudapi.SignatureInfo, error) {
	request := cloudapi.CreateDescribeSignaturesRequest()
	request.RegionId = client.RegionId
	request.SignatureId = id
	request.PageSize = requests.NewInteger(pageSizeSmall)
	request.PageNumber = requests.NewInteger(1)

	// The signature may not be visible to DescribeSignatures immediately after a create
	// or modify; retry the read-back for the same bounded window the SDKv2 resource did.
	// An empty list is a transient miss during read-back, and a real absence during a
	// standalone Read; the caller turns an exhausted not-found into a nil object. A real
	// API error is non-retryable and surfaces immediately.
	var object *cloudapi.SignatureInfo
	if err := retry(5*time.Minute, func() (bool, error) {
		raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
			return cloudApiClient.DescribeSignatures(request)
		})
		if err != nil {
			return false, err
		}
		response, _ := raw.(*cloudapi.DescribeSignaturesResponse)
		if len(response.SignatureInfos.SignatureInfo) == 0 {
			return true, errSignatureNotFound
		}
		object = &response.SignatureInfos.SignatureInfo[0]
		return false, nil
	}); err != nil {
		if err == errSignatureNotFound {
			return nil, nil
		}
		return nil, err
	}
	return object, nil
}

// signatureState maps a SignatureInfo onto the resource model. The API returns the
// sensitive key/secret back on read, so the state round-trips them like the SDKv2
// resource did; empty strings become null to match a Framework-computed absence.
func signatureState(sig *cloudapi.SignatureInfo) signatureResourceModel {
	str := func(v string) types.String {
		if v == "" {
			return types.StringNull()
		}
		return types.StringValue(v)
	}
	return signatureResourceModel{
		Id:              str(sig.SignatureId),
		SignatureName:   str(sig.SignatureName),
		SignatureKey:    str(sig.SignatureKey),
		SignatureSecret: str(sig.SignatureSecret),
		CreatedTime:     str(sig.CreatedTime),
		ModifiedTime:    str(sig.ModifiedTime),
	}
}

var errSignatureNotFound = fmt.Errorf("signature not found")

// isExpectedError reports whether the cloudapi client error carries one of the given
// SDK error codes. The cloudapi client returns *errors.ServerError, whose ErrorCode is
// the POP error code (e.g. NotFoundSignature, RepeatedCommit).
func isExpectedError(err error, codes ...string) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*errors.ServerError); ok {
		for _, code := range codes {
			if e.ErrorCode() == code || strings.Contains(e.Message(), code) {
				return true
			}
		}
		return false
	}
	// Wrapped errors from connectivity helpers keep the original message; fall back to
	// a substring match so the retry classification survives a wrap.
	msg := err.Error()
	for _, code := range codes {
		if strings.Contains(msg, code) {
			return true
		}
	}
	return false
}

// retry runs op until it returns done=true or no error, with the SDKv2 resource's
// exponential-backoff cadence. A true retry result with a nil error stops the loop;
// a non-retryable error surfaces immediately.
func retry(timeout time.Duration, op func() (retry bool, err error)) error {
	deadline := time.Now().Add(timeout)
	wait := 2 * time.Second
	for {
		retry, err := op()
		if err == nil {
			return nil
		}
		if !retry {
			return err
		}
		if time.Now().After(deadline) {
			return err
		}
		time.Sleep(wait)
		if wait < 10*time.Second {
			wait *= 2
		}
	}
}
