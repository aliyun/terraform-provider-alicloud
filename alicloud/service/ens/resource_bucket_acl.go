package ens

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/fwadapt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const bucketAclTypeName = "alicloud_ens_bucket_acl"

// ENS bucket ACL write/read is eventually consistent and non-deterministic right
// after a PutBucketAcl: within the same second GetBucketAcl can oscillate
// between the old and new ACL several times before settling. A single matching
// read is not enough to guarantee the next Read will see the correct value, so
// both the post-write wait and the plain Read require two consecutive polls that
// agree before returning. NotFound after a write is retried (the resource may be
// briefly invisible); during a refresh NotFound is surfaced so the caller can drop
// it from state.
const (
	bucketAclPollInterval   = 3 * time.Second
	bucketAclPollStep       = 5 * time.Second
	bucketAclWaitTimeout    = 5 * time.Minute
	bucketAclReadTimeout    = 2 * time.Minute
	bucketAclDeleteTimeout  = 5 * time.Minute
	bucketAclDefaultACL     = "private"
	bucketAclNotFoundString = "NoSuchBucket"
)

var (
	_ resource.Resource                = &bucketAclResource{}
	_ resource.ResourceWithImportState = &bucketAclResource{}
)

func NewBucketAclResource() resource.Resource {
	return &bucketAclResource{}
}

type bucketAclResource struct {
	fwadapt.ResourceBase
}

type bucketAclResourceModel struct {
	Id         types.String `tfsdk:"id"`
	BucketName types.String `tfsdk:"bucket_name"`
	BucketAcl  types.String `tfsdk:"bucket_acl"`
}

func (r *bucketAclResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = bucketAclTypeName
}

func (r *bucketAclResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a ENS Bucket Acl resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource ID in terraform of Bucket Acl. The value is the bucket name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bucket_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the ENS bucket. The bucket name is the identity of the BucketAcl resource and cannot be changed after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bucket_acl": schema.StringAttribute{
				Required:    true,
				Description: "The read and write permissions of the bucket. Valid values: `public-read-write`, `public-read`, `private`.",
				Validators: []validator.String{
					stringvalidator.OneOf("public-read-write", "public-read", "private"),
				},
			},
		},
	}
}

func (r *bucketAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan bucketAclResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := plan.BucketName.ValueString()
	acl := plan.BucketAcl.ValueString()

	if err := putBucketAcl(r.Client(), bucketName, acl, bucketAclWaitTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: calling PutBucketAcl", bucketAclTypeName),
			errs.WrapErrorf(err, "Creating %s: calling PutBucketAcl for bucket %s", bucketAclTypeName, bucketName).Error(),
		)
		return
	}

	// Wait for the write to be readable before Read persists it to state.
	if err := waitBucketAclConsistent(r.Client(), bucketName, acl, bucketAclWaitTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Creating %s: waiting for ACL to converge", bucketAclTypeName),
			errs.WrapErrorf(err, "Creating %s: ACL %s did not converge for bucket %s", bucketAclTypeName, acl, bucketName).Error(),
		)
		return
	}

	state := bucketAclResourceModel{
		Id:         types.StringValue(bucketName),
		BucketName: types.StringValue(bucketName),
		BucketAcl:  types.StringValue(acl),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *bucketAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state bucketAclResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := state.Id.ValueString()
	if bucketName == "" {
		bucketName = state.BucketName.ValueString()
	}

	acl, err := readBucketAclStable(r.Client(), bucketName, bucketAclReadTimeout)
	if err != nil {
		if isBucketAclNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Reading %s: calling GetBucketAcl", bucketAclTypeName),
			errs.WrapErrorf(err, "Reading %s: calling GetBucketAcl for bucket %s", bucketAclTypeName, bucketName).Error(),
		)
		return
	}

	state.Id = types.StringValue(bucketName)
	state.BucketName = types.StringValue(bucketName)
	state.BucketAcl = types.StringValue(acl)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *bucketAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan bucketAclResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := plan.BucketName.ValueString()
	acl := plan.BucketAcl.ValueString()

	if err := putBucketAcl(r.Client(), bucketName, acl, bucketAclWaitTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: calling PutBucketAcl", bucketAclTypeName),
			errs.WrapErrorf(err, "Updating %s: calling PutBucketAcl for bucket %s", bucketAclTypeName, bucketName).Error(),
		)
		return
	}

	if err := waitBucketAclConsistent(r.Client(), bucketName, acl, bucketAclWaitTimeout); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Updating %s: waiting for ACL to converge", bucketAclTypeName),
			errs.WrapErrorf(err, "Updating %s: ACL %s did not converge for bucket %s", bucketAclTypeName, acl, bucketName).Error(),
		)
		return
	}

	state := bucketAclResourceModel{
		Id:         types.StringValue(bucketName),
		BucketName: types.StringValue(bucketName),
		BucketAcl:  types.StringValue(acl),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *bucketAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state bucketAclResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	bucketName := state.BucketName.ValueString()
	if bucketName == "" {
		bucketName = state.Id.ValueString()
	}

	// BucketAcl is a singleton attribute of an ENS bucket and cannot be removed
	// (there is no DeleteBucketAcl API). On destroy we reset the ACL to the
	// service default "private" so the cloud object returns to its pre-managed
	// state, then drop the resource from Terraform state. A missing bucket is
	// treated as already gone.
	err := putBucketAcl(r.Client(), bucketName, bucketAclDefaultACL, bucketAclDeleteTimeout)
	if err != nil && !isBucketAclNotFound(err) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Deleting %s: calling PutBucketAcl to reset ACL", bucketAclTypeName),
			errs.WrapErrorf(err, "Deleting %s: calling PutBucketAcl to reset ACL for bucket %s", bucketAclTypeName, bucketName).Error(),
		)
		return
	}
}

func (r *bucketAclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getBucketAclOnce performs a single GetBucketAcl call and returns the ACL
// value extracted from the response. The response demo places the grant at the
// top-level BucketAcl key; the nested content.accessControlList.grant path is
// kept as a defensive fallback for SDK-wrapped forms.
func getBucketAclOnce(client *connectivity.AliyunClient, bucketName string) (string, error) {
	query := map[string]interface{}{"BucketName": bucketName}
	response, err := client.RpcPost("Ens", "2017-11-10", "GetBucketAcl", query, map[string]interface{}{}, true)
	if err != nil {
		return "", err
	}
	if acl := stringFromResponse(response); acl != "" {
		return acl, nil
	}
	return "", fmt.Errorf("the GetBucketAcl response carried no BucketAcl: %v", response)
}

func stringFromResponse(response map[string]interface{}) string {
	if v, ok := response["BucketAcl"].(string); ok && v != "" {
		return v
	}
	// Defensive fallback for a wrapped form some SDKs return.
	if content, ok := response["content"].(map[string]interface{}); ok {
		if acl, ok := content["accessControlList"].(map[string]interface{}); ok {
			if grant, ok := acl["grant"].(string); ok && grant != "" {
				return grant
			}
		}
	}
	return ""
}

// putBucketAcl calls PutBucketAcl with a bounded retry on transient errors
// (the ENS backend throttles or briefly rejects immediately after a write).
func putBucketAcl(client *connectivity.AliyunClient, bucketName, acl string, timeout time.Duration) error {
	request := map[string]interface{}{
		"BucketName": bucketName,
		"BucketAcl":  acl,
	}
	deadline := time.Now().Add(timeout)
	interval := bucketAclPollInterval
	step := bucketAclPollStep
	var lastErr error
	for {
		_, err := client.RpcPost("Ens", "2017-11-10", "PutBucketAcl", map[string]interface{}{}, request, true)
		if err == nil {
			return nil
		}
		lastErr = err
		if isBucketAclNotFound(err) {
			// The bucket no longer exists; there is nothing to put.
			return err
		}
		if !isBucketAclRetryable(err) {
			return err
		}
		if time.Now().After(deadline) {
			break
		}
		time.Sleep(interval)
		interval += step
	}
	return errs.WrapErrorf(lastErr, "PutBucketAcl for bucket %s timed out after %s", bucketName, timeout)
}

// waitBucketAclConsistent polls GetBucketAcl until the returned ACL matches the
// value just submitted via PutBucketAcl AND two consecutive reads return the
// same value. NotFound is retried (the resource may be briefly invisible after
// the write).
func waitBucketAclConsistent(client *connectivity.AliyunClient, bucketName, expected string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	interval := bucketAclPollInterval
	step := bucketAclPollStep
	var prevValue string
	consecutive := 0
	for {
		current, err := getBucketAclOnce(client, bucketName)
		if err != nil {
			if isBucketAclNotFound(err) {
				consecutive = 0
				prevValue = ""
				if time.Now().After(deadline) {
					return errs.WrapErrorf(err, "ENS bucket acl %s not converged (still not found) after %s", bucketName, timeout)
				}
				time.Sleep(interval)
				interval += step
				continue
			}
			if isBucketAclRetryable(err) {
				consecutive = 0
				prevValue = ""
				if time.Now().After(deadline) {
					return errs.WrapErrorf(err, "ENS bucket acl %s not converged (transient read error) after %s", bucketName, timeout)
				}
				time.Sleep(interval)
				interval += step
				continue
			}
			return err
		}
		if current != expected {
			consecutive = 0
			prevValue = ""
			if time.Now().After(deadline) {
				return fmt.Errorf("ENS bucket acl %s not converged: expected %s, last saw %s after %s", bucketName, expected, current, timeout)
			}
			time.Sleep(interval)
			interval += step
			continue
		}
		if current == prevValue {
			consecutive++
		} else {
			consecutive = 1
		}
		prevValue = current
		if consecutive >= 2 {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("ENS bucket acl %s matched %s once, timed out waiting for a second consecutive read after %s", bucketName, current, timeout)
		}
		time.Sleep(interval)
		interval += step
	}
}

// readBucketAclStable polls GetBucketAcl until two consecutive reads return the
// same ACL value, absorbing the ENS backend's non-deterministic post-write
// oscillation. It is used by Read (including import) so the value adopted into
// state is stable rather than a transient stale read. NotFound is not retried
// here: in refresh semantics a missing resource should be surfaced so the caller
// can drop it from state.
func readBucketAclStable(client *connectivity.AliyunClient, bucketName string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	interval := bucketAclPollInterval
	step := bucketAclPollStep
	var prevValue string
	consecutive := 0
	for {
		current, err := getBucketAclOnce(client, bucketName)
		if err != nil {
			if isBucketAclNotFound(err) {
				return "", err
			}
			if isBucketAclRetryable(err) {
				if time.Now().After(deadline) {
					return "", errs.WrapErrorf(err, "ENS bucket acl %s read not stable (transient error) after %s", bucketName, timeout)
				}
				time.Sleep(interval)
				interval += step
				continue
			}
			return "", err
		}
		if current == prevValue {
			consecutive++
		} else {
			consecutive = 1
		}
		prevValue = current
		if consecutive >= 2 {
			return current, nil
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("ENS bucket acl %s read not stable (last saw %s) after %s", bucketName, current, timeout)
		}
		time.Sleep(interval)
		interval += step
	}
}

// isBucketAclNotFound reports whether err indicates the bucket no longer
// exists. The code arrives on a *tea.SDKError from RpcPost; the substring
// fallback covers wrapped forms and the errs.NotFoundError sentinel.
func isBucketAclNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errs.NotFoundError(err) {
		return true
	}
	return strings.Contains(err.Error(), bucketAclNotFoundString) ||
		strings.Contains(err.Error(), "InvalidBucketName")
}

// isBucketAclRetryable reports whether err is a transient backend error worth a
// short retry (throttling, internal processing, briefly-busy). It mirrors the
// SDK line's NeedRetry on the common request-throttling and internal-error
// codes without importing package alicloud's NeedRetry.
func isBucketAclRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, frag := range []string{
		"Throttling",
		"RequestTimeout",
		"ServiceUnavailable",
		"InternalError",
		"Busy",
		"try again",
		"try it later",
	} {
		if strings.Contains(msg, frag) {
			return true
		}
	}
	return false
}
